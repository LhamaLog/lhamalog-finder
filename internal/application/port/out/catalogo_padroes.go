package out

import (
	"context"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// CatalogoPadroes lê o artefato imutável e executável da versão distribuída.
type CatalogoPadroes interface {
	Carregar(context.Context) (domain.CatalogoSeguranca, error)
}
