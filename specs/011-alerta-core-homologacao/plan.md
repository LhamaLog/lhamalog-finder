# Implementation Plan: Alerta ao Core e homologação Finder

**Branch**: `LLT-201-publicar-evento-de-ocorrencia-nova-ou-atualizada` | **Date**: 2026-11-26 (rodada 2); texto alinhado em 2026-11-27 (rodada 3, sem regerar o plano) | **Spec**: [specs/011-alerta-core-homologacao/spec.md](spec.md)

**Input**: Feature specification from `/specs/011-alerta-core-homologacao/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

No lhamalog-finder, a LLT-201 que falta é configurar e homologar o que o cliente já faz. O Finder continua entregando `ResumoAnalise` e consultando o achado pelas rotas do #20, e continua avisando a falha esgotada do OpenObserve em `POST /internal/v1/alertas/finder`. `FINDER_CORE_JWT` e `FINDER_MODULO_OCORRENCIA_JWT` podem carregar o mesmo JWT, emitido fora do código, com o segredo e o emissor do #20 e com as duas audiências no `aud` (a do Módulo e `lhamalog-core`). As duas bases permanecem `https://` no mesmo host. No teste local, esse host é a URL do Caddy no compose do Core, e a CA local entra só por `SSL_CERT_FILE`. Se o POST do alerta falha, o lote permanece pendente e o aviso volta no reprocessamento. A ocorrência, a gravidade (classe no campo `gravidade`) e o aviso do #24 ficam no Core (spec 016): um aviso por ocorrência por entrega, com a mudança mais forte. Depois da primeira gravação do alerta (201), a homologação confere no Core um documento na coleção `finder_alerts` e o aviso à tela no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, sem evento novo para a Notificação. A repetição (200) não avisa de novo. A precedência do aviso de ocorrência é a do FR-007, com literal `RECURRENT`. (fonte: spec.md, User Stories 1–3, FR-001 a FR-010 e SC-009; decisoes-joao-llt201.md, rota POST, JWT, TLS, evento, gravidade, pendências 1, 3, 5, 6, 13, 18 e 20, formulários 1 e 2, e Rodada 3; respostas-pendentes-llt201.md, pendências 2, 17 e 19)

Esta pasta não reimplementa lista, detalhe, consultas, resumos, ciclo de vida nem o canal do #24. Não cria evento nem envelope. Não edita a spec 009 nem a constituição. Não cria a coleção `finder_alerts` nem o literal `RECURRENT`: os dois são do Core. (fonte: decisoes-joao-llt201.md, base da branch, evento para a Notificação e Rodada 3; spec.md, Fora de escopo e FR-008)

## Technical Context

**Language/Version**: Go 1.27.4 (`go.mod`; README, pré-requisitos e imagem `golang:1.27.4`). (fonte: `go.mod`; README.md)

**Primary Dependencies**: as que o módulo já exige: `github.com/joho/godotenv` v1.5.1 e `github.com/redis/go-redis/v9` v9.23.0. HTTP do alerta e do Módulo usa a biblioteca padrão `net/http`. Indiretas já no `go.mod`: `github.com/cespare/xxhash/v2`, `go.uber.org/atomic`, `golang.org/x/sys`. Não entra dependência nova. Continuam fora do módulo, como o README lista: `go.mongodb.org/mongo-driver`, `go.opentelemetry.io/otel`, `golang.org/x/sync`, `github.com/testcontainers/testcontainers-go`. A CA local do teste usa `SSL_CERT_FILE`, que o pacote `crypto/x509` da biblioteca padrão já lê. (fonte: `go.mod`; README.md, seções Dependências e Restrições da constituição — candidatas ainda não adicionadas não são tratadas como presentes; decisoes-joao-llt201.md, pendência 18, formulário 2; respostas-pendentes-llt201.md, pendência 18)

**Storage**: este repositório não grava ocorrência nem o alerta. A fila de apontamento é Redis Streams; `FINDER_REDIS_STREAM` tem de ser `lhamalog:stream:logs`. O log é lido no OpenObserve. O Compose local sobe MongoDB 8.0.34, mas `internal/adapter/out/mongo` não abre conexão. A coleção `Ocorrencia` e o recibo idempotente são do Core, na pilha assumida mergeada. O alerta durável é a coleção Mongo `finder_alerts` no Core, chave `(tenantId, productId, batchId, motivo)`, distinta da chave do recibo. Esta feature não cria coleção, stream nem índice. (fonte: constituição, princípio IV e Restrições Técnicas; README.md, configuração e layout; contexto-llt201.md, §1 e §2, coleção `Ocorrencia` na spec 012; respostas-pendentes-llt201.md, pendência 7; origem: padrão do #20 + DF-FINDER §3; nome: decisoes-joao-llt201.md, Rodada 3)

**Testing**: `go test ./...`, `go test -race ./...` e `go vet ./...`, como o README e o `Makefile` (`make test`, `make lint`) já descrevem. A suíte padrão não exige Core no ar. A homologação desta feature é o roteiro ponta a ponta com o Core real (cinco cenários do DN-009-2, alerta, falha do POST, WebSocket do alerta e conferência no Core). Não há, nas fontes desta feature, uma suíte nova de teste automatizado. (fonte: README.md, seção Testes; spec.md, User Story 3, Roteiro de homologação e FR-006)

**Target Platform**: processo Go do Finder (`cmd/finder`), no mesmo desenho do README: binário nativo ou imagem Docker. O Core local da homologação responde em HTTP na porta 8080; o Finder só o alcança por `https://` no Caddy do compose do Core. (fonte: README.md, Compilar e executar; spec.md, User Story 2; decisoes-joao-llt201.md, pendência 13, formulário 2; respostas-pendentes-llt201.md, pendência 17)

**Project Type**: serviço com arquitetura hexagonal leve (`cmd/finder`, `internal/domain`, `internal/application`, `internal/adapter`, `internal/bootstrap`, `internal/config`). (fonte: README.md, abertura e Layout; constituição, Restrições Técnicas)

**Performance Goals**: a spec 011 não define vazão nova. Prazos já existentes e que esta feature não altera: `FINDER_CORE_TIMEOUT` padrão `2s`, que é também o teto em que a rota de alerta do Core responde; consulta ao Módulo `2s` e entrega `5s` quando o cliente não recebe outro prazo. O aviso "virou recorrente" no Core herda até 3 s depois do 201 da spec 015, SC-001. Os dois tetos valem em alvos diferentes e são o SC-009 da spec 011. Os resultados medíveis são SC-001 a SC-009. (fonte: spec.md, SC-009, Success Criteria e Assumptions; README.md, tabela de variáveis; `internal/adapter/out/modocorrencia/cliente.go`, prazos 2 s e 5 s; contexto-llt201.md, §3; respostas-pendentes-llt201.md, pendência 16; origem: spec 015, SC-001 + DF-FINDER §3 e §8)

**Constraints**: só base `https://` (sem usuário, consulta, fragmento ou caminho além de `/`); o cliente não segue redirecionamento; o Finder não emite nem inspeciona `iss`, `aud` ou `sub`; um JWT com as duas audiências pode ocupar as duas variáveis; segredo fora do repositório; CA local só no teste, via `SSL_CERT_FILE`, com a verificação ligada e sem código novo; sem edição de compose, spec 009, constituição, README ou `.env.example` nesta feature; pilha #19–#24 assumida mergeada antes da implementação; a rota de alerta no Core é a spec 016, não este repositório. (fonte: spec.md, FR-002, FR-004, FR-005 e Fora de escopo; `internal/adapter/out/confirmacao/http.go`; decisoes-joao-llt201.md, base da branch, JWT, TLS, pendências 18 e 20, formulários 2 e 1)

**Scale/Scope**: um processo Finder homologado contra o Core com a pilha mergeada e a rota da spec 016. O alerta operacional continua uma chamada depois de três tentativas de consulta ao OpenObserve, motivo único `consulta_openobserve_indisponivel`. Se essa chamada falha, o lote fica pendente e o aviso é reenviado no reprocessamento, sem número novo de tentativas dentro do POST. Não há, nas fontes, número de usuários, de lotes ou de instâncias para esta feature. (fonte: spec.md, FR-001, FR-009 e Key Entities; `internal/adapter/out/openobserve/recorrencia.go`; contexto-llt201.md, §3 e R9; respostas-pendentes-llt201.md, pendência 19; origem: DF-FINDER §2 e §3)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constituição lida: `.specify/memory/constitution.md`, versão 1.0.0, ratificada e emendada em 2026-10-29. Esta feature não edita esse arquivo. (fonte: spec.md, FR-008 e SC-006)

| Gate | Resultado antes da pesquisa | Fonte |
| --- | --- | --- |
| I. Finalidade: logs em ocorrências | Divergência registrada, sem emenda neste arquivo. O texto ainda diz que o Finder transforma o resultado em ocorrência. Para esta feature a ocorrência nasce no Módulo, pacote do Core, e o Finder entrega `ResumoAnalise`. O CRITICAL C1 permanece pendência: esta execução não emenda a constituição. O desenho não passa a fazer o Finder criar ocorrência para cumprir o texto antigo. | constituição, princípio I; spec.md, Divergências; contexto-llt201.md, R2; decisoes-joao-llt201.md, registro de que a Ocorrência fica dentro do Core; analyze-011.md, C1 |
| II. Origem dos logs no coletor | Passa. Esta feature não muda a origem OpenTelemetry nem a leitura no OpenObserve. | constituição, princípio II; spec.md, Fora de escopo |
| III. Categorias de ocorrência | Passa. Esta feature não cria classe fora de erro, alerta e informação. O motivo `consulta_openobserve_indisponivel` é o corpo do alerta HTTP já existente; não é uma quarta classe de ocorrência. A gravidade conferida no Core é a classe ERROR, WARNING ou INFO no campo `gravidade`. | constituição, princípio III; spec.md, FR-001 e FR-007; contexto-llt201.md, R9; decisoes-joao-llt201.md, pendência 3, formulário 1 |
| IV. Fila Redis e apontamento no OpenObserve | Passa. O fluxo normativo do Finder permanece. O canal `lhamalog:occurrence-events` é publicação do Core no #24; este repositório não o implementa. | constituição, princípio IV; contexto-llt201.md, §1; `git show` de `origin/LLT-137-atualizacao-em-tempo-real:src/main/resources/application.properties` (`lhamalog.occurrence.stream.channel`) |
| V. Incerteza sem invenção | Passa. As pendências do Finder que estavam abertas fecharam com fonte na rodada 2 (Caddy, `SSL_CERT_FILE`, reenvio no reprocessamento, um token com as duas audiências). Na rodada 3, a coleção é `finder_alerts` e o literal é `RECURRENT`. Os sufixos de `event=occurrence.alert.*` permanecem [PENDENTE] no Core e não são inventados. Segredo não entra no repositório. | constituição, princípio V; spec.md, Pendências; decisoes-joao-llt201.md, formulários 1 e 2 e Rodada 3; respostas-pendentes-llt201.md, pendências 7, 10, 17, 18, 19 e 20 |
| Restrições técnicas e README | Passa. Linguagem, hexagonal e dependências citadas são as do repositório. Nenhuma candidata do README é tratada como já presente. `SSL_CERT_FILE` é variável que o `crypto/x509` já lê; não é dependência nova. | constituição, Restrições Técnicas; README.md, Dependências; `go.mod`; respostas-pendentes-llt201.md, pendência 18 |
| Evidência e consulta | Passa. Fato, premissa (pilha mergeada) e resíduo sem fonte estão separados. Nenhum teste ou servidor foi executado nesta geração de plano. | constituição, Evidência e Consulta; spec.md, Evidência e consulta |

**Reavaliação após a Fase 1**: `research.md`, `data-model.md`, `contracts/` e `quickstart.md` não acrescentam classe de ocorrência, dependência, persistência no Finder, emissão de JWT no código, evento novo nem edição da constituição. Registram o token único com as duas audiências, o par `OCCURRENCE_INTERNAL_JWT_SECRET` / `OCCURRENCE_INTERNAL_JWT_SECRETS` sem unificar os nomes, o Caddy, `SSL_CERT_FILE`, o reenvio do alerta no reprocessamento, a gravidade igual à classe, a coleção `finder_alerts`, o literal `RECURRENT`, o aviso no canal `/ws/occurrences` com tipo `FINDER_ALERT` e corpo `{tenantId, productId, batchId, motivo}` só depois do 201, e os dois tetos do SC-009 (rota em 2 s; aviso "virou recorrente" em até 3 s depois do 201). O `quickstart.md` já é o roteiro. Os sufixos de `event=` continuam [PENDENTE]. O gate I (CRITICAL C1) continua a mesma divergência justificada abaixo, sem emenda. Os gates II a V e as restrições técnicas permanecem como na tabela. Os conflitos ADR_165, spec 015 FR-007 e ADR_062 ficam só registrados; não são resolvidos editando esses textos. (fonte: artefatos desta pasta; spec.md, FR-002, FR-007, FR-008, FR-010 e SC-009; decisoes-joao-llt201.md, Rodada 3; analyze-011.md, C1, MEDIUM I1 e MEDIUM I2)

## Project Structure

### Documentation (this feature)

```text
specs/011-alerta-core-homologacao/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/           # Phase 1 output (/speckit-plan command)
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

A spec e o checklist já existiam antes deste plano. `tasks.md` já existia de uma execução anterior e não é gerado nem marcado nesta execução.

### Source Code (repository root)

```text
cmd/finder/                                      # serviço já existente
internal/domain/
internal/application/handoff/                    # entrega do ResumoAnalise já existente
internal/application/port/out/
internal/adapter/in/queue/                       # Redis Streams
internal/adapter/out/core/aviso.go               # POST /internal/v1/alertas/finder já existente
internal/adapter/out/confirmacao/http.go         # HTTPS, Bearer, sem redirect
internal/adapter/out/modocorrencia/cliente.go    # consultas e resumos já existentes
internal/adapter/out/openobserve/recorrencia.go  # 3 tentativas e um alerta
internal/bootstrap/pipeline.go                   # lê FINDER_CORE_URL e FINDER_CORE_JWT
internal/config/
docker-compose.yml                               # repassa as variáveis; esta feature não edita
```

**Structure Decision**: o trabalho com fonte nesta feature fica na configuração de ambiente e na homologação dos clientes já presentes em `internal/adapter/out/core`, `internal/adapter/out/modocorrencia` e `internal/bootstrap`. Não há pacote novo, nem alteração desses arquivos, com fonte nesta spec: o cliente envia o JWT pronto e não monta claims. As duas variáveis de JWT podem receber o mesmo valor. `SSL_CERT_FILE` é ambiente do processo de teste, sem código novo no cliente. O Caddy e o compose ficam na spec 016; `docker-compose.yml` deste repositório só interpola `FINDER_CORE_*` e `FINDER_MODULO_OCORRENCIA_*`. (fonte: spec.md, User Story 1 cenário 4, FR-001, FR-002 e FR-005; `internal/bootstrap/pipeline.go`; `docker-compose.yml`; contexto-llt201.md, §3; decisoes-joao-llt201.md, pendências 13, 18 e 20, formulários 2 e 1)

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
| --- | --- | --- |
| Princípio I ainda atribui ao Finder a transformação do resultado em ocorrência. A emenda não é feita aqui. | A decisão de 26/11 e o R2 colocam a ocorrência no Módulo, que é o pacote `br.com.lhamalog.core.occurrence` dentro do Core. O Finder entrega `ResumoAnalise`. A spec 009, item 10, já deixa a emenda em aberto. | Fazer o Finder gravar a ocorrência contradiz o R2, a ADR_062 citada nesse requisito e a decisão de registrar o Módulo dentro do Core. Editar a constituição nesta pasta está fora do pedido. (fonte: constituição, princípio I; spec.md, Divergências; contexto-llt201.md, R2 e §8; decisoes-joao-llt201.md, registro de que a Ocorrência fica dentro do Core) |

## O que a implementação deste repositório cobre

1. Colocar em `FINDER_CORE_JWT` e, se for o caso, o mesmo valor em `FINDER_MODULO_OCORRENCIA_JWT`: um JWT HS256 emitido fora do Finder, com o segredo e o emissor do #20, `aud` contendo a audiência do Módulo e `lhamalog-core`, e `sub` na lista que o #20 já configura (padrão `lhamalog-finder`). O código do cliente permanece enviando `Authorization: Bearer` e não inspeciona claims. Fato: `scripts/InternalJwt.java` gera uma audiência por `--aud` e lê `OCCURRENCE_INTERNAL_JWT_SECRET` (singular); a propriedade do Core é `OCCURRENCE_INTERNAL_JWT_SECRETS` (plural). Esta feature não unifica esses nomes e não copia o script. Fato: no #20, `InternalServiceSecurityConfiguration` cobre todo `/internal/v1/**` e exige que o `aud` contenha a audiência do Módulo; o `aud` é conferido por "contém"; um token só com `lhamalog-core`, sem o ajuste por rota, recebe 401. O ajuste (propriedade `core-audience` só na rota de alerta; rotas do Módulo com a audiência do #20) é da spec 016. (fonte: spec.md, FR-002, FR-003 e Assumptions; tasks.md, T011; quickstart.md; decisoes-joao-llt201.md, JWT da rota de alerta, pendências 6 e 20, formulário 1; respostas-pendentes-llt201.md, resumo, fato de código)
2. Manter as rotas `POST /internal/v1/ocorrencias/consultas` e `POST /internal/v1/resumos-analise` exigindo a audiência `lhamalog-modulo-ocorrencia`. A audiência `lhamalog-core` não substitui essa. (fonte: spec.md, FR-003; contexto-llt201.md, §5 conflito 6; `application.properties` do #20; decisoes-joao-llt201.md, pendência 6, formulário 1)
3. Apontar `FINDER_CORE_URL` e `FINDER_MODULO_OCORRENCIA_URL` para `https://` do mesmo host. No teste local, esse host é a URL do Caddy no compose do Core, na frente do Core na porta 8080. A CA local entra só no teste, via `SSL_CERT_FILE`, com a verificação ligada. Este repositório não edita compose. (fonte: spec.md, FR-004, FR-005 e User Story 2; decisoes-joao-llt201.md, pendências 13 e 18, formulário 2; respostas-pendentes-llt201.md, pendência 17; origem: ADR_164 + arquitetura-modulos)
4. Homologar os cinco cenários do DN-009-2 com o cliente já existente e, no Core, conferir o FR-007 (precedência, gatilho de `before.recurring` / `after.recurring` e literal `RECURRENT`) e a gravidade igual à classe no campo `gravidade`. O Finder não publica esse aviso nem calcula gravidade. Os dois tetos medidos são o SC-009, em alvos diferentes: a rota de alerta responde dentro de 2 s (`FINDER_CORE_TIMEOUT`) e o aviso "virou recorrente" sai em até 3 s depois do 201. O Finder não implementa o teto de 3 s. O roteiro dessa medição é o `quickstart.md` já gerado. (fonte: spec.md, FR-006, FR-007, SC-009 e Roteiro de homologação; `specs/009-publicar-evento-ocorrencia/contracts/proposta-dn-009-2.md`, seção Cenários de homologação com o Módulo; decisoes-joao-llt201.md, pendências 3 e 5, formulário 1, e Rodada 3; respostas-pendentes-llt201.md, pendências 2, 5, 11 e 16; origem: spec 015, SC-001 + DF-FINDER §3 e §8; analyze-011.md, MEDIUM I1 e MEDIUM I2)
5. Se o POST do alerta falha, manter o lote pendente e reenviar o aviso no reprocessamento do lote, com a mesma identidade. Sem número novo de tentativas dentro do POST. Corpo incompleto, campo a mais ou motivo fora da lista é, no Core, 422 `DADOS_INVALIDOS`. JSON ilegível é 400 `REQUISICAO_INVALIDA`. Produto que não pertence ao tenant é 403 `ACESSO_NEGADO`, sem gravar o alerta. O cliente trata esses status como falha do POST. Esta feature não cria esses códigos. (fonte: spec.md, FR-001 e FR-009; respostas-pendentes-llt201.md, pendências 4, 12 e 19; origem: DF-FINDER §2 e §3 + padrão do #20; decisoes-joao-llt201.md, pendência 12, formulário 1; elaboração do 403: respostas-pendentes-llt201.md, pendência 12; origem: spec 013, FR-016)
6. Na homologação, depois da primeira gravação (201), conferir no Core um documento em `finder_alerts` e o aviso à tela no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, sem evento novo para a Notificação e sem publicar em `lhamalog:occurrence-events` nem em `lhamalog:notifications:events`. A repetição (200) não avisa de novo. Quem perde a corrida de duas primeiras gravações da mesma chave relê e responde 200 com a mesma identidade; o cliente trata 200 e 201 como sucesso. O Finder não publica esse aviso. (fonte: spec.md, FR-001 e FR-010; decisoes-joao-llt201.md, pendência 1 e pendência 9, formulário 1, e Rodada 3)
7. Deixar registrado nesta pasta que o Módulo é o pacote dentro do Core e que o HTTP Módulo → Core da spec 009 deixa de descrever o processo único. (fonte: spec.md, FR-008 e Responsabilidades)

## Fora deste plano

Lista e detalhe (#19, spec 012), consultas e resumos com recibo e JWT interno (#20, spec 013), ciclo de vida (#21, spec 014) e o canal `lhamalog:occurrence-events` com WebSocket `/ws/occurrences` (#24, spec 015). A rota `POST /internal/v1/alertas/finder` no servidor, a coleção `finder_alerts`, a deduplicação durável, o literal `RECURRENT`, o cálculo da gravidade e o aviso `FINDER_ALERT` no canal `/ws/occurrences` são a spec 016 do lhamalog-core. O código dessa rota fica no pacote `finderalert` (exemplo `br.com.lhamalog.core.finderalert`); a segurança continua na cadeia `/internal/v1/**`. O log de negócio dessa rota usa o prefixo `occurrence.alert.*`. Os sufixos de `event=` permanecem [PENDENTE] no Core. (fonte: decisoes-joao-llt201.md, base da branch, evento, mudanças, gravidade, pendências 1, 14 e 15, formulários 1 e 2, e Rodada 3; contexto-llt201.md, §1; spec.md, Fora de escopo; analyze-016.md, MEDIUM U2)

## Conflitos só registrados

Estes conflitos não são resolvidos de outro modo e não editam constituição, spec 015 nem ADR.

| Conflito | O que vale nesta feature | Fonte |
| --- | --- | --- |
| ADR_165 (ocorrências pela Notificação via stream) × "só o aviso do #24" | Vale o João. O aviso de ocorrência é o do #24. O alerta aceito avisa por WebSocket no sentido da ADR_051, sem evento novo. | decisoes-joao-llt201.md, pendência 1, formulário 1; respostas-pendentes-llt201.md, conflito 1 |
| 009-5 revisto (gravidade no aviso) × spec 015, FR-007 (aviso sem gravidade) | Vale o João e a pendência 2 C. A emenda da 015 fica listada e não é escrita aqui. | respostas-pendentes-llt201.md, pendência 2 e conflito 2; origem: Ocorrência.mmd + 009-5 revisto |
| ADR_062 (escala "critico"/"alto", matriz pendente) × gravidade = classe | Vale o João: pendência 3 D, gravidade = classe ERROR, WARNING ou INFO no campo `gravidade`. | decisoes-joao-llt201.md, pendência 3, formulário 1; respostas-pendentes-llt201.md, conflito 3 |

## Pendências — documentos necessários

Nenhum destes arquivos é criado neste plano.

| Documento | Onde | Por quê |
| --- | --- | --- |
| Emenda do princípio I | `.specify/memory/constitution.md` | O princípio I ainda diz que o Finder transforma o resultado em ocorrência. A spec 009, item 10, já aponta essa emenda. (fonte: constituição, princípio I; spec.md, Divergências; spec 009, questões em aberto, item 10, citado em spec.md) |
| Atualização da spec 009, inclusive `contracts/pipeline-aviso-referencia.md`, o status Draft, a questão 3, o item 10, o SC-002 e o SC-006 daquela spec | `specs/009-publicar-evento-ocorrencia/` | O Módulo ficou dentro do Core; o aviso HTTP Módulo → Core deixa de descrever o processo único; o cliente de alerta já tem rota e corpo; os termos da homologação desta feature são nova, reaberta e virou recorrente, um aviso por ocorrência por entrega, conferidos no Core; o alerta aceito avisa por WebSocket. (fonte: spec.md, Divergências; contexto-llt201.md, §5 conflito 2 e §8; decisoes-joao-llt201.md, evento para a Notificação e pendência 1, formulário 1) |
| ADR nova "Módulo Ocorrência integrado ao Core" | Drive e branch nova no lhamalog-docs | Registrar que a Ocorrência fica dentro do Core e substituir a parte da ADR_054 que a trata como serviço separado. (fonte: decisoes-joao-llt201.md, registro de que a Ocorrência fica dentro do Core; origem: contexto-llt201.md, §5 conflito 1 e §7 Q8 opção A) |
| Fechamento do trecho pendente da ADR_051 | Drive | A ADR_051 ainda deixa pendente o protocolo Finder → Core. O contrato já decidido, que a ADR pode registrar, é: depois do 201, canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`. Este plano não edita a ADR. (fonte: spec.md, Pendências — documentos necessários; decisoes-joao-llt201.md, pendência 1, formulário 1, e Rodada 3) |
| README e `.env.example` do Finder | Raiz deste repositório, numa execução que possa editá-los | Dizer que `FINDER_CORE_JWT` e `FINDER_MODULO_OCORRENCIA_JWT` podem carregar o mesmo JWT do segredo e do emissor do #20, com as duas audiências (`lhamalog-modulo-ocorrencia` e `lhamalog-core`), que as duas URLs `https://` apontam para o mesmo host, e que no teste a CA local entra só por `SSL_CERT_FILE`. Sem segredo no exemplo. (fonte: spec.md, Pendências — documentos necessários; decisoes-joao-llt201.md, pendências 18 e 20, formulários 2 e 1) |
| README, exemplo de ambiente e compose do Caddy | lhamalog-core, spec 016 | O proxy é o Caddy no compose do Core. O Finder só consome a URL `https://`. O arquivo não é escrito aqui. (fonte: decisoes-joao-llt201.md, pendência 13, formulário 2; respostas-pendentes-llt201.md, pendência 17) |
| Emenda da spec 015, FR-007 | lhamalog-core, spec 015 | O aviso do #24 passa a levar a gravidade. Esta execução não edita a 015. (fonte: respostas-pendentes-llt201.md, pendência 2 e conflito 2; origem: 009-5 revisto) |
| Alinhamento de `endpoints-ocorrencia.md`, `endpoints-finder.md` e `endpoints-core-notificacao.md` | lhamalog-docs | Citar o DN-009-2 e a rota de alerta quando essa documentação for atualizada. (fonte: contexto-llt201.md, §8) |

O `quickstart.md` desta pasta é o roteiro de homologação já gerado. A spec aponta para esse arquivo. Não substitui o README. (fonte: spec.md, Roteiro de homologação; analyze-011.md, MEDIUM I2)

## Evidência e consulta

- **Lido**: `decisoes-joao-llt201.md` (a seção Rodada 3 prevalece nos três pontos que estavam sem fonte: coleção `finder_alerts`, literal `RECURRENT` e aviso WebSocket; formulários 1 e 2 continuam no restante); `prompts/r3-decisoes.txt`; `analyze-011.md`; `respostas-pendentes-llt201.md`; `prompts/r2-decisoes.txt`; `.specify/memory/constitution.md` (só leitura); spec 011 já atualizada na rodada 3; `go.mod`; README.md; `docker-compose.yml` (variáveis `FINDER_CORE_*` e `FINDER_MODULO_OCORRENCIA_*`); `internal/adapter/out/core/aviso.go`; `internal/adapter/out/confirmacao/http.go`; `internal/adapter/out/modocorrencia/cliente.go`; `internal/adapter/out/openobserve/recorrencia.go`; `internal/bootstrap/pipeline.go`; cenários em `specs/009-publicar-evento-ocorrencia/contracts/proposta-dn-009-2.md`. No Core, só leitura, sem edição: spec `016-alerta-finder-evento-ocorrencia` e o fato de código citado em `respostas-pendentes-llt201.md` (`InternalServiceSecurityConfiguration`, `aud` por "contém", `scripts/InternalJwt.java` com um `--aud` e `OCCURRENCE_INTERNAL_JWT_SECRET` no singular). Os `git show` da geração anterior do plano (propriedades, filtro e publicador do #20/#24) permanecem a base factual onde esta rodada não os contradiz.
- **Executado**: nenhum teste, nenhum servidor, nenhum compose, nenhum `git checkout` / `git switch` / commit / push. Nenhum script PowerShell. A branch não foi trocada.
- **Em aberto**: sufixos de `event=occurrence.alert.*`, [PENDENTE] residual no Core, sem invenção da string exata de `event=`. Emenda do princípio I (CRITICAL C1), só registrada. Documentos da tabela acima.
