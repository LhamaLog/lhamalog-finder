# Data model: Alerta ao Core e homologação Finder

Entidades desta feature no Finder, alinhadas à rodada 3. Não há coleção, índice nem transição nova neste repositório. Onde o estado vive no Core, o texto só aponta o dono e a regra já fechada. A coleção do alerta é `finder_alerts` e o literal de `change` para "virou recorrente" é `RECURRENT`. Nenhum dos dois é criado aqui.

## Alerta operacional do Finder

Chamada já existente depois que a consulta ao OpenObserve esgota as tentativas. Não é o evento de ocorrência da LLT-201. Se o POST falha, o lote fica pendente e o aviso volta no reprocessamento. (fonte: spec.md, Key Entities; contexto-llt201.md, §3 e R9; respostas-pendentes-llt201.md, pendência 19)

| Campo | Regra |
| --- | --- |
| `tenantId` | Obrigatório no pedido. Vem da referência do lote. Vazio recusa a construção do cliente. (fonte: `internal/adapter/out/core/aviso.go`) |
| `productId` | Obrigatório, mesma origem. |
| `batchId` | Obrigatório. A chamada de aviso recusa lote diferente do que montou o cliente. |
| `motivo` | O Finder envia só `consulta_openobserve_indisponivel`. Outro motivo na chamada é recusado pelo cliente antes do HTTP. (fonte: spec.md, FR-001; `aviso.go`) |

Identidade do recibo de sucesso: os quatro campos iguais ao pedido e `aceito: true`, sem outro campo. Status HTTP 201 na primeira aceitação e 200 na repetição. O corpo não traz `repeticao`. (fonte: `internal/adapter/out/confirmacao/http.go`; spec.md, FR-001; respostas-pendentes-llt201.md, pendência 8; origem: DF-FINDER §3 + controller do #20 + cliente estrito)

Persistência durável dessa identidade, por `(tenantId, productId, batchId, motivo)`, é a coleção Mongo `finder_alerts` no Core, distinta da coleção do recibo e da coleção `Ocorrencia`. Não usar `AlertaFinder`. Quem perde a corrida para o índice único relê o registro e responde 200 com a mesma identidade; fica um documento só. O cliente trata 200 e 201 como sucesso quando a identidade confere. Este modelo não cria a coleção. (fonte: respostas-pendentes-llt201.md, pendências 7, 8 e 9; origem: padrão do #20 + DF-FINDER §3; decisoes-joao-llt201.md, pendência 9, formulário 1, e Rodada 3; spec.md, Edge Cases, FR-001 e FR-010)

Validação no Core, conferida na homologação e não implementada neste repositório: corpo incompleto, campo a mais ou motivo fora da lista → 422 `DADOS_INVALIDOS`, sem gravar; JSON ilegível → 400 `REQUISICAO_INVALIDA`; produto que não pertence ao tenant → 403 `ACESSO_NEGADO`, sem gravar o alerta. (fonte: spec.md, FR-001 e Edge Cases; respostas-pendentes-llt201.md, pendência 4; origem: DF-FINDER §3 + padrão do #20; decisoes-joao-llt201.md, pendência 12, formulário 1; elaboração do 403: respostas-pendentes-llt201.md, pendência 12; origem: spec 013, FR-016)

## Credencial do alerta e da consulta

Um JWT de ambiente. `FINDER_CORE_JWT` e `FINDER_MODULO_OCORRENCIA_JWT` podem carregar o mesmo valor. O Finder não guarda claim em struct: o valor inteiro vai no header `Authorization: Bearer`. O cliente não emite nem inspeciona o token. (fonte: spec.md, Key Entities e FR-002; `confirmacao/http.go`; decisoes-joao-llt201.md, pendência 20, formulário 1)

| Claim exigido na homologação | Valor com fonte |
| --- | --- |
| algoritmo | HS256, o do filtro interno do #20. (fonte: contexto-llt201.md, §1; `InternalServiceSecurityConfiguration.java`, via `git show`) |
| segredo | O de `lhamalog.occurrence.internal-auth.secrets`. Não entra no repositório. (fonte: spec.md, FR-002; `application.properties` do #20) |
| `iss` | O de `lhamalog.occurrence.internal-auth.issuer`. Padrão no #20: `lhamalog-internal`. (fonte: spec.md, User Story 1 cenário 2) |
| `aud` | Contém a audiência do Módulo (padrão `lhamalog-modulo-ocorrencia`) e `lhamalog-core`. (fonte: decisoes-joao-llt201.md, pendência 20, formulário 1) |
| `sub` | Um valor da lista `lhamalog.occurrence.internal-auth.subjects`. Padrão no #20: `lhamalog-finder`. A decisão de 26/11 não cria outro `sub`. (fonte: spec.md, Assumptions; `application.properties` do #20) |
| nomes do segredo | `OCCURRENCE_INTERNAL_JWT_SECRET` (singular, lido por `scripts/InternalJwt.java`) e `OCCURRENCE_INTERNAL_JWT_SECRETS` (plural, propriedade do Core). Não unificar os nomes e não copiar o script. (fonte: spec.md, FR-002; tasks.md, T011; quickstart.md) |

A rota de alerta exige `lhamalog-core`, lida da propriedade `core-audience`, usada só nessa rota. Ilustração de caminho: `lhamalog.occurrence.internal-auth.core-audience`. As rotas do Módulo exigem a audiência do #20. Fato: hoje `InternalServiceSecurityConfiguration` cobre todo `/internal/v1/**` e exige que o `aud` contenha a audiência do Módulo; o `aud` é conferido por "contém"; um token só com `lhamalog-core`, sem o ajuste por rota, recebe 401. `scripts/InternalJwt.java` gera uma audiência por `--aud`. Este modelo não altera o script. (fonte: decisoes-joao-llt201.md, pendências 6 e 20, formulário 1; respostas-pendentes-llt201.md, resumo, fato de código; spec.md, FR-003)

Validação de `exp` e `nbf`: o filtro do #20 declara esses claims. O Finder não os confere. (fonte: comentário de `InternalServiceSecurityConfiguration.java`, via `git show`; spec.md, User Story 1 cenário 4)

## ResumoAnalise e recibo

Os já definidos no DN-009-2. Esta feature não cria campo. Contrato: `specs/009-publicar-evento-ocorrencia/contracts/proposta-dn-009-2.md`. (fonte: spec.md, Key Entities; contexto-llt201.md, §3)

O que a homologação observa no servidor, já entregue pelo #20 e fora do modelo do Finder: recibo idempotente por (tenant, produto, batchId, execId); 201 na primeira entrega; 200 com `repeticao: true` no reenvio do resumo; `409 REVISAO_DESATUALIZADA`; `409 ENTREGA_EM_ANDAMENTO`; FECHADA reabre como ABERTA. O campo `repeticao` deste recibo não é o corpo do alerta. (fonte: contexto-llt201.md, §1; spec.md, FR-006; respostas-pendentes-llt201.md, pendência 8)

## Módulo de Ocorrência

Não é entidade persistida neste repositório. É o pacote `br.com.lhamalog.core.occurrence` dentro do Core, no mesmo processo. Por isso `FINDER_MODULO_OCORRENCIA_URL` e `FINDER_CORE_URL` usam o mesmo host na homologação. A rota de alerta, no Core, fica no pacote `finderalert` (exemplo `br.com.lhamalog.core.finderalert`); este repositório não o cria. (fonte: spec.md, Key Entities e FR-008; contexto-llt201.md, §1 e §8; decisoes-joao-llt201.md, pendência 15, formulário 1)

O cliente atual não recusa hosts diferentes. Essa recusa não entra no modelo. (fonte: spec.md, Edge Cases)

## Ocorrência, gravidade e aviso

Fora do modelo de dados do Finder. Dono: Core, spec 016, reutilizando o aviso do #24. O Finder não tem campo de gravidade nem enum de mudança. (fonte: spec.md, FR-007; decisoes-joao-llt201.md, gravidade e evento para a Notificação; contexto-llt201.md, R2)

O que a homologação confere no Core, sem gravar aqui:

| Regra | Valor |
| --- | --- |
| Precedência, gatilho e literal | Os do FR-007 da spec. O literal de "virou recorrente" é `RECURRENT`, distinto de `CREATED` e de `REOPENED`. Não usar `RECURRING`. (fonte: spec.md, FR-007; decisoes-joao-llt201.md, Rodada 3 e pendência 5, formulário 1; respostas-pendentes-llt201.md, pendência 11) |
| Gravidade | Classe ERROR, WARNING ou INFO copiada para o campo `gravidade`, tipo mantido, no documento e no aviso do #24. (fonte: decisoes-joao-llt201.md, pendência 3, formulário 1; respostas-pendentes-llt201.md, pendência 2; origem: Ocorrência.mmd + 009-5 revisto) |
| Prazo do aviso "virou recorrente" | Até 3 s depois do 201. Junto com o teto de 2 s da rota de alerta (`FINDER_CORE_TIMEOUT`), é o SC-009. Os dois valem em alvos diferentes. O Finder não implementa o teto de 3 s. (fonte: spec.md, SC-009; respostas-pendentes-llt201.md, pendência 16; origem: spec 015, SC-001 + DF-FINDER §3 e §8) |

Forma atual da mensagem no canal, só como fato do código do #24, antes da emenda da spec 015: `tenantId`, `productId`, `occurrenceId`, `occurrenceType`, `change` (`CREATED`, `REOPENED`), `at`. Sem gravidade. A homologação espera a gravidade no aviso e, quando a mudança for "virou recorrente", o literal `RECURRENT` (FR-007). A spec 015, FR-007, ainda não inclui a gravidade. A emenda não é escrita aqui. (fonte: `OccurrenceEvent.java` e `RedisOccurrenceEventPublisher.toMessage`, branch `LLT-137-atualizacao-em-tempo-real`, via `git show`; respostas-pendentes-llt201.md, conflito 2)

O alerta aceito, à parte, na primeira gravação (201), grava o documento em `finder_alerts` e avisa a tela no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, sem evento novo para a Notificação e sem publicar em `lhamalog:occurrence-events` nem em `lhamalog:notifications:events`. A repetição (200) não avisa de novo. Esse efeito não é campo deste modelo. (fonte: decisoes-joao-llt201.md, pendência 1, formulário 1, e Rodada 3; spec.md, FR-010)

## Transições que o Finder já executa

1. Consulta histórica ao OpenObserve, no máximo três tentativas quando a falha é retentável. (fonte: `internal/adapter/out/openobserve/recorrencia.go`)
2. Esgotadas as tentativas, um `POST` de alerta com a identidade acima.
3. Se o POST falha, o resultado é pendência de aviso (`ErrAvisoCorePendente`). O lote permanece pendente. O código não grava `recorrente=false` por essa degradação. O aviso é reenviado no reprocessamento do lote, com a mesma identidade. Não há número novo de tentativas dentro do POST. (fonte: spec.md, FR-009 e Edge Cases; respostas-pendentes-llt201.md, pendência 19; origem: DF-FINDER §2 e §3)
4. Sucesso do POST: 201 na primeira vez ou 200 na repetição, corpo estrito, sem `repeticao`. O Core grava uma vez em `finder_alerts` e, só no 201, avisa no canal `/ws/occurrences` com tipo `FINDER_ALERT`. O Finder não publica esse aviso. (fonte: respostas-pendentes-llt201.md, pendência 8; decisoes-joao-llt201.md, pendência 1, formulário 1, e Rodada 3)

A entrega do `ResumoAnalise` e o efeito no lote (pendente até recibo válido, reconhecimento da fila depois) permanecem os da spec 009. Esta feature não redesenha esse fluxo. (fonte: spec.md, User Story 3 cenário 2, item 4; `proposta-dn-009-2.md`, cenário 4)

## URL

| Campo | Regra já aplicada pelo cliente |
| --- | --- |
| esquema | `https` apenas. (fonte: spec.md, FR-004) |
| host | Obrigatório. Na homologação, o mesmo nas duas variáveis: a URL `https://` do Caddy no compose do Core. (fonte: spec.md, FR-005; decisoes-joao-llt201.md, pendência 13, formulário 2; respostas-pendentes-llt201.md, pendência 17) |
| usuário, consulta, fragmento | Proibidos. (fonte: `confirmacao/http.go` e `modocorrencia/cliente.go`) |
| caminho na base | Vazio ou só `/`. A rota é acrescentada pelo cliente. |
| redirecionamento | Não seguido. A URL configurada é o ponto HTTPS do Caddy. (fonte: spec.md, User Story 2 cenário 4) |
| certificado no teste | CA local só via `SSL_CERT_FILE`, verificação ligada, sem campo novo no cliente. (fonte: decisoes-joao-llt201.md, pendência 18, formulário 2; respostas-pendentes-llt201.md, pendência 18) |

O produto do proxy é Caddy. O arquivo de compose é o do Core. Não são campos deste modelo. (fonte: spec.md, FR-005)
