# Implementation Plan: Agrupar logs semelhantes na mesma ocorrência

**Branch**: `Development` (feature `008-agrupar-logs-ocorrencia`) | **Date**: 2026-11-16 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/008-agrupar-logs-ocorrencia/spec.md` (clarificado etapas 2b + **3b** + **3c** + **3e** + **3g** + **3h**; T1 revisto; chave 008-11 rev.)

## Summary

O **Finder** calcula a **assinatura** de cada log classificado como **nível OTel** + **mensagem normalizada** + **`service.name`** + **host** (máquina) (**DN-008-1 B**, 008-1 fech. 3b), aplicando o algoritmo canônico de normalização da mensagem (**DN-008-4 B**, etapa 3c) com marcadores **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`** (**DN-008-5 A**, etapa 3e); **`traceId`** e **`correlationId`** **não** entram na assinatura (**008-2 rev.**). O **valor** canônico é **texto** concatenando os quatro componentes com separador fixo **`|`** (**008-assinatura A**, etapa 3e; **DN-008-8 A**, etapa 3g). **No lote em curso**, logs equivalentes MUST agrupar de forma **determinística e local** pela assinatura (**FR-008**, **A1**); a **consulta ao OpenObserve** aplica **somente** a histórico/recorrência/similaridade **fora do lote** (**008-3 B**, **008-5/6**). Agrupa logs da mesma assinatura no **mesmo produto** num **AchadoAnalise** (**008-18 A**), **consulta** ocorrências existentes **uma chamada por achado** (**N2 B**, **DN-008-7 B**, etapa 3e) e agrega **entre lotes** (**008-12/008-14**), avalia **recorrente** via OpenObserve com até **3 tentativas** e, se esgotar, aviso ao **Core** (API interna, **batchId** + **motivo**; **sem DLQ** — **DN-008-3**, **DN-008-6** rev. 3h), e envia **exatamente um** `ResumoAnalise` **por lote** (HTTP+JWT — **N1 A**; rotas **DN-009-2**). O **Módulo de Ocorrência** cria ou atualiza pela chave **(assinatura, classe, produto, tenant)** (**008-11 rev.**), reabertura **FECHADA→ABERTA** (**008-4 B**), contador sem mudar status (**008-10 A**). **INFO** inclusivo (**T7 A**). Lote **incompleto** (**T4 A**) não agrupa. **Sem migração** legada documentada (**DN-008-2**). Meta **≥ 95% precisão por classe** fora desta feature (**DN-007-2 B**).

Artefatos: [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md).

## Technical Context

**Language/Version**: Go **1.27** (fonte: `go.mod` L3; README, imagem Docker Go 1.27.0).

**Primary Dependencies (presentes em `go.mod`)**: `github.com/joho/godotenv`, `github.com/redis/go-redis/v9` — consumo Redis Streams (constituição IV; specs 004/005).

**Primary Dependencies (README — candidatas, não em `go.mod`)**: `go.mongodb.org/mongo-driver/v2`, OpenTelemetry, `golang.org/x/sync`, testcontainers — **MUST NOT** ser assumidas instaladas (README § Dependências candidatas; constituição Restrições Técnicas). Persistência de **ocorrência** pelo Finder **não** é escopo (**T1 rev.**, FR-001). Existe `internal/adapter/out/mongo/doc.go` como reserva, sem papel de ocorrência nesta feature.

**Primary Dependencies (integração — parcialmente decidida)**: Cliente **HTTP + JWT** para **ResumoAnalise** e consulta de ocorrências (**N1 A**, **N2 B**); biblioteca JWT, rotas, timeout e política 5xx **NEEDS CLARIFICATION** neste repositório até contrato do Módulo (**DN-009-2**; fonte: respostas-joao.md §Etapa 3; origem: respostas-etapa3.md §10, ADR_028, ADR_054).

**Storage / consulta**: Conteúdo de log no **OpenObserve** (contrato-log-ingested-v1.txt §2–3; spec 005). `internal/adapter/out/openobserve/doc.go` declara que organização, stream, índice e credenciais **não foram informados** — forma da query de similaridade/recorrência permanece **NEEDS CLARIFICATION** operacional (sem inventar API HTTP/query). Chave idempotência ingested: `internal/domain/chave.go`, envelope v1 `internal/domain/envelope.go`.

**Testing**: `make test` → `go test -race ./...` (README L103). Aceite SC-001 com assinatura (**DN-008-1 B**), algoritmo **DN-008-4 B**, literais **DN-008-5 A**, forma **008-assinatura A**, separador **`|`** (**DN-008-8 A**) e escape **`|`** → **`\|`** (**DN-008-9 A**) fechados; escape de **`\`** literal nos componentes pendente (sem fonte). Handoff HTTP e consulta Módulo bloqueados até **DN-009-2** e **ADR N3** (granularidade da consulta **lógica** fechada: **por achado**). Base rotulada da spec **007**: formato **JSONL** (**007 formato A**) em `datasets/labeled-logs/` na raiz do repositório (**007 caminho C**) — só documentação cruzada; **MUST NOT** criar pasta/dados nesta feature (fonte: respostas-joao.md §Etapa 3b).

**Target Platform**: Serviço Finder no Docker Compose (README); OpenObserve local porta 5080 (`make run`).

**Project Type**: Módulo Go, arquitetura hexagonal leve — `cmd/finder/`, `internal/domain`, `internal/application`, `internal/adapter`, `internal/config` (README § Layout).

**Performance Goals**: Sem meta numérica com fonte para consulta OpenObserve (spec.md Fora de escopo).

**Constraints**: Assinatura = nível + mensagem normalizada + `service.name` + host (**DN-008-1 B**; ADR_013 mesma máquina). **Normalização**: **DN-008-4 B** + marcadores **DN-008-5 A**. **Valor** = concatenação texto (**008-assinatura A**; separador **`|`** — **DN-008-8 A**; escape **`|`** → **`\|`** — **DN-008-9 A**). Recorrência: **ADR_013**, janela **500 lotes** (**008-5/6**). Falha temporária OpenObserve: **3 tentativas**, depois aviso ao **Core** (API interna, **batchId** + **motivo**; **sem DLQ** — **DN-008-3**, **DN-008-6** rev. 3h); **MUST NOT** `recorrente=false` por degradação. Consulta Módulo **por achado**; envio **um ResumoAnalise/lote** (**DN-008-7 B**). Vetorial legado não reintroduzido (**ADR_168**, FR-008). `productId` inválido: falha crítica sem retry (**008-16 A**, ADR_049).

**Scale/Scope**: Uma ocorrência por achado ERROR/WARNING/INFO no Módulo (**T2 C**); unicidade **(assinatura, classe, produto, tenant)** (**008-11 rev.**). Depende de classificação (spec **006**). A **publicação de evento / notificação de ocorrência** ao Core permanece na spec **009** (spec.md fora de escopo). Em falha **temporária** da consulta OpenObserve, após **3 tentativas**, o Finder MUST avisar o **Core** pela **API interna** (**batchId** + **motivo**; **FR-002a**, **DN-008-3**, **DN-008-6** rev. 3h; fonte: respostas-joao.md §Etapa 3g, DLQ B). **ADR_051 [PENDENTE]**: rota, payload e evento **MUST NOT** inventar.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Princípio | Resultado | Leitura e fonte |
|-----------|-----------|-----------------|
| I. Finalidade: logs em ocorrências | **Pendência de emenda** | Finder agrega em **AchadoAnalise/ResumoAnalise**; Módulo materializa ocorrência (respostas-joao.md, T1 rev.; ADR_062). Divergência registrada; **não** editar constituição nesta execução. |
| II. Origem dos logs no coletor | Passa | Agregação sobre logs recuperados via referência ingested v1 e OpenObserve (FR-001; contrato §3). |
| III. Categorias de ocorrência | Passa | Apenas erro/alerta/informação → ERROR, WARNING, INFO (FR-006; respostas-joao.md taxonomia 14/11). |
| IV. Fila Redis → OpenObserve → agregação | Passa | Fluxo upstream preservado; esta feature estende pós-classificação sem tratar código de apontamento como conteúdo (`internal/domain/item.go`). |
| V. Incerteza sem invenção | Passa com lacunas explícitas | **DN-009-2** (API Módulo HTTP), escape de **`\`** literal (sem fonte), query OpenObserve operacional, contrato Finder→Core após falha OO (**ADR_051 [PENDENTE]**) e **ADR N3** sem inventar. **DN-008-1**–**DN-008-9 A** (etapas 3b–3h) e **008-assinatura A** fechados salvo escape **`\`**. **007-26 A** é processo da spec 007, não bloqueia este plano. SC-005: aceite ponta a ponta pendente até **DN-009-2**, **ADR_051** e **ADR N3**. |

**Resultado do gate (pré-Fase 0)**: aprovado para pesquisa e desenho com pendências e tarefas bloqueadas documentadas.

### Revisão pós-Fase 1

[data-model.md](./data-model.md) e [contracts/](./contracts/) refletem **DN-008-1 B**, **DN-008-4 B**, **DN-008-5 A**, **DN-008-6**, **DN-008-7 B** (consulta) / envio por lote, **008-assinatura A**, **DN-008-8 A** (`|`), N2 B, chave 008-11 rev., **DN-008-3**, escopo **FR-008** / **R16** (A1), handoff ampliado ([PENDENTE ADR] N3), **FR-002a** / aviso ao Core após 3 falhas OO (**ADR_051 [PENDENTE]**), contratos HTTP Módulo sem inventar rota; **DN-008-9 A** (**`|`** → **`\|`**); escape **`\`** pendente. Gate mantido (C1: emenda I pendente).

## Project Structure

### Documentation (this feature)

```text
specs/008-agrupar-logs-ocorrencia/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── assinatura-log.md
│   ├── agregar-achado-lote.md
│   ├── consulta-ocorrencias-modulo.md
│   ├── consulta-recorrencia-openobserve.md
│   └── resumo-analise-agrupado.md
└── tasks.md             # /speckit-tasks — fora deste comando
```

### Source Code (repository root) — alvo de implementação

Estrutura **existente** (README § Layout):

```text
cmd/finder/
internal/domain/
internal/application/
internal/application/port/out/
internal/adapter/in/queue/
internal/adapter/out/openobserve/   # doc.go — consulta OO a implementar (query NEEDS CLARIFICATION)
internal/adapter/out/modocorrencia/ # NEEDS CLARIFICATION — DN-009-2; pacote ainda ausente no repo
```

**Structure Decision**: estender hexágono existente; **sem** persistência de ocorrência no Finder (respostas-joao.md, T1 rev., FR-001).

## Fases de implementação (referência)

| Fase | Escopo | Bloqueios |
|------|--------|-----------|
| Assinatura e lote | `NormalizarMensagem` (**DN-008-4 B**, literais **DN-008-5 A**), componentes DN-008-1 B, serialização texto (**008-assinatura A**, **`|`**), escape **`|`** → **`\|`** (**DN-008-9 A**), agrupamento no lote, validação 008-16 | Escape **`\`** literal (sem fonte) |
| Entre lotes | Fusão contador/extremos após **uma consulta por achado** (**DN-008-7 B**) | Adaptador HTTP consulta + envio **DN-009-2** |
| Recorrência | Porta OpenObserve + marca `recorrente`; após 3 falhas, aviso Core (**DN-008-3**, **DN-008-6** rev. 3h; **ADR_051 [PENDENTE]**) | Query/index NEEDS CLARIFICATION |
| Handoff | **Um** `ResumoAnalise`/lote em memória (**DN-008-7**) | **N3 A** (nomes JSON finais); **DN-009-2** (envio real) |

## API do Módulo Ocorrência (rota, JWT, respostas HTTP, 5xx/timeout)

Pendência **DN-009-2** / `respostas-etapa3.md` §10; decisão João: **«API Módulo: pendência no plan até o contrato»** (respostas-joao.md, §Etapa 3). **Não** bloqueia o plano da 008 como um todo; bloqueia **somente** adaptadores HTTP (`consulta-ocorrencias-modulo`, envio `ResumoAnalise`) e validações integradas que exijam Módulo real.

| Tópico | O que há fonte | O que falta (sem inventar) |
|--------|----------------|----------------------------|
| Transporte | HTTP interno + JWT (**N1 A**); rede privada (**ADR_141**); JWT próprio (**ADR_028** [PENDENTE algoritmo/chaves], **ADR_054**) | Rota/path Finder → Módulo (consulta N2 B e entrega lote); claims; biblioteca |
| Semântica consulta | Chave **(assinatura, classe, produto, tenant)** (**008-11 rev.**); **uma chamada por achado** antes do envio (**N2 B**, **DN-008-7 B** — fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02)); fusão de contador/extremos no Finder | Corpo request/response HTTP (**DN-009-2**) |
| Semântica entrega | **Exatamente um** `ResumoAnalise` **por lote** (**DN-008-7**, ADR_062); idempotência **batchId + execId** (**N4 A**); idempotência no Módulo (**ADR_014**); Finder só conclui após enviar (seq. v6) | Confirmação HTTP; mapeamento 2xx ↔ upsert |
| Respostas HTTP | Idempotência Drive (**ADR_014**); 009-18 A (3 tentativas + DLQ) vale **publicação** no ecossistema, não copiada automaticamente para este handoff | Códigos 2xx/4xx/5xx; 202 vs 200; timeout; lista retentável (ADR_118 é OTLP Coletor) |
| Drive documentado | Única rota de ocorrência citada: `POST /ocorrencias/{id}/sugestoes` — **não** é handoff Finder (**respostas-etapa3.md** §10) | ADR/contrato do Módulo Ocorrência |

**Tarefas afetadas** (implementação futura em `tasks.md`): porta consulta ocorrências, porta envio `ResumoAnalise`, `internal/adapter/out/modocorrencia/`, quickstart §5–§8 com Módulo real, SC-005 ponta a ponta — **BLOQUEADAS** até contrato publicado.

## Tarefas bloqueadas (dentro do escopo 008)

| Bloqueio | Escopo afetado | Restante desbloqueado |
|----------|----------------|----------------------|
| **DN-009-2** | Adaptador HTTP+JWT (consulta **por achado** e envio **por lote**); URL; política retry/timeout/5xx; SC-005 com Módulo real | Assinatura com **`|`** (**DN-008-8 A**); agrupamento no lote; fusão **em memória**; política **3 tentativas** + aviso Core (**DN-008-3**, **DN-008-6** rev. 3h; contrato **ADR_051** pendente); contratos lógicos |
| **Escape `\`** | Escape de **barra invertida** literal nos componentes do `valor` | **DN-008-9 A** (**`|`** → **`\|`**) fechado; testes SC-001 com marcadores **DN-008-5 A** |
| **N3 A** | Nomes JSON finais e obrigatoriedade dos campos ampliados no `AchadoAnalise` | Requisitos e placeholders em [resumo-analise-agrupado.md](./contracts/resumo-analise-agrupado.md) |

**Fora do escopo 008 (só Pendências):** **DN-007-2** (meta 95% / recall-F1 — spec 007/006); **006-10 A** (Jira); **010-2 A**, **010-27 A** (spec 010).

## Decisões tomadas (etapa 3h)

| ID | Decisão | Nota de vigência |
|----|---------|------------------|
| **DLQ B** | **Não** há DLQ no fluxo `log.ingested` quando a consulta OpenObserve falha (contrato v1 não herda retry/DLQ do guia; **DN-006-1b B**). | (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22); origem: DLQ B; respostas-etapa3g.md §1) |
| **DN-008-3** (rev. 3h) | Até **3 tentativas**; após esgotar, Finder MUST avisar o **Core** pela **API interna** com **batchId** e **motivo** (**ADR_051 [PENDENTE]**). **Substitui** alcance anterior (3 tentativas + DLQ / reenvio stream — etapas 3b/3e). | (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22); origem: DLQ B) |
| **DN-008-6** (rev. 3h) | Mesmo gatilho: Finder→Core API interna (**batchId** + **motivo**); **não** DLQ nem alerta WebSocket pós-DLQ. **Substitui** alcance da etapa 3e (DLQ + Core/WebSocket). | (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22); origem: DLQ B) |
| **DN-008-9 A** | **`|`** literal em componente do `valor` MUST ser escapado como **`\|`** antes da concatenação. Escape de **`\`** literal: só pendência (sem fonte). | (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22); origem: DN-008-9 A) |

## Pendências

Itens **fora desta execução** ou que exigem artefato externo (sem criar ADR nesta pasta):

| ID | Descrição | Fonte | Ação nesta pasta |
|----|-----------|-------|------------------|
| **DN-009-2** | API Módulo: rota, JWT, HTTP, 5xx/timeout (consulta + envio) | respostas-joao.md §Etapa 3; respostas-etapa3.md §10 | § API do Módulo acima; contrato [consulta-ocorrencias-modulo.md](./contracts/consulta-ocorrencias-modulo.md) |
| **Contrato API Core (alertas)** | Rota, payload e evento Finder→Core após 3 falhas na consulta OpenObserve (**ADR_051 [PENDENTE]**; **batchId** + **motivo**) | respostas-joao.md §Etapa 3g, DLQ B | Fora desta pasta; **MUST NOT** inventar na 008 |
| **Escape `\`** | Escape de **barra invertida** literal nos componentes do `valor` (**DN-008-9 A** fechou só **`|`**) | spec.md Questões em aberto §3 | [assinatura-log.md](./contracts/assinatura-log.md) |
| **N3 A** | ADR: ampliar `AchadoAnalise` (assinatura, contagem, primeira/última, versão regras) | respostas-joao.md, N3 A | Contrato marca [PENDENTE ADR] |
| **006-10 A** | Alinhar Jira LLT-200 com specs 006/008 | respostas-joao.md, 006-10 A | Não alterar Jira |
| **010-2 A** | Lista de padrões de mascaramento em ADR | respostas-joao.md, 010-2 A | Spec 010 |
| **010-27 A** | Casos de teste de mascaramento / normalização (João) | respostas-joao.md, 010-27 A | Informa DN-008-4; spec 010 |
| **DN-007-2** | Meta **≥ 95% precisão por classe** (007-meta; **DN-007-2 B** etapa 3c) | respostas-joao.md §Etapa 3b; 006-2 B | Spec 006/007 apenas; **fora do escopo 008** |
| **007 formato A / caminho C** | Base rotulada **JSONL** em `datasets/labeled-logs/` (raiz do repo) | respostas-joao.md §Etapa 3b | Spec **007**; **não criar** pasta/dados na 008 |
| **Emenda constituição I** | Princípio I vs handoff ResumoAnalise (T1 rev.) | respostas-joao.md; ADR_062 | **Não editar** `.specify/memory/constitution.md` |
| **DN-006-1** | Origem de `execId` (chave N4 A) | spec 006 | Consultar antes de fechar porta de envio |
| **Query OpenObserve** | Organização, stream, índice, sintaxe de busca | `internal/adapter/out/openobserve/doc.go` | NEEDS CLARIFICATION; não inventar |
| **D1** | INFO vira ocorrência × Drive | respostas-joao.md, T7 A | Execução no Módulo; spec.md D1 |

**Resolvidas nas etapas 3b–3h (não são pendências de decisão):** **DN-008-1 B**; **DN-008-2**; **DN-008-3** (rev. 3h: **sem DLQ**; aviso Core **batchId** + **motivo**); **DN-008-4 B**; **DN-008-5 A**; **DN-008-6** (rev. 3h: Finder→Core API interna — fonte: respostas-joao.md §Etapa 3g, DLQ B); **DN-008-7 B**; **008-assinatura A**; **DN-008-8 A**; **DN-008-9 A** (**`|`** → **`\|`** — fonte: respostas-joao.md §Etapa 3g, DN-008-9 A).

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Princípio I (Finder grava ocorrência direto) | ADR_062 e T1 rev.: separação de módulos | Finder persiste MongoDB (T1 A antigo) — **substituído** 15/11 (respostas-joao.md, T1 rev.) |