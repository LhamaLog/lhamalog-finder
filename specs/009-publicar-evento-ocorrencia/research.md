# Research: Publicar evento de ocorrência (escopo Finder)

**Feature**: `009-publicar-evento-ocorrencia` | **Date**: 2026-11-16 (etapa 3h)

Consolidação de decisões das três fontes prioritárias para o plano. Nenhum item abaixo foi resolvido por suposição (constituição V).

---

## R1 — Papel do Finder na LLT-201 (T1 revisto)

- **Decision**: O **Finder não publica** evento de ocorrência nova/atualizada/recorrente. Envia **ResumoAnalise** ao **Módulo Ocorrência**; o Módulo cria/muda ocorrências e avisa o **Core** (**009-3 D**).
- **Rationale**: Alinha sequência varredura v6, ADR_062 e SSID_009.1 ao épico sem reatribuir publicação ao Finder.
- **Alternatives considered**: T1 A (João 14/11: Finder grava e publica) — **rejeitada** na revisão 15/11 23h47 (fonte: respostas-joao.md, T1 revisto; revisao-T1-drive.md §1; origem: ADR_062, Finder.mmd).

---

## R2 — Transporte Finder → Módulo Ocorrência

- **Decision**: **HTTP(S) interno com JWT** (**N1 A**), paridade com Módulo → Core (**009-3 D**). Detalhes de rota/resposta: **DN-009-2** (pendência João).
- **Rationale**: Fecha a pergunta N1 da revisão T1; não usar Redis Streams para handoff Finder → Módulo.
- **Alternatives considered**: Fila Redis (revisao-T1-drive.md tabela N1 opção B); ADR nova só para transporte (opção C) — **rejeitadas** por **N1 A**.

---

## R3 — Idempotência do ResumoAnalise

- **Decision**: Chave **`batchId`** quando o Finder só consome **log.ingested v1** sem `execId` (**DN-006-1b B**); **`batchId` + `execId`** quando o identificador de varredura estiver disponível (**N4 A**, **DN-006-1** — `execId` de Varredura `V<uuid>`/`M<uuid>`, **opcional** “quando houver”, ADR_146). Reenvio com mesma chave; Módulo não duplica efeito (**FR-F06**, FA-02, ADR_014).
- **Rationale**: Contrato 08/11 (`idempotencyId` = batchId); decisão etapa 3b alinha N4 ao fluxo atual.
- **Alternatives considered**: Chave sempre `batchId + execId` obrigatório — **rejeitada** por **DN-006-1b B**. `arquivoId + execId` histórico ADR_070 — **substituído** por **N4 A** (revisao-T1-drive.md N4).
- **Lacuna registrada**: transporte da chave (header vs corpo) — **sem fonte** nas três prioritárias; dependerá de **DN-009-2**.

---

## R4 — Momento e confirmação do handoff

- **Decision**: Enviar ao **fim da análise** do arquivo/lote; **aguardar confirmação** do Módulo antes de concluir varredura (**FR-F05**); falha após envio sem confirmação → reenvio mesma chave (FA-02).
- **Rationale**: sequência varredura v6 §2; **009-18** (lado Finder) em revisao-T1-drive.md.
- **Alternatives considered**: Concluir varredura sem confirmação — **sem fonte** que autorize.

---

## R5 — Conteúdo do ResumoAnalise

- **Decision**: **Contadores + lista AchadoAnalise** (**009-8**, ADR_062); achados **INFO** inclusos (**T7 A**, **006-11 A**); **gravidade ausente** no Finder (**ADR_062**); contador legado “spam ignorado” → **contador de INFO** (**spam: B**); nome canônico na ADR_102 [PENDENTE] até fechar lista — documentar como **contador de INFO** (**corrige analyze F3**, fonte: respostas-joao.md §Etapa 3, spam: B; origem: ADR_102, respostas-etapa3.md §9).
- **Rationale**: Divergência INFO × Drive consciente do João (16/11 00h01).
- **Alternatives considered**: Omitir INFO do handoff (SSID_007.1 / SSID_009.1) — **rejeitada** por **T7 A**.
- **Pendente**: campos N3 (assinatura, contagem, primeira/última, versão) — **N3 A** exige ADR nova; não preencher sem ADR.

---

## R6 — Pipeline downstream (referência, fora do Finder)

- **Decision**: Módulo → Core: **HTTP + JWT**, corpo alinhado a **envelope ADR_146** com campos LLT-201 no `payload` (**009-4 B**). No aviso HTTP: **id da ocorrência**, **tipo** (classe) e **gravidade** (**009-5 revisto**); no trecho **Core → Notificação**, campo **`ocorrencia`** no envelope permanece **id** opaco (**009-6 A**). Core → Notificação: **Redis Streams** `notificacao-eventos`, envelope **ADR_146** (**009-1 B**, **FR-009**). Notificação lê Ocorrência por **id** no MongoDB (**DN-009-1 A**) para resumo mascarado do e-mail (**LLT-218**).
- **Rationale**: Fecha [PENDENTE] ADR_051 para trecho Módulo → Core; reconcilia 009-5 A (“só id”) com diagrama Notificação v2 via **009-5 revisto**; separa idempotência HTTP (**id ocorrência + estado**, **009-9 A**) vs stream (**eventId** no consumidor, **009-10 B**).
- **Alternatives considered**: Finder → Core direto — **rejeitado** por T1 revisto. Aviso só com id sem tipo/gravidade (009-5 A literal) — **substituído** por **009-5 revisto** (fonte: respostas-joao.md §Etapa 3; origem: sequencia-notificacao.md v2 diagrama 1). Core monta texto do e-mail (DN-009-1 opção B) — **rejeitada** por **DN-009-1 A**.

---

## R7 — Semântica do aviso (natureza da mudança)

- **Decision**: Estados **nova**, **atualizada**, **recorrente** no payload do aviso; reabertura → **atualizada** (**009-14 A**); **recorrente** → **nova notificação** (**009-15 A**); **um aviso por ocorrência** por criação/mudança (**009-17**). Significado operacional de atualizada/recorrente no Finder/agrupamento: spec **008** (**009-22 A**).
- **Rationale**: LLT-201 critério 1 + SSID_009.1 passo 1; aceite Finder de estados bloqueado até 008.
- **Alternatives considered**: Aceitar 009 só para “nova” (**009-22** opção B) — **rejeitada** por **009-22 A**.

---

## R8 — Campo `tipo` (classe)

- **Decision**: **tipo** = classe **erro / alerta / informação** (ERROR, WARNING, INFO), distinta da natureza da mudança (**009-13 A**, constituição III).
- **Rationale**: ADR_170 e critério LLT-201.
- **Alternatives considered**: Tipo = estado nova/atualizada — **rejeitada** (fonte: respostas-documentadas.md, 009-13; origem: LLT-201).

---

## R9 — Efeito duplicado e idempotência ponta a ponta

- **Decision**: Duplicata indevida no produto = **2º alerta publicado pelo Core** em `notificacao-eventos` (**009-11 B**). Deduplicação HTTP no **Core** com chave **id ocorrência + estado** (**009-9 A**); na stream, idempotência no **consumidor** por **eventId** (**009-10 B**, guia §8).
- **Rationale**: Separa trechos HTTP vs Redis após **009-3 D**.
- **Alternatives considered**: Idempotência só no produtor — **rejeitada** por **009-10 B**.

---

## R10 — Falhas

- **Decision**: HTTP Módulo → Core: **tentar de novo e alertar** (**009-18 A**). Falha **consulta OpenObserve** no Finder (fluxo ingested): até **3 tentativas**; **sem DLQ** (**DLQ B**); depois **Finder → Core** API interna (**`batchId`**, **motivo** — **DN-008-3**/**DN-008-6**, contrato **ADR_051 [PENDENTE]**) (**FR-F07**). Falha **handoff** Finder→Módulo: não concluir varredura, reenvio com **mesma** chave (**FA-02**) — **distinto** de **009-18 A** (Módulo→Core) (**corrige analyze F1** para leitores de research; ajuste de `tasks.md` T036: `/speckit-tasks`).
- **Rationale**: ADR_014, ADR_046, seq. v6 §2 (revisao-T1-drive.md **009-18**).
- **Alternatives considered**: Outbox MongoDB no Módulo (**009-18** opção B) — **não escolhida** pelo João.

---

## R11 — Dependências entre features

- **Decision**: Handoff ResumoAnalise integra classificação (**006**, **T6 A** atualiza 006) e agrupamento entre lotes (**N2 B**, spec **008**). Esta spec **009** não substitui contratos da 008; reutiliza o mesmo handoff com coordenação explícita (**FR-F09**). **009 espera 008** para aceite de estados atualizada/recorrente (**009-22 A**).
- **Rationale**: Backlog E07.F3 ordem; **009-22 A**.
- **Alternatives considered**: Implementar publicação Finder antes da 008 — **rejeitada**.

---

## R12 — Meta de precisão (DN-007-2 B) e base 007

- **Decision**: **≥95% de precisão por classe** (ERROR, WARNING, INFO) na base de validação da spec **007**; **esta feature (009) não define gate métrico** — rastreabilidade a **LLT-199** / spec **006**; aceite de precisão em **SC-007** / **006-2 B**, não nos critérios de sucesso do handoff LLT-201 (fonte: respostas-joao.md §Etapa 3b, DN-007-2 B; origem: respostas-etapa3.md §1, 007-meta).
- **Rationale**: Separa classificação (006/007) do handoff ResumoAnalise (009).
- **Alternatives considered**: Gate 95% global (opção A em respostas-etapa3b.md) — **não escolhida** pelo João 16/11 09h02.
- **Documentação cruzada (007, não implementar aqui)**: formato base rotulada **JSONL** (007 formato **A**); caminho **`datasets/labeled-logs/`** na raiz do repositório (007 caminho **C**) — apenas registrar em specs/planos; **não** criar pasta (fonte: respostas-joao.md §Etapa 3b).

---

## R13 — Normalização da mensagem na assinatura (DN-008-4 B)

- **Decision**: Ao compor **assinatura** (ou campos derivados) em **AchadoAnalise**, o Finder MUST aplicar: aparar espaços nas extremidades; comparação **case-insensitive**; substituir **números**, **UUIDs**, **IPs**, **datas** e **hexadecimais** por marcadores com literais **DN-008-5 A**: `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>`. Valor serializado da assinatura: texto concatenando partes **DN-008-1 B** com **separador fixo** **`|`** (**DN-008-8 A**; dono canônico spec **008**) (**008 assinatura A**); **`|`** literal no componente → **`\|`** (**DN-008-9 A**). Algoritmo e ordem de aplicação **canônicos** na spec **008**; esta feature **reutiliza** sem regras divergentes (**FR-F11**, **N3 A**) (fonte: respostas-joao.md §Etapa 3b DN-008-4 B, §Etapa 3c 11h DN-008-5 A e 008 assinatura A, §Etapa 3e (João, 16/11 12h40) DN-008-8 A, §Etapa 3g DN-008-9 A; origem: DN-008-1 B, ADR_013).
- **Rationale**: Estabiliza agrupamento entre lotes (**N2 B**, **FR-F09**) sem duplicar semântica na 009.
- **Alternatives considered**: Normalização só na 008 sem obrigação no handoff — **rejeitada** por **FR-F11** e decisão 3c no spec.

---

## R15 — Literais dos marcadores (DN-008-5 A)

- **Decision**: Marcadores do normalizador de assinatura: `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` (um por categoria).
- **Rationale**: Fecha opção A da etapa 3c; mesma família que **DN-008-4 B**.
- **Alternatives considered**: Prefixo `__` (opção B em respostas-etapa3c.md §4) — **rejeitada** por **DN-008-5 A** (fonte: respostas-joao.md §Etapa 3c 11h01 e 11h; origem: DN-008-4 B).

---

## R16 — Aviso operacional após falha OpenObserve (DN-008-3 / DN-008-6, etapa 3h)

- **Decision**: Na **consulta ao OpenObserve** (recorrência **T3 B**), após **3 tentativas** em falha transitória, o **Finder** avisa o **Core** pela **API interna** com **`batchId`** e **motivo**; **sem DLQ** no **log.ingested v1** (**DLQ B**); **não** usar consumidor DLQ nem **WebSocket** pós-DLQ (substitui R16 anterior). Rota/payload/evento: **[PENDENTE ADR_051]**. Bloqueia handoff **ResumoAnalise** (**FR-F07**).
- **Rationale**: Contrato ingested não herda retry/DLQ da FilaOrdenada (respostas-etapa3g.md §1b); distinto do aviso LLT-201 Módulo→Core.
- **Alternatives considered**: DLQ + grupo `dlq-alerta` + WebSocket (guia FilaOrdenada) — **rejeitada** para ingested por **DLQ B** (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22); respostas-etapa3g.md §1; origem: contrato log.ingested v1, DN-006-1b B).

---

## R17 — Granularidade do handoff Finder → Módulo (DN-008-7)

- **Decision**: **Um** **ResumoAnalise** por lote (arquivo/lote analisado), com contadores + lista de **AchadoAnalise** e chave de idempotência — **não** uma chamada HTTP por achado.
- **Rationale**: ADR_062, Finder.mmd, **N1 A**, **N4 A**.
- **Alternatives considered**: Uma chamada por achado (DN-009-2 opção 10a-B) — **rejeitada** para handoff do resumo (fonte: respostas-etapa3c.md §6; origem: ADR_062).

---

## R18 — Consulta prévia entre lotes (DN-008-7 B)

- **Decision**: Para **N2 B**, o Finder faz **uma consulta ao Módulo Ocorrência por achado** (assinatura/classe elegível), não uma consulta agregada única por lote (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); origem: DN-008-7 B, N2 B).
- **Rationale**: Fecha lacuna que o Especialista marcou SEM FONTE em respostas-etapa3c.md §6; rotas JWT permanecem **DN-009-2**.
- **Alternatives considered**: Uma consulta por lote com lista de assinaturas (opção A) — **rejeitada** por **DN-008-7 B** (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); origem: N2 B, LLT-200).

---

## R19 — Serialização da assinatura (008 assinatura A)

- **Decision**: Valor da assinatura em **AchadoAnalise**: **texto** concatenando as partes acordadas (**DN-008-1 B**) com **separador fixo** **`|`** (**DN-008-8 A**; dono canônico spec **008**; propagado nesta feature — **FR-F11**); **`|`** literal no componente → **`\|`** (**DN-008-9 A**) (fonte: respostas-joao.md, DN-008-8 A; §Etapa 3g DN-008-9 A; origem: §Etapa 3e (João, 16/11 12h40), 008 assinatura A).
- **Rationale**: Legibilidade em logs vs hash; João escolheu opção A da etapa 3c; literal do separador fechado na etapa 3e.
- **Alternatives considered**: SHA-256 somente (opção B) — **rejeitada** por **008 assinatura A** (fonte: respostas-joao.md §Etapa 3c 11h; §Etapa 3e (João, 16/11 12h40), DN-008-8 A; origem: DN-008-1 B, 008-11).

---

## R14 — Estado do repositório vs desenho

- **Decision**: Implementação atual consome Redis Streams, valida envelope v1 (`internal/domain/envelope.go`), expõe health HTTP; **não há** `ResumoAnalise`, adaptador HTTP ao Módulo nem símbolo `finishScan` no código Go.
- **Rationale**: Observação em `go.mod`, layout `internal/` (README § Layout), grep em `*.go`.
- **Alternatives considered**: N/A — fato de repositório, não decisão de produto.

---

## Decisões etapa 3e–3h (antes NEEDS CLARIFICATION — fechadas; etapa 3h revisa alcance OpenObserve/DLQ e **DN-008-9 A**)

| ID | Tópico | Fonte |
|----|--------|-------|
| **DN-008-5 A** | Literais `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` | respostas-joao.md §Etapa 3c 11h; spec.md Clarifications 3e |
| **DLQ B** / **DN-008-3** / **DN-008-6** (rev.) | Sem DLQ ingested; 3× OpenObserve → Finder→Core API (**batchId**, motivo) | respostas-joao.md §Etapa 3g; respostas-etapa3g.md §1 |
| **DN-008-9 A** | Escape **`|`** nos componentes → **`\|`** | respostas-joao.md §Etapa 3g (João, 16/11 13h22) |
| **DN-008-7** | Envio: um ResumoAnalise por lote | respostas-etapa3c.md §6 |
| **DN-008-7 B** | Consulta N2 B: uma chamada por achado | respostas-joao.md §Etapa 3c (João, 16/11 11h02) |
| **008 assinatura A** / **DN-008-8 A** | Concatenação com separador **`|`** (dono canônico spec **008**) | respostas-joao.md §Etapa 3c 11h; **DN-008-8 A**: §Etapa 3e (João, 16/11 12h40) |

## Decisões etapa 3b–3c (antes NEEDS CLARIFICATION — fechadas)

| ID | Tópico | Fonte |
|----|--------|-------|
| **DN-007-2 B** | ≥95% precisão **por classe**; gate na **006**/**007**, não na 009 | respostas-joao.md §Etapa 3b; spec.md Clarifications 3c |
| **DN-008-4 B** | Normalização na assinatura (**FR-F11**); dono spec **008** | respostas-joao.md §Etapa 3b |
| 007 formato **A** / caminho **C** | Base rotulada **JSONL** em `datasets/labeled-logs/` (só documentar) | respostas-joao.md §Etapa 3b |
| **spam: B** | Contador INFO no ResumoAnalise (substitui spam ignorado); rótulo ADR_102 pendente | respostas-joao.md §Etapa 3 |

## Lacunas não resolvidas (registradas, não pesquisadas por suposição)

| ID | Tópico | Onde registrado | Escopo bloqueado |
|----|--------|-----------------|------------------|
| **DN-009-2** | API Finder → Módulo (rota, JWT, respostas, 5xx) para **POST resumo** e **consulta por achado** | spec.md; plan.md § API do Módulo | Adaptador HTTP (**FR-F01**, **FR-F09** implementação) |
| **N3 A** | Campos ampliados AchadoAnalise | plan.md Pendências | Campos no payload até ADR |
| — | Formato confirmação HTTP Módulo | contracts/resumo-analise-handoff.md | Parcialmente **DN-009-2** |
| — | URL/rota HTTP Core (Módulo → Core) | contracts/pipeline-aviso-referencia.md | Fora do Finder |
| **ADR_051 [PENDENTE]** | API **Finder → Core** após falha OpenObserve (**batchId**, motivo) | spec.md **FR-F07**; plan.md § Pendências | Cliente HTTP operacional no Finder; **sem** DLQ ingested |
| — | Escape de **`\`** nos componentes da assinatura | spec.md §Questões em aberto; data-model.md | **DN-008-9 A** cobre só **`|`** → **`\|`**; sem fonte para `\` |