# Contract: Captura do log

Contrato de comportamento entre o caso de uso da captura e as portas. Não é um schema HTTP do OpenObserve. Organização, stream, índice, campos, URL e credenciais estão fora, por FR-005.

## Porta de entrada

`Capturar logs apontados`

| Entrada | Saída |
|---------|--------|
| Um apontamento utilizável e um registro armazenado, com a instância no ar até o fim | Um log pronto, igual ao registro. Nova leitura pelo mesmo apontamento devolve esse registro |
| O mesmo apontamento indicado duas vezes | Um log pronto |
| Dois apontamentos distintos, os dois com registro, capturados ao mesmo tempo | Dois logs prontos, cada um com o próprio apontamento |
| Dois apontamentos distintos, um com registro e outro ausente, ao mesmo tempo | O existente fica pronto. O ausente conserva a própria falha |
| Fonte mockada sem apontamento | Zero logs prontos e zero falhas |
| Apontamento vazio | Falha inutilizável. Log não pronto |
| Apontamento utilizável sem registro | Falha ausente. Log não pronto |
| Armazenamento não alcançado | Falha inacessível, distinta da ausência. Log não pronto |
| Instância para no meio | Log não pronto. Sem retomada e sem prazo |
| Instância para depois do sucesso | Esta etapa não exige que o log pronto continue |

A fila não é entrada desta porta.

## Porta de saída: fonte mockada

| Operação | Obrigação |
|----------|-----------|
| Obter apontamentos | Entrega os valores desta execução, ou nenhum. Não lê a fila Redis |

O mecanismo concreto do mock não nomeia arquivo, variável ou tabela. No aceite, um dos valores seleciona um registro de `otel-logs.zip` já armazenado no OpenObserve. Esse arquivo ainda não foi entregue.

## Porta de saída: leitura do log

| Operação | Obrigação |
|----------|-----------|
| Ler um apontamento | Devolve um único registro, ou informa ausência, ou informa que o armazenamento não pode ser alcançado |
| Registro obtido | O conteúdo segue inteiro. A porta não oculta campo e não classifica ocorrência |

O adaptador que fala com o OpenObserve não é escrito enquanto organização, stream, índice, campos e credenciais não forem informados. O arquivo `otel-logs.zip` não é esta porta.

## Porta de saída: log pronto

| Operação | Obrigação |
|----------|-----------|
| Guardar | Associa o registro inteiro ao apontamento. Uma segunda guarda do mesmo apontamento não cria outro log |
| Ler | Com a instância no ar, devolve esse único registro |
| Falha ou interrupção | Não guarda log para esse apontamento |

O lugar não tem nome. A porta não promete leitura depois que a instância para.

## O que o Finder não expõe nesta feature

Nenhum endpoint HTTP novo. `GET /health` permanece como está. O contrato da fila da feature 001 permanece como está.
