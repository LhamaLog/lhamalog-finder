# Feature Specification: Consumir fila Redis

**Feature Branch**: `001-consumir-fila-redis`

**Created**: 2026-10-29

**Status**: Draft

**Input**: User description: "Será desenvolvido um mecanismo para coletar a fila do Redis ainda não implementada, no qual será consumida por workers e colocada em estado de processamento. A fila contém o local de apontamento do OpenObserve. O trabalho desta etapa é somente o consumo do Redis. O consumo deve funcionar com diversas instâncias do Finder, sem racing condition. O objetivo é consumo da fila stream Redis com alta performance e alta escalabilidade. Restrições do pedido: seguir a estrutura atual do projeto, não alterar locais não relacionados, não inventar informações, arquitetura clean, padrões do projeto, Golang, instalar libs somente relacionadas à tarefa. Incertezas devem ser consultadas antes de criar dados."

## Clarifications

### Session 2026-10-29

- Q: Quando o consumo termina e o item sai da pendência, o código de apontamento continua disponível para uma etapa futura? → A: Sim. O item sai da pendência e o código continua disponível para a etapa futura. Esta etapa não abre o OpenObserve e não define o nome do lugar onde o código fica.
- Q: Se a instância para depois de reservar o item e antes de confirmar a conclusão, outra instância pode reservar esse mesmo item? → A: Sim. A reserva incompleta acaba quando a instância para. Outra instância pode reservar o mesmo item. Nenhum prazo numérico é definido.
- Q: Se um worker para e a instância do Finder continua em execução, outro worker pode reservar o item que ainda não foi concluído? → A: Sim. A reserva incompleta acaba quando o worker para, mesmo com a instância no ar. Outro worker pode reservar o item. Nenhum prazo numérico é definido.
- Q: Com o worker e a instância no ar, a reserva exclusiva e a conclusão do consumo acontecem no mesmo passo? → A: Sim. No mesmo passo o item sai da pendência e o código fica disponível para a etapa futura. O estado reservado sem conclusão só permanece se o worker ou a instância para no meio desse passo.
- Q: Dois itens livres com o mesmo código de apontamento contam como dois consumos? → A: Não. Só o primeiro item com esse código é concluído. O segundo não é concluído e não permanece reservado.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Consumir um item e concluir o consumo (Priority: P1)

Uma instância do Finder lê um item ainda livre da fila Redis. O item traz somente o código de apontamento do log no OpenObserve. No mesmo passo, um worker reserva o item, confirma a conclusão e o retira da condição de pendente. O código de apontamento continua disponível para uma etapa futura. Esta etapa não busca o log no OpenObserve e não nomeia o lugar onde o código fica. O estado reservado sem conclusão só permanece se o worker ou a instância para no meio desse passo.

**Why this priority**: É o recorte pedido para esta etapa. A reserva exclusiva e a conclusão do consumo são o resultado mínimo, antes de abrir o log no OpenObserve.

**Independent Test**: Colocar um item cujo conteúdo é só o código de apontamento, iniciar uma instância do Finder e verificar que, no mesmo passo, o item sai da condição de pendente, o código continua disponível para uma etapa futura e o log no OpenObserve não é aberto. Em seguida, interromper o worker no meio do passo e verificar que o item não conta como concluído.

**Acceptance Scenarios**:

1. **Given** a fila contém um item livre e o conteúdo desse item é só o código de apontamento no OpenObserve, **When** um worker conclui o passo de consumo com a instância no ar, **Then** no mesmo passo o item sai da condição de pendente, não volta a ser entregue como item livre, o código continua disponível para uma etapa futura e o log no OpenObserve não é aberto.
2. **Given** o mesmo passo de um item está em andamento e o worker ainda não parou, **When** outro worker consome a fila, **Then** esse item não é entregue como item livre.
3. **Given** um worker parou no meio do passo, antes da confirmação, e a instância continua no ar, **When** outro worker consome a fila, **Then** o item não conta como concluído e esse outro worker pode reservá-lo, sem prazo numérico.
4. **Given** dois itens livres com o mesmo código de apontamento, **When** os passos de consumo desses itens terminam, **Then** um item é concluído, o código fica disponível uma vez para a etapa futura, e o outro item não é concluído e não permanece reservado.

---

### User Story 2 - Várias instâncias sem corrida (Priority: P2)

Várias instâncias do Finder consomem a mesma fila ao mesmo tempo. Cada item é reservado por no máximo uma instância.

**Why this priority**: O pedido exige consumo com diversas instâncias do Finder, sem racing condition. A história 1, sozinha, cobre uma instância.

**Independent Test**: Subir pelo menos duas instâncias do Finder contra a mesma fila, com um conjunto conhecido de itens, e conferir que cada item fica reservado por apenas uma instância por vez. Parar uma instância depois da reserva e antes da confirmação, e conferir que outra instância pode reservar esse item, sem prazo numérico.

**Acceptance Scenarios**:

1. **Given** duas ou mais instâncias do Finder ativas e uma fila com vários itens livres, **When** elas consomem a fila ao mesmo tempo, **Then** cada item é reservado por apenas uma instância.
2. **Given** o mesmo passo de um item está em andamento e o worker ainda não parou, **When** outra instância consome a fila nesse período, **Then** ela recebe outro item livre, ou nenhum, e não o item desse passo.
3. **Given** uma instância reservou um item e parou antes de confirmar a conclusão, **When** outra instância consome a fila, **Then** o item não conta como concluído e essa outra instância pode reservá-lo, sem prazo numérico.

---

### Edge Cases

- Fila sem itens livres: nenhum item entra em processamento.
- Redis indisponível ou falha na leitura: essa tentativa não conclui item. Nova tentativa, prazo e fila de erro não foram informados e ficam fora desta etapa.
- A instância para depois da reserva e antes da confirmação da conclusão: o item não conta como concluído, a reserva acaba e outra instância pode reservar o mesmo item. Nenhum prazo numérico se aplica. O código de apontamento fica disponível para a etapa futura só depois de uma conclusão.
- O worker para depois da reserva e antes da confirmação, com a instância ainda no ar: o item não conta como concluído, a reserva acaba e outro worker pode reservar o mesmo item. Nenhum prazo numérico se aplica.
- Dois ou mais itens com o mesmo código de apontamento: uma conclusão para esse código. Os outros itens não são concluídos e não permanecem reservados. O código fica disponível uma vez para a etapa futura.
- Item sem código de apontamento, ou item com conteúdo além desse código: está fora do contrato aceito em Q2 A. Nenhum tratamento adicional foi informado e nenhum foi criado.
- O código de apontamento permanece distinto do conteúdo do log (constituição, princípio IV).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: O Finder MUST consumir itens livres da fila Redis em stream.
- **FR-002**: O conteúdo do item MUST ser somente o código de apontamento, no OpenObserve, do log coletado. O Finder MUST NOT tratar esse código como o conteúdo do log.
- **FR-003**: Reserva exclusiva e conclusão do consumo MUST ocorrer no mesmo passo quando o worker e a instância seguem no ar até o fim desse passo.
- **FR-004**: Com mais de uma instância do Finder na mesma fila, o mesmo item MUST ser reservado por no máximo uma instância de cada vez.
- **FR-005**: No mesmo passo de consumo concluído, o Finder MUST retirar o item da condição de pendente e o código de apontamento MUST continuar disponível para uma etapa futura. Esta etapa MUST NOT abrir o log no OpenObserve e MUST NOT definir o nome do lugar onde o código fica. Agregar o log e classificar ocorrência ficam fora. O fluxo completo da constituição (fila, depois OpenObserve, depois agregação) não se encerra aqui.
- **FR-006**: Enquanto o mesmo passo está em andamento e o worker ainda não parou, o item MUST permanecer indisponível como item livre. O estado reservado sem conclusão MUST permanecer só se esse worker parar no meio do passo, mesmo com a instância no ar, ou se a instância parar no meio do passo. Nesse caso o item MUST NOT ser registrado como concluído, a reserva MUST acabar e o item MUST voltar a ficar livre para outro worker reservá-lo. Nenhum prazo numérico define essa liberação.
- **FR-007**: Nomes de stream, grupo e campos MUST NOT ser definidos por esta spec. O protótipo que usa o stream `fila:logs`, o campo `id` e a leitura a partir do id `110` MUST NOT ser tratado como contrato. Esses nomes são entrada externa, ainda não fornecida.
- **FR-008**: O trabalho desta feature MUST permanecer no consumo da fila, na reserva do item e na conclusão desse consumo. Partes do módulo não relacionadas a esse consumo MUST NOT ser alteradas. Novas dependências MUST ser somente as necessárias a este consumo.
- **FR-009**: Para o mesmo código de apontamento, o Finder MUST concluir um item. Os demais itens com esse código MUST NOT ser concluídos e MUST NOT permanecer reservados. O código MUST ficar disponível uma vez para a etapa futura. Se dois itens com esse código terminam no mesmo momento, o resultado exigido é uma conclusão e nenhuma reserva restante no outro. Qual dos dois é o concluído não está definido.

### Key Entities

- **Item da fila**: Conteúdo limitado ao código de apontamento do log no OpenObserve. O mesmo código pode aparecer em mais de um item. Desses itens, um é concluído. Nome de stream, grupo e campos não fazem parte desta spec.
- **Reserva**: Condição do item no meio do mesmo passo de consumo, indisponível para outro worker. No fim normal desse passo ela coincide com a conclusão. Ela permanece sem conclusão só se o worker ou a instância para no meio do passo. Nesse caso a reserva acaba e o item volta a ficar livre.
- **Conclusão do consumo**: Fim normal do mesmo passo de reserva e conclusão. O item sai da condição de pendente e o código de apontamento continua disponível para uma etapa futura. O nome do lugar onde o código fica não faz parte desta spec. A conclusão não inclui leitura do log no OpenObserve.
- **Instância do Finder**: Processo do módulo que consome a mesma fila que outras instâncias.
- **Worker**: Executor, dentro do Finder, que reserva o item e conclui o consumo. O repositório documenta `FINDER_WORKERS` com padrão 4. Esta spec não redefine esse número e não afirma que ele seja o grau de paralelismo deste consumo.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Em um conjunto conhecido de itens, com duas ou mais instâncias do Finder consumindo a mesma fila ao mesmo tempo, cada item tem no máximo uma reserva ativa. A contagem de reservas simultâneas do mesmo item é zero. Uma reserva nova, depois que o worker ou a instância anterior parou sem concluir, conta como reserva seguinte.
- **SC-004**: Quando o worker que reservou para antes da confirmação, mesmo com a instância no ar, ou quando a instância para, esse item não entra na contagem de concluídos e fica livre para uma reserva de outro worker. A liberação não depende de prazo numérico.
- **SC-002**: No passo concluído com o worker e a instância no ar, o observador vê o item sair da pendência e o código de apontamento disponível para a etapa futura, no mesmo passo. O estado reservado sem conclusão aparece quando o worker ou a instância para no meio do passo. Esta etapa não abre o log no OpenObserve.
- **SC-003**: Cada item concluído sai da condição de pendente uma vez e não é entregue de novo como item livre. O código desse item fica disponível para a etapa futura. Itens não concluídos ficam fora da contagem de concluídos. O log no OpenObserve permanece fechado nesta etapa.
- **SC-005**: Em um conjunto de itens livres com o mesmo código de apontamento, a contagem de conclusões desse código é um. Os demais itens desse código não permanecem reservados.

## Assumptions

Decisões e fatos abaixo têm fonte no pedido, na resposta do solicitante, na constituição ou em arquivo lido.

- Fonte: resposta do solicitante em 2026-10-29, Q1 B. O worker reserva o item. Ao terminar o consumo, o Finder confirma a conclusão e retira o item da condição de pendente, sem buscar o log no OpenObserve.
- Fonte: resposta do solicitante em 2026-10-29, Q2 A. O item carrega somente o código de apontamento no OpenObserve. Nomes de stream, grupo e campos continuam indefinidos. O protótipo não é o contrato. A spec não cria esses nomes. O planejamento do contrato concreto espera esses nomes.
- Fonte: resposta do solicitante em 2026-10-29, Q3 A. O aceite desta etapa é a ausência de consumo duplicado entre instâncias. Meta numérica de vazão, quantidade de instâncias além do mínimo de duas e tempo até a reserva ficam de fora até serem informadas.
- Fonte: resposta do solicitante em 2026-10-29, esclarecimento A. Ao sair da pendência, o código de apontamento continua disponível para a etapa futura. Esta etapa não abre o OpenObserve e não nomeia o lugar onde o código fica.
- Fonte: resposta do solicitante em 2026-10-29, esclarecimento A. Se a instância para depois da reserva e antes da confirmação, a reserva incompleta acaba. Outra instância pode reservar o mesmo item. Nenhum prazo numérico é definido. O item não conta como concluído.
- Fonte: resposta do solicitante em 2026-10-29, esclarecimento A. Se o worker para antes da confirmação e a instância continua no ar, a reserva incompleta acaba. Outro worker pode reservar o mesmo item. Nenhum prazo numérico é definido. O item não conta como concluído.
- Fonte: resposta do solicitante em 2026-10-29, esclarecimento A. Com o worker e a instância no ar, reserva e conclusão ocorrem no mesmo passo. O item sai da pendência e o código fica disponível para a etapa futura. O estado reservado sem conclusão só permanece se o worker ou a instância para no meio do passo.
- Fonte: resposta do solicitante em 2026-10-29, esclarecimento B. Dois itens livres com o mesmo código de apontamento não contam como dois consumos concluídos. Só um item com esse código é concluído. O outro não é concluído e não permanece reservado. Qual item vence quando os dois terminam no mesmo momento não foi definido.
- Fonte: pedido. A fila Redis deste consumo ainda não está implementada como mecanismo de coleta.
- Fonte: pedido e constituição (princípio IV). Esta etapa inclui consumo, reserva e conclusão desse consumo. Não inclui busca no OpenObserve nem agregação do log.
- Fonte: constituição (princípio V) e pedido. Nova tentativa depois de falha de leitura, prazo numérico, fila de erro e tratamento de item fora do contrato não foram informados e não foram criados. A liberação do item quando o worker ou a instância para antes da confirmação está definida no esclarecimento desta sessão.
- Fonte: README e constituição (Restrições Técnicas). O módulo é descrito como Go, arquitetura hexagonal leve, consumo de Redis Streams e análise em paralelo com pool de workers. Esta feature permanece dentro dessa descrição.
- Fonte: `go.mod`, lido em 2026-10-29. `github.com/redis/go-redis/v9 v9.23.0` aparece como dependência indireta. O README ainda lista essa biblioteca como candidata e não adicionada. A constituição manda não tratar candidata do README como já presente. Esta spec não decide qual dos dois textos prevalece.
- Fonte: `internal/config/config.go` e README. Existe `FINDER_WORKERS`, padrão 4. Não foi assumido que este consumo use esse valor.
- Fonte: protótipo em `internal/adapter/in/queue/redisStream.go`, lido em 2026-10-29. Ele lê `fila:logs`, o campo `id`, a partir do id `110`, e apaga a mensagem. Q2 A exclui esse protótipo do contrato.
