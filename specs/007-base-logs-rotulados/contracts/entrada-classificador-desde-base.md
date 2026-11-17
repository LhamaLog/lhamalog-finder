# Contract: Entrada do classificador a partir da base rotulada

Contrato **interno** de teste/avaliação: mapeia um registro da base 007 para a operação **`Classificar log`** da spec 006. **Não** passa pela fila Redis (007-23 A). **Não** consulta OpenObserve para obter texto — `level_text` e `message_text` já vêm materializados no artefato (007-5 B, 007-6 D).

## Operação

`Avaliar registro rotulado` (wrapper de teste sobre `Classificar log`)

### Entrada (por registro de log)

| Origem (base 007) | Destino (classificar-log) | Notas |
|-------------------|---------------------------|-------|
| `level_text` | `nivelCampo` | Opcional no contrato 006, **obrigatório** na prática quando o gold depende do nível |
| `message_text` | `mensagem` | Esperado **sem PII** (007-16 A); mascaramento 010 é ortogonal |
| — | `catalogo` | Omitir ou injetar fixture **só** em cenários só-catálogo futuros (DN-007-1 C) |

Campos `lhamalog_batch_id`, `lhamalog_record_id`, `tenantId`, `productId` **não** entram na porta `Classificar log` mínima; podem ser usados por testes de rastreio, métricas ou cenários de apontamento.

### Saída

| Campo | Uso na avaliação |
|-------|-------------------|
| `classe` | Predição (ERROR, WARNING, INFO) |
| `catalogoUsado` | Metadado para análise de cobertura de catálogo |

### Comparação gold vs predição

| `itemType` | Campo gold | Campo predito | Regra |
|------------|------------|---------------|-------|
| `log` | `class_label` | `classe` | Igualdade exata de string |
| `apontamento` | `class_label` (agregado) | Classe prevista do apontamento | **007 H2 D** — ver § Registro de apontamento abaixo |

Métricas agregadas (precisão, recall, F1 **por classe**): **006-2 B**; alvos na partição `validacao`: **≥ 95%** de precisão, recall e F1 **em cada classe** — **007-meta**, **DN-007-2 B**, **DN-006-2/007-3 B** (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h).

### Registro de apontamento

1. Gold **somente** no objeto apontamento (`class_label` agregado); membros **sem** `class_label` (DN-007-4 A; fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); origem: DN-007-4 A; respostas-etapa3c.md §3).
2. Para cada `member`, invocar `Classificar log` e obter `classe` predita por log (ADR_167; pista LLT-199 crit. 1 «todo log» — contexto, sem gold por membro).
3. **Classe prevista do apontamento** = a **mais grave** entre as predições dos `members` (**ERROR** > **WARNING** > **INFO**, ADR_170). Comparar essa prevista ao único `class_label` gold do apontamento (igualdade exata). **Unidade** da avaliação/métrica = **apontamento** (007 H2 D; FR-007c; fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), 007 H2 D; origem: ADR_170; respostas-etapa3f.md §D). **Exemplo:** gold ERROR, predições ERROR, ERROR, WARNING → prevista ERROR → acerto. **Contraexemplo:** gold WARNING, predições WARNING, WARNING, ERROR → prevista ERROR → erro (falso positivo para ERROR).
4. Este contrato **não** define agrupamento em achado (spec 008).

### Cenários de validação obrigatória (espelho spec)

| ID | Descrição | Fonte |
|----|-----------|-------|
| E1 | 100% dos registros expõem campos obrigatórios FR-004 | SC-004 / 007-21 A |
| E2 | Predição somente ERROR/WARNING/INFO | 006 FR-003 |
| E3 | Partição `validacao` usada para medição 007-meta | 007-20 B (nomes documentados) |
| E4 | Precisão, recall e F1 **≥ 95% por classe** na partição `validacao` | 006-2 B; DN-006-2/007-3 B; SC-005 |
| E5 | Apontamento: prevista = mais grave entre predições dos `members`; métrica por apontamento | 007 H2 D; FR-007c; DN-007-4 A |

## O que não cobre

- Handoff HTTP ao Módulo Ocorrência (DN-009-2; spec 006).
- Normalização de mensagem para assinatura (**DN-008-4**; spec 008).
- Leitura `log.ingested` da fila.

## Referência cruzada

- Contrato destino: [spec 006 `classificar-log.md`](../../006-modelo-classificacao-logs/contracts/classificar-log.md)
- Esquema de dados: [data-model.md](../data-model.md)
- Registro em arquivo: [registro-rotulado.md](./registro-rotulado.md)

## Fontes

- (fonte: respostas-joao.md, 007-5 B / 007-6 D / 007-22 B / 007-23 A / 007-meta / 006-2 B / §Etapa 3e 007 H2 D)
- (fonte: spec.md, User Story 2, FR-009)