package out

import (
	"context"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// CatalogoPadroes lê o artefato versionado. Adaptador de arquivo bloqueado P-010-2.
type CatalogoPadroes interface {
	Carregar(context.Context) (domain.CatalogoSeguranca, error)
}
