package application

import (
	"context"
	"errors"
	"log/slog"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// ServicoRecuperacao lê o lote no OpenObserve e informa se o efeito está concluído.
// Não confirma a entrada Redis.
type ServicoRecuperacao struct {
	Busca  out.Busca
	Acesso domain.ConfiguracaoAcesso
}

var _ out.Efeito = (*ServicoRecuperacao)(nil)

// Executar devolve a marca de efeito concluído. Erro só acompanha o contexto cancelado.
func (s *ServicoRecuperacao) Executar(ctx context.Context, ref domain.Referencia) (bool, error) {
	resultado, err := s.Recuperar(ctx, ref)
	if err != nil {
		return false, err
	}
	return resultado.EfeitoConcluido, nil
}

// Recuperar percorre as páginas do stream da referência e fecha o lote.
// Não há teto de páginas: a leitura segue enquanto a página volta cheia.
func (s *ServicoRecuperacao) Recuperar(ctx context.Context, ref domain.Referencia) (domain.ResultadoRecuperacao, error) {
	if err := ctx.Err(); err != nil {
		return domain.ResultadoRecuperacao{}, err
	}
	if !s.Acesso.Pronta() {
		resultado := domain.NovoResultadoRecuperacao(ref)
		resultado.Classificacao = domain.ClasseConfigAusente
		return s.publicar(ref, resultado), nil
	}
	if s.Acesso.CampoTenant == "" || s.Acesso.CampoProduto == "" {
		aval := domain.Avaliacao{Isolamento: true}
		return s.publicar(ref, domain.Fechar(ref, aval, nil)), nil
	}
	sql, inicio, fim, err := domain.ConsultaSQL(ref, s.Acesso)
	if err != nil {
		aval := domain.Avaliacao{Acesso: domain.AcessoResposta}
		return s.publicar(ref, domain.Fechar(ref, aval, nil)), nil
	}

	var docs []domain.DocumentoFisico
	from := 0
	paginas := 0
	for {
		if err := ctx.Err(); err != nil {
			return domain.ResultadoRecuperacao{}, err
		}
		pedido := out.PedidoPagina{
			URL:         s.Acesso.URL,
			Usuario:     s.Acesso.Usuario,
			Senha:       s.Acesso.Senha,
			Timeout:     s.Acesso.Timeout,
			Organizacao: ref.Organizacao,
			SQL:         sql,
			Inicio:      inicio,
			Fim:         fim,
			From:        from,
			Size:        s.Acesso.TamanhoPagina,
		}
		hits, err := s.Busca.Consultar(ctx, pedido)
		if err != nil {
			if ctx.Err() != nil {
				return domain.ResultadoRecuperacao{}, ctx.Err()
			}
			return s.publicar(ref, s.falhaDePagina(ref, docs, paginas, err)), nil
		}
		docs = append(docs, hits...)
		paginas++
		if len(hits) != s.Acesso.TamanhoPagina {
			break
		}
		from += s.Acesso.TamanhoPagina
	}
	return s.publicar(ref, s.fecharLeitura(ref, docs, paginas)), nil
}

func (s *ServicoRecuperacao) falhaDePagina(ref domain.Referencia, docs []domain.DocumentoFisico, paginas int, err error) domain.ResultadoRecuperacao {
	aval, _ := s.medir(ref, docs)
	aval.ContagemCalculada = false
	aval.Encontrado = 0
	aval.Acesso = tipoAcesso(err)
	aval.PaginasLidas = paginas
	return domain.Fechar(ref, aval, nil)
}

func (s *ServicoRecuperacao) fecharLeitura(ref domain.Referencia, docs []domain.DocumentoFisico, paginas int) domain.ResultadoRecuperacao {
	aval, registros := s.medir(ref, docs)
	aval.PaginasLidas = paginas
	return domain.Fechar(ref, aval, registros)
}

func (s *ServicoRecuperacao) medir(ref domain.Referencia, docs []domain.DocumentoFisico) (domain.Avaliacao, []domain.DocumentoFisico) {
	aval := domain.Avaliacao{
		Isolamento:        domain.Isolado(docs, s.Acesso.CampoTenant, s.Acesso.CampoProduto, ref.TenantID, ref.ProductID),
		ContagemCalculada: true,
	}
	proprios := domain.Proprios(docs, s.Acesso.CampoTenant, s.Acesso.CampoProduto, ref.TenantID, ref.ProductID)
	registros, divergente, semID := domain.Deduplicar(proprios, copiaDe(s.Acesso))
	aval.Divergente = divergente
	aval.SemIdentificador = semID
	aval.Encontrado = int64(len(registros))
	return aval, registros
}

func copiaDe(acesso domain.ConfiguracaoAcesso) []string {
	if acesso.CamposCopia == nil {
		return nil
	}
	return *acesso.CamposCopia
}

func tipoAcesso(err error) domain.FalhaAcesso {
	var falha *out.ErroBusca
	if !errors.As(err, &falha) || falha == nil {
		return domain.AcessoIndisponivel
	}
	switch falha.Tipo {
	case out.FalhaCredencial:
		return domain.AcessoCredencial
	case out.FalhaResposta:
		return domain.AcessoResposta
	case out.FalhaTimeout:
		return domain.AcessoTimeout
	case out.FalhaPagina:
		return domain.AcessoPagina
	default:
		return domain.AcessoIndisponivel
	}
}

func (s *ServicoRecuperacao) publicar(ref domain.Referencia, resultado domain.ResultadoRecuperacao) domain.ResultadoRecuperacao {
	args := []any{
		"classificacao", string(resultado.Classificacao),
		"idempotencyId", ref.IdempotencyID,
		"id", ref.EntradaID,
		"efeitoConcluido", resultado.EfeitoConcluido,
	}
	if resultado.ContagemCalculada {
		args = append(args, "esperado", resultado.RecordCount, "encontrado", resultado.Encontrado)
	}
	slog.Info("recuperação do lote", args...)
	return resultado
}
