package application

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

func TestExecutarConcluiEDisponibilizaUmaVez(t *testing.T) {
	fila := novaFila(domain.Item{ID: "1", Codigo: "apontamento-a", Estado: domain.Livre})
	codigos := &memCodigos{}
	disp := &memDisp{}
	s := &Servico{Fila: fila, Codigos: codigos, Disp: disp, Reservas: domain.NovoQuadro(nil)}
	saida, err := s.Executar(context.Background(), in.Pedido{
		Instancia: "a", WorkerNoAr: true, InstanciaNoAr: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if saida.Conclusoes != 1 || len(saida.CodigosDisponiveis) != 1 || saida.CodigosDisponiveis[0] != "apontamento-a" {
		t.Fatalf("saída: %+v", saida)
	}
	if !fila.codigos["apontamento-a"] || fila.concluidos["1"] != true {
		t.Fatalf("codigos=%v concluidos=%v", fila.codigos, fila.concluidos)
	}
	if fila.livres["1"] {
		t.Fatal("item concluído voltou a livre")
	}
}

func TestExecutarSegundoCodigoNaoConclui(t *testing.T) {
	fila := novaFila(
		domain.Item{ID: "1", Codigo: "mesmo", Estado: domain.Livre},
		domain.Item{ID: "2", Codigo: "mesmo", Estado: domain.Livre},
	)
	codigos := &memCodigos{}
	disp := &memDisp{}
	s := &Servico{Fila: fila, Codigos: codigos, Disp: disp, Reservas: domain.NovoQuadro(nil)}
	pedido := in.Pedido{Instancia: "a", WorkerNoAr: true, InstanciaNoAr: true}
	if _, err := s.Executar(context.Background(), pedido); err != nil {
		t.Fatal(err)
	}
	segunda, err := s.Executar(context.Background(), pedido)
	if err != nil {
		t.Fatal(err)
	}
	if segunda.Conclusoes != 0 || len(fila.codigos) != 1 {
		t.Fatalf("segunda=%+v codigos=%v", segunda, fila.codigos)
	}
	if fila.reservados["2"] || fila.livres["2"] || fila.concluidos["2"] {
		t.Fatalf("estado do segundo: livres=%v reservados=%v concluidos=%v encerrados=%v",
			fila.livres, fila.reservados, fila.concluidos, fila.encerrados)
	}
	if !fila.encerrados["2"] {
		t.Fatal("segundo item não foi encerrado sem conclusão")
	}
}

func TestExecutarParadaDoWorkerLibera(t *testing.T) {
	fila := novaFila(domain.Item{ID: "1", Codigo: "apontamento-a", Estado: domain.Livre})
	disp := &memDisp{}
	s := &Servico{Fila: fila, Codigos: &memCodigos{}, Disp: disp, Reservas: domain.NovoQuadro(nil)}
	saida, err := s.Executar(context.Background(), in.Pedido{
		Instancia: "a", WorkerNoAr: false, InstanciaNoAr: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if saida.Conclusoes != 0 || len(disp.codigos) != 0 || !fila.livres["1"] {
		t.Fatalf("saída=%+v livres=%v", saida, fila.livres)
	}
}

func TestExecutarFilaVaziaEFalhaDeLeitura(t *testing.T) {
	vazia := novaFila()
	s := &Servico{Fila: vazia, Codigos: &memCodigos{}, Disp: &memDisp{}, Reservas: domain.NovoQuadro(nil)}
	saida, err := s.Executar(context.Background(), in.Pedido{Instancia: "a", WorkerNoAr: true, InstanciaNoAr: true})
	if err != nil || saida.Conclusoes != 0 {
		t.Fatalf("vazia: %+v err=%v", saida, err)
	}
	falha := &filaFake{err: errors.New("falha de leitura")}
	s.Fila = falha
	saida, err = s.Executar(context.Background(), in.Pedido{Instancia: "a", WorkerNoAr: true, InstanciaNoAr: true})
	if err == nil || !saida.FalhaLeitura || saida.Conclusoes != 0 {
		t.Fatalf("falha: %+v err=%v", saida, err)
	}
}

func TestExecutarParcialDepoisCompletaSemSegundaContagem(t *testing.T) {
	fila := novaFila(domain.Item{ID: "1", Codigo: "apontamento-a", Estado: domain.Livre})
	fila.falhaConcluir = true
	s := &Servico{Fila: fila, Codigos: &memCodigos{}, Disp: &memDisp{}, Reservas: domain.NovoQuadro(nil)}
	pedido := in.Pedido{Instancia: "a", WorkerNoAr: true, InstanciaNoAr: true}
	if _, err := s.Executar(context.Background(), pedido); err == nil {
		t.Fatal("esperava falha ao completar")
	}
	if !fila.reservados["1"] || !fila.codigos["apontamento-a"] {
		t.Fatalf("parcial: reservados=%v codigos=%v", fila.reservados, fila.codigos)
	}
	saida, err := s.Executar(context.Background(), pedido)
	if err != nil {
		t.Fatal(err)
	}
	if saida.Conclusoes != 0 || fila.reservados["1"] || !fila.codigos["apontamento-a"] {
		t.Fatalf("saída=%+v reservados=%v codigos=%v", saida, fila.reservados, fila.codigos)
	}
}

func TestExecutarDescartaSemPrender(t *testing.T) {
	var buf bytes.Buffer
	anterior := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(anterior) })

	fila := novaFila(
		domain.Item{ID: "1", Estado: domain.Livre},
		domain.Item{ID: "2", Codigo: "a", ConteudoExtra: true, Estado: domain.Livre},
		domain.Item{ID: "3", Codigo: "ok", Estado: domain.Livre},
	)
	s := &Servico{Fila: fila, Codigos: &memCodigos{}, Disp: &memDisp{}, Reservas: domain.NovoQuadro(nil)}
	pedido := in.Pedido{Instancia: "a", WorkerNoAr: true, InstanciaNoAr: true}
	for _, id := range []domain.ID{"1", "2"} {
		saida, err := s.Executar(context.Background(), pedido)
		if err != nil || !saida.ForaDeContrato || saida.Conclusoes != 0 {
			t.Fatalf("%s: %+v err=%v", id, saida, err)
		}
		if fila.reservados[id] {
			t.Fatalf("%s permaneceu reservado", id)
		}
	}
	if strings.Contains(buf.String(), "item concluído") || !strings.Contains(buf.String(), "item descartado") {
		t.Fatalf("registro: %s", buf.String())
	}
	saida, err := s.Executar(context.Background(), pedido)
	if err != nil || saida.Conclusoes != 1 || saida.CodigosDisponiveis[0] != "ok" {
		t.Fatalf("aceito: %+v err=%v", saida, err)
	}
}

func TestExecutarConteudoExtraNaoConclui(t *testing.T) {
	fila := novaFila(domain.Item{ID: "1", Codigo: "apontamento-a", Estado: domain.Livre, ConteudoExtra: true})
	disp := &memDisp{}
	s := &Servico{Fila: fila, Codigos: &memCodigos{}, Disp: disp, Reservas: domain.NovoQuadro(nil)}
	saida, err := s.Executar(context.Background(), in.Pedido{Instancia: "a", WorkerNoAr: true, InstanciaNoAr: true})
	if err != nil || !saida.ForaDeContrato || saida.Conclusoes != 0 || len(disp.codigos) != 0 {
		t.Fatalf("saída=%+v disp=%v err=%v", saida, disp.codigos, err)
	}
}

type filaFake struct {
	espera        []domain.Item
	deNovo        []domain.Item
	tomados       map[domain.ID]domain.Item
	livres        map[domain.ID]bool
	reservados    map[domain.ID]bool
	concluidos    map[domain.ID]bool
	encerrados    map[domain.ID]bool
	codigos       map[domain.Codigo]bool
	falhaConcluir bool
	err           error
}

func novaFila(itens ...domain.Item) *filaFake {
	f := &filaFake{
		tomados:    map[domain.ID]domain.Item{},
		livres:     map[domain.ID]bool{},
		reservados: map[domain.ID]bool{},
		concluidos: map[domain.ID]bool{},
		encerrados: map[domain.ID]bool{},
		codigos:    map[domain.Codigo]bool{},
	}
	for _, item := range itens {
		f.espera = append(f.espera, item)
		f.livres[item.ID] = true
	}
	return f
}

func (f *filaFake) TomarLivre(context.Context) (domain.Item, bool, error) {
	if f.err != nil {
		return domain.Item{}, false, f.err
	}
	if len(f.deNovo) > 0 {
		item := f.deNovo[0]
		f.deNovo = f.deNovo[1:]
		return item, true, nil
	}
	if len(f.espera) == 0 {
		return domain.Item{}, false, nil
	}
	item := f.espera[0]
	f.espera = f.espera[1:]
	delete(f.livres, item.ID)
	f.reservados[item.ID] = true
	f.tomados[item.ID] = item
	return item, true, nil
}

func (f *filaFake) ConcluirPrimeira(_ context.Context, id domain.ID, codigo domain.Codigo) (bool, error) {
	if !f.reservados[id] {
		return false, errors.New("item não reservado")
	}
	if f.falhaConcluir {
		f.falhaConcluir = false
		f.codigos[codigo] = true
		f.deNovo = append(f.deNovo, f.tomados[id])
		return false, errors.New("falha ao concluir")
	}
	nova := !f.codigos[codigo]
	f.codigos[codigo] = true
	delete(f.reservados, id)
	if nova {
		f.concluidos[id] = true
	} else {
		f.encerrados[id] = true
	}
	return nova, nil
}

func (f *filaFake) Concluir(_ context.Context, id domain.ID) error {
	return f.mover(id, func() { f.concluidos[id] = true })
}

func (f *filaFake) Liberar(_ context.Context, id domain.ID) error {
	return f.mover(id, func() { f.livres[id] = true })
}

func (f *filaFake) EncerrarSemConclusao(_ context.Context, id domain.ID) error {
	return f.mover(id, func() { f.encerrados[id] = true })
}

func (f *filaFake) mover(id domain.ID, marcar func()) error {
	if !f.reservados[id] {
		return errors.New("item não reservado")
	}
	delete(f.reservados, id)
	marcar()
	return nil
}

type memCodigos struct {
	vistos map[domain.Codigo]bool
}

func (m *memCodigos) RegistrarPrimeira(_ context.Context, codigo domain.Codigo) (bool, error) {
	if m.vistos == nil {
		m.vistos = map[domain.Codigo]bool{}
	}
	if m.vistos[codigo] {
		return false, nil
	}
	m.vistos[codigo] = true
	return true, nil
}

func (m *memCodigos) JaConcluido(_ context.Context, codigo domain.Codigo) (bool, error) {
	return m.vistos[codigo], nil
}

type memDisp struct {
	codigos []domain.Codigo
}

func (m *memDisp) Disponibilizar(_ context.Context, codigo domain.Codigo) error {
	m.codigos = append(m.codigos, codigo)
	return nil
}

func TestInterpretacaoEntregaSemConfirmar(t *testing.T) {
	casos := []struct {
		nome     string
		contagem int
		recordID any
		trace    string
		comTrace bool
	}{
		{nome: "um registro com trace", contagem: 1, recordID: "rec-1", trace: "trace-1", comTrace: true},
		{nome: "varios sem trace", contagem: 4, recordID: nil, comTrace: false},
	}
	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			fila := &memEntradas{}
			chaves := &memChaves{}
			efeito := &efeitoFake{}
			fila.porLer = []domain.Entrada{{
				ID:     "1-0",
				Campos: camposLote(t, caso.contagem, caso.recordID, caso.trace, caso.comTrace),
			}}
			s := &ConsumoReferencia{Fila: fila, Chaves: chaves, Efeito: efeito}
			saida, err := s.Executar(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if saida.Confirmacoes != 0 || fila.confirmados["1-0"] != 0 || efeito.concluidos != 0 {
				t.Fatalf("confirmações=%d fila=%d efeitos=%d", saida.Confirmacoes, fila.confirmados["1-0"], efeito.concluidos)
			}
			if len(efeito.refs) != 1 {
				t.Fatalf("entregas=%d", len(efeito.refs))
			}
			ref := efeito.refs[0]
			if ref.TenantID != "tenant-1" || ref.ProductID != "prod-1" || ref.Organizacao != "org-1" || ref.Stream != "logs" {
				t.Fatalf("consulta: %+v", ref)
			}
			if ref.BatchID != "batch-1" || ref.IdempotencyID != "batch-1" || ref.CorrelationID != "corr-1" || ref.EntradaID != "1-0" {
				t.Fatalf("identidade: %+v", ref)
			}
			if ref.MinTimestampMicros != 100 || ref.MaxTimestampMicros != 200 || ref.RecordCount != int64(caso.contagem) {
				t.Fatalf("faixa: %+v", ref)
			}
			if caso.contagem == 1 {
				if ref.RecordID == nil || *ref.RecordID != "rec-1" {
					t.Fatalf("recordId: %+v", ref.RecordID)
				}
			} else if ref.RecordID != nil {
				t.Fatalf("lote múltiplo levou recordId %q", *ref.RecordID)
			}
			if caso.comTrace {
				if ref.TraceID == nil || *ref.TraceID != caso.trace {
					t.Fatalf("trace: %+v", ref.TraceID)
				}
			} else if ref.TraceID != nil {
				t.Fatal("trace ausente foi inventado")
			}
			if len(chaves.reservas) != 1 || chaves.reservas[0] != "batch-1" {
				t.Fatalf("reservas=%v", chaves.reservas)
			}
			if chaves.estado["batch-1"] != domain.EstadoEmProcessamento {
				t.Fatalf("estado=%s", chaves.estado["batch-1"])
			}
		})
	}
}

func TestEntradaInvalidaConfirmaUmaVezSemConcluir(t *testing.T) {
	var buf bytes.Buffer
	anterior := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(anterior) })

	meta := `{"tenantId":"tenant-1","productId":"prod-1","correlationId":"corr-1","senha":"segredo-redis"}`
	body := `{"signal":"LOGS","mensagem":"CONTEUDO-DE-LOG-SECRETO"}`
	fila := &memEntradas{chaves: &memChaves{estado: map[string]string{"batch-1": domain.EstadoEmProcessamento}}}
	fila.porLer = []domain.Entrada{{
		ID: "9-0",
		Campos: map[string]string{
			"idempotencyId": "batch-1",
			"event":         "trace.ingested",
			"origin":        "collector",
			"schemaVersion": "1",
			"publishedAt":   "2026-11-08T12:00:00Z",
			"metadata":      meta,
			"body":          body,
		},
	}}
	efeito := &efeitoFake{concluido: true}
	s := &ConsumoReferencia{Fila: fila, Chaves: fila.chaves, Efeito: efeito}
	saida, err := s.Executar(context.Background())
	if err != nil || !saida.Invalida || saida.Confirmacoes != 1 {
		t.Fatalf("saída=%+v err=%v", saida, err)
	}
	if _, err := s.Executar(context.Background()); err != nil {
		t.Fatal(err)
	}
	if fila.confirmados["9-0"] != 1 || len(fila.encerradas) != 1 || fila.encerradas[0].idem != "batch-1" {
		t.Fatalf("confirmações=%d encerradas=%v", fila.confirmados["9-0"], fila.encerradas)
	}
	if fila.pendentes["9-0"] {
		t.Fatal("entrada inválida permaneceu pendente")
	}
	if fila.chaves.conclusoes != 0 || fila.chaves.estado["batch-1"] == domain.EstadoConcluido {
		t.Fatalf("conclusões=%d estado=%s", fila.chaves.conclusoes, fila.chaves.estado["batch-1"])
	}
	if fila.chaves.estado["batch-1"] != domain.EstadoAusente {
		t.Fatalf("reserva=%s", fila.chaves.estado["batch-1"])
	}
	if len(efeito.refs) != 0 {
		t.Fatal("inválida foi entregue")
	}
	registro := buf.String()
	if !strings.Contains(registro, "9-0") || !strings.Contains(registro, "batch-1") || !strings.Contains(registro, "event diferente de log.ingested") {
		t.Fatalf("registro: %s", registro)
	}
	if strings.Contains(registro, "segredo-redis") || strings.Contains(registro, "CONTEUDO-DE-LOG-SECRETO") || strings.Contains(registro, meta) || strings.Contains(registro, body) {
		t.Fatalf("registro expõe conteúdo: %s", registro)
	}
}

func TestSegundaEntradaEsperaEEfeitoConcluidoNaoRepete(t *testing.T) {
	campos := camposLote(t, 1, "rec-1", "", false)
	fila := &memEntradas{}
	chaves := &memChaves{}
	efeito := &efeitoFake{}
	s := &ConsumoReferencia{Fila: fila, Chaves: chaves, Efeito: efeito}
	fila.porLer = []domain.Entrada{
		{ID: "1-0", Campos: campos},
		{ID: "2-0", Campos: campos},
	}
	primeira, err := s.Executar(context.Background())
	if err != nil || primeira.Confirmacoes != 0 || efeito.chamadas != 1 || efeito.concluidos != 0 {
		t.Fatalf("primeira=%+v chamadas=%d err=%v", primeira, efeito.chamadas, err)
	}
	segunda, err := s.Executar(context.Background())
	if err != nil || segunda.Confirmacoes != 0 || efeito.chamadas != 1 || fila.confirmados["2-0"] != 0 || fila.adiados["2-0"] != 1 {
		t.Fatalf("segunda=%+v chamadas=%d confirmados=%v adiados=%v err=%v", segunda, efeito.chamadas, fila.confirmados, fila.adiados, err)
	}
	if fila.pendentes["2-0"] != true {
		t.Fatal("segunda entrada não permaneceu pendente")
	}

	efeito.concluido = true
	fila.cursor = 0
	fila.porLer = []domain.Entrada{{ID: "1-0", Campos: campos, Ociosa: true}}
	concluida, err := s.Executar(context.Background())
	if err != nil || concluida.Confirmacoes != 1 || concluida.Efeitos != 1 || efeito.concluidos != 1 {
		t.Fatalf("conclusão=%+v efeitos=%d err=%v", concluida, efeito.concluidos, err)
	}
	fila.cursor = 0
	fila.porLer = []domain.Entrada{{ID: "2-0", Campos: campos, Ociosa: true}}
	delete(fila.adiados, "2-0")
	repetida, err := s.Executar(context.Background())
	if err != nil || repetida.Confirmacoes != 1 || efeito.concluidos != 1 || fila.confirmados["2-0"] != 1 {
		t.Fatalf("repetida=%+v efeitos=%d confirmados=%v err=%v", repetida, efeito.concluidos, fila.confirmados, err)
	}
}

func TestQuedaAntesEDepoisDoEfeito(t *testing.T) {
	campos := camposLote(t, 1, "rec-1", "", false)
	fila := &memEntradas{}
	chaves := &memChaves{}
	s := &ConsumoReferencia{Fila: fila, Chaves: chaves, Efeito: nil}
	fila.porLer = []domain.Entrada{{ID: "1-0", Campos: campos}}
	saida, err := s.Executar(context.Background())
	if err != nil || saida.Efeitos != 0 || saida.Confirmacoes != 0 || chaves.estado["batch-1"] != domain.EstadoEmProcessamento {
		t.Fatalf("antes=%+v estado=%s err=%v", saida, chaves.estado["batch-1"], err)
	}

	efeito := &efeitoFake{concluido: true}
	s.Efeito = efeito
	fila.falhaConfirmar = true
	fila.cursor = 0
	fila.porLer = []domain.Entrada{{ID: "1-0", Campos: campos, Ociosa: true}}
	if _, err := s.Executar(context.Background()); err == nil {
		t.Fatal("esperava falha antes da confirmação")
	}
	if efeito.concluidos != 1 || chaves.estado["batch-1"] != domain.EstadoConcluido || fila.confirmados["1-0"] != 0 {
		t.Fatalf("efeitos=%d estado=%s confirmados=%d", efeito.concluidos, chaves.estado["batch-1"], fila.confirmados["1-0"])
	}
	fila.cursor = 0
	fila.porLer = []domain.Entrada{{ID: "1-0", Campos: campos}}
	saida, err = s.Executar(context.Background())
	if err != nil || saida.Confirmacoes != 1 || efeito.concluidos != 1 {
		t.Fatalf("reentrega=%+v efeitos=%d err=%v", saida, efeito.concluidos, err)
	}
}

func TestReassumirEntradaAbandonada(t *testing.T) {
	campos := camposLote(t, 1, "rec-1", "", false)
	fila := &memEntradas{}
	chaves := &memChaves{estado: map[string]string{"batch-1": domain.EstadoEmProcessamento}}
	efeito := &efeitoFake{concluido: true}
	s := &ConsumoReferencia{Fila: fila, Chaves: chaves, Efeito: efeito}
	fila.porLer = []domain.Entrada{{ID: "4-0", Campos: campos, Ociosa: true, Reivindicada: true}}
	saida, err := s.Executar(context.Background())
	if err != nil || saida.Efeitos != 1 || saida.Confirmacoes != 1 || chaves.estado["batch-1"] != domain.EstadoConcluido {
		t.Fatalf("saída=%+v estado=%s err=%v", saida, chaves.estado["batch-1"], err)
	}
}

func TestCorrelationEIdRedisNaoDeduplicam(t *testing.T) {
	um := camposLote(t, 1, "rec-1", "", false)
	outro := camposLote(t, 1, "rec-2", "", false)
	outro["idempotencyId"] = "batch-2"
	body := map[string]any{}
	if err := json.Unmarshal([]byte(outro["body"]), &body); err != nil {
		t.Fatal(err)
	}
	aberto := body["openobserve"].(map[string]any)
	aberto["batchId"] = "batch-2"
	bruto, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	outro["body"] = string(bruto)
	fila := &memEntradas{porLer: []domain.Entrada{
		{ID: "1-0", Campos: um},
		{ID: "2-0", Campos: outro},
	}}
	chaves := &memChaves{}
	efeito := &efeitoFake{concluido: true}
	s := &ConsumoReferencia{Fila: fila, Chaves: chaves, Efeito: efeito}
	if _, err := s.Executar(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Executar(context.Background()); err != nil {
		t.Fatal(err)
	}
	if efeito.concluidos != 2 || len(chaves.reservas) != 2 || chaves.reservas[0] != "batch-1" || chaves.reservas[1] != "batch-2" {
		t.Fatalf("efeitos=%d reservas=%v", efeito.concluidos, chaves.reservas)
	}
}

type memEntradas struct {
	porLer         []domain.Entrada
	cursor         int
	confirmados    map[string]int
	adiados        map[string]int
	pendentes      map[string]bool
	encerradas     []parEncerrado
	falhaConfirmar bool
	chaves         *memChaves
}

type parEncerrado struct {
	id   string
	idem string
}

func (m *memEntradas) Ler(context.Context) (domain.Entrada, bool, error) {
	for m.cursor < len(m.porLer) {
		item := m.porLer[m.cursor]
		m.cursor++
		if m.confirmados[item.ID] > 0 {
			continue
		}
		if m.pendentes == nil {
			m.pendentes = map[string]bool{}
		}
		m.pendentes[item.ID] = true
		return item, true, nil
	}
	return domain.Entrada{}, false, nil
}

func (m *memEntradas) Confirmar(_ context.Context, id string) error {
	if m.falhaConfirmar {
		m.falhaConfirmar = false
		return errors.New("redis indisponível")
	}
	if m.confirmados == nil {
		m.confirmados = map[string]int{}
	}
	m.confirmados[id]++
	delete(m.pendentes, id)
	return nil
}

func (m *memEntradas) Adiar(_ context.Context, id string) error {
	if m.adiados == nil {
		m.adiados = map[string]int{}
	}
	m.adiados[id]++
	return nil
}

func (m *memEntradas) EncerrarInvalida(ctx context.Context, id, idem string) error {
	m.encerradas = append(m.encerradas, parEncerrado{id: id, idem: idem})
	if m.chaves != nil && idem != "" {
		if err := m.chaves.RemoverReserva(ctx, idem); err != nil {
			return err
		}
	}
	return m.Confirmar(ctx, id)
}

type memChaves struct {
	estado     map[string]string
	reservas   []string
	conclusoes int
	remocoes   int
}

func (m *memChaves) Reservar(_ context.Context, id string) (bool, error) {
	if m.estado == nil {
		m.estado = map[string]string{}
	}
	m.reservas = append(m.reservas, id)
	if m.estado[id] == "" || m.estado[id] == domain.EstadoAusente {
		m.estado[id] = domain.EstadoEmProcessamento
		return true, nil
	}
	return false, nil
}

func (m *memChaves) Concluir(_ context.Context, id string) error {
	if m.estado == nil {
		m.estado = map[string]string{}
	}
	m.conclusoes++
	if m.estado[id] != domain.EstadoConcluido {
		m.estado[id] = domain.EstadoConcluido
	}
	return nil
}

func (m *memChaves) Consultar(_ context.Context, id string) (string, error) {
	if m.estado[id] == "" {
		return domain.EstadoAusente, nil
	}
	return m.estado[id], nil
}

func (m *memChaves) RemoverReserva(_ context.Context, id string) error {
	if m.estado == nil {
		m.estado = map[string]string{}
	}
	m.remocoes++
	if m.estado[id] == domain.EstadoEmProcessamento {
		m.estado[id] = domain.EstadoAusente
	}
	return nil
}

type efeitoFake struct {
	refs       []domain.Referencia
	chamadas   int
	concluidos int
	concluido  bool
	err        error
}

func (e *efeitoFake) Executar(_ context.Context, ref domain.Referencia) (bool, error) {
	e.chamadas++
	e.refs = append(e.refs, ref)
	if e.err != nil {
		return false, e.err
	}
	if e.concluido {
		e.concluidos++
	}
	return e.concluido, nil
}

func camposLote(t *testing.T, contagem int, recordID any, trace string, comTrace bool) map[string]string {
	t.Helper()
	meta := map[string]any{
		"tenantId":      "tenant-1",
		"productId":     "prod-1",
		"correlationId": "corr-1",
	}
	if comTrace {
		meta["traceId"] = trace
	}
	body := map[string]any{
		"signal":    "LOGS",
		"productId": "prod-1",
		"openobserve": map[string]any{
			"organization":       "org-1",
			"stream":             "logs",
			"batchId":            "batch-1",
			"recordId":           recordID,
			"minTimestampMicros": 100,
			"maxTimestampMicros": 200,
		},
		"recordCount": contagem,
	}
	mb, err := json.Marshal(meta)
	if err != nil {
		t.Fatal(err)
	}
	bb, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	return map[string]string{
		"idempotencyId": "batch-1",
		"event":         "log.ingested",
		"origin":        "collector",
		"schemaVersion": "1",
		"publishedAt":   "2026-11-08T12:00:00Z",
		"metadata":      string(mb),
		"body":          string(bb),
	}
}
