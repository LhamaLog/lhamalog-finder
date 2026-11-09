# Contract: Classificar apontamento

Contrato entre o caso de uso e as portas. Não é um endpoint novo do Finder. A fila não é lida e a entrada não é confirmada.

## Porta de entrada

`Classificar apontamento`

| Entrada | Saída |
|---------|--------|
| Apontamento com logs presentes, tenant, produto e gravação disponível | Uma ocorrência. A classe é a mais grave entre os logs presentes. A consulta por essa identidade devolve esse documento, sem texto de log. |
| Três logs, classes informação, alerta e erro | Uma ocorrência de classe erro |
| Dois logs, informação e alerta | Uma ocorrência de classe alerta |
| Três esperados e dois presentes, informação e erro | Uma ocorrência de classe erro. O ausente não entra. |
| Nenhum log presente | Sem ocorrência. Motivo: log ausente. |
| Armazenamento inacessível | Sem ocorrência. Motivo distinto do log ausente. A lista que vier junto não é gravada. |
| Gravação não concluída | Sem ocorrência. Motivo distinto dos dois anteriores. |
| URI, banco ou coleção ausente | Sem ocorrência. Motivo: configuração ausente. Nenhuma URI é criada. |
| Tenant ou produto vazio | Sem ocorrência. Motivo: identidade incompleta. |
| Campo `error` e catálogo que diria informação | Classe do log: erro. O catálogo não entra. |
| Campo `informação` e a palavra inteira `error` na mensagem | Classe do log: erro. O catálogo não entra. |
| Campo `informação` e `errorless` na mensagem | Classe do log: informação. |
| Campo ausente e `traceback` na mensagem, sem catálogo | Classe do log: informação. |
| Sem nível no campo e na mensagem, primeira entrada válida `alerta` | Classe do log: alerta. |
| Sem identificação de catálogo | A ocorrência informa que só as regras de nível desta feature foram usadas. |
| Segunda chamada da mesma identidade | Outra ocorrência. A primeira permanece. |

Nenhuma linha envia o texto do log a um serviço externo de linguagem. Nenhuma linha publica evento, gera sugestão ou dispara notificação.

## Porta de saída: ocorrência

| Operação | Obrigação |
|----------|-----------|
| Inserir | Um documento com classe, identificação das regras, se o catálogo entrou, identidade do apontamento, tenant e produto. Sem texto de log. |
| Consultar | Pela identidade do apontamento, devolve os documentos inseridos para ela. |
| Repetir a identidade | Insere de novo. Não substitui. |
| Falha de gravação | Não entrega documento como ocorrência pronta. O motivo é falha de gravação. |
| Fora desta porta | Fila, confirmação da entrada, OpenObserve, sugestão, evento e notificação |

O adaptador não registra a mensagem nem o campo de nível.

## Registro operacional

Identidade do apontamento, classe e motivo quando não há ocorrência. Sem mensagem e sem texto de log.

## O que o Finder não expõe nesta feature

Nenhum endpoint HTTP novo. `GET /health` permanece como está. O consumo da fila não chama esta porta.
