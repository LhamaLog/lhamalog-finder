# Quickstart: Corrigir entrega da fila

Prova do comportamento descrito em [contracts/entrega-fila.md](./contracts/entrega-fila.md) e [data-model.md](./data-model.md). Não substitui o código de produção.

## Pré-requisitos

- Docker Compose na raiz do repositório.
- Redis do Compose publicado em `localhost:6379` (`REDIS_PORT` altera a porta).
- Nomes externos já lidos pelo processo. Este quickstart não cria nome novo. Exemplo do que o Compose de desenvolvimento já usa:

```text
FINDER_REDIS_ADDR=127.0.0.1:6379
FINDER_REDIS_STREAM=fila:logs
FINDER_REDIS_GROUP=finder
FINDER_REDIS_CODIGOS=finder:codigos
FINDER_HTTP_ADDR=:8085
```

O código só recebe stream, grupo e conjunto pelo ambiente. Sem eles, o processo recusa a fila.

## 1. Prova sem Redis

Na raiz:

```text
make test
```

Resultado esperado: `go test -race ./...` termina sem falha. Testes que exigem Redis não rodam neste alvo.

## 2. Prova com Redis

Subir o Redis:

```text
docker compose up -d redis
```

Esperar o healthcheck do serviço `redis`. Nesta instalação, `docker compose run` não aceita `--network`. A prova roda na rede do contêiner Redis, com a variável que torna Redis ausente uma falha:

```text
docker run --rm --network container:lhamalog-finder-redis -v "%CD%:/src" -w /src -e FINDER_REDIS_INTEGRATION=1 -e FINDER_REDIS_ADDR=127.0.0.1:6379 -e FINDER_REDIS_STREAM=fila:logs -e FINDER_REDIS_GROUP=finder -e FINDER_REDIS_CODIGOS=finder:codigos golang:1.27.0 go test -race -count=1 ./internal/adapter/in/queue/ ./internal/application/ ./cmd/finder/
```

Resultado esperado: a prova não é ignorada. Os cenários abaixo passam.

| Cenário | Resultado |
|---------|-----------|
| Um item aceito | Uma conclusão, código no conjunto uma vez, item fora da pendência |
| Queda depois do registro e antes da saída da pendência | A tentativa seguinte deixa uma conclusão e o item fora da pendência |
| Falha ao completar | O processo da prova segue, o item permanece parcial e a tentativa seguinte completa sem segunda conclusão |
| Instância some antes de qualquer efeito, com item atrás já livre | O item da frente é reservado antes do de trás, sem espera configurada |
| Duas instâncias | Itens diferentes podem estar reservados ao mesmo tempo; o mesmo código conta uma vez |
| Duas organizações, mesmo código | Uma conclusão |
| Item sem código e item com conteúdo extra, depois um aceito | Os dois primeiros descartados e observáveis, o aceito concluído, mapa de reserva sem os inválidos |
| Uma instância e mensagens já em ordem | A sequência de reservas é a ordem de entrada |
| Processo no ar | `GET /health` em `:8085` responde 200 e `{"status":"UP"}` |

## 3. Processo e healthcheck

```text
docker compose up --build --no-deps finder
```

O Compose já consulta `http://127.0.0.1:8085/health` no serviço `finder`, que compartilha a rede do Redis. Resultado esperado: o healthcheck deixa de falhar por conexão recusada enquanto o processo está no ar. Com o processo parado, a mesma consulta não é atendida.

## Fora desta prova

Leitura do log no OpenObserve, ordenação por nível dentro do Finder, fila de erro e prazo numérico.

## Resultado observado

Em 2026-11-05, com o Docker Desktop em execução e o serviço `redis` do Compose no ar:

- `docker compose run --rm --no-deps gotools test -race -count=1 ./...` terminou com código 0. Os pacotes `cmd/finder`, `internal/adapter/in/http`, `internal/adapter/in/queue`, `internal/application`, `internal/config` e `internal/domain` passaram. Sem `FINDER_REDIS_INTEGRATION`, a prova Redis não rodou.
- `docker run --rm --network container:lhamalog-finder-redis` com a imagem `golang:1.27.0`, `FINDER_REDIS_INTEGRATION=1`, `FINDER_REDIS_ADDR=127.0.0.1:6379` e os nomes `fila:logs`, `finder` e `finder:codigos`, executou `go test -race -count=1 ./internal/adapter/in/queue/ ./internal/application/ ./cmd/finder/`. Os três pacotes passaram. A prova não foi ignorada.
- A imagem `lhamalog-finder` foi construída e o processo subiu na rede do Redis. `GET http://127.0.0.1:8085/health` respondeu 200 e `{"status":"UP"}`. O log de arranque foi `{"level":"INFO","msg":"[LhamaLog] lhamalog-finder iniciado","versao":"dev"}`. Depois de remover o processo, a mesma consulta não foi atendida.

O comando `docker compose run` desta instalação não aceita `--network`. A prova com Redis usou `docker run --network container:lhamalog-finder-redis` no lugar da linha equivalente do quickstart.
