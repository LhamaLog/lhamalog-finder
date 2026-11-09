package out

import "context"

// ChaveLote guarda uma chave por idempotencyId no prefixo externo do Finder.
// Reservar cria em_processamento só se a chave não existe, sem prazo.
// Concluir grava concluido e 25 horas uma vez.
// RemoverReserva apaga a reserva e não grava concluido.
// A mesma chave não serve à Varredura.
type ChaveLote interface {
	Reservar(ctx context.Context, idempotencyID string) (bool, error)
	Concluir(ctx context.Context, idempotencyID string) error
	Consultar(ctx context.Context, idempotencyID string) (string, error)
	RemoverReserva(ctx context.Context, idempotencyID string) error
}
