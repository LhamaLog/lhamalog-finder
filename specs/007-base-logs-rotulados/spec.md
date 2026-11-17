# Feature Specification: Base de logs rotulados para treino e avaliação

**Feature Branch**: `007-base-logs-rotulados` (trabalho permanece na branch git `Development`; nenhuma branch dedicada foi criada nesta execução, conforme instrução do solicitante)

**Created**: 2026-11-14

**Status**: Draft (clarificado em 2026-11-16; `plan.md` e `tasks.md` presentes; **pronto para implementação MVP** (US1 + fundação), ressalvas **006-7 D** / subconjunto só-catálogo e pendência transversal **emenda do princípio I**)

**Input**: LLT-198 (Jira); item E07.F2.2 em `LhamaLog_backlog_pendente.docx` (14/10/2026); contexto do especialista em `/workspace/lhamalog/contexto-LLT-103.md` §2.1; respostas do João em `respostas-joao.md` (14–16/11/2026).

**Rastreabilidade**: Epic Finder — IA própria de classificação e agregação; Feature modelo de classificação (E07.F2); complementa LLT-197 (métricas de qualidade), LLT-199 (critério 2 — precisão na base de validação) e a spec `006-modelo-classificacao-logs` (FR-010 exclui meta de precisão não acordada; consumo da base para métricas 006-2 B).

## Clarifications

### Session 2026-11-16

- Q: 007-1 Qual é o tamanho mínimo da base? → A: Milhares (fonte: respostas-joao.md, 007-1; origem: perguntas-etapa2b.md)
- Q: 007-2 Qual é a cobertura mínima por classe? → A: Mínimo igual por classe (fonte: respostas-joao.md, 007-2; origem: perguntas-etapa2b.md)
- Q: 007-3 No MVP por regras, qual o uso da base? → A: Só avaliação das regras (fonte: respostas-joao.md, 007-3; origem: ADR_167)
- Q: 007-4 Onde e em que formato fica o conjunto? → A: Arquivos versionados no repositório (fonte: respostas-joao.md, 007-4; origem: perguntas-etapa2b.md)
- Q: 007-5 Quais campos cada registro MUST ter? → A: Nível + mensagem + `lhamalog_batch_id` / `lhamalog_record_id` (fonte: respostas-joao.md, 007-5; origem: contrato log.ingested v1 §7)
- Q: 007-6 Como nomear campo de nível e mensagem? → A: Mapeamento documentado nesta spec (fonte: respostas-joao.md, 007-6; origem: spec 006 research.md "Campo e mensagem já separados")
- Q: 007-7 Qual a granularidade do item rotulado? → A: Log e apontamento, ambos documentados (fonte: respostas-joao.md, 007-7; origem: T2 C)
- Q: 007-8 `tenant`/`produto` nos registros? → A: Obrigatórios, com valores de teste (fonte: respostas-joao.md, 007-8; origem: contrato log.ingested v1 §4)
- Q: 007-9 Grafia dos rótulos de classe? → A: ERROR / WARNING / INFO (fonte: respostas-joao.md, 007-9; origem: ADR_170; decisão taxonomia 14/11)
- Q: 007-10 Quem rotula? → A: Engenharia do Finder (fonte: respostas-joao.md, 007-10; origem: Jira LLT-198)
- Q: 007-11 Como validar o rótulo? → A: Só autoria + método documentado (fonte: respostas-joao.md, 007-11; origem: perguntas-etapa2b.md)
- Q: 007-12 Duplicata com rótulos diferentes? → A: Permitir e documentar (fonte: respostas-joao.md, 007-12; origem: perguntas-etapa2b.md)
- Q: 007-13 Registro sem nível ou sem mensagem? → A: Rejeitar (fonte: respostas-joao.md, 007-13; origem: spec 006 FR-006)
- Q: 007-14 Origem dos dados? → A: Anonimização de logs reais (fonte: respostas-joao.md, 007-14; origem: Jira LLT-198)
- Q: 007-15 Como auditar ausência de dado real de cliente? → A: Declaração de origem (fonte: respostas-joao.md, 007-15; origem: SC-001)
- Q: 007-16 PII fictícia ou casos LLT-202 na base? → A: Nenhuma PII, nem fictícia (fonte: respostas-joao.md, 007-16; origem: ADR_064)
- Q: 007-17 Rotular assinatura/agrupamento (LLT-200)? → A: Não; só classe (fonte: respostas-joao.md, 007-17; origem: FR-010)
- Q: 007-18 Casos em que só o catálogo decide? → A: Sim, subconjunto só-catálogo (fonte: respostas-joao.md, 007-18; origem: spec 006 US3)
- Q: 007-19 Relação com dados sintéticos da 006? → A: Substitui os sintéticos da 006 (fonte: respostas-joao.md, 007-19; origem: spec 006 plan.md)
- Q: 007-20 Divisão treino/validação/teste? → A: Split documentado (fonte: respostas-joao.md, 007-20; origem: perguntas-etapa2b.md)
- Q: 007-21 Amostra do SC-004 (formato)? → A: 100% dos registros (fonte: respostas-joao.md, 007-21; origem: SC-004)
- Q: 007-22 Consumidor principal da base? → A: Testes automatizados (fonte: respostas-joao.md, 007-22; origem: perguntas-etapa2b.md)
- Q: 007-23 Entrada direta no classificador sem fila? → A: Sim, desacoplada da fila (fonte: respostas-joao.md, 007-23; origem: spec 006 research.md)
- Q: 007-24 Volume/escala máximos? → A: Sem requisito de volume máximo (fonte: respostas-joao.md, 007-24; origem: perguntas-etapa2b.md)
- Q: 007-25 Onde publicar cobertura por classe? → A: Nesta spec e no checklist (fonte: respostas-joao.md, 007-25; origem: SC-003)
- Q: 007-26 Lacunas sem fonte bloqueiam `/speckit-plan`? → A: Sim; parar e consultar (fonte: respostas-joao.md, 007-26; origem: constituição princípio V)

### Session 2026-11-16 (etapa 3b)

- Q: 007-meta Qual é a meta numérica de precisão na base de validação (LLT-199 critério 2)? → A: **95% de precisão** (fonte: respostas-joao.md §Etapa 3; respostas-etapa3.md §1; origem: João 16/11 00h42; Jira LLT-199 critério 2)
- Q: DN-007-1 O que fazer com o subconjunto só-catálogo enquanto o catálogo não existe? → A: **Esperar o catálogo** para montar essa parte (opção C) (fonte: respostas-joao.md §Etapa 3; respostas-etapa3.md §2; origem: ADR_167 [PENDENTE]; 006-7 D)

### Session 2026-11-16 (etapa 3c)

- Q: DN-007-2 Qual é o escopo da meta de **95% de precisão** na base de validação (LLT-199 critério 2)? → A: **≥ 95% de precisão em cada classe** (ERROR, WARNING e INFO), não agregado global (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: DN-007-2 B; 006-2 B; Jira LLT-199 critério 2)
- Q: 007-formato Qual é o **formato de serialização** dos arquivos da base rotulada? → A: **JSONL** (um objeto JSON por linha) (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: 007 formato A; 007-4 A)
- Q: 007-caminho Qual é o **caminho físico** no repositório para os artefatos da base? → A: **`datasets/labeled-logs/`** na raiz do repositório `lhamalog-finder` (somente documentado nesta etapa; a pasta **não** é criada aqui) (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: 007 caminho C; 007-6 D)

### Session 2026-11-16 (etapa 3c — correções analyze)

- Q: C1 O Finder MUST agregar em ocorrência (constituição I) × **T1 revisto** (Módulo de Ocorrência cria o ocorrência)? → A: **Não alterar** a constituição nesta feature; manter **emenda do princípio I** como pendência formal; esta base **não** grava ocorrência (fonte: respostas-joao.md, T1 REVISTO 15/11; origem: revisao-T1-drive.md §2; analyze C1)
- Q: H2 O **Status** da spec ainda aponta para `/speckit-plan` com `plan.md` e `tasks.md` já existentes? → A: Status atualizado para fase **pós-plan / implementação MVP**, com ressalvas documentadas (fonte: analyze-007.md H2; artefatos `plan.md`, `tasks.md` nesta pasta)
- Q: H3 Metas numéricas de **recall/F1** para SC-005 / FR-007 além da precisão ≥95% por classe? → A: **Sem decisão** — permanece **DN-007-3** aberta; precisão por classe fechada em **DN-007-2 B** (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; respostas-etapa3.md §1; analyze H3). **Superado** na etapa 3e (**DN-007-3 B**); ver Session 2026-11-16 (etapa 3e).
- Q: M3 Contagem numérica exata para “milhares” e “mesmo mínimo por classe”? → A: **Sem número fixo** nas fontes; ordem de **milhares** no total e **mesma contagem mínima** entre ERROR, WARNING e INFO; valores numéricos **na entrega** (MANIFEST/checklist) (fonte: respostas-joao.md, 007-1 C / 007-2 B; analyze M3)
- Q: M1 Onde fica o gold `class_label` em **registro de apontamento** (único no agregado vs por `member`)? → A: **Sem decisão** nas fontes prioritárias — abrir **DN-007-4**; **007-7 C** só exige os dois tipos de item documentados (fonte: respostas-joao.md, 007-7 C; analyze M1; origem: perguntas-etapa2b.md §007-7). **Superado** na etapa 3e (**DN-007-4 A**); ver Session 2026-11-16 (etapa 3e).
- Q: L2 DN-007-1 e DN-007-2 ainda listadas em “Dúvidas novas (sem fonte)”? → A: Manter apenas referência curta “respondida”; detalhe em Clarifications etapas 3b/3c (fonte: respostas-joao.md §Etapa 3 / §Etapa 3b; analyze L2)

### Session 2026-11-16 (etapa 3c — analyze2 rodada 2)

- Q: I1 O checklist built-in `checklists/requirements.md` afirma ausência de lacunas e critérios plenos sem refletir **DN-007-3**, **DN-007-4**, split e gold em apontamento? → A: Revalidar itens de completude/mensurabilidade; **não** marcar pass pleno onde `revisao-etapa3c.md` mantém **CHK013** (split) pendente — **DN-007-3** / **DN-007-4** fechadas na etapa 3e (fonte: analyze-007-pos-correcoes.md I1; analyze-007-etapa3e.md I1; `checklists/revisao-etapa3c.md`; respostas-joao.md §Etapa 3c — DN-006-2/007-3 B: João 16/11 11h01 e 11h; DN-007-4 A: João, 16/11 11h02)

### Session 2026-11-16 (etapa 3e)

- Q: DN-007-3 Quais são as metas numéricas de **recall** e **F1** por classe (ERROR, WARNING, INFO) na validação LLT-199/197, além da precisão ≥ 95% por classe? → A: **≥ 95% de recall e ≥ 95% de F1 em cada classe** (ERROR, WARNING e INFO), não agregado global (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h01 e 11h; respostas-etapa3c.md §1; origem: DN-006-2/007-3 B; 006-2 B; Jira LLT-199 critério 2; LLT-51 item aberto supersedido neste alvo)
- Q: DN-007-4 Onde fica o rótulo gold `class_label` em **registro de apontamento** (agregado vs por `member`)? → A: **Um** `class_label` gold **por apontamento** (agregado); `members` MUST NOT carregar `class_label` divergente do apontamento (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); respostas-etapa3c.md §3; origem: DN-007-4 A; 007-7 C; pista Jira LLT-199 critério 1 «todo log» — não decide formato da base)

### Session 2026-11-16 (etapa 3f)

- Ajuste: **DN-007-4 A** — fonte canônica `respostas-joao.md §Etapa 3c (João, 16/11 11h02)` (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-007-4 A; linha 39 SUPERADO).
- Ajuste: `checklists/requirements.md` — **DN-007-3 B**, **DN-007-4 A** e **T043** não constam mais como lacuna/bloqueio documental; ressalva permanece só **split** (**CHK013**) e conteúdo só-catálogo (**006-7 D**) (fonte: analyze-007-etapa3e.md I1/I2; etapa 3e).
- Ajuste: Harmonizar notas da etapa 3c em Clarifications (**H3**, **M1**, **I1**) com estado pós-3e (fonte: analyze-007-etapa3e.md I3; respostas-joao.md §Etapa 3c, DN-006-2/007-3 B).

### Session 2026-11-16 (etapa 3g)

- Q: **007 H2** Como comparar o gold por apontamento com as predições **por log** na medição LLT-199/006-2 B? → A: **Opção D** — a classe **PREVISTA** do apontamento é a **mais grave** entre as classes previstas dos seus `members` (**ERROR** > **WARNING** > **INFO**, ADR_170); essa prevista MUST ser comparada ao **único** `class_label` gold do apontamento (**DN-007-4 A**); a **unidade** da avaliação é o **apontamento** (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), 007 H2 D; origem: ADR_170; respostas-etapa3f.md §D). **Exemplo:** gold ERROR e predições ERROR, ERROR, WARNING → prevista ERROR → **acerto**. **Contraexemplo:** gold WARNING e predições WARNING, WARNING, ERROR → prevista ERROR → **erro** (falso positivo para ERROR). **Implicação:** alinha avaliação à regra de classe da ocorrência (clarificação 09/11 citada pela ADR_170); a precisão de ERROR fica **mais exigente** que WARNING e INFO. **Contexto (não decide formato da base):** Jira LLT-199 critério 1 — «todo log processado recebe uma classe» (fonte: respostas-etapa3f.md §O que as fontes dizem; DN-007-4 A permanece no gold agregado).

## Decisões tomadas (com fonte)

| ID | Decisão | Fonte |
|----|---------|-------|
| **Objetivo** | Conjunto de logs rotulados para **avaliação** das regras no MVP (não treino de modelo no MVP). | (fonte: respostas-joao.md, 007-3; origem: ADR_167; Jira LLT-198) |
| **007-1** | Tamanho mínimo na ordem de **milhares** de registros. | (fonte: respostas-joao.md, 007-1; origem: perguntas-etapa2b.md) |
| **007-2** | **Mesmo mínimo por classe** (ERROR, WARNING, INFO). | (fonte: respostas-joao.md, 007-2; origem: perguntas-etapa2b.md) |
| **007-4** | **Arquivos versionados no repositório** (Git). | (fonte: respostas-joao.md, 007-4; origem: perguntas-etapa2b.md) |
| **007-formato** | Serialização dos registros: **JSONL** (um objeto JSON por linha). | (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: 007 formato A) |
| **007-caminho** | Diretório documentado: **`datasets/labeled-logs/`** na raiz de `lhamalog-finder` (criação física fora desta etapa de especificação). | (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: 007 caminho C) |
| **007-5** | Por registro: texto de **nível**, texto de **mensagem**, `lhamalog_batch_id`, `lhamalog_record_id`. | (fonte: respostas-joao.md, 007-5; origem: contrato log.ingested v1 §7; ADR_146 v2) |
| **007-6** | **Tabela de mapeamento** nesta spec (campos lógicos → nomes no artefato de dados); sem fixar nomes OTel na spec 006. | (fonte: respostas-joao.md, 007-6; origem: spec 006 research.md) |
| **007-7** | Dois tipos de item: **registro de log** e **apontamento** (vários logs); papéis documentados. | (fonte: respostas-joao.md, 007-7; origem: T2 C) |
| **007-8** | `tenantId` e `productId` **obrigatórios**, valores de teste documentados. | (fonte: respostas-joao.md, 007-8; origem: contrato §4) |
| **007-9** | Rótulos de classe: **ERROR**, **WARNING**, **INFO** (equivalentes a erro, alerta, informação). | (fonte: respostas-joao.md, 007-9; origem: ADR_170) |
| **007-10** | Rotulação: **engenharia do Finder**. | (fonte: respostas-joao.md, 007-10; origem: Jira LLT-198) |
| **007-11** | Validação: **autoria + método documentado** (sem revisão 100% obrigatória). | (fonte: respostas-joao.md, 007-11; origem: perguntas-etapa2b.md) |
| **007-12** | Duplicatas com rótulos diferentes: **permitidas**, com documentação. | (fonte: respostas-joao.md, 007-12; origem: perguntas-etapa2b.md) |
| **007-13** | Registro sem nível ou sem mensagem: **rejeitar**. | (fonte: respostas-joao.md, 007-13; origem: spec 006) |
| **007-14** | Origem: **anonimização de logs reais** (sem dados reais de cliente identificáveis). | (fonte: respostas-joao.md, 007-14; origem: Jira LLT-198) |
| **007-15** | Auditoria de origem: **declaração de origem** acompanhando o conjunto. | (fonte: respostas-joao.md, 007-15; origem: SC-001) |
| **007-16** | **Nenhuma PII** no conjunto, **nem fictícia** (sem casos LLT-202 nesta base). | (fonte: respostas-joao.md, 007-16; origem: ADR_064) |
| **007-17** | Escopo de rótulo: **só classe** (sem assinatura/agrupamento LLT-200). | (fonte: respostas-joao.md, 007-17; origem: LLT-200 fora de escopo) |
| **007-18** | Incluir **subconjunto só-catálogo** (nível não decide; catálogo decide). | (fonte: respostas-joao.md, 007-18; origem: spec 006 US3; catálogo: 006-7 D) |
| **007-meta** | Na base de validação (LLT-199 critério 2): **≥ 95% de precisão, ≥ 95% de recall e ≥ 95% de F1 em cada classe** (ERROR, WARNING e INFO). | (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; respostas-joao.md §Etapa 3c, João 16/11 11h; origem: DN-007-2 B; DN-006-2/007-3 B; Jira LLT-199; 006-2 B) |
| **DN-007-2** | Escopo dos 95% de **precisão**: **por classe**, não global. | (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: DN-007-2 B) |
| **DN-007-3** | **Recall** e **F1**: meta **≥ 95% em cada classe** (ERROR, WARNING e INFO), simétrica à precisão (DN-007-2 B). | (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h01 e 11h; origem: DN-006-2/007-3 B; 006-2 B) |
| **DN-007-4** | Gold em apontamento: **um** `class_label` no registro agregado; membros **sem** gold próprio. | (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); origem: DN-007-4 A; 007-7 C) |
| **007 H2** | Avaliação de apontamento: classe prevista = **mais grave** entre predições dos `members` (ERROR > WARNING > INFO); comparar ao gold agregado; **unidade** = apontamento. | (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), 007 H2 D; origem: ADR_170; respostas-etapa3f.md §D; DN-007-4 A) |
| **DN-007-1** | Subconjunto só-catálogo: **aguardar o catálogo** (006-7 D) antes de montar registros desse subconjunto; o restante da base não depende disso. | (fonte: respostas-joao.md §Etapa 3; respostas-etapa3.md §2; origem: ADR_167 [PENDENTE]) |
| **007-19** | Esta base **substitui** os dados sintéticos de prova da spec 006 (`plan.md`). | (fonte: respostas-joao.md, 007-19; origem: spec 006 plan.md) |
| **007-20** | **Split documentado** (partições treino/validação/teste definidas na documentação, para uso pós-MVP). | (fonte: respostas-joao.md, 007-20; origem: perguntas-etapa2b.md) |
| **007-21** | SC-004: verificação de formato em **100%** dos registros. | (fonte: respostas-joao.md, 007-21; origem: SC-004) |
| **007-22** | Consumidor principal: **testes automatizados** (e métricas LLT-197/199 quando aplicável). | (fonte: respostas-joao.md, 007-22; origem: perguntas-etapa2b.md) |
| **007-23** | Alimentação do classificador: **desacoplada da fila** (entrada direta na avaliação). | (fonte: respostas-joao.md, 007-23; origem: spec 006 research.md) |
| **007-24** | Sem requisito de **volume máximo**. | (fonte: respostas-joao.md, 007-24; origem: perguntas-etapa2b.md) |
| **007-25** | Contagens e cobertura por classe publicadas nesta **spec** e no **checklist**. | (fonte: respostas-joao.md, 007-25; origem: SC-003) |
| **007-26** | Lacunas sem fonte **bloqueiam** `/speckit-plan`. | (fonte: respostas-joao.md, 007-26; origem: constituição V) |
| **T2** | Uma ocorrência por achado (ERROR/WARNING/INFO). | (fonte: respostas-joao.md, T2 C; origem: ADR_035/040) |
| **T3** | Recorrência via consulta ao OpenObserve. | (fonte: respostas-joao.md, T3 B; origem: ADR_168) |
| **T4** | Lote incompleto não conta como sucesso de varredura. | (fonte: respostas-joao.md, T4 A; origem: spec 005) |
| **T5** | Contrato **log.ingested v1** é a única versão válida (`schemaVersion` "1"). | (fonte: respostas-joao.md, T5 D; origem: contrato 08/11 §4) |
| **T6** | Lacunas da 006 atualizam a própria spec 006 (não esta feature). | (fonte: respostas-joao.md, T6 A; origem: contexto-LLT-103) |
| **T7** | **INFO** gera ocorrência como ERROR e WARNING (decisão consciente, diverge do Drive). | (fonte: respostas-joao.md, T7 A e rodada 16/11 00h01; origem: ADR_170) |
| **T1 (revisto)** | **Módulo de Ocorrência** cria a ocorrência; Finder envia `ResumoAnalise` (não grava MongoDB nem publica evento de ocorrência). | (fonte: respostas-joao.md, T1 REVISTO 15/11; origem: revisao-T1-drive.md §2, ADR_062) |
| **Pendência — emenda princípio I** | Constituição I ainda descreve agregação no Finder; **T1 revisto** prevalece no escopo das specs; emenda em `.specify/memory/constitution.md` **fora** desta pasta (não editar aqui). | (fonte: analyze C1; respostas-joao.md, T1 REVISTO; origem: `plan.md` § Pendências) |
| **Taxonomia** | Classes vigentes: ERROR, WARNING, INFO (spam e "não classificado" não são rótulos válidos). | (fonte: respostas-joao.md, 007-9; origem: ADR_170; constituição III) |
| **006-2 B** | Métricas de avaliação: precisão, recall e F1 **por classe** quando esta base for usada (fora do escopo de entrega desta feature). | (fonte: respostas-joao.md, 006-2 B; origem: Jira LLT-197/199) |

### Mapeamento de campos (007-6 D)

A entrada do classificador (spec 006) usa **texto do campo de nível** e **texto da mensagem** já separados; os arquivos da base MUST expor esses valores nos campos lógicos abaixo. Os nomes das colunas/atributos nos arquivos versionados são definidos no artefato de dados e MUST constar nesta tabela (atualizada junto do conjunto):

| Campo lógico | Uso | Referência de alinhamento (sem impor nome OTel na 006) |
|--------------|-----|--------------------------------------------------------|
| `level_text` | Texto do campo de nível entregue ao classificador | Contrato de exemplo: `severityText` (contrato log.ingested v1, anexo JSON documental §6) |
| `message_text` | Texto da mensagem / corpo do log | Contrato de exemplo: corpo do registro (`body` no exemplo documental do contrato §6) |
| `lhamalog_batch_id` | Identificador de lote no OpenObserve | Contrato §7 |
| `lhamalog_record_id` | Identificador de registro no lote | Contrato §7 |
| `tenantId` | Isolamento multi-tenant (valor de teste) | Contrato §4 |
| `productId` | Produto (valor de teste) | Contrato §4 |
| `class_label` | Rótulo gold: ERROR, WARNING ou INFO | 007-9 A; ADR_170 |

### Cobertura alvo (007-1, 007-2, 007-25)

Valores numéricos finais MUST ser preenchidos quando o conjunto for entregue; até lá, aplicam-se as decisões abaixo:

| Métrica | Valor acordado |
|---------|----------------|
| Escala total | Ordem de **milhares** de registros — **sem inteiro mínimo** fixo nas fontes; preencher na entrega (007-1 C; analyze M3) |
| Por classe | **Mesma contagem mínima** entre ERROR, WARNING e INFO — **valor numérico** só na entrega (007-2 B; analyze M3) |
| Publicação | Contagens reais registradas em **checklists/requirements.md** e nesta seção (007-25 A). Lacuna de **split** qualitativo (**CHK013**) permanece em `revisao-etapa3c.md` |
| Precisão (LLT-199 critério 2) | **≥ 95% em cada classe** ERROR, WARNING e INFO (007-meta; DN-007-2 B) |
| Recall (006-2 B × LLT-199) | **≥ 95% em cada classe** ERROR, WARNING e INFO (DN-007-3 B) |
| F1 (006-2 B × LLT-199) | **≥ 95% em cada classe** ERROR, WARNING e INFO (DN-007-3 B) |
| Formato de arquivo | **JSONL** (007-formato A) |
| Caminho no repositório | **`datasets/labeled-logs/`** na raiz (007-caminho C) |

### Artefatos físicos (007-formato, 007-caminho, 007-6 D)

Os registros MUST ser persistidos como **JSONL** sob **`datasets/labeled-logs/`** na raiz do repositório `lhamalog-finder`. Nomes de arquivos (ex.: por partição do split ou por `itemType`) e nomes físicos de atributos JSON MUST constar na tabela de mapeamento quando os arquivos existirem; até a implementação, esta spec documenta apenas o diretório e o formato acordados (sem criar a pasta nesta etapa).

## Divergências documentadas

| Divergência | Lado A | Lado B | Prevalece | Por quê |
|-------------|--------|--------|-----------|---------|
| **INFO vira ocorrência** | Drive: SSID_007.1 FA01 e SSID_009.1 critério 1 — só ERRO/WARNING viram ocorrência; ADR_040 achados não-spam | João: **T7 A** e rodada 16/11 00h01 — INFO também vira ocorrência | **João** | Decisão consciente registrada em `respostas-joao.md`; ADR_170 deixa INFO em aberto no Drive antigo |
| **T1 — quem cria ocorrência** | João 14/11 **T1 A**: Finder grava e publica | Drive: ADR_062, arquitetura-modulos v2 — **Módulo Ocorrência** cria; Finder só `ResumoAnalise` | **T1 revisto (Drive)** | `respostas-joao.md` 15/11 23h47; alinhado a `revisao-T1-drive.md` §2 |
| **Taxonomia / spam** | ADR_019 (spam); SSID_007.1 legado | João + ADR_170: ERROR/WARNING/INFO | **João / ADR_170** | Esta base e a 006 usam três classes; emenda de ADRs fora do escopo |
| **LLT-199 critério (1)** | "Classe ou **não classificado**" | Taxonomia de três classes | **João / 007-9** | "Não classificado" não é rótulo na base |
| **007-14 × texto Jira LLT-198** | Jira: fictícios **ou** anonimizados | João **007-14 B**: só **anonimização de logs reais** | **João** | `respostas-joao.md`; LLT-198 não atualizado no Jira |
| **ADR_146 v2 "schemaVersion 2" × contrato 08/11** | Texto v2 em alguns trechos ainda menciona schema 2 | Contrato 08/11: **somente v1** válida | **Contrato v1 (T5 D)** | `contrato-log-ingested-v1.txt` §2–4; T5 D |
| **008-2 × contrato §4** | João **008-2 A**: `correlationId`/`traceId` na assinatura | Contrato: `correlationId` **aleatório por exportação**, não deduplica | **João na 008**; **risco** para agrupamento entre exportações | Registrar na spec 008; citado aqui porque a base espelha ids de lote/registro, não correlationId |
| **008-11 × 008-15** | 008-11 B inicial: chave sem classe; 16/11 **008-11** revisto: assinatura + **classe** + produto + tenant | 008-15 A: uma ocorrência por assinatura **e** classe | **Rodada 16/11** | `respostas-joao.md` pós-revisão T1 |
| **Jira LLT-201 × Drive (publicação)** | Jira: Finder avisa Ocorrências e Notificação | Drive: Finder → Módulo Ocorrência → Core → Notificação | **Drive + T1 revisto** | `revisao-T1-drive.md` §1 |
| **LLT-202 × ADR_064 (PII)** | Jira: PII conhecidos | Drive: só e-mail e linhas senha/token/secret | **Drive** | `respostas-documentadas.md` divergência 4; **irrelevante** para esta base (007-16 A proíbe PII) |
| **007-meta × 006-2 B** | Jira LLT-199 critério 2 e `respostas-etapa3c.md` §1: só **precisão** com meta acordada antes da etapa 3c | João **DN-006-2/007-3 B**: **≥ 95% recall e F1 por classe**, alinhado a **006-2 B** | **João etapa 3c** | `respostas-joao.md` §Etapa 3c; `respostas-etapa3c.md` §1 |
| **DN-007-4 × LLT-199 crit. 1** | Jira LLT-199 critério 1: «**todo log** processado recebe uma classe» (puxa para gold por membro na **classificação**) | **DN-007-4 A**: gold **um por apontamento** na base rotulada; membros sem `class_label` | **DN-007-4 A** na base; métricas de classificação seguem consumo da base (contrato 006) | `respostas-etapa3c.md` §3; (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-007-4 A) |
| **LLT-199 crit. 1 × 007 H2** | LLT-199 crit. 1: classificação **por log** («todo log») | **007 H2 D**: na **avaliação** com gold por apontamento, agregar predições dos `members` pela **classe mais grave** (ADR_170); **unidade** da métrica = apontamento para esses itens | **007 H2 D** para métricas na base; classificador ainda prediz **por log** (ADR_167) | (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), 007 H2 D; respostas-etapa3f.md §D; origem: ADR_170; Jira LLT-199) |
| **007-18 × catálogo [PENDENTE]** | **007-18 B**: exige subconjunto só-catálogo na base | **ADR_167** / **006-7 D**: catálogo e sintaxe ainda **[PENDENTE]** | **DN-007-1 C**: montar o subconjunto **após** o catálogo existir | `respostas-etapa3.md` §2: sem divergência Jira × Drive; LLT-198 sem critério de aceite sobre isso |

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Dispor conjunto rotulado para avaliação da classificação (Priority: P1)

Quem define ou valida regras (LLT-197, LLT-199) precisa de registros com rótulo **ERROR**, **WARNING** ou **INFO**, originados por **anonimização de logs reais**, sem dados reais de clientes, em volume na ordem de **milhares**, com **mesma contagem mínima por classe**.

**Why this priority**: LLT-199 critério 2 e LLT-197 dependem de base de validação; a spec 006 remete a LLT-198 e SC-007 cita métricas por classe quando a base existir.

**Independent Test**: Abrir os arquivos versionados no repositório, ler a declaração de origem e confirmar: (a) procedimento de anonimização descrito, (b) cada rótulo é ERROR, WARNING ou INFO, (c) contagens por classe iguais no mínimo acordado e total em milhares, (d) zero PII no conjunto.

**Acceptance Scenarios**:

1. **Given** o conjunto publicado, **When** o revisor lê a declaração de origem, **Then** ela descreve anonimização de logs reais e ausência de dados reais de clientes (007-14 B, 007-15 A).
2. **Given** um registro, **When** o revisor lê `class_label`, **Then** o valor é exatamente ERROR, WARNING ou INFO (007-9 A).
3. **Given** o conjunto completo, **When** o revisor busca spam ou "não classificado" como rótulo, **Then** nenhum registro usa essas classes (taxonomia; LLT-199 alinhado).
4. **Given** as contagens publicadas no checklist, **When** o revisor confere por classe, **Then** ERROR, WARNING e INFO têm a **mesma** contagem mínima e o total está na ordem de milhares (007-1 C, 007-2 B, 007-25 A).

---

### User Story 2 - Alinhar registros ao formato de entrada do classificador (Priority: P1)

A engenharia do Finder (007-10 B) monta registros com nível, mensagem, `lhamalog_batch_id`, `lhamalog_record_id`, `tenantId` e `productId` de teste, mapeados conforme a tabela desta spec (007-5 B, 007-6 D, 007-8 A). O consumidor principal são **testes automatizados** que alimentam o classificador **sem passar pela fila** (007-22 B, 007-23 A).

**Independent Test**: Em 100% dos registros (007-21 A), validar presença dos campos obrigatórios e que `level_text` + `message_text` alimentam a mesma entrada que a spec 006 espera.

**Acceptance Scenarios**:

1. **Given** um registro da base, **When** se monta a entrada do classificador, **Then** existem `level_text` e `message_text` distinguíveis (007-5, 007-6).
2. **Given** a intenção de espelhar captura real, **When** o revisor confere identificadores, **Then** cada registro traz `lhamalog_batch_id` e `lhamalog_record_id` (contrato §7; 007-5 B).
3. **Given** tenant e produto, **When** o revisor valida isolamento, **Then** `tenantId` e `productId` estão presentes com valores de teste documentados (007-8 A).
4. **Given** a tabela de mapeamento, **When** o revisor abre o arquivo de dados, **Then** os nomes físicos das colunas correspondem à tabela atualizada (007-6 D).
5. **Given** os artefatos versionados, **When** o revisor confere local e serialização, **Then** os registros estão em **JSONL** sob **`datasets/labeled-logs/`** na raiz do repositório (007-formato A, 007-caminho C).

---

### User Story 3 - Granularidade log vs apontamento e uso pós-MVP (Priority: P2)

A base inclui itens do tipo **log único** e **apontamento** (vários logs), com papéis documentados (007-7 C). No MVP por regras, a base serve **só para avaliação** (007-3 A). Um **split documentado** prepara partições para treino futuro (007-20 B), sem obrigar treino no MVP (ADR_167).

**Independent Test**: Verificar na documentação os dois tipos de item, o split nomeado e a declaração de que o MVP usa apenas avaliação.

**Acceptance Scenarios**:

1. **Given** o MVP por regras, **When** o PO lê a documentação do conjunto, **Then** o uso permitido no MVP é **somente avaliação das regras** (007-3 A).
2. **Given** a entrega da base, **When** o revisor lista partições, **Then** existe split documentado (nomes e critérios) para uso pós-MVP (007-20 B).
3. **Given** um item marcado como apontamento, **When** o revisor abre o registro, **Then** a estrutura de múltiplos logs está documentada e distinta do item log único (007-7 C).
4. **Given** um registro `itemType=apontamento`, **When** o revisor confere o gold, **Then** existe exatamente **um** `class_label` no apontamento e nenhum `member` traz `class_label` (DN-007-4 A).
5. **Given** documentação de avaliação, **When** o revisor lê a regra para apontamentos, **Then** a classe prevista do apontamento é a **mais grave** entre as predições dos `members` (ERROR > WARNING > INFO) e é comparada ao gold agregado; a unidade de contagem na métrica é o **apontamento** (007 H2 D; DN-007-4 A).

---

### Edge Cases

- Rótulo fora de ERROR/WARNING/INFO: registro **inválido** para a base (007-9; constituição III).
- Dado identificado como real de cliente não anonimizado: MUST NOT entrar (007-14 B; declaração de origem).
- Registro sem `level_text` ou sem `message_text`: **rejeitar** na validação do conjunto (007-13 A).
- Duplicata lógica com rótulos diferentes: **permitida** se documentada (motivo, versão, split) (007-12 B).
- Subconjunto só-catálogo: registros onde só o catálogo decide a classe; **entrega desse subconjunto adiada** até existir catálogo (DN-007-1 C; 006-7 D; ADR_167 [PENDENTE]); demais registros da base seguem independentes.
- Item `apontamento`: **um** `class_label` no agregado; `members` **sem** gold divergente (DN-007-4 A).
- Avaliação de apontamento: classificar **cada** `member` via classificador; derivar classe prevista do apontamento = **máxima gravidade** entre predições (ERROR > WARNING > INFO, ADR_170); comparar ao gold agregado — **unidade** = apontamento (007 H2 D). Ex.: gold ERROR + predições ERROR, ERROR, WARNING → acerto; gold WARNING + predições WARNING, WARNING, ERROR → erro (falso positivo ERROR) (fonte: respostas-etapa3f.md §D).
- Item `log` único: unidade de avaliação = **registro de log**; comparação direta `class_label` vs predição do mesmo registro.
- Validação LLT-199 critério 2 / 006-2 B: **≥ 95%** de precisão, recall e F1 **em cada classe** ERROR, WARNING e INFO na partição de validação (007-meta; DN-007-2 B; DN-007-3 B); para itens `apontamento`, métricas MUST usar a regra **007 H2 D** por apontamento.
- INFO na base: alinhado a T7 A e ADR_170; diverge de SSID_007.1 legado (ver Divergências).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: MUST entregar conjunto documentado com rótulos de classe, por **anonimização de logs reais**, sem dados reais de clientes identificáveis; declaração de origem obrigatória (007-14 B, 007-15 A; Jira LLT-198).
- **FR-002**: Cada `class_label` MUST ser **ERROR**, **WARNING** ou **INFO**. Spam e "não classificado" MUST NOT ser usados (007-9 A; ADR_170; constituição III).
- **FR-003**: O conjunto MUST ter volume na ordem de **milhares** de registros e **mesma contagem mínima** em cada classe, **sem** exigir número inteiro mínimo acordado nas fontes prioritárias (valores finais na entrega); contagens MUST ser publicadas nesta spec e no checklist (007-1 C, 007-2 B, 007-25 A; analyze M3).
- **FR-004**: Arquivos **versionados no repositório** sob **`datasets/labeled-logs/`** na raiz, em **JSONL**; cada linha MUST ser um objeto JSON com `level_text`, `message_text`, `lhamalog_batch_id`, `lhamalog_record_id`, `tenantId`, `productId` e `class_label`; mapeamento físico (nomes de arquivos e atributos) documentado na tabela de mapeamento e na seção Artefatos físicos (007-4 A, 007-formato A, 007-caminho C, 007-5 B, 007-6 D, 007-8 A).
- **FR-005**: No MVP por regras (ADR_167), o conjunto serve **somente para avaliação das regras**; MUST documentar split treino/validação/teste para uso pós-MVP (nomes e critérios qualitativos — **sem** proporção, seed ou estratificação numérica enquanto não houver fonte; 007-20 B).
- **FR-006**: Documentação MUST declarar autoria (engenharia do Finder) e **método de rotulação/validação** (007-10 B, 007-11 D).
- **FR-007**: Na validação LLT-199 critério 2, **precisão, recall e F1** MUST ser **≥ 95% em cada classe** ERROR, WARNING e INFO na partição de validação (007-meta; DN-007-2 B; DN-007-3 B; 006-2 B). Para `itemType=log`, a unidade de conta é o **registro de log** (comparação direta gold vs predição). Para `itemType=apontamento`, a unidade de conta é o **apontamento**, com agregação de predições conforme **FR-007c**.
- **FR-007b**: Registro `itemType=apontamento` MUST ter **um único** `class_label` gold no agregado; cada `member` MUST NOT incluir `class_label` (DN-007-4 A; 007-7 C).
- **FR-007c**: Na avaliação de `itemType=apontamento`, MUST obter predição **por** `member` (classificação por log, ADR_167); a classe **prevista** do apontamento MUST ser a **mais grave** entre essas predições (**ERROR** > **WARNING** > **INFO**, ADR_170); MUST comparar essa prevista ao único `class_label` gold do apontamento (007 H2 D; DN-007-4 A). Documentação MUST registrar exemplo, contraexemplo e que ERROR na métrica é mais exigente que WARNING/INFO (respostas-etapa3f.md §D).
- **FR-008**: Métricas além da precisão (LLT-197) seguem **006-2 B** (precisão, recall, F1 por classe) quando a base for consumida; não inventar métricas aqui.
- **FR-009**: Esta base **substitui** os dados sintéticos de prova da spec 006; consumo primário via **testes automatizados**, entrada **desacoplada da fila** (007-19 B, 007-22 B, 007-23 A).
- **FR-010**: Rotulagem MUST ser **só classe**; MUST NOT incluir assinatura/agrupamento (LLT-200) nem PII, real ou fictícia (007-17 A, 007-16 A).
- **FR-011**: MUST incluir **subconjunto só-catálogo** documentado; os **registros** desse subconjunto MUST NOT ser inventados antes do catálogo existir (**DN-007-1 C**); até lá, a spec MUST declarar o critério de inclusão alinhado à spec 006 US3 e a dependência de **006-7 D** (ADR_167 [PENDENTE]).
- **FR-012**: Duplicatas com rótulos diferentes: **permitidas** com documentação; registros sem nível ou mensagem: **rejeitar** (007-12 B, 007-13 A).
- **FR-013**: Sem requisito de volume máximo (007-24 A).

### Fora de escopo (com fonte)

- Treinar classificador ou usar LLM no MVP (007-3 A; ADR_167).
- Definir catálogo de expressões (006 FR-005; 006-7 D).
- Rotular assinatura, PII ou mascaramento (007-16 A, 007-17 A; LLT-200, LLT-202).
- Alterar ADR_019, ADR_167 ou SSID_007.1.
- Alterar `.specify/memory/constitution.md` (**emenda do princípio I** × T1 revisto — pendência transversal; analyze C1).
- Agrupar ocorrências ou publicar eventos (T1 revisto; specs 008/009).

### Key Entities

- **Conjunto rotulado (base)**: Artefatos **JSONL** em **`datasets/labeled-logs/`**, versionados no Git; escala milhares; split documentado; declaração de origem.
- **Registro de log**: Unidade com um par nível+mensagem e ids de lote/registro; `class_label`.
- **Registro de apontamento**: Unidade com vários logs agregados para cenários T2; estrutura e papel documentados (007-7 C); **um** `class_label` gold no agregado; `members` somente entrada (nível/mensagem/ids), sem gold (DN-007-4 A). Na avaliação: predição por `member`, **classe prevista do apontamento** = mais grave entre predições (007 H2 D); unidade de métrica = apontamento.
- **Rótulo de classe (`class_label`)**: ERROR, WARNING ou INFO.
- **Declaração de origem**: Texto que atesta anonimização e ausência de dado real de cliente (007-15 A).
- **Partição (split)**: Subconjunto nomeado (ex.: treino, validação, teste) para uso pós-MVP (007-20 B).
- **Subconjunto só-catálogo**: Parte da base reservada a cenários em que só o catálogo decide a classe; **conteúdo entregue após** catálogo (DN-007-1 C).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% dos registros cobertos pela declaração de origem (anonimização; sem dado real de cliente) (007-15 A).
- **SC-002**: 100% dos registros com `class_label` em {ERROR, WARNING, INFO}; 0% spam ou "não classificado" (007-9 A).
- **SC-003**: Contagens publicadas: mesma contagem mínima nas três classes; total na ordem de milhares (007-1 C, 007-2 B, 007-25 A).
- **SC-004**: 100% dos registros expõem campos obrigatórios de FR-004 e respeitam o mapeamento (007-21 A).
- **SC-005**: Na base de validação (LLT-199 critério 2): **precisão, recall e F1 ≥ 95% em cada classe** ERROR, WARNING e INFO (007-meta; DN-007-2 B; DN-007-3 B; FR-007, **FR-007c**; 006-2 B). Itens `apontamento` entram nas métricas **por apontamento** com classe prevista = mais grave entre predições dos `members` (007 H2 D).
- **SC-006**: Documentação declara uso MVP **somente avaliação** e descreve split para pós-MVP (007-3 A, 007-20 B).

## Dúvidas novas (sem fonte)

DN-007-1 — respondida em 2026-11-16 (etapa 3b), ver Clarifications → Session 2026-11-16 (etapa 3b).

DN-007-2 — respondida em 2026-11-16 (etapa 3c), ver Clarifications → Session 2026-11-16 (etapa 3c).

DN-007-3 — respondida em 2026-11-16 (etapa 3e), ver Clarifications → Session 2026-11-16 (etapa 3e).

DN-007-4 — respondida em 2026-11-16 (etapa 3e), ver Clarifications → Session 2026-11-16 (etapa 3e).

## Evidência e consulta

- **Fato**: Decisões em "Decisões tomadas", Clarifications 2026-11-16 (etapas 3b–3g): 007-meta, DN-007-1 C, DN-007-2 B, DN-007-3 B, DN-007-4 A, **007 H2 D**, 007-formato A, 007-caminho C; requisitos com IDs 007-x e transversais citados.
- **Premissa**: Nenhuma premissa nova sem fonte além das já respondidas nas fontes prioritárias.

## Questões em aberto

1. **Emenda do princípio I**: constituição × T1 revisto — fora desta pasta (analyze C1; `plan.md` § Pendências).