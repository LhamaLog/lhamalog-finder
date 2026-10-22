// Comando finder: ponto de entrada do módulo Finder do LhamaLog.
package main

import (
	"log/slog"
	"os"

	"github.com/ORG_PENDENTE/lhamalog-finder/internal/adapter/in/queue"
)

// versao é definida no build com -ldflags "-X main.versao=X"; padrão "dev".
var versao = "dev"

func main() {

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Log padrão do LhamaLog (ADR_149).
	log.Info("[LhamaLog] lhamalog-finder iniciado", "versao", versao)

	rdb := queue.InitRedis()
	defer rdb.Close()

	log.Info("finder encerrado")
}
