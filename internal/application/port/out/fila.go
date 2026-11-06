package out

import (
	"context"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// Fila é a porta da fila, sem comando de banco e sem nome de stream, grupo ou campo.
type Fila interface {
	TomarLivre(ctx context.Context) (domain.Item, bool, error)
	Concluir(ctx context.Context, id domain.ID) error
	// ConcluirPrimeira registra o código uma vez e tira o item da pendência na mesma operação.
	// O booleano é verdadeiro só quando o código entrou agora.
	ConcluirPrimeira(ctx context.Context, id domain.ID, codigo domain.Codigo) (bool, error)
	Liberar(ctx context.Context, id domain.ID) error
	EncerrarSemConclusao(ctx context.Context, id domain.ID) error
}
