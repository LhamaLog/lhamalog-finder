# Data Model: Modelo de classificação de logs

**Feature**: `006-modelo-classificacao-logs` | **Date**: 2026-11-16 (etapa 3h)

Entidades de domínio e payloads de handoff. Persistência de **Ocorrência** fica no **Módulo de Ocorrência** (T1 REVISTO).

---

## ClasseLog

Enum lógico de saída da classificação de **um** log.

| Valor | Constituição | Termos de nível (exemplos da spec) |
|-------|--------------|-------------------------------------|
| `ERROR` | erro | error, erro, fatal, … |
| `WARNING` | alerta | warn, warning, alerta, … |
| `INFO` | informação | info, informação, debug, trace, … |

- **Regra de precedência**: entre sinal do **campo de nível** e **palavras inteiras** na mensagem, vence a classe mais grave (ERROR > WARNING > INFO). (fonte: spec.md FR-004; User Story 1)
- **Mapeamento**: único para as três categorias constitucionais (fonte: respostas-joao.md, 006-6 A; origem: ADR_170)

---

## LogClassificado

| Campo | Obrigatório | Descrição | Validação |
|-------|-------------|-----------|-----------|
| `apontamento` | sim | Código do log no lote (`domain.Codigo`) | Não vazio para log contável |
| `classe` | sim | `ClasseLog` | Uma das três |
| `tenantId` | sim | Do envelope ingested v1 | `Referencia.TenantID` |
| `productId` | sim | Do envelope | `Referencia.ProductID` |
| `mensagem` | uso interno | Texto **mascarado** quando spec 010 aplicável | Não enviar a LLM (FR-002) |
| `nivelCampo` | opcional | Valor OTel de severidade/texto | Usado na regra FR-004 |
| `serviceName` | **sim** (log contável) | Atributo `service.name` OTel | Obrigatório para compor assinatura (DN-008-1 B); ausência → `assinatura_incompleta` (spec.md FR-007/FR-008; analyze E2/D1) |
| `host` | **sim** (log contável) | Host/máquina (host.name ou equivalente OTel) | Mesma regra que `serviceName` (fonte: respostas-joao.md, DN-008-1 B; origem: ADR_013) |
| `catalogoUsado` | sim (metadado) | Se catálogo decidiu a classe | false quando nível/texto decidiram |

**Estado**: imutável após classificação do log.

---

## Assinatura (valor lógico)

Identidade do achado no produto no escopo do lote (008-18 A).

| Componente | Descrição | Fonte |
|------------|-----------|-------|
| `nivel` | Classe do log (`ClasseLog`) após FR-004 | DN-008-1 B |
| `mensagemNormalizada` | Mensagem após normalização | DN-008-1 B; algoritmo **DN-008-4 B**; marcadores **DN-008-5 A**: `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` |
| `serviceName` | `service.name` OTel | DN-008-1 B |
| `host` | Host/máquina | DN-008-1 B; origem ADR_013 |

**Fora desta lista (006)**: `traceId`, `correlationId` (spec.md D5). **Não** usar `lhamalog_batch_id`/`record_id` na assinatura (inferência respostas-etapa3.md §4 — mudam por lote/registro).

**Semântica (DN-008-4 B)**: trim; case-insensitive; substituição por categoria (número, UUID, IP, data, hex). **Literais (DN-008-5 A)**: `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>`. Implementação canônica **única**: `internal/application/assinatura/normalizar_mensagem.go` (spec **008**, plan/tasks T016); a 006 **só invoca** — **proibido** cópia em `internal/domain/` (fonte: respostas-joao.md §Etapa 3b/3c, DN-008-4 B / DN-008-5 A; analyze2 **F1**).

**Valor agregado (`assinatura` string)**: componentes concatenados com separador **`|`** (**008 assinatura A**, **DN-008-8 A**); **`|`** literal dentro de um componente MUST ser escapado como **`\|`** antes da concatenação (**DN-008-9 A**). Escape de **`\`** literal: **pendência** (sem fonte). Ordem e implementação canônica: spec **008**; esta feature referencia. (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), DN-008-8 A; §Etapa 3g (João, 16/11 13h22), DN-008-9 A; respostas-etapa3c.md §9)

---

## Achado (agrupamento no lote)

| Campo | Descrição | Fonte |
|-------|-----------|-------|
| `assinatura` | Valor composto conforme seção **Assinatura** | DN-008-1 B |
| `classe` | Igual ao `nivel` na assinatura (008-15 A) | 008-15 A |
| `logs` | IDs/apontamentos dos membros | Escopo **só do lote** (FR-009) |
| `contagem` | Quantidade de logs | N3 A — [PENDENTE ADR] |
| `primeiraOcorrencia` | Timestamp do log mais antigo no lote | N3 A, 008-14 A — [PENDENTE ADR] |
| `ultimaOcorrencia` | Timestamp do log mais recente no lote | N3 A — [PENDENTE ADR] |

**Transições**: logs classificados → particionados por **assinatura** (classe já embutida) → um achado por partição.

---

## AchadoAnalise

Estrutura enviada no `ResumoAnalise` (contrato Drive + extensões João).

### Campos (base Drive)

| Campo | Tipo | Obrigatório | Notas |
|-------|------|-------------|-------|
| `tipo` | ERROR \| WARNING \| INFO | sim | 006-6 A; sem gravidade |
| `titulo` | string | presença exigida | ADR_040 / revisao-T1-drive.md §2 — **semântica mínima de preenchimento [PENDENTE N3 A]** (analyze E3) |
| `resumo` | string | presença exigida | Idem `titulo` |
| `numeroLinha` | número | conforme Drive | |
| `camposExtraido` | objeto | conforme Drive | |
| `dadosSensiveis` | flag/lista | conforme Drive | Contagem agregada no resumo (010-17) |
| `traceId` | string | opcional | OTel — **não** na assinatura 006 |
| `correlationId` | string | opcional | Envelope v1 — não na assinatura 006 |
| `hostName` | string | quando `host` na assinatura | **Mapeamento 1:1** ao componente `host` da assinatura quando este estiver presente; nome Drive ADR_040 vs `host` na assinatura DN-008-1 B (analyze **L1**) |
| `moduloId` | string | opcional | |
| `linhas` | array string | quando aplicável | Trechos **mascarados** (ADR_040, 010-25) |
| `recorrente` | boolean | sim | ADR_063 — no escopo mínimo desta feature (**FR-009**, sem entre-lotes): **MUST ser `false`**; recorrência histórica na spec **008** (fonte: spec.md FR-007; respostas-joao.md, N2 B; analyze E3) |

### Campos (extensão N3 A — [PENDENTE ADR])

| Campo | Descrição |
|-------|-----------|
| `assinatura` | Valor serializado (**008 assinatura A** — texto concatenado com separador **`|`**, **DN-008-8 A**; formalização JSON pendente **N3 A** ADR) |
| `contagem` | Logs no achado no lote |
| `primeiraOcorrencia` | Instantâneo no lote |
| `ultimaOcorrencia` | Instantâneo no lote |
| `versaoRegras` | Versão **única** regras+catálogo+base (006-3 C) |
| `catalogoEntrou` | boolean + referência de versão do catálogo quando 006-7 D fechar |

**Relacionamentos**: N `LogClassificado` → 1 `Achado` → 1 `AchadoAnalise`; no Módulo, 1 **Ocorrência** por achado (T2 C).

---

## ResumoAnalise

| Campo | Obrigatório | Descrição |
|-------|-------------|-----------|
| `batchId` | sim | UUID do lote ingested (`Referencia.BatchID`) |
| `execId` | opcional | Varredura `V<uuid>` / `M<uuid>` — **ausente** no pipeline só `log.ingested` v1 (DN-006-1b B) |
| `tenantId` | sim | Metadata ingested |
| `productId` | sim | Metadata ingested |
| `contadores` | sim | Ver **ContadoresResumo** abaixo |
| `achados` | sim | Lista de `AchadoAnalise` |
| `versaoRegras` | recomendado | 006-3 C quando não só por achado |
| `identificacaoRegras` | sim | Catálogo entrou? Versão única (FR-007) |

**Chave de idempotência lógica**:

- Com `execId`: `(batchId, execId)` — N4 A, DN-006-1.
- Sem `execId` (fluxo atual): **`batchId` somente** — DN-006-1b B.

Módulo ignora duplicata (FA-02, 006-9 A no módulo).

**Pré-condição de handoff (DN-008-7 B)**: antes do envio HTTP, o Finder realiza **uma consulta ao Módulo de Ocorrência por achado** (N2 B na spec **008**). Contrato/rota da consulta: **DN-009-2**. Envio: **exatamente um** `ResumoAnalise` por lote (fonte: respostas-joao.md §Etapa 3c DN-008-7 B; respostas-etapa3c.md §6; origem: ADR_062).

### ContadoresResumo (ADR_102 revisado)

| Contador | Descrição | Fonte |
|----------|-----------|-------|
| `linhasProcessadas` | Linhas analisadas no lote | ADR_102 |
| `erros` | Achados/logs ERROR conforme regra de contagem do lote | ADR_102 |
| `warnings` | Idem WARNING | ADR_102 |
| `info` | Contador de **INFO** (substitui «spam ignorado») | respostas-joao.md Etapa 3; ADR_170 |
| `dadosSensiveisMascarados` | Trechos mascarados | ADR_102 |

Nome JSON exato até publicação ADR_102 — alvo `info`, **não** `spamIgnorado` (fonte: spec.md etapa 3b).

---

## Ocorrência (Módulo de Ocorrência — referência)

Não persistida pelo Finder.

| Campo | Notas |
|-------|-------|
| Chave unicidade | assinatura + classe + produto + tenant (008-11 rodada 16/11) |
| `tipo` / classe | ERROR, WARNING ou INFO (T7 A) |
| `conteudo` | Exemplos mascarados de `AchadoAnalise.linhas` (ADR_040) |
| `gravidade` | Calculada no **Módulo** — fora do Finder (ADR_062) |

---

## Identificação das regras (006-3 C)

| Campo | Descrição |
|-------|-----------|
| `versaoUnica` | String versionada cobrindo regras de nível, catálogo (quando existir) e base 007 |
| `somenteNivel` | true quando catálogo ausente ou não usado (FR-007) |

---

## Avaliação offline (base spec 007)

| Regra | Descrição | Fonte |
|-------|-----------|-------|
| Unidade de gold | **Por apontamento** | **DN-007-4 A** (spec 007) |
| Classe prevista do apontamento | **Mais grave** entre as classes previstas dos logs do apontamento (**ERROR > WARNING > INFO**) | **007 H2 D** (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40); origem: ADR_170, respostas-etapa3f.md opção D) |
| Métricas e gate | Precisão, recall e F1 **por classe**; **≥ 95%** em ERROR, WARNING e INFO | **DN-007-2 B**, **DN-006-2 B** (SC-007, FR-010) |

---

## Falhas e motivos

| Motivo | Efeito | Fonte |
|--------|--------|-------|
| `log_ausente` | Sem achado no resumo de sucesso | spec.md Edge Cases |
| `armazenamento_inacessivel` | Sem `ResumoAnalise` de sucesso | FR-008 |
| `lote_incompleto` | Classificação parcial OK; **não** sucesso de lote (T4 A) | respostas-joao.md, T4 A |
| `envio_modulo_falhou` | Varredura não concluída; reenvio mesma chave — implementação HTTP **DN-009-2** | revisao-T1-drive.md §2 |
| `identidade_incompleta` | tenant/produto vazio no achado — rejeitar | 007-8 A |
| `produto_ausente_ou_inconsistente` | `productId` ausente no ingested v1 ou **divergente** entre metadata validada (token) e corpo do evento — **falha crítica permanente**, **sem retry** (ADR_049); MUST NOT sucesso de lote | respostas-joao.md, 008-16 A; contrato-log-ingested-v1.txt §4; analyze2 **E2** |
| `assinatura_incompleta` | Log contável sem `service.name` e/ou **host** exigidos por DN-008-1 B — MUST NOT agregar em achado; lote MUST NOT contar como varredura concluída com sucesso (paralelo a `identidade_incompleta`, 007-8 A) | respostas-joao.md, DN-008-1 B; spec.md FR-007/FR-008; analyze **E2**, **D1** |

### Rastreio — consulta OpenObserve (spec 008)

Após **3 tentativas** falhas na consulta OpenObserve (recorrência entre lotes), o Finder avisa o Core pela API interna (**`batchId`** + **motivo**); **sem** DLQ no `log.ingested` e **sem** alerta WebSocket pós-DLQ (**DLQ B**; **DN-008-3** / **DN-008-6** rev. 3h — **substituem** alcance 3b–3e com DLQ). Esta feature não define rota/payload/evento (**ADR_051 [PENDENTE]**). (fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22), DLQ B; origem: respostas-etapa3g.md §1b, contrato-log-ingested-v1.txt §8)

---

## Relação com tipos existentes no repositório

| Tipo Go | Uso nesta feature |
|---------|-------------------|
| `domain.Referencia` | Origem de batchId, tenant, produto, janela temporal |
| `domain.LogPronto` / `Registro` | Entrada à classificação (campos OTel até cadeia 003–005) |
| `domain.Codigo` | Apontamento no OpenObserve |
| `domain.Chave` | Idempotência **ingested** (≠ chave ResumoAnalise) |