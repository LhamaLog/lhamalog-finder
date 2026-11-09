---
description: "Task list for Recuperar logs do lote"
---

# Tasks: Recuperar logs do lote

**Input**: Design documents from `/specs/005-recuperar-logs-lote/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Incluídos. A spec pede provas unitárias e integradas com OpenObserve de teste e dados sintéticos. `quickstart.md` separa `make test` sem OpenObserve da prova com `FINDER_OPENOBSERVE_INTEGRATION=1`. Os testes de cada história vêm antes da implementação e devem falhar antes dela.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Configuração: `internal/config/config.go`
- Caso de uso: `internal/application/recuperar.go`
- Porta de busca: `internal/application/port/out/busca.go`
- Consulta e cliente: `internal/adapter/out/openobserve/`
- Domínio: `internal/domain/recuperacao.go`
- Testes no mesmo pacote dos arquivos que exercitam
- Ambiente: `.env.example`
- Fora desta lista: `internal/adapter/in/queue/redisStream.go`, `internal/adapter/in/http/health.go`, `internal/adapter/out/mongo/`, `internal/adapter/out/elastic/`. Esta feature não confirma a entrada Redis

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Manter o módulo Go já existente e nomear a configuração do OpenObserve sem valor padrão

- [ ] T001 Confirmar em `go.mod` que não entra biblioteca nova. A consulta usa `net/http`. Não executar `go get`. `github.com/redis/go-redis/v9 v9.23.0` permanece e não faz esta consulta
- [ ] T002 [P] Em `.env.example`, acrescentar os nomes `FINDER_OO_URL`, `FINDER_OO_USUARIO`, `FINDER_OO_SENHA`, `FINDER_OO_TIMEOUT`, `FINDER_OO_PAGINA`, `FINDER_OO_CAMPOS_COPIA`, `FINDER_OO_ORDEM`, `FINDER_OO_CAMPO_TENANT` e `FINDER_OO_CAMPO_PRODUTO` sem valor. Não gravar senha. Não definir `localhost:5080`, timeout, tamanho de página nem nome de campo como padrão

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Configuração sem padrão, tipos do resultado e a porta de busca, usados por todas as histórias

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T003 Em `internal/config/config.go`, ler os nove nomes de T002. `Carregar` não preenche URL, usuário, senha, timeout, tamanho de página, lista de cópia, expressão de ordenação nem nomes de tenant e produto. Lista ausente e lista presente vazia ficam distinguíveis. Senha vazia permanece vazia. Expor a leitura sem registrar a senha
- [ ] T004 Em `internal/config/config_test.go`, provar que configuração vazia não cria URL, usuário, senha, timeout, tamanho de página, lista, ordenação nem nomes de tenant e produto. Provar que `FINDER_OO_CAMPOS_COPIA` definido e vazio é lista presente. Não imprimir a senha no resultado do teste
- [ ] T005 [P] Criar em `internal/domain/recuperacao.go` os tipos da referência de entrada, da configuração de acesso, do documento físico e do resultado, como em `specs/005-recuperar-logs-lote/data-model.md`. O resultado carrega uma classificação, a marca de efeito concluído, o contexto do lote e, quando a contagem foi calculada, a quantidade esperada e a encontrada. As regras de deduplicação, completude e precedência ficam para as histórias
- [ ] T006 [P] Criar em `internal/application/port/out/busca.go` a porta que consulta uma página e devolve `hits` ou informa credencial recusada, resposta inválida, timeout ou indisponibilidade. A porta não confirma a entrada Redis, não oculta campo e não escolhe versão divergente
- [ ] T007 [P] Em `internal/adapter/out/openobserve/integracao_test.go`, a prova só executa com `FINDER_OPENOBSERVE_INTEGRATION=1`. Se a variável está definida e `FINDER_OO_URL` não responde, o teste falha. Se não está definida, o teste é ignorado. URL, credencial, timeout, página, lista, ordenação e nomes de campo vêm do ambiente do teste. Não copiar esses valores para o código de produção. Não usar log de produção

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Recuperar os logs do lote pela referência (Priority: P1) 🎯 MVP

**Goal**: Uma referência válida consulta o stream do OpenObserve e, numa página curta cuja contagem fecha, devolve os registros e a correlação. A entrada Redis não é confirmada.

**Independent Test**: Referência sintética de um registro e de vários, contra um substituto que já contém só os registros aceitos. A consulta usa organização, stream, lote e faixa. O lote unitário também filtra o identificador. O conteúdo não sai da entrada Redis e a confirmação desta etapa é zero.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T008 [P] [US1] Escrever em `internal/adapter/out/openobserve/consulta_test.go` a consulta de `specs/005-recuperar-logs-lote/contracts/recuperar-lote.md`: `POST /api/{organization}/_search` com `query.sql`, `query.start_time`, `query.end_time`, `query.from` e `query.size`. `start_time` é `minTimestampMicros`. `end_time` é `maxTimestampMicros + 1`. O SQL usa o stream da referência, não `lhamalog:stream:logs`, e filtra `lhamalog_batch_id`. Sem `recordId`, não filtra `lhamalog_record_id`. Com `recordId`, filtra. `from` da primeira página é 0. Nome de stream e de campo entre aspas duplas. Valor de texto entre aspas simples, com aspas simples interna repetida. Inteiro sem aspas. Valor que não possa ser aplicado assim produz resposta inválida e não troca lote, organização, stream nem faixa
- [ ] T009 [P] [US1] Escrever em `internal/application/recuperar_test.go` o sucesso de uma página curta cuja quantidade distinta é igual a `recordCount`. O resultado traz os registros, `tenantId`, `productId`, organização, stream, `batchId`, `idempotencyId`, `recordCount` e a quantidade distinta. A marca é efeito concluído. A contagem de confirmações Redis é zero

### Implementation for User Story 1

- [ ] T010 [P] [US1] Criar em `internal/adapter/out/openobserve/consulta.go` a montagem do corpo de T008, incluindo `ORDER BY` com `FINDER_OO_ORDEM`. Sem a expressão, não montar a consulta. Não gravar `_timestamp` nem `lhamalog_record_id` como ordem do código
- [ ] T011 [P] [US1] Criar em `internal/adapter/out/openobserve/cliente.go` o POST com `net/http`, usuário e senha em HTTP Basic, e o timeout de cada consulta no cliente. Não enviar o campo `timeout` do corpo. Não registrar senha nem `hits` em falha. Substituir o comentário de `internal/adapter/out/openobserve/doc.go` que diz que o pacote não consulta
- [ ] T012 [US1] Criar em `internal/application/recuperar.go` o caso de uso de uma página curta. URL, credencial, timeout, tamanho de página, lista de cópia ou ordenação ausentes: não consultar, não criar padrão e não marcar efeito concluído. Página curta com contagem distinta igual a `recordCount` e sem outra falha: efeito concluído e registros. Não confirmar a entrada Redis
- [ ] T013 [US1] Estender `internal/adapter/out/openobserve/integracao_test.go` com um lote sintético de uma página, de um registro e de vários. Conferir organização, stream, lote, faixa com fim no máximo mais 1 e o filtro de `recordId` só no lote unitário. A entrada Redis não é confirmada

**Checkpoint**: User Story 1 is testable. A single short page whose distinct count matches returns the records and leaves the Redis entry unconfirmed

---

## Phase 4: User Story 2 - Percorrer todas as páginas e deduplicar cópias físicas (Priority: P1)

**Goal**: A leitura avança por todas as páginas cheias e cada `lhamalog_record_id` válido aparece uma vez. Não há máximo de páginas.

**Independent Test**: Uma página, várias páginas e quantidade múltipla exata do tamanho de página. Cópias do mesmo registro na mesma página e entre páginas, com um identificador novo depois delas. Sem a lista de cópia, a recuperação não começa.

### Tests for User Story 2 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T014 [P] [US2] Estender `internal/application/recuperar_test.go` e `internal/domain/recuperacao_test.go` com os casos de `specs/005-recuperar-logs-lote/data-model.md`: `from` começa em 0 e a página seguinte soma `size`; página com menos documentos que `size` encerra; página cheia pede a seguinte, inclusive quando a quantidade é múltiplo exato de `size`; não existe máximo de páginas; `size` não é `recordCount`; o mesmo `lhamalog_record_id` com diferença só na lista de cópia vira um registro e o entregue é uma das cópias; diferença fora da lista não escolhe versão; documento sem `lhamalog_record_id` fica fora da contagem; identificador Redis e `correlationId` não deduplicam; lista ausente não consulta e não marca efeito concluído; lista presente vazia torna qualquer diferença conteúdo divergente

### Implementation for User Story 2

- [ ] T015 [US2] Em `internal/domain/recuperacao.go`, deduplicar por `lhamalog_record_id` com a lista de cópia. Diferença restrita à lista: um registro lógico, e o entregue é uma das cópias. Diferença fora da lista: nenhuma versão escolhida. Documento sem identificador: fora da contagem
- [ ] T016 [US2] Em `internal/application/recuperar.go`, repetir a porta de busca enquanto `hits` tiver o comprimento de `size`. Somar `size` a `from`. Não criar teto de páginas. Não parar em `recordCount`. Lista de cópia ausente não consulta
- [ ] T017 [US2] Estender `internal/adapter/out/openobserve/integracao_test.go` com cópias físicas e um identificador lógico que só aparece depois delas, em mais de uma página. Cada identificador aparece uma vez. O identificador que ficaria de fora se a leitura parasse em `recordCount` está no conjunto

**Checkpoint**: User Stories 1 and 2 work. Full pages are followed and physical copies collapse to one logical record

---

## Phase 5: User Story 3 - Declarar o lote completo só quando a contagem fecha (Priority: P1)

**Goal**: Igualdade com `recordCount` declara o lote completo só sem outra falha. Menor é incompleto, maior é excesso, e busca vazia com `recordCount` positivo não é sucesso.

**Independent Test**: Contagem igual, menor, maior e zero, mais um lote misto cuja referência conta só os aceitos. A causa da divergência não é afirmada.

### Tests for User Story 3 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T018 [P] [US3] Estender `internal/domain/recuperacao_test.go` com a comparação de `specs/005-recuperar-logs-lote/data-model.md`: quantidade distinta igual a `recordCount` e sem falha é lote completo com efeito concluído; menor é lote incompleto, sem efeito concluído, com as duas quantidades; maior é excesso, com efeito concluído e as duas quantidades; busca sem documentos, `recordCount` positivo e sem falha de acesso é lote incompleto com encontrado 0; a contagem não afirma atraso de visibilidade nem falha de persistência; lote misto completa pela quantidade aceita e não exige o recusado; não há segunda consulta `COUNT(DISTINCT ...)`

### Implementation for User Story 3

- [ ] T019 [US3] Em `internal/domain/recuperacao.go`, fechar o lote pela quantidade de identificadores distintos válidos contra `recordCount`. Igualdade só completa se isolamento, presença de identificador e ausência de conteúdo divergente tiverem passado. Informar esperado e encontrado quando a contagem foi calculada. Não emitir a causa da divergência
- [ ] T020 [US3] Estender `internal/adapter/out/openobserve/integracao_test.go` com busca vazia, contagem menor, contagem maior e lote misto sintético. A completude do misto usa só os registros aceitos

**Checkpoint**: User Stories 1, 2, and 3 work. A batch is complete only when the distinct count matches and no other failure is present

---

## Phase 6: User Story 4 - Manter o lote no tenant da referência (Priority: P1)

**Goal**: Com os nomes dos campos, a consulta filtra tenant e produto e cada documento é conferido. Sem os nomes, a consulta não sai e o resultado é isolamento com efeito concluído.

**Independent Test**: Dois tenants sintéticos. Documento do outro tenant ou produto não é resultado válido. Sem os nomes, a contagem de consultas é zero e o efeito está concluído.

### Tests for User Story 4 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T021 [P] [US4] Estender `internal/application/recuperar_test.go` e `internal/adapter/out/openobserve/consulta_test.go` com as regras de `specs/005-recuperar-logs-lote/data-model.md`: os dois nomes presentes filtram pelos valores da referência e cada documento é conferido; falta de um nome não envia a consulta, não entrega documento e classifica isolamento com efeito concluído; documento de outro tenant ou produto impede qualquer documento do lote de sair como válido; credencial sem filtro e sem conferência não entrega sucesso; atributo de dentro do log não substitui `tenantId` nem `productId`; conhecer organização, stream e `batchId` não autoriza

### Implementation for User Story 4

- [ ] T022 [US4] Em `internal/adapter/out/openobserve/consulta.go` e `internal/application/recuperar.go`, filtrar e conferir tenant e produto quando `FINDER_OO_CAMPO_TENANT` e `FINDER_OO_CAMPO_PRODUTO` estão os dois presentes. Se um falta, não consultar e devolver isolamento com efeito concluído. A credencial não substitui os dois controles
- [ ] T023 [US4] Estender `internal/adapter/out/openobserve/integracao_test.go` com dois tenants no mesmo `batchId`. Nenhum documento do outro tenant ou produto sai como resultado válido. A classificação é isolamento

**Checkpoint**: User Stories 1 through 4 work. A document from another tenant is not returned as a valid batch

---

## Phase 7: User Story 5 - Devolver falha classificada sem sucesso parcial (Priority: P2)

**Goal**: Falha chega classificada, sem conteúdo de log e sem credencial. Uma só classificação vale quando há mais de uma falha. O parcial não é lote pronto.

**Independent Test**: Timeout, credencial recusada, resposta inválida, falha na segunda página, lote incompleto, excesso, conteúdo divergente e documento sem identificador. Isolamento junto com documento sem identificador sai como isolamento. Só lote completo e falha que pede correção marcam efeito concluído.

### Tests for User Story 5 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T024 [P] [US5] Estender `internal/domain/recuperacao_test.go` e `internal/application/recuperar_test.go` com a precedência de `specs/005-recuperar-logs-lote/data-model.md`, uma só classificação, nesta ordem: falha de autenticação ou autorização; resposta inválida; inconsistência de identidade ou isolamento; documento sem identificador; conteúdo divergente; excesso de registros distintos; lote incompleto. Indisponibilidade, timeout e falha de página só saem quando nenhuma dessas sete ocorreu. Configuração ausente não é essas classes: não consulta e não marca efeito concluído. Página que falha depois de falha que pede correção emite essa falha, com efeito concluído, e o parcial não é lote pronto. Credencial recusada não é falha temporária. O resultado de falha não contém senha nem conteúdo de log

### Implementation for User Story 5

- [ ] T025 [US5] Em `internal/domain/recuperacao.go`, emitir uma classificação pela ordem de T024 e a marca de efeito concluído da tabela de `specs/005-recuperar-logs-lote/data-model.md`. Lote completo e falha que pede correção marcam efeito concluído. Indisponibilidade, timeout, falha de página e lote incompleto, sozinhos, não marcam
- [ ] T026 [US5] Em `internal/adapter/out/openobserve/cliente.go` e `internal/application/recuperar.go`, mapear credencial recusada, corpo que não pode ser interpretado, timeout e armazenamento inalcançável. Autenticação recusada ou resposta inválida encerra a leitura. Falha de documento não encerra sozinha, porque uma página seguinte pode trazer resposta inválida ou isolamento. Não confirmar a entrada Redis
- [ ] T027 [US5] Estender `internal/adapter/out/openobserve/integracao_test.go` com o OpenObserve inalcançável na segunda página. Sem falha que pede correção já vista, a classificação é indisponibilidade ou timeout e o parcial não é lote pronto. Com documento de outro tenant já visto, a classificação é isolamento e o efeito está concluído

**Checkpoint**: All user stories work. Failures return one classification and do not present a partial page set as a completed batch

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Ligar o caso de uso à porta de efeito e conferir o quickstart sem alargar o escopo

- [ ] T028 Em `internal/application/port/out/efeito.go`, quando o arquivo existir, delegar a execução a `internal/application/recuperar.go` e devolver a marca de efeito concluído. Se o arquivo ainda não existir, não criar o consumo Redis. Em nenhum caso editar `internal/adapter/in/queue/redisStream.go` para confirmar a entrada
- [ ] T029 [P] Em `internal/adapter/out/openobserve/doc.go`, descrever que o pacote consulta o stream da referência e não escolhe URL, senha, timeout, página, lista de cópia, ordenação nem nomes de tenant e produto
- [ ] T030 Executar a prova de `specs/005-recuperar-logs-lote/quickstart.md`. Registrar o comando observado e a saída. Não afirmar sucesso sem essa saída. Não usar log de produção e não declarar homologação

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
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Extends the page loop of US1 and stays testable with a fake search port
- **User Story 3 (P1)**: Can start after Foundational (Phase 2) - Uses the distinct set from US2 when both are present, and stays testable with a set already deduplicated
- **User Story 4 (P1)**: Can start after Foundational (Phase 2) - Adds the tenant filter to the query of US1 and stays testable with a fake search port
- **User Story 5 (P2)**: Can start after Foundational (Phase 2) - Adds precedence over the classifications of US3 and US4 and stays testable with a fake search port

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Query assembly before the use case that calls it
- Domain rule before the use case that applies it
- Story complete before treating the next priority as done

### Parallel Opportunities

- T002 can run in parallel with T001
- T005, T006, and T007 can run in parallel after T003
- T008 and T009 can run in parallel
- T010 and T011 can run in parallel after T008
- T014, T018, T021, and T024 share test files, so they do not run at the same time
- T029 can run in parallel with T028

---

## Parallel Example: User Story 1

```text
Task: "Escrever a consulta em internal/adapter/out/openobserve/consulta_test.go"
Task: "Escrever o sucesso de uma página em internal/application/recuperar_test.go"

Task: "Montar o corpo em internal/adapter/out/openobserve/consulta.go"
Task: "Enviar o POST em internal/adapter/out/openobserve/cliente.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Demo a single short page recovered from the reference, with the Redis entry still unconfirmed

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Demo (MVP!)
3. Add User Story 2 → Test independently → Demo full pages and dedup
4. Add User Story 3 → Test independently → Demo completeness
5. Add User Story 4 → Test independently → Demo tenant isolation
6. Add User Story 5 → Test independently → Demo one failure classification
7. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2
   - Developer C: User Story 4
3. User Stories 3 and 5 follow the classifications those stories introduce

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- Não criar URL, senha, timeout, tamanho de página, lista de cópia, ordenação nem nomes de tenant e produto no código
- Não criar máximo de páginas
- Não confirmar a entrada Redis nesta feature
