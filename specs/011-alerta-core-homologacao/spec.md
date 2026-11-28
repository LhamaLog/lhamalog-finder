# Feature Specification: Alerta ao Core e homologação Finder

**Feature Branch**: `LLT-201-publicar-evento-de-ocorrencia-nova-ou-atualizada`

**Created**: 2026-11-26

**Status**: Draft

**Input**: User description: "Terminar, no lhamalog-finder, a parte da LLT-201 que falta: o cliente da rota de alerta já existe; alinhar o JWT ao segredo do #20 com audiência `lhamalog-core`; manter só `https://` e usar proxy no teste local; homologar ponta a ponta com o Core; registrar que o Módulo de Ocorrência ficou dentro do Core. Não reescrever a spec 009."

**Diretório da spec**: `specs/011-alerta-core-homologacao` (número 011). O nome desta pasta e o nome da branch Git são independentes. Nenhuma branch nova é criada por esta spec. (fonte: pedido desta especificação)

Esta spec complementa `specs/009-publicar-evento-ocorrencia`. Não a substitui e não a edita. O contrato DN-009-2 permanece o da 009. A rota de alerta no Core é a spec 016 do lhamalog-core. A pilha #19, #20, #21 e #24 é premissa de base: esta spec não recria o que ela faz. (fonte: decisoes-joao-llt201.md, base da branch e rota POST; origem: contexto-llt201.md, §1)

## Clarifications

### Session 2026-11-26 (rodada 2 — LLT-201)

- Q: Pendência 1 — Qual é o efeito observável, na homologação, depois do aceite durável do alerta? → A: Opção A. O Core grava o alerta e avisa a tela por WebSocket, no sentido da ADR_051, sem evento novo e sem envelope próprio para a Notificação. (fonte: decisoes-joao-llt201.md, pendência 1, formulário 1). A recomendação D do Especialista não vale. Conflito ADR_165 (ocorrências pela Notificação via stream) × "só o aviso do #24": vale o João. (fonte: respostas-pendentes-llt201.md, conflito 1)
- Q: Pendência 2 — Onde a gravidade aparece no que a homologação confere no Core? → A: Opção C. No documento da ocorrência e no aviso do #24. (fonte: respostas-pendentes-llt201.md, pendência 2; origem: Ocorrência.mmd (Drive, DC-10 §8) + 009-5 revisto). A spec 015, FR-007 (aviso sem gravidade) fica como emenda pendente da 015; esta execução não edita a 015. Vale o João. (fonte: respostas-pendentes-llt201.md, conflito 2)
- Q: Pendência 3 — Qual regra de gravidade a homologação confere? → A: Opção D. Gravidade = classe: o Core copia a classe ERROR, WARNING ou INFO para o campo `gravidade` e mantém o tipo. O Finder não calcula. (fonte: decisoes-joao-llt201.md, pendência 3, formulário 1). Conflito ADR_062 (escala "critico"/"alto", matriz pendente) × esta regra: vale o João. (fonte: respostas-pendentes-llt201.md, conflito 3)
- Q: Pendência 5 — Quantos avisos de ocorrência a homologação espera por entrega? → A: Opção B. Um aviso por ocorrência por entrega, com a mudança mais forte: nova > reaberta > virou recorrente. Criação já recorrente → só `CREATED`. FECHADA que reabre e fica recorrente na mesma entrega → só `REOPENED`. (fonte: decisoes-joao-llt201.md, pendência 5, formulário 1; elaboração alinhada à opção B: respostas-pendentes-llt201.md, pendência 5)
- Q: Pendência 6 — Onde o Core configura a audiência `lhamalog-core` que o token do Finder precisa satisfazer? → A: Opção B. Propriedade `core-audience`, usada só na rota de alerta. Ilustração de caminho, não nome extra decidido: `lhamalog.occurrence.internal-auth.core-audience`. As rotas do Módulo seguem com a audiência do #20. (fonte: decisoes-joao-llt201.md, pendência 6, formulário 1)
- Q: Pendência 13 — Qual proxy TLS o teste local usa? → A: Opção A. Caddy no compose do Core. O proxy é do lado do Core. (fonte: decisoes-joao-llt201.md, pendência 13, formulário 2)
- Q: Pendência 17 — Qual proxy, em qual compose, do ponto de vista do Finder? → A: Opção A. O proxy, o produto e o arquivo são os que a spec 016 registra (Caddy no compose do Core). O Finder só recebe a URL `https://`. (fonte: respostas-pendentes-llt201.md, pendência 17; origem: decisoes-joao-llt201.md, TLS no teste local + ADR_164 + arquitetura-modulos; alinhado à pendência 13 A, formulário 2)
- Q: Pendência 18 — Como o Finder confia no certificado local? → A: Opção B. A CA local entra só no teste, pela variável `SSL_CERT_FILE`, sem código novo no cliente e sem desligar a verificação. (fonte: decisoes-joao-llt201.md, pendência 18, formulário 2; respostas-pendentes-llt201.md, pendência 18; origem: o pacote `crypto/x509` do Go lê `SSL_CERT_FILE`)
- Q: Pendência 19 — O que o Finder faz se o POST do alerta falha? → A: Opção C. O lote permanece pendente e o aviso é reenviado no reprocessamento do lote. (fonte: respostas-pendentes-llt201.md, pendência 19; origem: DF-FINDER §2 e §3). Compatível com 009-18 A, mas a opção escolhida é C.
- Q: Pendência 20 — Os dois papéis usam tokens separados? → A: Opção B. Um token com as duas audiências (a do Módulo e `lhamalog-core`). `FINDER_MODULO_OCORRENCIA_JWT` e `FINDER_CORE_JWT` podem carregar o mesmo JWT, que deve conter as duas audiências. O cliente não passa a emitir nem a inspecionar claims. A cadeia `/internal/v1/**` do #20 precisa aceitar esse token na rota de alerta pela audiência `lhamalog-core` (propriedade `core-audience`, pendência 6) e nas rotas do Módulo pela audiência do Módulo. (fonte: decisoes-joao-llt201.md, pendência 20, formulário 1). A recomendação A do Especialista (dois tokens) não vale. Fato de código: no #20, `InternalServiceSecurityConfiguration` cobre todo `/internal/v1/**` e exige que o `aud` contenha `lhamalog.occurrence.internal-auth.audience` (padrão `lhamalog-modulo-ocorrencia`); o `aud` é conferido por "contém"; um token só com `lhamalog-core`, sem o ajuste por rota, recebe 401. `scripts/InternalJwt.java` gera um token com uma única audiência (`--aud`). (fonte: respostas-pendentes-llt201.md, resumo, fato de código, e pendência 20)

### Session 2026-11-27 (rodada 3 — LLT-201)

- Q: Qual é o nome da coleção Mongo do alerta no Core? → A: `finder_alerts`. Fecha a pendência 7. A opção A (coleção nova, chave `(tenantId, productId, batchId, motivo)`) já estava fechada. Não usar `AlertaFinder`. (fonte: decisoes-joao-llt201.md, Rodada 3)
- Q: Qual é o literal de `Change` para "virou recorrente"? → A: `RECURRENT`. Fecha a pendência 10. Distinto de `CREATED` e de `REOPENED`. Não usar `RECURRING`. (fonte: decisoes-joao-llt201.md, Rodada 3)
- Q: Qual é o aviso WebSocket do alerta, e em qual resposta ele sai? → A: Depois do aceite durável (primeira gravação, 201), canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`. A repetição (200) não avisa de novo. O efeito da pendência 1 A permanece: gravar o alerta e avisar a tela, sem evento novo de Notificação e sem publicar em `lhamalog:occurrence-events` nem em `lhamalog:notifications:events`. (fonte: decisoes-joao-llt201.md, Rodada 3)
- Q: HIGH U1 — Qual HTTP separa a primeira aceitação da repetição, e o corpo traz `repeticao`? → A: 201 na primeira aceitação e 200 na repetição da mesma identidade. O corpo de sucesso não traz o campo `repeticao`. (fonte: respostas-pendentes-llt201.md, pendência 8; origem: DF-FINDER §3 + controller do #20 + cliente estrito; analyze-011.md, HIGH U1)
- Q: HIGH U2 — Quando a homologação espera o aviso "virou recorrente"? → A: Quando `before.recurring` não é TRUE e `after.recurring` é TRUE; ausente ou `null` conta como `false`. O literal do enum é `RECURRENT`. A precedência normativa é a do FR-007. (fonte: respostas-pendentes-llt201.md, pendência 11; origem: data-model da spec 012; literal: decisoes-joao-llt201.md, Rodada 3; analyze-011.md, HIGH U2)
- Q: HIGH U3 — O que a homologação espera se o produto não pertence ao tenant? → A: 403 `ACESSO_NEGADO`, sem gravar o alerta. (fonte: decisoes-joao-llt201.md, pendência 12, formulário 1; elaboração: respostas-pendentes-llt201.md, pendência 12; origem: spec 013, FR-016; analyze-011.md, HIGH U3)
- Q: MEDIUM U4 — Qual código vale para corpo inválido e para JSON ilegível? → A: Corpo incompleto, campo a mais ou motivo fora da lista → 422 `DADOS_INVALIDOS`. JSON ilegível → 400 `REQUISICAO_INVALIDA`. (fonte: respostas-pendentes-llt201.md, pendência 4; origem: DF-FINDER §3 + padrão do #20; analyze-011.md, MEDIUM U4)
- Q: MEDIUM U5 — O que acontece em duas primeiras gravações da mesma chave ao mesmo tempo? → A: Quem perde a corrida relê e responde 200 com a mesma identidade. O cliente trata 200 e 201 como sucesso. (fonte: decisoes-joao-llt201.md, pendência 9, formulário 1; analyze-011.md, MEDIUM U5)
- Q: MEDIUM I1 — Quais tetos de tempo a homologação mede? → A: A rota de alerta responde dentro de 2 s (`FINDER_CORE_TIMEOUT`). O aviso "virou recorrente" sai em até 3 s depois do 201. Os dois tetos são o SC-009. O plano aponta para esse critério. (fonte: respostas-pendentes-llt201.md, pendência 16; origem: spec 015, SC-001 + DF-FINDER §3 e §8; analyze-011.md, MEDIUM I1)
- Q: MEDIUM I2 — O quickstart desta pasta já existe? → A: Sim. O roteiro está em `specs/011-alerta-core-homologacao/quickstart.md`. Não é arquivo futuro. (fonte: `specs/011-alerta-core-homologacao/quickstart.md`; plan.md; analyze-011.md, MEDIUM I2)
- Q: MEDIUM I3 — Onde ficam o pacote do alerta e o prefixo de log, do ponto de vista do Finder? → A: Fora do Finder. Pacote `br.com.lhamalog.core.finderalert`. Prefixo `occurrence.alert.*`. Não adotar `finder.alert.*`. Os sufixos exatos de `event=` permanecem [PENDENTE] no Core. (fonte: decisoes-joao-llt201.md, pendência 15, formulário 1, e pendência 14, formulário 2; sufixos: analyze-016.md, MEDIUM U2, SEM FONTE; analyze-011.md, MEDIUM I3)
- Q: LOW D1 — Onde a precedência é normativa? → A: No FR-007. As demais ocorrências desta spec só apontam para ele. (fonte: decisoes-joao-llt201.md, pendência 5, formulário 1; elaboração: respostas-pendentes-llt201.md, pendência 5; analyze-011.md, LOW D1)
- Q: LOW L2 — Qual algoritmo o FR-002 nomeia? → A: HS256, o algoritmo já usado no filtro do #20. (fonte: plan.md, item 1 de "O que a implementação cobre"; `contracts/alerta-finder-cliente.md`; analyze-011.md, LOW L2)
- Q: LOW L3 — Quais nomes de segredo o FR-002 registra? → A: `OCCURRENCE_INTERNAL_JWT_SECRET` (singular, lido por `scripts/InternalJwt.java`) e `OCCURRENCE_INTERNAL_JWT_SECRETS` (plural, propriedade do Core). Não unificar os nomes e não copiar o script. (fonte: tasks.md, T011; `quickstart.md`; analyze-011.md, LOW L3)
- Q: MEDIUM I4 — Onde fica a conferência do WebSocket nas tarefas? → A: Na história 3, sem a marca [P] que a soltava da configuração do token (FR-010, US3 cenário 5). (fonte: spec.md, FR-010 e US3 cenário 5; analyze-011.md, MEDIUM I4)
- Q: LOW L1 — Qual item do roteiro a tarefa do JWT do Módulo cita? → A: O item 5 (JWT), não o item 4 (URLs). (fonte: spec.md, roteiro itens 4 e 5; analyze-011.md, LOW L1)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - O JWT do alerta usa o segredo do #20 e as duas audiências (Priority: P1)

O Finder já chama `POST /internal/v1/alertas/finder` com o token configurado em `FINDER_CORE_JWT`. Esse token, e o de `FINDER_MODULO_OCORRENCIA_JWT`, podem ser o mesmo JWT: segredo e emissor já usados nas rotas internas do #20, com as duas audiências no `aud` (a do Módulo, padrão `lhamalog-modulo-ocorrencia`, e `lhamalog-core`). O cliente não é reescrito para emitir ou inspecionar o token.

**Why this priority**: sem esse token o Core recusa o alerta. A pendência 20 fecha um único token com as duas audiências, e a pendência 6 fecha a audiência `lhamalog-core` só na rota de alerta. (fonte: decisoes-joao-llt201.md, JWT da rota de alerta, pendências 6 e 20, formulário 1)

**Independent Test**: configurar as duas variáveis com o mesmo JWT (segredo e emissor do #20, `aud` contendo as duas audiências), chamar a rota de alerta e as rotas do Módulo e receber o aceite de cada uma. Dá para verificar sem percorrer os cinco cenários do DN-009-2.

**Acceptance Scenarios**:

1. **Given** o cliente de alerta já existente, **When** se lê o que ele envia, **Then** o pedido continua `POST /internal/v1/alertas/finder` com corpo `{tenantId, productId, batchId, motivo}`, o único motivo enviado é `consulta_openobserve_indisponivel`, e o sucesso que o cliente aceita é 201 na primeira aceitação e 200 na repetição da mesma identidade, com `{aceito: true, tenantId, productId, batchId, motivo}` iguais ao pedido e sem o campo `repeticao`. (fonte: contexto-llt201.md, §3 e R9; origem: `internal/adapter/out/core/aviso.go`; respostas-pendentes-llt201.md, pendência 8; origem: DF-FINDER §3 + controller do #20 + cliente estrito)
2. **Given** as variáveis `FINDER_CORE_JWT` e `FINDER_MODULO_OCORRENCIA_JWT`, **When** o ambiente de homologação é configurado, **Then** as duas podem carregar o mesmo JWT HS256, emitido com o segredo do #20 (`lhamalog.occurrence.internal-auth.secrets`, ambiente `OCCURRENCE_INTERNAL_JWT_SECRETS`, plural) e com o emissor do #20 (padrão `lhamalog-internal`), e o `aud` desse JWT contém a audiência do Módulo e `lhamalog-core`. O script `scripts/InternalJwt.java` lê `OCCURRENCE_INTERNAL_JWT_SECRET` (singular); os dois nomes não se unificam e o script não é copiado. (fonte: decisoes-joao-llt201.md, pendência 20, formulário 1; origem: `application.properties` na branch `LLT-293-gravar-ocorrencias-a-partir-do-resumo-de-analise-do-finder-modulo-ocorrencia`; plan.md, item 1; tasks.md, T011)
3. **Given** as rotas de consulta e de resumo do #20 e a rota de alerta, **When** o Finder as chama com esse JWT, **Then** a audiência do Módulo (`lhamalog-modulo-ocorrencia`) continua a exigida nas rotas do Módulo, e a audiência `lhamalog-core` é a exigida na rota de alerta, lida pelo Core da propriedade `core-audience` (ilustração de caminho: `lhamalog.occurrence.internal-auth.core-audience`), usada só nessa rota. A audiência nova não substitui a audiência das rotas do Módulo. (fonte: decisoes-joao-llt201.md, pendências 6 e 20, formulário 1)
4. **Given** o cliente atual, **When** se confere se ele monta ou valida `iss`, `aud` ou `sub`, **Then** ele não faz isso: envia o valor de `FINDER_CORE_JWT` como `Authorization: Bearer` na rota de alerta e o valor de `FINDER_MODULO_OCORRENCIA_JWT` nas rotas do Módulo. Não há mudança de código do cliente para claims. O fato de o #20 conferir `aud` por "contém", e de `scripts/InternalJwt.java` gerar uma audiência por `--aud`, fica no Core; esta spec não altera o script nem o cliente. (fonte: contexto-llt201.md, §3; origem: `internal/adapter/out/confirmacao/http.go` e `internal/bootstrap/pipeline.go`; respostas-pendentes-llt201.md, resumo, fato de código)

---

### User Story 2 - O teste local usa HTTPS pelo Caddy do Core, sem liberar HTTP (Priority: P1)

No teste local o Core responde em HTTP na porta 8080. O Finder só chama base `https://`. O teste local usa o Caddy no compose do Core, na frente do Core. O cliente do Finder não passa a aceitar `http://`. A CA local desse proxy entra só no teste, por `SSL_CERT_FILE`.

**Why this priority**: a decisão escolhe o Caddy no compose do Core e recusa liberar HTTP no cliente. Sem a URL `https://` a homologação local não chega ao Core. (fonte: decisoes-joao-llt201.md, TLS no teste local, pendência 13, formulário 2; respostas-pendentes-llt201.md, pendência 17)

**Independent Test**: apontar o Finder para a URL `https://` do Caddy, com a CA local em `SSL_CERT_FILE`, e concluir uma chamada de alerta e uma entrega de resumo. Uma base `http://` é recusada pelo Finder antes de qualquer chamada. A verificação do certificado permanece ligada.

**Acceptance Scenarios**:

1. **Given** a base do Core ou a base do Módulo, **When** o Finder a usa, **Then** ela é `https://`, sem usuário, sem consulta, sem fragmento e sem caminho além de `/`. Base `http://` não é aceita. (fonte: contexto-llt201.md, §3 e §5 conflito 7; origem: `internal/adapter/out/confirmacao/http.go` e `internal/adapter/out/modocorrencia/cliente.go`)
2. **Given** o teste local, **When** o Finder aponta para o Core, **Then** a URL é a `https://` do Caddy que a spec 016 registra no compose do Core, e o Core local segue em HTTP na porta 8080. O Finder só recebe essa URL. O proxy, o produto e o arquivo são os do Core. Esta spec não altera compose. (fonte: decisoes-joao-llt201.md, pendência 13, formulário 2; respostas-pendentes-llt201.md, pendência 17; origem: ADR_164 + arquitetura-modulos)
3. **Given** o Módulo de Ocorrência dentro do Core, **When** a homologação configura as duas bases, **Then** `FINDER_MODULO_OCORRENCIA_URL` e `FINDER_CORE_URL` apontam para o mesmo host. (fonte: contexto-llt201.md, §8, quickstart do Finder)
4. **Given** o cliente de confirmação, **When** a base responde com redirecionamento, **Then** o cliente não segue o redirecionamento. A URL configurada é o próprio ponto HTTPS do Caddy. (fonte: `internal/adapter/out/confirmacao/http.go`, `CheckRedirect`)
5. **Given** o teste local com a CA do Caddy, **When** o processo do Finder sobe, **Then** a confiança nessa CA vale só nesse teste, pela variável `SSL_CERT_FILE`, sem código novo no cliente e sem desligar a verificação HTTPS. (fonte: decisoes-joao-llt201.md, pendência 18, formulário 2; respostas-pendentes-llt201.md, pendência 18)

---

### User Story 3 - Homologar o Finder com o Core de verdade (Priority: P1)

O cliente que consulta o achado e entrega o `ResumoAnalise` já está na Development (spec 009, PR #13). Esta história é o aceite com o servidor real do Core, depois do merge da pilha. O Finder não cria ocorrência, não calcula gravidade e não publica o evento de ocorrência. A homologação confere no Core o efeito do alerta, a gravidade e o aviso.

**Why this priority**: o contexto registra que falta a homologação com servidor real, inclusive os cenários DN-009-2 e os resultados citados como SC-002 e SC-006 da spec 009. (fonte: contexto-llt201.md, §1)

**Independent Test**: subir Finder e Core, com o Caddy no compose do Core e `SSL_CERT_FILE` no teste, e percorrer os cinco cenários do DN-009-2 usando só as rotas que o #20 já entrega. A conferência de ocorrência nova, reaberta e "virou recorrente" segue o FR-007. A do alerta é a gravação em `finder_alerts` e o aviso no canal `/ws/occurrences`, tipo `FINDER_ALERT`, só depois do 201. O Finder entra como quem entrega o resumo e quem envia o alerta.

**Acceptance Scenarios**:

1. **Given** as rotas `POST /internal/v1/ocorrencias/consultas` e `POST /internal/v1/resumos-analise` já entregues pelo #20, **When** a homologação roda, **Then** o Finder usa o cliente já existente. Esta spec não cria outra versão do contrato DN-009-2 nem outra implementação dessas chamadas. (fonte: contexto-llt201.md, §1 e §3; origem: `proposta-dn-009-2.md`; spec 009, status Draft com cliente implementado)
2. **Given** o contrato DN-009-2, **When** a homologação ponta a ponta termina, **Then** os cinco cenários abaixo passaram com o cliente do Finder e o Core:
   1. A consulta distingue ocorrência ausente de falha e rejeita acesso a outro tenant ou produto.
   2. A primeira entrega grava exatamente um efeito e retorna 201. Reenvio depois de perda de resposta, inclusive depois de reinício do Finder, retorna 200 e o mesmo `reciboId`, sem duplicar.
   3. Dois lotes que consultam a mesma revisão preservam todos os incrementos. Revisão desatualizada em um achado rejeita a entrega inteira antes de qualquer escrita.
   4. Confirmação inválida ou inesperada mantém o lote pendente no Finder. Um recibo válido antecede a conclusão do lote e o reconhecimento da fila.
   5. O JSON mantém as três classes, carrega somente trechos protegidos e conserva `linhas` vazias quando for o caso. A API rejeita tipos e contadores inválidos.
   (fonte: contexto-llt201.md, §1; origem: `specs/009-publicar-evento-ocorrencia/contracts/proposta-dn-009-2.md`, seção "Cenários de homologação com o Módulo")
3. **Given** uma entrega que o Core trata como ocorrência nova, como reaberta ou como "virou recorrente", **When** a homologação confere o resultado no Core, **Then** o Finder apenas entregou o `ResumoAnalise` pelo cliente já existente, e o Core segue o FR-007 (precedência, gatilho de `before.recurring` / `after.recurring` e literal `RECURRENT`). O Finder não publica evento de ocorrência e não calcula gravidade. (fonte: spec.md, FR-007; decisoes-joao-llt201.md, pendências 3 e 5, formulário 1, e Rodada 3; respostas-pendentes-llt201.md, pendências 2, 5 e 11)
4. **Given** o teste local com o Caddy, **When** a homologação inclui o alerta e a entrega, **Then** as duas chamadas usam URL `https://` no mesmo host, e a CA local está só em `SSL_CERT_FILE`. (fonte: decisoes-joao-llt201.md, pendências 13 e 18, formulário 2; respostas-pendentes-llt201.md, pendência 17)
5. **Given** um alerta aceito pelo Core pela primeira vez (201), **When** a homologação observa o efeito, **Then** o Core gravou um documento na coleção `finder_alerts` e avisou a tela no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, sem evento novo, sem envelope próprio para a Notificação e sem publicar em `lhamalog:occurrence-events` nem em `lhamalog:notifications:events`. A repetição da mesma identidade (200) não avisa de novo. O papel do Finder foi o POST. (fonte: decisoes-joao-llt201.md, pendência 1, formulário 1, e Rodada 3)
6. **Given** o POST do alerta falha, **When** o Finder trata a resposta, **Then** o lote permanece pendente e o aviso é reenviado no reprocessamento desse lote, com a mesma identidade. (fonte: respostas-pendentes-llt201.md, pendência 19; origem: DF-FINDER §2 e §3)

---

### Edge Cases

- Reenvio do mesmo alerta: 201 na primeira aceitação e 200 na repetição, mesma identidade e `aceito: true`, sem o campo `repeticao`. A deduplicação durável é do Core, na coleção `finder_alerts`. (fonte: contexto-llt201.md, R9; origem: `aviso.go`; respostas-pendentes-llt201.md, pendência 8; decisoes-joao-llt201.md, Rodada 3)
- Duas primeiras gravações da mesma chave ao mesmo tempo: quem perde a corrida relê e responde 200 com a mesma identidade. O cliente trata 200 e 201 como sucesso. (fonte: decisoes-joao-llt201.md, pendência 9, formulário 1)
- Corpo incompleto, campo a mais no pedido ou motivo fora da lista: o Core responde 422 `DADOS_INVALIDOS` e não grava. JSON ilegível: 400 `REQUISICAO_INVALIDA`. O cliente trata esses status como falha do POST. (fonte: respostas-pendentes-llt201.md, pendência 4; origem: DF-FINDER §3 + padrão do #20)
- Produto que não pertence ao tenant: 403 `ACESSO_NEGADO`, sem gravar o alerta. O lote permanece pendente. (fonte: decisoes-joao-llt201.md, pendência 12, formulário 1; elaboração: respostas-pendentes-llt201.md, pendência 12; origem: spec 013, FR-016)
- Corpo de sucesso com campo a mais, nome diferente ou identidade diferente: o cliente não considera aceite. (fonte: contexto-llt201.md, §3; origem: `aviso.go` e `aviso_test.go`)
- Resposta 202 ou corpo vazio no alerta: o cliente não considera aceite. (fonte: `aviso_test.go`)
- Falha da consulta ao OpenObserve: até 3 tentativas na consulta; depois uma chamada de alerta. Se essa chamada falha, o lote permanece pendente, o código não marca `recorrente=false` por degradação, e o aviso é reenviado no reprocessamento do lote. Não se acrescenta um número novo de tentativas dentro do POST. (fonte: respostas-pendentes-llt201.md, pendência 19; origem: DF-FINDER §2 e §3; `internal/adapter/out/openobserve/recorrencia.go`; spec 008, DN-008-3 revisto)
- As duas URLs em hosts diferentes: o cliente hoje aceita. A configuração correta da homologação é o mesmo host. Não há fonte para o cliente passar a recusar hosts diferentes. (fonte: contexto-llt201.md, §8)
- Gravidade, "virou recorrente" e o aviso à Notificação observados na homologação vêm do Core e seguem o FR-007. O Finder não calcula gravidade nem publica esse aviso. O alerta aceito, à parte, na primeira gravação (201), grava em `finder_alerts` e avisa a tela no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, sem evento novo para a Notificação. A repetição (200) não avisa de novo. (fonte: spec.md, FR-007 e FR-010; decisoes-joao-llt201.md, pendências 1, 3 e 5, formulário 1, e Rodada 3; respostas-pendentes-llt201.md, pendência 2)
- Token só com a audiência do Módulo: não autentica a rota de alerta. Token só com `lhamalog-core`: sem o ajuste por rota da pendência 6, a cadeia `/internal/v1/**` do #20 responde 401, porque exige que o `aud` contenha a audiência do Módulo. O JWT da homologação contém as duas. (fonte: respostas-pendentes-llt201.md, resumo, fato de código; decisoes-joao-llt201.md, pendências 6 e 20, formulário 1)
- Certificado local fora do teste, ou verificação desligada: não é o caminho adotado. A CA local vale só no teste, via `SSL_CERT_FILE`, com a verificação padrão ligada. (fonte: decisoes-joao-llt201.md, pendência 18, formulário 2)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: O Finder MUST continuar chamando `POST /internal/v1/alertas/finder` com o corpo `{tenantId, productId, batchId, motivo}` e MUST tratar como sucesso 201 na primeira aceitação e 200 na repetição da mesma identidade, com `{aceito: true}` e a mesma identidade, sem o campo `repeticao`. O motivo que o Finder envia MUST continuar `consulta_openobserve_indisponivel`. Corpo incompleto, campo a mais ou motivo fora da lista MUST ser, no Core, 422 `DADOS_INVALIDOS`. JSON ilegível MUST ser 400 `REQUISICAO_INVALIDA`. Produto que não pertence ao tenant MUST ser 403 `ACESSO_NEGADO`, sem gravar o alerta. Quem perde a corrida de duas primeiras gravações da mesma chave MUST reler e responder 200 com a mesma identidade; o cliente MUST tratar 200 e 201 como sucesso quando a identidade confere. Esta spec MUST NOT reimplementar esse cliente nem criar esses códigos. (fonte: contexto-llt201.md, §3 e R9; origem: `aviso.go`; respostas-pendentes-llt201.md, pendências 4 e 8; origem: DF-FINDER §3 + padrão do #20 + controller do #20 + cliente estrito; decisoes-joao-llt201.md, pendências 9 e 12, formulário 1; elaboração do 403: respostas-pendentes-llt201.md, pendência 12; origem: spec 013, FR-016)
- **FR-002**: O JWT colocado em `FINDER_CORE_JWT` MUST ser HS256, o algoritmo já usado no filtro do #20, emitido com o segredo e o emissor do #20, e MUST conter no `aud` a audiência do Módulo e `lhamalog-core`. `FINDER_MODULO_OCORRENCIA_JWT` MAY carregar o mesmo JWT. O script `scripts/InternalJwt.java` lê `OCCURRENCE_INTERNAL_JWT_SECRET` (singular). A propriedade do Core é `OCCURRENCE_INTERNAL_JWT_SECRETS` (plural). Esta spec MUST NOT unificar esses nomes e MUST NOT copiar o script. O Finder MUST NOT passar a emitir nem a inspecionar esse JWT no código. Segredo MUST NOT entrar no repositório. Esta decisão não muda o cliente além de aceitar, nas duas variáveis, esse mesmo valor. (fonte: decisoes-joao-llt201.md, pendência 20, formulário 1; constituição, princípio V, incerteza sem invenção — o segredo não é dado desta spec; plan.md, item 1 de "O que a implementação cobre"; tasks.md, T011; `quickstart.md`)
- **FR-003**: O token colocado em `FINDER_MODULO_OCORRENCIA_JWT` MUST continuar aceito nas rotas `POST /internal/v1/ocorrencias/consultas` e `POST /internal/v1/resumos-analise`, cuja audiência no #20 é `lhamalog-modulo-ocorrencia`. Esta spec MUST NOT trocar essa audiência. A rota de alerta MUST exigir `lhamalog-core`, lida pelo Core da propriedade `core-audience`, usada só nessa rota (ilustração de caminho: `lhamalog.occurrence.internal-auth.core-audience`). A cadeia `/internal/v1/**` do #20 MUST aceitar o token das duas audiências na rota de alerta por essa propriedade e nas rotas do Módulo pela audiência do Módulo. Fato: hoje `InternalServiceSecurityConfiguration` cobre todo `/internal/v1/**` e exige que o `aud` contenha a audiência do Módulo; o `aud` é conferido por "contém"; um token só com `lhamalog-core`, sem o ajuste por rota, recebe 401. `scripts/InternalJwt.java` gera uma audiência por `--aud`. Esta spec MUST NOT alterar esse script nem o cliente. (fonte: decisoes-joao-llt201.md, pendências 6 e 20, formulário 1; respostas-pendentes-llt201.md, resumo, fato de código)
- **FR-004**: `FINDER_CORE_URL` e `FINDER_MODULO_OCORRENCIA_URL` MUST ser `https://`, nas restrições de URL que o cliente já aplica. O Finder MUST NOT aceitar `http://` nessas bases. (fonte: decisoes-joao-llt201.md, TLS no teste local; origem: contexto-llt201.md, §7 Q7 opção A, que recusa liberar HTTP no cliente; `confirmacao/http.go` e `modocorrencia/cliente.go`)
- **FR-005**: No teste local, as duas URLs MUST apontar para a URL `https://` do Caddy no compose do Core (spec 016), na frente do Core. O Core local permanece em HTTP na porta 8080. As duas bases MUST usar o mesmo host, porque o Módulo é pacote dentro do Core. O Finder só recebe essa URL. O produto é Caddy e o arquivo de compose é o do Core. Esta spec MUST NOT editar compose. A CA local MUST ser confiável só no teste, via `SSL_CERT_FILE`, sem código novo no cliente e sem desligar a verificação. (fonte: decisoes-joao-llt201.md, pendências 13 e 18, formulário 2; respostas-pendentes-llt201.md, pendência 17; origem: ADR_164 + arquitetura-modulos)
- **FR-006**: A homologação MUST percorrer os cinco cenários do DN-009-2 citados na User Story 3, com o cliente já existente e o Core depois do merge da pilha. MUST NOT duplicar consulta nem entrega do `ResumoAnalise`. (fonte: contexto-llt201.md, §1; origem: `proposta-dn-009-2.md`)
- **FR-007**: Na homologação de ocorrência nova, reaberta e "virou recorrente", o Finder MUST só entregar o resumo. MUST NOT criar ocorrência, MUST NOT calcular gravidade e MUST NOT publicar evento de ocorrência. Este requisito é o texto normativo da precedência. No Core (spec 016), a homologação MUST confer um aviso por ocorrência por entrega, com a mudança mais forte (nova > reaberta > virou recorrente). Criação já recorrente MUST publicar só `CREATED`. FECHADA que reabre e fica recorrente na mesma entrega MUST publicar só `REOPENED`. O aviso "virou recorrente" MUST disparar quando `before.recurring` não é TRUE e `after.recurring` é TRUE; ausente ou `null` conta como `false`; essa mudança só sai se for a mais forte da entrega. O literal dessa mudança no enum `Change` MUST ser `RECURRENT`, distinto de `CREATED` e de `REOPENED`. A gravidade MUST ser igual à classe ERROR, WARNING ou INFO no campo `gravidade`, no documento e no aviso do #24, com o tipo mantido. O aviso de ocorrência reutilizado é o do #24. As demais menções desta spec à precedência apontam para este requisito. (fonte: decisoes-joao-llt201.md, pendências 3 e 5, formulário 1, e Rodada 3; respostas-pendentes-llt201.md, pendências 2, 5 e 11; origem: Ocorrência.mmd + 009-5 revisto; data-model da spec 012; constituição do Core, princípio IV)
- **FR-008**: Esta spec MUST registrar, só nesta pasta, que o Módulo de Ocorrência é o pacote `br.com.lhamalog.core.occurrence` dentro do Core e que o aviso HTTP do Módulo para o Core descrito na spec 009 deixa de ser o fluxo do processo único. MUST NOT editar a spec 009 nem a constituição. Onde esses textos divergirem, a pendência fica na seção Divergências. (fonte: contexto-llt201.md, §1, §5 conflito 2 e §8; pedido desta especificação)
- **FR-009**: Se o POST do alerta falha, o lote MUST permanecer pendente e o aviso MUST ser reenviado no reprocessamento do lote, com a mesma identidade, apoiado na deduplicação do Core. MUST NOT acrescentar um número novo de tentativas dentro desse POST. A opção escolhida é C. Ela é compatível com 009-18 A, e 009-18 A não é a regra aplicada a este POST. Corpo e identidade de sucesso permanecem os do cliente e do R9. (fonte: respostas-pendentes-llt201.md, pendência 19; origem: DF-FINDER §2 e §3)
- **FR-010**: Depois do aceite durável do alerta (primeira gravação, 201), a homologação MUST confer no Core um documento na coleção `finder_alerts` e o aviso à tela no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, sem evento novo, sem envelope próprio para a Notificação e sem publicar em `lhamalog:occurrence-events` nem em `lhamalog:notifications:events`. A repetição (200) MUST NOT avisar de novo. O Finder MUST NOT publicar esse aviso. (fonte: decisoes-joao-llt201.md, pendência 1, formulário 1, e Rodada 3)

### Key Entities *(include if feature involves data)*

- **Alerta operacional do Finder**: chamada já existente depois que a consulta ao OpenObserve esgota as tentativas. Identidade: tenant, produto, lote (`batchId`) e motivo `consulta_openobserve_indisponivel`. Não é o evento de ocorrência da LLT-201. A gravação durável é a coleção `finder_alerts` no Core. Se o POST falha, o lote fica pendente e o aviso volta no reprocessamento. (fonte: contexto-llt201.md, §3 e R9; respostas-pendentes-llt201.md, pendência 19; decisoes-joao-llt201.md, Rodada 3)
- **Credencial do alerta e da consulta**: um JWT HS256 de ambiente, segredo e emissor do #20, `aud` com a audiência do Módulo e `lhamalog-core`. `FINDER_CORE_JWT` e `FINDER_MODULO_OCORRENCIA_JWT` podem carregar esse mesmo valor. O script lê `OCCURRENCE_INTERNAL_JWT_SECRET` (singular); a propriedade do Core é `OCCURRENCE_INTERNAL_JWT_SECRETS` (plural). A rota de alerta exige `lhamalog-core` via `core-audience`; as rotas do Módulo exigem a audiência do #20. (fonte: decisoes-joao-llt201.md, pendências 6 e 20, formulário 1; plan.md, item 1; tasks.md, T011)
- **ResumoAnalise e recibo**: os já definidos no DN-009-2. Esta spec não cria campos novos. (fonte: contexto-llt201.md, §3; origem: `proposta-dn-009-2.md`)
- **Módulo de Ocorrência**: pacote dentro do Core, não serviço HTTP separado do Core. As duas URLs do Finder apontam para o mesmo host. (fonte: contexto-llt201.md, §1 e §8)
- **Gravidade conferida na homologação**: campo `gravidade` igual à classe ERROR, WARNING ou INFO, no documento da ocorrência e no aviso do #24, com o tipo mantido. Quem calcula é o Core. (fonte: decisoes-joao-llt201.md, pendência 3, formulário 1; respostas-pendentes-llt201.md, pendência 2)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Em 100% das chamadas de alerta e das chamadas de consulta ou resumo da homologação cujo JWT tem o segredo do #20, o emissor do #20 e as duas audiências no `aud`, o Finder recebe o aceite de cada rota. As duas variáveis podem carregar esse mesmo JWT. Um token cuja audiência é só `lhamalog-modulo-ocorrencia` não autentica a rota de alerta. Um token cuja audiência é só `lhamalog-core` não é a credencial das rotas do Módulo.
- **SC-002**: Em 100% das tentativas de configurar base `http://` para o Core ou para o Módulo, o Finder recusa a base e não faz a chamada. Nenhuma opção de liberar HTTP local é aceita.
- **SC-003**: No teste local, o Finder conclui uma chamada de alerta e uma entrega de resumo usando URL `https://` do mesmo host, através do Caddy no compose do Core, com a CA local só em `SSL_CERT_FILE` e com a verificação do certificado ligada.
- **SC-004**: Os cinco cenários do DN-009-2 passam uma vez, ponta a ponta, com o cliente já existente do Finder e o Core da pilha mergeada. As rotas usadas são as do #20.
- **SC-005**: Em 100% das homologações de ocorrência nova, reaberta e "virou recorrente", o Finder não publica evento de ocorrência e não calcula gravidade. No Core, a conferência segue o FR-007, inclusive o gatilho (`before.recurring` não é TRUE e `after.recurring` é TRUE; ausente ou `null` conta como `false`) e o literal `RECURRENT`.
- **SC-006**: A spec 009 e a constituição permanecem sem edição nesta feature. A divergência entre elas e o Módulo dentro do Core está registrada na seção Divergências, com os documentos a atualizar listados e não escritos aqui.
- **SC-007**: Em 100% das falhas do POST do alerta na homologação, o lote permanece pendente e o aviso é reenviado no reprocessamento do lote. O POST falho não confirma o lote.
- **SC-008**: Em 100% das primeiras gravações de alerta (201) na homologação, o Core gravou um documento em `finder_alerts` e avisou a tela no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, sem evento novo e sem envelope próprio para a Notificação. A repetição (200) não gera outro aviso.
- **SC-009**: A rota de alerta responde dentro de 2 s, o teto de `FINDER_CORE_TIMEOUT`. O aviso "virou recorrente" (`RECURRENT`, FR-007) sai em até 3 s depois do 201 da entrega. Os dois tetos valem em alvos diferentes. O Finder não implementa o teto de 3 s. (fonte: respostas-pendentes-llt201.md, pendência 16; origem: spec 015, SC-001 + DF-FINDER §3 e §8)

Os identificadores SC-002 e SC-006 desta spec não são os SC-002 e SC-006 da spec 009. Os da 009 aparecem na User Story 3 e na seção Divergências com o nome da pasta 009.

## Assumptions

Não há premissa inventada. O que segue é fato das fontes ou dependência já decidida.

- A implementação espera o merge de #19, #20, #21 e #24. Esta spec não descreve de novo lista, detalhe, consultas, resumos, ciclo de vida nem o canal `lhamalog:occurrence-events`. (fonte: decisoes-joao-llt201.md, base da branch)
- O cliente de alerta e o cliente DN-009-2 já estão na Development do Finder (PR #13, merge em 23/11). A spec 009 continua Draft. (fonte: contexto-llt201.md, §1 e §3)
- O `sub` que o #20 autoriza por padrão é `lhamalog-finder`. A decisão de 26/11 fala em segredo, emissor e audiências; não cria outro `sub`. Para o Core aceitar o token, o `sub` precisa estar na lista que o #20 já configura. (fonte: `application.properties` do #20, `subjects` padrão `lhamalog-finder`; contexto-llt201.md, §2)
- O cliente não segue redirecionamento e, sem transporte substituído, usa a verificação HTTPS padrão do processo. Teste automatizado pode injetar transporte. Não há variável de certificado própria do Finder. No teste local, a CA do Caddy entra por `SSL_CERT_FILE`, que o `crypto/x509` do Go já lê, sem código novo. (fonte: `confirmacao/http.go`; decisoes-joao-llt201.md, pendência 18, formulário 2)
- `FINDER_CORE_TIMEOUT` já existe, com padrão de 2 s. Esta spec não muda esse prazo. Esse teto, junto com os 3 s do aviso "virou recorrente", é o SC-009. (fonte: contexto-llt201.md, §3; respostas-pendentes-llt201.md, pendência 16)
- A decisão de 26/11 prevalece sobre o texto da spec 009 que ainda descreve aviso HTTP do Módulo para o Core. O registro da divergência está abaixo; a 009 não é editada aqui. (fonte: decisoes-joao-llt201.md, evento para a Notificação; contexto-llt201.md, §5 conflito 2)
- O efeito do alerta, a gravidade e a precedência do aviso de ocorrência são os que a spec 016 registra. Esta spec só os confere na homologação. A coleção é `finder_alerts`, o literal é `RECURRENT` e o aviso do alerta, depois do 201, é o do FR-010. (fonte: spec 016 do lhamalog-core, lida e não editada; decisoes-joao-llt201.md, formulários 1 e 2 e Rodada 3)

## Responsabilidades (registro só nesta pasta)

| Componente | Onde está | O que é dele neste fluxo |
|---|---|---|
| Finder | Este repositório | Entregar `ResumoAnalise` e consultar achado pelo cliente já existente. Avisar falha esgotada do OpenObserve em `POST /internal/v1/alertas/finder`. Se esse POST falha, manter o lote pendente e reenviar o aviso no reprocessamento. Configurar as duas variáveis de JWT, que podem carregar o mesmo token com as duas audiências, e as duas URLs `https://` no mesmo host. No teste, confiar na CA local só via `SSL_CERT_FILE`. Participar da homologação. Não criar ocorrência, não calcular gravidade, não publicar evento de ocorrência, não emitir nem inspecionar o JWT. |
| Módulo de Ocorrência | Pacote `br.com.lhamalog.core.occurrence` dentro do Core | Receber consulta e resumo nas rotas do #20, com a audiência do Módulo. Gravar ocorrência. Copiar a classe para o campo `gravidade` e manter o tipo. Publicar o aviso do #24 conforme o FR-007. |
| Core | Spec 016 do lhamalog-core | Atender a rota de alerta com a propriedade `core-audience` só nessa rota. Gravar o alerta na coleção `finder_alerts`. Depois do 201, avisar a tela no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, sem evento novo para a Notificação. A repetição (200) não avisa de novo. O pacote do alerta é `br.com.lhamalog.core.finderalert` e o prefixo de log é `occurrence.alert.*`; os sufixos de `event=` seguem [PENDENTE]. O Caddy no compose do Core termina o TLS local. (fonte: decisoes-joao-llt201.md, Rodada 3, pendência 14, formulário 2, e pendência 15, formulário 1) |
| Notificação | Fora deste repositório | Não é alcançada por um HTTP do Módulo para o Core nem por evento novo desta feature. O aviso de ocorrência reutilizado é o do #24. O alerta aceito, no 201, avisa no canal `/ws/occurrences` com tipo `FINDER_ALERT`. (fonte: decisoes-joao-llt201.md, evento para a Notificação, pendência 1, formulário 1, e Rodada 3) |

O desenho antigo Módulo → HTTP com JWT → Core → stream da Notificação, escrito em `specs/009-publicar-evento-ocorrencia/contracts/pipeline-aviso-referencia.md`, deixa de descrever o processo único. Módulo e Core são o mesmo processo. (fonte: contexto-llt201.md, §5 conflito 2 e §3)

## Roteiro de homologação

O arquivo deste roteiro já está em `specs/011-alerta-core-homologacao/quickstart.md`. O texto abaixo é o mesmo conteúdo, para leitura nesta spec. (fonte: `specs/011-alerta-core-homologacao/quickstart.md`; plan.md)

1. Core com a pilha #19–#24 mergeada e com a rota de alerta da spec 016 disponível. A rota de alerta lê `lhamalog-core` da propriedade `core-audience`, usada só nela. As rotas do Módulo seguem com a audiência do #20.
2. Caddy no compose do Core, na frente do Core local (HTTP na porta 8080). O Finder só recebe a URL `https://`. Esta spec não escreve o compose.
3. No teste, CA local só por `SSL_CERT_FILE`, com a verificação HTTPS ligada e sem código novo no cliente.
4. `FINDER_MODULO_OCORRENCIA_URL` e `FINDER_CORE_URL` iguais a `https://` desse Caddy, mesmo host, sem caminho na base.
5. `FINDER_MODULO_OCORRENCIA_JWT` e `FINDER_CORE_JWT` podem carregar o mesmo JWT, emitido com o segredo e o emissor do #20, cujo `aud` contém a audiência do Módulo e `lhamalog-core`. Fato: `scripts/InternalJwt.java` gera uma audiência por `--aud`; esta spec não altera o script nem manda o cliente inspecionar o token.
6. `sub` na lista do #20 (padrão `lhamalog-finder`). Segredo só no ambiente.
7. Confirmar que base `http://` é recusada.
8. Concluir uma chamada de alerta e uma entrega de resumo pela URL `https://`.
9. Percorrer os cinco cenários do DN-009-2 com o cliente já existente.
10. Conferir no Core o que o FR-007 exige: a precedência, o gatilho de "virou recorrente" e o literal `RECURRENT`. O Finder não publica esse evento e não calcula gravidade.
11. Conferir que, na primeira gravação do alerta (201), o Core gravou um documento na coleção `finder_alerts` e avisou a tela no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, sem evento novo para a Notificação e sem publicar em `lhamalog:occurrence-events` nem em `lhamalog:notifications:events`. A repetição (200) não avisa de novo. (fonte: decisoes-joao-llt201.md, Rodada 3)
12. Se o POST do alerta falha, o lote fica pendente e o aviso volta no reprocessamento do lote.

## Fora de escopo

- Reimplementar consulta ou entrega do `ResumoAnalise` (spec 009, PR #13).
- Fazer o Finder criar ocorrência, calcular gravidade ou publicar evento de ocorrência.
- Implementar `POST /internal/v1/alertas/finder` no Core (spec 016).
- Duplicar a pilha #19, #20, #21 e #24.
- Editar constituição, spec 009, spec 015, README, código, compose, ADRs, `.env.example` ou documentos do Drive e do lhamalog-docs.
- Criar evento novo ou envelope próprio para a Notificação.
- Liberar `http://` no cliente.
- Desligar a verificação TLS ou criar transporte próprio para a CA local.
- Fazer o cliente emitir ou inspecionar `iss`, `aud` ou `sub` para acomodar as duas audiências.
- Alterar `scripts/InternalJwt.java` ou unificar `OCCURRENCE_INTERNAL_JWT_SECRET` (singular, script) com `OCCURRENCE_INTERNAL_JWT_SECRETS` (plural, propriedade do Core).
- Criar, neste repositório, a coleção `finder_alerts`, o literal `RECURRENT`, o pacote `br.com.lhamalog.core.finderalert` ou o log com prefixo `occurrence.alert.*`. Esses nomes são do Core. Não adotar o prefixo `finder.alert.*`. Os sufixos de `event=` permanecem [PENDENTE] no Core. (fonte: decisoes-joao-llt201.md, Rodada 3, pendência 14, formulário 2, e pendência 15, formulário 1; sufixos: analyze-016.md, MEDIUM U2, SEM FONTE)

## Divergências com a constituição e com a spec 009

Estes textos não são alterados aqui. A decisão de 26/11 e o contexto desta rodada prevalecem sobre eles para o trabalho desta spec.

| Onde | O que o texto ainda diz | O que vale para esta spec | Pendência |
|---|---|---|---|
| Constituição, princípio I | O Finder transforma o resultado em ocorrência. | A ocorrência nasce no Módulo, que agora é pacote do Core. O Finder entrega `ResumoAnalise`. (fonte: contexto-llt201.md, R2; decisoes-joao-llt201.md, registro de que a Ocorrência fica dentro do Core) | Emenda da constituição, fora desta pasta. A spec 009 já deixa o item 10 das questões em aberto para essa emenda. |
| Spec 009, tabela de responsabilidades e `pipeline-aviso-referencia.md` | Módulo fora deste repositório; aviso Módulo → Core por HTTP com JWT; em seguida stream para a Notificação. | Módulo e Core são o mesmo processo. Esse HTTP deixa de ser o fluxo. O aviso de ocorrência reutilizado é o do #24, no Core. O alerta aceito avisa por WebSocket, no sentido da ADR_051. (fonte: contexto-llt201.md, §5 conflito 2 e §8; decisoes-joao-llt201.md, pendência 1, formulário 1) | Atualizar a spec 009 e o trecho HTTP desse contrato de referência. Não feito aqui. |
| Spec 009, questão em aberto 3 e FR-F07 | Rota, payload e evento do alerta ao Core estão pendentes na ADR_051. | O cliente já envia a rota, o corpo e o recibo descritos em FR-001. O R9 e a decisão de 26/11 fecham essa rota para a LLT-201. Depois do 201, o Core grava em `finder_alerts` e avisa no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`. (fonte: contexto-llt201.md, §3 e R9; decisoes-joao-llt201.md, rota POST, pendência 1, formulário 1, e Rodada 3) | Atualizar a spec 009 para não tratar o contrato do cliente como indefinido. Fechar o trecho pendente da ADR_051 no Drive, podendo registrar este contrato. |
| Spec 009, SC-002 e 009-22 A | Aceite de "atualizada" e "recorrente" bloqueado até a spec 008; reabertura mapeada como atualizada. | A homologação desta spec segue o FR-007. Aumento só de contagem não é mudança que esta decisão manda avisar. (fonte: spec.md, FR-007; decisoes-joao-llt201.md, pendência 5, formulário 1; respostas-pendentes-llt201.md, pendência 5) | Atualizar o texto do SC-002 na spec 009. A spec 008 não é editada aqui. |
| Spec 009, SC-006 e 009-15 A | Recorrente gera nova notificação, medida no módulo Notificação. | O Finder só se alinha na homologação. O aviso é o do #24, trabalho do Core. (fonte: decisoes-joao-llt201.md, evento para a Notificação; contexto-llt201.md, §1) | Atualizar o texto do SC-006 na spec 009 para não exigir um canal que a decisão de 26/11 não usa. |
| Spec 009, status | Draft. Aceite integrado ainda descrito como dependente da API do Módulo. | O cliente está na Development. O aceite com servidor real é esta spec, junto com a spec 016. (fonte: contexto-llt201.md, §1) | Atualizar o status da 009 quando essa atualização for autorizada. |
| ADR_165 | Ocorrências passam pela Notificação via stream. | Vale o João: evento para a Notificação é só o aviso do #24. O alerta aceito avisa por WebSocket no sentido da ADR_051, sem evento novo. (fonte: decisoes-joao-llt201.md, pendência 1, formulário 1; respostas-pendentes-llt201.md, conflito 1) | Só registrar. Não resolver de outro modo. |
| Spec 015, FR-007 | O aviso contém só o tipo de mudança, o id, o produto e o tipo, sem gravidade. | Vale o João e a pendência 2 C: o aviso do #24 leva a gravidade. (fonte: respostas-pendentes-llt201.md, pendência 2 e conflito 2; origem: 009-5 revisto) | Emenda da spec 015 no Core. Não escrita aqui. |
| ADR_062 | Gravidade = criticidade × tipo, escala "critico"/"alto", matriz pendente. | Vale o João: gravidade = classe ERROR, WARNING ou INFO no campo `gravidade`. (fonte: decisoes-joao-llt201.md, pendência 3, formulário 1; respostas-pendentes-llt201.md, conflito 3) | Só registrar. A matriz da ADR_062 não é adotada nesta homologação. |

## Pendências

As pendências 1, 2 e 3 desta pasta fecharam na sessão de 2026-11-26. Não há opção em aberto no lugar delas.

### Proxy local e certificado — fechados (mapa 13 A, 17 A e 18 B)

O teste local usa o Caddy no compose do Core. O Finder só recebe a URL `https://`. A CA local entra só no teste, via `SSL_CERT_FILE`, sem código novo e sem desligar a verificação. Era a Pendência 1 desta pasta. (fonte: decisoes-joao-llt201.md, pendências 13 e 18, formulário 2; respostas-pendentes-llt201.md, pendência 17)

### POST de alerta que falha — fechado (mapa 19 C)

Se o POST falha, o lote permanece pendente e o aviso é reenviado no reprocessamento do lote. Não há número novo de tentativas dentro do POST. Era a Pendência 2 desta pasta. (fonte: respostas-pendentes-llt201.md, pendência 19; origem: DF-FINDER §2 e §3)

### Audiências do JWT — fechado (mapa 20 B, com 6 B)

Um token com as duas audiências (a do Módulo e `lhamalog-core`). As duas variáveis podem carregar esse mesmo JWT. A propriedade `core-audience` vale só na rota de alerta. A recomendação A (dois tokens) não vale. Era a Pendência 3 desta pasta. (fonte: decisoes-joao-llt201.md, pendências 6 e 20, formulário 1)

Não resta decisão em aberto do Finder nesta pasta. A coleção do alerta é `finder_alerts` e o literal de "virou recorrente" é `RECURRENT`. Os dois fecharam na rodada 3. Não usar `AlertaFinder` nem `RECURRING`. (fonte: decisoes-joao-llt201.md, Rodada 3)

O residual sem fonte que esta spec não preenche: os sufixos de `event=occurrence.alert.*`. O prefixo está fechado (pendência 14 A) e fica fora do Finder, no pacote `br.com.lhamalog.core.finderalert`. A string exata de `event=` continua [PENDENTE] no Core. (fonte: decisoes-joao-llt201.md, pendência 14, formulário 2, e pendência 15, formulário 1; analyze-016.md, MEDIUM U2; sufixos: SEM FONTE)

O CRITICAL C1 (princípio I da constituição) permanece só registrado na seção Divergências e na tabela de documentos. Esta rodada não emenda a constituição.

## Pendências — documentos necessários

Nenhum destes documentos é criado nesta especificação.

| Documento | Onde | Por quê |
|---|---|---|
| Atualização da spec 009, inclusive `contracts/pipeline-aviso-referencia.md`, o status Draft, a questão 3 e o item 10 | `specs/009-publicar-evento-ocorrencia/` | O Módulo passou a ficar dentro do Core; o aviso HTTP Módulo → Core deixa de descrever o processo único; o cliente de alerta já tem rota e corpo; o item 10 da emenda constitucional continua aberto. (fonte: contexto-llt201.md, §5 conflito 8 e §8) |
| Emenda do princípio I da constituição | `.specify/memory/constitution.md` | O princípio I ainda diz que o Finder transforma o resultado em ocorrência. A spec 009, item 10, já aponta essa emenda. (fonte: constituição, princípio I; spec 009, questões em aberto, item 10) |
| ADR nova "Módulo Ocorrência integrado ao Core" | Drive e branch nova no lhamalog-docs | Registrar que a Ocorrência fica dentro do Core e substituir a parte da ADR_054 que a trata como serviço separado. (fonte: decisoes-joao-llt201.md, registro de que a Ocorrência fica dentro do Core; origem: contexto-llt201.md, §5 conflito 1 e §7 Q8 opção A) |
| Fechamento do trecho pendente da ADR_051 | Drive | A ADR_051 ainda deixa pendente o protocolo Finder → Core. O contrato já decidido, que a ADR pode registrar, é: depois do 201, canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`. Esta execução não edita a ADR. (fonte: decisoes-joao-llt201.md, Rodada 3 e pendência 1, formulário 1) |
| README e `.env.example` do Finder | Raiz deste repositório, numa execução que possa editá-los | Dizer que `FINDER_CORE_JWT` e `FINDER_MODULO_OCORRENCIA_JWT` podem carregar o mesmo JWT do segredo e do emissor do #20, com as duas audiências (`lhamalog-modulo-ocorrencia` e `lhamalog-core`), que as duas URLs `https://` apontam para o mesmo host, e que no teste a CA local entra só por `SSL_CERT_FILE`. Sem segredo no exemplo. (fonte: decisoes-joao-llt201.md, pendências 18 e 20, formulários 2 e 1) |
| README, exemplo de ambiente e compose do Caddy | lhamalog-core, spec 016 | O proxy é o Caddy no compose do Core. O Finder só consome a URL `https://`. O arquivo não é escrito aqui. (fonte: decisoes-joao-llt201.md, pendência 13, formulário 2; respostas-pendentes-llt201.md, pendência 17) |
| Emenda da spec 015, FR-007 | lhamalog-core, spec 015 | O aviso do #24 passa a levar a gravidade. Esta execução não edita a 015. (fonte: respostas-pendentes-llt201.md, pendência 2 e conflito 2; origem: 009-5 revisto) |
| Alinhamento de `endpoints-ocorrencia.md`, `endpoints-finder.md` e `endpoints-core-notificacao.md` | lhamalog-docs | Citar o DN-009-2 e a rota de alerta quando essa documentação for atualizada. (fonte: contexto-llt201.md, §8) |

## Evidência e consulta

- **Lido**: `decisoes-joao-llt201.md` (Rodada 3 prevalece nos três pontos que estavam sem fonte; formulários 1 e 2 continuam no restante); `prompts/r3-decisoes.txt`; `analyze-011.md`; `respostas-pendentes-llt201.md`; `prompts/r2-decisoes.txt`; constituição deste repositório (só leitura); spec 009 (`spec.md`, `contracts/proposta-dn-009-2.md`, `contracts/pipeline-aviso-referencia.md`); `internal/adapter/out/core/aviso.go` e `aviso_test.go`; `internal/adapter/out/confirmacao/http.go`; `internal/adapter/out/modocorrencia/cliente.go`; `internal/adapter/out/openobserve/recorrencia.go`; `internal/bootstrap/pipeline.go`. No Core, só leitura, sem edição: spec `016-alerta-finder-evento-ocorrencia` e o fato de código do #20 citado em `respostas-pendentes-llt201.md` (`InternalServiceSecurityConfiguration`, `aud` por "contém", `scripts/InternalJwt.java` com um `--aud`).
- **Executado**: nenhum teste, nenhum servidor, nenhum compose. Nenhum script PowerShell. A branch não foi trocada.
- **Em aberto**: sufixos de `event=occurrence.alert.*`, [PENDENTE] residual no Core, sem invenção aqui. Emenda do princípio I (CRITICAL C1), só registrada. Documentos da tabela acima, só listados.
