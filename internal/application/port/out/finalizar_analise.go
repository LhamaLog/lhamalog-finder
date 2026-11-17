package out

import (
	"context"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// FinalizarAnalise é o marco lógico finishScan da sequência; só após confirmação.
// O bootstrap real aguarda o contrato externo e o ciclo completo de varredura.
type FinalizarAnalise interface {
	ConcluirAnalise(context.Context, domain.ChaveResumo) error
}
