# Implementation Plan: Modelo de classificação de logs

**Branch**: `006-modelo-classificacao-logs` | **Date**: 2026-11-16 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/006-modelo-classificacao-logs/spec.md` (etapa 3h: DLQ B, rev. DN-008-3/DN-008-6, DN-008-9 A; etapa 3g: DN-006-3-ref, DN-008-8 A, 007 H2 D; etapas anteriores preservadas)

## Summary

O Finder classifica cada log já capturado **sem LLM**, usando nível (campo OTel) e termos de nível no texto (palavra inteira, classe mais grave entre os dois) e, só se nenhum sinal decidir, o catálogo opcional na primeira entrada válida (fonte: spec.md FR-004, FR-005, User Story 1; respostas-joao.md sessão 2026-11-09 vigente). Classes **ERROR**, **WARNING**, **INFO** apenas (fonte: respostas-joao.md, 006-6 A, 007-9 A; origem: ADR_170, constituição III).

Logs classificados agrupam-se em **achados** pela assinatura **nível + mensagem normalizada + `service.name` + host** no escopo do lote (fonte: respostas-joao.md, Etapa 3, DN-008-1 B; origem: ADR_013). O **valor serializado** da assinatura é **texto concatenado com separador `|`** (**008 assinatura A**, **DN-008-8 A**); **`|`** literal em componente → **`\|`** (**DN-008-9 A**; canônico na spec **008** — fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40); §Etapa 3g (João, 16/11 13h22), DN-008-9 A). Falha de consulta OpenObserve (recorrência entre lotes): **DN-008-3** / **DN-008-6** rev. 3h — até 3 tentativas, **sem DLQ** no `log.ingested`; aviso ao Core **batchId** + **motivo** (**ADR_051 [PENDENTE]** — fonte: respostas-joao.md §Etapa 3g, DLQ B; respostas-etapa3g.md §1b); implementação na spec **008**, rastreada pela 006. A **mensagem normalizada** segue **DN-008-4 B** com marcadores **DN-008-5 A**: `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` (algoritmo canônico na spec **008**, **`internal/application/assinatura/normalizar_mensagem.go`** — T016; a 006 **MUST** importar/delegar — fonte: respostas-joao.md §Etapa 3b/3c, DN-008-4 B / DN-008-5 A). Antes do envio do lote, o Finder consulta o Módulo de Ocorrência **uma vez por achado** (**DN-008-7 B**, alinhado a N2 B na spec **008** — fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02)). Por achado gera-se **AchadoAnalise**; ao fim do lote válido, **exatamente um** **ResumoAnalise** (contadores ADR_102 com contador de **INFO** em lugar de «spam ignorado» + achados) segue por **HTTP interno com JWT** ao **Módulo de Ocorrência** (fonte: respostas-etapa3c.md §6; origem: ADR_062, Finder.mmd), que cria as ocorrências (fonte: respostas-joao.md, T1 REVISTO, N1 A). O Finder **não** grava ocorrência no MongoDB, **não** publica evento de ocorrência, **não** calcula gravidade (fonte: spec.md FR-007, FR-009; respostas-joao.md, 006-6 A).

**Lote incompleto (T4 A)** não conta como sucesso de lote completo; **INFO** entra no resumo como achado (T7 A, 006-11 A). Chave de idempotência do envio: **`batchId` obrigatório**; **`batchId` + `execId`** quando `execId` estiver disponível da Varredura (`V<uuid>` / `M<uuid>` — DN-006-1); no fluxo atual só `log.ingested` v1 **sem** `execId` no corpo, chave **somente `batchId`** (DN-006-1b B, N4 A). Versão única regras+catálogo+base (006-3 C) vai no achado e/ou resumo; mudança de versão só em logs novos (006-4 A).

Métricas de aceite **006-2 B**: precisão, recall e F1 **por classe** com base rotulada (007, JSONL em `datasets/labeled-logs/` — documentado); comparação ao gold **por apontamento** com classe prevista do apontamento = **mais grave** entre os logs (**007 H2 D**, ADR_170 — fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40)). **Gate formal (SC-007, T047)**: em **cada** classe ERROR, WARNING, INFO — precisão **≥ 95%** (**DN-007-2 B**), recall **≥ 95%** e F1 **≥ 95%** (**DN-006-2 B**, etapa 3e — fonte: respostas-joao.md §Etapa 3c, João 16/11 11h e 11h01; respostas-etapa3c.md §1). **Desempenho (SC-002)**: ciclo **classificar + montar resumo** MUST completar em **≤ 2 s por lote** (**DN-006-3 B** — fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02)), medido na **VM-A real** com **todos os serviços** rodando e lote de **1.000 registros** (**DN-006-3-ref** — fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40)). Latência de consultas ao Módulo (**DN-008-7 B**) e handoff HTTP (**DN-009-2**) ficam fora do escopo do «2 s» (spec.md SC-002).

Artefatos de desenho: [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md).

## Technical Context

**Language/Version**: Go **1.27** (`go.mod`, linha 3; README, “Go 1.27.0” na imagem Docker).

**Primary Dependencies (presentes em `go.mod`)**: `github.com/joho/godotenv v1.5.1`, `github.com/redis/go-redis/v9 v9.23.0` (consumo Redis Streams — fluxo upstream, não reimplementado nesta feature).

**Primary Dependencies (README — candidatas, não em `go.mod`)**: `go.mongodb.org/mongo-driver/v2`, OpenTelemetry, `golang.org/x/sync`, testcontainers — **MUST NOT** ser assumidas instaladas (constituição, Restrições Técnicas; README § Dependências candidatas). Para **esta feature**, persistência de ocorrência no MongoDB pelo Finder **não** é escopo (T1 REVISTO).

**Primary Dependencies (pendência — handoff Módulo Ocorrência)**: Cliente HTTP + JWT (**N1 A**); biblioteca JWT, path, claims e política 5xx/timeout **sem** definição neste repositório até contrato do Módulo (**DN-009-2** / respostas-etapa3.md §10) (fonte: respostas-joao.md, Etapa 3 «API Módulo: pendência no plan até o contrato»; origem: ADR_028, ADR_051, ADR_054).

**Storage**: Logs no **OpenObserve** (recuperação nas specs 003–005); deduplicação/idempotência de ingested via Redis (`internal/domain/chave.go`, TTL 25h alinhado ao contrato §8). **Sem** gravação de ocorrência pelo Finder nesta feature.

**Testing**: `make test` → `go test -race ./...` (README). Avaliação offline (007): base rotulada em **JSONL** (**007 formato A**) sob `datasets/labeled-logs/` na raiz (**007 caminho C** — só documentado; pasta não criada nesta execução). Métricas LLT-197/199: precisão, recall e F1 **≥ 95% por classe** (**DN-007-2 B**, **DN-006-2 B**, SC-007, T047) com regra **007 H2 D**. Benchmark SC-002 (2 s/lote): **VM-A real**, todos os serviços, lote **1.000** registros (**DN-006-3-ref**).

**Target Platform**: Serviço Finder no Docker Compose (README). `GET /health` em `internal/adapter/in/http/health.go`; **nenhum** endpoint HTTP público novo para classificação (contratos internos em `contracts/`).

**Project Type**: Módulo Go, arquitetura hexagonal leve (`internal/domain`, `internal/application`, `internal/adapter`, `internal/config` — README § Layout).

**Performance Goals**: **SC-002 normativo** — **≤ 2 s por lote** para classificar todos os logs do lote + montar `ResumoAnalise` (**DN-006-3 B** — fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02)). Condições de medição (**DN-006-3-ref**): **VM-A real**, **todos os serviços** em execução, lote de **1.000 registros** (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40)). Exclui latência de consulta ao Módulo por achado (**DN-008-7 B**) e envio HTTP (**DN-009-2**) (spec.md SC-002).

**Constraints**: Sem LLM (FR-002). Somente `log.ingested` **schemaVersion `"1"`** (T5 D; `internal/domain/envelope.go` rejeita ≠1). **`productId` ausente ou inconsistente** metadata×corpo ingested v1 → falha crítica permanente **`produto_ausente_ou_inconsistente`**, sem retry (**008-16 A**, ADR_049 — analyze2 **E2**). Agrupamento **entre lotes** (OpenObserve / recorrência) **fora** do mínimo de FR-009; **consulta prévia ao Módulo por achado** (**DN-008-7 B**) **é** escopo de FR-007. Catálogo: sintaxe e dono **pendentes** (006-7 D). Assinatura (**DN-008-1 B**) + valor serializado (**008 assinatura A**) + normalização **DN-008-4 B** com marcadores **DN-008-5 A** via pacote canônico **008**.

**Scale/Scope**: Um **AchadoAnalise** por achado no lote; no produto, uma ocorrência por achado no Módulo (T2 C). Chave ocorrência no módulo: assinatura + classe + produto + tenant (008-11 rodada 16/11).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Princípio | Resultado | Leitura e fonte |
|-----------|-----------|-----------------|
| I. Finalidade: logs em ocorrências | **Pendência de emenda** | Constituição I exige que o Finder agregue “em ocorrência”. **T1 REVISTO**: o Finder agrega em **ResumoAnalise/AchadoAnalise** e o **Módulo de Ocorrência** materializa a ocorrência (fonte: respostas-joao.md, T1 REVISTO; origem: ADR_062). Ver § Pendências (emenda constituição). |
| II. Origem dos logs no coletor | Passa | Classificação sobre logs já capturados via OTel/OpenObserve (spec.md FR-001; contrato-log-ingested-v1.txt §3). |
| III. Categorias de ocorrência | Passa com mapeamento explícito | Constituição: erro, alerta, informação → **ERROR, WARNING, INFO** no `tipo` (spec.md taxonomia; respostas-joao.md, 006-6 A). |
| IV. Fila Redis → OpenObserve → agregação | Passa | Apontamento/código não é conteúdo do log (`internal/domain/item.go`); fluxo preservado nas specs anteriores. Esta feature não relê a fila na porta de classificação (contratos internos). |
| V. Incerteza sem invenção | Passa com lacunas nomeadas | Catálogo (006-7 D), contrato HTTP ampliado (N3 A), API Módulo (**DN-009-2**). **DN-008-4 B**, **DN-008-5 A**, **DN-008-8 A** (separador `|`, spec 008), **DN-008-9 A**, **DN-007-2 B**, **DN-006-2 B**, **DN-006-3 B**, **DN-006-3-ref**, **007 H2 D**, **DN-008-7 B**, **DN-008-3** / **DN-008-6** rev. 3h (rastreio; runtime na 008) fechados (etapa 3h). 007-26 A: lacunas sem fonte bloqueiam **partes** afetadas — **nenhuma** bloqueia o plano inteiro nesta rodada. |

**Resultado do gate (pré-Fase 0)**: aprovado para pesquisa e desenho, com emenda constitucional I pendente e bloqueios parciais documentados.

### Revisão pós-Fase 1

[data-model.md](./data-model.md) e [contracts/](./contracts/) refletem T1 REVISTO (Finder envia `ResumoAnalise`; Módulo cria/persiste ocorrência — **não** “Finder grava MongoDB”), DN-008-1 B, **DN-008-5 A**, **008 assinatura A** (separador **`|`**, **DN-008-8 A**; escape **`|`** → **`\|`**, **DN-008-9 A**), **DN-008-3** / **DN-008-6** rev. 3h (rastreio; sem DLQ no `log.ingested`), DN-008-7 B, DN-006-3 B + perfil **DN-006-3-ref** (VM-A, 1.000 registros), gate **DN-006-2 B**, comparação avaliação **007 H2 D**, DN-006-1b, contador INFO, handoff HTTP (corpo lógico) e ausência de persistência de ocorrência no Finder. Gate mantido; emenda constitucional I continua pendente.

## Project Structure

### Documentation (this feature)

```text
specs/006-modelo-classificacao-logs/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── classificar-log.md
│   ├── classificar-apontamento.md
│   └── resumo-analise-handoff.md
└── tasks.md             # /speckit-tasks — fora deste comando
```

### Source Code (repository root) — alvo de implementação

Estrutura **existente** (README): `cmd/finder/`, `internal/domain/`, `internal/application/`, `internal/application/port/{in,out}/`, `internal/adapter/in/{queue,http}/`, `internal/adapter/out/{mongo,openobserve}/`, `internal/config/`.

**Alvo coerente com esta feature** (nomes indicativos; detalhamento em tasks):

```text
internal/domain/                              # ClasseLog, Achado, assinatura (valor), classificação — sem normalizador duplicado
internal/application/                         # Orquestrar classificar → agrupar (lote) → montar ResumoAnalise
internal/application/assinatura/              # Dono canônico spec 008: normalizar_mensagem.go (DN-008-4 B) — T016 na 008
internal/application/port/in/                 # Portas: classificar lote / apontamento
internal/application/port/out/                # Portas: consultar Módulo por achado (DN-008-7 B); enviar ResumoAnalise — HTTP BLOQUEADO até DN-009-2
internal/adapter/out/modocorrencia/           # BLOQUEADO (DN-009-2): URL, JWT, códigos; consulta por achado depende de rota do Módulo
```

**Structure Decision**: estender o hexágono existente; **não** adicionar `adapter/out/mongo` para ocorrência nesta feature (fonte: respostas-joao.md, T1 REVISTO). **Normalização**: a 006 **reutiliza** `NormalizarMensagem` do caminho **`internal/application/assinatura/normalizar_mensagem.go`** (definido na spec **008**); **proibido** segunda implementação em `internal/domain/` ou outro pacote divergente (fonte: respostas-joao.md §Etapa 3b, DN-008-4 B; analyze2 **F1**).

## API do Módulo Ocorrência (rota, JWT, respostas HTTP, 5xx/timeout)

Pendência **DN-009-2** / `respostas-etapa3.md` §10; decisão João: **«API Módulo: pendência no plan até o contrato»** (respostas-joao.md, §Etapa 3). Não bloqueia o plano da 006 como um todo; bloqueia **somente** o adaptador HTTP e a confirmação de varredura que depende da resposta do Módulo.

| Tópico | O que há fonte | O que falta (sem inventar) |
|--------|----------------|----------------------------|
| Transporte | HTTP interno + JWT (**N1 A**); rede privada (ADR_141); JWT próprio (ADR_028, ADR_054) | Rota/path Finder → Módulo; claims; algoritmo/biblioteca JWT ([PENDENTE] ADR_028) |
| Semântica | Entregar `ResumoAnalise` ao fim do lote; idempotência no Módulo (ADR_014); Finder só conclui após confirmação (seq. v6 §2) | Corpo da resposta de confirmação; mapeamento HTTP ↔ FA-02 |
| Respostas HTTP | Idempotência Drive (ADR_014); 009-18 A (3 tentativas + alerta) vale a **publicação** de ocorrência no ecossistema, não copiada automaticamente para este handoff | Códigos 2xx/4xx/5xx; uso de 202 vs 200; timeout; lista retentável (ADR_118 é OTLP Coletor, não API interna) |
| Contrato corpo | JSON lógico em [contracts/resumo-analise-handoff.md](./contracts/resumo-analise-handoff.md) | Formato de transporte da chave (header vs corpo) |

**Tarefas de implementação afetadas**: porta `Enviar ResumoAnalise`, pacote `internal/adapter/out/modocorrencia/`, quickstart §5 integração real, SC-002 que exija confirmação HTTP — permanecem **BLOQUEADAS** até contrato publicado pelo Módulo Ocorrência (ou ADR dedicada citada em fonte).

### SC-002 e aceite parcial (analyze **E1**)

| Escopo de validação | O que prova | Bloqueio |
|---------------------|-------------|----------|
| **Parcial (desbloqueado)** | Classificar + agrupar (DN-008-1 B, **DN-008-4 B**, **DN-008-5 A**, **008 assinatura A**) + montar **um** `ResumoAnalise` em memória; contadores (incl. `info`); tenant/produto/versão; zero persistência/publicação pelo Finder; consulta Módulo **simulável** por achado (**DN-008-7 B**) | — |
| **Parcial (desempenho)** | Medir tempo classificar+montar resumo ≤ 2 s/lote (**DN-006-3 B**) no perfil **DN-006-3-ref** (VM-A real, todos os serviços, 1.000 registros) | Reproduzir ambiente completo pode ser parcial em dev local |
| **Parcial (fake/adaptador)** | Mesmo payload + simulação de confirmação na porta de saída **sem** contrato HTTP real (tasks T024/T032) | Não substitui aceite E2E |
| **Completo (SC-002 normativo)** | Confirmação HTTP do Módulo após envio JWT real, com semântica 2xx/4xx/5xx/timeout fechada | **DN-009-2** |

Sucesso de validação **não** deve ser lido como handoff produção enquanto **DN-009-2** estiver aberto (fonte: respostas-joao.md §Etapa 3; respostas-etapa3.md §10; quickstart §5).

## Tarefas bloqueadas (dentro do escopo 006)

| Bloqueio | Escopo afetado | Restante desbloqueado |
|----------|----------------|----------------------|
| **DN-009-2** | Adaptador HTTP+JWT; URL; política retry/timeout/5xx; validação SC-002 E2E com Módulo real; rota da **consulta por achado** (**DN-008-7 B**) | Classificar log; agrupar (DN-008-1 B + **DN-008-4 B** + **DN-008-5 A**); montar **um** `ResumoAnalise`; consulta Módulo com fake/stub |
| **006-7 D** | Cenários e testes que exijam sintaxe de catálogo concreta | Regras nível+texto; catálogo opcional opaco na porta |
| **N3 A** | Nomes JSON finais e obrigatoriedade dos campos ampliados até ADR | Requisitos funcionais e placeholders no contrato |

## Decisões tomadas (etapa 3h)

| ID | Decisão | Fonte | Nota |
|----|---------|-------|------|
| **DLQ B** | Sem DLQ no `log.ingested`; após 3 falhas OpenObserve, Finder avisa o Core (API interna, **batchId** + **motivo**) | (fonte: respostas-joao.md, DLQ B; origem: §Etapa 3g (João, 16/11 13h22); respostas-etapa3g.md §1b) | Contrato API Core **ADR_051 [PENDENTE]** |
| **DN-008-3** (rev. 3h) | **Substitui** alcance anterior (3 tentativas + DLQ no ingested) | (fonte: respostas-joao.md, DLQ B; origem: §Etapa 3g) | Implementação principal na spec **008** |
| **DN-008-6** (rev. 3h) | **Substitui** etapa 3c–3e (DLQ + alerta Core/WebSocket pós-DLQ); mesmo gatilho via API interna | (fonte: respostas-joao.md, DLQ B; origem: §Etapa 3g) | Distinto da publicação de evento de ocorrência (spec **009**) |
| **DN-008-9 A** | Escape **`|`** → **`\|`** nos componentes do `valor` antes da concatenação | (fonte: respostas-joao.md, DN-008-9 A; origem: §Etapa 3g (João, 16/11 13h22)) | Escape **`\`**: só pendência (sem fonte) |

## Pendências

Itens **fora desta execução** ou que exigem artefato externo (sem criar ADR nesta pasta):

| ID | Descrição | Fonte João / processo | Ação neste repo |
|----|-----------|----------------------|-----------------|
| **DN-009-2** | API Módulo: rota, JWT, HTTP, 5xx/timeout | respostas-joao.md §Etapa 3; respostas-etapa3.md §10 | § API do Módulo acima; **sem** fechar contrato HTTP |
| **006-7 D** | Sintaxe e dono do catálogo (ADR_167) | respostas-joao.md, 006-7 D | Não criar catálogo nem ADR |
| **006-10 A** | Alinhar Jira LLT-197 / LLT-199 com spec 006 | respostas-joao.md, 006-10 A | Não alterar Jira |
| **N3 A** | Formalizar ampliação `AchadoAnalise` em ADR | respostas-joao.md, N3 A | Requisitos na spec; contrato marca [PENDENTE ADR] |
| **T5 D (commit Coletor)** | Publicação `schemaVersion` "1" no repositório collector | respostas-joao.md, T5 D | Finder já rejeita ≠1; commit no collector fora deste repo |
| **010-2 A** | Lista oficial de padrões de mascaramento em ADR | respostas-joao.md, 010-2 A | Spec 010 |
| **010-27 A** | Casos de teste do mascaramento (010) e borda do normalizador (complementa **DN-008-5 A** na 008) | respostas-joao.md, 010-27 A | Specs 010 / 008 |
| **Emenda constituição I** | Princípio I vs handoff ResumoAnalise (T1 REVISTO) | respostas-joao.md, T1 REVISTO; ADR_062 | **Não editar** `.specify/memory/constitution.md` nesta execução |
| **Contrato API Core (alertas)** | Rota, payload e evento Finder→Core após 3 falhas OpenObserve (**batchId** + **motivo**) | respostas-joao.md §Etapa 3g, DLQ B; **ADR_051 [PENDENTE]** | Spec **008** implementa política; adaptador HTTP **MUST NOT** inventar rota/payload/evento |
| **Escape `\` na assinatura** | Barra invertida literal nos componentes do `valor` | Sem fonte | Complementa **DN-008-9 A** |

## Divergências documentadas

Outras divergências (D1, D5, D12, …): ver [spec.md](./spec.md) § Divergências documentadas.

## Complexity Tracking

> Preenchido porque o princípio I da constituição diverge do T1 REVISTO até emenda.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Princípio I (Finder “em ocorrência” direta) | ADR_062 e T1 REVISTO: separação de módulos; Finder envia análise, Módulo persiste ocorrência | Gravar ocorrência no Finder (T1 A antigo) foi **substituído** pelo João (fonte: respostas-joao.md, T1 REVISTO) |