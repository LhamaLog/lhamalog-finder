# Contract: Registro rotulado (artefato versionado)

Contrato **externo** entre o conjunto de dados no Git e consumidores (testes, métricas offline, spec 006). **Não** é endpoint HTTP. **Não** é mensagem Redis.

## Escopo

- Define campos lógicos e regras de validação para **registro de log** e **registro de apontamento**.
- Nomes físicos de arquivo e colunas MUST constar na tabela de mapeamento em [spec.md](../spec.md) (007-6 D) quando o conjunto for entregue.
- Formato de serialização: **JSONL** (um objeto JSON por linha) em **`datasets/labeled-logs/`** na raiz do repositório (007-formato A, 007-caminho C; fonte: respostas-joao.md §Etapa 3b).

## Registro de log (`itemType`: `log`)

### Campos obrigatórios

| Campo | Tipo lógico | Restrições |
|-------|-------------|------------|
| `itemType` | string | Valor exato `log` |
| `level_text` | string | Não vazio após trim |
| `message_text` | string | Não vazio após trim |
| `lhamalog_batch_id` | string | Não vazio |
| `lhamalog_record_id` | string | Não vazio |
| `tenantId` | string | Valor de teste documentado no conjunto |
| `productId` | string | Valor de teste documentado |
| `class_label` | string | Um de `ERROR`, `WARNING`, `INFO` |

### Campos opcionais

| Campo | Uso |
|-------|-----|
| `split` | Nome da partição (007-20 B) |
| `notas` | Duplicata, conflito intencional, referência de revisão (007-12 B) |

### Rejeição

| Condição | Resultado |
|----------|-----------|
| `class_label` ∉ {ERROR, WARNING, INFO} | Registro inválido |
| `level_text` ou `message_text` ausente/vazio | Rejeitar (007-13 A) |
| PII detectável (e-mail, senha, token, etc.) | MUST NOT publicar (007-16 A; ADR_064) |
| Spam / "não classificado" como rótulo | Registro inválido (007-9 A) |

## Registro de apontamento (`itemType`: `apontamento`)

### Campos obrigatórios

| Campo | Tipo lógico | Restrições |
|-------|-------------|------------|
| `itemType` | string | Valor exato `apontamento` |
| `members` | array | ≥ 1 objeto com campos de registro de log **sem** `class_label` (DN-007-4 A) |
| `class_label` | string | ERROR, WARNING ou INFO — **único** no apontamento (agregado); MUST NOT aparecer em `members` (DN-007-4 A; FR-007b) |
| `lhamalog_batch_id` | string | Cenário principal |
| `tenantId` / `productId` | string | Valores de teste |
| `documentacao` | string | Papel do item vs log único (007-7 C) |

### Regras

- MUST NOT inventar estrutura de hash/assinatura (007-17 A; agrupamento 008 é outra feature).
- Apontamentos seguem as mesmas regras de privacidade que registros de log.
- `class_label` em qualquer elemento de `members` → registro **inválido** (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); origem: DN-007-4 A; respostas-etapa3c.md §3).
- **Avaliação offline** (não altera o esquema do arquivo): predição por `member`; **classe prevista** do apontamento = **mais grave** entre predições (**ERROR** > **WARNING** > **INFO**, ADR_170); comparar ao `class_label` agregado; unidade de métrica = apontamento — **007 H2 D** / FR-007c ([entrada-classificador-desde-base.md](./entrada-classificador-desde-base.md) § Registro de apontamento; fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40); respostas-etapa3f.md §D).

## Conjunto e manifestos

Arquivos de nível de conjunto (exemplos de propósito, nomes físicos em tarefas):

| Artefato | Conteúdo mínimo |
|----------|-----------------|
| Declaração de origem | SC-001 / 007-15 A |
| Manifest de split (`SPLIT.md`) | Nomes de partições e **critérios qualitativos** (007-20 B; FR-005) — **sem** proporção/seed numéricos sem fonte (analyze M2) |
| Manifest de contagens | Totais por classe para SC-003 / checklist |

## Subconjunto só-catálogo

- Partição ou flag reservada conforme [data-model.md](../data-model.md) § `RegistroSoCatalogo`.
- **Conteúdo de registros**: **BLOQUEADO** até catálogo (DN-007-1 C; 006-7 D).

## Duplicatas

- Duplicatas lógicas com **rótulos diferentes** são permitidas se `notas` ou documentação do conjunto explicar o motivo (007-12 B).

## Fontes

- (fonte: spec.md, FR-002–FR-013, SC-001–SC-004, Edge Cases)
- (fonte: respostas-joao.md, 007-5 B … 007-13 A / 007-7 C / 007-12 B)
- (fonte: contrato-log-ingested-v1.txt, §4, §7)