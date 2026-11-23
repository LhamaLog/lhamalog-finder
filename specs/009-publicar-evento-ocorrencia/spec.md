# Feature Specification: Publicar evento de ocorrência nova ou atualizada (escopo Finder)

**Feature Branch**: `009-publicar-evento-ocorrencia` (trabalho permanece na branch git `Development`; nenhuma branch dedicada foi criada nesta execução, conforme instrução do solicitante)

**Created**: 2026-11-14

**Status** · Draft · Cliente Finder DN-009-2 implementado · Aceite integrado depende da spec 008 e da API do Módulo

**Input**: Jira **LLT-201**; Epic Finder — IA própria de classificação e agregação; Feature **Agregação em ocorrências (E07.F3)**; item **E07.F3.2**; contexto em `contexto-LLT-103.md` §2.3 e §3; clarificações João e Drive 15–16/11/2026.

**Rastreabilidade**: Depende da classificação (**LLT-199** / spec `006-modelo-classificacao-logs`) e do agrupamento **nova/atualizada/recorrente** (**LLT-200** / spec `008-agrupar-logs-ocorrencia`); complementa a proibição de publicar evento na spec 006 (**FR-009**, **SC-006**).

## Clarifications

### Session 2026-11-16

- Q: T1 — Quem grava a ocorrência e quem publica o evento de aviso? → A: Pelo Drive (T1 revisto 15/11 23h47): o **Módulo Ocorrência** cria a ocorrência; o **Finder não publica** evento de ocorrência — envia **ResumoAnalise** ao Módulo Ocorrência (**N1 A**: HTTP interno com JWT). O aviso ao **Core** é do Módulo Ocorrência (**009-3 D**). (fonte: respostas-joao.md, T1 revisto e N1 A; revisao-T1-drive.md §1 e §4; origem: ADR_062, SSID_009.1, arquitetura-modulos v2)
- Q: 009-1 — Destinatário do evento de ocorrência no fluxo de produto? → A: **B** — o **Core** recebe o aviso do Módulo Ocorrência e repassa à **Notificação** (stream de eventos de notificação). (fonte: respostas-documentadas.md, 009-1; revisao-T1-drive.md §1; origem: ADR_165, guia-redis-streams v3 §1)
- Q: 009-2 — Papel de “avisar Ocorrências” no LLT-201? → A: O Módulo Ocorrência **não recebe um “aviso” separado** — recebe o **ResumoAnalise**, cria/atualiza ocorrências e dispara o restante do fluxo; “avisar Ocorrências” no épico Finder = **enviar o ResumoAnalise**. (fonte: revisao-T1-drive.md, 009-2; origem: ADR_062, sequencia-varredura-finder.md v6 §2)
- Q: 009-3 — Stream, type, version e PAPEL do trecho até o Core? → A: **D** — **chamada HTTP interna com JWT** (Módulo Ocorrência → Core); fecha o [PENDENTE] da ADR_051 para esse trecho. Redis Streams com envelope ADR_146 aplicam-se ao trecho **Core → Notificação** (`notificacao-eventos`). (fonte: respostas-joao.md, 009-3 D; revisao-T1-drive.md §1; origem: ADR_051, SSID_009.1)
- Q: 009-4 — Formato do corpo do aviso HTTP (Módulo Ocorrência → Core)? → A: **B** — **envelope da ADR_146** transportado por HTTP (campos do critério LLT-201 no `payload`). (fonte: respostas-joao.md, 009-4 B; origem: ADR_146 v1, LLT-201 critério 1)
- Q: 009-5 — Conteúdo do campo “ocorrência” no evento? → A: **A** — **somente o id** da ocorrência (id próprio do Módulo Ocorrência). (fonte: respostas-joao.md, 009-5 A; origem: ADR_040)
- Q: 009-6 — Usar id da ocorrência sem fechar a chave concreta (C8)? → A: **A** — **sim**; id próprio é distinto da chave de idempotência da ocorrência ([PENDENTE] na ADR_040). (fonte: respostas-documentadas.md, 009-6; origem: ADR_040, ADR_035)
- Q: 009-7 — Entidade referenciada pelo campo “ocorrência” no evento? → A: Ocorrência criada pelo **Módulo Ocorrência** (substitui leitura T1 A); o evento referencia o **id** dessa entidade (**009-5 A**). (fonte: respostas-joao.md, T1 revisto + 009-5 A; revisao-T1-drive.md A.1; origem: ADR_040)
- Q: 009-8 — O ResumoAnalise (ADR_062) continua? → A: **Sim** — permanece o contrato Finder → Módulo Ocorrência (contadores + lista de AchadoAnalise). (fonte: revisao-T1-drive.md, 009-8; origem: ADR_062, ADR_102, Finder.mmd)
- Q: 009-9 — Chave de idempotência que impede o Core de publicar o mesmo aviso duas vezes (HTTP)? → A: **A** — **id da ocorrência + estado** (natureza da mudança: nova/atualizada/recorrente). No trecho **Core → Notificação** (Redis), vale **eventId** (ADR_146 / guia §8). (fonte: respostas-joao.md, 009-9 A; origem: ADR_146, guia-redis-streams v3 §8)
- Q: 009-10 — Idempotência no produtor, consumidor ou ambos? → A: **B** — **no consumidor** (padrão guia §8 / ADR_146); o **009-11 B** fixa o efeito duplicado a evitar no Core. (fonte: respostas-documentadas.md, 009-10; respostas-joao.md, 009-11 B; origem: guia-redis-streams v3 §8)
- Q: 009-11 — “Efeito duplicado” do SC-003? → A: **B** — um **2º alerta publicado pelo Core** em `notificacao-eventos`. (fonte: respostas-joao.md, 009-11 B)
- Q: 009-12 — Gravidade no aviso? → A: **Sim, no aviso** emitido pelo **Módulo Ocorrência** (SSID_009.1 passo 1); o Finder **não** calcula gravidade (ADR_062). (fonte: revisao-T1-drive.md, 009-12 revista; origem: SSID_009.1, ADR_062, ADR_165)
- Q: 009-13 — Significado do campo “tipo”? → A: **A** — classe **erro**, **alerta** ou **informação** (ERROR/WARNING/INFO); a natureza da mudança é campo separado (LLT-201 critério 1). (fonte: respostas-documentadas.md, 009-13; origem: ADR_170, LLT-201)
- Q: 009-14 — Reabertura (008-4 B) no evento? → A: **A** — só **nova/atualizada/recorrente**; reabertura mapeia como **atualizada**. (fonte: respostas-joao.md, 009-14 A; origem: LLT-201, ADR_082, 008-4 B)
- Q: 009-15 — Recorrente gera nova notificação? → A: **A** — **nova notificação** (fecha FA01 do ponto de vista do fluxo de aviso; política fina permanece na Notificação / LLT-216). (fonte: respostas-joao.md, 009-15 A; origem: SSID_009.1 FA01)
- Q: 009-16 — Momento do aviso no pipeline? → A: O **Módulo Ocorrência** avisa ao **criar** ou **mudar** a ocorrência; com **Sugestão automática**, o aviso ao Core **após** a Sugestão (SSID_009.1 passo 2). No Finder: envio do **ResumoAnalise** ao fim da análise do arquivo/lote (seq. v6 §2), **sem** publicar evento LLT-201. (fonte: revisao-T1-drive.md, 009-16; origem: SSID_009.1, LLT-201, sequencia-varredura-finder.md v6)
- Q: 009-17 — Um evento por mudança ou agregado? → A: **Um aviso por ocorrência** (“envia ao Core **cada** Ocorrência”, SSID_009.1 passo 1). (fonte: revisao-T1-drive.md, 009-17; origem: SSID_009.1)
- Q: 009-18 — Falha ao publicar / entregar? → A: **HTTP Módulo → Core (009-18 A):** tentar de novo e alertar. **Finder → Módulo (ResumoAnalise):** falha entre enviar e confirmar → reenvio com mesma chave (**FA-02**); falha **consulta OpenObserve** no ingested → **FR-F07** / **DN-008-3** (3 tentativas; **sem DLQ**; Finder→Core API — etapa 3h). *(Trecho «análise → DLQ FA-01» do seq. v6 aplica-se à **FilaOrdenada**, não ao **log.ingested v1** — **DLQ B**.)* (fonte: respostas-joao.md, 009-18 A; revisao-T1-drive.md, 009-18; origem: ADR_014, ADR_046, sequencia-varredura-finder.md v6)
- Q: 009-19 — Mesmo canal dos alertas ADR_051? → A: **Sim** — Módulo Ocorrência → Core por **HTTP interno com JWT** (mesmo padrão de alerta da ADR_051; emissor revisado para Módulo Ocorrência). (fonte: revisao-T1-drive.md, 009-19; origem: ADR_051, arquitetura-modulos v2)
- Q: 009-20 — Transição com a spec 006 (FR-009 proíbe evento)? → A: **T6 A** — lacunas da 006 são fechadas **atualizando a 006**; o Finder passa a produzir **ResumoAnalise/AchadoAnalise** em vez de publicar evento; esta feature define o handoff e referencia o pipeline completo. (fonte: respostas-joao.md, T6 A + T1 revisto; origem: spec 006 FR-009)
- Q: 009-21 — Evitar aviso ausente ou duplicado? → A: Idempotência do **ResumoAnalise** (**N4 A**: chave **`batchId` + `execId`** ou só **`batchId`** — **DN-006-1b B**); confirmação do Módulo antes de `finishScan()`; duplicata ignorada no Módulo (FA-02). Idempotência do aviso ao Core (**009-9 A**). (fonte: respostas-joao.md, N4 A; revisao-T1-drive.md §2; origem: ADR_014, ADR_065/070) — *chave sem execId refinada na Session etapa 3b (DN-006-1/1b).*
- Q: 009-22 — A 009 fecha aceite antes da 008? → A: **A** — **não**; a 009 **espera** a 008 para **atualizada** e **recorrente**. (fonte: respostas-joao.md, 009-22 A)
- Q: T2 — Grão da ocorrência? → A: **C** — **um achado elegível por classe ERROR, WARNING ou INFO** no ResumoAnalise (**T2 C** + **T7 A** / **006-11 A**). (fonte: respostas-joao.md, T2 C e T7 A; origem: 006-11 A)
- Q: T3 — Recorrência? → A: **B** — consulta ao OpenObserve (ADR_168). (fonte: respostas-joao.md, T3 B)
- Q: T4 — Lote parcial? → A: **A** — lote incompleto não conta como sucesso. (fonte: respostas-joao.md, T4 A)
- Q: T5 — Versão do log.ingested? → A: **D** — contrato **08/11 v1** é a única versão válida (`schemaVersion` = "1"). (fonte: respostas-joao.md, T5 D; origem: contrato-log-ingested-v1)
- Q: T6 — Lacunas da 006? → A: **A** — atualizar a própria spec 006. (fonte: respostas-joao.md, T6 A)
- Q: T7 — INFO gera ocorrência? → A: **A** — **sim**, por achado (**006-11 A**); decisão consciente divergindo do Drive (SSID_007.1 / SSID_009.1 / ADR_040). (fonte: respostas-joao.md, T7 A e rodada 16/11 00h01; origem: ADR_170, SSID_007.1)
- Q: N1 — Transporte Finder → Módulo Ocorrência? → A: **A** — HTTP interno com JWT. (fonte: respostas-joao.md, N1 A; origem: paridade com 009-3 D)
- Q: N2 — Agrupamento entre lotes? → A: **B** — o **Finder** junta logs entre lotes consultando ocorrências antes de enviar. (fonte: respostas-joao.md, N2 B)
- Q: N3 — Ampliar AchadoAnalise? → A: **A** — sim (assinatura, contagem, primeira/última vez, versão das regras), com **ADR nova**. (fonte: respostas-joao.md, N3 A)
- Q: N4 — Chave de idempotência do ResumoAnalise? → A: **A** — **batchId + execId**. (fonte: respostas-joao.md, N4 A; origem: contrato 08/11, ADR_070)

### Session 2026-11-16 (etapa 3b)

- Q: DN-009-1 De onde vem o resumo mascarado do e-mail (LLT-218)? → A: **A** — a **Notificação** lê a **Ocorrência** pelo **id** no **MongoDB compartilhado** (**ADR_054**). (fonte: respostas-joao.md §Etapa 3 / respostas-etapa3.md §7; origem: ADR_054, sequencia-notificacao.md v2)
- Q: 009-5 (revisto) O que o aviso do Módulo Ocorrência ao Core deve levar além do mínimo do critério Jira? → A: O aviso HTTP ao Core MUST incluir **id da ocorrência**, **tipo** (classe) e **gravidade** (substitui a leitura “só o id” como descrição do aviso inteiro); o campo **`ocorrencia`** no envelope permanece o **id** opaco (**009-6 A**). (fonte: respostas-joao.md §Etapa 3; origem: sequencia-notificacao.md v2 diagrama 1, ADR_040)
- Q: DN-006-1 Origem e obrigatoriedade do `execId` na chave do ResumoAnalise? → A: O **execId** vem da **Varredura** (`V<uuid>` automática, `M<uuid>` manual — **ADR_075**, **ADR_069**); **não é obrigatório** (“quando houver”, **ADR_146**; **N4 A** com “se houver”). (fonte: respostas-joao.md §Etapa 3 / respostas-etapa3.md §3; origem: ADR_075, ADR_069, ADR_070, ADR_146)
- Q: DN-006-1b Como o Finder obtém `execId` enquanto consome só `log.ingested` v1? → A: **B** — **sem `execId`** nesse fluxo: a chave de idempotência do **ResumoAnalise** fica **somente `batchId`** até existir contrato de `log.checked` ou outra fonte. (fonte: respostas-joao.md §Etapa 3 / respostas-etapa3.md §3; origem: contrato log.ingested v1, ADR_070)
- Q: API do Módulo Ocorrência (rota, JWT, respostas HTTP, 5xx/timeout) para o handoff Finder → Módulo? → A: **PENDÊNCIA** — fechar no **contrato/ADR do Módulo Ocorrência** antes do plan implementável; o Drive só fixa API interna + JWT próprio + rede privada (**ADR_028**, **ADR_054**). (fonte: respostas-joao.md §Etapa 3 / respostas-etapa3.md §10; origem: ADR_028, ADR_054, ADR_141)

### Session 2026-11-16 (etapa 3c)

- Q: DN-007-2 Como a meta **007-meta** (≥95% de precisão) se aplica a esta feature? → A: **B** — **≥95% de precisão por classe** (ERROR, WARNING, INFO) na base de validação da spec **007**; **esta spec (009) não define gate métrico** — a rastreabilidade a **LLT-199** / spec **006** permanece; o aceite de precisão fica em **SC-007** / **006-2 B** da spec **006**, não nos critérios de sucesso do handoff LLT-201. (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: respostas-etapa3.md §1, LLT-199, 007-meta)
- Q: DN-008-4 Qual algoritmo normaliza a mensagem na composição da **assinatura** transportada no **AchadoAnalise**? → A: **B** — aparar espaços nas extremidades; ignorar maiúsculas/minúsculas; substituir **números**, **UUIDs**, **IPs**, **datas** e **hexadecimais** por **marcadores** (um por categoria). Algoritmo **canônico** na spec **008**; esta feature **MUST** reutilizá-lo ao preencher assinatura ou equivalência (**N3 A**, **FR-F09**, dependência **006**/**008**), sem regras divergentes; literais exatos dos marcadores permanecem na spec **008** (**DN-008-5**). (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: DN-008-1 B, ADR_013)

### Session 2026-11-16 (etapa 3c — correções analyze)

- Q: C1 Como conciliar o princípio **I** da constituição (“Finder agrega em ocorrência”) com o **T1 revisto**? → A: Manter escopo desta feature em handoff **ResumoAnalise**; criação de ocorrência de negócio no **Módulo Ocorrência**; **emenda do princípio I** permanece pendência de governança **fora** desta pasta (não editar constituição aqui). (fonte: respostas-joao.md, T1 REVISTO 15/11 23h47; origem: revisao-T1-drive.md §1, ADR_062)
- Q: C2 O que significa “agregação em ocorrências” no **Finder** face ao princípio **II** da constituição? → A: No Finder, agregação = classificar logs, montar **AchadoAnalise**, enviar **ResumoAnalise** e, quando aplicável, **N2 B** (consulta prévia entre lotes); **não** inclui persistir ocorrência de negócio nem publicar aviso — isso fica no **Módulo Ocorrência** (**T1 revisto**, **ADR_062**). (fonte: respostas-joao.md, T1 REVISTO; origem: ADR_062, revisao-T1-drive.md §1)
- Q: E2 O que SC-004 deve medir sem misturar handoff Finder e aviso HTTP do Módulo? → A: Reformular **SC-004** como critério de **tipo** (classe) nos avisos **Módulo → Core**; taxonomia no handoff Finder fica em **SC-F003** / **FR-F08** (**009-13 A** no produto, não “avisos com origem Finder”). (fonte: respostas-documentadas.md, 009-13 A; origem: ADR_170, LLT-201 critério 1)
- Q: F2 Qual o grão da ocorrência quando **T2 C** fala em ERROR/WARNING e **T7 A** exige **INFO**? → A: **Um achado elegível por classe ERROR, WARNING ou INFO** no **ResumoAnalise** (**T2 C** + **T7 A** / **006-11 A**); INFO não fica fora do grão acordado. (fonte: respostas-joao.md §Etapa 3 / rodada 16/11 00h01, T2 C e T7 A; origem: 006-11 A)
- Q: F3 Como nomear o contador que substitui “spam ignorado” no **ResumoAnalise**? → A: **spam: B** — trocar o contador legado por contador de **INFO**; nome canônico na lista ADR_102 aguarda fechamento dessa ADR (até lá, documentar como **contador de INFO**). (fonte: respostas-joao.md §Etapa 3, spam: B; origem: respostas-etapa3.md §9, ADR_102, ADR_170)
- Q: D1 Qual narrativa canônica do conteúdo do aviso após **009-5 revisto**? → A: Requisitos e **FR-003** usam só **009-5 revisto** (id + tipo + gravidade no HTTP; campo `ocorrencia` no envelope Redis = id); leitura histórica **009-5 A** permanece apenas em **Divergências D12** e nota na tabela de decisões. (fonte: respostas-joao.md §Etapa 3, 009-5 revisto; origem: sequencia-notificacao.md v2)

### Session 2026-11-16 (etapa 3e)

- Q: DN-008-7 O handoff Finder → Módulo Ocorrência é um envio por lote ou um envio por achado? → A: **Um** **ResumoAnalise** por **lote** (arquivo/lote analisado), com contadores + lista de **AchadoAnalise** e chave de idempotência do resumo — não uma chamada HTTP por achado. (fonte: respostas-etapa3c.md §6; origem: ADR_062, Finder.mmd DC-10 §7, N4 A / `batchId`)
- Q: DN-008-7 B Como o Finder consulta ocorrências existentes antes do envio (**N2 B**)? → A: **B** — **uma consulta ao Módulo Ocorrência por achado** (por assinatura/classe elegível a consultar), não uma única consulta agregada por lote; **rotas, claims JWT e contrato HTTP** da consulta permanecem em **DN-009-2** / ADR do Módulo. (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); origem: N2 B, LLT-200)
- Q: DN-008-6 Após **3 tentativas** e **DLQ** em falha transitória do processamento do lote, o usuário recebe alerta? → A: *(registro §Etapa 3e — **SUPERADO** pela Session **etapa 3h**: falha **OpenObserve** esgotada → Finder avisa **Core** por API interna; **sem DLQ** no ingested; ver **FR-F07** e Decisões **DN-008-6 (rev.)**.)*
- Q: DN-008-5 Quais são os literais dos marcadores do normalizador (**DN-008-4 B**)? → A: **A** — `<NUM>`, `<UUID>`, `<IP>`, `<DATE>` e `<HEX>` (um marcador por categoria). Dono canônico e ordem de aplicação: spec **008**; esta feature **MUST** reutilizar os mesmos literais ao compor **mensagem normalizada** / **assinatura** (**FR-F11**). (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h01 e 11h; origem: DN-008-4 B, ADR_013)
- Q: 008 assinatura A Como serializar o valor da assinatura nos achados do resumo? → A: **A** — **texto** formado pela **concatenação** das partes acordadas (**DN-008-1 B**) com **separador fixo**; o literal exato do separador é decisão da spec **008** (não fixado nesta feature). (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h; origem: DN-008-1 B, 008-11)

### Session 2026-11-16 (etapa 3f)

- Ajuste: removida a divergência **D16** (Etapa 3c 11h01 «puladas» × 11h); as rodadas **não** se contradizem — **DN-006-3**, **DN-007-4** e **DN-008-7** foram respondidos em seguida (**SUPERADO** no arquivo de respostas). (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h01) — linha SUPERADO; respostas-joao.md §Etapa 3c (João, 16/11 11h02))
- Ajuste: fonte canônica **DN-008-7 B** (consulta **por achado**) passa a **respostas-joao.md §Etapa 3c (João, 16/11 11h02)** nas citações desta spec. (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); origem: N2 B, LLT-200)
- Ajuste: fonte canônica **DN-006-3 B** e **DN-007-4 A** → **respostas-joao.md §Etapa 3c (João, 16/11 11h02)** (mesma rodada que **DN-008-7 B**; sem tratar 11h01×11h como divergência). (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); origem: DN-006-3 B, DN-007-4 A)

### Session 2026-11-16 (etapa 3g)

- Ajuste: **DN-008-8 A** — na serialização **008 assinatura A**, o **separador fixo** entre as partes (**DN-008-1 B**) é o caractere **`|`** (pipe); dono canônico da decisão: spec **008**; esta feature **reutiliza** o literal (**FR-F11**). (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40); origem: DN-008-8 A, 008 assinatura A)

### Session 2026-11-16 (etapa 3h)

- Ajuste: **DLQ B** — **não** existe **DLQ** no fluxo **log.ingested v1**; as políticas retry/DLQ do guia (**FilaOrdenada**) **não** se aplicam ao ingested (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B; origem: respostas-etapa3g.md §1b, contrato log.ingested v1, DN-006-1b B).
- Ajuste: **DN-008-3** / **DN-008-6** (alcance revisto) — falha transitória na **consulta ao OpenObserve**: até **3 tentativas**; após esgotar, o Finder avisa o **Core** pela **API interna** com **`batchId`** e **motivo**; **sem** DLQ nem alerta via consumidor DLQ / **WebSocket** pós-DLQ (substitui leitura §Etapa 3c 11h). Rota, payload e evento da API do Core: **[PENDENTE]** (**ADR_051 [PENDENTE]**). (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B; origem: respostas-etapa3g.md §1, DN-008-3/008-6)
- Ajuste: **DN-008-9 A** — um caractere **`|`** literal **dentro** de um componente da assinatura é serializado como **`\|`** (barra invertida + pipe); dono canônico spec **008**; esta feature **reutiliza** (**FR-F11**). (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DN-008-9 A; origem: DN-008-8 A, 008 assinatura A)

### Session 2026-11-20

- O usuário aprovou a implementação de DN-009-2 e N3 A. O contrato define consulta por achado, entrega idempotente por lote, recibo explícito e revisão opaca para concorrência. O adaptador Finder foi verificado com servidor de teste. A homologação do Módulo e o aceite SC-002 permanecem em aberto. [Contrato aprovado](contracts/proposta-dn-009-2.md)

## Responsabilidades por componente (delimitação de escopo)

| Componente | Repositório / equipe | Responsabilidade neste fluxo (LLT-201 e dependentes) |
|------------|----------------------|------------------------------------------------------|
| **Finder** | **Este repositório** | Classificar e **agregar** logs em **AchadoAnalise** + **ResumoAnalise** (princípio **II** nesta feature: **não** persistir ocorrência de negócio — **T1 revisto**, **ADR_062**); consulta entre lotes (**N2 B**) **por achado** ao Módulo Ocorrência (**DN-008-7 B**; contrato em **DN-009-2**) quando a spec **008** estiver disponível; ao fim do lote, enviar **exatamente um** **ResumoAnalise** por lote (**DN-008-7**, **HTTP + JWT**, **N1 A**); chave **`batchId`**, ou **`batchId` + `execId`** quando houver (**N4 A**, **DN-006-1**, **DN-006-1b B**); contadores do resumo incluem **INFO** em lugar do legado “spam ignorado” (**spam: B**); incluir achados **INFO** no resumo (**T7 A**); ampliar AchadoAnalise conforme **N3 A** (depende de ADR); **não** publicar evento de ocorrência; **não** calcular gravidade; **não** chamar o Core para aviso de ocorrência LLT-201; após **3 falhas** na consulta **OpenObserve**, avisa o **Core** pela **API interna** (**`batchId`** + **motivo** — **DN-008-3**/**DN-008-6**; contrato **[PENDENTE ADR_051]**), distinto do aviso Módulo→Core. |
| **Módulo Ocorrência** | Fora deste repositório | Receber ResumoAnalise; criar/atualizar ocorrências; definir **gravidade**; converter AchadoAnalise → Ocorrência; disparar Sugestão automática quando aplicável; enviar **cada** ocorrência ao Core (**009-3 D**, **009-17**) com **id**, **tipo** (classe) e **gravidade** (**009-5 revisto**, **SSID_009.1**). |
| **Core** | Fora deste repositório | Receber aviso HTTP; deduplicar (**009-9 A**); publicar em **Redis Streams** (`notificacao-eventos`, envelope **ADR_146**, idempotência **eventId** no consumidor — **009-10 B**). |
| **Notificação** | Fora deste repositório | Consumir stream; **ler Ocorrência pelo id** no MongoDB compartilhado para montar e-mail com resumo mascarado (**DN-009-1 A**, **LLT-218**); destinatários e canais (tipo/gravidade no aviso — **009-5 revisto**); política de notificações repetidas (**LLT-216** / **008-17 C**); efeito de duplicata alvo do **009-11 B**. |

Fluxo acordado (Drive + decisões João): **Finder** → ResumoAnalise → **Módulo Ocorrência** → HTTP JWT → **Core** → stream → **Notificação**.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Handoff Finder → Módulo Ocorrência (Priority: P1)

Quando a análise de um arquivo/lote produz achados (incluindo **informação**, **T7 A**), o Finder envia **ResumoAnalise** ao Módulo Ocorrência de forma idempotente, para que o módulo crie ocorrências e acione o pipeline de aviso downstream.

**Why this priority**: É a única contribuição do Finder ao LLT-201 após o T1 revisto; sem handoff confiável não há ocorrência nem aviso.

**Independent Test**: Simular fim de análise com achados e verificar **um** envio HTTP com JWT por lote (**DN-008-7**), payload **ResumoAnalise**, chave **`batchId`** (fluxo `log.ingested` v1, **DN-006-1b B**) ou **`batchId` + `execId`** quando disponível, confirmação do Módulo e ausência de `finishScan()` antes da confirmação; quando **N2 B** aplicável, verificar consultas **por achado** (**DN-008-7 B**) antes do envio.

**Acceptance Scenarios**:

1. **Given** achados classificados (erro, alerta ou informação), **When** o arquivo/lote é analisado, **Then** o Finder envia **ResumoAnalise** ao Módulo Ocorrência (**FR-F01**) e **não** publica evento LLT-201 (**FR-F02**).
2. **Given** o Módulo confirma recebimento, **When** o Finder conclui a varredura, **Then** `finishScan()` só ocorre após confirmação (seq. v6 §2; **FR-F05**).
3. **Given** falha após envio sem confirmação, **When** o lote é reprocessado, **Then** o Finder reenvia com a mesma chave e o Módulo ignora duplicata (**FR-F06**).

---

### User Story 2 - Pipeline de aviso (referência; validação fora do Finder) (Priority: P1)

Quando o Módulo Ocorrência cria ou muda uma ocorrência, o produto emite aviso HTTP com tenant, produto, **id**, **tipo** (classe), **gravidade**, natureza da mudança (**nova/atualizada/recorrente**), de forma idempotente até a Notificação (**009-5 revisto**).

**Why this priority**: Cobre o critério (1) e (2) do LLT-201 no produto; o Finder depende deste fluxo para valor de negócio, mas **não implementa** o publicador.

**Independent Test**: Testes de aceite LLT-201 no ambiente integrado (Módulo Ocorrência + Core + Notificação), usando ResumoAnalise produzido pelo Finder.

**Acceptance Scenarios**:

1. **Given** ocorrência **nova** criada pelo Módulo, **When** o aviso HTTP é aceito pelo Core, **Then** o payload lógico contém **tenant**, **produto**, **id da ocorrência**, **tipo** (classe), **gravidade** e natureza **nova** (**FR-003**, **009-5 revisto**; J: LLT-201 critério 1 + SSID_009.1).
2. **Given** ocorrência **atualizada** ou **recorrente** (spec 008), **When** o Módulo envia aviso, **Then** o estado no payload é **atualizada** ou **recorrente**; reabertura mapeia como **atualizada** (**009-14 A**).
3. **Given** ocorrência **recorrente**, **When** a Notificação processa o evento, **Then** o fluxo prevê **nova notificação** (**009-15 A**), sem contradição com o aviso emitido pelo Módulo.

---

### User Story 3 - Idempotência ponta a ponta (Priority: P1)

Reentregas do mesmo fato não produzem efeito duplicado indevido no Core (**009-11 B**) nem duplicam handoff Finder → Módulo.

**Independent Test**: Reenviar ResumoAnalise com a mesma chave (**`batchId`** ou **`batchId` + `execId`**); reenviar aviso HTTP com mesmo **id ocorrência + estado**; redeliverar mensagem na stream e verificar consumidor.

**Acceptance Scenarios**:

1. **Given** mesmo ResumoAnalise já confirmado, **When** o Finder reenvia, **Then** o Módulo não duplica ocorrências (**FR-F06**, **SC-F002**).
2. **Given** mesmo aviso já processado pelo Core, **When** o HTTP se repete, **Then** o Core **não** publica 2º alerta em `notificacao-eventos` (**SC-003**, **009-11 B**).
3. **Given** mesma mensagem na stream, **When** o consumidor Notificação reprocessa, **Then** aplica idempotência por **eventId** (**009-10 B**, guia §8).

---

### Edge Cases

- **Épico Jira no Finder × fluxo Drive**: LLT-201 no Finder significa handoff ResumoAnalise, não publicação de evento (**divergência documentada**).
- **INFO vira ocorrência (João) × SSID_007.1/009.1 só ERRO/WARNING**: ResumoAnalise pode incluir achados INFO; o Módulo deve aceitar (**T7 A**); divergência registrada.
- **Aviso Módulo → Core**: **id** + **tipo** + **gravidade** no HTTP (**009-5 revisto**); campo **`ocorrencia`** no envelope Redis permanece o **id** opaco (**009-6 A**). Chave de unicidade da ocorrência (C8) continua **[PENDENTE]** na ADR_040.
- **E-mail LLT-218**: Notificação obtém trechos mascarados lendo a Ocorrência por **id** no MongoDB (**DN-009-1 A**); templates de corpo no Drive ainda **[PENDENTE]** (ressalva respostas-etapa3.md §7).
- **Handoff sem `execId`**: enquanto o Finder só consome **log.ingested** v1, idempotência do ResumoAnalise usa **só `batchId`** (**DN-006-1b B**).
- **API Finder → Módulo**: rota, claims JWT, códigos e timeout **[PENDENTE contrato]** (**DN-009-2**).
- **Gravidade**: presente no aviso do Módulo, ausente no Finder (**ADR_062**).
- **Sugestão automática**: aviso ao Core após Sugestão quando aplicável (SSID_009.1 passo 2).
- **Dependência spec 008**: estados **atualizada/recorrente** e agrupamento entre lotes (**N2 B**, **008-11** revisado com classe) bloqueiam SC-002 até a 008.
- **Assinatura e normalização**: o campo **assinatura** em **AchadoAnalise** (**N3 A**) usa **mensagem normalizada** conforme **DN-008-4 B**, marcadores **DN-008-5 A** (`<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` — canônico spec **008**), serialização **008 assinatura A** (concatenação **DN-008-1 B** com separador fixo **`|`** — **DN-008-8 A**; **`|`** nos componentes → **`\|`** — **DN-008-9 A**; escape de **`\`** **[PENDENTE]** sem fonte; dono canônico spec **008**); afeta agrupamento **FR-F09**/**N2 B**, não o formato do aviso HTTP.
- **Handoff vs consulta ao Módulo**: **um ResumoAnalise por lote** (**DN-008-7**); consulta prévia **por achado** (**DN-008-7 B**); rotas JWT da consulta e do POST do resumo: **DN-009-2**.
- **CorrelationId aleatório por exportação** (contrato 08/11 §4, **008-2 A**): risco para assinatura se correlationId integrar assinatura — registrar na 008; ResumoAnalise usa **`batchId`** (e **`execId`** quando houver) na chave (**N4 A**, **DN-006-1b B**).
- **Envelope ADR_146 v2 “schemaVersion 2”** × contrato ingested **v1** (**T5 D**): não misturar políticas; eventos de notificação seguem ADR_146 v1 vigente fora do ingested.
- **Falha consulta OpenObserve** (**T3 B**, **DN-008-3**): até **3 tentativas**; **sem DLQ** no **log.ingested v1** (**DLQ B**); após esgotar, Finder → **Core** por **API interna** (**`batchId`**, **motivo** — **DN-008-6**); contrato da API **[PENDENTE ADR_051]**; **sem** handoff **ResumoAnalise** de sucesso. Falha handoff ResumoAnalise: não marcar varrido (**FA-02**).

## Requirements *(mandatory)*

### Functional Requirements — Finder (escopo deste repositório)

- **FR-F01**: Ao concluir a análise de um arquivo/lote com achados elegíveis (**um achado por classe ERROR, WARNING ou INFO** — **T2 C**, **T7 A** / **006-11 A**), o Finder MUST enviar **exatamente um** **ResumoAnalise** por lote ao Módulo Ocorrência via **HTTP interno com JWT** (**N1 A**, **DN-008-7**), contendo contadores (inclui **contador de INFO** no lugar do legado “spam ignorado”, **spam: B**, até fechar lista na ADR_102) e **AchadoAnalise** conforme **ADR_062** / **ADR_102**, ampliado quando **N3 A** estiver formalizado em ADR. MUST NOT enviar um HTTP separado por achado no handoff do resumo.
- **FR-F02**: O Finder MUST **NOT** publicar evento de ocorrência nova/atualizada/recorrente ao Core nem à Notificação (**T1 revisto**, **ADR_062**).
- **FR-F03**: A chave de idempotência do ResumoAnalise MUST ser **`batchId`** quando o Finder não dispuser de **`execId`** (**DN-006-1b B**, leitura **log.ingested** v1, **T5 D**); MUST ser **`batchId` + `execId`** quando o identificador de varredura estiver disponível (**N4 A**, **DN-006-1**). Sem inventar origem de `execId` (**constituição V**).
- **FR-F04**: Achados **informação** MUST constar no ResumoAnalise quando classificados (**T7 A**, **006-11 A**), para o Módulo Ocorrência criar ocorrência.
- **FR-F05**: O Finder MUST aguardar **confirmação de recebimento** do Módulo Ocorrência antes de `finishScan()` (sequencia-varredura-finder.md v6 §2).
- **FR-F06**: Reenvio do ResumoAnalise após falha MUST usar a mesma chave; o Módulo MUST tratar duplicata sem efeito colateral indevido (FA-02, ADR_014).
- **FR-F07**: Falha transitória na **consulta ao OpenObserve** durante a análise do lote (recorrência **T3 B**; detalhes na spec **008**) MUST permitir até **3 tentativas**. No fluxo **log.ingested v1** (**DN-006-1b B**, **T5 D**) **MUST NOT** usar **DLQ** nem replicar retry/DLQ da **FilaOrdenada** (**DLQ B**, contrato 08/11). Após esgotar tentativas, o Finder MUST avisar o **Core** pela **API interna** com **`batchId`** e **motivo** (**DN-008-3**, **DN-008-6** — substitui alcance DLQ + WebSocket pós-DLQ). Rota, payload e evento dessa API: **[PENDENTE]** (**ADR_051 [PENDENTE]**). Bloqueia handoff **ResumoAnalise** até análise concluída com sucesso ou política acordada na spec **008**/**005** sem inventar contrato HTTP.
- **FR-F08**: O campo **tipo** em cada AchadoAnalise MUST ser **erro**, **alerta** ou **informação** (constituição III, **009-13 A**).
- **FR-F09**: O Finder MUST implementar agrupamento **entre lotes** consultando ocorrências existentes no Módulo Ocorrência **antes** de montar o **ResumoAnalise** (**N2 B**), com **uma consulta por achado** (**DN-008-7 B**), coordenado com a spec **008**; formato de rota/claims da consulta: **DN-009-2** (não fixado nesta spec).
- **FR-F10**: Esta feature MUST NOT alterar **ADR_019**, **ADR_167** nem **SSID_007.1** (apenas citação); taxonomia vigente é **ERROR/WARNING/INFO** (**ADR_170**).
- **FR-F11**: Ao compor **assinatura** (ou campos derivados) em **AchadoAnalise** (**N3 A**, **FR-F09**), o Finder MUST aplicar o normalizador de mensagem **DN-008-4 B** definido na spec **008** (trim, case-insensitive, marcadores por categoria para número, UUID, IP, data, hexadecimal) com literais **DN-008-5 A**: `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>`. A **assinatura** serializada MUST seguir **008 assinatura A** (texto concatenado das partes **DN-008-1 B** com separador fixo **`|`** — **DN-008-8 A**; dono canônico spec **008**). Um **`|`** literal **dentro** de um componente MUST ser escapado como **`\|`** (**DN-008-9 A**). MUST NOT publicar semântica divergente nesta feature. (fonte: respostas-joao.md §Etapa 3b/3c/3e/3g, DN-008-4 B / DN-008-5 A / 008 assinatura A / **DN-008-8 A** / **DN-008-9 A**; origem: DN-008-1 B, ADR_013)

### Functional Requirements — Produto (referência; implementação fora do Finder)

- **FR-001**: O **publicador** do aviso de ocorrência criação/mudança MUST ser o **Módulo Ocorrência** (**T1 revisto**, **SSID_009.1**), **não** o Finder.
- **FR-002**: A rota do aviso MUST ser **HTTP interno com JWT** (Módulo → Core, **009-3 D**); a **Notificação** é alcançada via **Core** e stream Redis (**009-1 B**, **ADR_165**).
- **FR-003**: O aviso HTTP Módulo → Core MUST incluir, no mínimo, **tenant**, **produto**, **id da ocorrência**, **tipo** (classe), **gravidade** e natureza **nova**, **atualizada** ou **recorrente** (**009-5 revisto**, **009-12** revista, **SSID_009.1**; alinhado a sequencia-notificacao.md v2 diagrama 1).
- **FR-004**: O campo **`ocorrencia`** no payload publicado pelo Core MUST ser o **id** opaco da ocorrência (**ADR_040**, **009-6 A**); o aviso HTTP upstream já transporta **tipo** e **gravidade** (**009-5 revisto**). Chave concreta de unicidade (C8) permanece **[PENDENTE]** na ADR_040.
- **FR-005**: A entrega MUST ser **idempotente** (J: LLT-201 critério 2): chave HTTP **id ocorrência + estado** (**009-9 A**); na stream, **eventId** e deduplicação no consumidor (**009-10 B**).
- **FR-006**: O Módulo MUST emitir aviso em **criação** e **mudança**, distinguindo **nova**, **atualizada** e **recorrente**; depende da spec **008** para significado operacional de atualizada/recorrente (**009-22 A**).
- **FR-007**: Falha no HTTP Módulo → Core MUST seguir **tentar de novo e alertar** (**009-18 A**).
- **FR-008**: Corpo HTTP SHOULD seguir **envelope ADR_146** com campos LLT-201 no `payload` (**009-4 B**).
- **FR-009**: Core → Notificação MUST usar **Redis Streams**, envelope **ADR_146**, papel `notificacao-eventos` (**ADR_143**, **ADR_165**).
- **FR-010**: **Recorrente** MUST acionar fluxo de **nova notificação** no conjunto acordado (**009-15 A**).
- **FR-011**: Reabertura (008-4 B) MUST aparecer no aviso como **atualizada** (**009-14 A**).
- **FR-012**: **Um aviso por ocorrência** por evento de criação/mudança tratado (**009-17**, SSID_009.1 passo 1).

### Fora de escopo (com fonte)

- Implementar **Módulo Ocorrência**, **Core** ou **Notificação** (fora deste repositório).
- Definir **contrato HTTP** Finder → Módulo enquanto **DN-009-2** estiver aberta (pendência João Etapa 3).
- **Publicar** na stream `notificacao-eventos` (Core) ou consumir Notificação.
- Definir **protocolo** genérico Finder/Coletor → Core além do handoff Finder → Módulo (**N1 A**).
- Resolver **chave concreta** da ocorrência além do id opaco (**ADR_040 [PENDENTE]**).
- Implementar agrupamento completo **008** (dependência explícita).
- Alterar **ADRs**, constitution, README, specs `001`–`008`, `.specify/`, código Go nesta execução.

### Key Entities *(include if feature involves data)*

- **ResumoAnalise**: Handoff Finder → Módulo; contadores + lista de **AchadoAnalise** (**ADR_062**, **ADR_102**, **009-8**); o contador legado “spam ignorado” é substituído por **contador de INFO** (**spam: B**); transporte **HTTP + JWT** (**N1 A**).
- **AchadoAnalise**: Achado por assinatura/classe no lote; campos Drive + extensões **N3 A** (assinatura, contagem, primeira/última vez, versão regras). A **assinatura** inclui **mensagem normalizada** (**DN-008-4 B**) e composição **DN-008-1 B**, serializada como **texto concatenado** (**008 assinatura A**; separador fixo **`|`** — **DN-008-8 A**; escape de **`|`** nos componentes: **`\|`** — **DN-008-9 A**; dono canônico spec **008**).
- **Mensagem normalizada**: Saída do algoritmo **DN-008-4 B** sobre o texto do log (após mascaramento quando aplicável — spec **010**); reutiliza implementação/spec **008**; literais dos marcadores: **DN-008-5 A** — `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` (canônico spec **008**).
- **Handoff ResumoAnalise**: **Um** payload por lote analisado (**DN-008-7**); distinto das consultas **por achado** (**DN-008-7 B**, **FR-F09**).
- **Ocorrência (negócio)**: Criada pelo Módulo Ocorrência; **id próprio** referenciado no aviso e lido pela Notificação (**009-5 revisto**, **009-6 A**, **DN-009-1 A**); trechos mascarados em `conteudo` (**ADR_040**).
- **Aviso de ocorrência (evento lógico LLT-201)**: Emitido pelo Módulo ao Core; físico HTTP (envelope **ADR_146**) + publicação Redis pelo Core.
- **Tipo (classe)**: Erro, alerta, informação (**009-13 A**, constituição III).
- **Natureza da mudança**: Nova, atualizada, recorrente; reabertura → atualizada (**009-14 A**).

## Success Criteria *(mandatory)*

### Measurable Outcomes — Finder

- **SC-F001**: Em testes acordados, **100%** dos lotes com achados produzem exatamente **um** **ResumoAnalise** por lote e por chave acordada (**`batchId`** ou **`batchId` + `execId`** — **DN-008-7**), com confirmação antes de `finishScan()`.
- **SC-F002**: Reenvio com mesma chave MUST NOT duplicar criação de ocorrência no Módulo (cenário integrado ou contrato duplicata FA-02).
- **SC-F003**: **0%** dos AchadoAnalise usam classe fora de erro, alerta ou informação (**FR-F08**).

### Measurable Outcomes — Produto (referência; aceite conjunto)

- **SC-001**: Cenário **nova**: **100%** dos avisos HTTP contêm **tenant**, **produto**, **id**, **tipo**, **gravidade** e natureza da mudança (**FR-003**, **009-5 revisto**).
- **SC-002**: Cenários **atualizada** e **recorrente** (**008**): **100%** distinguem o estado no payload (**FR-006**); **bloqueado** até spec 008 (**009-22 A**).
- **SC-003**: Reentrega duplicada MUST NOT gerar **2º alerta** no Core (**009-11 B**, **FR-005**).
- **SC-004**: **0%** dos avisos HTTP **Módulo Ocorrência → Core** usam **tipo** (classe) fora de erro, alerta ou informação (**009-13 A**, **FR-003**). Taxonomia no handoff Finder: **SC-F003** / **FR-F08** (não misturar com este critério).
- **SC-005**: Aceite formal de **SC-002** exige spec **008** concluída (**009-22 A**).
- **SC-006**: **Recorrente → nova notificação** (**009-15 A**) — mensurável no módulo Notificação (**008-17 C**); o aviso MUST sinalizar **recorrente**.

## Evidência e consulta

- **Fato**: Decisões 16/11 em `respostas-joao.md` (incl. **Etapa 3** 00h41, **Etapa 3b** 09h02, **Etapa 3c** 11h/11h01); `respostas-etapa3.md`; `respostas-etapa3c.md` (§5–§6); revisão T1 em `revisao-T1-drive.md`; LLT-201, **LLT-199** (precisão via 006/007, não gate nesta spec), **LLT-218**; ADR_054, 055, 059, 060, 062, 073, 051, 013, 143, 146 v1, 165; SSID_009.1; sequencia-notificacao.md v2; contrato log.ingested v1 (**T5 D**); Finder.mmd.
- **Premissa**: Nenhuma premissa inventada (constituição V).
- **Pergunta remanescente**: Ver **Dúvidas novas (sem fonte)** (**DN-009-2** API Módulo) e C8 abaixo.

## Dúvidas novas (sem fonte)

### DN-009-1 — respondida em 2026-11-16 (etapa 3b), ver Clarifications → Session 2026-11-16 (etapa 3b)

### DN-007-2 — respondida em 2026-11-16 (etapa 3c), ver Clarifications → Session 2026-11-16 (etapa 3c)

*Escopo nesta spec: referência cruzada a LLT-199/006; gate ≥95% por classe permanece na spec **006** / **007**, não em SC-F00x.*

### DN-008-4 — respondida em 2026-11-16 (etapa 3c), ver Clarifications → Session 2026-11-16 (etapa 3c)

### DN-008-5 — respondida em 2026-11-16 (etapa 3e), ver Clarifications → Session 2026-11-16 (etapa 3e)

*Literais **DN-008-5 A** reutilizados nesta spec (**FR-F11**); dono canônico: spec **008**.*

### DN-008-3 — respondida em 2026-11-16 (etapa 3h), ver Clarifications → Session 2026-11-16 (etapa 3h)

*Falha consulta OpenObserve: **3 tentativas**; **sem DLQ** no ingested; depois Finder → Core API (**batchId**, **motivo**). Substituí alcance «3 tentativas + DLQ» (§Etapa 3c).*

### DN-008-6 — respondida em 2026-11-16 (etapa 3h), ver Clarifications → Session 2026-11-16 (etapa 3h)

*Canal operacional após falha OpenObserve esgotada: Finder → **Core** API interna (**DN-008-3**); **não** consumidor DLQ nem **WebSocket** pós-DLQ (**FR-F07**). Substituí §Etapa 3e/3c. Não confundir com aviso LLT-201 Módulo → Core.*

### DN-008-9 — respondida em 2026-11-16 (etapa 3h), ver Clarifications → Session 2026-11-16 (etapa 3h)

*Escape de **`|`** nos componentes da assinatura: **`\|`**. Escape da barra invertida **`\`** nos componentes: **[PENDENTE]** sem fonte (não inventar).*

### DN-008-7 — respondida em 2026-11-16 (etapa 3e), ver Clarifications → Session 2026-11-16 (etapa 3e)

*Envio: um **ResumoAnalise** por lote; consulta **N2 B**: **por achado** (**DN-008-7 B**). Contrato HTTP: **DN-009-2**.*

### DN-009-2 — Contrato HTTP Finder → Módulo Ocorrência (rota, JWT, respostas, 5xx/timeout)

As questões abaixo registram as alternativas antes da aprovação. A decisão de 20/11/2026 define as duas rotas HTTPS, o corpo com chave de lote, os recibos `201` e `200` e um recálculo por revisão desatualizada. [Contrato aprovado](contracts/proposta-dn-009-2.md)

| Opção | Descrição | Prós | Contras |
|-------|-----------|------|---------|
| **10a-A** | Uma chamada HTTP por **ResumoAnalise** (lote inteiro) | Alinha a N4 e seq. v6 §2 | Payload grande |
| **10a-B** | Uma chamada por **achado** | Payload menor | Muitas round-trips; idempotência mais complexa |
| **10a-C** | Definir rota apenas numa **ADR do Módulo** | Dono claro | Bloqueia implementação até ADR |
| **10b-C** | Fechar pendências **ADR_028** (algoritmo, rotação) antes do plan | Segurança explícita | Atraso |
| **10c-A** | Sucesso = criada/atualizada; repetido = sucesso idempotente (**ADR_014**) | Clara para `finishScan()` | Sem 202 assíncrono |
| **10d-A** | 5xx/timeout: mesma regra **009-18 A** (3 tentativas + DLQ) | Consistente com falha de publicação | Sem alerta Core extra |
| **10d-B** | 10d-A + **alerta ao Core** | Operação visível | Decisão nova sem João |

O cliente Finder está implementado. A ativação depende da API do Módulo, da configuração JWT do ambiente e das demais integrações de produção.

## Decisões tomadas (com fonte)

| ID | Decisão | Fonte arquivo | Origem documental |
|----|---------|---------------|-------------------|
| T1 revisto | Finder **não** publica evento; envia **ResumoAnalise** ao Módulo Ocorrência; Módulo cria ocorrência e avisa Core | respostas-joao.md; revisao-T1-drive.md | ADR_062; sequencia-varredura-finder v6; SSID_009.1 |
| N1 A | Finder → Módulo: **HTTP + JWT** | respostas-joao.md | Paridade 009-3 D |
| N2 B | Agrupamento entre lotes no **Finder** (consulta prévia) | respostas-joao.md | LLT-200 / 008 |
| N3 A | Ampliar **AchadoAnalise** (assinatura, contagem, primeira/última, versão) via **ADR nova** | respostas-joao.md | Finder.mmd; lacunas Drive |
| N4 A | Idempotência ResumoAnalise: **batchId + execId** (ou só **batchId** sem execId — **DN-006-1b B**) | respostas-joao.md; Etapa 3 | contrato 08/11; ADR_070 |
| DN-006-1 | **execId** da Varredura (`V`/`M`+uuid); **opcional** (“quando houver”) | respostas-joao.md §Etapa 3; respostas-etapa3.md §3 | ADR_075, ADR_069, ADR_146 |
| DN-006-1b B | Sem **execId** enquanto só **log.ingested** v1: chave handoff = **batchId** | respostas-joao.md §Etapa 3 | contrato 08/11; ADR_070 |
| DN-009-1 A | E-mail: **Notificação** lê Ocorrência por **id** (MongoDB **ADR_054**) | respostas-joao.md §Etapa 3 | LLT-218; sequencia-notificacao v2 |
| T2 C | Um achado elegível por classe ERROR, WARNING ou INFO no ResumoAnalise (**T2 C** + **T7 A**) | respostas-joao.md | transversal; 006-11 A |
| spam: B | Contador “spam ignorado” (ADR_102) → **contador de INFO** no ResumoAnalise; nome canônico na ADR_102 [PENDENTE] | respostas-joao.md §Etapa 3 | ADR_102, ADR_170 |
| T3 B | Recorrência via OpenObserve | respostas-joao.md | ADR_168 |
| T4 A | Lote parcial não é sucesso | respostas-joao.md | spec 005 |
| T5 D | **log.ingested v1** única versão válida | respostas-joao.md | contrato-log-ingested-v1.txt |
| T6 A | Lacunas 006 na própria spec 006 | respostas-joao.md | — |
| T7 A | **INFO** gera ocorrência (achado no ResumoAnalise) | respostas-joao.md | ADR_170; diverge Drive |
| Taxonomia 14/11 | **ERROR/WARNING/INFO**; sem SPAM nem “não classificado” | respostas-joao.md | constituição III; ADR_170 |
| 009-1 B | Destino do aviso: **Core → Notificação** | respostas-documentadas.md; revisao-T1-drive.md | ADR_165; guia-redis v3 |
| 009-2 | “Avisar Ocorrências” = **ResumoAnalise** | revisao-T1-drive.md | ADR_062 |
| 009-3 D | Módulo → Core: **HTTP + JWT** | respostas-joao.md | ADR_051 |
| 009-4 B | Corpo HTTP: **envelope ADR_146** | respostas-joao.md | ADR_146 v1 |
| 009-5 revisto | Aviso Módulo → Core: **id + tipo + gravidade**; campo `ocorrencia` no envelope = **id** | respostas-joao.md §Etapa 3 | sequencia-notificacao v2; ADR_040 |
| 009-5 A (14–16/11) | Decisão anterior “aviso só id” — **substituída** pelo **009-5 revisto** | respostas-joao.md | histórico; ver D12 |
| API Módulo | Handoff Finder→Módulo: **pendência** até contrato (**DN-009-2**) | respostas-joao.md §Etapa 3 | ADR_028, ADR_054; respostas-etapa3 §10 |
| 009-6 A | Id próprio sem fechar chave C8 | respostas-documentadas.md | ADR_040 |
| 009-8 | **ResumoAnalise** mantido | revisao-T1-drive.md | ADR_062 |
| 009-9 A | Idempotência HTTP: **id + estado** | respostas-joao.md | — |
| 009-10 B | Idempotência stream: **consumidor** / eventId | respostas-documentadas.md | guia §8 |
| 009-11 B | Duplicata = **2º alerta no Core** | respostas-joao.md | — |
| 009-12 rev. | **Gravidade** no aviso do Módulo | revisao-T1-drive.md | SSID_009.1 passo 1 |
| 009-13 A | **Tipo** = classe ERROR/WARNING/INFO | respostas-documentadas.md | ADR_170 |
| 009-14 A | Reabertura → **atualizada** | respostas-joao.md | 008-4 B |
| 009-15 A | Recorrente → **nova notificação** | respostas-joao.md | SSID_009.1 FA01 |
| 009-16 rev. | Aviso na criação/mudança; Sugestão antes do Core se aplicável | revisao-T1-drive.md | SSID_009.1 |
| 009-17 | **Um aviso por ocorrência** | revisao-T1-drive.md | SSID_009.1 passo 1 |
| 009-18 A | Falha HTTP: **retry e alerta**; Finder: FA-01/FA-02 | respostas-joao.md; revisao-T1-drive.md | ADR_014; seq. v6 |
| 009-19 rev. | Mesmo canal alerta: Módulo → Core JWT | revisao-T1-drive.md | ADR_051; arquitetura v2 |
| 009-20 | Transição 006 via **T6 A**; Finder produz resumo, não evento | respostas-joao.md | 006 FR-009 |
| 009-21 | Anti-duplicata: **N4 A** + **009-9 A** + confirmação Módulo | respostas-joao.md; revisao-T1-drive.md | FA-02 |
| 009-22 A | **009 espera 008** | respostas-joao.md | backlog E07.F3 ordem |
| DN-007-2 B | Meta **≥95% precisão por classe** (ERROR, WARNING, INFO) na base 007 — **sem** SC próprio nesta feature; ver spec **006** / **007** | respostas-joao.md §Etapa 3b | respostas-etapa3.md §1; LLT-199; 006-2 B |
| DN-008-4 B | Normalização da mensagem na assinatura: trim, case-insensitive, marcadores número/UUID/IP/data/hex; dono spec **008**; 009 reutiliza (**FR-F11**) | respostas-joao.md §Etapa 3b | DN-008-1 B; ADR_013 |
| DN-008-5 A | Literais marcadores: `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` | respostas-joao.md §Etapa 3c 11h01/11h | DN-008-4 B |
| 008 assinatura A | Valor da assinatura: texto concatenado (**DN-008-1 B**); separador fixo **`|`** (**DN-008-8 A**; dono canônico spec **008**) | respostas-joao.md §Etapa 3c 11h; **DN-008-8 A**: §Etapa 3e (João, 16/11 12h40) | DN-008-1 B; 008-11 |
| DN-008-8 A | Separador **`|`** na serialização **008 assinatura A** (reuso **FR-F11**; dono spec **008**) | respostas-joao.md §Etapa 3e (João, 16/11 12h40) | 008 assinatura A |
| **DLQ B** | **Sem DLQ** no **log.ingested v1**; não copiar retry/DLQ da FilaOrdenada para o ingested | respostas-joao.md §Etapa 3g (João, 16/11 13h22) | contrato 08/11; respostas-etapa3g.md §1b; DN-006-1b B |
| **DN-008-3** (rev.) | Falha consulta **OpenObserve**: **3 tentativas**; depois Finder avisa **Core** (API interna, **`batchId`**, **motivo**); **sem DLQ** — **substitui** alcance «3 tentativas + DLQ» | respostas-joao.md §Etapa 3g, DLQ B | ADR_168 / T3 B; contrato ingested |
| **DN-008-6** (rev.) | Após esgotar tentativas OpenObserve: aviso **Finder → Core** (API interna); **não** WebSocket pós-DLQ — **substitui** §Etapa 3c 11h | respostas-joao.md §Etapa 3g; respostas-etapa3g.md §1 | ADR_051 [PENDENTE]; FR-F07 |
| **DN-008-9 A** | **`|`** literal no componente da assinatura → **`\|`** na serialização | respostas-joao.md §Etapa 3g (João, 16/11 13h22) | DN-008-8 A; 008 assinatura A |
| DN-008-7 | Handoff: **um ResumoAnalise por lote** | respostas-etapa3c.md §6 | ADR_062; Finder.mmd |
| DN-008-7 B | Consulta prévia (**N2 B**): **uma chamada por achado** ao Módulo Ocorrência | respostas-joao.md §Etapa 3c (João, 16/11 11h02) | N2 B; contrato **DN-009-2** |

## Divergências documentadas

| Id | Fontes | O que diverge | Qual lado vale | Por quê |
|----|--------|---------------|----------------|---------|
| D1 | **LLT-201** (épico Finder, “publicar evento… avisar Ocorrências e Notificação”) × **Drive** (ADR_062, SSID_009.1, arquitetura v2) | Finder publicaria e avisaria direto × Finder só **ResumoAnalise**; aviso Notificação via Módulo e Core | **Drive** (com T1 revisto alinhado ao Drive) | Sequência varredura v6; Módulo emissor SSID_009.1 |
| D2 | **LLT-201** critério 1 (sem gravidade) × **SSID_009.1** passo 1 | Evento mínimo Jira × aviso com **gravidade** | **Drive** para aviso do **Módulo** | ADR_062 atribui gravidade ao Módulo |
| D3 | **T7 A / 006-11 A** (João 16/11) × **SSID_007.1**, **SSID_009.1**, **ADR_040** (“só ERRO/WARNING”) | INFO vira ocorrência × Drive filtra INFO | **João** (decisão consciente) | Rodada 16/11 00h01 |
| D4 | **T1 A** (João 14/11, Finder grava/publica) × **Drive** | Persistência e publicação no Finder × Módulo cria e avisa | **T1 revisto / Drive** | 15/11 23h47 |
| D5 | **006-9 A** (upsert Finder) × **T1 revisto** | Upsert no Finder × criação no Módulo | **T1 revisto** | ADR_040, ADR_062 |
| D6 | **008-11 B** (chave sem classe) × **008-15 A** (uma ocorrência por assinatura **e classe**) | Chave de unicidade | **João** em ambos; **recomenda-se** chave **(assinatura, classe, produto, tenant)** na 008 | Rodada 16/11: 008-11 com classe |
| D7 | **ADR_146 v2** (“schemaVersion 2”, quebra v1) × **contrato 08/11** e **T5 D** | Versão do ingested | **Contrato 08/11 v1** | Nota de revisão 08/11; spec 004 |
| D8 | **ADR_019 / ADR_167 / SSID_007.1** (legado ERRO/WARNING/SPAM) × **ADR_170** e João 14/11 | Taxonomia antiga | **ADR_170 + João** | Não emendar ADRs nesta feature (**FR-F10**) |
| D9 | **008-2 A** (correlationId na assinatura) × **contrato 08/11 §4** (correlationId aleatório por exportação) | Estabilidade da assinatura entre exportações | Registrado como **risco**; sem decisão nova | Afeta 008 e conteúdo do AchadoAnalise, não o formato do aviso |
| D10 | **respostas-documentadas 009-1** (emissor Finder por T1 A) × **revisão T1** | Quem publica | **revisao-T1-drive** + João | Emissor = Módulo Ocorrência |
| D11 | **respostas-documentadas 009-12 A** (sem gravidade) × **revisão 009-12** | Gravidade no evento | **revisão 009-12** (Módulo emissor) | SSID_009.1 passo 1 |
| C8 | **ADR_035 / ADR_040 / ADR_161** | Chave concreta da ocorrência **[PENDENTE]** | Ainda **sem decisão** | 009-6 A permite id sem fechar chave |
| D12 | **009-5 A** (“só o id” no aviso) × **sequencia-notificacao.md v2** (diagrama 1: tipo e gravidade no alerta) | Conteúdo mínimo do HTTP ao Core | **009-5 revisto** (João, Etapa 3) | Notificação usa tipo/gravidade para destinatários (ADR_130/136/137/142); ressalva respostas-etapa3 §7 |
| D13 | **Recomendação DN-009-1 C** (spec 14–16/11, leitura `conteudo`) × **DN-009-1 A** (MongoDB por id) | Onde a Notificação lê o texto | **DN-009-1 A** (João) | Opções A e C equivalentes na prática com **ADR_054**; §7 registra ressalva |
| D14 | **LLT-218** (e-mail com resumo) × **Drive** (templates **[PENDENTE]** NOT-03) | Corpo do e-mail | **Jira** para obrigação; **Drive** para lacuna operacional | Não é contradição; falta template, não origem dos dados (**DN-009-1 A**) |
| D15 | **N4 A** texto “batchId + execId” × **contrato ingested v1** (sem execId no envelope) | Chave handoff no fluxo atual | **DN-006-1b B** | Até `log.checked` ou varredura expor execId; ressalva respostas-etapa3 §3 |

## Questões em aberto (remanescentes)

1. **Chave concreta** da ocorrência (C8, ADR_040 **[PENDENTE]**) — id opaco já permitido (**009-6 A**).
2. **DN-009-2** — implementar o servidor do Módulo, configurar o JWT e homologar a integração com o Finder. O cliente segue o [contrato aprovado](contracts/proposta-dn-009-2.md).
3. **API Finder → Core** após falha OpenObserve esgotada (**DN-008-3**/**DN-008-6**): rota, payload e evento **[PENDENTE]** (**ADR_051 [PENDENTE]**); **não** inventar contrato nesta feature.
4. **Commit de alinhamento** do Coletor ao contrato v1 (**T5 D**, contrato §12) — fato pendente, não decisão de design da 009.
5. **N3 A** — os campos de transporte são definidos pelo OpenAPI aprovado. A publicação e a conciliação da ADR do Módulo seguem externas a este repositório.
6. ~~**Separador fixo** da serialização **008 assinatura A**~~ — **DN-008-8 A**: literal **`|`** propagado para reuso; dono canônico spec **008** (ver Clarifications etapa 3g).
7. **Escape de `\`** nos componentes da assinatura (**DN-008-9**) — **[PENDENTE]** sem fonte; vale apenas **`|`** → **`\|`** (**DN-008-9 A**).
8. **Dono canônico** dos marcadores **DN-008-5 A** — documentados na spec **008**; 009 só reutiliza (**FR-F11**).
9. **Templates** de e-mail na Notificação (Drive NOT-03) — fora do Finder; leitura da Ocorrência decidida (**DN-009-1 A**).
10. **Emenda do princípio I** — a constituição ainda descreve agregação/persistência de ocorrência no Finder; prevalece **T1 revisto** (handoff + Módulo cria ocorrência) nesta spec; emenda em `.specify/memory/constitution.md` **fora** desta pasta (**C1** analyze).

---

*Não há marcadores `[NEEDS CLARIFICATION]` nesta versão do spec.*