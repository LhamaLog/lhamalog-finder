# Quickstart: Validação do agrupamento em ocorrências

**Feature**: `008-agrupar-logs-ocorrencia` | **Date**: 2026-11-16 (etapas 3b–3h)

Guia de **validação** (não implementação). Integração real com Módulo Ocorrência exige **DN-009-2**. Equivalência por mensagem segue **DN-008-4 B** com marcadores **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`** (**DN-008-5 A**); casos extras em **010-27 A**.

## Pré-requisitos

- Decisões em [spec.md](./spec.md) (T1 rev., N2 B, 008-11 rev., T7 A, T4 A, DN-008-1 B, DN-008-3, **DN-008-4 B**)
- Classificação (spec **006**)
- Contratos: [contracts/](./contracts/)
- Modelo: [data-model.md](./data-model.md)
- Infra opcional: `make run` (README) — Redis, OpenObserve, Finder

## 1. Assinatura (componentes + normalização DN-008-4 B)

**Objetivo**: FR-002, FR-002b / SC-001.

1. Validar **componentes fixos**: `nivel`, `service.name`, `host` idênticos → mesma base de assinatura (**DN-008-1 B**).
2. **`traceId`/`correlationId` diferentes** com demais componentes iguais → **mesma** assinatura (**008-2 rev.**).
3. Classes diferentes → chaves de ocorrência distintas mesmo com assinatura igual (**008-15 A**).
4. **Mensagem**: aparar espaços; ignorar maiúsculas/minúsculas; substituir número, UUID, IP, data e hex pelos marcadores **DN-008-5 A** — matriz adicional em **010-27 A**.
5. **Valor** da assinatura: texto concatenado (**008-assinatura A**) com separador **`|`** entre **nível**, **mensagem normalizada**, **`service.name`** e **host**, nessa ordem (**DN-008-8 A**). **`|`** literal nos componentes: escape **`\|`** (**DN-008-9 A**).

Referência: [assinatura-log.md](./contracts/assinatura-log.md).

## 2. Agrupamento no lote

**Objetivo**: User Story 1 — [agregar-achado-lote.md](./contracts/agregar-achado-lote.md).

1. Lote completo com N logs equivalentes (mesma assinatura e classe) e repetições no mesmo lote.
2. **Esperado**: um achado; `contagemLote` = N (**008-12**, **008-13**).
3. **Esperado**: `primeiraNoLote` / `ultimaNoLote` = min/max timestamp do log (**008-14**).

## 3. Lote incompleto (T4 A)

1. Simular recuperação com contagem distinta (spec **005**).
2. **Esperado**: **sem** agrupamento e **sem** `ResumoAnalise` de sucesso para esse lote.

## 4. Chave de unicidade

**Objetivo**: FR-005 / SC-001.

1. Mesma assinatura, **produtos**, **tenants** ou **classes** diferentes.
2. **Esperado**: achados distintos; no mock do Módulo, chaves distintas (**008-11 rev.**).

## 5. Fusão entre lotes (N2 B)

**Objetivo**: User Story 1 cenário 1.

1. **Mock** do Módulo (sem HTTP real — **DN-009-2**) retorna ocorrência existente com contador C e extremos — **uma chamada mock por achado** (**DN-008-7 B**).
2. Segundo lote com M logs equivalentes.
3. **Esperado** no `ResumoAnalise` em memória: `contagem` = C + M ([data-model.md](./data-model.md)); **um** envelope por lote ([consulta-ocorrencias-modulo.md](./contracts/consulta-ocorrencias-modulo.md), [resumo-analise-agrupado.md](./contracts/resumo-analise-agrupado.md)).

## 6. Reabertura (Módulo — teste de contrato)

**Objetivo**: SC-003 / **008-4 B**.

1. Mock com ocorrência **FECHADA** para a chave.
2. Finder envia achado equivalente.
3. **Esperado no Módulo** (fora do Finder): status **ABERTA**; Finder **não** altera status (**008-10 A**).

## 7. Recorrência OpenObserve

**Objetivo**: FR-003 / **008-5**, **008-6**, **SC-007**.

1. Adaptador quando existir — query **NEEDS CLARIFICATION** (`openobserve/doc.go`).
2. Cenário com ocorrências em ≥3 lotes distintos, mesma máquina, mesma assinatura.
3. **Esperado**: `recorrente: true` quando ADR_013 satisfeito na janela de **500 lotes**.
4. Repetição só no mesmo lote: contador sobe; recorrência **não** conta lotes distintos (**008-13** × **008-5**).

**Falha de consulta (DN-008-3, rev. etapa 3h)**:

1. Simular indisponibilidade temporária.
2. **Esperado**: até **3 tentativas** na consulta; **sem DLQ** no `log.ingested`; **não** registrar `recorrente=false` por degradação.
3. **Esperado** após 3 falhas: Finder avisa o **Core** pela **API interna** com **batchId** e **motivo** (**DN-008-6** rev. 3h; contrato **ADR_051 [PENDENTE]** — MUST NOT inventar rota/payload/evento).

Referência: [consulta-recorrencia-openobserve.md](./contracts/consulta-recorrencia-openobserve.md).

## 8. Handoff ampliado

**Objetivo**: [resumo-analise-agrupado.md](./contracts/resumo-analise-agrupado.md).

1. Montar **um** payload `ResumoAnalise` **por lote** em memória (**DN-008-7**); HTTP real **BLOQUEADO** (**DN-009-2**).
2. Chave `batchId` + `execId` (**N4 A**; **DN-006-1**).
3. **Esperado**: Finder **não** grava ocorrência localmente (FR-001); consulta prévia **por achado** já executada na §5 (**DN-008-7 B**).

## 9. INFO inclusive (T7 A)

1. Lote com logs INFO agrupáveis.
2. **Esperado**: achados `tipo` INFO no resumo; aceite Módulo fora deste repo.

## 10. Legado (DN-008-2)

1. Pipeline documentado (contrato v1 §11, T1 rev.).
2. **Esperado**: **não** exigir migração de ocorrências legadas; ressalva ambiente de teste (spec.md D2b).

## 11. O que não validar nesta feature

- Escape de **barra invertida** literal nos componentes do `valor` (pendente; **DN-008-9 A** fechou só **`|`** → **`\|`**)
- Contrato Finder→API Core após falha OpenObserve (rota/payload/evento — **ADR_051 [PENDENTE]**, sem fonte neste repo)
- API HTTP Módulo consulta/envio (**DN-009-2**)
- ADR N3 (nomes JSON finais)
- Notificação ADR_035 (**008-17 C**)
- Meta **≥ 95% precisão por classe** e base rotulada JSONL em `datasets/labeled-logs/` (**DN-007-2 B**, **007 formato A**, **007 caminho C** — spec **007**)

## Referências rápidas

| ID | Cenário |
|----|---------|
| DN-008-1 B | 1 |
| DN-008-4 B | 1 |
| DN-008-5 A | 1 |
| DN-008-8 A | 1 |
| DN-008-9 | 1, 11 |
| 008-assinatura A | 1 |
| DN-009-2 | 5, 8, 11 |
| DN-008-3 / DN-008-6 | 7 |
| DN-008-7 / DN-008-7 B | 5, 8 |
| DN-008-2 | 10 |
| 008-11 rev. | 4, 5 |
| 008-12 / 008-13 / 008-14 | 2, 5 |
| T4 A | 3 |
| N2 B | 5 |
| 008-4 B / 008-10 A | 6 |
| T7 A | 9 |

## Verificação desta execução

§1–§4: assinatura/partição/gate/chave verificados nos pacotes assinatura e agrupar. Testes canônicos cobrem cinco literais, equivalência e escape de |; matriz formal P-010-27 A e barra invertida literal continuam pendentes. §5–§7: fusão, estado FECHADA preservado e recorrência com 3 lotes distintos/janela de 500 verificados em portas; reabertura real e query histórica OO não executadas. §8: um resumo lógico por lote; HTTP/N3 A bloqueados. §9: INFO participa da mesma preparação e envio; sem traceId/correlationId na assinatura (008-2 A revogada, D1/D2). §10: sem migração de produção (DN-008-2). Consulta por achado e um envio por lote; execId opcional vem da Varredura, jamais criado neste caminho.

Comandos realmente observados e saídas completas: `evidencias/aceite-final-unit.txt`, `aceite-final-race-integracoes.txt` e `aceite-final-vet.txt` (exit 0). Executados em composição temporária dos patches das specs, com manifesto SHA-256. Para reproduzir: compor os pré-requisitos conforme relatório central, usar Go 1.27.2 e fornecer explicitamente os parâmetros de teste de Redis/OpenObserve; credenciais não estão nos relatórios. `make test` falhou no contêiner por certificado ao baixar módulos; o comando Go local com race passou. Não executar `go test ./...` isoladamente numa branch consumidora e interpretar falta de dependências como falha da composição.
