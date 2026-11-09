# Implementation Plan: Capturar referência de lote

**Branch**: `004-capturar-referencia-lote` | **Date**: 2026-11-08 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/004-capturar-referencia-lote/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

O Finder deixa de tratar a entrada da fila como um código único e passa a ler o stream `lhamalog:stream:logs`, interpretar o envelope `log.ingested` v1 e entregar a referência validada a uma porta da etapa de recuperação. A entrada válida só é confirmada depois que essa porta informa o efeito concluído. A entrada inválida é confirmada na hora, com o motivo registrado, sem estado concluído e sem nova tentativa. A idempotência é uma chave por `idempotencyId`, em processamento ou concluída. Só o estado concluído dura 25 horas. Outra entrada válida permanece pendente enquanto a chave está em processamento e o tempo de ociosidade configurado ainda não autoriza a reassunção.

O grupo nasce no início do que o stream ainda guarda (`0`, já usado em `internal/adapter/in/queue/redisStream.go`). A posição `$` não é usada. Quantidade por leitura, tempo de espera e ociosidade da pendência vêm da configuração. Se algum falta, o processo não abre o consumo e nenhum padrão numérico é criado. A busca no OpenObserve não entra.

A pesquisa está em [research.md](./research.md).

## Technical Context

**Language/Version**: Go 1.27, como em `go.mod` e no README

**Primary Dependencies**: `github.com/redis/go-redis/v9 v9.23.0`, já exigida em `go.mod`. Nenhuma biblioteca nova.

**Storage**: Redis Stream `lhamalog:stream:logs` e uma chave Redis por `idempotencyId`, no prefixo informado pelo ambiente. O conjunto `FINDER_REDIS_CODIGOS` não é o modelo desta feature. Não há fila intermediária, spool nem DLQ.

**Testing**: `make test` executa `go test -race ./...` sem Redis. A prova com Redis fica atrás de `FINDER_REDIS_INTEGRATION=1`, no mesmo critério já usado em `internal/adapter/in/queue/redis_integration_test.go`. Os números dessa prova são dados da execução, não padrão do processo.

**Target Platform**: Processo do Finder no Docker Compose, com o Redis já publicado. A consulta de saúde existente permanece. Esta feature não grava ocorrência e não chama o OpenObserve.

**Project Type**: Serviço Go de um módulo, entrada de fila e health HTTP já existente.

**Performance Goals**: Sem meta numérica. O aceite é o da spec: referência válida sem confirmação ao fim da interpretação (SC-001), lote unitário e múltiplo (SC-002), inválida confirmada uma vez (SC-003), um efeito dentro de 25 horas e segunda entrada pendente enquanto a chave está em processamento (SC-004), leitura do que já está no stream e grupo separado da Varredura (SC-005), Redis indisponível sem confirmar entrada inacabada (SC-006), pendência abandonada recuperada (SC-007), registro sem credencial e sem conteúdo de log (SC-008), consumo que não inicia sem os três parâmetros (SC-009).

**Constraints**: Não criar quantidade, tempo de espera, ociosidade, nome de grupo, nome de consumidor nem prefixo de idempotência no código. Não usar `$` na criação do grupo. Não confirmar entrada válida antes do efeito informado pela porta de recuperação. Não consultar o OpenObserve. Não registrar senha nem corpo de log. Não recriar entrada removida pelo corte aproximado de 100000. `FINDER_WORKERS` continua existindo e não é o grau de paralelismo deste consumo.

**Scale/Scope**: Duas ou mais instâncias do Finder no mesmo grupo, e um grupo distinto no papel da Varredura. O corte do Collector limita a quantidade aproximada de entradas; não é retenção de 25 horas.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Princípio | Resultado | Leitura |
|-----------|-----------|---------|
| I. Finalidade: logs em ocorrências | Justificado | Esta etapa não classifica ocorrência. A spec limita o trabalho à referência do lote. |
| II. Origem dos logs no coletor | Justificado | A entrada vem do Collector e aponta o lote. O conteúdo continua no OpenObserve. |
| III. Categorias de ocorrência | Justificado | Erro, alerta e informação ficam fora. Nenhuma classe nova. |
| IV. Fila Redis e apontamento no OpenObserve | Justificado | A fila entrega a referência, não o log. A confirmação válida espera o efeito da recuperação. A forma da busca fica na etapa seguinte. |
| V. Incerteza sem invenção | Passa | Stream, versão, campos, 25 horas, posição inicial e a recusa sem parâmetros vêm da spec e das respostas de 2026-11-08. Grupo, consumidor, prefixo e os três números continuam externos. O plano não os preenche. |

**Resultado do gate**: passa para o desenho. A implementação espera `/speckit-tasks`.

### Revisão depois do desenho

[research.md](./research.md), [data-model.md](./data-model.md), [contracts/referencia-lote.md](./contracts/referencia-lote.md) e [quickstart.md](./quickstart.md) repetem o stream do contrato, a chave de dois estados e a porta de efeito sem consulta ao OpenObserve. Os números e os nomes de grupo, consumidor e prefixo aparecem só como dados que a prova informa. Nenhum padrão numérico foi colocado no processo. O gate permanece aprovado.

## Project Structure

### Documentation (this feature)

```text
specs/004-capturar-referencia-lote/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── referencia-lote.md
└── tasks.md             # criado por /speckit-tasks, não por este comando
```

### Source Code (repository root)

```text
cmd/finder/main.go                         # recusa o arranque sem os parâmetros; liga a porta de efeito
internal/config/config.go                  # lê os parâmetros sem criar número, grupo, consumidor ou prefixo
internal/domain/                           # envelope, referência validada e estados da chave
internal/application/consumir.go           # interpretação, espera do efeito e confirmação
internal/application/port/out/             # fila, chave do lote e efeito da recuperação
internal/adapter/in/queue/redisStream.go   # leitura do stream de logs, pendência e chave
docker-compose.yml                         # stream do contrato no lugar de fila:logs
.env.example                               # mesmos nomes, sem senha e sem número padrão
```

Pacotes que este plano não altera no significado: `internal/adapter/in/http`, `internal/adapter/out/openobserve`, `internal/adapter/out/mongo`, `internal/adapter/out/elastic`.

**Structure Decision**: A hexagonal já usada em `cmd/finder` e `internal/` recebe o contrato novo no lugar do item de código único. O adaptador Redis existente é o ponto de leitura. A recuperação do conteúdo não ganha adaptador neste plano.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Etapa sem consulta ao OpenObserve e sem ocorrência (princípios I e IV) | A spec entrega a referência e espera o efeito; a busca e a classificação ficam de fora | Implementar a busca agora inventaria a etapa seguinte |
| Porta de efeito sem adaptador de busca | A confirmação válida depende desse efeito e a spec proíbe infraestrutura intermediária | Confirmar ao interpretar a referência contradiz a resposta síncrona de 2026-11-08 |
