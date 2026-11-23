package handoff_test

import (
	"context"
	"github.com/LhamaLog/lhamalog-finder/internal/adapter/out/modocorrencia"
	"github.com/LhamaLog/lhamalog-finder/internal/adapter/out/openobserve"
	"github.com/LhamaLog/lhamalog-finder/internal/application"
	"github.com/LhamaLog/lhamalog-finder/internal/application/agrupar"
	"github.com/LhamaLog/lhamalog-finder/internal/application/handoff"
	entrada "github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"strings"
	"testing"
	"time"
)

type moduloFluxo struct {
	eventos              *[]string
	consulta, envio, fim int
	t                    *testing.T
}

func (m *moduloFluxo) ConsultarAchado(context.Context, domain.IdentidadeAchado) ([]domain.EstadoOcorrencia, error) {
	*m.eventos = append(*m.eventos, "consulta")
	m.consulta++
	return []domain.EstadoOcorrencia{{ContagemAcumulada: 5, PrimeiraGlobal: time.Unix(1, 0), UltimaGlobal: time.Unix(2, 0), Status: "FECHADA", Revisao: "rev-5"}}, nil
}
func (m *moduloFluxo) Enviar(ctx context.Context, r domain.ResumoAnalise) (bool, error) {
	*m.eventos = append(*m.eventos, "envio")
	m.envio++
	if len(r.Achados) != 1 || r.Achados[0].Contagem != 8 || r.Achados[0].StatusAnterior != "FECHADA" || !r.Achados[0].Recorrente || r.Achados[0].RevisaoConsultada == nil || *r.Achados[0].RevisaoConsultada != "rev-5" {
		m.t.Fatal("fusao/recorrencia nao aplicada antes do envio")
	}
	b, err := modocorrencia.SerializarResumo(r)
	if err != nil || strings.Contains(string(b), "BRUTO_SENTINELA") {
		m.t.Fatal("handoff nao respeita fronteira segura")
	}
	return true, nil
}
func (m *moduloFluxo) ConcluirAnalise(context.Context, domain.ChaveResumo) error {
	*m.eventos = append(*m.eventos, "fim")
	m.fim++
	return nil
}

type historicoFluxo struct{ eventos *[]string }

func (h historicoFluxo) ConsultarHistorico(ctx context.Context, p out.PedidoRecorrencia) ([]domain.EvidenciaRecorrencia, error) {
	*h.eventos = append(*h.eventos, "historico")
	return []domain.EvidenciaRecorrencia{{BatchID: "a", Host: p.Host, Chave: p.Chave}, {BatchID: "b", Host: p.Host, Chave: p.Chave}, {BatchID: "c", Host: p.Host, Chave: p.Chave}}, nil
}
func TestFluxoLogico010006008009(t *testing.T) {
	eventos := []string{}
	m := &moduloFluxo{eventos: &eventos, t: t}
	scanner := application.ValidacaoSegurancaFunc(func(ctx context.Context, s []string) domain.ResultadoValidarSeguranca {
		eventos = append(eventos, "seguranca")
		r := domain.ResultadoValidarSeguranca{}
		for i := range s {
			r.LinhasMascaradas = append(r.LinhasMascaradas, domain.LinhaLogMascarada{Texto: "seguro INFO", IndiceLinha: i})
		}
		return r
	})
	a := application.ServicoAnalise{Seguranca: scanner}
	g := agrupar.Orquestrador{Montador: a, Consulta: m, Recorrencia: openobserve.Recorrencia{Historico: historicoFluxo{eventos: &eventos}}}
	h := handoff.OrquestradorEntrega{Preparador: g, Entrega: handoff.Entrega{Envio: m, Finalizar: m}}
	p := entrada.PedidoAnalise{Referencia: domain.Referencia{BatchID: "b", TenantID: "t", ProductID: "p", RecordCount: 3}, ProdutoBody: "p", SchemaVersion: "1", Completo: true, VersaoRegras: "v1", VersaoCatalogo: "teste"}
	for _, id := range []domain.Codigo{"r1", "r2", "r3"} {
		p.Logs = append(p.Logs, entrada.LogAnalise{Apontamento: id, Mensagem: "BRUTO_SENTINELA", ServiceName: "svc", Host: "h", Timestamp: time.Unix(10, 0)})
	}
	r, err := h.Executar(context.Background(), p)
	if err != nil || !r.Concluida || m.consulta != 1 || m.envio != 1 || m.fim != 1 || strings.Join(eventos, ",") != "seguranca,consulta,historico,envio,fim" {
		t.Fatalf("fluxo: %v eventos=%v", err, eventos)
	}
	// Integração lógica com portas fake; não prova API Módulo/Core nem catálogo oficial.
}
