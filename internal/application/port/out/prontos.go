package out

import (
	"context"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// LogsProntos guarda e lê o registro inteiro. O lugar não tem nome.
// Uma segunda guarda do mesmo apontamento não cria outro log.
// Com a instância parada, a leitura não é prometida.
type LogsProntos interface {
	Guardar(ctx context.Context, codigo domain.Codigo, conteudo domain.Registro) error
	Ler(ctx context.Context, codigo domain.Codigo, instanciaNoAr bool) (domain.Registro, bool, error)
}
