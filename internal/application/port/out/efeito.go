package out

import (
	"context"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// Efeito recebe a referência validada e o id da entrada, que acompanha a referência.
// concluido autoriza a confirmação. Falha, ou a porta ausente, mantém a entrada pendente.
// A consulta ao OpenObserve não faz parte desta porta.
type Efeito interface {
	Executar(ctx context.Context, ref domain.Referencia) (concluido bool, err error)
}
