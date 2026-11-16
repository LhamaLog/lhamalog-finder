# Checklist: Revisão etapa 3c — qualidade dos requisitos (feature 010)

**Purpose**: Validar completude, clareza, consistência entre spec/plan/tasks, critérios de aceite mensuráveis, cenários de exceção/recuperação, dependências e rastreabilidade às fontes — antes da implementação (revisor PR/gestor).

**Created**: 2026-11-16

**Feature**: [spec.md](../spec.md) | [plan.md](../plan.md) | [tasks.md](../tasks.md)

**Note**: Gerado por `/speckit-checklist` (modo não interativo, profundidade padrão). Reavaliação **etapa 3g** (2026-11-16) dos marcadores contra `spec.md`, `plan.md` e `tasks.md` atuais — sem regenerar itens.

**Review Ownership**: Em execuções normais, `[x]` indica critério de qualidade de requisitos aprovado pelo revisor. Neste arquivo, `[x]`/`[ ]` refletem a avaliação documental vigente após etapa 3g.

**Marker Semantics**: `[x]` = requisito/atendimento documental satisfeito nos artefatos lidos; `[ ]` = lacuna ou pendência externa ainda não fechada nos requisitos (nota `(pendente: …)`).

## Requirement Completeness

- [x] CHK001 — Os requisitos funcionais **FR-001** a **FR-015** estão enumerados na spec com escopo, fonte e pendências explícitas onde couber? [Completeness, Spec §Functional Requirements]
- [ ] CHK002 — A **lista oficial de regex** do catálogo está **definida e referenciável** como artefato ADR publicado (não só “comportamento interim”)? [Completeness, Gap, 010-2 A, FR-001, FR-008] (pendente: **010-2 A** / **P-010-2** — ADR ainda não criada/emendada)
- [ ] CHK003 — Os **casos de teste** acordados para **SC-001** e **SC-002** estão documentados com exemplos concretos nesta feature? [Completeness, Gap, 010-27 A, SC-005] (pendente: **010-27 A** — João ainda não entregou casos)
- [ ] CHK004 — O **contrato HTTP Finder → Módulo Ocorrência** (rota, JWT, códigos, timeout, 5xx) está especificado de forma fechável nos requisitos/plano? [Completeness, Gap, DN-009-2, plan §API do Módulo] (pendente: **DN-009-2** — API do Módulo)
- [ ] CHK005 — Os **campos ampliados de `AchadoAnalise`** (assinatura, contagem, primeira/última vez, versão das regras) estão definidos nos requisitos/contratos além de menção genérica? [Completeness, Gap, N3 A, data-model.md] (pendente: **N3 A** — ampliação depende de ADR)
- [ ] CHK006 — A decisão **006-7 D** (sintaxe/dono do catálogo de classificação — LLT-51 ou ADR nova) está **rastreada** nos artefatos 010 onde o catálogo de mascaramento se cruza com classificação? [Completeness, Traceability, 006-7 D] (pendente: **006-7 D** — não citado explicitamente; só **P-010-2**/ADR_064)
- [x] CHK007 — A dependência **009-22 A** (spec **009** só fecha aceite após **008**) está registrada onde a 010 referencia handoff/eventos downstream? [Completeness, Dependency, 009-22 A; plan.md §Integração; tasks.md T040/T042, §Pendências, §Dependencies]
- [x] CHK008 — As três histórias de usuário (US1–US3) e cenários de aceite Given/When/Then cobrem trechos mascarados, consumo externo e contador? [Completeness, Spec §User Scenarios]
- [x] CHK009 — As clarificações **010-1** a **010-28**, **T1 revisto**, **DN-010-1**, **DN-010-2**, **DN-007-2** e **DN-008-4** estão integradas na spec com formato de citação? [Completeness, Spec §Clarifications]

## Requirement Clarity

- [x] CHK010 — O **momento** do mascaramento (uma vez por lote, antes do MotorAnalise/classificação) está sem ambiguidade entre `montarContexto()` e MotorAnalise? [Clarity, FR-002, 010-3 A, 010-4 A]
- [x] CHK011 — O comportamento em **falha de `validarSeguranca()`** (seguir sem trechos no achado, sem fail-closed do lote) está definido de forma não contraditória? [Clarity, 010-16 C, FR-010]
- [ ] CHK012 — O **placeholder fixo de e-mail** está fechado com literal canônico único (não apenas “ex.: `[email]`”)? [Clarity, 010-11 A, FR-006] (pendente: **010-2 A** / **P-010-2** — literal até ADR)
- [x] CHK013 — O escopo **“só linhas de log”** versus atributos OTel/identificadores **ADR_040** está explícito e alinhado ao DPA? [Clarity, FR-004a, DN-010-2]
- [x] CHK014 — A ordem **mascaramento → classificação → normalização de mensagem para assinatura** está declarada sem a 010 definir algoritmo divergente da **008**? [Clarity, FR-015, DN-008-4 B]
- [x] CHK015 — O **caminho e formato** do artefato versionado de regex no monorepo estão especificados nos requisitos? [Clarity, FR-008, **DN-010-4 C**] — MUST ser os declarados na **ADR da lista** (**010-2 A**); valores concretos aguardam **P-010-2**, sem ambiguidade de *onde* a regra vive (fonte: respostas-joao.md, DN-010-4 C; origem: respostas-etapa3c.md §8)
- [ ] CHK016 — Os **nomes canônicos exatos** do contador “dados sensíveis mascarados” no handoff estão fechados nos requisitos/contratos? [Clarity, Ambiguity, ADR_102] (pendente: rótulos **[INFERIDO]** — `saida-mascarada-handoff.md` §Escopo)

## Requirement Consistency (spec / plan / tasks)

- [x] CHK017 — Spec, plan e tasks concordam que o **Módulo de Ocorrência cria** a Ocorrência e o Finder envia **`ResumoAnalise` mascarado** (T1 revisto)? [Consistency, FR-003, plan §Summary, tasks T027]
- [ ] CHK018 — A **constituição princípio I** (“Finder agrega em ocorrência”) e o **T1 revisto** estão reconciliados em requisito **emendado** (não só nota de pendência)? [Consistency, Conflict, emenda princípio I] (pendente: **emenda do princípio I** — só rastreada em plan/tasks T038; constituição inalterada)
- [x] CHK019 — **DN-007-2 B** (meta ≥95% por classe) aparece como **referência cruzada SC-007**, excluída do aceite LLT-202/010, de forma consistente em spec, plan e tasks (T044)? [Consistency, SC-007, plan §Referência cruzada]
- [x] CHK020 — **DN-008-4 B** e **DN-008-5 A** (literais **`<NUM>`** … **`<HEX>`**; dono spec **008**) estão tratados de forma alinhada entre spec **FR-015**, plan, tasks **T041/T043** (T043 **desbloqueada** etapa 3e) e contratos? [Consistency, DN-008-4 B, DN-008-5 A; fonte: respostas-joao.md §Etapa 3c, João 16/11 11h]
- [x] CHK021 — **009-5 revisto** (id + tipo + gravidade, sem texto de log) e **DN-010-1 A** (e-mail lê Ocorrência por id) não contradizem **FR-004** nem contratos? [Consistency, D10, D12]
- [x] CHK022 — Tarefas com marca **`BLOQUEADO (…)`** em `tasks.md` (**T013/T017** `P-010-2`; **T023/T030/T033/T040** `P-010-27`; **T042** `DN-009-2`; **T040** também cita cadeia **009-22 A**) espelham as pendências em spec §Pendências explícitas e plan §Pendências — **sem** `BLOQUEADO (DN-008-5)` (**T043** e **T046** desbloqueadas após **DN-008-5 A** / **DN-010-3 B**, etapa 3e)? [Consistency, tasks.md, SC-005; analyze etapa 3f **I2**]

## Acceptance Criteria Quality

- [x] CHK023 — **SC-001** a **SC-004** usam percentuais/contagens **mensuráveis** condicionados a dependências documentadas (não afirmam aceite fechado hoje)? [Measurability, SC-001–SC-004, SC-005]
- [x] CHK024 — **SC-005** deixa explícito que aceite formal de SC-001/002 permanece bloqueado por **P-010-27** e **P-010-2**? [Acceptance Criteria, SC-005, P-010-27, P-010-2]
- [x] CHK025 — **SC-006** (paridade Finder–Sugestão, mesma versão do catálogo) é objetivamente verificável quando o artefato existir? [Measurability, SC-006, FR-008]
- [x] CHK026 — Critérios Jira LLT-202 crit. 1/2 estão mapeados a FR/SC sem expandir PII além do Drive (**010-12 C**)? [Traceability, LLT-202, D2]

## Scenario Coverage

- [x] CHK027 — Fluxo **primário** (mascarar → classificar → montar achados/`linhas` → handoff) está coberto em requisitos e contrato `validar-seguranca.md` / `saida-mascarada-handoff.md`? [Coverage, US1, FR-002]
- [x] CHK028 — Fluxos **alternos** (PII fora do catálogo permanece bruto; sem limite de exemplos 010-8 C) estão nos requisitos? [Coverage, FR-009, Edge Cases]
- [x] CHK029 — Fluxos **de exceção** (falha mascaramento 010-16; lista regex pendente) estão documentados sem inventar retry de máscara? [Coverage, Exception Flow, FR-010]
- [x] CHK030 — O fluxo **e-mail LLT-218** está especificado até a fronteira Finder (mascarado em `conteudo` via Notificação por id), com template NOT-03 explicitamente fora de escopo? [Coverage, FR-004, D12]

## Edge Case Coverage

- [x] CHK031 — Riscos **RIPD** (ampliar padrões antes de dados reais) e **D2** (Jira PII amplo × ADR_064) estão registrados sem expandir escopo LLT-202? [Edge Case, 010-12, Spec §Edge/Divergências]
- [x] CHK032 — Divergência **INFO gera ocorrência** (T7 A vs Drive) está em **Divergências documentadas D1** sem exigir emenda ADR_019/167 nesta pasta? [Edge Case, FR-013, D1]
- [x] CHK033 — **Log bruto só OpenObserve** e proibição de reexposição em fila/handoff estão nas FR e alinhados ao contrato v1? [Edge Case, FR-011, 010-24, T5 D]

## Non-Functional Requirements

- [x] CHK034 — Ausência de **meta de desempenho** do mascaramento está explícita (**010-28 A**) em spec e plan? [NFR, Clarity]
- [x] CHK035 — **Princípio V** (sem invenção) refletido em lacunas **NEEDS CLARIFICATION** (biblioteca regex Go, path catálogo) sem decisão silenciosa? [NFR, plan §Technical Context, Constitution V]

## Dependencies & Assumptions

- [x] CHK036 — Dependência da spec **003** (captura integral) e não contradição **FR-012** documentada? [Dependency, FR-012, tasks T034]
- [x] CHK037 — Dependência da spec **008** para normalização (**DN-008-4 B**) e literais (**DN-008-5**) sem a 010 ser dona do algoritmo? [Dependency, FR-015, Spec §Fora de escopo]
- [x] CHK038 — **Base rotulada 007** (JSONL, `datasets/labeled-logs/`) documentada apenas como referência cruzada, sem exigir criação de pasta nesta feature? [Assumption, plan §Referência cruzada, 007 formato A / caminho C]
- [x] CHK039 — Premissa **Sugestão (Java)** e paridade de catálogo (**010-18/19 A**) registrada com ação de coordenação fora do repo (tasks T039)? [Assumption, FR-008]

## Ambiguities, Conflicts & Traceability

- [x] CHK040 — Tabela **Decisões tomadas** e **Divergências documentadas** permitem rastrear cada decisão 010/T1/DN-* à fonte `(respostas-joao.md, ID)`? [Traceability, Spec §Decisões, §Divergências]
- [x] CHK041 — **Dúvidas novas** na spec listam DN-010-1/2, DN-007-2, DN-008-4 como **respondidas**, sem decisões fantasma? [Traceability, Spec §Dúvidas novas]
- [x] CHK042 — Contratos em `contracts/` permanecem coerentes com FR e não afirmam transporte HTTP fechado enquanto **DN-009-2** estiver aberto? [Consistency, saida-mascarada-handoff.md §Transporte]

## Notes

- `/speckit-implement` pode usar este checklist como gate; não alterar marcadores sem nova revisão.
- `checklists/requirements.md` é artefato separado (speckit-specify/clarify) — não confundir com este arquivo.
- Pendências obrigatórias ainda **[ ]** na avaliação: **DN-009-2**, **N3 A**, **010-2 A** / **P-010-2**, **010-27 A**, **006-7 D**, **emenda do princípio I**; valores concretos **placeholder e-mail** e **rótulos contador ADR_102** (**CHK016**). **009-22 A** e regra **DN-010-4 C** (path/formato na ADR) constam em plan/tasks — **CHK007**, **CHK015** atendidos.