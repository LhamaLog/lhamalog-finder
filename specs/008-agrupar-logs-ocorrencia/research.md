# Research: Agrupar logs semelhantes na mesma ocorrência

**Feature**: `008-agrupar-logs-ocorrencia` | **Date**: 2026-11-16 (etapas 3b–3h)

Consolidação da Fase 0. Cada item cita fonte primária (João, revisão T1, respostas documentadas, contrato 08/11 ou artefato do repositório).

---

## R1 — Grão do achado e da ocorrência

- **Decision**: **Achado** = logs com a **mesma assinatura** no **produto**, materializado num `AchadoAnalise` (**008-18 A**). **Ocorrência** = uma por achado ERROR/WARNING/INFO (**T2 C**), unicidade **(assinatura, classe, produto, tenant)** (**008-11 rev.** 16/11, conciliado com **008-15 A**).
- **Rationale**: LLT-200 crit. 2 sem duplicata; João amplia chave com **classe** e **tenant** (spec.md D4/D5).
- **Alternatives considered**: Chave só (assinatura, produto, tenant) — **substituída** pela rodada 16/11.
- **Fonte**: (fonte: respostas-joao.md, 008-18 A / T2 C / 008-11 rev.; origem: Jira LLT-200, ADR_040)

---

## R2 — Quem agrega dentro e entre lotes

- **Decision**: **Finder** agrupa no **lote** e **consulta ocorrências existentes** para agregar **entre lotes** (contador, primeira/última vez) antes do handoff (**N2 B**). **Módulo Ocorrência** persiste o efeito na ocorrência (**008-12/008-14** no achado enviado).
- **Rationale**: reconcilia LLT-200 com T1 rev. (revisao-T1-drive.md A.1: AchadoAnalise sem contagem no Drive — extensão **N3**).
- **Alternatives considered**: N2 A (só Módulo entre lotes) — **não escolhida**.
- **Fonte**: (fonte: respostas-joao.md, N2 B; origem: revisao-T1-drive.md §2, ADR_062)

---

## R3 — Composição da assinatura (etapa 3b)

- **Decision**: Assinatura = **nível OTel** + **mensagem normalizada** + **`service.name`** + **host** (máquina) (**DN-008-1 B**, fech. 3b). **`traceId`** e **`correlationId`** **MUST NOT** integrar a assinatura (**008-2 rev.**; substitui **008-2 A**). Risco D2 do contrato §4 (correlationId aleatório) **não se aplica** à assinatura após esta decisão.
- **Rationale**: ADR_013 exige mesma máquina; lista fechada pelo João na etapa 3b.
- **Alternatives considered**: C) incluir traceId/correlationId (respostas-etapa3.md opção C) — **não escolhida** após **DN-008-1 B**.
- **Fonte**: (fonte: respostas-joao.md, DN-008-1 B / Etapa 3; origem: ADR_013, contrato-log-ingested-v1.txt §4)

**Decidido (etapa 3c):** algoritmo de **normalização da mensagem** — **DN-008-4 B** (trim, case-insensitive, marcadores para número, UUID, IP, data, hex; dono spec 008; consumidores 006/009/010).

**Decidido (etapa 3e):** literais dos marcadores — **DN-008-5 A**: **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`** (fonte: respostas-joao.md §Etapa 3c 11h01/11h; origem: DN-008-5 A; respostas-etapa3c.md §4). Serialização de `valor`: ver **R17** (**008-assinatura A**; separador **`|`** — **DN-008-8 A**, etapa 3g).

---

## R4 — Equivalência histórica e recorrência (OpenObserve)

- **Decision**: Similaridade/recorrência **fora do lote atual** via **consulta ao OpenObserve** (**T3 B**, **008-3 B**). Critério **recorrente**: **ADR_013** (3–4 vezes em arquivos/lotes distintos, mesma máquina; repetição no mesmo arquivo **não** conta — **008-5 A**). Janela: **últimos 500 lotes** (**008-6 A**). **ADR_074** substituída por este caminho (**008-7 B**, **ADR_168**).
- **Rationale**: Fim do banco vetorial legado (FR-008).
- **Alternatives considered**: Igualdade literal só no lote — **insuficiente** para **008-3 B**; Elasticsearch legado — **proibido**.
- **Fonte**: (fonte: respostas-joao.md, T3 B / 008-3 B / 008-5 A / 008-6 A; origem: ADR_013, ADR_063, ADR_168)

**NEEDS CLARIFICATION**: forma da consulta (query, filtros, índice) — `internal/adapter/out/openobserve/doc.go`; **não** inventar.

---

## R5 — Contador e extremos temporais

- **Decision**: Contador = **+1 por log agregado** (**008-12 A**); repetição no **mesmo lote** **conta** (**008-13 A**). **Primeira/última vez** = min/max do **timestamp do log** (**008-14 A**).
- **Rationale**: LLT-200 crit. 1; ADR_013 só distingue recorrência, não contador.
- **Alternatives considered**: Usar `receivedAt` do lote — **rejeitado**.
- **Fonte**: (fonte: respostas-joao.md, 008-12 A / 008-13 A / 008-14 A; origem: Jira LLT-200)

---

## R6 — Status, reabertura e recorrente

- **Decision**: **Resolvida** = **FECHADA** (**008-8 A**). Reabertura: só **FECHADA→ABERTA** (**008-4 B**). Atualizar contador **não** muda status salvo **008-4** (**008-10 A**). **Recorrente** é **atributo**, não status (FR-007, ADR_063). Fechamento pelo **operador** (**008-9 A**).
- **Rationale**: ADR_082.
- **Alternatives considered**: Recorrente como status — **rejeitado**.
- **Fonte**: (fonte: respostas-joao.md, 008-4 B / 008-8 A / 008-9 A / 008-10 A; origem: ADR_082, ADR_063)

---

## R7 — Classe INFO e taxonomia

- **Decision**: **INFO** gera achado e ocorrência como ERROR/WARNING (**T7 A**, decisão consciente 16/11).
- **Rationale**: Diverge SSID_007.1/009.1/ADR_040 — **vale João** (spec.md D1).
- **Alternatives considered**: Só ERROR/WARNING — **rejeitado**.
- **Fonte**: (fonte: respostas-joao.md, T7 A; origem: ADR_170)

---

## R8 — Entrada: lote, produto, contrato v1

- **Decision**: **Lote incompleto** (spec 005) **não** dispara agrupamento (**T4 A**). Envelope **`log.ingested` schemaVersion `"1"`** (**T5 D**; `internal/domain/envelope.go`). `productId` ausente ou metadata≠body: **falha crítica sem retry** (**008-16 A**).
- **Rationale**: spec 005 FR-007/008; ADR_049; contrato §4.
- **Fonte**: (fonte: respostas-joao.md, T4 A / T5 D / 008-16 A; origem: contrato-log-ingested-v1.txt §4, ADR_049)

---

## R9 — Handoff e idempotência

- **Decision**: Finder envia **exatamente um** `ResumoAnalise` **por lote** (**DN-008-7**, ADR_062) via HTTP+JWT (**N1 A**); **consulta prévia** ao Módulo para agregar entre lotes (**N2 B**): **uma chamada por achado** (**DN-008-7 B**, etapa 3e). Chave idempotência **`batchId` + `execId`** (**N4 A**). Campos de agrupamento no achado exigem **ADR N3** — pendentes nos contratos. Detalhes HTTP: **DN-009-2**.
- **Rationale**: ADR_070 (arquivoId histórico) vs contrato v1 batchId; consulta por achado concilia N2 B com chave 008-11 por achado.
- **Fonte**: (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; origem: ADR_062, N2 B, respostas-etapa3c.md §6)

---

## R10 — Notificação e escopo

- **Decision**: Política ADR_035 = **dependência** LLT-216/009, **fora** desta feature (**008-17 C**).
- **Fonte**: (fonte: respostas-documentadas.md, 008-17; origem: ADR_035, backlog E10.F1.1)

---

## R11 — Migração / legado spec 006

- **Decision**: **Não há legado documentado** em produção; **não** exigir migração (**DN-008-2**). **Ressalva**: ambiente de teste com ocorrências gravadas exige confirmação do João (spec.md D2b).
- **Rationale**: contrato v1 §11; T1 rev. (Finder não gravava ocorrência); LLT em "A fazer".
- **Fonte**: (fonte: respostas-joao.md, DN-008-2 / Etapa 3; origem: contrato-log-ingested-v1.txt §11, respostas-etapa3.md §5)

---

## R12 — Falha consulta OpenObserve

- **Decision**: Falha **temporária** na consulta (recorrência/similaridade): até **3 tentativas** na consulta; **sem DLQ** no `log.ingested` (**DLQ B**, etapa 3h); após esgotar, Finder MUST avisar o **Core** por **API interna** com **batchId** e **motivo** (**DN-008-3**, **DN-008-6** rev. 3h). Erro **crítico** sem retry. **MUST NOT** seguir o fluxo assumindo `recorrente=false`.
- **Rationale**: contrato v1 §8–§9 (não copiar retry/DLQ do guia; OpenObserve indisponível = temporário); **DN-006-1b B**; ADR_014 (não concluir antes do handoff completo).
- **Alternatives considered**: DLQ FilaOrdenada + consumidor `dlq-alerta` + WebSocket — **substituído** na etapa 3h (respostas-etapa3g.md §1b opção B).
- **Fonte**: (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B; origem: respostas-etapa3g.md §1, contrato v1 §8)

**Decidido (etapa 3h):** contrato Finder→Core para o aviso (**ADR_051 [PENDENTE]**): rota, payload e evento **MUST NOT** inventar neste repositório. Isto **não** substitui a spec **009** (publicação de evento/notificação de ocorrência); complementa **FR-002a** nesta feature.

---

## R13 — Normalização da mensagem (etapa 3c)

- **Decision**: **DN-008-4 B**: (1) aparar espaços nas extremidades; (2) comparação **case-insensitive**; (3) substituir ocorrências de **número**, **UUID**, **IP**, **data** e **hexadecimal** por **marcadores** — **um por categoria**, distintos entre si. Algoritmo **canônico** da spec **008**; specs **006**, **009** e **010** MUST reutilizá-lo ao derivar assinatura/equivalência a partir da mensagem.
- **Rationale**: FR-002b; fecha equivalência prática alinhada a ADR_013 (“praticamente idêntico”); literais dos tokens fechados em **DN-008-5 A** (etapa 3e); casos de teste alinhados a **010-27 A** quando publicados.
- **Alternatives considered**: Regra só na 006 ou herdar mascaramento 010 sem João — **rejeitado** (dono 008 explícito no João).
- **Fonte**: (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: DN-008-4 B)

---

## R14 — Meta de precisão do classificador (referência cruzada)

- **Decision**: Meta **≥ 95% de precisão em cada classe** (ERROR, WARNING, INFO) é critério de aceite da spec **007** / gates da **006** (**DN-007-2 B**, **007-meta**). **Fora** desta feature 008.
- **Rationale**: LLT-199 / 006-2 B; João separa treino/validação (007) de agrupamento (008).
- **Fonte**: (fonte: respostas-joao.md §Etapa 3b; origem: DN-007-2 B)

---

## R15 — Base rotulada (spec 007 — só documentação)

- **Decision**: Formato **JSONL** (**007 formato A**); caminho **`datasets/labeled-logs/`** na raiz do repositório (**007 caminho C**). **Não** criar pasta nem arquivos na execução da 008; implementação pertence à spec **007**.
- **Rationale**: decisão João etapa 3b; constituição V (não inventar dados).
- **Fonte**: (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: 007 formato A / 007 caminho C)

---

## R16 — Agrupamento local vs consulta OpenObserve (escopo)

- **Decision**: **No lote em processamento**, equivalência MUST seguir assinatura **DN-008-1 B** + normalização **DN-008-4 B** (determinístico, local). **OpenObserve** serve **somente** para histórico/recorrência/similaridade **fora do lote** (**008-3 B**, **008-5/6**); a consulta **MUST NOT** relaxar nem substituir o agrupamento local por assinatura (**FR-008**, achado analyze **A1**).
- **Rationale**: Evita ambiguidade entre “similaridade OO” e partição no lote (spec.md Clarifications etapa 3c analyze).
- **Alternatives considered**: Usar OO também para agrupar no lote corrente — **não autorizado** face a **DN-008-1 B** / **DN-008-4 B**.
- **Fonte**: (fonte: respostas-joao.md, 008-3 B / DN-008-1 B / DN-008-4 B; spec.md FR-008 etapa 3c analyze; origem: T3 B, ADR_168)

---

## R17 — Serialização canônica de `valor` da assinatura

- **Decision**: **`valor`** = **texto** concatenando **nível**, **mensagem normalizada**, **`service.name`** e **host** com o separador fixo **`|`** (barra vertical U+007C) entre componentes adjacentes (**008-assinatura A**, etapa 3e; **DN-008-8 A**, etapa 3g). **`|`** literal em componente: escape **`\|`** (**DN-008-9 A**, etapa 3h). Escape de **`\`** literal: pendente (sem fonte).
- **Rationale**: Opção A do Especialista §9 alinhada a Finder.mmd `log_signature` (string); hash (**opção B**) **não** escolhido; separador **`|`** fechado pelo João (16/11 12h40).
- **Alternatives considered**: SHA-256 do texto (**opção B** respostas-etapa3c.md §9) — **rejeitado** (João **008-assinatura A**). Outros separadores (`\t`, `NUL`, etc.) — **rejeitados** (não escolhidos em **DN-008-8 A**).
- **Fonte**: (fonte: respostas-joao.md, DN-008-8 A / §Etapa 3g (João, 16/11 13h22), DN-008-9 A; origem: §Etapa 3e (João, 16/11 12h40); 008-assinatura A; respostas-etapa3c.md §9; DN-008-1 B)