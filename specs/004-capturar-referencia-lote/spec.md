# Feature Specification: Capturar referência de lote

**Feature Branch**: `004-capturar-referencia-lote`

**Created**: 2026-11-08

**Status**: Draft

**Input**: User description: "O Finder deve consumir os eventos log.ingested v1 publicados no Redis, interpretar a referência de cada lote e encaminhá-la para a recuperação dos logs no OpenObserve. O Collector já persiste os registros válidos antes de publicar. Cada evento é um lote, não um log. A leitura direta do Redis é o escopo desta tarefa. Traces, a recuperação do conteúdo, a varredura, a publicação de checked, a Ocorrência e mudanças no Collector ficam de fora."

## Clarifications

### Session 2026-11-08

- Q: Quando a entrada Redis de um lote válido é confirmada? → A: No modelo síncrono. A confirmação espera a recuperação e o efeito que a etapa seguinte define como concluído. Interpretar a referência não confirma a entrada. A forma da busca permanece na etapa de recuperação.
- Q: O que acontece com uma referência inválida? → A: A entrada é confirmada, o motivo fica registrado e não há nova tentativa. Ela não permanece pendente e não segue como lote consultável.
- Q: Como a reserva e a conclusão da idempotência se separam? → A: Uma chave por `idempotencyId`, com estado em processamento ou concluído. Só o estado concluído dura 25 horas e impede repetir o efeito. O estado em processamento pode ser reassumido pelos critérios de pendência, sem esperar 25 horas.
- Q: O que a segunda entrada válida faz enquanto a chave daquele `idempotencyId` está em processamento e os critérios de pendência ainda não autorizam reassumir o lote? → A: Ela permanece pendente. Não inicia outro efeito e não é confirmada até a chave sair de em processamento.
- Q: Na primeira criação do grupo do Finder, a leitura ignora as entradas que já estão no stream ou inclui as que ainda estão lá? → A: Inclui as que ainda estão lá. O grupo começa no início do que o stream ainda guarda. A posição `$` não é usada.
- Q: Se a quantidade de entradas por leitura, o tempo de espera ou os critérios para reassumir uma pendência não forem informados, o Finder recusa o consumo ou pode começar mesmo assim? → A: Recusa o consumo. Nenhum valor padrão é criado.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Interpretar um lote válido e deixá-lo pronto para a etapa seguinte (Priority: P1)

O Finder lê uma entrada ainda não confirmada do stream de logs. A entrada é uma referência de lote, não o conteúdo dos logs. O Finder confere o envelope, separa metadata e body, e entrega à etapa de recuperação uma referência validada com a identidade do tenant vinda do metadata, o lugar do lote no OpenObserve, a faixa temporal, a quantidade de registros aceitos e o identificador da entrada. A interpretação não consulta o OpenObserve. A entrada só é confirmada depois que essa etapa termina o efeito que ela define como concluído.

**Why this priority**: Sem a referência validada, a etapa seguinte não sabe qual lote recuperar. É o resultado mínimo pedido.

**Independent Test**: Publicar uma entrada sintética de lote válido, com e sem traceId, em lote de um registro e em lote de vários registros. Conferir que a referência entregue traz os campos de consulta, que o tenant é o do metadata, que a interpretação não buscou conteúdo de log e que a entrada ainda não está confirmada.

**Acceptance Scenarios**:

1. **Given** uma entrada cujo envelope é `log.ingested`, origem `collector` e versão `"1"`, com metadata e body em JSON válido, signal `LOGS`, productId igual nos dois lados, `idempotencyId` igual ao `batchId`, faixa temporal com mínimo menor ou igual ao máximo e `recordCount` positivo, **When** o Finder a interpreta, **Then** a etapa de recuperação recebe tenantId e productId do metadata, organização, stream, batchId, a faixa em microssegundos, recordCount, idempotencyId, a correlação presente e o identificador da entrada Redis. A interpretação não consulta o OpenObserve, e a entrada segue sem confirmação.
2. **Given** a mesma entrada válida com `recordCount` igual a 1 e `recordId` preenchido, **When** o Finder a interpreta, **Then** a referência entregue inclui esse `recordId`.
3. **Given** a mesma entrada válida com `recordCount` maior que 1 e `recordId` nulo, **When** o Finder a interpreta, **Then** a referência entregue segue sem `recordId` e continua válida.
4. **Given** uma entrada válida sem `traceId` no metadata, **When** o Finder a interpreta, **Then** a ausência não a torna inválida e a referência segue para a etapa seguinte.
5. **Given** uma entrada válida com `traceId` no metadata, **When** o Finder a interpreta, **Then** esse `traceId` acompanha a referência e não é usado para deduplicar.

---

### User Story 2 - Impedir que referência inválida siga como lote consultável (Priority: P1)

Uma entrada que não pode ser consultada no OpenObserve não entra no conjunto de referências válidas. O motivo fica registrado sem credencial e sem conteúdo de log. A entrada é confirmada em seguida e não é tentada de novo. Ela não permanece pendente e não espera a recuperação.

**Why this priority**: Encaminhar uma referência incompleta ou incompatível faria a etapa seguinte buscar um lote que o contrato não sustenta. Deixá-la pendente repetiria a mesma rejeição.

**Independent Test**: Submeter, uma a uma, as entradas inválidas listadas nos cenários e conferir que nenhuma aparece como referência válida, que o motivo não contém credencial nem conteúdo de log, que cada entrada foi confirmada e que nenhuma delas continua pendente.

**Acceptance Scenarios**:

1. **Given** metadata ou body não é JSON válido, ou falta tenant, produto, organização, stream, batchId, faixa temporal ou recordCount, **When** o Finder avalia a entrada, **Then** ela não segue como referência válida e o motivo fica registrado sem credencial e sem conteúdo de log.
2. **Given** o evento não é `log.ingested`, a origem não é `collector`, a versão não é a string `"1"`, ou o signal não é `LOGS`, **When** o Finder avalia a entrada, **Then** ela não segue como referência válida. Menção a schema 2 não conta como versão aceita.
3. **Given** o productId do metadata difere do productId do body, o `idempotencyId` difere do `batchId`, ou o mínimo da faixa é maior que o máximo, **When** o Finder avalia a entrada, **Then** ela não segue como referência válida.
4. **Given** `recordCount` é zero ou negativo, **When** o Finder avalia a entrada, **Then** ela não segue como referência válida.
5. **Given** lote de um registro sem `recordId` preenchido, ou lote de vários registros com `recordId` preenchido, **When** o Finder avalia a entrada, **Then** ela não segue como referência válida.
6. **Given** `publishedAt` não é um instante ISO em UTC, **When** o Finder avalia a entrada, **Then** ela não segue como referência válida.
7. **Given** qualquer entrada inválida dos cenários anteriores, **When** o Finder termina a avaliação, **Then** a entrada é confirmada, não permanece pendente, não é tentada de novo e não recebe o estado concluído.

---

### User Story 3 - Não perder o lote nem repetir o efeito concluído (Priority: P2)

O mesmo `idempotencyId` representa o mesmo lote e ocupa uma única chave do Finder, em processamento ou concluída. O estado em processamento é a reserva: não dura 25 horas e não conta como conclusão. Ele pode ser reassumido quando a entrada atende aos critérios de pendência informados para a execução. O estado concluído é gravado quando a recuperação termina o efeito que ela define como concluído, dura 25 horas e, nesse prazo, impede repetir o efeito. A confirmação da entrada válida usa o identificador Redis e só ocorre depois desse efeito. Esta feature não promete entrega exatamente uma vez.

**Why this priority**: O contrato entrega pelo menos uma vez. Tratar a reserva como conclusão descarta um lote cujo efeito não aconteceu.

**Independent Test**: Publicar duas entradas Redis com o mesmo `idempotencyId`. Enquanto a primeira mantém a chave em processamento e os critérios de pendência do teste ainda não autorizam reassunção, conferir que a segunda permanece pendente, sem outro efeito e sem confirmação. Informar o efeito concluído da primeira e conferir que a segunda, ao sair dessa espera, não repete o efeito e é confirmada. Numa execução separada, interromper com a chave em processamento e conferir que a entrada pode ser reassumida pelos critérios de pendência do teste, sem esperar 25 horas. O teste observa o efeito informado pela etapa de recuperação e não descreve a consulta ao OpenObserve.

**Acceptance Scenarios**:

1. **Given** duas entradas Redis diferentes com o mesmo `idempotencyId` e a primeira já está concluída dentro das 25 horas, **When** a segunda é processada, **Then** o efeito da recuperação não se repete e a nova entrada é confirmada.
2. **Given** a chave está em processamento e o processo cai antes do efeito da recuperação, **When** a entrada atende aos critérios de pendência informados para a execução, **Then** essa reserva não conta como conclusão, o estado pode ser reassumido sem esperar 25 horas e o efeito ainda pode ser realizado.
3. **Given** o efeito da recuperação já ocorreu e a confirmação da entrada ainda não, **When** o processo cai e a entrada é reentregue, **Then** o trabalho concluído não é perdido, o efeito não é repetido e a entrada é confirmada.
4. **Given** a marca de conclusão de um `idempotencyId` completou 25 horas, **When** uma nova entrada com esse identificador chega, **Then** essa marca não a suprime mais.
5. **Given** o Finder e a Varredura recebem o mesmo `idempotencyId`, **When** cada um conclui o próprio processamento, **Then** a conclusão de um não impede a do outro.
6. **Given** a chave está em processamento, os critérios de pendência ainda não autorizam reassunção e existe uma segunda entrada válida com o mesmo `idempotencyId`, **When** essa segunda entrada é lida, **Then** ela permanece pendente, não inicia outro efeito e não é confirmada. Quando a chave passa a concluída, essa entrada é confirmada sem repetir o efeito.

---

### User Story 4 - Consumir no grupo do Finder e recuperar pendência (Priority: P2)

O Finder usa um grupo de consumo próprio, separado do grupo da Varredura. Os dois serviços podem observar as mesmas entradas. Instâncias do Finder que compartilham o grupo do Finder dividem as entradas entre si. Na criação, o grupo começa no início do que ainda está no stream, então as entradas já presentes são lidas. A posição `$` não é usada. Entradas abandonadas voltam ao fluxo pelos critérios de pendência fornecidos para aquela execução. Redis indisponível não confirma entrada inacabada.

**Why this priority**: Compartilhar o grupo com a Varredura dividiria os lotes entre os serviços. Perder uma pendência ou confirmá-la cedo demais descarta a referência.

**Independent Test**: Publicar entradas antes de criar o grupo do Finder e conferir que elas são lidas. Publicar outra entrada com o grupo do Finder e um grupo distinto no papel da Varredura e conferir que os dois a obtêm. Interromper uma instância do Finder antes do efeito da recuperação e conferir que a entrada permanece pendente. Aplicar critérios de abandono informados para o teste e conferir que a entrada abandonada pode ser reivindicada e não é descartada.

**Acceptance Scenarios**:

1. **Given** o grupo do Finder e o grupo da Varredura estão ativos no mesmo stream, **When** uma entrada é publicada, **Then** cada grupo pode receber essa entrada.
2. **Given** duas instâncias do Finder usam o grupo do Finder, **When** há várias entradas livres, **Then** cada entrada é entregue a uma instância desse grupo por vez.
3. **Given** o stream já contém entradas antes da criação do grupo do Finder, **When** o grupo passa a existir, **Then** essas entradas são entregues ao grupo. A posição `$` não é usada.
4. **Given** o Redis fica indisponível no meio do processamento, **When** a tentativa termina, **Then** nenhuma entrada ainda não concluída é confirmada.
5. **Given** uma entrada pendente ultrapassa os critérios de abandono informados para a execução, **When** a recuperação de pendentes ocorre, **Then** ela pode ser reivindicada e não é descartada por esta feature.
6. **Given** a recuperação já concluiu o efeito e a entrada válida foi confirmada, **When** uma falha acontece depois disso, **Then** essa falha pertence à etapa que ficou com o resultado, e esta feature não a trata como entrada ainda não confirmada.
7. **Given** falta a quantidade de entradas por leitura, o tempo de espera ou os critérios para reassumir uma pendência, **When** o Finder é iniciado, **Then** o consumo não começa, nenhuma entrada é lida e nenhum valor padrão é criado.

---

### Edge Cases

- Stream sem entrada nova: nenhuma referência é entregue e nenhuma entrada é confirmada.
- Redis indisponível: o estado já pendente permanece, e entrada não concluída não é confirmada.
- JSON inválido, campo necessário ausente, evento, origem, versão ou signal incompatível, divergência de produto, de batchId ou de faixa, recordCount não positivo, recordId incoerente com lote de um ou de vários registros: a entrada não segue como referência válida. O motivo fica registrado sem credencial e sem conteúdo de log. A entrada é confirmada, não é tentada de novo e não permanece pendente. Ela não recebe o estado concluído. Se a chave estava em processamento, ela é removida junto com a confirmação.
- Versão apresentada como schema 2: não é versão válida deste contrato.
- `traceId` ausente: o lote continua válido. `correlationId` e o identificador Redis não deduplicam.
- Atributos que o cliente teria enviado no log não substituem tenantId nem productId do metadata. O conteúdo do log não está na entrada.
- Lote misto no Collector: recordCount, faixa, bytes e agregados desta referência descrevem só os registros aceitos. O Finder não reconstrói o lote bruto recusado.
- Segunda entrada válida com a chave em processamento e sem autorização de reassunção: ela permanece pendente, não inicia outro efeito e não é confirmada. Quando a chave passa a concluída, essa entrada é confirmada sem repetir o efeito.
- Queda com a chave em processamento, antes do efeito da recuperação: o efeito não está concluído. A entrada continua pendente e pode ser reassumida pelos critérios de pendência informados para a execução, sem esperar 25 horas. Outra entrada válida do mesmo `idempotencyId` também permanece pendente até essa reassunção ou até a conclusão.
- Queda depois do efeito da recuperação e antes da confirmação: a chave já está concluída, o efeito não se repete e a reentrega confirma a entrada.
- Marca de conclusão com 25 horas completas: ela deixa de suprimir uma nova entrada com o mesmo `idempotencyId`.
- Entrada já removida pelo corte aproximado de 100000 entradas: esta feature não a recria. Esse corte não garante retenção de 25 horas.
- Evento `trace.ingested` e o stream de traces: fora desta feature. Se uma entrada no stream de logs não for `log.ingested`, ela cai na validação de evento.
- Credencial ou conteúdo de log em registro operacional: não atende esta feature.
- Quantidade de entradas por leitura, tempo de espera ou critérios para reassumir pendência ausentes: o consumo não começa e nenhum valor padrão é criado. Com os três informados, a execução usa esses valores.
- Entradas já presentes na criação do grupo: são entregues ao Finder. A posição `$` não é usada. Entrada já removida pelo corte aproximado continua irrecuperável por esta feature.
- Retry, spool ou DLQ do Collector: não existem nesta entrega e não são assumidos.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: O Finder MUST ler entradas do stream `lhamalog:stream:logs` por um grupo de consumo próprio, preservando o identificador de cada entrada Redis.
- **FR-002**: Conexão, autenticação, nome do grupo e identificação do consumidor MUST vir de configuração. Credencial MUST NOT aparecer no código nem em registro operacional. Registro operacional MUST NOT incluir conteúdo de log.
- **FR-003**: O grupo do Finder MUST ser distinto do grupo da Varredura, de modo que os dois serviços possam receber as mesmas entradas. Instâncias do Finder que usam o grupo do Finder MUST dividir as entradas entre si, cada entrada com no máximo uma entrega ativa nesse grupo.
- **FR-004**: Na criação do grupo do Finder, a posição inicial MUST ser o início do que ainda está no stream. Entradas já presentes MUST ser entregues ao grupo. A posição `$` MUST NOT ser usada.
- **FR-005**: A quantidade de entradas por leitura, o tempo de espera de uma leitura e os critérios para reivindicar entrada pendente abandonada MUST ser informados antes do consumo. Se qualquer um deles estiver ausente, o Finder MUST NOT iniciar o consumo e MUST NOT criar um valor padrão. Quando os três estão informados, a leitura e a recuperação de pendentes MUST usar esses valores. Esta spec MUST NOT atribuir esses números. A recuperação MUST NOT descartar a entrada só porque um número não estava no contrato.
- **FR-006**: Todos os campos do envelope MUST ser lidos como texto. `metadata` e `body` MUST ser interpretados como JSON separados. O body MUST ser tratado como referência e agregados do lote, não como o conteúdo dos logs.
- **FR-007**: O Finder MUST aceitar o envelope somente quando `event` é `log.ingested`, `origin` é `collector`, `schemaVersion` é a string `"1"` e `publishedAt` é um instante ISO em UTC. Schema 2, ou qualquer outra versão, MUST NOT ser aceita como este contrato.
- **FR-008**: Da metadata, o Finder MUST extrair `tenantId`, `productId` e `correlationId`. `traceId` pode estar ausente. A ausência de `traceId` MUST NOT invalidar o lote.
- **FR-009**: Do body, o Finder MUST extrair e preservar, quando presentes e legíveis, `receivedAt`, `signal`, `productId`, `openobserve.organization`, `openobserve.stream`, `openobserve.batchId`, `openobserve.recordId`, `openobserve.minTimestampMicros`, `openobserve.maxTimestampMicros`, `recordCount`, `bytes`, `serviceNames`, `priority`, `criticality`, `maxSeverityNumber` e `severityCounts`.
- **FR-010**: Para a referência ser entregue como válida, metadata e body MUST ser JSON válido e MUST atender a todas as condições a seguir: signal igual a `LOGS`; presença de tenant, produto, organização, stream, batchId, mínimo e máximo da faixa temporal e recordCount; productId do metadata igual ao productId do body; `idempotencyId` igual ao `batchId`; mínimo da faixa menor ou igual ao máximo; recordCount inteiro positivo; em lote de um registro, `recordId` preenchido; em lote de vários registros, `recordId` nulo.
- **FR-011**: `tenantId` e `productId` usados na referência entregue MUST ser os do metadata, originados do token validado pelo Collector. O Finder MUST NOT substituí-los por atributos enviados pelo cliente. `correlationId` e o identificador da entrada Redis MUST NOT ser a chave de idempotência.
- **FR-012**: A referência entregue à etapa de recuperação MUST conter tenantId, productId, organização, stream, batchId, recordId quando houver, mínimo e máximo da faixa em microssegundos, recordCount, idempotencyId e a correlação presente. O identificador da entrada Redis MUST permanecer associado a esse processamento para a confirmação correta. A interpretação MUST NOT consultar o OpenObserve. A forma da busca do conteúdo MUST permanecer na etapa de recuperação.
- **FR-013**: recordCount, bytes e os agregados MUST ser entendidos como referentes aos registros aceitos pelo Collector. O Finder MUST NOT tratá-los como a contagem do lote bruto enviado pelo cliente.
- **FR-014**: A confirmação de uma entrada válida MUST usar o identificador da entrada Redis, não o `idempotencyId`, e MUST ocorrer somente depois que a etapa de recuperação terminou o efeito que ela define como concluído. Interpretar a referência e entregá-la MUST NOT confirmar a entrada. Esta feature MUST NOT criar fila, spool ou outro armazenamento intermediário para essa espera. A forma da busca permanece na etapa de recuperação.
- **FR-015**: Referência inválida MUST NOT ser entregue como válida. O motivo MUST ser registrado sem credencial e sem conteúdo de log. A entrada MUST ser confirmada, MUST NOT ser tentada de novo e MUST NOT permanecer pendente. Essa confirmação MUST NOT gravar o estado concluído. Se a chave desse `idempotencyId` estiver em processamento, ela MUST ser removida junto com a confirmação, para deixar de estar em processamento. Esta feature MUST NOT criar retry, spool ou DLQ para a entrada inválida.
- **FR-016**: A deduplicação MUST usar uma única chave por `idempotencyId`, no namespace exclusivo do Finder, separado do da Varredura. A chave MUST estar em processamento ou concluída. A chave MUST ser criada somente se ainda não existir. O estado em processamento é a reserva: ele MUST NOT contar como conclusão e MUST NOT durar 25 horas como bloqueio do efeito. Quando a entrada pendente atende aos critérios de FR-005, esse estado MUST poder ser reassumido sem esperar 25 horas, e o efeito MUST poder ser tentado. Enquanto a chave está em processamento e esses critérios ainda não autorizam a reassunção, outra entrada válida com o mesmo `idempotencyId` MUST permanecer pendente, MUST NOT iniciar outro efeito e MUST NOT ser confirmada. Quando a chave passa a concluída, essa entrada MUST ser confirmada sem repetir o efeito. O estado concluído MUST ser gravado quando o efeito da recuperação estiver terminado, MUST durar 25 horas e, nesse prazo, MUST impedir repetir o efeito. Uma entrada nova com `idempotencyId` já concluído MUST ser confirmada sem repetir o efeito. Esta feature MUST NOT prometer entrega exatamente uma vez.
- **FR-017**: Com o Redis indisponível, o Finder MUST preservar o estado já pendente e MUST NOT confirmar entrada não concluída.
- **FR-018**: Falha de uma referência válida antes do efeito da recuperação MUST manter a entrada pendente e a chave em processamento. Entrada pendente abandonada MUST poder ser recuperada pelos critérios de FR-005, sem descarte indevido. Este requisito não se aplica à entrada inválida já encerrada por FR-015.
- **FR-019**: Depois que a entrada válida foi confirmada, falha posterior MUST pertencer à etapa que ficou com o resultado da recuperação.
- **FR-020**: O corte aproximado de 100000 entradas MUST ser entendido só como limite aproximado de quantidade. Ele MUST NOT ser tratado como retenção de 25 horas. Referência já removida por esse corte MUST NOT ser recriada por esta feature.
- **FR-021**: Esta feature MUST NOT definir a busca do conteúdo no OpenObserve, MUST NOT varrer dado sensível, MUST NOT publicar checked, MUST NOT integrar com Ocorrência, MUST NOT consumir traces e MUST NOT alterar o Collector. Ela MUST esperar o efeito concluído informado pela etapa de recuperação antes de confirmar uma entrada válida.
- **FR-022**: Enquanto o estado concluído de um `idempotencyId` não completar 25 horas, nova entrada com esse identificador MUST NOT repetir o efeito. Depois desse prazo, esse estado MUST NOT continuar suprimindo a entrada.

### Key Entities

- **Entrada Redis**: Uma publicação no stream `lhamalog:stream:logs`. Tem identificador próprio da entrada e o envelope em campos de texto. Não contém o conteúdo dos logs. O identificador da entrada confirma a entrada; não deduplica o lote.
- **Envelope**: `idempotencyId`, `event`, `origin`, `schemaVersion`, `publishedAt`, `metadata` e `body`. Os dois últimos são JSON em texto.
- **Metadata**: Identidade do lote. `tenantId` e `productId` vêm do token validado pelo Collector. `correlationId` correlaciona a exportação e não deduplica. `traceId` é condicional.
- **Referência de lote**: Conteúdo do body. Aponta organização, stream, batchId, recordId quando o lote tem um registro, faixa temporal em microssegundos, recordCount, bytes e agregados dos registros aceitos.
- **Referência validada**: Recorte entregue à etapa de recuperação, com a identidade do metadata, o apontamento no OpenObserve, a correlação presente e o identificador da entrada Redis. Não é o log. A entrega não confirma a entrada.
- **Chave do lote**: Uma por `idempotencyId`, exclusiva do Finder. Está em processamento ou concluída. Em processamento é a reserva e pode ser reassumida pelos critérios de pendência. Enquanto essa reassunção não é autorizada, outra entrada válida do mesmo identificador permanece pendente. Concluída significa que o efeito da recuperação terminou, dura 25 horas e impede repetir esse efeito. A remoção na entrada inválida tira a chave de em processamento sem gravar concluída.
- **Grupo do Finder**: Grupo de consumo exclusivo deste serviço. Independente do grupo da Varredura.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Em um conjunto de lotes sintéticos válidos, cada referência entregue contém o tenant e o produto do metadata, organização, stream, batchId, a faixa temporal, recordCount e idempotencyId, com o identificador da entrada associado. Ao fim da interpretação, a contagem de entradas já confirmadas é zero e a contagem de conteúdos de log entregues é zero.
- **SC-002**: Em lote de um registro, a referência entregue inclui `recordId`. Em lote de vários registros, a referência entregue não inclui `recordId`. Nos dois casos a referência permanece válida quando o restante do contrato é atendido.
- **SC-003**: Para cada caso inválido descrito nos cenários da história 2, a contagem de referências válidas produzidas por essa entrada é zero, a entrada é confirmada uma vez e a contagem de pendências restantes dela é zero. O registro do motivo não contém credencial nem conteúdo de log. A contagem de estados concluídos gravados por essa rejeição é zero.
- **SC-004**: Duas entradas Redis com o mesmo `idempotencyId`, processadas enquanto o estado concluído está dentro de 25 horas, produzem um efeito da recuperação. A segunda entrada é confirmada sem repetir esse efeito. Enquanto a chave está em processamento e os critérios de pendência do teste não autorizam reassunção, a contagem de efeitos iniciados para esse identificador não sobe por causa da segunda entrada e a contagem de confirmações dela é zero. Depois de 25 horas, esse estado não entra na contagem de supressões. Uma queda com a chave em processamento deixa a contagem de efeitos concluídos inalterada, e a entrada pode ser reassumida pelos critérios de pendência do teste sem esperar 25 horas.
- **SC-005**: Com o grupo do Finder e o grupo da Varredura ativos, uma entrada publicada é observável nos dois grupos. A contagem de entradas perdidas para o outro serviço por compartilharem o mesmo grupo é zero. Entradas já presentes no stream antes da criação do grupo do Finder são todas entregues a esse grupo. A contagem das que deixam de ser lidas por uso da posição `$` é zero.
- **SC-006**: Quando o Redis fica indisponível, a contagem de confirmações de entradas não concluídas é zero. Referência válida que falha antes do efeito da recuperação continua pendente. Entrada inválida já confirmada não volta a essa contagem.
- **SC-007**: Entrada pendente que atende aos critérios de abandono fornecidos para a execução volta a poder ser processada. A contagem de pendências descartadas por esta recuperação é zero.
- **SC-008**: Em uma execução com entradas válidas e inválidas, a contagem de credenciais e de conteúdos de log nos registros operacionais é zero. O tenant de cada referência válida é o tenant do metadata em 100% dos casos.
- **SC-009**: Quando falta a quantidade por leitura, o tempo de espera ou os critérios de pendência, a contagem de entradas lidas é zero e a contagem de confirmações é zero. Nenhum valor padrão entra na contagem de parâmetros usados.

## Assumptions

Decisões e fatos abaixo têm fonte no pedido, no contrato citado, na constituição ou em spec já escrita. O que essas fontes deixaram em aberto não recebeu valor novo.

- Fonte: pedido e contrato do Collector, seções 3, 4 e 6, lido em 2026-11-08. O Collector grava no OpenObserve os registros válidos e só depois publica uma entrada por lote. A entrada não traz o log. Agregados e recordCount descrevem o subconjunto aceito. Esta feature começa nessa entrada.
- Fonte: pedido e contrato, seção 4. O stream desta feature é `lhamalog:stream:logs`. O evento é `log.ingested`, a origem é `collector` e a versão vigente é a string `"1"`. Menções antigas a schema 2 não são outra versão válida. O stream de traces não entra.
- Fonte: contrato, seções 4 e 8, e o exemplo v1 da seção 10. `idempotencyId` é o `batchId` determinístico. A coerência exigida pelo pedido é a igualdade dos dois. `correlationId` não deduplica. O identificador Redis não substitui `idempotencyId`.
- Fonte: contrato, seção 5, e pedido. Lote unitário tem `recordId` preenchido. Lote múltiplo tem `recordId` nulo. Esta spec usa recordCount igual a 1 como lote unitário e recordCount maior que 1 como lote múltiplo, porque é o sentido desses nomes e do exemplo com recordCount 1000 e recordId nulo.
- Fonte: pedido. As validações de FR-010 são requisitos propostos para o Finder, derivados do fluxo de recuperação. O contrato não classifica formalmente cada campo como obrigatório. Campo fora dessa lista, quando presente e legível, é preservado. `traceId` pode faltar. Enumeração de prioridade e criticidade não é revalidada aqui.
- Fonte: pedido. `publishedAt` é interpretado como instante ISO UTC. Texto que não representa esse instante não forma envelope válido.
- Fonte: contrato, seção 8, e resposta do solicitante em 2026-11-08, Q3 B. A deduplicação do evento usa uma chave por `idempotencyId`, criada se ausente. O estado concluído dura 25 horas. O estado em processamento não usa esse prazo e pode ser reassumido pelos critérios de pendência. Isso não promete exactly-once e não substitui a deduplicação dos documentos por `lhamalog_record_id`, que pertence à etapa de recuperação. A garantia descrita no contrato é pelo menos uma vez.
- Fonte: resposta do solicitante em 2026-11-08, Q1 A. A entrada válida só é confirmada depois do efeito que a etapa de recuperação define como concluído. A interpretação não confirma a entrada e não cria armazenamento intermediário. A forma da busca não é definida aqui.
- Fonte: resposta do solicitante em 2026-11-08, Q2 A, e contrato, seção 11. A entrada inválida é confirmada, o motivo fica registrado e não há nova tentativa. Retry, spool e DLQ do Collector não são criados. A confirmação da inválida remove a chave em processamento e não grava o estado concluído. A spec 002, em 2026-11-05, já tirava da pendência o item fora do contrato antigo; esta sessão aplica esse desfecho ao contrato `log.ingested` v1.
- Fonte: resposta do solicitante em 2026-11-08, opção A. Uma segunda entrada válida, com a chave em processamento e sem autorização de reassunção, permanece pendente. Ela não inicia outro efeito e não é confirmada. Quando a chave passa a concluída, essa entrada é confirmada sem repetir o efeito.
- Fonte: pedido, decisão necessária 1, contrato, seção 8, e resposta do solicitante em 2026-11-08, opção B. Grupo próprio do Finder é requisito. O nome do grupo e o nome do consumidor não foram informados e não foram criados. A posição inicial é o início do que ainda está no stream. A posição `$` não é usada, porque ignoraria as entradas já presentes.
- Fonte: pedido, decisão necessária 5, e resposta do solicitante em 2026-11-08, opção A. Quantidade por leitura, tempo de espera e critérios de pendência não estão no contrato e não foram numerados aqui. Se algum deles falta, o consumo não começa e nenhum padrão é criado.
- Fonte: contrato, seções 4 e 9. O corte padrão aproxima 100000 entradas. Não garante 25 horas nem devolve referência já removida.
- Fonte: spec `001-consumir-fila-redis`. Ela deixou stream, grupo e campos sem nome e limitou o item ao código de apontamento. Esta feature nomeia stream e envelope porque o contrato do Collector agora os informa. O papel da entrada continua sendo apontar o lote, não carregar o log.
- Fonte: spec `002-corrigir-entrega-fila`. A ordenação da ADR 169 e os níveis Very High a Very Low não são requisito desta leitura. O pedido manda seguir a ordem e o contrato do stream `lhamalog:stream:logs`.
- Fonte: constituição, princípios I, II e IV, e pedido, fora de escopo. O fluxo normativo segue até a ocorrência, e esta feature cobre só a leitura e a interpretação da referência. A justificativa é o recorte explícito do pedido: recuperar o conteúdo, varrer, publicar checked e classificar ocorrência ficam para etapas posteriores.
- Fonte: constituição, princípio V. Código, teste e comando desta feature não foram executados durante a escrita desta spec. Nenhum resultado de execução é afirmado.
- Fonte: contrato, seção 1. O documento descreve o contrato que o consumidor precisa atender. Não é homologação de uma implementação já existente do Finder. Recomendações de implementação no pedido não foram tratadas como funcionalidade já presente.
- Fonte: guia histórico de Redis Streams, citado pelo contrato na seção 8, e resposta do solicitante em 2026-11-08, Q1 A. Esse guia usa outro envelope e outra organização de chaves. Políticas de retry e DLQ desse guia não foram copiadas. A confirmação da entrada válida ficou no efeito concluído da recuperação, não numa infraestrutura nova.

## Out of Scope

- A forma da busca do conteúdo dos logs no OpenObserve. Esta feature espera o efeito concluído informado por essa etapa e não descreve a consulta.
- Varredura de dados sensíveis.
- Publicação de checked.
- Integração com Ocorrência.
- Consumo de `trace.ingested` e do stream de traces.
- Alteração do Collector.
- Retry, spool, DLQ ou outro armazenamento intermediário.
- Deduplicação dos documentos por `lhamalog_record_id`.
