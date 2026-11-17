# Checklist de revisão (etapa 3c): qualidade dos requisitos — Agrupar logs na mesma ocorrência

**Purpose**: Validar completude, clareza, consistência (spec/plan/tasks/contratos), critérios de aceite mensuráveis, cenários de exceção/recuperação, dependências e rastreabilidade às fontes prioritárias — **antes da implementação** (revisor PR/gestor).

**Created**: 2026-11-16

**Feature**: [spec.md](../spec.md)

**Note**: Checklist gerado por `/speckit-checklist` (modo não interativo, etapa 3c). Foco: qualidade dos requisitos da feature 008 como um todo; profundidade padrão.

**Review Ownership**: `[x]` indica que o critério de **qualidade dos requisitos** foi avaliado e atendido nos artefatos atuais. Não indica implementação concluída.

**Avaliação**: 2026-11-16 — reavaliação pós-etapas **3e/3f** (leitura integral de `spec.md`, `plan.md`, `tasks.md` e contratos em `contracts/`; fontes prioritárias João — DN-008-5 A, DN-008-6, DN-008-7 B, 008-assinatura A; **DN-008-8** aberta).

---

## Requirement Completeness

- [x] CHK001 — A responsabilidade de **criar/persistir ocorrência** (Módulo) versus **agrupar e enviar ResumoAnalise** (Finder) está documentada com IDs e fontes? [Completeness, Spec §Decisões T1 rev., FR-001]
- [x] CHK002 — A **composição fechada da assinatura** (nível + mensagem normalizada + `service.name` + host) está especificada em requisito, entidade e contrato? [Completeness, Spec §FR-002, contracts/assinatura-log.md; DN-008-1 B]
- [x] CHK003 — O **algoritmo canônico de normalização da mensagem** (trim, case-insensitive, cinco categorias de marcadores) está como FR explícito com dono e consumidores? [Completeness, Spec §FR-002b; DN-008-4 B]
- [x] CHK004 — A **chave de unicidade** `(assinatura, classe, produto, tenant)` cobre FR, cenários de aceite e modelo de dados? [Completeness, Spec §FR-005, FR-011; 008-11 rev.]
- [x] CHK005 — Requisitos de **agregação entre lotes** (consulta antes do envio, fusão de contador e extremos) estão em FR e user stories? [Completeness, Spec §FR-001; N2 B]
- [x] CHK006 — Regras de **recorrência** (ADR_013, janela 500 lotes, OpenObserve, separação de status) estão documentadas? [Completeness, Spec §FR-003, FR-007, FR-008; 008-5/6]
- [x] CHK007 — Política de **falha temporária OpenObserve** (3 tentativas, mesma stream, DLQ; proibição de `recorrente=false` por degradação) está em requisito dedicado? [Completeness, Spec §FR-002a; DN-008-3]
- [x] CHK008 — **Reabertura** FECHADA→ABERTA e **não alteração de status só por contador** estão em FR e US3? [Completeness, Spec §FR-003, FR-010; 008-4, 008-10]
- [x] CHK009 — **Classe INFO** no agrupamento e divergência Drive estão registradas (decisão + fora de escopo onde couber)? [Completeness, Spec §FR-006; T7 A; D1]
- [x] CHK010 — **Lote incompleto** (spec 005) e **productId inválido** têm FR e edge cases? [Completeness, Spec §FR-013, FR-014; T4 A, 008-16 A]
- [ ] CHK011 — O **contrato HTTP completo** Finder→Módulo (rota, JWT, corpo, códigos, 5xx/timeout) está especificado além de lacunas nomeadas? [Gap, Spec §Questões em aberto; DN-009-2] (pendente: DN-009-2)
- [x] CHK012 — Os **literais exatos dos marcadores** de normalização estão definidos nos requisitos ou contrato versionado? [Completeness, Spec §FR-002b; DN-008-5 A; contracts/assinatura-log.md; 010-27 A matriz formal em P-010-27 A]
- [ ] CHK013 — O **schema ampliado de AchadoAnalise** (assinatura, contagem, primeira/última, versão regras) está fechado em requisito/contrato sem depender só de placeholder ADR? [Gap, Spec §N3 A; contracts/resumo-analise-agrupado.md] (pendente: N3 A)
- [x] CHK014 — **Migração de ocorrências legadas** e ressalva de ambiente de teste estão documentadas? [Completeness, Spec §FR-009, User Story 2; DN-008-2, D2b]
- [x] CHK015 — **Fora de escopo** lista classificação, base 007, evento/notificação/mascaramento, ADRs proibidas e meta 95%? [Completeness, Spec §Fora de escopo; DN-007-2 B]

## Requirement Clarity

- [x] CHK016 — Fica explícito que **`traceId`/`correlationId` não entram na assinatura** após revisão 3b (substitui 008-2 A)? [Clarity, Spec §FR-002; D2; 008-2 rev.]
- [x] CHK017 — A distinção entre **`nivel` na assinatura** e **`classe` na chave** está sem ambiguidade de duplicação? [Clarity, Spec §FR-002, FR-002 nota 008-15]
- [x] CHK018 — **Contador** (+1 por log agregado) e **primeira/última vez** (timestamp do log, não processamento) são mensuráveis em SC? [Clarity, Spec §SC-002; 008-12, 008-14]
- [x] CHK019 — **Repetição no mesmo lote** conta para contador mas não para recorrência ADR_013 está claro? [Clarity, Spec §Edge Cases; 008-13 A]
- [x] CHK020 — **Resolvida = FECHADA** e papel do **operador** no fechamento estão definidos? [Clarity, Spec §FR-012; 008-8, 008-9]
- [x] CHK021 — A **serialização canônica do valor `assinatura`** (hash vs string) está definida sem NEEDS CLARIFICATION? [Clarity, Spec §FR-002c, 008-assinatura A; contracts/assinatura-log.md §Serialização; separador exato permanece DN-008-8]
- [x] CHK022 — **Idempotência ResumoAnalise** (`batchId` + `execId` opcional) e origem do execId estão claros? [Clarity, Spec §N4 A; plan.md §Pendências DN-006-1b]
- [x] CHK023 — Ordem lógica das etapas de **DN-008-4 B** (trim → case → substituições) está explícita? [Clarity, Spec §FR-002b]

## Requirement Consistency

- [x] CHK024 — `spec.md`, `plan.md` e `tasks.md` concordam sobre **T1 rev.** (Finder MUST NOT gravar ocorrência)? [Consistency, FR-001; plan §Summary; tasks T031]
- [x] CHK025 — Assinatura e normalização são **consistentes** entre spec, plan, data-model/contratos e tarefas T015/T016/T047? [Consistency, plan §Constraints; tasks Phase 2]
- [x] CHK026 — **Chave 008-11 rev.** é a mesma em FR-005, consulta Módulo e particionamento de achados? [Consistency, contracts/consulta-ocorrencias-modulo.md; tasks T007/T019]
- [x] CHK027 — **SC-005** em spec alinha com bloqueios em plan (DN-009-2, N3) e tarefas BLOQUEADO T021/T024? [Consistency, Spec §SC-005; plan §Tarefas bloqueadas]
- [x] CHK028 — Tabela **Divergências documentadas** não contradiz decisões vencedoras nas Clarifications (ex.: D2 → DN-008-1 B)? [Consistency, Spec §Divergências]
- [ ] CHK029 — O **princípio I da constituição** (“logs em ocorrências”) está **reconciliado no texto constitucional** ou só como pendência registrada? [Conflict, plan §Constitution Check; emenda princípio I] (pendente: emenda princípio I)
- [x] CHK030 — **Contrato log.ingested v1** (T5 D) é o único envelope citado sem misturar ADR_146 v2 como válido? [Consistency, Spec §D7, FR-014 path; tasks T014]

## Acceptance Criteria Quality

- [x] CHK031 — **SC-001** (uma ocorrência por chave) é testável com assinatura fechada e cenário US1.5 (DN-008-4 B)? [Measurability, Spec §SC-001; User Story 1 cenário 5]
- [x] CHK032 — **SC-003** cobre reabertura e recorrência com critérios objetivos? [Measurability, Spec §SC-003]
- [x] CHK033 — **SC-007** materializa DN-008-3 em percentual de cenários de teste declarado? [Measurability, Spec §SC-007]
- [ ] CHK034 — **SC-001/SC-005** podem fechar aceite formal **ponta a ponta** sem as lacunas ainda nomeadas em SC-005? [Measurability, Spec §SC-005; DN-009-2; DN-008-8; N3 A] (pendente: DN-009-2, DN-008-8, N3 A — literais **DN-008-5 A** já fechados na etapa 3e)
- [x] CHK035 — User stories P1 têm **Independent Test** reproduzível sem HTTP real onde aplicável? [Measurability, Spec §User Story 1–3; tasks T049]

## Scenario Coverage

- [x] CHK036 — Fluxo **primário** (classificar → assinar → agrupar lote → consultar → ResumoAnalise → Módulo upsert) está coberto? [Coverage, Spec §US1; FR-001, FR-009]
- [x] CHK037 — Cenários **alternos** (produtos/classes/tenants distintos; mesma assinatura) estão nos acceptance scenarios? [Coverage, Spec §US1 cenários 3–4]
- [x] CHK038 — Fluxos **de exceção** (lote incompleto, produto inválido, OpenObserve indisponível) têm requisito e edge case? [Coverage, Spec §Edge Cases; FR-002a, FR-013, FR-014]
- [x] CHK039 — Fluxo **pós-resolução** (FECHADA + novo log) está em US3 e FR-003? [Coverage, Spec §US3]
- [ ] CHK040 — Dependência de **ordem entre features** (ex.: publicação de evento / spec 009 após agrupamento) está rastreada nesta spec? [Gap; 009-22 A] (pendente: 009-22 A — só fora de escopo 009/010, sem FR de sequência na 008)

## Edge Case Coverage

- [x] CHK041 — **Normalização** cobre maiúsculas, espaços e variantes numéricas/UUID/IP/data/hex por requisito, com vínculo a casos de teste futuros? [Edge Case, Spec §FR-002b; Edge Cases §DN-008-4]
- [x] CHK042 — **Ambiente de teste** com ocorrências pré-gravadas permanece condicionado a confirmação (D2b)? [Edge Case, Spec §D2b; DN-008-2]
- [x] CHK043 — **Alerta ao Core** além de DLQ na falha OpenObserve está documentado com fonte (sem invenção)? [Edge Case, Spec §FR-002a, SC-007; DN-008-6; respostas-joao.md §Etapa 3c 11h]
- [x] CHK044 — **Política ADR_035 / notificação duplicada** está fora de escopo com dependência explícita? [Edge Case, Spec §FR-015; 008-17 C]

## Non-Functional Requirements

- [x] CHK045 — Ausência de **meta de performance** OpenObserve está declarada como fora de escopo (sem número inventado)? [NFR, Spec §Fora de escopo; plan §Performance Goals]
- [x] CHK046 — **Isolamento tenant** na chave e limites do que fica fora (usuário, projeto, etc.) estão em FR? [NFR/Security boundary, Spec §FR-011]
- [x] CHK047 — Princípio **incerteza sem invenção (V)** aparece em SC-005, DN-008-5 e contratos com proibição de inventar? [NFR, Spec §FR-002b; plan §Constitution V]

## Dependencies & Assumptions

- [x] CHK048 — Dependência da **spec 006** (classificação, evolução ResumoAnalise, conflito D8) está documentada? [Dependency, Spec §FR-009; D8]
- [ ] CHK049 — Impacto de **006-7 D** (sintaxe catálogo / dono em LLT-51 ou ADR) nas regras consumidas pela 008 está rastreado? [Dependency, Gap; 006-7 D] (pendente: 006-7 D)
- [x] CHK050 — **010-2 A** e **010-27 A** como fontes de marcadores/casos estão registrados em plan/tasks sem criar ADR aqui? [Dependency, plan §Pendências; tasks T048/T051/T052]
- [x] CHK051 — **API Módulo Ocorrência (DN-009-2)** está listada como bloqueio com escopo desbloqueado documentado? [Dependency, plan §API do Módulo; tasks §Bloqueios]
- [x] CHK052 — **N3 A (ADR AchadoAnalise)** está como dependência externa com campos opcionais até ADR? [Dependency, Spec §Questões em aberto; tasks T005/T023]
- [x] CHK053 — **007 formato/caminho** (JSONL, `datasets/labeled-logs/`) aparece só referência cruzada sem exigir pasta na 008? [Dependency, plan §Testing; tasks T054; Etapa 3b]
- [x] CHK054 — **Emenda constitucional princípio I** tem ação/risco registrados sem alterar `.specify/memory/constitution.md`? [Assumption, plan §Constitution; tasks T053]

## Ambiguities & Conflicts

- [x] CHK055 — Divergências João × Drive (INFO, T1, chave, contrato v1) estão tabuladas com lado vencedor? [Conflict resolution, Spec §Divergências]
- [x] CHK056 — Substituição **008-2 A → assinatura sem trace/correlation** está refletida em clarifications e tabela de decisões (sem resíduo contraditório em FR)? [Conflict, Spec §D2; FR-002]
- [x] CHK057 — **Granularidade** do handoff e consulta prévia (lote vs achado) tem decisão do João refletida nos requisitos? [Clarity, Spec §FR-009, Edge Cases; DN-008-7 / DN-008-7 B; contracts/consulta-ocorrencias-modulo.md; rotas HTTP permanecem DN-009-2]
- [x] CHK058 — **DN-008-5 A** (literais dos cinco marcadores) está fechada e rastreada nos requisitos/contratos, sem suposição? [Traceability, Spec §FR-002b, Clarifications etapa 3e; contracts/assinatura-log.md]
- [ ] CHK059 — **006-10 A** (alinhar Jira LLT-200) está refletido como ação rastreável além de linha em tasks (requisito/checklist de produto)? [Gap, tasks T006; 006-10 A] (pendente: 006-10 A)

## Rastreabilidade e alinhamento plan/tasks

- [x] CHK060 — Tarefas de implementação citam **(fonte: …; origem: …)** nas fases críticas (assinatura, agrupamento, recorrência)? [Traceability, tasks.md §Format]
- [x] CHK061 — Fases bloqueadas (**DN-009-2**, **P-010-27 A**, **N3 A**) espelham SC-005 e plan §Tarefas bloqueadas? [Consistency, tasks §Bloqueios; Spec §SC-005]
- [x] CHK062 — Contratos em `contracts/` referenciam decisões 3b–3c (**DN-008-1 B**, **DN-008-4 B**) sem rotas inventadas? [Traceability, assinatura-log.md; consulta-ocorrencias-modulo.md]
- [x] CHK063 — **Quickstart** e research estão referenciados em plan/tasks para validação guiada pós-requisitos? [Completeness, plan §Documentation; tasks T027/T043]
- [x] CHK064 — Decisões **Etapa 3b João 09h02** (DN-007-2 B fora; DN-008-4 B dono 008) aparecem em spec e propagam a plan? [Traceability, respostas-joao.md §Etapa 3b; Spec §Session etapa 3c]

## Notes

- Reavaliação pós-**3e/3f**: pendências conhecidas citam IDs explícitos (**DN-009-2**, **DN-008-8**, **N3 A**, **P-010-27 A**, **006-7 D**, **009-22 A**, **006-10 A**, emenda princípio I). **DN-008-5 A**, **DN-008-6**, **DN-008-7 B** e **008-assinatura A** (salvo separador) atendidos nos artefatos atuais.
- `/speckit-implement` pode ler este arquivo como gate; não alterar marcadores automaticamente.
- Não confundir com [requirements.md](./requirements.md) (ciclo specify/clarify).