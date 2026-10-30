# Research: Consumir fila Redis

Leitura usada em 2026-10-29: `specs/001-consumir-fila-redis/spec.md`, `.specify/memory/constitution.md`, `README.md`, `go.mod`, `internal/adapter/in/queue/redisStream.go`, `internal/config/config.go`, `cmd/finder/main.go`. Nenhum teste desta feature foi executado nesta sessão.

## 1. Recorte da etapa

- **Decision**: O código desta etapa só reserva, conclui ou libera o item da fila e deixa o código de apontamento disponível uma vez. Não abre o OpenObserve, não agrega log e não classifica ocorrência.
- **Rationale**: Pedido do solicitante e FR-005. A constituição mantém o fluxo completo para uma etapa posterior.
- **Alternatives considered**: Implementar já a leitura do log e as classes erro, alerta e informação. Rejeitada porque a spec marca esse fluxo como fora desta etapa.

## 2. Protótipo atual

- **Decision**: `InitRedis` em `internal/adapter/in/queue/redisStream.go` não é a implementação pedida. Ele lê o stream `fila:logs`, o campo `id`, a partir do id `110`, e apaga a mensagem. O `main` hoje chama essa função. A substituição dessa chamada faz parte do consumo, e só ocorre quando o contrato de nomes existir.
- **Rationale**: Q2 A e FR-007. Apagar na leitura não é o passo único de reserva e conclusão, e não libera o item se o processo parar no meio.
- **Alternatives considered**: Evoluir o `XRead` + `XDel` do protótipo. Rejeitada pelo solicitante.

## 3. Nomes da fila

- **Decision**: Nenhum nome de stream, grupo ou campo é escolhido. O adaptador Redis permanece bloqueado.
- **Rationale**: FR-007 e a resposta Q2 A. A constituição, princípio V, proíbe criar esse contrato.
- **Alternatives considered**: Adotar `fila:logs`, campo `id` e id inicial `110`. Rejeitada pelo solicitante.

## 4. Onde o código concluído fica

- **Decision**: O caso de uso chama uma porta de saída “código disponível para a etapa futura”. Este plano não cria adaptador com nome de armazenamento, coleção, chave ou índice.
- **Rationale**: Esclarecimento de 2026-10-29: o código continua disponível e o lugar não é nomeado. Também não é o OpenObserve.
- **Alternatives considered**: Memória só do processo. Rejeitada como desenho final porque a etapa futura e as outras instâncias precisam ver uma única disponibilidade, e o solicitante não escolheu “só até a instância parar”. Nomear Redis, MongoDB ou Elasticsearch para esse fim. Rejeitado porque seria inventar o lugar.

## 5. Um código, uma conclusão

- **Decision**: A regra de domínio é uma conclusão por código de apontamento. O segundo item não é concluído e não permanece reservado. Se dois itens com o mesmo código terminam no mesmo momento, o resultado visível é uma conclusão e zero reservas restantes. Qual item é o concluído permanece indefinido.
- **Rationale**: FR-009 e SC-005, resposta B do solicitante.
- **Alternatives considered**: Concluir os dois itens. Rejeitada pelo solicitante. Escolher vencedor por ordem do identificador ou por relógio. Rejeitada porque a spec diz que o vencedor do empate não está definido.

## 6. Liberar sem prazo numérico

- **Decision**: Parada de worker com a instância no ar é observada por essa instância, que devolve o item à condição livre e não grava conclusão. Parada da instância inteira, vista por outra instância, sem prazo numérico, não tem mecanismo escolhido. O plano não introduz tempo de espera, batimento nem prazo de reivindicação.
- **Rationale**: FR-006 exige a liberação e proíbe prazo numérico. Dentro do processo, a parada do worker pode ser vista pelo próprio processo. Entre processos, este repositório não documenta um sinal de morte da instância que dispense prazo. Escolher um tempo inventaria o dado que a spec deixou de fora.
- **Alternatives considered**: Reivindicar o item depois de um tempo ocioso. Rejeitada porque cria prazo numérico. Apagar o item ao ler, como o protótipo. Rejeitada porque o item pararia de existir para a outra instância.

## 7. Dependência Redis

- **Decision**: Nenhuma biblioteca nova é adicionada por este plano. O import já presente é `github.com/redis/go-redis/v9 v9.23.0`, indireto em `go.mod`. O README ainda o lista como candidato. A conferência no Drive da equipe não foi executada aqui.
- **Rationale**: FR-008 e a constituição, Restrições Técnicas. O plano não trata o candidato do README como aprovação, nem remove o import que o protótipo já tem.
- **Alternatives considered**: Trocar de cliente Redis ou acrescentar outra biblioteca de concorrência. Rejeitada porque a tarefa ainda não tem o contrato de nomes que justificaria uma dependência nova.

## 8. Paralelismo

- **Decision**: `FINDER_WORKERS` não define o número de workers deste consumo.
- **Rationale**: A spec registra o padrão 4 e afirma que ele não é o grau de paralelismo desta etapa.
- **Alternatives considered**: Usar 4 workers porque o config já tem esse padrão. Rejeitada como suposição.

## 9. Ordem dos itens

- **Decision**: Não há requisito de ordem global. Outra instância pode receber outro item livre, ou nenhum.
- **Rationale**: Cenário 2 da história 2. A spec não pede ordem de entrada.
- **Alternatives considered**: Consumo estritamente na ordem de chegada com várias instâncias. Não foi pedido e apertaria o desenho de várias instâncias sem um requisito de aceite.

## 10. Prova

- **Decision**: A prova automatizada desta etapa, enquanto o contrato de nomes não existe, é teste de domínio e de caso de uso com substitutos das portas. A prova com duas instâncias e Redis real entra no quickstart como cenário bloqueado.
- **Rationale**: O README já define `go test -race ./...`. A FR-007 impede apontar o teste de integração para um stream nomeado.
- **Alternatives considered**: Declarar o teste de integração como executado. Rejeitada porque ele não foi executado e não tem nomes de fila.
