# Contract: Referência de lote

Contrato entre o caso de uso e as portas. Os campos do envelope estão em [data-model.md](../data-model.md). Este arquivo não repete o schema.

## Porta de entrada

`Interpretar entrada de log.ingested`

| Entrada | Saída |
|---------|--------|
| Envelope válido, chave ausente, efeito da recuperação concluído | Chave `concluido` por 25 horas, entrada confirmada pelo id Redis, um efeito |
| Envelope válido, interpretação terminada, efeito ainda não concluído | Entrada sem confirmação. Chave `em_processamento` se esta entrada criou a reserva |
| Envelope válido, chave `em_processamento`, ociosidade ainda não autoriza reassunção | Entrada permanece pendente. Nenhum efeito novo. Nenhuma confirmação |
| Envelope válido, chave `concluido` dentro de 25 horas | Entrada confirmada. Efeito não se repete. O prazo não é renovado |
| Envelope válido, chave `concluido` com 25 horas completas | A chave não suprime a entrada |
| Envelope inválido | Nenhuma referência válida. Motivo registrado sem senha e sem conteúdo de log. Entrada confirmada. Chave em processamento desse identificador, se existir, é removida. Estado `concluido` não é gravado |
| Redis indisponível no meio do processamento | Nenhuma confirmação de entrada inacabada |
| Quantidade, bloqueio ou ociosidade ausente | O consumo não abre. Nenhuma entrada é lida |

A confirmação usa o id da entrada. `idempotencyId` e `correlationId` não confirmam.

## Porta de saída: fila

| Operação | Obrigação |
|----------|-----------|
| Abrir grupo | Stream `lhamalog:stream:logs`, grupo externo, id inicial `0`. Não usa `$`. Grupo já existente não é apagado |
| Ler | Entrega entradas deste grupo, preservando o id. A quantidade e o bloqueio são os informados |
| Reivindicar pendência | Só depois da ociosidade informada. Não descarta a entrada. A quantidade da passagem é a mesma da leitura |
| Confirmar | Tira a pendência desse id |
| Redis ausente | A pendência já existente permanece |

## Porta de saída: chave do lote

| Operação | Obrigação |
|----------|-----------|
| Reservar | Cria `em_processamento` só se a chave não existe. Sem prazo de 25 horas |
| Concluir | Grava `concluido` e o prazo de 25 horas uma vez, quando o efeito terminou |
| Consultar | Diz se está ausente, em processamento ou concluída |
| Remover reserva | Apaga a chave em processamento. Não grava `concluido` |
| Isolamento | O prefixo é o do Finder. A mesma chave não serve à Varredura |

## Porta de saída: efeito da recuperação

| Operação | Obrigação |
|----------|-----------|
| Executar | Recebe a referência validada e o id da entrada. Devolve efeito concluído ou falha |
| Fora desta porta | Consulta ao OpenObserve, paginação e deduplicação por `lhamalog_record_id` |

Sem resultado de efeito concluído, a entrada válida não é confirmada.

## Registro operacional

Motivo de inválida, id da entrada e `idempotencyId` legível. Sem senha, sem metadata bruta, sem body bruto e sem conteúdo de log.

## O que o Finder não expõe nesta feature

Nenhum endpoint HTTP novo. A consulta de saúde já existente permanece com o mesmo significado.
