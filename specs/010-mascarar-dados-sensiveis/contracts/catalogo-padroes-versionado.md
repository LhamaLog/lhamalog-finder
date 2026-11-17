# Contract: Catálogo de padrões versionado (Finder + Sugestão)

Artefato compartilhado exigido por **FR-008**. Dono do **conteúdo** das regex: **ADR** (010-2 A) — este contrato define obrigações de leitura/versão, não a lista de padrões.

## Responsabilidades

| Módulo | Papel | Fonte |
|--------|-------|-------|
| **Finder** (este repo) | Lê o artefato; aplica em `validarSeguranca()`; registra `versao` no contexto/resumo | 010-18 A; 010-20 |
| **Sugestão** (Java, fora do repo) | Lê o **mesmo** artefato; reaplica regras no prompt LLM | 010-20; ADR_079; LLT-212 |

## Versão

| Regra | Descrição | Fonte |
|-------|-----------|-------|
| Paridade | Finder e Sugestão em um ambiente MUST us a **mesma** `versao` | 010-19 A |
| Registro | `versao` MUST aparecer no processamento do lote (ex. `ContextoLoteMascarado.versaoCatalogo`, `versaoRegras` no handoff quando aplicável) | FR-008; 010-19 |
| Divergência vs. versão combinada | Na **inicialização**, se `versao` lida ≠ versão **combinada** (config/deploy para paridade com Sugestão), o serviço MUST **recusar subir** — falha observável (log/métrica readiness) | **DN-010-3 B**; 010-19 A; FR-008; SC-006 |

## Localização

| Item | Status | Fonte |
|------|--------|-------|
| Mecanismo | Arquivo **versionado no repositório** (controle Git) | 010-18 A; ADR_079 |
| Caminho no monorepo | Declarado na **ADR da lista** (**010-2 A**) — **PENDENTE** até ADR publicar (**P-010-2**) | **DN-010-4 C** |
| Formato de serialização | Idem — conforme ADR da lista | **DN-010-4 C**; respostas-etapa3c.md §8 |

**Proibido nesta fase**: inventar path, formato ou exemplos de regex (constituição V; **DN-010-4 C**, P-010-2).

## Conteúdo mínimo lógico (quando ADR publicar)

| Seção | Propósito | Fonte |
|-------|-----------|-------|
| `versao` | Identificador publicado | 010-19 A |
| `padroes[]` | Definições regex/regra — **somente** conforme ADR oficial | 010-1; P-010-2 |
| `marcadorLinhaSegredo` | MUST alinhar a `"log com senha omitido"` | 010-10 A |
| `placeholderEmail` | Texto de substituição de e-mail | 010-11 A |

Até a ADR fechar, implementação MAY usar comportamento **documentado** na ADR_064 (interim) sem afirmar lista completa de regex.

## Validação na leitura (Finder)

| Falha | Comportamento sugerido pelo negócio | Fonte |
|-------|--------------------------------------|-------|
| Arquivo ausente ou ilegível | Tratar como falha de `validarSeguranca()` → 010-16 C | 010-16 C |
| `versao` lida ≠ versão **combinada** na inicialização | **Recusar subir** (processo não healthy); MUST NOT processar lotes até corrigir config/artefato | **DN-010-3 B**; respostas-joao.md §Etapa 3c, João 16/11 11h |
| `versao` diverge em runtime entre Finder e Sugestão (deploy descoordenado) | Prevenção operacional: **010-19 A** + política de deploy; evidência Finder em SC-006 | 010-19 A; SC-006 |

## Relação com outras decisões

- **Não** duplicar catálogo na spec 006 FR-009 (spec 010 / ADR).
- **Não** criar biblioteca comum Go/Java (ADR_145) — **arquivo** compartilhado (010-18 A supera cópia manual como mecanismo escolhido).

## Fontes consolidadas

- (fonte: respostas-joao.md, 010-2 A / 010-18 A / 010-19 A; §Etapa 3c, João 16/11 11h, DN-010-3 B / DN-010-4 C; respostas-etapa3c.md §§7–8)
- (fonte: respostas-documentadas.md, 010-1 / 010-18)
- (origem: ADR_064 [PENDENTE], ADR_079, ADR_145)