# Contract: Handoff ResumoAnalise (Finder → Módulo Ocorrência)

Interface **destinada a implementação neste repositório**. Esta feature substitui a publicação direta de evento LLT-201 pelo Finder (**FR-F02**, **T1 revisto**).

O cliente do Finder segue o [contrato DN-009-2 aprovado em 20/11/2026](proposta-dn-009-2.md). O [OpenAPI](proposta-dn-009-2.openapi.yaml) define os campos e as respostas. A homologação com a API do Módulo segue pendente.

## Transporte

| Item | Valor | Fonte |
|------|-------|-------|
| Protocolo | HTTPS interno, rede privada | N1 A; ADR_141; DN-009-2 |
| Autenticação | JWT próprio (padrão ADR_028 / ADR_054) | N1 A; respostas-etapa3.md §10 |
| Semântica | Entregar `ResumoAnalise` ao fim da análise do lote/arquivo | ADR_062, seq. v6 §2 |
| Granularidade handoff | Um resumo por lote, com a mesma chave em caso de reenvio | DN-008-7; DN-009-2 |
| URL base / path | Base HTTPS configurada para a implantação | DN-009-2 |
| Método e rota | `POST /internal/v1/resumos-analise` | DN-009-2 |

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

O contador `informacao` representa a classe `INFO` e a chave `dadosSensiveisMascarados` carrega o total devolvido pelo scanner.

### AchadoAnalise (`achados[]`)

| Campo | Tipo | Obrig. | Fonte |
|-------|------|--------|-------|
| `tipo` | string | sim | `ERROR` \| `WARNING` \| `INFO` — FR-F08, 009-13 A |
| `linhas` | string[] | sim | Mascaradas; podem ser vazias (010-16 C) |
| `recorrente` | boolean | sim | ADR_063; política operacional na spec **008** |
| `assinatura`, `hostName`, `contagem`, `primeiraOcorrencia`, `ultimaOcorrencia` | string, string, inteiro e datas UTC | sim | N3 A, DN-009-2 |
| `revisaoConsultada` | string ou null | sim | DN-009-2, token do estado ou ausência |
| `versaoRegras` | string | herdada do resumo | N3 A, DN-009-2 |
| Campos de apresentação | conforme o esquema OpenAPI | null quando indisponível | DN-009-2 |

A **`assinatura`** MUST incorporar **mensagem normalizada** **DN-008-4 B** com literais **DN-008-5 A** (`<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>`) e serialização **008 assinatura A** (texto concatenando partes **DN-008-1 B** com separador fixo **`|`** — **DN-008-8 A**; **`|`** nos componentes → **`\|`** — **DN-008-9 A**; escape de **`\`** nos componentes **[PENDENTE]** sem fonte; dono canônico spec **008**) (**FR-F11**; fonte: respostas-joao.md §Etapa 3b/3c 11h; §Etapa 3e (João, 16/11 12h40), DN-008-8 A; §Etapa 3g, DN-008-9 A; origem: DN-008-1 B).

**Falha consulta OpenObserve** (não confundir com campos do handoff): até **3 tentativas**; **sem DLQ** no **log.ingested v1** (**DLQ B**); após esgotar, **Finder → Core** por API interna (**`batchId`**, **motivo** — **DN-008-3**/**DN-008-6** rev.); rota/payload/evento **[PENDENTE ADR_051]**; **sem** handoff de sucesso (**FR-F07**; fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B; origem: respostas-etapa3g.md §1).

## Consulta prévia (**N2 B**, **FR-F09**)

| Item | Valor | Fonte |
|------|-------|-------|
| Momento | Antes de montar/enviar o **ResumoAnalise** do lote | N2 B, FR-F09 |
| Granularidade | Uma consulta por achado elegível em cada preparação; nova rodada em caso de revisão desatualizada | DN-008-7 B; DN-009-2 |
| Distinto do handoff | Consultas **não** substituem o **único** POST do resumo por lote (**DN-008-7**) | respostas-etapa3c.md §6 |
| Rota | `POST /internal/v1/ocorrencias/consultas` com Bearer JWT; estado nulo para ausência, revisão opaca quando há estado | DN-009-2 |

## Idempotência

| Item | Valor | Fonte |
|------|-------|-------|
| Chave (fluxo atual v1) | **`batchId`** somente | DN-006-1b B (fonte: respostas-joao.md §Etapa 3; origem: contrato 08/11) |
| Chave (com varredura) | **`batchId` + `execId`** quando houver | N4 A, DN-006-1 |
| Comportamento Módulo | Segunda entrega mesma chave: **sem** efeito colateral indevido | FR-F06, FA-02, ADR_014 |
| Transporte da chave | `batchId` e `execId` opcional no JSON, no escopo `tenantId`/`productId` | DN-009-2 |

## Resposta (confirmação)

| Código | Significado | Fonte |
|--------|-------------|-------|
| `201` ou `200` com recibo correspondente | Aceitação inicial ou repetição confirmada. O recibo deve confirmar a chave completa | FR-F05, SC-F001, DN-009-2 |
| `409 REVISAO_DESATUALIZADA` | Recalcular uma vez e reenviar com a mesma chave | DN-009-2 |
| Falha, timeout ou recibo inválido | Manter pendente; falhas temporárias são reenviáveis com a mesma chave | FR-F06, DN-009-2 |

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

Os testes verificam a serialização e as respostas HTTPS. A integração com Redis comprova que o recibo válido antecede a finalização e o `XACK`, inclusive em reenvio. [Comandos e evidências](../quickstart.md#verificação-do-contrato-dn-009-2)
