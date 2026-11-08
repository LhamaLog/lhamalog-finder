# Implementation Plan: Consumir logs no OpenObserve

**Branch**: `002-consumir-logs-openobserve` | **Date**: 2026-10-31 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/002-consumir-logs-openobserve/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

O Finder obtém apontamentos de uma fonte mockada, captura um log por apontamento no OpenObserve e o deixa pronto para processamento. O log pronto é o registro armazenado, inteiro e igual ao original. Enquanto a instância está no ar, o mesmo apontamento devolve esse único log. Depois que a instância para, esta etapa não exige que ele continue. O lugar não tem nome. A fila Redis não é lida. Não há agregação, categoria de ocorrência nem ocultação de dados sensíveis.

O desenho cabe na hexagonal já descrita no README, num caso de uso novo, separado do consumo da fila da feature 001. A prova automatizada desta fase usa substitutos das portas. O adaptador do OpenObserve, a troca concreta da base local e o aceite com `otel-logs.zip` não começam neste plano.

**Gate**: a constituição, princípio V, impede criar o que a spec deixou sem fonte. Faltam o arquivo `otel-logs.zip`, os nomes de organização, stream, índice, campos e credenciais, e a versão, a porta e a credencial da base local. O plano não os cria.

## Technical Context

**Language/Version**: Go 1.27, como em `go.mod` e no README (imagem Go 1.27.0)

**Primary Dependencies**: Nenhuma biblioteca nova. O cliente HTTP da biblioteca padrão basta quando o contrato do OpenObserve existir. `github.com/redis/go-redis/v9` permanece no consumo da fila e não entra nesta captura. Candidatas do README, inclusive `golang.org/x/sync` e o cliente de Elasticsearch, não são tratadas como presentes.

**Storage**: O log a capturar está no OpenObserve. O log pronto fica legível pelo apontamento só enquanto a instância está no ar. O lugar desse log pronto não tem nome e não é MongoDB, Redis nem arquivo. Elasticsearch sai da base local quando a troca puder ser escrita sem versão, porta ou credencial inventadas. MongoDB e Redis permanecem na base e esta captura não os usa.

**Testing**: `go test -race ./...`, pelo alvo já documentado no README (`make test` ou `docker compose run --rm --no-deps gotools test -race ./...`). Testes de domínio e de caso de uso usam substitutos. A prova com OpenObserve real e com um registro de `otel-logs.zip` fica bloqueada. Nenhum desses testes foi executado nesta sessão de planejamento.

**Target Platform**: Processo do Finder. A base local hoje sobe MongoDB, Redis e Elasticsearch, conforme `docker-compose.yml`, `Makefile` e README, lidos em 2026-10-31.

**Project Type**: Serviço Go de um módulo. Health HTTP, fila Redis e o caso de uso `application.Servico` ficam fora desta feature.

**Performance Goals**: Sem meta numérica. O aceite é duas capturas distintas ao mesmo tempo, com resultados próprios (SC-003), e um log pronto por apontamento repetido (SC-006).

**Constraints**: Permanecer em `cmd/finder` e `internal/` já existentes. Não reescrever a feature 001. Não alterar health, MongoDB nem o pacote `internal/adapter/out/elastic`. Não criar prazo, nova tentativa nem fila de erro. `FINDER_WORKERS` (padrão 4) não é o grau de paralelismo desta captura. Não nomear organização, stream, índice, campo, credencial, imagem, porta nem o conteúdo de `otel-logs.zip`.

**Scale/Scope**: Uma execução do Finder. Dois ou mais apontamentos distintos na mesma execução. Várias instâncias do Finder coordenando a mesma captura ficam fora. Volume e vazão não foram informados.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Princípio | Resultado | Leitura |
|-----------|-----------|---------|
| I. Finalidade: logs em ocorrências | Justificado | Esta etapa captura e prepara o log. Não gera ocorrência. O pedido e a FR-004 limitam o trabalho a isso. A ocorrência continua obrigatória numa etapa posterior. |
| II. Origem dos logs no coletor | Justificado em parte | O log capturado é o registro já armazenado no OpenObserve. Esta etapa não analisa nem agrega o log do LhamaLog-collector. |
| III. Categorias de ocorrência | Justificado | Erro, alerta e informação ficam fora. Dados sensíveis não são ocultados aqui. Nenhuma classe nova é criada. |
| IV. Fila Redis e apontamento no OpenObserve | Justificado em parte, incompleto | O pedido substitui a fila, nesta etapa, pela fonte mockada, e pede a captura sem agregação. A constituição não é emendada. O fluxo normativo continua sendo fila, depois OpenObserve, depois agregação. |
| V. Incerteza sem invenção | Falha de gate | Seguem sem fonte: o arquivo `otel-logs.zip` e os registros dele; organização, stream, índice, campos e credenciais; versão, porta e credencial do OpenObserve na base local; os passos que colocam um registro do arquivo no OpenObserve. O plano não os cria. |

**Resultado do gate**: ERROR para o adaptador do OpenObserve, para a edição de `docker-compose.yml`, `Makefile` e README, e para executar SC-001. O desenho de domínio, portas e testes com substitutos pode ser escrito.

### Revisão depois do desenho

`research.md`, `data-model.md`, `contracts/captura-log.md` e `quickstart.md` repetem os mesmos buracos. Nenhum nome de acesso, imagem, porta, credencial ou campo do arquivo foi criado. O gate segue em ERROR para o adaptador, a base local e o aceite com o arquivo.

## Project Structure

### Documentation (this feature)

```text
specs/002-consumir-logs-openobserve/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── captura-log.md
└── tasks.md             # criado por /speckit-tasks, não por este comando
```

### Source Code (repository root)

```text
cmd/finder/                                      # liga a captura à fonte mockada, sem ler a fila
internal/domain/                                 # apontamento, log pronto, falha, um log por apontamento
internal/application/                            # caso de uso novo da captura
internal/application/port/in/                    # entrada da captura
internal/application/port/out/                   # fonte mockada, leitura do log, logs prontos
internal/adapter/out/openobserve/                # adaptador bloqueado até o contrato externo
```

Pacotes que este plano não altera: `internal/adapter/in/http`, `internal/adapter/in/queue`, `internal/adapter/out/mongo`, `internal/adapter/out/elastic`, e o caso de uso `application.Servico`.

**Structure Decision**: O README já fixa a hexagonal leve nesses diretórios. A captura entra num caso de uso novo. O consumo da fila permanece o da feature 001. O adaptador do OpenObserve e a troca da base local esperam os dados que o gate recusou criar.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Etapa sem fila e sem ocorrência (princípios I e IV) | O solicitante limitou esta etapa ao apontamento mockado e à captura, com o log pronto e sem processamento | Ler a fila vazia ou classificar ocorrência inventaria o fluxo que a spec deixou de fora |
| Log pronto sem lugar nomeado | A resposta de 2026-10-31 pede leitura pelo mesmo apontamento só enquanto a instância está no ar, e proíbe nomear o lugar | Gravar em MongoDB, Redis ou arquivo nomearia um armazenamento que a spec não definiu |
