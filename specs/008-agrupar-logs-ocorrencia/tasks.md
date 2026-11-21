---
description: "Task list for feature 008-agrupar-logs-ocorrencia"
---

# Tasks: Agrupar logs semelhantes na mesma ocorrência

**Input**: Design documents from `/specs/008-agrupar-logs-ocorrencia/`

**Prerequisites**: [plan.md](./plan.md) (required), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md)

**Tests**: Não solicitados explicitamente na spec — tarefas de teste omitidas (validação manual via [quickstart.md](./quickstart.md)).

**Organization**: Fases por história de usuário (US1–US3, todas P1), conforme [spec.md](./spec.md).

**Última regeneração**: alinhada a `spec.md`/`plan.md`; **Etapa 3h**: **DLQ B** (rev. **DN-008-3** / **DN-008-6** — sem DLQ; aviso Core **batchId** + **motivo**), **DN-008-9 A** — fonte `respostas-joao.md §Etapa 3g (João, 16/11 13h22)`; **Etapa 3g**: **DN-008-8 A** (separador **`|`**).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Pode executar em paralelo (arquivos diferentes, sem dependência de tarefas incompletas)
- **[Story]**: US1, US2 ou US3 — apenas em fases de história
- Cada descrição inclui caminho de arquivo e rastreio `(fonte: …; origem: …)`
- **BLOQUEADO (&lt;pendência&gt;)**: somente na linha da tarefa afetada, após o ID (ex.: `T021 BLOQUEADO (DN-009-2)`)

---

## Pendências (decisões fora desta pasta — não inventar)

| ID | Escopo nesta feature | Fonte |
|----|----------------------|-------|
| **DN-009-2** | Adaptadores HTTP+JWT consulta/envio Módulo; quickstart §5–§8 com Módulo real; SC-005 ponta a ponta | (fonte: respostas-joao.md, Etapa 3; origem: respostas-etapa3.md §10; achado analyze **U1**) |
| **P-010-27 A** | Matriz formal de casos SC-001 (mascaramento/normalização; literais **DN-008-5 A** já fechados) | (fonte: respostas-joao.md, 010-27 A; origem: FR-002b / SC-001) |
| **Contrato API Core (alertas)** | Rota, payload e evento Finder→Core após 3 falhas OpenObserve (**ADR_051 [PENDENTE]**; **batchId** + **motivo**) | (fonte: respostas-joao.md §Etapa 3g, DLQ B; plan.md § Technical Context) |
| **Escape `\`** | Escape de barra invertida literal nos componentes do `valor` (sem fonte) | (fonte: spec.md Questões em aberto §3) |
| **N3 A** | Nomes JSON finais de `AchadoAnalise` ampliado — **ADR nova** (não criar ADR aqui) | (fonte: respostas-joao.md, N3 A; origem: ADR_062; achado analyze **U3**) |
| **Query OpenObserve** | Organização, stream, índice, sintaxe HTTP — **NEEDS CLARIFICATION** operacional | (fonte: `internal/adapter/out/openobserve/doc.go`; achado analyze **U4**) |
| **006-10 A** | Alinhar Jira LLT-200 com specs 006/008 — só registro | (fonte: respostas-joao.md, 006-10 A) |
| **010-2 A** | Lista de padrões de mascaramento em ADR (spec 010) — referência cruzada com **DN-008-4 B** | (fonte: respostas-joao.md, 010-2 A) |
| **DN-007-2 B** | Meta **≥ 95% precisão por classe** — spec 006/007; **fora do escopo 008** | (fonte: respostas-joao.md §Etapa 3b, DN-007-2 B; origem: 007-meta; 006-2 B) |
| **Emenda princípio I** | T1 rev.: Módulo Ocorrência cria ocorrência — **não editar** constituição nesta pasta | (fonte: respostas-joao.md, T1 rev.; origem: ADR_062; achado analyze **C1**) |

**Resolvidas nas etapas 3b–3h (desbloqueadas nas tarefas):** **DN-008-1 B**, **DN-008-2**, **DN-008-3** (rev. 3h: 3× → aviso Core API interna, **sem DLQ** — substitui alcance 3b/3e com DLQ; fonte: respostas-joao.md §Etapa 3g, DLQ B; origem: respostas-etapa3g.md §1b, contrato log.ingested v1 §8); **DN-008-4 B** (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02); **DN-008-5 A** (literais `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` — fonte: respostas-joao.md §Etapa 3c 11h01/11h; origem: DN-008-5 A); **DN-008-6** (rev. 3h: mesmo gatilho — Finder→Core **batchId** + **motivo**; **não** DLQ nem WebSocket pós-DLQ — substitui etapa 3e; fonte: respostas-joao.md §Etapa 3g, DLQ B); **DN-008-9 A** (**`|`** → **`\|`** — fonte: respostas-joao.md §Etapa 3g, DN-008-9 A); **DN-008-7** envio **um ResumoAnalise/lote** + consulta prévia **uma chamada por achado** (**DN-008-7 B** — fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02)); **008-assinatura A** (valor = texto concatenado); **DN-008-8 A** (separador **`|`** entre nível, mensagem normalizada, `service.name`, host — fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40); origem: DN-008-8 A); **DN-006-1** / **DN-006-1b** (execId opcional, N4 A); **DN-006-3 B** (referência cruzada spec 006 — fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02)); **DN-007-4 A** (referência cruzada spec 007 — fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02)); **DN-009-1** (fora do adaptador 008); **DN-007-2 B** e **007 formato A / caminho C** (`datasets/labeled-logs/`, JSONL na raiz — só referência cruzada **T054**, **fora do escopo 008**, sem criar pasta); **contador INFO** (**T7 A**, **T026**); **009-5 revisto** (referência cruzada spec 009, fora adaptador HTTP 008).

Divergências conscientes (João×Drive/Jira): ver [spec.md](./spec.md) § Divergências documentadas — sem entrada 11h×11h01 neste arquivo.

## Decisões tomadas (etapa 3h)

| ID | Decisão | Fonte |
|----|---------|-------|
| **DLQ B** | **Não** há DLQ no fluxo `log.ingested` quando a consulta OpenObserve falha (contrato v1 não herda retry/DLQ do guia; **DN-006-1b B**). | (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22); origem: DLQ B; respostas-etapa3g.md §1) |
| **DN-008-3** (rev. 3h) | Até **3 tentativas**; após esgotar, Finder MUST avisar o **Core** pela **API interna** com **batchId** e **motivo** (**ADR_051 [PENDENTE]**). **Substitui** alcance anterior (3 tentativas + DLQ / reenvio stream — etapas 3b/3e). | (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22); origem: DLQ B) |
| **DN-008-6** (rev. 3h) | Mesmo gatilho: Finder→Core API interna (**batchId** + **motivo**); **não** DLQ nem alerta WebSocket pós-DLQ. **Substitui** alcance da etapa 3e (DLQ + Core/WebSocket). | (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22); origem: DLQ B) |
| **DN-008-9 A** | **`|`** literal em componente do `valor` MUST ser escapado como **`\|`** antes da concatenação. Escape de **`\`** literal: só pendência (sem fonte). | (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22); origem: DN-008-9 A) |

Tarefas afetadas: **T039**, **T040**, **T041** (falha OpenObserve / aviso Core); **T008**, **T015**, **T047** (escape **`|`** — **DN-008-9 A** desbloqueado). Chamada HTTP real Finder→Core após 3 falhas: **BLOQUEADO** até contrato **ADR_051** (ver tabela Bloqueios — **T039** implementa política de tentativas e registro; adaptador **MUST NOT** inventar).

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Preparar extensões hexagonais e configuração mínima sem implementar regra de negócio.

- [X] T001 Criar esqueleto do pacote `internal/adapter/out/modocorrencia/` com `doc.go` descrevendo HTTP+JWT (**N1 A**) e lacunas **DN-009-2** / **DN-008-7** — origem: [plan.md](./plan.md) § Project Structure (fonte: respostas-joao.md, N1 A; origem: ADR_062)
- [X] T002 [P] Declarar arquivos de portas de saída `internal/application/port/out/assinatura.go`, `consulta_ocorrencia.go`, `consulta_recorrencia.go`, `resumo_analise.go` referenciando `specs/008-agrupar-logs-ocorrencia/contracts/` — origem: plan.md § Project Structure (fonte: plan.md, Estrutura hexagonal)
- [X] T003 [P] Estender `internal/config/config.go` e `internal/config/env.go` com placeholders de URL e credencial JWT do Módulo Ocorrência (sem valores fixos; rota **NEEDS CLARIFICATION** até **DN-009-2**) — origem: plan.md § Primary Dependencies (fonte: respostas-joao.md, N1 A; origem: ADR_028, ADR_054)
- [X] T004 [P] Estender `internal/config/config.go` com placeholders de OpenObserve (organização, stream, credenciais) alinhados a `internal/adapter/out/openobserve/doc.go` — origem: plan.md § Storage / consulta (fonte: plan.md, FR-008 / 008-3 B; achado analyze **U4** — query operacional aberta)
- [X] T005 DOCUMENTAÇÃO: Marcar **[PENDENTE ADR] N3 A** (ampliar `AchadoAnalise`: assinatura, contagem, primeira/última vez, versão regras) em `specs/008-agrupar-logs-ocorrencia/contracts/resumo-analise-agrupado.md`, sem criar ADR — origem: plan.md § Pendências (fonte: respostas-joao.md, N3 A; origem: ADR_062; achado analyze **U3**)
- [X] T006 DOCUMENTAÇÃO: Registrar ação **006-10 A** (alinhar Jira LLT-200 com specs 006/008) em `specs/008-agrupar-logs-ocorrencia/checklists/requirements.md` — origem: plan.md § Pendências (fonte: respostas-joao.md, 006-10 A)
- [X] T046 [P] DOCUMENTAÇÃO: Registrar pendência **DN-009-2** (rota, JWT, HTTP, 5xx/timeout) e granularidade fechada **DN-008-7 B** (consulta **por achado**) / envio **um ResumoAnalise por lote** em `specs/008-agrupar-logs-ocorrencia/contracts/consulta-ocorrencias-modulo.md` e `specs/008-agrupar-logs-ocorrencia/contracts/resumo-analise-agrupado.md` — origem: plan.md § API do Módulo Ocorrência (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; origem: ADR_062, N2 B; achado analyze **U1**)
- [X] T047 [P] DOCUMENTAÇÃO: Documentar **DN-008-4 B**, literais **DN-008-5 A** (`<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>`), serialização **008-assinatura A** com separador **`|`** (**DN-008-8 A**) e **DN-008-9 A** (escape **`|`** → **`\|`**); pendência escape **`\`** / matriz **P-010-27 A** em `specs/008-agrupar-logs-ocorrencia/contracts/assinatura-log.md` — origem: plan.md § Pendências (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), DN-008-8 A; §Etapa 3c 11h01/11h, DN-008-5 A; origem: FR-002b/c)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Entidades de domínio, validações transversais e gating de lote — bloqueia todas as histórias.

**⚠️ CRITICAL**: Nenhuma história inicia antes desta fase.

- [X] T007 Definir tipo `ChaveUnicidadeOcorrencia` (assinatura, classe, productId, tenantId) em `internal/domain/chave_ocorrencia.go` — origem: [data-model.md](./data-model.md) § ChaveUnicidadeOcorrencia (fonte: respostas-joao.md, 008-11 rev.; origem: LLT-200, FR-005)
- [X] T008 [P] Definir `AssinaturaLog` em `internal/domain/assinatura.go` com componentes **DN-008-1 B**; `traceId`/`correlationId` no log mas **fora** do cálculo (**008-2 rev.**); campo `valor` como texto **008-assinatura A** com separador fixo **`|`** (**DN-008-8 A**); escape **`|`** → **`\|`** (**DN-008-9 A**); escape **`\`**: pendente — origem: data-model.md § AssinaturaLog (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), DN-008-8 A; §Etapa 3c 11h, 008-assinatura A; origem: DN-008-1 B, FR-002c)
- [X] T009 [P] Definir `LogAgrupavel` e `AchadoLote` (partição por assinatura+classe, contagemLote, primeiraNoLote, ultimaNoLote) em `internal/domain/achado_lote.go` — origem: data-model.md § LogAgrupavel, § Achado (fonte: respostas-joao.md, 008-18 A / 008-12 A / 008-14 A)
- [X] T010 [P] Definir `EstadoAgregacaoEntreLotes` e `MarcaRecorrente` em `internal/domain/agregacao.go` — origem: data-model.md § EstadoAgregacaoEntreLotes, § MarcaRecorrente (fonte: respostas-joao.md, N2 B / 008-5 A / DN-008-3)
- [X] T011 Definir `AchadoAnaliseAgrupado` e `ResumoAnaliseEntrega` em `internal/domain/resumo_analise.go` com campos **N3 A** opcionais até ADR (nomes JSON `[PENDENTE ADR]`) — origem: data-model.md § AchadoAnalise, § ResumoAnalise (fonte: respostas-joao.md, N3 A / N4 A; origem: ADR_070; achado analyze **U3**)
- [X] T012 Implementar validação `productId` ausente ou metadata≠body como falha crítica sem retry em `internal/domain/validar_produto.go` (ou extensão de `internal/domain/envelope.go`) — origem: data-model.md, FR-013 (fonte: respostas-joao.md, 008-16 A; origem: ADR_049)
- [ ] T013 Integrar gate **lote incompleto**: não invocar agrupamento quando ingestão 005 falhar contagem em `internal/application/consumir.go` — origem: plan.md Summary, FR-014 (fonte: respostas-joao.md, T4 A; origem: spec 005 FR-007/008)
- [X] T014 Garantir uso exclusivo de envelope `log.ingested` schemaVersion `"1"` em `internal/domain/envelope.go` no caminho até agrupamento — origem: plan.md § Storage (fonte: respostas-joao.md, T5 D; origem: contrato-log-ingested-v1.txt)
- [X] T016 Implementar `NormalizarMensagem` em `internal/application/assinatura/normalizar_mensagem.go` conforme **DN-008-4 B** e marcadores **DN-008-5 A** — origem: plan.md § Fases / Assinatura (fonte: respostas-joao.md §Etapa 3b, DN-008-4 B; §Etapa 3c 11h01, DN-008-5 A; origem: spec.md FR-002b)
- [X] T015 Implementar cálculo canônico da assinatura e serialização de `valor` (**008-assinatura A**: nível + mensagem normalizada + `service.name` + host unidos por **`|`** entre componentes adjacentes, **DN-008-8 A**) em `internal/application/assinatura/calcular.go` conforme `specs/008-agrupar-logs-ocorrencia/contracts/assinatura-log.md` (depende de T016; escape **`|`** → **`\|`** — **DN-008-9 A**) — origem: plan.md § Fases / Assinatura (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), DN-008-8 A; origem: FR-002c; achado analyze **E3** resolvido etapa 3g)
- [X] T048 [P] Completar constantes/asserts dos cinco literais **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`** em `internal/application/assinatura/normalizar_mensagem.go` (**DN-008-5 A**; matriz formal **P-010-27 A** permanece em T051) — origem: plan.md § Fases / Assinatura (fonte: respostas-joao.md §Etapa 3c 11h01/11h, DN-008-5 A; origem: FR-002b)

**Checkpoint**: Domínio e portas prontos; normalização **DN-008-4 B** + literais **DN-008-5 A** (T016, T048); serialização `valor` com **`|`** (**DN-008-8 A**, T015) e escape **`|`** → **`\|`** (**DN-008-9 A**, T015); escape **`\`** pendente.

---

## Phase 3: User Story 1 - Unificar logs equivalentes em uma única ocorrência (Priority: P1) 🎯 MVP

**Goal**: Agrupar no lote, fusão entre lotes (**N2 B**), handoff `ResumoAnalise`; Módulo upsert pela chave 008-11 (**T1 rev.**) (fonte: spec.md, User Story 1; FR-001, FR-004, FR-009).

**Independent Test**: Dois lotes com logs equivalentes (mesma chave) → uma ocorrência, contador e extremos corretos (timestamps do log) — spec.md User Story 1.

### Implementation for User Story 1

- [X] T017 [P] [US1] Implementar porta `internal/application/port/out/assinatura.go` conforme `specs/008-agrupar-logs-ocorrencia/contracts/assinatura-log.md` — origem: plan.md § Project Structure (fonte: spec.md, FR-002)
- [X] T018 [P] [US1] Implementar porta `internal/application/port/out/consulta_ocorrencia.go` conforme `specs/008-agrupar-logs-ocorrencia/contracts/consulta-ocorrencias-modulo.md` — origem: [research.md](./research.md) R2 (fonte: respostas-joao.md, N2 B)
- [X] T019 [US1] Implementar particionamento de logs classificados em achados por `(assinatura, classe)` em `internal/application/agrupar/particionar_lote.go` — origem: [contracts/agregar-achado-lote.md](./contracts/agregar-achado-lote.md) (fonte: respostas-joao.md, 008-18 A / 008-15 A / 008-12 A / 008-13 A)
- [X] T020 [US1] Calcular `contagemLote`, `primeiraNoLote`, `ultimaNoLote` por min/max timestamp do log em `internal/application/agrupar/metricas_lote.go` — origem: data-model.md § Achado (fonte: respostas-joao.md, 008-14 A; origem: SC-002)
- [X] T049 [P] [US1] Implementar consulta em memória/fake do Módulo (**uma chamada por achado**, **DN-008-7 B**) para fusão local em `internal/application/agrupar/consulta_memoria.go` (sem HTTP; contrato lógico **desbloqueado**; HTTP real **DN-009-2**) — origem: quickstart.md §5 (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; origem: N2 B)
- [X] T021 [US1] Usar a consulta HTTPS/JWT da spec 009 por achado, preservando a revisão opaca para a entrega. Verificado na composição das branches (009 T050).
- [X] T022 [US1] Implementar fusão entre lotes: `contagem` = acumulada + `contagemLote`; extremos min/max em `internal/application/agrupar/fusao_entre_lotes.go` — origem: data-model.md § EstadoAgregacaoEntreLotes (fonte: respostas-joao.md, N2 B; 008-12 A)
- [X] T023 [US1] Montar `AchadoAnaliseAgrupado` pós-fusão (campos **N3 A** com placeholders até ADR; lógica em memória **desbloqueada**) em `internal/application/agrupar/montar_achado.go` — origem: [contracts/resumo-analise-agrupado.md](./contracts/resumo-analise-agrupado.md) (fonte: respostas-joao.md, N3 A; origem: SC-005 parcial; achado analyze **U3**)
- [X] T050 [US1] Montar **exatamente um** `ResumoAnaliseEntrega` **por lote** com idempotência `batchId` + `execId` opcional em `internal/application/agrupar/montar_resumo.go` (**DN-008-7**) — origem: data-model.md § ResumoAnalise (fonte: respostas-joao.md §Etapa 3c 11h, DN-008-7; origem: N4 A, ADR_062)
- [X] T024 [US1] Usar o POST de ResumoAnalise DN-009-2 da spec 009 e permitir recálculo após revisão desatualizada. Verificado na composição das branches (009 T018/T019).
- [X] T025 [US1] Orquestrar classificação 006 → assinatura → agrupamento **local determinístico** → **consulta por achado** (memória ou HTTP) → **um ResumoAnalise/lote** em `internal/application/agrupar/orquestrar.go` — origem: plan.md Summary (fonte: spec.md, FR-001 / FR-008 / FR-009; T1 rev.; DN-008-7 B; achado analyze **A1**)
- [X] T026 [US1] Incluir achados classe **INFO** no mesmo fluxo que ERROR/WARNING em `internal/application/agrupar/orquestrar.go` — origem: spec.md FR-006 (fonte: respostas-joao.md, T7 A; origem: ADR_170, D1)
- [X] T027 [US1] Documentar roteiro de validação quickstart §1/§2 para equivalência por **DN-008-4 B** e marcadores **DN-008-5 A** em `specs/008-agrupar-logs-ocorrencia/quickstart.md` — origem: spec.md SC-001 / FR-002b (fonte: respostas-joao.md §Etapa 3b/3c, DN-008-4 B / DN-008-5 A)
- [ ] T051 BLOQUEADO (P-010-27 A) [US1] Fechar matriz formal de casos SC-001 em `specs/008-agrupar-logs-ocorrencia/quickstart.md` §1 quando João publicar **010-27 A** (literais **DN-008-5 A** já fechados em T048) — origem: spec.md User Story 1 cenário 5 (fonte: respostas-joao.md, 010-27 A; origem: FR-002b)

**Checkpoint**: MVP de agrupamento e handoff **em memória** (**consulta por achado**, **envio por lote**); HTTP **DN-009-2**; **N3 A** aguardam ADR; SC-001 formal **P-010-27 A** (T051); serialização **`|`** (**DN-008-8 A**, T015) e escape **DN-008-9 A**; escape **`\`** pendente.

---

## Phase 4: User Story 2 - Evitar duplicatas pela chave acordada (Priority: P1)

**Goal**: Garantir uma ocorrência lógica por `(assinatura, classe, produto, tenant)`; sem migração legada em produção (**DN-008-2**) (FR-005, FR-009).

**Independent Test**: N≥2 logs mesma chave → consulta retorna uma ocorrência; produtos/classes/tenants distintos → chaves distintas — spec.md User Story 2.

### Implementation for User Story 2

- [X] T028 [P] [US2] Validar invariante “no máximo uma partição por chave no lote” em `internal/application/agrupar/particionar_lote.go` — origem: data-model.md § ChaveUnicidadeOcorrencia (fonte: spec.md, FR-005; 008-11 rev.)
- [X] T029 [US2] Rejeitar handoff se consulta retornar múltiplas ocorrências para mesma chave (erro de contrato) em `internal/application/agrupar/fusao_entre_lotes.go` — origem: spec.md User Story 2 cenário 1 (fonte: respostas-joao.md, 008-11 rev.; origem: LLT-200 crit. 2)
- [X] T030 [P] [US2] Separar achados quando mesma assinatura e `classe` ou `productId` ou `tenantId` diferem em `internal/application/agrupar/particionar_lote.go` — origem: spec.md User Story 1 cenário 3 (fonte: respostas-joao.md, 008-15 A / 008-11 rev.)
- [X] T031 [US2] Documentar no `internal/adapter/out/modocorrencia/doc.go` que unicidade persistida é do Módulo Ocorrência (Finder **MUST NOT** gravar ocorrência) — origem: FR-001 (fonte: respostas-joao.md, T1 rev.; origem: ADR_062; achado analyze **C1** — implementação segue T1 rev. até emenda I)
- [X] T032 DOCUMENTAÇÃO [US2]: Registrar **DN-008-2** (sem migração em produção documentada; ressalva D2b teste) em `specs/008-agrupar-logs-ocorrencia/research.md` — origem: spec.md User Story 2 cenário 2 (fonte: respostas-joao.md, DN-008-2; origem: contrato v1 §11, T1 rev.)

**Checkpoint**: Lógica de chave consistente; migração legada **não** exigida (**DN-008-2** fechado).

---

## Phase 5: User Story 3 - Tratar retorno após resolução (Priority: P1)

**Goal**: Recorrência via OpenObserve (008-5/6); reabertura FECHADA→ABERTA no Módulo; contador não altera status (008-10) (FR-003, FR-007, FR-010).

**Independent Test**: Ocorrência FECHADA + log equivalente → Módulo ABERTA; contador sobe com status ABERTA inalterado — spec.md User Story 3 (SC-003).

### Implementation for User Story 3

- [X] T033 [P] [US3] Implementar porta `internal/application/port/out/consulta_recorrencia.go` conforme `specs/008-agrupar-logs-ocorrencia/contracts/consulta-recorrencia-openobserve.md` — origem: plan.md § Storage (fonte: respostas-joao.md, T3 B / 008-3 B)
- [X] T034 [US3] Implementar avaliação ADR_013 (3–4 vezes em lotes distintos, mesma máquina; repetição no mesmo lote não conta para recorrência) em `internal/application/recorrencia/avaliar.go` — origem: research.md R4 (fonte: respostas-joao.md, 008-5 A / 008-13 A)
- [X] T035 [US3] Limitar janela da consulta aos últimos **500 lotes** em `internal/application/recorrencia/avaliar.go` — origem: [contracts/consulta-recorrencia-openobserve.md](./contracts/consulta-recorrencia-openobserve.md) (fonte: respostas-joao.md, 008-6 A; origem: ADR_063)
- [X] T036 [US3] Preencher `recorrente` boolean no `AchadoAnaliseAgrupado` sem mapear a status em `internal/application/agrupar/montar_achado.go` — origem: spec.md FR-007 (fonte: respostas-joao.md, 008-10 A; origem: ADR_063)
- [X] T037 [US3] Propagar `status` FECHADA da consulta ao Módulo sem alterar status no Finder em `internal/application/agrupar/fusao_entre_lotes.go` — origem: contracts/consulta-ocorrencias-modulo.md § Semântica (fonte: respostas-joao.md, 008-4 B / 008-8 A)
- [X] T038 [US3] Implementar adapter OpenObserve stub com semântica documentada em `internal/adapter/out/openobserve/recorrencia.go` (query HTTP/index **NEEDS CLARIFICATION** até fonte operacional — implementar só stub) — origem: `internal/adapter/out/openobserve/doc.go` (fonte: plan.md, FR-008; 008-7 B; achado analyze **U4**)
- [ ] T039 [US3] Implementar política **DN-008-3** + **DN-008-6** (rev. 3h) em `internal/adapter/out/openobserve/recorrencia.go`: até **3 tentativas** na consulta; **MUST NOT** `recorrente=false` por degradação; **MUST NOT** DLQ no `log.ingested`; após esgotar tentativas, registrar falha do lote (unidade = lote `log.ingested` v1) e preparar aviso ao **Core** (**batchId** + **motivo**); chamada real Finder→API Core **BLOQUEADO (Contrato API Core / ADR_051 [PENDENTE])** — **MUST NOT** inventar rota/payload/evento; publicação de evento/notificação ao Core permanece spec **009** — origem: spec.md FR-002a / SC-007 (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B; plan.md § Technical Context)
- [ ] T040 [US3] Integrar orquestração de recorrência e encadeamento pós-3-falhas alinhado a **T039** (3 tentativas; aviso Core **DN-008-6** rev. 3h; **sem** DLQ/WebSocket pós-DLQ) em `internal/application/agrupar/orquestrar.go` dependente de T039 — origem: plan.md § Fases / Recorrência (fonte: spec.md, FR-003 / FR-008 / FR-002a)

**Checkpoint**: Recorrência e reabertura modeladas; falha OpenObserve **DN-008-3** + aviso Core (**DN-008-6** rev. 3h); query HTTP OpenObserve operacional ainda **NEEDS CLARIFICATION** (T038).

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentação, dependências externas e validação guiada.

- [X] T041 [P] DOCUMENTAÇÃO: Registrar dependência ADR_035 / LLT-216 (**008-17 C**) e decisão **DN-008-6** (rev. 3h: Finder avisa Core após 3 falhas OO — **batchId** + **motivo**; **ADR_051 [PENDENTE]** — **MUST NOT** inventar rota/payload/evento; publicação de evento spec **009**) em `specs/008-agrupar-logs-ocorrencia/checklists/requirements.md` e `specs/008-agrupar-logs-ocorrencia/research.md` — origem: research.md R12 (fonte: respostas-joao.md §Etapa 3g, DLQ B; respostas-documentadas.md, 008-17)
- [X] T042 [P] DOCUMENTAÇÃO: Atualizar rastreio de divergências **D1**, **D2** (INFO; assinatura sem traceId/correlationId; **008-2 A revogada**, vigente **008-2 rev.**) em `specs/008-agrupar-logs-ocorrencia/quickstart.md` §9 e §1 — origem: spec.md Divergências (fonte: respostas-joao.md, T7 A / DN-008-1 B; achado analyze **I1** — rastreio em tasks, spec já clarificado)
- [ ] T043 Validar cenários 2–4 e 6–7 de `specs/008-agrupar-logs-ocorrencia/quickstart.md` contra implementação local (`make test` / `make run` conforme README); §5–§8 com Módulo real permanecem **DN-009-2** — origem: quickstart.md (fonte: spec.md, SC-001–SC-003 parcial / SC-005)
- [X] T044 Revisar alinhamento de campos **[PENDENTE ADR] N3** com leitura de `specs/006-modelo-classificacao-logs/data-model.md` (somente leitura cruzada; **não** editar spec 006) — origem: data-model.md § AchadoAnalise (fonte: respostas-joao.md, N3 A / 006-3 C; achado analyze **U3**)
- [X] T045 Documentar origem de `execId` (Varredura; opcional sem `log.ingested` — **DN-006-1** / **DN-006-1b**) na montagem idempotente em `internal/application/agrupar/montar_resumo.go` e comentário em `internal/adapter/out/modocorrencia/enviar_resumo.go` — origem: plan.md § Pendências resolvidas (fonte: respostas-joao.md, Etapa 3 / DN-006-1; origem: N4 A)
- [X] T052 [P] DOCUMENTAÇÃO: Registrar pendência **010-2 A** (lista de padrões de mascaramento em ADR na spec 010) e vínculo com categorias **DN-008-4 B** em `specs/008-agrupar-logs-ocorrencia/research.md` — origem: plan.md § Pendências (fonte: respostas-joao.md, 010-2 A; origem: ADR_064)
- [X] T053 [P] DOCUMENTAÇÃO: Registrar pendência **emenda do princípio I** (T1 rev.; Finder não cria ocorrência) em `specs/008-agrupar-logs-ocorrencia/research.md` sem alterar `.specify/memory/constitution.md` — origem: plan.md § Constitution Check (fonte: respostas-joao.md, T1 rev.; origem: ADR_062; achado analyze **C1**)
- [X] T054 [P] DOCUMENTAÇÃO: Referência cruzada **007 formato A / caminho C** (base JSONL em `datasets/labeled-logs/` na raiz — **fora do escopo 008**, sem criar pasta) em `specs/008-agrupar-logs-ocorrencia/research.md` — origem: plan.md § Pendências (fonte: respostas-joao.md §Etapa 3b, DN-007-2 B; origem: spec 007)
- [X] T055 [P] DOCUMENTAÇÃO: Registrar dependência obrigatória de evolução em `specs/006-modelo-classificacao-logs/` (**FR-009**, conflito **D8**, **T1 rev.**: a 006 MUST deixar de inserir ocorrência direta e alinhar handoff **ResumoAnalise**/classificação ao agrupamento 008; complementa **T006** Jira e leitura **T044**) em `specs/008-agrupar-logs-ocorrencia/research.md` — somente rastreio nesta pasta, **MUST NOT** editar spec 006 — origem: spec.md FR-009 / tabela D8 (fonte: respostas-joao.md, T1 rev.; origem: ADR_062, spec 006 FR-010; achado analyze **E1**)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Sem dependências — inicia imediato
- **Foundational (Phase 2)**: Depende de Phase 1 — **bloqueia** todas as histórias
- **User Stories (Phases 3–5)**: Dependem de Phase 2; US2 e US3 assumem fluxo US1 parcialmente implementado
- **Polish (Phase 6)**: Após histórias desejadas; T043 pode usar MVP US1

### User Story Dependencies

| História | Depende de | Independência de teste |
|----------|------------|-------------------------|
| **US1 (P1)** | Foundational | Dois lotes, mesma chave, contador/extremos (mock Módulo **T049**) |
| **US2 (P1)** | Foundational, particionamento US1 | Mesma chave → uma ocorrência; dimensões distintas → chaves distintas |
| **US3 (P1)** | Foundational, fusão US1 | FECHADA → reabertura no Módulo; `recorrente` separado de status |

### Bloqueios (tarefas)

| Pendência | Tarefas |
|-----------|---------|
| **DN-009-2** | T021, T024; quickstart §5–§8 com HTTP real (T043 parcial) |
| **Escape `\`** | Escape de barra invertida nos componentes do `valor` (T008/T015 parse simétrico; SC-005 parcial) |
| **P-010-27 A** | T051 (matriz formal SC-001; literais cobertos por T048) |
| **N3 A (ADR)** | Nomes JSON finais em T023/T044; SC-005 formal ponta a ponta com Módulo real |
| **Contrato API Core (alertas)** | Integração Finder→Core pós-3-falhas OO (**ADR_051 [PENDENTE]**) em **T039**/**T041** — **MUST NOT** inventar adaptador |

### Dependências entre features (rastreio)

| Dependência | Tarefas | Notas |
|-------------|---------|-------|
| **Evolução spec 006** (**FR-009**, **D8**, **T1 rev.**) | **T025**, **T044**, **T055**; apoio **T006** (**006-10 A** Jira) | 006 MUST evoluir em `specs/006-modelo-classificacao-logs/` — fora do escopo de edição 008 (fonte: spec.md FR-009/D8; achado analyze **E1**) |
| **Spec 009** (publicação evento Core) | **T039**, **T041** | Aviso Finder→Core após falha OO (**DN-008-6** rev. 3h) distinto da publicação de evento 008 |

### Parallel Opportunities

- **Phase 1**: T002, T003, T004, T046, T047 em paralelo
- **Phase 2**: T008, T009, T010 em paralelo após T007; **T016 → T048 → T015** (separador **`|`**, **DN-008-8 A**) antes de T019
- **US1**: T017, T018, T049 em paralelo; T019 após T016/T048/T015
- **US2**: T028 e T030 em paralelo
- **US3**: T033 em paralelo com T034 após Foundational
- **Polish**: T041, T042, T052, T053, T054, T055 em paralelo

### Parallel Example: User Story 1

```text
Paralelo: T017 (port/out/assinatura.go) + T018 (port/out/consulta_ocorrencia.go) + T049 (consulta_memoria.go)
Sequencial: T016 → T048 literais → T015 assinatura (**`|`**, **DN-008-8 A**) → T019 particionamento → T022 fusão → T050 montar_resumo → T025 orquestração (dependência 006: **T055** em paralelo com polish)
HTTP real: T021 e T024 após DN-009-2 (consulta **por achado**; envio **por lote** — DN-008-7 B já fechado)
SC-001 formal: T051 após P-010-27 A
```

---

## Implementation Strategy

### MVP First (User Story 1 apenas)

1. Completar Phase 1 + Phase 2 (**T016**, **T048**, **T015** com **`|`** **DN-008-8 A**)
2. Completar Phase 3 (US1) até T026 com **T049** (consulta memória **por achado**), **sem** T021/T024
3. **Parar e validar** quickstart §2–§4, §6–§7 (§5–§8 HTTP **DN-009-2**; §1 matriz formal **P-010-27 A**)
4. US2 e US3 incrementais

### Incremental Delivery

1. Setup + Foundational → base de domínio + **DN-008-4 B** ativo (T016)
2. US1 → agrupamento e handoff em memória (MVP)
3. US2 → invariantes de unicidade + nota DN-008-2
4. US3 → recorrência (**DN-008-3** implementado; stub OO **U4**)
5. Polish + desbloqueios conforme **DN-009-2**, **ADR_051**, escape **`\`**, **P-010-27 A**, ADR **N3**; rastreio 006 (**T055**)

### Escopo sugerido de MVP

- **Incluir**: T001–T015, T016, T048, T017–T020, T022–T027, T049–T050, T028–T032, T033–T040 (**DN-008-6** expectativa em T039, **C2**), documentação T005–T006, T046–T047, **T055**
- **Excluir até pendência**: T021, T024 (**DN-009-2**); T051 (**P-010-27 A**); aceite SC-005 ponta a ponta com Módulo real / **ADR_051** / escape **`\`**; adaptador Core inventado
- **Documentação obrigatória MVP**: T005, T006, T032, T046, T047, T052–T055

---

## Notes

- Finder **MUST NOT** persistir ocorrência (FR-001; **T1 rev.**)
- **DN-007-2 B** e base **007** (`datasets/labeled-logs/`, JSONL) não são implementação desta feature — só T054
- **010-2 A** é spec **010** — T052 registra pendência sem criar ADR
- **DN-008-8 A** / **DN-008-9 A**, **DN-009-2**, query OpenObserve, contrato API Core (**ADR_051 [PENDENTE]**), escape **`\`**: lacunas explícitas; **DN-008-3/5/6/7** fechados (3h); evolução spec **006** rastreada em **T055** (**FR-009**/**D8**)
- Formato de tarefa: checkbox, ID **T001–T055**, `[P]` quando aplicável, `[US#]` nas fases 3–5, caminhos explícitos
- Total: **55 tarefas** | Setup: 8 | Foundational: 11 (incl. T048) | US1: 14 | US2: 5 | US3: 8 | Polish: 9 (incl. **T055**)
## Ajuste de execução documentado

Novas notas e matrizes pedidas em arquivos de checklist são registradas em rastreabilidade.md; checklists e marcadores são somente leitura pela skill. Caminhos de código podem usar arquivos isolados e aliases dos donos comuns 006/008; o rastreio aponta a implementação real. Isso não converte aceite bloqueado em concluído.
