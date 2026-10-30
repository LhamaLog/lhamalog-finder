package out

import (
	"context"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// Fila é a porta da fila, sem comando de banco e sem nome de stream, grupo ou campo.
type Fila interface {
	TomarLivre(ctx context.Context) (domain.Item, bool, error)
	Concluir(ctx context.Context, id domain.ID) error
	Liberar(ctx context.Context, id domain.ID) error
	EncerrarSemConclusao(ctx context.Context, id domain.ID) error
}
