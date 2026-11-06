# Contract: Entrega da fila

Contrato de comportamento entre o caso de uso e as portas já existentes. Não nomeia stream, grupo nem campo. O mecanismo Redis está em [research.md](../research.md). O modelo está em [data-model.md](../data-model.md).

Este contrato substitui, para os pontos abaixo, o silêncio da feature 001 sobre item fora do contrato e sobre o sinal entre processos. O restante da feature 001 continua valendo onde este texto não muda o desfecho.

## Porta de entrada

`Consumir item livre`

| Entrada | Saída |
|---------|--------|
| Item aceito, worker e instância no ar, código ainda sem conclusão | Uma conclusão. Item fora da pendência. Código disponível uma vez. |
| O mesmo passo em andamento, instância no ar | Outra instância não recebe esse item como livre. |
| Parada antes de qualquer efeito da conclusão | Item não concluído, código ausente do conjunto, livre de novo na mesma identidade, sem prazo. Enquanto livre, segue à frente dos itens que já estavam atrás. |
| Algum efeito já ocorreu e ainda falta a saída da pendência | A tentativa completa a saída. A contagem permanece um. O código permanece disponível. Não é tratado como repetido que dispensa essa saída. |
| A tentativa de completar falha | Processo continua no ar. Item permanece parcial. Não conta como segunda conclusão. A iteração seguinte tenta de novo, sem prazo e sem repetir até o sucesso na mesma tentativa. |
| Conclusão já completa e o mesmo item aparece de novo | Não volta como livre e não gera segunda conclusão. |
| Segundo item com o mesmo código, inclusive de outra organização | Não concluído e não reservado. Uma conclusão no produto. |
| Item sem código, ou com conteúdo extra | Descarte. Fora da pendência, fora da reserva, sem conclusão, sem código no conjunto. O registro se distingue de item concluído. |
| Item aceito depois de um descarte | Pode ser reservado e concluído. |
| Fila sem item livre, ou falha antes de reservar | Nenhuma conclusão nessa tentativa. Sem prazo e sem fila de erro. |
| Uma instância, fila já ordenada | A sequência de reservas é a ordem já apresentada. |
| Duas ou mais instâncias | Cada uma reserva o próximo item livre que a fila lhe entrega. Podem reservar itens diferentes ao mesmo tempo. A ordem de término pode diferir. Nenhuma escolhe outro item por nível. |

O Finder não escolhe, não renomeia e não cria nível. Os nomes aceitos da ordenação anterior são Very High, High, Medium, Low e Very Low.

## Porta de saída: fila

| Operação | Obrigação |
|----------|-----------|
| Tomar item livre | No máximo um item. Antes de mensagem nova, entrega pendência própria e, em seguida, a pendência mais antiga de consumidor sem conexão viva, na identidade original. Não toma item de consumidor ainda conectado. |
| Concluir | O código entra uma vez no conjunto e o item sai da pendência na mesma operação. Se o código já estava no conjunto, só sai da pendência e a contagem não sobe. |
| Encerrar sem conclusão | Item sai da pendência e da reserva sem nova conclusão. Vale para código repetido e para descarte. |
| Liberar | Reserva em memória acaba sem conclusão e sem prazo. Não publica o item de novo no fim da fila. |

## Porta de saída: códigos e disponibilidade

| Operação | Obrigação |
|----------|-----------|
| Registrar a primeira conclusão | Aceita o código uma vez no produto. A disponibilidade para a etapa futura é esse mesmo registro. Não abre o OpenObserve. |
| Consultar | Diz se aquele código já tem conclusão. |

## Consulta de saúde

Nenhuma consulta nova. `GET /health` permanece o handler já existente: HTTP 200 e corpo `{"status":"UP"}`. O processo em execução atende o endereço já lido de `FINDER_HTTP_ADDR` (padrão `:8085`). O processo parado não atende.

## O que este contrato não cria

Organização no item, campo de nível, fila de erro, prazo numérico, nome novo de stream, grupo ou conjunto, e leitura do log no OpenObserve.
