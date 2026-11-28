# Contrato consumido: alerta do Finder

Interface que o Finder já chama. Este arquivo não é o servidor. A implementação da rota é a spec 016 do lhamalog-core, no pacote `finderalert` (exemplo `br.com.lhamalog.core.finderalert`). A segurança continua na cadeia `/internal/v1/**` do #20. O Finder não ganha segunda implementação. O log de negócio dessa rota, no Core, usa o prefixo `occurrence.alert.*`. (fonte: spec.md, FR-001 e Fora de escopo; decisoes-joao-llt201.md, rota POST, pendências 14 e 15, formulários 2 e 1)

## Pedido

- **Método e caminho**: `POST /internal/v1/alertas/finder`
- **URL base**: `FINDER_CORE_URL`, só `https://`, sem usuário, consulta, fragmento ou caminho além de `/`
- **Header**: `Content-Type: application/json` e `Authorization: Bearer` com o valor de `FINDER_CORE_JWT`
- **Prazo**: `FINDER_CORE_TIMEOUT`, padrão `2s`. A rota do Core responde dentro desse prazo. Esse teto, com o de 3 s do aviso "virou recorrente" depois do 201, é o SC-009. Os dois valem em alvos diferentes. O Finder não implementa o teto de 3 s. (fonte: spec.md, SC-009 e Assumptions; respostas-pendentes-llt201.md, pendência 16; origem: spec 015, SC-001 + DF-FINDER §3 e §8)
- **Corpo**:

```json
{
  "tenantId": "<tenant>",
  "productId": "<produto>",
  "batchId": "<lote>",
  "motivo": "consulta_openobserve_indisponivel"
}
```

O único motivo que este cliente envia é `consulta_openobserve_indisponivel`. (fonte: spec.md, FR-001; contexto-llt201.md, §3 e R9; `internal/adapter/out/core/aviso.go`)

## Sucesso que o cliente aceita

- **Primeira aceitação**: 201
- **Repetição da mesma identidade**: 200
- **Corpo**: exatamente `aceito` verdadeiro e os mesmos `tenantId`, `productId`, `batchId` e `motivo` do pedido. Sem campo `repeticao` e sem outro campo. O decodificador é estrito nos nomes pedidos. Campo a mais, nome diferente ou identidade diferente não é aceite. 202 e corpo vazio não são aceite.

```json
{
  "aceito": true,
  "tenantId": "<o mesmo>",
  "productId": "<o mesmo>",
  "batchId": "<o mesmo>",
  "motivo": "consulta_openobserve_indisponivel"
}
```

(fonte: spec.md, FR-001 e Edge Cases; `internal/adapter/out/confirmacao/http.go`; `internal/adapter/out/core/aviso_test.go`; respostas-pendentes-llt201.md, pendência 8; origem: DF-FINDER §3 + `InternalOccurrenceController.deliver` do #20 + cliente estrito)

Duas primeiras requisições da mesma chave: quem grava responde 201; quem perde a corrida relê e responde 200 com a mesma identidade. O cliente trata os dois status como sucesso quando a identidade confere. (fonte: decisoes-joao-llt201.md, pendência 9, formulário 1; respostas-pendentes-llt201.md, pendências 8 e 9)

## Credencial

O token é emitido fora do Finder. Segredo e emissor iguais aos das rotas internas do #20. O `aud` contém a audiência do Módulo (padrão `lhamalog-modulo-ocorrencia`) e `lhamalog-core`. `FINDER_MODULO_OCORRENCIA_JWT` pode carregar o mesmo JWT. `sub` na lista já configurada no #20 (padrão `lhamalog-finder`). Algoritmo HS256, o do filtro existente. O script `scripts/InternalJwt.java` lê `OCCURRENCE_INTERNAL_JWT_SECRET` (singular); a propriedade do Core é `OCCURRENCE_INTERNAL_JWT_SECRETS` (plural). Este contrato não unifica esses nomes e não copia o script. Segredo não aparece neste contrato nem no repositório. O cliente não emite nem inspeciona claims. (fonte: spec.md, FR-002; tasks.md, T011; quickstart.md; decisoes-joao-llt201.md, pendência 20, formulário 1; contexto-llt201.md, §1; `application.properties` do #20, via `git show`)

A rota exige `lhamalog-core`, lida da propriedade `core-audience`, usada só nesta rota. Ilustração de caminho: `lhamalog.occurrence.internal-auth.core-audience`. As rotas do Módulo continuam exigindo a audiência do Módulo. (fonte: decisoes-joao-llt201.md, pendência 6, formulário 1; spec.md, FR-003)

Fato de código: no #20, `InternalServiceSecurityConfiguration` cobre todo `/internal/v1/**` e exige que o `aud` contenha `lhamalog.occurrence.internal-auth.audience` (padrão `lhamalog-modulo-ocorrencia`); o `aud` é conferido por "contém"; um token só com `lhamalog-core`, sem o ajuste por rota, recebe 401. `scripts/InternalJwt.java` gera um token com uma única audiência (`--aud`). Este contrato não altera o script. Um token cuja audiência é só `lhamalog-modulo-ocorrencia` não autentica esta rota. (fonte: respostas-pendentes-llt201.md, resumo, fato de código; spec.md, SC-001)

## Recusas que o cliente trata como falha

Qualquer status diferente de 200 e 201 é falha para o cliente (`confirmacao/http.go`). O lote permanece pendente e o aviso é reenviado no reprocessamento. Códigos que o Core usa, e que este contrato não inventa:

| Situação | Resposta do Core |
| --- | --- |
| Corpo incompleto, campo a mais ou motivo fora da lista | 422 `{"codigo":"DADOS_INVALIDOS"}` |
| JSON ilegível | 400 `{"codigo":"REQUISICAO_INVALIDA"}` |
| Produto que não pertence ao tenant | 403 `{"codigo":"ACESSO_NEGADO"}`, sem gravar alerta |

(fonte: respostas-pendentes-llt201.md, pendências 4 e 19; origem: DF-FINDER §3 + padrão do #20; decisoes-joao-llt201.md, pendência 12, formulário 1; elaboração do 403: respostas-pendentes-llt201.md, pendência 12, spec 013, FR-016)

## O que este contrato não define no Finder

- Deduplicação durável pela chave `(tenantId, productId, batchId, motivo)`, na coleção Mongo `finder_alerts`, distinta da do recibo. Não usar `AlertaFinder`. (fonte: respostas-pendentes-llt201.md, pendência 7; origem: padrão do #20 + DF-FINDER §3; nome: decisoes-joao-llt201.md, Rodada 3)
- Efeito depois do aceite durável (201): o Core grava o alerta em `finder_alerts` e avisa a tela no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, sem evento novo, sem envelope próprio para a Notificação e sem publicar em `lhamalog:occurrence-events` nem em `lhamalog:notifications:events`. A repetição (200) não avisa de novo. O Finder não publica esse efeito. (fonte: decisoes-joao-llt201.md, pendência 1, formulário 1, e Rodada 3; spec.md, FR-010)
- Prefixo de log `occurrence.alert.*` e pacote `br.com.lhamalog.core.finderalert`: são do Core. Não adotar `finder.alert.*`. Os sufixos de `event=` permanecem [PENDENTE]; a string exata não é inventada. (fonte: decisoes-joao-llt201.md, pendência 14, formulário 2, e pendência 15, formulário 1; analyze-016.md, MEDIUM U2; sufixos: SEM FONTE)

## Transporte

Sem redirecionamento. Verificação TLS padrão do processo quando o transporte não é substituído em teste. No teste local, a base é a URL `https://` do Caddy no compose do Core, e a CA local entra só por `SSL_CERT_FILE`, sem código novo e sem desligar a verificação. (fonte: spec.md, User Story 2 e Assumptions; README.md; `confirmacao/http.go`; decisoes-joao-llt201.md, pendências 13 e 18, formulário 2; respostas-pendentes-llt201.md, pendências 17 e 18)
