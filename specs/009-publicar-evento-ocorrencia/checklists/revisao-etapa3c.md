# Checklist: Revisão etapa 3c — qualidade dos requisitos (feature 009)

**Purpose**: Revisar completude, clareza, consistência spec/plan/tasks, critérios mensuráveis, exceção/recuperação, dependências e rastreabilidade antes da implementação (público: revisor PR/gestor; profundidade: padrão).

**Created**: 2026-11-16

**Feature**: [spec.md](../spec.md) | [plan.md](../plan.md) | [tasks.md](../tasks.md)

**Note**: Checklist gerado por `/speckit-checklist` (modo não interativo). Itens avaliados contra os artefatos em 2026-11-16 etapa 3c.

**Review Ownership**: `[x]` neste arquivo indica que o critério de qualidade do requisito foi atendido na documentação atual; não significa implementação concluída.

## Requirement Completeness

- [x] CHK001 — O fluxo ponta a ponta (Finder → ResumoAnalise → Módulo → Core → Notificação) e a exclusão de publicação de evento pelo Finder estão documentados com requisitos dedicados? [Completeness, Spec §Responsabilidades, FR-F01/F02, FR-001]
- [x] CHK002 — Os requisitos do Finder para handoff (HTTP+JWT, chave idempotência, confirmação antes de concluir varredura, achados INFO, proibição Core/stream) estão listados sem lacunas funcionais explícitas além das pendências registradas? [Completeness, Spec FR-F01–FR-F11]
- [ ] CHK003 — O contrato HTTP Finder → Módulo Ocorrência (rota, claims JWT, códigos de resposta, timeout, transporte da chave) está especificado de forma implementável, não apenas como pendência? [Gap, Spec §DN-009-2; Plan §API do Módulo] (pendente: DN-009-2)
- [ ] CHK004 — Os campos ampliados de AchadoAnalise (assinatura, contagem, primeira/última vez, versão das regras) possuem definição de formato fechada nos artefatos da 009? [Gap, Spec N3 A; data-model; tasks T016] (pendente: N3 A)
- [x] CHK005 — O pipeline downstream (Módulo → Core HTTP, Core → Redis, Notificação lê MongoDB) está coberto por requisitos de referência (FR-001–FR-012) mesmo fora do escopo de código Finder? [Completeness, Spec User Story 2]
- [x] CHK006 — Dependências explícitas para specs 006, 008 e 010 e para decisões transversais (T1–T7, N1–N4) estão inventariadas em spec, plan e/ou tasks? [Completeness, Spec §Fora de escopo; Plan §Pendências; tasks §Pendências]
- [ ] CHK007 — A sintaxe e o dono do catálogo de classificação (decisão 006-7 D) estão rastreados nesta feature onde afetam versão de regras no handoff? [Gap, Spec 006-3/N3 A; tasks T021] (pendente: 006-7 D — LLT-51 ou ADR nova)
- [x] CHK008 — Literais exatos dos marcadores de normalização (DN-008-5) estão explicitamente fora desta spec com referência ao dono (spec 008)? [Completeness, Spec §DN-008-4, FR-F11]

## Requirement Clarity

- [x] CHK009 — A chave de idempotência do ResumoAnalise (`batchId` vs `batchId` + `execId` opcional) está definida sem ambiguidade para o fluxo log.ingested v1? [Clarity, Spec FR-F03, DN-006-1b B]
- [x] CHK010 — O conteúdo mínimo do aviso HTTP Módulo → Core (id, tipo, gravidade, tenant, produto, natureza da mudança) está distinguível do campo `ocorrencia` no envelope Redis (id opaco)? [Clarity, Spec 009-5 revisto, FR-003/FR-004]
- [x] CHK011 — O algoritmo canônico de normalização de mensagem na assinatura (DN-008-4 B) está descrito com regras verificáveis (trim, case, categorias de marcadores) e dono (spec 008)? [Clarity, Spec FR-F11, Clarifications etapa 3c]
- [x] CHK012 — Os papéis de gravidade (só Módulo) e tipo/classe (Finder em AchadoAnalise) estão separados de forma não ambígua? [Clarity, Spec §Edge Cases, ADR_062 citado]
- [x] CHK013 — O significado de “avisar Ocorrências” no LLT-201 está traduzido para handoff ResumoAnalise sem confundir com evento de notificação? [Clarity, Spec 009-2, D1]

## Requirement Consistency

- [x] CHK014 — Spec, plan e tasks concordam que o Finder não publica evento LLT-201 nem grava ocorrência de negócio (T1 revisto)? [Consistency, Spec FR-F02; Plan Summary; tasks Notes]
- [x] CHK015 — A dependência 009-22 A (009 espera 008 para atualizada/recorrente) está alinhada entre SC-002/SC-005, plan, tasks bloqueadas e status do spec? [Consistency, Spec SC-002; Plan; tasks T026/T031/T040]
- [x] CHK016 — Idempotência handoff (N4, FA-02) e idempotência aviso (009-9 A, 009-11 B) não se contradizem entre spec, data-model e contracts? [Consistency, Spec User Story 3]
- [x] CHK017 — A meta DN-007-2 B (≥95% precisão por classe) está consistentemente excluída dos SC desta feature e remetida a 006/007? [Consistency, Spec Clarifications etapa 3c; Plan Summary; tasks T047]
- [ ] CHK018 — O Princípio I da constituição (“Finder agrega em ocorrência”) está refletido nos requisitos sem conflito não tratado? [Conflict, Plan Constitution Check] (pendente: emenda do princípio I)
- [x] CHK019 — Decisões 009-5 revisto e 009-5 A histórica estão reconciliadas na tabela de decisões e em FR-003? [Consistency, Spec D12, §Decisões tomadas]

## Acceptance Criteria Quality

- [x] CHK020 — SC-F001–SC-F003 possuem critérios mensuráveis aplicáveis ao escopo Finder (percentuais/contagens explícitas)? [Measurability, Spec §Success Criteria — Finder]
- [ ] CHK021 — SC-002 e aceite formal SC-005 possuem critérios fechados nesta feature, não apenas bloqueio documentado? [Measurability, Spec SC-002/SC-005] (pendente: 009-22 A / spec 008)
- [x] CHK022 — Cenários de aceite das User Stories 1–3 mapeiam para FR/SC identificáveis (rastreio Given/When/Then)? [Acceptance Criteria, Spec User Scenarios]
- [x] CHK023 — O gate de precisão LLT-199 não é misturado indevidamente com SC de handoff LLT-201? [Measurability, Spec DN-007-2; Plan]

## Scenario Coverage

- [x] CHK024 — Fluxo primário de handoff com confirmação e sem `finishScan()` prematuro está coberto em requisitos e cenários? [Coverage, Spec US1, FR-F05]
- [x] CHK025 — Cenários de falha e recuperação Finder→Módulo (reenvio mesma chave, FA-02) e consulta OpenObserve (**FR-F07**, **DN-008-3**, sem DLQ ingested) estão especificados? [Exception/Recovery, Spec FR-F06/FR-F07, Edge Cases]
- [x] CHK026 — Falha HTTP Módulo→Core (009-18 A) e idempotência stream (009-10 B) estão documentados para aceite integrado? [Coverage, Spec FR-007, US2/US3]
- [ ] CHK027 — Requisitos de mascaramento em `AchadoAnalise.linhas` estão fechados o suficiente para aceite do handoff sem ADR 010-2 e casos João 010-27? [Coverage, Spec §Key Entities; Plan §010-2/010-27] (pendente: 010-2 A; 010-27 A)

## Edge Case Coverage

- [x] CHK028 — Divergências João×Drive (INFO, épico LLT-201, correlationId, schemaVersion ingested) estão registradas como requisitos/edge cases, não omitidas? [Edge Cases, Spec §Divergências D1/D3/D7/D9]
- [x] CHK029 — Handoff sem execId (DN-006-1b B) e chave só batchId estão tratados como caso nominal no fluxo atual? [Edge Case, Spec FR-F03, Edge Cases]
- [x] CHK030 — Sugestão automática atrasando aviso ao Core (SSID_009.1 passo 2) consta nos requisitos de referência? [Edge Case, Spec FR-006 / 009-16 rev.]
- [x] CHK031 — C8 (chave concreta ocorrência ADR_040 [PENDENTE]) está documentada como limitação sem inventar chave? [Edge Case, Spec §Questões em aberto, FR-004]

## Non-Functional Requirements

- [x] CHK032 — Restrições de rede e JWT interno (ADR_028/141) estão citadas sem afirmar detalhes não decididos? [NFR, Plan §API; Spec DN-009-2]
- [x] CHK033 — Ausência de meta de latência para handoff/pipeline está declarada onde aplicável (sem SLA inventado)? [NFR, Plan Technical Context Performance Goals]
- [x] CHK034 — Princípio constituição V (sem invenção) aparece em FR-F03 e em registro de dúvidas/pendências? [NFR, Spec FR-F03; Plan Constitution V]

## Dependencies & Assumptions

- [x] CHK035 — Premissa “nenhuma premissa inventada” e fontes prioritárias (João, revisao-T1, respostas-documentadas, etapa 3/3b) estão refletidas na seção Evidência do spec? [Dependencies, Spec §Evidência]
- [x] CHK036 — Bloqueio de implementação do adaptador HTTP está espelhado em plan, tasks (T018/T019/T034) e contrato handoff? [Traceability, Plan; tasks; contracts/resumo-analise-handoff.md — via plan pós-Fase 1]
- [ ] CHK037 — Lista oficial de padrões de mascaramento (010-2 A) está referenciada como dependência bloqueante onde o handoff exige linhas mascaradas? [Dependency, Plan §Pendências; tasks T006/T007] (pendente: 010-2 A)
- [x] CHK038 — Base rotulada 007 (JSONL, caminho datasets/labeled-logs/) está só documentada, sem exigir criação de dados nesta feature? [Dependency, Plan Summary etapa 3b]

## Ambiguities & Conflicts

- [x] CHK039 — Tabela de divergências documentadas cobre conflitos materialmente relevantes (D1–D15) com “qual lado vale”? [Traceability, Spec §Divergências]
- [x] CHK040 — DN-009-2 permanece com opções e recomendação marcada como não-decisão, sem fechamento indevido? [Ambiguity, Spec §DN-009-2]

## Traceability spec ↔ plan ↔ tasks

- [x] CHK041 — Cada FR-F01–FR-F11 do spec possui eco em plan (constraints/structure) ou tasks (IDs T008–T026, T041)? [Traceability, tasks T041; plan Constraints]
- [x] CHK042 — User Stories 1–3 do spec correspondem a fases/tasks US1–US3 na tasks.md com mesmos independent tests (quickstart)? [Traceability, tasks Organization; spec User Scenarios]
- [x] CHK043 — Decisões etapa 3b–3c (DN-006-1/1b, DN-009-1 A, 009-5 revisto, DN-007-2 B, DN-008-4 B) estão aplicadas nos três artefatos? [Traceability, Spec Clarifications; Plan; tasks §Desbloqueadas]
- [x] CHK044 — Itens BLOQUEADO em tasks.md citam IDs de pendência corretos (DN-009-2, N3 A, 009-22 A, P-010-27)? [Consistency, tasks §Pendências]
- [ ] CHK045 — Tarefa de registrar emenda constituição I (T043) reflete requisito já satisfeito nos artefatos ou apenas pendência operacional? [Traceability, tasks T043; plan emenda] (pendente: emenda do princípio I — requisito de governança, não texto emendado)

## Notes

- `/speckit-implement` pode ler este checklist; não alterar marcadores automaticamente.
- Revisão etapa 3c: prioridade fontes João §Etapa 3b (16/11 09h02) e §Etapa 3 (00h41).
- Pendências externas conhecidas citadas na avaliação: **DN-009-2**, **N3 A**, **010-2 A**, **010-27 A**, **006-7 D**, **009-22 A**, **emenda do princípio I**.