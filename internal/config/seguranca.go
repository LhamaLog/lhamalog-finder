package config

import (
	"context"
	"os"

	"github.com/LhamaLog/lhamalog-finder/internal/adapter/out/catalogo"
	"github.com/LhamaLog/lhamalog-finder/internal/application"
)

// InicializarSeguranca deve preceder a fila e a prontidão do executável.
// A versão esperada é obrigatória; somente o caminho tem padrão publicado.
func InicializarSeguranca(ctx context.Context) (application.Seguranca, error) {
	caminho := os.Getenv("FINDER_CATALOGO_SEGURANCA")
	if caminho == "" {
		caminho = "config/seguranca/catalogo.json"
	}
	return application.InicializarSeguranca(ctx, catalogo.Arquivo{Caminho: caminho}, os.Getenv("FINDER_VERSAO_REGRAS"))
}
