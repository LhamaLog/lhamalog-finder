# Feature Specification: Corrigir entrega da fila

**Feature Branch**: `002-corrigir-entrega-fila`

**Created**: 2026-11-05

**Status**: Draft

**Input**: User description: "Ajuste no pull request da revisão do consumo da fila. Corrigir a conclusão que não fecha junto com a saída da pendência, a pendência que não volta ao fluxo, a prioridade conforme a ADR 169, e a deduplicação sem organização. Junto: item fora do contrato não pode ficar preso; o processo em execução precisa atender a consulta de saúde já prometida."

## Clarifications

### Session 2026-11-05

- Q: O código de apontamento é único em todo o produto, ou único por organização? → A: Único no produto inteiro. A deduplicação continua só pelo código. Duas organizações com o mesmo código são uma conclusão. O item não ganha organização por causa desta decisão.
- Q: Com duas instâncias consumindo a mesma fila ao mesmo tempo, a ordem já recebida exige uma reserva de cada vez na sequência da fila, ou cada instância pode reservar um item diferente ao mesmo tempo sem escolher outro por causa do nível? → A: Cada instância reserva o próximo item livre que a fila já ordenada lhe entrega e não escolhe outro por nível. Duas instâncias podem reservar itens diferentes ao mesmo tempo. A ordem em que terminam pode diferir da ordem da fila.
- Q: Quando um item reservado por uma instância que parou volta a ficar livre, ele volta ao lugar original na fila já ordenada ou vai para o fim da fila? → A: Volta ao lugar original. Enquanto estiver livre, é entregue antes dos itens que já estavam atrás dele. Um item mais novo e menos prioritário não passa na frente.
- Q: Se a tentativa de completar uma conclusão interrompida falha antes dos efeitos que faltam, o processo continua no ar e tenta de novo sem prazo, ou encerra a execução? → A: O processo continua no ar. A tentativa seguinte, sem prazo numérico, completa os efeitos que faltam. Enquanto isso o item permanece parcial e não conta como segunda conclusão.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Conclusão interrompida termina num único desfecho (Priority: P1)

Uma instância reserva um item aceito e começa a conclusão. A conclusão só está completa quando as três coisas são verdadeiras ao mesmo tempo: o item saiu da pendência, existe uma conclusão daquele código no produto, e o código ficou disponível para a etapa futura. Se o processo para depois que alguma dessas coisas já aconteceu e antes que as três estejam feitas, a retomada termina esse mesmo desfecho uma vez. A retomada não abre uma segunda conclusão e não descarta o item como repetido deixando o código indisponível. Se a tentativa de completar os efeitos que faltam falha, o processo continua no ar. O item permanece parcial. Uma tentativa seguinte, sem prazo numérico, completa esses efeitos.

**Why this priority**: A revisão bloqueou o merge neste ponto. Uma marcação antecipada, seguida de queda, faz a reentrega parecer repetida e encerra o item sem garantir o desfecho completo.

**Independent Test**: Concluir um item aceito e interromper o processo depois do primeiro efeito da conclusão e antes dos três efeitos estarem presentes. Retomar e conferir uma conclusão, o item fora da pendência e o código disponível uma vez. Repetir a entrega desse item e conferir que a contagem de conclusões não sobe. Numa outra tentativa, falhar ao completar os efeitos que faltam e conferir que o processo segue no ar, o item permanece parcial e a tentativa seguinte chega ao mesmo desfecho único.

**Acceptance Scenarios**:

1. **Given** um item aceito foi reservado e nenhum efeito da conclusão aconteceu, **When** a instância para antes da confirmação, **Then** o item não conta como concluído, o código não fica disponível por causa dessa tentativa e o item volta a ficar livre, sem prazo numérico.
2. **Given** um item aceito já teve pelo menos um efeito da conclusão e ainda falta outro, **When** o processo para e a entrega é retomada, **Then** o item sai da pendência, a contagem de conclusões daquele código permanece um e o código fica disponível uma vez para a etapa futura.
3. **Given** a conclusão já está completa, **When** o mesmo item é entregue de novo, **Then** ele não volta como item livre, não gera segunda conclusão e o código continua disponível uma vez.
4. **Given** dois itens aceitos do mesmo código, ainda que de organizações diferentes, e um deles foi interrompido no meio da conclusão, **When** os dois chegam ao fim da retomada, **Then** existe uma conclusão e o outro item não permanece reservado.
5. **Given** um item aceito já teve pelo menos um efeito da conclusão e a tentativa de completar o restante falha, **When** o processo segue no ar, **Then** o item permanece parcial, não conta como segunda conclusão, e uma tentativa seguinte, sem prazo numérico, completa os efeitos que faltam.

---

### User Story 2 - Item deixado por instância parada volta ao fluxo (Priority: P2)

Uma instância reserva um item e para antes de qualquer efeito da conclusão. O item não pode permanecer preso a essa instância. Outra instância em execução, ou a mesma instância depois de iniciar de novo, pode reservá-lo. Ao voltar a ficar livre, o item ocupa o lugar original na fila já ordenada. Enquanto estiver livre, é entregue antes dos itens que já estavam atrás dele. Isso vale no arranque e enquanto outras instâncias seguem consumindo. Nenhum prazo numérico define essa volta.

**Why this priority**: A spec anterior já exige esse desfecho. A revisão mostrou que o item reservado na fila compartilhada não volta ao fluxo quando o processo morre depois de lê-lo.

**Independent Test**: Com duas instâncias na mesma fila, reservar um item na primeira que esteja à frente de outros itens livres, encerrar essa instância antes de qualquer efeito da conclusão e conferir que a segunda reserva esse item, antes dos que estavam atrás, sem espera numérica. Reiniciar a primeira e conferir que um item deixado por ela também volta ao lugar original.

**Acceptance Scenarios**:

1. **Given** uma instância reservou um item e parou antes de qualquer efeito da conclusão, **When** outra instância consome a fila, **Then** o item não conta como concluído e essa outra instância pode reservá-lo, sem prazo numérico.
2. **Given** uma instância reservou um item e o processo foi encerrado antes de qualquer efeito da conclusão, **When** essa instância inicia de novo, **Then** o item deixado por ela volta a poder ser reservado e não permanece preso ao processo encerrado.
3. **Given** o mesmo passo ainda está em andamento e a instância segue no ar, **When** outra instância consome a fila, **Then** ela não recebe esse item como livre.
4. **Given** o item deixado pela instância parada estava à frente de outros itens livres, **When** ele volta a ficar livre, **Then** a próxima reserva é esse item, antes dos que estavam atrás dele. Um item mais novo não é reservado enquanto esse seguir livre.

---

### User Story 3 - Consumir a fila já ordenada (Priority: P3)

A fila chega ao Finder já ordenada. A ordem é a da ADR 169: todos os itens do nível mais prioritário vêm primeiro; no mesmo nível, o horário de geração mais antigo vem primeiro. Os níveis aceitos, do mais prioritário ao menos prioritário, são Very High, High, Medium, Low e Very Low. Medium é o nível padrão dessa ordenação. O Finder não escolhe o nível, não reordena e não cria nível. Com uma instância, ele reserva o próximo item livre na ordem em que a fila já ordenada o apresenta. Com duas ou mais instâncias, cada uma reserva o próximo item livre que a fila lhe entrega. Elas podem reservar itens diferentes ao mesmo tempo. A ordem em que terminam pode diferir da ordem da fila.

**Why this priority**: A tarefa pede consumo por prioridade. A ADR 169 fecha esse ponto: a ordenação ocorre antes do Finder. Sem este recorte, a revisão seguiria pedindo uma seleção de nível que a decisão proíbe.

**Independent Test**: Com uma instância, entregar uma fila já ordenada pelos cinco níveis, com empate resolvido pelo horário de geração mais antigo, e conferir que a sequência de reservas segue essa ordem. Com duas instâncias, conferir que cada uma recebe o próximo item livre que a fila lhe entrega, que as duas podem reservar itens diferentes ao mesmo tempo e que o término fora da ordem da fila não reprova o cenário.

**Acceptance Scenarios**:

1. **Given** a fila já ordenada tem itens livres em mais de um nível, **When** uma instância reserva esses itens, **Then** a sequência de reservas é a ordem já apresentada, com o nível mais prioritário antes do menos prioritário.
2. **Given** dois itens livres do mesmo nível e horários de geração diferentes, **When** uma instância reserva esses itens, **Then** o de horário de geração mais antigo é reservado primeiro.
3. **Given** a fila já ordenada apresenta o próximo item livre, **When** o Finder consome, **Then** ele reserva esse item e não escolhe outro por causa do nome do nível.
4. **Given** duas instâncias consomem a mesma fila já ordenada ao mesmo tempo, **When** ambas reservam, **Then** cada uma fica com o próximo item livre que a fila lhe entrega, as duas podem reservar itens diferentes ao mesmo tempo e o término de uma antes da outra não reprova a ordem.

---

### User Story 4 - O mesmo código conta uma vez no produto (Priority: P4)

O código de apontamento é único no produto inteiro. Dois itens com o mesmo código são uma conclusão, mesmo quando vêm de organizações diferentes. O segundo não é concluído e não permanece reservado. O item não passa a carregar organização por causa desta regra.

**Why this priority**: A revisão apontou deduplicação sem organização. A resposta desta sessão fixa a premissa que faltava: o código já distingue a conclusão no produto, e a organização fica de fora.

**Independent Test**: Entregar o mesmo código em nome de duas organizações e conferir uma conclusão, o código disponível uma vez e o outro item fora da reserva.

**Acceptance Scenarios**:

1. **Given** duas organizações apresentam o mesmo código de apontamento em itens aceitos, **When** os passos de consumo terminam, **Then** existe uma conclusão, o código fica disponível uma vez e o outro item não é concluído e não permanece reservado.
2. **Given** o mesmo código já foi concluído, **When** outra organização entrega um item aceito com esse código, **Then** esse item não é concluído e não permanece reservado.

---

### User Story 5 - Item fora do contrato não fica preso (Priority: P5)

Um item sem código de apontamento, ou com conteúdo além desse código, está fora do contrato. Ele não entra em conclusão e não deixa o código disponível. Também não permanece reservado nem pendente. O descarte fica observável: um observador distingue esse item de um item concluído e de um item ainda pendente. Itens aceitos que vêm depois seguem o fluxo.

**Why this priority**: A revisão mostrou que o item recusado fica reservado e pendente até o processo acabar. Com a volta das pendências ao fluxo, esse item voltaria sem fim.

**Independent Test**: Entregar um item sem código e outro com conteúdo extra, seguidos de um item aceito. Conferir que os dois primeiros não estão pendentes nem reservados, não contam como conclusão e o descarte é observável, e que o item aceito conclui.

**Acceptance Scenarios**:

1. **Given** a fila entrega um item sem código de apontamento, **When** o Finder o toma, **Then** o item é descartado, a reserva acaba, ele sai da pendência, não conta como concluído e o descarte é observável.
2. **Given** a fila entrega um item com conteúdo além do código de apontamento, **When** o Finder o toma, **Then** o desfecho é o mesmo do item sem código, e esse conteúdo não vira conclusão.
3. **Given** um item fora do contrato foi descartado e existe um item aceito livre em seguida, **When** o consumo continua, **Then** o item aceito pode ser reservado e concluído.

---

### User Story 6 - Processo em execução atende a consulta de saúde (Priority: P6)

Enquanto o processo do Finder está em execução, a consulta de saúde já prometida é atendida. Esta feature não cria outra consulta e não muda o significado da existente. Com o processo parado, a consulta não é atendida.

**Why this priority**: O ambiente local já usa essa consulta para saber se o processo está no ar. O contrato da feature anterior a mantém. O processo atual não a atende.

**Independent Test**: Iniciar o processo e obter resposta positiva da consulta de saúde já prometida. Encerrar o processo e conferir que a mesma consulta não é atendida.

**Acceptance Scenarios**:

1. **Given** o processo do Finder está em execução, **When** a consulta de saúde já prometida é feita, **Then** ela é atendida.
2. **Given** o processo do Finder não está em execução, **When** a mesma consulta é feita, **Then** ela não é atendida.
3. **Given** esta feature entra em vigor, **When** se observa o que o Finder expõe, **Then** não há consulta de saúde nova e o significado da consulta já prometida permanece.

---

### Edge Cases

- Parada antes de qualquer efeito da conclusão: o item não é concluído e volta a ficar livre, sem prazo numérico.
- Parada depois de um efeito e antes dos três: a retomada completa o mesmo desfecho uma vez. Não há segunda conclusão e o código não fica indisponível.
- Falha ao completar os efeitos que faltam: o processo continua no ar e o item permanece parcial. Essa falha não conta como segunda conclusão e não encerra o processo. A tentativa seguinte, sem prazo numérico, completa esses efeitos. A falha não exige repetição imediata até o sucesso.
- Conclusão já completa e o mesmo item entregue de novo: a contagem de conclusões não aumenta.
- Duas instâncias retomam o mesmo item interrompido: o resultado continua uma conclusão, ou nenhuma, conforme o ponto da interrupção. Não há duas conclusões.
- Instância parada com item reservado e outras instâncias seguindo no ar: o item volta ao lugar original, sem prazo numérico. Enquanto livre, é entregue antes dos itens que já estavam atrás dele.
- Instância reiniciada com item que ela mesma deixou reservado: esse item volta ao lugar original e pode ser reservado antes dos que estavam atrás.
- Item fora do contrato no meio de itens aceitos: os aceitos seguem; o item fora do contrato não fica pendente.
- Vários itens fora do contrato seguidos: cada um é descartado e observável, e nenhum permanece reservado.
- Fila já ordenada com um único nível: com uma instância, a ordem observada é a do horário de geração mais antigo primeiro.
- Duas ou mais instâncias ao mesmo tempo: cada uma reserva o próximo item livre que a fila lhe entrega. Reservas simultâneas de itens diferentes são válidas. A ordem de término pode diferir da fila. Nenhuma escolhe outro item por nível.
- Fila sem item livre: nenhuma conclusão nessa tentativa.
- Falha ao obter o próximo item, antes de reservá-lo: essa tentativa não conclui item. Prazo numérico e fila de erro continuam fora.
- Duas organizações com o mesmo código de apontamento: uma conclusão no produto. O outro item não permanece reservado. O item não ganha organização.
- O código de apontamento permanece distinto do conteúdo do log.
- Esta etapa não abre o log no OpenObserve.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Uma conclusão completa MUST reunir, para o mesmo item aceito, três efeitos: o item fora da pendência, uma conclusão daquele código no produto e o código disponível uma vez para a etapa futura.
- **FR-002**: Se a instância para antes de qualquer efeito da conclusão, o item MUST NOT contar como concluído, o código MUST NOT ficar disponível por causa dessa tentativa, a reserva MUST acabar e o item MUST voltar a ficar livre. Nenhum prazo numérico define essa volta.
- **FR-003**: Se já existe algum efeito da conclusão e ainda falta outro, a retomada MUST completar os efeitos que faltam desse mesmo item. A contagem de conclusões MUST permanecer um. O item MUST NOT ser encerrado como repetido enquanto o código não estiver disponível. Se essa tentativa falha antes de completar os efeitos que faltam, o processo MUST continuar no ar, o item MUST permanecer parcial e MUST NOT contar como segunda conclusão. Uma tentativa seguinte, sem prazo numérico, MUST completar os efeitos que faltam. A falha MUST NOT exigir repetição imediata até o sucesso.
- **FR-004**: Com a conclusão já completa, uma nova entrega do mesmo item MUST NOT criar outra conclusão e MUST NOT devolver o item como livre.
- **FR-005**: Um item reservado por uma instância que parou antes de qualquer efeito da conclusão MUST voltar a poder ser reservado por outra instância em execução e pela mesma instância depois que ela inicia de novo. O item MUST NOT permanecer preso ao processo encerrado. Ao voltar a ficar livre, MUST ocupar o lugar original na fila já ordenada. Enquanto estiver livre, MUST ser entregue antes dos itens que já estavam atrás dele. Um item mais novo MUST NOT ser reservado na frente dele enquanto ele seguir livre. Nenhum prazo numérico define essa volta.
- **FR-006**: Enquanto o passo segue em andamento e a instância está no ar, o item MUST permanecer indisponível como item livre para outra instância.
- **FR-007**: O Finder MUST reservar o próximo item livre na ordem em que a fila já ordenada o apresenta. Essa ordem é a da ADR 169: o nível mais prioritário antes do menos prioritário e, no mesmo nível, o horário de geração mais antigo primeiro. Com uma instância, a sequência de reservas MUST seguir essa ordem. Com duas ou mais instâncias ao mesmo tempo, cada uma MUST reservar o próximo item livre que a fila lhe entrega e MUST NOT escolher outro por causa do nível. Reservas simultâneas de itens diferentes são válidas. A ordem em que essas reservas terminam pode diferir da ordem da fila.
- **FR-008**: Os nomes de nível aceitos MUST ser, nesta ordem, Very High, High, Medium, Low e Very Low. A grafia faz parte do nome. Medium permanece o nível padrão da ordenação feita antes do Finder. O Finder MUST NOT escolher, renomear, traduzir, abreviar ou criar nível.
- **FR-009**: Um item sem código de apontamento, ou com conteúdo além desse código, MUST ser descartado. O descarte MUST acabar a reserva, MUST tirar o item da pendência, MUST NOT contar como conclusão e MUST NOT deixar o código disponível. O descarte MUST ser observável, de modo que se distinga de uma conclusão e de um item ainda pendente.
- **FR-010**: Um item fora do contrato MUST NOT impedir a reserva e a conclusão de um item aceito seguinte.
- **FR-011**: Enquanto o processo do Finder está em execução, a consulta de saúde já prometida MUST ser atendida. Esta feature MUST NOT criar outra consulta e MUST NOT alterar o significado da existente. Com o processo parado, a consulta não é atendida.
- **FR-012**: O código de apontamento é único no produto inteiro. Itens com o mesmo código MUST ter uma só conclusão, inclusive quando vêm de organizações diferentes. A organização MUST NOT entrar nessa unicidade. O item MUST NOT ganhar organização por causa desta regra.
- **FR-013**: Esta feature MUST permanecer na entrega do consumo já especificado: conclusão recuperável, volta do item deixado por instância parada, preservação da ordem já recebida, descarte do item fora do contrato e atendimento da consulta de saúde já prometida. Ela MUST NOT abrir o log no OpenObserve, MUST NOT nomear stream, grupo ou campos e MUST NOT introduzir prazo numérico.

### Key Entities

- **Item da fila**: Registro livre, reservado ou concluído. O conteúdo aceito continua sendo o código de apontamento do log. Item sem esse código, ou com conteúdo além dele, é fora do contrato e descartável.
- **Conclusão do consumo**: Desfecho único de um item aceito. Só está completa com o item fora da pendência, uma conclusão daquele código no produto e o código disponível uma vez para a etapa futura.
- **Fila já ordenada**: Fila entregue ao Finder na ordem da ADR 169. Com uma instância, a sequência de reservas segue essa ordem. Com várias instâncias, cada uma reserva o próximo item livre que a fila lhe entrega. Elas podem reservar itens diferentes ao mesmo tempo, e a ordem de término pode diferir da fila. Um item que volta a ficar livre ocupa o lugar original. O Finder não escolhe o nível.
- **Nível de prioridade**: Um dos cinco nomes aceitos. Medium é o padrão da ordenação anterior ao Finder. O nível não é uma escolha do Finder.
- **Descarte**: Saída de um item fora do contrato, sem conclusão e sem reserva, observável por quem acompanha o consumo.
- **Consulta de saúde**: Pergunta já prometida sobre o processo estar em execução. Não é uma capacidade nova.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Num conjunto de itens aceitos interrompidos depois do primeiro efeito da conclusão e antes dos três efeitos, a retomada deixa cada código com uma conclusão, zero itens ainda pendentes nesse conjunto e o código disponível uma vez. Uma entrega seguinte não aumenta a contagem. Se a tentativa de completar falha, o processo segue no ar, a contagem de conclusões não sobe nessa falha, o item permanece parcial e a tentativa seguinte chega a esse mesmo desfecho.
- **SC-002**: Quando a instância para depois da reserva e antes de qualquer efeito da conclusão, esse item não entra na contagem de concluídos. Uma instância em execução, ou a mesma instância após iniciar de novo, consegue reservá-lo. Se havia itens livres atrás dele, a próxima reserva é esse item, e a contagem de itens de trás reservados enquanto ele seguia livre é zero. A observação dessa volta não depende de prazo numérico. Itens presos ao processo encerrado, nesse conjunto, são zero.
- **SC-003**: Com uma instância, numa fila já ordenada por nível e, no empate, pelo horário de geração mais antigo, a sequência de reservas é essa ordem. A contagem de itens reservados fora dela é zero. Com duas ou mais instâncias, a contagem de itens escolhidos por nível, fora do próximo item que a fila entregou àquela instância, é zero. Reservas simultâneas de itens diferentes contam como sucesso. A ordem de término pode diferir da ordem da fila e não reprova o cenário.
- **SC-004**: Num conjunto de itens fora do contrato, depois da tentativa de consumo, zero permanecem pendentes ou reservados, zero entram na contagem de concluídos e cada descarte é distinguível de uma conclusão. Um item aceito entregue depois deles pode ser concluído.
- **SC-005**: Com o processo em execução, a consulta de saúde já prometida é atendida em cada tentativa desse cenário. Com o processo parado, a mesma consulta não é atendida. O número de consultas novas introduzidas por esta feature é zero.
- **SC-006**: Num conjunto de itens aceitos com o mesmo código, vindos de duas organizações, a contagem de conclusões desse código é um. Os demais itens desse código não permanecem reservados. O item aceito não passa a trazer organização.

## Assumptions

- Fonte: spec `specs/001-consumir-fila-redis/spec.md`, lida em 2026-11-05. Esta feature ajusta a entrega daquela. Continuam valendo o item aceito como código de apontamento, a reserva exclusiva, a conclusão no mesmo passo quando o processo segue até o fim, uma conclusão por código no produto, a ausência de prazo numérico e a etapa sem abrir o OpenObserve.
- Fonte: revisão do pull request e leitura de `internal/application/consumir.go` e `internal/adapter/in/queue/redisStream.go` em 2026-11-05. Hoje a marcação do código pode ocorrer antes de o item sair da pendência. Se o processo para nesse intervalo, a reentrega trata o código como já concluído e encerra o item sem repetir o desfecho. Esta spec exige o desfecho único e não escolhe o mecanismo.
- Fonte: a mesma revisão e o FR-006 da feature 001. O item lido e não confirmado permanece com a instância que morreu e não volta ao fluxo. O contrato da feature 001 deixou o sinal entre processos fora do texto. Esta spec exige o desfecho, no arranque e com outras instâncias no ar, e mantém a proibição de prazo numérico. O sinal em si não é nomeado aqui.
- Fonte: ADR 169, "Nomes aceitos dos níveis de prioridade", status Aceito, 2026-11-05, documento informado no pedido. A Varredura ordena a fila. O Finder não escolhe o nível e consome a fila já ordenada. Os nomes deixam de ser provisórios: Very High, High, Medium, Low e Very Low. Medium segue como padrão. No empate, o horário de geração mais antigo vem primeiro. Esta spec não reabre a quantidade de níveis nem passa a ordenação para o Finder.
- Fonte: a mesma ADR e a tarefa que pede consumo por prioridade. Atender a prioridade, neste módulo, é preservar a ordem já recebida. Não é selecionar um nível dentro do Finder.
- Fonte: resposta do solicitante em 2026-11-05, esclarecimento da ordem com várias instâncias, opção A. Cada instância reserva o próximo item livre que a fila lhe entrega. Duas instâncias podem reservar itens diferentes ao mesmo tempo. A ordem de término pode diferir da ordem da fila.
- Fonte: resposta do solicitante em 2026-11-05, esclarecimento do lugar do item recuperado, opção A. O item volta ao lugar original. Enquanto livre, é entregue antes dos itens que já estavam atrás dele. Um item mais novo e menos prioritário não passa na frente.
- Fonte: resposta do solicitante em 2026-11-05, esclarecimento da falha ao completar a conclusão, opção A. O processo continua no ar. A tentativa seguinte, sem prazo numérico, completa os efeitos que faltam. O item permanece parcial e não conta como segunda conclusão. A falha não exige repetição imediata até o sucesso.
- Fonte: pedido desta spec, que mandou definir descarte, fila de erro ou retry observável. O item fora do contrato não se torna aceito ao ser tentado de novo, então retry repetiria o mesmo item. Uma fila de erro nomearia um destino que não foi informado. O desfecho adotado é descarte observável, com reserva limpa e saída da pendência, sem conclusão. A feature 001 tinha deixado esse item sem tratamento; esta spec substitui esse silêncio.
- Fonte: `specs/001-consumir-fila-redis/contracts/consumo-fila.md`, lido em 2026-11-05. Nenhuma consulta HTTP nova. A consulta de saúde já prometida é `GET /health`. O pedido informa que o ambiente local a usa para saber se o processo está no ar e que o processo atual não a atende. Esta spec restaura o atendimento e não altera o contrato dessa consulta.
- Fonte: revisão e feature 001, FR-007. Nomes de stream, grupo e campos continuam de fora.
- Fonte: resposta do solicitante em 2026-11-05, Q1 A. O código de apontamento é único no produto inteiro. Duas organizações com o mesmo código são uma conclusão. A deduplicação permanece só pelo código. O item não ganha organização. A constituição, princípio V, pedia consulta porque a spec 001 e a ADR 169 não traziam essa premissa; a resposta fecha o ponto.
- Textos operacionais que ainda descrevem o consumo anterior passam, no planejamento, a descrever o comportamento desta spec. Esta spec não cria página, variável nem comando novo.
- A sugestão da revisão de estados intermediários e de operações atômicas é exemplo de meio. Não é requisito.
