# Research: Base de logs rotulados para treino e avaliação

**Feature**: `007-base-logs-rotulados` | **Date**: 2026-11-16

Consolidação da Fase 0. Cada item cita fonte primária (João, revisão T1, respostas documentadas, contrato 08/11 ou artefato do repositório).

---

## R1 — Objetivo e uso no MVP

- **Decision**: A base serve **somente para avaliação das regras** no MVP por regras; **não** treinar modelo nem usar LLM (ADR_167).
- **Rationale**: 007-3 A; alinhamento à spec 006 (métricas offline, SC-007).
- **Alternatives considered**: Treino e avaliação já no MVP (007-3 C) — **rejeitado**.
- **Fonte**: (fonte: respostas-joao.md, 007-3 A; origem: ADR_167; Jira LLT-198)

---

## R2 — Escala e balanceamento por classe

- **Decision**: Ordem de **milhares** de registros; **mesma contagem mínima** para ERROR, WARNING e INFO; contagens publicadas na spec e no checklist (007-25 A).
- **Rationale**: Evita viés de classe fraca na validação; SC-003.
- **Alternatives considered**: Mínimos distintos por classe (007-2 C) — **rejeitado**.
- **Fonte**: (fonte: respostas-joao.md, 007-1 C / 007-2 B / 007-25 A; origem: perguntas-etapa2b.md)

---

## R3 — Local e versionamento

- **Decision**: Conjunto em **arquivos versionados no repositório Git** (não Drive isolado, não só fixtures efêmeras).
- **Rationale**: 007-4 A; rastreabilidade com testes automatizados (007-22 B).
- **Alternatives considered**: Dataset só no Drive (007-4 C) — **rejeitado**.
- **Fonte**: (fonte: respostas-joao.md, 007-4 A; origem: perguntas-etapa2b.md 007-4)

- **Decision (etapa 3c)**: Formato **JSONL**; caminho **`datasets/labeled-logs/`** na raiz do repositório `lhamalog-finder`.
- **Rationale**: 007-formato A; 007-caminho C (`respostas-joao.md` §Etapa 3b, João 16/11 09h02).
- **Fonte**: (fonte: respostas-joao.md §Etapa 3b; origem: 007 formato A / 007 caminho C; 007-4 A)

Nomes de arquivos `.jsonl` e atributos JSON físicos seguem **007-6 D** (atualizar tabela na spec quando existirem).

---

## R4 — Campos por registro e mapeamento lógico

- **Decision**: Por registro: `level_text`, `message_text`, `lhamalog_batch_id`, `lhamalog_record_id`, `tenantId`, `productId`, `class_label`; mapeamento físico documentado na spec (tabela 007-6).
- **Rationale**: Entrada alinhada ao classificador 006 (nível + mensagem já separados; spec 006 research “Campo e mensagem já separados”); ids espelham contrato OpenObserve §7.
- **Alternatives considered**: Só nível+mensagem (007-5 A) — **rejeitado**; documento OTel completo (007-5 C) — **rejeitado**.
- **Fonte**: (fonte: respostas-joao.md, 007-5 B / 007-6 D / 007-8 A; origem: contrato-log-ingested-v1.txt §4, §7)

---

## R5 — Taxonomia de rótulos

- **Decision**: `class_label` ∈ {**ERROR**, **WARNING**, **INFO**}; sem spam nem “não classificado”.
- **Rationale**: 007-9 A; constituição III via ADR_170; coerência LLT-199 critério 1 após decisão João.
- **Alternatives considered**: Enum legado com SPAM (ADR_019) — **rejeitado** para esta base.
- **Fonte**: (fonte: respostas-joao.md, 007-9 A; origem: ADR_170; constituição III)

---

## R6 — Origem, privacidade e auditoria

- **Decision**: Origem por **anonimização de logs reais** (não fictícios puros); **declaração de origem**; **nenhuma PII**, nem fictícia (sem casos LLT-202 na base).
- **Rationale**: 007-14 B prevalece sobre texto aberto do Jira LLT-198; ADR_064; SC-001.
- **Alternatives considered**: Dados fictícios (Jira literal) — **rejeitado** pelo João.
- **Fonte**: (fonte: respostas-joao.md, 007-14 B / 007-15 A / 007-16 A; origem: Jira LLT-198; ADR_064)

---

## R7 — Granularidade: log único vs apontamento

- **Decision**: Dois tipos de item documentados: **registro de log** (um par nível+mensagem) e **registro de apontamento** (vários logs agregados para cenários de achado — T2 C).
- **Rationale**: 007-7 C; prepara testes que não passam pela fila mas espelham agrupamento futuro (spec 008).
- **Alternatives considered**: Só log único (007-7 A) — **rejeitado**.
- **Fonte**: (fonte: respostas-joao.md, 007-7 C; origem: T2 C)

---

## R8 — Consumo e relação com a spec 006

- **Decision**: Consumidor principal: **testes automatizados**; entrada **desacoplada da fila**; esta base **substitui** sintéticos de prova da 006 (007-19 B).
- **Rationale**: FR-009; 006 quickstart cenário 7.
- **Alternatives considered**: Avaliação só manual — **rejeitado** (007-22 B).
- **Fonte**: (fonte: respostas-joao.md, 007-19 B / 007-22 B / 007-23 A; origem: spec 006 plan.md / research.md)

Contrato de ponte: [contracts/entrada-classificador-desde-base.md](./contracts/entrada-classificador-desde-base.md) → [spec 006 `classificar-log.md`](../006-modelo-classificacao-logs/contracts/classificar-log.md).

---

## R9 — Split treino/validação/teste

- **Decision**: **Split documentado** (nomes de partições + **critérios qualitativos** em `SPLIT.md`) para uso **pós-MVP**; MVP usa partições relevantes à **avaliação** (partição `validacao` para **007-meta** / SC-005 — fonte: contracts `entrada-classificador-desde-base.md` E3). **Sem** proporção, seed ou estratificação numérica nas fontes — não inventar (analyze M2; spec FR-005).
- **Rationale**: 007-20 B; 007-3 A limita treino no MVP; 007-26 A — lacuna não bloqueia o plano inteiro.
- **Alternatives considered**: Sem split — **rejeitado**; fixar 70/15/15 sem fonte — **rejeitado** (princípio V).
- **Fonte**: (fonte: respostas-joao.md, 007-20 B; origem: perguntas-etapa2b.md; spec.md FR-005)

---

## R10 — Meta de precisão 95% (007-meta, DN-007-2)

- **Decision**: Alvo **≥ 95% de precisão em cada classe** ERROR, WARNING e INFO na base de validação (LLT-199 critério 2).
- **Rationale**: 007-meta (João 16/11 00h42) + **DN-007-2 B** (João 16/11 09h02); alinha a **006-2 B** (métricas por classe).
- **Alternatives considered**: Precisão global agregada — **rejeitada** (DN-007-2 B).
- **Fonte**: (fonte: respostas-joao.md §Etapa 3 e §Etapa 3b; origem: Jira LLT-199 critério 2; DN-007-2 B)

- **Decision (etapa 3e — recall/F1)**: **≥ 95% de recall e ≥ 95% de F1 em cada classe** ERROR, WARNING e INFO na partição de validação, simétrico à precisão (DN-006-2/007-3 B).
- **Rationale**: **006-2 B** mede precisão, recall e F1 por classe; João fixa o mesmo alvo numérico que a precisão (etapa 3c). Especialista marcou **SEM FONTE** em Jira/Drive para recall/F1 (`respostas-etapa3c.md` §1) — decisão João prevalece; divergência Jira/Drive × João em [spec.md](./spec.md) § Divergências (**007-meta × 006-2 B**).
- **Alternatives considered**: A) só reportar recall/F1 sem meta; C) metas menores; D) decidir na LLT-51 — **rejeitadas** (opção B).
- **Fonte**: (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h; origem: DN-006-2/007-3 B; 006-2 B)

---

## R11 — Subconjunto só-catálogo

- **Decision**: **Esperar o catálogo** para montar registros desse subconjunto; documentar critério de inclusão alinhado à 006 US3.
- **Rationale**: DN-007-1 C; ADR_167 [PENDENTE]; 006-7 D.
- **Alternatives considered**: Subconjunto vazio com placeholder (A) — **rejeitado** em favor de C.
- **Fonte**: (fonte: respostas-joao.md, §Etapa 3 DN-007-1 C; origem: respostas-etapa3.md §2)

---

## R12 — Validação de qualidade do conjunto

- **Decision**: Validar rótulo por **autoria + método documentado** (sem revisão 100% humana obrigatória); rejeitar registro sem nível ou mensagem; duplicatas com rótulos diferentes **permitidas** se documentadas; verificação de formato em **100%** dos registros (007-21 A).
- **Rationale**: SC-004; edge cases spec.
- **Fonte**: (fonte: respostas-joao.md, 007-11 D / 007-12 B / 007-13 A / 007-21 A; origem: spec 006 FR-006)

---

## R13 — Contrato log.ingested v1

- **Decision**: Referência normativa para ids e tenant/produto: **somente v1** válida (`schemaVersion` `"1"`).
- **Rationale**: T5 D; contrato 08/11 §2–4; alinhamento com `internal/domain/envelope.go` (outras features).
- **Fonte**: (fonte: respostas-joao.md, T5 D; origem: contrato-log-ingested-v1.txt §2–4)

---

## R14 — Rotulação e escopo do rótulo

- **Decision**: Rotulação pela **engenharia do Finder**; escopo **só classe** — sem assinatura LLT-200.
- **Rationale**: 007-10 B, 007-17 A.
- **Fonte**: (fonte: respostas-joao.md, 007-10 B / 007-17 A; origem: Jira LLT-198; LLT-200 fora de escopo)

Implicação: **DN-008-4** (normalização para assinatura na 008) **não** bloqueia montagem da base 007.

---

## R15 — Métricas de avaliação (006-2 B)

- **Decision**: Quando a base for consumida, reportar **precisão, recall e F1 por classe** (fora do escopo de implementação automática obrigatória nesta entrega, salvo testes acordados em tasks). Registro `log`: um par gold/predição por item. Registro `apontamento`: seguir **007 H2 D** (R20) — unidade de contagem na métrica = apontamento após agregar predições dos `members`.
- **Rationale**: 006-2 B; spec.md FR-008, **FR-007c**.
- **Fonte**: (fonte: respostas-joao.md, 006-2 B; origem: Jira LLT-197/199); (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), 007 H2 D; origem: ADR_170)

Aceite na partição `validacao`: precisão, recall e F1 **≥ 95% por classe** (**007-meta**; DN-007-2 B; DN-006-2/007-3 B).

---

## R16 — Serialização JSONL no repositório Go

- **Decision**: Persistir cada registro como **um objeto JSON por linha** (JSONL) em `datasets/labeled-logs/`; leitura/validação via **`encoding/json`** (biblioteca padrão Go) linha a linha — sem dependência nova em `go.mod`.
- **Rationale**: 007-formato A, 007-caminho C (etapa 3c); alinhamento a `go test` / `make test` (README § Comandos).
- **Alternatives considered**: Parquet ou CSV (não ofertados nas fontes João) — **rejeitado**.
- **Fonte**: (fonte: respostas-joao.md §Etapa 3b; origem: 007 formato A / 007 caminho C); (fonte: `go.mod`, linha 1–8; README § Layout)

---

## R17 — Normalização DN-008-4 (escopo 008, referência cruzada)

- **Decision**: A regra **DN-008-4 B** (trim, case-insensitive, marcadores para números/UUIDs/IPs/datas/hex) aplica-se à **assinatura de agrupamento** na spec **008**, **não** à montagem nem ao esquema gold da base 007 (007-17 A).
- **Rationale**: Resposta João 16/11 09h02; evita duplicar normalizador na base quando o rótulo é só classe.
- **Fonte**: (fonte: respostas-joao.md §Etapa 3b, DN-008-4 B; origem: spec 008)

---

## R18 — Volume máximo (FR-013)

- **Decision**: **Nenhum** requisito de volume ou escala máximo para o conjunto; crescimento limitado apenas por prática de engenharia e revisão de PR.
- **Rationale**: 007-24 A; FR-013 / SC não exigem teto — satisfação **documental** (analyze M6).
- **Alternatives considered**: Cap fixo (ex. 100k linhas) — **rejeitado** (sem fonte).
- **Fonte**: (fonte: respostas-joao.md, 007-24 A; origem: perguntas-etapa2b.md; spec.md FR-013)

---

## R19 — Gold em apontamento (DN-007-4)

- **Decision**: **Um** `class_label` gold **por apontamento** (nível agregado); cada `member` MUST NOT incluir `class_label` (DN-007-4 A).
- **Rationale**: 007-7 C exige os dois `itemType` documentados; gold no agregado evita modo dual. Jira LLT-199 critério 1 («todo log») fala da **classificação**, não do formato da base — pista não normativa (`respostas-etapa3c.md` §3). **DN-007-4 A** fechada na etapa 3e (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02)); divergência com Jira em [spec.md](./spec.md) § Divergências (**DN-007-4 × LLT-199 crit. 1**).
- **Alternatives considered**: B) gold por membro; C) modos coexistindo — **rejeitadas** (DN-007-4 A).
- **Fonte**: (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); origem: DN-007-4 A; 007-7 C; respostas-etapa3c.md §3)

---

## R20 — Avaliação de apontamento: classe prevista = mais grave (007 H2 D)

- **Decision**: Para `itemType=apontamento`, obter `classe` predita **por cada** `member` (classificação por log, ADR_167); a **classe prevista do apontamento** MUST ser a **mais grave** entre essas predições (**ERROR** > **WARNING** > **INFO**, ADR_170); MUST comparar essa prevista ao único `class_label` gold do apontamento (igualdade exata). **Unidade** da avaliação/métrica = **apontamento** (não por log isolado na contagem).
- **Rationale**: Coerência com a regra de classe da ocorrência quando há divergência entre logs (clarificação 09/11 citada pela ADR_170); alinha métricas à produção sem alterar gold agregado (DN-007-4 A). Exemplo: gold ERROR + predições ERROR, ERROR, WARNING → prevista ERROR → acerto. Contraexemplo: gold WARNING + predições WARNING, WARNING, ERROR → prevista ERROR → erro (falso positivo para ERROR) (`respostas-etapa3f.md` §D).
- **Alternatives considered**: A) unanimidade; B) maioria; C) por log (cada membro herda gold do apontamento) — **rejeitadas** (João escolheu D; fonte: respostas-etapa3f.md §Opções).
- **Fonte**: (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), 007 H2 D; origem: ADR_170; respostas-etapa3f.md §D; DN-007-4 A)

## Evidência desta execução

Origem e autoria não foram fornecidas; o usuário informou não haver base disponível para esta execução. Não foi produzida declaração de anonimização nem contagem real. ORIGEM.md é somente formulário pendente, sem atestação.

007-19 B: a futura entrega T020 deve substituir dados sintéticos usados como prova pela 006 (coordenação T050/006-10 A, sem editar Jira). Até lá, fixtures verificam código, sem prova de precisão de produto.

message_text da base deve ser sem PII (007-16 A); a heurística de rejeição do loader não certifica anonimização. P-010-2/P-010-27 são dependências ortogonais de mascaramento runtime. DN-008-4 B não bloqueia montagem da base (007-17 A). R18/FR-013: sem máximo de volume documental; nenhum limite arbitrário de corpus introduzido.
