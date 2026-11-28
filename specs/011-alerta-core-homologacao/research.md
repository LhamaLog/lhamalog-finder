# Research: Alerta ao Core e homologação Finder

Fase 0 do plano em `specs/011-alerta-core-homologacao/plan.md`, atualizada na rodada 3 (27/11) sem regerar o plano. Cada decisão cita a fonte. Onde a fonte não escolhe o detalhe concreto, o item permanece [PENDENTE] residual, sem opção nova e sem valor inventado. Não há [NEEDS CLARIFICATION] aberto.

Premissa de base, já decidida: a implementação espera o merge de #19, #20, #21 e #24. Este texto não redesenha o que essa pilha entrega. (fonte: decisoes-joao-llt201.md, base da branch)

## D1 — Escopo no Finder

- **Decision**: o acréscimo da LLT-201 neste repositório é a configuração do JWT (um token com as duas audiências, nas duas variáveis se for o caso) e das duas URLs `https://`, mais a homologação do cliente já existente. O Finder não cria ocorrência, não calcula gravidade e não publica evento de ocorrência.
- **Rationale**: o cliente do DN-009-2 já está na Development (PR #13, merge em 23/11). A rota de alerta no Core, o aviso, a gravidade e o WebSocket do alerta são a spec 016. (fonte: spec.md, User Story 3, FR-006, FR-007 e FR-010; contexto-llt201.md, §1 e §3; decisoes-joao-llt201.md, evento para a Notificação, mudanças que geram evento, gravidade e pendência 1, formulário 1)
- **Alternatives considered**: reimplementar consulta e entrega (recusado: spec.md, Fora de escopo); fazer o Finder publicar no canal do #24 (recusado: decisoes-joao-llt201.md, evento para a Notificação — só o aviso do #24, no Core).

## D2 — JWT do alerta, sem código que emita token

- **Decision**: `FINDER_CORE_JWT` é um JWT HS256 emitido fora do processo Finder. Segredo: o de `lhamalog.occurrence.internal-auth.secrets` (`OCCURRENCE_INTERNAL_JWT_SECRETS`). Emissor: o de `lhamalog.occurrence.internal-auth.issuer` (padrão `lhamalog-internal`). O `aud` contém a audiência do Módulo (padrão `lhamalog-modulo-ocorrencia`) e `lhamalog-core`. `FINDER_MODULO_OCORRENCIA_JWT` pode carregar esse mesmo JWT. `sub` permanece o que o #20 já autoriza (padrão `lhamalog-finder`). O cliente só copia o valor para `Authorization: Bearer` e não inspeciona claims.
- **Rationale**: a pendência 20 fecha um token com as duas audiências. A pendência 6 fecha a audiência `lhamalog-core` numa propriedade `core-audience` usada só na rota de alerta; as rotas do Módulo seguem com a audiência do #20. Ilustração de caminho, não nome extra decidido: `lhamalog.occurrence.internal-auth.core-audience`. O filtro do #20 valida HS256, `iss`, `aud` como coleção que contém a audiência configurada, e `sub` na lista. Fato de código: `InternalServiceSecurityConfiguration` cobre todo `/internal/v1/**` e exige que o `aud` contenha `lhamalog.occurrence.internal-auth.audience` (padrão `lhamalog-modulo-ocorrencia`); o `aud` é conferido por "contém"; um token só com `lhamalog-core`, sem o ajuste por rota da pendência 6, recebe 401. `scripts/InternalJwt.java` gera um token com uma única audiência (`--aud`). Esta feature não altera o script nem manda o cliente inspecionar o token. Uma invocação com um único `--aud` não é, por si, o token da pendência 20. O script lê o segredo de `OCCURRENCE_INTERNAL_JWT_SECRET` (nome no singular, só no script) e não grava o segredo em argumento. O Finder não chama esse script. (fonte: decisoes-joao-llt201.md, JWT da rota de alerta, pendências 6 e 20, formulário 1; respostas-pendentes-llt201.md, resumo, fato de código, e pendência 20; spec.md, FR-002, FR-003 e Assumptions; `internal/adapter/out/confirmacao/http.go`)
- **Alternatives considered**: o Finder passar a montar o JWT (recusado: spec.md, FR-002); dois tokens, um por audiência (recomendação A do Especialista na pendência 20, recusada: decisoes-joao-llt201.md, pendência 20, formulário 1); usar só a audiência `lhamalog-core` nas duas variáveis (as rotas do Módulo deixam de autenticar; spec.md, FR-003).

O nome de ambiente do script (`OCCURRENCE_INTERNAL_JWT_SECRET`) e o da propriedade do Core (`OCCURRENCE_INTERNAL_JWT_SECRETS`) são os que os arquivos mostram. Não há fonte neste repositório que unifique os dois nomes. O comentário do script diz que o valor é o mesmo segredo do Core. (fonte: `scripts/InternalJwt.java` e `application.properties` do #20, via `git show`, preservado da geração anterior; respostas-pendentes-llt201.md, resumo, fato de código)

## D3 — Audiência das rotas do #20 permanece

- **Decision**: as rotas `POST /internal/v1/ocorrencias/consultas` e `POST /internal/v1/resumos-analise` continuam exigindo a audiência `lhamalog-modulo-ocorrencia`. A rota de alerta exige `lhamalog-core`, lida da propriedade `core-audience`, usada só nessa rota. A audiência `lhamalog-core` não substitui a das rotas do Módulo. O mesmo JWT, com as duas audiências no `aud`, autentica os dois papéis porque o #20 confere `aud` por "contém", depois do ajuste por rota.
- **Rationale**: o padrão do #20 é a audiência do Módulo, e o filtro atual cobre `/internal/v1/**` com essa audiência. Aceitar `lhamalog-core` só na rota de alerta, sem trocar a audiência das outras duas rotas, é mudança do Core (spec 016). Este plano não desenha esse filtro. Sem o ajuste, um token só com `lhamalog-core` recebe 401. (fonte: spec.md, FR-003; contexto-llt201.md, §5 conflito 6; decisoes-joao-llt201.md, pendências 6 e 20, formulário 1; respostas-pendentes-llt201.md, resumo, fato de código)

## D4 — HTTPS e Caddy, sem liberar HTTP

- **Decision**: `FINDER_CORE_URL` e `FINDER_MODULO_OCORRENCIA_URL` continuam `https://`, nas regras que `confirmacao.Novo` e `NovoCliente` já aplicam. Base `http://` segue recusada. No teste local as duas apontam para a URL `https://` do Caddy no compose do Core, que permanece em HTTP na porta 8080, no mesmo host, porque o Módulo é pacote do Core. O proxy, o produto e o arquivo são os que a spec 016 registra. Este repositório não ganha serviço de proxy e não edita `docker-compose.yml`. O cliente não segue redirecionamento (`CheckRedirect` devolve `http.ErrUseLastResponse`). A CA local entra só no teste, pela variável `SSL_CERT_FILE`, sem código novo no cliente e sem desligar a verificação. Não há variável de certificado própria do Finder: o `crypto/x509` do Go já lê `SSL_CERT_FILE`. O caminho do arquivo é o que o Caddy da spec 016 expõe; este plano não nomeia outro caminho.
- **Rationale**: a pendência 13 escolhe Caddy no compose do Core. A pendência 17 confirma que o Finder só recebe a URL `https://`. A pendência 18 escolhe a CA local só no teste, via `SSL_CERT_FILE`. (fonte: decisoes-joao-llt201.md, TLS no teste local, pendências 13 e 18, formulário 2; respostas-pendentes-llt201.md, pendências 17 e 18; origem: ADR_164 + arquitetura-modulos, e o pacote `crypto/x509` do Go; spec.md, FR-004, FR-005 e User Story 2; `internal/adapter/out/confirmacao/http.go`; `internal/adapter/out/modocorrencia/cliente.go`; `docker-compose.yml`)
- **Alternatives considered**: aceitar `http://` no cliente (opção B da Q7 do contexto, recusada pela decisão de TLS); encerrar TLS no processo do Core sem proxy (a decisão escolhe proxy no compose); desligar a verificação do certificado (a pendência 18 escolhe `SSL_CERT_FILE` com a verificação ligada).

## D5 — Contrato do alerta já está no cliente

- **Decision**: o pedido permanece `POST /internal/v1/alertas/finder` com `{tenantId, productId, batchId, motivo}`. O único motivo que o Finder envia é `consulta_openobserve_indisponivel`. Sucesso para o cliente: 200 ou 201, corpo com `aceito: true` e a mesma identidade, sem campo a mais. A primeira aceitação no Core é 201; a repetição da mesma identidade é 200, sem campo `repeticao`. Campo a mais, nome diferente, identidade diferente, 202 ou corpo vazio não são aceite. A deduplicação durável é do Core, na coleção `finder_alerts`, cuja chave é `(tenantId, productId, batchId, motivo)`. Quem perde a corrida de duas primeiras gravações relê e responde sucesso (200 com a mesma identidade); o cliente já trata 200 e 201 como sucesso. Se o POST falha, o lote permanece pendente e o aviso é reenviado no reprocessamento do lote, com a mesma identidade. Não há número novo de tentativas dentro do POST. Status diferente de 200/201 é falha para o cliente, inclusive 422 `DADOS_INVALIDOS` (corpo incompleto, campo a mais ou motivo fora da lista), 400 `REQUISICAO_INVALIDA` (JSON ilegível) e 403 `ACESSO_NEGADO` (produto que não pertence ao tenant, sem gravar o alerta). São respostas do Core e não códigos novos deste repositório.
- **Rationale**: o corpo e o aceite estrito já estão em `aviso.go`, `aviso_test.go` e no R9. A spec proíbe reimplementar o cliente. O HTTP 201/200 sem `repeticao` é a pendência 8, porque o cliente recusa campo a mais. A falha que mantém o lote pendente é a pendência 19. (fonte: spec.md, FR-001, FR-009 e Edge Cases; contexto-llt201.md, §3 e R9; `internal/adapter/out/core/aviso.go`; `internal/adapter/out/confirmacao/http.go`; respostas-pendentes-llt201.md, pendências 4, 7, 8, 9 e 19; origem: DF-FINDER §2 e §3 + controller do #20 + cliente estrito; decisoes-joao-llt201.md, pendências 9 e 12, formulário 1)
- **Alternatives considered**: novo corpo ou novo motivo (sem fonte; spec.md, FR-001); campo `repeticao` no corpo do alerta (recusado: respostas-pendentes-llt201.md, pendência 8 — o cliente estrito recusa campo a mais); repetir o POST um número novo de vezes dentro da chamada (recusado: spec.md, FR-009).

## D6 — DN-009-2 não é redesenhado

- **Decision**: a homologação usa o contrato e os cinco cenários já escritos em `specs/009-publicar-evento-ocorrencia/contracts/proposta-dn-009-2.md`. O cliente é `internal/adapter/out/modocorrencia`. Não há segunda versão do contrato nesta pasta. O recibo do resumo continua o do #20: 201 na primeira entrega, 200 com `repeticao: true` no reenvio, `409 REVISAO_DESATUALIZADA` e `409 ENTREGA_EM_ANDAMENTO`. Esse campo `repeticao` é do recibo do resumo. Não entra no corpo do alerta (D5).
- **Rationale**: a spec 011 diz que o DN-009-2 permanece o da 009 e que o #20 já entrega o lado servidor. (fonte: spec.md, abertura e FR-006; contexto-llt201.md, §1; `proposta-dn-009-2.md`, seção Cenários de homologação com o Módulo; respostas-pendentes-llt201.md, pendência 8, que distingue o cliente estrito do alerta)

## D7 — Aviso de ocorrência e gravidade ficam no Core

- **Decision**: o Finder só entrega o `ResumoAnalise`. A homologação confere no Core um aviso por ocorrência por entrega, com a mudança mais forte: nova > reaberta > virou recorrente. Criação já recorrente publica só `CREATED`. FECHADA que reabre e fica recorrente na mesma entrega publica só `REOPENED`. O aviso "virou recorrente" dispara quando `before.recurring` não é TRUE e `after.recurring` é TRUE; ausente ou `null` conta como `false`, e essa mudança só sai se for a mais forte da entrega. A constante é inglesa, no enum `Change`, no estilo de `CREATED` e `REOPENED`. O literal é `RECURRENT`. Não usar `RECURRING`. A precedência normativa está no FR-007. A gravidade é a classe ERROR, WARNING ou INFO copiada para o campo `gravidade`, com o tipo mantido, no documento da ocorrência e no aviso do #24. O canal reutilizado é o do #24: pub/sub `lhamalog:occurrence-events` e WebSocket `/ws/occurrences`. O aviso "virou recorrente" sai em até 3 s depois do 201. Esse teto, com o de 2 s da rota de alerta, é o SC-009 (D9). Este repositório não publica a mensagem e não calcula a gravidade.
- **Rationale**: a decisão manda usar só o aviso do #24 e calcular a gravidade agora, no Módulo dentro do Core, como a classe. A forma atual no código do #24 ainda é `tenantId`, `productId`, `occurrenceId`, `occurrenceType`, `change` (`CREATED` ou `REOPENED`) e `at`, sem gravidade. A spec 015, FR-007, ainda descreve o aviso sem gravidade; a emenda fica pendente e não é escrita aqui. O literal fechado na rodada 3 é `RECURRENT`. (fonte: decisoes-joao-llt201.md, evento, mudanças, gravidade, pendências 3 e 5, formulário 1; respostas-pendentes-llt201.md, pendências 2, 5, 10, 11 e 16; origem: Ocorrência.mmd + 009-5 revisto, constituição do Core princípio IV, enum `Change` do #24, data-model da spec 012, spec 015 SC-001; spec.md, FR-007; `git show` de `OccurrenceEvent.java` e `RedisOccurrenceEventPublisher.java` na branch `LLT-137-atualizacao-em-tempo-real`, fato do código atual)
- **Alternatives considered**: envelope ADR_146 ou stream `lhamalog:notifications:events` (a decisão fecha o aviso do #24); o Finder calcular gravidade (recusado: contexto-llt201.md, R2; spec.md, FR-007); adotar o literal `RECURRING` (recusado: decisoes-joao-llt201.md, Rodada 3 — o literal é `RECURRENT`).

Há diferença de nome entre o research da spec 015 (`type` no JSON do canal) e o `toMessage` do publicador (`occurrenceType`). O plano do Finder não resolve esse nome: o Finder não publica a mensagem. (fonte: `specs/015-tempo-real-ocorrencias/research.md` na branch `LLT-137-atualizacao-em-tempo-real`, via `git grep` na geração anterior; `RedisOccurrenceEventPublisher.toMessage`, via `git show`)

## D8 — Mesmo host, sem validação nova de host

- **Decision**: a homologação configura as duas URLs no mesmo host. O cliente de hoje aceita hosts diferentes. Não há mudança de código para recusar hosts diferentes.
- **Rationale**: a spec registra a ausência de fonte para essa recusa. (fonte: spec.md, Edge Cases e FR-005; contexto-llt201.md, §8)

## D9 — Prazo do alerta e prazo do aviso

- **Decision**: `FINDER_CORE_TIMEOUT` permanece, padrão `2s`. Esta feature não muda o prazo. Esse é o teto em que a rota de alerta do Core responde. O aviso "virou recorrente" herda até 3 s depois do 201 da spec 015, SC-001. Os dois valem, em alvos diferentes, e são o SC-009. O plano aponta para esse critério. O Finder não implementa o teto de 3 s. O roteiro que mede os dois tetos é o `quickstart.md` já gerado.
- **Rationale**: o prazo do cliente já está documentado. A pendência 16 separa os alvos. A spec 011 registra os dois no SC-009. (fonte: spec.md, SC-009, Assumptions e Roteiro de homologação; README.md, tabela de variáveis; contexto-llt201.md, §3; respostas-pendentes-llt201.md, pendência 16; origem: spec 015, SC-001 + DF-FINDER §3 e §8; analyze-011.md, MEDIUM I1 e MEDIUM I2)

## D10 — Constituição, princípio I

- **Decision**: o plano registra a divergência e não edita `.specify/memory/constitution.md`. A emenda é documento pendente.
- **Rationale**: a constituição manda registrar emenda de finalidade neste arquivo antes de orientar código; o pedido desta feature proíbe essa edição e manda listar a emenda. Enquanto a emenda não existir, o trabalho desta spec segue a decisão de 26/11 e o R2, e não faz o Finder criar ocorrência. (fonte: constituição, princípio I e Governance; spec.md, Divergências e FR-008; decisoes-joao-llt201.md, registro de que a Ocorrência fica dentro do Core)

## D11 — Pacote do Módulo

- **Decision**: o registro nesta pasta usa o nome já declarado na pilha: `br.com.lhamalog.core.occurrence`, dentro do Core, no mesmo processo. As duas URLs do Finder apontam para o mesmo host por isso. O código da rota de alerta, no Core, fica no pacote `finderalert` (exemplo `br.com.lhamalog.core.finderalert`). A segurança continua na cadeia `/internal/v1/**` do #20. O log de negócio dessa rota usa o prefixo `occurrence.alert.*`. Este repositório não cria esse pacote nem esse log.
- **Rationale**: o #20 declara o pacote do Módulo. O João nomeou o pacote do alerta e o prefixo do log. A ADR nova que formaliza o Módulo dentro do Core não é escrita aqui. (fonte: contexto-llt201.md, §1 e §8; decisoes-joao-llt201.md, registro de que a Ocorrência fica dentro do Core, pendências 14 e 15, formulários 2 e 1; spec.md, FR-008 e Key Entities)

## D12 — Efeito do alerta depois do aceite

- **Decision**: depois do aceite durável (primeira gravação, 201), a homologação confere no Core um documento em `finder_alerts` e o aviso à tela no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, sem evento novo e sem envelope próprio para a Notificação, e sem publicar em `lhamalog:occurrence-events` nem em `lhamalog:notifications:events`. A repetição (200) não avisa de novo. O Finder só faz o POST.
- **Rationale**: a pendência 1 fecha a opção A. A Rodada 3 fecha canal, tipo e corpo. A recomendação D do Especialista (só guardar) não vale. O conflito com a ADR_165 fica só registrado: vale o João. (fonte: decisoes-joao-llt201.md, pendência 1, formulário 1, e Rodada 3; respostas-pendentes-llt201.md, conflito 1; spec.md, FR-010)
- **Alternatives considered**: só gravar, sem aviso (recomendação D, recusada pelo João); evento novo ou envelope para a Notificação (recusado: decisoes-joao-llt201.md, evento para a Notificação e pendência 1, formulário 1).

## P1 — Proxy local: fechado (Caddy e `SSL_CERT_FILE`)

Fechado na rodada 2. Produto: Caddy no compose do Core. O Finder só recebe a URL `https://`. A CA local vale só no teste, via `SSL_CERT_FILE`, com a verificação ligada e sem código novo. Detalhe em D4. Não resta opção em aberto no lugar desta pesquisa. (fonte: decisoes-joao-llt201.md, pendências 13 e 18, formulário 2; respostas-pendentes-llt201.md, pendência 17; origem: ADR_164 + arquitetura-modulos)

## P2 — Retentativa do POST de alerta: fechada (reenvio no reprocessamento)

Fechado na rodada 2. Se o POST falha, o lote permanece pendente e o aviso é reenviado no reprocessamento do lote, com a mesma identidade, apoiado na deduplicação do Core. O código já devolve `ErrAvisoCorePendente` e não marca `recorrente=false` por degradação. Não se acrescenta um número novo de tentativas dentro do POST. A opção escolhida é C. Ela é compatível com 009-18 A, e 009-18 A não é a regra aplicada a este POST. Detalhe em D5. (fonte: respostas-pendentes-llt201.md, pendência 19; origem: DF-FINDER §2 e §3; spec.md, FR-009; `internal/adapter/out/openobserve/recorrencia.go`)

## P3 — Audiências do JWT: fechada (um token com as duas)

Fechado na rodada 2. Um token com as duas audiências (a do Módulo e `lhamalog-core`). As duas variáveis podem carregar esse mesmo JWT. A propriedade `core-audience` vale só na rota de alerta. A recomendação A (dois tokens) não vale. A cadeia `/internal/v1/**` precisa aceitar esse token na rota de alerta por `core-audience` e nas rotas do Módulo pela audiência do Módulo. Detalhe em D2 e D3. (fonte: decisoes-joao-llt201.md, pendências 6 e 20, formulário 1; respostas-pendentes-llt201.md, resumo, fato de código)

## Residuais sem fonte

O nome da coleção (`finder_alerts`) e o literal (`RECURRENT`) fecharam na rodada 3. Não há pergunta nova.

- Sufixos de `event=occurrence.alert.*`: [PENDENTE]. O prefixo está fechado (pendência 14 A) e fica fora do Finder. A string exata de `event=` não é inventada. (fonte: decisoes-joao-llt201.md, pendência 14, formulário 2; analyze-016.md, MEDIUM U2; sufixos: SEM FONTE)

## Conflitos só registrados

1. ADR_165 (ocorrências pela Notificação via stream) × decisão do João "só o aviso do #24". Vale o João. O alerta aceito avisa por WebSocket no sentido da ADR_051. (fonte: decisoes-joao-llt201.md, pendência 1, formulário 1; respostas-pendentes-llt201.md, conflito 1)
2. 009-5 revisto (gravidade no aviso) × spec 015, FR-007 (aviso sem gravidade). Vale o João. A emenda da 015 fica como documento pendente, sem editar a 015. (fonte: respostas-pendentes-llt201.md, pendência 2 e conflito 2)
3. ADR_062 (escala "critico"/"alto", matriz pendente) × gravidade = classe ERROR/WARNING/INFO. Vale o João, pendência 3 D. (fonte: decisoes-joao-llt201.md, pendência 3, formulário 1; respostas-pendentes-llt201.md, conflito 3)

## O que não virou pesquisa de tecnologia nova

Não há biblioteca, framework nem banco novo a escolher no Finder. O Technical Context do plano usa só o que `go.mod`, o README e o código já mostram, mais `SSL_CERT_FILE` da biblioteca padrão. (fonte: constituição, Restrições Técnicas; README.md, Dependências; respostas-pendentes-llt201.md, pendência 18)
