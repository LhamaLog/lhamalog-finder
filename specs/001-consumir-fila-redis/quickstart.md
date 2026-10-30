# Quickstart: Consumir fila Redis

Guia de prova do comportamento já decidido. Não contém implementação. Os cenários com Redis real não foram executados e não podem ser executados até o solicitante informar os nomes de stream, grupo e campos.

Detalhe de campos e transições: [data-model.md](./data-model.md). Operações das portas: [contracts/consumo-fila.md](./contracts/consumo-fila.md).

## Pré-requisitos já presentes no repositório

- Docker, como no README, para `make test`
- Go 1.27.0 na imagem `gotools`, se o make não estiver instalado
- A spec `specs/001-consumir-fila-redis/spec.md`

Não há, nesta sessão, nome de fila, grupo ou campo para montar uma mensagem real.

## Prova que pode ser escrita antes do contrato de nomes

Testes de domínio e de caso de uso, com substitutos das três portas de saída. Comando já documentado no README, a partir da raiz do repositório:

```bash
make test
```

Sem make:

```bash
docker compose run --rm --no-deps gotools test ./...
```

O README equivale esse alvo a `go test -race ./...` no container Go 1.27.0. Este quickstart não afirma que esse comando foi executado para esta feature.

Resultados esperados quando esses testes existirem:

| Cenário | Esperado |
|---------|----------|
| Item com só o código, worker e instância no ar | Uma conclusão, item fora da pendência, código disponível uma vez, OpenObserve fechado |
| Passo em andamento e worker no ar | Outro worker não recebe o mesmo item |
| Worker para no meio do passo | Item não concluído e livre, sem prazo |
| Dois itens com o mesmo código, em sequência | Uma conclusão. O segundo não fica concluído nem reservado |
| Dois itens com o mesmo código terminam juntos | Uma conclusão e zero reservas restantes. O teste não exige qual item venceu |
| Fila vazia ou falha de leitura | Zero conclusões nessa tentativa |
| Item sem código ou com conteúdo extra | Nenhuma conclusão criada por regra nova |

## Prova bloqueada

Duas instâncias reais do Finder na mesma fila Redis, incluindo parar uma instância no meio do passo e ver a outra reservar o item.

Bloqueio: FR-007. Sem stream, grupo e campos, o quickstart não monta a mensagem nem aponta o processo para a fila. O lugar do código concluído também está sem nome, então “disponível para a etapa futura” ainda não tem onde ser lido fora do substituto de teste.

Quando esses nomes forem informados, a prova usa o mesmo `make test` se o teste de integração entrar no módulo, ou o arranque já descrito no README (`make run`) se a prova for manual. O endereço Redis do Compose local é o serviço `redis` publicado em `6379`, fato do `docker-compose.yml`. Isso não autoriza um nome de stream.
