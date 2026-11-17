# Feature Specification: Mascarar dados sensíveis antes de gravar a ocorrência

**Feature Branch**: `010-mascarar-dados-sensiveis` (trabalho permanece na branch git `Development`; nenhuma branch dedicada foi criada nesta execução, conforme instrução do solicitante)

**Created**: 2026-11-14

**Status**: Draft (clarificação 010-1–010-28; etapas 3b/3c/3e/3h; **010 H1**, **DN-010-3 B**, **DN-010-4 C**, **DN-008-5 A** em 2026-11-16)

**Input**: Jira **LLT-202** (texto exato do solicitante); Epic Finder — IA própria de classificação e agregação; Feature **Agregação em ocorrências (E07.F3)**; item **E07.F3.3** em `LhamaLog_backlog_pendente.docx` (14/10/2026); contexto em `/workspace/lhamalog/contexto-LLT-103.md` §2.4 e §3; instruções do João em 14/11/2026.

**Rastreabilidade**: Complementa o que a spec `006-modelo-classificacao-logs` deixa **fora** do catálogo de mascaramento (**FR-009** na 006). Relaciona-se a captura do log na spec `003-consumir-logs-openobserve` (texto integral na etapa de captura). Pré-condição em **SSID_007.1** (Finder já mascarou — ADR_064). Item aberto de produto **LLT-51**. Epic **LLT-134** / feature E07.F3.

## Clarifications

### Session 2026-11-16

Clarificação não interativa: respostas 010-1 a 010-28 de `respostas-joao.md`, `respostas-documentadas.md` e `revisao-T1-drive.md` (010-6, 010-17, 010-21, 010-23, 010-25), com T1 revisto (Módulo de Ocorrência cria a ocorrência; Finder envia `ResumoAnalise`).

- Q: **010-1** Qual é o catálogo de padrões de segredo e PII? → A: Conteúdo da ADR_064 (e-mails substituídos; linhas com senha/token/secret → marcador); lista **única e compartilhada** Finder e Sugestão (ADR_079); lista concreta de regex segue **[PENDENTE]** na ADR_064 até a ADR oficial (010-2). (fonte: respostas-documentadas.md, 010-1; origem: ADR_064, ADR_079)
- Q: **010-2** Onde fica a lista oficial de padrões? → A: Numa **ADR** (nova ou emenda da ADR_064) — **PENDÊNCIA**: ADR ainda não criada/emendada nesta execução; a spec **referencia** a ADR sem listar regex. (fonte: respostas-joao.md, 010-2; origem: ADR_064 [PENDENTE])
- Q: **010-3** Momento do mascaramento no pipeline? → A: **Antes da classificação/análise** (`validarSeguranca()` antes do MotorAnalise; alinhado ao título LLT-202 “antes de gravar”). (fonte: respostas-documentadas.md, 010-3; origem: ADR_064, SSID_007.1)
- Q: **010-4** `montarContexto()` e “antes do MotorAnalise” são o mesmo estágio? → A: **Mesmo estágio** (uma única passada de `validarSeguranca()`). (fonte: respostas-documentadas.md, 010-4; origem: ADR_016, ADR_064)
- Q: **010-5** O classificador vê texto bruto ou mascarado? → A: **Mascarado**. (fonte: respostas-documentadas.md, 010-5; origem: ADR_064, SSID_007.1)
- Q: **010-6** O que é “trecho guardado” na ocorrência? → A: **Exemplos mascarados** (`AchadoAnalise.linhas` no Finder; Módulo de Ocorrência persiste em `conteudo` da Ocorrência). (fonte: respostas-documentadas.md, 010-6 + revisao-T1-drive.md, 010-6; origem: ADR_040, ADR_167, E08.F1.3/LLT-205)
- Q: **010-8** Limite de exemplos de log por ocorrência? → A: **Sem limite nesta história**. (fonte: respostas-joao.md, 010-8; origem: backlog E08.F1.3)
- Q: **010-9** O que vai à Sugestão? → A: **Amostras mascaradas** (`conteudo` mascarado da Ocorrência e metadados). (fonte: respostas-documentadas.md, 010-9; origem: ADR_039, ADR_079, LLT-212)
- Q: **010-10** Texto do marcador da linha removida? → A: **`"log com senha omitido"`** (texto único para senha/token/secret). (fonte: respostas-joao.md, 010-10; origem: ADR_064 exemplo)
- Q: **010-11** Formato da substituição de e-mail? → A: **Texto fixo** (placeholder; ex.: `[email]`). (fonte: respostas-joao.md, 010-11; origem: ADR_064 “substituídos”)
- Q: **010-12** Outras PII (CPF, telefone…)? → A: **Só o que estiver no catálogo** (hoje, só e-mail + linhas senha/token/secret na ADR_064). (fonte: respostas-documentadas.md, 010-12; origem: ADR_064, RIPD — vale o Drive)
- Q: **010-13** Dado sensível fora do catálogo? → A: **Permanece bruto**; alcance limitado e documentado (DPA). (fonte: respostas-documentadas.md, 010-13; origem: DPA Anexo II/IV)
- Q: **010-14** Unidade do mascaramento? → A: **Linha a linha**. (fonte: respostas-documentadas.md, 010-14; origem: ADR_064, ADR_011)
- Q: **010-15** Frequência por lote? → A: **Uma vez por lote** (arquivo/lote antes da análise). (fonte: respostas-documentadas.md, 010-15; origem: ADR_064, ADR_016)
- Q: **010-16** Falha de `validarSeguranca()`? → A: **Segue o processamento sem os trechos** (não bloqueia o lote inteiro). (fonte: respostas-joao.md, 010-16)
- Q: **010-17** Onde fica `dadosSensiveisIdentificados`? → A: Nos **contadores** do `ResumoAnalise` e do **resumo do arquivo** (“dados sensíveis mascarados”). (fonte: revisao-T1-drive.md, 010-17; origem: ADR_102, ADR_064)
- Q: **010-18** Como Finder e Sugestão compartilham a lista de regex? → A: **Arquivo versionado** lido pelos dois módulos. (fonte: respostas-joao.md, 010-18; origem: ADR_079)
- Q: **010-19** Paridade de versão do catálogo? → A: **Sempre a mesma versão** entre Finder e Sugestão. (fonte: respostas-joao.md, 010-19; origem: ADR_079)
- Q: **010-20** Critério 2 (“telas, e-mails e LLM”) por módulo? → A: Finder **produz** texto mascarado; Sugestão **reaplica** as mesmas regras no prompt; telas/e-mails usam trechos já mascarados. (fonte: respostas-documentadas.md, 010-20; origem: ADR_167, ADR_079)
- Q: **010-21** Quem garante mascarado em tela e e-mail? → A: O **Finder** garante que o `ResumoAnalise` (e `AchadoAnalise.linhas`) **já sai mascarado**; Módulo Ocorrência/Core/Notificação **não** revalidam. (fonte: revisao-T1-drive.md, 010-21; origem: ADR_167)
- Q: **010-22** Qual canal é o “e-mail” do critério 2? → A: **Core → Notificação** (e-mail externo pela Notificação). (fonte: respostas-documentadas.md, 010-22; origem: ADR_165, LLT-218)
- Q: **010-23** Trechos mascarados no aviso HTTP Módulo Ocorrência → Core? → A: **Não** incluir texto de log; na **etapa 3b** o aviso passa a levar **id + tipo + gravidade** (009-5 revisto; ver Session 2026-11-16 etapa 3b). (fonte: revisao-T1-drive.md, 010-23 + respostas-joao.md Etapa 3; origem: 009-5 revisto, respostas-etapa3.md §7)
- Q: **010-24** O log bruto continua acessível? → A: **Só no OpenObserve**; o Finder não reexpõe o bruto. (fonte: respostas-documentadas.md, 010-24; origem: ADR_146 v2, contrato 08/11 §2)
- Q: **010-25** A spec 006 passa a “guardar” trecho mascarado? → A: A **006 não persiste** texto de log; o **Finder** envia linhas mascaradas em **`AchadoAnalise.linhas`**; o **Módulo de Ocorrência grava** em `conteudo` (esta spec 010 descreve o mascaramento; persistência na Ocorrência é efeito do contrato Drive, não emenda da 006 nesta execução). (fonte: revisao-T1-drive.md, 010-25; origem: ADR_040, ADR_167)
- Q: **010-26** Escopo da LLT-51? → A: **Fora** desta história (LLT-202). (fonte: respostas-documentadas.md, 010-26; origem: Jira LLT-51, LLT-202)
- Q: **010-27** Quem define casos de teste do mascaramento (SC-001/002)? → A: **O João** — **PENDÊNCIA explícita**: aguardando exemplos/casos dele; **não** criar conjunto de testes nesta spec. (fonte: respostas-joao.md, 010-27)
- Q: **010-28** Meta de desempenho do mascaramento? → A: **Sem meta** nesta feature. (fonte: respostas-joao.md, 010-28)

**Transversais aplicáveis a esta feature** (registradas para coerência com o pipeline):

- Q: **T1 (revisto)** Quem cria a ocorrência e quem publica evento à Notificação? → A: **Módulo de Ocorrência** cria a Ocorrência e avisa o **Core** (HTTP interno JWT, 009-3 D); Finder envia **`ResumoAnalise`** e **não** publica evento de ocorrência. (fonte: respostas-joao.md, T1 revisto 15/11; origem: ADR_062, SSID_009.1, revisao-T1-drive.md)
- Q: **T5** Versão do contrato `log.ingested`? → A: **v1** (08/11) é a única válida. (fonte: respostas-joao.md, T5 D)
- Q: **T6** Lacunas da 006? → A: Atualizar a **própria 006** quando couber (fora do escopo de escrita desta execução). (fonte: respostas-joao.md, T6 A)
- Q: **T7 / 006-11** INFO gera ocorrência? → A: **Sim**, por achado INFO (decisão consciente **divergente** do Drive SSID_007.1/009.1/ADR_040). (fonte: respostas-joao.md, rodada 16/11 00h01; origem: ADR_170, SSID_007.1)

### Session 2026-11-16 (etapa 3b)

Clarificação não interativa: **DN-010-1**, **DN-010-2** e **009-5 revisto** em `respostas-joao.md` (Etapa 3, 16/11 00h41) e `respostas-etapa3.md` (§§7–8).

- Q: **DN-010-1** De onde vem o resumo mascarado do e-mail (LLT-218)? → A: A **Notificação lê a Ocorrência pelo id** no MongoDB compartilhado; o `conteudo` já está mascarado (mesma decisão que **DN-009-1 A**). (fonte: respostas-joao.md, DN-009-1 A; origem: respostas-etapa3.md §7, LLT-218, ADR_040, ADR_054)
- Q: **009-5 revisto** O que o aviso Módulo Ocorrência → Core carrega além do identificador? → A: **id + tipo + gravidade**, **sem** texto de log (mantém 010-23: evento/aviso **não** inclui trechos de log). (fonte: respostas-joao.md, Etapa 3; origem: respostas-etapa3.md §7, diagrama Notificação v2)
- Q: **DN-010-2** Atributos OTel entram no mascaramento do LhamaLog? → A: **Não** — o mascaramento do LhamaLog vale **só para as linhas de log**; atributos OTel sensíveis são responsabilidade do **Cliente na origem** (Coletor/collector). **Ressalva ADR_040:** `traceId`, `correlationId`, `hostName`, `moduloId` copiam para a Ocorrência **como identificadores** (sem passar pelo catálogo de linha). **Alerta RIPD:** ampliar mascaramento (CPF, cartão, telefone) antes de dados reais permanece registrado (010-12). (fonte: respostas-joao.md, DN-010-2 / Etapa 3; origem: respostas-etapa3.md §8, DPA Anexo IV item 6, DPA Anexo II, ADR_064, ADR_040, RIPD)

### Session 2026-11-16 (etapa 3c)

Clarificação não interativa: **DN-007-2** e **DN-008-4** em `respostas-joao.md` (§Etapa 3b, João 16/11 09h02).

- Q: **DN-007-2** Como a meta **007-meta** (≥95% de precisão) se relaciona com esta feature (LLT-202)? → A: **B** — **≥95% de precisão por classe** (ERROR, WARNING, INFO) na base de validação da spec **007**; **esta spec (010) não define gate métrico** — o mascaramento permanece regra de **SC-001**–**SC-004** e FR-001–FR-006; o aceite de **precisão de classificação** fica nas specs **006**/**007** (**006-2 B**, **SC-007** / **007-meta**), não nos critérios de sucesso do Jira LLT-202. (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: respostas-etapa3.md §1, LLT-199)
- Q: **DN-008-4** Qual algoritmo normaliza a **mensagem** na composição da **assinatura** (e equivalências em 006/009)? → A: **B** — aparar espaços nas extremidades; ignorar maiúsculas/minúsculas; substituir **números**, **UUIDs**, **IPs**, **datas** e **hexadecimais** por **marcadores** (um por categoria, distintos entre categorias). Algoritmo **canônico** na spec **008** (dona); **006**, **009** e **010** MUST reutilizá-lo sem regra divergente ao derivar assinatura/equivalência a partir da mensagem. Nesta spec: o texto da mensagem usado na normalização MUST ser o de linha **já mascarado** por `validarSeguranca()` (**010-5 A**, **FR-002**), aplicado **antes** da classificação e **antes** da normalização. Literais dos marcadores de normalização: **DN-008-5 A** (`<NUM>` `<UUID>` `<IP>` `<DATE>` `<HEX>`, spec **008**) — registrados nesta spec só como referência cruzada (**FR-015**). (fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: DN-008-1 B, ADR_013)
- Q: **U1** / **DN-010-3** *(pré-3e)* — política de mismatch: **superada** pela Session **2026-11-16 (etapa 3e)** (**DN-010-3 B**).
- Q: **U2** / **DN-010-4** *(pré-3e)* — caminho/formato: **superada** pela Session **2026-11-16 (etapa 3e)** (**DN-010-4 C**).

### Session 2026-11-16 (etapa 3e)

Clarificação não interativa: **DN-010-3 B**, **DN-010-4 C** e **DN-008-5 A** em `respostas-joao.md` (§Etapa 3c, João 16/11 11h) e `respostas-etapa3c.md` (§§4, 7–8).

- Q: **DN-010-3** O que o Finder (e/ou o deploy) faz quando a **versão** lida do catálogo **não** coincide com a versão **combinada** para paridade com a Sugestão? → A: **B** — o **serviço recusa subir** (falha na **inicialização**) se a versão do catálogo for **diferente da combinada** (**010-19 A**: mesma versão obrigatória Finder e Sugestão). (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h; origem: ADR_079, 010-19 A; respostas-etapa3c.md §7)
- Q: **DN-010-4** Qual o **caminho** no monorepo e o **formato** do arquivo versionado do catálogo compartilhado Finder–Sugestão? → A: **C** — **caminho e formato definidos na ADR da lista** (**010-2 A**); permanecem **pendência ligada a P-010-2** até a ADR publicar metadados; esta spec **referencia** a ADR e **não** fixa path/formato próprios (**010-18 A** = mecanismo arquivo versionado no repositório). (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h; origem: ADR_079, 010-2 A, 006-3 C; respostas-etapa3c.md §8)
- Q: **DN-008-5** Quais os textos exatos dos **marcadores do normalizador** de mensagem na assinatura? → A: **A** — **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`** (algoritmo canônico na spec **008**, **DN-008-4 B**; nesta spec o texto normalizado parte de linha **já mascarada** — **FR-015**. **Não** confundir com o marcador de **mascaramento** de linha **`"log com senha omitido"`** — **010-10 A** / **FR-006**). (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h; origem: DN-008-4 B, ADR_013; respostas-etapa3c.md §4)

### Session 2026-11-16 (etapa 3f)

Clarificação não interativa (ajustes etapa 3f; sem perguntas).

- Ajuste: removida a entrada **D13** e qualquer tratamento de **11h01 × 11h** como divergência entre respostas do João — não há conflito; a linha “puladas” da §Etapa 3c (João, 16/11 11h01) está **SUPERADA** e as respostas definitivas de **DN-006-3 B**, **DN-007-4 A** e **DN-008-7 B** constam na mesma rodada às **11h02**. (fonte: respostas-joao.md l.39–51; instrução João etapa 3f)
- Ajuste: fonte canônica de referência cruzada **DN-006-3 B**, **DN-007-4 A**, **DN-008-7 B** (specs 006/007/008; impacto indireto em 010) → **`respostas-joao.md` §Etapa 3c (João, 16/11 11h02)**. (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h02, DN-006-3 B / DN-007-4 A / DN-008-7 B; l.49–51; instrução João etapa 3f)
- Ajuste: `checklists/requirements.md` — item «No [NEEDS CLARIFICATION] markers remain» desmarcado; lacunas explícitas em `plan.md` e `data-model.md` (biblioteca regex Go, placeholder de e-mail, tipos Go, path/formato catálogo **P-010-2**) permanecem sem resolução nesta etapa. (fonte: analyze-010-etapa3e.md **F1**; instrução João etapa 3f)

### Session 2026-11-16 (etapa 3c — correções analyze)

Clarificação não interativa: correções dos achados **analyze-010.md** (etapa 3c), com fonte nas respostas já documentadas.

- Q: **I1** Quem publica o aviso HTTP ao Core e o que ele transporta em relação a trechos de log? → A: O **Módulo de Ocorrência** publica o aviso HTTP **interno** ao Core (**009-3 D**); o payload MUST incluir **id**, **tipo** e **gravidade** da ocorrência e MUST **not** incluir trechos de log (**009-5 revisto**). O **Finder** MUST **not** implementar esse aviso; garante apenas texto mascarado no `ResumoAnalise` / consumo downstream (010-20, 010-21). (fonte: respostas-joao.md, 009-3 D / 009-5 revisto Etapa 3; origem: respostas-etapa3.md §7, T1 revisto)
- Q: **I2** Como unificar `dadosSensiveisIdentificados`, o rótulo “dados sensíveis mascarados” e nomes no handoff? → A: **`dadosSensiveisIdentificados`** é o **total numérico** retornado por `validarSeguranca()` (ADR_064). Nos **contadores** do `ResumoAnalise` e do resumo do arquivo, o **rótulo humano** é **“dados sensíveis mascarados”** (010-17; ADR_102 [lista exata de contadores INFERIDA]). No mapeamento AchadoAnalise → Ocorrência (Drive), o campo de achado aparece como **`dadosSensiveis`** — sem fixar chave JSON definitiva nesta spec até ADR/N3; **não** confundir com `dadosSensiveisMascarados` usado só como rótulo em contratos de handoff. (fonte: revisao-T1-drive.md, 010-17; origem: ADR_102, ADR_064)
- Q: **E2** Qual evidência mínima verificável no Finder para **SC-006**? → A: No **Finder**, evidência MUST ser o **identificador de versão** do artefato de catálogo **lido** na carga do lote (ou na inicialização), **registrado de forma auditável** (010-18 A). A **paridade Finder–Sugestão** em todos os ambientes depende de coordenação **fora deste repositório** (Sugestão Java; 010-19 A); o aceite **cross-módulo** não é gate único desta pasta — só a leitura/versionamento no Finder. (fonte: respostas-joao.md, 010-18 A / 010-19 A; origem: ADR_079)
- Q: **A1** O que significa “execuções acordadas” em **SC-001**–**SC-004**? → A: O conjunto de **casos de teste definidos pelo João** (**010-27 A**; pendência **P-010-27**), não um catálogo inventado nesta spec. (fonte: respostas-joao.md, 010-27 A; origem: LLT-202)

### Session 2026-11-16 (etapa 3h)

Clarificação não interativa (ajustes etapa 3h; sem perguntas).

- Ajuste **010 H1**: A spec **010** **reaproveita** o normalizador canônico da spec **008** (`internal/application/assinatura/normalizar_mensagem.go`, **T016**) — **MUST NOT** criar algoritmo/`normalizar_mensagem` próprio (**FR-015**). Fonte canônica **João 16/11 13h22** (respostas-joao.md §Etapa 3g); substitui referência à orientação do Gestor LhamaLog 16/11 12h37 “sujeita à aprovação do João nas tasks”. (fonte: respostas-joao.md, 010 H1; origem: Etapa 3g, DN-008-4 B)
- Ajuste **DN-008-9 A**: **Não alterado** nesta pasta — a spec **010** não documenta composição textual da assinatura (**008 assinatura A**, separador **DN-008-8 A**); escape **`|`** → **`\|`** nos componentes é dona spec **008** (**DN-008-9 A** respondida em respostas-joao.md §Etapa 3g). (fonte: respostas-joao.md, DN-008-9 A; origem: Etapa 3g)
- Ajuste **DLQ B**: **Não alterado** nesta pasta — a spec **010** não cita **DLQ**, **DN-008-3** nem **DN-008-6**; após 3 falhas na consulta OpenObserve o Finder avisa o Core por API interna (`batchId`, motivo) sem DLQ no **log.ingested v1** permanece nas specs **008**/**009**. (fonte: respostas-joao.md, DLQ B; origem: respostas-etapa3g.md §1, contrato log.ingested v1, DN-006-1b B)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Trechos guardados sem segredos e PII expostos (Priority: P1)

Quando logs são processados para compor o que será **repassado** ao Módulo de Ocorrência e **persistido** como exemplos mascarados, padrões cobertos pelo catálogo (ADR oficial — ver PENDÊNCIA 010-2) não permanecem legíveis no mesmo formato do log bruto.

**Why this priority**: Critério de aceite (1) do Jira LLT-202; pré-condição SSID_007.1; ADR_064 e ADR_167 colocam o mascaramento no Finder.

**Independent Test**: Após o João entregar os **casos de teste** acordados (PENDÊNCIA 010-27) e a **ADR** fechar a lista de regex (PENDÊNCIA 010-2), processar entradas que contenham padrões cobertos e verificar `AchadoAnalise.linhas` e `conteudo` da Ocorrência.

**Acceptance Scenarios**:

1. **Given** trecho com padrão de **segredo** coberto pelo catálogo (FR-001), **When** o Finder aplica mascaramento **antes** da classificação (FR-002), **Then** as **`linhas`** enviadas no achado MUST NOT expor o segredo no formato original (J: LLT-202 critério 1; ADR_064; 010-5 A).
2. **Given** trecho com **e-mail**, **When** o mascaramento é aplicado, **Then** o e-mail MUST ser substituído por **placeholder fixo** (FR-006; 010-11 A; ADR_064).
3. **Given** linha com **senha**, **token** ou **secret**, **When** mascarada, **Then** a linha inteira MUST ser trocada pelo marcador **`"log com senha omitido"`** (FR-006; 010-10 A).
4. **Given** o Módulo de Ocorrência recebe o achado, **When** persiste a Ocorrência, **Then** o campo **`conteudo`** MUST conter apenas os **exemplos mascarados** recebidos em `AchadoAnalise.linhas` (010-6; 010-25; ADR_040).
5. **Given** o pipeline aplica classificação e, downstream, **normalização de mensagem** para assinatura (**DN-008-4 B**, spec **008**), **When** o Finder processa um lote, **Then** `validarSeguranca()` MUST ter sido aplicado **antes** desses estágios sobre as **linhas de log** (**FR-015**, **010-5 A**).

---

### User Story 2 - Consumo externo só vê texto mascarado (Priority: P1)

Telas, e-mails e fluxos que usam **LLM** (Sugestão) recebem **apenas** versão **mascarada** dos trechos no escopo desta feature, alinhada às mesmas regras na Sugestão (ADR_079; arquivo versionado — FR-008).

**Why this priority**: Critério de aceite (2) do Jira; ADR_167 liga trechos repassados à Ocorrência e à Sugestão; ADR_079 exige regras compartilhadas.

**Independent Test**: Rastrear um trecho desde `AchadoAnalise.linhas` até tela (E08.F1.3), prompt da Sugestão e canal de e-mail (Core → Notificação); verificar ausência de texto bruto sensível coberto pelo catálogo.

**Acceptance Scenarios**:

1. **Given** exemplos mascarados em `conteudo`, **When** o detalhe da ocorrência exibe “exemplos de log” (E08.F1.3), **Then** **somente** texto mascarado MUST ser apresentado (J: LLT-202 critério 2).
2. **Given** ADR_167 (Finder sem LLM), **When** o critério Jira cita “na LLM”, **Then** a obrigação MUST aplicar-se à **Sugestão**, que reaplica as mesmas regras (ADR_079; 010-20).
3. **Given** aviso Módulo Ocorrência → Core (009-3 D), **When** publicado, **Then** o payload MUST incluir **id**, **tipo** e **gravidade** da ocorrência e MUST **not** incluir trechos de log (009-5 revisto; 010-23; respostas-etapa3.md §7).
4. **Given** LLT-218 exige e-mail com resumo mascarado, **When** a Notificação formata o e-mail, **Then** MUST obter o texto mascarado **consultando a Ocorrência pelo id** (DN-010-1 A / DN-009-1 A); MUST **not** depender de trechos de log no aviso HTTP ao Core.
5. **Given** a **versão combinada** do catálogo configurada para paridade com a Sugestão (**010-19 A**), **When** na **inicialização** a versão lida do artefato for **diferente**, **Then** o serviço MUST **recusar subir** (**DN-010-3 B**, **FR-008**).

---

### User Story 3 - Rastreio de dados sensíveis identificados (Priority: P2)

O processamento registra **quantos** dados sensíveis foram identificados no lote, nos contadores do `ResumoAnalise` e do resumo do arquivo.

**Why this priority**: ADR_064 retorno; ADR_102 contadores; decisão 010-17.

**Independent Test**: Processar log de teste com quantidade conhecida de ocorrências de padrões cobertos e comparar contador “dados sensíveis mascarados”.

**Acceptance Scenarios**:

1. **Given** `validarSeguranca()` retorna total (ADR_064), **When** o lote é concluído, **Then** o total MUST aparecer nos **contadores** do ResumoAnalise/resumo do arquivo (FR-005; 010-17).

---

### Edge Cases

- **Lista de regex [PENDENTE na ADR]**: até a ADR oficial (010-2), padrão não catalogado permanece **bruto** (010-13 A); CPF/cartão/telefone **fora** do catálogo atual — risco RIPD R03 documentado (010-12).
- **Falha de mascaramento** (010-16 C): fluxo segue **sem** incluir trechos no achado; classificação/agregação podem prosseguir com o restante.
- **Remoção de linha** com token/senha pode ocultar contexto útil do erro **antes** da classificação (risco ADR_064 — aceito com 010-5 A).
- **Compartilhamento**: arquivo versionado no repositório (010-18 A); **caminho** e **formato** conforme **ADR da lista** (**DN-010-4 C**, **P-010-2**). **Mesma versão** obrigatória Finder e Sugestão (010-19 A); em **mismatch** com a versão combinada, o serviço **recusa subir** (**DN-010-3 B**).
- **Bruto**: apenas OpenObserve (010-24); fila e Ocorrência não reexpõem bruto.
- **Sem limite** de exemplos por ocorrência nesta história (010-8 C).
- **LLT-51**: minimização ampla **fora de escopo** (010-26 A).
- **Atributos OTel sensíveis** (e-mail em atributo, PII em resource/scope): **fora** do mascaramento desta feature — responsabilidade do **Cliente na origem** (DPA Anexo IV item 6; DN-010-2). Na Ocorrência, **traceId**, **correlationId**, **hostName** e **moduloId** seguem **ADR_040** como identificadores, copiados **como enviados** (não são alvo do catálogo linha a linha).
- **E-mail (LLT-218)**: resumo mascarado vem do **`conteudo`** da Ocorrência; Notificação **lê por id** após o aviso ao Core (DN-010-1 A). Template de e-mail no Drive ainda **[PENDENTE]** (respostas-etapa3.md §7) — impacto na spec **Notificação**, não altera regra de mascaramento no Finder.
- **RIPD**: recomendação de ampliar padrões antes de dados reais — **alerta** mantido (010-12; respostas-etapa3.md §8); não expande escopo LLT-202 nesta história.
- **Ordem no pipeline**: `validarSeguranca()` (**FR-002**) antes do classificador (**010-5 A**) e antes de qualquer **normalização de mensagem** para assinatura (**DN-008-4 B**, **FR-015**); algoritmo canônico na spec **008**; literais dos marcadores de normalização **`<NUM>` `<UUID>` `<IP>` `<DATE>` `<HEX>` (**DN-008-5 A**) — distintos do marcador de mascaramento **`"log com senha omitido"`** (**010-10 A**).
- **Meta de precisão de classificação**: **≥95% por classe** (**DN-007-2 B**) é gate das specs **006**/**007**, **não** desta feature LLT-202.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: O conjunto de **padrões** aplicados MUST ser o definido na **ADR oficial** do catálogo (**PENDÊNCIA 010-2**: ADR ainda não criada/emendada). **Comportamento vigente** até a lista fechar: e-mails substituídos; linhas com senha/token/secret removidas com marcador; lista **compartilhada** Finder e Sugestão (ADR_064, ADR_079). **Escopo de PII “conhecidos”** no produto atual: **só o que estiver no catálogo** (010-12 C), não o texto amplo do Jira sem ADR.
- **FR-002**: O mascaramento MUST ocorrer **uma vez por lote**, em **montarContexto()** / **antes do MotorAnalise** (mesmo estágio — 010-3 A, 010-4 A, 010-15 A). O **classificador** MUST operar sobre texto **mascarado** (010-5 A).
- **FR-003**: O texto que conta como “trecho **guardado**” MUST ser **exemplos mascarados**: o Finder envia em **`AchadoAnalise.linhas`**; o **Módulo de Ocorrência** persiste em **`conteudo`** da Ocorrência (010-6; 010-25; ADR_040). **Sem limite** de quantidade nesta história (010-8 C). A spec **006** continua **sem** obrigação de persistir texto de log; a gravação é responsabilidade do Módulo de Ocorrência ao converter o achado.
- **FR-004**: **Telas** (ex.: E08.F1.3), **e-mails** via Core → Notificação (010-22 A) e **entrada LLM** na Sugestão MUST usar **exclusivamente** texto **mascarado** produzido conforme FR-001 (J: LLT-202 critério 2). O **Módulo de Ocorrência**, no aviso HTTP **interno** ao Core (**009-3 D**), MUST transportar **id**, **tipo** e **gravidade** da ocorrência e MUST **not** carregar trechos de log (**009-5 revisto**; 010-23). O **Finder** MUST **not** publicar esse aviso (T1 revisto; D6). Para o e-mail (LLT-218), a Notificação MUST montar o resumo lendo a Ocorrência **pelo id** (conteúdo já mascarado — DN-010-1 A). **Não** exige LLM no Finder (ADR_167).
- **FR-004a**: O mascaramento desta feature MUST aplicar-se **somente às linhas de log** tratadas por `validarSeguranca()` (ADR_064; DN-010-2). Atributos OTel sensíveis MUST ser tratados pelo **Cliente na origem** antes do envio (DPA Anexo IV item 6; DPA Anexo II). Identificadores copiados para a Ocorrência conforme **ADR_040** (**traceId**, **correlationId**, **hostName**, **moduloId**) MUST **not** ser reinterpretados como “linha” para FR-006, salvo decisão futura se o Finder passar outros atributos nos trechos (ressalva respostas-etapa3.md §8).
- **FR-005**: O total **`dadosSensiveisIdentificados`** (retorno numérico de `validarSeguranca()` — ADR_064) MUST ser exposto nos **contadores** do `ResumoAnalise` e do **resumo do arquivo** com o **rótulo humano** “dados sensíveis mascarados” (010-17; ADR_102). O nome de campo no achado no Drive é **`dadosSensiveis`** (revisao-T1-drive.md); chave JSON final aguarda ADR/N3 — **não** usar `dadosSensiveisMascarados` como nome de campo canônico nesta spec.
- **FR-006**: **E-mail**: substituição por **placeholder fixo** (010-11 A). **Senha/token/secret**: **remoção da linha inteira** e marcador **`"log com senha omitido"`** (010-10 A). Unidade **linha a linha** (010-14 C).
- **FR-007**: O **ScannerSegurança** / **`validarSeguranca()`** permanecem os nomes citados nas ADRs_016 e ADR_064 **sem** fixar pacote ou linguagem nesta spec.
- **FR-008**: A lista de regex MUST residir em **arquivo versionado** no repositório (mecanismo — 010-18 A; ADR_079), lido por **Finder** e **Sugestão**, com **obrigatoriedade da mesma versão** em ambos (010-19 A). **Caminho** no monorepo e **formato** de serialização MUST ser os declarados na **ADR da lista** (**DN-010-4 C**, **010-2 A**; publicação em **P-010-2**). Se a ADR ainda não publicar path/formato, a implementação do adaptador MUST aguardar **P-010-2** — esta spec **não** inventa path nem formato. Se, na **inicialização**, a versão lida do catálogo for **diferente da versão combinada** (config/deploy acordada para paridade com a Sugestão), o serviço MUST **recusar subir** (**DN-010-3 B**). Conteúdo da lista = ADR oficial (**P-010-2** / 010-2 A).
- **FR-009**: Dado sensível **fora do catálogo** MUST permanecer **como enviado** (010-13 A; DPA). Esta feature MUST documentar o alcance limitado; **não** bloqueia o lote nem aplica marcador genérico automático.
- **FR-010**: Se **`validarSeguranca()`** falhar, o Finder MUST **continuar o processamento sem incluir trechos mascarados** no achado (010-16 C) — não é fail-closed para o lote inteiro.
- **FR-011**: O **log bruto** MUST permanecer acessível **somente no OpenObserve**; o Finder MUST NOT reexpô-lo em APIs, filas ou Ocorrência (010-24 C; contrato 08/11; ADR_146 v2).
- **FR-012**: Esta feature MUST NOT contradizer a spec 003: captura entrega registro **inteiro**; mascaramento na **etapa de processamento** do Finder.
- **FR-013**: Esta feature **não** altera **ADR_019**, **ADR_167** nem **SSID_007.1** por escrita; divergências de taxonomia/fluxo legado ficam em **Divergências documentadas**.
- **FR-014**: Escopo **LLT-51** MUST permanecer **fora** desta história (010-26 A).
- **FR-015**: O texto de **linha** (mensagem/corpo usado na análise) entregue ao **MotorAnalise** e a qualquer derivação de **mensagem normalizada** para assinatura (**DN-008-4 B**, spec **008**) MUST ser o produzido **após** `validarSeguranca()` na mesma passada do lote (**FR-002**, **010-5 A**). A spec **010** MUST **reaproveitar** o normalizador canônico da **008** (**010 H1**) e MUST **not** definir algoritmo de normalização divergente nem módulo `normalizar_mensagem` próprio; consumidores **006** e **009** reutilizam o mesmo canônico. Para referência cruzada (dono **008**): marcadores canônicos **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`** (**DN-008-5 A**) — **não** são o marcador de mascaramento de linha **`"log com senha omitido"`** (**010-10 A** / **FR-006**). (fonte: respostas-joao.md, 010 H1; origem: §Etapa 3g, DN-008-4 B)

### Fora de escopo (com fonte)

- Definir **classificação** por tipo erro/alerta/informação (spec 006; LLT-199) — salvo usar texto mascarado já decidido (010-5).
- **Criar** ou **emendar** ADRs nesta execução (incluindo catálogo 010-2).
- **Agrupar** logs, **publicar evento** (Finder não publica — T1 revisto), calcular **gravidade** no Finder, **notificar** (outras specs).
- **Capturar** log no OpenObserve (spec 003).
- Resolver **todo** o escopo **LLT-51** (010-26).
- Alterar **constitution**, **README** ou specs `001`–`009` **por escrita** nesta pasta.
- Escrever **código**, **casos de teste** do mascaramento (PENDÊNCIA 010-27 — João) ou **metas de desempenho** (010-28 A).
- **Emendar** ADR_019, ADR_167 ou SSID_007.1.
- Definir **gate métrico** **≥95% de precisão por classe** (ERROR, WARNING, INFO) — spec **007** / **006** (**DN-007-2 B**; **007-meta**).
- Ser **dona** do algoritmo de **normalização da mensagem** para assinatura — spec **008** (**DN-008-4 B**); esta spec só garante **texto mascarado antes** do uso (**FR-015**).

### Key Entities *(include if feature involves data)*

- **Trecho de log mascarado**: Saída linha a linha de `validarSeguranca()`; transportado em **`AchadoAnalise.linhas`** (Finder) e persistido em **`Ocorrência.conteudo`** (Módulo de Ocorrência).
- **Padrão de segredo / PII**: Definição oficial na **ADR** (PENDÊNCIA); comportamento interim ADR_064.
- **Marcador de omissão**: `"log com senha omitido"` (010-10 A).
- **Placeholder de e-mail**: Texto fixo acordado na ADR/lista (ex.: `[email]`) (010-11 A).
- **`dadosSensiveisIdentificados`**: Total numérico por lote (ADR_064); exposto nos contadores com rótulo **“dados sensíveis mascarados”** (010-17; ADR_102). Distinto do campo de achado **`dadosSensiveis`** no mapeamento Drive (revisao-T1-drive.md) e de rótulos em contratos de handoff — ver glossário em Clarifications (I2).
- **Arquivo versionado de regex**: Artefato compartilhado Finder–Sugestão (010-18 A); mesma versão obrigatória (010-19 A). **Caminho** e **formato**: declarados na **ADR da lista** (**DN-010-4 C**, **P-010-2**). **Mismatch** com versão combinada na inicialização: serviço **recusa subir** (**DN-010-3 B**).
- **ResumoAnalise / AchadoAnalise**: Contrato Finder → Módulo de Ocorrência (T1 revisto); achado inclui `linhas` mascaradas; Finder não grava Ocorrência.
- **Identificadores OTel na Ocorrência** (ADR_040): `traceId`, `correlationId`, `hostName`, `moduloId` — copiados como identificadores, **fora** do catálogo linha a linha (DN-010-2).
- **Canal e-mail (LLT-218)**: Notificação obtém resumo mascarado via leitura da Ocorrência por **id** (MongoDB compartilhado — ADR_054; DN-010-1 A).
- **Mensagem normalizada (assinatura)**: Entrada MUST ser linha **mascarada** quando derivada do log processado pelo Finder; algoritmo **DN-008-4 B** definido na spec **008** (trim, case-insensitive); marcadores **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`** (**DN-008-5 A**). Distinto de **`"log com senha omitido"`** (mascaramento — 010-10 A).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Quando existirem **casos de teste definidos pelo João** (PENDÊNCIA 010-27) e catálogo fechado na ADR (PENDÊNCIA 010-2), **100%** das **execuções acordadas** — isto é, do conjunto de casos entregue pelo João (**010-27 A**) — MUST NOT expor segredos/PII **cobertos** no formato original em `linhas`/`conteudo` (J: LLT-202 critério 1).
- **SC-002**: Em pontos **acordados** — conjunto de casos **010-27 A** — (**tela**, **e-mail** Core→Notificação com resumo lido da Ocorrência por id — DN-010-1 A, **Sugestão/LLM**), **100%** das amostras auditadas MUST usar **somente** texto mascarado nas **linhas** no escopo de FR-001/FR-004a (J: LLT-202 critério 2). Atributos OTel fora de linha seguem FR-004a (origem do Cliente).
- **SC-003**: Para e-mail (placeholder) e linha senha/token/secret (marcador fixo), **100%** das execuções de teste **acordadas** (casos João — **010-27 A**) MUST reproduzir FR-006 **quando** o padrão estiver no catálogo (ADR_064; 010-10, 010-11).
- **SC-004**: Em casos de teste com contagem conhecida (conjunto **010-27 A**), o contador “dados sensíveis mascarados” MUST coincidir com o esperado em **100%** das execuções acordadas (010-17; ADR_064).
- **SC-005**: **Enquanto** PENDÊNCIA **010-27** (casos de teste) ou **010-2** (ADR do catálogo) não fecharem, **SC-001**, **SC-002** e aceite formal completo permanecem **bloqueados** — não por lacuna de regra de negócio, mas por dependência externa documentada.
- **SC-006**: Com arquivo versionado (010-18 A) e **mesma versão obrigatória** entre Finder e Sugestão (010-19 A; FR-008): no **Finder**, evidência MUST registrar o **identificador de versão** do catálogo lido, de forma auditável, em **100%** das execuções de verificação de configuração acordadas (**010-27 A** quando aplicável). Se a versão lida na inicialização for **diferente da combinada**, o serviço MUST **recusar subir** e MUST deixar falha **observável** (log/métrica de readiness — **DN-010-3 B**). A confirmação de que a **Sugestão** leu a **mesma** versão é **coordenação cross-repositório** (módulo Java) e complementa, mas **não substitui**, a evidência no Finder. **Caminho** e **formato** do artefato MUST seguir a **ADR da lista** quando **P-010-2** fechar (**DN-010-4 C**).
- **SC-007** *(referência cruzada; não é aceite LLT-202)*: Meta **≥95% de precisão por classe** (ERROR, WARNING, INFO) na validação documentada na spec **007** (**DN-007-2 B**; **007-meta**); avaliação de **classificação** nas specs **006**/**007**, **não** substitui **SC-001**–**SC-004** de mascaramento.

## Evidência e consulta

- **Fato**: Decisões 010-1–010-28 integradas; T1 revisto; fontes Jira/Drive/ADR citadas por ID.
- **Premissa**: Nenhuma premissa nova sem fonte (constituição V).
- **Pendência**: ADR do catálogo (010-2); casos de teste do João (010-27).
- **Dúvidas etapa 3b**: DN-010-1 e DN-010-2 **respondidas** (ver Clarifications, Session etapa 3b).
- **Dúvidas etapa 3c**: **DN-007-2** e **DN-008-4** **respondidas** (ver Clarifications, Session etapa 3c).
- **Dúvidas etapa 3e**: **DN-010-3**, **DN-010-4** e **DN-008-5** (referência cruzada 008) **respondidas** (ver Clarifications, Session etapa 3e).
- **Dúvidas abertas**: nenhuma **DN-010-*** em aberto nesta pasta após etapa 3e.
- **Constituição**: **P-emenda-I** e **P-emenda-III** (não editar constituição nesta pasta).

## Dúvidas respondidas (histórico)

- **DN-010-1** — respondida em 2026-11-16 (etapa 3b), ver Clarifications → Session 2026-11-16 (etapa 3b).
- **DN-010-2** — respondida em 2026-11-16 (etapa 3b), ver Clarifications → Session 2026-11-16 (etapa 3b).
- **DN-007-2** — respondida em 2026-11-16 (etapa 3c), ver Clarifications → Session 2026-11-16 (etapa 3c).
- **DN-008-4** — respondida em 2026-11-16 (etapa 3c), ver Clarifications → Session 2026-11-16 (etapa 3c).
- **DN-010-3** — respondida em 2026-11-16 (etapa 3e), ver Clarifications → Session 2026-11-16 (etapa 3e).
- **DN-010-4** — respondida em 2026-11-16 (etapa 3e), ver Clarifications → Session 2026-11-16 (etapa 3e).
- **DN-008-5** — respondida em 2026-11-16 (etapa 3e), referência cruzada normalização (dono spec **008**); ver Clarifications → Session 2026-11-16 (etapa 3e).

## Dúvidas novas (sem fonte)

*Nenhuma DN desta feature em aberto após etapa 3e (2026-11-16). Histórico de opções **DN-010-3** / **DN-010-4** (analyze U1/U2) substituído pelas decisões **DN-010-3 B** e **DN-010-4 C** em Clarifications → Session etapa 3e.*

## Pendências explícitas

| ID | Descrição | Responsável / próximo passo |
|----|-----------|----------------------------|
| **P-010-2** | Lista oficial de regex em **ADR** (João 010-2 A); **caminho e formato** do artefato (**DN-010-4 C**) | Criar/emendar ADR **fora** desta pasta; spec referencia “conforme ADR”; adaptador FR-008 aguarda metadados publicados |
| **P-010-27** | **Casos de teste** do mascaramento para SC-001/002 | João define exemplos **antes** do plano; spec **não** lista casos |
| **P-emenda-I** | **Emenda do Princípio I** da constituição (Finder agrega “em ocorrência” × **T1 revisto**: Módulo de Ocorrência materializa Ocorrência) | Rastrear **fora** desta pasta; **não** editar `.specify/memory/constitution.md` aqui (analyze **C1**; plan T038) |
| **P-emenda-III** | Alinhar texto literal do **Princípio III** (taxonomia legada) a **ADR_170** / classes ERROR-WARNING-INFO | Fora desta pasta; mitigado em **D8** e FR-013 nesta spec (analyze **G1**) |

## Decisões tomadas (com fonte)

Cada item: **arquivo-fonte**, **ID**, **origem** quando houver.

### Transversais (vigentes para o pipeline desta feature)

| ID | Decisão | Fonte | Origem |
|----|---------|-------|--------|
| **T1 revisto** | Módulo de Ocorrência **cria** a Ocorrência; Finder envia `ResumoAnalise`; Finder **não** publica evento de ocorrência | respostas-joao.md | ADR_062, SSID_009.1, revisao-T1-drive.md |
| **T2** | Uma ocorrência por achado (grão) | respostas-joao.md, T2 C | ADR_035/040 |
| **T3** | Recorrência via consulta OpenObserve (sem banco vetorial) | respostas-joao.md, T3 B | ADR_168 |
| **T4** | Lote parcial reprocessável | respostas-joao.md, T4 A | ADR_014 |
| **T5** | Contrato `log.ingested` **v1** (08/11) vale | respostas-joao.md, T5 D | contrato-log-ingested-v1.txt |
| **T6** | Lacunas da 006 entram atualizando a 006 | respostas-joao.md, T6 A | — |
| **T7** | INFO gera ocorrência por achado (diverge Drive) | respostas-joao.md, 16/11 00h01 | SSID_007.1, SSID_009.1, ADR_040 |
| **Taxonomia** | Classes **ERROR / WARNING / INFO** no Finder; sem spam/“não classificado” | respostas-joao.md; 007-9 A | ADR_170, constituição III |

### Feature 010 (010-1 a 010-28)

| ID | Decisão | Fonte | Origem |
|----|---------|-------|--------|
| **010-1** | Catálogo = ADR_064 + compartilhamento ADR_079; regex [PENDENTE] na ADR | respostas-documentadas.md | ADR_064, ADR_079 |
| **010-2** | Lista oficial na **ADR** (não na spec 010) | respostas-joao.md | ADR_064 [PENDENTE] |
| **010-3** | Mascarar **antes** da classificação | respostas-documentadas.md | ADR_064 |
| **010-4** | `montarContexto()` = antes MotorAnalise | respostas-documentadas.md | ADR_016, ADR_064 |
| **010-5** | Classificador vê texto **mascarado** | respostas-documentadas.md | ADR_064 |
| **010-6** | Exemplos em `AchadoAnalise.linhas` → `conteudo` | respostas-documentadas.md + revisao-T1-drive.md | ADR_040, E08.F1.3 |
| **010-8** | Sem limite de exemplos nesta história | respostas-joao.md | LLT-205 / E08.F1.3 |
| **010-9** | Sugestão recebe amostras mascaradas | respostas-documentadas.md | ADR_039, ADR_079 |
| **010-10** | Marcador `"log com senha omitido"` | respostas-joao.md | ADR_064 |
| **010-11** | E-mail → placeholder fixo | respostas-joao.md | ADR_064 |
| **010-12** | Outras PII só se no catálogo | respostas-documentadas.md | ADR_064, RIPD |
| **010-13** | Fora do catálogo → **bruto** | respostas-documentadas.md | DPA |
| **010-14** | Mascaramento **linha a linha** | respostas-documentadas.md | ADR_064 |
| **010-15** | **Uma vez** por lote | respostas-documentadas.md | ADR_064 |
| **010-16** | Falha: **segue sem trechos** | respostas-joao.md | — |
| **010-17** | Contador no ResumoAnalise/resumo arquivo | revisao-T1-drive.md | ADR_102 |
| **010-18** | Lista: **arquivo versionado** | respostas-joao.md | ADR_079 |
| **010-19** | **Mesma versão** Finder e Sugestão | respostas-joao.md | ADR_079 |
| **010-20** | Finder produz; Sugestão reaplica regras | respostas-documentadas.md | ADR_167, ADR_079 |
| **010-21** | Finder garante `ResumoAnalise` mascarado | revisao-T1-drive.md | ADR_167 |
| **010-22** | E-mail = canal Core → Notificação | respostas-documentadas.md | ADR_165, LLT-218 |
| **010-23** | Aviso ao Core: **id + tipo + gravidade**, **sem** texto de log | respostas-joao.md, Etapa 3 + revisao-T1-drive.md | 009-5 revisto (respostas-etapa3.md §7) |
| **DN-010-1** | E-mail: Notificação lê Ocorrência **pelo id** | respostas-joao.md, DN-009-1 A | respostas-etapa3.md §7, LLT-218, ADR_054 |
| **DN-010-2** | Mascaramento LhamaLog **só linhas**; OTel sensível na **origem** (Cliente) | respostas-joao.md, Etapa 3 | respostas-etapa3.md §8, DPA Anexo IV/II, ADR_064, ADR_040 |
| **010-24** | Bruto só OpenObserve | respostas-documentadas.md | ADR_146 v2, contrato §2 |
| **010-25** | 006 não persiste; Módulo grava `conteudo` | revisao-T1-drive.md | ADR_040 |
| **010-26** | LLT-51 fora desta história | respostas-documentadas.md | LLT-202 |
| **010-27** | Casos de teste: **João** (pendente) | respostas-joao.md | — |
| **010-28** | **Sem** meta de tempo | respostas-joao.md | — |
| **DN-007-2** | **≥95% precisão por classe** (ERROR, WARNING, INFO) — gate **006**/**007**, **não** LLT-202/010 | respostas-joao.md §Etapa 3b, João 16/11 09h02 | respostas-etapa3.md §1, 007-meta, LLT-199 |
| **DN-008-4** | Normalização mensagem na assinatura: trim, case-insensitive, marcadores número/UUID/IP/data/hex; **dono** spec **008**; 010 garante texto **mascarado antes** (**FR-015**) | respostas-joao.md §Etapa 3b, João 16/11 09h02 | DN-008-1 B, ADR_013 |
| **010 H1** | A **010 reaproveita** o normalizador canônico da **008** — **MUST NOT** criar segundo `normalizar_mensagem` nesta feature | respostas-joao.md §Etapa 3g, João 16/11 13h22 | DN-008-4 B, FR-015 |
| **DN-008-5** | Literais **`<NUM>` `<UUID>` `<IP>` `<DATE>` `<HEX>`** (normalizador; distinto de 010-10 A) — referência cruzada; dono spec **008** | respostas-joao.md §Etapa 3c, João 16/11 11h | respostas-etapa3c.md §4, DN-008-4 B, ADR_013 |
| **DN-010-3** | Versão do catálogo ≠ versão combinada → serviço **recusa subir** | respostas-joao.md §Etapa 3c, João 16/11 11h | respostas-etapa3c.md §7, ADR_079, 010-19 A |
| **DN-010-4** | **Caminho** e **formato** do catálogo só na **ADR da lista** (010-2 A); pendência **P-010-2** | respostas-joao.md §Etapa 3c, João 16/11 11h | respostas-etapa3c.md §8, ADR_079, 006-3 C |

### Referência cruzada (specs 006 / 007 / 008 — dono nas respectivas pastas)

| ID | Decisão (resumo) | Fonte | Origem |
|----|------------------|-------|--------|
| **DN-006-3** | **B** — meta **2 s por lote**, com tamanho de lote e máquina de referência declarados | respostas-joao.md §Etapa 3c (João, 16/11 11h02) | spec 006 |
| **DN-007-4** | **A** — rótulo **gold por apontamento** | respostas-joao.md §Etapa 3c (João, 16/11 11h02) | spec 007 |
| **DN-008-7** | **B** — **uma consulta** ao Módulo de Ocorrência por achado | respostas-joao.md §Etapa 3c (João, 16/11 11h02) | spec 008 |

### Jira LLT-202 (fatos mantidos)

- Esconder senhas, tokens e dados pessoais nos trechos guardados; critérios (1) padrões mascarados e (2) só mascarado em tela, e-mail e LLM — interpretados com Drive (010-12; LLM só Sugestão).

## Divergências documentadas

| Id | Fontes | O que diverge | Vale / tratamento nesta spec |
|----|--------|----------------|------------------------------|
| **D1** | João T7/006-11 (16/11) × SSID_007.1 crit.1 × SSID_009.1 crit.1 × ADR_040 | INFO **vira** ocorrência vs. Drive “só ERRO/WARNING” | **Vale o João** (decisão consciente); Finder envia achados INFO no ResumoAnalise; impacto no Módulo Ocorrência — fora da escrita desta pasta |
| **D2** | LLT-202 “PII conhecidos” × ADR_064 × RIPD R03 | Jira amplo vs. Drive só e-mail + senha/token/secret | **Vale o Drive** (010-12 C); RIPD exige ampliar **antes de dados reais** — risco registrado, não escopo LLT-202 |
| **D3** | ADR_146 v2 “única cópia do conteúdo” × ADR_035/040 `conteudo` mascarado | Parece proibir texto fora do OpenObserve | **Leitura adotada** (respostas-documentadas.md): bruto só OpenObserve; Ocorrência guarda **trechos mascarados** (010-6, 010-24) |
| **D4** | ADR_146 v2 schemaVersion 2 × contrato 08/11 v1 | Versão de envelope da fila | **Vale contrato v1** para `log.ingested` (T5 D); v2 não altera decisões de mascaramento da 010 |
| **D5** | T1 A antiga (João 14/11) × Drive ADR_062/040 | Quem grava ocorrência | **Substituída** por T1 revisto (15/11); 010-6/21/25 usam contrato Drive |
| **D6** | LLT-201 épico Finder “publicar evento” × arquitetura v2 / SSID_009.1 | Quem avisa Notificação | **Vale o Drive**: Módulo Ocorrência → Core (009-3 D); Finder só ResumoAnalise |
| **D7** | 008-2 A (correlationId na assinatura) × contrato 08/11 §4 | correlationId **aleatório** por exportação | **Vale o João**; **risco** de assinatura instável — registrar na spec 008, não altera regras de mascaramento |
| **D8** | SSID_007.1 / ADR_019 / ADR_167 legado × ADR_170 + João 14/11 | Taxonomia ERRO/WARNING/SPAM vs ERROR/WARNING/INFO | **Vale João / ADR_170** para implementação Finder; **não emendar** ADR_019, ADR_167, SSID_007.1 nesta execução (FR-013) |
| **D9** | respostas-documentadas 010-21 (T1 A Finder grava) × revisao-T1-drive | Quem garante tela/e-mail | **Vale revisão T1**: Finder garante saída mascarada no ResumoAnalise (010-21) |
| **D10** | 009-5 A inicial (“só id”) × diagrama Notificação v2 (tipo e gravidade no alerta) | Payload mínimo do aviso ao Core | **Vale João, 009-5 revisto** (16/11 Etapa 3): **id + tipo + gravidade**, sem texto de log (respostas-etapa3.md §7) |
| **D11** | LLT-202 crit.1/2 (amplo) × DPA/ADR_064 (só linhas) | Atributos OTel sensíveis | **Vale o Drive** (DN-010-2): mascaramento LhamaLog **só linhas**; Cliente na origem (DPA Anexo IV item 6). **Ressalva ADR_040** para identificadores. **Alerta RIPD** (ampliar antes de dados reais) — já 010-12 |
| **D12** | E-mail LLT-218 × diagrama Notificação v2 (templates [PENDENTE]) | Quem monta corpo do e-mail | **Vale DN-010-1 A**: Notificação lê Ocorrência por id; lacuna de **template** permanece na spec Notificação (respostas-etapa3.md §7), não contradiz mascaramento no Finder |