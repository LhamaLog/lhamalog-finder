# Checklist: Revisão etapa 3c — qualidade dos requisitos (feature 007)

**Purpose**: Validar completude, clareza, consistência entre spec/plan/tasks, critérios de aceite mensuráveis, cenários de exceção, dependências e rastreabilidade da feature `007-base-logs-rotulados` antes da implementação.

**Created**: 2026-11-16

**Feature**: [spec.md](../spec.md) · [plan.md](../plan.md) · [tasks.md](../tasks.md)

**Note**: Gerado por `/speckit-checklist` (modo não interativo; foco integral da feature; público revisor PR/gestor; profundidade padrão). **Reavaliado** em 2026-11-16 (etapa 3e/3f) contra `spec.md`, `plan.md` e `tasks.md` atuais — `[x]` = requisito documentado de forma satisfatória; `[ ]` = lacuna ou pendência externa explícita.

**Review Ownership**: Checklist de qualidade de requisitos (não verificação de implementação).

## Requirement Completeness

- [x] CHK001 — Os FR-001 a FR-013 cobrem objetivo MVP (avaliação), formato, campos, privacidade, métricas e subconjunto só-catálogo? [Completeness, Spec § Requirements]
- [x] CHK002 — Formato **JSONL** e caminho **`datasets/labeled-logs/`** estão exigidos em spec, plan e tasks? [Completeness, Spec § FR-004; Clarifications etapa 3c; plan.md § Summary; tasks.md Path Conventions]
- [x] CHK003 — Campos obrigatórios por registro (`level_text`, `message_text`, ids, tenant, produto, `class_label`) estão listados com tabela de mapeamento 007-6 D? [Completeness, Spec § Mapeamento de campos; FR-004]
- [x] CHK004 — Declaração de origem, método de rotulação e política de duplicatas estão previstos como entregáveis documentais? [Completeness, Spec § FR-001, FR-006, FR-012; US1/US3]
- [x] CHK005 — Subconjunto **só-catálogo**: critério de inclusão e adiamento de **conteúdo** até catálogo estão em spec/plan/tasks? [Completeness, Spec § FR-011, DN-007-1 C; plan.md § Tarefas bloqueadas; tasks T033/T051]
- [x] CHK006 — Metas numéricas de **recall** e **F1 por classe** estão definidas nos requisitos de aceite além do reporte 006-2 B? [Completeness, Spec § FR-007; SC-005; 007-meta; DN-007-3 B; etapa 3e]
- [x] CHK007 — Integração com classificador 006 (nível+mensagem, desacoplada da fila) está referenciada em contratos e tasks? [Completeness, contracts/entrada-classificador-desde-base.md; FR-009; tasks T015–T016]
- [x] CHK008 — Pendências externas conhecidas (DN-009-2, N3 A, 010-2 A, 010-27 A, 006-7 D, emenda princípio I) constam em plan.md § Pendências com impacto na 007? [Completeness, plan.md § Pendências; tasks § Pendências]

## Requirement Clarity

- [x] CHK009 — Meta de precisão LLT-199 critério 2 está quantificada como **≥ 95% em cada classe** ERROR, WARNING e INFO? [Clarity, Spec § FR-007, SC-005; DN-007-2 B]
- [ ] CHK010 — O tamanho mínimo **“milhares”** e a **contagem mínima numérica por classe** estão especificados com valor fixo antes da entrega? [Clarity, Spec § Cobertura alvo; SC-003] (pendente: valores finais só na entrega — tabela com placeholders)
- [x] CHK011 — Rejeição vs permissão (sem nível/mensagem, rótulo inválido, duplicata com rótulos diferentes) está sem ambiguidade nos edge cases? [Clarity, Spec § Edge Cases; FR-012, FR-002]
- [x] CHK012 — Escopo “rótulo só classe” (sem assinatura LLT-200, sem PII) está explícito em FR e fora de escopo? [Clarity, Spec § FR-010; Fora de escopo]
- [ ] CHK013 — Critério de **split** (proporção, seed, estratificação) está definido além de “nomes e critérios documentados”? [Ambiguity, data-model.md § ParticaoSplit; research.md R9; tasks T032] (pendente: NEEDS CLARIFICATION — 007-20 B não fixa algoritmo)
- [x] CHK014 — Modelo gold em **apontamento** (`class_label` único vs por `member`) está fechado nos requisitos? [Clarity, Spec § FR-007b; US3 cenário 4; data-model.md § RegistroApontamentoRotulado; contracts/registro-rotulado.md; DN-007-4 A; etapa 3e]

## Requirement Consistency

- [x] CHK015 — Spec, plan e tasks concordam em **JSONL**, **`datasets/labeled-logs/`** e precisão **por classe** (etapa 3c)? [Consistency, etapa 3b/3c; tasks desbloqueados T003/T007/T044/T055]
- [x] CHK016 — Taxonomia ERROR/WARNING/INFO é coerente entre spec, constituição III (via ADR_170) e proibição de spam/não classificado? [Consistency, Spec § FR-002; Divergências]
- [x] CHK017 — Uso MVP **somente avaliação** (007-3 A) não conflita com split documentado para pós-MVP? [Consistency, Spec § FR-005, SC-006; plan.md § Summary]
- [x] CHK018 — `checklists/requirements.md` reflete lacunas equivalentes (**split** CHK013, conteúdo só-catálogo **006-7 D**) sem afirmar critérios plenos; **DN-007-3 B** / **DN-007-4 A** fechadas na etapa 3e (analyze2 I1; etapa 3f; Spec § Session etapa 3e)
- [x] CHK019 — 007-meta (precisão, recall e F1 **por classe**) e 006-2 B estão harmonizados na tabela de divergências após **DN-007-3 B**? [Consistency, Spec § Divergências 007-meta × 006-2 B; etapa 3e]
- [ ] CHK020 — Requisitos de governança constitucional refletem **T1 revisto** sem contradição não documentada? [Consistency, plan.md § Constitution Check] (pendente: emenda do princípio I — constituição ainda descreve Finder agregando ocorrência)

## Acceptance Criteria Quality

- [x] CHK021 — SC-001 a SC-004 possuem critérios objetivos (100% cobertura, taxonomia, contagens, campos)? [Measurability, Spec § Success Criteria]
- [x] CHK022 — SC-005 define limiar mensurável de **precisão por classe** alinhado a FR-007? [Measurability, Spec § SC-005; DN-007-2 B]
- [x] CHK023 — SC-005 pode ser considerado **completo** para LLT-199/197 incluindo limiares de recall e F1? [Measurability, Spec § SC-005; FR-007; DN-007-3 B; 006-2 B; etapa 3e]
- [x] CHK024 — Cenários Given/When/Then das user stories são verificáveis em termos de requisitos (não de código)? [Measurability, Spec § User Scenarios]
- [ ] CHK025 — Tabela de contagens em `checklists/requirements.md` (007-25 A) exige números concretos nos requisitos ou só na entrega? [Clarity, checklists/requirements.md § Cobertura; Spec § SC-003] (pendente: placeholders — critério documentado, valores numéricos ausentes)

## Scenario Coverage

- [x] CHK026 — Fluxo primário (base para avaliação com anonimização e balanceamento) está coberto em US1 e FR-001–FR-003? [Coverage, Spec § US1]
- [x] CHK027 — Fluxo de alinhamento ao classificador e testes sem fila está coberto em US2? [Coverage, Spec § US2; 007-22 B, 007-23 A]
- [x] CHK028 — Granularidade log/apontamento e split pós-MVP estão cobertos em US3? [Coverage, Spec § US3]
- [x] CHK029 — Cenários de registro inválido, PII e subconjunto só-catálogo adiado estão em Edge Cases? [Coverage, Spec § Edge Cases]
- [ ] CHK030 — Dependência de ordem entre features do pipeline (ex.: **009-22 A** — 009 espera 008) está documentada onde a base impacta métricas do ecossistema? [Coverage, Gap, tasks/plan] (pendente: 009-22 A não rastreado nos artefatos 007)

## Edge Case Coverage

- [x] CHK031 — Comportamento para rótulo fora do enum e registro incompleto está especificado? [Edge Case, Spec § Edge Cases; 007-13 A]
- [x] CHK032 — Duplicata com rótulos diferentes: requisito de documentação explícita? [Edge Case, Spec § FR-012; 007-12 B]
- [x] CHK033 — Divergência INFO como ocorrência (T7) vs Drive legado está registrada sem apagar requisito da base? [Edge Case, Spec § Divergências]
- [x] CHK034 — Subconjunto só-catálogo sem inventar registros antes de **006-7 D** está tratado? [Edge Case, Spec § FR-011; DN-007-1 C] (pendente externa: 006-7 D para conteúdo — requisito atendido)

## Non-Functional Requirements

- [x] CHK035 — Privacidade (sem PII, anonimização, ADR_064) e ausência de mascaramento LLT-202 nesta base estão nos requisitos? [NFR, Spec § FR-001, FR-010; 007-16 A]
- [x] CHK036 — Sem requisito de volume máximo (007-24 A) documentado? [NFR, Spec § FR-013]
- [x] CHK037 — Desacoplamento da fila Redis para avaliação da base está como restrição? [NFR, Spec § FR-009; plan.md § Constraints]
- [ ] CHK038 — Impacto de **API Módulo Ocorrência** (DN-009-2) sobre requisitos da 007 está delimitado como “nenhum” de forma rastreável? [NFR, plan.md § API do Módulo] (pendente: DN-009-2 — contrato HTTP/JWT ainda sem fonte)

## Dependencies & Assumptions

- [x] CHK039 — Dependência do catálogo (**006-7 D**, ADR_167) para conteúdo só-catálogo está nomeada e não inventada? [Dependency, Spec § FR-011; plan.md § Tarefas bloqueadas] (pendente externa: 006-7 D)
- [x] CHK040 — Substituição dos sintéticos da 006 (007-19 B) e consumo por testes automatizados estão como premissa rastreada? [Dependency, Spec § FR-009; tasks T049–T050]
- [ ] CHK041 — Pendência **N3 A** (ADR AchadoAnalise) está refletida nos requisitos da 007 sem exigir campos de assinatura na base? [Dependency, plan.md § Pendências; 007-17 A] (pendente: N3 A — ADR formal não publicada)
- [ ] CHK042 — Dependências **010-2 A** / **010-27 A** (mascaramento) estão documentadas como ortogonais à base sem PII? [Dependency, plan.md § Pendências; tasks T005/T039] (pendente: 010-2 A, 010-27 A)
- [x] CHK043 — Contrato upstream **log.ingested v1** (T5 D) citado para ids sem exigir consumo de fila na feature? [Assumption, Spec § Decisões T5; plan.md § Constraints]
- [x] CHK044 — Princípio V (007-26 A): lacunas novas bloqueiam plan — **DN-007-3 B** fechada na etapa 3e com fonte, sem invenção pré-decisão? [Assumption, Spec § Clarifications etapa 3e; plan.md § Constitution Check V; respostas-joao.md §Etapa 3c, DN-006-2/007-3 B]

## Ambiguities & Conflicts

- [x] CHK045 — Tabela **Divergências documentadas** cobre INFO, T1, catálogo, métricas e contrato v1? [Traceability, Spec § Divergências]
- [x] CHK046 — Decisões etapa 3b/3c (007-meta, DN-007-1/2, formato, caminho) têm IDs nas seções Clarifications e Decisões tomadas? [Traceability, Spec § Clarifications; § Decisões tomadas]
- [x] CHK047 — Tasks citam fonte `(fonte: respostas-joao.md, <ID>)` de forma consistente com spec? [Traceability, tasks.md cabeçalho Notes]
- [ ] CHK048 — Resolução formal do conflito **constituição I × T1 revisto** consta como requisito de emenda, não como decisão fechada? [Conflict, plan.md § Pendências emenda princípio I] (pendente: emenda do princípio I)
- [x] CHK049 — **DN-007-3 B** (recall/F1 por classe) está rastreada em Clarifications, Decisões tomadas e Dúvidas novas (respondida), com fonte etapa 3e? [Traceability, Spec § Session etapa 3e; § Decisões tomadas DN-007-3]

## Rastreabilidade tasks ↔ requisitos

- [x] CHK050 — Cada fase de tasks mapeia FR/US/SC correspondentes nas descrições? [Traceability, tasks.md por fase]
- [x] CHK051 — Tarefas BLOQUEADAS (**T034**, **T054**) citam apenas **006-7 D** / conteúdo só-catálogo; **T043** desbloqueada na etapa 3e (**DN-007-3 B**), sem bloqueio documental residual? [Consistency, tasks.md § Bloqueios explícitos; § Desbloqueado etapa 3e]
- [ ] CHK052 — Tarefa **T056** (rastreio emenda princípio I) está coerente com estado atual do plan (já listado) — requisito de revisão gestor satisfeito? [Traceability, tasks T056; plan.md § Pendências] (pendente: emenda do princípio I ainda aberta fora da pasta)

## Notes

- **2026-11-16 etapa 3e/3f (reavaliação)**: **CHK006**, **CHK014**, **CHK023** e **CHK051** passam a `[x]` após **DN-007-3 B** e **DN-007-4 A** (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h01 e 11h; DN-007-4 A: §Etapa 3c (João, 16/11 11h02)). **T043** desbloqueada em `tasks.md`.
- **2026-11-16 analyze2 (I1)**: `checklists/requirements.md` — lacunas equivalentes remanescentes: **CHK013** (split), **006-7 D** (só-catálogo); métricas e gold apontamento fechados na etapa 3e.
- `/speckit-implement` pode ler este arquivo como gate; marcadores desta reavaliação com base em leitura integral de spec.md, plan.md e tasks.md (fontes prioritárias João/etapa 3e).
- Pendências externas / requisitos ainda abertos neste checklist: DN-009-2, N3 A, 010-2 A, 010-27 A, 006-7 D, 009-22 A, emenda do princípio I; **CHK013** (split qualitativo sem proporção/seed/estratificação — 007-20 B).