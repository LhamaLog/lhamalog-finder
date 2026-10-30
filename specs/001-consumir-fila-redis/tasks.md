---
description: "Task list for Consumir fila Redis"
---

# Tasks: Consumir fila Redis

**Input**: Design documents from `/specs/001-consumir-fila-redis/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Incluídos. `research.md` (decisão 10) e `quickstart.md` pedem prova de domínio e de caso de uso com substitutos das portas. Prova com Redis real não entra nesta lista.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Domínio: `internal/domain/`
- Caso de uso: `internal/application/`
- Portas: `internal/application/port/in/` e `internal/application/port/out/`
- Fila: `internal/adapter/in/queue/` permanece bloqueada
- Testes no mesmo pacote, como `internal/config/config_test.go`
- Fora desta lista: `internal/adapter/in/http/`, `internal/adapter/out/mongo/`, `internal/adapter/out/elastic/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Manter o módulo Go já existente, sem biblioteca nova

- [X] T001 Não adicionar módulo em go.mod. O import já presente é `github.com/redis/go-redis/v9 v9.23.0`, indireto. Não executar `go get`. Não criar pacote fora de `internal/domain/`, `internal/application/`, `internal/application/port/in/`, `internal/application/port/out/` e `internal/adapter/in/queue/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Tipos e portas que as duas histórias usam

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T002 Criar o item da fila e o código de apontamento em `internal/domain/item.go`. A identidade do item é distinta do código. O conteúdo aceito é só o código de apontamento. Item sem código, ou com conteúdo além do código, está fora do contrato e não ganha tratamento novo. Não definir formato, tamanho nem alfabeto do código
- [X] T003 [P] Declarar a porta de entrada `Consumir item livre` em `internal/application/port/in/consumir.go`, com as saídas de `specs/001-consumir-fila-redis/contracts/consumo-fila.md`
- [X] T004 [P] Declarar a porta de saída da fila em `internal/application/port/out/fila.go`, com tomar item livre, concluir e liberar, sem comando Redis e sem nome de stream, grupo ou campo
- [X] T005 [P] Declarar a porta de códigos concluídos em `internal/application/port/out/codigos.go`. Registrar a primeira conclusão de um código aceita uma vez e rejeita as seguintes. O lugar físico não tem nome
- [X] T006 [P] Declarar a porta de código disponível em `internal/application/port/out/disponibilidade.go`. Disponibilizar guarda o código para a etapa futura, uma vez, sem abrir o OpenObserve. O lugar não tem nome

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Consumir um item e concluir o consumo (Priority: P1) 🎯 MVP

**Goal**: No mesmo passo, com worker e instância no ar, o item sai da pendência e o código fica disponível uma vez. Parada do worker no meio do passo não conclui o item. Um código tem uma conclusão.

**Independent Test**: Colocar um item cujo conteúdo é só o código de apontamento e verificar que, no mesmo passo, o item sai da pendência, o código continua disponível uma vez e o OpenObserve não é aberto. Interromper o worker no meio do passo e verificar que o item não conta como concluído.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T007 [P] [US1] Escrever os testes de domínio em `internal/domain/passo_test.go` para estas regras literais de `data-model.md`: "Com o worker e a instância no ar, o passo livre termina em conclusão, se o código ainda não foi concluído"; "Código já concluído não gera segunda conclusão nem reserva restante"; "Parada no meio do passo não grava conclusão"; "Dois itens com o mesmo valor que terminam juntos produzem uma conclusão e nenhuma reserva restante. Qual item vence não está no modelo"; "Fila sem item livre, ou falha de leitura, não gravam conclusão"
- [X] T008 [P] [US1] Escrever o teste do caso de uso com substitutos das portas em `internal/application/consumir_test.go`: um item com só o código, worker e instância no ar, produz uma conclusão, item fora da pendência e código disponível uma vez, sem chamada ao OpenObserve

### Implementation for User Story 1

- [X] T009 [US1] Implementar as transições do passo em `internal/domain/passo.go`: livre para em passo; em passo para concluído quando o worker e a instância seguem no ar e o código ainda não foi concluído; em passo para livre quando o worker para no meio, com a instância no ar, sem prazo numérico; em passo para não concluído e não reservado quando já existe conclusão daquele código. Não escolher qual item vence se os dois terminam juntos: o resultado exigido é uma conclusão e nenhuma reserva restante
- [X] T010 [US1] Implementar o caso de uso em `internal/application/consumir.go`, chamando as portas de `internal/application/port/out/fila.go`, `internal/application/port/out/codigos.go` e `internal/application/port/out/disponibilidade.go`. No fim normal, reserva e conclusão são o mesmo passo. Não abrir o OpenObserve e não nomear armazenamento

**Checkpoint**: User Story 1 is testable with fakes, without Redis names

---

## Phase 4: User Story 2 - Várias instâncias sem corrida (Priority: P2)

**Goal**: O mesmo item tem no máximo uma reserva ativa. Outra instância recebe outro item livre, ou nenhum. Se a instância para no meio do passo, o item não é concluído e pode ser reservado de novo, sem prazo numérico.

**Independent Test**: Dois tomadores contra a mesma fila, com substitutos. Cada item fica reservado por um só enquanto o passo segue. A parada de uma instância, recebida como fato e não como prazo, libera o item para o outro tomador.

### Tests for User Story 2 ⚠️

- [X] T011 [P] [US2] Escrever o teste de exclusão em `internal/domain/reserva_test.go`: "Com mais de uma instância do Finder na mesma fila, o mesmo item MUST ser reservado por no máximo uma instância de cada vez"; durante o passo, o outro tomador recebe outro item livre ou nenhum; "A instância para depois da reserva e antes da confirmação da conclusão: o item não conta como concluído, a reserva acaba e outra instância pode reservar o mesmo item. Nenhum prazo numérico se aplica."

### Implementation for User Story 2

- [X] T012 [US2] Implementar a reserva exclusiva em `internal/domain/reserva.go`. A liberação por parada da instância só ocorre quando esse fato é informado ao domínio. Não criar prazo, batimento nem tempo de espera
- [X] T013 [US2] Usar `internal/domain/reserva.go` em `internal/application/consumir.go` para o segundo tomador, sem alterar o fim normal da User Story 1
- [X] T014 [US2] Não modificar `internal/adapter/in/queue/redisStream.go` nem `cmd/finder/main.go`. O adaptador Redis continua bloqueado até o solicitante informar nomes de stream, grupo e campos, o lugar do código concluído, qual item vence no empate e o sinal de parada da instância sem prazo numérico. Não usar o stream `fila:logs`, o campo `id` nem o id `110`

**Checkpoint**: User Stories 1 and 2 pass with fakes. Real Redis stays unwired

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Prova do quickstart e limite do que não entra nesta feature

- [X] T015 Executar a prova com substitutos de `specs/001-consumir-fila-redis/quickstart.md` com `make test` na raiz do repositório. Registrar a saída observada. Não declarar sucesso sem essa saída. Não criar teste de integração Redis com nome de stream inventado
- [X] T016 [P] Confirmar que `internal/adapter/in/http/health.go`, `internal/adapter/out/mongo/` e `internal/adapter/out/elastic/` permanecem sem alteração desta feature

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Uses `internal/application/consumir.go` from US1 and stays testable with fakes, without the Redis adapter

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Domain rules before the use case
- Use case before any adapter
- T014 stays blocked. It does not unlock Redis wiring

### Parallel Opportunities

- T003, T004, T005 and T006 touch different files and can run together after T002
- T007 and T008 can run together
- US1 and US2 can start together after Phase 2, from T007 and T011, until T013 needs `internal/application/consumir.go`
- T016 can run beside T015

---

## Parallel Example: User Story 1

```bash
# Testes da história 1, arquivos diferentes:
Task: "T007 internal/domain/passo_test.go"
Task: "T008 internal/application/consumir_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: `make test` for the US1 fakes
5. Do not wire Redis. Names, storage place, tie-break and cross-instance stop signal are still missing

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently with fakes (MVP)
3. Add User Story 2 → Test the exclusive reservation with fakes
4. Leave T014 blocked until the requester supplies the four missing facts

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2 tests and `internal/domain/reserva.go`
3. T013 waits for `internal/application/consumir.go` from US1
4. Nobody implements `internal/adapter/in/queue/redisStream.go` in this list

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- Não inventar stream, grupo, campo, lugar do código, vencedor do empate nem prazo numérico
- `FINDER_WORKERS` não é o número de workers deste consumo
