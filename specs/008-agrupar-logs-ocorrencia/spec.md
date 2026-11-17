# Feature Specification: Agrupar logs semelhantes na mesma ocorrência

**Feature Branch**: `008-agrupar-logs-ocorrencia` (trabalho permanece na branch git `Development`; nenhuma branch dedicada foi criada nesta execução, conforme instrução do solicitante)

**Created**: 2026-11-14

**Status**: Draft (clarificado em 2026-11-16, etapas 2b + 3b + 3c + 3c analyze + 3e + 3f + 3g + 3h, modo não interativo)

**Input**: Jira **LLT-200** (texto do solicitante); Epic Finder — IA própria de classificação e agregação; Feature **Agregação em ocorrências (E07.F3)**; item **E07.F3.1** em `LhamaLog_backlog_pendente.docx` (14/10/2026); contexto em `/workspace/lhamalog/contexto-LLT-103.md` §2.2 e §3; decisões do João (14–16/11/2026).

**Rastreabilidade**: Relacionado a **SSID_009**, **SSID_010**, **LLT-24**; depende da classificação da spec `006-modelo-classificacao-logs`; complementa o que a spec 006 deixa explicitamente fora (agrupamento por assinatura — spec 006 FR-009, FR-010).

## Clarifications

### Session 2026-11-16

- Q: T1 — Quem cria/grava a ocorrência? → A: O **Módulo de Ocorrência** cria a ocorrência, define gravidade, persiste, cuida de status e reabertura; o Finder envia **ResumoAnalise** (substitui T1 A). (fonte: respostas-joao.md, T1 revisto; origem: revisao-T1-drive.md A.1, ADR_062, sequencia-varredura-finder v6)
- Q: N2 — Quem junta logs iguais de lotes diferentes (contador, primeira/última vez)? → A: O **Finder**, consultando ocorrências existentes antes de enviar o resumo. (fonte: respostas-joao.md, N2; origem: rodada pós-revisão T1 16/11)
- Q: N3 — Ampliar AchadoAnalise com assinatura, contagem, primeira/última vez e versão das regras? → A: Sim, via **ADR nova** (pendência de criação; esta spec não cria a ADR). (fonte: respostas-joao.md, N3; origem: Finder.mmd × respostas João 008-12/008-14)
- Q: N4 — Chave de idempotência do ResumoAnalise: `arquivoId + execId` ou batchId? → A: **batchId + execId** (quando houver execId). (fonte: respostas-joao.md, N4; origem: contrato log.ingested v1, ADR_070)
- Q: T2 — Grão da ocorrência? → A: **Uma por achado** ERROR/WARNING/INFO (ADR_035/040), reconciliado com unicidade por assinatura via 008-18 A. (fonte: respostas-joao.md, T2; origem: ADR_035, ADR_040)
- Q: T3 — Semelhança/recorrência sem banco vetorial (ADR_168)? → A: **Consulta ao OpenObserve**. (fonte: respostas-joao.md, T3; origem: ADR_168, 008-3 B)
- Q: T4 — Lote/apontamento parcial? → A: **Spec 005**: lote incompleto não é sucesso; nada é entregue ao fluxo de agrupamento. (fonte: respostas-joao.md, T4; origem: spec 005 FR-007/008)
- Q: T7 — INFO vira ocorrência? → A: **Sim**, grava e agrupa como as demais (decisão consciente, diverge do Drive). (fonte: respostas-joao.md, T7/006-11 rodada 16/11; origem: SSID_007.1, SSID_009.1, ADR_040)
- Q: 008-18 — Como "achado" (T2) se concilia com sem duplicata por assinatura+produto? → A: **Achado = logs com a mesma assinatura no produto** (agrupados num AchadoAnalise no Finder). (fonte: respostas-joao.md, 008-18; origem: Jira LLT-200, ADR_040)
- Q: 008-1 — Campos e normalização da assinatura? → A: **Campos OTel explícitos** (lista a fechar; ver pendência). **↳ REVOGADA (vigência):** superseded by **DN-008-1 B** (etapa 3b) — ver Session 2026-11-16 (etapa 3b) e FR-002. (fonte: respostas-joao.md, 008-1; origem: Jira LLT-200)
- Q: 008-2 — correlationId/traceId fazem parte da assinatura? → A: **Sim, fazem parte** (risco: contrato 08/11 §4 — correlationId aleatório por exportação). **↳ REVOGADA (vigência):** decisão **008-2 A** substituída por **DN-008-1 B** / **008-2 rev.** na etapa 3b — ver Divergência **D2** e Clarifications etapa 3b. (fonte: respostas-joao.md, 008-2; origem: ADR_011, contrato §4 — histórico 15/11)
- Q: 008-3 — "Logs equivalentes" = assinatura idêntica? → A: **Similaridade via consulta ao OpenObserve** (não só igualdade literal no primeiro match). (fonte: respostas-joao.md, 008-3; origem: T3 B, ADR_168)
- Q: 008-4 — Retorno após resolvida: recorrente, reaberta ou ambos? → A: **Só reabertura FECHADA→ABERTA** (sem combinar com recorrente como substituto de status). (fonte: respostas-joao.md, 008-4; origem: ADR_082)
- Q: 008-5 — Critério de recorrência na consulta ao OpenObserve? → A: **Critério ADR_013** (3–4 vezes em arquivos/lotes distintos, mesma máquina; repetição no mesmo arquivo não conta para recorrência). (fonte: respostas-joao.md, 008-5; origem: ADR_013, ADR_063)
- Q: 008-6 — Janela da consulta de recorrência? → A: **Últimos 500 lotes** (análogo aos 500 arquivos da ADR_063). (fonte: respostas-joao.md, 008-6; origem: ADR_063)
- Q: 008-7 — Critério de semelhança da ADR_074 após ADR_168? → A: **Migrado ao OpenObserve** (detalhes operacionais em 008-3, 008-5, 008-6). (fonte: respostas-joao.md, T3 B + mapeamento perguntas-etapa2b; origem: ADR_074, ADR_168)
- Q: 008-8 — O que é "resolvida"? → A: Status **FECHADA** (ADR_082). (fonte: respostas-joao.md, 008-8; origem: ADR_082)
- Q: 008-9 — Quem fecha a ocorrência? → A: **Operador**, pela tela/Módulo Ocorrência. (fonte: respostas-joao.md, 008-9; origem: ADR_082)
- Q: 008-10 — Status muda ao atualizar o contador? → A: **Não**, salvo a regra do 008-4 (reabertura). (fonte: respostas-joao.md, 008-10; origem: ADR_082)
- Q: 008-11 — Chave de unicidade da ocorrência? → A: **(assinatura, classe, produto, tenant)** — revisão 16/11 (substitui assinatura+produto+tenant sem classe). (fonte: respostas-joao.md, 008-11 revisto; origem: rodada 16/11, concilia com 008-15 A)
- Q: 008-12 — Unidade do contador? → A: **Um log agregado** (cada log equivalente soma 1). (fonte: respostas-joao.md, 008-12; origem: Jira LLT-200 critério 1)
- Q: 008-13 — Repetição no mesmo lote/arquivo conta no contador? → A: **Sim** conta no contador (independente da ADR_013 para recorrência). (fonte: respostas-joao.md, 008-13; origem: ADR_013 × LLT-200)
- Q: 008-14 — Semântica de primeira/última vez? → A: **Timestamp do log** (não hora de processamento nem só receivedAt do lote). (fonte: respostas-joao.md, 008-14; origem: Jira LLT-200, contrato §5)
- Q: 008-15 — A assinatura inclui a classe? → A: **Sim; uma ocorrência por (assinatura, classe)** na chave de unicidade (classe pode também entrar na lista OTel do 008-1). (fonte: respostas-joao.md, 008-15 A; origem: ADR_035 tipo, T2 C)
- Q: 008-16 — Produto ausente ou inconsistente entre metadata e body? → A: **Rejeitar como falha crítica**, sem retry (padrão ADR_049). (fonte: respostas-joao.md, 008-16 A; origem: contrato 08/11 §4, ADR_049)
- Q: 008-17 — Política de re-notificação (ADR_035) entra na 008? → A: **Só documentar dependência** (motor de Notificação / LLT-216). (fonte: respostas-documentadas.md, 008-17; origem: backlog E10.F1.1, ADR_035)

### Session 2026-11-16 (etapa 3b)

- Q: DN-008-1 Quais campos compõem a assinatura de equivalência de logs? → A: **Nível + mensagem normalizada + `service.name` + host** (máquina). **Substitui 008-2 A**: `traceId` e `correlationId` **não** entram na assinatura; o risco do contrato v1 §4 (correlationId aleatório por exportação) **deixa de se aplicar à assinatura**. (fonte: respostas-joao.md, Etapa 3 / respostas-etapa3.md §4; origem: ADR_013, contrato log.ingested v1 §4)
- Q: DN-008-2 Como tratar ocorrências legadas da spec 006? → A: **Não há legado documentado; nada a migrar** (contrato v1 §11; LLT em "A fazer"; T1 revisto — Finder não grava ocorrência). **Ressalva:** ambiente de teste com ocorrências gravadas exige confirmação do João; então voltam as opções fundir/conviver/reprocessar. (fonte: respostas-joao.md, Etapa 3 / respostas-etapa3.md §5; origem: contrato v1 §11, T1 revisto)
- Q: DN-008-3 O que fazer se a consulta ao OpenObserve falhar (recorrência/similaridade)? → A: **Falha temporária**: até **3 tentativas** com reenvio no fim da mesma stream, depois **DLQ**; **erro crítico sem retry**. **MUST NOT** seguir o fluxo assumindo `recorrente=false`. **↳ SUPERADO (vigência):** etapa **3h** / **DLQ B** — ver Session 2026-11-16 (etapa 3h) e Decisões **DN-008-3** rev. 3h. (fonte: respostas-joao.md, Etapa 3 / respostas-etapa3.md §6; origem: ADR_146, contrato v1 §9, ADR_014)
- Q: API do Módulo Ocorrência (consulta do Finder: rota, JWT, respostas, 5xx/timeout) → A: **Pendência** até contrato/ADR do Módulo Ocorrência (N2 B e N1 A permanecem; detalhes HTTP não definidos no Drive). (fonte: respostas-joao.md, Etapa 3 / respostas-etapa3.md §10; origem: ADR_028, ADR_054, ADR_014)

### Session 2026-11-16 (etapa 3c)

- Q: DN-007-2 Qual é o escopo da meta de 95% de precisão na validação do classificador? → A: **Fora do escopo desta spec** como critério de aceite; referência cruzada: **≥ 95% de precisão em cada classe** (ERROR, WARNING, INFO) na spec **007** / gates da **006** (**DN-007-2 B**, **007-meta**). (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: DN-007-2 B; Jira LLT-199 critério 2; 006-2 B)
- Q: DN-008-4 Qual é o algoritmo de normalização da mensagem na assinatura? → A: **Aparar espaços**; **comparação sem distinção de maiúsculas/minúsculas**; **substituir** ocorrências de **números, UUIDs, IPs, datas e hexadecimais** por **marcadores** (categorias distintas). A spec **008** é **dona** do algoritmo; **006**, **009** e **010** MUST reutilizá-lo ao derivar equivalência/assinatura a partir da mensagem. (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: DN-008-4 B)

### Session 2026-11-16 (etapa 3c — correções analyze)

- Q: C1 O Princípio I da constituição (“Finder transforma em ocorrência”) conflita com o T1 revisto? → A: A implementação MUST seguir **T1 revisto** (Finder envia **ResumoAnalise**; **Módulo de Ocorrência** cria/persiste a ocorrência). A **emenda formal do princípio I** permanece **pendência externa** — MUST NOT editar a constituição nesta feature. (fonte: respostas-joao.md, T1 revisto; origem: ADR_062, revisao-T1-drive.md A.1; achado analyze **C1**)
- Q: I1 Qual decisão vigora sobre traceId/correlationId na assinatura após as etapas 3b e 3c? → A: Vigente **008-2 rev.** / **DN-008-1 B**: `traceId` e `correlationId` **não** entram na assinatura; o registro histórico **008-2 A** na Session 2026-11-16 inicial está **revogado** (ver **D2**). (fonte: respostas-joao.md §Etapa 3, DN-008-1 B; origem: contrato v1 §4; achado analyze **I1**)
- Q: A1 Como FR-008 (OpenObserve) se relaciona com o agrupamento determinístico por assinatura no lote? → A: **No lote em processamento**, logs equivalentes MUST agrupar pela assinatura **DN-008-1 B** + normalização **DN-008-4 B** (determinístico, local). A **consulta ao OpenObserve** aplica somente a **histórico/recorrência/similaridade fora do lote** (**008-3 B**, **008-5/6**); MUST NOT substituir nem relaxar o agrupamento local por assinatura. (fonte: respostas-joao.md, 008-3 B / DN-008-1 B; origem: T3 B, ADR_168; achado analyze **A1**)
- Q: U2 Os literais dos marcadores de normalização estão fechados? → A: A **semântica** está fechada (**DN-008-4 B**); os **literais exatos** dos cinco marcadores permanecem em **DN-008-5**, alinhados a **010-27 A** quando publicados — MUST NOT inventar (constituição V). (fonte: respostas-joao.md, 010-27 A; origem: spec §DN-008-5; achado analyze **U2**)
- Q: U1 API do Módulo Ocorrência (rota, JWT, HTTP, 5xx/timeout) → A: **Pendência DN-009-2** / questão em aberto §2; detalhes HTTP só após contrato/ADR do Módulo (João: «API Módulo: pendência no plan até o contrato»). (fonte: respostas-joao.md §Etapa 3; origem: respostas-etapa3.md §10; achado analyze **U1** — decisão não reaberta nesta sessão)
- Q: U3 Campos ampliados em AchadoAnalise → A: **Pendência N3 A** (ADR nova fora desta feature); nomes JSON finais permanecem `[PENDENTE ADR]`. (fonte: respostas-joao.md, N3 A; origem: ADR_062; achado analyze **U3**)
- Q: G1 Deve haver alerta ao Core além da DLQ quando a consulta OpenObserve esgota retries? → A: **Fechado na etapa 3e** — **DN-008-6**: após 3 falhas transitórias → **DLQ** → **alerta do lote** pelo **Core/WebSocket** (ver Session 2026-11-16 (etapa 3e)). **↳ SUPERADO (vigência):** etapa **3h** / **DLQ B** — ver Session 2026-11-16 (etapa 3h). (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h; achado analyze **G1**; correção etapa 3f / analyze-008-etapa3e.md **I1**)
- Q: G2 Qual a granularidade HTTP da consulta ao Módulo (lote inteiro vs por achado)? → A: **Fechado na etapa 3e** — **DN-008-7** / **DN-008-7 B** (envio **um ResumoAnalise/lote**; consulta prévia **uma chamada por achado**; **DN-008-7 B** em respostas-joao.md §Etapa 3c (João, 16/11 11h02)); contrato HTTP permanece **DN-009-2** (ver Session 2026-11-16 (etapa 3e)). O registro «sem decisão» desta sessão analyze está **superado**. (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; achado analyze **G2**; correção etapa 3f / analyze-008-etapa3e.md **I1**)

### Session 2026-11-16 (etapa 3e)

- Q: DN-008-5 Quais são os literais exatos dos marcadores do normalizador? → A: **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`** (cinco tokens distintos, um por categoria). (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h01 e 11h; respostas-etapa3c.md §4; origem: DN-008-5 A; DN-008-4 B)
- Q: DN-008-6 Após esgotar 3 tentativas e DLQ na consulta OpenObserve, o usuário recebe alerta além da DLQ? → A: **Sim**: após 3 falhas transitórias → **DLQ**, o **Core** envia **alerta do lote** ao usuário por **WebSocket** (não push/e-mail). **↳ SUPERADO (vigência):** etapa **3h** / **DLQ B** — Finder avisa Core (API interna, **batchId** + **motivo**); ver Session 2026-11-16 (etapa 3h). (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h; respostas-etapa3c.md §5; origem: ADR_059, ADR_073, ADR_046, ADR_051, ADR_055, ADR_060; confirmação João «DN-008-6 sim»)
- Q: DN-008-7 Qual a granularidade do handoff e da consulta prévia ao Módulo Ocorrência? → A: **Envio:** **um** `ResumoAnalise` **por lote** (ADR_062, Finder.mmd, N4 A). **Consulta prévia (N2 B):** **uma chamada por achado** (**DN-008-7 B**). Rotas, JWT e códigos HTTP da consulta permanecem em **DN-009-2**. (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; respostas-etapa3c.md §6; origem: ADR_062, N2 B, N1 A)
- Q: 008-assinatura A Como serializar o **valor** canônico da assinatura para comparação e chave? → A: **Texto** com os campos **nível + mensagem normalizada + `service.name` + host** concatenados por **um separador fixo** (opção A). O **separador exato** não tem fonte — ver **DN-008-8** (MUST NOT escolher por suposição). (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h01 e 11h; respostas-etapa3c.md §9; origem: DN-008-1 B, Finder.mmd `log_signature`)

### Session 2026-11-16 (etapa 3f)

- Q: Remover divergência documentada entre registros Etapa 3c 11h01 × 11h (D11 / «puladas»)? → A: **Não há divergência** — a linha «puladas» em respostas-joao.md está **SUPERADA** (16/11, linha 39); **D11** retirada de «Divergências documentadas»; **DN-008-7 B** cita **§Etapa 3c (João, 16/11 11h02)**. (fonte: respostas-joao.md, Etapa 3c 11h01 linha 39; instrução etapa 3f)
- Q: Harmonizar achado analyze **I1** (G1/G2 «sem decisão» vs etapa 3e)? → A: **DN-008-6** e **DN-008-7** fechadas na **etapa 3e**; entradas **G1/G2** na sessão «etapa 3c — correções analyze» atualizadas como **superadas** ( **DN-008-8** inalterada). (fonte: analyze-008-etapa3e.md, achado **I1**; respostas-joao.md §Etapa 3c 11h, DN-008-6 / DN-008-7 B)

### Session 2026-11-16 (etapa 3g)

- Q: DN-008-8 Qual é o separador fixo na serialização do valor canônico da assinatura (**008-assinatura A**)? → A: O caractere **`|`** (barra vertical U+007C), entre **nível**, **mensagem normalizada**, **`service.name`** e **host**, nessa ordem (ex.: `ERROR|msg normalizada|meu-serviço|host-01`). (fonte: respostas-joao.md, DN-008-8 A; origem: §Etapa 3e (João, 16/11 12h40); instrução etapa 3g)

### Session 2026-11-16 (etapa 3h)

- Q: DLQ B — Há DLQ no `log.ingested` quando a consulta ao OpenObserve falha? → A: **Não**; o contrato **log.ingested v1** não herda retry/DLQ do guia da FilaOrdenada (**DN-006-1b B**; respostas-etapa3g.md §1). **DN-008-3** (rev.): até **3 tentativas** na consulta; se esgotar, o Finder avisa o **Core** pela **API interna** com **batchId** e **motivo** — **sem DLQ**. (fonte: respostas-joao.md, §Etapa 3g (João, 16/11 13h22); origem: DLQ B)
- Q: DLQ B — Qual o alcance vigente de **DN-008-6**? → A: **Substitui o alcance da etapa 3e** (DLQ + alerta WebSocket pós-DLQ): após **3 falhas** na consulta OpenObserve, o Finder MUST notificar o **Core** pela **API interna** (**batchId** + **motivo**). Contrato da API Core (rota, payload, evento): **ADR_051 [PENDENTE]** — MUST NOT inventar. (fonte: respostas-joao.md, §Etapa 3g (João, 16/11 13h22); origem: DLQ B)
- Q: DN-008-9 A — Como tratar **`|`** literal dentro de um componente do `valor`? → A: Escapar como **`\|`** (barra invertida seguida de barra vertical). Escape de **barra invertida** literal nos componentes: **sem fonte** — apenas pendência (MUST NOT decidir). (fonte: respostas-joao.md, §Etapa 3g (João, 16/11 13h22); origem: DN-008-9 A)

## Decisões tomadas (com fonte)

### Transversais (impacto nesta feature)

| ID | Decisão | Fonte arquivo | Origem documental |
|----|---------|---------------|-------------------|
| T1 (rev.) | Finder envia **ResumoAnalise**; **Módulo Ocorrência** cria/persiste ocorrências, status e reabertura | respostas-joao.md | ADR_062; revisao-T1-drive.md A.1; sequencia-varredura-finder v6 |
| T2 C | Grão: **uma ocorrência por achado** (tipo ERROR/WARNING/INFO), fundido por assinatura conforme 008-18 A | respostas-joao.md | ADR_035, ADR_040 |
| T3 B | Semelhança/recorrência histórica: **consulta ao OpenObserve** (pós-ADR_168) | respostas-joao.md | ADR_168, ADR_074 |
| T4 A | **Lote incompleto** (spec 005): não é sucesso; **não** prosseguir agrupamento com aquele lote | respostas-joao.md | spec 005 FR-007/008 |
| T5 D | Contrato **log.ingested v1** (08/11) é a versão válida do evento de entrada | respostas-joao.md | contrato-log-ingested-v1.txt |
| T7 A | Classe **INFO** gera achado e ocorrência como ERROR/WARNING (agrupamento inclusivo) | respostas-joao.md | ADR_170; decisão consciente 16/11 |
| N2 B | **Finder** agrega logs **entre lotes** (contador, primeira/última vez), **consultando** ocorrências antes do envio | respostas-joao.md | LLT-200; rodada 16/11 |
| N3 A | Ampliar **AchadoAnalise** com assinatura, contagem, primeira/última vez, versão das regras — **ADR pendente** | respostas-joao.md | ADR_062, Finder.mmd |
| N4 A | Chave idempotência **ResumoAnalise**: `batchId` (+ `execId` se houver) | respostas-joao.md | contrato v1; ADR_070 (histórico arquivoId) |

### Taxonomia (João 14/11)

- Classes válidas: **erro**, **alerta**, **informação** (ERROR, WARNING, INFO). **Spam** e **"não classificado"** MUST NOT existir. (fonte: respostas-joao.md, constituição III; origem: ADR_170)

### Feature 008 (LLT-200)

| ID | Decisão | Fonte arquivo | Origem documental |
|----|---------|---------------|-------------------|
| 008-18 A | **Achado** = conjunto de logs com a **mesma assinatura** no **produto**, materializado num `AchadoAnalise` | respostas-joao.md | Jira LLT-200; ADR_040 |
| 008-1 B (fech. 3b) | Assinatura = **nível** + **mensagem normalizada** + **`service.name`** + **host** (máquina) | respostas-joao.md Etapa 3 | Jira LLT-200; ADR_013; DN-008-1 B |
| **DN-008-4 B** (fech. 3c) | **Normalização da mensagem**: aparar espaços; ignorar maiúsculas/minúsculas; trocar números, UUIDs, IPs, datas e hexadecimais por marcadores. **Dono:** spec 008; **consumidores:** 006, 008, 009, 010 | respostas-joao.md §Etapa 3b, João 16/11 09h02 | DN-008-4 B; ADR_013 ("praticamente idêntico") |
| **DN-008-5 A** (fech. 3e) | Literais dos marcadores: **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`** | respostas-joao.md §Etapa 3c 11h01/11h | DN-008-5 A; DN-008-4 B |
| **DN-008-3** (fech. 3b; **rev. 3h**) | Falha **temporária** na consulta OpenObserve: até **3 tentativas**; **sem DLQ** no `log.ingested` (**DLQ B**). Após esgotar, Finder MUST avisar o **Core** por **API interna** com **batchId** e **motivo** (**ADR_051 [PENDENTE]**). MUST NOT `recorrente=false` por degradação. **Substitui** alcance anterior com DLQ (etapas 3b/3e). | respostas-joao.md §Etapa 3g (João, 16/11 13h22); origem: DLQ B; respostas-etapa3g.md §1 |
| **DN-008-6** (fech. 3h; **substitui 3e**) | Após **3 falhas** na consulta OpenObserve, Finder MUST notificar o **Core** pela **API interna** (**batchId** + **motivo**); **não** DLQ nem alerta WebSocket pós-DLQ neste fluxo. Contrato Finder→Core: **ADR_051 [PENDENTE]**. | respostas-joao.md §Etapa 3g (João, 16/11 13h22); origem: DLQ B |
| **DN-008-7** (fech. 3e, parcial) | **Envio** um **ResumoAnalise**/lote; **consulta prévia N2 B** **por achado** (DN-008-7 B). Contrato HTTP: **DN-009-2** | respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; respostas-etapa3c.md §6 | ADR_062, Finder.mmd, N2 B, N4 A |
| **DN-008-8 A** (fech. 3g) | Separador na serialização de `valor`: **`|`** entre nível, mensagem normalizada, `service.name`, host (**008-assinatura A**) | respostas-joao.md §Etapa 3e (João, 16/11 12h40) | DN-008-8 A |
| **DN-008-9 A** (fech. 3h) | **`|`** literal dentro de qualquer componente do `valor` MUST ser escapado como **`\|`** antes da concatenação | respostas-joao.md §Etapa 3g (João, 16/11 13h22) | DN-008-9 A |
| **008-assinatura A** (fech. 3e + 3g) | Valor da assinatura = texto concatenado (campos **DN-008-1 B**) com separador **`|`** (**DN-008-8 A**) | respostas-joao.md §Etapa 3c 11h01/11h; DN-008-8 A | DN-008-1 B; respostas-etapa3c.md §9 |
| 008-2 rev. | **008-2 A substituída (etapa 3b):** `traceId` e `correlationId` **não** fazem parte da assinatura | respostas-joao.md Etapa 3 | contrato v1 §4 (risco não aplica à assinatura) |
| 008-3 B | Equivalência histórica: **similaridade via OpenObserve**, não só hash local | respostas-joao.md | T3 B, ADR_168 |
| 008-4 B | Ocorrência **resolvida** que volta: **FECHADA → ABERTA** (reabertura pelo Módulo Ocorrência) | respostas-joao.md | ADR_082, Jira LLT-200 crit. 3 |
| 008-5 A | Marca **recorrente**: critério **ADR_013** na consulta OpenObserve | respostas-joao.md | ADR_013, ADR_063 |
| 008-6 A | Janela da consulta: **últimos 500 lotes** | respostas-joao.md | ADR_063 (500 arquivos) |
| 008-7 B | Critério ADR_074 **[PENDENTE]** **substituído** por consulta OpenObserve (008-3/5/6) | respostas-joao.md, perguntas-etapa2b | ADR_074, ADR_168 |
| 008-8 A | **Resolvida** = status **FECHADA** | respostas-joao.md | ADR_082 |
| 008-9 A | **Operador** fecha ocorrência (UI Módulo Ocorrência) | respostas-joao.md | ADR_082 |
| 008-10 A | Atualizar contador **não** altera status (exceto reabertura 008-4) | respostas-joao.md | ADR_082 |
| 008-11 rev. | Unicidade: **(assinatura, classe, produto, tenant)** | respostas-joao.md 16/11 | LLT-200 (ampliado), ADR_074, tenant contrato §4 |
| 008-12 A | Contador: **+1 por log** agregado | respostas-joao.md | Jira LLT-200 crit. 1 |
| 008-13 A | Repetição no **mesmo lote** **conta** no contador | respostas-joao.md | ADR_013 (só recorrência) × LLT-200 |
| 008-14 A | **Primeira/última vez** = extremos do **timestamp do log** | respostas-joao.md | Jira LLT-200; contrato §5 |
| 008-15 A | **Uma ocorrência por (assinatura, classe)**; classe na chave | respostas-joao.md | ADR_035 `tipo` |
| 008-16 A | `productId` ausente ou metadata≠body: **falha crítica**, sem retry | respostas-joao.md | contrato §4; ADR_049 |
| 008-17 C | Política ADR_035 de não repetir notificação: **fora da 008**; dependência LLT-216/009 | respostas-documentadas.md | ADR_035; E10.F1.1 |

### Jira e ADRs (inalterados como contexto)

- **LLT-200** critérios 1–3: contador, primeira/última vez; sem duplicata (chave ampliada pelo João); retorno após resolvida. (fonte: Jira)
- **ADR_082**: status ABERTA, EM_ANALISE, FECHADA; FECHADA→ABERTA. (fonte: Drive)
- **ADR_063**: **recorrente** é atributo, não status. (fonte: Drive)
- **ADR_062**: Finder **não** conhece Ocorrência; envia ResumoAnalise. (fonte: Drive)
- **N1 A** (entrega ResumoAnalise): HTTP interno com JWT — pertinente ao handoff; detalhe de transporte registrado em respostas-joao 16/11 (implementação em contrato/ADR, fora desta spec).

### Responsabilidades após clarificação

- **Finder**: calcular assinatura (**nível + mensagem normalizada + service.name + host**, 008-1 fech. 3b; normalização **DN-008-4 B**; marcadores **DN-008-5 A**; serialização texto **008-assinatura A** com separador **`|`** — **DN-008-8 A**); agrupar logs no **mesmo lote** num `AchadoAnalise`; **consultar** ocorrências existentes **por achado** (N2 B, **DN-008-7 B**; rotas/JWT em **DN-009-2**) e **agregar entre lotes**; preencher contagem e primeira/última vez no achado (N3, quando ADR existir); detectar **recorrente** via OpenObserve (008-3/5/6), com até **3 tentativas** na consulta e, se esgotar, aviso ao **Core** por API interna (**batchId** + **motivo**; **DN-008-3**, **DN-008-6** rev. 3h; **sem DLQ**); enviar **um ResumoAnalise por lote** (HTTP JWT, N1 A, **DN-008-7**); rejeitar lote com produto inválido (008-16 A); respeitar T4 A (não processar lote incompleto).
- **Módulo Ocorrência**: criar ou atualizar ocorrência pela chave 008-11; aplicar contador/extremos recebidos no achado; manter **recorrente** separado de status; **reabrir** FECHADA→ABERTA quando aplicável (008-4); permitir fechamento pelo operador (008-9); aceitar achados **INFO** (T7 A).

## Divergências documentadas

| Id | Fontes | O que diverge | Lado que vale nesta feature | Por quê |
|----|--------|---------------|----------------------------|---------|
| D1 | João T7 A, 006-11; SSID_007.1 crit.1; SSID_009.1 crit.1; ADR_040 | **INFO vira ocorrência** × Drive (só ERRO/WARNING) | **João (T7 A)** | Decisão consciente 16/11; ADR_170 abre caminho; execução no Módulo Ocorrência |
| D2 | João 008-2 A (15/11); João DN-008-1 B (16/11 Etapa 3); contrato v1 §4 | **correlationId/traceId na assinatura** (008-2 A) × **assinatura sem esses campos** (DN-008-1 B) | **João DN-008-1 B (etapa 3b)** | Fecha a lista da assinatura alinhada a ADR_013 (mesma máquina); o aviso do contrato §4 sobre correlationId **não se aplica à composição da assinatura** |
| D2b | respostas-etapa3.md §5; contrato v1 §11 | **Sem legado em produção** × cenários hipotéticos de migração (relatório etapa 3) | **Drive + T1 rev. (DN-008-2)** | **Ressalva:** dados em ambiente de teste não cobertos pela §11 — confirmação do João se existirem ocorrências de teste |
| D3 | João T1 rev.; T1 A antiga; spec 006 clarificação Q2 | Finder gravava no MongoDB × **Módulo Ocorrência cria** | **João T1 rev. + Drive (ADR_062)** | T1 A substituída 15/11 |
| D4 | João 008-11/008-15; Jira LLT-200 crit. 2 | Unicidade **(assinatura, classe, produto, tenant)** × Jira **assinatura + produto** | **João** | Critério Jira satisfeito de forma mais restrita (sem duplicata lógica) |
| D5 | João 008-11 rev. × 008-11 B anterior | Chave **com classe** × chave sem classe (15/11) | **João 16/11** | Concilia 008-11 com 008-15 A |
| D6 | ADR_019, ADR_167, SSID_007.1 (texto antigo); João 14/11 | Taxonomia com spam × três classes | **João** | Emenda documental fora desta spec |
| D7 | ADR_146 **v2** (schemaVersion `"2"`); contrato 08/11 **v1**; João T5 D | Versão do envelope do consumidor × contrato publicado | **Contrato 08/11 v1 (T5 D)** | Única versão válida para `log.ingested` neste pipeline |
| D8 | spec 006 FR-010 (só inserir); LLT-200 (atualizar) | Inserção por apontamento × upsert por assinatura | **LLT-200 + T1 rev.** | 006 deve emitir ResumoAnalise; fusão no Módulo Ocorrência (não no Finder MongoDB) |
| D9 | respostas-documentadas 008-17; ADR_035 [PENDENTE] | Agrupamento na 008 × política de notificação duplicada | **008-17 C** | Notificação (LLT-216), não agrupamento |
| D10 | ADR_062 (Finder não envia gravidade); SSID_009.1 passo 1 | Gravidade no aviso × Finder sem criticidade | **Drive com T1 rev.** | Módulo Ocorrência calcula e envia ao Core |

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Unificar logs equivalentes em uma única ocorrência (Priority: P1)

O sistema recebe logs já classificados (erro, alerta ou informação — T7 A). O Finder agrupa logs da **mesma assinatura** no **mesmo produto** num único **achado**, atualiza contagem e extremos temporais e, **entre lotes**, consulta ocorrências existentes antes de enviar o **ResumoAnalise** (N2 B). O Módulo Ocorrência **atualiza** a ocorrência existente para a chave **(assinatura, classe, produto, tenant)** em vez de criar duplicata (008-11, FR-004/FR-005).

**Why this priority**: Núcleo LLT-200 / E07.F3.1.

**Independent Test**: Enviar dois lotes com logs equivalentes (mesma assinatura, classe, produto, tenant); verificar **uma** ocorrência com contador, primeira e última vez corretos (timestamps dos logs).

**Acceptance Scenarios**:

1. **Given** ocorrência existente para chave **(A, ERROR, P, T)**, **When** chega novo log equivalente no próximo lote, **Then** o Finder consulta, agrega no AchadoAnalise e o Módulo incrementa contador e ajusta primeira/última vez (008-12, 008-14).
2. **Given** N logs equivalentes na mesma chave, **Then** existe **exatamente uma** ocorrência consultável (LLT-200 crit. 2; chave 008-11).
3. **Given** mesma assinatura em **produtos** ou **classes** diferentes, **Then** ocorrências **distintas** (008-15 A, 008-11).
4. **Given** lote **incompleto** (spec 005), **When** o Finder recebe o apontamento, **Then** **não** agrupa nem envia resumo para aquele lote (T4 A).
5. **Given** logs com a mesma mensagem salvo **maiúsculas**, **espaços nas pontas** ou **números/UUID/IP/data/hex** distintos (**DN-008-4 B**), **When** nível, `service.name`, host, classe, produto e tenant coincidem, **Then** **mesma assinatura** e **uma** ocorrência (FR-002b, SC-001).

---

### User Story 2 - Evitar duplicatas pela chave acordada (Priority: P1)

Consumidores não encontram várias ocorrências para o mesmo problema lógico na mesma chave de unicidade.

**Acceptance Scenarios**:

1. **Given** N ≥ 2 logs na mesma chave 008-11, **Then** consulta retorna **uma** ocorrência.
2. **Given** pipeline documentado (contrato v1 §11, LLT em "A fazer", T1 rev.), **When** esta feature entra em vigor, **Then** **não há migração de ocorrências legadas** (DN-008-2), salvo confirmação explícita para ambiente de teste (D2b).

---

### User Story 3 - Tratar retorno após resolução (Priority: P1)

**Resolvida** = **FECHADA** (008-8). Novo log da mesma chave após FECHADA: **reabertura FECHADA→ABERTA** (008-4), sem mudar status só por contador (008-10). Atributo **recorrente** segue ADR_013 via consulta OpenObserve (008-5/6) e permanece separado de status (ADR_063).

**Acceptance Scenarios**:

1. **Given** ocorrência **FECHADA** para chave **(A, W, P, T)**, **When** chega log equivalente, **Then** status passa a **ABERTA** (008-4) e **recorrente** é avaliado conforme 008-5 (sem tratar recorrente como status).
2. **Given** ocorrência **ABERTA** com contador atualizado, **Then** status permanece **ABERTA** (008-10).

---

### Edge Cases

- **Normalização da mensagem** (**DN-008-4 B**, **DN-008-5 A**): aparar espaços; case-insensitive; substituições pelos marcadores **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`**; matriz de casos em **010-27 A** / contrato **assinatura-log**.
- **`|`** no texto de componente após normalização: separador de `valor` é **`|`** (**DN-008-8 A**); **`|`** literal MUST ser escapado como **`\|`** (**DN-008-9 A**). Escape de **barra invertida** literal nos componentes: **pendente** (sem fonte).
- **Consulta OpenObserve indisponível** (recorrência/similaridade): falha **temporária** — até **3 tentativas** na consulta; **sem DLQ** no `log.ingested` (**DLQ B**). Após esgotar, o Finder MUST avisar o **Core** pela **API interna** com **batchId** e **motivo** (**DN-008-3**, **DN-008-6** rev. 3h; contrato **ADR_051 [PENDENTE]**); MUST NOT degradar para `recorrente=false`. Erro **crítico** (ex.: parsing): **sem retry** (ADR_049).
- **productId** ausente ou metadata≠body: rejeição crítica, sem retry (008-16 A).
- **Classe INFO**: gera ocorrência e entra na chave (T7 A); diverge de SSID_007.1/009.1 (D1).
- **Mesma assinatura, classes diferentes**: ocorrências distintas (008-15 A).
- **Repetição no mesmo lote**: conta no contador (008-13 A); não conta para recorrência ADR_013 no mesmo arquivo/lote.
- **API consulta/atualização Módulo Ocorrência** (N2 B): **consulta prévia por achado** (**DN-008-7 B**); **envio** **um ResumoAnalise por lote** (**DN-008-7**). Rota, claims JWT, códigos HTTP e timeout **pendentes** (**DN-009-2**; respostas-etapa3.md §10).
- **ADR N3 (campos no AchadoAnalise)**: pendente; até lá, entrega do achado ao Módulo depende do contrato ampliado.
- **Ambiente de teste** com ocorrências já gravadas: fora da garantia documental §11 — ver D2b.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: O **Finder** MUST agrupar logs por assinatura no lote, **consultar** ocorrências existentes e **agregar entre lotes**, preenchendo contagem e primeira/última vez no `AchadoAnalise` (N2 B, N3 A quando ADR vigente). O **Módulo Ocorrência** MUST **criar ou atualizar** a ocorrência persistida pela chave **(assinatura, classe, produto, tenant)** (T1 rev., 008-11). O Finder MUST NOT gravar ocorrência no MongoDB.
- **FR-002**: A assinatura MUST ser calculada como **nível** + **mensagem normalizada** + **`service.name`** + **host** (008-1 B fech. etapa 3b, DN-008-1 B; ADR_013 mesma máquina). **traceId** e **correlationId** MUST NOT integrar a assinatura (**008-2 A substituída**). A **classe** de negócio (erro/alerta/informação) entra na **chave de unicidade** 008-11, não necessariamente duplicando o nível OTel se já parte da assinatura (008-15 A).
- **FR-002b**: A **mensagem normalizada** MUST aplicar, nesta ordem lógica: **(1)** aparar espaços nas extremidades; **(2)** comparar e agrupar sem distinção de maiúsculas/minúsculas; **(3)** substituir ocorrências de **números**, **UUIDs**, **endereços IP**, **datas** e **valores hexadecimais** pelos marcadores **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`** (**DN-008-5 A**), distintos entre categorias. O algoritmo MUST ser o **canônico** desta feature (**DN-008-4 B**); specs **006**, **009** e **010** MUST NOT definir regra divergente ao derivar assinatura/equivalência a partir da mensagem. Casos de teste MUST alinhar a **010-27 A** e ao contrato **assinatura-log** sem alterar os cinco literais acima.
- **FR-002c**: O **valor canônico** da assinatura (comparação e chave lógica antes da tupla 008-11) MUST ser **texto** formado concatenando, **nesta ordem**, **nível**, **mensagem normalizada**, **`service.name`** e **host**, com o separador fixo **`|`** (barra vertical, **DN-008-8 A**, **008-assinatura A**). **`|`** literal dentro de qualquer componente MUST ser escapado como **`\|`** (**DN-008-9 A**). Regra de escape de **barra invertida** literal nos componentes: **pendente** (sem fonte; MUST NOT inventar).
- **FR-002a**: Em falha **temporária** da consulta ao OpenObserve (similaridade/recorrência), o Finder MUST aplicar até **3 tentativas** na consulta e, se esgotar, MUST avisar o **Core** pela **API interna** com **batchId** e **motivo** (**DN-008-3**, **DN-008-6** rev. 3h; **DLQ B**; contrato v1 §8; **ADR_051 [PENDENTE]** — MUST NOT inventar rota/payload/evento). **MUST NOT** usar **DLQ** no fluxo `log.ingested`. MUST NOT concluir o lote com `recorrente=false` por indisponibilidade da consulta. Erro **crítico** (parsing/leitura) MUST NOT ter retry (**ADR_049**; 008-16 A).
- **FR-003**: Quando ocorrência **FECHADA** (008-8) recebe log da mesma chave 008-11, o Módulo Ocorrência MUST aplicar **FECHADA → ABERTA** (008-4 B). O Finder MUST avaliar **recorrente** conforme **ADR_013** em consulta ao OpenObserve na janela dos **últimos 500 lotes** (008-5 A, 008-6 A, T3 B). **Recorrente** MUST NOT ser status (FR-007).
- **FR-004**: Para chave 008-11 existente, o Módulo Ocorrência MUST **atualizar** contador, primeira e última vez a partir do achado (008-12, 008-14), não criar duplicata (LLT-200 crit. 1–2).
- **FR-005**: O sistema MUST NOT manter duas ou mais ocorrências distintas para a mesma chave **(assinatura, classe, produto, tenant)** (008-11, 008-15).
- **FR-006**: Classes MUST ser **erro**, **alerta** ou **informação** apenas (João 14/11; constituição III). Achados **INFO** MUST ser aceitos pelo Módulo Ocorrência (T7 A).
- **FR-007**: **Recorrente** MUST NOT ser tratado como status; status válidos: **ABERTA**, **EM_ANALISE**, **FECHADA** (ADR_082, ADR_063).
- **FR-008**: **Agrupamento no lote em curso** MUST usar assinatura determinística (**FR-002**, **DN-008-1 B**, **DN-008-4 B**) — independente do OpenObserve (**A1**, etapa 3c analyze). Equivalência **histórica** além desse lote MUST usar **consulta ao OpenObserve** (008-3 B, 008-7 B, T3 B), com critério de recorrência ADR_013 na janela 008-6; a consulta MUST NOT relaxar nem substituir o agrupamento local por assinatura. O mecanismo de banco vetorial/Elasticsearch no formato ADR_152 MUST NOT ser reintroduzido (ADR_168).
- **FR-009**: O Finder MUST enviar **exatamente um** `ResumoAnalise` **por lote** processado com achados agrupados (T1 rev., 008-18 A, **DN-008-7**; ADR_062). Antes do envio, para agregar entre lotes (**N2 B**), o Finder MUST **consultar** ocorrências existentes com **uma chamada por achado** (**DN-008-7 B**); detalhes HTTP da consulta MUST seguir **DN-009-2** quando publicados. A spec 006 MUST evoluir para **não** inserir ocorrência direta (conflito D8). **Migração de ocorrências legadas:** MUST NOT ser exigida em produção documentada (DN-008-2; contrato v1 §11), com ressalva D2b para teste.
- **FR-010**: Atualização de contador MUST NOT alterar status da ocorrência, exceto reabertura 008-4 (008-10 A).
- **FR-011**: Isolamento MUST incluir **tenant** na chave (008-11 rev.) além de produto; usuário, projeto, tipo e varredura (SSID_010) permanecem **fora** da chave salvo decisão futura.
- **FR-012**: **Resolvida** = **FECHADA**; fechamento pelo **operador** no Módulo Ocorrência (008-8 A, 008-9 A).
- **FR-013**: Entrada com **productId** ausente ou inconsistente entre metadata validada e corpo MUST ser **rejeitada** como falha crítica sem retry (008-16 A, ADR_049).
- **FR-014**: Lote de ingestão **incompleto** (spec 005) MUST NOT disparar agrupamento nem ResumoAnalise para esse lote (T4 A).
- **FR-015**: Política de **não repetir notificação** (ADR_035) MUST ser apenas **dependência** documentada para LLT-216/009 (008-17 C), não requisito desta feature.

### Fora de escopo (com fonte)

- Definir ou treinar **modelo de classificação** (spec 006; LLT-199); meta **≥ 95% de precisão por classe** (**DN-007-2 B**, **007-meta**) — aceite na spec **007**/gates na **006**, não nesta feature.
- Montar **base rotulada** (LLT-198; spec 007).
- **Publicar evento** ao Core, **notificar**, **mascarar**, **sugerir** (spec 009/010; ADR_051).
- Implementar **BancoVetorial** / Elasticsearch legado (ADR_168).
- Alterar **ADR_019**, **ADR_167**, **SSID_007.1** (só citação e divergências).
- **Criar ou editar ADRs** (incluindo ampliação AchadoAnalise — N3 A pendente).
- Resolver **política de notificação duplicada** (008-17 C → LLT-216).
- Meta de desempenho da consulta OpenObserve (sem fonte).

### Key Entities *(include if feature involves data)*

- **Log (classificado)**: Entrada com classe erro, alerta ou informação; o **nível** OTel participa da assinatura; participa do achado (008-18).
- **Assinatura**: Identificador lógico de equivalência = **nível + mensagem normalizada + service.name + host** (DN-008-1 B). **Mensagem normalizada** = saída do algoritmo **DN-008-4 B** com marcadores **DN-008-5 A**. **Valor serializado** = os quatro componentes em ordem, unidos por **`|`** (**DN-008-8 A**, **008-assinatura A**); **`|`** literal nos componentes: escape **`\|`** (**DN-008-9 A**). `traceId`/`correlationId` permanecem no log/ocorrência (ADR_040) mas **fora** da assinatura (008-2 rev.).
- **AchadoAnalise**: Agregação de logs de mesma assinatura no produto (008-18); deve carregar assinatura, contagem, primeira/última vez, versão das regras (**N3 A**, ADR pendente).
- **ResumoAnalise**: Contadores do lote + lista de achados; **um envelope por lote** (DN-008-7); entrega Finder → Módulo Ocorrência (T1 rev.); idempotência `batchId` + `execId` (N4 A).
- **Ocorrência (agrupada)**: Única por **(assinatura, classe, produto, tenant)**; atributos contador, primeira/última vez, recorrente (atributo), status ADR_082; criada/atualizada pelo **Módulo Ocorrência**.
- **Produto / Tenant**: Produto no corpo e token; tenant em metadata confiável (contrato §4); ambos na chave de unicidade.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Com a assinatura definida (FR-002, DN-008-1 B), N logs da mesma chave 008-11 resultam em **exatamente uma** ocorrência em 100% dos testes (LLT-200 crit. 2).
- **SC-002**: Sequência de M logs na mesma chave: contador = M agregações acordadas (008-12); primeira/última vez = min/max **timestamp do log** (008-14).
- **SC-003**: Ocorrência **FECHADA** + novo log equivalente: 100% aplicam **FECHADA→ABERTA** (008-4); recorrente conforme ADR_013 quando consulta retorna padrão (008-5).
- **SC-004**: **0%** de ocorrências com classe fora de erro/alerta/informação (João 14/11).
- **SC-005**: Enquanto a **ADR N3** (AchadoAnalise/contrato Finder→Módulo), o **contrato HTTP** da consulta (**DN-009-2**) ou o contrato Finder→Core para aviso após falha OpenObserve (**ADR_051 [PENDENTE]**) estiverem pendentes, o aceite formal de ponta a ponta não fecha (constituição V). A **composição da assinatura** (FR-002), o **algoritmo de normalização** e os **cinco literais** (FR-002b, **DN-008-4 B**, **DN-008-5 A**), o **separador** **`|`** (**DN-008-8 A**), o **escape** de **`|`** (**DN-008-9 A**) e a **forma de serialização** (**008-assinatura A**) estão fechados (etapas 3b–3h), salvo escape de **barra invertida** nos componentes (sem fonte).
- **SC-007**: Falha temporária na consulta OpenObserve: 100% dos cenários de teste MUST seguir **3 tentativas → aviso ao Core (API interna, batchId + motivo)** (**DN-008-3**, **DN-008-6** rev. 3h; **sem DLQ**) e MUST NOT registrar `recorrente=false` por degradação.
- **SC-006**: Dependência ADR_035/LLT-216 documentada; **não** mensurável nesta feature (008-17 C).

## Evidência e consulta

- **Fato**: Decisões com fonte em `respostas-joao.md`, `revisao-T1-drive.md`, `respostas-documentadas.md`, Jira LLT-200, ADRs citadas, contrato 08/11.
- **Premissa**: Nenhuma premissa inventada; lacunas sem fonte em Questões em aberto ou Dúvidas novas.
- **Ressalva checklist H1 (analyze rodada 2)**: o `spec.md` não contém marcadores literais `[NEEDS CLARIFICATION]`; o **plano** (`plan.md`) e demais artefatos podem manter **NEEDS CLARIFICATION** / **PENDENTE** até alinhamento pós-3h — especialmente **DN-009-2**, **ADR_051 [PENDENTE]**, detalhes operacionais OpenObserve (esta execução atualiza `spec.md` e `contracts/assinatura-log.md`).
- **Perguntas resolvidas nesta sessão**: **etapa 3h:** **DLQ B** (rev. **DN-008-3** / **DN-008-6**); **DN-008-9 A**. **Etapa 3g:** **DN-008-8 A** (separador **`|`**). **Etapa 3e:** **DN-008-5**, **DN-008-7** (envio/consulta), **008-assinatura A** (**DN-008-6** vigência substituída na 3h). **Etapa 3f:** remoção de **D11**; harmonização **G1/G2** / **I1**. Sessões anteriores: 008-1 a 008-18; etapa 3b DN-008-1–3; etapa 3c **DN-008-4**; etapa 3c analyze C1, I1, A1.

## Questões em aberto

1. **ADR** para ampliar `AchadoAnalise` / contrato Finder → Módulo Ocorrência (**N3 A**) — pendência explícita do João (achado analyze **U3**).
2. **API do Módulo Ocorrência** para consulta/atualização pelo Finder (**N2 B**, **DN-009-2**): rota, JWT (claims/algoritmo ADR_028 [PENDENTE]), respostas HTTP, política 5xx/timeout — **pendente até contrato do Módulo**. **Granularidade da consulta prévia** fechada: **por achado** (**DN-008-7 B**); **envio** fechado: **um ResumoAnalise por lote** (**DN-008-7**). (respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; respostas-etapa3.md §6, §10)
3. **Escape de barra invertida (`\`)** literal nos componentes do `valor` (**DN-008-9 A** fechou só **`|`** → **`\|`**; parse simétrico da `\` **sem fonte**).
4. **Contrato API Core** para aviso após 3 falhas na consulta OpenObserve (**ADR_051 [PENDENTE]**): rota, payload, evento — MUST NOT inventar (fonte: respostas-joao.md §Etapa 3g, DLQ B).
5. **Emenda do princípio I** da constituição (T1 revisto × texto vigente do princípio I): pendência **fora desta pasta** — MUST NOT editar constituição aqui (achado analyze **C1**).

## Dúvidas novas (sem fonte)

- **DN-008-1** — respondida em 2026-11-16 (etapa 3b), ver Clarifications → Session 2026-11-16 (etapa 3b).
- **DN-008-2** — respondida em 2026-11-16 (etapa 3b), ver Clarifications → Session 2026-11-16 (etapa 3b).
- **DN-008-3** — respondida em 2026-11-16 (etapa 3b); **alcance revisto** em 2026-11-16 (etapa 3h, **DLQ B**), ver Clarifications → Session 2026-11-16 (etapa 3h).
- **DN-008-4** — respondida em 2026-11-16 (etapa 3c), ver Clarifications → Session 2026-11-16 (etapa 3c).
- **DN-008-5** — respondida em 2026-11-16 (etapa 3e), ver Clarifications → Session 2026-11-16 (etapa 3e).
- **DN-008-6** — respondida em 2026-11-16 (etapa 3e); **alcance substituído** em 2026-11-16 (etapa 3h, **DLQ B**), ver Clarifications → Session 2026-11-16 (etapa 3h).
- **DN-008-7** — respondida em 2026-11-16 (etapa 3e), ver Clarifications → Session 2026-11-16 (etapa 3e); **DN-008-7 B** (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02)).
- **DN-008-8** — respondida em 2026-11-16 (etapa 3g), ver Clarifications → Session 2026-11-16 (etapa 3g); **DN-008-8 A** (fonte: respostas-joao.md, DN-008-8 A; origem: §Etapa 3e (João, 16/11 12h40)).
- **DN-008-9** — respondida em 2026-11-16 (etapa 3h), ver Clarifications → Session 2026-11-16 (etapa 3h); **DN-008-9 A**: **`|`** → **`\|`** (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22); origem: DN-008-9 A).

---

**Pendência ADR (não criar nesta execução):** ampliar `AchadoAnalise` com assinatura, contagem, primeira/última vez, versão das regras (**N3 A**).

**Pendência constitucional (não editar nesta execução):** emenda do **princípio I** (T1 revisto: Módulo de Ocorrência cria a ocorrência; Finder envia ResumoAnalise) — ver Questões em aberto §5 e achado analyze **C1**.