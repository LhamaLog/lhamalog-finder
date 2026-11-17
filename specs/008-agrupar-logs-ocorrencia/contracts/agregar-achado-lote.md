# Contract: Agregar achados no lote

Operação **interna**: particionar logs classificados do lote corrente.

## Entrada

| Campo | Descrição |
|-------|-----------|
| `logs[]` | Logs classificados com `assinatura` calculada |
| `batchId` | Lote ingested v1 |
| `productId` / `tenantId` | Isolamento |

## Saída

| Campo | Descrição |
|-------|-----------|
| `achados[]` | Um elemento por partição `(assinatura, classe)` |

Cada achado contém:

| Campo | Regra |
|-------|-------|
| `contagemLote` | Número de logs na partição (**008-12 A** — cada log +1) |
| `primeiraNoLote` / `ultimaNoLote` | Min/max **timestamp do log** (**008-14 A**) |
| Repetições no mesmo lote | **Contam** no contador (**008-13 A**) |

## Regras

1. **Achado** = logs com a **mesma assinatura** no **produto** (**008-18 A**).
2. Mesma assinatura, **classes** diferentes → achados **distintos** (**008-15 A**).
3. Logs **INFO** entram como demais classes (**T7 A**).
4. Lote **incompleto** → operação **não** é invocada (**T4 A**).
5. Agrupamento no lote MUST usar apenas assinatura **DN-008-1 B** + **DN-008-4 B** (determinístico, local); **MUST NOT** depender de consulta OpenObserve nesta etapa (**FR-008**, achado analyze **A1**).

## Transições

```text
Classificar (006) → Calcular assinatura (DN-008-5 A) → Particionar → (por achado, DN-008-7 B) consulta entre lotes → (um ResumoAnalise/lote, DN-008-7)
```

## Fontes

- (fonte: respostas-joao.md, 008-18 A / 008-12 A / 008-13 A / 008-14 A / 008-15 A / T7 A / T4 A; origem: Jira LLT-200)
- (fonte: spec.md, FR-001, User Story 1)