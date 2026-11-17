package openobserve

import (
	"context"
	"errors"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/application/recorrencia"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

var ErrQueryHistoricaPendente = errors.New("query_historica_openobserve_pendente")
var ErrAvisoCorePendente = errors.New("consulta_openobserve_esgotada_aviso_core_pendente")

// HistoricoHTTPPendente não usa a busca de lote 005 como query histórica inventada.
type HistoricoHTTPPendente struct{}

func (HistoricoHTTPPendente) ConsultarHistorico(context.Context, out.PedidoRecorrencia) ([]domain.EvidenciaRecorrencia, error) {
	return nil, ErrQueryHistoricaPendente
}

type Recorrencia struct {
	Historico out.ConsultaHistorico
	Aviso     out.AvisoFalhaOpenObserve
}

func (s Recorrencia) ConsultarRecorrencia(ctx context.Context, p out.PedidoRecorrencia) (domain.MarcaRecorrente, error) {
	if s.Historico == nil {
		return domain.MarcaRecorrente{}, ErrQueryHistoricaPendente
	}
	p.LimiteLotes = recorrencia.JanelaLotes
	for tentativa := 0; tentativa < 3; tentativa++ {
		if err := ctx.Err(); err != nil {
			return domain.MarcaRecorrente{}, err
		}
		historico, err := s.Historico.ConsultarHistorico(ctx, p)
		if err == nil {
			v := recorrencia.Avaliar(historico, p.Chave, p.Host)
			return domain.MarcaRecorrente{Valor: &v}, nil
		}
		if ctx.Err() != nil {
			return domain.MarcaRecorrente{}, ctx.Err()
		}
		var falha interface{ Retentavel() bool }
		if !errors.As(err, &falha) || !falha.Retentavel() {
			if errors.Is(err, ErrQueryHistoricaPendente) {
				return domain.MarcaRecorrente{}, ErrQueryHistoricaPendente
			}
			return domain.MarcaRecorrente{}, out.FalhaRecorrencia{Temporaria: false}
		}
	}
	// Nenhum ACK, DLQ, WebSocket ou recorrente=false por degradação.
	// O cliente HTTP Core permanece bloqueado ADR_051; esta é só a porta semântica.
	if s.Aviso == nil {
		return domain.MarcaRecorrente{}, ErrAvisoCorePendente
	}
	if err := s.Aviso.AvisarFalhaOpenObserve(ctx, p.BatchID, "consulta_openobserve_indisponivel"); err != nil {
		return domain.MarcaRecorrente{}, ErrAvisoCorePendente
	}
	return domain.MarcaRecorrente{}, out.FalhaRecorrencia{Temporaria: true}
}
