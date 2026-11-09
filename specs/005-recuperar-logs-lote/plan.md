# Implementation Plan: Recuperar logs do lote

**Branch**: `005-recuperar-logs-lote` | **Date**: 2026-11-08 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/005-recuperar-logs-lote/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

O Finder recebe a referência validada e recupera os logs desse lote no OpenObserve. A consulta é `POST /api/{organization}/_search`, no stream da referência, com filtro de `lhamalog_batch_id`, faixa de `minTimestampMicros` até `maxTimestampMicros + 1` e, no lote unitário, `lhamalog_record_id`. Quando os nomes dos campos de tenant e produto estão configurados, a consulta também os filtra e cada documento é conferido. A paginação começa em zero e avança pelo tamanho de página até uma página curta. Não há máximo de páginas.

A deduplicação usa `lhamalog_record_id`. Diferença só na lista configurada de campos de cópia é uma cópia; qualquer outra diferença é conteúdo divergente, sem versão escolhida. O lote só fica completo quando a contagem distinta fecha e não há falha de maior precedência. O resultado informa se é efeito concluído. Esta etapa não confirma a entrada Redis.

A pesquisa está em [research.md](./research.md).

## Technical Context

**Language/Version**: Go 1.27, como em `go.mod` e no README

**Primary Dependencies**: Biblioteca padrão `net/http`. Nenhuma biblioteca nova. `github.com/redis/go-redis/v9` permanece no módulo e não faz esta consulta.

**Storage**: OpenObserve já publicado no Compose, imagem `openobserve/openobserve:v1.1.2`, porta do host `5080`. A organização e o stream vêm da referência. O stream Redis `lhamalog:stream:logs` não é o lugar da consulta. Não há persistência, cache nomeado nem DLQ para os logs recuperados.

**Testing**: `make test` executa `go test -race ./...` sem OpenObserve. A prova com armazenamento fica atrás de `FINDER_OPENOBSERVE_INTEGRATION=1` e usa dados sintéticos. Os números e os nomes de campos dessa prova são dados da execução, não padrão do processo.

**Target Platform**: Processo do Finder no Docker Compose. A consulta de saúde existente permanece. Esta feature não classifica ocorrência e não confirma a entrada Redis.

**Project Type**: Serviço Go de um módulo. A recuperação é caso de uso interno, chamado pela porta de efeito da feature `004-capturar-referencia-lote`. Nenhum endpoint HTTP novo.

**Performance Goals**: Sem meta numérica. O aceite é o da spec: consulta fiel à referência (SC-001 e SC-002), paginação sem máximo e sem parar em `recordCount` (SC-003), um registro lógico por identificador (SC-004), completude (SC-005), isolamento (SC-006), falha sem parcial concluído (SC-007), lote misto só com os aceitos (SC-008), configuração ausente sem consulta (SC-009), efeito concluído só no lote completo e nas falhas que pedem correção (SC-010).

**Constraints**: Não criar URL, usuário, senha, timeout, tamanho de página, lista de campos de cópia, nomes de tenant e produto nem expressão de ordenação no código. Não criar máximo de páginas. Não registrar credencial nem conteúdo de log em falha. Não usar o identificador Redis nem `correlationId` para deduplicar. Não confirmar a entrada Redis. Não varrer dado sensível, não publicar checked, não integrar Ocorrência e não recuperar trace.

**Scale/Scope**: Um lote por chamada da porta de efeito. A quantidade de documentos físicos pode passar de `recordCount` por causa de reenvio. A leitura segue enquanto a página voltar cheia.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Princípio | Resultado | Leitura |
|-----------|-----------|---------|
| I. Finalidade: logs em ocorrências | Justificado | Esta etapa entrega os logs do lote. Erro, alerta e informação ficam fora, pelo recorte da spec. |
| II. Origem dos logs no coletor | Passa | Os logs vêm do que o Collector gravou no OpenObserve. A referência não é o conteúdo. |
| III. Categorias de ocorrência | Justificado | Nenhuma classe nova de ocorrência. A classificação desta feature é o resultado da recuperação. |
| IV. Fila Redis e apontamento no OpenObserve | Passa | A fila já entregou a referência. Esta etapa resolve o conteúdo e informa o efeito. A confirmação continua na primeira tarefa. |
| V. Incerteza sem invenção | Passa | URL, credencial, timeout, página, lista de cópia, nomes de tenant e produto e a expressão de ordenação ficam externos. O plano não escolhe esses valores. O máximo de páginas não é criado, pela resposta B de 2026-11-08. |

**Resultado do gate**: passa para o desenho. A implementação espera `/speckit-tasks`.

### Revisão depois do desenho

[research.md](./research.md), [data-model.md](./data-model.md), [contracts/recuperar-lote.md](./contracts/recuperar-lote.md) e [quickstart.md](./quickstart.md) repetem a consulta da spec, a deduplicação, a precedência de uma classificação e a ausência de máximo de páginas. Os valores da prova aparecem só no quickstart. Nenhum padrão numérico nem nome de campo foi colocado no processo. O gate permanece aprovado.

## Project Structure

### Documentation (this feature)

```text
specs/005-recuperar-logs-lote/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── recuperar-lote.md
└── tasks.md             # criado por /speckit-tasks, não por este comando
```

### Source Code (repository root)

```text
internal/domain/recuperacao.go             # deduplicação, completude e precedência da classificação
internal/application/recuperar.go          # paginação e montagem do resultado
internal/application/port/out/busca.go     # porta da consulta, sem HTTP
internal/adapter/out/openobserve/          # POST /api/{organization}/_search
internal/config/config.go                  # lê os parâmetros sem criar número, URL, senha ou nome de campo
.env.example                               # nomes novos, sem valor e sem senha
```

A porta de efeito de `specs/004-capturar-referencia-lote` chama este caso de uso. Este plano não lê o stream Redis e não confirma a entrada. Pacotes que este plano não altera no significado: `internal/adapter/in/http`, `internal/adapter/in/queue`, `internal/adapter/out/mongo`, `internal/adapter/out/elastic`.

**Structure Decision**: A hexagonal já usada em `internal/` recebe a consulta no adaptador que hoje só reserva o pacote, em `internal/adapter/out/openobserve/doc.go`. A regra de lote fica no domínio, testável sem OpenObserve. O cliente HTTP usa a biblioteca padrão.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Etapa sem classificar ocorrência (princípios I e III) | A spec para nos logs recuperados e no efeito informado à primeira tarefa | Classificar ocorrência agora entraria no que a spec deixou fora |
| Expressão de ordenação externa, sem coluna no código | A spec exige ordem estável e não nomeou a coluna. A documentação pública do OpenObserve cita `_timestamp` como ordem padrão | Gravar `_timestamp` no código escolheria a coluna que a spec não nomeou |
| Leitura sem máximo de páginas | Resposta B de 2026-11-08 | Criar um teto inventaria um limite que o solicitante recusou |
