---
description: "Task list for Capturar referência de lote"
---

# Tasks: Capturar referência de lote

**Input**: Design documents from `/specs/004-capturar-referencia-lote/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Incluídos. A spec pede provas unitárias e integradas com Redis de teste e dados sintéticos. `quickstart.md` separa `make test` sem Redis da prova com `FINDER_REDIS_INTEGRATION=1`. Os testes de cada história vêm antes da implementação e devem falhar antes dela.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Processo: `cmd/finder/main.go`
- Configuração: `internal/config/config.go`
- Caso de uso: `internal/application/consumir.go`
- Portas: `internal/application/port/out/`
- Fila: `internal/adapter/in/queue/redisStream.go`
- Domínio: `internal/domain/envelope.go` e `internal/domain/chave.go`
- Testes no mesmo pacote dos arquivos que exercitam
- Ambiente: `docker-compose.yml` e `.env.example`
- Fora desta lista: `internal/adapter/in/http/health.go`, `internal/adapter/out/openobserve/`, `internal/adapter/out/mongo/`, `internal/adapter/out/elastic/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Manter o módulo Go já existente e apontar o ambiente para o stream do contrato, sem número padrão

- [X] T001 Confirmar em `go.mod` que a dependência permanece `github.com/redis/go-redis/v9 v9.23.0`. Não executar `go get` e não acrescentar biblioteca
- [X] T002 [P] Em `docker-compose.yml` e `.env.example`, trocar o stream usado por este módulo de `fila:logs` para `lhamalog:stream:logs`. Acrescentar os nomes `FINDER_REDIS_CONSUMIDOR`, `FINDER_REDIS_IDEMPOTENCIA`, `FINDER_REDIS_COUNT`, `FINDER_REDIS_BLOCK` e `FINDER_REDIS_PENDENCIA_OCIOSA` sem valor numérico padrão. Não gravar senha. Não reutilizar `finder:codigos` como prefixo da chave desta feature

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Configuração sem padrão e a guarda da prova Redis, usadas por todas as histórias

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T003 Em `internal/config/config.go`, ler stream, grupo, consumidor, prefixo da chave, quantidade por leitura, tempo de bloqueio, ociosidade da pendência e senha. `Carregar` não preenche quantidade, bloqueio, ociosidade, grupo, consumidor nem prefixo. A senha vazia permanece vazia. Expor uma validação que recusa stream vazio ou diferente de `lhamalog:stream:logs`, grupo vazio, consumidor vazio, prefixo vazio, e quantidade, bloqueio ou ociosidade ausentes ou não positivos. `FINDER_WORKERS` continua com o padrão 4 e não é o grau de paralelismo deste consumo
- [X] T004 Em `internal/config/config_test.go`, provar que configuração vazia não cria quantidade, bloqueio, ociosidade, grupo, consumidor nem prefixo, e que a validação de T003 aponta cada ausência. Um `.env` temporário pode citar nomes de exemplo; esses nomes não viram padrão de `Carregar`. Não registrar a senha no resultado do teste
- [X] T005 [P] Em `internal/adapter/in/queue/redis_integration_test.go`, a prova Redis só executa com `FINDER_REDIS_INTEGRATION=1`. Se a variável está definida e `FINDER_REDIS_ADDR` não responde, o teste falha. Se não está definida, o teste é ignorado. Stream, grupo, consumidor, prefixo, quantidade, bloqueio e ociosidade vêm do ambiente do teste. Não copiar esses valores para o código de produção

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Interpretar um lote válido e deixá-lo pronto para a etapa seguinte (Priority: P1) 🎯 MVP

**Goal**: Um envelope `log.ingested` v1 válido vira referência entregue à etapa de recuperação. A interpretação não consulta o OpenObserve e não confirma a entrada.

**Independent Test**: Publicar um lote sintético válido, com e sem `traceId`, em lote de um registro e de vários. A referência traz os campos de consulta, o tenant é o do metadata, a interpretação não busca conteúdo de log e a entrada segue sem confirmação.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T006 [P] [US1] Escrever em `internal/domain/envelope_test.go` o envelope válido de `specs/004-capturar-referencia-lote/data-model.md`: `event` `log.ingested`, `origin` `collector`, `schemaVersion` texto `"1"`, `publishedAt` instante ISO UTC, `signal` `LOGS`, `productId` igual nos dois lados, `idempotencyId` igual a `openobserve.batchId`, faixa com mínimo menor ou igual ao máximo, `recordCount` inteiro positivo. `recordCount` 1 inclui `recordId` preenchido. `recordCount` maior que 1 traz `recordId` nulo. Com `traceId` e sem `traceId` a referência continua válida. `traceId` não deduplica
- [X] T007 [P] [US1] Escrever em `internal/application/consumir_test.go` o caso de `specs/004-capturar-referencia-lote/contracts/referencia-lote.md`: ao fim da interpretação a porta de efeito recebe `tenantId` e `productId` do metadata, organização, stream do OpenObserve, `batchId`, a faixa, `recordCount`, `idempotencyId`, a correlação presente e o id da entrada Redis. A contagem de confirmações é zero e nenhuma consulta ao OpenObserve ocorre

### Implementation for User Story 1

- [X] T008 [P] [US1] Criar em `internal/domain/envelope.go` a leitura dos campos de texto do envelope e a referência validada. `metadata` e `body` são JSON em texto. O body não é o conteúdo dos logs. `tenantId` e `productId` entregues são os do metadata. `correlationId` acompanha e não deduplica. Ausência de `traceId`, `receivedAt`, `bytes`, `serviceNames`, `priority`, `criticality`, `maxSeverityNumber` ou `severityCounts` não invalida. Enumeração de prioridade e criticidade não é revalidada. `recordCount`, `bytes` e os agregados descrevem os registros aceitos pelo Collector
- [X] T009 [P] [US1] Criar em `internal/application/port/out/efeito.go` a porta que recebe a referência validada e o id da entrada. O resultado desta história é a entrega. A porta não confirma a entrada e não descreve a consulta ao OpenObserve
- [X] T010 [US1] Em `internal/adapter/in/queue/redisStream.go`, deixar de tratar a mensagem como um único código. Ler os campos do envelope como texto, preservar o id Redis e não confirmar a entrada válida na interpretação
- [X] T011 [US1] Em `internal/application/consumir.go`, montar a referência validada e entregá-la à porta de T009. Não confirmar a entrada. Não abrir o OpenObserve. Não usar `correlationId` nem o id Redis como chave de idempotência

**Checkpoint**: User Story 1 is testable. A valid batch reference is delivered and the Redis entry stays unconfirmed

---

## Phase 4: User Story 2 - Impedir que referência inválida siga como lote consultável (Priority: P1)

**Goal**: Referência inválida não segue para a recuperação. O motivo fica registrado sem senha e sem conteúdo de log. A entrada é confirmada uma vez e não recebe o estado concluído.

**Independent Test**: Cada entrada inválida da história 2 da spec produz zero referências válidas, é confirmada, não permanece pendente e o registro não contém senha nem conteúdo de log.

### Tests for User Story 2 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T012 [P] [US2] Estender `internal/domain/envelope_test.go` com JSON inválido, campo necessário ausente, `event` diferente de `log.ingested`, `origin` diferente de `collector`, `schemaVersion` diferente do texto `"1"`, `signal` diferente de `LOGS`, `productId` divergente, `idempotencyId` diferente de `batchId`, mínimo da faixa maior que o máximo, `recordCount` zero ou negativo, lote de um registro sem `recordId` preenchido, lote de vários registros com `recordId` preenchido, e `publishedAt` que não é instante ISO UTC. Schema 2 não é versão aceita
- [X] T013 [P] [US2] Estender `internal/application/consumir_test.go` com o contrato de entrada inválida: não há referência válida, a entrada é confirmada uma vez, não é tentada de novo, não permanece pendente e não grava estado concluído. O registro leva o id da entrada, o `idempotencyId` quando legível e o motivo, sem senha, sem metadata bruta, sem body bruto e sem conteúdo de log
- [X] T014 [P] [US2] Em `internal/adapter/in/queue/redis_integration_test.go`, publicar uma entrada sintética inválida e conferir que ela sai da pendência, não produz referência válida e não deixa chave concluída

### Implementation for User Story 2

- [X] T015 [P] [US2] Em `internal/domain/envelope.go`, devolver o motivo das rejeições de T012 sem incluir o texto bruto de `metadata` ou `body`
- [X] T016 [P] [US2] Em `internal/adapter/in/queue/redisStream.go`, confirmar a entrada inválida pelo id Redis. Se a chave desse `idempotencyId` estiver em processamento, removê-la. Não gravar estado concluído. Não criar retry, spool ou DLQ
- [X] T017 [US2] Em `internal/application/consumir.go`, encerrar a entrada inválida por esse caminho e registrar o motivo com `log/slog` sem senha e sem conteúdo de log

**Checkpoint**: User Stories 1 and 2 work. Invalid references are confirmed once and never delivered as valid

---

## Phase 5: User Story 3 - Não perder o lote nem repetir o efeito concluído (Priority: P2)

**Goal**: Uma chave por `idempotencyId`, em processamento ou concluída. Só o estado concluído dura 25 horas. A segunda entrada válida espera. O efeito da recuperação, quando informado como concluído, autoriza a confirmação.

**Independent Test**: Com a chave em processamento e sem ociosidade autorizada, a segunda entrada permanece pendente, sem outro efeito e sem confirmação. Depois do efeito concluído, essa entrada é confirmada sem repetir o efeito. O prazo observado na chave concluída é 25 horas.

### Tests for User Story 3 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T018 [P] [US3] Escrever em `internal/domain/chave_test.go` os estados de `specs/004-capturar-referencia-lote/data-model.md`: ausente, `em_processamento` sem prazo, e `concluido` com 25 horas contadas uma vez. Passagem repetida para concluído não renova o prazo. Depois de 25 horas a chave não suprime. Remoção da reserva não grava `concluido`
- [X] T019 [P] [US3] Estender `internal/application/consumir_test.go`: segunda entrada válida com chave `em_processamento` e ociosidade ainda não autorizada permanece pendente, não inicia outro efeito e não é confirmada. Efeito informado como concluído confirma sem repetir. Queda antes do efeito deixa a contagem de efeitos concluídos inalterada. Queda depois do efeito e antes da confirmação não repete o efeito e a reentrega confirma. O substituto da porta informa o efeito. O teste não descreve a consulta ao OpenObserve
- [X] T020 [P] [US3] Em `internal/adapter/in/queue/redis_integration_test.go`, conferir uma chave por `idempotencyId` no prefixo do ambiente, valor `em_processamento` sem prazo de 25 horas, valor `concluido` com prazo de 25 horas lido pelo tempo restante, e duas entradas Redis diferentes com o mesmo identificador sem segundo efeito enquanto a primeira não concluiu

### Implementation for User Story 3

- [X] T021 [P] [US3] Criar em `internal/domain/chave.go` a regra da chave: criar `em_processamento` só se ausente e sem prazo; gravar `concluido` e 25 horas só quando o efeito terminou; não renovar o prazo se já estiver concluída; remover a reserva sem gravar `concluido`
- [X] T022 [P] [US3] Criar em `internal/application/port/out/chave.go` as operações reservar, concluir, consultar e remover reserva. O prefixo é externo. A chave do Finder não serve à Varredura
- [X] T023 [US3] Em `internal/adapter/in/queue/redisStream.go`, gravar a chave sob `FINDER_REDIS_IDEMPOTENCIA`. Não usar o conjunto `FINDER_REDIS_CODIGOS` para esse estado. Enquanto a chave está `em_processamento` e a ociosidade configurada não autoriza reassunção, a outra entrada permanece pendente e não é lida de novo na mesma passagem. Quando a chave passa a `concluido`, confirmar essa entrada sem repetir o efeito
- [X] T024 [US3] Em `internal/application/consumir.go`, chamar a porta de efeito só para a entrada autorizada a executar o efeito. Confirmar a entrada válida pelo id Redis somente depois do efeito concluído. Porta ausente ou falha antes desse resultado mantém a entrada pendente e a chave `em_processamento`. Não criar fila, spool nem outro armazenamento intermediário

**Checkpoint**: User Stories 1, 2, and 3 work. A completed effect is not repeated within 25 hours, and an in-progress key does not confirm the second entry

---

## Phase 6: User Story 4 - Consumir no grupo do Finder e recuperar pendência (Priority: P2)

**Goal**: O grupo do Finder nasce no início do que ainda está no stream, separado do grupo da Varredura. Pendência abandonada volta pela ociosidade informada. Sem os três parâmetros, o consumo não abre.

**Independent Test**: Entradas publicadas antes do grupo são lidas. O grupo do Finder e outro grupo observam a mesma entrada. Sem quantidade, bloqueio ou ociosidade, nenhuma entrada é lida. Depois da ociosidade informada, a pendência pode ser reivindicada e não é descartada.

### Tests for User Story 4 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T025 [P] [US4] Estender `internal/config/config_test.go` com três arranques, cada um sem quantidade, sem bloqueio ou sem ociosidade. Em cada um a validação recusa o consumo e nenhum valor padrão é criado
- [X] T026 [P] [US4] Em `internal/adapter/in/queue/redis_integration_test.go`, cobrir os cenários de `specs/004-capturar-referencia-lote/quickstart.md`: entradas já presentes são lidas e a criação não usa `$`; o grupo configurado e um segundo grupo observam a mesma entrada; pendência que atinge a ociosidade informada pode ser reivindicada e não é descartada; Redis indisponível no meio não confirma entrada inacabada. A quantidade da reivindicação é a mesma quantidade da leitura

### Implementation for User Story 4

- [X] T027 [US4] Em `internal/adapter/in/queue/redisStream.go`, criar o grupo em `lhamalog:stream:logs` com id `0` quando ele ainda não existe. Não usar `$`. Não apagar grupo já existente. Ler com a quantidade e o bloqueio informados. A leitura de pendência já deste consumidor não espera o tempo de bloqueio. Reivindicar só depois da ociosidade informada, usando essa mesma quantidade, no lugar da regra de conexão viva com ociosidade zero. Redis indisponível não confirma entrada inacabada
- [X] T028 [US4] Em `cmd/finder/main.go`, chamar a validação de T003 antes de abrir a fila. Se ela recusar, encerrar o processo sem ler entrada e sem criar número. Não registrar a senha. Passar a senha do ambiente ao cliente Redis; senha vazia continua sem senha, como o Compose atual. Não alterar `internal/adapter/in/http/health.go`

**Checkpoint**: All user stories work. The group reads existing entries, pending entries can be reclaimed, and missing parameters refuse startup

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Conferir os limites que atravessam as histórias, sem afirmar execução que não ocorreu

- [X] T029 [P] Não criar cliente de busca em `internal/adapter/out/openobserve/`. A porta em `internal/application/port/out/efeito.go` permanece sem adaptador de consulta
- [X] T030 Conferir `specs/004-capturar-referencia-lote/quickstart.md` contra o comportamento implementado. Não escrever resultado de comando que não tenha sido executado nesta tarefa
- [X] T031 [P] Em `internal/application/consumir.go`, emitir com `log/slog` uma linha por entrada lida, e em `cmd/finder/main.go` apontar `slog.Default` para o handler JSON que já escreve em `os.Stdout`, para o feedback aparecer no terminal. Leitura sem entrada não escreve linha. Falha de leitura mantém só o `log.Error` já existente em `cmd/finder/main.go`, sem segunda linha. Os desfechos, um por entrada, são: `referência entregue` quando a referência válida foi montada e a entrada segue sem confirmação, inclusive com a porta de efeito ausente ou com o efeito chamado e não concluído; `entrada inválida` na linha que T017 já grava, sem duplicá-la; `efeito concluído` quando a chave passa a concluída e a entrada é confirmada; `confirmada sem repetir efeito` em `confirmarSemEfeito`; `entrada pendente` quando a chave está em processamento e a entrada é adiada. Cada linha leva o id da entrada Redis, o `idempotencyId` quando legível e o desfecho. Não incluir senha, metadata bruta, body bruto nem conteúdo de log. Em `internal/application/consumir_test.go`, provar cada desfecho e que a leitura sem entrada não gera linha. A guarda de T013 continua: o registro não contém senha, metadata bruta, body bruto nem conteúdo de log

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: sem dependência
- **Foundational (Phase 2)**: depende do Setup e bloqueia as histórias
- **User Stories (Phase 3+)**: dependem do Foundational
- **Polish (Phase 7)**: depende das histórias desejadas

### User Story Dependencies

- **User Story 1 (P1)**: começa depois do Foundational. Não depende das outras histórias
- **User Story 2 (P1)**: começa depois do Foundational. Usa o envelope da User Story 1 e é testável com entradas inválidas próprias
- **User Story 3 (P2)**: começa depois da User Story 1. A confirmação válida e a chave usam a referência já interpretada. A remoção da reserva inválida encontra a User Story 2
- **User Story 4 (P2)**: a recusa de arranque e o id `0` podem ser provados depois do Foundational. A reivindicação da chave em processamento encontra a User Story 3 em `internal/adapter/in/queue/redisStream.go`, então essa parte segue a User Story 3

### Within Each User Story

- Testes antes da implementação, e falhando antes dela
- Modelo antes do caso de uso
- Caso de uso antes da ligação em `cmd/finder/main.go`
- A história fica testável no checkpoint antes da próxima prioridade

### Parallel Opportunities

- T002 pode seguir junto com T001
- T005 pode seguir junto com T003. T004 vem depois de T003 e afirma a validação
- T006 e T007 podem ser escritos juntos
- T008 e T009 podem ser implementados juntos
- T012, T013 e T014 podem ser escritos juntos
- T015 e T016 podem ser implementados juntos
- T018, T019 e T020 podem ser escritos juntos
- T021 e T022 podem ser implementados juntos
- T025 e T026 podem ser escritos juntos
- T029 pode seguir enquanto T030 confere o quickstart
- T031 usa arquivos diferentes de T029 e T030 e pode seguir em paralelo com elas. Dentro de T031, o teste em `internal/application/consumir_test.go` vem antes da linha nova em `internal/application/consumir.go`

---

## Parallel Example: User Story 1

```text
T006 Escrever envelope válido em internal/domain/envelope_test.go
T007 Escrever entrega sem confirmação em internal/application/consumir_test.go

T008 Criar o envelope em internal/domain/envelope.go
T009 Criar a porta de efeito em internal/application/port/out/efeito.go
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Completar Phase 1: Setup
2. Completar Phase 2: Foundational
3. Completar Phase 3: User Story 1
4. Parar e validar a referência entregue sem confirmação e sem consulta ao OpenObserve

### Incremental Delivery

1. Setup + Foundational
2. User Story 1: lote válido interpretado
3. User Story 2: inválida confirmada uma vez e fora da recuperação
4. User Story 3: chave, 25 horas e espera da segunda entrada
5. User Story 4: grupo no início do stream, pendência e recusa sem parâmetro
6. Polish: busca no OpenObserve continua ausente
7. T031: cada entrada lida deixa um desfecho visível no stdout do processo

### Parallel Team Strategy

1. Todos fecham Setup e Foundational
2. Depois disso, uma pessoa pode escrever os testes da User Story 1 enquanto outra escreve os da User Story 2
3. User Story 3 e a reivindicação da User Story 4 compartilham `internal/adapter/in/queue/redisStream.go` e ficam em sequência

---

## Notes

- [P] marca arquivos diferentes, sem dependência incompleta
- Nenhum tarefa grava quantidade, bloqueio, ociosidade, grupo, consumidor ou prefixo como padrão de produção
- A posição `$` não é usada na criação do grupo
- `internal/adapter/in/http/health.go` não muda de significado
- Commit ao fim de cada tarefa ou grupo lógico
- Parar em qualquer checkpoint para validar a história sozinha
