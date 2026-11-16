# Revisão etapa 3c Checklist: Modelo de classificação de logs

**Purpose**: Validar qualidade dos requisitos da feature 006 como um todo (completude, clareza, consistência spec/plan/tasks, critérios mensuráveis, exceção/recuperação, dependências e rastreabilidade) antes da implementação.
**Created**: 2026-11-16
**Feature**: [spec.md](../spec.md)
**Audience**: Revisor pré-implementação (PR/gestor)
**Depth**: Padrão
**Etapa**: 3c (DN-007-2 B, DN-008-4 B aplicados em 2026-11-16)

**Note**: Itens gerados por `/speckit-checklist` (modo não interativo). Avaliação automática contra spec.md, plan.md e tasks.md na mesma execução.
**Review Ownership**: `[x]` indica critério de qualidade de requisitos **satisfeito** na avaliação desta rodada; `[ ]` permanece pendente para o revisor ou até fechamento de lacuna citada.

## Requirement Completeness

- [x] CHK001 — O fluxo vigente (Finder envia `ResumoAnalise`/`AchadoAnalise`; Módulo Ocorrência cria ocorrência) está documentado em requisitos obrigatórios e user stories? [Completeness, Spec §Clarifications T1 REVISTO, FR-007; fonte: respostas-joao.md, T1 REVISTO; origem: ADR_062]
- [x] CHK002 — As três classes ERROR, WARNING e INFO (inclusive INFO como achado/ocorrência) estão exigidas de ponta a ponta na spec? [Completeness, Spec §FR-003, T7 A/006-11 A; fonte: respostas-joao.md, T7 A; origem: ADR_170]
- [x] CHK003 — A composição da assinatura do achado (nível + mensagem normalizada + `service.name` + host) está fechada nos FR e entidades? [Completeness, Spec §FR-007, DN-008-1 B; fonte: respostas-joao.md, Etapa 3; origem: ADR_013]
- [x] CHK004 — O algoritmo de normalização da mensagem (DN-008-4 B) e o dono canônico (spec 008) estão exigidos sem contradizer a assinatura? [Completeness, Spec §FR-013, Session etapa 3c; fonte: respostas-joao.md §Etapa 3b, DN-008-4 B; origem: ADR_013]
- [x] CHK005 — Regras de chave de idempotência (`batchId` obrigatório; `execId` quando houver; só `batchId` no fluxo atual) estão em FR dedicado? [Completeness, Spec §FR-011, DN-006-1/DN-006-1b B; fonte: respostas-joao.md, Etapa 3; origem: contrato-log-ingested-v1.txt]
- [ ] CHK006 — Os requisitos de handoff HTTP ao Módulo Ocorrência (rota, claims JWT, códigos HTTP, timeout, política 5xx) estão completos o bastante para implementar o adaptador sem lacuna? [Gap, Spec §FR-007, Plan §API do Módulo; pendente: DN-009-2 API do Módulo]
- [ ] CHK007 — A sintaxe, o formato e o dono do catálogo de classificação estão definidos nos requisitos desta feature? [Gap, Spec §FR-005, Pendências 006-7 D; pendente: 006-7 D]
- [ ] CHK008 — O contrato formal ampliado de `AchadoAnalise` (assinatura, contagem, primeira/última, versão) está fechado em requisitos sem depender só de placeholders? [Gap, Spec §FR-007, N3 A; pendente: N3 A]
- [x] CHK009 — Gate de avaliação offline (precisão ≥95% **por classe**) e vínculo à base rotulada 007 (JSONL, caminho documentado) estão nos requisitos de aceite? [Completeness, Spec §SC-007, FR-010; fonte: respostas-joao.md §Etapa 3b, DN-007-2 B; origem: LLT-199]

## Requirement Clarity

- [x] CHK010 — A precedência campo de nível × termos na mensagem × catálogo está especificada sem ambiguidade (palavra inteira; mais grave vence)? [Clarity, Spec §FR-004, User Story 1/3]
- [x] CHK011 — O critério de sucesso de lote incompleto (T4 A) distingue-se explicitamente de falha de envio ao Módulo? [Clarity, Spec §FR-008, Edge Cases; fonte: respostas-joao.md, T4 A]
- [x] CHK012 — A meta numérica **007-meta** está quantificada por classe (ERROR, WARNING, INFO) e ligada a SC/T047? [Clarity, Spec §SC-007; fonte: respostas-joao.md §Etapa 3b, DN-007-2 B]
- [x] CHK013 — Os passos do normalizador DN-008-4 B (trim, case-insensitive, categorias número/UUID/IP/data/hex) estão descritos de forma implementável na spec? [Clarity, Spec §FR-013, Session etapa 3c]
- [x] CHK014 — Os literais exatos dos marcadores de normalização estão definidos nos requisitos desta feature ou referência fechada? [Ambiguity, Spec §FR-013, Key Entities; fonte: respostas-joao.md §Etapa 3c, DN-008-5 A]
- [x] CHK015 — O escopo «agrupamento somente no lote» versus consulta entre lotes (N2 B) está delimitado nos FR sem misturar responsabilidades? [Clarity, Spec §FR-009, Assumptions; fonte: respostas-joao.md, N2 B]

## Requirement Consistency

- [x] CHK016 — spec.md, plan.md e tasks.md concordam sobre T1 REVISTO (sem MongoDB de ocorrência no Finder)? [Consistency, Plan §Constitution Check, tasks §Fora do escopo]
- [x] CHK017 — Contador «spam ignorado» versus contador INFO está alinhado entre spec, data-model/contracts e tasks (T007/T029)? [Consistency, Spec §FR-007; fonte: respostas-joao.md, Etapa 3; origem: ADR_102, ADR_170]
- [x] CHK018 — DN-008-1 B (sem traceId/correlationId na assinatura 006) está consistente com a divergência D5 e com plan/research? [Consistency, Spec §D5, research.md R14]
- [x] CHK019 — Pendência DN-009-2 está referenciada de forma coerente em spec (FR-007), plan (§API), tasks (T031/T045) e contract handoff? [Consistency, contracts/resumo-analise-handoff.md]
- [ ] CHK020 — O princípio I da constituição («Finder agrega em ocorrência») está reconciliado nos requisitos com T1 REVISTO sem contradição não nomeada? [Conflict, Plan §Constitution Check; pendente: emenda do princípio I]
- [x] CHK021 — Versão única regras+catálogo+base (006-3 C) e política só novos logs (006-4 A) são consistentes entre FR-010, FR-012 e tasks T008/T048? [Consistency, Spec §FR-012; fonte: respostas-joao.md, 006-3 C/006-4 A]

## Acceptance Criteria Quality

- [x] CHK022 — SC-007 expressa gate mensurável (≥95% precisão por classe) ligado a procedimento em tasks (T047)? [Measurability, Spec §SC-007, tasks §T047]
- [x] CHK023 — Limiares numéricos de recall e F1 por classe estão definidos para aceite LLT-199 alinhado a 006-2 B? [Measurability, Spec §SC-007, DN-006-2 B; fonte: respostas-joao.md §Etapa 3c, 11h/11h01]
- [x] CHK024 — User Stories 1–3 possuem testes independentes descritos e mapeáveis a tasks de teste? [Measurability, Spec §User Scenarios; tasks §Phases 3–5]
- [ ] CHK025 — Requisitos de mascaramento de linhas em `AchadoAnalise` referenciam lista oficial de padrões com fonte fechada? [Traceability, Spec §Edge Cases, tasks T044; pendente: 010-2 A]
- [ ] CHK026 — Casos de teste determinísticos de normalização/mascaramento exigidos pelo João estão especificados ou rastreados até artefato fechado? [Measurability, tasks T052; pendente: 010-27 A]

## Scenario Coverage

- [x] CHK027 — Cenários primários de classificação (nível, texto, catálogo, INFO default) cobrem FR-004–FR-006 nas acceptance scenarios? [Coverage, Spec §User Story 1/3]
- [x] CHK028 — Cenários de montagem de achados e envio de resumo (várias assinaturas, INFO por achado, confirmação Módulo) estão nos requisitos? [Coverage, Spec §User Story 2]
- [ ] CHK029 — Requisitos de confirmação HTTP real e conclusão de varredura cobrem resposta do Módulo quando o contrato HTTP existir? [Coverage, Gap; pendente: DN-009-2 API do Módulo]
- [x] CHK030 — Rejeição de `schemaVersion` ≠ `"1"` no consumo ingested está nos requisitos? [Coverage, Spec §FR-001, Edge Cases; fonte: respostas-joao.md, T5 D]

## Edge Case Coverage

- [x] CHK031 — Apontamento sem logs, armazenamento inacessível e lote incompleto têm requisitos de falha distinguível? [Edge Case, Spec §FR-008, Edge Cases]
- [x] CHK032 — Reenvio com mesma chave de idempotência e papel do Módulo na deduplicação estão documentados? [Edge Case, Spec §Edge Cases, 006-9 A; fonte: respostas-joao.md, 006-9 A; origem: ADR_014]
- [x] CHK033 — Ausência de `execId` no pipeline atual e transição futura para `batchId`+`execId` estão nas bordas e FR-011? [Edge Case, Spec §D8, FR-011]
- [ ] CHK034 — Dependência de agrupamento/recorrência entre lotes (N2 B) e ordem de aceite downstream (009 espera 008) está explicitada nos requisitos 006 para limites de aceite? [Coverage, Spec §Assumptions/FR-009; pendente: 009-22 A — spec 008 antes de fechar aceite 009/recorrência]

## Non-Functional Requirements

- [x] CHK035 — Proibição de LLM no classificador está em FR e success criteria? [NFR, Spec §FR-002, SC-001]
- [x] CHK036 — Rastreio de versão das regras no achado/resumo sem gravar na ocorrência pelo Finder está nos requisitos? [NFR, Spec §FR-007, 006-3 C]
- [x] CHK037 — Meta de tempo SC-002 está qualificada (ciclo classificar+montar resumo, sem persistência MongoDB Finder)? [Clarity, Spec §SC-002]

## Dependencies & Assumptions

- [x] CHK038 — Tabela ou seção de pendências explícitas lista 006-7 D, N3 A, API Módulo, DN-006-2, DN-008-5 com ação bloqueada? [Dependencies, Spec §Pendências explícitas, Plan §Pendências]
- [ ] CHK039 — Alinhamento Jira LLT-197/LLT-199 com FR/SC está previsto como requisito de processo rastreável? [Dependency, Spec §006-10 A; pendente: 006-10 A]
- [x] CHK040 — Assunção de fluxo normativo Finder→HTTP JWT→Módulo cita ADR_062 e decisões João com IDs? [Assumption, Spec §Assumptions, §Decisões tomadas]

## Traceability & Conflicts

- [x] CHK041 — Decisões etapa 3b/3c (DN-006-1b, DN-008-1 B, DN-007-2 B, DN-008-4 B, contador INFO) aparecem em spec, plan e tasks com mesmos IDs? [Traceability, tasks §Desbloqueados 3b–3c]
- [x] CHK042 — Divergências conscientes (INFO×Drive, T1×constituição histórica, assinatura×008-2 A) estão documentadas sem apagar histórico? [Traceability, Spec §Divergências D1/D2/D5]
- [x] CHK043 — DN-006-2 está registrada como dúvida nova com opções e recomendação, sem decisão inventada? [Traceability, Spec §Dúvidas novas; fonte parcial: 006-2 B/DN-007-2 B]
- [x] CHK044 — tasks.md marca bloqueios (DN-009-2, 006-7 D, DN-008-5/010-27) nas mesmas linhas que as pendências da spec/plan? [Consistency, tasks §Pendências, §Bloqueios]

## Ambiguities & Conflicts (residual)

- [x] CHK045 — Gravidade e aviso ao Core (009-5 revisto) permanecem fora do Finder com papéis claros? [Clarity, Spec §FR-007, D10; fonte: respostas-joao.md, Etapa 3]
- [ ] CHK046 — Emenda do princípio I da constituição está tratada como pendência de governança com proposta em plan, não como requisito já fechado? [Governance, Plan §Pendências emenda; pendente: emenda do princípio I]

## Notes

- Avaliação: 2026-11-16, etapa 3c; fontes prioritárias: respostas-joao.md §Etapa 3b/3, respostas-etapa3.md, artefatos 006 spec/plan/tasks.
- Pendências conhecidas forçam `[ ]` nos itens que exigem fechamento correspondente (João).
- `/speckit-implement` lê estado dos checkboxes; este arquivo não altera `requirements.md`.
- Revisão humana pode reabrir itens marcados `[x]` se novas decisões mudarem requisitos.