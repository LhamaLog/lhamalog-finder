# Research: Modelo de classificação de logs

**Feature**: `006-modelo-classificacao-logs` | **Date**: 2026-11-16 (etapa 3h)

Consolidação da Fase 0. Cada item cita fonte primária (João, revisão T1, respostas documentadas ou artefato do repositório).

---

## R1 — Taxonomia de classes

- **Decision**: Usar **ERROR**, **WARNING**, **INFO** no campo `tipo` do achado; sem spam e sem “não classificado”.
- **Rationale**: Alinha constituição III (erro/alerta/informação) à ADR_170; decisões 006-6 A e 007-9 A.
- **Alternatives considered**: Enum ERRO/WARNING/SPAM (ADR_019 / Finder.mmd legado) — **rejeitado** (fonte: respostas-joao.md, 007-9 A; spec.md taxonomia 14/11).
- **Fonte**: (fonte: respostas-joao.md, 006-6 A / 007-9 A; origem: ADR_170, constituição III)

---

## R2 — Regra nível + texto + catálogo

- **Decision**: Campo de nível e texto da mensagem **juntos**; classe do log = mais grave dos dois; termos só como **palavra inteira**; catálogo **somente** se nível não reconhecível; primeira entrada válida do catálogo; default **INFO** (incluído no `ResumoAnalise`, T7 A).
- **Rationale**: Spec original 2026-11-09 permanece vigente nos cenários de aceite; T7 A inclui INFO como achado.
- **Alternatives considered**: Catálogo vence nível — **rejeitado** (sessão 2026-11-09 substituída apenas onde indicado na spec).
- **Fonte**: (fonte: spec.md, User Story 1 e 3, FR-004–FR-006; respostas-joao.md, T7 A / 006-11 A)

---

## R3 — Quem persiste ocorrência e contrato de saída

- **Decision**: Finder produz **AchadoAnalise** + **ResumoAnalise** e envia ao **Módulo de Ocorrência**; Módulo cria/atualiza ocorrências; Finder **não** grava MongoDB nem publica evento LLT-201.
- **Rationale**: T1 REVISTO fecha divergência com ADR_062; sequência Drive (revisao-T1-drive.md §1).
- **Alternatives considered**: T1 A (Finder grava MongoDB) — **substituída** (fonte: respostas-joao.md, T1 REVISTO).
- **Fonte**: (fonte: respostas-joao.md, T1 REVISTO / N1 A; origem: ADR_062, revisao-T1-drive.md §2)

---

## R4 — Transporte Finder → Módulo Ocorrência

- **Decision**: **HTTP interno com JWT** (mesmo padrão 009-3 D / ADR_051). Rota, claims, códigos e timeout **não** fechados nesta feature (**DN-009-2**).
- **Rationale**: N1 A pós-revisão T1; João §Etapa 3: pendência no plan até contrato do Módulo.
- **Alternatives considered**: Fila Redis Streams (N1 opção B) — **não escolhida**.
- **Fonte**: (fonte: respostas-joao.md, N1 A / Etapa 3 API Módulo; origem: ADR_051, respostas-etapa3.md §10)

---

## R5 — Grão: achado, ocorrência, upsert

- **Decision**: Um **AchadoAnalise** por achado no lote; no Módulo, **uma ocorrência por achado** (T2 C); chave de unicidade no módulo: **assinatura + classe + produto + tenant** (rodada 16/11 00h01). **Upsert** ao criar ocorrência é do **Módulo**, sem efeito no Finder (006-9 A após T1 revisto).
- **Rationale**: ADR_040/035; 008-15 A + revisão 008-11.
- **Alternatives considered**: Uma ocorrência por apontamento inteiro — **substituída** por T2 C.
- **Fonte**: (fonte: respostas-joao.md, T2 C / 008-15 A / 008-11 rodada 00h01 / 006-9 A; origem: ADR_040)

---

## R6 — Assinatura do achado (campos)

- **Decision**: Assinatura nesta feature = **nível (classe do log)** + **mensagem normalizada** + **`service.name`** + **host/máquina** (**DN-008-1 B**). **Não** incluir `traceId` nem `correlationId` na assinatura da 006 (risco contrato v1 §4 documentado em spec D5; 008-2 A permanece histórico transversal).
- **Rationale**: João Etapa 3; ADR_013 (“mesma máquina”); fecha lacuna de campos pendentes ADR_063/074 para o MVP desta feature.
- **Decision (normalização da mensagem)**: **DN-008-4 B** — trim, case-insensitive, marcadores número/UUID/IP/data/hex; literais **DN-008-5 A**: `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` (detalhe **R17**).
- **Decision (serialização do valor da assinatura)**: **008 assinatura A** — texto com componentes concatenados por separador **`|`** (**DN-008-8 A**); **`|`** literal em componente → **`\|`** (**DN-008-9 A**); escape de **`\`** literal: **pendência** (sem fonte). Canônico na spec **008**. (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), DN-008-8 A; §Etapa 3g (João, 16/11 13h22), DN-008-9 A; respostas-etapa3c.md §9)
- **Fonte**: (fonte: respostas-joao.md, Etapa 3 DN-008-1 B / §Etapa 3b DN-008-4 B / §Etapa 3c DN-008-5 A; origem: ADR_013, respostas-etapa3.md §4; spec.md FR-007, FR-013)

---

## R7 — Lote incompleto e falhas

- **Decision**: Lote incompleto (**T4 A**) **não** é sucesso; distinguir log ausente, armazenamento inacessível, falha de envio HTTP; reenvio `ResumoAnalise` com mesma chave; confirmação do Módulo antes de concluir varredura (FA-02 Drive) — confirmação HTTP **bloqueada até DN-009-2** para implementação real.
- **Decision (`assinatura_incompleta` — analyze E2/D1)**: Log contável sem `service.name` ou **host** (DN-008-1 B) → falha distinguishável; sem agrupamento; sem sucesso de lote completo (spec.md FR-007/FR-008; 007-8 A como padrão).
- **Decision (`produto_ausente_ou_inconsistente` — 008-16 A, analyze2 E2)**: `productId` **ausente** no `log.ingested` v1 ou **inconsistente** entre metadata validada (token) e corpo repetido → **falha crítica permanente**, **sem retry** de processamento do lote (ADR_049); MUST NOT sucesso de varredura; distinto de `identidade_incompleta` (tenant/produto vazio no achado — 007-8 A).
- **Rationale**: spec 005 + T4 A; revisao-T1-drive.md §2; contrato ingested §4.
- **Fonte**: (fonte: respostas-joao.md, T4 A / 008-16 A; origem: revisao-T1-drive.md §2, ADR_014, ADR_049, contrato-log-ingested-v1.txt §4)

---

## R8 — Versão única das regras

- **Decision**: Versão **única** combinando regras + catálogo + base rotulada (006-3 C) em **AchadoAnalise** e/ou **ResumoAnalise**; ao mudar versão, **só novos logs** (006-4 A).
- **Rationale**: LLT-199 critério 3; T1 revisto.
- **Fonte**: (fonte: respostas-joao.md, 006-3 C / 006-4 A; origem: LLT-199, LLT-51)

---

## R9 — Contrato `log.ingested` v1

- **Decision**: Aceitar **apenas** `schemaVersion` `"1"`; `batchId` em `openobserve.batchId` / `idempotencyId`; rejeitar eventos fora do contrato (espelhado em `InterpretarEnvelope`).
- **Rationale**: T5 D; contrato 08/11 prevalece sobre ADR_146 v2 “schema 2” (spec D3).
- **Fonte**: (fonte: respostas-joao.md, T5 D; origem: contrato-log-ingested-v1.txt §2–4; `internal/domain/envelope.go`)

---

## R10 — Chave de idempotência do ResumoAnalise

- **Decision**: **`batchId` obrigatório**; compor **`batchId` + `execId`** quando `execId` disponível (**N4 A**, **DN-006-1**). `execId` origem Varredura: `V<uuid>` automático, `M<uuid>` manual (ADR_075); **opcional** (“quando houver”, ADR_146). Enquanto o Finder consome **somente** `log.ingested` v1 **sem** canal que entregue `execId`, chave **somente `batchId`** (**DN-006-1b B**).
- **Rationale**: Corpo v1 não traz `execId`; João escolheu B na Etapa 3; alinha ADR_070/146 ao pipeline atual.
- **Fonte**: (fonte: respostas-joao.md, N4 A / Etapa 3 DN-006-1b B; origem: respostas-etapa3.md §3, contrato-log-ingested-v1.txt §§2–5, ADR_075)

---

## R11 — Campos AchadoAnalise / contadores ResumoAnalise

- **Decision (base Drive)**: Campos legados ADR_040 / revisao-T1-drive.md §2 (`titulo`, `resumo`, `numeroLinha`, `camposExtraido`, `dadosSensiveis`, `traceId`, `correlationId`, `hostName`, `moduloId`, `linhas`, `tipo`, `recorrente`).
- **Decision (escopo mínimo 006 — analyze E3)**: `recorrente` **= `false`** sem consulta entre lotes (**FR-009**; **N2 B** na spec 008). Campos Drive `titulo`/`resumo`/`numeroLinha`/`camposExtraido`: **presença** no payload; derivação mínima **[PENDENTE N3 A]** até ADR.
- **Decision (`hostName` vs `host` — analyze L1)**: `hostName` no `AchadoAnalise` espelha 1:1 o `host` da assinatura quando disponível (ADR_040; DN-008-1 B).
- **Decision (contadores)**: ADR_102 com contador de **INFO** **em lugar de** «spam ignorado» (João Etapa 3 opção B; ADR_170 sem SPAM).
- **Decision (extensão João)**: **assinatura**, **contagem**, **primeira/última ocorrência no lote**, **versão única** (N3 A) — formalização em ADR **pendente**.
- **Fonte**: (fonte: respostas-joao.md, N3 A / Etapa 3 contadores; origem: ADR_040, ADR_102, ADR_170)

---

## R12 — Catálogo de classificação

- **Decision**: **Não criar** catálogo nesta feature; opcional na entrada; sintaxe/dono **pendente** (006-7 D).
- **Fonte**: (fonte: respostas-joao.md, 006-7 D; origem: ADR_167)

---

## R13 — Métricas de qualidade (aceite)

- **Decision**: Documentar **precisão, recall e F1 por classe** (006-2 B) com base spec 007. Formato da base: **JSONL** (**007 formato A**); caminho documentado: `datasets/labeled-logs/` na raiz do repositório (**007 caminho C** — pasta **não** criada nesta feature). Gold **por apontamento** (**DN-007-4 A**); **classe prevista do apontamento** = **mais grave** entre as classes previstas dos logs do apontamento (**ERROR > WARNING > INFO** — **007 H2 D**, ADR_170). **Gate formal (SC-007, T047)**: em **cada** classe ERROR, WARNING, INFO — precisão **≥ 95%** (**DN-007-2 B**), recall **≥ 95%** e F1 **≥ 95%** (**DN-006-2 B**, etapa 3e — fonte: respostas-joao.md §Etapa 3c, João 16/11 11h e 11h01; respostas-etapa3c.md §1).
- **Impacto 006**: gate SC-007 / T047 **desbloqueado** para as três métricas por classe; regra de comparação à base 007 **fechada** (etapa 3g).
- **Fonte**: (fonte: respostas-joao.md §Etapa 3b/3c, DN-007-2 B / DN-006-2 B / 007 formato A / 007 caminho C; §Etapa 3e (João, 16/11 12h40), 007 H2 D; origem: LLT-197, LLT-199, ADR_170, respostas-etapa3f.md opção D)

---

## R14 — INFO gera ocorrência (divergência Drive)

- **Decision (produto)**: Achados **INFO** entram no `ResumoAnalise`; Módulo deve aceitar (T7 A, 006-11 A).
- **Fonte**: (fonte: respostas-joao.md, T7 A; origem: revisao-T1-drive.md A.1, ADR_170)

---

## R15 — Consulta ao Módulo e envio (N2 B, DN-008-7)

- **Decision (envio)**: **Um** `ResumoAnalise` **por lote** ao Módulo de Ocorrência (fonte: respostas-etapa3c.md §6; origem: ADR_062, Finder.mmd).
- **Decision (consulta prévia N2 B)**: **Uma chamada ao Módulo por achado** antes de montar o envio do lote (**DN-008-7 B**). Rota/contrato da consulta **DN-009-2** (sem inventar path).
- **Decision (entre lotes)**: Agrupamento/recorrência histórica via OpenObserve permanece na spec **008**; FR-009 exclui agrupamento **entre lotes** do mínimo desta feature.
- **Fonte**: (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; origem: LLT-200, ADR_062, respostas-etapa3c.md §6)

## R15b — Desempenho do ciclo classificar + montar resumo (DN-006-3)

- **Decision**: Meta **≤ 2 s por lote** para classificar todos os logs do lote + montar `ResumoAnalise` (**DN-006-3 B**, **SC-002**). Perfil **DN-006-3-ref**: **VM-A real**, **todos os serviços**, lote **1.000** registros. Não inclui latência de consultas **DN-008-7 B** nem handoff **DN-009-2** (spec.md SC-002).
- **Fonte**: (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-006-3 B; §Etapa 3e (João, 16/11 12h40), DN-006-3-ref)

---

## R16 — Estado do repositório (implementação)

- **Decision**: Não há pacotes de classificação nem cliente HTTP para Módulo Ocorrência no `internal/` atual; envelope v1 e captura existem (`domain/envelope.go`, `application/capturar.go`). Mongo adapter é só `doc.go`.
- **Fonte**: (fonte: README.md § Layout; listagem `internal/**/*.go`)

---

## R17 — Normalização da mensagem (DN-008-4 B, DN-008-5 A)

- **Decision**: **DN-008-4 B** — trim; case-insensitive; substituição por categorias número, UUID, IP, data, hex. Literais **DN-008-5 A**: `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>`. Implementação **canônica única** na spec **008**: **`internal/application/assinatura/normalizar_mensagem.go`** (plan/tasks T016). A spec **006** MUST **importar/delegar** e **MUST NOT** criar cópia paralela (analyze2 **F1**). Casos de borda adicionais: **010-27 A** (spec 010).
- **Impacto**: testes determinísticos de agrupamento na 006 podem assertar marcadores **DN-008-5 A** via pacote 008.
- **Fonte**: (fonte: respostas-joao.md §Etapa 3b DN-008-4 B / §Etapa 3c DN-008-5 A; origem: ADR_013, specs/008 plan/tasks T016)

---

## R18 — Falha de consulta OpenObserve (DLQ B, DN-008-3 / DN-008-6 rev. 3h)

- **Decision**: **Não** há DLQ no fluxo `log.ingested` — o contrato v1 **não** herda retry/DLQ do guia da FilaOrdenada (**DN-006-1b B**). Após **3 tentativas** na consulta OpenObserve (recorrência entre lotes, spec **008**), o Finder avisa o **Core** pela **API interna** com **`batchId`** e **motivo**; **sem** DLQ no ingested e **sem** alerta WebSocket pós-DLQ.
- **Rationale**: João **DLQ B** corrige alcance anterior (etapas 3b–3e: 3 tentativas + DLQ / alerta Core após DLQ — respostas-joao.md Etapa 3; respostas-etapa3g.md §1b).
- **Alternatives considered**: DLQ no ingested espelhando FilaOrdenada — **rejeitada** (fonte: respostas-joao.md §Etapa 3g, DLQ B; origem: contrato-log-ingested-v1.txt §8, respostas-etapa3g.md §1b).
- **Escopo 006**: **rastreio** e referências em FR-009 / contratos; runtime e adaptador HTTP → **spec 008**. Rota, payload e evento da API Core: **ADR_051 [PENDENTE]** — MUST NOT inventar.
- **Fonte**: (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B; origem: respostas-etapa3g.md §1b, contrato-log-ingested-v1.txt §8)

---

## Resolução de NEEDS CLARIFICATION

| Item | Status |
|------|--------|
| Transporte handoff (N1 A) | **Resolvido** — HTTP+JWT; detalhes HTTP **DN-009-2** |
| Taxonomia | **Resolvido** — ERROR/WARNING/INFO |
| Chave batchId | **Resolvido** — Referencia / ingested v1 |
| Chave execId | **Resolvido** — DN-006-1 + DN-006-1b B (só batchId no pipeline atual) |
| Campos assinatura | **Resolvido** — DN-008-1 B |
| Algoritmo normalização | **Resolvido** — **DN-008-4 B** + literais **DN-008-5 A** |
| Serialização assinatura | **Resolvido** — **008 assinatura A**; separador **`|`**, **DN-008-8 A**; escape **`|`** → **`\|`**, **DN-008-9 A** (spec 008); escape **`\`**: pendência |
| Falha consulta OpenObserve | **Resolvido (rastreio 006)** — **DN-008-3** / **DN-008-6** rev. 3h: 3×, sem DLQ; aviso Core **batchId** + **motivo**; runtime spec **008**; API **ADR_051 [PENDENTE]** |
| Consulta Módulo / envio | **Resolvido** — consulta **por achado** (**DN-008-7 B**); envio **por lote** (respostas-etapa3c.md §6); HTTP **DN-009-2** |
| Desempenho SC-002 | **Resolvido** — meta **2 s/lote** (**DN-006-3 B**); perfil **DN-006-3-ref** (VM-A, 1.000 registros) |
| Campos ampliados AchadoAnalise | **Pendência N3 ADR** |
| Catálogo | **Pendência 006-7 D** |
| API Módulo (rota/JWT/códigos) | **DN-009-2** — bloqueia adaptador HTTP e consulta real |
| Aceite métrico formal | **Resolvido** — precisão/recall/F1 **≥ 95%/classe** (**DN-007-2 B**, **DN-006-2 B**) |
| Comparação avaliação × base 007 | **Resolvido** — gold por apontamento; classe prevista = mais grave entre logs (**007 H2 D**, ADR_170) |

Nenhum item **bloqueia o plano inteiro** sob 007-26 A nesta rodada (fonte: respostas-joao.md, Etapa 3/3e/3g).

## Evidência desta execução

006-7 D segue aberto: a lista ordenada é opaca, com predicado injetado; não há regex, arquivo ou dono de catálogo definido. T036 continua bloqueada; T037 opera apenas sobre essa fronteira lógica.

DLQ B / DN-008-3 / DN-008-6 revisados substituem o alcance antigo com DLQ e WebSocket: três falhas temporárias OpenObserve exigem aviso operacional Core (batchId/motivo). Esse aviso não é o evento LLT-201; cliente real exige ADR_051. DN-008-9 A fecha escape de `|`; escape de barra invertida literal continua pendente.

A avaliação offline usa `labeleddataset.Avaliar` e `ValidarGate` da 007, com um par gold/predição por apontamento e maior classe entre membros. Não há base real disponível; não se afirma ≥95%.
