# Implementation Plan: Modelo de classificação de logs

**Branch**: `006-modelo-classificacao-logs` | **Date**: 2026-11-09 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/006-modelo-classificacao-logs/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

O Finder recebe um apontamento cujos logs já foram capturados. Não lê a fila de novo e não confirma a entrada. Cada log presente recebe uma classe: informação, alerta ou erro. O campo de nível e as palavras inteiras da mensagem valem juntos, e fica a classe mais grave. O catálogo só entra quando nenhum dos dois decide, na primeira entrada válida. A ocorrência gravada é uma por chamada, com a classe mais grave entre os logs presentes. Log ausente no meio do conjunto não impede a gravação e não entra na classe. Sem nenhum log presente, com o armazenamento inacessível ou com a gravação falha, não há ocorrência pronta, e os três motivos se distinguem.

A ocorrência no MongoDB guarda a classe, a identificação das regras, se o catálogo entrou, a identidade do apontamento, o tenant e o produto. Não guarda texto de log. A segunda gravação do mesmo apontamento insere outra ocorrência.

A pesquisa está em [research.md](./research.md).

## Technical Context

**Language/Version**: Go 1.27, como em `go.mod` e no README

**Primary Dependencies**: Biblioteca padrão para a regra de classe. `go.mongodb.org/mongo-driver/v2` entra nesta feature. O README lista esse módulo como candidato e ele ainda não está em `go.mod`. Nenhuma biblioteca de LLM, de expressão regular ou de banco vetorial.

**Storage**: MongoDB já publicado no Compose, imagem `mongo:8.0.34`. URI, banco e coleção vêm da execução. O código não escolhe esses três. A ocorrência é inserida, não substituída. OpenObserve e Redis não são lidos aqui.

**Testing**: `make test` executa `go test -race ./...` sem MongoDB. A prova com MongoDB fica atrás de `FINDER_MONGO_INTEGRATION=1` e usa dados sintéticos. Os nomes de URI, banco e coleção dessa prova não são padrão do processo.

**Target Platform**: Processo do Finder no Docker Compose. `GET /health` permanece. Esta feature não consome a fila.

**Project Type**: Serviço Go de um módulo. A classificação é caso de uso interno. Nenhum endpoint HTTP novo.

**Performance Goals**: O aceite observável é o da spec. Uma chamada com a gravação disponível deixa uma ocorrência consultável em até 2 segundos (SC-002). Esse prazo não vira timeout interno. A contagem de envios do texto a um serviço externo de linguagem é zero (SC-001).

**Constraints**: Três classes apenas. Sem spam e sem "não classificado". Sem texto de log na ocorrência e sem texto de log em falha. Sem reler a fila e sem confirmar a entrada. Sem agrupar por assinatura, sem gravidade, sem evento, sem sugestão, sem notificação e sem catálogo de mascaramento. Sem criar o catálogo de classificação. Sem URI, banco ou coleção padrão no código.

**Scale/Scope**: Um apontamento por chamada. Vários logs nesse apontamento viram uma ocorrência. Duas chamadas do mesmo apontamento viram duas ocorrências.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Princípio | Resultado | Leitura |
|-----------|-----------|---------|
| I. Finalidade: logs em ocorrências | Passa | A chamada transforma os logs presentes do apontamento em uma ocorrência. |
| II. Origem dos logs no coletor | Passa | A entrada já é o log capturado. Esta feature não coleta de novo. |
| III. Categorias de ocorrência | Passa | As classes são informação, alerta e erro. Spam não entra. |
| IV. Fila Redis e apontamento no OpenObserve | Justificado | A fila e o armazenamento já entregaram o apontamento nas specs `001` a `005`. Esta spec manda não ler a fila de novo e não confirmar a entrada. |
| V. Incerteza sem invenção | Passa | URI, banco, coleção e o catálogo ficam externos. O nome do campo de nível no armazenamento não é escolhido aqui: a chamada já traz o texto do campo e o texto da mensagem separados. |

**Resultado do gate**: passa para o desenho. A implementação espera `/speckit-tasks`.

### Revisão depois do desenho

[research.md](./research.md), [data-model.md](./data-model.md), [contracts/classificar-apontamento.md](./contracts/classificar-apontamento.md) e [quickstart.md](./quickstart.md) repetem as três classes, a palavra inteira, a gravidade, uma ocorrência por apontamento e a gravação sem o texto. URI, banco e coleção aparecem só como dados da prova. O gate permanece aprovado.

## Project Structure

### Documentation (this feature)

```text
specs/006-modelo-classificacao-logs/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── classificar-apontamento.md
└── tasks.md             # criado por /speckit-tasks, não por este comando
```

### Source Code (repository root)

```text
internal/domain/classificacao.go          # classe do log, palavra inteira e classe da ocorrência
internal/application/classificar.go       # um apontamento, uma gravação, três falhas
internal/application/port/in/classificar.go
internal/application/port/out/ocorrencia.go
internal/adapter/out/mongo/ocorrencia.go  # insere e consulta, sem o texto do log
internal/config/config.go                 # lê URI, banco e coleção sem criar padrão
```

O consumidor em `internal/adapter/in/queue` não passa a chamar este caso de uso. Pacotes que este plano não altera no significado: `internal/adapter/in/http`, `internal/adapter/in/queue`, `internal/adapter/out/elastic`, `internal/adapter/out/openobserve`.

**Structure Decision**: A regra cabe no domínio e é testável sem MongoDB. A porta de saída só insere e lê a ocorrência. O adaptador usa o driver que o README já nomeou e que esta feature adiciona. O caso de uso não abre Redis nem OpenObserve.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Esta feature não lê a fila (princípio IV) | A spec, na resposta B de quantidade, manda não ler a fila de novo e não confirmar a entrada. O apontamento chega pronto. | Ler o stream de novo reabriria o consumo e a confirmação que as specs `001` a `005` já separaram. |
| Driver MongoDB novo | A spec exige a ocorrência no MongoDB. O pacote `internal/adapter/out/mongo` hoje só reserva o diretório. O README nomeia `go.mongodb.org/mongo-driver/v2` como candidato ainda ausente de `go.mod`. | Gravar em memória deixaria a ocorrência só enquanto o processo está no ar, e a spec pede consulta depois da gravação no MongoDB. |
| Finder grava a ocorrência | A resposta Q2 A de 2026-11-09 atribui a gravação a esta feature. A ADR_062 entrega isso ao módulo de Ocorrência. A constituição não definia a persistência. | Enviar só um resumo deixaria a ocorrência sem gravação, que é o resultado pedido. |
