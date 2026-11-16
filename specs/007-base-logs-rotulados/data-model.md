# Data Model: Base de logs rotulados

**Feature**: `007-base-logs-rotulados` | **Date**: 2026-11-16

Entidades lógicas do conjunto versionado e documentação associada. Persistência = arquivos no Git + manifests (007-4 A).

**Serialização física** (007-formato A, 007-caminho C): cada entidade `RegistroLogRotulado` ou `RegistroApontamentoRotulado` MUST ocupar **uma linha** em arquivo **`.jsonl`** sob `datasets/labeled-logs/` na raiz do repositório (fonte: respostas-joao.md §Etapa 3b; spec.md § Artefatos físicos). Manifests (`ORIGEM.md`, split, contagens) são arquivos de documentação no mesmo diretório, não JSONL.

---

## ConjuntoRotulado (base)

| Campo | Obrigatório | Descrição | Validação / fonte |
|-------|-------------|-----------|-------------------|
| `versaoConjunto` | recomendado | Identificador semver ou data da entrega | Documentação; 006-3 C quando cruzar com versão de regras |
| `declaracaoOrigem` | sim | Texto em arquivo dedicado (ex. `ORIGEM.md`) | 007-15 A; SC-001 |
| `metodoRotulacao` | sim | Autoria + procedimento | 007-10 B, 007-11 D |
| `contagensPorClasse` | sim (na entrega) | ERROR, WARNING, INFO | Mesmo mínimo cada; milhares no total — 007-1/2/25 |
| `split` | sim (documentado) | Partições nomeadas | 007-20 B |
| `mapeamentoFisico` | sim | Nomes de colunas/arquivos | Tabela em [spec.md](./spec.md) § Mapeamento — 007-6 D |

**Estado**: evolui com commits Git; duplicatas lógicas com rótulos diferentes permitidas se documentadas (007-12 B).

---

## ClassLabel

Enum gold:

| Valor | Constituição | Notas |
|-------|--------------|-------|
| `ERROR` | erro | 007-9 A |
| `WARNING` | alerta | |
| `INFO` | informação | T7 A — rótulo válido na base mesmo divergindo SSID_007.1 legado |

**Validação**: qualquer outro valor → registro **inválido** (spec Edge Cases; constituição III).

---

## RegistroLogRotulado

Unidade **log único** (007-7 C — tipo `log`).

| Campo lógico | Obrigatório | Descrição | Validação |
|--------------|-------------|-----------|-----------|
| `itemType` | sim | Constante `log` | Discriminador de granularidade |
| `level_text` | sim | Texto do campo de nível para o classificador | Rejeitar se vazio — 007-13 A |
| `message_text` | sim | Corpo / mensagem | Rejeitar se vazio — 007-13 A |
| `lhamalog_batch_id` | sim | Id de lote OpenObserve | Contrato §7 |
| `lhamalog_record_id` | sim | Id de registro | Contrato §7 |
| `tenantId` | sim | Valor de **teste** documentado | 007-8 A; contrato §4 |
| `productId` | sim | Valor de **teste** documentado | 007-8 A |
| `class_label` | sim | `ClassLabel` gold | FR-002 |
| `split` | opcional | Nome da partição | 007-20 B |
| `notas` | opcional | Duplicata, conflito de rótulo, etc. | 007-12 B |

**Relacionamentos**: alimenta operação `Classificar log` via [entrada-classificador-desde-base.md](./contracts/entrada-classificador-desde-base.md).

---

## RegistroApontamentoRotulado

Unidade **apontamento** — vários logs sob um mesmo contexto de avaliação (007-7 C). Estrutura lógica mínima **sem** inventar formato de agrupamento em produção:

| Campo lógico | Obrigatório | Descrição |
|--------------|-------------|-----------|
| `itemType` | sim | Constante `apontamento` |
| `members` | sim | Array não vazio de objetos com os **mesmos campos** de `RegistroLogRotulado` **exceto** `class_label` — MUST NOT incluir `class_label` no membro (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); origem: DN-007-4 A; respostas-etapa3c.md §3) |
| `class_label` | sim | Gold **único** no registro de apontamento (agregado); ERROR, WARNING ou INFO — DN-007-4 A |

Validação estrutural: presença de `class_label` em qualquer `member` → registro **inválido** (DN-007-4 A; FR-007b).

| `tenantId` / `productId` | sim | No apontamento ou herdados uniformes nos membros |
| `lhamalog_batch_id` | sim | Referência principal do cenário |
| `documentacao` | sim | Texto que explica o papel do item vs log único | 007-7 C |

**Transições**: usado em testes que simulam vários logs antes do agrupamento 008; **não** implica leitura Redis.

**Avaliação (007 H2 D)**: predição por `member`; classe prevista do apontamento = **mais grave** entre predições (ERROR > WARNING > INFO, ADR_170); comparar ao `class_label` gold agregado; unidade de métrica = apontamento (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), 007 H2 D; origem: ADR_170; respostas-etapa3f.md §D; DN-007-4 A). Ex.: gold ERROR + predições ERROR, ERROR, WARNING → acerto; gold WARNING + predições WARNING, WARNING, ERROR → erro.

---

## RegistroSoCatalogo (placeholder)

Subconjunto reservado a cenários em que **só o catálogo** decide a classe (007-18 B, spec 006 US3).

| Campo | Estado |
|-------|--------|
| Conteúdo de exemplos | **BLOQUEADO** até catálogo existir — DN-007-1 C |
| `criterioInclusao` | Documentado em spec/README do subconjunto; alinhado 006 US3 |
| `particao` | Nome reservado no split (ex. `so_catalogo`) — nome exato em tarefas |

**Estado**: manifest pode existir vazio ou só com critério textual; **MUST NOT** inventar registros gold (FR-011).

---

## DeclaracaoOrigem

Artefato textual (não necessariamente um registro por linha da base).

| Seção obrigatória | Conteúdo |
|-------------------|----------|
| Procedimento | Anonimização de logs reais |
| Exclusões | Sem dado real de cliente identificável |
| PII | Ausência total, incluindo fictícia (007-16 A) |
| Responsável | Engenharia Finder (007-10 B) |
| Data / versão | Rastreio de entrega |

---

## ParticaoSplit

| Nome lógico | Uso típico (não normativo até documentar) |
|-------------|-------------------------------------------|
| `treinamento` | Pós-MVP — 007-3 A |
| `validacao` | Métricas LLT-199 / 007-meta |
| `teste` | Held-out adicional |
| `so_catalogo` | Reservada — DN-007-1 C |

Critérios de atribuição: MUST documentar **nomes** de partições e **critérios qualitativos** em `SPLIT.md` (007-20 B; spec FR-005; analyze M2). **Sem fonte** para proporção numérica, seed ou regra de estratificação — MUST NOT inventar na spec/plan; preencher só na entrega se o João indicar, senão descrever qualitativamente (ex.: estratificação por classe sem percentuais fixos).

---

## Relação com entidades 006

| Entidade 006 | Relação com 007 |
|--------------|-----------------|
| Entrada `Classificar log` | `level_text` + `message_text` (+ `catalogo` opcional em cenários só-catálogo futuros) |
| Métricas offline | Log único: `class_label` vs `classe` predita. Apontamento: prevista = mais grave entre predições dos `members`; comparar ao gold agregado (007 H2 D; FR-007c) |
| Dados sintéticos 006 plan | **Substituídos** por este conjunto (007-19 B) |

---

## Validação estrutural (resumo)

| Regra | Fonte |
|-------|-------|
| 100% registros com campos FR-004 | SC-004 / 007-21 A |
| 0% PII | SC-001 implícito / 007-16 A |
| Contagens publicadas | SC-003 / 007-25 A |
| Precisão, recall e F1 ≥ 95% por classe | SC-005 / 007-meta / DN-007-2 B / DN-006-2/007-3 B |
| Apontamento: um gold, members sem `class_label` | FR-007b / DN-007-4 A |
| Apontamento: prevista = mais grave entre predições; métrica por apontamento | FR-007c / 007 H2 D |

Detalhe operacional: [contracts/registro-rotulado.md](./contracts/registro-rotulado.md).