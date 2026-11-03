package out

import (
	"context"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// FonteApontamentos entrega os apontamentos desta execução, ou nenhum.
// Não lê a fila Redis. Não nomeia arquivo, variável ou tabela.
type FonteApontamentos interface {
	Obter(ctx context.Context) ([]domain.Codigo, error)
}
