# Quickstart: Capturar referência de lote

Prova do comportamento de [contracts/referencia-lote.md](./contracts/referencia-lote.md) e [data-model.md](./data-model.md). Não substitui o código de produção. Os números e os nomes de grupo, consumidor e prefixo abaixo são dados desta prova. O processo não os traz como padrão.

Esta sessão de planejamento não executou os comandos. Nenhum resultado de teste é afirmado aqui.

## Pré-requisitos

- Docker Compose na raiz do repositório.
- Redis do Compose, serviço `redis`, contêiner `lhamalog-finder-redis`, publicado em `localhost:6379`.
- Dados desta prova:

```text
FINDER_REDIS_ADDR=127.0.0.1:6379
FINDER_REDIS_STREAM=lhamalog:stream:logs
FINDER_REDIS_GROUP=finder
FINDER_REDIS_CONSUMIDOR=finder-prova
FINDER_REDIS_IDEMPOTENCIA=finder:log-ingested:
FINDER_REDIS_COUNT=2
FINDER_REDIS_BLOCK=200ms
FINDER_REDIS_PENDENCIA_OCIOSA=1s
```

`fila:logs` e `finder:codigos` não são a entrada desta prova. A senha fica de fora porque o Redis do Compose sobe sem senha.

## 1. Prova sem Redis

Na raiz:

```text
make test
```

Resultado esperado: `go test -race ./...` termina sem falha. Testes que exigem Redis não rodam neste alvo. A prova de domínio cobre envelope válido, JSON inválido, versão diferente de `"1"`, signal diferente de `LOGS`, divergência de produto, de batchId e de faixa, lote de um registro e de vários, metadata com e sem `traceId`, e duas entradas com o mesmo `idempotencyId`.

## 2. Arranque sem parâmetro

Repetir o arranque três vezes, cada uma sem um dos três: quantidade, bloqueio ou ociosidade.

Resultado esperado: o consumo não abre. Nenhuma entrada é lida. Nenhum desses três valores aparece criado pelo processo.

## 3. Prova com Redis

Subir o Redis:

```text
docker compose up -d redis
```

Esperar o healthcheck do serviço `redis`. A prova usa a rede do contêiner e os dados da seção de pré-requisitos, mais `FINDER_REDIS_INTEGRATION=1`. O substituto da porta de efeito informa conclusão ou falha. A prova não chama o OpenObserve.

Resultado esperado, com entradas sintéticas publicadas em `lhamalog:stream:logs`:

| Cenário | Resultado |
|---------|-----------|
| Entradas publicadas antes da criação do grupo | São lidas. A criação não usa `$` |
| Envelope válido, efeito ainda não informado | Referência pronta, entrada sem confirmação, chave `em_processamento` |
| Efeito informado como concluído | Chave `concluido` com prazo de 25 horas, entrada confirmada pelo id Redis |
| Segunda entrada com a mesma chave em processamento, antes da ociosidade | Permanece pendente, sem outro efeito e sem confirmação |
| Chave já `concluido` dentro de 25 horas | A outra entrada é confirmada e o efeito não se repete |
| Queda antes do efeito | A entrada segue pendente. Depois da ociosidade informada, pode ser reassumida |
| Queda depois do efeito e antes da confirmação | O efeito não se repete e a reentrega confirma a entrada |
| Envelope inválido | Não segue como referência válida. É confirmada uma vez. O registro não tem senha nem conteúdo de log |
| Grupo `finder` e um segundo grupo na mesma entrada | Os dois grupos observam a entrada |
| Redis indisponível no meio | Entrada inacabada não é confirmada |

O prazo de 25 horas é conferido pelo tempo restante da chave concluída, não por uma espera de 25 horas.

## 4. Fora desta prova

Consulta ao OpenObserve, varredura, checked, ocorrência, `trace.ingested` e alteração do Collector.

## Processo já existente

A consulta de saúde não muda de significado. Esta prova não a redefine.
