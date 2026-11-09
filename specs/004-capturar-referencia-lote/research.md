# Research: Capturar referência de lote

Decisões desta feature. Cada uma tem fonte na spec, no código lido em 2026-11-08 ou na resposta do solicitante. Nenhuma preenche número, grupo, consumidor ou prefixo que a spec deixou externo.

## Stream do contrato

- **Decision**: O consumo desta feature lê `lhamalog:stream:logs`. Se `FINDER_REDIS_STREAM` vier vazio ou com outro nome, o processo não abre o consumo. O valor `fila:logs`, hoje padrão em `docker-compose.yml` e `.env.example`, deixa de ser o stream deste módulo.
- **Rationale**: FR-001 nomeia o stream. O código atual aceita qualquer stream informado e o Compose aponta para `fila:logs`. Aceitar outro nome leria a fila antiga.
- **Alternatives considered**: Manter `fila:logs` como padrão. Rejeitado porque a spec substitui esse nome. Tratar stream vazio como o nome do contrato dentro do código, e ao mesmo tempo aceitar outro valor no ambiente. Rejeitado porque dois streams possíveis escondem o contrato.

## Posição inicial do grupo

- **Decision**: A criação do grupo usa o id `0`, que já está em `XGroupCreateMkStream` no adaptador atual. A posição `$` não é usada. Grupo que já existe (`BUSYGROUP`) não é apagado nem recriado.
- **Rationale**: A resposta de 2026-11-08, opção B, manda entregar as entradas que ainda estão no stream. O id `0` é o que o Redis usa para isso e já é o argumento do código. Apagar um grupo existente para refazê-lo descartaria pendências.
- **Alternatives considered**: `$`, que ignora o que já está no stream. Rejeitado pela resposta B. Recriar o grupo quando ele já nasceu em `$`. Rejeitado porque destruiria o estado do grupo sem pedido.

## Parâmetros de leitura e de pendência

- **Decision**: Três dados externos, sem padrão no código: quantidade de entradas por leitura, tempo de bloqueio da leitura de entradas novas, e tempo de ociosidade para uma pendência ser reassumida. Se algum falta, não é número ou não é positivo, o processo não abre o consumo. A quantidade também limita quantas pendências abandonadas são reivindicadas em uma passagem. A leitura das pendências já deste consumidor não usa o tempo de bloqueio, para não esperar antes de observar o que já está pendente.
- **Rationale**: FR-005 e a resposta de 2026-11-08, opção A, proíbem padrão numérico. O código atual fixa `Count: 1`, bloqueio `-1` ou `0` e ociosidade `0`. Esses literais deixam de ser a regra. Reusar a quantidade na reivindicação evita um quarto número.
- **Alternatives considered**: Conservar `Count: 1` e ociosidade zero. Rejeitado porque são padrões inventados. Exigir uma quantidade diferente para a reivindicação. Rejeitado pelo mesmo motivo.

## O que conta como pendência abandonada

- **Decision**: Uma pendência pode ser reassumida quando a ociosidade configurada foi atingida. A regra atual, que reivindica com ociosidade zero a mensagem de um consumidor sem conexão em `CLIENT LIST`, deixa de ser a regra deste consumo.
- **Rationale**: A spec manda usar o critério informado. A ociosidade zero mais a lista de conexões é outro critério, e o código a descreve como "sem prazo numérico". Esta feature tem prazo informado.
- **Alternatives considered**: Manter as duas regras juntas. Rejeitado porque a conexão viva passaria a ser exigência que a spec não tem. Esperar 25 horas para reassumir. Rejeitado porque esse prazo é só do estado concluído.

## Chave do lote

- **Decision**: Uma chave por `idempotencyId`, sob o prefixo externo `FINDER_REDIS_IDEMPOTENCIA`. Sem esse prefixo, o processo não abre o consumo. O valor é `em_processamento` ou `concluido`. A criação só ocorre se a chave não existe. `em_processamento` não recebe prazo. A passagem para `concluido` grava o valor e o prazo de 25 horas, sem renovar esse prazo se a chave já estava concluída. Entrada inválida remove a chave só quando ela está em processamento e em seguida confirma a entrada. O conjunto `FINDER_REDIS_CODIGOS` não guarda esse estado.
- **Rationale**: FR-016 e a resposta Q3 B. O conjunto atual marca um código para sempre e confirma na mesma operação. Isso trata a reserva como conclusão e não separa o Finder da Varredura por um prefixo desta feature.
- **Alternatives considered**: Reusar `finder:codigos`. Rejeitado porque é um conjunto, sem os dois estados e sem 25 horas. Escolher o prefixo no código. Rejeitado porque a spec não nomeou o namespace. Colocar 25 horas já na reserva. Rejeitado porque a queda no meio precisa poder tentar o efeito de novo.

## Segunda entrada com a chave em processamento

- **Decision**: A entrada fica pendente, sem segundo efeito e sem confirmação, até a chave sair de `em_processamento`. O leitor não a confirma e não a trata de novo na mesma passagem, para não girar em claro. Quando a chave está `concluido`, a entrada pendente é confirmada sem repetir o efeito.
- **Rationale**: Resposta de 2026-11-08, opção A, e FR-014. Confirmar antes do efeito perderia a entrada se a primeira tentativa caísse. Executar em paralelo repetiria o efeito.
- **Alternatives considered**: Confirmar a segunda na hora. Rejeitado pela opção A. Devolver a entrada para o fim do stream. Rejeitado porque a spec não cria reenvio nem outra fila.

## Efeito da recuperação

- **Decision**: O caso de uso chama uma porta com a referência validada e o id da entrada. Só um resultado de efeito concluído autoriza gravar `concluido` e confirmar a entrada. Esta feature não implementa a porta com consulta ao OpenObserve. Porta ausente ou falha antes desse resultado mantém a entrada pendente e a chave em processamento. O teste informa o resultado por um substituto.
- **Rationale**: Resposta síncrona de 2026-11-08. A spec proíbe confirmar na interpretação e proíbe criar fila intermediária. Um substituto que sempre conclui, ligado no processo real, confirmaria sem recuperação.
- **Alternatives considered**: Confirmar ao entregar a referência em memória. Rejeitado. Gravar a referência em banco novo e confirmar. Rejeitado porque seria infraestrutura intermediária.

## Confirmação, autenticação e registro

- **Decision**: A confirmação usa o id da entrada Redis. A senha, quando o ambiente a tiver, vem de configuração e não entra em registro. Registro de inválida leva o id da entrada, o `idempotencyId` quando ele pôde ser lido, e o motivo. Não leva senha, metadata bruta, body bruto nem conteúdo de log. Grupo e consumidor continuam externos; sem eles o processo não abre o consumo. O gerador atual `hostname-pid` deixa de suprir o consumidor ausente.
- **Rationale**: FR-002, FR-014 e FR-015. O gerador atual cria um nome que a spec não informou. O cliente Redis atual não tem campo de senha; sem o campo externo a autenticação não pode ser configurada. Senha vazia continua válida para o Redis local, que hoje sobe sem senha.
- **Alternatives considered**: Manter `hostname-pid` como padrão. Rejeitado por FR-002. Exigir senha sempre. Rejeitado porque o Redis do Compose não tem senha e a spec não exige que ela exista, só que não esteja no código.

## Biblioteca e corte do stream

- **Decision**: Nenhuma dependência nova. O corte aproximado de 100000 entradas é fato do publicador. O Finder não o altera e não reconstrói entrada já removida. O prazo de 25 horas da chave concluída não é lido desse corte.
- **Rationale**: `go.mod` já exige `go-redis`. FR-020 separa quantidade de retenção.
- **Alternatives considered**: Trocar de cliente Redis. Rejeitado porque o cliente atual já faz leitura em grupo, confirmação e reivindicação.
