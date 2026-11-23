# Quickstart: Validação do modelo de classificação de logs

**Feature**: `006-modelo-classificacao-logs` | **Date**: 2026-11-16 (etapa 3h)

Guia de **validação** (não implementação). Prova ponta a ponta com Módulo de Ocorrência real exige contrato **DN-009-2** ainda pendente.

## Pré-requisitos

- Go 1.27 via Docker, conforme [README.md](../../../README.md)
- Decisões etapa 3h (DLQ B, rev. **DN-008-3** / **DN-008-6**, **DN-008-9 A**) e anteriores na spec [spec.md](../spec.md)
- Contratos: [classificar-log.md](./contracts/classificar-log.md), [classificar-apontamento.md](./contracts/classificar-apontamento.md), [resumo-analise-handoff.md](./contracts/resumo-analise-handoff.md)
- Modelo: [data-model.md](../data-model.md); plano: [plan.md](../plan.md)

## 1. Regra de classe (unidade)

**Objetivo**: SC-001 / User Story 1 — sem LLM.

1. Invocar `Classificar log` conforme [classificar-log.md](./contracts/classificar-log.md).
2. Executar tabela de cenários da spec (nível + mensagem + catálogo).
3. **Esperado**: somente ERROR, WARNING ou INFO; `errorless` não vira ERROR; `traceback` sozinho não vira ERROR.

**Comando** (quando existirem testes): `make test` ou `go test ./internal/domain/...` — alvo em `tasks.md`.

## 2. Agrupamento em achados (lote)

**Objetivo**: User Story 2 — um AchadoAnalise por assinatura (classe no componente `nivel`).

1. Assinatura: **nível + mensagem normalizada + `service.name` + host** (DN-008-1 B).
2. Entrada: três logs mesma assinatura ERROR + um WARNING com assinatura distinta.
3. **Esperado**: dois achados; contadores `erros`, `warnings`, `info` coerentes (ADR_102 revisado).

**DN-008-4 B** / **DN-008-5 A**: normalização via **`internal/application/assinatura/normalizar_mensagem.go`** (spec **008**). Assertar marcadores `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>` (fonte: respostas-joao.md §Etapa 3c). Valor da assinatura: texto concatenado (**008 assinatura A**), separador **`|`** (**DN-008-8 A**, spec **008**); **`|`** literal em componente → **`\|`** (**DN-008-9 A** — fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22)); comparar agrupamento pelo valor serializado.

**010-27 A**: casos de borda adicionais do normalizador (spec 010).

## 3. INFO no resumo

**Objetivo**: T7 A, 006-11 A.

1. Dois logs INFO, assinaturas distintas.
2. **Esperado**: dois achados `tipo` INFO; contador `info` ≥ 2; **nenhuma** gravação de ocorrência no Finder.

## 4. Lote incompleto (T4 A)

1. Simular lote com contagem esperada maior que logs recuperados.
2. **Esperado**: motivo `lote_incompleto`; **não** sucesso de lote completo.

## 4c. Produto ausente ou inconsistente (008-16 A)

**Objetivo**: FR-008 / SC-003 — falha crítica permanente, sem retry (analyze2 **E2**).

1. Evento `log.ingested` v1 com `productId` ausente **ou** metadata validada (token) ≠ `productId` no corpo repetido para agrupamento (contrato §4).
2. **Esperado**: motivo `produto_ausente_ou_inconsistente`; **não** sucesso de lote; **sem** retry de processamento do lote (ADR_049). Distinto de `identidade_incompleta` (007-8 A).

(fonte: respostas-joao.md, 008-16 A; contrato-log-ingested-v1.txt §4)

## 4b. Assinatura incompleta (DN-008-1 B)

**Objetivo**: FR-008 / `assinatura_incompleta` (analyze **E2**, **G1** cobertura de teste fica em `tasks.md`).

1. Log contável com classe definida, **sem** `service.name` ou **sem** host.
2. **Esperado**: motivo `assinatura_incompleta`; log **não** entra em achado; **não** sucesso de lote completo (fonte: spec.md FR-007/FR-008; respostas-joao.md, DN-008-1 B).

## 4d. Consulta ao Módulo por achado (DN-008-7 B)

**Objetivo**: FR-007 / User Story 2 — N2 B + consulta antes do envio.

1. Lote com **N** achados distintos após classificação e agrupamento.
2. **Esperado**: **N** consultas ao Módulo de Ocorrência (**uma por achado**), depois **um único** `ResumoAnalise` na entrega HTTP (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B; respostas-etapa3c.md §6).
3. Implementação/rota real: **BLOQUEADO (DN-009-2)** — validar com adaptador fake que conte chamadas por achado.

## 5. Handoff HTTP (integração)

**Objetivo**: SC-002 — **aceite parcial** vs **E2E real** (analyze **E1**; [plan.md](../plan.md) § SC-002 e aceite parcial).

| Nível | Passos | Critério de sucesso desta rodada |
|-------|--------|----------------------------------|
| **Parcial** | Validar `ResumoAnalise` **em memória** contra [resumo-analise-handoff.md](./contracts/resumo-analise-handoff.md); opcional adaptador fake (T024/T032) | Payload + idempotência lógica corretos; **não** equivale a SC-002 completo |
| **E2E real** | **BLOQUEADO (DN-009-2)** — JWT, rota, códigos HTTP, timeout | Confirmação explícita do Módulo antes de concluir varredura |

1. Chave no fluxo atual: **somente `batchId`** (DN-006-1b B); com `execId`, `batchId`+`execId` (DN-006-1).
2. Após contrato DN-009-2: ambiente integrado; 2xx + confirmação; reenvio mesma chave (FA-02).
3. **Esperado**: zero publicação de evento de ocorrência pelo Finder.

## 6. Envelope ingested v1

1. `schemaVersion` `"1"` e `batchId` válido (`internal/domain/envelope.go`).
2. Rejeitar `schemaVersion` ≠ `"1"` (T5 D).

## 7. Métricas de aceite (offline)

**Objetivo**: SC-007 / 006-2 B / DN-006-2 B (etapa 3e).

1. Com base rotulada da spec **007** (**JSONL**, caminho documentado `datasets/labeled-logs/` na raiz — **007 formato A**, **007 caminho C**), calcular precisão, recall e F1 **por classe**; gold **por apontamento**; classe prevista do apontamento = **mais grave** entre os logs (**007 H2 D**, ADR_170).
2. **Gate**: em **cada** classe ERROR, WARNING, INFO — precisão **≥ 95%** (**DN-007-2 B**), recall **≥ 95%** e F1 **≥ 95%** (**DN-006-2 B** — fonte: respostas-joao.md §Etapa 3c, João 16/11 11h e 11h01).

## 7b. Desempenho classificar + montar resumo (SC-002, DN-006-3 B)

**Objetivo**: meta **≤ 2 s por lote** para classificar logs do lote + montar `ResumoAnalise` (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-006-3 B).

1. Ambiente: **VM-A real** com **todos os serviços** rodando; lote de **1.000 registros** (**DN-006-3-ref** — fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40)).
2. Medir tempo do ciclo **sem** incluir consultas ao Módulo (**DN-008-7 B**) nem envio HTTP (**DN-009-2**) (spec.md SC-002).
3. **Esperado**: ≤ 2 s por lote (**DN-006-3 B**).

## 8. O que não validar nesta feature

- Persistência MongoDB de ocorrência pelo Finder — o Finder envia `ResumoAnalise`; o **Módulo de Ocorrência** cria/persiste (T1 REVISTO; analyze-006-etapa3g **F1**)
- Gravidade, sugestão, notificação, evento Core
- Política pós-3 falhas na consulta OpenObserve (**DN-008-3** / **DN-008-6** rev. 3h, **DLQ B**) — runtime na spec **008**; contrato Finder→Core **ADR_051 [PENDENTE]** (fonte: respostas-joao.md §Etapa 3g, DLQ B)
- Sintaxe do catálogo (006-7 D), ADR N3, API Módulo (DN-009-2)

## Referências rápidas

| ID | Uso na validação |
|----|------------------|
| T4 A | Cenário 4 |
| T5 D | Cenário 6 |
| T7 A / 006-11 A | Cenário 3 |
| DN-008-1 B | Cenário 2 |
| DN-006-1b B | Cenário 5 (chave) |
| N1 A / N4 A | Cenário 5 (pós DN-009-2) |
| 006-2 B / DN-007-2 B / DN-006-2 B / 007 H2 D | Cenário 7 |
| DN-006-3 B / DN-006-3-ref | Cenário 7b |
| DN-008-7 B | Cenário 4d |
| DN-009-2 | Cenário 5 bloqueado |
| DN-008-4 B / DN-008-5 A | Cenário 2 |
| 008 assinatura A | Cenário 2 |
| DN-008-8 A / DN-008-9 A | Cenário 2 |
| 008-16 A | Cenário 4c |
| DLQ B / DN-008-3 rev. 3h / DN-008-6 rev. 3h | §8 (rastreio; spec **008**) |

(fonte: respostas-joao.md; detalhes em [research.md](../research.md))

## Verificação desta execução

§1–§4d: classificacao_test.go e analisar_lote_test.go cobrem classes, partições (3 ERROR equivalentes + 2 INFO distintos), lote/identidade/produto/assinatura inválidos, N consultas/um envio, confirmação, catálogo decisório e snapshot de versão. §6: envelope v1 conferido nos testes existentes/004. §5: HTTP bloqueado DN-009-2. §7: procedimento offline em rastreabilidade.md/007; base real indisponível. §7b: `TestAceiteDesempenhoVMA` é opt-in `FINDER_ACEITE_VMA=1`, exige atestar VM-A real com todos os serviços e 1.000 registros; perfil-vma.txt mostra SKIP, não sucesso de ≤2s. §8: três falhas OO → porta Core batchId/motivo, sem DLQ/WebSocket; cliente ADR_051 pendente.

Comandos realmente observados e saídas completas: `evidencias/aceite-final-unit.txt`, `aceite-final-race-integracoes.txt` e `aceite-final-vet.txt` (exit 0). Executados em composição temporária dos patches das specs, com manifesto SHA-256. Para reproduzir: compor os pré-requisitos conforme relatório central, usar Go 1.27.2 e fornecer explicitamente os parâmetros de teste de Redis/OpenObserve; credenciais não estão nos relatórios. `make test` falhou no contêiner por certificado ao baixar módulos; o comando Go local com race passou. Não executar `go test ./...` isoladamente numa branch consumidora e interpretar falta de dependências como falha da composição.
