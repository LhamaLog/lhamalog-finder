# Data Model: Publicar evento de ocorrência (escopo Finder)

**Feature**: `009-publicar-evento-ocorrencia` | **Date**: 2026-11-16 (etapa 3h)

Entidades do **handoff Finder → Módulo Ocorrência** e modelo **lógico** do aviso LLT-201 no produto (persistência e publicação **fora** deste repositório). Classificação: spec **006**; agrupamento: spec **008**.

---

## ChaveIdempotenciaHandoff

| Componente | Descrição | Fonte |
|------------|-----------|-------|
| `batchId` | UUID do lote ingested (`Referencia.BatchID` / `idempotencyId` no envelope) | N4 A; DN-006-1b B; `internal/domain/envelope.go` |
| `execId` | Identificador Varredura `V<uuid>` ou `M<uuid>` | DN-006-1 (fonte: respostas-joao.md §Etapa 3; origem: ADR_075, ADR_069) — **opcional** |

**Regra de composição**

- Fluxo atual (**log.ingested** v1, sem `execId` no envelope): chave = **`batchId` somente** (**DN-006-1b B**).
- Quando `execId` estiver disponível ao Finder: chave = **`batchId` + `execId`** (**N4 A**).

**Invariante**: segunda entrega com a mesma chave **não** deve duplicar efeito no Módulo (**FR-F06**, SC-F002, ADR_014).

---

## ResumoAnalise (mensagem Finder → Módulo)

Handoff ao fim da análise do arquivo/lote (**ADR_062**, revisao-T1-drive.md §2).

| Campo | Obrig. | Descrição | Fonte |
|-------|--------|-----------|-------|
| `batchId` | sim | Lote ingested | N4 A, DN-006-1b B, contrato v1 |
| `execId` | não | Presente apenas quando o Finder dispuser do id de varredura | DN-006-1, N4 A |
| `tenantId` | sim | Isolamento multi-tenant | contrato metadata |
| `productId` | sim | Produto analisado | contrato §4 |
| `versaoRegras` | sim | Versão única das regras de classificação | 006-3 C (revisao-T1-drive §4) |
| `identificacaoRegras` | opcional | Metadados do catálogo (ex. só nível OTel) | 006-3 C |
| `contadores` | sim | Contadores do lote (linhas, erros, warnings, **contador de INFO**, mascarados, …); legado “spam ignorado” → **contador de INFO** (**spam: B**); nome canônico na lista **ADR_102** [PENDENTE] (**corrige analyze F3**, fonte: respostas-joao.md §Etapa 3, spam: B; origem: ADR_102, ADR_170) | ADR_102 |
| `achados` | sim | `AchadoAnalise[]` — pode ser vazio se lote sem achados elegíveis | ADR_062 |

**O que não contém**: `gravidade`, `criticidade`, id de ocorrência persistida, envelope Redis LLT-201 (**FR-F02**, ADR_062).

**Granularidade**: **exatamente um** payload por lote analisado (**DN-008-7**; fonte: respostas-etapa3c.md §6; origem: ADR_062). **MUST NOT** enviar HTTP separado por achado no handoff do resumo (**FR-F01**).

**Relação com código**: montado após análise do lote (e consultas **N2 B**, quando aplicável); entrega via porta `ModuloOcorrencia` — implementação HTTP **bloqueada** por **DN-009-2**.

---

## AchadoAnalise (elemento de `achados[]`)

| Campo | Obrig. | Descrição | Fonte |
|-------|--------|-----------|-------|
| `tipo` | sim | `ERROR` \| `WARNING` \| `INFO` | FR-F08, 009-13 A |
| `linhas` | sim | Trechos mascarados (podem ser `[]` se 010-16 C) | ADR_040; 010-21 rev. |
| `recorrente` | sim | Booleano de marca recorrente no achado; semântica operacional na spec **008** | ADR_063 |
| Campos Drive | conforme ADR | `titulo`, `resumo`, `numeroLinha`, `traceId`, `correlationId`, … | Finder.mmd, ADR_040 |

### Campos N3 A — [PENDENTE ADR]

| Campo | Descrição | Fonte |
|-------|-----------|-------|
| `assinatura` | Chave de agrupamento; valor = texto concatenando partes **DN-008-1 B** com separador fixo **`|`** (**008 assinatura A**; **DN-008-8 A**); **`|`** nos componentes → **`\|`** (**DN-008-9 A**); escape de **`\`** **[PENDENTE]**; dono canônico spec **008** | N3 A; DN-008-1 B; 008 assinatura A; DN-008-8 A; DN-008-9 A (fonte: respostas-joao.md §Etapa 3c 11h; §Etapa 3e DN-008-8 A; §Etapa 3g DN-008-9 A) |
| `contagem` | Após fusão entre lotes | N3 A, N2 B |
| `primeiraOcorrencia` / `ultimaOcorrencia` | Extremos temporais | N3 A, 008-14 A |
| `versaoRegras` | Por achado se não só no resumo | N3 A |

Até ADR publicada, handoff **pode** omitir estes campos no aceite inicial (**corrige analyze U2**); **FR-F09** / **N2 B** permanecem requisito de desenho nesta feature, com aceite parcial do payload até ADR + spec **008**; agrupamento operacional depende dos campos **N3 A**. Semântica **DN-008-4 B** / **DN-008-5 A** / **008 assinatura A** / **DN-008-8 A** fechada (etapas 3c–3g, **FR-F11**); literais de marcador fechados (**DN-008-5 A**); separador fixo **`|`** propagado (**DN-008-8 A**; dono canônico spec **008**) (fonte: respostas-joao.md, N3 A; §Etapa 3c 11h; §Etapa 3e DN-008-8 A).

### Mensagem normalizada (entrada da assinatura)

| Etapa | Regra | Fonte |
|-------|-------|-------|
| Entrada | Texto do log após mascaramento quando aplicável (spec **010**) | ADR_167; FR-F11 |
| Trim | Remover espaços nas extremidades | DN-008-4 B (fonte: respostas-joao.md §Etapa 3b) |
| Case | Ignorar maiúsculas/minúsculas na comparação | DN-008-4 B |
| Marcadores | Substituir por literais **DN-008-5 A**: `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` | DN-008-4 B; DN-008-5 A (fonte: respostas-joao.md §Etapa 3c 11h) |
| Dono do algoritmo | Ordem de aplicação e definição canônica do separador **`|`** (**DN-008-8 A**) na spec **008**; 009 reutiliza o literal | FR-F11; spec.md Clarifications 3g |

O Finder MUST NOT publicar regras de marcador divergentes da **008** nesta feature (**FR-F11**).

---

## ConsultaOcorrenciaExistente (Finder → Módulo, pré-handoff)

Consulta acionada por **N2 B** antes de montar o **ResumoAnalise**, coordenada com spec **008** (**FR-F09**).

| Aspecto | Regra | Fonte |
|---------|-------|-------|
| Granularidade | **Uma** chamada ao Módulo **por achado** (assinatura/classe elegível) | **DN-008-7 B** (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); origem: N2 B) |
| Distinto do handoff | **Não** substitui o envio único do **ResumoAnalise** por lote (**DN-008-7**) | respostas-etapa3c.md §6 |
| Contrato HTTP | Rota, JWT, corpo de resposta: **DN-009-2** | respostas-etapa3.md §10 |

**Finder MUST NOT** assumir uma única consulta agregada por lote para **N2 B** (**DN-008-7 B**).

---

## ConfirmacaoHandoff

Resposta lógica do Módulo Ocorrência ao Finder (**FR-F05**).

| Campo | Descrição | Fonte |
|-------|-----------|-------|
| Confirmação explícita | Handoff aceito | seq. v6 §2 |
| Corpo HTTP / códigos | **Pendente DN-009-2** — Drive não fixa além de “confirma o recebimento” | respostas-etapa3.md §10 |

**Transição de estado (Finder)**: análise concluída → handoff enviado → **aguardando confirmação** → varredura concluída.

**Ponte `finishScan()` (corrige analyze A2)**: na sequência varredura v6 §2, **`finishScan()`** é o marco de conclusão da varredura **somente após** confirmação do Módulo (**FR-F05**). No repositório Go atual não há símbolo homônimo (grep `finishScan` em `*.go` vazio); na implementação, mapear para a porta de aplicação que encerra o ciclo do lote (ex. callback na orquestração pós-`ModuloOcorrencia.Confirmar`) — **sem inventar** nome de função até o código existir (fonte: sequencia-varredura-finder.md v6 §2; plan.md Technical Context).

---

## OcorrenciaNegocio (Módulo Ocorrência — lógico)

Entidade criada/atualizada a partir do ResumoAnalise (**T1 revisto**).

| Atributo | Regra | Fonte |
|----------|-------|-------|
| `id` | Id próprio opaco; referenciado no aviso e lido pela Notificação | 009-6 A, ADR_040, DN-009-1 A |
| Chave unicidade (C8) | **[PENDENTE]** na ADR_040 | spec.md C8; 009-6 A |
| `gravidade` | Calculada no Módulo | 009-12 rev., ADR_062 |
| `tipo` (classe) | ERROR/WARNING/INFO | 009-13 A |
| `conteudo` | Trechos mascarados para leitura posterior (Notificação/e-mail) | ADR_040, ADR_167, DN-009-1 A |
| Natureza mudança | nova \| atualizada \| recorrente | FR-003, 009-14 A |

**Finder MUST NOT** criar nem atualizar esta entidade (**FR-F02**).

---

## AvisoOcorrencia (evento lógico LLT-201)

Mensagem emitida pelo **Módulo Ocorrência** ao **Core** (**FR-001**, SSID_009.1 passo 1).

### Payload mínimo HTTP (critério LLT-201 + **009-5 revisto**)

| Campo | Descrição | Fonte |
|-------|-----------|-------|
| `tenant` | Tenant | FR-003 |
| `produto` | Produto | FR-003 |
| `ocorrencia` / id | **Id** opaco da ocorrência | 009-5 revisto; 009-6 A |
| `tipo` | Classe ERROR/WARNING/INFO | 009-13 A, 009-5 revisto |
| `gravidade` | No aviso do Módulo | 009-12 rev., 009-5 revisto, SSID_009.1 |
| natureza da mudança | nova \| atualizada \| recorrente | FR-003, 009-14 A |
| marca `recorrente` | Quando aplicável | SSID_009.1 passo 1 |

**Sem texto bruto de log** no aviso (**010-23** via 009-5 A histórico + DN-009-1 A).

### Payload stream (Core → Notificação)

| Campo | Descrição | Fonte |
|-------|-----------|-------|
| `ocorrencia` (payload/envelope) | **Somente id** opaco | 009-6 A |
| `tipo`, `gravidade` | Já consumidos pelo Core do HTTP upstream; refletidos no envelope conforme ADR_146 | 009-4 B, 009-5 revisto |

**Transporte físico**: HTTP + JWT com envelope **ADR_146** no corpo (**009-4 B**). Detalhes: [pipeline-aviso-referencia.md](./contracts/pipeline-aviso-referencia.md).

---

## ChaveIdempotenciaAvisoHTTP (Core)

| Componente | Descrição | Fonte |
|------------|-----------|-------|
| `ocorrenciaId` | Id da ocorrência | 009-5 revisto / 009-6 A |
| `estado` | Natureza da mudança | 009-9 A |

**Efeito duplicado a evitar**: segundo alerta em `notificacao-eventos` (**009-11 B**, SC-003).

---

## EventoNotificacao (Core → Notificação)

| Aspecto | Valor | Fonte |
|---------|-------|-------|
| Stream | `notificacao-eventos` | FR-009, ADR_165 |
| Envelope | ADR_146 v1 | 009-4 B (trecho stream); ADR_146; **T5 D** (não confundir com ingested v1) |
| Idempotência consumidor | `eventId` | 009-10 B, guia §8 |
| Texto do e-mail | Notificação lê **Ocorrência** por id no MongoDB (**DN-009-1 A**) | respostas-joao.md §Etapa 3; ADR_054 |

---

## Relacionamentos

```text
1 Referencia (lote v1) --> N LogAnalisado --> N AchadoAnalise
N AchadoAnalise --consulta por achado (N2 B)--> Módulo Ocorrência  [DN-008-7 B; contrato DN-009-2]
N AchadoAnalise --> 1 ResumoAnalise --HTTP+JWT (1x por lote)--> Módulo Ocorrência  [DN-008-7]
Módulo Ocorrência --> 0..N OcorrenciaNegocio
1 OcorrenciaNegocio --aviso HTTP--> Core --stream--> Notificacao
Notificacao --leitura id--> OcorrenciaNegocio (MongoDB compartilhado)
```

---

## Validações (Finder)

| Regra | Fonte |
|-------|-------|
| Não publicar evento LLT-201 / não chamar Core para aviso | FR-F02, T1 revisto |
| Classe fora ERROR/WARNING/INFO: proibida | FR-F08, SC-F003 |
| Lote incompleto: sem handoff de sucesso | T4 A |
| Mesma chave no reenvio após falha | FR-F06 |
| INFO presente quando classificado | FR-F04, T7 A |
| Coordenar agrupamento entre lotes: consulta **por achado** antes do envio | FR-F09, N2 B, DN-008-7 B |
| Enviar **um** ResumoAnalise por lote | FR-F01, DN-008-7 |
| Aplicar **DN-008-4 B** / **DN-008-5 A** / **008 assinatura A** (separador fixo **`|`** — **DN-008-8 A**; escape **`|`** → **`\|`** — **DN-008-9 A**) ao compor assinatura | FR-F11 |
| Falha consulta OpenObserve: até **3 tentativas**; **sem DLQ** ingested (**DLQ B**); depois Finder→Core API (**batchId**, motivo); contrato **[PENDENTE ADR_051]** | FR-F07, DN-008-3, DN-008-6 (rev.) |

## Evidência desta execução

Tipos ResumoAnalise/AchadoAnalise/ChaveResumo são reutilizados da 006. C8 (chave concreta de ocorrência ADR_040) e N3 A seguem pendentes; 009-6 A permanece rastreado no modelo, sem nova decisão. Semântica externa de dadosSensiveis não fechada; nil é desconhecido. Normalização/composição e agrupamento usam exclusivamente a 008; não existe normalizador paralelo no domínio da 009.
