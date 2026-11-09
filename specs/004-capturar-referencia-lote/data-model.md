# Data Model: Capturar referência de lote

Modelo da spec. Nomes de grupo, consumidor e prefixo, e os três números de leitura e pendência, não são campos com valor próprio. São configuração externa.

## Entrada Redis

Publicação em `lhamalog:stream:logs`.

| Campo | Regra |
|-------|--------|
| Id da entrada | Identidade da mensagem no stream. Confirma a entrada. Não deduplica o lote. |
| Valores | Todos texto. Os nomes do envelope são `idempotencyId`, `event`, `origin`, `schemaVersion`, `publishedAt`, `metadata`, `body`. |

Relação: uma entrada produz no máximo uma referência validada. Várias entradas podem apontar o mesmo `idempotencyId`.

## Envelope

| Campo | Regra |
|-------|--------|
| `idempotencyId` | Igual a `openobserve.batchId`. Chave do lote. |
| `event` | `log.ingested` |
| `origin` | `collector` |
| `schemaVersion` | Texto `"1"`. Outro valor, inclusive schema 2, não é este contrato. |
| `publishedAt` | Instante ISO em UTC. |
| `metadata` | Texto JSON. |
| `body` | Texto JSON. Referência e agregados. Não é o conteúdo dos logs. |

## Metadata

| Campo | Regra |
|-------|--------|
| `tenantId` | Identidade do tenant. Vem do token validado pelo Collector. É a identidade usada na referência entregue. |
| `productId` | Identidade do produto, também do token. Igual ao `productId` do body. |
| `correlationId` | Correlação da exportação. Não deduplica. |
| `traceId` | Opcional. A ausência não invalida. Não deduplica. |

## Referência de lote

Conteúdo interpretado do body.

| Campo | Regra |
|-------|--------|
| `receivedAt` | Preservado quando legível. A ausência não invalida. |
| `signal` | `LOGS` para a referência seguir. |
| `productId` | Igual ao da metadata. |
| `openobserve.organization` | Obrigatório para a referência válida. |
| `openobserve.stream` | Obrigatório. Não é o nome do stream Redis. |
| `openobserve.batchId` | Obrigatório. Igual ao `idempotencyId`. |
| `openobserve.recordId` | Preenchido quando `recordCount` é 1. Nulo quando `recordCount` é maior que 1. |
| `openobserve.minTimestampMicros` | Obrigatório. Inteiro. Menor ou igual ao máximo. |
| `openobserve.maxTimestampMicros` | Obrigatório. Inteiro. |
| `recordCount` | Inteiro positivo. Conta os registros aceitos pelo Collector, não o lote bruto. |
| `bytes` | Preservado quando legível. A ausência não invalida. |
| `serviceNames`, `priority`, `criticality`, `maxSeverityNumber`, `severityCounts` | Preservados quando legíveis. A ausência não invalida. A enumeração não é revalidada. |

## Referência validada

Recorte entregue à porta de efeito. A entrega não confirma a entrada.

Obrigatórios: `tenantId`, `productId`, organização, stream do OpenObserve, `batchId`, faixa em microssegundos, `recordCount`, `idempotencyId`, correlação presente e o id da entrada Redis. `recordId` entra no lote de um registro. `traceId`, `receivedAt`, `bytes` e os agregados acompanham quando presentes.

## Chave do lote

Uma chave por `idempotencyId`, no prefixo externo do Finder. Não é compartilhada com a Varredura.

| Estado | Significado | Prazo | Transição |
|--------|-------------|-------|-----------|
| ausente | Nenhum processamento deste identificador | — | A primeira entrada válida autorizada cria `em_processamento` |
| `em_processamento` | Reserva. Não é conclusão | Nenhum | Efeito concluído passa para `concluido`. Entrada inválida remove a chave. Queda mantém o estado |
| `concluido` | Efeito da recuperação terminou | 25 horas a partir da passagem para este estado | Nova entrada com o mesmo identificador é confirmada sem repetir o efeito. Depois de 25 horas a chave não suprime mais |

Enquanto está `em_processamento` e a ociosidade configurada não autoriza reassunção, outra entrada válida do mesmo identificador permanece pendente.

## Grupo do Finder

Grupo de consumo no stream de logs, com nome externo, distinto do grupo da Varredura. Nasce no id `0` quando ainda não existe. Instâncias que usam esse grupo dividem as entradas. Cada entrada tem no máximo uma entrega ativa nesse grupo.

## Configuração exigida para abrir o consumo

Sem valor padrão no código:

| Dado | Efeito se ausente ou inválido |
|------|-------------------------------|
| Stream diferente de `lhamalog:stream:logs`, ou vazio | O consumo não abre |
| Grupo | O consumo não abre |
| Consumidor | O consumo não abre |
| Prefixo da chave | O consumo não abre |
| Quantidade por leitura | O consumo não abre |
| Tempo de bloqueio da leitura nova | O consumo não abre |
| Ociosidade para reassumir pendência | O consumo não abre |

Senha do Redis, quando existir, vem do ambiente e não é registrada. Vazia significa Redis sem senha, como o Compose atual.

## O que este modelo não guarda

Conteúdo do log, consulta ao OpenObserve, ocorrência, checked, trace e fila de erro.
