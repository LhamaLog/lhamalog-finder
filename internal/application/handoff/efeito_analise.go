package handoff

import (
	"context"
	"errors"

	entrada "github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// RecuperadorLote expõe os documentos deduplicados pela recuperação da spec 005.
type RecuperadorLote interface {
	Recuperar(context.Context, domain.Referencia) (domain.ResultadoRecuperacao, error)
}

// EfeitoAnalise liga o consumo de log.ingested v1 à entrega do resumo.
// ConsumoReferencia continua responsável pela reserva e confirmação Redis.
type EfeitoAnalise struct {
	Recuperacao    RecuperadorLote
	Mapear         func(domain.DocumentoFisico) (entrada.LogAnalise, error)
	Entrega        OrquestradorEntrega
	VersaoRegras   string
	VersaoCatalogo string
	Catalogo       []domain.RegraClassificacao
}

var _ out.Efeito = EfeitoAnalise{}

func (e EfeitoAnalise) Executar(ctx context.Context, ref domain.Referencia) (bool, error) {
	if e.Recuperacao == nil || e.Mapear == nil {
		return false, errors.New("efeito_analise_indisponivel")
	}
	lote, err := e.Recuperacao.Recuperar(ctx, ref)
	if err != nil {
		return false, err
	}
	if lote.Classificacao != domain.ClasseLoteCompleto {
		// Recuperação incompleta continua pendente. Rejeições definitivas
		// conservam a decisão da 005, sem produzir um resumo de sucesso.
		return lote.EfeitoConcluido, nil
	}
	if !lote.EfeitoConcluido || ref.RecordCount <= 0 || int64(len(lote.Registros)) != ref.RecordCount {
		return false, errors.New("lote_incompleto")
	}

	// InterpretarEnvelope já validou a versão e a igualdade entre produto
	// do corpo e do metadata. log.ingested v1 não fornece execId.
	pedido := entrada.PedidoAnalise{
		Referencia:     ref,
		ProdutoBody:    ref.ProductID,
		SchemaVersion:  "1",
		Completo:       true,
		VersaoRegras:   e.VersaoRegras,
		VersaoCatalogo: e.VersaoCatalogo,
		Catalogo:       e.Catalogo,
		Logs:           make([]entrada.LogAnalise, 0, len(lote.Registros)),
	}
	for _, documento := range lote.Registros {
		if err := ctx.Err(); err != nil {
			return false, err
		}
		log, err := e.Mapear(documento)
		if err != nil {
			return false, err
		}
		pedido.Logs = append(pedido.Logs, log)
	}

	entrega, err := e.Entrega.Executar(ctx, pedido)
	if err != nil {
		return false, err
	}
	return entrega.Confirmada && entrega.Concluida, nil
}
