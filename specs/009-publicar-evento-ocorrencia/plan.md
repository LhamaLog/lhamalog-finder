# Implementation Plan: Publicar evento de ocorrência (escopo Finder — handoff ResumoAnalise)

**Branch**: `009-publicar-evento-ocorrencia` | **Date**: 2026-11-16 (etapa 3h) | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/009-publicar-evento-ocorrencia/spec.md` (clarificado 2026-11-16 etapas 3b–3h; T1 revisto; aceite **SC-002** bloqueado até spec **008** — **009-22 A**).

## Summary

Após o T1 revisto (15/11 23h47), o **Finder não publica** o evento lógico da **LLT-201** (nova/atualizada/recorrente). Sua contribuição é o **handoff** **ResumoAnalise** → **Módulo Ocorrência** por **HTTP interno com JWT** (**N1 A**), com idempotência **`batchId`** no fluxo atual **`log.ingested` v1** (**DN-006-1b B**, **T5 D**) ou **`batchId` + `execId`** quando o identificador de varredura estiver disponível (**N4 A**, **DN-006-1**), achados **INFO** incluídos (**T7 A**), confirmação obrigatória antes de concluir a varredura (**FR-F05**, sequência varredura v6 §2) e **proibição** de chamar Core ou stream de notificação (**FR-F02**, **T1 revisto**). O Módulo Ocorrência cria ocorrências, calcula **gravidade**, envia **cada** ocorrência ao Core (**HTTP + JWT**, **009-3 D**) com **id**, **tipo** (classe) e **gravidade** (**009-5 revisto**); o Core publica em **Redis Streams** (`notificacao-eventos`, envelope **ADR_146**) para a Notificação (**009-1 B**), que lê trechos mascarados da Ocorrência pelo **id** no MongoDB compartilhado (**DN-009-1 A**). Este plano documenta handoff Finder, entidades de mensagem e pipeline downstream como **referência de aceite integrado**; a implementação do **adaptador HTTP** Finder → Módulo permanece **bloqueada** até **DN-009-2** (contrato do Módulo), sem bloquear o restante do desenho nem das tarefas de montagem lógica do ResumoAnalise.

Artefatos: [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md).

**Métricas 006/007 e base rotulada (normativa na spec)**: ver [spec.md](./spec.md) §Decisões tomadas (**DN-007-2 B**, 007 formato **A**, caminho **C**) e Clarifications etapas 3c–3h — este plano só referencia; detalhes em [research.md](./research.md) R12 (**corrige analyze D2**). **Não** criar pasta `datasets/` nem dados.

**FR-F07 / DN-008-3 / DN-008-6 (rastreio de plano, etapa 3h)**: falha transitória na **consulta OpenObserve** durante a análise: até **3 tentativas**; **sem DLQ** no **log.ingested v1** (**DLQ B**, contrato 08/11). Após esgotar, o **Finder** avisa o **Core** pela **API interna** (**`batchId`**, **motivo**); contrato HTTP **[PENDENTE ADR_051]** — **não** consumidor DLQ nem **WebSocket** pós-DLQ (substitui §Etapa 3c). Sem handoff **ResumoAnalise** até análise concluída; coordenação com spec **008**/**005** sem inventar rota (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22); respostas-etapa3g.md §1; origem: T3 B, DN-006-1b B).

### Decisões tomadas (etapa 3h — propagadas de spec.md)

| ID | Tópico | Fonte |
|----|--------|-------|
| **DLQ B** | **Sem DLQ** no **log.ingested v1**; não copiar retry/DLQ da FilaOrdenada para o ingested | (fonte: respostas-joao.md, DLQ B; origem: respostas-etapa3g.md §1b, contrato log.ingested v1, DN-006-1b B) |
| **DN-008-3** (rev.) | Falha consulta **OpenObserve**: **3 tentativas**; depois Finder avisa **Core** (API interna, **`batchId`**, **motivo**); **sem DLQ** — **substitui** alcance «3 tentativas + DLQ» (§Etapa 3c) | (fonte: respostas-joao.md, DLQ B; origem: respostas-etapa3g.md §1) |
| **DN-008-6** (rev.) | Após esgotar tentativas OpenObserve: aviso **Finder → Core** (API interna); **não** WebSocket pós-DLQ — **substitui** §Etapa 3c/3e | (fonte: respostas-joao.md, DLQ B; origem: respostas-etapa3g.md §1, ADR_051 [PENDENTE], FR-F07) |
| **DN-008-9 A** | **`|`** literal no componente da assinatura → **`\|`** na serialização | (fonte: respostas-joao.md, DN-008-9 A; origem: DN-008-8 A, 008 assinatura A) |

## Technical Context

**Language/Version**: Go **1.27** (fonte: `go.mod` linha 3; README § Como iniciar, imagem Docker Go 1.27.0).

**Primary Dependencies (presentes em `go.mod`)**: `github.com/joho/godotenv`, `github.com/redis/go-redis/v9` — consumo de fila/stream de apontamentos (constituição IV; specs 004/005). **Não** há adaptador de saída HTTP/JWT para handoff: `net/http` aparece só em health (`internal/adapter/in/http/health.go`, `cmd/finder/`) — sem cliente ao Módulo Ocorrência (fonte: `go.mod`; grep `*.go`).

**Primary Dependencies (README — candidatas, não em `go.mod`)**: `go.mongodb.org/mongo-driver/v2`, OpenTelemetry, `golang.org/x/sync`, testcontainers — **MUST NOT** ser tratadas como instaladas (README linhas 146–155; constituição, Restrições Técnicas). Handoff ResumoAnalise exigirá adaptador HTTP + JWT (**bloqueio de implementação**: **DN-009-2**; padrão de produto: **N1 A**, **ADR_028** — sem biblioteca nem variáveis nomeadas neste repo).

**Storage**: O Finder lê referências **log.ingested v1** (`schemaVersion` `"1"` — **T5 D**, `internal/domain/envelope.go` linhas 73–75) e conteúdo via **OpenObserve** (`internal/adapter/out/openobserve/doc.go`). Stub Mongo em `internal/adapter/out/mongo/doc.go`. **Não** persiste ocorrência de negócio (**T1 revisto**, **FR-F02**; fonte: ADR_062).

**Testing**: `make test` → `go test -race ./...` (README linhas 103–104). Aceite **SC-F001–SC-F003** neste repo foca handoff e classes; **SC-001–SC-003** (LLT-201 no produto) exigem ambiente Módulo + Core + Notificação ([quickstart.md](./quickstart.md)). Equivalente sequência v6 **`finishScan()`** citado em **FR-F05** **não** aparece no código Go atual (grep `finishScan` em `*.go` vazio) — porta de conclusão de varredura **a nomear na implementação**, sem inventar identificador.

**Target Platform**: Serviço Finder no Docker Compose (`make run`; README). Health em `GET /health` (`internal/adapter/in/http/health.go`).

**Project Type**: Módulo Go, arquitetura hexagonal leve (`cmd/finder/`, `internal/domain`, `internal/application`, `internal/adapter`, `internal/config` — README § Layout).

**Performance Goals**: Sem meta numérica com fonte para latência do handoff ou do pipeline de aviso (LLT-201 não impõe SLA nas fontes prioritárias).

**Constraints**: Idempotência handoff **`batchId`** (**DN-006-1b B**, contrato 08/11) ou **`batchId` + `execId`** quando houver (**N4 A**, **DN-006-1**). **Um** **ResumoAnalise** por lote (**DN-008-7**); consulta prévia ao Módulo **por achado** (**DN-008-7 B**, **N2 B**, **FR-F09**). Lote incompleto não conta como sucesso (**T4 A**). Assinatura: composição **DN-008-1 B**, normalização **DN-008-4 B**, marcadores **DN-008-5 A** (`<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` — canônico spec **008**), serialização **008 assinatura A** (concatenação com separador fixo **`|`** — **DN-008-8 A**; **`|`** nos componentes → **`\|`** — **DN-008-9 A**; escape de **`\`** **[PENDENTE]** sem fonte; dono canônico spec **008**; **FR-F11**). Falha consulta OpenObserve: **FR-F07** / **DLQ B** / **DN-008-3**/**DN-008-6** (rev.). Extensões **N3 A** aguardam ADR. Finder **não** envia gravidade nem publica envelope LLT-201 (**ADR_062**, **FR-F02**). Gate **≥95% precisão por classe** (**DN-007-2 B**) na spec **006**/**007**, não nesta spec. Taxonomia **ERROR/WARNING/INFO** (**009-13 A**, constituição III; **FR-F10**).

**Scale/Scope**: Depende de classificação (spec **006**) e agrupamento (spec **008**). Estados **atualizada/recorrente** no aviso downstream dependem da **008** (**009-22 A**, **SC-002** bloqueado).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Princípio | Resultado | Leitura e fonte |
|-----------|-----------|-----------------|
| I. Finalidade: logs em ocorrências | **Pendência de emenda** | Constituição diz que o Finder agrega em **ocorrência**; T1 revisto: materialização no **Módulo Ocorrência** via **ResumoAnalise** (fonte: respostas-joao.md, T1 revisto; origem: ADR_062, revisao-T1-drive.md §1). **Não** editar constituição nesta pasta — ver § Pendências. |
| II. Origem dos logs no coletor | Passa | Análise sobre logs do Coletor via OTel/OpenObserve (spec 006/005; constituição II). |
| III. Categorias de ocorrência | Passa | `tipo` em AchadoAnalise: erro, alerta, informação → ERROR, WARNING, INFO (**FR-F08**, **009-13 A**). |
| IV. Fila Redis → OpenObserve → agregação | Passa | Fluxo upstream inalterado; handoff é **pós-análise** do lote (revisao-T1-drive.md §2; sequência varredura v6 §2). |
| V. Incerteza sem invenção | Passa com lacunas explícitas | **DN-009-2** e **ADR_051 [PENDENTE]** (Finder→Core pós-OpenObserve) registradas sem suposição; etapas 3b–3h fecharam **DN-006-1**/**1b**, **DN-007-2 B**, **DN-008-4 B**/**DN-008-5 A**, **DLQ B**, **DN-008-3**/**DN-008-6** (rev.), **DN-008-9 A**, **DN-008-7**/**DN-008-7 B**, **008 assinatura A**, **DN-008-8 A** (fonte: respostas-joao.md §Etapa 3c–3g; **DLQ B**/**DN-008-9 A**: §Etapa 3g 16/11 13h22; respostas-etapa3g.md §1). Nenhuma dúvida **nova** sem registro em spec.md bloqueia o plano como um todo. |

**Resultado do gate (pré-Fase 0)**: aprovado para pesquisa e desenho com pendências documentadas.

### Revisão pós-Fase 1

[data-model.md](./data-model.md) separa entidades Finder (handoff) vs produto (aviso). [contracts/resumo-analise-handoff.md](./contracts/resumo-analise-handoff.md) é a interface implementável neste repo (**implementação HTTP bloqueada** por **DN-009-2**). [contracts/pipeline-aviso-referencia.md](./contracts/pipeline-aviso-referencia.md) documenta Módulo → Core → Notificação sem inventar URL HTTP Core nem campos além das fontes. Gate mantido.

## API do Módulo Ocorrência (pendência DN-009-2)

Contrato HTTP **Finder → Módulo Ocorrência** (rota/path, claims JWT, códigos de resposta, timeout e semântica 5xx) **não** está definido no Drive nem no Jira além do padrão genérico (fonte: respostas-joao.md §Etapa 3 — *API Módulo: pendência no plan até o contrato*; origem: respostas-etapa3.md §10).

**O que as fontes fixam hoje**

| Aspecto | Valor documentado | Fonte |
|---------|-------------------|-------|
| Transporte | HTTP(S) interno, rede privada entre módulos | N1 A (fonte: respostas-joao.md, N1 A; origem: ADR_141) |
| Autenticação | JWT próprio por módulo (não Clerk/M2M) | ADR_028, ADR_054 (respostas-etapa3.md §10) |
| Granularidade envio | **Exatamente um** **ResumoAnalise** por lote analisado — não HTTP por achado | **DN-008-7** (fonte: respostas-etapa3c.md §6; origem: ADR_062, Finder.mmd) |
| Granularidade consulta (**N2 B**) | **Uma** consulta ao Módulo **por achado** (assinatura/classe elegível) antes do envio do resumo | **DN-008-7 B** (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); origem: N2 B) |
| Semântica | Confirmação antes de concluir varredura; idempotência ADR_014 | seq. v6 §2; ADR_014 |
| Falha publicação (paralelo) | **009-18 A**: tentar de novo e alertar — aplicável ao trecho Módulo→Core; para Finder→Módulo, FA-02 (reenvio mesma chave) | respostas-joao.md, 009-18 A |
| Rota conhecida (outro uso) | `POST /ocorrencias/{id}/sugestoes` (Sugestão manual) — **não** é handoff ResumoAnalise | ADR_083, ADR_092 (respostas-etapa3.md §10) |

**O que permanece em aberto (DN-009-2)**

- Path, método exato, headers e corpo de confirmação para **POST ResumoAnalise** e para **consulta por achado** (**DN-008-7 B**).
- Algoritmo/rotacao JWT (**[PENDENTE]** ADR_028).
- Mapeamento 4xx/5xx/timeout e política de retry além do princípio FA-02 (opções 10a–10d em spec.md §DN-009-2).

**Impacto no plano**: tarefas de **adaptador HTTP** (`internal/adapter/out/modocorrencia/`, cliente JWT, **entrega do resumo** e **consulta por achado**, testes contra URL real) **bloqueadas** até contrato/ADR do Módulo Ocorrência. Tarefas de **modelagem ResumoAnalise**, portas de aplicação, validação de taxonomia, normalizador **FR-F11** e documentação **seguem**.

## Project Structure

### Documentation (this feature)

```text
specs/009-publicar-evento-ocorrencia/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── resumo-analise-handoff.md      # Finder → Módulo (escopo repo)
│   └── pipeline-aviso-referencia.md   # Módulo → Core → Notificação (referência LLT-201)
└── tasks.md                           # /speckit-tasks — fora deste comando
```

### Source Code (repository root) — alvo de implementação

Estrutura **existente** (README § Layout):

```text
cmd/finder/
internal/domain/                 # envelope v1, Referencia.BatchID, chaves
internal/application/            # orquestração consumo → análise → (futuro) handoff
internal/application/port/out/   # porta de saída: entregar ResumoAnalise (a criar)
internal/adapter/in/queue/       # Redis Streams (entrada)
internal/adapter/in/http/        # health
internal/adapter/out/openobserve/
internal/adapter/out/mongo/      # doc.go — stub
internal/adapter/out/modocorrencia/  # bloqueado por DN-009-2 — adaptador HTTP+JWT (não existe)
```

**Structure Decision**: uma porta de saída **ModuloOcorrencia** (entrega + confirmação); **sem** adaptador Core/Notificação no Finder (**FR-F02**; fonte: spec.md, T1 revisto).

## Pendências

Itens **fora desta execução**, **outro repositório**, ou que exigem artefato externo (**sem** criar ADR nesta pasta):

| ID | Descrição | Fonte | Efeito nesta feature |
|----|-----------|-------|----------------------|
| **DN-009-2** | Contrato HTTP Finder → Módulo (rota, JWT, respostas, 5xx/timeout) | respostas-joao.md §Etapa 3; respostas-etapa3.md §10 | **Bloqueia** adaptador HTTP (**FR-F01** wiring real) |
| **N3 A** | ADR nova: ampliar `AchadoAnalise` | respostas-joao.md, N3 A | Campos [PENDENTE ADR] no contrato handoff |
| **006-10 A** | Alinhar Jira LLT-201/200 com specs 006/008/009 | respostas-joao.md, 006-10 A | Fora desta pasta |
| **010-2 A** | Lista de padrões de mascaramento em ADR | respostas-joao.md, 010-2 A | Spec 010 |
| **010-27 A** | Casos de teste de mascaramento (João) | respostas-joao.md, 010-27 A | Spec 010 |
| **008** | Aceite **SC-002** / estados atualizada e recorrente | respostas-joao.md, 009-22 A | Bloqueio de **aceite formal**, não de desenho do handoff |
| **C8** | Chave concreta da ocorrência [PENDENTE] ADR_040 | respostas-documentadas.md, 009-6; spec.md | Id opaco permitido no aviso |
| **Separador assinatura** | Literal **`|`** em **008 assinatura A** (**DN-008-8 A**) | respostas-joao.md §Etapa 3e (João, 16/11 12h40), DN-008-8 A; dono canônico spec **008** | Propagado para reuso (**FR-F11**); não inventar outro literal |
| **Emenda constituição I** | “Agregar em ocorrência” vs handoff | constituição L24–27; T1 revisto | Pendência de emenda, não editar aqui |
| **T6 A** | Lacunas da spec 006 fechadas na própria 006 | respostas-joao.md, T6 A | Spec 006 |
| **D3** | INFO vira ocorrência (João) × SSID_007.1/009.1 | respostas-joao.md, T7 A | ResumoAnalise inclui INFO; Módulo executa |
| **ADR_028** | Algoritmo, rotação e variáveis JWT interno [PENDENTE] | respostas-etapa3.md §10 | Paralelo a **DN-009-2** |
| **Templates NOT-03** | Corpo e-mail Notificação [PENDENTE] | respostas-etapa3.md §7 | Fora do Finder; origem de dados **DN-009-1 A** fechada |
| **ADR_051 [PENDENTE]** | API **Finder → Core** após falha OpenObserve esgotada (**`batchId`**, **motivo**) — rota, payload, evento | respostas-joao.md §Etapa 3g, DLQ B; spec.md **FR-F07** | Implementação aviso operacional (**T049**); sem inventar contrato |
| **Escape `\` na assinatura** | Sem fonte; vale **DN-008-9 A** (**`|`** → **`\|`**) apenas | respostas-etapa3g.md §2; spec.md §Questões em aberto | Não fixar regra de escape da barra invertida até fonte |

## Tarefas bloqueadas vs liberadas (resumo)

| Área | Liberado para plano/tarefas | Bloqueado até pendência |
|------|-----------------------------|-------------------------|
| Modelo lógico ResumoAnalise / AchadoAnalise | Sim (com lacunas **N3 A**) | Campos N3 até ADR |
| Porta `ModuloOcorrencia` (interface) | Sim | — |
| Adaptador HTTP+JWT + env vars | Não | **DN-009-2**, **ADR_028** |
| Agrupamento entre lotes (**N2 B**, **DN-008-7 B**) e normalização (**FR-F11** / **DN-008-4 B** / **DN-008-5 A** / **008 assinatura A** / **DN-008-8 A**) | Desenho, normalizador com marcadores fechados, separador **`|`** e **uma** consulta **por achado** (porta lógica) podem avançar **sem** ADR **N3 A** | Chamadas HTTP reais de consulta e handoff até **DN-009-2**; campo estruturado `assinatura` no payload até ADR **N3 A** |
| Pipeline Módulo→Core→Notificação (doc) | Sim (referência) | URLs Core |

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Princípio I (Finder materializa ocorrência localmente) | ADR_062 + T1 revisto: separação de módulos e SSID_009.1 | Finder publica evento e grava ocorrência (T1 A antigo) — **substituído** (fonte: respostas-joao.md, T1 revisto; revisao-T1-drive.md §1) |
| Épico LLT-201 no repositório Finder | Jira coloca “publicar evento” no épico Finder | Documentar handoff + pipeline de referência; aceite LLT-201 integrado fora do binário Finder (**D1**, spec.md) |