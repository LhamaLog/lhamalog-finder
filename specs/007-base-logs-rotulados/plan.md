# Implementation Plan: Base de logs rotulados para treino e avaliação

**Branch**: `007-base-logs-rotulados` (trabalho git permanece em `Development`, conforme spec) | **Date**: 2026-11-16 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/007-base-logs-rotulados/spec.md` (etapas 3b–3e e **3g**: 007-meta, DN-007-2 B, DN-006-2/007-3 B, DN-007-4 A, **007 H2 D**, 007-formato, 007-caminho, DN-007-1 C)

## Summary

Entregar um **conjunto rotulado** em **arquivos versionados no Git** (007-4 A), na ordem de **milhares** de registros com **mesma contagem mínima** por classe ERROR, WARNING e INFO (007-1 C, 007-2 B, 007-9 A), por **anonimização de logs reais**, sem PII e com **declaração de origem** (007-14 B, 007-15 A, 007-16 A). No MVP por regras (ADR_167), o uso é **somente avaliação** das regras de classificação (007-3 A); **substitui** os dados sintéticos de prova da spec 006 (007-19 B). Consumo principal via **testes automatizados**, alimentando o classificador **sem fila Redis** (007-22 B, 007-23 A).

Meta na partição de validação (LLT-199 critério 2): **≥ 95% de precisão, recall e F1 em cada classe** (ERROR, WARNING, INFO) — **007-meta**, **DN-007-2 B** e **DN-006-2/007-3 B** (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h; respostas-etapa3c.md §1; origem: 006-2 B). Artefatos em **JSONL** sob **`datasets/labeled-logs/`** na raiz do repositório (documentado; pasta **não** criada nesta etapa de spec). Em apontamento: **um** `class_label` gold no agregado; `members` **sem** gold (**DN-007-4 A**; fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); respostas-etapa3c.md §3). Na **avaliação** de apontamento: predição **por** `member`; **classe prevista** do apontamento = **mais grave** entre predições (**ERROR** > **WARNING** > **INFO**, ADR_170); comparar ao gold agregado; **unidade** da métrica = apontamento (**007 H2 D**; FR-007c; fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), 007 H2 D; origem: ADR_170; respostas-etapa3f.md §D). Subconjunto **só-catálogo**: critério documentado, **registros adiados** até existir catálogo (**DN-007-1 C**; 006-7 D; ADR_167 [PENDENTE]).

Artefatos de desenho: [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md).

## Technical Context

**Language/Version**: Go **1.27** (`go.mod`, linha 3; README, Go 1.27.0 na imagem Docker).

**Primary Dependencies (presentes em `go.mod`)**: `github.com/joho/godotenv`, `github.com/redis/go-redis/v9` — **não** no caminho crítico desta feature (007-23 A: avaliação desacoplada da fila).

**Primary Dependencies (README — candidatas, não em `go.mod`)**: mongo-driver, OpenTelemetry, `golang.org/x/sync`, testcontainers — **MUST NOT** assumidas (constituição, Restrições Técnicas; README § Dependências candidatas). Esta feature **não** exige MongoDB nem OpenObserve em runtime de teste da base (entrada direta no classificador).

**Primary Dependencies (implementação da base)**: Leitor de arquivos versionados + validador de esquema; integração com testes Go existentes (`make test`). Biblioteca de serialização: **`encoding/json`** lendo **JSONL** (uma linha = um objeto JSON) — **007-formato A** (fonte: respostas-joao.md §Etapa 3b).

**Storage**: Artefatos de dados e documentação no **próprio repositório** (007-4 A), diretório **`datasets/labeled-logs/`** na raiz de `lhamalog-finder` (**007-caminho C**). Criação física do diretório fica para a implementação; nomes de arquivos e atributos JSON na tabela de mapeamento da spec (007-6 D).

**Testing**: `make test` → `go test -race ./...` (README). Validação de 100% dos registros (SC-004 / 007-21 A) via testes de contrato/fixture. Métricas LLT-197/199 (006-2 B): precisão, recall e F1 **≥ 95% por classe** na partição `validacao` (007-meta; DN-007-2 B; DN-006-2/007-3 B). Itens `log`: comparar `class_label` ao `classe` predito. Itens `apontamento`: agregar predições dos `members` pela **classe mais grave** antes de comparar ao gold (**007 H2 D**; contrato [entrada-classificador-desde-base.md](./contracts/entrada-classificador-desde-base.md)).

**Target Platform**: Linux / ambiente Docker de desenvolvimento (README). **Nenhum** endpoint HTTP novo no Finder para esta feature.

**Project Type**: Módulo Go, arquitetura hexagonal leve (`cmd/finder/`, `internal/domain`, `internal/application`, `internal/adapter`, `internal/config` — README § Layout).

**Performance Goals**: Sem requisito de volume máximo (007-24 A). Ordem de milhares de registros — sem meta de latência inventada.

**Constraints**: Sem PII, nem fictícia (007-16 A; ADR_064). Rótulo **só classe** — sem assinatura/agrupamento LLT-200 (007-17 A). Rejeitar registro sem `level_text` ou `message_text` (007-13 A). Contrato upstream **log.ingested v1** (`schemaVersion` `"1"`) relevante como alinhamento de ids (T5 D; `contrato-log-ingested-v1.txt` §4, §7), não como consumo de fila nesta feature.

**Scale/Scope**: Milhares de registros; mínimo **igual** por ERROR/WARNING/INFO; split treino/validação/teste **documentado** para pós-MVP (007-20 B); itens **log único** e **apontamento** (007-7 C).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Princípio | Resultado | Leitura e fonte |
|-----------|-----------|-----------------|
| I. Finalidade: logs em ocorrências | **Passa (escopo indireto); conflito documentado** | Constituição I: o Finder MUST agregar em **ocorrência** (`.specify/memory/constitution.md`, princípio I). **T1 REVISTO**: o **Módulo de Ocorrência** cria a ocorrência; o Finder envia `ResumoAnalise` (fonte: respostas-joao.md, T1 REVISTO 15/11; revisao-T1-drive.md §2). Esta feature **não** grava ocorrência — fornece base para **avaliar** regras do fluxo produtivo. **Pendência**: **emenda do princípio I** na constituição (não editada nesta pasta). |
| II. Origem dos logs no coletor | **Passa** | Registros espelham campos derivados de captura OTel/OpenObserve (`lhamalog_batch_id`, `lhamalog_record_id`; contrato §7); origem declarada como anonimização (007-14 B). |
| III. Categorias de ocorrência | **Passa** | `class_label` ∈ {ERROR, WARNING, INFO} (007-9 A; constituição III via ADR_170). |
| IV. Fila Redis → OpenObserve | **Passa (fora do caminho da base)** | Consumo da base **desacoplado da fila** (007-23 A); ids de lote/registro mantêm rastreio sem tratar código de apontamento Redis como conteúdo (coerente com constituição IV para o pipeline produtivo). |
| V. Incerteza sem invenção | **Passa com lacunas nomeadas** | Etapas 3b–3e e **3g**: formato **JSONL**, caminho **`datasets/labeled-logs/`**, metas **por classe** (precisão/recall/F1 ≥ 95%); gold apontamento **DN-007-4 A**; regra de avaliação apontamento **007 H2 D** (analyze-007-etapa3f.md H2 **fechado**); pendências externas listadas abaixo. 007-26 A: nenhuma lacuna **nova** bloqueia o plano inteiro. |

**Divergências conscientes (sem editar constituição/ADR nesta pasta)**:

- **Princípio I × T1 revisto** → pendência **emenda do princípio I** (tabela § Pendências).
- Taxonomia e INFO como ocorrência (**T7 A**) divergem de SSID_007.1 legado — spec § Divergências; a base **inclui** rótulos INFO (007-9 A).

**Resultado do gate (pré e pós-Fase 1)**: aprovado para implementação do conjunto, validação estrutural e cenários **SC-005** (metas por classe fechadas etapa 3e); **conteúdo** do subconjunto só-catálogo permanece bloqueado (**DN-007-1 C**) conforme tabela abaixo.

## Project Structure

### Documentation (this feature)

```text
specs/007-base-logs-rotulados/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── registro-rotulado.md
│   └── entrada-classificador-desde-base.md
└── tasks.md             # /speckit-tasks — fora deste comando
```

### Source Code (repository root) — alvo de implementação

Estrutura **existente** (README): `cmd/finder/`, `internal/domain/`, `internal/application/`, `internal/adapter/`, `internal/config/`.

**Alvo coerente com esta feature** (007-caminho C, 007-formato A):

```text
datasets/labeled-logs/                       # JSONL + docs versionados 007-4 A
  ORIGEM.md                                  # declaração 007-15 A
  README.md ou SPLIT.md                      # split 007-20 B
  *.jsonl                                    # um objeto JSON por linha (007-formato A)
  ...                                        # manifest de contagens 007-25 A
internal/...                                 # loader + validação + testes 007-22 B (pacote em tasks)
```

**Structure Decision**: dados como artefatos Git separados do código; loader hexagonal opcional em `internal/` consumido por testes. **Não** adicionar adaptador HTTP ao Módulo Ocorrência nesta feature.

## API do Módulo Ocorrência (rota, JWT, respostas HTTP, 5xx/timeout)

Esta seção registra a pendência **DN-009-2** / `respostas-etapa3.md` §10 para o ecossistema; **não se aplica** à entrega da base rotulada (007 não envia `ResumoAnalise`).

| Tópico | O que há fonte | O que falta |
|--------|----------------|-------------|
| Transporte | HTTP interno + JWT (**N1 A**; ADR_028, ADR_054, ADR_141) | Rota/path, claims, códigos HTTP, timeout/5xx — **sem** documento (respostas-etapa3.md §10) |
| Decisão João | **API Módulo: pendência no plan até o contrato** (respostas-joao.md, §Etapa 3) | Contrato publicado pelo Módulo Ocorrência |
| Retentativa | 009-18 A (3 tentativas + alerta) aplica-se a **publicação** de ocorrência, não a leitura da base 007 | — |

**Impacto nesta feature**: **nenhuma tarefa** de adaptador HTTP. Tarefas equivalentes nas specs **006** / **009** permanecem **BLOQUEADAS** até contrato (ver § Pendências).

Divergências conscientes (incl. Jira × João em métricas e gold): ver [spec.md](./spec.md) § Divergências documentadas.

## Tarefas bloqueadas (dentro do escopo 007)

| Bloqueio | Escopo afetado | Restante desbloqueado |
|----------|----------------|----------------------|
| **DN-007-1 C** | **Conteúdo** do subconjunto só-catálogo | Demais registros; documentação do critério de inclusão (FR-011); itens `apontamento` com gold agregado (DN-007-4 A) |
| **006-7 D** | Exemplos que dependem de sintaxe de catálogo | Registros decididos só por nível/mensagem |

## Pendências

Itens **fora desta execução** ou que exigem artefato externo (sem criar/editar ADR nesta pasta):

| ID | Descrição | Fonte | Ação nesta feature |
|----|-----------|-------|---------------------|
| **DN-009-2** | API Módulo Ocorrência: rota, JWT, respostas, 5xx/timeout | respostas-joao.md §Etapa 3; respostas-etapa3.md §10 | Só referência na seção acima; **sem** código HTTP aqui |
| **Emenda princípio I** | Constituição I vs T1 revisto (Módulo cria ocorrência) | respostas-joao.md, T1 REVISTO; revisao-T1-drive.md §2 | Registrar emenda em `.specify/memory/constitution.md` **fora** desta feature |
| **DN-008-4** | Normalização para assinatura: aparar espaços, ignorar maiúsculas/minúsculas, substituir números/UUIDs/IPs/datas/hex por marcadores | respostas-joao.md §Etapa 3b, João 16/11 09h02 | **Decidido** para a spec **008**; **fora do escopo 007** (007-17 A: rótulo só classe, sem assinatura) |
| **006-7 D** | Catálogo: sintaxe e dono | respostas-joao.md, 006-7 D; ADR_167 [PENDENTE] | Aguardar para subconjunto só-catálogo |
| **N3 A** | Ampliação formal de `AchadoAnalise` em ADR | respostas-joao.md, N3 A | Não alterar ADR aqui |
| **010-2 A** | Lista de padrões de mascaramento em ADR | respostas-joao.md, 010-2 A | Spec 010 |
| **010-27 A** | Casos de teste de mascaramento (João) | respostas-joao.md, 010-27 A | Spec 010 |
| **006-10 A** | Alinhar Jira LLT-197/199 com spec 006 | respostas-joao.md, 006-10 A | Não editar Jira neste repo |
| **ADR_019 / ADR_167 / SSID_007.1** | Legado spam / treino LLM / FA01 INFO | instrução de escopo | Apenas citação; sem edição |

## Complexity Tracking

> Vazio — sem violação constitucional que exija justificativa de complexidade extra nesta feature.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| — | — | — |

## Fase 0 e Fase 1 (status)

| Fase | Artefato | Status |
|------|----------|--------|
| 0 | [research.md](./research.md) | Concluído — 007-formato A, 007-caminho C; **007-meta** + DN-007-2 B + **DN-006-2/007-3 B** + **DN-007-4 A** (etapa 3e); **007 H2 D** / R20 (etapa 3g); R9 split qualitativo; R18 FR-013 |
| 1 | [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md) | Concluído — gold apontamento agregado; avaliação apontamento **007 H2 D**; métricas SC-005 fechadas; API Módulo Ocorrência **fora** do escopo (DN-009-2) |

## Lacunas documentadas (não bloqueiam o plano inteiro)

| Tópico | Tratamento | Fonte |
|--------|------------|-------|
| Proporção/seed/estratificação do split | Critérios **qualitativos** em `SPLIT.md` (FR-005); **sem** percentuais/seed nas fontes — MUST NOT inventar | 007-20 B; analyze M2; research.md R9 |

## Evidência desta execução

T024/T031/T041/T042/T045 dependem da implementação real da operação ClassificarLog da 006, não apenas do mapeamento T015. Essa implementação é pré-requisito da composição de validação, sem código 006 duplicado nesta branch. N3 A, DN-009-2 e emenda do princípio I permanecem externos e não exigem adaptador HTTP nesta feature.
