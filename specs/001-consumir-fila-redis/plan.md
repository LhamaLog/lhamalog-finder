# Implementation Plan: Consumir fila Redis

**Branch**: `001-consumir-fila-redis` | **Date**: 2026-10-29 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-consumir-fila-redis/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

O Finder passa a consumir um item livre da fila Redis em stream, no mesmo passo em que reserva e conclui esse consumo, quando o worker e a instância seguem no ar. O item sai da pendência e o código de apontamento fica disponível uma vez para uma etapa futura. O log no OpenObserve não é aberto. Com várias instâncias, o mesmo item tem no máximo uma reserva ativa. Se o worker ou a instância para no meio do passo, o item não conta como concluído e volta a ficar livre, sem prazo numérico. Dois itens com o mesmo código geram uma conclusão.

O desenho cabe na arquitetura hexagonal já descrita no README: regra no domínio, caso de uso na aplicação, fila na entrada e portas de saída para a disponibilidade do código e para a memória de códigos já concluídos. O protótipo `internal/adapter/in/queue/redisStream.go` não é o contrato.

**Gate**: a implementação do adaptador Redis e do registro compartilhado não começa neste plano. Faltam os nomes de stream, grupo e campos (FR-007), o lugar onde o código fica depois da conclusão, e qual item vence quando dois códigos iguais terminam juntos (FR-009). A constituição, princípio V, impede criar esses dados.

## Technical Context

**Language/Version**: Go 1.27, como em `go.mod` e no README (imagem Go 1.27.0)

**Primary Dependencies**: O pacote `internal/adapter/in/queue/redisStream.go` já importa `github.com/redis/go-redis/v9 v9.23.0`, listado em `go.mod` como dependência indireta. O README ainda marca essa biblioteca como candidata e não adicionada. A checagem do Drive da equipe não foi feita nesta sessão. Nenhuma biblioteca nova entra neste plano.

**Storage**: Fila Redis em stream para o item livre, a reserva e a saída da pendência. O lugar em que o código concluído fica disponível para a etapa futura está sem nome. Nenhum nome de stream, grupo, campo, coleção ou índice é escolhido aqui.

**Testing**: `go test -race ./...`, pelo alvo já documentado no README (`make test` ou `docker compose run --rm --no-deps gotools test ./...`). Testes de domínio e de caso de uso usam substitutos em memória. Prova contra um Redis real fica bloqueada até existir o contrato de nomes.

**Target Platform**: Processo do Finder, o mesmo módulo que o README sobe com Docker Compose ao lado de Redis, MongoDB e Elasticsearch. Esta etapa usa a fila. Não grava ocorrência no MongoDB nem no Elasticsearch e não chama o OpenObserve.

**Project Type**: Serviço Go de um módulo, com entrada de fila e health HTTP já existente. O health em `internal/adapter/in/http` fica fora desta feature.

**Performance Goals**: Sem meta numérica. A spec (Q3 A) fixa o aceite em zero reservas simultâneas do mesmo item (SC-001), no passo único de conclusão (SC-002, SC-003) e em uma conclusão por código (SC-005).

**Constraints**: Permanecer em `cmd/finder` e `internal/` já existentes. Não alterar health, MongoDB, Elasticsearch nem o pacote de configuração além do estritamente necessário para ligar o consumo, e essa ligação espera o contrato de nomes. Não usar o stream `fila:logs`, o campo `id` nem o id `110` do protótipo. Não criar prazo numérico, nova tentativa após falha de leitura, nem fila de erro. `FINDER_WORKERS` (padrão 4) existe e não é o grau de paralelismo deste consumo.

**Scale/Scope**: Duas ou mais instâncias do Finder na mesma fila. Vários workers por instância só no que a spec já exige: um worker parado libera o item mesmo com a instância no ar. Volume de itens e vazão não foram informados.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Princípio | Resultado | Leitura |
|-----------|-----------|---------|
| I. Finalidade: logs em ocorrências | Justificado | Esta etapa não gera ocorrência. O pedido e a FR-005 limitam o trabalho ao consumo da fila. A ocorrência continua obrigatória no fluxo do módulo, numa etapa posterior. |
| II. Origem dos logs no coletor | Justificado | O código de apontamento não é o log. Esta etapa não analisa log do LhamaLog-collector. |
| III. Categorias de ocorrência | Justificado | Erro, alerta e informação ficam fora. Nenhuma classe nova é criada. |
| IV. Fila Redis e apontamento no OpenObserve | Justificado em parte, incompleto | A fila traz o código de apontamento. A resolução no OpenObserve e a agregação do log não entram nesta etapa, por pedido explícito. O código precisa continuar disponível, e o lugar disso não foi nomeado. |
| V. Incerteza sem invenção | Falha de gate | Três dados seguem sem fonte: nomes de stream, grupo e campos; lugar do código concluído; item vencedor no empate de códigos iguais. O plano não os cria. |

**Resultado do gate**: ERROR para começar o adaptador Redis e o registro compartilhado. O desenho de domínio, portas e testes com substitutos pode ser escrito. A ligação com Redis e com o destino do código espera resposta do solicitante.

### Revisão depois do desenho

`research.md`, `data-model.md`, `contracts/consumo-fila.md` e `quickstart.md` repetem os mesmos três buracos. Nenhum nome de fila, lugar do código ou vencedor de empate foi criado. O gate segue em ERROR para a implementação do adaptador.

## Project Structure

### Documentation (this feature)

```text
specs/001-consumir-fila-redis/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── consumo-fila.md
└── tasks.md             # criado por /speckit-tasks, não por este comando
```

### Source Code (repository root)

```text
cmd/finder/                              # wiring do consumo, no lugar de InitRedis
internal/domain/                         # item, reserva, conclusão, uma conclusão por código
internal/application/                    # caso de uso do passo de consumo
internal/application/port/in/            # entrada do caso de uso
internal/application/port/out/           # fila, códigos já concluídos, disponibilidade do código
internal/adapter/in/queue/               # adaptador Redis, bloqueado até o contrato de nomes
internal/config/                         # só se a ligação do consumo precisar de configuração já autorizada
```

Pacotes que este plano não altera: `internal/adapter/in/http`, `internal/adapter/out/mongo`, `internal/adapter/out/elastic`.

**Structure Decision**: O README já fixa a hexagonal leve nesses diretórios. O consumo novo entra neles. O arquivo `internal/adapter/in/queue/redisStream.go` deixa de ser o comportamento do Finder quando o adaptador real for ligado. Até lá, o plano não reescreve esse arquivo com nomes inventados.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Etapa sem ocorrência e sem OpenObserve (princípios I e IV) | O solicitante limitou o trabalho ao consumo da fila e à disponibilidade do código | Fazer agora a leitura do log e a classificação inventaria o fluxo que a spec deixou de fora |
| Portas de saída sem adaptador nomeado | O código concluído e a unicidade do código precisam de um lugar compartilhado entre instâncias | Guardar isso só na memória de um processo faria duas instâncias concluírem o mesmo código e perderia o código quando o processo parasse. O lugar continua sem nome |
