package handoff

import (
	"context"
	"errors"
	entrada "github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// PreparadorAnalise é satisfeito pela 008: segurança/classificação → agrupamento
// → consulta por achado/recorrência. A 009 não duplica essas consultas ou regras.
type PreparadorAnalise interface {
	Preparar(context.Context, entrada.PedidoAnalise) (domain.ResumoAnalise, error)
}
type OrquestradorEntrega struct {
	Preparador PreparadorAnalise
	Entrega    Entrega
}

func (o OrquestradorEntrega) Executar(ctx context.Context, p entrada.PedidoAnalise) (ResultadoEntrega, error) {
	if !p.Completo {
		return ResultadoEntrega{}, errors.New("lote_incompleto")
	}
	if o.Preparador == nil {
		return ResultadoEntrega{}, errors.New("analise_indisponivel")
	}
	resumo, err := o.Preparador.Preparar(ctx, p)
	if err != nil {
		return ResultadoEntrega{}, err
	}
	if resumo.BatchID != p.Referencia.BatchID || resumo.TenantID != p.Referencia.TenantID || resumo.ProductID != p.Referencia.ProductID {
		return ResultadoEntrega{}, errors.New("identidade_incompleta")
	}
	pendente, err := NovaPendente(resumo)
	if err != nil {
		return ResultadoEntrega{}, err
	}
	return o.Entrega.Enviar(ctx, pendente)
}
