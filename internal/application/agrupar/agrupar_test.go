package agrupar

import (
	"context"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"testing"
	"time"
)

func TestParticionarIsolaChaveEContaTimestampLog(t *testing.T) {
	a := domain.LogClassificado{Apontamento: "a", Classe: domain.ERROR, TenantID: "t", ProductID: "p", ServiceName: "svc", Host: "h", Mensagem: "falha 10", Timestamp: time.Unix(10, 0)}
	b := a
	b.Apontamento = "b"
	b.Mensagem = "FALHA 20"
	b.Timestamp = time.Unix(5, 0)
	c := a
	c.TenantID = "outro"
	r, err := ParticionarLote([]domain.LogClassificado{a, b, c})
	if err != nil || len(r) != 2 {
		t.Fatalf("particoes=%d err=%v", len(r), err)
	}
	if r[0].ContagemLote != 2 || !r[0].PrimeiraNoLote.Equal(b.Timestamp) || !r[0].UltimaNoLote.Equal(a.Timestamp) {
		t.Fatal("metricas do lote")
	}
}
func TestFusaoEntreLotesNaoAlteraStatus(t *testing.T) {
	a := domain.AchadoAnalise{Tipo: domain.INFO, Assinatura: "s", Contagem: 2, PrimeiraOcorrencia: time.Unix(10, 0), UltimaOcorrencia: time.Unix(20, 0)}
	e := domain.EstadoOcorrencia{ContagemAcumulada: 5, PrimeiraGlobal: time.Unix(1, 0), UltimaGlobal: time.Unix(15, 0), Status: "FECHADA"}
	r, err := Fundir(a, []domain.EstadoOcorrencia{e})
	if err != nil || r.Contagem != 7 || r.StatusAnterior != "FECHADA" || !r.PrimeiraOcorrencia.Equal(e.PrimeiraGlobal) || !r.UltimaOcorrencia.Equal(a.UltimaOcorrencia) {
		t.Fatal("fusao")
	}
	if _, err := Fundir(a, []domain.EstadoOcorrencia{e, e}); err == nil {
		t.Fatal("duplicata de ocorrencia aceita")
	}
}

type montadorTeste struct{ r domain.ResumoAnalise }

func (m montadorTeste) Montar(context.Context, PedidoAnalise) (domain.ResumoAnalise, error) {
	return m.r, nil
}

// A interface usa o alias da porta de entrada comum; nenhuma leitura Redis.
type consultaTeste struct{ n int }

func (c *consultaTeste) ConsultarAchado(context.Context, domain.IdentidadeAchado) ([]domain.EstadoOcorrencia, error) {
	c.n++
	return nil, nil
}

type recorrenciaTeste struct{ falhar bool }

func (c recorrenciaTeste) ConsultarRecorrencia(context.Context, PedidoRecorrencia) (domain.MarcaRecorrente, error) {
	if c.falhar {
		return domain.MarcaRecorrente{}, ErrRecorrenciaIndisponivel
	}
	v := true
	return domain.MarcaRecorrente{Valor: &v}, nil
}

type envioTeste struct{ n int }

func (e *envioTeste) Enviar(context.Context, domain.ResumoAnalise) (bool, error) {
	e.n++
	return true, nil
}
func TestOrquestrarConsultaPorAchadoEnviaUmResumo(t *testing.T) {
	c := &consultaTeste{}
	e := &envioTeste{}
	m := montadorTeste{r: domain.ResumoAnalise{BatchID: "b", TenantID: "t", ProductID: "p", Achados: []domain.AchadoAnalise{{Tipo: domain.INFO, Assinatura: "s1", HostName: "h"}, {Tipo: domain.ERROR, Assinatura: "s2", HostName: "h"}}}}
	o := Orquestrador{Montador: m, Consulta: c, Recorrencia: recorrenciaTeste{}, Envio: e}
	r, err := o.Executar(context.Background(), PedidoAnalise{Completo: true})
	if err != nil || !r.Concluida || c.n != 2 || e.n != 1 {
		t.Fatalf("orquestracao: %v consultas=%d envios=%d", err, c.n, e.n)
	}
	o.Recorrencia = recorrenciaTeste{falhar: true}
	r, err = o.Executar(context.Background(), PedidoAnalise{Completo: true})
	if err == nil || r.Concluida || e.n != 1 {
		t.Fatal("recorrencia indisponivel virou sucesso")
	}
	r, err = o.Executar(context.Background(), PedidoAnalise{Completo: false})
	if err == nil || r.Concluida || e.n != 1 {
		t.Fatal("lote incompleto chegou ao handoff")
	}
}

func TestConsultaMemoriaCopiaEstadoEContaChamadas(t *testing.T) {
	chave := domain.IdentidadeAchado{Assinatura: "a", Classe: domain.ERROR, TenantID: "t", ProductID: "p"}
	c := NovaConsultaMemoria(map[domain.IdentidadeAchado][]domain.EstadoOcorrencia{chave: {{ContagemAcumulada: 3, Status: "ABERTA"}}})
	r, err := c.ConsultarAchado(context.Background(), chave)
	if err != nil || len(r) != 1 {
		t.Fatal("estado ausente")
	}
	r[0].ContagemAcumulada = 99
	outro, err := c.ConsultarAchado(context.Background(), chave)
	if err != nil || outro[0].ContagemAcumulada != 3 || c.Chamadas(chave) != 2 {
		t.Fatal("estado compartilhado ou contagem errada")
	}
}
