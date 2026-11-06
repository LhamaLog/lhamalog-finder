# Data Model: Consumir fila Redis

Nomes de stream, grupo, campo e do lugar do código concluído não fazem parte deste modelo. Eles continuam entrada externa.

## Item da fila

| Campo | Regra |
|-------|--------|
| Identidade do item | Distinta do código de apontamento. Dois itens podem carregar o mesmo código. A identidade de fio (nome do campo, id da mensagem) não está definida. |
| Código de apontamento | Único conteúdo aceito. Não é o log. Item sem código, ou com conteúdo além do código, está fora do contrato e não ganha tratamento novo. |

## Código de apontamento

| Campo | Regra |
|-------|--------|
| Valor | Posição do log no OpenObserve. O plano não define formato, tamanho nem alfabeto. |
| Conclusões | No máximo uma conclusão para o mesmo valor. Essa conclusão deixa o código disponível uma vez para a etapa futura. |
| Empate | Dois itens com o mesmo valor que terminam juntos produzem uma conclusão e nenhuma reserva restante. Qual item vence não está no modelo. |

## Reserva

Estado do meio do mesmo passo de consumo.

| Transição | Quando | Efeito |
|-----------|--------|--------|
| Livre → em passo | Um worker toma o item e o worker ainda não parou | O item não é entregue a outro worker nem a outra instância |
| Em passo → concluído | O worker e a instância seguem no ar até o fim do passo, e ainda não existe conclusão para aquele código | O item sai da pendência. O código fica disponível uma vez |
| Em passo → livre | O worker para no meio do passo, com a instância no ar, ou a instância para no meio do passo | O item não é concluído. A reserva acaba. Outro worker pode tomá-lo. Sem prazo numérico |
| Em passo → não concluído e não reservado | Já existe conclusão para o mesmo código | O item não entra na contagem de concluídos e não permanece reservado |

No fim normal, reserva e conclusão são o mesmo passo. A reserva sem conclusão só permanece no caso de parada no meio do passo, e essa permanência acaba com a liberação.

## Conclusão do consumo

| Campo | Regra |
|-------|--------|
| Item | Sai da pendência e não volta como item livre |
| Código | Disponível para a etapa futura, uma vez por valor |
| Fora | Leitura do OpenObserve, agregação do log, classe da ocorrência, nome do lugar |

## Relações

- Uma instância do Finder executa um ou mais workers. O número não é `FINDER_WORKERS` por decisão da spec.
- Um worker executa o passo de um item por vez.
- Várias instâncias leem a mesma fila. O mesmo item tem no máximo uma reserva ativa.
- Vários itens podem apontar para o mesmo código. Só um deles é concluído.

## Validação que o domínio aplica

1. Conteúdo do item é só o código de apontamento.
2. Com o worker e a instância no ar, o passo livre termina em conclusão, se o código ainda não foi concluído.
3. Código já concluído não gera segunda conclusão nem reserva restante.
4. Parada no meio do passo não grava conclusão.
5. Falha de leitura e fila sem item livre não gravam conclusão. Nova tentativa e fila de erro não são regras deste modelo.
