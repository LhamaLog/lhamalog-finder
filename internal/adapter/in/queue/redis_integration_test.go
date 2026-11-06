package queue

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/application"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
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
