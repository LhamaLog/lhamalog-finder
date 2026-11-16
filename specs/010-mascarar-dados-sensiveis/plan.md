# Implementation Plan: Mascarar dados sensíveis antes de gravar a ocorrência

**Branch**: `010-mascarar-dados-sensiveis` | **Date**: 2026-11-16 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/010-mascarar-dados-sensiveis/spec.md` (clarificação 010-1–010-28, T1 revisto; etapa 3b: DN-010-1/DN-010-2, 009-5 revisto; etapa 3c: DN-007-2 B, DN-008-4 B, 007 formato A / caminho C; etapa 3e: **DN-010-3 B**, **DN-010-4 C**, **DN-008-5 A**; etapa 3h: **010 H1**)

## Summary

Antes do **MotorAnalise** / classificação, o Finder aplica **`validarSeguranca()`** uma vez por lote em **`montarContexto()`**, linha a linha, sobre o texto recuperado do OpenObserve (fonte: spec.md FR-002, FR-006; respostas-documentadas.md, 010-3 A / 010-4 A / 010-14 C / 010-15 A; origem: ADR_064, ADR_016). O **classificador** e o restante do pipeline no Finder operam só sobre texto **mascarado** (010-5 A). **E-mails** viram placeholder fixo (010-11 A); linhas com **senha/token/secret** são trocadas pelo marcador **`"log com senha omitido"`** (010-10 A). O total numérico **`dadosSensiveisIdentificados`** (retorno de `validarSeguranca()` — ADR_064) MUST aparecer nos **contadores** do `ResumoAnalise` e do resumo do arquivo com rótulo humano **“dados sensíveis mascarados”** (010-17; **I2**). Chave JSON no handoff `[INFERIDO]` como `contadores.dadosSensiveisMascarados` até ADR_102/N3 fecharem — **não** confundir com campo de achado **`dadosSensiveis`** no Drive (fonte: spec.md FR-005, Clarifications I2; revisao-T1-drive.md §2). Trechos guardados saem em **`AchadoAnalise.linhas`** já mascarados; o **Módulo de Ocorrência** persiste em **`conteudo`** (010-6, 010-25; T1 revisto). Lista de regex em **arquivo versionado** compartilhado com Sugestão, **mesma versão** obrigatória (010-18 A, 010-19 A; ADR_079) — conteúdo concreto **PENDENTE ADR** (010-2 A). **Caminho** no monorepo e **formato** do artefato MUST seguir a **ADR da lista** (**DN-010-4 C**, **P-010-2**) — esta spec **não** fixa path/formato. **Mismatch** entre versão lida e versão **combinada** (deploy): serviço **recusa subir** na inicialização (**DN-010-3 B**, **FR-008**). Falha de `validarSeguranca()`: continuar **sem** incluir trechos no achado (010-16 C). **Sem** meta de desempenho (010-28 A). Escopo **só linhas de log** (**DN-010-2**); atributos OTel sensíveis são **Cliente na origem** (respostas-etapa3.md §8). E-mail LLT-218: Notificação lê Ocorrência **pelo id** (**DN-010-1 A** / **DN-009-1 A**). Aviso Módulo → Core: **id + tipo + gravidade**, sem texto de log (**009-5 revisto**; mantém 010-23).

Ordem com **normalização de mensagem** para assinatura (**DN-008-4 B**, dona spec **008**): `validarSeguranca()` MUST preceder classificação e qualquer uso da mensagem na assinatura; texto de entrada à normalização MUST ser linha **já mascarada** (**FR-015**, **010-5 A**). A spec **010** MUST **reutilizar** o normalizador canônico entregue pela spec **008** em `internal/application/assinatura/normalizar_mensagem.go` (**T016** em `specs/008-agrupar-logs-ocorrencia/tasks.md`) — **MUST NOT** criar outro módulo/algoritmo de normalização nesta feature (**010 H1**, etapa 3h). (fonte: respostas-joao.md, 010 H1; origem: §Etapa 3g, João 16/11 13h22, DN-008-4 B). Literais dos marcadores de normalização: **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`** (**DN-008-5 A**) — referência cruzada na spec **008**; distintos do marcador de mascaramento **`"log com senha omitido"`** (010-10 A). Meta **≥95% de precisão por classe** (ERROR, WARNING, INFO) é gate das specs **006**/**007** (**DN-007-2 B**; **007-meta**), **não** aceite LLT-202/010.

Artefatos: [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md). Decisões 010-* canônicas: [spec.md § Decisões tomadas](./spec.md#decisões-tomadas-com-fonte) (analyze **D1**).

## Technical Context

**Language/Version**: Go **1.27** (`go.mod`, linha 3; README, Go 1.27.0 na imagem Docker).

**Primary Dependencies (presentes em `go.mod`)**: `github.com/joho/godotenv`, `github.com/redis/go-redis/v9` (consumo Redis Streams — upstream; mascaramento não altera o contrato da fila).

**Primary Dependencies (README — candidatas, não em `go.mod`)**: `go.mongodb.org/mongo-driver/v2`, OpenTelemetry, `golang.org/x/sync`, testcontainers — **MUST NOT** ser assumidas instaladas (constituição, Restrições Técnicas; README § Dependências candidatas, linhas 146–155).

**Primary Dependencies (mascaramento)**: Motor de regex — **NEEDS CLARIFICATION**: biblioteca específica (`regexp` stdlib vs. extensão) não decidida nas três fontes prioritárias; ADR_064 cita “regex” sem pacote Go (fonte: respostas-documentadas.md, 010-1; origem: ADR_064 [PENDENTE]).

**Storage**: Conteúdo bruto **somente** no OpenObserve (010-24 C; `contrato-log-ingested-v1.txt` §2, §5; ADR_146 v2). Finder **não** reexpõe bruto em filas, APIs ou handoff. Artefato de catálogo: **arquivo versionado no repositório** (010-18 A); **caminho** e **formato** declarados na **ADR da lista** (**DN-010-4 C**, **010-2 A**) — implementação do adaptador **aguarda P-010-2** se a ADR ainda não publicar metadados (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h, DN-010-4 C; origem: ADR_079, 006-3 C).

**Testing**: `make test` → `go test -race ./...` (README). Casos de aceite SC-001/002 **bloqueados** até casos do João (010-27 A) e ADR do catálogo (010-2 A) — spec.md SC-005.

**Target Platform**: Serviço Finder no Docker Compose (README). `GET /health` existente em `internal/adapter/in/http/health.go`.

**Project Type**: Módulo Go, arquitetura hexagonal leve (`cmd/finder/`, `internal/domain`, `internal/application`, `internal/adapter`, `internal/config` — README § Layout, linhas 132–141).

**Performance Goals**: **Sem meta** nesta feature (010-28 A; fonte: respostas-joao.md, 010-28 A).

**Constraints**: Uma passada por lote (010-15 A). PII fora do catálogo permanece **bruto** (010-13 A; DPA). Não emendar ADR_019, ADR_167, SSID_007.1 (FR-013). Sugestão (Java) fora deste repositório — só contrato do artefato compartilhado (FR-008). **DN-010-2**: mascaramento **só linhas**; identificadores ADR_040 (`traceId`, `correlationId`, `hostName`, `moduloId`) copiados como enviados, fora do catálogo linha (spec.md FR-004a).

**Scale/Scope**: Implementação atual cobre consumo Redis e domínio de fila (`internal/application/consumir.go`, `internal/domain/item.go`); **não há** `ScannerSegurança`, `montarContexto()` nem `ResumoAnalise` no código Go observado — escopo desta feature é desenhar e depois implementar a etapa no hexágono (fonte: estado do repo; alinhado spec 009 research R8).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Princípio | Resultado | Leitura e fonte |
|-----------|-----------|-----------------|
| I. Finalidade: logs em ocorrências | **Pendência de emenda** (inalterada por esta feature) | Constituição I: Finder agrega “em ocorrência”. **T1 revisto**: materialização no Módulo via `ResumoAnalise`; mascaramento garante trechos em `AchadoAnalise.linhas` (fonte: respostas-joao.md, T1 revisto 15/11; spec.md FR-003). **P-emenda-I** em spec § Pendências (analyze **C1**); **não** editar `.specify/memory/constitution.md` nesta pasta. |
| II. Origem dos logs no coletor | Passa | Mascaramento na **processamento** após leitura OpenObserve; captura entrega registro inteiro (FR-012; spec 003). |
| III. Categorias de ocorrência | Passa com **pendência de emenda** | Mascaramento independente de ERROR/WARNING/INFO; classificador usa texto mascarado (010-5 A). Taxonomia implementada ERROR/WARNING/INFO (ADR_170) vs texto literal da constituição (erro/alerta/informação) — mitigado em spec **D8**/FR-013; **P-emenda-III** fora desta pasta (analyze **G1**). |
| IV. Fila Redis → OpenObserve → agregação | Passa | Fila carrega só apontamento (`internal/domain/item.go`, `Aceito()`); bruto não trafega na fila (010-24; contrato §2). |
| V. Incerteza sem invenção | Passa com lacunas explícitas | **DN-010-1**, **DN-010-2**, **DN-007-2**, **DN-008-4**, **DN-010-3 B**, **DN-010-4 C**, **DN-008-5 A** **respondidos** (etapas 3b/3c/3e); **010 H1** (reuso normalizador **008**, etapa 3h). Valores concretos de **path/formato** do catálogo e regex aguardam **publicação da ADR** (**P-010-2**, não inventar path). Placeholder exato de e-mail e biblioteca regex Go **NEEDS CLARIFICATION** até ADR — ver § Pendências. **DN-009-2**: contrato HTTP — § API do Módulo; **não** inventar rota/JWT. SC-005: aceite formal parcialmente bloqueado por P-010-2 / P-010-27. |

**Resultado do gate (pré-Fase 0)**: aprovado para pesquisa e desenho.

### Revisão pós-Fase 1

[data-model.md](./data-model.md) e [contracts/](./contracts/) descrevem `validarSeguranca()`, artefato versionado (leitura na inicialização; **DN-010-3 B** fail-fast em mismatch de versão), handoff mascarado e ordem versus **DN-008-4 B** / literais **DN-008-5 A** (referência à spec **008**), sem inventar regex nem path de catálogo (**DN-010-4 C** → **P-010-2**) nem HTTP do Módulo. **DN-010-1** define consumo downstream (Notificação), não implementação Finder. Gate mantido (lacunas explícitas conforme princípio V).

## API do Módulo Ocorrência (pendência DN-009-2)

Contrato HTTP **Finder → Módulo Ocorrência** (rota/path, claims JWT, códigos de resposta, timeout e semântica 5xx) **não** está definido no Drive nem no Jira além do padrão genérico (fonte: respostas-joao.md §Etapa 3 — *API Módulo: pendência no plan até o contrato*; origem: respostas-etapa3.md §10).

**O que as fontes fixam hoje**

| Aspecto | Valor documentado | Fonte |
|---------|-------------------|-------|
| Transporte | HTTP(S) interno, rede privada entre módulos | N1 A (fonte: respostas-joao.md, N1 A; origem: ADR_141) |
| Autenticação | JWT próprio por módulo (não Clerk/M2M) | ADR_028, ADR_054 (respostas-etapa3.md §10) |
| Semântica | Um envio por **ResumoAnalise** (lote); confirmação antes de concluir varredura; idempotência ADR_014 | ADR_062, seq. v6 §2; ADR_014 |
| Conteúdo mascarado | `achados[].linhas` e contadores **já mascarados** no Finder antes do POST | 010-6, 010-21; [saida-mascarada-handoff.md](./contracts/saida-mascarada-handoff.md) |
| Falha de envio (paralelo produto) | **009-18 A**: tentar de novo e alertar — princípio análogo para Finder→Módulo (reenvio mesma chave) | respostas-joao.md, 009-18 A |
| Rota conhecida (outro uso) | `POST /ocorrencias/{id}/sugestoes` (Sugestão manual) — **não** é handoff ResumoAnalise | ADR_083, ADR_092 (respostas-etapa3.md §10) |

**O que permanece em aberto (DN-009-2)**

- Path, método exato, headers e corpo de confirmação.
- Algoritmo/rotação JWT (**[PENDENTE]** ADR_028).
- Mapeamento 4xx/5xx/timeout e política de retry (opções 10a–10d, respostas-etapa3.md §10).

**Impacto no plano 010**: o **desenho** de `validarSeguranca()`, catálogo versionado e obrigações de `AchadoAnalise.linhas` **segue liberado**. Tarefas de **adaptador HTTP** (`internal/adapter/out/modocorrencia/`, testes contra URL real) permanecem **bloqueadas** até contrato — detalhamento em spec **009**; esta feature só define **o que** vai mascarado no payload (contrato `saida-mascarada-handoff.md`).

## Project Structure

### Documentation (this feature)

```text
specs/010-mascarar-dados-sensiveis/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── validar-seguranca.md
│   ├── catalogo-padroes-versionado.md
│   └── saida-mascarada-handoff.md
└── tasks.md             # /speckit-tasks — fora deste comando
```

### Source Code (repository root) — alvo de implementação

Estrutura **existente** (README § Layout):

```text
cmd/finder/
internal/domain/                 # item, envelope v1, chaves
internal/application/            # consumir.go — sem montarContexto/validarSeguranca hoje
internal/application/port/out/   # futuras: catalogo, (handoff — spec 009)
internal/adapter/in/queue/       # Redis Streams
internal/adapter/in/http/        # health
internal/adapter/out/openobserve/
internal/adapter/out/mongo/        # doc.go — stub
# internal/adapter/out/modocorrencia/  # handoff HTTP — bloqueado DN-009-2 (spec 009)
```

**Alvo coerente com esta feature** (nomes ADR; pacotes Go **NEEDS CLARIFICATION** até tasks):

```text
internal/domain/                    # LinhaLogMascarada, ResultadoValidarSeguranca, versão catálogo
internal/application/               # montarContexto(): validarSeguranca → MotorAnalise
internal/application/assinatura/    # NormalizarMensagem (dono spec 008, T016) — 010 só integra chamada (FR-015); sem normalizar_mensagem.go duplicado
internal/application/port/out/      # leitura artefato padrões (catálogo versionado)
# Artefato versionado (path/formato conforme ADR P-010-2):   # FR-008, DN-010-4 C
```

**Structure Decision**: inserir mascaramento **antes** de classificação/agregação; **não** duplicar validação no Módulo Ocorrência (010-21; revisao-T1-drive.md). **Não** duplicar o normalizador **DN-008-4 B**: dono da implementação é a spec **008** (`normalizar_mensagem.go`, **T016**); a **010** encadeia mascaramento → chamada a esse artefato (**FR-015**).

## Integração do normalizador canônico (etapa 3h — **010 H1**)

| Aspecto | Obrigação nesta feature | Fonte |
|---------|-------------------------|-------|
| Dono do algoritmo **DN-008-4 B** | Spec **008** — `internal/application/assinatura/normalizar_mensagem.go`, tarefa **T016** (`specs/008-agrupar-logs-ocorrencia/tasks.md`) | (fonte: respostas-joao.md §Etapa 3b, DN-008-4 B; origem: ADR_013) |
| Papel da **010** | Garantir **ordem** (`validarSeguranca()` antes da normalização) e **invocar** o normalizador canônico — **MUST NOT** reimplementar **DN-008-4 B** em pacote próprio desta feature | (fonte: respostas-joao.md, 010 H1; §Etapa 3g, João 16/11 13h22; §Etapa 3b, DN-008-4 B; spec.md FR-015) |
| Escopo do normalizador (produto) | O normalizador vale de forma única para as specs **006**, **008**, **009** e **010** (consumidores reutilizam o canônico da **008**) | (fonte: respostas-joao.md §Etapa 3b, DN-008-4 B; origem: DN-008-1 B) |
| **010 H1** (reuso canônico) | A **010 reaproveita** o normalizador da **008** — decisão formal do João 16/11 13h22 | (fonte: respostas-joao.md §Etapa 3g, João 16/11 13h22, 010 H1; origem: DN-008-4 B) |
| Ordem entre features | Integração mascaramento + normalização na **010** depende da **008** estar entregue (**T016** / literais **T048**); cadeia **009-22 A** (**009** após **008**) | (fonte: respostas-joao.md, 009-22 A; origem: Rodada pós-revisão T1, 16/11 00h01) |

## Pendências

**Fonte canônica** das pendências de produto (**P-010-2**, **P-010-27**, **P-emenda-I**, **P-emenda-III**): [spec.md § Pendências explícitas](./spec.md#pendências-explícitas) (corrige analyze **D2** — evitar três tabelas idênticas).

Pendências **adicionais** registradas neste plano (desenho/implementação Finder; sem criar ADR nesta pasta):

| ID | Descrição | Fonte | Efeito nesta feature |
|----|-----------|-------|----------------------|
| **DN-009-2** | Contrato HTTP Finder → Módulo (rota, JWT, respostas, 5xx/timeout) | respostas-joao.md §Etapa 3; respostas-etapa3.md §10 | **Bloqueia** apenas integração HTTP do handoff (spec 009); ver § API do Módulo |
| **N3 A** | Ampliar `AchadoAnalise` em **ADR** | respostas-joao.md, N3 A | Chaves JSON finais [PENDENTE]; não bloqueia núcleo de máscara |
| **006-10 A** | Alinhar Jira LLT-202 / epic | respostas-joao.md, 006-10 A | Fora desta pasta |
| **ADR_028** | JWT interno [PENDENTE] | respostas-etapa3.md §10 | Paralelo a **DN-009-2** |
| **P-010-2** (metadados artefato) | Valores concretos de **caminho** e **formato** do catálogo no repositório | **DN-010-4 C**; 010-2 A; respostas-joao.md §Etapa 3c, João 16/11 11h | Aguardar ADR da lista — **não** inventar path/formato nesta pasta |
| **Biblioteca regex Go** | `regexp` stdlib vs extensão | analyze **U2**; ADR_064 | **NEEDS CLARIFICATION** até ADR |
| **Placeholder e-mail exato** | Texto fixo canônico | 010-11 A; P-010-2 | **NEEDS CLARIFICATION** até ADR |

## Tarefas bloqueadas vs liberadas (resumo)

| Área | Liberado para plano/tarefas | Bloqueado até pendência |
|------|-----------------------------|-------------------------|
| `validarSeguranca()` / contratos internos | Sim | — |
| Catálogo versionado (interface, não regex) | Sim (conteúdo + política **DN-010-3 B**) | Path/formato concretos (**P-010-2** / ADR) |
| `AchadoAnalise.linhas` mascaradas | Sim | — |
| Aceite formal SC-001/002 | Desenho | **P-010-27**, **P-010-2** |
| Entrega HTTP `ResumoAnalise` | Modelo de payload mascarado | **DN-009-2**, **ADR_028** (spec 009) |
| Assinatura / `mensagem` normalizada para agrupamento | Ordem mascaramento → normalização (**DN-008-4 B**, **FR-015**); **010** integra chamada a `normalizar_mensagem.go` (**008**, **T016**) — sem cópia do algoritmo (**010 H1**) | Entrega **T016**/**T048** na spec **008** antes de fechar integração na **010** |
| Inicialização / readiness catálogo | **DN-010-3 B**: recusar subir se versão ≠ combinada | Versão combinada em config de deploy (coordenação Sugestão — 010-19 A) |
| Métricas **≥95% precisão por classe** (ERROR, WARNING, INFO) | Documentado como **SC-007** referência cruzada (**DN-007-2 B**); mascaramento não tem gate métrico próprio | Avaliação formal na spec **007** / base `datasets/labeled-logs/` (**007 formato A**, **007 caminho C** — só documentar; pasta **não** criada nesta execução) |

## Referência cruzada — base rotulada (spec 007)

Decisões **transversais** documentadas aqui para coerência de pipeline; **dona** da base e gates de precisão é a spec **007**, não LLT-202.

| Aspecto | Valor | Fonte |
|---------|-------|-------|
| Formato | **JSONL** (**007 formato A**) | respostas-joao.md §Etapa 3b, João 16/11 09h02 |
| Caminho no repositório | `datasets/labeled-logs/` na **raiz** do repositório (**007 caminho C**) | respostas-joao.md §Etapa 3b, João 16/11 09h02 |
| Criação de pasta/dados | **Proibida** nesta execução — apenas documentação | instrução etapa 3c |
| Relação com 010 | Base **007-16 A** (sem PII); validação de **classificação** em texto que, no Finder real, MUST vir **mascarado** quando derivado de linhas processadas (**010-5 A**, **FR-015**) | respostas-joao.md 007-16 A; spec.md DN-007-2 B |
| Meta precisão | **≥95% por classe** ERROR/WARNING/INFO (**DN-007-2 B**; **007-meta**) | respostas-joao.md §Etapa 3b; respostas-etapa3.md §1 |

(fonte: respostas-joao.md §Etapa 3b, DN-007-2 B / 007 formato A / 007 caminho C; origem: LLT-199, spec 007)

## Decisões etapa 3e (resumo)

| ID | Decisão | Fonte |
|----|---------|-------|
| **DN-010-3 B** | Versão lida ≠ versão combinada → **recusar subir** (inicialização) | respostas-joao.md §Etapa 3c, João 16/11 11h; respostas-etapa3c.md §7; ADR_079, 010-19 A |
| **DN-010-4 C** | **Caminho** e **formato** do artefato só na **ADR** da lista (010-2 A); pendência **P-010-2** | respostas-joao.md §Etapa 3c, João 16/11 11h; respostas-etapa3c.md §8; ADR_079, 006-3 C |
| **DN-008-5 A** | Marcadores normalização: `<NUM>` `<UUID>` `<IP>` `<DATE>` `<HEX>` (dono spec **008**) | respostas-joao.md §Etapa 3c, João 16/11 11h; respostas-etapa3c.md §4; DN-008-4 B, ADR_013 |

## Decisões etapa 3h (resumo)

| ID | Decisão | Fonte |
|----|---------|-------|
| **010 H1** | A **010 reaproveita** o normalizador canônico da **008** — **MUST NOT** criar segundo `normalizar_mensagem` nesta feature | (fonte: respostas-joao.md, 010 H1; origem: §Etapa 3g, João 16/11 13h22, DN-008-4 B, FR-015) |

**Sem alteração nesta pasta (etapa 3h):** **DN-008-9 A** (composição/escape da assinatura — dono spec **008**); **DLQ B** / **DN-008-3** / **DN-008-6** (não citados na spec **010**). (fonte: respostas-joao.md §Etapa 3g; instrução etapa 3h)

## Divergências documentadas

Canônico completo em [spec.md § Divergências documentadas](./spec.md#divergências-documentadas).

Decisões **DN-006-3 B**, **DN-007-4 A**, **DN-008-7 B** (specs 006/007/008) **não** são reescritas aqui — permanecem nas respectivas pastas; fonte canônica: `respostas-joao.md` §Etapa 3c (João, 16/11 11h02); impacto indireto em 010 só via pipeline (**FR-015**, base 007).

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Princípio I (Finder “em ocorrência” direta) | T1 revisto + mascaramento na saída `ResumoAnalise` | Finder gravar `conteudo` localmente (T1 A antiga) — **substituído** (revisao-T1-drive.md; 010-25) |