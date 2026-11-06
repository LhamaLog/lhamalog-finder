package out

import (
	"context"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// Disponibilidade guarda o código para a etapa futura, uma vez.
// Esta porta não abre o OpenObserve. O lugar não tem nome.
type Disponibilidade interface {
	Disponibilizar(ctx context.Context, codigo domain.Codigo) error
}
