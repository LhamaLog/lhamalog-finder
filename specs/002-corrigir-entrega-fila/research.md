# Research: Corrigir entrega da fila

## 1. Conclusão e saída da pendência na mesma operação

**Decision**: A conclusão de um item aceito passa a ser uma operação só no Redis: registrar o código no conjunto já configurado e tirar o item da pendência. O código disponível para a etapa futura é esse mesmo registro. Não há terceira escrita. Se o código já está no conjunto e o item ainda está pendente, a operação só tira o item da pendência e não aumenta a contagem de conclusões.

**Rationale**: A spec pede três efeitos observáveis no mesmo desfecho: item fora da pendência, uma conclusão do código no produto e código disponível uma vez. No código atual, `RegistrarPrimeira` e `Disponibilizar` gravam o mesmo conjunto, e o `XACK` vem no meio. A janela que a revisão apontou é essa separação. Juntar o registro e a saída da pendência remove a janela para o processo que cai no meio. O estado parcial que já tiver ficado para trás (código no conjunto, item ainda pendente) é curado pela mesma operação, sem segunda conclusão.

**Alternatives considered**: Manter as três chamadas e só inverter a ordem. Continua existindo intervalo entre elas. Criar uma chave de "processando" e outra de "concluído". A spec cita isso como exemplo de meio, não como requisito, e criaria um lugar que ninguém nomeou.

## 2. Falha ao completar a conclusão

**Decision**: Se a operação de conclusão falha, o processo continua no ar. O item permanece pendente e, se o código já entrou no conjunto, permanece parcial. A iteração acaba. A iteração seguinte faz uma tentativa. Não há espera numérica, nem laço até o sucesso dentro da mesma tentativa, nem `os.Exit` por essa falha. Outro item livre pode ser consumido numa iteração em que a conclusão pendente não seja a única leitura.

**Rationale**: A resposta da sessão escolheu manter o processo e completar depois, sem prazo. Encerrar a execução e repetir até o sucesso na mesma tentativa foram as opções recusadas. O laço atual de `cmd/finder/main.go` já volta à próxima leitura; a mudança é não tratar essa falha como fim do processo.

**Alternatives considered**: Dormir antes de tentar de novo. Isso introduziria o prazo numérico que a spec proíbe. Encerrar o processo. Foi a opção B recusada.

## 3. Item de instância parada volta ao lugar original

**Decision**: A liberação não usa tempo de ociosidade. O adaptador identifica o consumidor pelo nome que já gera (máquina e processo) e publica esse nome na conexão Redis. Antes de ler mensagem nova, procura a pendência mais antiga de um consumidor que não tem conexão viva com esse nome e a reserva para si, no mesmo identificador da fila. O argumento de ociosidade exigido pela operação de reserva alheia fica em zero: ele não é uma regra do Finder. A regra é a conexão ausente. Em seguida, se não houver essa pendência, lê a próxima mensagem nova. A pendência do próprio consumidor vivo é lida antes das duas, para a tentativa seguinte da conclusão parcial.

**Rationale**: A spec exige que outra instância, ou a mesma depois de reiniciar, reserve o item sem prazo numérico, e que ele continue à frente dos itens que já estavam atrás. O identificador da fila não muda quando a reserva muda de consumidor, então o item não vai para o fim. Consumidor ainda conectado não perde o item (FR-006). Processo reiniciado troca o identificador de processo, a conexão antiga some e a pendência antiga pode ser reservada. Duas instâncias podem reservar itens diferentes ao mesmo tempo quando ambas leem mensagens novas; a ordem de término pode diferir.

**Alternatives considered**: `XAUTOCLAIM` com tempo mínimo de ociosidade. O número passaria a definir a liberação, o que a spec proíbe. Confirmar a mensagem e publicá-la de novo. O item iria para o fim, que foi a opção recusada. Apagar o consumidor morto sem reservar a mensagem. A mensagem deixaria o grupo sem o desfecho e sem o lugar original garantido para outra instância.

**Consequência**: Se o sistema operacional ainda mostrar a conexão depois da queda, o item permanece com ela até a conexão sumir. O Finder não configura um tempo para abreviar isso.

## 4. Prioridade

**Decision**: O Finder não lê, não grava e não ordena nível. A ordem da fila é a ordem em que as mensagens já estão na fila. Com uma instância, a sequência de reservas segue essa ordem. Com várias, cada uma reserva a próxima mensagem livre que a fila lhe entrega. Os nomes Very High, High, Medium, Low e Very Low permanecem decisão da ordenação anterior, citada na ADR 169. Nenhum campo novo entra no item.

**Rationale**: A ADR 169 diz que o Finder não escolhe o nível e consome a fila já ordenada. A spec repetiu isso e definiu o comportamento com várias instâncias. Ordenar dentro do Finder contradiria a ADR e exigiria um campo que a spec não nomeia.

**Alternatives considered**: Selecionar pelo nome do nível dentro do adaptador. Foi o pedido que a ADR fechou ao contrário. Guardar o nível no item. Criaria um campo.

## 5. Unicidade do código

**Decision**: O conjunto já configurado continua indexado só pelo código de apontamento. Duas organizações com o mesmo código produzem uma conclusão. O item não ganha organização.

**Rationale**: Resposta Q1 A da spec.

**Alternatives considered**: Incluir organização na chave. Foi a opção B recusada.

## 6. Item fora do contrato

**Decision**: Item sem código, ou com conteúdo além do código, não entra no conjunto. A reserva em memória é apagada e o item sai da pendência. O processo registra um descarte distinguível da linha de item concluído. Não há fila de erro. O item aceito seguinte segue no fluxo. Falha nessa saída da pendência deixa o item pendente e a iteração seguinte tenta o descarte de novo, com o processo no ar.

**Rationale**: A spec pediu descarte observável e limpeza da reserva. Retry não torna o item aceito. Uma fila de erro nomearia um destino que não foi informado.

**Alternatives considered**: Deixar o item pendente até o processo acabar. É o defeito da revisão. Fila de erro. Destino sem nome.

## 7. Consulta de saúde

**Decision**: `cmd/finder/main.go` passa a servir o roteador que já existe em `internal/adapter/in/http`, no endereço que `FINDER_HTTP_ADDR` já define, com padrão `:8085`. O corpo permanece o que `GET /health` já escreve. Nenhuma consulta nova. O processo encerrado não atende.

**Rationale**: O contrato da feature 001 e o Compose já usam essa consulta. O handler existe e o processo não o liga. Restaurar o atendimento preserva o significado.

**Alternatives considered**: Tirar o healthcheck do Compose e apagar a promessa do contrato. Mudaria o significado que a spec mandou manter.

## 8. Dependências e nomes

**Decision**: Nenhuma biblioteca nova. Continua `github.com/redis/go-redis/v9 v9.23.0`, já exigida em `go.mod`. Stream, grupo e conjunto continuam os valores que o processo já lê do ambiente (`FINDER_REDIS_STREAM`, `FINDER_REDIS_GROUP`, `FINDER_REDIS_CODIGOS`). Este plano não escolhe esses nomes nem cria variável nova.

**Rationale**: A spec proíbe nomear stream, grupo e campos. O processo já falha com `ErrNomesAusentes` quando eles não vêm de fora. O Compose de desenvolvimento já os informa. Reabri-los inventaria um contrato que a feature 001 deixou externo.

**Alternatives considered**: Fixar `fila:logs` dentro do código. O protótipo antigo deixou de ser contrato; o valor segue externo.

## 9. Prova

**Decision**: `make test` continua a prova sem Redis, com `-race`. A prova com Redis real roda só quando o quickstart liga o Redis do Compose e define a variável que torna a ausência de Redis uma falha, não um teste ignorado. Essa prova cobre conclusão única, estado parcial, pendência de consumidor sem conexão no lugar original, duas instâncias, código repetido, item inválido e a consulta de saúde do processo no ar.

**Rationale**: A revisão não executou testes e apontou a falta de prova com Redis real. O alvo `make test` usa `gotools` sem o Redis. Misturar as duas provas faria o alvo atual depender da infra. Separar a prova do quickstart mantém o alvo e impede que a prova com Redis passe por ter sido ignorada.

**Alternatives considered**: Exigir Redis em todo `make test`. Quebraria o alvo atual quando a infra estiver desligada. Biblioteca de Redis embutida. Seria dependência nova e não é o Redis do Compose.
