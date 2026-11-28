# Quickstart: homologação Finder com o Core

Roteiro de validação da spec 011, já gerado nesta pasta e alinhado à rodada 3 (27/11) em cima do roteiro da rodada 2. A spec aponta para este arquivo; ele não é arquivo futuro. Não é implementação e não foi executado na geração deste plano. Comandos de teste abaixo são os que o README já documenta para a suíte existente; a homologação ponta a ponta depende do Core com a pilha mergeada e da rota da spec 016. (fonte: spec.md, Roteiro de homologação e Evidência e consulta; analyze-011.md, MEDIUM I2; constituição, princípio V)

Contratos: [alerta-finder-cliente.md](contracts/alerta-finder-cliente.md) e [dn-009-2-homologacao.md](contracts/dn-009-2-homologacao.md). Entidades: [data-model.md](data-model.md). Decisões: [research.md](research.md) D2, D4, D5, D7, D9, D12 e P1–P3 fechadas. A coleção conferida é `finder_alerts` e o literal conferido é `RECURRENT`. Residual que este roteiro não preenche: sufixos de `event=occurrence.alert.*` ([research.md](research.md), Residuais sem fonte).

## Pré-requisitos

1. Core com #19, #20, #21 e #24 mergeados e com `POST /internal/v1/alertas/finder` da spec 016 disponível. A rota lê `lhamalog-core` da propriedade `core-audience`, usada só nela. As rotas do Módulo seguem com a audiência do #20. A cadeia `/internal/v1/**` aceita um token com as duas audiências: na rota de alerta por `core-audience`, nas rotas do Módulo pela audiência do Módulo. Fato: sem esse ajuste, `InternalServiceSecurityConfiguration` exige que o `aud` contenha a audiência do Módulo em todo `/internal/v1/**`, e um token só com `lhamalog-core` recebe 401. (fonte: spec.md, Roteiro item 1; decisoes-joao-llt201.md, base da branch, rota POST, pendências 6 e 20, formulário 1; respostas-pendentes-llt201.md, resumo, fato de código)
2. Caddy no compose do Core, na frente do Core local, que segue em HTTP na porta 8080. O Finder só recebe a URL `https://` que a spec 016 registrar. Este repositório não é editado para criar o proxy. (fonte: spec.md, FR-005; decisoes-joao-llt201.md, pendência 13, formulário 2; respostas-pendentes-llt201.md, pendência 17; origem: ADR_164 + arquitetura-modulos)
3. No teste, a CA local desse Caddy só por `SSL_CERT_FILE`, com a verificação HTTPS ligada e sem código novo no cliente. O caminho do arquivo é o que o Caddy da spec 016 expõe. (fonte: decisoes-joao-llt201.md, pendência 18, formulário 2; respostas-pendentes-llt201.md, pendência 18)
4. Go 1.27.4 se for rodar a suíte já existente do Finder. (fonte: README.md, Pré-requisitos; `go.mod`)
5. Segredo só no ambiente, com pelo menos 32 bytes, o mesmo valor que o Core lê em `OCCURRENCE_INTERNAL_JWT_SECRETS`. O script do Core, se for usado como referência, lê o nome `OCCURRENCE_INTERNAL_JWT_SECRET`. Os dois nomes são os dos arquivos; este roteiro não os funde num nome novo. O script gera uma audiência por `--aud`. Esta feature não altera o script. Uma invocação com um único `--aud` não é o token da pendência 20. (fonte: `application.properties` e `scripts/InternalJwt.java` do #20, via `git show` e respostas-pendentes-llt201.md, resumo, fato de código; spec.md, FR-002 e FR-003)

## Configuração

As duas URLs são `https://` desse Caddy, o mesmo host, sem caminho na base. A CA local do teste está em `SSL_CERT_FILE`. (fonte: spec.md, Roteiro itens 2 a 4; FR-004 e FR-005)

| Variável | Valor na homologação |
| --- | --- |
| `FINDER_MODULO_OCORRENCIA_URL` | `https://` do Caddy, mesmo host de `FINDER_CORE_URL` |
| `FINDER_CORE_URL` | a mesma origem |
| `FINDER_MODULO_OCORRENCIA_JWT` | pode ser o mesmo JWT de `FINDER_CORE_JWT` |
| `FINDER_CORE_JWT` | JWT HS256; segredo e emissor do #20; `aud` contém a audiência do Módulo e `lhamalog-core`; `sub` na lista do #20 (padrão `lhamalog-finder`). O script lê `OCCURRENCE_INTERNAL_JWT_SECRET` (singular); a propriedade do Core é `OCCURRENCE_INTERNAL_JWT_SECRETS` (plural). Não unificar os nomes. |
| `FINDER_CORE_TIMEOUT` | permanece; padrão `2s` se ausente. A rota de alerta responde dentro desse prazo. Esse teto é uma das duas medidas do SC-009. |
| `SSL_CERT_FILE` | só no teste: CA local do Caddy. Verificação ligada. Sem código novo. |

(fonte: spec.md, Roteiro itens 3 a 6; research.md D2, D3, D4 e D9; decisoes-joao-llt201.md, pendências 18 e 20, formulários 2 e 1; respostas-pendentes-llt201.md, pendência 16)

Emissão fora do Finder. Não copiar `scripts/InternalJwt.java` para este repositório e não gravar o token no Git. O cliente não inspeciona `iss`, `aud` nem `sub`. (fonte: spec.md, FR-002 e User Story 1 cenário 4)

## O que conferir antes do fluxo

1. Base `http://` em `FINDER_CORE_URL` ou em `FINDER_MODULO_OCORRENCIA_URL` faz o processo recusar a configuração (`pipeline recusado` / erro de URL). Nenhuma chamada sai. Resultado esperado do SC-002. (fonte: spec.md, SC-002 e User Story 2; `internal/bootstrap/pipeline.go`; README.md, configuração)
2. Com as URLs `https://`, `SSL_CERT_FILE` no teste e o JWT presente nas duas variáveis, a validação de prontidão constrói o cliente de alerta e o cliente do Módulo sem chamada externa. (fonte: `internal/bootstrap/pipeline.go`, comentário “Sem chamadas externas”; spec.md, User Story 2 cenário 5)

## Suíte já existente

Na raiz do Finder, sem Core, a regressão do cliente que esta feature não reescreve:

```bash
go test ./...
go vet ./...
```

`make test` acrescenta `-race`, como o README. Isso não substitui a homologação com servidor real. Não há variável de integração Finder↔Core documentada no README para ligar os cinco cenários. (fonte: README.md, seção Testes; spec.md, User Story 3 — o aceite é com o servidor real)

## Cenários ponta a ponta

Com Finder e Core no ar, Caddy com URL `https://` conhecida, `SSL_CERT_FILE` no teste e o JWT das duas audiências nas duas variáveis.

### Alerta

1. Provocar o caminho em que a consulta ao OpenObserve esgota as três tentativas. (fonte: `recorrencia.go`)
2. Esperado: um `POST /internal/v1/alertas/finder` com o corpo de [alerta-finder-cliente.md](contracts/alerta-finder-cliente.md) e o Bearer de `FINDER_CORE_JWT`.
3. Esperado no Finder, na primeira aceitação: 201 com `aceito: true` e a mesma identidade, sem campo `repeticao`. Na repetição da mesma identidade: 200, o mesmo corpo, sem segundo alerta. A rota responde dentro de 2 s (`FINDER_CORE_TIMEOUT`), um dos tetos do SC-009. Isso cobre o SC-001 quando o token tem o segredo, o emissor e as duas audiências. (fonte: spec.md, SC-001 e SC-009; respostas-pendentes-llt201.md, pendências 8 e 16)
4. Esperado no Core, na primeira gravação (201): um documento na coleção `finder_alerts` e o aviso à tela no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, sem evento novo, sem envelope para a Notificação e sem publicar em `lhamalog:occurrence-events` nem em `lhamalog:notifications:events`. A repetição (200) não avisa de novo. Quem perde a corrida de duas primeiras gravações da mesma chave relê e responde 200 com a mesma identidade; o cliente trata 200 e 201 como sucesso. Isso cobre o SC-008. (fonte: spec.md, SC-008, FR-001 e FR-010; decisoes-joao-llt201.md, pendência 1 e pendência 9, formulário 1, e Rodada 3)
5. Um token só com audiência `lhamalog-modulo-ocorrencia` não autentica a rota de alerta. Um token só com `lhamalog-core` não é a credencial das rotas do Módulo; sem o ajuste por rota, a cadeia atual responde 401. (fonte: spec.md, SC-001 e Edge Cases; respostas-pendentes-llt201.md, resumo, fato de código)
6. Se o POST falha, o lote permanece pendente e o aviso é reenviado no reprocessamento desse lote, com a mesma identidade. O POST falho não confirma o lote. Não há número novo de tentativas dentro do POST. Isso cobre o SC-007. Status diferente de 200/201 conta como essa falha, inclusive 422 `DADOS_INVALIDOS` (corpo incompleto, campo a mais ou motivo fora da lista), 400 `REQUISICAO_INVALIDA` (JSON ilegível) e 403 `ACESSO_NEGADO` (produto que não pertence ao tenant, sem gravar o alerta). (fonte: spec.md, SC-007, FR-001 e FR-009; respostas-pendentes-llt201.md, pendências 4, 12 e 19; decisoes-joao-llt201.md, pendência 12, formulário 1)

### Entrega no mesmo host

1. Uma entrega de `ResumoAnalise` pela URL `https://` do mesmo host do alerta, com a CA local só em `SSL_CERT_FILE` e a verificação ligada. Esperado: a rota do #20 responde e o Finder aceita o recibo do DN-009-2. Isso, junto com o alerta, cobre o SC-003. (fonte: spec.md, SC-003 e User Story 3 cenário 4; decisoes-joao-llt201.md, pendências 13 e 18, formulário 2)
2. Percorrer os cinco cenários listados em [dn-009-2-homologacao.md](contracts/dn-009-2-homologacao.md), com o cliente já existente. Esperado: cada cenário passa uma vez. Isso cobre o SC-004. (fonte: spec.md, SC-004 e FR-006)

### Conferência que não é do Finder

Para uma entrega que o Core trata como ocorrência nova, como reaberta ou como "virou recorrente":

- O Finder apenas entregou o resumo.
- O Finder não publica em `lhamalog:occurrence-events` e não calcula gravidade.
- O Core segue o FR-007: precedência, gatilho (`before.recurring` não é TRUE e `after.recurring` é TRUE; ausente ou `null` conta como `false`) e literal `RECURRENT`.
- O aviso "virou recorrente" sai em até 3 s depois do 201. Esse teto, com o de 2 s da rota de alerta medido no cenário do alerta, é o SC-009. Os dois valem em alvos diferentes. O Finder não implementa o teto de 3 s.
- A gravidade no documento e no aviso do #24 é a classe no campo `gravidade`, com o tipo mantido.

Resultado esperado do SC-005: em todas essas homologações o processo Finder não é a origem do evento nem da gravidade, e o Core segue o FR-007. (fonte: spec.md, SC-005, SC-009 e FR-007; decisoes-joao-llt201.md, pendências 3 e 5, formulário 1, e Rodada 3; respostas-pendentes-llt201.md, pendências 2, 11 e 16)

## Resultado que este arquivo não produz

SC-006 é documental: a spec 009 e a constituição permanecem sem edição. A divergência está em `spec.md`, seção Divergências, e a emenda está em `plan.md`, seção Pendências — documentos necessários. Os conflitos ADR_165, spec 015 FR-007 e ADR_062 ficam só registrados. (fonte: spec.md, SC-006 e Divergências; plan.md, Conflitos só registrados)
