---
description: "Task list for Modelo de classificação de logs"
---

# Tasks: Modelo de classificação de logs

**Input**: Design documents from `/specs/006-modelo-classificacao-logs/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Cenários de teste alinhados ao teste independente de cada história em `spec.md` e às seções de `quickstart.md`. Executar antes da implementação correspondente; devem falhar até o código existir.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions
- **BLOQUEADO (<pendência>)**: na própria linha da tarefa; não implementar até a pendência indicada

## Path Conventions (plan.md § Project Structure)

| Área | Caminho alvo |
|------|----------------|
| Domínio | `internal/domain/` — `classe_log.go`, `classificacao.go`, `achado.go`, `resumo_analise.go`, `assinatura.go`, `versao_regras.go`, `catalogo.go`, `envelope.go` (existente) — **sem** `normalizar_mensagem.go` (dono canônico na 008) |
| Assinatura / normalização (canônico spec **008**) | `internal/application/assinatura/normalizar_mensagem.go` — definido em `specs/008-agrupar-logs-ocorrencia/plan.md` e tasks **T016**; a 006 **só reutiliza** |
| Casos de uso | `internal/application/` — `classificar_log.go`, `analisar_lote.go` |
| Portas entrada | `internal/application/port/in/analisar_lote.go` |
| Portas saída | `internal/application/port/out/enviar_resumo.go`, `internal/application/port/out/assinatura.go`, `internal/application/port/out/consultar_achado_modulo.go` (**DN-008-7 B**) |
| Adapter handoff | `internal/adapter/out/modocorrencia/` — `cliente.go` |
| Configuração | `internal/config/config.go`, `internal/config/env.go`, `.env.example` |
| Testes | Mesmo pacote dos arquivos (`*_test.go`) |

**Fora do escopo desta feature (não acoplar)**: gravar ocorrência em `internal/adapter/out/mongo/` (T1 REVISTO); novo endpoint HTTP público; leitura da fila só para classificar (`internal/adapter/in/queue/redisStream.go` permanece fluxo upstream); **reimplementar** `NormalizarMensagem` fora do pacote canônico 008.

---

## Pendências

Itens de processo ou decisão externa citados nesta lista de tarefas (não criar ADR nem alterar Jira nesta feature):

| ID | Descrição | Fonte | Tarefas bloqueadas |
|----|-----------|-------|-------------------|
| **DN-009-2** | API Módulo Ocorrência: rota, claims JWT, códigos HTTP, 5xx/timeout (envio **e** consulta por achado) | (fonte: respostas-joao.md, Etapa 3 / API Módulo; origem: respostas-etapa3.md §10, ADR_028, ADR_051) | T031, T055; confirmação HTTP real em SC-002; validação quickstart §5 com Módulo real |
| **006-7 D** | Sintaxe e dono do catálogo (ADR_167 / LLT-51) | (fonte: respostas-joao.md, 006-7 D; origem: ADR_167) | T036 |
| **N3 A** | ADR formalizando ampliação `AchadoAnalise` | (fonte: respostas-joao.md, N3 A; origem: ADR_040) | T040 (fechamento formal apenas; não bloqueia runtime provisório) |
| **006-10 A** | Alinhar Jira LLT-197 / LLT-199 | (fonte: respostas-joao.md, 006-10 A) | T041 (ação manual fora do repo) |
| **010-2 A** | Lista oficial de padrões de mascaramento em ADR (spec 010) | (fonte: respostas-joao.md, 010-2 A) | T044 (rastreio documental; implementação mascaramento fora da 006) |
| **Emenda constituição I** | Princípio I vs handoff ResumoAnalise (T1 REVISTO) | (fonte: respostas-joao.md, T1 REVISTO; origem: ADR_062) | T043 (proposta em plan.md; não editar constitution aqui) |
| **ADR_051 [PENDENTE]** | Contrato API interna Finder→Core após 3 falhas na consulta OpenObserve (**batchId** + **motivo**); MUST NOT inventar rota/payload/evento | (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B; origem: respostas-etapa3g.md §1b) | Nenhuma tarefa de adaptador HTTP Core nesta feature (runtime **spec 008**); rastreio **DN-008-3** / **DN-008-6** rev. 3h |
| **Escape `\` na assinatura** | Regra para barra invertida literal nos componentes do `valor` | Sem fonte | Complementa **DN-008-9 A** (só **`|`** fechado — **T027**) |

**Decisões tomadas (etapa 3h)** — substituem alcance anterior onde indicado (fonte: respostas-joao.md §Etapa 3g; respostas-etapa3g.md §1b):

| ID | Decisão | Nota |
|----|---------|------|
| **DLQ B** | Sem DLQ no `log.ingested`; após 3 falhas OpenObserve, Finder avisa o Core (API interna, **batchId** + **motivo**) | Contrato **ADR_051 [PENDENTE]** |
| **DN-008-3** (rev. 3h) | **Substitui** etapas 3b–3e (3 tentativas + DLQ no ingested) | Implementação principal na **spec 008** |
| **DN-008-6** (rev. 3h) | **Substitui** etapa 3c–3e (DLQ + alerta Core/WebSocket pós-DLQ); mesmo gatilho via API interna | Distinto de publicação de ocorrência (**spec 009**) |
| **DN-008-9 A** | Escape **`|`** → **`\|`** nos componentes antes da concatenação | Escape **`\`**: só pendência (sem fonte) |

**Desbloqueados na etapa 3b–3c–3e–3g–3h** (aplicados nas tarefas abaixo): DN-006-1, DN-006-1b B, DN-008-1 B, **DN-007-2 B** (≥95% precisão **por classe**), **007 formato A** (JSONL), **007 caminho C** (`datasets/labeled-logs/` na raiz — documentado, pasta não criada), **DN-008-4 B** (fonte: respostas-joao.md §Etapa 3b João 16/11 09h02; origem: ADR_013; dono impl. **008** — analyze2 **F1**); **DN-006-2 B** (recall e F1 ≥95% **por classe** no gate **T047** — fonte: respostas-joao.md §Etapa 3c, João 16/11 11h; origem: LLT-199, 006-2 B); **DN-006-3 B** (meta **≤ 2 s** classificar+montar resumo por lote em **SC-002** — fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02)); **DN-006-3-ref** (VM-A real, todos os serviços, lote **1.000** registros — fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40)); **DN-008-5 A** (marcadores `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` — **T052**); **DN-008-7 B** (consulta ao Módulo **uma vez por achado**; envio **um** `ResumoAnalise` por lote — **T054**/**T021**, fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); origem: ADR_062, respostas-etapa3c.md §6); **008 assinatura A** + **DN-008-8 A** (separador **`|`**, canônico spec **008**, **T027**); **DN-008-9 A** (escape **`|`** → **`\|`**, **T027** — fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DN-008-9 A); **007 H2 D** (avaliação: classe prevista do apontamento = mais grave entre logs — **T047**).

**Complemento não bloqueante:** **010-27 A** (casos de borda do normalizador na spec **010**) acrescenta cenários a **T052** quando a 010 fechar; literais **DN-008-5 A** já fecham o mínimo normativo.

**Rastreados na 006 sem implementação de runtime** (decisões etapa 3; analyze **I1**): DN-008-2 (sem legado), **DN-008-3** / **DN-008-6** rev. 3h (falha OpenObserve: 3×, **sem DLQ**; aviso Core **batchId** + **motivo** — spec **008**; contrato API **ADR_051 [PENDENTE]**), DN-009-1 A (Notificação lê ocorrência — spec 009), DN-010-1/2 (mascaramento — spec 010; ver T044), 009-5 revisto, contador **INFO**.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Estrutura mínima para handoff HTTP+JWT e versão das regras, sem MongoDB de ocorrência nem LLM (fonte: plan.md § Summary; respostas-joao.md, T1 REVISTO/N1 A)

- [ ] T001 Criar pacote `internal/adapter/out/modocorrencia/doc.go` com comentário de fronteira HTTP+JWT ao Módulo Ocorrência, conforme plan.md § Source Code (fonte: plan.md § Project Structure; respostas-joao.md, N1 A; origem: ADR_051)
- [ ] T002 [P] Em `internal/config/env.go`, adicionar apenas **nomes** de variáveis de ambiente para URL e credenciais JWT do Módulo Ocorrência, sem valores padrão nem inventar host/porta/rota (fonte: research.md R4; respostas-joao.md, N1 A; origem: ADR_051, ADR_028) — valores concretos de path/claims aguardam DN-009-2 somente no adaptador T031
- [ ] T003 [P] Espelhar em `.env.example` os mesmos nomes de T002 com valores vazios e comentário referenciando `specs/006-modelo-classificacao-logs/contracts/resumo-analise-handoff.md` (fonte: plan.md § Technical Context; respostas-joao.md, N1 A)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Tipos de domínio, portas, versão única das regras e validação `log.ingested` v1 — bloqueiam todas as histórias (fonte: plan.md § Project Structure; data-model.md; FR-001, T5 D)

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T004 Em `internal/domain/classe_log.go`, definir `ClasseLog` com valores exclusivos `ERROR`, `WARNING`, `INFO` e ordem de gravidade ERROR > WARNING > INFO (fonte: data-model.md § ClasseLog; respostas-joao.md, 006-6 A/007-9 A; origem: ADR_170, constituição III)
- [ ] T005 [P] Em `internal/domain/classificacao.go`, definir `LogClassificado` com campos obrigatórios `apontamento`, `classe`, `tenantId`, `productId`, `catalogoUsado`, `serviceName`, `host` e opcionais `nivelCampo`, `mensagem` conforme data-model.md (validação: classe uma das três; apontamento não vazio para log contável) (fonte: data-model.md § LogClassificado; spec.md FR-003; DN-008-1 B)
- [ ] T006 [P] Em `internal/domain/achado.go`, definir `Achado` com `assinatura`, `classe`, `logs`, `contagem`, `primeiraOcorrencia`, `ultimaOcorrencia` (fonte: data-model.md § Achado; respostas-joao.md, N3 A/008-18 A; origem: ADR_040) — campos N3 podem usar tipos provisórios até ADR
- [ ] T007 [P] Em `internal/domain/resumo_analise.go`, definir `AchadoAnalise`, `ResumoAnalise`, `IdentificacaoRegras` e contadores ADR_102 revisados: `linhasProcessadas`, `erros`, `warnings`, **`info`** (não `spamIgnorado`), `dadosSensiveisMascarados` (fonte: data-model.md § ContadoresResumo; contracts/resumo-analise-handoff.md; respostas-joao.md, Etapa 3 contadores B/N3 A/006-3 C; origem: ADR_102, ADR_170)
- [ ] T008 Em `internal/domain/versao_regras.go`, definir identificador de versão **única** regras+catálogo+base com `somenteNivel` e `catalogoEntrou` para `identificacaoRegras` (fonte: data-model.md § Identificação das regras; respostas-joao.md, 006-3 C; origem: LLT-199)
- [ ] T009 [P] Em `internal/application/port/in/analisar_lote.go`, declarar porta `Analisar lote` conforme `specs/006-modelo-classificacao-logs/contracts/classificar-apontamento.md` (entrada: `referencia`, `logs`, `versaoRegras`, `catalogo` opcional; saída sucesso/falha com motivos distinguíveis) (fonte: contracts/classificar-apontamento.md; spec.md FR-001, FR-008)
- [ ] T010 [P] Em `internal/application/port/out/enviar_resumo.go`, declarar porta `Enviar` com confirmação explícita e reenvio mesma chave (fonte: contracts/classificar-apontamento.md § Porta de saída; respostas-joao.md, N1 A/N4 A/006-9 A; origem: ADR_014, FA-02)
- [ ] T011 [P] Em `internal/application/port/out/assinatura.go`, declarar porta de cálculo de assinatura no lote com componentes **nível + mensagem normalizada + `service.name` + host** (DN-008-1 B); normalização MUST delegar a `internal/application/assinatura/normalizar_mensagem.go` (DN-008-4 B, dono 008) com marcadores **DN-008-5 A** (`<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>`); valor serializado **008 assinatura A** (texto concatenado — separador **`|`**, **DN-008-8 A**; escape **`|`** → **`\|`**, **DN-008-9 A**, spec **008**); sem `traceId`/`correlationId` (fonte: plan.md § Structure Decision; respostas-joao.md §Etapa 3c, João 16/11 11h, DN-008-1 B/DN-008-4 B/DN-008-5 A/008 assinatura A; §Etapa 3e (João, 16/11 12h40), DN-008-8 A; §Etapa 3g (João, 16/11 13h22), DN-008-9 A; origem: ADR_013, spec.md D5)
- [ ] T012 Em `internal/config/config.go`, carregar versão única das regras (006-3 C) e parâmetros de T002/T003; ausência de variável **não** preenche padrão inventado (fonte: plan.md § Constraints; respostas-joao.md, 006-3 C/006-4 A; origem: LLT-51)
- [ ] T013 Em `internal/domain/envelope.go`, garantir rejeição de `schemaVersion` ≠ `"1"` no fluxo que alimenta análise (fonte: spec.md FR-001/T5 D; research.md R9; `contrato-log-ingested-v1.txt` §2–4)
- [ ] T014 [P] Em `internal/domain/envelope_test.go`, testes para `schemaVersion` `"1"` aceito e ≠ `"1"` rejeitado (fonte: quickstart.md §6; respostas-joao.md, T5 D)
- [ ] T015 Em `internal/domain/resumo_analise.go`, documentar e implementar `execId` **opcional** (`V<uuid>`/`M<uuid>` da Varredura — DN-006-1); chave lógica `batchId` obrigatório; `batchId`+`execId` quando houver; fluxo só `log.ingested` v1 **sem** `execId` no corpo → chave somente `batchId` (DN-006-1b B) (fonte: respostas-joao.md, Etapa 3; origem: contrato-log-ingested-v1.txt, ADR_070, ADR_146)

**Checkpoint**: Foundation ready — User Story 1 pode começar sem agrupamento nem HTTP

---

## Phase 3: User Story 1 - Classificar o log capturado sem serviço de linguagem externo (Priority: P1) 🎯 MVP

**Goal**: Classe por log via campo OTel + palavras inteiras na mensagem; mais grave vence; sem LLM; versão das regras rastreável (fonte: spec.md User Story 1; FR-002–FR-004)

**Independent Test**: Log com nível de erro → classe ERROR; zero envio de texto a serviço externo de linguagem; metadado de versão das regras presente no fluxo de classificação.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T016 [P] [US1] Em `internal/domain/classificacao_test.go`, cobrir tabela de `specs/006-modelo-classificacao-logs/contracts/classificar-log.md` e cenários 1–7 da User Story 1 (ERROR/WARNING/INFO; `errorless`; `traceback`; palavra inteira com delimitadores não alfanuméricos) (fonte: contracts/classificar-log.md; spec.md User Story 1; FR-004)
- [ ] T017 [P] [US1] Em `internal/domain/classificacao_test.go`, provar que nenhum teste aciona rede nem chamada a API de linguagem (fonte: spec.md FR-002; SC-001)

### Implementation for User Story 1

- [ ] T018 [US1] Em `internal/domain/classificacao.go`, implementar operação pura `Classificar log` (entrada/saída de `contracts/classificar-log.md`): mapear termos de nível do campo; extrair palavras inteiras da mensagem; classe = mais grave entre campo e texto; sem catálogo neste método (fonte: spec.md FR-004; respostas-joao.md sessão 2026-11-09 vigente)
- [ ] T019 [US1] Em `internal/application/classificar_log.go`, orquestrar classificação de um log capturado para `LogClassificado`, anexando `tenantId`/`productId` da referência ingested, `serviceName`/`host` OTel e `versaoRegras` de config (fonte: spec.md FR-001/FR-003; data-model.md § LogClassificado; DN-008-1 B)
- [ ] T020 [US1] Em `internal/application/classificar_log_test.go`, caso apontamento sem logs disponíveis: nenhum `LogClassificado`; motivo distinguível `log_ausente` (fonte: spec.md User Story 1 cenário 9; data-model.md § Falhas)

**Checkpoint**: MVP classificação log-a-log; executar `go test ./internal/domain/...` conforme quickstart.md §1

---

## Phase 4: User Story 2 - Produzir achados e enviar ResumoAnalise ao Módulo Ocorrência (Priority: P1)

**Goal**: Agrupar por assinatura no lote; montar `AchadoAnalise` e `ResumoAnalise`; handoff HTTP+JWT (adaptador real pendente DN-009-2); INFO no resumo; lote incompleto não é sucesso; Finder não grava ocorrência nem publica evento (fonte: spec.md User Story 2; FR-007–FR-011)

**Independent Test**: Três logs mesma assinatura ERROR → um `AchadoAnalise` ERROR no resumo; validação em memória do JSON de handoff; zero persistência de ocorrência no Finder.

### Tests for User Story 2 ⚠️

- [ ] T021 [P] [US2] Em `internal/application/analisar_lote_test.go`, três logs mesma assinatura ERROR → um achado ERROR; dois INFO assinaturas distintas → dois achados INFO; lote com **N** achados distintos → **N** invocações da porta de consulta ao Módulo (**DN-008-7 B**, fake) e **um** `ResumoAnalise` montado (fonte: spec.md User Story 2 cenários 1–3 e 7; respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; respostas-etapa3c.md §6; origem: ADR_062)
- [ ] T022 [P] [US2] Em `internal/application/analisar_lote_test.go`, lote incompleto (T4 A): motivo `lote_incompleto`; não tratar como sucesso de lote completo / sem handoff final equivalente (fonte: respostas-joao.md, T4 A; spec.md FR-008)
- [ ] T023 [P] [US2] Em `internal/application/analisar_lote_test.go`, motivos distinguíveis: `log_ausente`, `armazenamento_inacessivel`, `identidade_incompleta` (tenant/produto vazio — 007-8 A), **`assinatura_incompleta`** (log contável sem `service.name` ou host — DN-008-1 B), **`produto_ausente_ou_inconsistente`** (`productId` ausente ou metadata≠corpo ingested v1 — **008-16 A**, **sem retry** de processamento do lote), `envio_modulo_falhou` (corrige analyze2 **E2** no escopo teste) (fonte: data-model.md § Falhas; spec.md FR-008/SC-003; contracts/classificar-apontamento.md; respostas-joao.md, 008-16 A; origem: ADR_049, contrato-log-ingested-v1.txt §4)
- [ ] T024 [P] [US2] Em `internal/application/analisar_lote_test.go`, fake da porta `enviar_resumo.go`: confirmação 2xx simulada; segunda chamada mesma chave `batchId` (ou `batchId`+`execId`) sem efeito colateral no fake (FA-02) — sem depender de T031/DN-009-2 (fonte: contracts/resumo-analise-handoff.md § Idempotência; respostas-joao.md, 006-9 A no módulo)
- [ ] T025 [US2] Em `internal/application/analisar_lote_test.go`, assert de chave de idempotência: somente `batchId` no fluxo sem `execId`; com `execId` injetado, `batchId`+`execId` (fonte: respostas-joao.md, DN-006-1/DN-006-1b B/N4 A; spec.md FR-011)

### Implementation for User Story 2

- [ ] T026 [US2] Em `internal/application/analisar_lote.go`, classificar todos os logs presentes via `classificar_log.go`, particionar em achados por **assinatura** no escopo do lote apenas (classe embutida no componente `nivel` da assinatura — 008-15 A) (fonte: spec.md FR-007/FR-009; respostas-joao.md, 008-15 A/T2 C/DN-008-1 B)
- [ ] T051 [US2] **Pré-requisito:** tarefa **T016** da spec **008** (`specs/008-agrupar-logs-ocorrencia/tasks.md` — implementação canônica de `internal/application/assinatura/normalizar_mensagem.go`; DN-008-4 B). Conectar o fluxo de assinatura da 006 a essa implementação **canônica**: porta T011 e `internal/domain/assinatura.go` MUST **importar/delegar** a esse pacote; **MUST NOT** criar `internal/domain/normalizar_mensagem.go` nem cópia divergente do algoritmo (corrige analyze2 **F1**) (fonte: respostas-joao.md §Etapa 3b, DN-008-4 B; spec.md FR-013; plan.md § Structure Decision)
- [ ] T027 [P] [US2] Em `internal/domain/assinatura.go`, montar valor de assinatura com `nivel`, `mensagemNormalizada` (via delegação T051 → `internal/application/assinatura/normalizar_mensagem.go`, marcadores **DN-008-5 A**), `serviceName`, `host` a partir de `LogClassificado`; serializar **008 assinatura A** como texto dos componentes **concatenados** com separador **`|`** (**DN-008-8 A**, canônico na spec **008**), aplicando escape **`|`** → **`\|`** (**DN-008-9 A**) em cada componente antes da junção; rejeitar agregação e sinalizar **`assinatura_incompleta`** se `serviceName` ou `host` ausentes em log contável; integrar com porta T011 (fonte: data-model.md § Assinatura; respostas-joao.md §Etapa 3c, João 16/11 11h, DN-008-1 B/DN-008-4 B/DN-008-5 A/008 assinatura A/008-18 A; §Etapa 3e (João, 16/11 12h40), DN-008-8 A; §Etapa 3g (João, 16/11 13h22), DN-008-9 A; origem: ADR_013, respostas-etapa3c.md §9)
- [ ] T053 [US2] Em `internal/application/analisar_lote.go` e/ou validação em `internal/domain/envelope.go`, detectar `productId` **ausente** ou **inconsistente** entre metadata validada (token) e corpo do `log.ingested` v1; retornar falha **`produto_ausente_ou_inconsistente`**, falha crítica permanente **sem retry** de processamento do lote (008-16 A, FR-008; distinto de `identidade_incompleta`) (corrige analyze2 **E2** no escopo implementação) (fonte: respostas-joao.md, 008-16 A; origem: contrato-log-ingested-v1.txt §4, ADR_049)
- [ ] T028 [US2] Em `internal/application/analisar_lote.go`, montar `AchadoAnalise` com `tipo` ERROR|WARNING|INFO, **`recorrente=false`** no escopo mínimo (FR-009/N2 B), `hostName` 1:1 ao `host` da assinatura quando presente; campos Drive obrigatórios com semântica mínima de `titulo`/`resumo`/`numeroLinha`/`camposExtraido` **provisória** até **N3 A**; `linhas` recebe trechos **já mascarados** vindos da spec **010** (origem: ADR_040, **010-25**; revisao-T1-drive.md §5), com lista oficial de padrões pendente **010-2 A** (**P-010-2**), conforme `spec.md` Edge Cases, `data-model.md` § AchadoAnalise e **T044**; extensões (`assinatura`, `contagem`, `primeiraOcorrencia`, `ultimaOcorrencia`, `versaoRegras`) (fonte: contracts/resumo-analise-handoff.md; spec.md FR-007 etapa 3c **E3**; respostas-joao.md, N3 A/006-6 A; origem: ADR_040, ADR_063)
- [ ] T029 [US2] Em `internal/application/analisar_lote.go`, agregar contadores ADR_102 no `ResumoAnalise` incluindo contador **`info`** (substitui spam ignorado) a partir dos logs/achados do lote (fonte: contracts/resumo-analise-handoff.md § Contadores; research.md R11; respostas-joao.md, Etapa 3)
- [ ] T030 [US2] Em `internal/application/analisar_lote.go`, popular `batchId` de `domain.Referencia.BatchID`; incluir `execId` no payload e na chave somente quando disponível (DN-006-1/DN-006-1b B) (fonte: data-model.md § ResumoAnalise; respostas-joao.md, Etapa 3; origem: ADR_146)
- [ ] T031 BLOQUEADO (DN-009-2) [US2] Em `internal/adapter/out/modocorrencia/cliente.go`, implementar adapter HTTP+JWT da porta `enviar_resumo.go` (URL, claims, códigos, timeout, retry) serializando JSON de `contracts/resumo-analise-handoff.md` (fonte: respostas-joao.md, N1 A/Etapa 3 API Módulo; origem: ADR_051, respostas-etapa3.md §10)
- [ ] T032 [US2] Em `internal/application/analisar_lote.go`, após confirmação da porta de envio (fake ou T031 quando desbloqueado), sinalizar conclusão da varredura; em falha, não concluir e permitir reenvio mesma chave; montar e validar `ResumoAnalise` em memória independente de T031; documentar no teste/comentário que **SC-002 completo** (handoff HTTP real) exige **DN-009-2** — não confundir aceite em memória/fake com E2E Módulo; meta de tempo **≤ 2 s** (**DN-006-3 B**) mede só classificar+montar resumo no perfil **DN-006-3-ref** (VM-A real, todos os serviços, 1.000 registros) — **exclui** latência de consultas por achado (**DN-008-7 B**) e envio HTTP (**DN-009-2**) (corrige analyze **E1** no escopo tasks) (fonte: revisao-T1-drive.md §2; spec.md User Story 2 cenário 4/SC-002; respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-006-3 B/DN-008-7 B; §Etapa 3e (João, 16/11 12h40), DN-006-3-ref; plan.md § SC-002 e aceite parcial)
- [ ] T033 [US2] Garantir em `internal/application/analisar_lote.go` e adaptadores que **não** há chamada a `internal/adapter/out/mongo/` para ocorrência, publicação de evento Core nem cálculo de gravidade; `tipo` no achado sem gravidade (009-5 revisto no Módulo) (fonte: spec.md FR-007/FR-009/SC-006; respostas-joao.md, T1 REVISTO/Etapa 3)
- [ ] T054 [US2] Em `internal/application/port/out/consultar_achado_modulo.go`, declarar porta de consulta prévia ao Módulo de Ocorrência **por achado** (identidade: assinatura + classe + tenant + produto alinhados a **008-11** / **008-15 A**); em `internal/application/analisar_lote.go`, invocar **uma vez por achado** antes de montar o envio final do lote, com fake nos testes (**DN-008-7 B**); envio permanece **um** `ResumoAnalise` por lote (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; respostas-etapa3c.md §6; origem: ADR_062, N2 B spec 008)
- [ ] T055 BLOQUEADO (DN-009-2) [US2] Em `internal/adapter/out/modocorrencia/cliente.go`, implementar chamada HTTP+JWT da consulta por achado (**DN-008-7 B**) alinhada à mesma pendência de rota/claims/códigos de T031 — sem inventar path (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; plan.md § API do Módulo Ocorrência / DN-009-2)

**Checkpoint**: Resumo + achados validáveis em memória; consulta por achado exercitável com fake; `make test` em `application`; handoff e consulta HTTP reais após DN-009-2 (quickstart.md §5)

---

## Phase 5: User Story 3 - Usar o catálogo só quando o nível não decide (Priority: P2)

**Goal**: Catálogo opcional só quando campo e mensagem não trazem nível reconhecível; primeira entrada válida; default INFO no resumo; feature não cria catálogo (fonte: spec.md User Story 3; FR-005, FR-006)

**Independent Test**: Nível ERROR vence catálogo INFO; sem nível, primeira entrada WARNING; sem catálogo → INFO.

### Tests for User Story 3 ⚠️

- [ ] T034 [P] [US3] Em `internal/domain/classificacao_test.go`, cenários User Story 3 e `contracts/classificar-log.md`: catálogo não vence nível; primeira entrada válida; sem catálogo → INFO com achado elegível ao resumo (fonte: spec.md User Story 3; FR-005/FR-006)
- [ ] T035 [P] [US3] Em `internal/application/analisar_lote_test.go`, `identificacaoRegras.catalogoEntrou` true somente quando catálogo foi fonte decisória de algum log (fonte: data-model.md § Identificação das regras; spec.md FR-007)

### Implementation for User Story 3

- [ ] T036 BLOQUEADO (006-7 D) [US3] Em `internal/domain/catalogo.go`, definir tipo de entrada de catálogo `{padrao, classe}` com sintaxe **aguardando** ADR_167/LLT-51 — até lá, aceitar lista ordenada opaca documentada em comentário (fonte: respostas-joao.md, 006-7 D; origem: ADR_167)
- [ ] T037 [US3] Em `internal/domain/classificacao.go`, estender `Classificar log` para consultar catálogo somente se nível/texto não reconhecidos; ignorar entradas cuja classe ∉ {ERROR,WARNING,INFO}; primeira correspondência na ordem fornecida (fonte: spec.md FR-005; contracts/classificar-log.md)
- [ ] T038 [US3] Em `internal/application/classificar_log.go` e `analisar_lote.go`, propagar catálogo opcional da entrada e atualizar `catalogoUsado` / `identificacaoRegras` (fonte: data-model.md § LogClassificado; respostas-joao.md, 006-3 C)
- [ ] T039 [US3] Em `internal/application/analisar_lote.go`, provar que nenhum código cria ou persiste catálogo (contagem zero de catálogos criados) (fonte: spec.md FR-005/SC-004; respostas-joao.md, 006-7 D)

**Checkpoint**: User Story 3 completa sem fechar sintaxe do catálogo (006-7 D permanece pendência de processo)

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Métricas offline, validação quickstart, pendências documentais e dependências externas (fonte: plan.md § Pendências; spec.md SC-007)

### Documentação / pendências (sem código de produção novo; sem criar ADR)

- [ ] T040 DOCUMENTAÇÃO/PENDÊNCIA: Atualizar `specs/006-modelo-classificacao-logs/contracts/resumo-analise-handoff.md` com matriz de rastreabilidade campo a campo para formalização **N3 A** (ampliação `AchadoAnalise` — assinatura, contagem, primeira/última, versão; marcar `[PENDENTE ADR]`), sem criar arquivo em `adr/` (fonte: respostas-joao.md, N3 A; origem: ADR_040)
- [ ] T041 DOCUMENTAÇÃO/PENDÊNCIA: Produzir nota de alinhamento Jira LLT-197 / LLT-199 com FR-010 e SC-007 em `specs/006-modelo-classificacao-logs/checklists/requirements.md` — **não** alterar issues Jira (006-10 A); manter nota histórica da iteração 2 alinhada a **T1 REVISTO** (analyze-006-etapa3g **F1**) (fonte: respostas-joao.md, 006-10 A/T6 A; T1 REVISTO)
- [ ] T042 DOCUMENTAÇÃO/PENDÊNCIA: Registrar em `specs/006-modelo-classificacao-logs/research.md` reforço **006-7 D** (sintaxe/dono catálogo) e impacto em T036/T037; registrar etapa **3h** (**DLQ B**, **DN-008-3** / **DN-008-6** rev. 3h em §R18 — **substituem** alcance com DLQ/WebSocket pós-DLQ; **DN-008-9 A** em tabela de pendências) (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B / DN-008-9 A; origem: respostas-etapa3g.md §1b, ADR_167)
- [ ] T043 DOCUMENTAÇÃO/PENDÊNCIA: Registrar em `specs/006-modelo-classificacao-logs/plan.md` § Pendências proposta de redação para emenda do princípio I da constituição (handoff ResumoAnalise) — **não** editar `.specify/memory/constitution.md` (fonte: plan.md § Constitution Check; respostas-joao.md, T1 REVISTO)
- [ ] T044 DOCUMENTAÇÃO/PENDÊNCIA: Em `specs/006-modelo-classificacao-logs/data-model.md`, referenciar dependência do campo `AchadoAnalise.linhas` à lista oficial de padrões de mascaramento **010-2 A** (spec 010) — sem implementar padrões nesta feature (fonte: plan.md § Pendências; respostas-joao.md, 010-2 A; origem: ADR_040)
- [ ] T045 DOCUMENTAÇÃO/PENDÊNCIA: Em `specs/006-modelo-classificacao-logs/contracts/resumo-analise-handoff.md` e `quickstart.md` §5, manter referência explícita à pendência **DN-009-2** (sem fechar rota JWT) alinhada ao `plan.md` § API do Módulo Ocorrência; em `quickstart.md` §8/fora de escopo, manter rastreio **DLQ B** / **DN-008-3** / **DN-008-6** rev. 3h e **ADR_051 [PENDENTE]** (sem DLQ no `log.ingested`; aviso Core **batchId** + **motivo**) (fonte: respostas-joao.md, Etapa 3 / API Módulo / §Etapa 3g DLQ B; origem: respostas-etapa3.md §10, respostas-etapa3g.md §1b)

### Implementação / validação transversal

- [ ] T046 [P] Em `specs/006-modelo-classificacao-logs/quickstart.md`, após implementação, conferir passos §1–§6 contra pacotes de teste reais (`make test` / `go test ./internal/...`) sem afirmar resultado não executado (fonte: quickstart.md; plan.md § Testing)
- [ ] T047 [P] Documentar em `specs/006-modelo-classificacao-logs/checklists/requirements.md` procedimento offline de precisão, recall e F1 **por classe** (006-2 B) com base rotulada spec 007 em **JSONL** em `datasets/labeled-logs/` na raiz do repositório (**007 formato A**, **007 caminho C** — só documentado; pasta não criada); comparar ao gold **por apontamento** com classe prevista do apontamento = **mais grave** entre os logs (**007 H2 D**, ADR_170); **assert** gate em **cada** classe ERROR, WARNING, INFO: precisão **≥ 95%** (**DN-007-2 B**), recall **≥ 95%** e F1 **≥ 95%** (**DN-006-2 B**, etapa 3e, SC-007/FR-010) (fonte: spec.md SC-007/FR-010; respostas-joao.md §Etapa 3e (João, 16/11 12h40), 007 H2 D / §Etapa 3c DN-006-2 B / §Etapa 3b DN-007-2 B; origem: LLT-199, ADR_170, respostas-etapa3f.md opção D, plan.md § Technical Context)
- [ ] T048 Garantir que mudança de `versaoRegras` (006-4 A) aplica-se somente a logs novos no wiring de `internal/config/config.go` / `analisar_lote.go` (fonte: respostas-joao.md, 006-4 A; origem: LLT-51)
- [ ] T049 [P] Auditar que `internal/adapter/in/http/health.go` permanece único endpoint HTTP público exposto; classificação não adiciona rota HTTP (fonte: plan.md § Target Platform; spec.md contratos internos)
- [ ] T050 [P] Atualizar `specs/006-modelo-classificacao-logs/contracts/resumo-analise-handoff.md` e, se necessário, comentários em `internal/domain/resumo_analise.go` quanto a `execId` opcional e chaves `batchId` / `batchId`+`execId` (DN-006-1, DN-006-1b B, N4 A) (fonte: respostas-joao.md, Etapa 3; origem: contrato-log-ingested-v1.txt)
- [ ] T052 [P] Em `internal/application/assinatura/normalizar_mensagem_test.go` (pacote canônico spec **008**), testes determinísticos com literais **DN-008-5 A**: `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h, DN-008-5 A; respostas-etapa3c.md §4; spec.md FR-013); acrescentar casos de borda de **010-27 A** quando a spec **010** fechar — a 006 **não** duplica o normalizador (analyze2 **F1**)
- [ ] T056 [P] Em `specs/006-modelo-classificacao-logs/quickstart.md` e/ou checklist, documentar meta normativa **SC-002**: ciclo classificar+montar `ResumoAnalise` **≤ 2 s por lote** (**DN-006-3 B**) no perfil **DN-006-3-ref** (**VM-A real**, **todos os serviços**, lote **1.000** registros); habilitar assert automatizado de tempo conforme quickstart §7b (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), DN-006-3-ref / §Etapa 3c DN-006-3 B; spec.md SC-002)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup — **BLOCKS** all user stories
- **User Story 1 (Phase 3)**: After Foundational — **MVP** independente (sem HTTP nem agrupamento)
- **User Story 2 (Phase 4)**: After Foundational; integra classificação US1; **T051** (delegação ao normalizador canônico 008) com **pré-requisito** **T016** em `specs/008-agrupar-logs-ocorrencia/tasks.md` (implementação canônica `internal/application/assinatura/normalizar_mensagem.go`; fonte: respostas-joao.md §Etapa 3b DN-008-4 B) **antes** de T027; **T053** (008-16 A) pode paralelizar com T026 após envelope; **T054** (consulta por achado, **DN-008-7 B**) após T026/T027 e antes de T032; T031 e T055 bloqueados até **DN-009-2**
- **User Story 3 (Phase 5)**: After US1 domain rules; estende `classificacao.go` / lote; sintaxe catálogo T036 bloqueada em 006-7 D
- **Polish (Phase 6)**: After desired user stories; tarefas DOC podem rodar em paralelo ao fechamento de US2

### User Story Dependencies

| Story | Depende de | Teste independente |
|-------|------------|-------------------|
| US1 (P1) | Fase 2 | Classificação ERROR + sem LLM + versão regras |
| US2 (P1) | Fase 2 + classificação US1 | Resumo + achados + porta fake; sem Mongo ocorrência |
| US3 (P2) | US1 regras de nível | Precedência catálogo; não cria catálogo |

### Bloqueios explícitos (tarefas)

| Pendência | Tarefas |
|-----------|---------|
| **DN-009-2** | T031, T055 |
| **006-7 D** | T036 |

### Parallel Opportunities

- Fase 1: T002 ∥ T003 após T001
- Fase 2: T005–T007, T009–T011, T014 em paralelo após T004
- US1: T016 ∥ T017; depois T018–T020 sequencial
- US2: T021–T025 em paralelo; **T051 após T026, antes de T027**; T053 ∥ T028 após T013; T027 ∥ T028 após T051; **T054** após T027, antes de T032
- US3: T034 ∥ T035
- Polish: T040–T045, T046, T049–T050, **T052**, T056 em paralelo; **T047** com gate precisão/recall/F1 **007 H2 D** (**DN-007-2 B** + **DN-006-2 B**)

---

## Parallel Example: User Story 1

```text
# Testes em paralelo:
T016 domain/classificacao_test.go — tabela classificar-log.md
T017 domain/classificacao_test.go — zero rede/LLM

# Depois:
T018 domain/classificacao.go — regra pura
T019 application/classificar_log.go — orquestração
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup  
2. Complete Phase 2: Foundational (CRITICAL)  
3. Complete Phase 3: User Story 1  
4. **STOP and VALIDATE**: `go test ./internal/domain/...` + cenários quickstart §1  
5. Demo: um log → ERROR/WARNING/INFO sem LLM  

### Incremental Delivery

1. Setup + Foundational → base de tipos e portas  
2. US1 → classificação log-a-log (MVP)  
3. US2 → achados + `ResumoAnalise` em memória + fake envio; HTTP real após DN-009-2  
4. US3 → catálogo opcional (com 006-7 D ainda aberto)  
5. Polish → DOC pendências, quickstart, métricas 006-2 B offline  

### Parallel Team Strategy

1. Time fecha Fase 2 junto  
2. Dev A: US1 (`domain/classificacao.go`)  
3. Dev B: US2 — T051 delegação a `internal/application/assinatura/normalizar_mensagem.go` (DN-008-4 B, dono 008) **antes** de T027 `assinatura.go`; T053 (008-16 A)  
4. Dev C: US3 após T018  
5. Desbloquear **T031** e **T055** quando **DN-009-2** fechar; benchmark **≤ 2 s** (**T056**) no perfil **DN-006-3-ref**  

---

## Notes

- [P] = arquivos distintos, sem dependência de tarefa incompleta na mesma fase  
- Cada tarefa de implementação cita fonte João/spec/plan para rastreio  
- Não commitar segredos JWT; não gravar texto integral de log em handoff  
- `internal/adapter/out/mongo/doc.go` permanece fora do fluxo de ocorrência (T1 REVISTO)  
- Agrupamento **entre lotes** e recorrência OpenObserve (N2 B): spec **008**, fora do mínimo **FR-009**; após **3 falhas** na consulta OpenObserve, Finder avisa o **Core** pela API interna (**batchId** + **motivo**) — **sem** DLQ no `log.ingested` e **sem** WebSocket pós-DLQ (**DLQ B**, **DN-008-3** / **DN-008-6** rev. 3h); adaptador HTTP Core **BLOQUEADO** até **ADR_051 [PENDENTE]** (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B; origem: respostas-etapa3g.md §1b, contrato-log-ingested-v1.txt §8)  
- **Consulta prévia ao Módulo por achado** (**DN-008-7 B**) é escopo **FR-007** (envio **um** `ResumoAnalise` por lote — fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; respostas-etapa3c.md §6; ADR_062)  
- Literais de marcador **DN-008-5 A** e testes **T052** no pacote canônico 008; **010-27 A** complementa bordas na 010  
- Serialização da assinatura: escape **`|`** → **`\|`** (**DN-008-9 A**, **T027**); escape de **`\`** literal: pendência sem fonte  
- Normalizador: **uma** implementação em `internal/application/assinatura/` (spec 008); 006 só delega (analyze2 **F1**)