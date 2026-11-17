package out

import "github.com/LhamaLog/lhamalog-finder/internal/domain"

// Assinador deve delegar ao algoritmo canônico em application/assinatura (spec 008).
type Assinador interface {
	Calcular(domain.LogClassificado) (string, error)
}
