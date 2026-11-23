package handoff_test

import (
	"context"
	"errors"
	"testing"

	"github.com/LhamaLog/lhamalog-finder/internal/application/handoff"
	entrada "github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

type concorrenciaTeste struct {
	nPreparar, nAgrupar, nEnviar, nConcluir int
	falhas                                  []error
}

func (m *concorrenciaTeste) Montar(_ context.Context, p entrada.PedidoAnalise) (domain.ResumoAnalise, error) {
	m.nPreparar++
	return domain.ResumoAnalise{BatchID: p.Referencia.BatchID, TenantID: p.Referencia.TenantID, ProductID: p.Referencia.ProductID, ExecID: p.ExecID, VersaoRegras: "v"}, nil
}
func (m *concorrenciaTeste) Enviar(_ context.Context, r domain.ResumoAnalise) (bool, error) {
	m.nEnviar++
	if m.nEnviar <= len(m.falhas) {
		return false, m.falhas[m.nEnviar-1]
	}
	return true, nil
}
func (m *concorrenciaTeste) ConcluirAnalise(context.Context, domain.ChaveResumo) error {
	m.nConcluir++
	return nil
}
func TestHandoffRecalculaUmaVezApenasEmRevisaoDesatualizada(t *testing.T) {
	for _, caso := range []struct {
		nome       string
		falhas     []error
		tentativas int
		conclui    bool
	}{
		{"conflito resolvido", []error{out.ErrRevisaoDesatualizada}, 2, true},
		{"conflito persistente", []error{out.ErrRevisaoDesatualizada, out.ErrRevisaoDesatualizada}, 2, false},
		{"falha indeterminada", []error{errors.New("falha")}, 1, false},
	} {
		t.Run(caso.nome, func(t *testing.T) {
			m := &concorrenciaTeste{falhas: caso.falhas}
			o := handoff.OrquestradorEntrega{Preparador: m, Entrega: handoff.Entrega{Envio: m, Finalizar: m}}
			p := entrada.PedidoAnalise{Completo: true, Referencia: domain.Referencia{BatchID: "b", TenantID: "t", ProductID: "p"}}
			r, err := o.Executar(context.Background(), p)
			if r.Concluida != caso.conclui || (err == nil) != caso.conclui || m.nPreparar != 1 || m.nAgrupar != caso.tentativas || m.nEnviar != caso.tentativas || (m.nConcluir == 1) != caso.conclui {
				t.Fatalf("resultado %v, erro %v, chamadas %+v", r, err, m)
			}
		})
	}
}

func (p *concorrenciaTeste) Agrupar(_ context.Context, r domain.ResumoAnalise) (domain.ResumoAnalise, error) {
	p.nAgrupar++
	return r, nil
}
