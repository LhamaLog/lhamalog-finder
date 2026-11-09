# Research: Recuperar logs do lote

Decisões desta feature. Cada uma tem fonte na spec, no código lido em 2026-11-08, na documentação pública do OpenObserve lida nesta sessão ou na resposta do solicitante. Nenhuma preenche URL, credencial, timeout, tamanho de página, nome de campo ou expressão de ordenação que a spec deixou externos.

Nenhuma consulta foi enviada ao OpenObserve local durante esta escrita.

## Cliente HTTP

- **Decision**: O adaptador usa `net/http` da biblioteca padrão. Não entra SDK nem cliente de terceiros. O pacote `internal/adapter/out/openobserve` deixa de ser só a reserva de `doc.go` e passa a executar a consulta. `go.mod` não ganha dependência.
- **Rationale**: `go.mod` hoje exige `godotenv` e `go-redis`. O README descreve o módulo em Go com arquitetura hexagonal. A spec 003 já tinha recusado um SDK enquanto o contrato não existia. O contrato agora está na spec 005.
- **Alternatives considered**: Adotar um SDK do OpenObserve. Rejeitado porque não está no `go.mod` e a consulta cabe num POST. Reusar o cliente Redis. Rejeitado porque o stream Redis não é o stream do lote.

## Corpo da consulta

- **Decision**: `POST {FINDER_OO_URL}/api/{organization}/_search`. O corpo leva `query.sql`, `query.start_time`, `query.end_time`, `query.from` e `query.size`. `start_time` é `minTimestampMicros`. `end_time` é `maxTimestampMicros + 1`. `from` começa em 0 e soma `size` a cada página. `size` é o tamanho de página configurado. O SQL consulta o stream da referência, filtra `lhamalog_batch_id` e, se a referência tiver `recordId`, filtra `lhamalog_record_id`. Com os nomes de tenant e produto configurados, o SQL também filtra esses campos pelos valores da referência. A resposta usada é o arranjo `hits`.
- **Rationale**: FR-003 e FR-004. A página [Search API](https://openobserve.ai/docs/reference/api/search/search/), lida em 2026-11-08, descreve esse POST, os cinco campos do `query` e `hits` como as linhas devolvidas. A imagem local é `openobserve/openobserve:v1.1.2`, em `docker-compose.yml`. Esta sessão não confrontou essa imagem com a página.
- **Alternatives considered**: Mandar o tempo só dentro do SQL e omitir `start_time` e `end_time`. Rejeitado porque a spec e a página pedem a faixa na consulta. Usar `size` igual a `recordCount`. Rejeitado por FR-004. Acrescentar `search_type` ou amostragem. Rejeitado porque a spec não pede.

## Valores dentro do SQL

- **Decision**: Identificador de stream e nome de campo entram entre aspas duplas. Valor de texto entra entre aspas simples, com a aspas simples interna repetida. Inteiro entra sem aspas. Valor ou nome que não possa ser aplicado assim produz resposta inválida, sem consulta ou sem tratar a resposta como sucesso, e marca efeito concluído. A organização entra só no caminho do POST.
- **Rationale**: FR-003. A mesma página de busca diz que valor de texto usa aspas simples, valor inteiro não usa aspas e nome de campo pode usar aspas duplas. A spec proíbe emendá-los de modo que troquem o lote, a organização, o stream ou a faixa.
- **Alternatives considered**: Concatenar o valor cru. Rejeitado pela spec. Recusar todo valor que contenha espaço. Rejeitado porque a spec só recusa o valor que não possa ser dado do filtro.

## Ordenação

- **Decision**: A expressão de `ORDER BY` vem de `FINDER_OO_ORDEM`. Sem ela, a recuperação não começa e não marca efeito concluído. O código não grava `_timestamp` nem `lhamalog_record_id` como ordem. A prova fornece a expressão e mostra que uma cópia física não esconde outro identificador lógico.
- **Rationale**: FR-005 exige ordem estável e a spec não nomeou a coluna. A página de busca diz que a ordem padrão é `_timestamp` e que `ORDER BY` a substitui. Adotar `_timestamp` no código escolheria a coluna omitida. Cópias do mesmo instante podem empatar nessa coluna.
- **Alternatives considered**: Confiar na ordem padrão `_timestamp`. Rejeitado pelo motivo acima. Inventar ordem por `lhamalog_record_id` mais um desempate físico. Rejeitado porque o desempate não foi nomeado. Criar um máximo de páginas para compensar ordem instável. Rejeitado pela resposta B.

## Paginação sem teto

- **Decision**: A leitura pede outra página enquanto `hits` tiver o mesmo comprimento de `size`. A página menor encerra. Não existe máximo de páginas no código nem na configuração. Se uma consulta estoura o timeout do cliente ou o armazenamento não responde, a classificação é indisponibilidade ou timeout, salvo se uma falha que pede correção já tiver sido encontrada.
- **Rationale**: FR-004 e a resposta do solicitante em 2026-11-08, opção B. FR-015 e FR-016 dizem quando o timeout deixa de ser a classificação emitida.
- **Alternatives considered**: Parar em `recordCount` páginas. Rejeitado porque cópias físicas podem ocupar essas páginas e deixar um identificador lógico de fora. Tratar página cheia eterna como efeito concluído. Rejeitado porque a resposta B não cria esse desfecho.

## Quando parar depois de uma falha de documento

- **Decision**: Autenticação recusada ou resposta inválida encerra a leitura. Isolamento, documento sem identificador, conteúdo divergente e contagem ainda não fecham a leitura sozinhos, porque uma página seguinte pode trazer resposta inválida, que tem precedência maior, ou isolamento, que precede as falhas de documento. Se a página seguinte falha e uma falha que pede correção já foi vista, a classificação emitida é essa falha, na ordem de FR-015, e os documentos lidos não saem como lote pronto.
- **Rationale**: A ordem de FR-015 e o cenário 9 da história 5. Parar no primeiro documento sem identificador esconderia um documento de outro tenant numa página posterior.
- **Alternatives considered**: Parar no primeiro problema de documento. Rejeitado pela precedência. Continuar depois de resposta inválida. Rejeitado porque nada na ordem restante a ultrapassa, exceto uma falha de autenticação, que já encerra a chamada.

## Deduplicação e cópia física

- **Decision**: A chave lógica é `lhamalog_record_id`. `FINDER_OO_CAMPOS_COPIA` é a lista. Variável ausente: a recuperação não começa e não marca efeito concluído. Variável presente e vazia: qualquer diferença é conteúdo divergente. Diferença só nos campos da lista: um registro lógico, e o entregue é uma das cópias. Documento sem o identificador não entra na contagem. O identificador Redis e `correlationId` não deduplicam.
- **Rationale**: FR-002, FR-006, FR-007 e FR-008, e a resposta A de 2026-11-08 sobre cópia e divergência. A lista não é nomeada na spec.
- **Alternatives considered**: Tratar qualquer diferença como divergência e dispensar a lista. Rejeitado pela resposta A. Escolher sempre o primeiro documento do mesmo identificador. Rejeitado porque a spec proíbe escolher versão quando o conteúdo diverge.

## Isolamento

- **Decision**: `FINDER_OO_CAMPO_TENANT` e `FINDER_OO_CAMPO_PRODUTO` são os nomes persistidos. Os dois presentes: o SQL filtra pelos valores da referência e cada `hit` é conferido de novo. Falta um dos nomes: a consulta não é enviada, nenhum documento sai como válido e a classificação é isolamento, com efeito concluído. Credencial não substitui filtro nem conferência. Atributo de dentro do log não substitui `tenantId` nem `productId` da referência.
- **Rationale**: FR-012, FR-013 e a resposta A de 2026-11-08. Os nomes continuam fora do código. O exemplo SQL do contrato, só com `batchId`, não basta.
- **Alternatives considered**: Confiar na credencial sozinha. Rejeitado pela resposta A. Tratar nome ausente como configuração que impede o início, igual à URL. Rejeitado porque a spec marca esse caso como isolamento com efeito concluído, não como recuperação que não começou.

## Efeito concluído

- **Decision**: O caso de uso devolve uma classificação e o marca de efeito concluído. Efeito concluído: lote completo e a falha que pede correção escolhida pela precedência (autenticação ou autorização, resposta inválida, isolamento, documento sem identificador, conteúdo divergente, excesso). Sem efeito concluído: indisponibilidade, timeout, falha de página e lote incompleto, quando nenhuma falha que pede correção ocorreu; também URL, credencial, timeout, tamanho de página, lista de cópia ou expressão de ordenação ausentes. Esta etapa não confirma a entrada Redis. Não há segunda consulta `COUNT(DISTINCT ...)`.
- **Rationale**: FR-015, FR-018 e FR-019, e as respostas B de 2026-11-08 sobre o efeito e A sobre a precedência. FR-011 permite a contagem distinta só como conferência. O sucesso precisa dos registros, então a contagem não substitui a paginação.
- **Alternatives considered**: Confirmar a entrada neste adaptador. Rejeitado por FR-018 e pelo contrato da feature 004. Usar `COUNT(DISTINCT lhamalog_record_id)` no lugar da leitura. Rejeitado por FR-011.

## Credencial e timeout

- **Decision**: Usuário e senha vêm do ambiente e saem em HTTP Basic no cabeçalho `Authorization`. Não têm padrão no código, não entram em log e não entram no resultado de falha. O timeout é o prazo do cliente HTTP para cada POST. Não é um prazo total da paginação e não é o campo `timeout` do corpo, que esta feature não envia.
- **Rationale**: FR-002 e FR-013. A página de busca lida não descreve o cabeçalho; descreve um `timeout` opcional no corpo, com padrão do servidor. O Finder precisa de um prazo por consulta, que a spec manda configurar, e de não gravar a senha. HTTP Basic é a forma desta decisão. A imagem local, no trecho lido de `docker-compose.yml`, não define usuário nem senha. O README diz que este módulo não os define.
- **Alternatives considered**: Gravar a senha padrão da imagem no Finder. Rejeitado pelo README e por FR-002. Enviar também `timeout` no corpo com o mesmo número. Rejeitado porque a spec não pede esse campo e o prazo do cliente já encerra a chamada.

## Prova

- **Decision**: `make test` cobre domínio e caso de uso com substituto da porta de busca. A prova contra o OpenObserve do Compose só roda com `FINDER_OPENOBSERVE_INTEGRATION=1`, dados sintéticos e os parâmetros informados no quickstart. Ela não usa log de produção e não declara homologação.
- **Rationale**: FR-020. `make test` já é `go test -race ./...` sem serviço externo, como no README e na feature 004. O pacote OpenObserve hoje não consulta.
- **Alternatives considered**: Exigir o OpenObserve em `make test`. Rejeitado porque o alvo atual não sobe o armazenamento. Tratar a prova sintética como homologação dos consumidores. Rejeitado por FR-020.
