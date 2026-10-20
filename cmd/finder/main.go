// Comando finder: ponto de entrada do módulo Finder do LhamaLog.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	httpadp "github.com/LhamaLog/lhamalog-finder/internal/adapter/in/http"
	"github.com/LhamaLog/lhamalog-finder/internal/config"
)

// versao é definida no build com -ldflags "-X main.versao=X"; padrão "dev".
var versao = "dev"

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Carregar()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: httpadp.NovoRoteador(), ReadHeaderTimeout: 5 * time.Second}
	ln, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		log.Error("falha ao abrir a porta HTTP", "erro", err)
		os.Exit(1)
	}
	porta := ln.Addr().(*net.TCPAddr).Port
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("falha no servidor HTTP", "erro", err)
			stop()
		}
	}()
	// Log padrão do LhamaLog (ADR_149).
	log.Info("[LhamaLog] lhamalog-finder iniciado", "versao", versao, "porta", porta, "go", runtime.Version(), "workers", cfg.Workers)

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	log.Info("finder encerrado")
}
