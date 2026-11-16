# Specification Quality Checklist: Base de logs rotulados para treino e avaliação

**Purpose**: Validate specification completeness and quality before proceeding to planning
**Created**: 2026-11-14
**Feature**: [spec.md](../spec.md)

## Content Quality

- [x] No implementation details (languages, frameworks, APIs)
- [x] Focused on user value and business needs
- [x] Written for non-technical stakeholders
- [x] All mandatory sections completed

## Requirement Completeness

- [ ] No [NEEDS CLARIFICATION] markers remain — **spec.md** sem marcadores literais; lacuna equivalente: split (proporção/seed/estratificação; `revisao-etapa3c.md` **CHK013**)
- [x] Requirements are testable and unambiguous — ressalva **FR-005** (split só critérios qualitativos; 007-20 B); métricas recall/F1 (**DN-007-3 B**, etapa 3e); gold em apontamento (**DN-007-4 A**, etapa 3e); avaliação apontamento (**007 H2 D** / **FR-007c**, etapa 3g)
- [x] Success criteria are measurable — **SC-005**: precisão, recall e F1 **≥95% por classe** (**DN-007-2 B**, **DN-007-3 B**)
- [x] Success criteria are technology-agnostic (no implementation details)
- [x] All acceptance scenarios are defined
- [x] Edge cases are identified
- [x] Scope is clearly bounded
- [x] Dependencies and assumptions identified

## Feature Readiness

- [ ] All functional requirements have clear acceptance criteria — parte bloqueada: conteúdo só-catálogo (**006-7 D** / **DN-007-1 C**); **FR-007**/**FR-007b**/**FR-007c** e gold + regra de avaliação apontamento fechados (etapas 3e/3g)
- [x] User scenarios cover primary flows
- [x] Feature meets measurable outcomes defined in Success Criteria — **SC-005** completo para precisão/recall/F1 por classe (etapa 3e); contagens numéricas finais de volume só na entrega (**CHK010**, **CHK025**)
- [x] No implementation details leak into specification

## Cobertura da base (007-25 A — preencher na entrega)

| Classe | Contagem mínima (igual entre classes) | Contagem entregue |
|--------|--------------------------------------|-------------------|
| ERROR | _A definir na entrega (ordem milhares no total)_ | |
| WARNING | _idem_ | |
| INFO | _idem_ | |
| **Total** | _Ordem de milhares (007-1 C)_ | |

## Notes

- Iteração 2026-11-16: clarify não interativo; respostas 007-1 a 007-26 de `respostas-joao.md` aplicadas; transversais T1 (revisto), T2–T7 e taxonomia ADR_170.
- Iteração 2026-11-16 etapa 3b: **007-meta** 95% de precisão (João); **DN-007-1 C** esperar catálogo; checklist revalidado após clarify.
- Iteração 2026-11-16 etapa 3c: **DN-007-2 B** precisão **por classe**; **007-formato** JSONL; **007-caminho** `datasets/labeled-logs/`; checklist revalidado após clarify.
- Iteração 2026-11-16 etapa 3c — correções analyze: **Status** pós-plan; **C1** emenda princípio I; **M3** volume sem número fixo; **DN-007-4** gold em apontamento; checklist revalidado (spec sem marcadores `[NEEDS CLARIFICATION]`).
- Iteração 2026-11-16 etapa 3c — **analyze2 (I1)**: itens de completude/mensurabilidade **desmarcados** onde há pendência documentada em `checklists/revisao-etapa3c.md` (**CHK013** split; histórico antes da etapa 3e).
- Iteração 2026-11-16 etapa 3e: **DN-007-3 B** recall/F1 **≥95% por classe**; **DN-007-4 A** gold **um por apontamento**; checklist revalidado após clarify não interativo.
- Iteração 2026-11-16 etapa 3f: **DN-007-3 B**, **DN-007-4 A** e **T043** (sincronizar LLT-199 crit. 2 no checklist) **não** são lacuna/bloqueio documental; fonte canônica **DN-007-4 A** = `respostas-joao.md §Etapa 3c (João, 16/11 11h02)` (fonte: analyze-007-etapa3e.md I1/I2; respostas-joao.md linha 39 SUPERADO).
- Iteração 2026-11-16 etapa 3g: **007 H2 D** — classe prevista do apontamento = mais grave entre predições dos `members` (ADR_170); unidade de avaliação = apontamento; fechado em **FR-007c** e contrato `entrada-classificador-desde-base.md` (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40); analyze-007-etapa3f.md H2).
- Ressalva: **emenda princípio I** — pendência fora desta pasta.
- DN-007-1 respondida (esperar catálogo); subconjunto só-catálogo adiado até 006-7 D; não bloqueia SC-001–SC-004 nem SC-006.
- Gate de qualidade detalhado: `checklists/revisao-etapa3c.md` (revisão PR/gestor); este checklist built-in reflete **readiness de spec** com ressalvas acima (split qualitativo; conteúdo só-catálogo).