# Quickstart: Validar handoff e pipeline LLT-201 (009)

Guia de **validação** — não inclui implementação Go. Escopo Finder: [contracts/resumo-analise-handoff.md](./contracts/resumo-analise-handoff.md). Pipeline produto: [contracts/pipeline-aviso-referencia.md](./contracts/pipeline-aviso-referencia.md).

## Pré-requisitos

| Item | Fonte |
|------|-------|
| Spec **006** (classificação ERROR/WARNING/INFO) | FR-F04, T7 A |
| Spec **008** (agrupamento, N2 B) para cenários atualizada/recorrente | 009-22 A, SC-002 |
| Contrato **log.ingested v1** (`schemaVersion` `"1"`) | T5 D; `internal/domain/envelope.go` |
| Ambiente com Módulo Ocorrência, Core e Notificação para US2/US3 produto | spec.md US2 |

**Pendência**: URL, JWT e respostas do Módulo Ocorrência para **POST ResumoAnalise** e **consulta por achado** — **DN-009-2** (bloqueia testes de integração do adaptador HTTP Finder; fonte: respostas-joao.md §Etapa 3). Granularidade do handoff (**DN-008-7**) e da consulta (**DN-008-7 B**) já fechadas.

**Precisão e base 007**: ver [spec.md](./spec.md) §Decisões tomadas (**DN-007-2 B**, 007 formato **A**, caminho **C**); este quickstart não repete o gate (**corrige analyze D2**).

---

## Cenário 1 — Handoff feliz (US1, SC-F001)

**Objetivo**: Lote com achados produz um handoff idempotente e confirmação antes de concluir varredura.

1. Ingerir lote válido v1 com achados ERROR, WARNING e **INFO** (T7 A).
2. Completar análise do lote no Finder (classificação 006).
3. Verificar **exatamente uma** requisição HTTP+JWT de handoff ao Módulo por lote (**DN-008-7**, **FR-F01**) com `ResumoAnalise` contendo **`batchId`** (**DN-006-1b B** — sem `execId` no fluxo v1 atual), contadores e `achados[]` com `tipo` válidos.
4. Simular resposta 2xx com confirmação do Módulo (formato exato: **DN-009-2**).
5. **Esperado**: varredura só conclui **após** passo 4 (FR-F05, SC-F001); **nenhuma** publicação Core/Notificação pelo Finder (FR-F02).

**Variante**: quando `execId` existir, chave **`batchId` + `execId`** (N4 A, DN-006-1).

---

## Cenário 2 — Proibição de publicar evento (FR-F02)

**Objetivo**: Garantir que o Finder não emite aviso LLT-201.

1. Repetir cenário 1 com tráfego instrumentado em Core e Redis `notificacao-eventos`.
2. **Esperado**: **0** chamadas do Finder ao Core para aviso; **0** mensagens LLT-201 originadas no processo Finder.

---

## Cenário 3 — Reenvio idempotente (US3, SC-F002)

**Objetivo**: Mesma chave (`batchId` ou `batchId + execId`) sem duplicar ocorrências no Módulo.

1. Concluir handoff com confirmação (cenário 1).
2. Reprocessar ou reenviar o mesmo `ResumoAnalise` com **mesma** chave (falha simulada antes da confirmação — FA-02).
3. **Esperado**: Módulo ignora duplicata; contagem de ocorrências criadas não aumenta (FR-F06, SC-F002).

---

## Cenário 4 — Assinatura e normalização (FR-F11, DN-008-4 B)

**Objetivo**: Achados usam mensagem normalizada canônica da spec **008** antes do handoff.

1. Montar achados com variações de caixa, espaços e literais substituíveis (número, UUID, IP, data, hex).
2. **Esperado**: saída do normalizador **DN-008-4 B** com marcadores **DN-008-5 A** (`<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>`) idêntica à regra da spec **008**; quando o campo `assinatura` existir (**N3 A** / ADR), MUST seguir **008 assinatura A** (concatenação **DN-008-1 B** + separador fixo **`|`** — **DN-008-8 A**); **`|`** literal dentro de um componente MUST aparecer como **`\|`** (**DN-008-9 A**; fonte: respostas-joao.md §Etapa 3g (João, 16/11 13h22); origem: DN-008-8 A).
3. **Não** inventar literais de marcador, separador divergente de **`|`** nem regra de escape de **`\`** (sem fonte — ver pendência em spec.md) nesta feature (**FR-F11**).

---

## Cenário 5 — Classes inválidas (SC-F003)

**Objetivo**: 0% de achados fora ERROR/WARNING/INFO.

1. Tentar incluir classe fora da taxonomia (ex. SPAM) no payload de teste.
2. **Esperado**: rejeição na montagem do handoff ou no contrato de classificação 006 — não enviar ao Módulo (FR-F08, constituição III).

---

## Cenário 6 — Falha na consulta OpenObserve (FR-F07, DN-008-3)

**Objetivo**: Política de recorrência/consulta antes de qualquer handoff.

1. Simular falha transitória na **consulta ao OpenObserve** durante a análise do lote (não no HTTP ao Módulo).
2. **Esperado**: até **3 tentativas**; **sem DLQ** no **log.ingested v1** (**DLQ B**); **sem** publicação em stream DLQ nem alerta via consumidor DLQ / **WebSocket** pós-DLQ; após esgotar, o **Finder** chama a **API interna do Core** com **`batchId`** e **motivo** (**DN-008-6** rev.); contrato rota/payload/evento **[PENDENTE ADR_051]**; **sem** `ResumoAnalise` de sucesso (fonte: respostas-joao.md §Etapa 3g, DLQ B; respostas-etapa3g.md §1; rastreio em [plan.md](./plan.md) §Decisões tomadas etapa 3h).

---

## Cenário 7 — Lote incompleto (T4 A)

**Objetivo**: Sem handoff de sucesso.

1. Usar lote marcado incompleto conforme spec **005**.
2. **Esperado**: **sem** `ResumoAnalise` de sucesso nem conclusão de varredura como sucesso (T4 A).

---

## Cenário 8 — Pipeline integrado: ocorrência **nova** (US2, SC-001) — ambiente produto

**Objetivo**: Critério (1) LLT-201 após handoff real.

1. Produzir `ResumoAnalise` real (cenário 1).
2. No Módulo: criar ocorrência; enviar aviso HTTP ao Core.
3. Core publica no stream `notificacao-eventos`.
4. **Esperado**: payload HTTP lógico com **tenant**, **produto**, **id** ocorrência, **tipo** (classe), natureza **nova**, **gravidade** (**009-5 revisto**, FR-003).

---

## Cenário 9 — Consulta prévia por achado (FR-F09, DN-008-7 B)

**Objetivo**: **N2 B** antes do handoff único por lote.

1. Lote com múltiplos achados elegíveis a consulta (spec **008** disponível ou mock).
2. **Esperado**: **N** consultas HTTP ao Módulo (**uma por achado**), **antes** do único **POST** `ResumoAnalise` (**DN-008-7**); formato das consultas: **DN-009-2** (bloqueia teste real até contrato).

---

## Cenário 10 — Estados **atualizada** / **recorrente** (SC-002) — bloqueado até 008

**Objetivo**: Distinguir estados no aviso (FR-006).

1. Pré-requisito: spec **008** implementada e agrupamento entre lotes (N2 B).
2. Segundo lote com mesma assinatura+classe+produto+tenant (008-11 rev. João 16/11).
3. **Esperado**: aviso com natureza **atualizada** ou **recorrente** conforme regras 008; reabertura → **atualizada** (009-14 A).
4. **Nota**: aceite formal **SC-002** = **009-22 A** — não fechar aceite 009 até 008.

---

## Cenário 11 — Duplicata no Core (SC-003)

**Objetivo**: Reentrega HTTP não gera 2º alerta no stream.

1. Processar aviso com id ocorrência + estado **nova**.
2. Reenviar **mesmo** aviso HTTP.
3. **Esperado**: Core **não** publica segundo alerta em `notificacao-eventos` (009-11 B, 009-9 A).

---

## Cenário 12 — Duplicata na stream (009-10 B)

**Objetivo**: Consumidor Notificação idempotente por `eventId`.

1. Redeliver mesma mensagem ADR_146 no stream.
2. **Esperado**: sem segundo efeito de notificação indevido no consumidor (guia §8).

---

## Cenário 13 — E-mail com resumo mascarado (DN-009-1 A)

**Objetivo**: LLT-218 sem texto no aviso HTTP.

1. Após Notificação consumir evento com id ocorrência.
2. **Esperado**: Notificação lê Ocorrência no MongoDB compartilhado (**DN-009-1 A**, ADR_054) e monta e-mail com trechos de `conteudo` mascarado; templates NOT-03 podem permanecer [PENDENTE].

---

## Evidências a registrar

- Comando ou trace da chamada HTTP handoff (URL quando **DN-009-2** fechar, status, correlation).
- Corpo `ResumoAnalise` (sem segredos; linhas mascaradas).
- Confirmação do Módulo.
- Para integrado: payload HTTP Core e entrada no stream (referência, não código Finder).

## Referências

- [data-model.md](./data-model.md)
- [research.md](./research.md)
- [plan.md](./plan.md) — § API do Módulo Ocorrência, § Pendências