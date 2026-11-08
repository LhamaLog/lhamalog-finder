package out

import (
	"context"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// EstadoLeitura distingue registro obtido, ausência e armazenamento inacessível.
type EstadoLeitura int

const (
	LeituraPronta EstadoLeitura = iota
	LeituraAusente
	LeituraInacessivel
)

// LeituraLog lê um único registro. O conteúdo segue inteiro.
// A porta não oculta campo e não classifica ocorrência.
type LeituraLog interface {
	Ler(ctx context.Context, codigo domain.Codigo) (domain.Registro, EstadoLeitura, error)
}
