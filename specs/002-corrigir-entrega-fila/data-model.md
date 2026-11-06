# Data Model: Corrigir entrega da fila

Os nomes de stream, grupo, conjunto e campos continuam fora deste modelo. O conjunto dos códigos é o que o ambiente já entrega ao processo.

## Item da fila

| Campo | Regra |
|-------|--------|
| Identidade | Distinta do código de apontamento. É a identidade da mensagem na fila e não muda quando outra instância passa a reservar o item. |
| Código | Apontamento do log. Único no produto. Não é o conteúdo do log. Ausente, ou acompanhado de outro conteúdo, torna o item fora do contrato. |
| Conteúdo extra | Verdadeiro quando a mensagem traz mais de um valor, ou um valor que não é o código sozinho. |
| Estado | `livre`, `em_passo`, `concluido` ou `nao_concluido`, como o domínio já nomeia. |

O item não ganha organização nem nível de prioridade.

## Condição parcial

Não é um quinto estado gravado. É a combinação observável: o código já está no conjunto e o item ainda está pendente. A contagem de conclusões desse código permanece um. A tentativa seguinte tira o item da pendência sem novo registro.

## Conclusão

Só está completa quando as três observações são verdadeiras ao mesmo tempo:

1. O item não está pendente.
2. O código tem uma conclusão no produto.
3. O código está disponível uma vez para a etapa futura.

As observações 2 e 3 são o mesmo registro no conjunto. A operação de conclusão grava esse registro e tira o item da pendência juntas. Uma entrega seguinte do mesmo item não aumenta a contagem.

## Código concluído

| Regra | Efeito |
|-------|--------|
| Primeiro item aceito com esse código | Entra no conjunto. Conta uma conclusão. Fica disponível para a etapa futura. |
| Outro item com o mesmo código, inclusive de outra organização | Não entra de novo. Não conta outra conclusão. Não permanece reservado. |
| Parada antes de qualquer efeito | O código não entra no conjunto por causa dessa tentativa. |

## Descarte

Item fora do contrato. Não é gravado no conjunto. Depois do descarte não está pendente nem na reserva em memória do processo. O registro observável é diferente do registro de item concluído. Não há destino adicional.

## Fila já ordenada

A ordem é a ordem das mensagens já presentes. O Finder não reordena. Com uma instância, a sequência de reservas segue essa ordem. Com várias, cada uma reserva a próxima mensagem livre que recebe. Um item que volta a ficar livre conserva a identidade original e, enquanto livre, é entregue antes das mensagens que já estavam atrás dele.

Níveis aceitos, só como ordem já recebida: Very High, High, Medium, Low, Very Low. Medium é o padrão dessa ordenação anterior. O Finder não escolhe nem grava o nível.

## Consulta de saúde

Não é entidade da fila. É a consulta `GET /health` já implementada. Resposta atual: corpo `{"status":"UP"}` e status HTTP 200. Este ajuste não muda esse significado.

## Transições

| De | Quando | Para |
|----|--------|------|
| Livre | Instância reserva o item | Em passo, ainda sem efeito de conclusão |
| Em passo | Parada antes de qualquer efeito | Livre, na mesma identidade, sem código no conjunto |
| Em passo | Operação de conclusão bem-sucedida e código novo | Concluído, fora da pendência, código disponível uma vez |
| Em passo | Código já estava no conjunto | Não concluído, fora da pendência, uma conclusão só |
| Em passo | Registro feito e saída da pendência não feita | Parcial (derivado), contagem permanece um |
| Parcial | Tentativa seguinte consegue sair da pendência | Concluído, sem nova contagem |
| Parcial | Tentativa de completar falha | Parcial, processo continua no ar |
| Fora do contrato | Descarte | Fora da pendência e fora da reserva, sem conclusão |
| Concluído | Nova entrega do mesmo item | Continua fora dos itens livres |

Enquanto o passo segue e a instância está no ar, outra instância não recebe esse item como livre.

## Reserva em memória

`porID` guarda o item reservado por este processo. Todo caminho que tira o item da pendência também o apaga dali: conclusão, código repetido e descarte. Item inválido não permanece no mapa até o processo acabar.
