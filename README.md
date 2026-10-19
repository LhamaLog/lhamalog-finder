# lhamalog-finder

Módulo Finder do LhamaLog, em Go: consome as execuções da fila Redis Streams, encontra problemas nos logs em paralelo (worker pool) e mascara dados sensíveis. Arquitetura hexagonal leve.

> **[PENDENTE] org do GitHub:** o módulo está como `github.com/ORG_PENDENTE/lhamalog-finder`.
> Ao definir a org: `go mod edit -module github.com/<org>/lhamalog-finder` e troque os imports
> (`grep -rl ORG_PENDENTE . | xargs sed -i 's#ORG_PENDENTE#<org>#g'`).

## Pré-requisitos
- Go 1.27.x
- Docker + Docker Compose
- make

## Instalando o make

Os comandos `make` abaixo precisam do GNU Make instalado.

**Windows (PowerShell):**

```powershell
winget install ezwinports.make
# alternativa: choco install make   (ou use Git Bash / WSL)
```

Depois feche e abra o terminal e confira com `make --version`.

**Linux:**

```bash
# Debian/Ubuntu
sudo apt update && sudo apt install -y make
# Fedora
sudo dnf install -y make
# Arch
sudo pacman -S make
make --version
```

macOS (opcional): `xcode-select --install`.

**Sem make no Windows:** use os equivalentes diretos:

```powershell
go run ./cmd/finder
go build ./...
go test ./...
```

## Como executar
1. Extraia este zip na raiz do repositório.
2. Ajuste o módulo/org (acima).
3. `make infra-up` sobe MongoDB, Redis e Elasticsearch locais.
4. `make run` (equivale a `go run ./cmd/finder`). Health em http://localhost:8085/health
5. `make test` (equivale a `go test -race ./...`).

Equivalentes sem make: `go run ./cmd/finder`, `go build -o bin/finder ./cmd/finder`, `go test ./...`, `go vet ./...`.

Outros alvos: `make build` (bin/finder), `make lint` (go vet + golangci-lint opcional),
`make docker` (imagem multi-stage), `make infra-down`.

Variáveis: `FINDER_HTTP_ADDR` (padrão `:8085`), `FINDER_WORKERS` (padrão 4).
Sem GitHub Actions: rode `make test` localmente antes de cada push.

## Log de inicialização

Depois de abrir a porta HTTP, o Finder registra a linha padrão do LhamaLog (ADR_149) via `slog` (JSON):

```
{"level":"INFO","msg":"[LhamaLog] lhamalog-finder iniciado","versao":"dev","porta":8085,"go":"go1.27.2","workers":4}
```

A versão é `dev` por padrão; para gravar outra no build: `go build -ldflags "-X main.versao=0.0.1" -o bin/finder ./cmd/finder`.

## Layout
```
cmd/finder/                   main (wiring, shutdown gracioso)
internal/domain/              entidades e regras
internal/application/         casos de uso
internal/application/port/in  portas de entrada
internal/application/port/out portas de saída
internal/adapter/in/queue/    consumidor Redis Streams
internal/adapter/in/http/     health interno
internal/adapter/out/{mongo,elastic,llm}/
internal/config/              configuração por variáveis de ambiente
```

## Dependências candidatas (ainda não adicionadas)
Antes de adicionar qualquer dependência, ver Drive da equipe > Arquitetura > Dependências específicas dos módulos.
- github.com/redis/go-redis/v9 — Redis Streams
- go.mongodb.org/mongo-driver/v2 — MongoDB
- github.com/elastic/go-elasticsearch/v9 — Elasticsearch
- go.opentelemetry.io/otel (+ sdk, exporters OTLP) — observabilidade
- golang.org/x/sync (errgroup, semaphore) — concorrência
- github.com/testcontainers/testcontainers-go — testes de integração
- LLM / BancoVetorial — [PENDENTE]


## Padrão de commit

Mensagens no formato `CHAVE-JIRA: descricao` (ex.: `LLT-82: inicializacao do modulo: Collector`), conforme a ADR_148. Detalhes em [CONTRIBUTING.md](CONTRIBUTING.md).
