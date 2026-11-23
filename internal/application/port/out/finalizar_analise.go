package out

import (
	"context"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// FinalizarAnalise é o marco lógico finishScan da sequência; só após confirmação.
// DF-FINDER-2026-11-22: v1 conclui a chave Redis; contexto V/M real usa HTTPS.
type FinalizarAnalise interface {
	ConcluirAnalise(context.Context, domain.ChaveResumo) error
}
