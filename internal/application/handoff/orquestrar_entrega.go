package handoff

import (
	"context"
	"errors"
	entrada "github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// PreparadorAnalise é satisfeito pela 008: segurança/classificação → agrupamento
// → consulta por achado/recorrência. A 009 não duplica essas consultas ou regras.
type PreparadorAnalise interface {
	Montar(context.Context, entrada.PedidoAnalise) (domain.ResumoAnalise, error)
	Agrupar(context.Context, domain.ResumoAnalise) (domain.ResumoAnalise, error)
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
	local, err := o.Preparador.Montar(ctx, p)
	if err != nil {
		return ResultadoEntrega{}, err
	}
	r, err := o.agruparEEnviar(ctx, p, local)
	if errors.Is(err, out.ErrRevisaoDesatualizada) {
		// O 409 garante que a escrita ainda não ocorreu. Recalcular uma vez
		// mantém a chave e usa a revisão de uma nova consulta por achado.
		return o.agruparEEnviar(ctx, p, local)
	}
	return r, err
}

func (o OrquestradorEntrega) agruparEEnviar(ctx context.Context, p entrada.PedidoAnalise, local domain.ResumoAnalise) (ResultadoEntrega, error) {
	if err := ctx.Err(); err != nil {
		return ResultadoEntrega{}, err
	}
	resumo, err := o.Preparador.Agrupar(ctx, local)
	if err != nil {
		return ResultadoEntrega{}, err
	}
	if resumo.BatchID != p.Referencia.BatchID || resumo.TenantID != p.Referencia.TenantID || resumo.ProductID != p.Referencia.ProductID || resumo.ExecID != p.ExecID {
		return ResultadoEntrega{}, errors.New("identidade_incompleta")
	}
	pendente, err := NovaPendente(resumo)
	if err != nil {
		return ResultadoEntrega{}, err
	}
	return o.Entrega.Enviar(ctx, pendente)
}
