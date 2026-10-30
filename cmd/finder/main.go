// Comando finder: ponto de entrada do módulo Finder do LhamaLog.
package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"

	"github.com/LhamaLog/lhamalog-finder/internal/adapter/in/queue"
	"github.com/LhamaLog/lhamalog-finder/internal/application"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/config"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// versao é definida no build com -ldflags "-X main.versao=X"; padrão "dev".
var versao = "dev"

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log.Info("[LhamaLog] lhamalog-finder iniciado", "versao", versao)

	cfg := config.Carregar()
	fila, err := queue.Nova(queue.Opcoes{
		Addr:    cfg.RedisAddr,
		Stream:  cfg.RedisStream,
		Group:   cfg.RedisGroup,
		Codigos: cfg.RedisCodigos,
	})
	if err != nil {
		if errors.Is(err, queue.ErrNomesAusentes) {
			log.Error("consumo da fila sem nomes externos", "erro", err)
		} else {
			log.Error("falha ao abrir a fila", "erro", err)
		}
		os.Exit(1)
	}
	defer fila.Fechar()

	svc := &application.Servico{
		Fila:     fila,
		Codigos:  fila,
		Disp:     fila,
		Reservas: domain.NovoQuadro(nil),
	}
	ctx, parar := signal.NotifyContext(context.Background(), os.Interrupt)
	defer parar()

	for {
		saida, err := svc.Executar(ctx, in.Pedido{
			Instancia:     fila.Instancia(),
			WorkerNoAr:    true,
			InstanciaNoAr: true,
		})
		if ctx.Err() != nil {
			break
		}
		if err != nil {
			log.Error("falha de leitura da fila", "erro", err, "conclusoes", saida.Conclusoes)
			os.Exit(1)
		}
		if saida.Conclusoes > 0 {
			log.Info("item concluído", "codigos", saida.CodigosDisponiveis)
		}
	}

	log.Info("finder encerrado")
}
