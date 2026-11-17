# Contract (referência): Aviso de ocorrência — Módulo Ocorrência → Core → Notificação

**Não implementável no repositório lhamalog-finder** (**FR-F02**, spec.md tabela de responsabilidades). Documenta o restante da **LLT-201** para testes integrados e rastreio de requisitos (**User Story 2**).

## Visão do fluxo

```text
Módulo Ocorrência --HTTP+JWT+ADR_146--> Core --Redis notificacao-eventos--> Notificação
                                                      |
Notificação --leitura id--> Ocorrência (MongoDB compartilhado, DN-009-1 A)
```

(fonte: revisao-T1-drive.md §1; respostas-joao.md, 009-3 D, 009-1 B, DN-009-1 A)

---

## Trecho 1 — Módulo Ocorrência → Core (HTTP)

| Item | Valor | Fonte |
|------|-------|-------|
| Emissor | **Módulo Ocorrência** (não Finder) | T1 revisto, FR-001, SSID_009.1 |
| Protocolo | HTTP(S) interno com JWT | 009-3 D (fonte: respostas-joao.md; origem: ADR_051) |
| Momento | Ao **criar** ou **mudar** ocorrência; após Sugestão automática quando aplicável | 009-16 rev., SSID_009.1 passos 1–2 |
| Granularidade | **Um aviso por ocorrência** por evento de criação/mudança tratado | 009-17 |
| URL / rota Core | **Sem fonte** nas três prioritárias | respostas-etapa3.md §10 |

### Corpo HTTP

- **Decision**: Envelope **ADR_146** transportado por HTTP; campos LLT-201 no `payload` (**009-4 B**).
- Campos mínimos lógicos no `payload` (**009-5 revisto** — substitui leitura “só id” do aviso HTTP):

| Campo | Descrição | Fonte |
|-------|-----------|-------|
| `tenant` | Tenant | FR-003 |
| `produto` | Produto | FR-003 |
| `ocorrencia` / id | **Id** opaco da ocorrência | 009-5 revisto; 009-6 A |
| `tipo` | Classe ERROR/WARNING/INFO | 009-13 A, 009-5 revisto |
| natureza da mudança | `nova` \| `atualizada` \| `recorrente` | FR-003, 009-14 A |
| `gravidade` | Calculada no Módulo | 009-12 rev., 009-5 revisto, SSID_009.1 |
| marca recorrente | Quando aplicável | SSID_009.1 passo 1 |

**Não incluir no aviso**: texto bruto de log (**010-23**; trechos ficam na Ocorrência — **DN-009-1 A**).

### Idempotência (Core)

| Item | Valor | Fonte |
|------|-------|-------|
| Chave deduplicação HTTP | **id da ocorrência + estado** (natureza da mudança) | 009-9 A |
| Efeito duplicado a evitar | Segundo alerta publicado pelo **Core** em `notificacao-eventos` | 009-11 B, SC-003 |

### Falha

- **Decision**: Tentar de novo e alertar (**009-18 A**).
- (fonte: respostas-joao.md, 009-18 A)

---

## Trecho 2 — Core → Notificação (Redis Streams)

| Item | Valor | Fonte |
|------|-------|-------|
| Stream | `notificacao-eventos` (papel ADR_143) | FR-009, ADR_165 |
| Envelope | **ADR_146 v1** (não confundir com `log.ingested` v1 — **T5 D**, D7) | ADR_146; spec.md |
| Campo `ocorrencia` | **Somente id** opaco no payload publicado | 009-6 A |
| Destino | Módulo Notificação consome stream | 009-1 B |

### Idempotência (consumidor)

| Item | Valor | Fonte |
|------|-------|-------|
| Onde | **Consumidor** (padrão guia §8) | 009-10 B |
| Chave | `eventId` no envelope | guia-redis-streams v3 §8 |

### Conteúdo do e-mail (LLT-218)

- **Decision**: A **Notificação** lê a **Ocorrência** pelo **id** no **MongoDB compartilhado** (**DN-009-1 A**) e usa `conteudo` mascarado (ADR_040, ADR_167). Templates de corpo: **[PENDENTE]** NOT-03 (respostas-etapa3.md §7).

---

## Estados nova / atualizada / recorrente

| Estado | Significado no aviso | Fonte |
|--------|----------------------|-------|
| `nova` | Ocorrência criada | LLT-201 |
| `atualizada` | Ocorrência alterada; **reabertura** mapeia aqui | 009-14 A, 008-4 B |
| `recorrente` | Marca recorrente; **nova notificação** no fluxo acordado | 009-15 A, SSID_009.1 FA01 |

Aceite **100%** em cenários **atualizada** e **recorrente** (**SC-002**) bloqueado até spec **008** (**009-22 A**).

---

## Aviso operacional Finder → Core após falha OpenObserve (não é aviso LLT-201)

| Item | Valor | Fonte |
|------|-------|-------|
| Quando | Após **3 tentativas** esgotadas na **consulta OpenObserve** (recorrência **T3 B**) | **DN-008-3**, **FR-F07** |
| Canal | **Finder** → **Core** por **API interna** (**`batchId`**, **motivo**); **sem DLQ** no ingested (**DLQ B**) | respostas-joao.md §Etapa 3g; **DN-008-6** (rev.) |
| Contrato | Rota, payload, evento **[PENDENTE]** (**ADR_051 [PENDENTE]**); **não** inventar nesta feature | spec.md §Questões em aberto; respostas-joao.md §Etapa 3g, DLQ B |
| Distinção | **Não** confundir com HTTP Módulo → Core de ocorrência (**009-3 D**); **não** DLQ ingested, consumidor `dlq-alerta` nem **WebSocket** pós-DLQ (substitui fluxo FilaOrdenada) | spec.md **FR-F07**; respostas-etapa3g.md §1b |

---

## O que o Finder valida indiretamente

- `AchadoAnalise.tipo` válido — alimenta `tipo` no aviso (**SC-F003**, **SC-004**).
- **Assinatura** com **DN-008-4 B** / **DN-008-5 A** / **008 assinatura A** (separador fixo **`|`** — **DN-008-8 A**; **`|`** nos componentes → **`\|`** — **DN-008-9 A**; fonte: respostas-joao.md §Etapa 3e DN-008-8 A; §Etapa 3g DN-008-9 A) — agrupamento estável antes do handoff (**FR-F11**); não altera campos do aviso HTTP.
- **Um** handoff por lote (**DN-008-7**); consultas **por achado** (**DN-008-7 B**) quando **N2 B** aplicável.
- `AchadoAnalise.linhas` mascaradas — suporte a LLT-218 via persistência no Módulo (**DN-009-1 A**).
- `recorrente` no achado — sinal para Módulo; política fina na Notificação (**008-17 C**, LLT-216).

---

## Fontes consolidadas

- (fonte: respostas-joao.md, 009-3 D, 009-4 B, 009-5 revisto, 009-9 A, 009-11 B, 009-14 A, 009-15 A, 009-18 A, DN-009-1 A)
- (fonte: respostas-documentadas.md, 009-1 B, 009-6 A, 009-10 B, 009-13 A)
- (fonte: revisao-T1-drive.md, §1 e §4)
- (fonte: respostas-etapa3.md, §7, §10)
- (fonte: spec.md, FR-001–FR-012, SC-001–SC-006)

## Evidência desta execução

Revisão de responsabilidade (documental, sem evidência executada downstream): Módulo→Core leva tenant, produto, id da ocorrência, tipo, gravidade e natureza da mudança (009-5 revisto). Core→Notificação usa notificacao-eventos, ADR_146 v1/009-10 B, referencia ocorrência por id; a Notificação lê conteúdo no MongoDB ADR_054 para LLT-218 (DN-009-1 A).

Idempotência Módulo→Core: id da ocorrência + estado; duplicata não pode produzir segundo alerta (009-11 B). Essa chave não é batchId/execId do handoff Finder. Falha 009-18 A desse trecho não define a política do handoff nem o aviso operacional FR-F07. Nenhum desses serviços externos foi executado nesta validação.
