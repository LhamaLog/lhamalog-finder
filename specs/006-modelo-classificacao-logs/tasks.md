---
description: "Task list for Modelo de classificação de logs"
---

# Tasks: Modelo de classificação de logs

**Input**: Design documents from `/specs/006-modelo-classificacao-logs/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Incluídos. A spec define o teste independente de cada história e `quickstart.md` separa `make test` sem MongoDB da prova com `FINDER_MONGO_INTEGRATION=1`. Os testes de cada história vêm antes da implementação e devem falhar antes dela.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Configuração: `internal/config/config.go`
- Caso de uso: `internal/application/classificar.go`
- Porta de entrada: `internal/application/port/in/classificar.go`
- Porta de ocorrência: `internal/application/port/out/ocorrencia.go`
- Gravação: `internal/adapter/out/mongo/`
- Domínio: `internal/domain/classificacao.go`
- Testes no mesmo pacote dos arquivos que exercitam
- Ambiente: `.env.example`
- Fora desta lista: `internal/adapter/in/queue/redisStream.go`, `internal/adapter/in/http/health.go`, `internal/adapter/out/elastic/`, `internal/adapter/out/openobserve/`. Esta feature não lê a fila e não confirma a entrada

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Acrescentar o driver MongoDB nomeado no README e os nomes de configuração, sem valor padrão

- [ ] T001 Acrescentar `go.mongodb.org/mongo-driver/v2` em `go.mod` e `go.sum`. Não acrescentar biblioteca de LLM, de expressão regular nem de banco vetorial. `github.com/redis/go-redis/v9 v9.23.0` permanece e não grava a ocorrência
- [ ] T002 [P] Em `.env.example`, acrescentar os nomes `FINDER_MONGO_URI`, `FINDER_MONGO_BANCO` e `FINDER_MONGO_COLECAO` sem valor. Não gravar senha. Não definir `localhost:27017`, banco nem coleção como padrão

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Configuração sem padrão, tipos do apontamento e as portas, usados por todas as histórias

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T003 Em `internal/config/config.go`, ler `FINDER_MONGO_URI`, `FINDER_MONGO_BANCO` e `FINDER_MONGO_COLECAO`. `Carregar` não preenche os três quando a variável está ausente ou vazia. Não registrar a URI
- [ ] T004 Em `internal/config/config_test.go`, provar que configuração vazia não cria URI, banco nem coleção. Não imprimir a URI no resultado do teste
- [ ] T005 [P] Criar em `internal/domain/classificacao.go` os tipos de `specs/006-modelo-classificacao-logs/data-model.md`: apontamento, log, catálogo, classe e resultado. A classe só pode ser informação, alerta ou erro. As regras de palavra inteira, gravidade e catálogo ficam para as histórias
- [ ] T006 [P] Criar em `internal/application/port/out/ocorrencia.go` a porta que insere um documento e consulta pela identidade do apontamento. A porta não lê a fila, não confirma a entrada e não chama o OpenObserve
- [ ] T007 [P] Criar em `internal/application/port/in/classificar.go` a porta `Classificar apontamento` de `specs/006-modelo-classificacao-logs/contracts/classificar-apontamento.md`. A entrada traz identidade, tenant, produto, a marca de armazenamento inacessível e a lista de logs
- [ ] T008 [P] Em `internal/adapter/out/mongo/integracao_test.go`, a prova só executa com `FINDER_MONGO_INTEGRATION=1`. Se a variável está definida e `FINDER_MONGO_URI` não responde, o teste falha. Se não está definida, o teste é ignorado. URI, banco e coleção vêm do ambiente do teste. Não copiar esses valores para o código de produção. Não usar log de produção

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Classificar o log capturado sem serviço de linguagem externo (Priority: P1) 🎯 MVP

**Goal**: Um log presente recebe informação, alerta ou erro a partir do campo de nível e das palavras inteiras da mensagem. O texto não sai para um serviço externo. Sem catálogo, a ocorrência de prova guarda essa classe e informa que só as regras de nível desta feature foram usadas.

**Independent Test**: Entregar um log já capturado com campo de nível de erro e conferir que a classe é erro, que nenhum texto é enviado a um serviço externo e que o resultado informa as regras de nível.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T009 [P] [US1] Escrever em `internal/domain/classificacao_test.go` a lista de `specs/006-modelo-classificacao-logs/data-model.md`: erro, error e fatal são erro; alerta, warn e warning são alerta; informação, info, debug e trace são informação. O campo é o valor inteiro, sem distinção de maiúsculas e com espaços das bordas removidos. " Error " é erro. Campo informação com a palavra inteira error na mensagem é erro. Campo fatal com a palavra inteira info é erro. errorless com campo informação permanece informação. traceback sem campo e sem catálogo é informação. A palavra conta quando o caractere ao lado não é letra nem dígito. "error", "error," e "(error)" contam. Não há chamada de rede
- [ ] T010 [P] [US1] Escrever em `internal/application/classificar_test.go` um apontamento com um log presente, campo error e mensagem fatal. A ocorrência inserida na porta falsa tem classe erro, tenant, produto e identidade. `Catálogo entrou` é falso. A identificação informa que só as regras de nível desta feature foram usadas. O documento não contém a mensagem nem o campo de nível. A contagem de envios do texto a um serviço externo é zero

### Implementation for User Story 1

- [ ] T011 [US1] Em `internal/domain/classificacao.go`, calcular a classe do log como a mais grave entre o campo e as palavras inteiras da mensagem. Erro vence alerta. Alerta vence informação. Sem campo na lista e sem palavra inteira, a classe é informação. Não consultar catálogo nesta tarefa
- [ ] T012 [US1] Em `internal/application/classificar.go`, classificar um log presente e inserir pela porta de `internal/application/port/out/ocorrencia.go`. Sem identificação de catálogo, gravar que só as regras de nível desta feature foram usadas. Não ler `internal/adapter/in/queue/redisStream.go` e não confirmar a entrada

**Checkpoint**: User Story 1 is testable. One present log with field error yields class erro and no log text is stored

---

## Phase 4: User Story 2 - Gravar uma ocorrência por apontamento, sem o texto do log (Priority: P1)

**Goal**: Um apontamento, com um ou vários logs, grava uma ocorrência no MongoDB. A classe é a mais grave entre os logs presentes. Log ausente no meio do conjunto não entra e não impede a gravação. Sem nenhum log, com armazenamento inacessível ou com falha de gravação, não há ocorrência pronta.

**Independent Test**: Três logs presentes com campo error, MongoDB disponível, consulta pela identidade devolve uma ocorrência de classe erro, sem texto. A fila não é lida e a entrada não é confirmada.

### Tests for User Story 2 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T013 [P] [US2] Estender `internal/application/classificar_test.go` com os casos de `specs/006-modelo-classificacao-logs/data-model.md`: três logs informação, alerta e erro viram uma ocorrência de classe erro; informação e alerta viram alerta; três esperados e dois presentes, informação e erro, viram erro e o ausente não entra; nenhum presente é log ausente e não insere; armazenamento inacessível não insere mesmo com lista preenchida; falha da porta é falha de gravação; URI, banco ou coleção ausente é configuração ausente e não cria padrão; tenant ou produto vazio é identidade incompleta e não insere. Esses cinco motivos ficam distinguíveis. Segunda chamada da mesma identidade insere outra e a primeira permanece
- [ ] T014 [P] [US2] Escrever em `internal/adapter/out/mongo/ocorrencia_test.go` que o documento inserido só tem classe, identificação das regras, se o catálogo entrou, identidade do apontamento, tenant e produto. Não há campo de mensagem, de campo de nível, de texto de log, de gravidade, de sugestão, de evento nem de notificação

### Implementation for User Story 2

- [ ] T015 [US2] Criar em `internal/adapter/out/mongo/ocorrencia.go` a inserção e a consulta com `go.mongodb.org/mongo-driver/v2`. Inserir não substitui documento da mesma identidade. A consulta pela identidade devolve os documentos inseridos. Não registrar mensagem nem campo de nível. Substituir o comentário de `internal/adapter/out/mongo/doc.go` que diz que o pacote está por implementar
- [ ] T016 [US2] Em `internal/application/classificar.go`, aplicar o conjunto: a classe da ocorrência é a mais grave entre os presentes; ausente é ignorado e o texto que vier com ele também; nenhum presente devolve log ausente; armazenamento inacessível devolve esse motivo e não grava; erro da porta devolve falha de gravação; tenant ou produto vazio devolve identidade incompleta. URI, banco ou coleção ausentes devolvem configuração ausente e não criam padrão. Log ausente, armazenamento inacessível, falha de gravação, configuração ausente e identidade incompleta ficam distinguíveis
- [ ] T017 [US2] Estender `internal/adapter/out/mongo/integracao_test.go` com um apontamento sintético de um log, campo error e mensagem com a palavra fatal. A consulta devolve uma ocorrência de classe erro, com tenant, produto e identificação das regras, sem mensagem e sem campo de nível. A segunda inserção da mesma identidade não apaga a primeira

**Checkpoint**: User Stories 1 and 2 work. One apontamento with several present logs stores one occurrence, and a missing log does not block it

---

## Phase 5: User Story 3 - Usar o catálogo só quando o nível não decide (Priority: P2)

**Goal**: O catálogo fornecido só classifica quando o campo e a mensagem não têm nível reconhecível. A primeira entrada válida decide. Entrada de classe fora das três é ignorada. Nenhum catálogo é criado por esta feature.

**Independent Test**: Campo erro com catálogo que diria informação continua erro. Sem nível, a primeira entrada alerta vence a seguinte. Sem catálogo e sem nível, a classe é informação e a contagem de catálogos criados é zero.

### Tests for User Story 3 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T018 [US3] Estender `internal/domain/classificacao_test.go` e `internal/application/classificar_test.go` com o catálogo de `specs/006-modelo-classificacao-logs/data-model.md`. Campo error e catálogo que diria informação: classe erro e catálogo não entra. Campo informação e a palavra inteira error, com catálogo que diria informação: classe erro e catálogo não entra. Sem nível no campo e na mensagem, primeira entrada válida alerta: classe alerta. Entrada cuja classe não seja informação, alerta ou erro é ignorada. Várias entradas que casam: vale a primeira, na ordem fornecida. Frase com espaço casa o texto inteiro, e as bordas da frase não podem estar coladas a letra ou dígito. Sem catálogo fornecido, a contagem de catálogos criados é zero e a classe sem sinal é informação. `Catálogo entrou` só é verdadeiro se a classe de algum log presente saiu de uma entrada

### Implementation for User Story 3

- [ ] T019 [US3] Em `internal/domain/classificacao.go`, consultar o catálogo só quando o campo não está na lista e a mensagem não tem palavra inteira da lista. Ignorar classe fora de informação, alerta e erro. A primeira palavra inteira, na ordem fornecida, decide. Sem acerto, a classe é informação. Não criar entrada
- [ ] T020 [US3] Em `internal/application/classificar.go`, gravar a identificação fornecida com o catálogo. Sem identificação, a ocorrência informa que só as regras de nível desta feature foram usadas. `Catálogo entrou` segue a regra de T018

**Checkpoint**: All user stories work. The catalog never overrides a recognized level, and the occurrence records whether it was used

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Confirmar que a fila, a saúde e os outros adaptadores ficam fora, e que a prova do quickstart cabe no que foi gerado

- [ ] T021 Confirmar que `internal/adapter/in/queue/redisStream.go`, `internal/adapter/in/http/health.go`, `internal/adapter/out/elastic/` e `internal/adapter/out/openobserve/` não passam a chamar `internal/application/classificar.go`. A contagem de confirmações da entrada Redis por esta feature é zero. `GET /health` segue `{"status":"UP"}`
- [ ] T022 [P] Conferir `specs/006-modelo-classificacao-logs/quickstart.md` contra o código: `make test` cobre a prova sem MongoDB, e a prova com MongoDB continua atrás de `FINDER_MONGO_INTEGRATION=1`. Não afirmar resultado de comando que não tenha sido executado

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
- **Polish (Final Phase)**: Depends on the desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Starts after User Story 1 - Uses the log class from `internal/domain/classificacao.go` and adds the apontamento, the motives and the MongoDB insert
- **User Story 3 (P2)**: Starts after User Story 2 - The catalog changes the same domain rule and the same occurrence flag

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Domain rule before the use case
- Use case before the MongoDB integration proof
- Story complete before moving to the next

### Parallel Opportunities

- T002 can run in parallel with T001
- T005, T006, T007 and T008 can run in parallel after T003 starts, because they touch other files
- T009 and T010 can run in parallel
- T013 and T014 can run in parallel
- User stories do not run in parallel: US2 edits the files US1 creates, and US3 edits those same files again

---

## Parallel Example: User Story 1

```text
Task: "T009 domain level tests in internal/domain/classificacao_test.go"
Task: "T010 one-log application test in internal/application/classificar_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Demo one log classified as erro, alerta or informação, with no log text stored

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Demo (MVP)
3. Add User Story 2 → Test independently → Demo the one occurrence per apontamento
4. Add User Story 3 → Test independently → Demo the catalog only when the level does not decide
5. Each story adds value without breaking the previous class rules

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. After Foundational, one person implements User Story 1
3. The next person takes User Story 2 only after the class function of User Story 1 exists
4. User Story 3 follows User Story 2

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate the story independently
- Avoid: vague tasks, same file conflicts, a catalog entry created inside the module
