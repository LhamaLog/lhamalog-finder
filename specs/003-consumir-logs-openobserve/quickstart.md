# Quickstart: Consumir logs no OpenObserve

Guia de validação. Não substitui `tasks.md`. Nenhum comando abaixo foi executado nesta sessão de planejamento.

Contrato: [captura-log.md](./contracts/captura-log.md). Modelo: [data-model.md](./data-model.md).

## Pré-requisitos

- Docker em execução e o alvo de teste já descrito no README.
- Go 1.27.0 na imagem `gotools`. Não é preciso instalar Go na máquina.
- Para SC-001 e SC-004: o arquivo `otel-logs.zip` entregue, um registro desse arquivo armazenado no OpenObserve, e os nomes externos de acesso. Estes três itens não estavam disponíveis em 2026-10-31. Sem eles, esses dois critérios ficam bloqueados.

## Prova que pode ser escrita agora

Testes de domínio e de caso de uso, com substitutos das portas. O registro do substituto é opaco. Ele não é um registro de `otel-logs.zip` e não cumpre SC-001.

```bash
make test
```

Equivalente sem make:

```bash
docker compose run --rm --no-deps gotools test -race ./...
```

| Cenário | Resultado esperado | Critério |
|---------|--------------------|----------|
| Um apontamento e um registro no substituto, instância no ar | Uma leitura seguinte devolve o mesmo registro. Nenhuma ocorrência. A fila não é chamada | SC-001 só na forma. O arquivo real continua obrigatório para o aceite |
| Fonte sem apontamento | Zero logs prontos e zero falhas | SC-005 |
| Apontamento vazio | Falha inutilizável. Log não pronto | SC-002 |
| Apontamento utilizável sem registro | Falha ausente, distinta da anterior | SC-002 |
| Leitura que não alcança o armazenamento | Falha inacessível, distinta da ausência | SC-002 |
| Dois apontamentos distintos ao mesmo tempo, um sem registro | O existente fica pronto. A falha do outro permanece | SC-003 |
| O mesmo apontamento duas vezes | Um log pronto | SC-006 |
| Instância para no meio | Log não pronto | borda da spec |

## Prova bloqueada

| Cenário | O que falta |
|---------|-------------|
| SC-001 com registro real | `otel-logs.zip` e o contrato de acesso ao OpenObserve |
| SC-004, base local com OpenObserve no ar | Docker em execução. A definição está em `docker-compose.yml`: imagem `openobserve/openobserve:v1.1.2`, porta 5080 |
| Captura pelo processo `cmd/finder` contra OpenObserve | O mesmo contrato de acesso. O `main` atual executa o consumo da fila e esta etapa não o reescreve |

Subir a base atual não valida esta feature:

```bash
make infra-up
```

Esse alvo sobe MongoDB, Redis e OpenObserve. Ele não coloca um registro de `otel-logs.zip` no OpenObserve.
