---
description: "Task list for Corrigir entrega da fila"
---

# Tasks: Corrigir entrega da fila

**Input**: Design documents from `/specs/002-corrigir-entrega-fila/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Incluídos. `research.md` (decisão 9) e `quickstart.md` pedem `make test` sem Redis e prova com Redis real quando `FINDER_REDIS_INTEGRATION=1`. Os testes de cada história vêm antes da implementação e devem falhar antes dela.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Processo: `cmd/finder/main.go`
- Caso de uso: `internal/application/consumir.go`
- Fila: `internal/adapter/in/queue/redisStream.go`
- Saúde: `internal/adapter/in/http/health.go` permanece com o mesmo corpo
- Domínio: `internal/domain/item.go`
- Testes no mesmo pacote, como `internal/config/config_test.go`
- Fora desta lista: `internal/adapter/out/mongo/`, `internal/adapter/out/elastic/`, `internal/adapter/out/openobserve/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Manter o módulo Go já existente, sem biblioteca nova

- [X] T001 Não adicionar módulo em `go.mod`. A dependência permanece `github.com/redis/go-redis/v9 v9.23.0`. Não executar `go get`. Não criar pacote fora de `cmd/finder/`, `internal/application/`, `internal/adapter/in/queue/`, `internal/adapter/in/http/` e `internal/domain/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Porta de prova e a regra de estado que todas as histórias usam

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T002 [P] Criar em `internal/adapter/in/queue/redis_integration_test.go` a guarda da prova com Redis: executa só com `FINDER_REDIS_INTEGRATION=1`; se a variável está definida e o endereço de `FINDER_REDIS_ADDR` não responde, o teste falha; se a variável não está definida, o teste é ignorado. Stream, grupo e conjunto vêm de `FINDER_REDIS_STREAM`, `FINDER_REDIS_GROUP` e `FINDER_REDIS_CODIGOS`. Não gravar esses nomes como padrão dentro do código de produção
- [X] T003 [P] Em `internal/domain/item.go`, manter os estados literais `livre`, `em_passo`, `concluido` e `nao_concluido`. Não criar um quinto estado. A condição parcial de `specs/002-corrigir-entrega-fila/data-model.md` continua sendo só a combinação "código no conjunto e item ainda pendente". A identidade do item permanece distinta do código. O conteúdo aceito continua sendo somente o código de apontamento

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Conclusão interrompida termina num único desfecho (Priority: P1) 🎯 MVP

**Goal**: Registro do código e saída da pendência acontecem juntos. Estado parcial termina numa conclusão só. Falha nessa tentativa não encerra o processo.

**Independent Test**: Interromper depois do primeiro efeito e retomar: uma conclusão, item fora da pendência, código disponível uma vez. Falhar ao completar: o processo segue, o item permanece parcial e a tentativa seguinte chega ao mesmo desfecho.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T004 [P] [US1] Escrever em `internal/application/consumir_test.go` os casos de `specs/002-corrigir-entrega-fila/contracts/entrega-fila.md`: código já registrado e item ainda pendente, na tentativa seguinte, produz uma conclusão, item fora da pendência e código disponível uma vez; falha ao completar deixa o item parcial, não aumenta a contagem e não exige repetição até o sucesso na mesma chamada; segundo item com o mesmo código não gera segunda conclusão e não permanece reservado
- [X] T005 [P] [US1] Escrever em `internal/adapter/in/queue/redis_integration_test.go` o caso em que o código já está no conjunto e a mensagem segue pendente: depois da operação de conclusão há um membro no conjunto e a mensagem não está pendente. Uma falha dessa operação deixa a mensagem pendente para a chamada seguinte

### Implementation for User Story 1

- [X] T006 [US1] Em `internal/adapter/in/queue/redisStream.go`, gravar o código no conjunto e tirar a mensagem da pendência numa única operação. Se o código já é membro, só tirar da pendência e não aumentar a contagem. Não publicar a mensagem de novo
- [X] T007 [US1] Em `internal/application/consumir.go`, o caminho aceito usa essa operação única. Código já presente e item ainda pendente completa a saída da pendência; não encerra como repetido deixando o item pendente. Erro nessa operação devolve o erro e conserva a condição parcial. Não abrir o OpenObserve
- [X] T008 [US1] Em `cmd/finder/main.go`, falha de `Executar` ao completar a conclusão não chama `os.Exit`. O laço segue para a tentativa seguinte. Não acrescentar espera numérica

**Checkpoint**: User Story 1 is testable. A partial completion finishes once, and a failed attempt leaves the process running

---

## Phase 4: User Story 2 - Item deixado por instância parada volta ao fluxo (Priority: P2)

**Goal**: Item reservado por instância que parou antes de qualquer efeito volta à mesma identidade, à frente dos itens que já estavam atrás, sem prazo numérico.

**Independent Test**: A primeira instância lê um item e para sem efeito. A segunda reserva esse item antes de um item que já estava atrás. Instância ainda conectada não perde o item.

### Tests for User Story 2 ⚠️

- [X] T009 [US2] Escrever em `internal/adapter/in/queue/redis_integration_test.go` o caso de `specs/002-corrigir-entrega-fila/data-model.md`: a identidade do item não muda quando outra instância o reserva; o item da frente, deixado sem efeito, é reservado antes do item que já estava atrás; consumidor ainda com conexão viva não é espoliado. A prova não configura tempo de ociosidade como regra de liberação

### Implementation for User Story 2

- [X] T010 [US2] Em `internal/adapter/in/queue/redisStream.go`, publicar o nome do consumidor na conexão. Antes de ler mensagem nova, ler a pendência deste consumidor e, depois, reservar a pendência mais antiga de consumidor sem conexão viva, no mesmo identificador. O argumento de ociosidade exigido pela operação fica em zero e não vira regra do Finder. Não reservar consumidor que ainda tem conexão. Não publicar o item no fim da fila

**Checkpoint**: User Stories 1 and 2 work. A dead instance's item returns at its original place

---

## Phase 5: User Story 3 - Consumir a fila já ordenada (Priority: P3)

**Goal**: O Finder reserva na ordem já recebida. Não escolhe, não grava e não ordena nível.

**Independent Test**: Com uma instância, a sequência de reservas é a ordem de entrada. Com duas, cada uma pode reservar um item diferente ao mesmo tempo. Nenhuma escolhe outro item por causa do nome do nível.

### Tests for User Story 3 ⚠️

- [X] T011 [US3] Escrever em `internal/adapter/in/queue/redis_integration_test.go` três mensagens cujo conteúdo é só o código, acrescentadas numa ordem conhecida: com uma instância, a sequência de reservas é essa ordem. Com duas instâncias, duas mensagens novas diferentes podem estar reservadas ao mesmo tempo, e a ordem de término não reprova o caso. O teste não envia nível e não espera reordenação

### Implementation for User Story 3

- [X] T012 [US3] Em `internal/adapter/in/queue/redisStream.go` e `internal/domain/item.go`, não acrescentar campo de nível e não ordenar a fila. `itemDaMensagem` continua aceitando só um código. Os nomes Very High, High, Medium, Low e Very Low não são lidos nem gravados. A ordem observada é a ordem em que as mensagens já estão na fila

**Checkpoint**: User Story 3 preserves the order already in the queue and adds no priority field

---

## Phase 6: User Story 4 - O mesmo código conta uma vez no produto (Priority: P4)

**Goal**: O mesmo código, mesmo vindo de duas organizações, tem uma conclusão. O item não ganha organização.

**Independent Test**: Dois itens aceitos com o mesmo código produzem um membro no conjunto, o segundo sai da reserva e o membro é só o código.

### Tests for User Story 4 ⚠️

- [X] T013 [US4] Escrever em `internal/adapter/in/queue/redis_integration_test.go` dois itens aceitos com o mesmo código: o conjunto fica com um membro, esse membro é só o texto do código, o segundo item sai da pendência e não permanece reservado. Não há sufixo de organização

### Implementation for User Story 4

- [X] T014 [US4] Em `internal/adapter/in/queue/redisStream.go`, o membro do conjunto é somente o código. Não ler nem gravar organização. Em `internal/domain/item.go`, não acrescentar campo de organização

**Checkpoint**: User Story 4 keeps one completion per code in the product

---

## Phase 7: User Story 5 - Item fora do contrato não fica preso (Priority: P5)

**Goal**: Item sem código, ou com conteúdo além do código, é descartado, sai da pendência e da reserva, e não entra no conjunto. O descarte se distingue de uma conclusão.

**Independent Test**: Um item sem código e outro com dois valores, seguidos de um item aceito. Os dois primeiros não ficam pendentes nem no mapa de reserva, não contam como conclusão, e o aceito conclui.

### Tests for User Story 5 ⚠️

- [X] T015 [P] [US5] Escrever em `internal/application/consumir_test.go` o descarte de item sem código e de item com conteúdo extra: sem conclusão, sem código no conjunto, fora da reserva, com registro observável diferente de item concluído, e o item aceito seguinte pode concluir
- [X] T016 [P] [US5] Escrever em `internal/adapter/in/queue/redis_integration_test.go` uma mensagem sem valores e outra com dois valores: ambas saem da pendência, não entram no conjunto e não permanecem em `porID`. A mensagem aceita seguinte conclui

### Implementation for User Story 5

- [X] T017 [US5] Em `internal/application/consumir.go`, quando `!item.Aceito()`, chamar `EncerrarSemConclusao` e registrar descarte distinguível de `item concluído`. Em `internal/adapter/in/queue/redisStream.go`, esse caminho apaga o item de `porID` e não cria fila de erro. Falha nessa saída deixa o item pendente para a tentativa seguinte, com o processo no ar

**Checkpoint**: Invalid items no longer stay pending or reserved

---

## Phase 8: User Story 6 - Processo em execução atende a consulta de saúde (Priority: P6)

**Goal**: Com o processo no ar, `GET /health` responde como o handler já existente. Nenhuma consulta nova.

**Independent Test**: Processo em execução atende a consulta. Processo parado não atende. O corpo permanece `{"status":"UP"}`.

### Tests for User Story 6 ⚠️

- [X] T018 [P] [US6] Manter em `internal/adapter/in/http/health_test.go` a prova de `GET /health` com status 200 e corpo `{"status":"UP"}`. Não alterar esse corpo em `internal/adapter/in/http/health.go`. Acrescentar em `cmd/finder/saude_test.go` a prova de que o servidor montado com `NovoRoteador` atende essa consulta no endereço de `FINDER_HTTP_ADDR`

### Implementation for User Story 6

- [X] T019 [US6] Em `cmd/finder/main.go`, servir o roteador de `internal/adapter/in/http` no endereço de `internal/config/config.go` (`FINDER_HTTP_ADDR`, padrão `:8085`). Encerrar o servidor quando o contexto do processo for cancelado. Não criar outra consulta

**Checkpoint**: The running process answers the existing health check

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: Textos operacionais e a prova do quickstart

- [X] T020 [P] Alinhar `README.md` ao processo de `cmd/finder/main.go`: a consulta de saúde em `:8085` e o registro de arranque que o processo realmente emite
- [X] T021 [P] Reescrever `.env.example` como linhas `KEY=VALUE` copiáveis para as chaves lidas em `internal/config/config.go`, sem descrição colada depois de `:`
- [X] T022 Executar os comandos de `specs/002-corrigir-entrega-fila/quickstart.md` e registrar nesse arquivo o comando e o resultado observado. `make test` passa sem Redis. A prova com `FINDER_REDIS_INTEGRATION=1` não é ignorada

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: começa imediatamente
- **Foundational (Phase 2)**: depende do Setup e bloqueia todas as histórias
- **User Stories (Phase 3+)**: dependem do Foundational
- **Polish (Phase 9)**: depende das histórias desejadas

### User Story Dependencies

- **User Story 1 (P1)**: começa depois do Foundational. Não depende das outras histórias
- **User Story 2 (P2)**: começa depois do Foundational. A implementação em `internal/adapter/in/queue/redisStream.go` vem depois de T006, porque é o mesmo arquivo
- **User Story 3 (P3)**: a prova e a ausência de nível podem ser escritas depois do Foundational. T012 vem depois de T010, no mesmo arquivo da fila
- **User Story 4 (P4)**: T014 vem depois de T006, no mesmo script de conclusão
- **User Story 5 (P5)**: T017 vem depois de T007, no mesmo caso de uso
- **User Story 6 (P6)**: T018 pode começar depois do Foundational. T019 vem depois de T008, porque os dois editam `cmd/finder/main.go`

### Within Each User Story

- Testes antes da implementação, e falhando antes dela
- A história fica demonstrável no checkpoint, sem precisar das histórias seguintes

### Parallel Opportunities

- T002 e T003 podem correr juntos
- T004 e T005 podem correr juntos
- T015 e T016 podem correr juntos
- T018 pode correr junto com os testes de outras histórias que não editem `cmd/finder/`
- T020 e T021 podem correr juntos
- Histórias em arquivos diferentes podem avançar em paralelo depois das dependências do mesmo arquivo

---

## Parallel Example: User Story 1

```text
T004: testes de conclusão parcial e falha em internal/application/consumir_test.go
T005: prova Redis do conjunto e da pendência em internal/adapter/in/queue/redis_integration_test.go
```

Depois, em ordem: T006, T007 e T008.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Completar Phase 1 e Phase 2
2. Completar Phase 3
3. Parar e validar a conclusão única e a falha que não encerra o processo
4. Só então seguir para a pendência da instância parada

### Incremental Delivery

1. Setup + Foundational
2. User Story 1: desfecho único
3. User Story 2: item da instância parada volta ao lugar original
4. User Story 3: ordem já recebida, sem nível no Finder
5. User Story 4: um código, uma conclusão no produto
6. User Story 5: item inválido não fica preso
7. User Story 6: consulta de saúde
8. Polish: README, `.env.example` e quickstart

### Parallel Team Strategy

1. Todos fecham Setup e Foundational
2. Em seguida, quem mexe em `internal/adapter/in/queue/redisStream.go` segue T006, T010, T012 e T014 nessa ordem
3. Quem mexe em `internal/application/consumir.go` segue T007 e depois T017
4. Quem mexe em `cmd/finder/main.go` segue T008 e depois T019

---

## Notes

- [P] só vale para arquivos diferentes e sem dependência incompleta
- Não criar prazo numérico, fila de erro, campo de organização nem campo de nível
- Não alterar o corpo de `GET /health`
- Cada checkpoint deve poder ser demonstrado sem as histórias posteriores
