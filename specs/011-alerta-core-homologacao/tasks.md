---
description: "Task list for Alerta ao Core e homologação Finder"
---

# Tasks: Alerta ao Core e homologação Finder

**Input**: Design documents from `/specs/011-alerta-core-homologacao/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/, quickstart.md

**Tests**: Não há suíte nova. O plano registra que a homologação é o roteiro ponta a ponta e que `go test ./...`, `go test -race ./...` e `go vet ./...` já existem (`Makefile`: `make test`, `make lint`). Esta lista não pede arquivo de teste novo. (fonte: plan.md, Technical Context, Testing; spec.md, User Story 3 e FR-006)

**Organização**: Fases por história. A implementação espera #19, #20, #21 e #24 mergeados e não recria lista, detalhe, consultas, resumos, recibo, ciclo de vida, `CREATED`, `REOPENED`, o canal `lhamalog:occurrence-events` nem o WebSocket `/ws/occurrences`. (fonte: decisoes-joao-llt201.md, base da branch e evento para a Notificação; origem: contexto-llt201.md, §1; plan.md, Summary e Fora deste plano; spec.md, Fora de escopo)

**Regra desta lista**: a decisão da rodada 2 está escrita na tarefa, com fonte. A rodada 3 fecha o nome `finder_alerts` (T033) e o literal `RECURRENT` (T020), e move a conferência do WebSocket (T010) para a história 3, sem a marca [P]. Caixas marcadas nesta execução são só as tarefas concluídas no Finder, sem servidor real. Nenhuma tarefa traz texto de trabalho encerrado.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: pode andar em paralelo (arquivos diferentes, sem depender de tarefa ainda aberta)
- **[Story]**: história da spec (US1, US2, US3). Setup, Foundational e Polish não levam história

## Path Conventions

- Serviço único já existente: `cmd/finder`, `internal/domain`, `internal/application`, `internal/adapter`, `internal/bootstrap`, `internal/config` (plan.md, Project Structure)
- Clientes que esta feature configura e homologa, sem segunda implementação: `internal/adapter/out/core/aviso.go`, `internal/adapter/out/confirmacao/http.go`, `internal/adapter/out/modocorrencia/cliente.go`, `internal/bootstrap/pipeline.go`
- `docker-compose.yml` só interpola `FINDER_CORE_*` e `FINDER_MODULO_OCORRENCIA_*`. O Caddy não entra neste arquivo (plan.md, Structure Decision; spec.md, FR-005)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirmar o módulo já existente e o que esta feature não recria. Não há projeto novo nem dependência nova.

- [x] T001 Confirmar, sem criar código, que o trabalho desta feature fica na configuração de ambiente e na homologação dos clientes já presentes em `internal/adapter/out/core/aviso.go`, `internal/adapter/out/modocorrencia/cliente.go` e `internal/bootstrap/pipeline.go`. A base da implementação é a pilha #19 ← #20 ← #21 ← #24 já mergeada. Esta lista não recria `POST /internal/v1/ocorrencias/consultas`, `POST /internal/v1/resumos-analise`, o recibo idempotente, `409 REVISAO_DESATUALIZADA`, `409 ENTREGA_EM_ANDAMENTO`, o ciclo de vida, os avisos `CREATED` e `REOPENED`, o canal `lhamalog:occurrence-events` nem `/ws/occurrences`. A rota servidora `POST /internal/v1/alertas/finder` permanece a spec 016 do lhamalog-core, no pacote `finderalert` (exemplo `br.com.lhamalog.core.finderalert`), com segurança na cadeia `/internal/v1/**` e log de negócio no prefixo `occurrence.alert.*`. Este repositório não cria esse pacote nem esse log. Não adotar o prefixo `finder.alert.*`. Os sufixos da string exata de `event=occurrence.alert.*` ficam na pendência residual desta lista, no Core: esta tarefa não os inventa e não espera por eles. (fonte: plan.md, Summary, Project Structure e Fora deste plano; spec.md, abertura, Assumptions e Fora de escopo; decisoes-joao-llt201.md, base da branch, rota POST, pendências 14 e 15, formulários 2 e 1; analyze-016.md, MEDIUM U2; sufixos: SEM FONTE; origem: contexto-llt201.md, §1)
- [x] T002 [P] Manter `go.mod` na linguagem e nas dependências que o plano já cita: Go 1.27.4, `github.com/joho/godotenv` v1.5.1, `github.com/redis/go-redis/v9` v9.23.0 e `net/http` da biblioteca padrão. Não acrescentar `go.mongodb.org/mongo-driver`, `go.opentelemetry.io/otel`, `golang.org/x/sync` nem `github.com/testcontainers/testcontainers-go`. Não criar pacote novo. `SSL_CERT_FILE` é lida por `crypto/x509` da biblioteca padrão; não é dependência nova e não entra em `go.mod`. (fonte: plan.md, Technical Context e Constitution Check, Restrições técnicas e README; constituição, Restrições Técnicas; research.md, O que não virou pesquisa de tecnologia nova; decisoes-joao-llt201.md, pendência 18, formulário 2)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Preservar o cliente já entregue antes de qualquer história. Nenhuma história começa alterando esses arquivos.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T003 Manter em `internal/adapter/out/core/aviso.go` o pedido `POST /internal/v1/alertas/finder` com corpo `{tenantId, productId, batchId, motivo}`, motivo único `consulta_openobserve_indisponivel`, e sucesso somente 201 na primeira aceitação ou 200 na repetição da mesma identidade, com `{aceito: true}` e os mesmos quatro campos, sem campo `repeticao` e sem outro campo. Campo a mais, nome diferente, identidade diferente, 202 e corpo vazio continuam fora do aceite. Qualquer status diferente de 200 e 201 é falha para o cliente, inclusive 422 `DADOS_INVALIDOS` (corpo incompleto, campo a mais ou motivo fora da lista), 400 `REQUISICAO_INVALIDA` (JSON ilegível) e 403 `ACESSO_NEGADO` (produto que não pertence ao tenant). Não criar esses códigos neste repositório. A deduplicação durável fica no Core, coleção `finder_alerts`, chave `(tenantId, productId, batchId, motivo)`, distinta da chave do recibo. Quem perde a corrida de duas primeiras gravações relê e responde 200 com a mesma identidade; o cliente já trata 200 e 201 como sucesso. Não usar `AlertaFinder`. O nome `finder_alerts` está decidido e fica registrado em T033; esta tarefa não o deixa em aberto. (fonte: spec.md, FR-001, FR-009, US1 cenário 1 e Edge Cases; plan.md, O que a implementação cobre não inclui reescrita deste arquivo; research.md, D5; contracts/alerta-finder-cliente.md; data-model.md, Alerta operacional; decisoes-joao-llt201.md, pendências 9 e 12, formulário 1, e Rodada 3; respostas-pendentes-llt201.md, pendências 4, 7, 8 e 9; origem: DF-FINDER §3 + controller do #20 + cliente estrito; contexto-llt201.md, §3 e R9)
- [x] T004 [P] Manter em `internal/adapter/out/confirmacao/http.go` a função `Novo`: base só `https`, host obrigatório, sem usuário, consulta, fragmento ou caminho além de `/`, token no header `Authorization: Bearer`, e `CheckRedirect` devolvendo `http.ErrUseLastResponse`. O Finder não passa a montar nem a inspecionar `iss`, `aud` ou `sub`. A CA local do teste não altera este arquivo: entra só por `SSL_CERT_FILE` no processo, com a verificação padrão ligada. (fonte: spec.md, US1 cenário 4, US2 cenários 1, 4 e 5, FR-002 e FR-004; plan.md, Constraints e Structure Decision; research.md, D2 e D4; data-model.md, Credencial do alerta e URL; decisoes-joao-llt201.md, pendência 18, formulário 2)
- [x] T005 [P] Manter em `internal/adapter/out/modocorrencia/cliente.go` as rotas `POST /internal/v1/ocorrencias/consultas` e `POST /internal/v1/resumos-analise`, a base `https` com as mesmas restrições de URL, o Bearer de `FINDER_MODULO_OCORRENCIA_JWT` e os prazos padrão de 2 s na consulta e 5 s na entrega. A audiência exigida nestas rotas continua `lhamalog-modulo-ocorrencia`. O JWT pode ser o mesmo de `FINDER_CORE_JWT` quando o `aud` contém as duas audiências; `lhamalog-core` não substitui a audiência destas rotas. O contrato continua `specs/009-publicar-evento-ocorrencia/contracts/proposta-dn-009-2.md`. Esta pasta não publica outra versão: `specs/011-alerta-core-homologacao/contracts/dn-009-2-homologacao.md` só referencia. O campo `repeticao: true` do recibo do resumo não entra no corpo do alerta. (fonte: spec.md, FR-003, FR-004 e FR-006; plan.md, O que a implementação cobre item 2; research.md, D3 e D6; contracts/dn-009-2-homologacao.md; decisoes-joao-llt201.md, pendências 6 e 20, formulário 1; respostas-pendentes-llt201.md, pendência 8; origem: contexto-llt201.md, §1 e §3)
- [x] T006 [P] Manter `docker-compose.yml` só interpolando `FINDER_MODULO_OCORRENCIA_URL`, `FINDER_MODULO_OCORRENCIA_JWT`, `FINDER_CORE_URL`, `FINDER_CORE_JWT` e `FINDER_CORE_TIMEOUT`. Não acrescentar serviço de proxy neste arquivo. O proxy do teste local é o Caddy no compose do Core (spec 016). Este repositório não edita esse compose. (fonte: spec.md, FR-005 e US2 cenário 2; plan.md, Project Structure e Structure Decision; research.md, D4; decisoes-joao-llt201.md, pendência 13, formulário 2; respostas-pendentes-llt201.md, pendência 17; origem: ADR_164 + arquitetura-modulos)
- [x] T007 [P] Manter `FINDER_CORE_TIMEOUT` com padrão `2s`, lido em `internal/bootstrap/pipeline.go` (`prazo`). Não alterar esse prazo nem os prazos de `internal/adapter/out/modocorrencia/cliente.go`. Esse `2s` é o teto em que a rota de alerta do Core responde. O aviso "virou recorrente" no Core herda até 3 s depois do 201 da spec 015, SC-001. Os dois tetos valem em alvos diferentes. O Finder não implementa o teto de 3 s. (fonte: spec.md, SC-009 e Assumptions; plan.md, Performance Goals; research.md, D9; respostas-pendentes-llt201.md, pendência 16; origem: spec 015, SC-001 + DF-FINDER §3 e §8; contexto-llt201.md, §3)
- [x] T008 [P] Manter o Finder como quem entrega `ResumoAnalise` e consulta o achado. Não criar publicador, canal, envelope nem campo de gravidade neste repositório. O aviso de ocorrência reutilizado é o do #24, no Core (spec 016): pub/sub `lhamalog:occurrence-events` e WebSocket `/ws/occurrences`, um aviso por ocorrência por entrega, com a mudança mais forte. A gravidade conferida no Core é a classe ERROR, WARNING ou INFO no campo `gravidade`, no documento e nesse aviso, com o tipo mantido. O alerta aceito, à parte, na primeira gravação (201), grava na coleção `finder_alerts` e avisa a tela no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, sem evento novo para a Notificação e sem publicar em `lhamalog:occurrence-events` nem em `lhamalog:notifications:events`. A repetição (200) não avisa de novo. A conferência desse efeito é T010, na história 3, depois do token. Esta tarefa não cria esse aviso. O registro de que o Módulo é o pacote `br.com.lhamalog.core.occurrence` já está em `specs/011-alerta-core-homologacao/spec.md`, seção Responsabilidades; esta tarefa não edita a spec 009 nem `.specify/memory/constitution.md`. (fonte: spec.md, FR-007, FR-008, FR-010 e Responsabilidades; plan.md, Summary; research.md, D1, D7, D11 e D12; decisoes-joao-llt201.md, evento para a Notificação, pendências 1, 3 e 5, formulário 1, e Rodada 3; respostas-pendentes-llt201.md, pendências 2 e 5; origem: contexto-llt201.md, §1, R2 e §5 conflito 2)

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - O JWT do alerta usa o segredo do #20 e as duas audiências (Priority: P1) 🎯 MVP

**Goal**: `FINDER_CORE_JWT` e `FINDER_MODULO_OCORRENCIA_JWT` podem carregar o mesmo JWT, emitido fora do Finder, com o segredo e o emissor do #20 e com as duas audiências no `aud` (a do Módulo e `lhamalog-core`). O cliente não é reescrito para emitir ou inspecionar claims. A propriedade `core-audience` vale só na rota de alerta; as rotas do Módulo seguem com a audiência do #20. (fonte: spec.md, User Story 1; plan.md, O que a implementação cobre itens 1 e 2; research.md, D2 e D3)

**Independent Test**: Configurar as duas variáveis com o mesmo JWT (segredo e emissor do #20, `aud` contendo as duas audiências), chamar a rota de alerta e as rotas do Módulo e receber o aceite de cada uma. Dá para verificar sem percorrer os cinco cenários do DN-009-2. (fonte: spec.md, User Story 1, Independent Test)

### Decisão aplicada na história 1

- [x] T009 [P] [US1] Aplicar um único JWT com as duas audiências. `FINDER_CORE_JWT` e `FINDER_MODULO_OCORRENCIA_JWT` podem carregar o mesmo valor: HS256, segredo e emissor do #20, `aud` contendo a audiência do Módulo (padrão `lhamalog-modulo-ocorrencia`) e `lhamalog-core`. A recomendação de dois tokens não vale. Não alterar `internal/adapter/out/confirmacao/http.go` nem `internal/adapter/out/modocorrencia/cliente.go` para emitir ou inspecionar `aud`. No Core (spec 016, fora deste repositório), a cadeia `/internal/v1/**` do #20 precisa aceitar esse token na rota de alerta pela audiência `lhamalog-core`, lida da propriedade `core-audience` usada só nessa rota, e nas rotas do Módulo pela audiência do Módulo. Ilustração de caminho, não nome extra decidido: `lhamalog.occurrence.internal-auth.core-audience`. Fato de código: `InternalServiceSecurityConfiguration` cobre todo `/internal/v1/**` e exige que o `aud` contenha `lhamalog.occurrence.internal-auth.audience` (padrão `lhamalog-modulo-ocorrencia`); o `aud` é conferido por "contém"; um token só com `lhamalog-core`, sem o ajuste por rota, recebe 401. `scripts/InternalJwt.java` gera uma audiência por `--aud`. Esta tarefa não altera o script. (fonte: spec.md, FR-002, FR-003 e Pendências, audiências do JWT; plan.md, O que a implementação cobre itens 1 e 2; research.md, D2, D3 e P3; data-model.md, Credencial do alerta e da consulta; contracts/alerta-finder-cliente.md; decisoes-joao-llt201.md, pendências 6 e 20, formulário 1; respostas-pendentes-llt201.md, resumo, fato de código, e pendência 20)

### Implementation for User Story 1

- [ ] T011 [P] [US1] Colocar em `FINDER_CORE_JWT`, lida por `internal/bootstrap/pipeline.go` e, quando houver arquivo de ambiente local, por `internal/config/env.go` (`AplicarArquivoEnv`), um JWT HS256 emitido fora do processo Finder. Segredo: o de `lhamalog.occurrence.internal-auth.secrets` (ambiente do Core `OCCURRENCE_INTERNAL_JWT_SECRETS`). Emissor: o de `lhamalog.occurrence.internal-auth.issuer` (padrão `lhamalog-internal`). O `aud` contém a audiência do Módulo e `lhamalog-core`. `sub`: um valor da lista `lhamalog.occurrence.internal-auth.subjects` (padrão `lhamalog-finder`); a decisão de 26/11 não cria outro `sub`. O mesmo valor pode ir para `FINDER_MODULO_OCORRENCIA_JWT` (T012). `scripts/InternalJwt.java` gera uma audiência por `--aud` e lê `OCCURRENCE_INTERNAL_JWT_SECRET` (singular); a propriedade do Core é `OCCURRENCE_INTERNAL_JWT_SECRETS` (plural). Não unificar esses nomes, não copiar o script para este repositório e não gravar segredo nem token no Git. Uma invocação com um único `--aud` não é, por si, o token desta tarefa. O cliente em `internal/adapter/out/confirmacao/http.go` continua enviando o valor inteiro em `Authorization: Bearer`. (fonte: spec.md, US1 cenário 2, FR-002, Assumptions e Roteiro itens 5 e 6; plan.md, O que a implementação cobre item 1; research.md, D2; data-model.md, Credencial do alerta; quickstart.md, Configuração; decisoes-joao-llt201.md, JWT da rota de alerta e pendência 20, formulário 1; respostas-pendentes-llt201.md, resumo, fato de código; origem: contexto-llt201.md, §7 Q6 opção A e §2)
- [ ] T012 [P] [US1] Colocar em `FINDER_MODULO_OCORRENCIA_JWT` (`internal/config/analise.go`, constantes `EnvModuloOcorrenciaURL` e `EnvModuloOcorrenciaJWT`) o JWT das rotas `POST /internal/v1/ocorrencias/consultas` e `POST /internal/v1/resumos-analise`. Pode ser o mesmo valor de `FINDER_CORE_JWT`: mesmo segredo, mesmo emissor, `aud` com `lhamalog-modulo-ocorrencia` e `lhamalog-core`, `sub` na lista do #20 (padrão `lhamalog-finder`). A audiência exigida nestas rotas continua `lhamalog-modulo-ocorrencia`. A audiência `lhamalog-core` não substitui essa. A propriedade `core-audience` não é lida nestas rotas. Segredo só no ambiente. (fonte: spec.md, US1 cenário 3, FR-003 e Roteiro item 5 (JWT); plan.md, O que a implementação cobre item 2; research.md, D3; data-model.md, Credencial da consulta e da entrega; contracts/dn-009-2-homologacao.md; decisoes-joao-llt201.md, pendências 6 e 20, formulário 1; origem: contexto-llt201.md, §5 conflito 6)
- [x] T013 [P] [US1] Conferir, sem alterar código, que `internal/adapter/out/confirmacao/http.go` e `internal/bootstrap/pipeline.go` não montam nem validam `iss`, `aud` ou `sub`: o header da rota de alerta leva o valor de `FINDER_CORE_JWT`, e o das rotas do Módulo leva o valor de `FINDER_MODULO_OCORRENCIA_JWT`. Não acrescentar biblioteca de JWT. A documentação dessa configuração é a tarefa de DOCUMENTAÇÃO/pendência sobre `README.md` e `.env.example` (T024), não uma edição destes arquivos Go. (fonte: spec.md, US1 cenário 4 e FR-002; plan.md, Structure Decision; research.md, D2; constituição, princípio V)

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - O teste local usa HTTPS pelo Caddy do Core, sem liberar HTTP (Priority: P1)

**Goal**: No teste local o Finder aponta `FINDER_CORE_URL` e `FINDER_MODULO_OCORRENCIA_URL` para o `https://` do Caddy no compose do Core, que permanece em HTTP na porta 8080. As duas bases usam o mesmo host. O cliente não aceita `http://`. A CA local entra só no teste, por `SSL_CERT_FILE`, com a verificação ligada. (fonte: spec.md, User Story 2 e FR-004 e FR-005; plan.md, O que a implementação cobre item 3; research.md, D4)

**Independent Test**: Apontar o Finder para a URL `https://` do Caddy, com a CA local em `SSL_CERT_FILE`, e concluir uma chamada de alerta e uma entrega de resumo. Uma base `http://` é recusada pelo Finder antes de qualquer chamada. A verificação do certificado permanece ligada. (fonte: spec.md, User Story 2, Independent Test)

### Decisão aplicada na história 2

- [ ] T014 [P] [US2] Aplicar o proxy do teste local: Caddy no compose do Core. O produto é Caddy. O arquivo de compose é o do Core, o que a spec 016 registrar. O Finder só recebe a URL `https://` desse Caddy, na frente do Core local em HTTP na porta 8080. Não editar `docker-compose.yml` deste repositório e não acrescentar proxy no compose do Finder. Não encerrar TLS no processo do Core sem proxy e não liberar `http://` no cliente. (fonte: spec.md, FR-005 e US2 cenário 2; plan.md, Constraints e O que a implementação cobre item 3; research.md, P1 e D4; quickstart.md, Pré-requisitos item 2; decisoes-joao-llt201.md, pendência 13, formulário 2; respostas-pendentes-llt201.md, pendência 17; origem: ADR_164 + arquitetura-modulos; contexto-llt201.md, §7 Q7 opção A e §5 conflito 7)
- [ ] T015 [P] [US2] Aplicar a confiança no certificado local: a CA do Caddy entra só no teste, pela variável de ambiente `SSL_CERT_FILE`, com a verificação HTTPS ligada. Não criar variável de certificado no Finder, não alterar `internal/adapter/out/confirmacao/http.go` e não desligar a verificação. O `crypto/x509` da biblioteca padrão já lê `SSL_CERT_FILE`. O caminho do arquivo é o que o Caddy da spec 016 expõe; esta tarefa não nomeia outro caminho. Fora do teste, a CA local não entra. (fonte: spec.md, US2 cenário 5, Assumptions e Edge Cases; research.md, P1 e D4; data-model.md, URL; contracts/alerta-finder-cliente.md, Transporte; quickstart.md, Configuração; decisoes-joao-llt201.md, pendência 18, formulário 2; respostas-pendentes-llt201.md, pendência 18; origem: o pacote `crypto/x509` do Go lê `SSL_CERT_FILE`)

### Implementation for User Story 2

- [x] T016 [P] [US2] Confirmar a recusa de base `http://` em `FINDER_CORE_URL` e em `FINDER_MODULO_OCORRENCIA_URL` pela prontidão já existente em `internal/bootstrap/pipeline.go`, com as regras de `internal/adapter/out/confirmacao/http.go` (`Novo`) e `internal/adapter/out/modocorrencia/cliente.go` (`NovoCliente`). Nenhuma chamada sai. Resultado esperado do SC-002. Não acrescentar opção de liberar HTTP e não criar suíte nova; o caso `http://` já aparece em `internal/bootstrap/pipeline_test.go`. (fonte: spec.md, US2 cenário 1, FR-004 e SC-002; quickstart.md, O que conferir antes do fluxo item 1; research.md, D4; plan.md, Testing; decisoes-joao-llt201.md, TLS no teste local)
- [ ] T017 [US2] Apontar `FINDER_CORE_URL` e `FINDER_MODULO_OCORRENCIA_URL` para a URL `https://` do Caddy que a spec 016 registrar no compose do Core, o mesmo host, sem usuário, consulta, fragmento ou caminho além de `/`. O Core local segue em HTTP na porta 8080. A URL configurada é o próprio ponto HTTPS: o cliente não segue redirecionamento (`internal/adapter/out/confirmacao/http.go`). No teste, exportar a CA local só em `SSL_CERT_FILE`, com a verificação ligada e sem código novo. Não editar `docker-compose.yml`. O cliente de hoje aceita hosts diferentes; não acrescentar recusa de hosts diferentes em `internal/adapter/out/modocorrencia/cliente.go` nem em `internal/adapter/out/confirmacao/http.go`. A configuração correta da homologação continua o mesmo host. (fonte: spec.md, US2 cenários 2, 3, 4 e 5, FR-005 e Edge Cases; plan.md, O que a implementação cobre item 3; research.md, D4 e D8; data-model.md, Módulo de Ocorrência e URL; quickstart.md, Configuração; decisoes-joao-llt201.md, pendências 13 e 18, formulário 2; respostas-pendentes-llt201.md, pendência 17; origem: contexto-llt201.md, §8)

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Homologar o Finder com o Core de verdade (Priority: P1)

**Goal**: Percorrer o alerta e os cinco cenários do DN-009-2 com o cliente já existente e o Core da pilha mergeada mais a rota da spec 016. O Finder não cria ocorrência, não calcula gravidade e não publica o evento de ocorrência. A homologação confere no Core o efeito do alerta (FR-010), a gravidade, a precedência do FR-007 (literal `RECURRENT`) e o reenvio quando o POST falha. (fonte: spec.md, User Story 3, FR-006, FR-007, FR-009 e FR-010; plan.md, O que a implementação cobre itens 4, 5 e 6; research.md, D1, D6, D7 e D12)

**Independent Test**: Subir Finder e Core, com o Caddy no compose do Core e `SSL_CERT_FILE` no teste, e percorrer os cinco cenários do DN-009-2 usando só as rotas que o #20 já entrega. A conferência de ocorrência nova, reaberta e "virou recorrente", da gravidade e do aviso por WebSocket do alerta é a do Core (spec 016); o Finder entra como quem entrega o resumo e quem envia o alerta. (fonte: spec.md, User Story 3, Independent Test)

### Decisão aplicada na história 3

- [x] T018 [P] [US3] Aplicar a falha do POST do alerta: se a chamada falha, o lote permanece pendente e o aviso é reenviado no reprocessamento desse lote, com a mesma identidade, apoiado na deduplicação do Core. Não acrescentar um número novo de tentativas dentro do POST em `internal/adapter/out/openobserve/recorrencia.go`. O caminho já existente consulta o OpenObserve até 3 vezes e, esgotadas, faz uma chamada; se ela falha, devolve `ErrAvisoCorePendente` e não marca `recorrente=false` por essa degradação. A opção escolhida é C. Ela é compatível com 009-18 A, e 009-18 A não é a regra aplicada a este POST. Status diferente de 200/201 conta como essa falha, inclusive 422 `DADOS_INVALIDOS`, 400 `REQUISICAO_INVALIDA` e 403 `ACESSO_NEGADO`. (fonte: spec.md, FR-009, SC-007, Edge Cases e US3 cenário 6; research.md, P2 e D5; data-model.md, Transições que o Finder já executa; quickstart.md, Cenários ponta a ponta, Alerta item 6; respostas-pendentes-llt201.md, pendências 4, 12 e 19; origem: DF-FINDER §2 e §3; contexto-llt201.md, R8 e R9)
- [x] T019 [P] [US3] Aplicar a gravidade que a homologação confere no Core: copiar a classe ERROR, WARNING ou INFO para o campo `gravidade` e manter o tipo. A gravidade aparece no documento da ocorrência e no aviso do #24. O Finder não calcula gravidade e não ganha campo em `internal/adapter/out/modocorrencia/` nem em `internal/domain`. O conflito ADR_062 (escala "critico"/"alto", matriz pendente) fica só registrado: vale a classe ERROR/WARNING/INFO. O conflito com a spec 015, FR-007 (aviso sem gravidade) fica só registrado: vale a gravidade no aviso; a emenda da 015 é T034 e não é escrita aqui. (fonte: spec.md, FR-007, SC-005 e Edge Cases; data-model.md, Ocorrência, gravidade e aviso; research.md, D7; quickstart.md, Conferência que não é do Finder; decisoes-joao-llt201.md, pendência 3, formulário 1; respostas-pendentes-llt201.md, pendência 2 e conflito 3; origem: Ocorrência.mmd (Drive) + 009-5 revisto; contexto-llt201.md, R2 e §5 conflito 3)
- [ ] T020 [P] [US3] Conferir no Core o literal `RECURRENT` no enum `Change` para "virou recorrente", distinto de `CREATED` e de `REOPENED`. Não adotar `RECURRING`. Não criar essa constante em `internal/`. A precedência e o gatilho são os do FR-007; esta tarefa não os reescreve. O Finder não publica em `lhamalog:occurrence-events`. (fonte: spec.md, FR-007; decisoes-joao-llt201.md, Rodada 3 e pendência 5, formulário 1; respostas-pendentes-llt201.md, pendências 5 e 11; origem: constituição do Core, princípio IV + enum `Change` do #24)

### Implementation for User Story 3

- [ ] T010 [US3] Conferir, depois da configuração do token (T011) e não em paralelo com ela, o efeito do aceite durável. Na primeira gravação (201), o Core grava um documento na coleção `finder_alerts` e avisa a tela no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, sem evento novo, sem envelope próprio para a Notificação e sem publicar em `lhamalog:occurrence-events` nem em `lhamalog:notifications:events`. A repetição (200) não avisa de novo. O aceite que o Finder exige continua 201 na primeira aceitação ou 200 na repetição, com `{aceito: true}` e a mesma identidade, sem o campo `repeticao`, em `internal/adapter/out/core/aviso.go`. O Finder não publica esse aviso. A recomendação de só guardar, sem aviso, não vale. O conflito ADR_165 fica só registrado: vale o João. Sem a marca [P]: esta conferência espera o JWT das duas audiências (FR-010, US3 cenário 5). (fonte: spec.md, FR-010, US3 cenário 5 e SC-008; contracts/alerta-finder-cliente.md; research.md, D12; quickstart.md, Cenários ponta a ponta, Alerta item 4; decisoes-joao-llt201.md, pendência 1, formulário 1, e Rodada 3; respostas-pendentes-llt201.md, conflito 1; analyze-011.md, MEDIUM I4)
- [ ] T021 [US3] Homologar o alerta com o cliente já existente, depois de T011 e, no teste local, de T017. Provocar o caminho de `internal/adapter/out/openobserve/recorrencia.go` em que a consulta ao OpenObserve esgota as três tentativas. Esperado: um `POST /internal/v1/alertas/finder` com o corpo de `specs/011-alerta-core-homologacao/contracts/alerta-finder-cliente.md` e o Bearer de `FINDER_CORE_JWT`. Esperado no Finder, na primeira aceitação: 201 com `aceito: true` e a mesma identidade, sem campo `repeticao`. Na repetição da mesma identidade: 200, o mesmo corpo. Duas primeiras requisições da mesma chave: quem grava responde 201; quem perde a corrida relê e responde 200; o cliente trata os dois como sucesso quando a identidade confere. Isso cobre SC-001 e a parte de alerta do SC-003. Um token só com audiência `lhamalog-modulo-ocorrencia` não autentica a rota de alerta. Um token só com `lhamalog-core` não é a credencial das rotas do Módulo; sem o ajuste por rota da pendência 6, a cadeia atual responde 401. Não alterar `aviso.go`. Se o POST falha, seguir T018: o lote permanece pendente e o aviso volta no reprocessamento, sem número novo de tentativas dentro do POST (SC-007). Depois da primeira gravação (201), seguir T010: conferir no Core um documento em `finder_alerts` e o aviso no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}` (SC-008). A repetição (200) não avisa de novo. (fonte: spec.md, US1, US3 cenário 5, SC-001, SC-003, SC-007, SC-008 e Roteiro itens 8, 11 e 12; quickstart.md, Cenários ponta a ponta, Alerta; contracts/alerta-finder-cliente.md; plan.md, O que a implementação cobre itens 1, 5 e 6; decisoes-joao-llt201.md, pendências 1 e 9, formulário 1; respostas-pendentes-llt201.md, pendências 8, 19 e resumo, fato de código; origem: contexto-llt201.md, §3 e R9)
- [ ] T022 [US3] Homologar uma entrega de `ResumoAnalise` pela URL `https://` do mesmo host do alerta e, em seguida, os cinco cenários do DN-009-2, com `internal/adapter/out/modocorrencia/cliente.go` e o Core da pilha mergeada. No teste local, esta tarefa espera T012 e T017, com a CA local só em `SSL_CERT_FILE` e a verificação ligada. Cenários, já escritos em `specs/009-publicar-evento-ocorrencia/contracts/proposta-dn-009-2.md` e repetidos em `specs/011-alerta-core-homologacao/contracts/dn-009-2-homologacao.md`: 1) a consulta distingue ocorrência ausente de falha e rejeita acesso a outro tenant ou produto; 2) a primeira entrega grava exatamente um efeito e retorna 201; reenvio depois de perda de resposta, inclusive depois de reinício do Finder, retorna 200 e o mesmo `reciboId`, sem duplicar; 3) dois lotes que consultam a mesma revisão preservam todos os incrementos; revisão desatualizada em um achado rejeita a entrega inteira antes de qualquer escrita; 4) confirmação inválida ou inesperada mantém o lote pendente no Finder; um recibo válido antecede a conclusão do lote e o reconhecimento da fila; 5) o JSON mantém as três classes, carrega somente trechos protegidos e conserva `linhas` vazias quando for o caso; a API rejeita tipos e contadores inválidos. Isso cobre o restante do SC-003 e o SC-004. Não criar segunda implementação dessas chamadas. (fonte: spec.md, US3 cenários 1, 2 e 4, FR-006 e SC-004; plan.md, O que a implementação cobre item 4; research.md, D6; quickstart.md, Entrega no mesmo host; origem: contexto-llt201.md, §1; `proposta-dn-009-2.md`, seção Cenários de homologação com o Módulo)
- [ ] T023 [US3] Na mesma homologação, para uma entrega que o Core trata como ocorrência nova, como reaberta ou como "virou recorrente", limitar o Finder à entrega do `ResumoAnalise` pelo cliente de T022. O processo Finder não publica em `lhamalog:occurrence-events` e não calcula gravidade (SC-005). Conferir no Core o FR-007: precedência, gatilho (`before.recurring` não é TRUE e `after.recurring` é TRUE; ausente ou `null` conta como `false`) e literal `RECURRENT` (T020). O aviso "virou recorrente" sai em até 3 s depois do 201 (SC-009). A gravidade no documento e no aviso é a classe ERROR, WARNING ou INFO no campo `gravidade`, tipo mantido (T019). Não adotar `RECURRING`. (fonte: spec.md, US3 cenário 3, FR-007, SC-005, SC-009 e Roteiro item 10; quickstart.md, Conferência que não é do Finder; research.md, D7 e D9; decisoes-joao-llt201.md, pendências 3 e 5, formulário 1, e Rodada 3; respostas-pendentes-llt201.md, pendências 2, 5, 11 e 16)

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentos necessários, resíduos sem fonte, suíte já existente e o roteiro. Nenhuma tarefa desta fase cria ADR, README, `.env.example`, compose, spec 009, spec 015, constituição nem arquivo no Drive ou no lhamalog-docs.

- [ ] T024 [P] DOCUMENTAÇÃO/pendência: quando uma execução autorizada puder editar a raiz deste repositório, atualizar `README.md` (tabela de variáveis) e `.env.example` para dizer que `FINDER_CORE_JWT` e `FINDER_MODULO_OCORRENCIA_JWT` podem carregar o mesmo JWT do segredo e do emissor do #20, com as duas audiências (`lhamalog-modulo-ocorrencia` e `lhamalog-core`), que `FINDER_CORE_URL` e `FINDER_MODULO_OCORRENCIA_URL` são `https://` no mesmo host, e que no teste a CA local entra só por `SSL_CERT_FILE`. Sem segredo no exemplo. Esta lista não escreve esses arquivos. (fonte: plan.md, Pendências — documentos necessários; spec.md, US1 cenário 4, FR-002, FR-005 e Fora de escopo; decisoes-joao-llt201.md, pendências 18 e 20, formulários 2 e 1; origem: contexto-llt201.md, §8)
- [ ] T025 [P] DOCUMENTAÇÃO/pendência: emenda do princípio I em `.specify/memory/constitution.md`. O princípio I ainda diz que o Finder transforma o resultado em ocorrência; para esta feature a ocorrência nasce no pacote `br.com.lhamalog.core.occurrence` e o Finder entrega `ResumoAnalise`. FR-008 proíbe editar a constituição nesta feature. Não escrever o arquivo. Este item é o CRITICAL C1 de analyze-011.md: esta rodada só o lista e não emenda a constituição. A divergência permanece em `specs/011-alerta-core-homologacao/spec.md`, seção Divergências, e em `specs/011-alerta-core-homologacao/plan.md`, Complexity Tracking. (fonte: plan.md, Constitution Check gate I, Complexity Tracking e Pendências — documentos necessários; spec.md, FR-008, SC-006 e Divergências; constituição, princípio I e Governance; research.md, D10; spec 009, questões em aberto, item 10, citado em spec.md; origem: contexto-llt201.md, R2; decisoes-joao-llt201.md, registro de que a Ocorrência fica dentro do Core)
- [ ] T026 [P] DOCUMENTAÇÃO/pendência: atualização de `specs/009-publicar-evento-ocorrencia/`, inclusive `specs/009-publicar-evento-ocorrencia/contracts/pipeline-aviso-referencia.md`, o status Draft, a questão 3, o item 10, o SC-002 e o SC-006 daquela spec. O Módulo ficou dentro do Core; o HTTP Módulo → Core deixa de descrever o processo único; o cliente de alerta já tem rota e corpo; a homologação desta feature confere nova, reaberta e virou recorrente no Core, um aviso por ocorrência por entrega; o alerta aceito avisa por WebSocket. FR-008 proíbe editar a spec 009 nesta feature. Não escrever esses arquivos. (fonte: plan.md, Pendências — documentos necessários; spec.md, Divergências e FR-008; origem: contexto-llt201.md, §5 conflito 2 e §8; decisoes-joao-llt201.md, evento para a Notificação e pendência 1, formulário 1)
- [ ] T027 [P] DOCUMENTAÇÃO/pendência: ADR nova "Módulo Ocorrência integrado ao Core", no Drive e em branch nova no lhamalog-docs. Registra que a Ocorrência fica dentro do Core e substitui a parte da ADR_054 que a trata como serviço separado. Não criar a ADR, não editar o Drive e não criar branch no lhamalog-docs. (fonte: plan.md, Pendências — documentos necessários; spec.md, FR-008 e tabela de documentos; decisoes-joao-llt201.md, registro de que a Ocorrência fica dentro do Core; origem: contexto-llt201.md, §5 conflito 1 e §7 Q8 opção A)
- [ ] T028 [P] DOCUMENTAÇÃO/pendência: fechamento do trecho pendente da ADR_051 no Drive. A ADR ainda deixa pendente o protocolo Finder → Core. O contrato já decidido, que a ADR pode registrar, é: depois do 201, canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`. O efeito da pendência 1 A permanece: gravar o alerta e avisar a tela, sem evento novo de Notificação e sem publicar em `lhamalog:occurrence-events` nem em `lhamalog:notifications:events`. Não criar nem editar a ADR. (fonte: plan.md, Pendências — documentos necessários; spec.md, Divergências e tabela de documentos; decisoes-joao-llt201.md, pendência 1, formulário 1, e Rodada 3; origem: contexto-llt201.md, §8 e R9)
- [ ] T029 [P] DOCUMENTAÇÃO/pendência: README, exemplo de ambiente e compose do Caddy no lhamalog-core, spec 016. O proxy é o Caddy no compose do Core. O Finder só consome a URL `https://`. Não criar esses arquivos neste repositório e não editar `docker-compose.yml`. (fonte: plan.md, Pendências — documentos necessários; spec.md, FR-005; decisoes-joao-llt201.md, pendência 13, formulário 2; respostas-pendentes-llt201.md, pendência 17; origem: ADR_164 + arquitetura-modulos; contexto-llt201.md, §8)
- [ ] T030 [P] DOCUMENTAÇÃO/pendência: alinhar `endpoints-ocorrencia.md`, `endpoints-finder.md` e `endpoints-core-notificacao.md` no lhamalog-docs para citar o DN-009-2 e a rota de alerta. O contexto cita a pasta `planejamento/LLT-145-endpoints-dashboard/` para `endpoints-ocorrencia.md`. Não escrever esses arquivos. (fonte: plan.md, Pendências — documentos necessários; spec.md, tabela de documentos; origem: contexto-llt201.md, §6 e §8)
- [x] T031 [P] Rodar a suíte já existente na raiz, sem Core: `make test` (`go test -race ./...` via `Makefile`) e `make lint` (`go vet ./...`). Isso não substitui a homologação com servidor real e não cria teste novo. Não há variável de integração Finder↔Core no README para ligar os cinco cenários. (fonte: plan.md, Testing; quickstart.md, Suíte já existente; README.md, seção Testes; spec.md, User Story 3)
- [ ] T032 Validar `specs/011-alerta-core-homologacao/quickstart.md` na ordem do arquivo. A recusa de `http://` (T016) e a suíte de T031 não esperam o Caddy. A chamada local de alerta (T021), a entrega e os cinco cenários (T022) e a conferência no Core (T023) usam a URL `https://` do Caddy e `SSL_CERT_FILE` (T014, T015, T017). Falha do POST segue T018. Gravidade segue T019. O literal conferido é `RECURRENT` (T020, FR-007). O nome conferido da coleção é `finder_alerts` (T033, T010). SC-006: `specs/009-publicar-evento-ocorrencia/` e `.specify/memory/constitution.md` permanecem sem edição nesta feature; a divergência fica em `specs/011-alerta-core-homologacao/spec.md`, seção Divergências. SC-007, SC-008 e SC-009 entram na conferência. (fonte: spec.md, SC-001 a SC-009 e Roteiro de homologação; quickstart.md; plan.md, Evidência e consulta)
- [x] T033 [P] Registrar, sem criar coleção neste repositório, que o nome da coleção Mongo do alerta no Core é `finder_alerts`, chave `(tenantId, productId, batchId, motivo)`, distinta da chave do recibo e da coleção `Ocorrencia`. Não usar `AlertaFinder`. A homologação (T010, T021) confere a gravação nesse nome. (fonte: decisoes-joao-llt201.md, Rodada 3; spec.md, FR-010; plan.md, Technical Context, Storage; respostas-pendentes-llt201.md, pendência 7; origem: padrão do #20 + DF-FINDER §3)
- [ ] T034 [P] DOCUMENTAÇÃO/pendência: emenda da spec 015, FR-007, no lhamalog-core. O aviso do #24 passa a levar a gravidade (a classe no campo `gravidade`). O texto atual da 015 ainda descreve o aviso sem gravidade. Não editar a spec 015 neste repositório. O conflito com o 009-5 revisto fica só registrado: vale o João. (fonte: plan.md, Pendências — documentos necessários e Conflitos só registrados; spec.md, Divergências; respostas-pendentes-llt201.md, pendência 2 e conflito 2; origem: 009-5 revisto)

---

## Pendências — documentos necessários

Nenhum destes arquivos é criado por esta lista. As tarefas T024 a T030 e T034 só registram a pendência.

| Documento | Onde | Por quê | Tarefa |
| --- | --- | --- | --- |
| Emenda do princípio I | `.specify/memory/constitution.md` | O princípio I ainda diz que o Finder transforma o resultado em ocorrência. A spec 009, item 10, já aponta essa emenda. FR-008 proíbe a edição nesta feature. (fonte: constituição, princípio I; spec.md, Divergências; plan.md, Complexity Tracking) | T025 |
| Atualização da spec 009, inclusive `contracts/pipeline-aviso-referencia.md`, o status Draft, a questão 3, o item 10, o SC-002 e o SC-006 daquela spec | `specs/009-publicar-evento-ocorrencia/` | O Módulo ficou dentro do Core; o aviso HTTP Módulo → Core deixa de descrever o processo único; o cliente de alerta já tem rota e corpo; os termos da homologação são nova, reaberta e virou recorrente, um aviso por ocorrência por entrega, conferidos no Core; o alerta aceito avisa por WebSocket. (fonte: spec.md, Divergências; contexto-llt201.md, §5 conflito 2 e §8; decisoes-joao-llt201.md, evento para a Notificação e pendência 1, formulário 1) | T026 |
| ADR nova "Módulo Ocorrência integrado ao Core" | Drive e branch nova no lhamalog-docs | Registrar que a Ocorrência fica dentro do Core e substituir a parte da ADR_054 que a trata como serviço separado. (fonte: decisoes-joao-llt201.md, registro de que a Ocorrência fica dentro do Core; origem: contexto-llt201.md, §5 conflito 1 e §7 Q8 opção A) | T027 |
| Fechamento do trecho pendente da ADR_051 | Drive | A ADR_051 ainda deixa pendente o protocolo Finder → Core. O contrato já decidido, que a ADR pode registrar, é: depois do 201, canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`. Esta lista não cria nem edita a ADR. (fonte: spec.md, tabela de documentos; decisoes-joao-llt201.md, pendência 1, formulário 1, e Rodada 3; origem: contexto-llt201.md, §8 e R9) | T028 |
| README e `.env.example` do Finder | Raiz deste repositório, numa execução que possa editá-los | Dizer que `FINDER_CORE_JWT` e `FINDER_MODULO_OCORRENCIA_JWT` podem carregar o mesmo JWT do segredo e do emissor do #20, com as duas audiências (`lhamalog-modulo-ocorrencia` e `lhamalog-core`), que as duas URLs `https://` apontam para o mesmo host, e que no teste a CA local entra só por `SSL_CERT_FILE`. Sem segredo no exemplo. (fonte: spec.md, FR-002 e FR-005; decisoes-joao-llt201.md, pendências 18 e 20, formulários 2 e 1) | T024 |
| README, exemplo de ambiente e compose do Caddy | lhamalog-core, spec 016 | O proxy é o Caddy no compose do Core. O Finder só consome a URL `https://`. (fonte: decisoes-joao-llt201.md, pendência 13, formulário 2; respostas-pendentes-llt201.md, pendência 17; origem: ADR_164 + arquitetura-modulos) | T029 |
| Emenda da spec 015, FR-007 | lhamalog-core, spec 015 | O aviso do #24 passa a levar a gravidade. Esta lista não edita a 015. (fonte: respostas-pendentes-llt201.md, pendência 2 e conflito 2; origem: 009-5 revisto) | T034 |
| Alinhamento de `endpoints-ocorrencia.md`, `endpoints-finder.md` e `endpoints-core-notificacao.md` | lhamalog-docs; o contexto cita `planejamento/LLT-145-endpoints-dashboard/` para o primeiro | Citar o DN-009-2 e a rota de alerta quando essa documentação for atualizada. (fonte: contexto-llt201.md, §6 e §8; plan.md, Pendências — documentos necessários) | T030 |

`specs/011-alerta-core-homologacao/quickstart.md` já é o roteiro. Não substitui o README. (fonte: plan.md, Pendências — documentos necessários, linha do quickstart)

Não criar ADR_146 nem documento de envelope novo. O aviso de ocorrência desta feature é o do #24, no Core. O alerta aceito, depois do 201, avisa no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`. (fonte: decisoes-joao-llt201.md, evento para a Notificação, pendência 1, formulário 1, e Rodada 3; spec.md, FR-007 e FR-010)

## Pendência residual sem fonte

Os sufixos da string exata de `event=occurrence.alert.*` permanecem [PENDENTE] no Core (SEM FONTE, MEDIUM U2 da spec 016). O prefixo `occurrence.alert.*` está fechado (pendência 14 A) e fica fora do Finder, no pacote `br.com.lhamalog.core.finderalert`. Esta lista não inventa esses sufixos. Nenhuma tarefa da coleção (`finder_alerts`, T033), do literal (`RECURRENT`, T020) ou do WebSocket (T010) fica bloqueada por esse residual. (fonte: decisoes-joao-llt201.md, pendência 14, formulário 2, e pendência 15, formulário 1; analyze-016.md, MEDIUM U2; spec.md, Pendências; plan.md, Fora deste plano)

## Conflitos só registrados

Estes conflitos não são resolvidos de outro modo e não editam constituição, spec 015 nem ADR.

| Conflito | O que vale nesta feature | Fonte |
| --- | --- | --- |
| ADR_165 (ocorrências pela Notificação via stream) × "só o aviso do #24" | Vale o João. O aviso de ocorrência é o do #24. O alerta aceito avisa por WebSocket no sentido da ADR_051, sem evento novo. | decisoes-joao-llt201.md, pendência 1, formulário 1; respostas-pendentes-llt201.md, conflito 1 |
| 009-5 revisto (gravidade no aviso) × spec 015, FR-007 (aviso sem gravidade) | Vale o João e a pendência 2 C. A emenda da 015 é T034 e não é escrita aqui. | respostas-pendentes-llt201.md, pendência 2 e conflito 2; origem: Ocorrência.mmd + 009-5 revisto |
| ADR_062 (escala "critico"/"alto", matriz pendente) × gravidade = classe | Vale o João: pendência 3 D, gravidade = classe ERROR, WARNING ou INFO no campo `gravidade`. | decisoes-joao-llt201.md, pendência 3, formulário 1; respostas-pendentes-llt201.md, conflito 3 |

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: sem dependência de código novo.
- **Foundational (Phase 2)**: depende do Setup. Bloqueia as histórias no sentido de preservar o cliente antes de configurar ambiente ou homologar.
- **User Stories (Phase 3+)**: dependem da Phase 2. US1, US2 e US3 são P1. US1 (token) não depende dos cinco cenários. A chamada local de US2 e US3 usa a URL `https://` do Caddy (T014, T017) e `SSL_CERT_FILE` (T015). US3 usa os tokens de US1 e a URL de US2.
- **Polish (Phase 6)**: T031 pode rodar sem as histórias. T032 espera as histórias nas partes que cada uma cita. T024 a T030 e T034 não autorizam criar os documentos nesta lista. T033 nomeia `finder_alerts` e não trava as outras tarefas.

### User Story Dependencies

- **User Story 1 (P1)**: depois da Phase 2. Não depende de US2 nem de US3. T011 e T012 podem receber o mesmo JWT, com as duas audiências. A conferência do WebSocket não está nesta história.
- **User Story 2 (P1)**: depois da Phase 2. T016 não espera T014 nem T015. T017 aplica o Caddy e `SSL_CERT_FILE` já decididos em T014 e T015. T015 não altera código.
- **User Story 3 (P1)**: depois da Phase 2. T010 espera T011 (token) e não leva [P]. T021 espera T011 e, no teste local, T017, e usa T010 na conferência do 201. T022 espera T012 e, no teste local, T017. T023 espera T022 e confere o FR-007, inclusive o literal `RECURRENT` de T020. T018 não acrescenta retentativa em `recorrencia.go`. T019 não cria campo de gravidade no Finder.

### Within Each User Story

- T020 usa o literal `RECURRENT` e T033 usa a coleção `finder_alerts`. Nenhum dos dois inventa outro valor.
- Não há suíte nova antes da implementação: o plano não pede teste automatizado novo.
- Configuração de ambiente antes da chamada de homologação.
- História conferida pelo Independent Test da spec antes de tratar a seguinte como incremento.

### Parallel Opportunities

- T002 em paralelo com T001.
- T003, T004, T005, T006, T007 e T008 em paralelo entre si, depois do Setup.
- T009, T011, T012 e T013 em paralelo entre si. T011 e T012 gravam variáveis diferentes e podem receber o mesmo JWT. T010 não entra neste paralelo: fica na história 3 e espera T011.
- T014, T015 e T016 em paralelo. T017 aplica a URL do Caddy e `SSL_CERT_FILE`.
- T018, T019 e T020 em paralelo entre si, depois da Phase 2. T020 fornece o literal `RECURRENT` que T023 confere. T010 não é paralelo a T011.
- Depois da Phase 2, a configuração de US1 (T011, T012) e a recusa de `http://` (T016) podem andar juntas. A homologação local (T017, T021, T022, T023) usa a URL `https://` do Caddy.
- T024, T025, T026, T027, T028, T029, T030, T031, T033 e T034 em paralelo. T032 não, porque percorre o roteiro. T033 registra o nome `finder_alerts` que T010 e T032 conferem.

---

## Parallel Example: User Story 1

```bash
# Regra do token, já decidida:
Task: "T009 um JWT com as duas audiências nas duas variáveis"

# Configuração de ambiente, o mesmo JWT pode ocupar as duas:
Task: "T011 FINDER_CORE_JWT com aud do Módulo e lhamalog-core"
Task: "T012 FINDER_MODULO_OCORRENCIA_JWT pode ser o mesmo valor"
Task: "T013 conferir Bearer sem montar claims em confirmacao/http.go"
```

## Parallel Example: User Story 2

```bash
Task: "T014 Caddy no compose do Core; Finder só recebe a URL https"
Task: "T015 CA local só no teste, via SSL_CERT_FILE"
Task: "T016 recusa de http:// na prontidão do pipeline"
```

## Parallel Example: User Story 3

```bash
Task: "T018 POST que falha mantém o lote pendente e reenvia no reprocessamento"
Task: "T019 gravidade igual à classe, no documento e no aviso do #24"
Task: "T020 literal RECURRENT; precedência no FR-007, conferida em T023"
```

T010 não entra nesse paralelo. A conferência do canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, espera T011 (token das duas audiências) e fica na história 3, sem a marca [P] (FR-010, US3 cenário 5). (fonte: analyze-011.md, MEDIUM I4; spec.md, FR-010; decisoes-joao-llt201.md, Rodada 3)

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Phase 1 e Phase 2.
2. Phase 3, história 1. T011 e T012 podem colocar o mesmo JWT nas duas variáveis. O `aud` contém a audiência do Módulo e `lhamalog-core`. O segredo não entra no repositório. O cliente não inspeciona claims.
3. Validar a história 1 pelo Independent Test da spec: a rota de alerta e as rotas do Módulo aceitam esse token. Os cinco cenários do DN-009-2 não fazem parte desse teste.
4. No teste local, a URL é a `https://` do Caddy no compose do Core. Sem essa URL registrada pela spec 016, a história 1 não inventa host nem libera `http://`.

### Incremental Delivery

1. Setup + Foundational.
2. US1: um JWT com as duas audiências, cliente intacto.
3. US2: recusa de `http://` já na prontidão (T016). O apontamento (T017) usa o Caddy do Core e `SSL_CERT_FILE`.
4. US3: alerta, cinco cenários do DN-009-2, conferência no Core do FR-007 (literal `RECURRENT`), gravidade igual à classe, WebSocket `FINDER_ALERT` depois do 201 (T010, depois do token) e reenvio do POST no reprocessamento do lote. O Finder não publica o aviso e não calcula gravidade. A coleção conferida é `finder_alerts`.
5. Polish: suíte já existente e lista de documentos, sem escrevê-los aqui.

### Parallel Team Strategy

1. Todos fazem Phase 1 e Phase 2.
2. Em seguida:
   - Quem pega US1 faz T009, T011, T012 e T013. A conferência do WebSocket não está nesta história.
   - Quem pega US2 faz T016 e, com a URL do Caddy, T017. T014 e T015 não alteram código deste repositório.
   - Quem pega US3 espera os tokens e, no teste local, a URL do Caddy. T010 espera T011 e não leva [P]. T018 e T019 não geram código neste repositório. T020 confere o literal `RECURRENT`.

---

## Notes

- [P] = arquivos diferentes, sem depender de tarefa ainda aberta.
- História no rótulo só nas fases US1, US2 e US3.
- T020 confere `RECURRENT` e T033 registra `finder_alerts`. T010 está na história 3, sem [P], e espera o token de T011. T012 cita o item 5 do roteiro (JWT).
- Os sufixos de `event=occurrence.alert.*` permanecem [PENDENTE] no Core. Esta lista não inventa a string exata de `event=` e não bloqueia T020, T033 nem T010.
- Não recriar #19, #20, #21 nem #24. Não criar evento nem envelope próprio. Não editar spec 009, spec 015, constituição, README, `.env.example`, `docker-compose.yml`, ADR, Drive nem lhamalog-docs nesta lista.
- Nenhuma caixa desta lista está marcada.
- Identificadores Go permanecem os do código já existente. Não há tipo novo com fonte nesta feature (plan.md, Structure Decision; data-model.md).
