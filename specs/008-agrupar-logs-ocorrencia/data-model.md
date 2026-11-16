# Data Model: Agrupar logs semelhantes na mesma ocorrência

**Feature**: `008-agrupar-logs-ocorrencia` | **Date**: 2026-11-16 (etapas 3b–3h)

Entidades do fluxo de agrupamento no **Finder** e chaves persistidas no **Módulo Ocorrência** (fora deste repo). Classificação prévia: spec **006**.

---

## ChaveUnicidadeOcorrencia

Identificador lógico de **uma** ocorrência agrupada no Módulo.

| Componente | Descrição | Fonte |
|------------|-----------|-------|
| `assinatura` | Valor estável derivado dos quatro componentes OTel | DN-008-1 B, 008-1 fech. 3b |
| `classe` | `ERROR` \| `WARNING` \| `INFO` | 008-15 A, FR-006 |
| `productId` | Produto do token/contrato | 008-11 rev., contrato §4 |
| `tenantId` | Metadata confiável | 008-11 rev., FR-011 |

**Invariante**: no máximo **uma** ocorrência por tupla (FR-005).

**Nota**: `nivel` na assinatura pode coincidir semanticamente com `classe`; a dimensão `classe` na chave permanece explícita (**008-15 A**, **008-11 rev.**).

---

## AssinaturaLog

| Campo / componente | Obrigatório | Descrição |
|--------------------|-------------|-----------|
| `nivel` | sim | Nível OTel do log (componente da assinatura) | DN-008-1 B |
| `mensagemNormalizada` | sim | Mensagem após **DN-008-4 B** com marcadores **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`** (**DN-008-5 A**) | DN-008-1 B, DN-008-4 B, DN-008-5 A |
| `serviceName` | sim | Atributo `service.name` OTel | DN-008-1 B |
| `host` | sim | Host/máquina (ex.: `host.name` ou equivalente ADR_013) | DN-008-1 B |
| `valor` | após cálculo | Texto: **nível**|**mensagemNormalizada**|**serviceName**|**host** com **`|`** (**DN-008-8 A**, **008-assinatura A**); **`|`** em componentes → **`\|`** (**DN-008-9 A**); escape de **`\`** literal: pendente (sem fonte) | spec FR-002, FR-002c, research R17 |
| `versaoRegra` | sim quando existir | Versão da regra de assinatura / 006-3 C | N3 A [PENDENTE ADR] |

**Fora da assinatura:** `traceId`, `correlationId` permanecem no log e podem ir à ocorrência (**ADR_040**) mas **MUST NOT** entrar no cálculo (**008-2 rev.**).

**Validação**: logs sem `productId` coerente **não** entram no agrupamento — falha crítica (**008-16 A**).

---

## LogAgrupavel

Extensão lógica de **LogClassificado** (spec 006).

| Campo | Descrição |
|-------|-----------|
| `classe` | ERROR \| WARNING \| INFO |
| `assinatura` | Calculada após **DN-008-4 B** / **DN-008-5 A**; `valor` conforme **008-assinatura A** |
| `timestampLog` | Instantâneo do registro OTel para **008-14 A** |
| `apontamento` | Referência no lote |
| `tenantId` / `productId` | Isolamento |
| `nivel`, `serviceName`, `host`, `mensagem` | Entrada do cálculo de assinatura |

---

## Achado (partição no lote)

| Campo | Descrição | Fonte |
|-------|-----------|-------|
| `chave` | `(assinatura, classe)` dentro do produto do lote | 008-18 A, 008-15 A |
| `logs` | Membros do lote | 008-12 |
| `contagemLote` | Contagem **só do lote** antes de fusão | 008-12, 008-13 |
| `primeiraNoLote` / `ultimaNoLote` | Extremos timestamp no lote | 008-14 A |

**Transição**: logs classificados → particionar por `(assinatura, classe)` → um achado por partição.

---

## EstadoAgregacaoEntreLotes

Resultado da consulta ao Módulo Ocorrência (**N2 B**) para uma `ChaveUnicidadeOcorrencia` — **uma consulta por achado** (**DN-008-7 B**, etapa 3e). Transporte HTTP **BLOQUEADO** até **DN-009-2**.

| Campo | Descrição | Fonte |
|-------|-----------|-------|
| `existe` | Se há ocorrência para a chave | LLT-200 crit. 2 |
| `contagemAcumulada` | Contador conhecido (se API retornar) | N2 B, 008-12 |
| `primeiraGlobal` / `ultimaGlobal` | Extremos conhecidos | 008-14 A |
| `status` | ABERTA \| EM_ANALISE \| FECHADA | ADR_082, 008-8 A |

**Fusão no Finder (antes do handoff)**: `contagem` = contagemAcumulada (se existir) + contagemLote; `primeiraOcorrencia` = min(primeiraGlobal, primeiraNoLote); `ultimaOcorrencia` = max(ultimaGlobal, ultimaNoLote). Contrato: [consulta-ocorrencias-modulo.md](./contracts/consulta-ocorrencias-modulo.md).

---

## MarcaRecorrente

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `recorrente` | boolean | Resultado da consulta OpenObserve (**ADR_013**) |
| `evidencia` | opcional | Detalhes da janela 500 lotes — **NEEDS CLARIFICATION** (query) |

**Regra**: repetição no mesmo lote **conta** no contador (**008-13**) mas **não** conta para recorrência ADR_013 no mesmo arquivo/lote (**008-5 A**).

**Falha de consulta**: até **3 tentativas** (**DN-008-3** rev. 3h); após esgotar, aviso ao **Core** (API interna, **batchId** + **motivo** — **DN-008-6** rev. 3h; **sem DLQ**); **MUST NOT** `recorrente=false` por indisponibilidade.

---

## AchadoAnalise (agrupado)

Payload no `ResumoAnalise`. Base Drive + extensões João (ver spec 006 data-model).

### Campos adicionais desta feature (N3 A — [PENDENTE ADR])

| Campo | Descrição |
|-------|-----------|
| `assinatura` | Valor estável |
| `contagem` | Após fusão entre lotes |
| `primeiraOcorrencia` | Timestamp log |
| `ultimaOcorrencia` | Timestamp log |
| `versaoRegras` | 006-3 C |

Demais campos Drive (`tipo`, `linhas`, `recorrente`, …) permanecem conforme ADR_040 / Finder.mmd (revisao-T1-drive.md §2).

---

## ResumoAnalise

**Um envelope por lote** processado com achados agrupados (**DN-008-7**, ADR_062). Idempotência **batchId + execId** (**N4 A**; **DN-006-1** na spec 006). Detalhes: [resumo-analise-agrupado.md](./contracts/resumo-analise-agrupado.md).

---

## Ocorrencia (persistida — Módulo Ocorrência)

Modelo **lógico** (persistência fora deste repositório).

| Atributo | Regra |
|----------|-------|
| Chave | `ChaveUnicidadeOcorrencia` |
| Contador | Atualizado a partir do achado (**008-12**) |
| Primeira/última vez | Do achado (**008-14**) |
| `recorrente` | Atributo, não status (ADR_063) |
| Status | ABERTA, EM_ANALISE, FECHADA |
| Reabertura | FECHADA → ABERTA no upsert equivalente (**008-4**) |

**Finder MUST NOT** criar nem atualizar documento de ocorrência (FR-001).

---

## Relacionamentos

```text
N LogAgrupavel --(mesma assinatura+classe)--> 1 Achado (lote)
1 Achado + ConsultaModulo + ConsultaOO --> 1 AchadoAnalise
N AchadoAnalise --> 1 ResumoAnalise
1 AchadoAnalise --handoff--> 0..1 Ocorrencia (upsert no Módulo)
```

---

## Validações de negócio

| Regra | Fonte |
|-------|-------|
| Lote incompleto: sem agrupamento | T4 A |
| Duas ocorrências mesma chave: proibido | FR-005, 008-11 |
| Classe fora das três: proibido | FR-006 |
| Algoritmo DN-008-4 B e marcadores **DN-008-5 A** fechados | DN-008-4 B, DN-008-5 A, FR-002b |
| Valor assinatura: forma texto (**008-assinatura A**); separador **`|`** (**DN-008-8 A**); escape **`|`** → **`\|`** (**DN-008-9 A**) | FR-002c, research R17 |
| Consulta Módulo **por achado**; envio **um ResumoAnalise/lote** | DN-008-7 B, FR-009 |
| Composição assinatura fechada; aceite E2E parcial até DN-009-2 / ADR N3 / ADR_051 / escape `\` | DN-008-1 B, SC-005 |