# Feature Specification: Recuperar logs do lote

**Feature Branch**: `005-recuperar-logs-lote`

**Created**: 2026-11-08

**Status**: Draft

**Input**: User description: "Recuperar no OpenObserve os logs completos indicados pela referência validada de log.ingested v1. A busca usa organização, stream, batchId e a faixa temporal da referência, percorre todas as páginas, deduplica por identificador lógico e só declara o lote completo quando a contagem distinta coincide com recordCount e as demais validações passam. O tenant da referência é a autoridade de isolamento. Falhas e divergências voltam classificadas ao fluxo da primeira tarefa, sem confirmação da entrada Redis por esta etapa. Spans, consumo Redis, ACK, varredura, checked, Ocorrência e alterações no Collector ficam de fora."

## Clarifications

### Session 2026-11-08

- Q: Quais resultados desta recuperação a primeira tarefa pode tratar como efeito concluído para confirmar a entrada? → A: O lote completo e as falhas que pedem correção. Indisponibilidade, timeout, falha no meio da paginação e lote incompleto continuam pendentes. As falhas que pedem correção são autenticação ou autorização, isolamento ou identidade, resposta inválida, excesso de registros distintos, conteúdo divergente e documento sem identificador. Esta etapa informa o efeito e não confirma a entrada.
- Q: Quando dois documentos têm o mesmo identificador de log e não são idênticos, o que conta como cópia física e o que conta como conteúdo divergente? → A: A lista dos campos que só identificam a cópia física é configuração obrigatória. Sem ela, a recuperação não começa e não marca efeito concluído. Com ela, diferença restrita a esses campos é uma cópia; qualquer outra diferença é conteúdo divergente, e nenhuma versão é escolhida.
- Q: Quando os nomes dos campos de tenant e produto já estão disponíveis, como a recuperação confirma que os documentos são desse tenant e desse produto? → A: A consulta filtra pelos campos confirmados de tenant e produto, e cada documento retornado é conferido com a referência. A credencial é necessária e não substitui esses dois controles.
- Q: Se o mesmo lote reúne mais de uma falha, qual classificação a recuperação devolve? → A: Uma só classificação, nesta ordem: autenticação ou autorização; resposta inválida; isolamento ou identidade; documento sem identificador; conteúdo divergente; excesso de registros distintos; lote incompleto. Indisponibilidade, timeout e falha de página ficam pendentes e só valem quando nenhuma falha que pede correção ocorreu.
- Q: O que a recuperação faz se as páginas continuam cheias e a leitura não encontra um fim? → A: Não há máximo. A leitura só para quando uma página volta com menos documentos do que o tamanho pedido.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Recuperar os logs do lote pela referência (Priority: P1)

O Finder recebe a referência já validada pela primeira tarefa e recupera, no armazenamento, os logs daquele lote. A referência não traz o conteúdo. A consulta usa a organização e o stream indicados nela, filtra o lote e cobre a faixa temporal, inclusive o último instante. Em lote de um registro, a consulta também restringe o identificador daquele registro. O resultado de sucesso devolve os registros distintos e os dados que ligam esse resultado ao lote. Esta etapa não confirma a entrada Redis.

**Why this priority**: Sem os logs do lote, a referência não chega a conteúdo recuperado. É o resultado mínimo pedido.

**Independent Test**: Entregar uma referência sintética válida, de lote com um registro e de lote com vários, contra um armazenamento de teste que já contém só os registros aceitos. Conferir organização, stream, lote, faixa e o filtro adicional do lote unitário. Conferir que o conteúdo não saiu da entrada Redis e que esta etapa não confirmou essa entrada.

**Acceptance Scenarios**:

1. **Given** uma referência validada de vários registros, com organização, stream, batchId, faixa e recordCount positivo, e sem recordId, **When** a recuperação consulta o armazenamento, **Then** a consulta usa essa organização e esse stream, filtra `lhamalog_batch_id` pelo batchId, começa em `minTimestampMicros` e termina em `maxTimestampMicros + 1`, e não aplica filtro por identificador de registro.
2. **Given** uma referência validada de um registro, com recordId preenchido, **When** a recuperação consulta o armazenamento, **Then** a consulta aplica também o filtro por `lhamalog_record_id` desse registro, além do lote e da faixa.
3. **Given** o stream da referência é diferente do stream Redis `lhamalog:stream:logs`, **When** a recuperação consulta o armazenamento, **Then** o lugar consultado é o stream da referência.
4. **Given** a consulta percorre o lote e as validações desta feature passam, **When** o resultado de sucesso é entregue, **Then** ele contém os registros distintos recuperados e a correlação do lote: tenantId, productId, organização, stream, batchId, idempotencyId, recordCount e a quantidade distinta encontrada. A entrada Redis segue sem confirmação desta etapa.
5. **Given** um valor da referência não pode ser aplicado só como dado do filtro pretendido, **When** a consulta seria montada, **Then** a classificação é resposta inválida, nenhum registro é entregue como sucesso e a consulta não troca o lote, a organização, o stream nem a faixa.

---

### User Story 2 - Percorrer todas as páginas e deduplicar cópias físicas (Priority: P1)

A recuperação pede o lote em páginas, a partir do início, com o tamanho de página configurado para a execução. A página seguinte avança por esse tamanho. A leitura continua até uma página trazer menos documentos do que o tamanho pedido. Se a página volta cheia, a seguinte é consultada para saber se o lote acabou. Não há número máximo de páginas. Cópias físicas do mesmo registro, na mesma página ou em páginas diferentes, viram um único registro lógico, identificado por `lhamalog_record_id`. A lista dos campos que só identificam a cópia física vem da configuração da execução. Sem essa lista, a recuperação não começa. Com ela, diferença só nesses campos é uma cópia; qualquer outra diferença é conteúdo divergente, e nenhuma versão é escolhida. A quantidade de documentos físicos pode ser maior que recordCount. Limitar a leitura a recordCount não é a forma de recuperar o lote.

**Why this priority**: Parar na contagem esperada, ou tratar uma cópia como outro registro, entrega um lote errado mesmo quando a consulta da história 1 está correta.

**Independent Test**: Fornecer a lista de campos que só identificam a cópia física. Montar um lote sintético com uma página, com várias páginas e com quantidade de documentos múltipla exata do tamanho de página. Incluir cópias do mesmo registro na mesma página e em páginas diferentes, com um registro lógico novo depois de uma sequência de cópias. Conferir que cada identificador lógico aparece uma vez e que o registro que estaria de fora, se a leitura parasse em recordCount, está no resultado. Repetir sem a lista e conferir que a recuperação não começa.

**Acceptance Scenarios**:

1. **Given** o lote cabe em uma página, **When** a recuperação lê, **Then** ela começa no deslocamento zero, usa o tamanho de página configurado e para quando a página devolve menos documentos do que esse tamanho.
2. **Given** o lote ocupa mais de uma página, **When** a recuperação lê, **Then** cada página seguinte aumenta o deslocamento pelo tamanho de página, e a leitura só termina numa página incompleta.
3. **Given** a quantidade de documentos físicos é múltipla exata do tamanho de página, **When** a última página cheia é lida, **Then** a página seguinte é consultada antes de a recuperação tratar o lote como encerrado.
4. **Given** o mesmo `lhamalog_record_id` aparece mais de uma vez na mesma página e em páginas diferentes, com o mesmo conteúdo de log e com diferença só em campos da lista configurada, **When** a deduplicação termina, **Then** esse identificador entra uma única vez no resultado lógico. O registro entregue é uma das cópias. O identificador Redis e o correlationId não são a chave dessa deduplicação.
5. **Given** há mais documentos físicos do que recordCount e um identificador lógico novo aparece depois de cópias já suficientes para atingir recordCount, **When** a recuperação termina, **Then** esse identificador novo entra no conjunto distinto. A leitura não usa recordCount como tamanho único da consulta.
6. **Given** um documento do lote não traz `lhamalog_record_id`, **When** a deduplicação o encontra, **Then** ele não entra na contagem de registros lógicos válidos e o lote não é declarado completo. A falha sinaliza documento sem o identificador necessário.
7. **Given** dois documentos com o mesmo `lhamalog_record_id` diferem em campo que não está na lista configurada, **When** a deduplicação os encontra, **Then** nenhuma das versões é escolhida como o registro lógico e o lote não é declarado completo. A falha é conteúdo divergente, distinta de lote incompleto e de excesso de identificadores.
8. **Given** a lista de campos que só identificam a cópia física não foi fornecida, **When** a recuperação é pedida, **Then** ela não começa, nenhum valor padrão é criado e o resultado não marca efeito concluído. Uma lista fornecida vazia permanece válida: qualquer diferença conta como conteúdo divergente.

---

### User Story 3 - Declarar o lote completo só quando a contagem fecha (Priority: P1)

Ao fim da paginação, a quantidade de identificadores distintos válidos é comparada com recordCount. Igualdade só declara o lote completo se as outras validações também passaram. Quantidade menor é lote incompleto e informa o esperado e o encontrado. Quantidade maior é excesso. Uma busca vazia, com recordCount positivo, não é sucesso. A contagem não afirma, sozinha, se a divergência veio de atraso de visibilidade ou de falha de persistência.

**Why this priority**: Devolver um conjunto sem essa conferência faria a etapa seguinte tratar um lote parcial ou inflado como o lote aceito pelo Collector.

**Independent Test**: Repetir a recuperação com lote cuja contagem distinta é igual, menor, maior e zero, e com um lote misto em que recordCount conta só os registros aceitos. Conferir a classificação, as duas quantidades quando houver divergência, e que a causa da divergência não foi afirmada pela contagem.

**Acceptance Scenarios**:

1. **Given** a paginação terminou, a quantidade de identificadores distintos válidos é igual a recordCount e isolamento, identificadores e conteúdo passaram, **When** o resultado é fechado, **Then** o lote é declarado completo.
2. **Given** a quantidade distinta válida é menor que recordCount, **When** o resultado é fechado, **Then** a classificação é lote incompleto e informa a quantidade esperada e a quantidade encontrada. O lote não é declarado completo.
3. **Given** a quantidade distinta válida é maior que recordCount, **When** o resultado é fechado, **Then** a classificação é excesso de registros distintos. O lote não é declarado completo.
4. **Given** recordCount é positivo e a busca não devolve documentos, **When** a consulta termina sem falha de acesso, **Then** o resultado não é sucesso. A classificação é lote incompleto, com encontrado igual a zero.
5. **Given** a contagem distinta difere de recordCount, **When** a falha é registrada, **Then** o registro não afirma que a causa foi atraso de visibilidade nem falha de persistência.
6. **Given** o Collector aceitou só parte de um lote misto e a referência traz o recordCount desse subconjunto, **When** os identificadores distintos recuperados coincidem com esse recordCount e as outras validações passam, **Then** o lote é declarado completo. Os registros recusados pelo Collector não são exigidos.
7. **Given** uma contagem distinta é usada como conferência adicional, **When** ela é calculada, **Then** ela aplica o mesmo lote, o mesmo identificador de registro quando a referência o tiver, e a mesma faixa com fim em `maxTimestampMicros + 1`. Essa contagem não substitui a entrega dos registros no sucesso.

---

### User Story 4 - Manter o lote no tenant da referência (Priority: P1)

O tenant e o produto esperados são os da referência validada. A recuperação não reconstrói essa identidade a partir dos atributos dos logs. Com os nomes reais dos campos já fornecidos, a consulta filtra por esses campos e cada documento retornado é conferido com a referência. A credencial é necessária e não substitui esses dois controles. Os documentos devolvidos como sucesso pertencem a esse tenant e a esse produto. Documento de outro tenant ou produto não entra num resultado válido, e o lote não é dado como concluído em silêncio. Falha de isolamento é distinta de indisponibilidade e de lote incompleto. Conhecer organização, stream ou batchId não autoriza o acesso.

**Why this priority**: Uma consulta só por lote, como no exemplo do contrato, não demonstra sozinha o isolamento entre tenants. Entregar o documento errado vazaria o conteúdo.

**Independent Test**: Fornecer os nomes reais dos campos de tenant e produto. Preparar dois tenants com dados sintéticos e uma referência de um deles. Conferir que a consulta filtra esses campos, que cada documento é conferido com a referência e que documento do outro tenant ou produto não sai como resultado válido. Conferir que a classificação é de isolamento e que uma falha de armazenamento continua classificável à parte. Repetir sem os nomes e conferir que a consulta não roda e que não há sucesso.

**Acceptance Scenarios**:

1. **Given** a referência traz tenantId e productId, **When** a recuperação avalia os documentos, **Then** essa identidade é a da referência. Atributos encontrados dentro dos logs não a substituem.
2. **Given** os nomes reais dos campos de tenant e produto estão fornecidos e um documento recuperado pelo lote pertence a outro tenant ou a outro produto, **When** a avaliação termina, **Then** nenhum documento desse lote é entregue como resultado válido e o lote não é declarado completo. A classificação é inconsistência de identidade ou isolamento, distinta de lote incompleto e de indisponibilidade.
3. **Given** os nomes reais dos campos de tenant e produto ainda não foram fornecidos, **When** a recuperação é pedida, **Then** a consulta que filtraria esses campos não é executada e nenhum documento é entregue como resultado válido. A classificação é inconsistência de identidade ou isolamento, distinta de indisponibilidade, e o resultado marca efeito concluído.
4. **Given** a credencial é recusada, **When** a consulta executa, **Then** a classificação é falha de autenticação ou autorização, sem conteúdo de log e sem credencial no registro. Essa falha não é classificada como temporária.
5. **Given** a execução só conhece organização, stream e batchId, ou possui credencial sem o filtro e sem a conferência de cada documento, **When** a autorização é avaliada, **Then** esse conhecimento não basta para entregar os documentos como sucesso.
6. **Given** os nomes reais dos campos de tenant e produto estão fornecidos, **When** a consulta é montada e os documentos voltam, **Then** a consulta filtra tenant e produto pelos valores da referência e cada documento é conferido com esses valores.

---

### User Story 5 - Devolver falha classificada sem sucesso parcial (Priority: P2)

Quando a recuperação não fecha um lote completo válido, o fluxo da primeira tarefa recebe uma classificação, o contexto do lote e as contagens quando existirem. Se houver mais de uma falha, a classificação é uma só, na primeira que se aplicar: autenticação ou autorização, resposta inválida, isolamento ou identidade, documento sem identificador, conteúdo divergente, excesso de registros distintos e lote incompleto. Indisponibilidade, timeout e falha de página só saem quando nenhuma falha que pede correção ocorreu. O conjunto parcial de uma paginação interrompida não é apresentado como recuperação concluída. Timeout, indisponibilidade, falha no meio da paginação e lote incompleto deixam a entrada pendente. Lote completo e falha que pede correção são o efeito concluído que autoriza a primeira tarefa a confirmar a entrada. Esta etapa informa esse efeito e não confirma a entrada. Número de tentativas, intervalos, prazo total e o destino definitivo de lote incompleto não são criados aqui.

**Why this priority**: As histórias 1 a 4 cobrem o lote recuperado e as conferências. Esta história impede que uma falha no caminho seja lida como lote pronto e define quais resultados a primeira tarefa pode confirmar.

**Independent Test**: Provocar timeout, recusa de autorização, resposta inválida, falha na segunda página, lote incompleto, excesso, conteúdo divergente e documento sem identificador. Conferir a classificação, que nenhum caso entrega conjunto parcial como lote completo, que esta etapa não confirma a entrada e que só o lote completo e as falhas que pedem correção saem marcados como efeito concluído. Num lote com isolamento e também documento sem identificador, conferir que a classificação única é isolamento.

**Acceptance Scenarios**:

1. **Given** o armazenamento está indisponível ou a consulta excede o tempo configurado, **When** a tentativa termina, **Then** a classificação é indisponibilidade ou timeout, o lote não é declarado completo e os documentos já lidos não são entregues como recuperação concluída.
2. **Given** uma página posterior falha depois de uma página já lida e nenhuma falha que pede correção está presente, **When** a tentativa termina, **Then** a classificação é falha no meio da paginação, o conjunto parcial não é entregue como recuperação concluída, o lote não é declarado completo e o resultado não marca efeito concluído.
3. **Given** a resposta do armazenamento não pode ser interpretada, **When** a tentativa termina, **Then** a classificação é resposta inválida, distinta de lote incompleto e de indisponibilidade.
4. **Given** qualquer falha desta etapa, **When** o resultado chega ao fluxo da primeira tarefa, **Then** ele traz a classificação e o contexto do lote, com as contagens esperada e encontrada quando a falha for de completude ou de excesso. O resultado de falha não traz conteúdo de log nem credencial.
5. **Given** a recuperação termina em lote completo ou em falha que pede correção, **When** o resultado chega ao fluxo da primeira tarefa, **Then** ele marca efeito concluído. Falha que pede correção é autenticação ou autorização, inconsistência de identidade ou isolamento, resposta inválida, excesso de registros distintos, conteúdo divergente ou documento sem identificador. Esta etapa não confirma a entrada Redis.
6. **Given** a recuperação termina em indisponibilidade, timeout, falha no meio da paginação ou lote incompleto, e nenhuma falha que pede correção está presente, **When** o resultado chega ao fluxo da primeira tarefa, **Then** ele não marca efeito concluído. A entrada permanece pendente para a primeira tarefa. Esta etapa não confirma a entrada Redis.
7. **Given** o lote tem documento sem identificador e também conteúdo divergente ou contagem diferente de recordCount, **When** o resultado é fechado, **Then** a classificação única é documento sem identificador. As quantidades esperada e encontrada seguem informadas quando a contagem foi calculada, e o resultado marca efeito concluído.
8. **Given** o lote tem documento de outro tenant ou produto e também documento sem identificador, conteúdo divergente ou contagem diferente, **When** o resultado é fechado, **Then** a classificação única é inconsistência de identidade ou isolamento. Nenhum documento é entregue como resultado válido e o resultado marca efeito concluído.
9. **Given** a paginação falha ou a consulta expira depois que uma falha que pede correção já foi encontrada, **When** o resultado é fechado, **Then** a classificação é essa falha que pede correção, na ordem desta história. O conjunto parcial não é entregue como recuperação concluída e o resultado marca efeito concluído.

---

### Edge Cases

- URL base, autenticação, timeout, tamanho de página ou lista de campos que só identificam a cópia física ausentes: a recuperação não começa e nenhum valor padrão é criado. Nenhum lote é declarado completo e não há efeito concluído. Lista fornecida vazia é configuração presente.
- Stream Redis `lhamalog:stream:logs` usado no lugar do stream da referência: não atende esta feature.
- Faixa em que o mínimo é igual ao máximo: o fim da consulta é esse instante mais um. O instante do lote continua coberto.
- Página final cheia: a página seguinte é consultada. A paginação não termina na página cheia. Não há número máximo de páginas. A leitura só para quando uma página volta com menos documentos do que o tamanho pedido. O timeout de uma consulta continua sendo falha temporária.
- Cópias com diferença só nos campos configurados: um registro lógico, e o entregue é uma das cópias. Diferença fora dessa lista: conteúdo divergente, nenhuma versão escolhida. A ordenação da consulta precisa ser estável o bastante para que a troca de página não perca outro identificador lógico. O teste com cópias é a prova exigida.
- Documento sem `lhamalog_record_id` no meio de um lote cuja contagem válida já iguala recordCount: o lote não é completo.
- Mesmo identificador com conteúdo divergente: nenhuma versão é escolhida e o lote não é completo. A classificação não é reduzida a incompleto nem a excesso.
- Busca vazia com recordCount positivo: lote incompleto, não sucesso e não indisponibilidade.
- Contagem distinta maior ou menor que recordCount: o lote não é completo. A contagem não declara a causa.
- Documento de outro tenant ou produto no mesmo lote: resultado de isolamento, sem documentos entregues como válidos.
- Nomes dos campos de tenant e produto ainda não fornecidos: a consulta não roda, nenhum documento é entregue e a classificação é isolamento, com efeito concluído. Com os nomes fornecidos, a consulta filtra esses campos e cada documento é conferido. A credencial não substitui esses controles.
- Credencial recusada: falha de autenticação ou autorização, não temporária, sem repetição definida por esta feature.
- Falha no meio da paginação, sem falha que pede correção: o parcial não é recuperação concluída e a entrada permanece pendente. Se uma falha que pede correção já foi encontrada, a classificação segue a ordem: autenticação ou autorização, resposta inválida, isolamento ou identidade, documento sem identificador, conteúdo divergente, excesso e lote incompleto.
- Lote misto: recordCount e a faixa representam só o subconjunto aceito. Registro recusado não é exigido para a completude.
- Deduplicação do evento pela primeira tarefa: não deduplica estes documentos. Esta deduplicação também não cumpre a idempotência do evento.
- Lote completo e falha que pede correção: o resultado marca efeito concluído. A primeira tarefa pode confirmar a entrada. Esta etapa não a confirma.
- Indisponibilidade, timeout, falha no meio da paginação e lote incompleto, sem falha que pede correção junto: o resultado não marca efeito concluído. A entrada permanece pendente.
- Configuração de acesso ausente, inclusive a lista de campos da cópia física: a recuperação não começa, nenhum valor padrão é criado e não há efeito concluído de lote.
- Tentativas, intervalos, prazo total e tratamento definitivo de lote incompleto: não estão no contrato e não são criados. DLQ não é assumida.
- Credencial ou conteúdo de log em registro operacional de falha: não atende esta feature.
- Traces, spans e busca livre: fora desta feature. Um lote de logs não dispara recuperação de trace.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: A recuperação MUST partir da referência validada entregue pela primeira tarefa. Ela MUST usar tenantId, productId, organização, stream, batchId, minTimestampMicros, maxTimestampMicros, recordCount e idempotencyId. Quando a referência incluir recordId, esse identificador MUST entrar na consulta. O Finder MUST NOT reconstruir tenantId nem productId a partir dos atributos dos logs. O conteúdo dos logs MUST NOT ser lido da entrada Redis.
- **FR-002**: URL base, autenticação, timeout, tamanho de página e a lista de campos que só identificam a cópia física MUST ser fornecidos antes da recuperação. Credencial MUST NOT aparecer no código nem em registro operacional. Se qualquer um desses itens estiver ausente, a recuperação MUST NOT começar, MUST NOT criar um valor padrão e MUST NOT marcar efeito concluído. Lista fornecida com zero campos é configuração presente: qualquer diferença entre documentos do mesmo identificador é conteúdo divergente. Esta spec MUST NOT nomear os campos dessa lista. Organização e stream MUST ser os da referência. O stream Redis `lhamalog:stream:logs` MUST NOT ser usado como stream de armazenamento.
- **FR-003**: A consulta MUST ser `POST /api/{organization}/_search`, com a organização da referência. Ela MUST consultar o stream da referência, filtrar `lhamalog_batch_id` pelo batchId, usar início `minTimestampMicros` e fim `maxTimestampMicros + 1`, e usar os nomes normalizados em minúsculas. Em referência com recordId, a consulta MUST filtrar também `lhamalog_record_id`. Quando os nomes reais dos campos de tenant e produto estão fornecidos, a consulta MUST filtrar também esses campos pelos valores da referência. Cada valor da referência MUST entrar só como dado do filtro pretendido. Valor que não possa ser aplicado assim MUST produzir resposta inválida, sem registros entregues como sucesso e sem trocar lote, organização, stream ou faixa.
- **FR-004**: A paginação MUST começar no deslocamento zero e no tamanho de página configurado. O deslocamento da página seguinte MUST aumentar por esse tamanho. A leitura MUST continuar até uma página devolver menos documentos do que o tamanho pedido. Página cheia MUST ser seguida de outra consulta. Esta spec MUST NOT criar um número máximo de páginas. O tamanho da página MUST NOT ser substituído por recordCount como única estratégia de recuperação.
- **FR-005**: A ordem da consulta MUST ser estável entre as páginas, de modo que cópias físicas do mesmo registro não façam a paginação omitir outro identificador lógico. A prova com cópias físicas MUST fazer parte do aceite deste requisito.
- **FR-006**: Os documentos MUST ser deduplicados por `lhamalog_record_id` ao longo de todas as páginas. Diferença restrita aos campos da lista configurada MUST contar como uma cópia e como um registro lógico. O registro entregue MUST ser uma das cópias. Diferença em qualquer campo fora dessa lista MUST ser conteúdo divergente. O identificador Redis e o correlationId MUST NOT deduplicar os logs. Essa deduplicação MUST ser independente da idempotência do evento feita pela primeira tarefa.
- **FR-007**: Documento sem `lhamalog_record_id` MUST NOT entrar na contagem de registros lógicos válidos. A recuperação MUST NOT declarar o lote completo. A classificação emitida segue FR-015.
- **FR-008**: Documentos com o mesmo `lhamalog_record_id` e diferença fora da lista configurada MUST NOT ter uma versão escolhida. O lote MUST NOT ser declarado completo. Nenhuma das versões MUST ser entregue como o registro lógico válido. A classificação emitida segue FR-015.
- **FR-009**: Ao fim da paginação, a quantidade de identificadores distintos válidos MUST ser comparada com recordCount. Igualdade MUST declarar o lote completo somente se isolamento, presença de identificador e ausência de conteúdo divergente também tiverem passado. Quantidade menor é lote incompleto, e quantidade maior é excesso de registros distintos, quando essa for a classificação de FR-015. O resultado MUST informar a quantidade esperada e a encontrada quando a contagem tiver sido calculada. Busca sem documentos, com recordCount positivo e sem falha de acesso nem falha de maior precedência, MUST ser lote incompleto e MUST NOT ser sucesso. A contagem MUST NOT ser apresentada como identificação da causa da divergência.
- **FR-010**: recordCount MUST ser tratado como a quantidade de registros aceitos pelo Collector. A completude MUST NOT exigir registros recusados de um lote misto.
- **FR-011**: Se uma contagem distinta for usada como conferência, ela MUST aplicar os mesmos filtros e a mesma faixa de FR-003. Ela MUST NOT substituir a entrega dos registros no resultado de sucesso.
- **FR-012**: O tenant e o produto esperados MUST ser os da referência. Esta spec MUST NOT nomear os campos persistidos de tenant e produto. Quando esses nomes estão fornecidos, a consulta MUST filtrá-los pelos valores da referência e cada documento retornado MUST ser conferido com esses valores. A credencial MUST ser exigida e MUST NOT substituir o filtro nem a conferência. Documento que não pertence a esse tenant ou produto MUST impedir a entrega de qualquer documento do lote como resultado válido. Enquanto os nomes não estiverem fornecidos, a consulta MUST NOT ser executada, nenhum documento MUST ser entregue como resultado válido e a classificação MUST ser inconsistência de identidade ou isolamento, com efeito concluído. Essa classificação MUST ser distinta de indisponibilidade e de lote incompleto. O lote MUST NOT ser declarado completo nesse caso.
- **FR-013**: Conhecer organização, stream ou batchId MUST NOT ser tratado como autorização de acesso. Falha de autenticação ou autorização MUST ser classificação própria, MUST NOT incluir credencial nem conteúdo de log, e MUST NOT ser classificada como indisponibilidade ou timeout.
- **FR-014**: O resultado de sucesso MUST devolver os registros distintos e a correlação do lote: tenantId, productId, organização, stream, batchId, idempotencyId, recordCount e a quantidade distinta encontrada. O registro entregue MUST ser o documento armazenado, sem ocultar campos nesta etapa.
- **FR-015**: O resultado de falha MUST devolver uma classificação, o contexto do lote e as quantidades esperada e encontrada quando a contagem tiver sido calculada. Ele MUST NOT incluir conteúdo de log nem credencial. As classificações MUST distinguir, no mínimo: indisponibilidade ou timeout; falha de autenticação ou autorização; resposta inválida; lote incompleto; excesso de registros distintos; inconsistência de identidade ou isolamento; conteúdo divergente; documento sem identificador. Quando mais de uma falha está presente, a classificação emitida MUST ser uma só, a primeira que se aplicar nesta ordem: falha de autenticação ou autorização; resposta inválida; inconsistência de identidade ou isolamento; documento sem identificador; conteúdo divergente; excesso de registros distintos; lote incompleto. Indisponibilidade, timeout e falha no meio da paginação MUST ser a classificação somente quando nenhuma falha que pede correção estiver presente.
- **FR-016**: Indisponibilidade, timeout ou falha de uma página, quando for a classificação de FR-015, MUST ser falha temporária, MUST NOT declarar o lote completo e MUST NOT entregar o conjunto parcial como recuperação concluída. Se uma falha que pede correção já tiver sido encontrada, a classificação MUST ser essa falha e o conjunto parcial MUST NOT ser entregue como recuperação concluída.
- **FR-017**: Esta etapa MUST NOT definir número de tentativas, intervalo, prazo total da recuperação nem tratamento definitivo de lote incompleto. Ela MUST NOT assumir DLQ. A classificação de autenticação ou autorização MUST permanecer distinta da falha temporária, para que repetição infinita não seja a leitura dessa falha.
- **FR-018**: O resultado MUST alimentar o ciclo de processamento da primeira tarefa e MUST informar se esse resultado é efeito concluído. Esta etapa MUST NOT confirmar a entrada Redis.
- **FR-019**: O efeito concluído, que autoriza a primeira tarefa a confirmar a entrada, MUST ser o lote completo e qualquer falha que pede correção escolhida por FR-015: autenticação ou autorização, inconsistência de identidade ou isolamento, resposta inválida, excesso de registros distintos, conteúdo divergente e documento sem identificador. As quantidades MUST continuar informadas quando a contagem tiver sido calculada. Indisponibilidade, timeout, falha no meio da paginação e lote incompleto, quando nenhuma falha que pede correção estiver presente, MUST NOT marcar efeito concluído e MUST deixar a entrada pendente. Ausência de URL base, autenticação, timeout, tamanho de página ou da lista de campos da cópia física MUST NOT marcar efeito concluído.
- **FR-020**: A prova integrada MUST usar ambiente de teste e dados sintéticos. Logs de produção MUST NOT ser usados. Atender os cenários desta spec MUST NOT ser apresentado como homologação de produção nem como integração homologada com os consumidores externos.
- **FR-021**: Esta feature MUST NOT consumir a fila Redis, MUST NOT confirmar entrada, MUST NOT varrer dado sensível, MUST NOT publicar checked, MUST NOT integrar com Ocorrência, MUST NOT oferecer busca livre, MUST NOT recuperar traces ou spans e MUST NOT alterar o Collector.

### Key Entities

- **Referência validada**: Recorte recebido da primeira tarefa. Identifica tenant, produto, organização, stream, lote, faixa em microssegundos, quantidade esperada de registros aceitos e a correlação do evento. Pode trazer o identificador do registro quando o lote é unitário. Não contém os logs.
- **Documento físico**: Um item armazenado que satisfaz a consulta. Reenvios podem gerar vários documentos físicos para o mesmo identificador lógico.
- **Registro lógico**: O log distinto depois da deduplicação por `lhamalog_record_id`. É o que compõe a contagem comparada com recordCount e o que o sucesso entrega.
- **Resultado da recuperação**: Sucesso, com os registros lógicos e a correlação do lote, ou falha classificada, com contexto e contagens quando couber, sem conteúdo de log e sem credencial. Informa se é efeito concluído. Não confirma a entrada Redis. Lote completo e falha que pede correção são efeito concluído. Indisponibilidade, timeout, falha no meio da paginação e lote incompleto, sozinhos, não são.
- **Configuração de acesso**: URL base, autenticação, timeout, tamanho de página e a lista de campos que só identificam a cópia física. Não guarda segredos nem os nomes desses campos na especificação. Lista ausente impede o início. Lista vazia é configuração presente. Organização e stream não fazem parte dela; vêm da referência.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Em um conjunto de referências sintéticas válidas, 100% das consultas usam a organização, o stream, o batchId e a faixa da própria referência, com o fim da faixa em um microssegundo além do máximo. A contagem de conteúdos de log obtidos da entrada Redis é zero. A contagem de confirmações de entrada Redis feitas por esta etapa é zero.
- **SC-002**: Em lote unitário, a consulta inclui o identificador do registro. Em lote de vários registros, a consulta não inclui esse filtro. Um valor que não possa ser aplicado só como dado do filtro pretendido produz resposta inválida. A contagem de consultas que trocaram o lote, a organização, o stream ou a faixa é zero.
- **SC-003**: Nos três arranjos de página — uma página, várias páginas e quantidade múltipla exata do tamanho de página — a contagem de identificadores lógicos perdidos por encerrar a leitura cedo é zero. A contagem de paradas causadas por um número máximo de páginas é zero. Quando há cópias físicas e um identificador novo além de recordCount, esse identificador entra no conjunto distinto.
- **SC-004**: Depois da deduplicação, cada `lhamalog_record_id` válido aparece uma vez. Em lote cuja quantidade distinta iguala recordCount e cujas outras validações passam, o lote é declarado completo. Documento sem identificador e par divergente do mesmo identificador deixam a contagem de lotes completos em zero para esse caso.
- **SC-005**: Para busca vazia, contagem menor e contagem maior, a contagem de lotes declarados completos é zero. Incompleto e excesso são distinguíveis, e as quantidades esperada e encontrada estão presentes. A contagem de falhas que afirmam a causa da divergência é zero.
- **SC-006**: Em prova com dois tenants sintéticos e com os nomes dos campos fornecidos, a consulta filtra tenant e produto e cada documento é conferido. A contagem de documentos de outro tenant ou outro produto entregues como resultado válido é zero. Sem esses nomes, a contagem de consultas executadas é zero, a contagem de resultados válidos é zero e o resultado marca efeito concluído por isolamento. Falha de isolamento, indisponibilidade e lote incompleto são distinguíveis entre si. Credencial sem o filtro e sem a conferência deixa a contagem de resultados válidos em zero.
- **SC-007**: Em timeout, indisponibilidade, autorização recusada, resposta inválida e falha no meio da paginação, a contagem de conjuntos parciais entregues como recuperação concluída é zero e a contagem de lotes declarados completos é zero. Autorização recusada não entra na contagem de falhas temporárias. Registro operacional dessas falhas contém zero credenciais e zero conteúdos de log.
- **SC-008**: No lote misto sintético, a completude usa somente a quantidade de registros aceitos. A contagem de lotes reprovados por ausência dos registros recusados é zero.
- **SC-009**: Quando falta URL base, autenticação, timeout, tamanho de página ou a lista de campos da cópia física, a contagem de consultas iniciadas é zero, a contagem de valores padrão criados é zero e a contagem de efeitos concluídos é zero. Com a lista fornecida, diferença só nos campos listados deixa um registro lógico; diferença fora dela deixa a contagem de lotes completos em zero para esse identificador.
- **SC-010**: Em um conjunto com lote completo, falha de autorização, isolamento, resposta inválida, excesso, conteúdo divergente e documento sem identificador, todos esses resultados marcam efeito concluído. Em indisponibilidade, timeout, falha no meio da paginação e lote incompleto sem falha que pede correção, a contagem de efeitos concluídos é zero. Quando isolamento e documento sem identificador ocorrem no mesmo lote, a classificação contada é isolamento. A contagem de confirmações de entrada Redis feitas por esta etapa é zero.

## Assumptions

Decisões e fatos abaixo têm fonte no pedido, na spec já escrita ou em arquivo lido nesta sessão. O que essas fontes deixaram em aberto não recebeu valor novo, salvo o comportamento observável que o próprio pedido já exige.

- Fonte: pedido de 2026-11-08, que cita o contrato do Collector, especialmente as seções 3, 5, 6, 7, 8, 9 e 11. O documento não foi aberto nesta sessão. Os filtros, a faixa com fim em máximo mais um, os nomes `lhamalog_batch_id` e `lhamalog_record_id`, a paginação e a deduplicação repetem o que o pedido atribui a esse contrato.
- Fonte: spec `004-capturar-referencia-lote`, FR-012 e FR-014, lida em 2026-11-08, e resposta do solicitante em 2026-11-08, Q1 B. A referência validada já traz tenant e produto do metadata, organização, stream, batchId, recordId quando o lote é unitário, faixa, recordCount, idempotencyId e a correlação. A interpretação não consulta o armazenamento. A confirmação da entrada válida espera o efeito concluído definido em FR-019: lote completo e falha que pede correção. Indisponibilidade, timeout, falha no meio da paginação e lote incompleto, sozinhos, deixam a entrada pendente. Esta etapa informa o efeito e não confirma a entrada.
- Fonte: spec `004-capturar-referencia-lote`. Lote de um registro traz recordId. Lote de vários registros não traz. recordCount descreve os registros aceitos, não o lote bruto recusado. A deduplicação por `lhamalog_record_id` foi deixada para esta etapa e não substitui a chave por idempotencyId.
- Fonte: pedido, seções 2 e 9, e o precedente da spec 004 de não criar número ausente no contrato. Timeout e tamanho de página precisam existir antes da execução e não foram numerados aqui. URL base e autenticação também. Ausência impede o início e não gera padrão.
- Fonte: resposta do solicitante em 2026-11-08, opção B. Não há número máximo de páginas. A leitura só para quando uma página devolve menos documentos do que o tamanho pedido. O timeout configurado de cada consulta permanece falha temporária.
- Fonte: resposta do solicitante em 2026-11-08, opção A. A lista de campos que só identificam a cópia física é configuração obrigatória. Os nomes não foram informados e não foram criados nesta spec. Sem a lista, a recuperação não começa e não marca efeito concluído. Lista vazia é configuração presente e torna qualquer diferença conteúdo divergente. Com campos na lista, diferença restrita a eles é uma cópia e o registro entregue é uma das cópias. Qualquer outra diferença é conteúdo divergente, sem versão escolhida.
- Fonte: pedido, seções 4 e 6, constituição, princípio V, e resposta do solicitante em 2026-11-08, opção A. Os nomes reais dos campos de tenant e produto não foram informados e não foram criados. A regra, quando os nomes existem, é filtrar a consulta por eles e conferir cada documento com a referência. A credencial é necessária e não substitui esses controles. Sem os nomes, a consulta não roda, nenhum documento é entregue e a classificação é isolamento, com efeito concluído. Organização, stream ou batchId não autorizam o acesso.
- Fonte: pedido, seções 5 e 6, e a resposta desta sessão sobre cópia e divergência. Cópia física é reenvio do mesmo identificador com diferença só na lista configurada. Conteúdo divergente é qualquer outra diferença sob o mesmo identificador. A prova fornece a lista e constrói os dois casos.
- Fonte: pedido, seção 8, e spec 004, resposta de 2026-11-08 sobre o modelo síncrono. O resultado volta ao ciclo da primeira tarefa no processamento que espera a recuperação. Esta etapa não cria persistência, cache nomeado nem outro armazenamento intermediário para os logs. O registro de sucesso é o documento armazenado, inteiro. Ocultar dados sensíveis continua fora, como na spec `003-consumir-logs-openobserve`.
- Fonte: pedido, seções 8 e 9, e resposta do solicitante em 2026-11-08, opção A. O resultado traz uma classificação. Quando há mais de uma falha, vale a primeira desta ordem: autenticação ou autorização, resposta inválida, isolamento ou identidade, documento sem identificador, conteúdo divergente, excesso e lote incompleto. Indisponibilidade, timeout e falha de página só valem quando nenhuma falha que pede correção ocorreu. Número de tentativas, intervalos, prazo total e tratamento definitivo de lote incompleto não estão, e não foram criados. DLQ não é assumida.
- Fonte: spec `003-consumir-logs-openobserve` e `internal/adapter/out/openobserve/doc.go`, lidos em 2026-11-08. A spec anterior cobre captura de um log por apontamento mockado e proíbe nomear organização, stream, campos e credenciais. O pacote do adaptador reserva a leitura e declara que não consulta o armazenamento. Esta feature define a recuperação do lote pela referência. Ela não afirma que o módulo atual já a executa.
- Fonte: pedido, premissas. A integração com os consumidores externos não está homologada na documentação citada. O aceite desta feature não declara essa homologação.
- Fonte: constituição, princípios I, II e IV, e pedido, fora de escopo. O fluxo normativo segue até a ocorrência. Esta feature cobre a recuperação dos logs do lote. A justificativa do recorte é o pedido: fila, ACK, varredura, checked, Ocorrência, traces e Collector ficam para outras etapas.
- Fonte: constituição, princípio V. Código, teste e comando desta feature não foram executados durante a escrita desta spec. Nenhum resultado de execução é afirmado.

## Out of Scope

- Consumo do stream Redis e o ato de confirmar a entrada. A primeira tarefa confirma quando o resultado marca o efeito concluído de FR-019. Esta etapa só informa esse efeito.
- Varredura de dados sensíveis ou de prompt injection, e ocultação de campos do log recuperado.
- Publicação de checked.
- Integração com Ocorrência.
- Busca livre fora do lote indicado pela referência.
- Recuperação de traces e spans.
- Alteração do Collector.
- Número de tentativas, intervalos, prazo total, DLQ e tratamento definitivo de lote incompleto.
- Homologação da integração com consumidores externos.
