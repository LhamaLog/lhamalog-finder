package out

import (
	"context"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// Codigos guarda uma conclusão por código. O lugar físico não tem nome.
type Codigos interface {
	RegistrarPrimeira(ctx context.Context, codigo domain.Codigo) (bool, error)
	JaConcluido(ctx context.Context, codigo domain.Codigo) (bool, error)
}
