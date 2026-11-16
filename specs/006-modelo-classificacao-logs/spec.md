# Feature Specification: Modelo de classificação de logs

**Feature Branch**: `006-modelo-classificacao-logs`

**Created**: 2026-11-09

**Status**: Draft (atualizada 2026-11-16 — etapa 3h: DLQ B, rev. DN-008-3/DN-008-6, DN-008-9 A; etapa 3g: DN-006-3-ref, DN-008-8 A, 007 H2 D; etapas anteriores preservadas em Clarifications)

**Input**: Classificar logs capturados sem LLM (nível + texto + catálogo opcional), agrupar em achados, produzir `AchadoAnalise` e `ResumoAnalise` e enviar ao Módulo de Ocorrência por HTTP interno com JWT. O Finder **não** grava ocorrência no MongoDB nem publica evento de ocorrência. (fonte: respostas-joao.md, T1 REVISTO; origem: revisao-T1-drive.md §2, ADR_062)

**Taxonomia de classes nesta spec**: os rótulos de negócio **erro**, **alerta** e **informação** (constituição, spec original) correspondem a **ERROR**, **WARNING** e **INFO** (ADR_170, LLT-199). O campo `tipo` no achado usa ERROR/WARNING/INFO (fonte: respostas-joao.md, 006-6 A; origem: ADR_170).

## Clarifications

### Session 2026-11-09

*(Várias respostas desta sessão foram **substituídas** em 2026-11-16; cada item abaixo indica a decisão vigente.)*

- Q: Qual conjunto de classes vale nesta feature? → A: Informação, alerta e erro. Não existe spam nem "não classificado". A constituição deste repositório permanece como está. A quantidade ficou uma ocorrência por apontamento, não uma por log. **Substituído (grão):** uma ocorrência **por achado** (T2 C), inclusive INFO (T7 A, 006-11 A). (fonte: respostas-joao.md, T2 C/T7 A/006-11 A; origem: ADR_040)
- Q: Quem grava a ocorrência e o que ela contém? → A: O Finder grava no MongoDB a classe, a versão das regras e a identidade do apontamento, com o tenant e o produto. A ocorrência fica consultável sem copiar o texto dos logs. Mascaramento, gravidade, evento e sugestão continuam fora. **Substituído (T1 REVISTO):** o Finder **não** grava ocorrência; envia `ResumoAnalise` com `AchadoAnalise` ao Módulo de Ocorrência, que cria a ocorrência (fonte: respostas-joao.md, T1 REVISTO; origem: ADR_062, revisao-T1-drive.md A.1)
- Q: Quando o nível do log e o catálogo apontam classes diferentes, qual vale? → A: Vale o nível do log. O catálogo só classifica quando o log não traz nível reconhecível. A resposta D seguinte inclui o campo e o texto nesse nível; a classe do log é a mais grave entre os dois. O catálogo continua de fora quando um dos dois já decide. **Vigente** (regras de classificação por log inalteradas).
- Q: Quando um apontamento da fila já trouxe vários logs capturados, quantas ocorrências esta feature grava? → A: O apontamento inteiro vira uma única ocorrência, mesmo quando há vários logs. **Substituído:** um achado por assinatura (dentro do lote) → um `AchadoAnalise` → o Módulo Ocorrência cria **uma ocorrência por achado** (T2 C). (fonte: respostas-joao.md, T2 C; origem: ADR_035/ADR_040)
- Q: Quando os logs de um mesmo apontamento recebem classes diferentes, qual classe fica na ocorrência única? → A: A mais grave. Erro vence alerta, e alerta vence informação. **Substituído:** não há mais “classe da ocorrência única por apontamento”; cada achado tem sua própria classe; logs com assinaturas/classes diferentes geram achados distintos (T2 C, 008-15 A). (fonte: respostas-joao.md, 008-15 A)
- Q: O nível que decide a classe de um log é um campo próprio desse log, ou uma palavra igual encontrada no texto da mensagem? → A: Campo e texto valem juntos. A classe do log é a mais grave entre os dois. **Vigente**.
- Q: Quando a mensagem traz um termo de nível no meio de uma palavra maior, esse termo conta para a classe? → A: Só a palavra inteira. "error" conta; "errorless" e "traceback" não. **Vigente**.
- Q: Quando um apontamento tem vários logs e um deles não está disponível, a feature grava a ocorrência com os demais? → A: Grava uma ocorrência só com os logs presentes. A classe é a mais grave entre eles. **Substituído (T4 A):** lote incompleto **não** conta como sucesso; o Finder não deve tratar o lote como concluído com sucesso nem enviar `ResumoAnalise` como entrega final de lote completo. Classificação dos logs presentes pode ocorrer, mas o critério de sucesso do lote é falha. (fonte: respostas-joao.md, T4 A; origem: spec 005)

### Session 2026-11-16

- Q: T1 Quem grava a ocorrência e publica o evento? → A: **T1 REVISTO** — Módulo de Ocorrência cria a ocorrência; Finder envia `ResumoAnalise` por HTTP interno com JWT (N1 A). Finder não publica evento de ocorrência. (fonte: respostas-joao.md, T1 REVISTO/N1 A; origem: ADR_062, revisao-T1-drive.md §1–2)
- Q: T2 Grão da ocorrência? → A: Uma ocorrência por **achado** (ERROR, WARNING ou INFO). (fonte: respostas-joao.md, T2 C; origem: ADR_040)
- Q: T4 Lote parcial? → A: Lote incompleto não é sucesso. (fonte: respostas-joao.md, T4 A)
- Q: T5 Versão do log.ingested? → A: Contrato 08/11 — somente **v1** vale; `schemaVersion` "1". (fonte: respostas-joao.md, T5 D; origem: contrato-log-ingested-v1.txt)
- Q: T6 Lacunas 006 × LLT-197/199? → A: Atualizar a própria spec 006. (fonte: respostas-joao.md, T6 A)
- Q: T7 INFO gera ocorrência? → A: Sim, decisão consciente mesmo divergindo do Drive (SSID_007.1, SSID_009.1, ADR_040). (fonte: respostas-joao.md, T7 A e rodada 16/11 00h01; origem: ADR_170)
- Q: 006-2 Métricas de qualidade? → A: Precisão e recall/F1 **por classe**. (fonte: respostas-joao.md, 006-2 B; origem: LLT-197, LLT-199)
- Q: 006-3 Versão das regras? → A: Versão **única** regras+catálogo+base, registrada no `AchadoAnalise` e/ou `ResumoAnalise` (não na ocorrência gravada pelo Finder). (fonte: respostas-joao.md, 006-3 C e rodada 16/11; origem: LLT-199 critério 3, ADR_062)
- Q: 006-4 Reprocessar ao mudar versão? → A: Só novos logs (sem reprocessamento automático retroativo nesta feature). (fonte: respostas-joao.md, 006-4 A; origem: LLT-51)
- Q: 006-6 tipo/severidade LLT-199? → A: O próprio **tipo** (ERROR/WARNING/INFO); gravidade fica fora (Módulo Ocorrência, ADR_062). (fonte: respostas-joao.md, 006-6 A; origem: ADR_062)
- Q: 006-7 Sintaxe e dono do catálogo? → A: Definir na LLT-51 ou ADR nova — **pendência explícita**; esta execução **não** cria ADR nem fecha sintaxe. (fonte: respostas-joao.md, 006-7 D; origem: ADR_167 [PENDENTE])
- Q: 006-9 Upsert por achado? → A: Opção A do João, **sem efeito no Finder** após T1 revisto; upsert/idempotência ao criar ocorrência é do Módulo Ocorrência. (fonte: respostas-joao.md, 006-9 A; origem: revisao-T1-drive.md A.1, ADR_040)
- Q: 006-10 Alinhar Jira LLT-197/199? → A: Sim, junto da atualização da 006 — **pendência explícita**; esta execução **não** altera o Jira. (fonte: respostas-joao.md, 006-10 A)
- Q: 006-11 Grão das ocorrências INFO? → A: Por achado INFO, igual a ERROR/WARNING (T2 C + T7 A). (fonte: respostas-joao.md, 006-11 A)
- Q: N1 Transporte Finder → Módulo Ocorrência? → A: HTTP interno com JWT (como 009-3 D no padrão ADR_051). (fonte: respostas-joao.md, N1 A; origem: ADR_051, revisao-T1-drive.md)
- Q: N3 Ampliar AchadoAnalise (assinatura, contagem, primeira/última, versão)? → A: Sim, com ADR — **ADR não deve ser criada nesta execução**; requisitos funcionais ficam na spec; formalização do contrato pendente. (fonte: respostas-joao.md, N3 A)
- Q: N4 Chave do ResumoAnalise? → A: **`batchId` obrigatório** + **`execId` quando houver** (alinha ao log.ingested v1 e DN-006-1b B); no fluxo só `log.ingested` v1, chave **somente `batchId`**. (fonte: respostas-joao.md, N4 A / DN-006-1b B; origem: contrato-log-ingested-v1.txt, ADR_070, ADR_146)

### Session 2026-11-16 (etapa 3b)

- Q: DN-006-1 De onde vem o `execId` na chave do `ResumoAnalise` e é sempre obrigatório? → A: Da **Varredura** (`V<uuid>` automática, `M<uuid>` manual); **opcional** — chave `batchId` + `execId` **quando houver** `execId` (fonte: respostas-joao.md, Etapa 3 / respostas-etapa3.md §3; origem: ADR_075, ADR_069, ADR_070, ADR_146)
- Q: DN-006-1b Como o Finder obtém `execId` enquanto consome só `log.ingested` v1? → A: **B** — sem `execId` neste fluxo; chave de idempotência **somente `batchId`** até `log.checked` ou outro canal entregar `execId` (fonte: respostas-joao.md, Etapa 3 / respostas-etapa3.md §3; origem: contrato-log-ingested-v1.txt §§2–5)
- Q: DN-008-1 Quais campos compõem a assinatura do achado no escopo desta feature? → A: **B** — nível + mensagem normalizada + `service.name` + host/máquina (fonte: respostas-joao.md, Etapa 3 / respostas-etapa3.md §4; origem: ADR_013; sem `traceId`/`correlationId` nesta lista)
- Q: Contador «spam ignorado» (ADR_102) com taxonomia sem SPAM (ADR_170)? → A: **B** — trocar por contador de **INFO** no `ResumoAnalise` (fonte: respostas-joao.md, Etapa 3 / respostas-etapa3.md §9; origem: ADR_102, ADR_170)
- Q: API do Módulo Ocorrência (rota, JWT, HTTP, 5xx/timeout) para o handoff? → A: **Pendência no plan** até contrato do Módulo sair; Drive só fixa API interna + JWT próprio (fonte: respostas-joao.md, Etapa 3 / respostas-etapa3.md §10; origem: ADR_028, ADR_054, ADR_141)
- Q: 009-5 (aviso ao Core) — o que o Finder precisa alinhar? → A: **Revisto** — aviso leva **id + tipo + gravidade** (publicado pelo Módulo Ocorrência, fora do Finder); Finder envia `tipo` no `AchadoAnalise`, sem gravidade (fonte: respostas-joao.md, Etapa 3; origem: sequencia-notificacao.md v2, ADR_062)

### Session 2026-11-16 (etapa 3c)

- Q: DN-007-2 Como a meta **007-meta** (≥95% de precisão) se aplica ao gate de avaliação desta feature? → A: **B** — **≥95% de precisão por classe** (ERROR, WARNING e INFO), medida na base de validação da spec 007; a 006 referencia o gate em **SC-007** e tarefa **T047** (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: respostas-etapa3.md §1, LLT-199, 006-2 B)
- Q: DN-008-4 Qual algoritmo normaliza a mensagem na assinatura do achado? → A: **B** — aparar espaços nas extremidades; ignorar maiúsculas/minúsculas; substituir **números**, **UUIDs**, **IPs**, **datas** e **hexadecimais** por **marcadores** (um por categoria). Algoritmo **canônico** na spec **008**; specs **006**, **009** e **010** MUST reutilizá-lo ao derivar assinatura ou equivalência (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: DN-008-1 B, ADR_013)
- Q: 008-16 A O Finder deve tratar `productId` ausente ou inconsistente entre metadata validada e corpo do `log.ingested` v1? → A: **Rejeitar** como **falha crítica permanente**, **sem retry** (padrão ADR_049); motivo distinguível em **FR-008** / **SC-003** (fonte: respostas-joao.md, 008-16 A; origem: contrato-log-ingested-v1.txt §4, ADR_049)
- Q: F1 analyze2 Onde reside a implementação canônica do normalizador **DN-008-4 B** e o que a spec 006 faz? → A: Dono **spec 008** — `internal/application/assinatura/normalizar_mensagem.go` (definido em `specs/008-agrupar-logs-ocorrencia/plan.md` e **tasks.md** T016); a 006 **MUST** reutilizar (delegar/importar), **MUST NOT** reimplementar cópia divergente (ex. em `internal/domain/`) (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02, DN-008-4 B; origem: ADR_013, specs/008-agrupar-logs-ocorrencia/plan.md)

### Session 2026-11-16 (etapa 3c — correções analyze)

- Q: F2 A linha N4 na sessão 2026-11-16 carece da ressalva «quando houver» para `execId`? → A: Texto de N4 alinhado a **FR-011** e **DN-006-1b B** — `execId` opcional; chave só `batchId` no pipeline atual. (fonte: respostas-joao.md, N4 A / DN-006-1b B; origem: respostas-etapa3.md §3, contrato-log-ingested-v1.txt)
- Q: E2 Como tratar log sem `service.name` ou host na assinatura (DN-008-1 B)? → A: Componentes **obrigatórios** para agregar em achado; log contável sem um deles → falha distinguishável **`assinatura_incompleta`**; lote MUST NOT contar como varredura concluída com sucesso (padrão **007-8 A** / `identidade_incompleta`). (fonte: respostas-joao.md, DN-008-1 B; origem: ADR_013, respostas-etapa3.md §4)
- Q: E3 Como preencher `recorrente` e campos Drive (`titulo`, `resumo`, …) no escopo mínimo «só lote» (FR-009)? → A: **`recorrente` = `false`** nesta feature (sem consulta entre lotes — **FR-009**, **N2 B** na spec **008**); avaliação de recorrência via OpenObserve na spec **008** (**008-3/5/6**; **DN-008-3** rev. 3h: até 3 tentativas, **sem DLQ** no `log.ingested`; após esgotar, Finder→Core API interna **batchId** + **motivo** — MUST NOT marcar `recorrente=false` por degradação transitória da consulta). Campos Drive obrigatórios no payload; **`hostName`** alinhado 1:1 ao `host` da assinatura quando presente; **`titulo`/`resumo`/`numeroLinha`/`camposExtraido`**: presença exigida, **semântica de preenchimento mínimo pendente N3 A** (revisao-T1-drive.md §2, ADR_040). (fonte: respostas-joao.md, N2 B / N3 A; §Etapa 3g (João, 16/11 13h22), DLQ B; origem: respostas-etapa3.md §6, ADR_040, ADR_063)
- Q: E4 Como fechar aceite LLT-199 quanto a recall/F1 no gate por classe (DN-006-2 / DN-007-2)? → A: Gate formal de aceite (**SC-007**, **FR-010**, T047): em **cada** classe ERROR, WARNING e INFO — **precisão ≥ 95%** (**DN-007-2 B**), **recall ≥ 95%** e **F1 ≥ 95%** (**DN-006-2 B**). Alinhado ao normativo vigente; substitui a redação provisória «somente precisão» da analyze etapa 3c. (fonte: respostas-joao.md, DN-006-2 B / DN-007-2 B; origem: LLT-199, 006-2 B)
- Q: I1 A meta de **2 s** em SC-002 é normativa e testável? → A: **Sem fonte** para ambiente/volume; nota herda spec 09/11 — tratada como **não normativa** até **DN-006-3** (ver § Dúvidas novas). (fonte: analyze I1; analogia processual **010-28 A** — sem meta de tempo — não decide)
- Q: D2 A ordem textual FR-011/FR-012/FR-013 está invertida? → A: Reordenar blocos FR sem alterar numeração (editorial). (fonte: analyze D2)
- Q: D3 Regras nível+texto duplicadas entre User Story 1 e FR-004? → A: US1 referencia **FR-004** como norma única de classificação por log (editorial). (fonte: analyze D3; spec sessão 2026-11-09)
- Q: C1 Conflito princípio I × T1 REVISTO? → A: **Pendência «emenda do princípio I»** — não alterar constituição nesta pasta; spec/plan/tasks rastreiam proposta (**T043**). (fonte: respostas-joao.md, T1 REVISTO; origem: ADR_062, plan.md § Constitution Check)

### Session 2026-11-16 (etapa 3e)

- Q: DN-006-2 Quais limiares numéricos de **recall** e **F1** entram no gate de aceite por classe? → A: **B** — meta **≥ 95% de recall** e **≥ 95% de F1** em **cada** classe ERROR, WARNING e INFO, **além** da precisão **≥ 95% por classe** já fixada (**DN-007-2 B**). **Substitui** a exceção documentada em E4 (etapa 3c analyze) nos requisitos **FR-010** e **SC-007**. (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h e 11h01; respostas-etapa3c.md §1; origem: LLT-199, 006-2 B)
- Q: DN-006-3 Qual meta de tempo normativa para classificar e montar o resumo de um lote? → A: **B** — **2 s por lote** para o ciclo **classificar + montar `ResumoAnalise`** (**SC-002**); **tamanho de lote** e **máquina de referência** citados na redação da opção B **sem valores nas fontes** — **pendência de documentação** (**DN-006-3-ref**, sem inventar números). (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-006-3 B; respostas-etapa3c.md §2 opção B)
- Q: DN-008-5 Quais literais de marcador o normalizador **DN-008-4 B** MUST usar? → A: **A** — `<NUM>`, `<UUID>`, `<IP>`, `<DATE>` e `<HEX>` (um por categoria). Norma canônica na spec **008** / T052; a 006 **referencia** e reutiliza a implementação (**FR-013**). (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h e 11h01; respostas-etapa3c.md §4)
- Q: DN-008-7 Qual o grão do envio ao Módulo e da consulta prévia (N2 B)? → A: **Envio** — **um** `ResumoAnalise` **por lote** (fonte: respostas-etapa3c.md §6; origem: ADR_062, Finder.mmd). **Consulta prévia ao Módulo** — **uma chamada por achado** (**DN-008-7 B**). (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; respostas-etapa3c.md §6)
- Q: 008 assinatura A Como serializar o **valor** da assinatura agregada? → A: **A** — texto com os componentes (**nível + mensagem normalizada + `service.name` + host**) **concatenados** por um **separador** fixo; o **separador exato** não tem fonte nesta feature — decisão e DN na spec **008** (esta spec **só referencia**). (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h e 11h01; respostas-etapa3c.md §9)

### Session 2026-11-16 (etapa 3f)

- Q: A divergência D13 (11h01 «puladas» × 11h) ainda vale após SUPERADO em respostas-joao.md? → A: **Não** — removida de **Divergências documentadas** e demais menções; DN-006-3 B, DN-007-4 A e DN-008-7 B respondidas em **11h02**, sem conflito entre seções. (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h01 linha SUPERADO; instrução etapa 3f)
- Q: Qual fonte canônica para **DN-006-3 B**, **DN-007-4 A** e **DN-008-7 B** na tabela de decisões? → A: **respostas-joao.md §Etapa 3c (João, 16/11 11h02)** (com **respostas-etapa3c.md §2** / **§6** e **ADR_062** onde já citados para 006-3/008-7). (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); respostas-etapa3c.md §2/§6)
- Q: Achado **I1** (nota «consulta ocorrências prévias» fora de FR mínima) — como alinhar a **FR-007** / **DN-008-7 B**? → A: Redação corrigida na nota final de **tasks.md**: consulta **por achado** ao Módulo é **FR-007**; fora do mínimo permanecem agrupamento **entre lotes** e recorrência OpenObserve (spec **008**). (fonte: analyze-006-etapa3e.md I1; respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; respostas-etapa3c.md §6; ADR_062)

### Session 2026-11-16 (etapa 3g)

- Q: H2 (analyze etapa 3f) O bloco **E4** em Clarifications ainda contradiz **FR-010**/**SC-007** ao descrever gate somente com precisão ≥95%? → A: Redação de **E4** (sessão etapa 3c — correções analyze) atualizada para gate **precisão, recall e F1 ≥ 95% por classe** (ERROR, WARNING, INFO), coerente com **DN-006-2 B** e **DN-007-2 B**. (fonte: respostas-joao.md, DN-006-2 B / DN-007-2 B; origem: analyze-006-etapa3f.md H2, LLT-199)
- Q: DN-006-3-ref Qual perfil de medição para os **2 s** de **SC-002** (**DN-006-3 B**)? → A: **VM-A real** com **todos os serviços** rodando; lote de **1.000 registros**. (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40); origem: DN-006-3 B)
- Q: DN-008-8 A Qual o separador da assinatura serializada (**008 assinatura A**)? → A: **`|`** (decisão e implementação canônica na spec **008**; a 006 referencia). (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40); origem: 008 assinatura A)
- Q: 007 H2 D Como comparar predição à base rotulada quando o gold é por apontamento? → A: Classe prevista do apontamento = **mais grave** entre as classes previstas dos logs (**ERROR > WARNING > INFO**), comparada ao gold **por apontamento**. (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40); origem: ADR_170, respostas-etapa3f.md opção D)

### Session 2026-11-16 (etapa 3h)

- Q: DLQ B — Existe DLQ no fluxo `log.ingested` quando a consulta OpenObserve falha? → A: **Não** — o contrato `log.ingested` v1 não herda retry/DLQ do guia (**DN-006-1b B**). Após **3 tentativas** na consulta OpenObserve, o Finder avisa o **Core** pela **API interna** com **`batchId`** e **motivo**; **sem** DLQ no ingested e **sem** alerta WebSocket pós-DLQ. **Rev. DN-008-3** / **DN-008-6** substituem o alcance anterior (3 tentativas + DLQ / alerta Core após DLQ — etapas 3b–3e). Rota, payload e evento da API Core: **pendência** (**ADR_051 [PENDENTE]** — MUST NOT inventar). (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B; origem: respostas-etapa3g.md §1b, contrato-log-ingested-v1.txt §8)
- Q: DN-008-9 A — Como tratar o caractere `|` literal dentro de um componente do `valor` da assinatura? → A: Escapar como **`\|`** (barra invertida + barra vertical) **antes** da concatenação com separador **`|`** (**DN-008-8 A**); canônico na spec **008**. Escape de **`\`** literal nos componentes: **pendência** (sem fonte — MUST NOT decidir). (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DN-008-9 A)
- Q: F1 (analyze-006-etapa3g) — A nota histórica do checklist ainda diz que o Finder grava no MongoDB? → A: Redação corrigida em `checklists/requirements.md` (iteração 2) para alinhar a **T1 REVISTO** — Finder envia `ResumoAnalise`; Módulo de Ocorrência cria/persiste. (fonte: respostas-joao.md, T1 REVISTO; origem: analyze-006-etapa3g.md F1)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Classificar o log capturado sem serviço de linguagem externo (Priority: P1)

O Finder parte do log que a captura já deixou disponível (fluxo real: fila → OpenObserve, contrato `log.ingested` v1). O apontamento identifica o conjunto e não é o conteúdo. O modelo lê o campo de nível e o texto da mensagem (já mascarado quando aplicável — fora desta feature). A classe do log segue **FR-004** (nível + texto, palavra inteira, mais grave vence): ERROR vence WARNING, e WARNING vence INFO. O texto não é enviado a um serviço externo de linguagem. A versão única das regras usadas fica no `AchadoAnalise`/`ResumoAnalise`.

**Why this priority**: Sem a classe atribuída no próprio Finder, não há achado classificável para o Módulo Ocorrência. É o resultado mínimo do pedido e da feature LLT-133.

**Independent Test**: Entregar um log já capturado com nível de erro e conferir que a classe é ERROR, que nenhum texto desse log é enviado a um serviço externo de linguagem e que o achado informa a versão única das regras.

**Acceptance Scenarios**:

1. **Given** um log já capturado cujo campo de nível é erro, error ou fatal, e a mensagem não tem termo de nível, **When** a classificação termina, **Then** a classe é ERROR. O texto do log não é enviado a um serviço externo de linguagem. O achado informa a versão das regras.
2. **Given** um log já capturado cujo campo de nível é alerta, warn ou warning, e a mensagem não tem termo de nível mais grave, **When** a classificação termina, **Then** a classe é WARNING.
3. **Given** um log já capturado cujo campo de nível é informação, info, debug ou trace, e a mensagem não tem termo de nível mais grave, **When** a classificação termina, **Then** a classe é INFO.
4. **Given** o campo de nível é informação e a mensagem contém a palavra inteira error, **When** a classificação termina, **Then** a classe do log é ERROR. O catálogo não altera essa classe.
5. **Given** o campo de nível é fatal e a mensagem contém a palavra inteira info, **When** a classificação termina, **Then** a classe do log é ERROR.
6. **Given** o campo de nível é informação e a mensagem contém errorless, **When** a classificação termina, **Then** a classe do log é INFO.
7. **Given** o campo de nível está ausente e a mensagem contém traceback, **When** a classificação termina, **Then** traceback não conta como INFO. Sem outro termo de nível e sem catálogo, a classe é INFO.
8. **Given** a fila só traz o apontamento do conjunto, **When** a classificação é pedida, **Then** a classe é calculada sobre os logs capturados. O apontamento não é tratado como o conteúdo.
9. **Given** o apontamento não tem nenhum log capturado disponível, **When** a classificação é pedida, **Then** nenhum achado é incluído no resumo de sucesso e a falha distingue log ausente de armazenamento inacessível.

---

### User Story 2 - Produzir achados e enviar ResumoAnalise ao Módulo Ocorrência (Priority: P1)

O Finder agrupa logs classificados em **achados** (mesma assinatura no produto, escopo da classificação/agrupamento interno ao lote). A assinatura nesta feature é **nível + mensagem normalizada + `service.name` + host** (DN-008-1 B); o **valor serializado** da assinatura segue **008 assinatura A** (texto concatenado com separador **`|`** — **DN-008-8 A**; `|` literal em componente → **`\|`** — **DN-008-9 A**; canônico na spec **008**). A **mensagem normalizada** segue o algoritmo **DN-008-4 B** com marcadores **DN-008-5 A** (`<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` — implementação canônica na spec **008**, `internal/application/assinatura/normalizar_mensagem.go` — ver plan/tasks da 008); a 006 **só reutiliza**, sem segunda implementação (fonte: respostas-joao.md §Etapa 3b/3c, DN-008-4 B / DN-008-5 A). Para cada achado há um `AchadoAnalise` com `tipo` ERROR, WARNING ou INFO — **inclusive INFO** (T7 A). Antes de montar o envio, o Finder consulta o Módulo de Ocorrência **uma vez por achado** (**DN-008-7 B**, alinhado a N2 B na spec **008**). Ao fim da análise do lote, o Finder monta **um único** `ResumoAnalise` por lote (contadores ADR_102 com contador de **INFO** no lugar de «spam ignorado» + lista de achados) e envia ao Módulo por **HTTP interno com JWT** (fonte: respostas-etapa3c.md §6; origem: ADR_062, Finder.mmd; rota, claims, códigos HTTP e timeout: **pendência no plan** até contrato do Módulo — **DN-009-2**). Chave de idempotência: **`batchId` obrigatório**; **`execId` quando houver** (Varredura); no fluxo atual só `log.ingested` v1, **somente `batchId`** (DN-006-1, DN-006-1b B, N4 A). O Módulo confirma o recebimento; o Finder só então conclui a varredura do arquivo conforme sequência do Drive. O Finder **não** grava ocorrência no MongoDB, **não** calcula gravidade e **não** publica evento para Core/Notificação (o Módulo publica aviso com id + tipo + gravidade — 009-5 revisto). Versão única regras+catálogo+base (006-3 C) vai no achado e/ou no resumo. Desempenho: o ciclo **classificar + montar resumo** de um lote MUST completar em **≤ 2 s** (**DN-006-3 B**, **SC-002**), medido na **VM-A real** com **todos os serviços** rodando e lote de **1.000 registros** (**DN-006-3-ref**).

**Why this priority**: ADR_062 e T1 revisto definem este como o contrato de saída da classificação.

**Independent Test**: Entregar um lote com três logs de mesma assinatura e classe ERROR e conferir um único `AchadoAnalise` ERROR no `ResumoAnalise`, envio HTTP com JWT, confirmação do módulo, e ausência de gravação de ocorrência pelo Finder.

**Acceptance Scenarios**:

1. **Given** três logs capturados com a mesma assinatura e classe ERROR, **When** o envio termina, **Then** o `ResumoAnalise` contém um `AchadoAnalise` com `tipo` ERROR, versão das regras, tenant e produto, sem persistência de ocorrência pelo Finder.
2. **Given** logs com assinaturas ou classes diferentes, **When** o envio termina, **Then** há um `AchadoAnalise` por achado (T2 C), cada um com sua classe.
3. **Given** apenas logs INFO com assinaturas distintas, **When** o envio termina, **Then** cada achado INFO gera entrada no resumo (006-11 A, T7 A).
4. **Given** o Módulo Ocorrência confirma o recebimento, **When** o fluxo termina, **Then** o Finder pode marcar a varredura do arquivo como concluída; **Given** falha de envio antes da confirmação, **Then** reenvio com a **mesma chave de idempotência** (`batchId` ou `batchId`+`execId`, conforme montada) (comportamento FA-02 do Drive, sem duplicar efeito no módulo).
5. **Given** o lote está incompleto (T4 A), **When** o processamento termina, **Then** o resultado **não** é tratado como sucesso de lote completo (sem `ResumoAnalise` final de sucesso equivalente a lote fechado).
6. **Given** envio concluído, **When** o resultado é conferido, **Then** não há evento publicado pelo Finder, sugestão gerada nem notificação disparada por esta feature.
7. **Given** um lote com vários achados distintos, **When** o handoff é preparado, **Then** o Finder consulta o Módulo **uma vez por achado** (**DN-008-7 B**) e envia **um único** `ResumoAnalise` por lote (ADR_062).

---

### User Story 3 - Usar o catálogo só quando o nível não decide (Priority: P2)

O modelo considera o campo de nível e os termos de nível no texto da mensagem. A classe do log é a mais grave entre esses dois sinais. O catálogo de expressões e palavras-chave só entra quando nem o campo nem o texto trazem termo de nível reconhecível. O catálogo **não** é criado por esta feature; sua sintaxe e dono estão pendentes (006-7 D). Sem esses sinais e sem correspondência no catálogo, a classe é INFO e o achado correspondente ainda entra no `ResumoAnalise`.

**Why this priority**: ADR_167 escolheu regras e padrões para o MVP; o catálogo concreto permanece em aberto.

**Independent Test**: Repetir classificação com nível ERROR e catálogo apontando INFO, com nível ausente e catálogo apontando WARNING, e sem nível nem catálogo; conferir classes e que nenhum catálogo foi criado.

**Acceptance Scenarios**:

1. **Given** o campo de nível é error e o catálogo fornecido apontaria INFO, **When** a classificação termina, **Then** a classe é ERROR e o achado correspondente usa `tipo` ERROR.
2. **Given** o campo é INFO, a mensagem contém error e o catálogo apontaria INFO, **When** a classificação termina, **Then** a classe do log é ERROR; o catálogo não entra.
3. **Given** nem o campo nem a mensagem trazem termo de nível e a primeira entrada correspondente do catálogo aponta WARNING, **When** a classificação termina, **Then** a classe é WARNING.
4. **Given** não há catálogo e nem o campo nem a mensagem trazem termo de nível, **When** a classificação termina, **Then** a classe é INFO, nenhum catálogo é criado e o achado entra no `ResumoAnalise`.
5. **Given** várias entradas do catálogo correspondem e não há nível reconhecível, **When** a classificação termina, **Then** vale a primeira entrada correspondente, na ordem fornecida.

---

### Edge Cases

- Apontamento sem nenhum log capturado: nenhum achado no resumo de sucesso; falha distingue log ausente.
- **Lote incompleto (T4 A):** não é sucesso; distinguir de falha de envio ao Módulo Ocorrência.
- Armazenamento do log inacessível: nenhum `ResumoAnalise` de sucesso; falha distinta.
- Falha de envio HTTP ao Módulo Ocorrência: arquivo não marcado como varrido até confirmação; reenvio com mesma chave (Drive, ADR_014).
- Campo de nível vazio/ausente e mensagem sem termo: catálogo se fornecido; senão INFO; achado incluído no resumo (T7 A).
- Palavra inteira, múltiplos logs, catálogo inválido: regras da sessão 2026-11-09 permanecem.
- **Reenvio com a mesma chave** (`batchId` ou `batchId`+`execId`): Módulo Ocorrência ignora duplicata (006-9 A no módulo, não no Finder).
- **Sem `execId` no pipeline atual:** idempotência usa só `batchId` (DN-006-1b B); quando a Varredura passar `execId`, compor `batchId`+`execId` (DN-006-1).
- **Log sem `service.name` ou host** para assinatura (DN-008-1 B): falha **`assinatura_incompleta`**; não sucesso de lote completo (paralelo a `identidade_incompleta` — 007-8 A).
- **`productId` ausente** no evento `log.ingested` v1 ou **inconsistente** entre metadata validada (token) e corpo repetido para agrupamento: **falha crítica permanente**, **sem retry** — motivo distinguível **`produto_ausente_ou_inconsistente`** (**008-16 A**, ADR_049, contrato §4); não sucesso de lote; distinto de `identidade_incompleta` (007-8 A: tenant/produto vazio no achado).
- Texto com segredo: trechos mascarados vão em `AchadoAnalise.linhas` conforme ADR_040; definição de mascaramento é spec 010.
- Consumo de `log.ingested`: rejeitar ou não processar eventos com `schemaVersion` diferente de `"1"` (T5 D).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: O Finder MUST classificar cada log já capturado de um apontamento. O apontamento MUST identificar esse conjunto e MUST NOT ser tratado como o conteúdo. O processamento MUST alinhar-se ao fluxo real fila/OpenObserve e ao contrato `log.ingested` **v1** (`schemaVersion` "1") (T5 D).
- **FR-002**: A classificação MUST ocorrer sem enviar o texto do log a um serviço externo de linguagem.
- **FR-003**: A classe de cada log MUST ser ERROR, WARNING ou INFO. Spam e "não classificado" MUST NOT ser classes. Cada log classificado MUST contribuir para um **achado**; o grão de ocorrência no produto é **um por achado** (T2 C), não um por apontamento.
- **FR-004**: *(inalterado — mapear erro/alerta/informação para ERROR/WARNING/INFO.)* O campo de nível e o texto da mensagem MUST contar juntos; palavra inteira; a classe do log é a mais grave entre campo e texto.
- **FR-005**: O catálogo MUST ser usado somente quando o nível não é reconhecível. Esta feature MUST NOT criar o catálogo. Sintaxe e dono: **pendência** 006-7 D (LLT-51 ou ADR nova).
- **FR-006**: Sem termo de nível reconhecível e sem catálogo válido, a classe MUST ser INFO; o achado MUST ser incluído no `ResumoAnalise` quando o lote for tratado como sucesso.
- **FR-007**: O Finder MUST produzir `AchadoAnalise` por achado com: `tipo` (ERROR/WARNING/INFO), identidade do achado (**assinatura** = nível + mensagem normalizada + `service.name` + host — DN-008-1 B; **valor** serializado conforme **008 assinatura A**, separador **`|`** — **DN-008-8 A**, canônico na spec **008**), tenant, produto, versão **única** regras+catálogo+base (006-3 C), e demais campos exigidos pelo contrato Drive (revisao-T1-drive.md §2, ADR_040) e extensões (assinatura, contagem, primeira/última no lote — N3 A, pendência de ADR). Para logs contáveis, `service.name` e **host/máquina** MUST estar presentes para compor assinatura; caso contrário MUST registrar falha **`assinatura_incompleta`** (sem valor sentinela inventado). Antes do envio do lote, o Finder MUST consultar o Módulo de Ocorrência **uma vez por achado** (**DN-008-7 B**; contrato/rota: **DN-009-2**). No escopo mínimo desta feature (**FR-009**, sem agrupamento **entre lotes**), `recorrente` MUST ser **`false`**; recorrência histórica via OpenObserve permanece na spec **008** (008-3/5/6; **DN-008-3** rev. 3h e **DN-008-6** rev. 3h — até **3 tentativas**, **sem DLQ** no `log.ingested`; após esgotar, Finder avisa o Core pela API interna com **batchId** e **motivo** [**ADR_051 [PENDENTE]**]; MUST NOT assumir `recorrente=false` por falha transitória de consulta OpenObserve nem usar DLQ/WebSocket pós-DLQ do guia da FilaOrdenada). Campos Drive `titulo`, `resumo`, `numeroLinha`, `camposExtraido` MUST constar no payload; **semântica de preenchimento mínimo** permanece **pendente N3 A** até ADR. `hostName` MUST alinhar 1:1 ao componente `host` da assinatura quando este estiver disponível (ADR_040). O Finder MUST NOT gravar ocorrência no MongoDB. MUST montar **exatamente um** `ResumoAnalise` **por lote** (ADR_062) com contadores alinhados à ADR_102, usando contador de **INFO** em lugar de «spam ignorado» (etapa 3b, ADR_170), e enviá-lo ao Módulo por HTTP interno com JWT (N1 A). Rota, claims JWT, códigos de resposta e política 5xx/timeout do Módulo **não** estão fechados nesta spec — **pendência no plan** até contrato do Módulo (respostas-etapa3.md §10). MUST NOT incluir gravidade nem criticidade (ADR_062); o Módulo calcula gravidade no aviso ao Core (009-5 revisto: id + tipo + gravidade). Identificação das regras MUST indicar se o catálogo entrou e qual versão única foi usada. Ciclo **classificar + montar resumo** do lote: ver **SC-002** (**DN-006-3 B**).
- **FR-008**: Se não há log presente, se o armazenamento é inacessível, se o envio ao Módulo Ocorrência não é confirmado, se algum log contável carece de `service.name` ou host exigidos pela assinatura (**`assinatura_incompleta`**), ou se o **`productId`** está **ausente** ou **inconsistente** entre metadata validada e corpo do `log.ingested` v1 (**`produto_ausente_ou_inconsistente`** — **008-16 A**, ADR_049, contrato §4), o Finder MUST NOT apresentar o lote como varredura concluída com sucesso. Falhas **008-16 A** MUST ser tratadas como **críticas permanentes** — **MUST NOT** aplicar retry de processamento do lote (padrão ADR_049). Motivos MUST ser distinguíveis (incluindo `identidade_incompleta` — 007-8 A). **Lote incompleto MUST NOT contar como sucesso** (T4 A). Logs presentes em lote incompleto: classificação pode ocorrer, mas não satisfaz sucesso de lote.
- **FR-009**: Esta feature MUST NOT agrupar **entre lotes** (LLT-200 / 008 — consulta OpenObserve entre lotes fora do mínimo desta FR), MUST NOT publicar evento de ocorrência, MUST NOT calcular gravidade, MUST NOT gerar sugestão, MUST NOT notificar, MUST NOT definir catálogo de mascaramento. Agrupamento **dentro do lote** por assinatura para formar achado é escopo da 006 em conjunto com 008-18 A. Consulta ao Módulo **por achado** (**DN-008-7 B**, N2 B no handoff) **é** escopo de **FR-007**, distinta de agrupamento entre lotes.
- **FR-010**: A classe de um log MUST ser calculável log a log. **Upsert por achado (006-9 A) não é responsabilidade do Finder** após T1 revisto. O Módulo Ocorrência cria/atualiza ocorrências sem duplicar por chave. Métricas de qualidade para aceite LLT-197/199: precisão, recall e F1 **por classe** (006-2 B) com base rotulada da spec 007 (**007 formato A** JSONL, **007 caminho C** `datasets/labeled-logs/` na raiz do repositório — só documentado, pasta não criada nesta feature). Na avaliação offline, o gold é **por apontamento** (spec 007, **DN-007-4 A**); a **classe prevista do apontamento** para comparação MUST ser a **mais grave** entre as classes previstas dos logs do apontamento (**ERROR > WARNING > INFO** — **007 H2 D**; origem: ADR_170; detalhes: respostas-etapa3f.md opção D). **Gate formal de aceite nesta spec** na partição de validação (**007-meta**, SC-007, T047): em **cada** classe ERROR, WARNING e INFO — **precisão ≥ 95%** (**DN-007-2 B**), **recall ≥ 95%** e **F1 ≥ 95%** (**DN-006-2 B**, etapa 3e). Relatório MUST documentar as três métricas por classe mesmo fora do gate offline.
- **FR-011**: A chave de idempotência do `ResumoAnalise` MUST usar **`batchId`** do `log.ingested` v1 (`idempotencyId` / `domain.Referencia.BatchID`). **`execId` MUST ser incluído quando disponível** (origem Varredura: `V<uuid>` ou `M<uuid>` — DN-006-1); **enquanto o Finder consume apenas `log.ingested` v1 sem canal que entregue `execId`**, a chave MUST ser **somente `batchId`** (DN-006-1b B). Alinha a N4 A («batchId + execId, se houver») e ADR_146 («execId quando houver»).
- **FR-012**: Ao mudar a versão única das regras (006-3 C), MUST aplicar-se **somente a novos logs** (006-4 A), sem reprocessamento automático retroativo definido nesta feature.
- **FR-013**: Ao serializar o **valor** da assinatura (**008 assinatura A**, **DN-008-8 A**), cada componente MUST escapar **`|`** literal como **`\|`** antes da concatenação (**DN-008-9 A**; canônico na spec **008**); escape de **`\`** literal nos componentes permanece **pendente** (sem fonte). Ao compor a assinatura (FR-007, DN-008-1 B), o Finder MUST aplicar o normalizador de mensagem **DN-008-4 B** pela implementação **canônica** da spec **008** — arquivo de referência `internal/application/assinatura/normalizar_mensagem.go` (definido em `specs/008-agrupar-logs-ocorrencia/plan.md` e **tasks.md** T016): aparar espaços; comparação case-insensitive; substituir número, UUID, IP, data e hexadecimal pelos marcadores **DN-008-5 A**: `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>`. A spec **006** MUST **reutilizar** essa implementação (delegação/import compartilhado) e **MUST NOT** manter algoritmo paralelo ou cópia divergente (ex. em `internal/domain/normalizar_mensagem.go`). MUST NOT publicar regras de normalização divergentes. Casos de borda adicionais: **010-27 A** (spec 010). (fonte: respostas-joao.md §Etapa 3b/3c, DN-008-4 B / DN-008-5 A; origem: ADR_013)

### Key Entities

- **Log capturado**: Entrada da classificação; recuperado via apontamento/OpenObserve.
- **Nível reconhecível / FR-004**: Como na spec de 2026-11-09.
- **Catálogo de classificação**: Opcional, fornecido na execução; não criado aqui; sintaxe pendente (006-7 D).
- **Achado**: Conjunto de logs com a mesma assinatura no produto no escopo do lote (008-18 A); **assinatura** = nível + mensagem normalizada + `service.name` + host (DN-008-1 B); **valor** agregado = texto concatenado (**008 assinatura A**); separador **`|`** (**DN-008-8 A**); escape **`|`** → **`\|`** nos componentes (**DN-008-9 A**; canônico na spec **008**). Unidade que gera um `AchadoAnalise`, uma consulta prévia ao Módulo (**DN-008-7 B**) e, no Módulo, uma ocorrência (T2 C).
- **Mensagem normalizada**: Saída do algoritmo **DN-008-4 B** sobre o texto do log (após mascaramento aplicável — spec 010); **única** implementação canônica na spec **008** (`internal/application/assinatura/normalizar_mensagem.go`); a 006 só invoca essa implementação; marcadores **DN-008-5 A**: `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>`.
- **AchadoAnalise**: Estrutura enviada pelo Finder; inclui `tipo` ERROR/WARNING/INFO (006-6 A), versão única das regras (006-3 C), campos do Drive (ADR_040, revisao-T1-drive.md §2) e extensões decididas (assinatura, contagem, primeira/última, versão — N3 A, contrato formal pendente de ADR). No escopo mínimo 006, `recorrente` inicia **`false`**; `hostName` espelha `host` da assinatura quando disponível; derivação de `titulo`/`resumo` pendente **N3 A**.
- **ResumoAnalise**: Contadores do lote (ADR_102: linhas processadas, erros, warnings, **INFO**, dados sensíveis mascarados — sem contador «spam ignorado») + lista de `AchadoAnalise`; **um envio por lote** ao Módulo Ocorrência (ADR_062); chave `batchId` ou `batchId`+`execId` conforme FR-011.
- **execId**: Identificador da execução de varredura gerado na Varredura (`V<uuid>` / `M<uuid>`), opcional no handoff; ausente no corpo `log.ingested` v1 (DN-006-1).
- **Ocorrência**: Entidade do Módulo Ocorrência, **criada fora do Finder** (T1 REVISTO). Chave de unicidade no módulo: assinatura + classe + produto + tenant (008-11 revisado 16/11), coerente com uma ocorrência por (assinatura, classe) (008-15 A).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Em classificação concluída, cada log tem uma única classe ERROR, WARNING ou INFO; versão única das regras identificável no achado/resumo; zero envios de texto a LLM.
- **SC-002**: Após envio confirmado pelo Módulo Ocorrência, o observador verifica `ResumoAnalise` com achados corretos por assinatura/classe (DN-008-1 B), contadores coerentes (incluindo INFO, sem spam ignorado), com tenant, produto e versão, sem gravação de ocorrência pelo Finder. **Desempenho (normativo):** o ciclo **classificar todos os logs do lote + montar o `ResumoAnalise`** MUST completar em **≤ 2 s por lote** (**DN-006-3 B**). **Perfil de medição (**DN-006-3-ref**):** **VM-A real** com **todos os serviços** em execução; lote de **1.000 registros**. (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40); origem: DN-006-3 B). Não inclui latência de consultas ao Módulo (**DN-008-7 B**) nem handoff HTTP até contrato **DN-009-2**.
- **SC-003**: Ausência de log, armazenamento inacessível, lote incompleto (T4 A), **`assinatura_incompleta`**, **`produto_ausente_ou_inconsistente`** (**008-16 A**, sem retry) ou falha de envio: nenhum sucesso de lote completo; motivos distinguíveis.
- **SC-004**: Sem catálogo criado pela feature; nível error → ERROR; sem nível → INFO; achados ERROR no resumo quando aplicável.
- **SC-005**: Cenários de precedência nível/texto/catálogo da spec original permanecem válidos com rótulos ERROR/WARNING/INFO.
- **SC-006**: Zero eventos de ocorrência publicados pelo Finder nesta feature.
- **SC-007**: Métricas de avaliação offline (fora do runtime de produção): precisão, recall e F1 **por classe** documentados (006-2 B) com base rotulada da spec 007; comparação ao gold **por apontamento** usando classe prevista do apontamento = **mais grave** entre os logs (**007 H2 D**, ADR_170). **Gate formal nesta spec** na partição de validação (**007-meta**; procedimento e asserts em T047): em **cada** classe ERROR, WARNING e INFO — **precisão ≥ 95%** (**DN-007-2 B**), **recall ≥ 95%** e **F1 ≥ 95%** (**DN-006-2 B**, etapa 3e).

## Pendências explícitas (não executar nesta etapa)

| ID | Descrição | Decisão do João | Ação bloqueada nesta execução |
|----|-----------|-----------------|-------------------------------|
| 006-7 D | Sintaxe e dono do catálogo ADR_167 | Definir na LLT-51 ou ADR nova | Não criar ADR |
| 006-10 A | Alinhar status Jira LLT-197/LLT-199 | Sim, com a 006 | Não alterar Jira |
| N3 A | Formalizar ampliação de `AchadoAnalise` em ADR (incl. semântica `titulo`/`resumo`) | Sim, com ADR | Não criar ADR |
| T5 D | Commit do Coletor publicando `schemaVersion` "1" | v1 vale no contrato | Fora do repositório Finder |
| API Módulo (etapa 3) / **DN-009-2** | Rota, claims JWT, códigos HTTP, 5xx/timeout do handoff Finder → Módulo Ocorrência | Pendência no **plan** até contrato do Módulo | SC-002 E2E real bloqueado |
| **Emenda princípio I** | Constituição I («Finder agrega em ocorrência») × **T1 REVISTO** (Módulo Ocorrência cria ocorrência) | Governança — proposta em **plan.md** / **T043** | Não editar `.specify/memory/constitution.md` nesta pasta |
| **010-27 A** | Casos de teste adicionais do normalizador | João define na 010 | Complementa **DN-008-5 A** em T052 |
| **Contrato API Core (alertas)** | Rota, payload e evento Finder→Core após 3 falhas na consulta OpenObserve (**batchId** + **motivo**) | **ADR_051 [PENDENTE]** (fonte: respostas-joao.md §Etapa 3g, DLQ B) | Implementação na spec **008**; a 006 só referencia **DN-008-3** / **DN-008-6** rev. 3h |
| **Escape `\` na assinatura** | Regra para barra invertida literal nos componentes do `valor` | Sem fonte | Complementa **DN-008-9 A** (só **`|`** fechado) |

## Fora de escopo

- Gravar ou consultar ocorrência no MongoDB pelo Finder (substituído por envio de `ResumoAnalise`).
- Publicar evento de ocorrência ao Core (Módulo Ocorrência, 009).
- Calcular gravidade, sugestão, notificação, mascaramento (specs 009, 010).
- Criar catálogo de classificação ou ADR de contrato AchadoAnalise.
- Alinhar issues Jira LLT-197/199.
- Agrupamento **entre lotes** e upsert de ocorrência (Módulo Ocorrência / 008; 006-9 A sem efeito no Finder).

## Assumptions

- Decisões 2026-11-16 prevalecem sobre 2026-11-09 onde indicado na sessão de clarificações.
- Fluxo normativo: Finder → HTTP(JWT) → Módulo Ocorrência → cria ocorrências (ADR_062).
- Handoff: **um** `ResumoAnalise` por lote (ADR_062); consulta ao Módulo **por achado** antes do envio (**DN-008-7 B**, N2 B). Recorrência histórica OpenObserve: spec **008** (fora do mínimo **FR-009**).
- Constituição: classes erro/alerta/informação; princípio V — sem invenção; lacunas bloqueiam plano (007-26 A aplicável ao conjunto de features).

## Decisões tomadas (com fonte)

| ID | Decisão | Arquivo-fonte | Origem documental |
|----|---------|---------------|-------------------|
| T1 REVISTO | Finder não grava ocorrência; envia `ResumoAnalise`; Módulo Ocorrência cria ocorrência | respostas-joao.md | ADR_062; revisao-T1-drive.md §2 |
| N1 A | HTTP interno com JWT Finder → Módulo Ocorrência | respostas-joao.md | ADR_051; 009-3 D |
| T2 C | Uma ocorrência por achado (no módulo), um `AchadoAnalise` por achado no Finder | respostas-joao.md | ADR_040, ADR_035 |
| T3 B | Recorrência por consulta OpenObserve (não detalha 006) | respostas-joao.md | ADR_168 |
| T4 A | Lote incompleto não é sucesso | respostas-joao.md | spec 005 |
| T5 D | Somente `log.ingested` v1 (`schemaVersion` "1") | respostas-joao.md | contrato-log-ingested-v1.txt |
| T6 A | Lacunas LLT-197/199 na spec 006 | respostas-joao.md | LLT-197, LLT-199 |
| T7 A | INFO gera ocorrência (achado INFO no resumo) | respostas-joao.md | ADR_170; diverge Drive — ver Divergências |
| 006-2 B | Métricas: precisão/recall/F1 por classe | respostas-joao.md | LLT-197, LLT-199 |
| 006-3 C | Versão única regras+catálogo+base no achado/resumo | respostas-joao.md | LLT-199 crit. 3 |
| 006-4 A | Só novos logs ao mudar versão | respostas-joao.md | LLT-51 |
| 006-6 A | `tipo` = ERROR/WARNING/INFO; sem gravidade no Finder | respostas-joao.md | ADR_062 |
| 006-7 D | Catálogo: LLT-51 ou ADR nova (**pendência**) | respostas-joao.md | ADR_167 |
| 006-9 A | Upsert por achado: **sem efeito no Finder** | respostas-joao.md; revisao-T1-drive.md A.1 | ADR_040 |
| 006-10 A | Alinhar Jira (**pendência**, não feito aqui) | respostas-joao.md | LLT-197, LLT-199 |
| 006-11 A | INFO por achado, como ERROR/WARNING | respostas-joao.md | T2 C + T7 A |
| N3 A | Ampliar AchadoAnalise (ADR **pendente**, não criada) | respostas-joao.md | ADR_040; Finder.mmd |
| N4 A | Chave ResumoAnalise: `batchId` + `execId` **se houver**; só `batchId` enquanto sem `execId` (DN-006-1b B) | respostas-joao.md; respostas-etapa3.md §3 | contrato v1; ADR_070; ADR_146 |
| DN-006-1 | `execId` da Varredura (V/M + uuid); **opcional** | respostas-etapa3.md §3 | ADR_075, ADR_069, ADR_070, ADR_146 |
| DN-006-1b B | Sem `execId` no fluxo só `log.ingested`; chave só `batchId` | respostas-joao.md, Etapa 3 | contrato-log-ingested-v1.txt |
| DN-008-1 B | Assinatura: nível + mensagem normalizada + `service.name` + host | respostas-joao.md, Etapa 3 | ADR_013 |
| Contadores etapa 3 | «spam ignorado» → contador de **INFO** | respostas-joao.md, Etapa 3 | ADR_102, ADR_170 |
| 009-5 revisto | Aviso ao Core: id + tipo + gravidade (Módulo Ocorrência) | respostas-joao.md, Etapa 3 | sequencia-notificacao.md v2 |
| 008-11 (16/11) | Chave ocorrência: assinatura + classe + produto + tenant | respostas-joao.md rodada 00h01 | reconcilia 008-15 A |
| Taxonomia 14/11 | ERROR/WARNING/INFO (não ERRO/WARNING/SPAM ADR_019) | respostas-joao.md 007-9 A; 006-6 A | ADR_170 |
| DN-007-2 B | Meta **007-meta**: precisão **≥ 95% por classe** (ERROR, WARNING, INFO); gate SC-007 / T047 | respostas-joao.md §Etapa 3b | respostas-etapa3.md §1; LLT-199 |
| DN-008-4 B | Normalização da mensagem na assinatura: trim, case-insensitive, marcadores número/UUID/IP/data/hex; dono spec **008**, impl. `internal/application/assinatura/normalizar_mensagem.go`; 006 só reutiliza | respostas-joao.md §Etapa 3b | DN-008-1 B; ADR_013; specs/008 plan/tasks T016 |
| 008-16 A | `productId` ausente ou metadata≠body: falha crítica permanente, sem retry; motivo `produto_ausente_ou_inconsistente` | respostas-joao.md, 008-16 A | contrato-log-ingested-v1.txt §4; ADR_049 |
| DN-006-2 B | Gate recall **≥ 95%** e F1 **≥ 95%** por classe (além precisão **DN-007-2 B**) | respostas-joao.md §Etapa 3c, 11h e 11h01 | respostas-etapa3c.md §1; LLT-199; 006-2 B |
| DN-006-3 B | Meta **≤ 2 s** por lote (classificar + montar resumo) em **SC-002** | respostas-joao.md §Etapa 3c (João, 16/11 11h02) | respostas-etapa3c.md §2 opção B |
| DN-008-5 A | Marcadores `<NUM>` `<UUID>` `<IP>` `<DATE>` `<HEX>` | respostas-joao.md §Etapa 3c, 11h e 11h01 | respostas-etapa3c.md §4; impl. canônica spec 008 |
| DN-008-7 B | Consulta ao Módulo **por achado**; envio **um** `ResumoAnalise` por lote | respostas-joao.md §Etapa 3c (João, 16/11 11h02) | respostas-etapa3c.md §6; ADR_062; Finder.mmd |
| 008 assinatura A | Valor da assinatura = texto concatenado com separador | respostas-joao.md §Etapa 3c, 11h e 11h01 | respostas-etapa3c.md §9 |
| DN-006-3-ref | Medição **SC-002**: VM-A real, todos os serviços, lote **1.000** registros | respostas-joao.md §Etapa 3e (João, 16/11 12h40) | DN-006-3 B |
| DN-008-8 A | Separador da assinatura: **`|`** | respostas-joao.md §Etapa 3e (João, 16/11 12h40) | spec **008** (dona) |
| 007 H2 D | Avaliação: classe prevista do apontamento = mais grave entre logs vs gold por apontamento | respostas-joao.md §Etapa 3e (João, 16/11 12h40) | ADR_170; respostas-etapa3f.md opção D |
| **DLQ B** | **Não** há DLQ no fluxo `log.ingested` quando a consulta OpenObserve falha | respostas-joao.md §Etapa 3g (João, 16/11 13h22) | contrato-log-ingested-v1.txt §8; respostas-etapa3g.md §1b; **DN-006-1b B** |
| **DN-008-3** (rev. 3h) | Até **3 tentativas**; após esgotar, Finder→Core API interna (**batchId** + **motivo**); **substitui** alcance 3b–3e (3 tentativas + DLQ) | respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B | spec **008** (runtime); **ADR_051 [PENDENTE]** |
| **DN-008-6** (rev. 3h) | Mesmo gatilho que **DN-008-3** rev. 3h; **não** DLQ nem alerta WebSocket pós-DLQ; **substitui** etapa 3e | respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B | **ADR_051 [PENDENTE]** |
| **DN-008-9 A** | **`|`** literal em componente → **`\|`** antes da concatenação; escape de **`\`**: pendência | respostas-joao.md §Etapa 3g (João, 16/11 13h22), DN-008-9 A | spec **008** (serialização canônica) |

### Decisões tomadas (etapa 3h)

| ID | Decisão | Fonte | Nota |
|----|---------|-------|------|
| **DLQ B** | Sem DLQ no `log.ingested`; após 3 falhas OpenObserve, Finder avisa o Core (API interna, **batchId** + **motivo**) | (fonte: respostas-joao.md, DLQ B; origem: §Etapa 3g (João, 16/11 13h22); respostas-etapa3g.md §1b) | Contrato API Core **ADR_051 [PENDENTE]** |
| **DN-008-3** (rev. 3h) | **Substitui** alcance anterior (3 tentativas + DLQ / stream ingested) | (fonte: respostas-joao.md, DLQ B; origem: §Etapa 3g) | Implementação principal na spec **008** |
| **DN-008-6** (rev. 3h) | **Substitui** etapa 3e (DLQ + Core/WebSocket pós-DLQ); mesmo gatilho via API interna | (fonte: respostas-joao.md, DLQ B; origem: §Etapa 3g) | Distinto da publicação de evento de ocorrência (spec **009**) |
| **DN-008-9 A** | Escape **`|`** → **`\|`** nos componentes do `valor` | (fonte: respostas-joao.md, DN-008-9 A; origem: §Etapa 3g (João, 16/11 13h22)) | Escape **`\`**: só pendência |

## Divergências documentadas

### D1 — INFO vira ocorrência (João) × Drive (só ERRO/WARNING)

- **Lado A (vale para o produto):** T7 A e 006-11 A — achados INFO no `ResumoAnalise`; Módulo Ocorrência deve aceitar e criar ocorrência INFO. (fonte: respostas-joao.md, T7 A; origem: decisão consciente 16/11 00h01)
- **Lado B (Drive):** SSID_007.1 critério 1, SSID_009.1 (“Toda Ocorrência enviada é ERRO ou WARNING”), ADR_040 (achados não spam). (origem: revisao-T1-drive.md A.1)
- **Por quê vale A:** respostas do João prevalecem; ADR_170 já adota ERROR/WARNING/INFO e deixa “INFO gera ocorrência” em aberto no Drive — fechado pelo João.

### D2 — T1 A (Finder grava MongoDB) × ADR_062 / arquitetura v2

- **Resolvido por T1 REVISTO (16/11):** vale o fluxo Drive; spec 006 deixa de exigir gravação pelo Finder.

### D3 — ADR_146 v2 “schemaVersion 2” × contrato 08/11 v1

- **Vale:** contrato `log.ingested` **v1**, `schemaVersion` `"1"` (T5 D). ADR_146 v2 não altera esta feature para o evento ingested. (fonte: respostas-joao.md, T5 D; origem: contrato-log-ingested-v1.txt)

### D4 — LLT-201 (Jira) × Drive (emissor e destino do aviso)

- **Vale Drive:** Finder só “avisa Ocorrências” via `ResumoAnalise`; aviso à Notificação sai do Módulo Ocorrência → Core (HTTP JWT 009-3 D), payload **id + tipo + gravidade** (009-5 revisto). Jira não alterado (006-10 A pendente). (origem: revisao-T1-drive.md §1; respostas-joao.md, Etapa 3)

### D5 — Lista da assinatura (DN-008-1 B) × 008-2 A (correlationId/traceId) × contrato v1

- **Vale nesta spec (006):** assinatura = **nível + mensagem normalizada + `service.name` + host** (DN-008-1 B, etapa 3b). **Não** inclui `traceId` nem `correlationId` aqui.
- **Risco residual (008-2 A em outras specs):** contrato 08/11 §4 — `correlationId` aleatório por exportação; incluí-lo na assinatura duplicaria achados entre exportações. João manteve 008-2 A no histórico transversal; **006 não aplica** trace/correlation na assinatura. (fonte: respostas-joao.md, DN-008-1 B; respostas-etapa3.md §4; origem: contrato-log-ingested-v1.txt §4, ADR_063/074 [PENDENTE])

### D8 — `execId` (Varredura/ADR_070) × corpo `log.ingested` v1 (sem `execId`)

- **Lado A (chave):** N4 A + DN-006-1 — quando houver `execId` da Varredura (`V<uuid>`/`M<uuid>`), compor `batchId`+`execId`. (fonte: respostas-etapa3.md §3; origem: ADR_075, ADR_070)
- **Lado B (pipeline atual):** contrato v1 não traz `execId`; João **DN-006-1b B** — chave **só `batchId`** até contrato `log.checked` ou equivalente. (fonte: respostas-joao.md, Etapa 3; origem: contrato-log-ingested-v1.txt)

### D9 — Contador «spam ignorado» (ADR_102) × taxonomia sem SPAM (ADR_170) e INFO como achado (T7 A)

- **Vale:** contador de **INFO** no `ResumoAnalise`, não «spam ignorado». (fonte: respostas-joao.md, Etapa 3; origem: ADR_102 [INFERIDO], ADR_170)
- **Ressalva Drive:** ADR_019, 035, 040, 066 ainda citam `spamIdentificados` em diagramas legados — rótulo desta feature segue João + ADR_170.

### D10 — 009-5 A histórico («só o id») × diagrama Notificação v2 (tipo + gravidade no alerta)

- **Vale (etapa 3b):** **009-5 revisto** — aviso ao Core leva **id + tipo + gravidade** (Módulo Ocorrência); Finder envia `tipo` no achado, sem gravidade. Coerente com destinatários por cargo (ADR_130/136/137/142). (fonte: respostas-joao.md, Etapa 3; origem: respostas-etapa3.md §7)
- **Fora da 006:** origem do texto do e-mail (DN-009-1 A: Notificação lê ocorrência pelo id) — spec 009/010.

### D11 — Handoff HTTP ao Módulo × Drive sem rota/códigos fechados

- **Vale N1 A:** HTTP interno + JWT em rede privada (ADR_028, ADR_141).
- **Lacuna:** rota, claims, códigos e timeout **pendentes no plan** até contrato do Módulo (respostas-etapa3.md §10). Não inventar path nem semântica 202 vs 200 nesta spec.

### D12 — Princípio I da constituição × T1 REVISTO

- **Lado A (constituição vigente, não editada nesta pasta):** o Finder agrega logs capturados **em ocorrência**.
- **Lado B (produto, João):** o Finder envia `ResumoAnalise`; o **Módulo de Ocorrência** cria a ocorrência (T1 REVISTO, ADR_062).
- **Por quê vale B para esta spec:** decisões João prevalecem; conformidade formal exige **emenda do princípio I** (pendência governança — **plan.md** § Constitution Check, **T043**).

### D6 — 008-11 B histórico (sem classe) × 008-15 A e revisão 008-11 (16/11)

- **Vale (16/11):** chave com **assinatura + classe + produto + tenant**, alinhada a uma ocorrência por (assinatura, classe). (fonte: respostas-joao.md, 008-11 rodada 00h01; 008-15 A)

### D7 — 009-12 (sem gravidade no Finder) × SSID_009.1 (com gravidade no aviso)

- **Sem conflito na 006:** Finder não envia gravidade (ADR_062); gravidade no aviso ao Core é calculada e enviada pelo Módulo Ocorrência. (origem: revisao-T1-drive.md A.2, 009-12 revista)

## Dúvidas novas (sem fonte)

- **DN-006-1** — respondida em 2026-11-16 (etapa 3b), ver Clarifications → Session 2026-11-16 (etapa 3b).
- **DN-006-1b** — respondida em 2026-11-16 (etapa 3b), ver Clarifications → Session 2026-11-16 (etapa 3b).
- **DN-007-2** — respondida em 2026-11-16 (etapa 3c), ver Clarifications → Session 2026-11-16 (etapa 3c).
- **DN-008-4** — respondida em 2026-11-16 (etapa 3c), ver Clarifications → Session 2026-11-16 (etapa 3c).
- **DN-006-2** — respondida em 2026-11-16 (etapa 3e), ver Clarifications → Session 2026-11-16 (etapa 3e).
- **DN-006-3** — respondida em 2026-11-16 (etapa 3e) quanto à meta **2 s** (**DN-006-3 B**).
- **DN-006-3-ref** — respondida em 2026-11-16 (etapa 3g), ver Clarifications → Session 2026-11-16 (etapa 3g).
- **DN-008-5** — respondida em 2026-11-16 (etapa 3e), ver Clarifications → Session 2026-11-16 (etapa 3e).
- **DN-008-7** — respondida em 2026-11-16 (etapa 3e) quanto a envio por lote e consulta por achado; ver Clarifications → Session 2026-11-16 (etapa 3e).
- **DN-008-9** — respondida em 2026-11-16 (etapa 3h), ver Clarifications → Session 2026-11-16 (etapa 3h).

*Demais pendências consolidadas na tabela § Pendências explícitas.*