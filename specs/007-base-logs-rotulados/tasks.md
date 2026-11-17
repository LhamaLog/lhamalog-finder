---
description: "Task list for Base de logs rotulados para treino e avaliação"
---

# Tasks: Base de logs rotulados para treino e avaliação

**Input**: Design documents from `/specs/007-base-logs-rotulados/`

**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md)

**Tests**: Validação estrutural em **100%** dos registros (007-21 A; SC-004) e consumo via testes automatizados (007-22 B). Escrever testes antes da implementação correspondente quando aplicável.

**Organization**: Fases Setup → Foundational → user stories (P1, P1, P2) → Polish. Cada tarefa cita origem (FR/US, plan, decisão com fonte).

**Correção analyze (etapa 3c)**: M5 — rastreios T004, T008, T056 em `plan.md`; H1/M4 — T057; M6 — T058 (fonte: analyze-007.md).

**Etapa 3e**: **DN-006-2/007-3 B** (`respostas-joao.md` §Etapa 3c, João 16/11 11h01 e 11h) e **DN-007-4 A** (`respostas-joao.md` §Etapa 3c (João, 16/11 11h02)) — desbloqueiam **T031**, **T035**, **T043** e asserts recall/F1 em **T042**.

**Etapa 3f**: Não há divergência entre §Etapa 3c (João, 16/11 11h01) e §Etapa 3c (João, 16/11 11h) — nota «puladas» **SUPERADA** (`respostas-joao.md` linha 39); **DN-007-4 A** (e **DN-006-3 B** / **DN-008-7 B**, fora deste artefato) com fonte canônica `respostas-joao.md §Etapa 3c (João, 16/11 11h02)`. Divergências de produto/Jira: apenas [spec.md](./spec.md) § Divergências documentadas (sem tabela 11h×11h01 neste arquivo).

**Etapa 3g**: Achado **H1** (`analyze-007-etapa3f.md`) — pré-condição **T059**: implementação da spec **006** (`classificar-log.md`) antes de integração **Classificar log** e gates **≥ 95%** (fonte: [plan.md](./plan.md) § Testing; spec.md FR-007, SC-005; respostas-joao.md, 007-22 B / 007-23 A; origem: spec 006).

**Etapa 3g (continuação)**: Achado **H2** (`analyze-007-etapa3f.md`) — **007 H2 D** nas tarefas de métricas/avaliação **T041**, **T042**, **T043**, **T044**, **T045**, **T053**, **T055**: para `apontamento`, predição **por** `member`; classe **prevista** = mais grave entre predições (**ERROR** > **WARNING** > **INFO**, ADR_170); comparar ao único gold agregado (**DN-007-4 A**); unidade = **apontamento**; itens `log` mantêm par gold×predição por registro (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), 007 H2 D; origem: ADR_170; respostas-etapa3f.md §D; spec.md FR-007c).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Paralelizável (arquivos distintos, sem dependência de tarefas incompletas)
- **[US#]**: História em fases de user story
- **BLOQUEADO (`<pendência>`)**: Não implementar até a pendência indicada (ex.: `006-7 D`)

## Path Conventions (plan.md § Project Structure)

| Área | Caminho alvo |
|------|----------------|
| Artefatos versionados (Git) | `datasets/labeled-logs/` — **007-caminho C**, formato **JSONL** (007-formato A); documentado em spec/plan § Artefatos físicos |
| Declaração / split / manifest | `datasets/labeled-logs/ORIGEM.md`, `datasets/labeled-logs/SPLIT.md`, `datasets/labeled-logs/MANIFEST-contagens.md` (nomes lógicos; físicos na tabela de mapeamento) |
| Loader / validação | `internal/labeleddataset/` — `doc.go`, tipos, loader, validador, métricas offline (plan.md § Source Code — `internal/…`) |
| Testes Go | `internal/labeleddataset/*_test.go`; integração com classificador 006 via contrato em `specs/007-base-logs-rotulados/contracts/entrada-classificador-desde-base.md` |
| Documentação feature | `specs/007-base-logs-rotulados/` (spec, checklist, quickstart) |

**Fora do escopo desta feature (não acoplar)**: adaptador HTTP Módulo Ocorrência (DN-009-2); fila Redis `log.ingested` (007-23 A); normalização de assinatura (DN-008-4 / spec 008); treino LLM (ADR_167; 007-3 A).

---

## Pendências (referência — plan.md § Pendências)

| ID | Impacto nesta feature | Tarefas bloqueadas |
|----|----------------------|-------------------|
| **DN-007-1 C** / **006-7 D** | Conteúdo do subconjunto só-catálogo | T034, T054 |
| **DN-009-2** | Nenhuma tarefa HTTP na 007 | — (somente referência; rastreio T008 em `plan.md`) |
| **DN-008-4** | Fora do escopo (007-17 A) | — (nota R14 em T052) |
| **010-27 A** | Ortogonal à base (mascaramento spec 010) | — (rastreio em T005, T039; não bloqueia validação 007) |
| **N3 A**, **010-2 A** | ADR fora desta pasta | T004, T005, T039 (rastreio T004 em `plan.md`) |
| **006-10 A** | Coordenação Jira spec 006 | T006, T040 (documentação) |
| **Spec 006 — operação Classificar log** (`classificar-log.md`) | Integração de testes e métricas SC-005 / FR-007 | T024, T031, T041, T042, T045 — gate **T059** |
| **emenda princípio I** | Constituição × T1 revisto | — (rastreio T056 em `plan.md`) |

### Desbloqueado na etapa 3c (fonte nas tarefas afetadas)

| Decisão | Efeito nas tarefas |
|---------|-------------------|
| **007 formato A** + **007 caminho C** (JSONL, `datasets/labeled-logs/`) | **T003**, **T007**, **T013**, **T020**, **T026**–**T027** sem bloqueio de especificação de formato |
| **DN-007-2 B** + **007-meta** (≥ 95% precisão **por classe**) | **T042**, **T044**, **T045**, **T053**, **T055** (precisão) |
| **DN-007-1 C** (critério só-catálogo) | **T033**, **T051** documentação **desbloqueada**; conteúdo **T034**/**T054** ainda bloqueado em **006-7 D** |
| **DN-008-4 B** (normalização spec 008) | Fora do escopo 007 — **T052** nota; não bloqueia base |

### Desbloqueado na etapa 3e (fonte nas tarefas afetadas)

| Decisão | Efeito nas tarefas |
|---------|-------------------|
| **DN-006-2/007-3 B** (≥ 95% recall e F1 **por classe**) | **T042**, **T043**, **T044**, **T045**, **T053**, **T055** (métricas LLT-199 critério 2 completas) |
| **DN-007-4 A** (gold **um** `class_label` por apontamento; `members` sem gold) (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-007-4 A; linha 39 SUPERADO) | **T011**, **T030**, **T031**, **T035** |

### Desbloqueado na etapa 3g (fonte nas tarefas afetadas)

| Decisão | Efeito nas tarefas |
|---------|-------------------|
| **007 H2 D** (prevista do apontamento = mais grave entre predições dos `members`; gold×prevista no agregado; unidade = apontamento) (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), 007 H2 D; origem: ADR_170; respostas-etapa3f.md §D) | **T041**, **T042**, **T043**, **T044**, **T045**, **T053**, **T055** (comparação gold × predição e gate **≥ 95%** por classe) |

Decisões **desbloqueadas** na etapa 3b (fonte: respostas-joao.md §Etapa 3; respostas-etapa3.md): **DN-007-1 C** aguardar catálogo para **conteúdo** só-catálogo; **007-meta** ≥ 95% precisão (escopo por classe fechado em 3c: **DN-007-2 B**).

Divergências conscientes: ver [spec.md](./spec.md) § Divergências documentadas.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Estrutura de diretórios, rastreio de pendências externas e alinhamento documental **sem** criar ADR nem tocar outras specs (fonte: plan.md § Summary; 007-4 A).

**Independent Test**: `plan.md` § Pendências e este arquivo listam **DN-007-1 C**, **006-7 D** e pendências ADR/Jira referenciadas (metas SC-005 e gold apontamento fechados na etapa 3e).

- [X] T001 Criar `internal/labeleddataset/doc.go` descrevendo fronteira: leitura de artefatos Git, validação FR-004, ponte para testes do classificador **sem** fila Redis (fonte: plan.md § Project Structure; respostas-joao.md, 007-23 A; FR-009)
- [X] T002 [P] Criar esqueleto de diretório `datasets/labeled-logs/` com `.gitkeep` ou README stub apontando para contratos em `specs/007-base-logs-rotulados/contracts/` na **implementação** (fonte: plan.md § Technical Context Storage; 007-4 A; 007-caminho C — pasta não criada na etapa de spec)
- [X] T003 Atualizar tabela § Mapeamento de campos em `specs/007-base-logs-rotulados/spec.md` com **caminho físico** `datasets/labeled-logs/` e formato **JSONL** (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02, 007-caminho C / 007-formato A; 007-6 D)
- [X] T004 [P] DOCUMENTAÇÃO/PENDÊNCIA — Rastreio **N3 A** já registrado em `specs/007-base-logs-rotulados/plan.md` § Pendências — **não** criar nem editar ADR (fonte: respostas-joao.md, N3 A; plan.md § Pendências; corrige analyze M5)
- [X] T005 [P] DOCUMENTAÇÃO/PENDÊNCIA — Rastrear **010-2 A** (lista de padrões de mascaramento em ADR) e **010-27 A** (casos de teste do João) como dependências da spec **010**; anotar em `specs/007-base-logs-rotulados/research.md` que `message_text` na base 007 é **sem PII** (007-16 A) e mascaramento runtime é ortogonal (fonte: respostas-joao.md, 010-2 A / 010-27 A; contracts/entrada-classificador-desde-base.md)
- [X] T006 [P] DOCUMENTAÇÃO/PENDÊNCIA — Registrar coordenação **006-10 A** (alinhar Jira LLT-197/LLT-199 com spec 006) como ação **fora** desta pasta; em `specs/007-base-logs-rotulados/quickstart.md` §9 citar que métricas 006-2 B consomem esta base quando entregue (fonte: respostas-joao.md, 006-10 A; FR-008)
- [X] T007 Documentar em `specs/007-base-logs-rotulados/spec.md` formato **JSONL** e caminho **`datasets/labeled-logs/`** (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02, 007 formato A / 007 caminho C)
- [X] T008 [P] DOCUMENTAÇÃO/PENDÊNCIA — Referência **DN-009-2** já em `specs/007-base-logs-rotulados/plan.md` § API do Módulo Ocorrência — **sem** tarefas de adaptador HTTP nesta feature (fonte: respostas-joao.md §Etapa 3; plan.md § API do Módulo; corrige analyze M5)
- [X] T056 [P] DOCUMENTAÇÃO/PENDÊNCIA — Rastreio **emenda do princípio I** já em `specs/007-base-logs-rotulados/plan.md` § Pendências — **não** editar `.specify/memory/constitution.md` nesta feature (fonte: plan.md § Constitution Check; respostas-joao.md, T1 REVISTO 15/11; revisao-T1-drive.md §2; corrige analyze C1/M5)
- [X] T059 [P] PRÉ-CONDIÇÃO/DEPENDÊNCIA — Registrar que **T024**, **T031**, **T041**, **T042** e **T045** dependem da **implementação** da spec **006** (operação **Classificar log** conforme `specs/006-modelo-classificacao-logs/contracts/classificar-log.md`), não apenas do mapeamento **T015**; anotar em `specs/007-base-logs-rotulados/plan.md` § Pendências e `specs/007-base-logs-rotulados/quickstart.md` §3 — **sem** implementar código da spec 006 nesta pasta (fonte: analyze-007-etapa3f.md H1; plan.md § Testing; spec.md FR-007, SC-005; respostas-joao.md, 007-22 B / 007-23 A; origem: spec 006 contrato classificar-log.md)

**Checkpoint**: Pacote `internal/labeleddataset` iniciado; diretório de dados reservado; pendências ADR/Jira/DN rastreadas; pré-condição classificador 006 rastreada em **T059**.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Tipos lógicos, validador estrutural e carga de registros — bloqueia todas as histórias (fonte: data-model.md; contracts/registro-rotulado.md; FR-004).

**⚠️ CRITICAL**: Nenhuma user story fecha sem esta fase.

**Independent Test**: Validador rejeita `class_label` inválido, campos vazios (007-13 A) e aceita ERROR/WARNING/INFO (007-9 A).

- [X] T009 Implementar tipos `ClassLabel`, `RegistroLogRotulado`, `RegistroApontamentoRotulado` em `internal/labeleddataset/types.go` alinhados a `specs/007-base-logs-rotulados/data-model.md` (`itemType` `log` | `apontamento`; campos FR-004) (fonte: data-model.md § RegistroLogRotulado / RegistroApontamentoRotulado; FR-002, FR-004)
- [X] T010 [P] Implementar `ValidateRegistro` em `internal/labeleddataset/validate.go`: rejeitar `level_text` ou `message_text` vazios (007-13 A); `class_label` ∈ {ERROR, WARNING, INFO}; proibir spam / "não classificado" (fonte: spec.md Edge Cases; FR-002, FR-012)
- [X] T011 [P] Em `internal/labeleddataset/validate.go`, validar apontamento: `members` não vazio, `documentacao` presente, mesmas regras de privacidade que log único; **exatamente um** `class_label` gold no agregado; cada `member` MUST NOT incluir `class_label` (**DN-007-4 A**) (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-007-4 A; origem: contracts/registro-rotulado.md § Registro de apontamento; 007-7 C; spec.md FR-007b)
- [X] T012 [P] Em `internal/labeleddataset/validate_test.go`, cobrir tabela de rejeição de `contracts/registro-rotulado.md` e cenários SC-002 (fonte: SC-004 / 007-21 A; contracts/registro-rotulado.md)
- [X] T013 Implementar `LoadConjunto` em `internal/labeleddataset/load.go` lendo **JSONL** em `datasets/labeled-logs/` conforme mapeamento (fonte: respostas-joao.md §Etapa 3b, 007 formato A / 007-caminho C; plan.md § Primary Dependencies)
- [X] T014 [P] Implementar `IterRegistros` / interface de leitura em memória em `internal/labeleddataset/load.go` para testes e fixtures **sem** depender de arquivos físicos completos até T020 (fonte: 007-22 B; quickstart.md §3)
- [X] T015 Em `internal/labeleddataset/map_to_classifier.go`, mapear registro `log` para entrada de `specs/006-modelo-classificacao-logs/contracts/classificar-log.md` (`level_text` → `nivelCampo`, `message_text` → `mensagem`) conforme `specs/007-base-logs-rotulados/contracts/entrada-classificador-desde-base.md` (fonte: respostas-joao.md, 007-5 B / 007-6 D / 007-23 A; User Story 2)
- [X] T016 [P] Em `internal/labeleddataset/map_to_classifier_test.go`, provar equaldade gold `class_label` vs predição esperada em cenário **fixture em memória** (sem Redis, sem OpenObserve) (fonte: contracts/entrada-classificador-desde-base.md E1–E2; SC-004)

**Checkpoint**: Validação estrutural testável; ponte para classificador 006 definida.

---

## Phase 3: User Story 1 - Dispor conjunto rotulado para avaliação da classificação (Priority: P1) 🎯 MVP

**Goal**: Milhares de registros, balanceamento ERROR/WARNING/INFO, anonimização, declaração de origem, zero PII (fonte: spec.md User Story 1; FR-001–FR-003, SC-001–SC-003).

**Independent Test**: Ler `datasets/labeled-logs/ORIGEM.md`, manifest de contagens e amostra de registros — rótulos válidos, contagens iguais no mínimo, total em milhares (quickstart.md §1–§2).

### Tests for User Story 1

- [X] T017 [P] [US1] Em `internal/labeleddataset/validate_test.go`, teste de conjunto: 0% rótulos fora de ERROR/WARNING/INFO; falha se PII heurística simples detectada em `message_text` (fonte: 007-16 A; ADR_064; SC-002)
- [ ] T018 [P] [US1] Em `internal/labeleddataset/manifest_test.go`, assertir que contagens publicadas satisfazem **mesmo mínimo** por classe e total na ordem de milhares (fonte: 007-1 C, 007-2 B; SC-003)

### Implementation for User Story 1

- [ ] T019 [US1] Redigir `datasets/labeled-logs/ORIGEM.md` com seções obrigatórias de `specs/007-base-logs-rotulados/data-model.md` § DeclaracaoOrigem (anonimização; sem dado real de cliente; sem PII fictícia; responsável engenharia Finder) (fonte: respostas-joao.md, 007-14 B / 007-15 A / 007-10 B; SC-001)
- [ ] T020 [US1] Popular arquivos **JSONL** de registros anonimizados em `datasets/labeled-logs/` até volume e balanceamento acordados — depende T013 (fonte: FR-001, FR-003, FR-004; 007-1 C, 007-2 B; 007-formato A, 007-caminho C)
- [ ] T021 [P] [US1] Criar `datasets/labeled-logs/MANIFEST-contagens.md` e atualizar `specs/007-base-logs-rotulados/spec.md` § Cobertura alvo com totais reais quando T020 existir (fonte: respostas-joao.md, 007-25 A; SC-003)
- [ ] T022 [P] [US1] Atualizar `specs/007-base-logs-rotulados/checklists/requirements.md` com contagens por classe e referência a `MANIFEST-contagens.md` (fonte: 007-25 A; FR-003)
- [ ] T023 [US1] Documentar em `datasets/labeled-logs/ORIGEM.md` método de rotulação e autoria (007-11 D) e política de duplicatas com rótulos diferentes (007-12 B) (fonte: FR-006, FR-012)

**Checkpoint**: Conjunto publicável cobre SC-001–SC-003 quando T020 completo; até lá, documentação e validador prontos.

---

## Phase 4: User Story 2 - Alinhar registros ao formato de entrada do classificador (Priority: P1)

**Goal**: Campos obrigatórios em 100% dos registros; testes alimentam classificador sem fila (fonte: spec.md User Story 2; FR-004, FR-009; SC-004).

**Independent Test**: Para cada registro `itemType=log`, mapeamento T015 + classificador 006 completa sem Redis (quickstart.md §3–§4).

### Tests for User Story 2

- [X] T024 [P] [US2] Em `internal/labeleddataset/integration_classificador_test.go`, varrer 100% dos registros carregados (ou fixtures) aplicando `map_to_classifier` + `Classificar log` 006 — 0 falhas estruturais (fonte: 007-21 A; SC-004; contracts/entrada-classificador-desde-base.md E1)
- [X] T025 [P] [US2] Em `internal/labeleddataset/integration_classificador_test.go`, provar ausência de import/uso de `github.com/redis/go-redis` e ausência de chamadas à fila no caminho de avaliação (fonte: respostas-joao.md, 007-23 A; FR-009)

### Implementation for User Story 2

- [ ] T026 [US2] Garantir que cada registro em `datasets/labeled-logs/` inclui `lhamalog_batch_id`, `lhamalog_record_id`, `tenantId`, `productId` com valores de **teste** documentados em `datasets/labeled-logs/ORIGEM.md` ou `SPLIT.md` (fonte: respostas-joao.md, 007-5 B / 007-8 A; contrato-log-ingested-v1.txt §4, §7; T5 D)
- [ ] T027 [US2] Sincronizar nomes físicos de colunas/atributos JSON na tabela § Mapeamento de `specs/007-base-logs-rotulados/spec.md` com o artefato real em `datasets/labeled-logs/` após T020 (fonte: 007-6 D; User Story 2 cenário 4)
- [X] T028 [P] [US2] Em `internal/labeleddataset/validate.go`, validar presença de todos os campos obrigatórios FR-004 em 100% dos itens retornados por `LoadConjunto` / iterador (fonte: FR-004; 007-21 A)
- [X] T029 [US2] Documentar em `specs/007-base-logs-rotulados/quickstart.md` §3 o pacote de teste Go alvo (`internal/labeleddataset/...`) e comando `make test` (fonte: quickstart.md; plan.md § Testing)

**Checkpoint**: SC-004 atingível via testes; substituição dos sintéticos 006 preparada (complemento T050).

---

## Phase 5: User Story 3 - Granularidade log vs apontamento e uso pós-MVP (Priority: P2)

**Goal**: Itens log e apontamento documentados; split nomeado; MVP **somente avaliação** (fonte: spec.md User Story 3; FR-005, FR-011; SC-006).

**Independent Test**: `datasets/labeled-logs/SPLIT.md` descreve partições; exemplos `apontamento` com `documentacao`; texto declara 007-3 A (quickstart.md §5–§6).

### Tests for User Story 3

- [X] T030 [P] [US3] Em `internal/labeleddataset/validate_test.go`, casos `itemType=apontamento` com `members` válidos e inválidos; rejeitar `member` com `class_label` e apontamento sem gold agregado (**DN-007-4 A**) (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-007-4 A; 007-7 C; contracts/registro-rotulado.md)
- [X] T031 [US3] Em `internal/labeleddataset/map_to_classifier_test.go`, cenário por apontamento: **um** `class_label` gold no agregado; mapeamento/classificação dos `members` **sem** gold próprio (**DN-007-4 A**) (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-007-4 A; origem: contracts/entrada-classificador-desde-base.md § Registro de apontamento; spec.md FR-007b)

### Implementation for User Story 3

- [X] T032 [US3] Redigir `datasets/labeled-logs/SPLIT.md` com partições nomeadas (ex.: `treinamento`, `validacao`, `teste`, reserva `so_catalogo`) e critérios **qualitativos** — **sem** proporção, seed ou estratificação numérica (fonte: respostas-joao.md, 007-20 B; plan.md § Lacunas documentadas; analyze M2)
- [X] T033 [US3] Incluir em `datasets/labeled-logs/SPLIT.md` critério textual do subconjunto só-catálogo alinhado spec 006 US3; **não** adicionar registros gold nessa partição até catálogo existir (fonte: DN-007-1 C; respostas-joao.md §Etapa 3; FR-011 — **desbloqueado** etapa 3b/3c)
- [ ] T034 BLOQUEADO (006-7 D) — Popular registros na partição `so_catalogo` em `datasets/labeled-logs/` quando catálogo e sintaxe **006-7 D** existirem (fonte: respostas-joao.md, 007-18 B; ADR_167 [PENDENTE]; DN-007-1 C)
- [ ] T035 [P] [US3] Incluir exemplos `itemType=apontamento` em `datasets/labeled-logs/` com **um** `class_label` gold no agregado, `documentacao` presente e `members` **sem** `class_label` (**DN-007-4 A**) (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-007-4 A; respostas-etapa3c.md §3; 007-7 C)
- [X] T036 [US3] Documentar em `datasets/labeled-logs/SPLIT.md` e `ORIGEM.md` que o **uso no MVP** é somente **avaliação das regras** (007-3 A; ADR_167), não treino (fonte: respostas-joao.md, 007-3 A; SC-006)

**Checkpoint**: SC-006 atendido documentalmente; gold apontamento **DN-007-4 A** refletido em T011/T030–T035; só-catálogo bloqueado em conteúdo (T034, T054).

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Métricas offline, quickstart, substituição sintéticos 006, pendências documentais (fonte: plan.md § Pendências; SC-005 / **FR-007c** / **007 H2 D**; FR-008, FR-009).

### Documentação / pendências (sem código de produção novo; sem criar ADR)

- [X] T037 DOCUMENTAÇÃO/PENDÊNCIA — Em `specs/007-base-logs-rotulados/research.md`, acrescentar nota R8 reforçando **007-19 B** (substitui sintéticos 006) e link para tarefa T050 (fonte: respostas-joao.md, 007-19 B)
- [X] T038 [P] DOCUMENTAÇÃO/PENDÊNCIA — Em `specs/007-base-logs-rotulados/spec.md` § Fora de escopo, citar que emenda **ADR_019 / ADR_167 / SSID_007.1** fica fora desta pasta (fonte: plan.md § Pendências)
- [X] T039 [P] DOCUMENTAÇÃO/PENDÊNCIA — Registrar dependência ortogonal **010-2 A** / **010-27 A** em `specs/007-base-logs-rotulados/checklists/requirements.md` (mascaramento não é critério desta base; 007-16 A) (fonte: respostas-joao.md, 010-2 A / 010-27 A)
- [X] T040 DOCUMENTAÇÃO/PENDÊNCIA — Reiterar rastreio **006-10 A** para atualização Jira quando base alimentar métricas 006 — sem editar Jira neste repositório (fonte: respostas-joao.md, 006-10 A; FR-008)
- [X] T057 [P] DOCUMENTAÇÃO — Alinhar `specs/007-base-logs-rotulados/checklists/requirements.md` à fase pós-plan/3e: nota de que lacunas nomeadas (**split** qualitativo) permanecem em `data-model.md` / `research.md` / `checklists/revisao-etapa3c.md` (CHK013, CHK018); **DN-007-3** / **DN-007-4** respondidas na etapa 3e — não como `[NEEDS CLARIFICATION]` na `spec.md` (fonte: analyze-007.md H1, M4; checklists/revisao-etapa3c.md CHK018; respostas-joao.md §Etapa 3c, João 16/11 11h)
- [X] T058 [P] DOCUMENTAÇÃO — Registrar em `specs/007-base-logs-rotulados/quickstart.md` §10 e `research.md` R18 que **FR-013** / 007-24 A (sem volume máximo) está satisfeito por documentação — sem tarefa de código (fonte: respostas-joao.md, 007-24 A; spec.md FR-013; analyze M6)

### Métricas e validação transversal

- [X] T041 [P] Implementar em `internal/labeleddataset/metrics.go` cálculo de precisão, recall e F1 **por classe** conforme **006-2 B** sobre partição `validacao`: `itemType=log` — um par gold `class_label` × `classe` predita por registro; `itemType=apontamento` — classificar cada `member`, prevista = **mais grave** entre predições (**ERROR** > **WARNING** > **INFO**, ADR_170), comparar ao único gold agregado; **unidade** de contagem = **apontamento** (**007 H2 D**; FR-007c) (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), 007 H2 D; origem: ADR_170; respostas-etapa3f.md §D; 006-2 B; contracts/entrada-classificador-desde-base.md E4–E5)
- [X] T042 Implementar gate de aceite em `internal/labeleddataset/metrics_test.go`: **precisão, recall e F1 ≥ 95% por classe** (SC-005 / FR-007 / **FR-007c**; 007-meta; **DN-007-2 B**; **DN-006-2/007-3 B**) usando comparação gold×predição conforme **T041**/`metrics.go` — incluindo apontamentos com **007 H2 D** (ex.: gold ERROR + predições ERROR, ERROR, WARNING → acerto; gold WARNING + predições WARNING, WARNING, ERROR → erro) (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), 007 H2 D; §Etapa 3c, DN-006-2/007-3 B; origem: ADR_170 / 006-2 B; Jira LLT-199 critério 2)
- [X] T043 Atualizar `specs/007-base-logs-rotulados/checklists/requirements.md` com critério LLT-199 critério 2 **completo**: precisão, recall e F1 **≥ 95% em cada classe** ERROR, WARNING e INFO (**DN-006-2/007-3 B**); apontamentos avaliados **por apontamento** após **007 H2 D** (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h01 e 11h, DN-006-2/007-3 B; §Etapa 3e, 007 H2 D; respostas-etapa3c.md §1; spec.md SC-005, FR-007c)
- [X] T044 Atualizar `specs/007-base-logs-rotulados/quickstart.md` §7 para asserts **≥ 95% por classe** em precisão, recall e F1 e regra **007 H2 D** na contagem de apontamentos (prevista = mais grave entre predições dos `members`) (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h, DN-007-2 B; DN-006-2/007-3 B; §Etapa 3e, 007 H2 D; 007-meta; origem: ADR_170)
- [X] T045 [P] Implementar em `internal/labeleddataset/metrics.go` relatório na partição `validacao` com assert **≥ 95% por classe** em precisão, recall e F1 (ERROR, WARNING, INFO), alinhado a **T041** — gold×predição por log e **007 H2 D** por apontamento (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h, DN-007-2 B; DN-006-2/007-3 B; §Etapa 3e, 007 H2 D; 007-meta; origem: ADR_170)
- [ ] T046 [P] Executar e documentar em `specs/007-base-logs-rotulados/quickstart.md` validação dos cenários §1–§6 contra `make test` / `go test ./internal/labeleddataset/...` quando implementação existir (fonte: quickstart.md; plan.md § Testing)
- [X] T047 Garantir que nenhum pacote novo em `internal/labeleddataset/` adiciona endpoint HTTP público (fonte: plan.md § Target Platform; DN-009-2 só referência)
- [X] T048 [P] Auditar que testes da base não importam adaptador `internal/adapter/out/modocorrencia/` nem publicam `ResumoAnalise` (fonte: respostas-joao.md, T1 REVISTO; plan.md § API do Módulo)

### Substituição sintéticos 006

- [X] T049 [P] DOCUMENTAÇÃO/PENDÊNCIA — Em `specs/007-base-logs-rotulados/quickstart.md` §9, descrever coordenação com spec **006** para substituição de sintéticos e consumo de `datasets/labeled-logs/` — **não** editar arquivos fora de `specs/007-base-logs-rotulados/` (fonte: respostas-joao.md, 007-19 B / T6 A; FR-009)
- [ ] T050 Coordenar em `specs/007-base-logs-rotulados/quickstart.md` §9 a remoção de dependência de dados sintéticos de prova da spec 006 quando T020 entregar conjunto — redirecionamento de código na spec **006** fica fora desta pasta (fonte: respostas-joao.md, 007-19 B; FR-009)

### Subconjunto só-catálogo (critério documentado; conteúdo bloqueado)

- [X] T051 [P] Em `datasets/labeled-logs/SPLIT.md`, manifest vazio ou seção **critério de inclusão** para partição `so_catalogo` sem registros inventados (fonte: DN-007-1 C; FR-011 — **desbloqueado** etapa 3b; respostas-joao.md §Etapa 3)
- [ ] T054 BLOQUEADO (006-7 D) — Após catálogo, revisar `specs/007-base-logs-rotulados/contracts/entrada-classificador-desde-base.md` para injeção opcional de `catalogo` nos cenários só-catálogo (fonte: contracts/entrada-classificador-desde-base.md; 006-7 D; DN-007-1 C)

### Referências cruzadas DN-008-4 (fora do escopo)

- [X] T052 [P] DOCUMENTAÇÃO/PENDÊNCIA — Nota em `specs/007-base-logs-rotulados/research.md` R14: **DN-008-4 B** decidido (normalização spec 008) **não** bloqueia montagem da base (007-17 A) (fonte: respostas-joao.md §Etapa 3b, DN-008-4 B; plan.md § Pendências)
- [X] T053 [P] DOCUMENTAÇÃO/PENDÊNCIA — Em `specs/007-base-logs-rotulados/checklists/requirements.md`, registrar **007-meta** completo: precisão, recall e F1 **≥ 95% por classe** (**DN-007-2 B**; **DN-006-2/007-3 B**) e avaliação de apontamento **007 H2 D** / **FR-007c** (fonte: 007-25 A; spec.md SC-005; respostas-joao.md §Etapa 3c, João 16/11 11h; §Etapa 3e, 007 H2 D; origem: ADR_170)
- [X] T055 Atualizar `specs/007-base-logs-rotulados/spec.md` § SC-005 com precisão, recall e F1 **≥ 95% por classe** e regra de apontamento **007 H2 D** (**FR-007c**) alinhada a **T041**/**T042** (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h, DN-007-2 B; DN-006-2/007-3 B; §Etapa 3e, 007 H2 D; origem: ADR_170; etapas 3e–3g)

**Checkpoint**: Métricas 006-2 B reportáveis; SC-005 (precisão/recall/F1 por classe) com **007 H2 D** (etapa 3g); gold apontamento **DN-007-4 A**; base substitui sintéticos quando T020 completo.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Início imediato; especificação formato/caminho em **T003**/**T007** — **T013**/**T020** dependem de implementação, não de lacuna de spec
- **Foundational (Phase 2)**: Após Setup; **bloqueia** user stories
- **User Story 1 (Phase 3)**: Após Foundational — **MVP** (conjunto + origem)
- **User Story 2 (Phase 4)**: Após Foundational; integra T015; **T024**/**T031** após pré-condição **T059** (classificador 006 implementado); paralelo a US1 após validador
- **User Story 3 (Phase 5)**: Após tipos apontamento; gold apontamento **DN-007-4 A** (T031, T035); **T031** (classificação) após **T059**; só-catálogo bloqueado em **006-7 D** (T034, T054)
- **Polish (Phase 6)**: Após US1/US2 desejados; métricas **T041**/**T042**/**T045** após **T059** e integração classificador; T057/T058 documentais em paralelo cedo

### User Story Dependencies

| Story | Depende de | Teste independente |
|-------|------------|-------------------|
| US1 (P1) | Fase 2 | ORIGEM + taxonomia + balanceamento (SC-001–SC-003) |
| US2 (P1) | Fase 2 + mapeamento T015 + pré-condição **T059** | 100% campos + classificador sem fila (SC-004) |
| US3 (P2) | Fase 2 | SPLIT + apontamento documentado (SC-006); exemplos gold apontamento (**DN-007-4 A**, T035) |

### Bloqueios explícitos

| Pendência | Tarefas | Motivo |
|-----------|---------|--------|
| **006-7 D** / **DN-007-1 C** (conteúdo) | T034, T054 | Aguardar catálogo; critério T033/T051 **desbloqueado** (etapa 3b/3c) |
| **DN-009-2** | — | Nenhuma tarefa HTTP na 007 |
| **DN-008-4** | — | Fora do escopo 007 (T052 nota) |
| **N3 A**, **010-2 A**, **010-27 A** | T005, T039 | ADR/casos mascaramento fora desta pasta (rastreio T004 em plan) |
| **006-10 A** | T006, T040 | Alinhamento Jira fora do repo |
| **Implementação spec 006** (`classificar-log.md`) | T024, T031, T041, T042, T045 | Pré-condição **T059**; consumo real de **Classificar log**, não só contrato/mapeamento 007 (fonte: analyze-007-etapa3f.md H1; respostas-joao.md, 007-22 B / 007-23 A; origem: spec 006) |
| **emenda princípio I** | — | Rastreio em plan (T056) |

### Parallel Opportunities

- Fase 1: T005–T006 em paralelo após T001; T002 ∥ T001; T004/T008/T056 em paralelo (rastreio documental)
- Fase 2: T010–T012, T014, T016 em paralelo após T009
- US1: T017 ∥ T018; T021 ∥ T022
- US2: T024 ∥ T025; T028 ∥ T027 (após T020)
- US3: T030 ∥ T032; T036 após T032
- Polish: T038–T041, T045–T046, T048, T051–T052, T057–T058 em paralelo

---

## Parallel Example: User Story 2

```text
# Testes em paralelo:
T024 integration_classificador_test.go — 100% estrutura + classificador
T025 integration_classificador_test.go — zero Redis

# Depois:
T026 tenantId/productId/batch/record nos artefatos JSONL
T027 sincronizar spec § Mapeamento
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Completar Fase 1: Setup  
2. Completar Fase 2: Foundational (CRITICAL)  
3. Completar Fase 3: User Story 1 (declaração + validador + manifest; dados físicos via T020 após T013)  
4. **STOP and VALIDATE**: quickstart §1–§2 + `go test ./internal/labeleddataset/...`  
5. Entregar conjunto mínimo versionado quando volume e balanceamento estiverem registrados  

### Incremental Delivery

1. Setup + Foundational → validador e ponte 006  
2. US1 → conjunto + ORIGEM + SC-001–SC-003  
3. US2 → SC-004 + testes sem fila  
4. US3 → split + apontamentos documentais (**DN-007-4 A**) + SC-006  
5. Polish → métricas 006-2 B (SC-005; **007 H2 D** etapa 3g), substituição 006  

### Parallel Team Strategy

1. Time fecha Fase 2 junto  
2. Dev A: US1 dados + manifest  
3. Dev B: US2 integração classificador  
4. Dev C: US3 SPLIT + T030–T035 (gold **DN-007-4 A**)  
5. Polish: T042–T045/T043 checklist com metas **DN-006-2/007-3 B** e avaliação apontamento **007 H2 D**  

---

## Notes

- [P] = arquivos distintos, sem conflito na mesma fase  
- Cada tarefa cita fonte — formato `(fonte: respostas-joao.md, <ID>; origem: …)`  
- Não commitar dados reais de cliente; só anonimização documentada (007-14 B)  
- Rotulação **só classe** — sem assinatura (007-17 A)  
- **007-meta**: precisão, recall e F1 **≥ 95% por classe** (**DN-007-2 B**; **DN-006-2/007-3 B**, etapa 3e)  
- **007 H2 D**: apontamento — predição por `member`; prevista = mais grave (**ERROR** > **WARNING** > **INFO**, ADR_170); gold×prevista no agregado; unidade = apontamento (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), 007 H2 D; origem: ADR_170; respostas-etapa3f.md §D; FR-007c)  
- Gold apontamento: **DN-007-4 A** (um `class_label` no agregado; fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02))  
- Pendência transversal: **emenda do princípio I** (constituição × T1 revisto) — somente rastreio T056  
- Não criar nem editar ADR; implementação em `internal/labeleddataset/` e `datasets/labeled-logs/` conforme plan; documentação alterável só em `specs/007-base-logs-rotulados/` nas tarefas DOCUMENTAÇÃO/PENDÊNCIA
## Ajuste de execução documentado

Novas notas e matrizes pedidas em arquivos de checklist são registradas em rastreabilidade.md; checklists e marcadores são somente leitura pela skill. Caminhos de código podem usar arquivos isolados e aliases dos donos comuns 006/008; o rastreio aponta a implementação real. Isso não converte aceite bloqueado em concluído.
