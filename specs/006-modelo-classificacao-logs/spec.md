# Feature Specification: Modelo de classificação de logs

**Feature Branch**: `006-modelo-classificacao-logs`

**Created**: 2026-11-09

**Status**: Draft

**Input**: User description: "Definir e implementar um modelo capaz de classificar logs de forma eficiente e rápida sem LLM. A classificação categoriza o log como informação, alerta ou erro. A captura já está em parte implementada. O Finder acessa a fila para identificar o log, coleta o log capturado pelo Coletor e salvo no armazenamento, classifica com o modelo escolhido e devolve o resultado como ocorrência salva no MongoDB."

## Clarifications

### Session 2026-11-09

- Q: Qual conjunto de classes vale nesta feature? → A: Informação, alerta e erro. Não existe spam nem "não classificado". A constituição deste repositório permanece como está. A quantidade ficou uma ocorrência por apontamento, não uma por log.
- Q: Quem grava a ocorrência e o que ela contém? → A: O Finder grava no MongoDB a classe, a versão das regras e a identidade do apontamento, com o tenant e o produto. A ocorrência fica consultável sem copiar o texto dos logs. Mascaramento, gravidade, evento e sugestão continuam fora.
- Q: Quando o nível do log e o catálogo apontam classes diferentes, qual vale? → A: Vale o nível do log. O catálogo só classifica quando o log não traz nível reconhecível. A resposta D seguinte inclui o campo e o texto nesse nível; a classe do log é a mais grave entre os dois. O catálogo continua de fora quando um dos dois já decide.
- Q: Quando um apontamento da fila já trouxe vários logs capturados, quantas ocorrências esta feature grava? → A: O apontamento inteiro vira uma única ocorrência, mesmo quando há vários logs.
- Q: Quando os logs de um mesmo apontamento recebem classes diferentes, qual classe fica na ocorrência única? → A: A mais grave. Erro vence alerta, e alerta vence informação.
- Q: O nível que decide a classe de um log é um campo próprio desse log, ou uma palavra igual encontrada no texto da mensagem? → A: Campo e texto valem juntos. A classe do log é a mais grave entre os dois.
- Q: Quando a mensagem traz um termo de nível no meio de uma palavra maior, esse termo conta para a classe? → A: Só a palavra inteira. "error" conta; "errorless" e "traceback" não.
- Q: Quando um apontamento tem vários logs e um deles não está disponível, a feature grava a ocorrência com os demais? → A: Grava uma ocorrência só com os logs presentes. A classe é a mais grave entre eles.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Classificar o log capturado sem serviço de linguagem externo (Priority: P1)

O Finder parte do log que a captura já deixou disponível. O apontamento identifica esse conjunto e não é o conteúdo. O modelo lê o campo de nível e o texto da mensagem. A classe do log é a mais grave entre os dois: erro vence alerta, e alerta vence informação. O texto não é enviado a um serviço externo de linguagem. A versão das regras usadas fica no resultado.

**Why this priority**: Sem a classe atribuída no próprio Finder, não há ocorrência classificada. É o resultado mínimo do pedido e da feature LLT-133.

**Independent Test**: Entregar um log já capturado com nível de erro e conferir que a classe é erro, que nenhum texto desse log é enviado a um serviço externo de linguagem e que o resultado informa quais regras foram usadas.

**Acceptance Scenarios**:

1. **Given** um log já capturado cujo campo de nível é erro, error ou fatal, e a mensagem não tem termo de nível, **When** a classificação termina, **Then** a classe é erro. O texto do log não é enviado a um serviço externo de linguagem. O resultado informa as regras usadas.
2. **Given** um log já capturado cujo campo de nível é alerta, warn ou warning, e a mensagem não tem termo de nível mais grave, **When** a classificação termina, **Then** a classe é alerta.
3. **Given** um log já capturado cujo campo de nível é informação, info, debug ou trace, e a mensagem não tem termo de nível mais grave, **When** a classificação termina, **Then** a classe é informação.
4. **Given** o campo de nível é informação e a mensagem contém a palavra inteira error, **When** a classificação termina, **Then** a classe do log é erro. O catálogo não altera essa classe. O texto não é enviado a um serviço externo de linguagem.
5. **Given** o campo de nível é fatal e a mensagem contém a palavra inteira info, **When** a classificação termina, **Then** a classe do log é erro.
6. **Given** o campo de nível é informação e a mensagem contém errorless, **When** a classificação termina, **Then** a classe do log é informação. "errorless" não conta como erro.
7. **Given** o campo de nível está ausente e a mensagem contém traceback, **When** a classificação termina, **Then** traceback não conta como informação. Sem outro termo de nível e sem catálogo, a classe é informação.
8. **Given** a fila só traz o apontamento do conjunto, **When** a classificação é pedida, **Then** a classe é calculada sobre os logs capturados. O apontamento não é tratado como o conteúdo.
9. **Given** o apontamento não tem nenhum log capturado disponível, **When** a classificação é pedida, **Then** nenhuma ocorrência é gravada e a falha distingue log ausente de armazenamento inacessível.

---

### User Story 2 - Gravar uma ocorrência por apontamento, sem o texto do log (Priority: P1)

O Finder grava uma ocorrência no MongoDB para o apontamento inteiro. Um apontamento com vários logs capturados continua gerando uma ocorrência, não uma por log. A ocorrência guarda a classe, a identificação das regras, a identidade do apontamento, o tenant e o produto. Ela não guarda o texto dos logs. A classe da ocorrência é a mais grave entre os logs desse apontamento: erro vence alerta, e alerta vence informação. Se todos têm a mesma classe, a ocorrência usa essa classe. Depois que o processamento termina, a ocorrência continua consultável. Agrupar logs parecidos por assinatura, publicar evento, calcular gravidade, abrir sugestão, notificar e mascarar ficam fora. Esta feature não lê a fila de novo e não confirma a entrada.

**Why this priority**: O pedido pede o resultado em forma de ocorrência salva. A resposta B de 2026-11-09 define que o apontamento inteiro é essa ocorrência.

**Independent Test**: Entregar um apontamento com três logs já capturados, todos com nível de erro, e o MongoDB disponível. Consultar pela identidade do apontamento e conferir que existe uma ocorrência, com classe erro, sem o texto dos logs. Conferir que a fila não foi lida de novo e que a entrada não foi confirmada.

**Acceptance Scenarios**:

1. **Given** um apontamento já trouxe três logs capturados e o MongoDB está disponível, **When** a gravação termina, **Then** uma consulta pela identidade desse apontamento devolve uma ocorrência, com a identificação das regras, o tenant e o produto, e sem o texto dos logs.
2. **Given** todos os logs desse apontamento recebem a mesma classe, **When** a gravação termina, **Then** a classe da ocorrência é essa classe.
3. **Given** o apontamento tem um log de informação, um de alerta e um de erro, **When** a gravação termina, **Then** continua uma ocorrência e a classe é erro.
4. **Given** o apontamento tem um log de informação e um de alerta, e nenhum de erro, **When** a gravação termina, **Then** a classe da ocorrência é alerta.
5. **Given** a gravação no MongoDB não pode ser concluída, **When** a tentativa termina, **Then** esse apontamento não é apresentado como ocorrência pronta e o motivo da falha de gravação fica distinguível da falha de log ausente.
6. **Given** a gravação foi concluída, **When** o resultado é conferido, **Then** não há evento publicado, sugestão gerada nem notificação disparada. A fila não é lida de novo e a entrada não é confirmada. Juntar logs de mesma assinatura continua fora.
7. **Given** o apontamento esperava três logs e só dois estão presentes, um de informação e um de erro, **When** a gravação termina, **Then** existe uma ocorrência e a classe é erro. O log ausente não entra na classe.

---

### User Story 3 - Usar o catálogo só quando o nível não decide (Priority: P2)

O modelo considera o campo de nível e os termos de nível no texto da mensagem. A classe do log é a mais grave entre esses dois sinais. O catálogo de expressões e palavras-chave só entra quando nem o campo nem o texto trazem termo de nível reconhecível. O catálogo não é criado por esta feature. Sem esses sinais e sem correspondência no catálogo, a classe é informação e a ocorrência ainda é gravada.

**Why this priority**: A ADR_167 já escolheu regras e padrões para o MVP e deixou o catálogo concreto em aberto. A resposta de 2026-11-09 definiu que o nível vence o catálogo.

**Independent Test**: Repetir a classificação com nível de erro e catálogo apontando informação, com nível ausente e catálogo apontando alerta, e sem nível nem catálogo. Conferir as três classes e que nenhum catálogo foi criado no lugar do que faltou.

**Acceptance Scenarios**:

1. **Given** o campo de nível é erro e o catálogo fornecido apontaria informação, **When** a classificação termina, **Then** a classe é erro. O catálogo não altera essa classe.
2. **Given** o campo de nível é informação, a mensagem contém error e o catálogo apontaria informação, **When** a classificação termina, **Then** a classe do log é erro. O catálogo não entra.
3. **Given** nem o campo nem a mensagem trazem termo de nível e a primeira entrada correspondente do catálogo aponta alerta, **When** a classificação termina, **Then** a classe é alerta.
4. **Given** não há catálogo e nem o campo nem a mensagem trazem termo de nível, **When** a classificação termina, **Then** a classe é informação, nenhum catálogo é criado e a ocorrência é gravada como nas demais classes.
5. **Given** nem o campo nem a mensagem trazem termo de nível e mais de uma entrada do catálogo corresponde ao log, **When** a classificação termina, **Then** vale a primeira entrada correspondente, na ordem em que o catálogo foi fornecido.

---

### Edge Cases

- Apontamento sem nenhum log capturado: nenhuma ocorrência é gravada. A falha informa que o log está ausente.
- Apontamento com alguns logs presentes e outros ausentes: uma ocorrência só com os presentes. A classe é a mais grave entre eles. O ausente não entra na classe e não impede a gravação.
- Armazenamento do log inacessível: nenhuma ocorrência é gravada. A falha é distinta da de log ausente e da de gravação.
- MongoDB indisponível: a classe pode ter sido calculada, mas a ocorrência não fica pronta. A falha é distinta das duas anteriores.
- Campo de nível vazio, ausente ou fora da lista, e mensagem sem termo de nível: o campo não decide. O catálogo entra. Sem correspondência, a classe é informação e a ocorrência é gravada.
- Campo e texto com termos de gravidade diferente: a classe do log é a mais grave. Erro vence alerta, e alerta vence informação. O catálogo não entra.
- Campo informação e mensagem com error: a classe do log é erro.
- Campo ausente e mensagem com a palavra inteira warn: a classe do log é alerta. O catálogo não entra.
- Palavra inteira: o termo no texto conta quando não está colado a letra ou dígito. "error", "error," e "(error)" contam. "errorless" e "traceback" não contam. Maiúsculas não mudam o reconhecimento.
- Várias entradas do catálogo correspondem e não há nível reconhecível: vale a primeira, na ordem fornecida.
- Entrada de catálogo com classe diferente de informação, alerta ou erro: essa entrada é ignorada. Ela não cria uma quarta classe.
- Vários logs no mesmo apontamento: uma ocorrência. A ocorrência não copia o texto desses logs.
- Logs do mesmo apontamento com a mesma classe: a ocorrência usa essa classe.
- Logs do mesmo apontamento com classes diferentes: continua uma ocorrência. A classe gravada é a mais grave: erro vence alerta, e alerta vence informação.
- O mesmo apontamento gravado outra vez: outra ocorrência, sem atualizar a anterior. Esta feature não funde, não conta repetição e não marca recorrência.
- Texto com segredo: a ocorrência não recebe o texto do log. O catálogo de mascaramento não é definido aqui. O trecho bruto não fica pronto para tela, e-mail ou sugestão por esta feature.
- Serviço externo de linguagem disponível: a classificação não o usa.
- Span ou outro sinal que não seja o log capturado: fica fora. A classe não é calculada sobre ele.
- Maiúsculas, minúsculas e espaços nas bordas do nível: não mudam o reconhecimento. " Error " é erro.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: O Finder MUST classificar cada log já capturado de um apontamento. O apontamento MUST identificar esse conjunto e MUST NOT ser tratado como o conteúdo. Esta feature MUST NOT ler a fila de novo e MUST NOT confirmar a entrada.
- **FR-002**: A classificação MUST ocorrer sem enviar o texto do log a um serviço externo de linguagem e sem esperar resposta desse serviço. Um classificador treinado com histórico e a recorrência por similaridade ficam fora desta feature.
- **FR-003**: A classe de cada log, e a classe da ocorrência, MUST ser uma destas: informação, alerta ou erro. Spam e "não classificado" MUST NOT ser classes desta feature. Um apontamento MUST gerar uma ocorrência, tenha ele um log ou vários.
- **FR-004**: O campo de nível e o texto da mensagem MUST contar juntos. A correspondência, sem distinção de maiúsculas e ignorando espaços nas bordas do campo, é: erro, error e fatal são erro; alerta, warn e warning são alerta; informação, info, debug e trace são informação. No texto, o termo MUST contar só como palavra inteira: conta quando não está colado a letra ou dígito. "error", "error," e "(error)" contam. "errorless" e "traceback" MUST NOT contar. A classe do log MUST ser a mais grave entre o campo e as palavras inteiras encontradas na mensagem. Erro vence alerta, e alerta vence informação. Campo vazio, ausente ou fora dessa lista MUST NOT decidir sozinho quando a mensagem também não tem palavra inteira dessa lista.
- **FR-005**: O catálogo de expressões e palavras-chave MUST ser usado somente quando o campo não é reconhecível e a mensagem não tem palavra inteira de nível da lista de FR-004. Esta feature MUST NOT criar esse catálogo. Quando várias entradas correspondem, MUST valer a primeira, na ordem fornecida. Entrada cuja classe não seja informação, alerta ou erro MUST ser ignorada.
- **FR-006**: Sem termo de nível no campo, sem termo de nível na mensagem e sem correspondência válida no catálogo, a classe desse log MUST ser informação. Nenhuma classe extra pode ser apresentada como se um catálogo tivesse sido criado.
- **FR-007**: O Finder MUST gravar uma ocorrência no MongoDB para o apontamento, com a classe, a identificação das regras usadas, a identidade do apontamento, o tenant e o produto. A ocorrência MUST NOT conter o texto dos logs. A classe da ocorrência MUST ser a mais grave entre as classes dos logs desse apontamento: erro vence alerta, e alerta vence informação. Quando todas são iguais, a ocorrência MUST usar essa classe. A identificação das regras MUST dizer se um catálogo entrou e MUST guardar a identificação fornecida com esse catálogo. Se nenhuma identificação foi fornecida, a ocorrência MUST informar que só as regras de nível desta feature foram usadas.
- **FR-008**: Se o apontamento não tem nenhum log presente, se o armazenamento desses logs não pode ser alcançado ou se a gravação no MongoDB não conclui, o Finder MUST NOT apresentar ocorrência pronta. Os três motivos MUST ser distinguíveis. Se ao menos um log do apontamento está presente e outro está ausente, o Finder MUST gravar uma ocorrência só com os logs presentes. A classe MUST ser a mais grave entre os presentes. O log ausente MUST NOT entrar na classe e MUST NOT impedir essa gravação.
- **FR-009**: Esta feature MUST NOT agrupar logs semelhantes, MUST NOT publicar evento de ocorrência, MUST NOT calcular gravidade, MUST NOT gerar sugestão, MUST NOT notificar e MUST NOT definir o catálogo de mascaramento. Um trecho bruto de log MUST NOT ser tratado como pronto para tela, e-mail ou sugestão.
- **FR-010**: A classe de um log MUST poder ser conferida sem depender da classe de outro log. A gravação MUST produzir uma ocorrência por apontamento, não uma por log. Duas gravações bem-sucedidas do mesmo apontamento MUST produzir duas ocorrências, sem atualizar a primeira. Meta numérica de precisão sobre base rotulada não faz parte desta feature, porque essa meta não foi acordada.

### Key Entities

- **Log capturado**: Registro já recuperado a partir do apontamento da fila. É a entrada da classificação. Não é o apontamento e não é a ocorrência.
- **Nível reconhecível**: Termo da lista de FR-004 no campo de nível, ou a mesma palavra inteira no texto da mensagem. No texto, a palavra não pode estar colada a letra ou dígito. A classe do log fica com o mais grave entre o campo e essas palavras. Fora dessa lista, esse sinal não decide a classe.
- **Catálogo de classificação**: Lista opcional, fornecida para a execução, de expressões ou palavras-chave com uma das três classes. Só participa quando o nível não é reconhecível. Não é criado por esta feature.
- **Apontamento**: Identificador do conjunto de logs já capturados. Não é o conteúdo desses logs. Uma gravação bem-sucedida produz uma ocorrência para ele.
- **Ocorrência**: Registro gravado pelo Finder no MongoDB para um apontamento. Guarda uma classe, a identificação das regras, a identidade do apontamento, o tenant e o produto. Não guarda o texto dos logs. A classe é a mais grave entre os logs do apontamento: erro vence alerta, e alerta vence informação. Gravidade, status de tratamento, sugestão, agrupamento por assinatura e notificação não fazem parte desta entidade nesta feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Em uma classificação concluída, o observador identifica uma única classe entre informação, alerta e erro, e identifica quais regras foram usadas. A contagem de envios do texto desse log a um serviço externo de linguagem é zero.
- **SC-002**: Para um apontamento já capturado, com um ou vários logs e com a gravação disponível, uma ocorrência pode ser consultada de novo em até 2 segundos depois do pedido de classificação. A consulta devolve uma ocorrência, com a classe, a identificação das regras, a identidade do apontamento, o tenant e o produto, e não devolve o texto dos logs.
- **SC-003**: Quando o apontamento não tem nenhum log presente, quando o armazenamento não pode ser alcançado e quando a gravação falha, nenhuma ocorrência pronta é consultável. Em cada um, o observador distingue o motivo. Quando o apontamento tem logs presentes e outros ausentes, uma ocorrência é consultável, a classe é a mais grave entre os presentes e o log ausente não muda essa classe.
- **SC-004**: Em uma execução sem catálogo fornecido, a contagem de catálogos criados por esta feature é zero. Um log com nível error recebe erro. Um log sem nível reconhecível recebe informação. Um apontamento cujos logs são todos de nível error produz uma ocorrência de classe erro.
- **SC-005**: Quando o campo é error e o catálogo apontaria informação, a classe desse log é erro. Quando o campo é informação e a mensagem contém a palavra inteira error, a classe desse log é erro, mesmo que o catálogo diga informação. Quando a mensagem contém errorless e o campo é informação, a classe permanece informação. Quando nem o campo nem uma palavra inteira da mensagem têm termo de nível e a primeira entrada correspondente do catálogo é alerta, a classe desse log é alerta. Um apontamento com um log de informação, um de alerta e um de erro produz uma ocorrência de classe erro.
- **SC-006**: Ao fim de uma classificação concluída, a contagem de eventos publicados, sugestões geradas e notificações disparadas por esta feature é zero.

## Assumptions

Fato, premissa e pergunta em aberto ficam separados. As respostas de 2026-11-09 estão na sessão de clarificações. Premissa sem fonte no pedido, no repositório ou nas fontes lidas fica marcada como premissa.

- Fonte: respostas do solicitante em 2026-11-09. Q1 A: as classes são informação, alerta e erro; spam e "não classificado" não existem nesta feature. Q2 A: o Finder grava no MongoDB a classe, a identificação das regras e a identidade do apontamento, além do tenant e do produto; o texto dos logs não é copiado. Q3 A: o nível vence o catálogo. Opção B: um apontamento grava uma ocorrência, mesmo com vários logs. Opção A seguinte: a classe dessa ocorrência é a mais grave; erro vence alerta, e alerta vence informação. Opção D: o campo de nível e o texto da mensagem valem juntos; a classe do log é a mais grave entre os dois. O catálogo só entra quando nenhum dos dois traz termo de nível. Opção A seguinte: no texto, só a palavra inteira conta. "error" conta; "errorless" e "traceback" não. A palavra conta quando não está colada a letra ou dígito. Opção B seguinte: se falta um log e outros estão presentes, a ocorrência é gravada só com os presentes e a classe é a mais grave entre eles. Sem nenhum log presente, não há ocorrência.
- Fonte: pedido de 2026-11-09. A captura já existe em parte. Esta feature começa no log capturado. Não redefine leitura da fila, confirmação da entrada, paginação nem deduplicação. Essas regras permanecem nas specs `001` a `005`.
- Fonte: constituição, princípios I, III e IV, ratificada em 2026-10-29. O Finder transforma log coletado em ocorrência. As classes da constituição são erro, alerta e informação, as mesmas da resposta Q1 A. O apontamento da fila não é o conteúdo. A constituição não definia a persistência; a resposta Q2 A define a gravação no MongoDB por esta feature, sem o texto do log.
- Fonte: ADR_167, aceita em 2026-10-29, lida no Drive. No MVP a análise não usa serviço externo de linguagem. A classificação é por regras e padrões: nível do log, expressões e palavras-chave. A recorrência por similaridade e o classificador treinado ficam fora. O catálogo concreto continua pendente nessa ADR e não é criado aqui.
- Fonte: ADR_019, aceita em 2026-10-02. Ela usa erro, alerta e spam. A resposta Q1 A não adota essa taxonomia. A constituição não é emendada.
- Fonte: LLT-133, LLT-197 e LLT-199, lidos em 2026-11-09. A identificação das regras na ocorrência cobre o pedido de versão do LLT-199. A meta de precisão não entra, porque não foi acordada. A base rotulada é o LLT-198 e fica fora.
- Fonte: LLT-134, LLT-200, LLT-201 e LLT-202. Agrupar logs, publicar evento e mascarar ficam fora. A resposta Q2 A não trouxe o texto do log, então o mascaramento não entrou.
- Fonte: ADR_062. Ela manda o módulo de Ocorrência criar a ocorrência e a gravidade. A resposta Q2 A atribui a gravação desta ocorrência ao Finder e deixa a gravidade fora. Não emenda a ADR.
- Premissa desta spec, para o aceite ser verificável: a lista de níveis reconhecíveis é a de FR-004. Outros textos de nível não decidem a classe. Várias correspondências no catálogo usam a ordem fornecida, a primeira válida. Sem sinal suficiente, a classe é informação, porque Q1 A não deixa uma quarta classe e exige ocorrência na classificação bem-sucedida.
- Premissa desta spec, alinhada à resposta B: cada gravação bem-sucedida de um apontamento cria uma ocorrência. A segunda gravação do mesmo apontamento não atualiza a primeira. Fusão por assinatura e recorrência continuam no LLT-200. Esta feature não lê a fila de novo e não confirma a entrada; isso permanece nas specs `001` a `005`.
- Fonte: pedido, "rápido". O limite de 2 segundos em SC-002 é premissa desta spec para um log já capturado, com o MongoDB disponível, observada por quem pede a classificação. Não é medida de capacidade interna.
- Fonte: constituição, princípio V. Nenhum teste desta feature foi executado durante a escrita desta spec.
