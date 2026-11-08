---
description: "Task list for Consumir logs no OpenObserve"
---

# Tasks: Consumir logs no OpenObserve

**Input**: Design documents from `/specs/002-consumir-logs-openobserve/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Incluídos. `research.md` (decisão 11) e `quickstart.md` pedem prova de domínio e de caso de uso com substitutos das portas. Prova com OpenObserve real e com `otel-logs.zip` não entra como tarefa executável.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Domínio: `internal/domain/captura.go`
- Caso de uso: `internal/application/capturar.go`
- Portas: `internal/application/port/in/capturar.go`, `internal/application/port/out/apontamentos.go`, `internal/application/port/out/leitura.go`, `internal/application/port/out/prontos.go`
- Testes no mesmo pacote, como `internal/config/config_test.go`
- Fora desta lista: `internal/adapter/in/http/`, `internal/adapter/in/queue/`, `internal/adapter/out/mongo/`, `internal/adapter/out/elastic/`, `internal/application/consumir.go`, `cmd/finder/main.go`
- `internal/adapter/out/openobserve/` não é criado nesta lista

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Manter o módulo Go já existente, sem biblioteca nova

- [X] T001 Não adicionar dependência em `go.mod` e não executar `go get`. Não criar pacote fora de `internal/domain/`, `internal/application/`, `internal/application/port/in/` e `internal/application/port/out/`. Não criar `internal/adapter/out/openobserve/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Tipos e portas que as histórias usam. O caso de uso da fila permanece intacto

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T002 Criar apontamento, registro armazenado, log pronto e falha em `internal/domain/captura.go`, reutilizando `domain.Codigo` de `internal/domain/item.go` sem alterar esse arquivo. Regras literais de `data-model.md`: "Valor: Seleciona um único log no OpenObserve. Não é o conteúdo do log. Formato, tamanho e alfabeto não estão definidos."; "Vazio: Inutilizável."; "Conteúdo: O que o OpenObserve devolve para o apontamento. Os campos internos não são descritos."; "Conteúdo: Igual ao registro armazenado. Nenhum campo é ocultado ou reescrito."; "Motivo: Um destes três: inutilizável, ausente, inacessível. Os três são distinguíveis."
- [X] T003 [P] Declarar a porta de entrada `Capturar logs apontados` em `internal/application/port/in/capturar.go`, com as saídas de `specs/002-consumir-logs-openobserve/contracts/captura-log.md`. A fila não é entrada desta porta
- [X] T004 [P] Declarar a porta da fonte mockada em `internal/application/port/out/apontamentos.go`. Obter apontamentos entrega os valores desta execução, ou nenhum, e não lê a fila Redis. Não nomear arquivo, variável ou tabela
- [X] T005 [P] Declarar a porta de leitura em `internal/application/port/out/leitura.go`. Ler um apontamento devolve um único registro, ou informa ausência, ou informa que o armazenamento não pode ser alcançado. O conteúdo segue inteiro. A porta não oculta campo e não classifica ocorrência
- [X] T006 [P] Declarar a porta do log pronto em `internal/application/port/out/prontos.go`. Guardar associa o registro inteiro ao apontamento. Uma segunda guarda do mesmo apontamento não cria outro log. Ler, com a instância no ar, devolve esse único registro. Falha ou interrupção não guarda log. O lugar não tem nome e a porta não promete leitura depois que a instância para

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Capturar o log apontado e deixá-lo pronto (Priority: P1) 🎯 MVP

**Goal**: Um apontamento utilizável deixa um log pronto, igual ao registro armazenado. Com a instância no ar, o mesmo apontamento devolve esse registro. Fonte vazia não gera log nem falha. Parada no meio não deixa log pronto. Depois que a instância para, a etapa não exige que o log continue.

**Independent Test**: Com substitutos, apontar um registro opaco, concluir a captura com a instância no ar e ler de novo pelo mesmo apontamento. O registro volta inteiro. A fila não é chamada. Nenhuma ocorrência é classificada. Esse registro opaco não cumpre SC-001.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T007 [P] [US1] Escrever os testes de domínio em `internal/domain/captura_test.go` para estas regras literais de `data-model.md`: "Sem apontamento: Nenhum log pronto e nenhuma falha"; "Em captura → pronto: O conteúdo fica igual ao armazenado e legível pelo apontamento"; "Leitura: Com a instância no ar, o mesmo apontamento devolve esse único conteúdo"; "Parada no meio da captura não deixa log pronto"; "Depois que a instância para, o modelo não exige que o log continue"; "O conteúdo pronto é o conteúdo armazenado, sem campo oculto"
- [X] T008 [P] [US1] Escrever o teste do caso de uso com substitutos das portas em `internal/application/capturar_test.go`: um apontamento e um registro opaco, instância no ar até o fim, uma leitura seguinte devolve o mesmo registro, a porta da fila não é chamada e nenhuma ocorrência é classificada. Incluir fonte sem apontamento: zero logs prontos e zero falhas

### Implementation for User Story 1

- [X] T009 [US1] Implementar em `internal/domain/captura.go` as transições "Indicado → em captura" e "Em captura → pronto", a leitura com a instância no ar, a interrupção no meio e a ausência de exigência depois que a instância para. O conteúdo pronto permanece igual ao armazenado, sem campo oculto e sem categoria
- [X] T010 [US1] Implementar o caso de uso em `internal/application/capturar.go`, chamando `internal/application/port/out/apontamentos.go`, `internal/application/port/out/leitura.go` e `internal/application/port/out/prontos.go`. Não chamar `internal/application/consumir.go` nem a fila. Não nomear armazenamento

**Checkpoint**: User Story 1 is testable with fakes. SC-001 with otel-logs.zip stays blocked

---

## Phase 4: User Story 2 - Falha de captura com detalhes completos (Priority: P2)

**Goal**: Apontamento vazio, log ausente e armazenamento inacessível produzem falhas distinguíveis. O log não fica pronto. Não há nova tentativa, prazo nem fila de erro.

**Independent Test**: Três capturas com substitutos, uma para cada motivo. Cada falha mostra o apontamento e o motivo. Nenhum desses logs fica pronto.

### Tests for User Story 2 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T011 [P] [US2] Acrescentar em `internal/domain/captura_test.go` as regras literais: "Vazio: Inutilizável. A captura não começa e a falha leva esse motivo."; "Ausente: O apontamento é utilizável e o OpenObserve não tem o log. A falha leva o motivo log ausente. Nenhum registro vazio é criado no lugar."; "Em captura → inacessível: O armazenamento não pode ser alcançado. Falha com esse motivo, distinta da ausência"; "Motivo ausente e motivo inacessível não se substituem."; "Falha não produz log pronto"
- [X] T012 [P] [US2] Acrescentar em `internal/application/capturar_test.go` três casos: apontamento vazio não consulta a leitura; apontamento utilizável sem registro gera falha ausente; leitura que não alcança o armazenamento gera falha inacessível. Cada falha carrega o apontamento e o motivo

### Implementation for User Story 2

- [X] T013 [US2] Implementar em `internal/domain/captura.go` os três motivos distinguíveis e a regra "Indicado → inutilizável: O valor é vazio. Falha com esse motivo. O OpenObserve não é consultado para esse valor". Não criar prazo, nova tentativa nem fila de erro
- [X] T014 [US2] Fazer `internal/application/capturar.go` devolver esses três motivos sem substituir um pelo outro e sem guardar log pronto. Não alterar o fim feliz da User Story 1

**Checkpoint**: User Stories 1 and 2 pass with fakes

---

## Phase 5: User Story 4 - Base local com OpenObserve (Priority: P2)

**Goal**: A base local oferece OpenObserve e deixa de incluir Elasticsearch, com MongoDB e Redis ainda presentes.

**Independent Test**: Bloqueado. Exige versão, porta e credencial do OpenObserve, que não foram informadas. `docker-compose.yml` e `Makefile` ainda sobem Elasticsearch.

### Implementation for User Story 4

- [X] T015 [US4] Não editar `docker-compose.yml`, `Makefile` nem `README.md`. Não escolher imagem, porta nem credencial. Não migrar dado do Elasticsearch. Não apagar `internal/adapter/out/elastic/`. MongoDB e Redis permanecem. A história continua bloqueada até esses dados externos existirem

**Checkpoint**: User Story 4 stays blocked. No local base change in this list

---

## Phase 6: User Story 3 - Capturas paralelas com resultado próprio (Priority: P3)

**Goal**: Dois apontamentos distintos podem ser capturados ao mesmo tempo, cada um com o próprio resultado. O mesmo apontamento, duas vezes, deixa um único log pronto.

**Independent Test**: Disparar duas capturas sobrepostas com substitutos, uma com registro e outra ausente. A primeira fica pronta e a segunda conserva a falha. Repetir o mesmo apontamento e obter um log pronto.

### Tests for User Story 3 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T016 [P] [US3] Acrescentar em `internal/domain/captura_test.go` a regra literal "O mesmo apontamento tem no máximo um log pronto na execução" e "Repetido: Várias indicações do mesmo valor na mesma execução contam como um apontamento"
- [X] T017 [P] [US3] Acrescentar em `internal/application/capturar_test.go` duas capturas de apontamentos distintos ao mesmo tempo, uma com registro e outra ausente: o existente fica pronto e a falha do outro permanece. O substituto deve mostrar que as duas leituras se sobrepõem, sem prazo numérico de produto. Incluir o mesmo apontamento duas vezes: um log pronto

### Implementation for User Story 3

- [X] T018 [US3] Implementar em `internal/domain/captura.go` a convergência de indicações repetidas num único log pronto, na mesma execução, em sequência ou sobrepostas
- [X] T019 [US3] Implementar em `internal/application/capturar.go` a captura simultânea de apontamentos distintos com a biblioteca padrão, sem `golang.org/x/sync` e sem usar `FINDER_WORKERS`. Uma falha não apaga o resultado da outra. Não alterar os fins da User Story 1 e da User Story 2

**Checkpoint**: User Stories 1, 2, and 3 pass with fakes

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Prova do quickstart com substitutos e limite do que não entra nesta feature

- [ ] T020 Executar a prova com substitutos de `specs/002-consumir-logs-openobserve/quickstart.md` com `make test` na raiz do repositório. Registrar a saída observada. Não declarar sucesso sem essa saída. Não executar SC-001 com registro inventado de `otel-logs.zip` nem SC-004 com a base atual
- [X] T021 [P] Confirmar que `cmd/finder/main.go`, `internal/application/consumir.go`, `internal/adapter/in/queue/`, `internal/adapter/in/http/health.go`, `internal/adapter/out/mongo/` e `internal/adapter/out/elastic/` permanecem sem alteração desta feature. Não criar `internal/adapter/out/openobserve/`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on the desired user stories that are not blocked

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Extends `internal/domain/captura.go` and `internal/application/capturar.go` from US1 and stays testable with fakes
- **User Story 4 (P2)**: Blocked. No code task until version, port, and credentials exist. Independent of US1
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Extends the same use case and stays testable with fakes. Does not need the OpenObserve adapter

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Domain transitions before the use case
- Story complete before moving to the next priority
- Do not treat a fake record as a record from `otel-logs.zip`

### Parallel Opportunities

- T003, T004, T005, and T006 can run in parallel after T002
- T007 and T008 can run in parallel
- T011 and T012 can run in parallel
- T016 and T017 can run in parallel
- US4 does not block US1, US2, or US3

---

## Parallel Example: User Story 1

```bash
# Testes da User Story 1, em arquivos diferentes:
Task: "Regras de log pronto em internal/domain/captura_test.go"
Task: "Caso de uso com substitutos em internal/application/capturar_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 with fakes
5. Do not demo SC-001 until `otel-logs.zip` and the OpenObserve access names exist

### Incremental Delivery

1. Setup + Foundational → ports ready
2. User Story 1 → one ready log, readable again while the instance is up
3. User Story 2 → three distinguishable failures
4. User Story 3 → overlapping captures and one ready log per repeated pointer
5. User Story 4 stays blocked with the current base

### Parallel Team Strategy

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2 after the US1 files exist, or in parallel only on new test files before editing `captura.go`
   - Developer C: does not start User Story 4
3. User Story 3 follows the use case from US1 and US2

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Verify tests fail before implementing
- The OpenObserve adapter, the local base change, and SC-001 stay out of this list
- Commit after each task or logical group
- Stop at any checkpoint to validate the story with fakes

---

## Phase 8: Convergence

- [X] T022 CRITICAL: Retirar da captura o contrato inventado do OpenObserve em `internal/adapter/out/openobserve/cliente.go`, `internal/config/config.go` e `docker-compose.yml` (organização `default`, URL, usuário, senha, caminho `api/{org}/_search` e SQL). Esses nomes seguem sem fonte do solicitante per Constitution V e FR-005 (contradicts)
- [X] T023 Fazer um apontamento selecionar um único log armazenado. `Ler` não pode tratar o apontamento como nome de stream e devolver `hits[0]` com `size` 1 per FR-002 (contradicts)
- [X] T024 Não classificar apontamento não vazio como inacessível por falha de `streamAceito`. Formato, tamanho e alfabeto não estão definidos; só o vazio é inutilizável e o armazenamento não é consultado para esse valor per FR-006 (contradicts)
- [X] T025 Não devolver o corpo cru como log pronto quando a resposta não é a consulta esperada. Sucesso só com o registro armazenado, inteiro per FR-003 (contradicts)
