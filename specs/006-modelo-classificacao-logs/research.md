# Research: Modelo de classificação de logs

Decisões desta feature. Cada uma tem fonte na spec, no `go.mod`, no README ou no `docker-compose.yml` lidos em 2026-11-09. Nenhuma preenche URI, banco, coleção ou o catálogo que a spec deixou externos.

Nenhuma gravação foi feita no MongoDB local durante esta escrita.

## Onde a classe é calculada

- **Decision**: A classe de cada log e a classe da ocorrência ficam em `internal/domain`. O caso de uso em `internal/application` percorre os logs presentes, pede a classe e grava uma ocorrência. O adaptador MongoDB não classifica.
- **Rationale**: A spec separa a classe do lugar da gravação. A prova de palavra inteira, de gravidade e de catálogo precisa rodar sem MongoDB, como `make test` já faz nas features anteriores.
- **Alternatives considered**: Classificar dentro do adaptador. Rejeitado porque a regra ficaria presa ao driver. Chamar um serviço externo de linguagem. Rejeitado por FR-002.

## Campo e mensagem já separados

- **Decision**: A chamada entrega, para cada log presente, o texto do campo de nível e o texto da mensagem. Esta feature não escolhe o nome desses campos no OpenObserve e não consulta o armazenamento.
- **Rationale**: FR-004 fala no campo e na mensagem, sem nomeá-los. A spec 005 já trata nomes de campo não informados como configuração externa. Inventar `severity` ou `body` escolheria o nome omitido.
- **Alternatives considered**: Ler `severity_text` e `body` do modelo OpenTelemetry. Rejeitado porque a spec não usa esses nomes. Abrir o OpenObserve de novo. Rejeitado porque a captura já entregou o log e esta feature não redefine essa leitura.

## Palavra inteira

- **Decision**: No texto da mensagem, um termo da lista conta quando o trecho é igual ao termo, sem distinção de maiúsculas, e o caractere imediatamente antes e o imediatamente depois não são letra nem dígito. Início e fim do texto contam como borda. "error", "error," e "(error)" contam. "errorless" e "traceback" não contam. O campo de nível é o valor inteiro, com espaços das bordas removidos, comparado à mesma lista. Não há motor de expressão regular.
- **Rationale**: FR-004 e a resposta A de 2026-11-09. Letra e dígito cobrem o "colado" da spec. Um motor de expressão regular traria sintaxe, expressão inválida e limite de tempo que a spec não definiu.
- **Alternatives considered**: Fronteira `\b` de expressão regular. Rejeitado pelo motivo acima e porque `\b` não trata letra acentuada como a spec trata "letra". Contar o termo no meio de outra palavra. Rejeitado pela resposta A.

## Gravidade

- **Decision**: A ordem é erro, depois alerta, depois informação. A classe do log é a mais grave entre o campo e as palavras inteiras da mensagem. A classe da ocorrência é a mais grave entre os logs presentes. Log ausente não entra. Catálogo não entra se o campo ou uma palavra inteira já decidiu.
- **Rationale**: Resposta A da classe da ocorrência e resposta D do campo com o texto, ambas em 2026-11-09. FR-005 reserva o catálogo ao caso sem nível reconhecível.
- **Alternatives considered**: Primeira palavra da mensagem. Rejeitado porque a resposta D pede a mais grave. Deixar o catálogo vencer o nível. Rejeitado pela resposta Q3 A.

## Catálogo

- **Decision**: Cada entrada fornecida tem uma classe e um texto. A classe fora de informação, alerta e erro é ignorada. As demais são provadas na ordem, e a primeira que aparece como palavra inteira na mensagem, pela mesma regra acima, decide. Frase com espaço precisa aparecer com esse espaço, e as bordas da frase inteira é que não podem estar coladas a letra ou dígito. Sem correspondência, a classe do log é informação. Esta feature não cria entrada.
- **Rationale**: FR-005 e FR-006. A clarificação deixou a correspondência interna para o plano. Usar a mesma palavra inteira evita um segundo critério. A ADR_167 cita expressões regulares, mas não fixa a sintaxe nem o que fazer com expressão inválida. O catálogo concreto continua pendente nessa ADR.
- **Alternatives considered**: Compilar cada entrada como expressão regular. Rejeitado porque a spec não define dialeto, flag nem falha de compilação. Substring solta. Rejeitado porque a resposta A acabou de recusar o termo preso a outra palavra para o nível.

## Uma ocorrência e logs ausentes

- **Decision**: A chamada traz a lista de logs do apontamento e, em cada um, se está presente. Presente entra na classe. Ausente é ignorado. Com ao menos um presente, há uma inserção. Com nenhum presente, o resultado é log ausente e não há inserção. Se a chamada informa armazenamento inacessível, não há inserção, mesmo que a lista venha junto. Falha de gravação também não deixa ocorrência pronta. Os três motivos saem distintos.
- **Rationale**: FR-003, FR-008 e a resposta B de 2026-11-09. A classe mais grave só olha os presentes, então um erro que não chegou não pode ser inventado.
- **Alternatives considered**: Recusar o apontamento inteiro quando falta um log. Rejeitado pela resposta B. Guardar na ocorrência que faltou log. Rejeitado porque essa opção era a C, não a escolhida.

## Gravação

- **Decision**: O adaptador insere um documento por chamada bem-sucedida. Não atualiza o documento anterior do mesmo apontamento. A consulta pela identidade devolve os documentos inseridos. URI, banco e coleção vêm de `FINDER_MONGO_URI`, `FINDER_MONGO_BANCO` e `FINDER_MONGO_COLECAO`. Sem um deles, a chamada não grava e o motivo é configuração ausente, distinto dos três de FR-008. O driver é `go.mongodb.org/mongo-driver/v2`, acrescentado agora. `go.mod` hoje não o tem.
- **Rationale**: FR-007 e FR-010. O README nomeia esse driver como candidato ainda não adicionado. A imagem local é `mongo:8.0.34` em `docker-compose.yml`. Escolher banco ou coleção no código repetiria o que a spec não nomeou. A constituição não definia a persistência; a resposta Q2 A define a gravação por esta feature.
- **Alternatives considered**: Upsert pela identidade do apontamento. Rejeitado por FR-010. Adiar o driver e guardar em memória. Rejeitado porque a consulta pedida é a da ocorrência no MongoDB. Usar o banco `lhamalog_collector` do Compose como padrão. Rejeitado porque esse nome é do inicializador local, não da spec.

## O que a ocorrência guarda

- **Decision**: Classe, identificação das regras, se o catálogo entrou, identidade do apontamento, tenant e produto. O catálogo entrou quando a classe de algum log presente saiu de uma entrada dele. Se a chamada traz identificação do catálogo, ela é guardada. Se não traz, a identificação informa que só as regras de nível desta feature foram usadas. O texto do campo e o texto da mensagem não são gravados nem escritos em falha.
- **Rationale**: FR-007 e Q2 A. O flag e a frase cobrem os dois lados da mesma frase da spec: dizer se o catálogo entrou, guardar a identificação fornecida, e, sem identificação, informar que só as regras de nível foram usadas.
- **Alternatives considered**: Gravar a mensagem mascarada. Rejeitado porque o mascaramento é o LLT-202 e a resposta Q2 A exclui o texto. Gravar a classe de cada log. Rejeitado porque a ocorrência tem uma classe.

## Tenant e produto

- **Decision**: Tenant e produto são os do apontamento recebido. Não são lidos de dentro da mensagem. Se um dos dois vier vazio, não há ocorrência pronta e o motivo é identidade incompleta, distinto de log ausente, de armazenamento inacessível e de falha de gravação.
- **Rationale**: FR-007 exige os dois na ocorrência. A spec 005 já trata o tenant e o produto da referência como autoridade, não o atributo dentro do log. Esta feature não reabre essa leitura.
- **Alternatives considered**: Procurar tenant na mensagem. Rejeitado porque a mensagem pode ser de outro tenant e a ocorrência não deve copiar esse texto.

## Prazo de 2 segundos

- **Decision**: SC-002 é medida da prova, do pedido até a consulta da ocorrência, com a gravação disponível. O código não cria um timeout de 2 segundos.
- **Rationale**: A spec marca esse número como premissa de observação, não como prazo interno. Criá-lo falharia uma gravação lenta e inventaria um cancelamento.
- **Alternatives considered**: `context` com 2 segundos em toda chamada. Rejeitado pelo motivo acima.

## Fila e saúde

- **Decision**: `internal/adapter/in/queue` e `GET /health` permanecem. Este caso de uso não é ligado ao consumo da fila.
- **Rationale**: A spec manda não ler a fila de novo e não confirmar a entrada. Ligar os dois neste plano confirmaria a entrada ou leria o stream, que ficou nas specs anteriores.
- **Alternatives considered**: Classificar dentro de `consumir`. Rejeitado pelo recorte acima.
