package out

import (
	"context"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// Entradas lê o stream, confirma pelo id Redis e adia a entrada que permanece pendente.
// EncerrarInvalida confirma uma vez e, se a chave está em processamento, remove a reserva.
// Essa confirmação não grava o estado concluído.
type Entradas interface {
	Ler(ctx context.Context) (domain.Entrada, bool, error)
	Confirmar(ctx context.Context, id string) error
	Adiar(ctx context.Context, id string) error
	EncerrarInvalida(ctx context.Context, id, idempotencyID string) error
}
