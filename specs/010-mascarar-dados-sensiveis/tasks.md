# Tasks: Mascarar dados sensíveis antes de gravar a ocorrência

**Input**: Design documents from `specs/010-mascarar-dados-sensiveis/`

**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md)

**Tests**: Casos de aceite **não** são inventados nesta feature — **P-010-27** (010-27 A). Tarefas de aceite formal usam `BLOQUEADO (P-010-27)`.

**Organization**: Fases Setup → Foundational → user stories (P1, P1, P2) → Polish. Cada tarefa cita origem (FR/US, plan, decisão com fonte).

**Regenerado**: 2026-11-16 (etapa **3e** — **DN-010-3 B** desbloqueia T046; **DN-008-5 A** desbloqueia T043; **DN-010-4 C** realinha T013 para **P-010-2**; pendências/DN respondidas 3b/3c/3e; etapa **3f** — fontes cruzadas 006/007/008, sem divergência 11h×11h01 neste arquivo; etapa **3h** — **010 H1**, **DN-008-9 A**)

**Etapa 3f**: Rodadas §Etapa 3c 11h01 e 11h **não** divergem — «puladas» **SUPERADO** (`respostas-joao.md` l.39). **DN-006-3 B**, **DN-007-4 A** e **DN-008-7 B** (referência cruzada specs 006/007/008) → fonte `respostas-joao.md §Etapa 3c (João, 16/11 11h02)`. Divergências João×Drive/Jira: [spec.md](./spec.md) § Divergências documentadas — **sem** tabela 11h×11h01 nem **D13** neste arquivo.

**Etapa 3h**: **010 H1** — a **010 reaproveita** o normalizador canônico da **008** (fonte canônica **João 16/11 13h22**, `respostas-joao.md §Etapa 3g`; substitui orientação do Gestor LhamaLog 16/11 12h37 sujeita à aprovação do João). **DN-008-9 A** respondida — escape **`|`** → **`\|`** nos componentes da assinatura é dono spec **008**; **sem** **BLOQUEADO (DN-008-9)** nesta feature. **DLQ B** / **DN-008-3** / **DN-008-6** rev. 3h **não** citados neste arquivo (runtime specs **008**/**009**). (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22); origem: 010 H1, DN-008-9 A, DLQ B; respostas-etapa3g.md §1)

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Paralelizável (arquivos distintos, sem dependência de tarefas incompletas)
- **[US#]**: História em fases de user story
- **BLOQUEADO (…)** na linha: apenas tarefas com pendência externa ou DN aberta (formato exato exigido pelo João)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Preparar rastreio de pendências e alinhamento documental da feature (somente escopo 010).

**Independent Test**: `plan.md` § Pendências e `spec.md` § Pendências explícitas refletem **P-010-2** (regex + metadados catálogo **DN-010-4 C**), P-010-27, N3 A, DN-009-2; DN-010-1, DN-010-2, DN-007-2, DN-008-4, **DN-010-3 B**, **DN-008-5 A**, **DN-010-4 C**, **010 H1** e 007 formato/caminho constam como **respondidas** (etapas 3b/3c/3e/3h).

- [X] T001 Confirmar pré-requisitos de documentação listados no cabeçalho deste arquivo contra `specs/010-mascarar-dados-sensiveis/plan.md` § Project Structure (origem: plan.md § Project Structure)
- [X] T002 [P] DOCUMENTAÇÃO/PENDÊNCIA — Rastrear criação ou emenda de **ADR** para lista oficial de regex do catálogo (**P-010-2**, 010-2 A); referenciar `specs/010-mascarar-dados-sensiveis/contracts/catalogo-padroes-versionado.md` e **não** criar ADR neste repositório (fonte: respostas-joao.md, 010-2 A; origem: ADR_064 [PENDENTE]; FR-001, FR-008)
- [X] T003 [P] DOCUMENTAÇÃO/PENDÊNCIA — Rastrear emenda de **ADR** para ampliação de `AchadoAnalise` (campos de agrupamento, **N3 A**); alinhar expectativa com `specs/010-mascarar-dados-sensiveis/contracts/saida-mascarada-handoff.md` sem editar ADR aqui (fonte: respostas-joao.md, N3 A; origem: backlog 006/008)
- [X] T004 [P] DOCUMENTAÇÃO/PENDÊNCIA — Registrar coordenação **006-10 A** (alinhar Jira **LLT-202**/epic junto da atualização da spec `006-modelo-classificacao-logs`) como ação **fora** desta pasta; vínculo mascaramento ↔ classificação via `mensagem` mascarada (fonte: respostas-joao.md, 006-10 A; spec.md T6 A)
- [X] T005 [P] DOCUMENTAÇÃO/PENDÊNCIA — Manter registro explícito **P-010-27** / 010-27 A (João define casos SC-001/002) em `specs/010-mascarar-dados-sensiveis/quickstart.md` e `spec.md` SC-005 — não listar casos de teste (fonte: respostas-joao.md, 010-27 A)
- [X] T006 [P] DOCUMENTAÇÃO — Confirmar que `specs/010-mascarar-dados-sensiveis/contracts/saida-mascarada-handoff.md` § Downstream documenta **DN-010-1 A** (Notificação lê Ocorrência pelo id; `conteudo` mascarado) — **desbloqueado** etapa 3b (fonte: respostas-joao.md, DN-009-1 A; origem: respostas-etapa3.md §7, LLT-218; FR-004)
- [X] T007 [P] DOCUMENTAÇÃO — Confirmar que `specs/010-mascarar-dados-sensiveis/contracts/validar-seguranca.md` e `data-model.md` refletem **DN-010-2** (mascaramento **só linhas**; identificadores ADR_040 fora do catálogo linha) — **desbloqueado** etapa 3b (fonte: respostas-joao.md, DN-010-2; origem: respostas-etapa3.md §8; FR-004a)

**Checkpoint**: Pendências ADR/Jira registradas; decisões 3b/3c refletidas nos contratos; implementação pode prosseguir no núcleo `validarSeguranca()`.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Tipos de domínio, porta do catálogo, `validarSeguranca()` e `montarContexto()` — pré-requisito de **todas** as histórias (origem: plan.md § Summary, § Project Structure; FR-002, FR-007).

**⚠️ CRITICAL**: Nenhuma user story fecha sem esta fase.

**Independent Test**: Lote processado executa **uma** chamada a `validarSeguranca()` em `montarContexto()` antes de qualquer passo de classificação (quickstart Cenário 1; 010-15 A).

- [X] T008 Definir constante de negócio `MarcadorOmisso` = `"log com senha omitido"` em `internal/domain/` conforme `specs/010-mascarar-dados-sensiveis/data-model.md` (fonte: respostas-joao.md, 010-10 A; FR-006)
- [X] T009 [P] Implementar struct/tipo `LinhaLogMascarada` em `internal/domain/` (campos obrigatórios `texto`; recomendados `indiceLinha`, `linhaRemovida`; **sem** mascaramento de atributos OTel — DN-010-2) (origem: data-model.md LinhaLogMascarada; 010-14 C; FR-004a)
- [X] T010 [P] Implementar struct/tipo `ResultadoValidarSeguranca` em `internal/domain/` (`linhasMascaradas`, `totalSensivel`, `erro` opcional) (origem: data-model.md; FR-007)
- [X] T011 [P] Implementar struct/tipo `ContextoLoteMascarado` em `internal/domain/` (`batchId`, `linhas`, `versaoCatalogo`, `dadosSensiveisIdentificados`, `falhaSeguranca`) (origem: data-model.md; 010-15 A, 010-17)
- [X] T012 [P] Declarar interface de saída `CatalogoPadroes` (leitura do artefato versionado) em `internal/application/port/out/` alinhada a `specs/010-mascarar-dados-sensiveis/contracts/catalogo-padroes-versionado.md` (fonte: respostas-joao.md, 010-18 A; FR-008)
- [ ] T013 BLOQUEADO (P-010-2) Implementar adaptador que satisfaz `CatalogoPadroes` lendo artefato versionado — **caminho** e **formato** MUST ser os declarados na **ADR da lista** quando **P-010-2** publicar metadados (**DN-010-4 C**, **010-2 A**); até lá **não** inventar path nem serialização — em pacote `internal/adapter/` conforme layout README (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h, DN-010-4 C; origem: respostas-etapa3c.md §8, ADR_079, 006-3 C; 010-18 A; FR-008; constituição V)
- [ ] T014 Implementar função/método `validarSeguranca()` em `internal/application/` conforme `specs/010-mascarar-dados-sensiveis/contracts/validar-seguranca.md` (entrada `linhasBrutas`, saída `ResultadoValidarSeguranca`; regras interim ADR_064; escopo linha — DN-010-2) (fonte: respostas-documentadas.md, 010-3 A / 010-4 A; FR-002, FR-006, FR-004a)
- [X] T015 Implementar `montarContexto()` em `internal/application/` invocando `validarSeguranca()` **exatamente uma vez** por lote antes do MotorAnalise/classificação (fonte: respostas-documentadas.md, 010-4 A / 010-15 A; plan.md § Summary)
- [X] T016 Implementar política **010-16 C** em `internal/application/`: erro de `validarSeguranca()` → continuar lote **sem** trechos em achados (`falhaSeguranca`, `linhas` vazias ou omitidas) (fonte: respostas-joao.md, 010-16 C; FR-010)
- [ ] T017 BLOQUEADO (P-010-2) Ligar regras `padroes[]` do catálogo a `validarSeguranca()` em `internal/application/` quando ADR publicar lista — até lá, comportamento interim só para e-mail e senha/token/secret **sem** afirmar regex completas (fonte: respostas-joao.md, 010-2 A; research.md R1; FR-001)

**Checkpoint**: Pipeline chama `montarContexto()` → `validarSeguranca()`; contador `totalSensivel` disponível no domínio.

---

## Phase 3: User Story 1 — Trechos guardados sem segredos e PII expostos (Priority: P1) 🎯 MVP

**Goal**: `AchadoAnalise.linhas` e persistência downstream em `conteudo` recebem **apenas** exemplos mascarados para padrões **no catálogo** (fonte: spec.md US1; FR-003, FR-006).

**Independent Test**: Após **P-010-27** + ADR **P-010-2**, processar casos do João e verificar ausência de segredo/PII coberto no formato original em `linhas`/`conteudo` (SC-001).

### Implementation for User Story 1

- [X] T018 [US1] Garantir que montagem de achados use **somente** saída de `validarSeguranca()` para preencher `AchadoAnalise.linhas` em `internal/application/` (fonte: respostas-documentadas.md, 010-6; revisao-T1-drive.md, 010-25; FR-003)
- [ ] T019 [P] [US1] Aplicar substituição de **e-mail** por placeholder fixo em `validarSeguranca()` em `internal/application/` (interim até ADR; valor canônico em `data-model.md` PlaceholderEmail) (fonte: respostas-joao.md, 010-11 A; FR-006; pendência canônica **P-010-2**)
- [ ] T020 [P] [US1] Aplicar troca de linha inteira por `MarcadorOmisso` quando padrão senha/token/secret catalogado em `validarSeguranca()` em `internal/application/` (fonte: respostas-joao.md, 010-10 A; FR-006)
- [ ] T021 [US1] Preservar texto **inalterado** para dado sensível **fora** do catálogo em `validarSeguranca()` em `internal/application/` (fonte: respostas-documentadas.md, 010-13 A; FR-009)
- [X] T022 [US1] Proibir uso de linhas brutas OpenObserve na montagem de `linhas` após falha ou bypass de `validarSeguranca()` em `internal/application/` (fonte: revisao-T1-drive.md, 010-21; 010-5 A)
- [ ] T023 [US1] BLOQUEADO (P-010-27) Executar validação de aceite SC-001 / SC-003 conforme `specs/010-mascarar-dados-sensiveis/quickstart.md` Cenário 2 (fonte: respostas-joao.md, 010-27 A; spec.md SC-001, SC-003, SC-005)

**Checkpoint**: Com catálogo interim, linhas mascaradas seguem FR-006 para comportamento documentado em ADR_064; aceite formal aguarda João + ADR.

---

## Phase 4: User Story 2 — Consumo externo só vê texto mascarado (Priority: P1)

**Goal**: Finder **produz** texto mascarado; consumidores não recebem bruto catalogado via handoff Finder (fonte: spec.md US2; FR-004, FR-011; 010-20, 010-21).

**Independent Test**: Payload `ResumoAnalise` auditado — sem bruto; classificador só vê mascarado; paridade `versao` catálogo quando implementado (quickstart Cenários 1, 4, 6, 8).

### Implementation for User Story 2

- [X] T024 [US2] Assegurar que MotorAnalise/classificador em `internal/application/` consuma **exclusivamente** linhas pós-`montarContexto()` (fonte: respostas-documentadas.md, 010-5 A; FR-002)
- [X] T025 [P] [US2] Propagar `versaoCatalogo` do lote em `ContextoLoteMascarado` e metadados de handoff interno em `internal/domain/` / `internal/application/` (fonte: respostas-joao.md, 010-19 A; FR-008; SC-006)
- [X] T026 [US2] Montar `ResumoAnalise` em `internal/application/` **sem** campos com conteúdo integral bruto do OpenObserve (fonte: respostas-documentadas.md, 010-24 C; FR-011; contrato-log-ingested-v1.txt §2)
- [ ] T027 [US2] Alinhar payload interno ao Módulo Ocorrência com `specs/010-mascarar-dados-sensiveis/contracts/saida-mascarada-handoff.md` (`achados[].linhas` mascaradas; contadores; sem gravidade Finder) (fonte: revisao-T1-drive.md, 010-21 / 010-23; T1 revisto)
- [X] T028 [P] [US2] DOCUMENTAÇÃO — Garantir em `specs/010-mascarar-dados-sensiveis/research.md` § R7 que o Finder **não** implementa aviso HTTP Módulo→Core; obrigação **009-5 revisto** / **010-23** (id + tipo + gravidade, **sem** texto de log) é do **Módulo de Ocorrência** (spec **009**). T028 descreve **garantia de saída Finder** (handoff mascarado), não o evento downstream (fonte: respostas-joao.md, Etapa 3 / clarify I1; origem: respostas-etapa3.md §7; FR-004; analyze **F1**)
- [X] T029 [US2] DOCUMENTAÇÃO — Registrar em `specs/010-mascarar-dados-sensiveis/research.md` § R7 que critério LLT-218 no Finder limita-se a `ResumoAnalise`/`linhas` mascarados; e-mail e tela consomem Ocorrência/Notificação **fora** deste repo (**DN-010-1 A**) — sem implementar Notificação aqui (fonte: respostas-joao.md, DN-009-1 A; 010-22 A; FR-004; analyze **F1**)
- [ ] T030 [US2] BLOQUEADO (P-010-27) Validar SC-002 (tela, e-mail, Sugestão/LLM) conforme quickstart após casos do João (fonte: respostas-joao.md, 010-27 A; spec.md SC-002, SC-005)
- [X] T041 [US2] Garantir em `internal/application/` que `validarSeguranca()`/`montarContexto()` concluam **antes** de classificação e de qualquer uso de `mensagem` derivada do log para assinatura/normalização; entrada MUST ser linha **mascarada** (fonte: respostas-joao.md §Etapa 3b, DN-008-4 B; origem: spec.md FR-015, plan.md § Summary; quickstart Cenário 8)
- [X] T043 [US2] Ao compor `assinatura` em `internal/application/`, **REUTILIZAR/delegar** ao normalizador canônico `internal/application/assinatura/normalizar_mensagem.go` (pré-requisito `specs/008-agrupar-logs-ocorrencia/tasks.md` **T016**; algoritmo **DN-008-4 B**; literais **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`** — **DN-008-5 A**) sobre texto de linha **já mascarado** — **MUST NOT** criar `normalizar_mensagem` próprio nem reimplementar **DN-008-4 B** nesta feature. **010 H1**: a **010 reaproveita** o normalizador canônico da **008** (fonte: respostas-joao.md §Etapa 3g, João 16/11 13h22, 010 H1); normalizador único para specs **006**, **008**, **009** e **010** (**DN-008-4 B**). Composição textual **008 assinatura A** e escape **DN-008-9 A** (**`|`** → **`\|`** nos componentes) — dona spec **008** (`specs/008-agrupar-logs-ocorrencia/tasks.md` **T051**); **010** só integra normalização de mensagem — **sem** **BLOQUEADO (DN-008-9)** (fonte: respostas-joao.md, DN-008-9 A; origem: §Etapa 3g, João 16/11 13h22). (fonte: respostas-joao.md §Etapa 3b, DN-008-4 B; §Etapa 3c, DN-008-5 A; origem: respostas-etapa3c.md §4, ADR_013; FR-015; analyze etapa 3f **H1**)
- [ ] T046 [US2] Implementar política **DN-010-3 B** no bootstrap/leitura do catálogo em `internal/application/`: se versão lida do artefato ≠ versão **combinada** (config/deploy; **010-19 A**), o serviço MUST **recusar subir** na inicialização (falha observável em log/readiness); documentar em `specs/010-mascarar-dados-sensiveis/contracts/catalogo-padroes-versionado.md` (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h, DN-010-3 B; origem: respostas-etapa3c.md §7, ADR_079, 010-19 A; FR-008, SC-006)

**Checkpoint**: Handoff Finder cumpre FR-004/FR-011; ordem mascaramento → normalização (T041); integração assinatura com literais **DN-008-5 A** (T043) e fail-fast de versão **DN-010-3 B** (T046) **desbloqueados** etapa 3e; **010 H1** e **DN-008-9 A** (composição na **008**) fechados etapa 3h em **T043**; T013 aguarda P-010-2 e T042 aguarda a homologação do cliente DN-009-2 com o catálogo e o Módulo reais.

---

## Phase 5: User Story 3 — Rastreio de dados sensíveis identificados (Priority: P2)

**Goal**: Total **`dadosSensiveisIdentificados`** nos contadores do `ResumoAnalise` e resumo do arquivo, com rótulo humano **“dados sensíveis mascarados”** (fonte: spec.md US3; FR-005; revisao-T1-drive.md, 010-17; clarify **I2**).

**Independent Test**: Lote com contagem conhecida (após P-010-27) — contador coincide (SC-004).

### Implementation for User Story 3

- [ ] T031 [US3] Mapear `ResultadoValidarSeguranca.totalSensivel` → contador **`dadosSensiveisIdentificados`** no `ResumoAnalise` e resumo do arquivo, com rótulo humano **“dados sensíveis mascarados”**; em `contracts/saida-mascarada-handoff.md` usar chave JSON inferida `contadores.dadosSensiveisMascarados` **somente** como rótulo de handoff até **N3 A**/**ADR_102** (fonte: revisao-T1-drive.md, 010-17; spec.md clarify I2; FR-005; ADR_102)
- [X] T032 [US3] Incluir `dadosSensiveisIdentificados` em `ContextoLoteMascarado` ao concluir lote em `internal/application/` (origem: data-model.md; 010-17)
- [ ] T033 [US3] BLOQUEADO (P-010-27) Validar SC-004 via `specs/010-mascarar-dados-sensiveis/quickstart.md` Cenário 3 (fonte: respostas-joao.md, 010-27 A; spec.md SC-004)

**Checkpoint**: Contador exposto no handoff conforme `saida-mascarada-handoff.md`; aceite numérico aguarda casos João.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Integrações 003/006/009, checklist, quickstart e aceite formal.

- [X] T034 [P] Verificar que captura de log (spec 003) permanece registro integral e mascaramento só na etapa de processamento — revisar integração em `internal/application/consumir.go` sem contradizer FR-012 (fonte: spec.md FR-012; research.md R9)
- [X] T035 Confirmar que implementação **não** altera ADR_019, ADR_167 nem SSID_007.1 por escrita de código ou docs fora da pasta 010 (fonte: spec.md FR-013; research.md R10)
- [ ] T036 [P] Executar cenários **não bloqueados** de `specs/010-mascarar-dados-sensiveis/quickstart.md` (Cenários 1, 4, 5, 6, 7, 8, 9) e registrar resultado em notas de implementação na pasta 010 (010-28 A: sem meta de desempenho)
- [X] T037 [P] Revisar checklist `specs/010-mascarar-dados-sensiveis/checklists/requirements.md` contra FR-001–FR-015 após implementação; marcar exceções explícitas para **SC-005**, **P-010-27** e **P-010-2** (regex + metadados catálogo **DN-010-4 C**) onde “testable and unambiguous” não se aplica até a ADR fechar (fonte: spec.md SC-005; DN-010-3 B e DN-010-4 C respondidos etapa 3e)
- [X] T038 DOCUMENTAÇÃO/PENDÊNCIA — Rastrear emenda do **Princípio I** da constituição (Finder vs Módulo Ocorrência, T1 revisto) **fora** deste repositório; não editar `.specify/memory/constitution.md` (origem: plan.md § Pendências, P-emenda-I; analyze **C1**)
- [ ] T039 DOCUMENTAÇÃO/PENDÊNCIA — Confirmar paridade Finder–Sugestão do artefato versionado (010-18 A, 010-19 A) com time Sugestão (Java, fora do repo); registrar em `specs/010-mascarar-dados-sensiveis/research.md` § R6 se necessário
- [ ] T040 BLOQUEADO (P-010-2; P-010-27) Fechar aceite formal **SC-005** (SC-001/002 completos) — não marcar feature “done” em produto até pendências João/ADR; handoff de ponta a ponta condicionado à cadeia specs **008→009** (**009-22 A**: aceite da **009** só fecha após **008** — sem inventar aceite local da **009**) (fonte: spec.md SC-005; respostas-joao.md, 009-22 A; origem: Rodada pós-revisão T1, 16/11 00h01)
- [ ] T042 Parcial. O cliente HTTPS/JWT foi implementado na spec 009 e testado na composição. Validar o handoff mascarado com o catálogo oficial, os casos P-010-27 e o Módulo real após a spec 008, conforme 009-22 A.
- [X] T044 [P] DOCUMENTAÇÃO — Manter `specs/010-mascarar-dados-sensiveis/plan.md` § Referência cruzada — base rotulada e `research.md` § R12 alinhados a **DN-007-2 B**, **007 formato A** (JSONL) e **007 caminho C** (`datasets/labeled-logs/` na raiz do repositório; **não** criar pasta nesta feature) (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; spec.md SC-007)
- [X] T045 [P] DOCUMENTAÇÃO — Verificar que implementação e documentação na pasta 010 **não** expandem escopo **LLT-51** (minimização ampla); registrar verificação em `specs/010-mascarar-dados-sensiveis/research.md` § R8 e nota em `quickstart.md` (fonte: respostas-documentadas.md, 010-26 A; FR-014; analyze **E1**)

**Checkpoint**: Núcleo de mascaramento e payload mascarado prontos; integração HTTP e aceite formal dependem de pendências listadas abaixo.

---

## Pendências

### Decisões tomadas (etapa 3h)

| ID | Decisão | Fonte |
|----|---------|-------|
| **010 H1** | A **010 reaproveita** o normalizador canônico da **008** — **MUST NOT** criar `normalizar_mensagem` próprio (**T043**, **FR-015**) | (fonte: respostas-joao.md, 010 H1; origem: §Etapa 3g, João 16/11 13h22, DN-008-4 B) |
| **DN-008-9 A** | **`|`** literal em componente da assinatura MUST ser escapado como **`\|`** antes da concatenação (dono spec **008**, **T051**). Escape de **`\`** literal: só pendência (sem fonte) — ver tabela abaixo | (fonte: respostas-joao.md, DN-008-9 A; origem: §Etapa 3g, João 16/11 13h22) |

**Sem alteração de runtime nesta feature (etapa 3h):** **DLQ B**, **DN-008-3** e **DN-008-6** rev. 3h (após **3 falhas** na consulta OpenObserve: Finder avisa o **Core** por API interna com **batchId** e **motivo**; **sem** DLQ no **log.ingested v1** e **sem** WebSocket pós-DLQ — substitui alcance anterior com DLQ) — implementação e contrato **ADR_051 [PENDENTE]** nas specs **008**/**009**; esta lista **não** cita DLQ nem bloqueia tarefas por **DN-008-9**. (fonte: respostas-joao.md, DLQ B; origem: respostas-etapa3g.md §1, contrato log.ingested v1, DN-006-1b B)

### Pendências abertas

| ID | Descrição | Efeito nas tarefas |
|----|-----------|-------------------|
| **P-010-2** / 010-2 A | Lista oficial de regex e metadados do artefato (**caminho**/**formato** — **DN-010-4 C**) na **ADR** | T002, T013, T017, T019, T040; conteúdo de `padroes[]` |
| **P-010-27** / 010-27 A | Casos de teste João (SC-001/002) | T023, T030, T033, T040 |
| **DN-009-2** | Contrato HTTPS aprovado; cliente implementado na 009 | T042 aguarda integração e aceite com serviços reais |
| **009-22 A** / spec **008** | Aceite da spec **009** só fecha após **008** (**009 espera a 008**) | Cadeia **008→009** condiciona handoff HTTP (T042) e entrega canônica na **008** (T043 integra sem inventar algoritmo); não inventar aceite da **009** na **010** (fonte: respostas-joao.md, 009-22 A; origem: perguntas-etapa2b.md 009-22 A) |
| **DN-007-2 B** / 007 formato A / 007 caminho C | Meta ≥95% por classe; base JSONL em `datasets/labeled-logs/` | **Sem** tarefa bloqueada além da doc T044 (gate na spec **007**) |
| **N3 A** | ADR ampliação `AchadoAnalise` | T003 (rastreio); campos handoff [PENDENTE ADR] |
| **006-10 A** | Alinhar Jira LLT-202 | T004 (fora desta pasta) |
| **Escape `\` na assinatura** | Regra para barra invertida literal nos componentes do `valor` da assinatura | Sem fonte — referência cruzada **T043** / spec **008**; **não** bloqueia normalizador **010** (**DN-008-9 A** cobre só **`|`** → **`\|`**) |

**Respondidas (etapas 3b/3c/3e/3h — desbloqueiam onde indicado)**: **DN-010-1 A**, **DN-010-2**, **009-5 revisto**, **DN-007-2 B**, **DN-008-4 B**, **DN-010-3 B** (T046), **DN-010-4 C** (T013 → **P-010-2**), **DN-008-5 A** (T043), **010 H1** (**T043**, etapa 3h), **DN-008-9 A** (respondida — escopo composição/escape na spec **008**; **não** bloqueia **010**), **007 formato A**, **007 caminho C** (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; §Etapa 3c, João 16/11 11h; Etapa 3 16/11 00h41; §Etapa 3g, João 16/11 13h22; origem: respostas-etapa3c.md §§4, 7–8). Referência cruzada 006/007/008: **DN-006-3 B**, **DN-007-4 A**, **DN-008-7 B** — fonte `respostas-joao.md` §Etapa 3c (João, 16/11 11h02).

---

## Dependencies & Execution Order

### Phase Dependencies

- **Cadeia specs 008→009 (handoff 010)**: o aceite da spec **009** só fecha após a spec **008** (**009-22 A**, João: **009 espera a 008**); handoff mascarado Finder→Módulo (T027–T029, T042) e ordem mascaramento→normalização **008** (T043, FR-015) dependem dessa cadeia — sem inventar aceite da **009** nesta feature (fonte: respostas-joao.md, 009-22 A; origem: Rodada pós-revisão T1, 16/11 00h01; analyze-010-final.md **H1**).
- **Phase 1 (Setup)**: Sem dependências — pode iniciar imediatamente (tarefas documentais/pendências).
- **Phase 2 (Foundational)**: Depende de T001; **bloqueia** todas as user stories. **T013** depende de **P-010-2** (path/formato na ADR — **DN-010-4 C**); **T017** de **P-010-2** (regex ADR).
- **Phase 3 (US1)**: Depende da Phase 2; MVP alvo.
- **Phase 4 (US2)**: Depende da Phase 2; T041 **desbloqueada** (**DN-008-4 B**); T043 **desbloqueada** (**DN-008-5 A**); T046 **desbloqueada** (**DN-010-3 B**); T042 (Polish) de **DN-009-2**.
- **Phase 5 (US3)**: Depende de T014/T031 (contador); pode iniciar após T010/T014.
- **Phase 6 (Polish)**: Depois das histórias desejadas; T040 de **P-010-2** + **P-010-27**; T034–T039, T044, T045 em paralelo quando couber.

### User Story Dependencies

| Story | Prioridade | Depende de | Bloqueios externos |
|-------|------------|------------|-------------------|
| US1 | P1 | Phase 2 | P-010-2 (regex/placeholder canônico em T017/T019), P-010-27 (SC-001/003) |
| US2 | P1 | Phase 2 | P-010-27 (SC-002); **P-010-2** em T013 (Phase 2); DN-009-2 só em T042; cadeia **008→009** (**009-22 A**: aceite **009** após **008**) condiciona handoff T042 e canônico na **008** (T043 integra **DN-008-5 A**) (fonte: respostas-joao.md, 009-22 A) |
| US3 | P2 | Phase 2, contador | P-010-27 (SC-004) |

### Parallel Opportunities

- **Phase 1**: T002–T007 em paralelo `[P]`.
- **Phase 2**: T009–T012 em paralelo `[P]` após T008 (ou em paralelo se structs independentes).
- **US1**: T019–T020 em paralelo `[P]` (mesma função, regras distintas — coordenar arquivo).
- **US2**: T025 com T028 em paralelo `[P]` se arquivos distintos.
- **Polish**: T034, T036, T037, T039, T044, T045 em paralelo `[P]`.

### Parallel Example: User Story 1

```text
# Após Phase 2:
T019 — placeholder e-mail em internal/application/ (validarSeguranca)
T020 — marcador linha segredo em internal/application/ (validarSeguranca)
T021 — fora do catálogo inalterado (mesmo arquivo; serializar se conflito)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Completar Phase 1 (pendências rastreadas).
2. Completar Phase 2 (interim ADR_064 sem lista regex definitiva; adaptador catálogo após **P-010-2** / metadados ADR — **DN-010-4 C**).
3. Completar Phase 3 (US1) até T022.
4. **Parar e validar** quickstart Cenários 1, 4, 5, 8 (parcial) — **não** declarar SC-001 fechado (P-010-27).

### Incremental Delivery

1. Foundation → US1 (trechos mascarados) → US2 (handoff/consumo) → US3 (contador) → Polish.
2. Cada incremento mantém FR-011 (sem bruto no handoff).

### Parallel Team Strategy

- Dev A: T014–T022 (validarSeguranca + US1).
- Dev B: T024–T029, T041, T043, T046 (classificador + handoff US2 + versão catálogo + integração normalização **DN-008-5 A**); evitar T013 até **P-010-2**.
- Dev C: T031–T032 (contadores US3) após T014.

---

## Notes

- **Sem** meta de desempenho (010-28 A).
- **LLT-51** fora de escopo — verificação explícita em **T045** (FR-014).
- Path/formato do artefato versionado: só na **ADR** (**DN-010-4 C**, **P-010-2**) — não inventar (spec.md FR-008; **DN-010-3 B**: mismatch → recusar subir).
- Biblioteca Go de regex: **NEEDS CLARIFICATION** nas fontes prioritárias — decisão na implementação só se ADR ou João fixar pacote (plan.md Technical Context).
- Gate **≥95% precisão por classe**: spec **007** / **006** (**DN-007-2 B**), não LLT-202 — ver T044 e SC-007 no spec.
- **010 H1** / **DN-008-9 A** (etapa 3h): reuso do normalizador **008** em **T043**; composição da assinatura e escape **`|`** → **`\|`** na spec **008** — escape de **`\`** sem fonte permanece pendência (tabela § Pendências).
- Falha OpenObserve / aviso Core (**DLQ B**, **DN-008-3**/**DN-008-6** rev. 3h): **fora** das tarefas desta feature (specs **008**/**009**; **ADR_051 [PENDENTE]**).
- Linhas com `BLOQUEADO (…)` não devem ser tratadas como concluídas nem gerar decisões locais.

---

## Rastreabilidade resumida

| Tarefas | Origem principal |
|---------|------------------|
| T002, T017, T019, T040 | 010-2 A, P-010-2, FR-001, FR-008 |
| T013 | DN-010-4 C, P-010-2, FR-008, 010-18 A |
| T003 | N3 A |
| T004 | 006-10 A, T6 A |
| T005, T023, T030, T033, T040 | 010-27 A, P-010-27, SC-005 |
| T006, T029 | DN-010-1 A / DN-009-1 A, FR-004 |
| T007, T009, T014 | DN-010-2, FR-004a |
| T008–T022 | FR-002, FR-003, FR-006, FR-009, FR-010, US1 |
| T024–T030, T041, T043, T046 | FR-004, FR-011, US2, 010-5, 010-24, 009-5 revisto, DN-008-4 B, DN-008-5 A, DN-010-3 B |
| T028, T029 | analyze F1, FR-004, 009-5 revisto |
| T031–T033 | FR-005, US3, 010-17, clarify I2 |
| T034–T037, T045 | FR-012, FR-013, FR-014; T037 analyze F2 |
| T038 | T1 revisto, P-emenda-I, analyze C1 |
| T039, T044 | 010-18/010-19, DN-007-2 B, SC-007 |
| T042 | DN-009-2, plan.md § API do Módulo Ocorrência |
## Ajuste de execução documentado

Novas notas e matrizes pedidas em arquivos de checklist são registradas em rastreabilidade.md; checklists e marcadores são somente leitura pela skill. Caminhos de código podem usar arquivos isolados e aliases dos donos comuns 006/008; o rastreio aponta a implementação real. Isso não converte aceite bloqueado em concluído.
