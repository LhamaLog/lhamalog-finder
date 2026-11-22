package application

import (
	"context"
	"errors"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

var ErrCatalogoSeguranca = errors.New("catalogo_seguranca_indisponivel")
var ErrVersaoCatalogo = errors.New("versao_catalogo_divergente")

type ValidacaoSeguranca interface {
	ValidarSeguranca(context.Context, []string) domain.ResultadoValidarSeguranca
}
type ValidacaoSegurancaFunc func(context.Context, []string) domain.ResultadoValidarSeguranca

func (f ValidacaoSegurancaFunc) ValidarSeguranca(ctx context.Context, s []string) domain.ResultadoValidarSeguranca {
	return f(ctx, s)
}

type Seguranca struct{ Catalogo domain.CatalogoSeguranca }

// InicializarSeguranca é a barreira anterior à prontidão do pipeline.
// O bootstrap do processo aguarda o adaptador/artefato P-010-2; não o simula.
func InicializarSeguranca(ctx context.Context, fonte out.CatalogoPadroes, esperada string) (Seguranca, error) {
	if fonte == nil {
		return Seguranca{}, ErrCatalogoSeguranca
	}
	cat, err := fonte.Carregar(ctx)
	if err != nil {
		return Seguranca{}, ErrCatalogoSeguranca
	}
	if esperada == "" || cat.Versao == "" || cat.Versao != esperada {
		return Seguranca{}, ErrVersaoCatalogo
	}
	if cat.Transformar == nil {
		return Seguranca{}, ErrCatalogoSeguranca
	}
	return Seguranca{Catalogo: cat}, nil
}
func (s Seguranca) ValidarSeguranca(ctx context.Context, linhas []string) domain.ResultadoValidarSeguranca {
	falhar := func() domain.ResultadoValidarSeguranca {
		return domain.ResultadoValidarSeguranca{Erro: ErrCatalogoSeguranca}
	}
	if s.Catalogo.Versao == "" || s.Catalogo.Transformar == nil {
		return falhar()
	}
	r := domain.ResultadoValidarSeguranca{LinhasMascaradas: make([]domain.LinhaLogMascarada, 0, len(linhas))}
	for i, linha := range linhas {
		if ctx.Err() != nil {
			return falhar()
		}
		segura, n, err := s.Catalogo.Transformar(linha)
		if err != nil || n < 0 {
			return falhar()
		}
		segura.IndiceLinha = i
		r.LinhasMascaradas = append(r.LinhasMascaradas, segura)
		r.TotalSensivel += n
	}
	return r
}

// MontarContexto chama o scanner exatamente uma vez, antes da análise.
// 010-16 C: erro retira todos os trechos, inclusive saída parcial do scanner.
func MontarContexto(ctx context.Context, batch string, linhas []string, versao string, v ValidacaoSeguranca) domain.ContextoLoteMascarado {
	c := domain.ContextoLoteMascarado{BatchID: batch, VersaoCatalogo: versao}
	if v == nil {
		c.FalhaSeguranca = true
		return c
	}
	r := v.ValidarSeguranca(ctx, append([]string(nil), linhas...))
	if r.Erro != nil || len(r.LinhasMascaradas) != len(linhas) || r.TotalSensivel < 0 {
		c.FalhaSeguranca = true
		return c
	}
	c.Linhas = append([]domain.LinhaLogMascarada(nil), r.LinhasMascaradas...)
	c.DadosSensiveisIdentificados = r.TotalSensivel
	return c
}
