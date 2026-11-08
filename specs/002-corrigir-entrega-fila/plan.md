# Implementation Plan: Corrigir entrega da fila

**Branch**: `002-corrigir-entrega-fila` | **Date**: 2026-11-05 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/002-corrigir-entrega-fila/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

O ajuste corrige a entrega do consumo já existente. A conclusão de um item aceito grava o código uma vez e tira o item da pendência na mesma operação. Se o processo cai no meio, a tentativa seguinte completa o que falta sem segunda conclusão. Se a instância some antes de qualquer efeito, outra instância, ou a mesma ao voltar, reserva o item na identidade original, sem prazo numérico, e esse item segue à frente dos que já estavam atrás. O Finder não escolhe prioridade: consome a fila na ordem já recebida. O código é único no produto. Item fora do contrato é descartado e sai da reserva. O processo em execução volta a atender `GET /health`.

O desenho permanece na hexagonal já usada. A pesquisa está em [research.md](./research.md).

## Technical Context

**Language/Version**: Go 1.27, como em `go.mod`

**Primary Dependencies**: `github.com/redis/go-redis/v9 v9.23.0`, já exigida em `go.mod`. Nenhuma biblioteca nova.

**Storage**: Fila Redis em stream e o conjunto de códigos que o ambiente já configura. Sem chave nova, sem fila de erro e sem campo de organização ou de nível.

**Testing**: `make test` executa `go test -race ./...` sem Redis. A prova com Redis real está no [quickstart.md](./quickstart.md) e falha se o Redis não estiver acessível quando `FINDER_REDIS_INTEGRATION=1`.

**Target Platform**: Processo do Finder no Docker Compose, ao lado do Redis já publicado em `localhost:6379`. Não grava ocorrência e não chama o OpenObserve.

**Project Type**: Serviço Go de um módulo, com entrada de fila e a consulta de saúde já existente.

**Performance Goals**: Sem meta numérica. O aceite é o das spec: uma conclusão por código (SC-001, SC-006), pendência de instância parada de volta ao lugar original (SC-002), ordem preservada conforme o número de instâncias (SC-003), descarte sem resíduo (SC-004) e consulta de saúde (SC-005).

**Constraints**: Não nomear stream, grupo ou campo. Não criar prazo numérico, fila de erro nem seleção de nível. Não alterar o corpo de `GET /health`. `FINDER_WORKERS` continua existindo e não é o grau de paralelismo deste consumo. Nomes Redis seguem o ambiente; o Compose de desenvolvimento já informa `fila:logs`, `finder` e `finder:codigos`.

**Scale/Scope**: Duas ou mais instâncias na mesma fila. O ajuste cobre os caminhos de `cmd/finder`, o caso de uso do consumo e o adaptador da fila. Não reabre a captura no OpenObserve.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

`.specify/memory/constitution.md` neste checkout é o template sem princípios preenchidos. O gate usa os princípios I a V já aplicados em `specs/001-consumir-fila-redis/plan.md`, da constituição 1.0.0 ratificada em 2026-10-29. Nenhum princípio novo foi criado a partir do template.

| Princípio | Resultado | Leitura |
|-----------|-----------|---------|
| I. Finalidade: logs em ocorrências | Justificado | Esta etapa não gera ocorrência. A spec limita o trabalho à entrega do consumo. |
| II. Origem dos logs no coletor | Justificado | O código de apontamento não vira log e o coletor não é analisado aqui. |
| III. Categorias de ocorrência | Justificado | Erro, alerta e informação ficam fora. Nenhuma classe nova. |
| IV. Fila Redis e apontamento no OpenObserve | Justificado | A fila continua entregando o código. A abertura do log fica para etapa posterior. O código permanece disponível uma vez. |
| V. Incerteza sem invenção | Passa | Não há nome novo de fila, campo, organização, nível ou fila de erro. Não há prazo numérico. A prioridade segue a ADR 169: o Finder não ordena. A unicidade global do código veio da resposta Q1 A. |

**Resultado do gate**: passa para o desenho. A implementação espera `/speckit-tasks`.

### Revisão depois do desenho

[research.md](./research.md), [data-model.md](./data-model.md), [contracts/entrega-fila.md](./contracts/entrega-fila.md) e [quickstart.md](./quickstart.md) usam o Redis, o conjunto e a consulta de saúde que o repositório já tem. Não criam destino, prazo nem campo. O gate permanece aprovado.

## Project Structure

### Documentation (this feature)

```text
specs/002-corrigir-entrega-fila/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── entrega-fila.md
└── tasks.md             # criado por /speckit-tasks, não por este comando
```

### Source Code (repository root)

```text
cmd/finder/main.go                       # liga a consulta de saúde e não encerra o processo na falha da conclusão
internal/application/consumir.go         # descarte, conclusão única e tentativa seguinte do estado parcial
internal/adapter/in/queue/redisStream.go # operação única de conclusão, pendência própria e pendência de consumidor sem conexão
internal/adapter/in/http/health.go       # sem mudança de significado
internal/domain/                         # estados já existentes; parcial é combinação, não um estado novo
```

Textos operacionais que ainda descrevem o processo sem a consulta de saúde, ou o item inválido sem descarte, são alinhados na implementação. Este plano não cria variável nem comando novo.

**Structure Decision**: O ajuste entra nos arquivos que já fazem o consumo. O handler de saúde permanece o arquivo atual. Pacotes de MongoDB, Elasticsearch e OpenObserve ficam fora.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Reservar pendência de consumidor sem conexão, em vez de esperar um tempo | FR-005 exige a volta ao lugar original sem prazo numérico | Tempo mínimo de ociosidade faria o número definir a liberação |
| Registro do código e saída da pendência na mesma operação | FR-001 e FR-003 exigem um desfecho e a cura do parcial | Escritas separadas mantêm a janela entre marcar e confirmar |
| Processo continua depois da falha da conclusão | Resposta da spec, opção A | Encerrar o processo ou repetir até o sucesso na mesma tentativa foram recusados |
