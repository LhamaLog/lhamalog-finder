---
description: "Task list for feature 009-publicar-evento-ocorrencia"
---

# Tasks: Publicar evento de ocorrência (escopo Finder — handoff ResumoAnalise)

**Input**: Design documents from `/specs/009-publicar-evento-ocorrencia/`

**Prerequisites**: [plan.md](./plan.md) (required), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md)

**Tests**: Não solicitados explicitamente na spec — tarefas de `go test` omitidas; validação manual via [quickstart.md](./quickstart.md).

**Organization**: Fases por história de usuário (US1–US3, todas P1), conforme [spec.md](./spec.md).

**Etapa 3f**: Rodadas §Etapa 3c 11h01 e 11h **não** divergem — «puladas» **SUPERADO** (`respostas-joao.md` linha 39). **DN-006-3 B**, **DN-007-4 A** (referência cruzada specs 006/007) e **DN-008-7 B** nas tarefas desta lista → fonte `respostas-joao.md §Etapa 3c (João, 16/11 11h02)`. Divergências João×Drive/Jira: [spec.md](./spec.md) § Divergências documentadas — **sem** tabela 11h×11h01 nem **D16** neste arquivo.

**Etapa 3h**: **DLQ B** — **sem DLQ** no **log.ingested v1**; **DN-008-3**/**DN-008-6** (rev.) — após **3 falhas** na consulta **OpenObserve**, Finder avisa **Core** por API interna (**`batchId`**, **motivo**); contrato **[PENDENTE ADR_051]**; **DN-008-9 A** — escape **`|`** → **`\|`** nos componentes da assinatura (fonte: `respostas-joao.md §Etapa 3g (João, 16/11 13h22)`; origem: `respostas-etapa3g.md §1` / §1b).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Pode executar em paralelo (arquivos diferentes, sem dependência de tarefas incompletas)
- **[Story]**: US1, US2 ou US3 — apenas em fases de história
- **BLOQUEADO (…)**: não implementar até a pendência indicada (marcador **somente** em tarefas bloqueadas)
- Cada descrição inclui caminho de arquivo e rastreio `(fonte: …; origem: …)`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Esqueleto do adaptador ao Módulo Ocorrência, configuração e pendências documentais — sem publicar evento LLT-201.

- [X] T001 Criar `internal/adapter/out/modocorrencia/doc.go` documentando HTTP+JWT (N1 A), **DN-008-7** (um **ResumoAnalise** por lote no POST) vs **DN-008-7 B** (consulta **por achado**, contrato em **DN-009-2**), proibições FR-F02 e bloqueio **DN-009-2** (rota, claims, confirmação) — origem: plan.md § API do Módulo Ocorrência (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; respostas-etapa3c.md §6, DN-008-7; origem: ADR_141, ADR_062)
- [X] T002 [P] Declarar porta de saída em `internal/application/port/out/modulo_ocorrencia.go` (`EntregarResumoAnalise` — **DN-008-7**; `ConsultarOcorrenciaPorAchado` — **DN-008-7 B** / **FR-F09**; wiring HTTP bloqueado **DN-009-2**) referenciando [contracts/resumo-analise-handoff.md](./contracts/resumo-analise-handoff.md) — origem: plan.md § Structure Decision (fonte: spec.md, FR-F01 / FR-F09; research.md R2; respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B)
- [X] T003 [P] Estender `internal/config/config.go` e `internal/config/env.go` com placeholders de URL base e credencial JWT do Módulo Ocorrência (sem valores fixos; wiring efetivo bloqueado por **DN-009-2** / **ADR_028**) — origem: plan.md § Primary Dependencies (fonte: respostas-joao.md, N1 A; origem: ADR_028)
- [X] T004 DOCUMENTAÇÃO/PENDÊNCIA: Registrar **N3 A** (ADR nova — ampliar `AchadoAnalise`), **DN-008-4 B** / **DN-008-5 A** (`<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` — canônico spec **008**), **008 assinatura A** (concatenação **DN-008-1 B**; separador fixo **`|`** — **DN-008-8 A**; escape **`|`** nos componentes → **`\|`** — **DN-008-9 A**; dono canônico spec **008**), **DN-008-7** (um resumo por lote) e **DN-008-7 B** (consulta por achado) em `specs/009-publicar-evento-ocorrencia/contracts/resumo-analise-handoff.md` e `specs/009-publicar-evento-ocorrencia/data-model.md` — **[PENDENTE ADR]** campos N3 — **não** criar ADR nesta feature — origem: plan.md § Pendências (fonte: respostas-joao.md §Etapa 3c 11h, DN-008-5 A / 008 assinatura A; respostas-joao.md §Etapa 3e (João, 16/11 12h40), DN-008-8 A; respostas-joao.md §Etapa 3g (João, 16/11 13h22), DN-008-9 A; respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; respostas-etapa3c.md §6; N3 A §Etapa 3b)
- [X] T005 DOCUMENTAÇÃO/PENDÊNCIA: Registrar ação **006-10 A** (alinhar Jira LLT-201/200 com specs 006/008/009 e transição T6 A) em `specs/009-publicar-evento-ocorrencia/checklists/requirements.md` — **não** alterar issues Jira — origem: plan.md § Pendências (fonte: respostas-joao.md, 006-10 A / T6 A)
- [X] T006 DOCUMENTAÇÃO/PENDÊNCIA: Registrar dependência **010-2 A** (lista de padrões de mascaramento em ADR futura) para `AchadoAnalise.linhas` em `specs/009-publicar-evento-ocorrencia/checklists/requirements.md` — origem: plan.md § Pendências (fonte: respostas-joao.md, 010-2 A; origem: spec 010)
- [X] T007 DOCUMENTAÇÃO/PENDÊNCIA: Registrar **P-010-27** (casos de teste do João na spec 010) bloqueando aceite formal de linhas mascaradas no handoff em `specs/009-publicar-evento-ocorrencia/quickstart.md` § Evidências — origem: plan.md § Pendências (fonte: respostas-joao.md, 010-27 A)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Tipos de domínio do handoff, validação de taxonomia e gating transversal — bloqueia todas as histórias.

**⚠️ CRITICAL**: Nenhuma história inicia antes desta fase.

- [X] T008 Definir `ChaveIdempotenciaHandoff` (`batchId`, `execId` opcional) em `internal/domain/idempotencia_handoff.go` alinhado a [data-model.md](./data-model.md) § ChaveIdempotenciaHandoff — origem: plan.md § Constraints (fonte: respostas-joao.md, N4 A / DN-006-1b B; origem: contrato 08/11, ADR_070)
- [X] T009 [P] Definir struct `ResumoAnalise` (contadores com **informacao** / contador INFO — §Etapa 3, `achados[]`, metadados) em `internal/domain/resumo_analise.go` conforme [data-model.md](./data-model.md) § ResumoAnalise — origem: data-model.md (fonte: revisao-T1-drive.md, 009-8; respostas-joao.md §Etapa 3 contador INFO; origem: ADR_062)
- [X] T010 [P] Definir struct `AchadoAnalise` com `tipo` ∈ `ERROR|WARNING|INFO`, `linhas`, `recorrente` e campos Drive documentados em `internal/domain/achado_analise.go` — origem: data-model.md § AchadoAnalise (fonte: spec.md, FR-F08 / FR-F04; respostas-documentadas.md, 009-13 A)
- [X] T011 [P] Definir tipo `ConfirmacaoHandoff` em `internal/domain/confirmacao_handoff.go` (campo `aceito` lógico; mapeamento HTTP/códigos **pendente DN-009-2**) — origem: data-model.md § ConfirmacaoHandoff (fonte: spec.md, FR-F05; plan.md § DN-009-2)
- [X] T012 Implementar validação que rejeita `tipo` fora de ERROR/WARNING/INFO em `internal/domain/achado_analise.go` (ou `internal/domain/validar_taxonomia.go`) — origem: spec.md SC-F003 (fonte: respostas-joao.md, taxonomia 14/11; origem: ADR_170)
- [X] T013 Garantir leitura de `Referencia.BatchID` para `batchId` em `internal/domain/envelope.go` no caminho de montagem do handoff — origem: data-model.md (fonte: respostas-joao.md, N4 A / T5 D / DN-006-1b B; origem: `internal/domain/envelope.go`, contrato-log-ingested-v1)
- [X] T014 Implementar gate **lote incompleto**: não montar handoff de sucesso quando T4 A aplicável em `internal/application/handoff/validar_lote.go` — origem: plan.md § Constraints (fonte: respostas-joao.md, T4 A; origem: spec 005)
- [X] T015 Implementar regra de chave em `internal/domain/resumo_analise.go`: **`batchId` somente** no fluxo **log.ingested v1** (**DN-006-1b B**); incluir `execId` opcional `V|M+uuid` quando disponível (**DN-006-1**, **N4 A**) — origem: data-model.md § ChaveIdempotenciaHandoff (fonte: respostas-joao.md §Etapa 3, DN-006-1 / DN-006-1b B; origem: ADR_075, ADR_069)
- [X] T016 Serializar assinatura, contagem, extremos em UTC e revisão consultada. A versão das regras no resumo vale para todos os achados. Verificado em `internal/adapter/out/modocorrencia/serializar_resumo_test.go` conforme N3 A e DN-009-2.
- [X] T048 [P] Reutilizar normalizador canônico da spec 008 em `internal/application/assinatura/normalizar_mensagem.go`, por meio do montador 006 e preparador 008; sem criar `internal/domain/normalizar_mensagem.go`. Pré-requisito: 008 T016/T048; verificar literais DN-008-5 A e FR-F11 na composição.
- [X] T051 [P] Reutilizar composição textual canônica da spec 008 em `internal/application/assinatura/calcular.go` via resumo preparado; sem criar `internal/domain/compor_assinatura.go`. DN-008-8 A / DN-008-9 A; barra invertida literal continua pendente.
- [ ] T049 Integrar **FR-F07**/**DN-008-3** (rev.) e **DN-008-6** (rev.): até **3 tentativas** na consulta **OpenObserve**; **sem DLQ** / **sem** WebSocket pós-DLQ no ingested (**DLQ B**); após esgotar, acionar aviso **Finder → Core** por API interna (**`batchId`** + **motivo** — gating e porta em `internal/application/consumir.go`; **não** invoca handoff em `internal/application/handoff/orquestrar_entrega.go` até análise concluída); **BLOQUEADO (ADR_051 [PENDENTE])** cliente HTTP em `internal/adapter/out/core/aviso_falha_openobserve.go` (rota/payload/evento **não** inventar) — origem: plan.md § Decisões tomadas etapa 3h / FR-F07 (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B; respostas-etapa3g.md §1; spec.md **FR-F07**)

**Checkpoint**: Domínio e porta declarados (**T002** inclui consulta por achado na interface); chave handoff liberada (batchId / batchId+execId); normalizador **T048** + compositor **T051** liberados (**FR-F11**, **DN-008-5 A**, **008 assinatura A**, **DN-008-9 A**); gating **FR-F07**/**DN-008-6** (**T049**) antes do handoff; campos N3 no payload verificados em T016.

---

## Phase 3: User Story 1 - Handoff Finder → Módulo Ocorrência (Priority: P1) 🎯 MVP

**Goal**: Enviar **ResumoAnalise** idempotente com JWT ao fim da análise; aguardar confirmação; **não** publicar evento LLT-201 (fonte: spec.md, User Story 1; FR-F01, FR-F02).

**Independent Test**: quickstart.md § Cenário 1 — **um** handoff **ResumoAnalise** por lote (**DN-008-7**); consultas **por achado** quando aplicável (**DN-008-7 B**); HTTP+JWT; chave **`batchId`** (v1) ou **`batchId` + `execId`**; confirmação antes de concluir varredura; **0** avisos LLT-201 Finder→Core/stream (Cenário 2 — distinto de **FR-F07** operacional em Cenário 6).

### Implementation for User Story 1

- [X] T017 [P] [US1] Implementar serialização JSON do handoff alinhada a [contracts/resumo-analise-handoff.md](./contracts/resumo-analise-handoff.md) em `internal/adapter/out/modocorrencia/serializar_resumo.go` (**um** payload **ResumoAnalise** por lote — **DN-008-7**) — origem: contracts/ (fonte: respostas-etapa3c.md §6, DN-008-7; respostas-joao.md, N1 A / N4 A / T7 A / DN-006-1b B)
- [X] T018 [US1] Entregar um resumo por lote via HTTPS/JWT em `internal/adapter/out/modocorrencia/cliente.go` conforme DN-009-2, com tratamento de erro, timeout e reenvio com a mesma identidade.
- [X] T019 [US1] Validar recibo de aceitação inicial `201` ou repetida `200` contra tenant/produto/lote/execução antes de concluir. Verificado em `internal/adapter/out/modocorrencia/cliente_test.go` e no teste Redis/HTTPS.
- [X] T020 [US1] Implementar serviço `MontarResumoAnalise` pós-classificação (contadores incl. INFO — T7 A / contador INFO §Etapa 3; `achados[]`; **T048**/**T051** para assinatura derivada quando aplicável — campo struct **T016** bloqueado) em `internal/application/handoff/montar_resumo.go` — origem: research.md R5 (fonte: respostas-joao.md, T7 A / 006-11 A; §Etapa 3c 11h, 008 assinatura A; spec.md **FR-F11**)
- [X] T021 [US1] Popular `versaoRegras` / metadados 006-3 C no resumo em `internal/application/handoff/montar_resumo.go` — origem: data-model.md § ResumoAnalise (fonte: respostas-joao.md, 006-3; revisao-T1-drive.md §4)
- [X] T022 [US1] Orquestrar análise concluída (**T049** / **FR-F07** — **sem** DLQ ingested) → validar lote (T4) → consultas **por achado** via porta (**DN-008-7 B**; HTTP real **T050** bloqueado) quando **N2 B** aplicável → montar resumo → **uma** entrega handoff (**DN-008-7**) em `internal/application/handoff/orquestrar_entrega.go` — origem: plan.md Summary (fonte: spec.md, **FR-F01** / **FR-F07** / **FR-F09**; respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; seq. v6 §2)
- [X] T023 [US1] Popular `execId` em `internal/application/handoff/montar_resumo.go` **somente quando** o Finder dispuser do id de varredura (**DN-006-1**); omitir no consumo atual **log.ingested v1** (**DN-006-1b B**) — origem: FR-F03 (fonte: respostas-joao.md §Etapa 3, DN-006-1 / DN-006-1b B; origem: ADR_075)
- [X] T024 [US1] Implementar estado **aguardando confirmação** antes de concluir varredura em `internal/application/handoff/orquestrar_entrega.go` (equivalente `finishScan` — nome no código a fixar na implementação) — origem: data-model.md § ConfirmacaoHandoff (fonte: spec.md, FR-F05; research.md R4)
- [X] T025 [US1] Garantir ausência de publicação LLT-201 (stream `notificacao-eventos` / aviso ocorrência): revisar `internal/application/` e adapters; **permitido** apenas aviso operacional **FR-F07** (**DN-008-3**/**DN-008-6** rev.) — adaptador HTTP **BLOQUEADO (ADR_051 [PENDENTE])** conforme **T049** — origem: plan.md § Structure Decision (fonte: spec.md, FR-F02 / **FR-F07** / T1 revisto; respostas-joao.md §Etapa 3g, DLQ B; origem: ADR_062)
- [ ] T026 BLOQUEADO (009-22 A) [US1] Fechar agrupamento **entre lotes** (regras **008** / estados atualizada-recorrente) após consultas **DN-008-7 B** em `internal/application/handoff/orquestrar_entrega.go` — depende spec **008** (grão consulta e envio já fixados: **DN-008-7 B** + **DN-008-7**) — origem: spec.md, FR-F09 (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; N2 B / 009-22 A)
- [ ] T027 BLOQUEADO (P-010-27) [US1] Validar cenário quickstart §1 com exemplos de `AchadoAnalise.linhas` mascaradas só após casos do João na spec 010 — não inventar asserts — origem: quickstart.md § Cenário 1 (fonte: respostas-joao.md, 010-27 A)
- [X] T050 [US1] Consultar o estado via HTTPS/JWT por achado elegível em `internal/adapter/out/modocorrencia/cliente.go`. Mapear ausência, revisão e falha conforme DN-008-7 B e DN-009-2, com teste de integração para revisão desatualizada.

**Checkpoint** · DTO, cliente HTTPS e validação de recibo verificados. O aceite do agrupamento (T026), as linhas mascaradas oficiais (T027) e a homologação do Módulo continuam pendentes.

---

## Phase 4: User Story 2 - Pipeline de aviso (referência; fora do Finder) (Priority: P1)

**Goal**: Documentar e rastrear Módulo→Core→Notificação para aceite integrado LLT-201 sem código Finder (fonte: spec.md, User Story 2; FR-001–FR-012).

**Independent Test**: quickstart.md § Cenários 6–9 em ambiente Módulo+Core+Notificação; Finder só produz ResumoAnalise (Cenário 6).

### Documentação e rastreio (somente specs/009)

- [X] T028 [P] [US2] Revisar cobertura de FR-003, FR-004 e **009-5 revisto** (id + tipo + gravidade no HTTP Módulo→Core; campo `ocorrencia` no stream = id) em `specs/009-publicar-evento-ocorrencia/contracts/pipeline-aviso-referencia.md` — origem: plan.md pós-Fase 1 (fonte: respostas-joao.md §Etapa 3, 009-5 revisto; revisao-T1-drive.md, 009-12 rev.)
- [X] T029 [P] [US2] Revisar trecho Core→Notificação (`notificacao-eventos`, ADR_146 v1, 009-10 B) em `specs/009-publicar-evento-ocorrencia/contracts/pipeline-aviso-referencia.md` — origem: research.md R6 (fonte: respostas-documentadas.md, 009-1 B / 009-10 B; origem: ADR_165)
- [X] T030 [US2] Atualizar `specs/009-publicar-evento-ocorrencia/quickstart.md` § Cenário 6 com checklist tenant/produto/id/tipo/gravidade/natureza (SC-001, **009-5 revisto**) — origem: quickstart.md (fonte: spec.md, FR-003 / SC-001)
- [ ] T031 BLOQUEADO (009-22 A) [US2] Manter § Cenário 7 quickstart marcado como aceite **SC-002** bloqueado até spec **008** em `specs/009-publicar-evento-ocorrencia/quickstart.md` — origem: plan.md Summary (fonte: respostas-joao.md, 009-22 A)
- [X] T032 [US2] DOCUMENTAÇÃO: Registrar decisão **DN-009-1 A** (Notificação lê Ocorrência por **id** no MongoDB **ADR_054** para e-mail LLT-218) em `specs/009-publicar-evento-ocorrencia/checklists/requirements.md` — origem: spec.md § Clarifications etapa 3b (fonte: respostas-joao.md §Etapa 3, DN-009-1 A; origem: ADR_054, sequencia-notificacao v2)
- [X] T033 [US2] DOCUMENTAÇÃO: Citar divergência **D1** (épico Jira × handoff Drive) em `specs/009-publicar-evento-ocorrencia/checklists/requirements.md` — origem: spec.md § Divergências (fonte: spec.md, D1; respostas-joao.md, T1 revisto)

**Checkpoint**: Contrato de pipeline referência completo; aceite estados atualizada/recorrente permanece bloqueado (008).

---

## Phase 5: User Story 3 - Idempotência ponta a ponta (Priority: P1)

**Goal**: Reentregas com mesma chave sem duplicar handoff; documentar idempotência HTTP/stream downstream (fonte: spec.md, User Story 3; FR-F06, SC-F002, SC-003).

**Independent Test**: quickstart.md § Cenários 3, 8, 9; US1 Cenário 1 como pré-requisito.

### Implementation for User Story 3 (Finder)

- [X] T034 [US3] Serializar `batchId` e `execId` opcional no corpo, no escopo de tenant e produto, e preservar a chave no reenvio (FR-F03). Verificado em `internal/adapter/out/modocorrencia/cliente_test.go` e no teste Redis/HTTPS.
- [X] T035 [US3] Implementar reenvio com **mesma** chave após 5xx/timeout sem concluir varredura em `internal/application/handoff/reenvio.go` — origem: research.md R10 (fonte: spec.md, FR-F06; revisao-T1-drive.md FA-02)
- [X] T036 [US3] Em falha **handoff Finder→Módulo** (4xx/5xx/timeout sem confirmação): **não** marcar varredura concluída; aplicar **FA-02** (reenvio mesma chave via **T035**) — **009-18 A** aplica-se só ao HTTP **Módulo→Core** (doc **T038**/US2), não ao aviso operacional Finder→Core (**DN-008-3**/**DN-008-6**) — em `internal/application/handoff/orquestrar_entrega.go` — origem: spec.md Clarifications 009-18 (fonte: respostas-joao.md, 009-18 A; revisao-T1-drive.md §009-18; corrige analyze **F1**; spec.md **FR-F05** / **FR-F06**)
- [X] T037 [US3] Fechar semântica de idempotência handoff em `internal/application/handoff/orquestrar_entrega.go` e `internal/domain/idempotencia_handoff.go` alinhada a T015/T023 (**batchId** padrão v1; **execId** opcional) — origem: FR-F03 / FR-F06 (fonte: respostas-joao.md §Etapa 3, DN-006-1 / DN-006-1b B; origem: ADR_014)

### Documentação idempotência produto (referência)

- [X] T038 [P] [US3] Documentar chave HTTP **id ocorrência + estado** e efeito duplicado 009-11 B em `specs/009-publicar-evento-ocorrencia/contracts/pipeline-aviso-referencia.md` § Idempotência — origem: data-model.md § ChaveIdempotenciaAvisoHTTP (fonte: respostas-joao.md, 009-9 A / 009-11 B)
- [X] T039 [P] [US3] Atualizar `specs/009-publicar-evento-ocorrencia/quickstart.md` § Cenários 3, 8 e 9 com passos SC-F002 / SC-003 / 009-10 B — origem: quickstart.md (fonte: spec.md, User Story 3)
- [ ] T040 BLOQUEADO (009-22 A) [US3] Fechar evidência de SC-002 (cenários atualizada/recorrente) no checklist até spec **008** — `specs/009-publicar-evento-ocorrencia/checklists/requirements.md` (fonte: respostas-joao.md, 009-22 A)

**Checkpoint** · A chave no corpo e o reenvio estão verificados no cliente. SC-002 continua condicionado à spec 008 e ao ambiente de produto.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Consistência entre artefatos, constituição e dependências 006/008/010.

- [X] T041 [P] Atualizar matriz de rastreio FR-F01–FR-F11 ↔ tarefas (incl. **FR-F07**/**DN-008-6** → **T049**; **FR-F09**/**DN-008-7 B** → **T002**/**T022**/**T050**; **FR-F11** → **T048**/**T051**) e notar status **Draft**, **SC-002** bloqueado (**009-22 A**), adaptador bloqueado (**DN-009-2**) em `specs/009-publicar-evento-ocorrencia/checklists/requirements.md` — origem: spec.md § Requirements (fonte: plan.md Constitution Check; respostas-joao.md §Etapa 3c 11h; 009-22 A; spec.md **FR-F11**)
- [X] T042 [P] Revisar alinhamento [research.md](./research.md) R1–R16 com decisões etapas 3–3h (**DLQ B**, **DN-008-3**/**DN-008-6** rev., **DN-008-9 A**, DN-006-1/1b, DN-009-1 A, 009-5 revisto, contador INFO, **DN-007-2 B**, **DN-008-4 B**/**DN-008-5 A**, **DN-008-8 A**, **DN-008-7**/**DN-008-7 B**, **008 assinatura A**) — origem: plan.md § Decisões tomadas etapa 3h (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22); respostas-etapa3g.md §1; **DN-008-7 B**: §Etapa 3c (João, 16/11 11h02); respostas-etapa3c.md §5–§6)
- [X] T043 Registrar pendência **emenda constituição I** (Princípio I vs handoff) em `specs/009-publicar-evento-ocorrencia/checklists/requirements.md` — **não** editar `.specify/memory/constitution.md` — origem: plan.md § Pendências (fonte: respostas-joao.md, T1 revisto)
- [X] T044 DOCUMENTAÇÃO: Registrar **C8** (chave concreta ocorrência ADR_040 [PENDENTE]) e 009-6 A em `specs/009-publicar-evento-ocorrencia/data-model.md` / checklist se ainda não rastreado em tasks — origem: spec.md § Questões em aberto (fonte: respostas-documentadas.md, 009-6 A)
- [ ] T045 Executar roteiro manual `specs/009-publicar-evento-ocorrencia/quickstart.md` Cenários 1–5 no ambiente disponível e registrar lacunas em `specs/009-publicar-evento-ocorrencia/checklists/requirements.md` — origem: quickstart.md (fonte: spec.md, SC-F001–SC-F003; Cenário 4 **DN-008-4 B**/**DN-008-5 A**/**008 assinatura A** liberados no desenho via **T048**/**T051**)
- [X] T046 DOCUMENTAÇÃO/PENDÊNCIA: Registrar bloqueio do adaptador HTTP (**DN-009-2**, **ADR_028**) para **POST ResumoAnalise** (**DN-008-7**) e **consulta por achado** (**DN-008-7 B** — **T050**) até contrato publicado fora do repo — `specs/009-publicar-evento-ocorrencia/checklists/requirements.md` (fonte: plan.md § API do Módulo Ocorrência; respostas-joao.md §Etapa 3; **DN-008-7 B**: §Etapa 3c (João, 16/11 11h02))
- [X] T047 [P] DOCUMENTAÇÃO: Registrar **DN-007-2 B** (gate ≥95% precisão **por classe** na spec **006**/**007**, não SC-F00x), base rotulada **JSONL** (007 formato **A**) e caminho documentado **`datasets/labeled-logs/`** na raiz (007 caminho **C** — **não** criar pasta) em `specs/009-publicar-evento-ocorrencia/checklists/requirements.md` — origem: plan.md Summary (fonte: respostas-joao.md §Etapa 3b, DN-007-2 B / 007 formato A / 007 caminho C)

---

## Pendências

Itens externos ou sem decisão que afetam tarefas desta lista (**sem** criar ADR nesta pasta):

| ID | Descrição | Tarefas bloqueadas | Fonte |
|----|-----------|-------------------|-------|
| **DN-009-2** | Contrato HTTP Finder → Módulo (POST **ResumoAnalise** **DN-008-7** + consulta **por achado** **DN-008-7 B**) | T018, T019, T034, **T050** | respostas-joao.md §Etapa 3; **DN-008-7 B**: §Etapa 3c (João, 16/11 11h02); plan.md § API do Módulo Ocorrência |
| **N3 A** | ADR nova — ampliar `AchadoAnalise` | T016 | respostas-joao.md, N3 A |
| ~~**Separador 008 assinatura A**~~ | ~~Literal do separador~~ — fechado **DN-008-8 A**: **`|`** (dono canônico spec **008**) | — | respostas-joao.md §Etapa 3e (João, 16/11 12h40), DN-008-8 A |
| **P-010-27** / **010-27 A** | Casos de teste de mascaramento (João, spec 010) | T027 | respostas-joao.md, 010-27 A |
| **009-22 A** / spec **008** | Aceite SC-002 (atualizada/recorrente, FR-F09) | T026, T031, T040 | respostas-joao.md, 009-22 A |
| **010-2 A** | Lista de padrões de mascaramento em ADR futura | aceite formal de `linhas` (T027) | respostas-joao.md, 010-2 A |
| **006-10 A** | Alinhamento Jira (fora desta pasta) | — (somente doc T005) | respostas-joao.md, 006-10 A |
| **ADR_028** | JWT interno (algoritmo, rotação, env) | paralelo a **DN-009-2** | respostas-etapa3.md §10 |
| **ADR_051 [PENDENTE]** | API **Finder → Core** após **3×** falha consulta **OpenObserve** (**`batchId`**, **motivo**) — rota, payload, evento | cliente HTTP em **T049** (gating **FR-F07** liberado) | respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B; plan.md § Decisões tomadas etapa 3h |
| **Escape `\` na assinatura** | Sem fonte; **DN-008-9 A** cobre só **`|`** → **`\|`** | — (doc **T004**/**T051**; não bloqueia compositor) | respostas-etapa3g.md §2; spec.md §Questões em aberto |

### Desbloqueadas na etapa 3b–3h (sem marcador BLOQUEADO)

Decisões abaixo **liberam** desenho e tarefas que antes dependiam só delas (fonte citada em cada tarefa correspondente):

| ID | Efeito nas tarefas 009 |
|----|------------------------|
| **DN-006-1** / **DN-006-1b B** | T008, T013, T015, T017, T023, T034, T037 — chave `batchId` / `execId` opcional |
| **DN-007-2 B** | T047 — gate métrico fora desta feature; nenhum SC-F bloqueado |
| **DN-008-4 B** | T004, T048, T051, T016 (só estrutura N3), T020, T041, T045 Cenário 4 — **FR-F11** |
| **DN-008-5 A** | T048 (literais `<NUM>`…`<HEX>`), T004, T051, T020 — **FR-F11** |
| **008 assinatura A** / **DN-008-8 A** | T051, T020, T004 — concatenação **DN-008-1 B**; separador **`|`** (dono canônico spec **008**) |
| **DN-008-6** (rev.) | T049 — após 3× falha OpenObserve: Finder→Core API (**batchId**, motivo); **sem DLQ**; ADR_051 [PENDENTE] |
| **DN-008-9 A** | T051 — escape **`|`** → **`\|`** na composição da assinatura |
| **DN-008-7** | T017, T018, T022 — **um** **ResumoAnalise** / **um** POST handoff por lote |
| **DN-008-7 B** | T002, T022 (porta/mock), **T050** (consulta HTTPS verificada); **T026** aceite **008** ainda bloqueada |
| **FR-F07** / **DLQ B** | T049 — gating análise antes do handoff; **sem DLQ** ingested (spec **005** / etapa 3h) |
| **007 formato A / caminho C** | T047 — documentação `datasets/labeled-logs/` (JSONL), sem implementação |
| **DN-009-1 A**, **009-5 revisto**, contador INFO | T009, T020, T028, T032 |
| **DN-008-1 B**, **DN-008-2**, **DN-008-3** (rev.) | coordenação spec **008**; **T026** permanece bloqueada por **009-22 A** |

**Decisões etapa 3b–3h aplicadas (resumo):** DN-006-1, DN-006-1b B, DN-009-1 A, 009-5 revisto, contador INFO (§Etapa 3), **DN-007-2 B**, **DN-008-4 B**, **DN-008-5 A**, **DLQ B**, **DN-008-3**/**DN-008-6** (rev.), **DN-008-9 A**, **DN-008-7**, **DN-008-7 B**, **008 assinatura A**, **DN-008-8 A**, **007 formato A / caminho C** (documentação apenas).

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: sem dependências
- **Foundational (Phase 2)**: depende da Phase 1 — **bloqueia** US1–US3 (**T049** antes de **T022** handoff)
- **US1 (Phase 3)**: depende da Phase 2 — **MVP** (adaptador HTTP após **DN-009-2**)
- **US2 (Phase 4)**: pode paralelizar com US1 após Phase 2 (somente docs em `specs/009-publicar-evento-ocorrencia/`)
- **US3 (Phase 5)**: T035–T037 após orquestração T022; T034 após **DN-009-2**
- **Polish (Phase 6)**: após histórias desejadas

### User Story Dependencies

| História | Depende de | Independência |
|----------|------------|---------------|
| **US1** | spec **006** (classificação); **008** para FR-F09; **DN-009-2** para HTTP real | Testável com mock Módulo + quickstart §1–2 (domínio/serialização) |
| **US2** | Handoff real opcional; **008** para SC-002 | Testável só em ambiente produto |
| **US3** | US1 orquestração; **DN-009-2** para T034; **008** para SC-002 | Quickstart §3 lógico sem HTTP; §7 bloqueado |

### Parallel Opportunities

- Phase 1: T002, T003 em paralelo após T001
- Phase 2: T009, T010, T011, T048, T051 em paralelo
- US1: T017 paralelo a T020/T021
- US2: T028, T029 em paralelo
- US3: T038, T039 em paralelo
- Polish: T041, T042, T047 em paralelo

### Parallel Example: User Story 1

```text
# Após Phase 2:
T017 serializar_resumo.go  ||  T020 montar_resumo.go  ||  T021 versaoRegras
# Após DN-009-2: T050 consultar_ocorrencia (por achado) + T018 entregar_resumo (um por lote) → T019 confirmacao → T022 orquestrar → T024 confirmação antes de varredura
```

---

## Implementation Strategy

### MVP First (User Story 1 apenas)

1. Completar Phase 1 + Phase 2 (T016 verificada; **T048** + **T051** + **T049** **FR-F07**/**DN-008-6**)
2. Completar Phase 3 T017, T020–T025, T023 — handoff lógico (**DN-008-7** um resumo/lote; orquestração consulta **DN-008-7 B** com mock), após **T049**
3. Validar quickstart Cenários 1–2 com mock (sem adaptador real até **DN-009-2**)
4. **PARAR** antes de agrupamento N2 (T026) e aceite 008/010

### Entrega incremental

1. Setup + Foundational → base domínio/porta + idempotência liberada
2. US1 → montagem/orquestração; HTTP quando **DN-009-2** fechar
3. US2 → contrato pipeline e checklists (**009-5 revisto**, DN-009-1 A)
4. US3 → reenvio aplicação + doc idempotência Core/Notificação
5. Polish → quickstart, **DN-007-2 B** / base 007 (T047) e pendências Jira/ADR

### Escopo fora deste repositório (não criar tarefas de código)

- Módulo Ocorrência, Core, Notificação
- ADRs novas (N3 A, 010-2 A) — tarefas **DOCUMENTAÇÃO/PENDÊNCIA**
- Alteração specs 006/008/010 — coordenação via T005 / T006 / T007

---

## Notes

- Proibido publicar em `notificacao-eventos` ou aviso LLT-201 ao Core (**FR-F02**, T1 revisto); exceção operacional **FR-F07** (**DN-008-3**/**DN-008-6** rev., **DLQ B**) — HTTP real **ADR_051 [PENDENTE]** (**T049**)
- `log.ingested` permanece **v1** apenas (T5 D) — não misturar com ADR_146 v2 no ingested (D7)
- INFO em ResumoAnalise (T7 A) apesar de divergência Drive (D3)
- Não alterar ADR_019, ADR_167, SSID_007.1 (FR-F10)
- Formato checklist: todas as tarefas usam `- [ ] Tnnn` com caminho de arquivo na descrição

---

## Resumo de contagem

| Métrica | Valor |
|---------|-------|
| **Total de tarefas** | 52 |
| **Setup** | 7 |
| **Foundational** | 12 |
| **US1** | 12 |
| **US2** | 6 |
| **US3** | 7 |
| **Polish** | 7 |
| **Convergence** | 1 |
| **Marcadas BLOQUEADO** | 5 |
| **Paralelizáveis [P]** | 14 |

**MVP sugerido**: Phase 1 + Phase 2 (inclui **T016** e **T048**, **T051**, **T049**) + US1 T017, T020–T025, T023; validar quickstart §1–2 e §6 (**FR-F07**/**DN-008-6**) com mock; cliente e recibo verificados em **T018**/**T050**/**T019**.

**Critérios de teste independente por história**: spec.md — US1 quickstart §1; US2 §6 (integrado); US3 §3 + §8–9; **FR-F07** quickstart §6 (**T049**).

**Validação de formato** · 52 tarefas no formato Spec Kit. Veja a execução e os bloqueios em [rastreabilidade.md](rastreabilidade.md).

## Ajuste de execução documentado

T048/T051 foram alinhadas ao dono canônico 008, como exigem FR-F11 e a ordem compartilhada; não foi criado algoritmo duplicado. Destinos de notas de execução originalmente em checklists são registrados em rastreabilidade.md (checklists somente leitura pela skill). Itens bloqueados permanecem desmarcados.

## Phase 7: Convergence

- [X] T052 Implementar e testar o efeito de análise que encadeia recuperação de lote completo, mapeamento e entrega antes da confirmação Redis, reutilizando as portas existentes. Fonte: FR-F01, FR-F05, FR-F07, 008 T013 e sequência Varredura/Finder v6 (missing). O bootstrap com serviços reais continua dependente dos respectivos contratos.
