# lhamalog-finder

Módulo Finder do LhamaLog, em Go: consome as execuções da fila Redis Streams, encontra problemas nos logs em paralelo (worker pool) e mascara dados sensíveis. Arquitetura hexagonal leve.

## Pré-requisitos

- Docker Engine com o plugin Compose (`docker compose`)
- make (opcional)

O Go **1.27.0** já está na imagem. Não é preciso instalar a linguagem na máquina para compilar, testar ou executar o Finder.

## Instalando o Docker

Confira no final com `docker --version` e `docker compose version`.

**Windows:**

1. O Docker Desktop usa o WSL 2. Se o instalador pedir, aceite a configuração do WSL.
2. Instale o [Docker Desktop para Windows](https://docs.docker.com/desktop/setup/install/windows-install/):

```powershell
winget install -e --id Docker.DockerDesktop
```

1. Abra o Docker Desktop e espere o status ficar em execução.
2. Feche e abra o terminal antes de usar `docker`.

**Linux (Ubuntu):**

```bash
sudo apt update
sudo apt install -y ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo usermod -aG docker "$USER"
```

Saia e entre de novo na sessão para o grupo `docker` valer. No Debian, use o [repositório Debian](https://docs.docker.com/engine/install/debian/) (o trecho acima é o do Ubuntu). Nas outras distros, siga o guia do [Docker Engine](https://docs.docker.com/engine/install/).

**macOS:** instale o [Docker Desktop para Mac](https://docs.docker.com/desktop/setup/install/mac-install/) e abra o aplicativo. Com Homebrew: `brew install --cask docker`.

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

**Sem make:** os mesmos passos estão na seção seguinte, com `docker compose`.

## Como iniciar

Na raiz do repositório, com o Docker em execução:

```bash
make run
```

O alvo sobe MongoDB, Redis e Elasticsearch, espera o Redis aceitar conexões e então inicia o Finder em primeiro plano, compilado com Go 1.27.0. A primeira execução baixa as imagens e pode demorar.

O log fica nesse terminal. Interrompa com Ctrl+C. Para remover os contêineres: `make infra-down`.

Sem make:

```bash
docker compose up -d mongo redis elasticsearch
docker compose up -d --wait redis
docker compose up --build --no-deps finder
```

Só a infraestrutura, sem o Finder: `make infra-up` (ou `docker compose up -d mongo redis elasticsearch`).

Outros alvos, todos dentro do Docker:


| Comando                                                | O que faz                                    |
| ------------------------------------------------------ | -------------------------------------------- |
| `make build` ou `make docker`                          | gera a imagem `lhamalog-finder`              |
| `make test`                                            | `go test -race ./...` no container Go 1.27.0 |
| `make lint`                                            | `go vet ./...` no mesmo container            |
| `docker compose run --rm --no-deps gotools test ./...` | testes sem make                              |


Para subir só o Finder, sem MongoDB, Redis e Elasticsearch:

```bash
docker build -t lhamalog-finder .
docker run --rm -p 8085:8085 lhamalog-finder
```

Variáveis: `FINDER_PORT` (porta no host, padrão `8085`), `FINDER_WORKERS` (padrão `4`). Dentro do contêiner o processo escuta em `:8085`.
Sem GitHub Actions: rode `make test` localmente antes de cada push.

## Log de inicialização

Depois de abrir a porta HTTP, o Finder registra a linha padrão do LhamaLog (ADR_149) via `slog` (JSON):

```
{"level":"INFO","msg":"[LhamaLog] lhamalog-finder iniciado","versao":"dev","porta":8085,"go":"go1.27.0","workers":4}
```

A versão é `dev` por padrão. A linha `go` reflete o Go 1.27.0 da imagem.

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