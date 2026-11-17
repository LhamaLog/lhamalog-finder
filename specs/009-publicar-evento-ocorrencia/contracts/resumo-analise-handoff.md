# Contract: Handoff ResumoAnalise (Finder → Módulo Ocorrência)

Interface **destinada a implementação neste repositório**. Esta feature substitui a publicação direta de evento LLT-201 pelo Finder (**FR-F02**, **T1 revisto**).

> **Bloqueio de implementação (DN-009-2):** rota, claims JWT, códigos HTTP, timeout e corpo de confirmação **não** estão definidos (fonte: respostas-joao.md §Etapa 3 — *API Módulo: pendência no plan até o contrato*; origem: respostas-etapa3.md §10). Até o contrato do Módulo Ocorrência, **não** implementar o adaptador HTTP em `internal/adapter/out/modocorrencia/`. O restante (porta, montagem de payload, testes de domínio) pode avançar.

## Transporte

| Item | Valor | Fonte |
|------|-------|-------|
| Protocolo | HTTP(S) interno, rede privada | N1 A; ADR_141 |
| Autenticação | JWT próprio (padrão ADR_028 / ADR_054) | N1 A; respostas-etapa3.md §10 |
| Semântica | Entregar `ResumoAnalise` ao fim da análise do lote/arquivo | ADR_062, seq. v6 §2 |
| Granularidade handoff | **Exatamente uma** chamada por lote — **um** `ResumoAnalise` por lote analisado (**DN-008-7**) | respostas-etapa3c.md §6; ADR_062; spec.md **FR-F01** |
| URL base / path | **Pendente DN-009-2** | respostas-etapa3.md §10 |
| Método HTTP exato | **Pendente DN-009-2** | seq. v6 §2 — Drive só diz “envia e confirma” |

## Corpo (JSON lógico)

Alinhado a [data-model.md](../data-model.md). Exemplo estrutural (campos mínimos com fonte):

```json
{
  "batchId": "<uuid lote ingested>",
  "tenantId": "<metadata>",
  "productId": "<metadata>",
  "versaoRegras": "<versão única 006-3 C>",
  "contadores": {
    "linhasProcessadas": 0,
    "erros": 0,
    "warnings": 0,
    "informacao": 0,
    "dadosSensiveisMascarados": 0
  },
  "achados": []
}
```

`execId` **omitido** no fluxo **log.ingested v1** atual (**DN-006-1b B**). Incluir quando disponível (**DN-006-1**, **N4 A**):

```json
{
  "batchId": "<uuid>",
  "execId": "V<uuid> ou M<uuid>",
  ...
}
```

**Notas**:

- Contador de achados **INFO**: substitui o legado “spam ignorado” (**spam: B**); no exemplo JSON use chave provisória `informacao` até a **ADR_102** publicar o nome canônico (**corrige analyze F3**, fonte: respostas-joao.md §Etapa 3, spam: B; origem: ADR_102).
- Rótulos exatos dos demais contadores **devem** alinhar à **ADR_102** quando publicada; lista detalhada no Drive permanece **[INFERIDO]** até decisão fechada (revisao-T1-drive.md §2) — não fixar rótulos adicionais aqui (constituição V).

### AchadoAnalise (`achados[]`)

| Campo | Tipo | Obrig. | Fonte |
|-------|------|--------|-------|
| `tipo` | string | sim | `ERROR` \| `WARNING` \| `INFO` — FR-F08, 009-13 A |
| `linhas` | string[] | sim | Mascaradas; podem ser vazias (010-16 C) |
| `recorrente` | boolean | sim | ADR_063; política operacional na spec **008** |
| Campos Drive | vários | conforme ADR_040 | Finder.mmd |
| `assinatura`, `contagem`, `primeiraOcorrencia`, `ultimaOcorrencia`, `versaoRegras` | — | **[PENDENTE ADR N3 A]** | respostas-joao.md, N3 A |

Quando presente, **`assinatura`** MUST incorporar **mensagem normalizada** **DN-008-4 B** com literais **DN-008-5 A** (`<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>`) e serialização **008 assinatura A** (texto concatenando partes **DN-008-1 B** com separador fixo **`|`** — **DN-008-8 A**; **`|`** nos componentes → **`\|`** — **DN-008-9 A**; escape de **`\`** nos componentes **[PENDENTE]** sem fonte; dono canônico spec **008**) (**FR-F11**; fonte: respostas-joao.md §Etapa 3b/3c 11h; §Etapa 3e (João, 16/11 12h40), DN-008-8 A; §Etapa 3g, DN-008-9 A; origem: DN-008-1 B).

**Falha consulta OpenObserve** (não confundir com campos do handoff): até **3 tentativas**; **sem DLQ** no **log.ingested v1** (**DLQ B**); após esgotar, **Finder → Core** por API interna (**`batchId`**, **motivo** — **DN-008-3**/**DN-008-6** rev.); rota/payload/evento **[PENDENTE ADR_051]**; **sem** handoff de sucesso (**FR-F07**; fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B; origem: respostas-etapa3g.md §1).

## Consulta prévia (**N2 B**, **FR-F09**)

| Item | Valor | Fonte |
|------|-------|-------|
| Momento | Antes de montar/enviar o **ResumoAnalise** do lote | N2 B, FR-F09 |
| Granularidade | **Uma** chamada HTTP ao Módulo **por achado** (assinatura/classe elegível) | **DN-008-7 B** (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); origem: N2 B) |
| Distinto do handoff | Consultas **não** substituem o **único** POST do resumo por lote (**DN-008-7**) | respostas-etapa3c.md §6 |
| Path / JWT / corpo | **Pendente DN-009-2** | respostas-etapa3.md §10 |

## Idempotência

| Item | Valor | Fonte |
|------|-------|-------|
| Chave (fluxo atual v1) | **`batchId`** somente | DN-006-1b B (fonte: respostas-joao.md §Etapa 3; origem: contrato 08/11) |
| Chave (com varredura) | **`batchId` + `execId`** quando houver | N4 A, DN-006-1 |
| Comportamento Módulo | Segunda entrega mesma chave: **sem** efeito colateral indevido | FR-F06, FA-02, ADR_014 |
| Transporte da chave | **Pendente DN-009-2** (header dedicado vs corpo) | — |

## Resposta (confirmação)

| Código | Significado | Fonte |
|--------|-------------|-------|
| 2xx + confirmação explícita | Handoff aceito; Finder pode concluir varredura | FR-F05, SC-F001 |
| 4xx | Erro de contrato; **não** concluir varredura | FA-02; detalhes **DN-009-2** |
| 5xx / timeout | Reenviável com **mesma** chave | FR-F06; retry sugerido 10d-A (recomendação DN-009-2, não decisão) |

Corpo da confirmação (JSON) e códigos exatos: **Pendente DN-009-2** — além de “confirma o recebimento” (seq. v6 §2).

## Comportamentos proibidos (Finder)

- Publicar em Redis `notificacao-eventos` ou qualquer stream de aviso LLT-201 (**FR-F02**).
- Chamar API Core para aviso de ocorrência (**T1 revisto**).
- Incluir `gravidade` ou `criticidade` no payload (**ADR_062**).
- Incluir **id** de ocorrência persistida (quem gera id é o Módulo).

## Dependências

- Classificação e montagem dos achados: spec **006** (**T6 A**).
- Agrupamento entre lotes: consulta **por achado** (**DN-008-7 B**, **N2 B**, **FR-F09**); handoff **um resumo por lote** (**DN-008-7**); assinatura (**DN-008-1 B**, **008 assinatura A**; serialização com separador fixo **`|`** — **DN-008-8 A**; fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40)); normalização (**DN-008-4 B**, **DN-008-5 A**, **FR-F11**).
- Lote incompleto: **não** enviar handoff de sucesso (**T4 A**).

## Fontes consolidadas

- (fonte: respostas-joao.md, N1 A / N4 A / DN-006-1b B / T7 A / T1 revisto / 009-20 T6 A / §Etapa 3)
- (fonte: revisao-T1-drive.md, §1 fluxo; §2 contrato Finder → Módulo)
- (fonte: respostas-etapa3.md, §3 DN-006-1; §10 API Módulo)
- (fonte: spec.md, FR-F01–FR-F11, SC-F001–SC-F003, DN-009-2, DN-008-4 B)

## Evidência desta execução

O serializador implementado é um DTO **provisório**, não contrato HTTP publicado: comum 006 conta `info`, DTO 009 emite `informacao` conforme esta spec; conciliação ADR_102/DN-009-2 pendente. Campos N3 A não são serializados. dadosSensiveis=nil é null/desconhecido, sem assumir false. Titulo/resumo/numeroLinha/camposExtraido têm estrutura, sem conteúdo de apresentação inventado.

Snapshot do resumo e chave são copiados no estado Pendente. Reenviar conserva exatamente esse conteúdo e batchId/execId. Confirmação lógica deve anteceder callback de conclusão; não define qual resposta HTTP equivale a aceitação, nem header/corpo da chave (DN-009-2).
