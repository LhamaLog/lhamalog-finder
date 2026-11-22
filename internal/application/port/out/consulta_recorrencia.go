package out

import (
	"context"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

type PedidoRecorrencia struct {
	BatchID, Host string
	Chave         domain.IdentidadeAchado
	LimiteLotes   int
}

// ConsultarHistorico deve devolver os últimos lotes, do mais recente ao antigo,
// incluindo os lotes sem equivalência para aplicar a janela de 500 lotes distintos.
// O adaptador aplica DF-FINDER-2026-11-22 e mantém snapshot por referência.
type ConsultaHistorico interface {
	ConsultarHistorico(context.Context, PedidoRecorrencia) ([]domain.EvidenciaRecorrencia, error)
}
type ConsultaRecorrencia interface {
	ConsultarRecorrencia(context.Context, PedidoRecorrencia) (domain.MarcaRecorrente, error)
}
type AvisoFalhaOpenObserve interface {
	AvisarFalhaOpenObserve(ctx context.Context, batchID, motivo string) error
}
type FalhaRecorrencia struct{ Temporaria bool }

func (FalhaRecorrencia) Error() string { return "consulta_openobserve_falhou" }

func (f FalhaRecorrencia) Retentavel() bool { return f.Temporaria }
