# Feature Specification: Consumir logs no OpenObserve

**Feature Branch**: `002-consumir-logs-openobserve`

**Created**: 2026-10-31

**Status**: Draft

**Input**: User description: "Ajustar o módulo para consumir logs no OpenObserve de forma genérica, sem usar a fila nesta etapa. A base local de desenvolvimento deixa de incluir o Elasticsearch e passa a incluir o OpenObserve. O apontamento do log a consumir é mockado, porque a fila está vazia e o uso dela não entra agora. Com esse apontamento, o Finder coleta o log armazenado no OpenObserve e o deixa pronto para processamento. O processamento do log não entra nesta etapa. O arquivo otel-logs.zip é o exemplo de logs. Aceites do pedido: código em Go, consumo do log, falha com detalhes completos, execução paralela, clean code, estrutura atual de pastas. Restrições: sem biblioteca descontinuada, sem dependência desnecessária, sem alterar componente fora do escopo. Premissas devem ficar explícitas. Não inventar resultado de teste nem afirmar execução sem tê-la feito."

## Clarifications

### Session 2026-10-31

- Q: Nesta etapa, cada indicação mockada deve levar o Finder a coletar um único log no OpenObserve, ou uma indicação deve trazer todos os logs do exemplo de uma vez? → A: Cada indicação seleciona um único log. Vários logs exigem várias indicações, que podem ser capturadas ao mesmo tempo.
- Q: Quando a captura de um log termina com sucesso, esse log pronto precisa continuar acessível só enquanto o Finder estiver em execução, ou também depois que ele parar? → A: Enquanto a instância está no ar, o log pronto pode ser lido de novo pelo mesmo apontamento. Depois que ela para, esta etapa não exige que ele continue. O lugar não é nomeado.
- Q: Se o mesmo apontamento aparece duas vezes na mesma execução, quantos logs prontos esse apontamento deve ter? → A: Um log pronto. Indicações repetidas do mesmo apontamento convergem nesse único log.
- Q: Para o aceite desta etapa, o log capturado precisa ser um registro do arquivo otel-logs.zip, ou basta um log conhecido já armazenado no OpenObserve? → A: O aceite só vale se o log capturado for um registro do otel-logs.zip. O arquivo precisa ser entregue antes do planejamento.
- Q: O log pronto deve ser o registro do otel-logs.zip como foi armazenado, ou a preparação desta etapa já deve ocultar dados sensíveis? → A: O log pronto é o registro armazenado, inteiro, igual ao registro do arquivo. Ocultar dados sensíveis fica na etapa de processamento.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Capturar o log apontado e deixá-lo pronto (Priority: P1)

O Finder recebe, de uma fonte mockada, o apontamento de um log armazenado no OpenObserve. No aceite, esse log é um registro do arquivo `otel-logs.zip`. Esse apontamento indica qual log coletar. Não é o conteúdo do log. O Finder coleta esse log no OpenObserve e o deixa pronto para uma etapa futura de processamento. O log pronto é o registro armazenado, inteiro. Ocultar dados sensíveis fica na etapa de processamento. Enquanto a instância está no ar, esse log pronto pode ser lido de novo pelo mesmo apontamento. Depois que a instância para, esta etapa não exige que ele continue. O lugar não é nomeado. A fila Redis não é lida. Nesta etapa o Finder não agrega o log e não classifica ocorrência.

**Why this priority**: É o resultado pedido. Sem a captura e a preparação, a troca da base local e o detalhe da falha não entregam o log pronto para processamento.

**Independent Test**: Armazenar no OpenObserve um registro do arquivo `otel-logs.zip`, fazer o mock apontar só para ele, executar a captura com a instância no ar até o fim e verificar que esse registro fica pronto para processamento. Em seguida, ler de novo pelo mesmo apontamento, ainda com a instância no ar, e verificar que o mesmo registro volta. A fila não é lida e nenhuma ocorrência é classificada. Este teste espera a entrega do arquivo.

**Acceptance Scenarios**:

1. **Given** o OpenObserve contém um registro do arquivo `otel-logs.zip` e o mock aponta somente para ele, **When** a captura termina com a instância no ar e o mesmo apontamento é lido de novo, **Then** esse registro fica pronto para processamento, a segunda leitura devolve o mesmo registro, a fila não é lida e nenhuma ocorrência é classificada como erro, alerta ou informação.
2. **Given** a captura terminou com sucesso e a instância parou em seguida, **When** esta etapa é conferida, **Then** ela não exige que o log pronto continue disponível.
3. **Given** o mock não fornece apontamento, **When** a etapa executa, **Then** nenhum log fica pronto para processamento e nenhuma falha de captura é registrada.
4. **Given** o mock aponta para um registro armazenado do arquivo `otel-logs.zip`, **When** a captura termina com sucesso, **Then** o que fica pronto é esse registro, inteiro e igual ao armazenado, e não uma ocorrência processada nem um registro com dados sensíveis ocultos.

---

### User Story 2 - Falha de captura com detalhes completos (Priority: P2)

Quando a captura não conclui, o log não fica pronto para processamento. A falha mostra o apontamento usado e o motivo, de forma que um observador distinga apontamento inutilizável, log ausente e OpenObserve inacessível.

**Why this priority**: O pedido exige que a falha apresente detalhes completos. A história 1, sozinha, cobre só o caminho em que o log é capturado.

**Independent Test**: Repetir a captura três vezes, uma com apontamento inutilizável, uma com log ausente e uma com OpenObserve inacessível, e verificar que cada falha identifica o apontamento e o motivo correspondente, e que nenhum desses logs fica pronto.

**Acceptance Scenarios**:

1. **Given** o mock aponta para um log que não está armazenado, **When** a captura executa, **Then** esse log não fica pronto e a falha mostra esse apontamento e que o log está ausente.
2. **Given** o OpenObserve não pode ser alcançado, **When** a captura executa, **Then** o log não fica pronto e a falha mostra o apontamento e que o armazenamento não pôde ser alcançado, distinta da falha de log ausente.
3. **Given** o mock fornece um apontamento inutilizável, **When** a captura executa, **Then** nenhum log fica pronto e a falha mostra que o apontamento é inutilizável.

---

### User Story 3 - Capturas paralelas com resultado próprio (Priority: P3)

Dois ou mais apontamentos distintos podem ser capturados ao mesmo tempo. Cada um termina no próprio resultado: log pronto ou falha com os próprios detalhes.

**Why this priority**: O pedido pede execução paralela. As histórias 1 e 2 continuam válidas com uma captura de cada vez.

**Independent Test**: Disparar ao mesmo tempo duas capturas com apontamentos distintos, uma com log armazenado e outra com log ausente, e verificar que a primeira deixa o log pronto e a segunda conserva a própria falha, sem apagar o resultado da primeira.

**Acceptance Scenarios**:

1. **Given** dois apontamentos distintos selecionam dois logs armazenados, **When** as capturas ocorrem ao mesmo tempo, **Then** os dois logs ficam prontos, cada um associado ao próprio apontamento.
2. **Given** dois apontamentos distintos e só um dos logs está armazenado, **When** as capturas ocorrem ao mesmo tempo, **Then** o log existente fica pronto e o ausente conserva a própria falha com detalhes completos, sem retirar o log que ficou pronto.
3. **Given** duas indicações com o mesmo apontamento para um log armazenado, **When** as capturas terminam com a instância no ar, **Then** existe um log pronto para esse apontamento e a leitura devolve esse único log.

---

### User Story 4 - Base local com OpenObserve (Priority: P2)

A base local de desenvolvimento passa a oferecer o OpenObserve como lugar onde o log de exemplo fica armazenado e de onde o Finder o captura. O Elasticsearch deixa de fazer parte dessa base. MongoDB e Redis continuam na base. Esta etapa não passa a depender da fila.

**Why this priority**: O pedido pede a troca na base local para que o log possa ser armazenado e capturado no OpenObserve. A história 1 descreve a captura; esta história descreve o lugar local em que isso pode ser feito.

**Independent Test**: Subir a base local desta etapa e verificar que o OpenObserve está disponível para armazenar e ler um log, e que o Elasticsearch não faz parte da base.

**Acceptance Scenarios**:

1. **Given** a base local desta etapa está no ar, **When** uma pessoa verifica o que a base oferece, **Then** o OpenObserve está disponível para armazenar e ler logs e o Elasticsearch não faz parte da base.
2. **Given** a base local está no ar com MongoDB e Redis, **When** a captura desta etapa executa, **Then** ela usa o apontamento mockado e o OpenObserve, e não lê a fila.

---

### Edge Cases

- Mock sem apontamento: nenhum log fica pronto e isso não conta como falha de captura.
- Apontamento inutilizável: falha com detalhes completos. O log não fica pronto.
- OpenObserve inacessível: falha com detalhes completos, distinta da de log ausente. O log não fica pronto. Nova tentativa, prazo e fila de erro não foram informados e ficam fora desta etapa.
- Log ausente para o apontamento: falha com detalhes completos. O log não fica pronto. A etapa não cria um log vazio no lugar do ausente.
- A instância para no meio da captura: o log não fica pronto. Nenhum prazo numérico se aplica. Retomada não foi informada e fica fora desta etapa.
- Indicação que corresponderia a vários logs ou ao exemplo inteiro: fora desta etapa. Cada indicação seleciona um único log.
- Dois apontamentos iguais no mesmo conjunto, ao mesmo tempo ou em sequência: um log pronto para esse apontamento. A leitura, com a instância no ar, devolve esse único log.
- Log pronto com a instância no ar: uma nova leitura pelo mesmo apontamento devolve esse log, inteiro e igual ao registro armazenado. O lugar não é nomeado. Dados sensíveis ocultos ficam fora desta etapa.
- Log pronto depois que a instância para: esta etapa não exige que ele continue disponível.
- A fila contém itens: esta etapa não os lê. O fluxo normativo da constituição (fila, depois OpenObserve, depois agregação) não se encerra aqui.
- Dados já existentes no Elasticsearch local: a retirada do Elasticsearch da base não inclui migração. Migração não foi pedida e fica fora.
- O arquivo `otel-logs.zip` não estava no repositório nem na pasta LhamaLog Workspace quando esta spec foi escrita. O aceite exige um registro desse arquivo, e o arquivo precisa ser entregue antes do planejamento. Os registros internos não são descritos aqui. O Finder captura no OpenObserve. O arquivo não é a fonte da captura.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Nesta etapa o Finder MUST obter o apontamento do log a partir de uma fonte mockada. O Finder MUST NOT ler a fila Redis para obter esse apontamento.
- **FR-002**: Cada apontamento MUST selecionar um único log armazenado no OpenObserve. Vários logs exigem vários apontamentos. Um apontamento MUST NOT selecionar mais de um log nem o exemplo inteiro. O apontamento MUST NOT ser tratado como o conteúdo do log.
- **FR-003**: Quando a captura termina com a instância no ar, o Finder MUST deixar esse log pronto para processamento, inteiro e igual ao registro armazenado, associado ao apontamento que o selecionou. Enquanto a instância está no ar, uma nova leitura por esse apontamento MUST devolver o mesmo log. Depois que a instância para, esta etapa MUST NOT exigir que o log pronto continue. O lugar MUST NOT ser nomeado.
- **FR-004**: Esta etapa MUST NOT agregar o log, MUST NOT classificar ocorrência como erro, alerta ou informação, MUST NOT ocultar dados sensíveis e MUST NOT definir persistência de ocorrência. O processamento do log, inclusive ocultar dados sensíveis, fica para uma etapa futura.
- **FR-005**: Organização, stream, índice, campos e credenciais do OpenObserve MUST NOT ser nomeados por esta spec. São entrada externa, ainda não fornecida. Nomes de protótipo de outra etapa MUST NOT ser tratados como contrato desta captura.
- **FR-006**: Se o apontamento é inutilizável, se o log apontado está ausente ou se o OpenObserve não pode ser alcançado, o Finder MUST NOT deixar esse log pronto. A falha MUST apresentar o apontamento usado e o motivo, de modo que esses três motivos sejam distinguíveis entre si. Uma mensagem genérica que esconda o motivo não atende.
- **FR-007**: Dois ou mais apontamentos distintos MUST poder ser capturados ao mesmo tempo. Cada captura MUST terminar no próprio resultado. A falha de uma MUST conservar os próprios detalhes e MUST NOT apagar o resultado de outra.
- **FR-011**: Se o mesmo apontamento aparece mais de uma vez na mesma execução, o Finder MUST deixar um único log pronto para esse apontamento. Uma leitura por esse apontamento, com a instância no ar, MUST devolver esse único log.
- **FR-008**: A base local de desenvolvimento MUST incluir o OpenObserve como lugar de armazenamento e leitura do log a capturar e MUST NOT incluir o Elasticsearch. MongoDB e Redis permanecem na base. Esta etapa MUST NOT passar a depender da leitura da fila.
- **FR-009**: O aceite desta etapa MUST usar um registro do arquivo `otel-logs.zip` como o log capturado. O arquivo MUST ser entregue antes do planejamento. O Finder MUST capturar esse registro no OpenObserve, a partir do apontamento. O arquivo MUST NOT ser a fonte da captura. Esta spec não descreve os registros internos do arquivo, porque ele não estava disponível na escrita.
- **FR-010**: O trabalho desta feature MUST permanecer na base local, no apontamento mockado, na captura e na preparação do log. Partes do módulo não relacionadas a isso MUST NOT ser alteradas. Dependência nova MUST ser somente a necessária a esta captura.

### Key Entities

- **Apontamento mockado**: Indica um único log a coletar no OpenObserve. Não é o conteúdo do log. Vários logs exigem vários apontamentos. Organização, stream, índice, campos e credenciais não fazem parte desta spec.
- **Log armazenado**: Registro já presente no OpenObserve e selecionado por um apontamento. No aceite, esse registro vem do arquivo `otel-logs.zip`. Os campos internos desse arquivo não fazem parte desta spec. O arquivo precisa ser entregue antes do planejamento.
- **Log pronto**: Registro armazenado capturado por inteiro, igual ao que foi armazenado, associado ao apontamento que o selecionou, aguardando a etapa futura de processamento. Não é ocorrência, não tem categoria e não tem dados sensíveis ocultos. Enquanto a instância está no ar, o mesmo apontamento devolve esse log. O lugar não tem nome. Depois que a instância para, esta etapa não exige que ele continue.
- **Falha de captura**: Resultado em que o log não fica pronto. Contém o apontamento usado e um motivo distinguível: apontamento inutilizável, log ausente ou OpenObserve inacessível.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Com um registro do arquivo `otel-logs.zip` armazenado no OpenObserve e um mock que aponta só para ele, ao fim da captura com a instância no ar, uma nova leitura pelo mesmo apontamento devolve esse registro pronto, igual ao armazenado. A contagem de campos ocultos por esta etapa é zero. A contagem de ocorrências classificadas é zero. A fila não é lida. Depois que a instância para, a ausência do log pronto não reprova esta etapa. Sem a entrega do arquivo, este critério não pode ser executado.
- **SC-002**: Em um conjunto com uma falha de apontamento inutilizável, uma de log ausente e uma de OpenObserve inacessível, nenhum desses logs fica pronto. Em cada falha o observador identifica o apontamento e qual dos três motivos ocorreu. Falha sem esse motivo não entra na contagem de atendidas.
- **SC-003**: Duas capturas com apontamentos distintos, executadas ao mesmo tempo, produzem dois resultados distinguíveis. Se uma falha e a outra conclui, os dois resultados permanecem visíveis ao fim. Nenhuma meta numérica de vazão é definida.
- **SC-006**: Em uma execução com duas indicações do mesmo apontamento para um log armazenado, a contagem de logs prontos desse apontamento é um. A leitura pelo apontamento, com a instância no ar, devolve esse log.
- **SC-004**: Com a base local desta etapa no ar, um log pode ser armazenado e lido no OpenObserve, e o Elasticsearch não faz parte dessa base.
- **SC-005**: Quando o mock não fornece apontamento, a contagem de logs prontos é zero e a contagem de falhas de captura é zero.

## Assumptions

Decisões e fatos abaixo têm fonte no pedido, na constituição ou em arquivo lido. O que o pedido não definiu fica fora, sem regra criada no lugar.

- Fonte: pedido, etapas 2 e 3. A fila está vazia, o uso dela não foi desenvolvido para esta etapa e não entra. O apontamento vem de uma fonte mockada. Isso não emenda a constituição. O fluxo normativo continua sendo fila, depois OpenObserve, depois agregação. Esta etapa cobre o apontamento mockado e a captura. A justificativa é o recorte explícito do pedido: a fila não será usada agora, e o processamento do log não será desenvolvido.
- Fonte: constituição, princípio IV, e spec `001-consumir-fila-redis`. O código de apontamento não é o conteúdo do log. Esta etapa reutiliza esse papel para o apontamento mockado. Os nomes de stream, grupo e campos que a spec 001 deixou indefinidos continuam indefinidos aqui.
- Fonte: pedido, etapa 4, e constituição, princípios I e III. "Deixar em estado de processamento" significa o log capturado, inteiro, associado ao apontamento, pronto para a etapa futura. Agregação e categoria de ocorrência ficam de fora. A constituição não define persistência da ocorrência.
- Fonte: resposta do solicitante em 2026-10-31, opção A. O log pronto é o registro armazenado, inteiro, igual ao registro do arquivo. Ocultar dados sensíveis fica na etapa de processamento. O README descreve mascaramento no módulo. Esta etapa não o aplica. Isso não emenda a constituição nem o README.
- Fonte: resposta do solicitante em 2026-10-31, opção A. Enquanto a instância está no ar, o log pronto pode ser lido de novo pelo mesmo apontamento. Depois que ela para, esta etapa não exige que ele continue. O lugar não é nomeado.
- Fonte: resposta do solicitante em 2026-10-31, opção A. Cada indicação seleciona um único log. Vários logs exigem várias indicações, que podem ser capturadas ao mesmo tempo. Uma indicação não traz o exemplo inteiro nem mais de um log.
- Fonte: pedido, "falhas deve apresentar detalhes completos". Os motivos distinguíveis desta etapa são apontamento inutilizável, log ausente e OpenObserve inacessível. Nova tentativa, prazo e fila de erro não foram informados e não foram criados.
- Fonte: resposta do solicitante em 2026-10-31, opção A. Indicações repetidas do mesmo apontamento na mesma execução convergem em um único log pronto. A leitura por esse apontamento, com a instância no ar, devolve esse log.
- Fonte: pedido, etapa 1, `docker-compose.yml`, `Makefile` e README, lidos em 2026-10-31. A base local sobe MongoDB, Redis e Elasticsearch. Esta etapa retira o Elasticsearch dessa base e inclui o OpenObserve. MongoDB e Redis permanecem. Versão, porta e credencial do OpenObserve não foram informadas e não foram criadas. Migração de dado já existente no Elasticsearch não foi pedida.
- Fonte: resposta do solicitante em 2026-10-31, opção B. O aceite só vale se o log capturado for um registro do `otel-logs.zip`. O arquivo precisa ser entregue antes do planejamento. Busca em 2026-10-31 no repositório e na pasta LhamaLog Workspace não encontrou esse arquivo. O Finder captura o registro no OpenObserve. O arquivo não é a fonte da captura. Os registros internos não são descritos aqui.
- Fonte: pedido e constituição, princípio V. Código, teste e comando desta feature não foram executados durante a escrita desta spec. Nenhum resultado de execução é afirmado.
- Fonte: pedido. O trabalho fica na base local, no apontamento mockado, na captura e na preparação. Health, consumo de fila já especificado na feature 001 e o placeholder de Elasticsearch no módulo ficam fora, salvo o necessário para a base deixar de incluir o Elasticsearch. A feature 001 não é reescrita.
- Fonte: pedido, README e constituição (Restrições Técnicas), lidos em 2026-10-31. O módulo é descrito como Go, arquitetura hexagonal leve. O pedido pede código em Go, clean code e a estrutura atual de pastas, sem biblioteca descontinuada e sem dependência além da necessária a esta captura. Esses fatos orientam o planejamento. Não são desenho desta spec. Dependências listadas no README como candidatas e ainda não adicionadas não são tratadas como já presentes.
- Fonte: pedido, grafia "Open Observer" e "open observer". A constituição e o repositório nomeiam a mesma ferramenta como OpenObserve. Esta spec usa OpenObserve.
