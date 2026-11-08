# Research: Consumir logs no OpenObserve

Leitura usada em 2026-10-31: `specs/002-consumir-logs-openobserve/spec.md`, `.specify/memory/constitution.md`, `README.md`, `go.mod`, `docker-compose.yml`, `Makefile`, `cmd/finder/main.go`, `internal/config/config.go`, `internal/application/consumir.go`, `internal/domain/item.go`. Nenhum teste desta feature foi executado nesta sessão. O arquivo `otel-logs.zip` não estava no repositório nem na pasta LhamaLog Workspace.

## 1. Recorte da etapa

- **Decision**: O caso de uso novo obtém apontamentos da fonte mockada, captura um log por apontamento e o deixa pronto. Não lê a fila, não agrega, não classifica ocorrência e não oculta dados sensíveis.
- **Rationale**: FR-001, FR-004 e as respostas de 2026-10-31. A constituição mantém o fluxo completo para uma etapa posterior. O caso de uso `application.Servico` permanece o consumo da fila.
- **Alternatives considered**: Estender `Servico.Executar` para abrir o OpenObserve depois de concluir o item. Rejeitada porque esta etapa não lê a fila e a feature 001 fica fora do escopo.

## 2. Um apontamento, um log

- **Decision**: Cada apontamento seleciona um único log. Vários logs exigem vários apontamentos. A captura não pede o exemplo inteiro nem um conjunto de tamanho indefinido.
- **Rationale**: Resposta A de 2026-10-31 e FR-002.
- **Alternatives considered**: Uma indicação trazer todos os registros do arquivo. Rejeitada pelo solicitante.

## 3. Log pronto

- **Decision**: O log pronto é o registro armazenado, inteiro, associado ao apontamento. Enquanto a instância está no ar, uma nova leitura por esse apontamento devolve o mesmo registro. Depois que a instância para, esta etapa não exige que ele continue. O lugar não recebe nome de banco, arquivo ou coleção. A memória do processo satisfaz essa regra sem nomear armazenamento.
- **Rationale**: FR-003 e a resposta A de 2026-10-31. A constituição não define persistência da ocorrência.
- **Alternatives considered**: Gravar no MongoDB ou no Redis. Rejeitada porque nomearia o lugar. Exigir o log depois que o processo para. Rejeitada pelo solicitante.

## 4. Apontamento repetido

- **Decision**: Indicações repetidas do mesmo apontamento na mesma execução viram uma captura. O resultado é um log pronto. A leitura devolve esse único registro.
- **Rationale**: FR-011, SC-006 e a resposta A de 2026-10-31.
- **Alternatives considered**: Um log pronto por aparição. Rejeitada pelo solicitante.

## 5. Falha

- **Decision**: Três motivos são distintos: apontamento inutilizável, log ausente e OpenObserve inacessível. A falha carrega o apontamento e o motivo. O log não fica pronto. Não há nova tentativa, prazo nem fila de erro. Apontamento vazio é inutilizável, no mesmo critério já usado por `Item.Aceito` para código vazio. Outros formatos inválidos não são criados.
- **Rationale**: FR-006 e o domínio existente, em que código vazio não é aceito. A spec não definiu alfabeto nem formato.
- **Alternatives considered**: Tratar apontamento vazio como log ausente. Rejeitada porque a spec separa os dois motivos. Inventar código HTTP ou corpo de erro. Rejeitada pelo princípio V.

## 6. Paralelismo

- **Decision**: Apontamentos distintos da mesma execução podem ser capturados ao mesmo tempo, com a biblioteca padrão. `FINDER_WORKERS` não define a quantidade. Não há meta de vazão. Coordenação entre várias instâncias do Finder fica fora.
- **Rationale**: FR-007 e SC-003. O README lista `golang.org/x/sync` como candidata ainda não adicionada.
- **Alternatives considered**: Adotar o padrão 4 de `FINDER_WORKERS`. Rejeitada porque a spec não liga esse número a esta captura. Acrescentar `golang.org/x/sync`. Rejeitada porque a biblioteca padrão cobre a espera das capturas.

## 7. OpenObserve

- **Decision**: Nenhum nome de organização, stream, índice, campo, URL ou credencial é escolhido. O adaptador não é escrito. Quando o contrato externo existir, o cliente HTTP da biblioteca padrão é o candidato. Nenhuma biblioteca nova entra neste plano.
- **Rationale**: FR-005, FR-010 e o princípio V. O `go.mod` não contém cliente de OpenObserve.
- **Alternatives considered**: Escolher um SDK ou um caminho de API público. Rejeitada porque a spec proíbe nomear esse contrato e o princípio V proíbe criá-lo.

## 8. Base local

- **Decision**: O resultado pedido é a base sem Elasticsearch e com OpenObserve, mantendo MongoDB e Redis. A edição de `docker-compose.yml`, `Makefile` e README fica bloqueada. Versão, porta e credencial não foram informadas. Migração de dado já existente no Elasticsearch fica fora. O pacote `internal/adapter/out/elastic` não é apagado.
- **Rationale**: FR-008 e FR-010. `docker-compose.yml` hoje usa `docker.elastic.co/elasticsearch/elasticsearch:9.4.5`. Copiar um número de versão ou uma porta de outra fonte inventaria dado ausente no pedido.
- **Alternatives considered**: Publicar uma imagem e uma porta correntes do OpenObserve. Rejeitada pelo princípio V.

## 9. Arquivo de exemplo

- **Decision**: O aceite SC-001 só vale com um registro de `otel-logs.zip` já armazenado no OpenObserve. O Finder não lê o arquivo. O arquivo não estava disponível em 2026-10-31. Os passos que colocam o registro no OpenObserve não são definidos. Os campos internos não são descritos. Testes com substitutos usam um registro opaco e não contam como SC-001.
- **Rationale**: FR-009 e a resposta B de 2026-10-31.
- **Alternatives considered**: Aceitar qualquer log conhecido. Rejeitada pelo solicitante. Fazer a base local carregar o arquivo. Rejeitada porque o solicitante não escolheu essa opção e o arquivo não está no repositório.

## 10. Igualdade do registro

- **Decision**: O log pronto é igual ao registro armazenado. Esta etapa não remove nem reescreve campos. A contagem de campos ocultos no aceite é zero. A lista do que seria sensível não é criada.
- **Rationale**: Resposta A de 2026-10-31 e FR-004. O README descreve mascaramento no módulo. Esta etapa não o aplica e não emenda o README.
- **Alternatives considered**: Ocultar dados sensíveis na preparação. Rejeitada pelo solicitante.

## 11. Prova

- **Decision**: Enquanto o arquivo e o contrato do OpenObserve não existem, a prova automatizada é teste de domínio e de caso de uso com substitutos. O quickstart marca a prova com OpenObserve real e com o arquivo como bloqueada. Nenhum comando de teste foi executado nesta sessão.
- **Rationale**: O README já define `go test -race ./...`. A FR-009 impede tratar um registro inventado como o aceite.
- **Alternatives considered**: Declarar SC-001 executado com um JSON de exemplo. Rejeitada pelo princípio V e pela resposta B.
