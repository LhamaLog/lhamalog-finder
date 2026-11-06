# Contract: Consumo da fila

Contrato de comportamento entre o caso de uso do Finder e as portas de saída. Não é um schema de mensagem Redis. Nomes de stream, grupo e campos estão fora, por FR-007.

## Porta de entrada

`Consumir item livre`

| Entrada | Saída |
|---------|--------|
| Um item livre cujo conteúdo é só o código de apontamento, com worker e instância no ar | Conclusão no mesmo passo, item fora da pendência, código disponível uma vez |
| O mesmo passo ainda em andamento, worker no ar | O item não é entregue de novo |
| Worker para no meio do passo, instância no ar | Item não concluído e livre de novo, sem prazo |
| Instância para no meio do passo | Item não concluído e livre para outra instância, sem prazo. O sinal entre processos que cumpre isso sem prazo não está neste contrato |
| Segundo item com código já concluído | Não concluído e não reservado |
| Dois itens de mesmo código terminam juntos | Uma conclusão e nenhuma reserva restante. O vencedor não está no contrato |
| Fila sem item livre, ou falha de leitura | Nenhuma conclusão nesta tentativa |
| Item sem código, ou com conteúdo extra | Fora do contrato. Sem tratamento novo |

## Porta de saída: fila

Operações vistas pelo caso de uso, sem comando de banco:

| Operação | Obrigação |
|----------|-----------|
| Tomar item livre | Entrega no máximo um item ainda não reservado por outra instância |
| Concluir | O item deixa a pendência e não volta como livre |
| Liberar | A reserva acaba e o item volta a ficar livre, sem conclusão e sem prazo numérico |

O adaptador que fala com o Redis não é escrito enquanto stream, grupo e campos não forem informados. O protótipo que apaga `fila:logs` não satisfaz esta porta.

## Porta de saída: códigos concluídos

| Operação | Obrigação |
|----------|-----------|
| Registrar primeira conclusão de um código | Aceita uma vez e rejeita as seguintes |
| Consultar | Diz se aquele código já tem conclusão |

O registro é compartilhado entre instâncias. O lugar físico não tem nome neste contrato.

## Porta de saída: código disponível

| Operação | Obrigação |
|----------|-----------|
| Disponibilizar | Guarda o código para a etapa futura, uma vez, sem abrir o OpenObserve |

O lugar também não tem nome. Esta etapa não lê o log.

## O que o Finder não expõe nesta feature

Nenhum endpoint HTTP novo. `GET /health` permanece como está.
