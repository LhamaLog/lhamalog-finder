# lhamalog-finder

Módulo Finder do LhamaLog, em Go: consome referências de lote na fila Redis Streams, recupera os logs no OpenObserve, classifica, mascara dados sensíveis e publica o efeito do pipeline. Arquitetura hexagonal leve.

O número em `FINDER_WORKERS` permanece na configuração, com padrão `4`. A leitura do stream segue o consumo ligado em `cmd/finder`.

## Pré-requisitos

- **Go 1.27.4** para compilar, executar e testar na máquina. É a versão de `go.mod`, da imagem `golang:1.27.4` do `Dockerfile` e do serviço `gotools` no Compose.
- **Docker Engine** com o plugin Compose (`docker compose`), para a infraestrutura local e para o caminho que compila dentro do contêiner.
- **GNU Make** (opcional). Os mesmos passos têm equivalente com `docker compose` e com `go`.
- Dependências locais definidas em `docker-compose.yml`:
  - Redis `redis:8.2.10`, porta do host `${REDIS_PORT:-6379}`
  - OpenObserve `openobserve/openobserve:v1.1.2`, porta do host `${OPENOBSERVE_PORT:-5080}`
  - MongoDB `mongo:8.0.34`, porta do host `${MONGO_PORT:-27017}`

O caminho Docker compila o Finder na imagem. A linguagem na máquina só é necessária para `go build`, `go test` e `go vet` fora do contêiner.

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

**Sem make:** os mesmos passos estão nas seções seguintes, com `docker compose` e com `go`.

## Configuração

Na raiz do repositório:

```bash
cp .env.example .env
```

`.env` está no `.gitignore`. Grave senhas e JWTs só nesse arquivo, no ambiente do processo ou no gerenciador de segredos. Este README e o `.env.example` deixam esses campos vazios.

`cmd/finder` chama `config.AplicarArquivoEnv()` antes de ler a configuração. O arquivo é o `.env` do diretório atual. Arquivo ausente segue a execução. Variável já exportada no terminal permanece; o arquivo não a substitui. Se o mesmo nome aparece mais de uma vez no arquivo, vale a última ocorrência. O exemplo repete `FINDER_VERSAO_REGRAS`, `FINDER_MODULO_OCORRENCIA_URL` e `FINDER_MODULO_OCORRENCIA_JWT`.

O processo encerra com `consumo recusado`, `seguranca recusada` ou `pipeline recusado` enquanto a configuração obrigatória estiver incompleta. Os padrões abaixo são os de `internal/config` e de `internal/bootstrap`. Onde a célula diz "sem padrão", o `.env.example` deixa o valor vazio de propósito.

| Variável | Padrão no código | Papel |
| --- | --- | --- |
| `FINDER_HTTP_ADDR` | `:8085` | Endereço da consulta `GET /health` |
| `FINDER_WORKERS` | `4` | Inteiro positivo. Valor ausente ou inválido mantém `4` |
| `FINDER_REDIS_ADDR` | `localhost:6379` | Endereço do Redis. O exemplo traz `127.0.0.1:6379` |
| `FINDER_REDIS_STREAM` | sem padrão | Tem de ser `lhamalog:stream:logs` |
| `FINDER_REDIS_GROUP` | sem padrão | Nome do consumer group. O exemplo traz `finder` |
| `FINDER_REDIS_CONSUMIDOR` | sem padrão | Nome do consumidor. Vazio recusa o consumo |
| `FINDER_REDIS_IDEMPOTENCIA` | sem padrão | Prefixo da chave de idempotência. Vazio recusa o consumo |
| `FINDER_REDIS_COUNT` | sem padrão | Quantidade positiva de entradas por leitura |
| `FINDER_REDIS_BLOCK` | sem padrão | Bloqueio positivo, em duração Go (`200ms`, `1s`) |
| `FINDER_REDIS_PENDENCIA_OCIOSA` | sem padrão | Ociosidade positiva da pendência, em duração Go |
| `FINDER_REDIS_PASSWORD` | vazio | Senha do Redis. Ausente no `.env.example`. Vazia permanece vazia. O Redis do Compose sobe sem senha |
| `FINDER_REDIS_CODIGOS` | sem padrão | Lido para `Config`. A validação do consumo não exige o valor. O exemplo traz `finder:codigos`. A prova antiga de fila falha se estiver vazio com `FINDER_REDIS_INTEGRATION=1` |
| `FINDER_OO_URL` | sem padrão | URL base do OpenObserve, com esquema |
| `FINDER_OO_USUARIO` | sem padrão | Usuário da autenticação Basic |
| `FINDER_OO_SENHA` | sem padrão | Senha da autenticação Basic |
| `FINDER_OO_TIMEOUT` | sem padrão | Prazo positivo por consulta, em duração Go |
| `FINDER_OO_PAGINA` | sem padrão | Quantidade positiva de documentos físicos por página |
| `FINDER_OO_CAMPOS_COPIA` | ausente | A variável precisa existir para a recuperação consultar. Presente e vazia compara o conteúdo inteiro. Vários campos separam-se por vírgula |
| `FINDER_OO_ORDEM` | sem padrão | Expressão SQL de ordenação estável entre páginas |
| `FINDER_OO_CAMPO_TENANT` | sem padrão | Nome do campo de tenant no documento. Vazio recusa o pipeline |
| `FINDER_OO_CAMPO_PRODUTO` | sem padrão | Nome do campo de produto no documento. Vazio recusa o pipeline |
| `FINDER_CATALOGO_SEGURANCA` | `config/seguranca/catalogo.json` | Caminho do catálogo de mascaramento |
| `FINDER_VERSAO_REGRAS` | sem padrão | Tem de ser igual ao campo `versao` do catálogo. No arquivo versionado, esse valor é `finder-2026-11-22.1` |
| `FINDER_MODULO_OCORRENCIA_URL` | sem padrão | URL `https` do Módulo de Ocorrência, sem usuário, consulta ou fragmento na URL |
| `FINDER_MODULO_OCORRENCIA_JWT` | sem padrão | JWT desse módulo, sem espaço em branco |
| `FINDER_CORE_URL` | sem padrão | URL `https` do Core, nas mesmas restrições de URL |
| `FINDER_CORE_JWT` | sem padrão | JWT do Core |
| `FINDER_CORE_TIMEOUT` | `2s` | Duração positiva. O exemplo já traz `2s` |
| `FINDER_VARREDURA_URL` | sem padrão | URL `https` da conclusão de varredura. Os dois lados vazios deixam essa etapa de fora |
| `FINDER_VARREDURA_JWT` | sem padrão | JWT da varredura. Um lado preenchido e o outro vazio recusa o pipeline |
| `FINDER_VARREDURA_TIMEOUT` | `2s` | Duração positiva. O exemplo já traz `2s` |

Forma local para os nomes que o exemplo deixa vazios e o processo exige. Troque os marcadores pelos valores do seu ambiente. A senha do OpenObserve e os JWTs ficam só no `.env`.

```text
FINDER_REDIS_CONSUMIDOR=finder-local
FINDER_REDIS_IDEMPOTENCIA=finder:log-ingested:
FINDER_REDIS_COUNT=2
FINDER_REDIS_BLOCK=200ms
FINDER_REDIS_PENDENCIA_OCIOSA=1s
FINDER_OO_URL=http://127.0.0.1:5080
FINDER_OO_USUARIO=<usuario-da-instancia-local>
FINDER_OO_SENHA=<senha-da-instancia-local>
FINDER_OO_TIMEOUT=2s
FINDER_OO_PAGINA=2
FINDER_OO_CAMPOS_COPIA=
FINDER_OO_ORDEM=lhamalog_record_id
FINDER_OO_CAMPO_TENANT=tenantid
FINDER_OO_CAMPO_PRODUTO=productid
FINDER_VERSAO_REGRAS=finder-2026-11-22.1
FINDER_MODULO_OCORRENCIA_URL=https://modulo.exemplo
FINDER_MODULO_OCORRENCIA_JWT=<jwt-do-ambiente>
FINDER_CORE_URL=https://core.exemplo
FINDER_CORE_JWT=<jwt-do-ambiente>
```

`FINDER_OO_ORDEM`, `FINDER_OO_CAMPO_TENANT` e `FINDER_OO_CAMPO_PRODUTO` nesse bloco são os nomes usados pela prova de integração do binário. O código de produção não os assume.

`FINDER_PORT`, `REDIS_PORT`, `OPENOBSERVE_PORT`, `MONGO_PORT`, `MONGO_ROOT_USERNAME`, `MONGO_ROOT_PASSWORD`, `MONGO_APP_DATABASE`, `MONGO_APP_USERNAME` e `MONGO_APP_PASSWORD` são interpolados pelo Compose. O processo Go não os lê. `FINDER_PORT` (padrão `8085`) publica a consulta de saúde no serviço `redis`, porque o contêiner `finder` usa `network_mode: service:redis`.

O serviço `finder` recebe as variáveis listadas em `environment` no `docker-compose.yml`. `FINDER_REDIS_ADDR` e `FINDER_REDIS_PASSWORD` ficam de fora dessa lista: dentro do Compose o processo usa `localhost:6379` e senha vazia, na rede do Redis. `FINDER_OO_CAMPOS_COPIA` está declarado nesse arquivo sem valor. A recuperação só consulta quando essa variável existe no processo; no binário nativo, a linha vazia do `.env` a deixa presente.

Credenciais locais do Mongo escritas no Compose, usadas só na inicialização do banco: usuário raiz `root`, senha `root`, banco `lhamalog_collector`, usuário da aplicação `lhamalog`, senha `lhamalog`. O pacote `internal/adapter/out/mongo` ainda não abre conexão.

Organização, usuário e senha do OpenObserve não são definidos por este módulo. A interface local escuta em `http://127.0.0.1:5080`.

## Subir as dependências

Na raiz, com o Docker em execução:

```bash
make infra-up
```

Sem make:

```bash
docker compose up -d mongo redis openobserve
```

O alvo sobe MongoDB 8.0.34, Redis 8.2.10 e OpenObserve v1.1.2. A primeira execução baixa as imagens. Para esperar o healthcheck do Redis:

```bash
docker compose up -d --wait redis
```

Para remover os contêineres e a rede: `make infra-down` (ou `docker compose down`). Os volumes nomeados permanecem.

## Compilar e executar

O binário de serviço é `./cmd/finder`. Execute a partir da raiz, para o `.env` e o caminho relativo `config/seguranca/catalogo.json` resolverem.

Com Go 1.27.4 na máquina, o mesmo comando do `Dockerfile`:

```bash
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.versao=dev" -o bin/finder ./cmd/finder
./bin/finder
```

`bin/` está no `.gitignore`. Sem `-X main.versao=...`, a versão registrada é `dev`.

Sem gerar arquivo:

```bash
go run ./cmd/finder
```

Pelo Compose, na raiz, com o `.env` preenchido para as chaves que o serviço interpola:

```bash
make run
```

Sem make:

```bash
docker compose up -d mongo redis openobserve
docker compose up -d --wait redis
docker compose up --build --no-deps finder
```

O log fica nesse terminal. Interrompa com Ctrl+C (`os.Interrupt`). O processo registra `finder encerrado` depois de fechar a consulta de saúde.

Imagem sem o restante da stack:

```bash
docker build -t lhamalog-finder .
docker run --rm -p 8085:8085 --env-file .env lhamalog-finder
```

A imagem exclui `.env` no build. `--env-file` entrega as variáveis ao processo. O Redis e o OpenObserve precisam estar alcançáveis nos endereços desse arquivo.

Com o processo no ar, `GET /health` responde `{"status":"UP"}`. No Compose, a porta do host é `FINDER_PORT` (padrão `8085`).

Outros comandos, também a partir da raiz:

```bash
go run ./cmd/avaliar-base -diretorio datasets/labeled-logs
go run ./cmd/capturar-base -destino /caminho/novo -por-classe 1020
```

`avaliar-base` carrega a base versionada, exige o gate de precisão/recall/F1 e escreve o relatório em JSON na saída padrão. `capturar-base` exige `FINDER_REDIS_ADDR` e usa `FINDER_REDIS_PASSWORD` quando houver senha. `-destino` é um diretório novo. `-por-classe` padrão é `1020`.

| Comando | O que faz |
| --- | --- |
| `make build` ou `make docker` | gera a imagem `lhamalog-finder` |
| `make run` | sobe Mongo, Redis e OpenObserve e inicia o Finder em primeiro plano |
| `make infra-up` / `make infra-down` | sobe ou derruba a infraestrutura |
| `make test` | `go test -race ./...` no contêiner Go 1.27.4 |
| `make lint` | `go vet ./...` no mesmo contêiner |
| `make clean` | remove o diretório `bin/` |

Sem GitHub Actions: rode a suíte de testes localmente antes de cada push. O passo está em [CONTRIBUTING.md](CONTRIBUTING.md).

## Testes

Na raiz, com Go 1.27.4. A suíte padrão pula as provas que dependem de Redis, OpenObserve ou da VM de aceite enquanto as variáveis abaixo estão ausentes.

```bash
go test ./...
go test -race ./...
go test -cover ./...
go test -coverprofile=cobertura.out ./...
go tool cover -func=cobertura.out
go vet ./...
```

`-race` e `-cover` podem ir no mesmo `go test`. `cobertura.out` é artefato local; `*.out` e `coverage*` estão no `.gitignore`. O `Makefile` cobre race (`make test`) e vet (`make lint`). Cobertura e `go test` sem `-race` rodam pelo comando `go`, ou pelo serviço `gotools`:

```bash
docker compose run --rm --no-deps gotools test ./...
docker compose run --rm --no-deps gotools test -race ./...
docker compose run --rm --no-deps gotools test -cover ./...
docker compose run --rm --no-deps gotools vet ./...
```

`make test` é o primeiro desses com `-race`. O diretório de trabalho no contêiner é `/src`, com o repositório montado.

`TestBaseRealRotulada` (`internal/labeleddataset/base_real_test.go`) entra em `go test ./...`. Ele carrega `datasets/labeled-logs` por `CarregarBasePublicada` e falha se o manifesto ou os JSONL publicados não validam. A variável `FINDER_DATASET_INTEGRATION`, citada em [specs/007-base-logs-rotulados](specs/007-base-logs-rotulados/quickstart.md), não é lida por esse teste.

### Provas opt-in

Com a variável habilitada e o serviço fora do ar, ou com parâmetro obrigatório ausente, a prova falha. Ela é ignorada só quando a variável de habilitação está ausente (ou, no caso do Redis e do aceite, com valor diferente de `1`).

| Variável | Efeito |
| --- | --- |
| `FINDER_REDIS_INTEGRATION=1` | Executa as provas de Redis em `internal/adapter/in/queue/redis_integration_test.go` (`nomesRedis` e `provaLote`) e `TestHandoffHTTPComRedisMantemPendenteAteRecibo`. Outros testes do mesmo pacote de fila seguem na suíte padrão. A prova de fila usa `FINDER_REDIS_ADDR` (se vazio, `127.0.0.1:6379` nessa prova), `FINDER_REDIS_STREAM`, `FINDER_REDIS_GROUP` e `FINDER_REDIS_CODIGOS`. A prova de lote lê stream, grupo, consumidor, prefixo, quantidade, bloqueio, ociosidade e senha por `config.Carregar` (`FINDER_REDIS_PASSWORD`). A prova de handoff lê a senha em `FINDER_REDIS_SENHA` e também `FINDER_REDIS_ADDR`, `FINDER_REDIS_STREAM`, `FINDER_REDIS_GROUP` e `FINDER_REDIS_IDEMPOTENCIA` |
| `FINDER_OPENOBSERVE_INTEGRATION=1` | Presença do nome, mesmo vazio, habilita `TestOpenObserveIntegracao` e `TestMapearDocumentoOpenObserveIntegracao`. Valor não vazio habilita também `TestHistoricoOpenObserveReal`. Exigem `FINDER_OO_URL` respondendo e `FINDER_OO_USUARIO`, `FINDER_OO_SENHA`, `FINDER_OO_TIMEOUT`, `FINDER_OO_PAGINA`, `FINDER_OO_ORDEM`, `FINDER_OO_CAMPO_TENANT`, `FINDER_OO_CAMPO_PRODUTO`, com `FINDER_OO_CAMPOS_COPIA` definido (pode ser vazio). O histórico usa URL, usuário e senha |
| `FINDER_REDIS_INTEGRATION=1` e `FINDER_OPENOBSERVE_INTEGRATION` não vazio | No Linux, `cmd/finder/pipeline_integracao_test.go` compila o binário e executa o pipeline. Usa `FINDER_REDIS_ADDR`, `FINDER_REDIS_PASSWORD`, `FINDER_OO_URL`, `FINDER_OO_USUARIO` e `FINDER_OO_SENHA`. Em outro sistema operacional a prova é ignorada |
| `FINDER_ACEITE_VMA=1` | Roda `TestAceiteDesempenhoVMA`. O próprio teste descreve a exigência de VM-A real, com os serviços e um lote de 1.000 registros. A flag não comprova esse perfil; o corpo mede um ciclo local de 1.000 logs e falha acima de 2s |

Exemplo de prova local, com Redis e OpenObserve do Compose já no ar e os parâmetros de teste exportados:

```bash
export FINDER_REDIS_INTEGRATION=1
export FINDER_OPENOBSERVE_INTEGRATION=1
go test -race -count=1 ./internal/adapter/in/queue/ ./internal/adapter/out/openobserve/ ./internal/application/handoff/ ./cmd/finder/
```

Os quickstarts de prova estão em [specs/004-capturar-referencia-lote/quickstart.md](specs/004-capturar-referencia-lote/quickstart.md), [specs/005-recuperar-logs-lote/quickstart.md](specs/005-recuperar-logs-lote/quickstart.md), [specs/006-modelo-classificacao-logs/quickstart.md](specs/006-modelo-classificacao-logs/quickstart.md) e [specs/007-base-logs-rotulados/quickstart.md](specs/007-base-logs-rotulados/quickstart.md).

## Recuperação do lote no OpenObserve

O Finder usa a organização, o stream e os identificadores da referência recebida para consultar o lote. A configuração abaixo é lida do ambiente ou do arquivo `.env`. Os valores são definidos pelo ambiente de execução.

| Variável | Uso |
| --- | --- |
| `FINDER_OO_URL` | URL base do OpenObserve, com esquema HTTP ou HTTPS |
| `FINDER_OO_USUARIO` | Usuário da autenticação Basic |
| `FINDER_OO_SENHA` | Senha da autenticação Basic |
| `FINDER_OO_TIMEOUT` | Prazo positivo por consulta, em formato de duração Go |
| `FINDER_OO_PAGINA` | Quantidade positiva de documentos físicos por página |
| `FINDER_OO_CAMPOS_COPIA` | Lista de campos de cópia física separados por vírgula. Ausente impede a consulta. Presente e vazia compara o conteúdo completo |
| `FINDER_OO_ORDEM` | Expressão SQL de ordenação estável para todas as páginas do lote |
| `FINDER_OO_CAMPO_TENANT` | Nome do campo confiável de tenant no topo do documento armazenado |
| `FINDER_OO_CAMPO_PRODUTO` | Nome do campo confiável de produto no topo do documento armazenado |

A recuperação filtra lote, tenant e produto, pagina até uma página curta e deduplica por `lhamalog_record_id`. O resultado só disponibiliza registros quando a contagem distinta fecha e nenhuma falha de maior precedência foi encontrada. As regras e a configuração da prova sintética estão em [specs/005-recuperar-logs-lote](specs/005-recuperar-logs-lote/quickstart.md).

No processo atual, o efeito ligado em `cmd/finder` é o pipeline de `internal/bootstrap`. A confirmação da entrada no Redis continua no consumo, depois do efeito. O contrato e os limites da recuperação estão em [contracts/recuperar-lote.md](specs/005-recuperar-logs-lote/contracts/recuperar-lote.md).

## Log de inicialização

Ao iniciar, o Finder registra esta linha via `slog` (JSON):

```text
{"level":"INFO","msg":"[LhamaLog] lhamalog-finder iniciado","versao":"dev"}
```

A versão é `dev` por padrão, ou o valor de `-X main.versao`. No interrupt, o processo registra `finder encerrado`.

O consumo da referência (`ConsumoReferencia`) registra, no mesmo JSON: `entrada inválida`, `entrada pendente`, `referência entregue`, `efeito concluído` e `confirmada sem repetir efeito`. O consumo antigo de item fora do contrato registra `item descartado`.

## Layout

```text
cmd/finder/                   serviço (wiring, saúde, shutdown no interrupt)
cmd/avaliar-base/             avalia a base rotulada versionada
cmd/capturar-base/            captura diagnóstica controlada no Redis
internal/domain/              entidades e regras
internal/application/         casos de uso
internal/application/port/in  portas de entrada
internal/application/port/out portas de saída
internal/adapter/in/queue/    consumidor Redis Streams
internal/adapter/in/http/     health interno
internal/adapter/out/openobserve/
internal/adapter/out/mongo/   reserva; sem cliente Mongo
internal/bootstrap/           pipeline de produção
internal/config/              configuração por variáveis de ambiente
config/seguranca/catalogo.json
datasets/labeled-logs/        base rotulada publicada
```

## Dependências

`go.mod` já exige `github.com/redis/go-redis/v9` e `github.com/joho/godotenv`.

Antes de adicionar qualquer outra dependência, ver Drive da equipe > Arquitetura > Dependências específicas dos módulos. Ainda fora do módulo:

- go.mongodb.org/mongo-driver — MongoDB
- go.opentelemetry.io/otel (+ sdk, exporters OTLP) — observabilidade
- golang.org/x/sync (errgroup, semaphore) — concorrência
- github.com/testcontainers/testcontainers-go — testes de integração
- LLM / BancoVetorial — pendente

## Padrão de commit

Mensagens no formato `CHAVE-JIRA: descricao` (ex.: `LLT-82: inicializacao do modulo: Collector`), conforme a ADR_148. Detalhes em [CONTRIBUTING.md](CONTRIBUTING.md).
