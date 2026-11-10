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

	"github.com/LhamaLog/lhamalog-finder/internal/adapter/in/queue"
	"github.com/LhamaLog/lhamalog-finder/internal/application"
	"github.com/LhamaLog/lhamalog-finder/internal/config"
)

// versao é definida no build com -ldflags "-X main.versao=X"; padrão "dev".
var versao = "dev"

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(log)
	log.Info("[LhamaLog] lhamalog-finder iniciado", "versao", versao)

	if err := config.AplicarArquivoEnv(); err != nil {
		log.Error("falha ao ler .env", "erro", err)
		os.Exit(1)
	}
	cfg := config.Carregar()
	if err := cfg.ValidarConsumo(); err != nil {
		log.Error("consumo recusado", "erro", err)
		os.Exit(1)
	}
	lote, err := queue.NovaLote(queue.OpcoesLote{
		Addr:         cfg.RedisAddr,
		Senha:        cfg.RedisSenha,
		Stream:       cfg.RedisStream,
		Group:        cfg.RedisGroup,
		Consumidor:   cfg.RedisConsumidor,
		Idempotencia: cfg.RedisIdempotencia,
		Count:        int64(cfg.RedisCount),
		Block:        cfg.RedisBlock,
		Ociosidade:   cfg.RedisPendenciaOciosa,
	})
	if err != nil {
		log.Error("falha ao abrir a fila", "erro", err)
		os.Exit(1)
	}
	defer lote.Fechar()

	ouvinte, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		log.Error("falha ao abrir a consulta de saúde", "erro", err, "addr", cfg.HTTPAddr)
		os.Exit(1)
	}
	srv := servidor(cfg.HTTPAddr)
	go func() {
		if err := srv.Serve(ouvinte); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("consulta de saúde encerrada", "erro", err)
		}
	}()

	// A porta de efeito fica ausente: esta feature não consulta o OpenObserve.
	// FINDER_WORKERS não define o paralelismo desta leitura.
	consumo := &application.ConsumoReferencia{Fila: lote, Chaves: lote}
	ctx, parar := signal.NotifyContext(context.Background(), os.Interrupt)
	defer parar()

	for {
		_, err := consumo.Executar(ctx)
		if ctx.Err() != nil {
			break
		}
		if err != nil {
			log.Error("falha no consumo do lote", "erro", err)
			continue
		}
	}

	_ = srv.Shutdown(context.Background())
	log.Info("finder encerrado")
}
