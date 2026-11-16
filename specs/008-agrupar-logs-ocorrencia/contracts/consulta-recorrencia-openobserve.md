# Contract: Consultar recorrência / similaridade no OpenObserve

Porta de saída do Finder para equivalência **histórica** fora do lote corrente (**T3 B**, **008-3 B**).

## Estado no repositório

`internal/adapter/out/openobserve/doc.go` declara que organização, stream, índice, campos e credenciais **não foram informados** e o pacote **não consulta** o armazenamento. Este contrato define **semântica de negócio**, não implementação de query.

## Entrada lógica

| Campo | Descrição |
|-------|-----------|
| `assinatura` ou critério derivado | Identidade do achado (componentes DN-008-1 B) |
| `classe` | ERROR \| WARNING \| INFO |
| `productId` / `tenantId` | Isolamento |
| `host` | Mesma máquina (ADR_013) |
| Janela | **Últimos 500 lotes** (**008-6 A**) |

## Critério `recorrente`

- Seguir **ADR_013**: ocorre **3–4 vezes** em **arquivos/lotes distintos**, **mesma máquina** (**008-5 A**).
- Repetição no **mesmo arquivo/lote** **não** conta para recorrência (**008-5 A**, ADR_063).
- **Recorrente** é **booleano no AchadoAnalise**, não status de ocorrência (ADR_063, FR-007).

## Saída

| Campo | Descrição |
|-------|-----------|
| `recorrente` | true/false |
| Evidências | **NEEDS CLARIFICATION** (detalhes da query até config OpenObserve) |

## Substituição ADR_074 / ADR_168

Similaridade não usa banco vetorial legado (**008-7 B**; FR-008).

## Falhas (decidido — DN-008-3, rev. etapa 3h)

| Tipo | Comportamento Finder |
|------|----------------------|
| Falha **temporária** (ex.: OpenObserve indisponível — contrato v1 §9) | Até **3 tentativas** na consulta; se esgotar, avisar o **Core** pela **API interna** com **batchId** e **motivo** (**sem DLQ** no `log.ingested` — **DLQ B**) |
| Erro **crítico** | **Sem retry** (alinhado 008-16 A, ADR_049) |
| Degradação | **MUST NOT** concluir com `recorrente=false` por indisponibilidade da consulta |

## Após 3 falhas — aviso ao Core (**DN-008-6**, rev. etapa 3h)

Após esgotar as **3 tentativas** na consulta, o Finder MUST notificar o **Core** pela **API interna**, enviando **batchId** e **motivo**. **MUST NOT** publicar em **DLQ** nem acionar fluxo de alerta **WebSocket pós-DLQ** (modelo da FilaOrdenada; contrato **log.ingested v1** §8; **DN-006-1b B**).

**Contrato Finder→Core:** **ADR_051 [PENDENTE]** — rota, payload e nome de evento **não** definidos nas fontes prioritárias; **MUST NOT** inventar nesta spec (distinto da publicação de evento de ocorrência, spec **009**).

(fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22); origem: DLQ B; respostas-etapa3g.md §1)

## Fontes

- (fonte: respostas-joao.md, DN-008-3 / T3 B / 008-3 B / 008-5 A / 008-6 A; §Etapa 3g DLQ B; origem: ADR_013, ADR_063, ADR_168, contrato v1 §8–§9)
- (fonte: internal/adapter/out/openobserve/doc.go)