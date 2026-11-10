package queue

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/application"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/config"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"github.com/redis/go-redis/v9"
)

func TestDescarteMapaVazioNaoFicaNoMapa(t *testing.T) {
	f := &Fila{
		fluxo: &fluxoSequencia{itens: []domain.Item{
			itemDaMensagem("1-0", map[string]any{}),
			itemDaMensagem("2-0", map[string]any{"c": "ok"}),
		}},
		porID: map[domain.ID]domain.Item{},
	}
	svc := &application.Servico{Fila: f, Reservas: domain.NovoQuadro(nil)}
	pedido := in.Pedido{Instancia: "a", WorkerNoAr: true, InstanciaNoAr: true}
	saida, err := svc.Executar(context.Background(), pedido)
	if err != nil || !saida.ForaDeContrato || saida.Conclusoes != 0 {
		t.Fatalf("vazio: %+v err=%v", saida, err)
	}
	if _, ok := f.porID["1-0"]; ok {
		t.Fatal("item sem valores permaneceu no mapa")
	}
	seguinte, err := svc.Executar(context.Background(), pedido)
	if err != nil || seguinte.Conclusoes != 1 {
		t.Fatalf("seguinte: %+v err=%v", seguinte, err)
	}
}

func TestConclusaoParcialEFalhaNoRedis(t *testing.T) {
	addr, stream, group, codigos := nomesRedis(t)
	f := abrirFila(t, addr, stream, group, codigos, "parcial")
	publicar(t, addr, stream, map[string]any{"c": "apontamento-a"})
	item, ok, err := f.TomarLivre(context.Background())
	if err != nil || !ok {
		t.Fatalf("tomar: ok=%v err=%v", ok, err)
	}
	rdb := cliente(t, addr)
	if err := rdb.SAdd(context.Background(), codigos, "apontamento-a").Err(); err != nil {
		t.Fatal(err)
	}
	nova, err := f.ConcluirPrimeira(context.Background(), item.ID, item.Codigo)
	if err != nil {
		t.Fatal(err)
	}
	if nova || pendencias(t, addr, stream, group) != 0 || membros(t, addr, codigos) != 1 {
		t.Fatalf("nova=%v pendencias=%d membros=%d", nova, pendencias(t, addr, stream, group), membros(t, addr, codigos))
	}

	f2 := abrirFila(t, addr, stream, group, codigos, "falha")
	publicar(t, addr, stream, map[string]any{"c": "apontamento-b"})
	item, ok, err = f2.TomarLivre(context.Background())
	if err != nil || !ok {
		t.Fatalf("tomar falha: ok=%v err=%v", ok, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := f2.ConcluirPrimeira(ctx, item.ID, item.Codigo); err == nil {
		t.Fatal("esperava falha")
	}
	if pendencias(t, addr, stream, group) == 0 {
		t.Fatal("a falha tirou a mensagem da pendência")
	}
	nova, err = f2.ConcluirPrimeira(context.Background(), item.ID, item.Codigo)
	if err != nil || !nova {
		t.Fatalf("retomada nova=%v err=%v", nova, err)
	}
	if pendencias(t, addr, stream, group) != 0 {
		t.Fatal("a tentativa seguinte não tirou a pendência")
	}
}

func TestInstanciaParadaDevolveNoLugarOriginal(t *testing.T) {
	addr, stream, group, codigos := nomesRedis(t)
	a := abrirFila(t, addr, stream, group, codigos, "inst-a")
	publicar(t, addr, stream, map[string]any{"c": "frente"})
	publicar(t, addr, stream, map[string]any{"c": "atras"})
	frente, ok, err := a.TomarLivre(context.Background())
	if err != nil || !ok || frente.Codigo != "frente" {
		t.Fatalf("frente: %+v ok=%v err=%v", frente, ok, err)
	}
	if err := a.Fechar(); err != nil {
		t.Fatal(err)
	}
	b := abrirFila(t, addr, stream, group, codigos, "inst-b")
	deNovo, ok, err := b.TomarLivre(context.Background())
	if err != nil || !ok {
		t.Fatalf("retomada: ok=%v err=%v", ok, err)
	}
	if deNovo.ID != frente.ID || deNovo.Codigo != "frente" {
		t.Fatalf("lugar: %+v, esperado %s", deNovo, frente.ID)
	}
}

func TestConsumidorVivoNaoPerdeOItem(t *testing.T) {
	addr, stream, group, codigos := nomesRedis(t)
	a := abrirFila(t, addr, stream, group, codigos, "vivo-a")
	b := abrirFila(t, addr, stream, group, codigos, "vivo-b")
	publicar(t, addr, stream, map[string]any{"c": "um"})
	publicar(t, addr, stream, map[string]any{"c": "dois"})
	primeiro, ok, err := a.TomarLivre(context.Background())
	if err != nil || !ok {
		t.Fatalf("a: ok=%v err=%v", ok, err)
	}
	segundo, ok, err := b.TomarLivre(context.Background())
	if err != nil || !ok {
		t.Fatalf("b: ok=%v err=%v", ok, err)
	}
	if primeiro.ID == segundo.ID {
		t.Fatalf("o mesmo item foi para as duas instâncias: %s", primeiro.ID)
	}
}

func TestOrdemDeEntradaComUmaInstancia(t *testing.T) {
	addr, stream, group, codigos := nomesRedis(t)
	f := abrirFila(t, addr, stream, group, codigos, "ordem")
	for _, codigo := range []string{"c", "a", "b"} {
		publicar(t, addr, stream, map[string]any{"c": codigo})
	}
	var ordem []domain.Codigo
	for i := 0; i < 3; i++ {
		item, ok, err := f.TomarLivre(context.Background())
		if err != nil || !ok {
			t.Fatalf("tomar %d: ok=%v err=%v", i, ok, err)
		}
		ordem = append(ordem, item.Codigo)
		if _, err := f.ConcluirPrimeira(context.Background(), item.ID, item.Codigo); err != nil {
			t.Fatal(err)
		}
	}
	if ordem[0] != "c" || ordem[1] != "a" || ordem[2] != "b" {
		t.Fatalf("ordem = %v", ordem)
	}
}

func TestDuasInstanciasReservamItensDiferentes(t *testing.T) {
	addr, stream, group, codigos := nomesRedis(t)
	a := abrirFila(t, addr, stream, group, codigos, "par-a")
	b := abrirFila(t, addr, stream, group, codigos, "par-b")
	publicar(t, addr, stream, map[string]any{"c": "1"})
	publicar(t, addr, stream, map[string]any{"c": "2"})
	ia, ok, err := a.TomarLivre(context.Background())
	if err != nil || !ok {
		t.Fatal(err)
	}
	ib, ok, err := b.TomarLivre(context.Background())
	if err != nil || !ok {
		t.Fatal(err)
	}
	if ia.ID == ib.ID {
		t.Fatal("reservaram o mesmo item")
	}
}

func TestMesmoCodigoUmaConclusao(t *testing.T) {
	addr, stream, group, codigos := nomesRedis(t)
	f := abrirFila(t, addr, stream, group, codigos, "unico")
	publicar(t, addr, stream, map[string]any{"c": "mesmo"})
	publicar(t, addr, stream, map[string]any{"c": "mesmo"})
	primeira, err := concluirProximo(t, f)
	if err != nil || !primeira {
		t.Fatalf("primeira=%v err=%v", primeira, err)
	}
	segunda, err := concluirProximo(t, f)
	if err != nil || segunda {
		t.Fatalf("segunda=%v err=%v", segunda, err)
	}
	if membros(t, addr, codigos) != 1 || pendencias(t, addr, stream, group) != 0 {
		t.Fatalf("membros=%d pendencias=%d", membros(t, addr, codigos), pendencias(t, addr, stream, group))
	}
	rdb := cliente(t, addr)
	lista, err := rdb.SMembers(context.Background(), codigos).Result()
	if err != nil {
		t.Fatal(err)
	}
	if len(lista) != 1 || lista[0] != "mesmo" {
		t.Fatalf("conjunto = %v", lista)
	}
}

func TestDescarteInvalidoEItemAceito(t *testing.T) {
	addr, stream, group, codigos := nomesRedis(t)
	f := abrirFila(t, addr, stream, group, codigos, "descarte")
	publicar(t, addr, stream, map[string]any{"c": ""})
	publicar(t, addr, stream, map[string]any{"a": "1", "b": "2"})
	publicar(t, addr, stream, map[string]any{"c": "valido"})
	svc := &application.Servico{Fila: f, Reservas: domain.NovoQuadro(nil)}
	pedido := in.Pedido{Instancia: "a", WorkerNoAr: true, InstanciaNoAr: true}
	for i := 0; i < 2; i++ {
		saida, err := svc.Executar(context.Background(), pedido)
		if err != nil || !saida.ForaDeContrato || saida.Conclusoes != 0 {
			t.Fatalf("descarte %d: %+v err=%v", i, saida, err)
		}
	}
	if len(f.porID) != 0 {
		t.Fatalf("mapa = %v", f.porID)
	}
	saida, err := svc.Executar(context.Background(), pedido)
	if err != nil || saida.Conclusoes != 1 || len(saida.CodigosDisponiveis) != 1 || saida.CodigosDisponiveis[0] != "valido" {
		t.Fatalf("aceito: %+v err=%v", saida, err)
	}
	if membros(t, addr, codigos) != 1 || pendencias(t, addr, stream, group) != 0 {
		t.Fatalf("membros=%d pendencias=%d", membros(t, addr, codigos), pendencias(t, addr, stream, group))
	}
}

func concluirProximo(t *testing.T, f *Fila) (bool, error) {
	t.Helper()
	item, ok, err := f.TomarLivre(context.Background())
	if err != nil || !ok {
		t.Fatalf("tomar: ok=%v err=%v", ok, err)
	}
	return f.ConcluirPrimeira(context.Background(), item.ID, item.Codigo)
}

func nomesRedis(t *testing.T) (addr, stream, group, codigos string) {
	t.Helper()
	if os.Getenv("FINDER_REDIS_INTEGRATION") != "1" {
		t.Skip("FINDER_REDIS_INTEGRATION não definido")
	}
	addr = os.Getenv("FINDER_REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	stream = os.Getenv("FINDER_REDIS_STREAM")
	group = os.Getenv("FINDER_REDIS_GROUP")
	codigos = os.Getenv("FINDER_REDIS_CODIGOS")
	if stream == "" || group == "" || codigos == "" {
		t.Fatal("FINDER_REDIS_STREAM, FINDER_REDIS_GROUP ou FINDER_REDIS_CODIGOS ausente")
	}
	sufixo := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	stream += ":" + sufixo
	group += ":" + sufixo
	codigos += ":" + sufixo
	rdb := cliente(t, addr)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Fatalf("redis: %v", err)
	}
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = rdb.Del(c, stream, codigos).Err()
	})
	return addr, stream, group, codigos
}

func abrirFila(t *testing.T, addr, stream, group, codigos, consumidor string) *Fila {
	t.Helper()
	f, err := Nova(Opcoes{
		Addr:       addr,
		Stream:     stream,
		Group:      group,
		Codigos:    codigos,
		Consumidor: consumidor,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Fechar() })
	return f
}

func publicar(t *testing.T, addr, stream string, valores map[string]any) {
	t.Helper()
	rdb := cliente(t, addr)
	if _, err := rdb.XAdd(context.Background(), &redis.XAddArgs{Stream: stream, Values: valores}).Result(); err != nil {
		t.Fatal(err)
	}
}

func cliente(t *testing.T, addr string) *redis.Client {
	t.Helper()
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	t.Cleanup(func() { _ = rdb.Close() })
	return rdb
}

func pendencias(t *testing.T, addr, stream, group string) int64 {
	t.Helper()
	p, err := cliente(t, addr).XPending(context.Background(), stream, group).Result()
	if err != nil {
		t.Fatal(err)
	}
	return p.Count
}

func membros(t *testing.T, addr, codigos string) int64 {
	t.Helper()
	n, err := cliente(t, addr).SCard(context.Background(), codigos).Result()
	if err != nil {
		t.Fatal(err)
	}
	return n
}

type fluxoSequencia struct {
	itens []domain.Item
	i     int
}

func (f *fluxoSequencia) Ler(context.Context) (domain.Item, bool, error) {
	if f.i >= len(f.itens) {
		return domain.Item{}, false, nil
	}
	item := f.itens[f.i]
	f.i++
	return item, true, nil
}

func (f *fluxoSequencia) Ack(context.Context, domain.ID) error { return nil }

func (f *fluxoSequencia) ConcluirCodigo(context.Context, domain.ID, domain.Codigo) (bool, error) {
	return true, nil
}

func (f *fluxoSequencia) Marcar(context.Context, domain.Codigo) (bool, error) {
	return true, nil
}

func (f *fluxoSequencia) Existe(context.Context, domain.Codigo) (bool, error) {
	return false, nil
}

func (f *fluxoSequencia) Fechar() error { return nil }

func TestGrupoLeEntradaAnteriorSemDolar(t *testing.T) {
	op := provaLote(t)
	id := publicarID(t, op, map[string]any{"event": "antes"})
	lote := abrirLote(t, op)
	if ultimo := ultimoIDEntregue(t, op); ultimo != "0-0" {
		t.Fatalf("posição inicial = %s", ultimo)
	}
	ctx, cancel := context.WithTimeout(context.Background(), op.Block+time.Second)
	defer cancel()
	entrada, ok, err := lote.Ler(ctx)
	if err != nil || !ok || entrada.ID != id {
		t.Fatalf("entrada=%+v ok=%v err=%v", entrada, ok, err)
	}
}

func TestDoisGruposObservamAMesmaEntrada(t *testing.T) {
	op := provaLote(t)
	varredura := op
	varredura.Group = op.Group + "-varredura"
	varredura.Consumidor = op.Consumidor + "-varredura"
	abrirLote(t, varredura)
	finder := abrirLote(t, op)
	nomes := map[string]bool{}
	for _, grupo := range infoGrupos(t, op) {
		nomes[grupo.Name] = true
	}
	if !nomes[op.Group] || !nomes[varredura.Group] {
		t.Fatalf("grupos = %v", nomes)
	}
	id := publicarID(t, op, map[string]any{"event": "compartilhada"})
	ctx, cancel := context.WithTimeout(context.Background(), op.Block+time.Second)
	defer cancel()
	primeira, ok, err := finder.Ler(ctx)
	if err != nil || !ok || primeira.ID != id {
		t.Fatalf("finder=%+v ok=%v err=%v", primeira, ok, err)
	}
	segunda, ok, err := abrirLote(t, varredura).Ler(ctx)
	if err != nil || !ok || segunda.ID != id {
		t.Fatalf("varredura=%+v ok=%v err=%v", segunda, ok, err)
	}
}

func TestInvalidaSaiDaPendenciaSemChaveConcluida(t *testing.T) {
	op := provaLote(t)
	lote := abrirLote(t, op)
	rdb := clienteSenha(t, op.Addr, op.Senha)
	if err := rdb.Set(context.Background(), nomeDaChave(op.Idempotencia, "batch-1"), domain.EstadoEmProcessamento, 0).Err(); err != nil {
		t.Fatal(err)
	}
	publicarID(t, op, map[string]any{
		"idempotencyId": "batch-1",
		"event":         "trace.ingested",
		"origin":        "collector",
		"schemaVersion": "1",
		"publishedAt":   "2026-11-08T12:00:00Z",
		"metadata":      `{"senha":"segredo-redis"}`,
		"body":          `CONTEUDO-DE-LOG-SECRETO`,
	})
	efeito := &efeitoIntegracao{}
	svc := &application.ConsumoReferencia{Fila: lote, Chaves: lote, Efeito: efeito}
	ctx, cancel := context.WithTimeout(context.Background(), op.Block+time.Second)
	defer cancel()
	saida, err := svc.Executar(ctx)
	if err != nil || !saida.Invalida || saida.Confirmacoes != 1 || len(efeito.refs) != 0 {
		t.Fatalf("saída=%+v refs=%d err=%v", saida, len(efeito.refs), err)
	}
	if pendencias(t, op.Addr, op.Stream, op.Group) != 0 {
		t.Fatal("entrada inválida permaneceu pendente")
	}
	valor, err := rdb.Get(context.Background(), nomeDaChave(op.Idempotencia, "batch-1")).Result()
	if err != redis.Nil || valor == domain.EstadoConcluido {
		t.Fatalf("chave=%q err=%v", valor, err)
	}
}

func TestChaveEmProcessamentoSemPrazo(t *testing.T) {
	op := provaLote(t)
	lote := abrirLote(t, op)
	publicarID(t, op, camposIntegracao("batch-1"))
	efeito := &efeitoIntegracao{}
	svc := &application.ConsumoReferencia{Fila: lote, Chaves: lote, Efeito: efeito}
	ctx, cancel := context.WithTimeout(context.Background(), op.Block+time.Second)
	defer cancel()
	if _, err := svc.Executar(ctx); err != nil {
		t.Fatal(err)
	}
	rdb := clienteSenha(t, op.Addr, op.Senha)
	chave := nomeDaChave(op.Idempotencia, "batch-1")
	valor, err := rdb.Get(context.Background(), chave).Result()
	if err != nil || valor != domain.EstadoEmProcessamento {
		t.Fatalf("valor=%q err=%v", valor, err)
	}
	ttl, err := rdb.PTTL(context.Background(), chave).Result()
	if err != nil || ttl != -1 {
		t.Fatalf("prazo da reserva = %s err=%v", ttl, err)
	}
	if pendencias(t, op.Addr, op.Stream, op.Group) != 1 || len(efeito.refs) != 1 {
		t.Fatalf("pendencias=%d efeitos=%d", pendencias(t, op.Addr, op.Stream, op.Group), len(efeito.refs))
	}
}

func TestChaveConcluidaDuraVinteCincoHorasSemRenovar(t *testing.T) {
	op := provaLote(t)
	lote := abrirLote(t, op)
	publicarID(t, op, camposIntegracao("batch-1"))
	efeito := &efeitoIntegracao{concluido: true}
	svc := &application.ConsumoReferencia{Fila: lote, Chaves: lote, Efeito: efeito}
	ctx, cancel := context.WithTimeout(context.Background(), op.Block+time.Second)
	defer cancel()
	saida, err := svc.Executar(ctx)
	if err != nil || saida.Efeitos != 1 || saida.Confirmacoes != 1 {
		t.Fatalf("saída=%+v err=%v", saida, err)
	}
	rdb := clienteSenha(t, op.Addr, op.Senha)
	chave := nomeDaChave(op.Idempotencia, "batch-1")
	valor, err := rdb.Get(context.Background(), chave).Result()
	if err != nil || valor != domain.EstadoConcluido {
		t.Fatalf("valor=%q err=%v", valor, err)
	}
	ttl, err := rdb.PTTL(context.Background(), chave).Result()
	if err != nil || ttl <= 0 || ttl > domain.PrazoConcluido || ttl < domain.PrazoConcluido-5*time.Second {
		t.Fatalf("prazo concluído = %s", ttl)
	}
	if pendencias(t, op.Addr, op.Stream, op.Group) != 0 {
		t.Fatal("entrada concluída permaneceu pendente")
	}
	time.Sleep(1100 * time.Millisecond)
	if err := lote.Concluir(context.Background(), "batch-1"); err != nil {
		t.Fatal(err)
	}
	deNovo, err := rdb.PTTL(context.Background(), chave).Result()
	if err != nil || deNovo >= ttl {
		t.Fatalf("prazo renovado: antes=%s depois=%s", ttl, deNovo)
	}
}

func TestSegundaEntradaNaoRepeteEfeito(t *testing.T) {
	op := provaLote(t)
	lote := abrirLote(t, op)
	publicarID(t, op, camposIntegracao("batch-1"))
	publicarID(t, op, camposIntegracao("batch-1"))
	efeito := &efeitoIntegracao{}
	svc := &application.ConsumoReferencia{Fila: lote, Chaves: lote, Efeito: efeito}
	ctx, cancel := context.WithTimeout(context.Background(), 4*op.Block+2*time.Second)
	defer cancel()
	for i := 0; i < 4; i++ {
		if _, err := svc.Executar(ctx); err != nil && ctx.Err() == nil {
			t.Fatal(err)
		}
	}
	if len(efeito.refs) != 1 {
		t.Fatalf("efeitos=%d", len(efeito.refs))
	}
	if pendencias(t, op.Addr, op.Stream, op.Group) != 2 {
		t.Fatalf("pendencias=%d", pendencias(t, op.Addr, op.Stream, op.Group))
	}
}

func TestRedisFechadoNaoConfirma(t *testing.T) {
	op := provaLote(t)
	lote := abrirLote(t, op)
	publicarID(t, op, map[string]any{"event": "no-meio"})
	ctx, cancel := context.WithTimeout(context.Background(), op.Block+time.Second)
	defer cancel()
	entrada, ok, err := lote.Ler(ctx)
	if err != nil || !ok {
		t.Fatalf("ler ok=%v err=%v", ok, err)
	}
	if err := lote.Fechar(); err != nil {
		t.Fatal(err)
	}
	if err := lote.Confirmar(context.Background(), entrada.ID); err == nil {
		t.Fatal("confirmou com redis fechado")
	}
	if pendencias(t, op.Addr, op.Stream, op.Group) == 0 {
		t.Fatal("entrada inacabada foi confirmada")
	}
}

func TestReivindicacaoRespeitaQuantidade(t *testing.T) {
	op := provaLote(t)
	a := abrirLote(t, op)
	total := op.Count + 1
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(total+2)*op.Block+2*time.Second)
	defer cancel()
	for i := int64(0); i < total; i++ {
		publicarID(t, op, map[string]any{"n": strconv.FormatInt(i, 10)})
	}
	vistos := map[string]struct{}{}
	for len(vistos) < int(total) {
		entrada, ok, err := a.Ler(ctx)
		if err != nil || !ok {
			t.Fatalf("ler %d ok=%v err=%v", len(vistos), ok, err)
		}
		if _, repetida := vistos[entrada.ID]; repetida {
			t.Fatalf("entrada %s foi lida de novo na mesma passagem", entrada.ID)
		}
		vistos[entrada.ID] = struct{}{}
		if err := a.Adiar(context.Background(), entrada.ID); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(op.Ociosidade + 400*time.Millisecond)
	bOp := op
	bOp.Consumidor = op.Consumidor + "-b"
	b := abrirLote(t, bOp)
	entrada, ok, err := b.Ler(context.Background())
	if err != nil || !ok || entrada.ID == "" {
		t.Fatalf("reivindicar ok=%v entrada=%+v err=%v", ok, entrada, err)
	}
	if n := pendenciasConsumidor(t, op, bOp.Consumidor); n != op.Count {
		t.Fatalf("reivindicadas=%d quantidade=%d", n, op.Count)
	}
	if n := pendenciasConsumidor(t, op, op.Consumidor); n < 1 {
		t.Fatalf("pendências restantes=%d", n)
	}
}

type opcoesProva struct {
	Addr         string
	Senha        string
	Stream       string
	Group        string
	Consumidor   string
	Idempotencia string
	Count        int64
	Block        time.Duration
	Ociosidade   time.Duration
}

func (o opcoesProva) lote() OpcoesLote {
	return OpcoesLote{
		Addr:         o.Addr,
		Senha:        o.Senha,
		Stream:       o.Stream,
		Group:        o.Group,
		Consumidor:   o.Consumidor,
		Idempotencia: o.Idempotencia,
		Count:        o.Count,
		Block:        o.Block,
		Ociosidade:   o.Ociosidade,
	}
}

func provaLote(t *testing.T) opcoesProva {
	t.Helper()
	if os.Getenv("FINDER_REDIS_INTEGRATION") != "1" {
		t.Skip("FINDER_REDIS_INTEGRATION não definido")
	}
	addr := os.Getenv("FINDER_REDIS_ADDR")
	if addr == "" {
		t.Fatal("FINDER_REDIS_ADDR não responde")
	}
	cfg := config.Carregar()
	if cfg.RedisStream == "" || cfg.RedisGroup == "" || cfg.RedisConsumidor == "" || cfg.RedisIdempotencia == "" || cfg.RedisCount <= 0 || cfg.RedisBlock <= 0 || cfg.RedisPendenciaOciosa <= 0 {
		t.Fatal("stream, grupo, consumidor, prefixo, quantidade, bloqueio ou ociosidade ausente no ambiente do teste")
	}
	rdb := clienteSenha(t, addr, cfg.RedisSenha)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Fatalf("redis: %v", err)
	}
	sufixo := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	op := opcoesProva{
		Addr:         addr,
		Senha:        cfg.RedisSenha,
		Stream:       cfg.RedisStream + ":" + sufixo,
		Group:        cfg.RedisGroup + ":" + sufixo,
		Consumidor:   cfg.RedisConsumidor + ":" + sufixo,
		Idempotencia: cfg.RedisIdempotencia + sufixo + ":",
		Count:        int64(cfg.RedisCount),
		Block:        cfg.RedisBlock,
		Ociosidade:   cfg.RedisPendenciaOciosa,
	}
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		var chaves []string
		iter := rdb.Scan(c, 0, op.Idempotencia+"*", 50).Iterator()
		for iter.Next(c) {
			chaves = append(chaves, iter.Val())
		}
		if len(chaves) > 0 {
			_ = rdb.Del(c, chaves...).Err()
		}
		_ = rdb.Del(c, op.Stream).Err()
	})
	return op
}

func abrirLote(t *testing.T, op opcoesProva) *Lote {
	t.Helper()
	lote, err := NovaLote(op.lote())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = lote.Fechar() })
	return lote
}

func publicarID(t *testing.T, op opcoesProva, valores map[string]any) string {
	t.Helper()
	id, err := clienteSenha(t, op.Addr, op.Senha).XAdd(context.Background(), &redis.XAddArgs{
		Stream: op.Stream,
		Values: valores,
	}).Result()
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func clienteSenha(t *testing.T, addr, senha string) *redis.Client {
	t.Helper()
	rdb := redis.NewClient(&redis.Options{Addr: addr, Password: senha})
	t.Cleanup(func() { _ = rdb.Close() })
	return rdb
}

func ultimoIDEntregue(t *testing.T, op opcoesProva) string {
	t.Helper()
	for _, grupo := range infoGrupos(t, op) {
		if grupo.Name == op.Group {
			return grupo.LastDeliveredID
		}
	}
	t.Fatal("grupo não encontrado")
	return ""
}

func infoGrupos(t *testing.T, op opcoesProva) []redis.XInfoGroup {
	t.Helper()
	grupos, err := clienteSenha(t, op.Addr, op.Senha).XInfoGroups(context.Background(), op.Stream).Result()
	if err != nil {
		t.Fatal(err)
	}
	return grupos
}

func pendenciasConsumidor(t *testing.T, op opcoesProva, consumidor string) int64 {
	t.Helper()
	lista, err := clienteSenha(t, op.Addr, op.Senha).XPendingExt(context.Background(), &redis.XPendingExtArgs{
		Stream:   op.Stream,
		Group:    op.Group,
		Start:    "-",
		End:      "+",
		Count:    op.Count + 10,
		Consumer: consumidor,
	}).Result()
	if err == redis.Nil {
		return 0
	}
	if err != nil {
		t.Fatal(err)
	}
	return int64(len(lista))
}

func camposIntegracao(batch string) map[string]any {
	meta := `{"tenantId":"tenant-1","productId":"prod-1","correlationId":"corr-1"}`
	body := `{"signal":"LOGS","productId":"prod-1","recordCount":1,"openobserve":{"organization":"org-1","stream":"logs","batchId":"` + batch + `","recordId":"rec-1","minTimestampMicros":100,"maxTimestampMicros":200}}`
	return map[string]any{
		"idempotencyId": batch,
		"event":         "log.ingested",
		"origin":        "collector",
		"schemaVersion": "1",
		"publishedAt":   "2026-11-08T12:00:00Z",
		"metadata":      meta,
		"body":          body,
	}
}

type efeitoIntegracao struct {
	refs      []domain.Referencia
	concluido bool
}

func (e *efeitoIntegracao) Executar(_ context.Context, ref domain.Referencia) (bool, error) {
	e.refs = append(e.refs, ref)
	return e.concluido, nil
}
