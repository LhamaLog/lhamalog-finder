# Contract: Handoff ResumoAnalise (Finder → Módulo Ocorrência)

Interface **externa** entre este repositório e o **Módulo de Ocorrência**. Formato alinhado ao Drive; extensões do João marcadas até ADR (N3 A).

**Estado do contrato HTTP**: **BLOQUEADO (DN-009-2)** — corpo lógico abaixo é **intenção de desenho**, não contrato fechado (princípio V; analyze **U1**). Rota, JWT, respostas HTTP e timeout aguardam contrato do Módulo (fonte: respostas-joao.md §Etapa 3; respostas-etapa3.md §10). Ver [plan.md](../plan.md) § API do Módulo Ocorrência.

## Transporte

| Item | Valor | Fonte |
|------|-------|-------|
| Protocolo | HTTP(S) interno | N1 A |
| Autenticação | JWT (padrão ADR_051 / ADR_028) | N1 A, 009-3 D |
| Operação semântica | Após **uma consulta ao Módulo por achado** (**DN-008-7 B**), entregar **um** `ResumoAnalise` ao fim da análise do lote | respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; ADR_062, seq. v6 §2; respostas-etapa3c.md §6 |
| URL / método / path | **BLOQUEADO (DN-009-2)** | respostas-etapa3.md §10 |
| Claims JWT | **BLOQUEADO (DN-009-2)** | ADR_028 [PENDENTE] |

## Corpo (JSON lógico)

```json
{
  "batchId": "<uuid lote ingested>",
  "execId": "<opcional — ausente no fluxo só log.ingested v1; DN-006-1b B>",
  "tenantId": "<metadata>",
  "productId": "<metadata>",
  "versaoRegras": "<versão única 006-3 C>",
  "identificacaoRegras": {
    "catalogoEntrou": false,
    "somenteNivel": true
  },
  "contadores": {
    "linhasProcessadas": 0,
    "erros": 0,
    "warnings": 0,
    "info": 0,
    "dadosSensiveisMascarados": 0
  },
  "achados": []
}
```

### Contadores

Alinhamento ADR_102 com revisão etapa 3b: contador **`info`**, **não** `spamIgnorado` (fonte: respostas-joao.md, Etapa 3; origem: ADR_102, ADR_170). Nomes finais podem ajustar com ADR publicada.

### AchadoAnalise (elemento de `achados[]`)

| Campo | Tipo | Notas |
|-------|------|-------|
| `tipo` | string | `ERROR` \| `WARNING` \| `INFO` |
| `linhas` | string[] | Mascaradas; podem estar vazias (010-16 C) |
| `recorrente` | boolean | ADR_063 — no escopo 006: **`false`** (FR-009; analyze **E3**) |

Campos **Drive**: `titulo`, `resumo`, `numeroLinha`, `camposExtraido`, `dadosSensiveis`, `traceId`, `correlationId`, `hostName`, `moduloId`. Presença exigida; semântica mínima de `titulo`/`resumo` **[PENDENTE N3 A]** (analyze **E3**). **`hostName`**: valor 1:1 com `host` da assinatura quando presente (analyze **L1**; ADR_040).

Campos **N3 A [PENDENTE ADR]**: `assinatura` — componentes DN-008-1 B; **valor** serializado **008 assinatura A** (texto concatenado; separador **`|`**, **DN-008-8 A**; escape **`|`** → **`\|`**, **DN-008-9 A**); `contagem`, `primeiraOcorrencia`, `ultimaOcorrencia`, `versaoRegras` por achado se aplicável.

## Idempotência

- **Chave de negócio**: `batchId` ou `batchId` + `execId` quando houver (N4 A, DN-006-1, DN-006-1b B).
- **Comportamento esperado do Módulo**: segunda entrega com mesma chave **não** duplica efeito (FA-02, ADR_014).
- **Formato de transporte da chave**: **BLOQUEADO (DN-009-2)** — header vs corpo.

## Resposta

> **Somente intenção** — cada linha abaixo é **[PENDENTE DN-009-2]** até contrato publicado pelo Módulo Ocorrência (analyze **U1**). Não usar como SLA nem contrato de implementação.

| Código | Significado | Status |
|--------|-------------|--------|
| 2xx + confirmação explícita | Handoff aceito; Finder pode concluir varredura | **[PENDENTE DN-009-2]** |
| 4xx | Erro de contrato; não concluir varredura | **[PENDENTE DN-009-2]** |
| 5xx / timeout | Falha reenviável; mesma chave no retry | **[PENDENTE DN-009-2]** |

Semântica exata (200 vs 202, corpo, timeout, lista retentável): **BLOQUEADO (DN-009-2)** (respostas-etapa3.md §10, opções 10a–10d).

## O que o Finder não envia

- `gravidade`, `criticidade` (ADR_062)
- Evento Redis / envelope ADR_146 de ocorrência (FR-009, SC-006)
- Texto bruto não mascarado (spec 010 quando ativa)

## Pré-condição (consulta por achado)

Antes do POST/PUT de envio (contrato **DN-009-2**), o Finder MUST ter executado **uma consulta ao Módulo por cada** `AchadoAnalise` do lote (**DN-008-7 B**). Falha de consulta: distinguir de falha de envio; política de retry **DN-009-2** / spec **008** (não inventar aqui).

## Fontes

- (fonte: respostas-joao.md, N1 A / N3 A / N4 A / 006-3 C / T7 A / T1 REVISTO / Etapa 3 contadores / DN-006-1b B / DN-008-1 B / §Etapa 3c (João, 16/11 11h02) DN-008-7 B / DN-008-5 A / 008 assinatura A / §Etapa 3e (João, 16/11 12h40) DN-008-8 A / §Etapa 3g (João, 16/11 13h22) DN-008-9 A)
- (fonte: revisao-T1-drive.md, §2 tabela “Contrato Finder → Módulo Ocorrência”)
- (fonte: spec.md, FR-007, FR-011)
- (fonte: respostas-etapa3.md, §10)

## Evidência desta execução

Matriz N3 A **[PENDENTE ADR]**; modelo de memória e DTO provisórios, sem decidir nomes finais:

| Campo lógico | Origem/prova | Situação externa |
|---|---|---|
| batchId / execId | Referencia.BatchID / id V ou M injetado pela Varredura; chave estruturada | execId omitido quando ausente; localização HTTP da chave depende DN-009-2 |
| tenantId / productId | Referência validada; produto body conferido | Sem identidade inventada |
| tipo / hostName / linhas | Classificação e contexto 010; host da assinatura | Somente três classes e trechos pós-segurança |
| assinatura | Compositor canônico 008 | [PENDENTE ADR] N3; excluído do JSON provisório |
| contagem | Partição local; fusão 008 quando utilizada | [PENDENTE ADR] N3 |
| primeiraOcorrencia / ultimaOcorrencia | Min/max timestamps do log | [PENDENTE ADR] N3 |
| versaoRegras do achado | Snapshot da entrada do lote | [PENDENTE ADR] N3 |
| titulo / resumo / numeroLinha / camposExtraido | Estrutura disponível | Semântica final pendente; sem texto de apresentação fabricado |
| dadosSensiveis | Ponteiro lógico desconhecido | nil/null até semântica acordada; não presume false |
| contadores.info | Quantidade de logs INFO | DTO 009 usa `informacao` provisoriamente; reconciliação ADR_102/DN-009-2 pendente |
| dadosSensiveisMascarados | Total informado pelo contexto 010 | Nome final/semântica N3/ADR_102 pendentes |

Um resumo por lote; uma consulta por achado antes do envio. Sem rota, claims, status, timeout ou retry HTTP inventados (DN-009-2). Em falha histórica: até três tentativas temporárias → porta de aviso Core com batchId/motivo; sem DLQ/WebSocket no ingested; cliente Core depende ADR_051.
