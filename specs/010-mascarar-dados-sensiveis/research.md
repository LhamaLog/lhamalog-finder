# Research: Mascarar dados sensíveis antes de gravar a ocorrência

**Feature**: `010-mascarar-dados-sensiveis` | **Date**: 2026-11-16

Consolidação da Fase 0. Cada item cita fonte primária. Nada abaixo foi resolvido por suposição (constituição V; **007-26 A**).

**Tabela canônica de decisões 010-***: [spec.md § Decisões tomadas](./spec.md#decisões-tomadas-com-fonte) — este arquivo guarda *rationale* e alternativas rejeitadas, não duplica a matriz completa (corrige analyze **D1**).

---

## R1 — Catálogo de padrões (interim vs. oficial)

- **Decision**: Comportamento até ADR oficial: **e-mails substituídos**; **linhas** com padrão senha/token/secret **removidas** com marcador; lista **única** Finder + Sugestão (ADR_079). Conteúdo regex = **ADR oficial** — **[PENDENTE]** (010-2 A).
- **Rationale**: ADR_064 decisão + consequências; João escolheu ADR como dono da lista (não spec 010).
- **Alternatives considered**: Listar regex na spec 010 (010-2 opção B) — **rejeitado** (fonte: respostas-joao.md, 010-2 A).
- **Fonte**: (fonte: respostas-documentadas.md, 010-1; respostas-joao.md, 010-2 A; origem: ADR_064, ADR_079)

---

## R2 — Momento e estágio no pipeline

- **Decision**: **Uma vez por lote**, em **`montarContexto()`**, mesma passada que “antes do MotorAnalise” (`validarSeguranca()` única).
- **Rationale**: SSID_007.1 pré-condição; ADR_016 + ADR_064; título LLT-202 “antes de gravar” compatível (010-3).
- **Alternatives considered**: Mascarar só no handoff — **rejeitado** (classificador veria bruto; 010-5 A).
- **Fonte**: (fonte: respostas-documentadas.md, 010-3 A / 010-4 A / 010-15 A; origem: ADR_016, ADR_064, SSID_007.1)

---

## R3 — Entrada e saída do mascaramento

- **Decision**: Unidade **linha a linha** (010-14 C). Classificador e MotorAnalise recebem **linhas mascaradas** (010-5 A). “Trecho guardado” = **exemplos mascarados** em `AchadoAnalise.linhas` → Módulo grava `conteudo` (010-6 rev., 010-25).
- **Rationale**: ADR_011 (linha = log); ADR_040 conversão; T1 revisto.
- **Alternatives considered**: 010-6 opção A (só metadados) — **rejeitado** pelo João (exemplos mascarados B, texto ajustado revisão T1).
- **Fonte**: (fonte: respostas-documentadas.md, 010-5 A / 010-14 C; revisao-T1-drive.md, 010-6 / 010-25; origem: ADR_040, ADR_167)

**Escopo de entrada (etapa 3b)**: **DN-010-2** — mascaramento LhamaLog vale **somente** para **linhas de log**; atributos OTel sensíveis são responsabilidade do **Cliente na origem** (Coletor). Identificadores **traceId**, **correlationId**, **hostName**, **moduloId** na Ocorrência seguem **ADR_040** sem passar pelo catálogo de linha (fonte: respostas-joao.md §Etapa 3 / DN-010-2; origem: respostas-etapa3.md §8, DPA Anexo IV item 6, ADR_064, ADR_040).

---

## R4 — Transformações concretas (interim ADR_064)

- **Decision**: E-mail → **placeholder fixo** (010-11 A; ex. `[email]` na spec, texto final na ADR). Senha/token/secret → linha inteira substituída por **`"log com senha omitido"`** (010-10 A). Fora do catálogo → **permanece bruto** (010-13 A). Outras PII (CPF, telefone) → **fora** até entrarem no catálogo (010-12 C).
- **Rationale**: João fechou marcador e estratégia de e-mail; Drive limita escopo vs. Jira amplo (D2 no spec.md).
- **Alternatives considered**: Fail-closed no lote em falha (010-16 A) — **rejeitado**; segue sem trechos (010-16 C).
- **Fonte**: (fonte: respostas-joao.md, 010-10 A / 010-11 A / 010-16 C; respostas-documentadas.md, 010-12 C / 010-13 A; origem: ADR_064, DPA Anexo II/IV, RIPD R03)

---

## R5 — Contadores e rastreio

- **Decision**: Total `dadosSensiveisIdentificados` nos **contadores** do `ResumoAnalise` e do **resumo do arquivo**, rótulo semântico “dados sensíveis mascarados” (010-17).
- **Rationale**: ADR_064 retorno; ADR_102; revisão T1 alinha com handoff 009.
- **Alternatives considered**: Contador só na ocorrência Finder (legado T1 A) — **rejeitado** pós-revisão T1.
- **Fonte**: (fonte: revisao-T1-drive.md, 010-17; origem: ADR_102, ADR_064)

---

## R6 — Compartilhamento Finder–Sugestão

- **Decision**: **Arquivo versionado** no repositório, lido pelos dois módulos; **mesma versão obrigatória** (010-18 A, 010-19 A). Sugestão **reaplica** regras no prompt (010-20); Finder **produz** texto mascarado — Módulo/Core/Notificação **não** revalidam (010-21 rev.).
- **Rationale**: ADR_079; ADR_167 sem LLM no Finder; paridade evita fuga no prompt LLT-212.
- **Alternatives considered**: Serviço de configuração (010-18 B) — **rejeitado**; cópia manual (C) — **rejeitado**.
- **Fonte**: (fonte: respostas-joao.md, 010-18 A / 010-19 A; respostas-documentadas.md, 010-20; revisao-T1-drive.md, 010-21; origem: ADR_079, ADR_167)

**Path e formato** (etapa **3e**): **DN-010-4 C** — caminho e formato MUST ser declarados na **ADR da lista** (**010-2 A**); pendência **P-010-2** até publicação — **não** fixar nesta pasta (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h, DN-010-4 C; origem: ADR_079, 006-3 C; respostas-etapa3c.md §8). Biblioteca regex Go permanece **NEEDS CLARIFICATION** até ADR.

**Mismatch de versão** (etapa **3e**): **DN-010-3 B** — na **inicialização**, se versão lida ≠ versão **combinada** (config deploy; **010-19 A**), serviço **recusa subir** (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h, DN-010-3 B; origem: ADR_079; respostas-etapa3c.md §7).

---

## R7 — Consumidores externos e bruto

- **Decision**: Telas (E08.F1.3), e-mail **Core → Notificação** (010-22), LLM (**Sugestão**) só veem mascarado coberto pelo catálogo. Aviso HTTP Módulo → Core **sem** texto de log; payload inclui **id + tipo + gravidade** (**009-5 revisto**, 010-23). Bruto **só** OpenObserve; Finder não reexpõe (010-24 C).
- **Rationale**: LLT-202 crit. 2; ADR_146 v2 + contrato 08/11; leitura D3 no spec.md.
- **Alternatives considered**: Evento com amostras mascaradas (010-23 B) — **rejeitado** (009-5 revisto).
- **Fonte**: (fonte: respostas-documentadas.md, 010-20 / 010-22 / 010-24; respostas-joao.md §Etapa 3 009-5 revisto; origem: ADR_165, LLT-218, contrato-log-ingested-v1.txt §2)

**E-mail LLT-218 (etapa 3b)**: **DN-010-1 A** / **DN-009-1 A** — Notificação lê a Ocorrência **pelo id** no MongoDB compartilhado; `conteudo` já mascarado (fonte: respostas-joao.md §Etapa 3; origem: respostas-etapa3.md §7, ADR_054, ADR_040). Template de e-mail no Drive permanece [PENDENTE] — spec Notificação.

---

## R8 — Escopo LLT-51, desempenho e testes

- **Decision**: LLT-51 (minimização ampla) **fora** desta história (010-26 A). **Sem** meta de tempo de mascaramento (010-28 A). Casos de teste SC-001/002: **João**, antes de aceite formal (010-27 A) — spec **não** lista conjunto.
- **Rationale**: Separação LLT-202 vs. LLT-51 no backlog; João veta inventário na spec.
- **Alternatives considered**: Casos derivados só da ADR_064 (010-27 B) — **rejeitado**; base 007 com PII (007-16 A) incompatível.
- **Fonte**: (fonte: respostas-joao.md, 010-27 A / 010-28 A; respostas-documentadas.md, 010-26 A; origem: LLT-202, LLT-51)

---

## R9 — Integração com código e specs irmãs

- **Decision**: Mascaramento na **etapa de processamento** do Finder; spec **003** intacta na captura (FR-012). **006** não persiste trechos; coordenação via `mensagem`/linhas mascaradas e `AchadoAnalise.linhas` (T6 A). Handoff alinhado a spec **009** (`linhas` mascaradas, contador); transporte HTTP **DN-009-2**.
- **Rationale**: Estado atual do repo sem implementação de análise; contratos 006/009 referenciam spec 010.
- **Alternatives considered**: Emendar 006 FR-009 para catálogo — **rejeitado** (010-2 A coloca catálogo na ADR).
- **Fonte**: (fonte: spec.md FR-012, FR-003; revisao-T1-drive.md, 010-25; `internal/application/consumir.go` — sem mascaramento hoje)

---

## R10 — Divergências conscientes (sem emendar ADR_019 / ADR_167 / SSID_007.1)

- **Decision**: INFO gera achado com mesmo pipeline de máscara (T7 A). Taxonomia ERROR/WARNING/INFO (ADR_170) — **não** reabrir ADR_019/SSID_007.1 nesta feature (FR-013).
- **Rationale**: spec.md Divergências D1, D8; João 16/11 00h01.
- **Fonte**: (fonte: respostas-joao.md, T7 A; spec.md FR-013, D1/D8)

---

## R11 — Integração com specs 008/007/009 (etapa 3c)

- **Decision**: **DN-008-4 B** — algoritmo **canônico** de normalização da **mensagem** para assinatura: aparar espaços; case-insensitive; substituir **números**, **UUIDs**, **IPs**, **datas** e **hexadecimais** por **marcadores** (um por categoria, distintos entre categorias). **Dona**: spec **008** (**DN-008-1 B**). O normalizador é **único** para as specs **006**, **008**, **009** e **010** — consumidores **006**/**009**/**010** MUST reutilizar o artefato da **008** **sem divergência**. **Artefato canônico de implementação** (repo Finder): `internal/application/assinatura/normalizar_mensagem.go`, tarefa **T016** em `specs/008-agrupar-logs-ocorrencia/tasks.md` (literais **DN-008-5 A** alinhados em **T048**). Nesta feature: entrada à normalização MUST ser linha **já mascarada** por `validarSeguranca()` (**FR-015**, **010-5 A**); a **010** MUST **invocar** esse pacote — **MUST NOT** criar segundo `normalizar_mensagem` ou reimplementar o algoritmo (**010 H1**, etapa 3h). Literais: **DN-008-5 A** — **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`** (**010** só referência cruzada — distinto de **`"log com senha omitido"`**, 010-10 A).
- **Rationale**: Assinatura estável para agrupamento (ADR_013); mascaramento não pode ser contornado pela normalização em bruto; uma única implementação evita drift entre specs no mesmo monorepo.
- **Alternatives considered**: Definir normalização na 010 — **rejeitado** (dona 008; instrução etapa 3c). Implementar **DN-008-4 B** diretamente na pasta da feature 010 antes da **008** — **rejeitado** (analyze **H1**; dependência explícita de **T016**).
- **010 H1** (etapa 3h): a **010 reaproveita** o normalizador canônico da **008** (`normalizar_mensagem.go`, **T016**). (fonte: respostas-joao.md, 010 H1; origem: §Etapa 3g, João 16/11 13h22, DN-008-4 B, FR-015)
- **Fonte**: (fonte: respostas-joao.md §Etapa 3g, João 16/11 13h22, 010 H1; §Etapa 3b, DN-008-4 B; §Etapa 3c, DN-008-5 A; origem: DN-008-1 B, ADR_013, respostas-etapa3c.md §4, spec.md FR-015; specs/008-agrupar-logs-ocorrencia/tasks.md T016)

- **Decision**: **DN-007-2 B** — meta **≥95% de precisão por classe** (ERROR, WARNING, INFO) na validação da spec **007**; **esta feature (LLT-202) não define gate métrico** — aceite de mascaramento permanece **SC-001**–**SC-004** e FR-001–FR-006; precisão de **classificação** nas specs **006**/**007** (**006-2 B**, **007-meta**, **SC-007** referência).
- **Rationale**: Separação LLT-202 (mascaramento) vs. LLT-199 (modelo).
- **Fonte**: (fonte: respostas-joao.md §Etapa 3b, DN-007-2 B; origem: respostas-etapa3.md §1, LLT-199)

- **Decision**: **DN-009-2** — contrato HTTP Finder → Módulo **ainda** pendente; bloqueia só adaptador HTTP; conteúdo mascarado em `saida-mascarada-handoff.md` **liberado**.
- **Fonte**: (fonte: respostas-joao.md §Etapa 3; respostas-etapa3.md §10)

---

## R12 — Base rotulada spec 007 (documentação cruzada)

- **Decision**: Formato **JSONL** (**007 formato A**); caminho documentado `datasets/labeled-logs/` na raiz do repositório (**007 caminho C**). **Não** criar pasta nem arquivos nesta execução do plano 010.
- **Rationale**: Coerência de pipeline — classificação avaliada na 007 usa entradas alinhadas ao Finder; linhas derivadas de log real passam por mascaramento antes da normalização (**FR-015**).
- **Alternatives considered**: Outro formato/caminho (opções 007) — **rejeitado** pelo João (etapa 3b).
- **Fonte**: (fonte: respostas-joao.md §Etapa 3b, 007 formato A / 007 caminho C; origem: spec 007, 007-4 A versionados no repo)

---

## R13 — Catálogo: versão e localização (etapa 3e)

- **Decision**: **DN-010-3 B** — serviço recusa subir se versão lida ≠ versão combinada (inicialização). **DN-010-4 C** — caminho e formato do artefato só na ADR (**010-2 A**); **P-010-2** até ADR publicar.
- **Rationale**: Paridade obrigatória Finder–Sugestão (010-19 A); evita processar com catálogo divergente; path centralizado na ADR compartilhada (Sugestão fora do repo — respostas-etapa3c.md §8).
- **Alternatives considered**: DN-010-3 A/C (parar lote / seguir com alerta) — **rejeitadas** (João **B**); DN-010-4 A/B (path fixo no Finder) — **rejeitadas** (João **C**).
- **Fonte**: (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h, DN-010-3 B / DN-010-4 C; origem: ADR_079, 010-19 A; respostas-etapa3c.md §§7–8)

---

## Itens não resolvidos (registro)

| ID | Assunto | Onde |
|----|---------|------|
| **P-010-2** | Regex + **path/formato** na ADR | plan.md § Pendências; **DN-010-4 C** fixou *dono*, não valores |
| **P-010-27** | Casos de teste João | plan.md § Pendências |
| **DN-009-2** | API HTTP Finder → Módulo | plan.md § API do Módulo |
| Biblioteca regex Go | Implementação `validarSeguranca()` | plan.md Technical Context |
| Placeholder e-mail exato | String canônica | **P-010-2** / ADR |

**DN-010-1**, **DN-010-2**, **DN-007-2**, **DN-008-4**, **DN-010-3**, **DN-010-4**, **DN-008-5**, **010 H1**: **respondidas** (etapas 3b/3c/3e/3h) — ver R6, R11, R13 e spec.md Clarifications etapas 3e/3h.

Nenhum item acima **bloqueia** o desenho do contrato interno `validarSeguranca()` e obrigações de saída mascarada (**SC-005** bloqueia aceite **formal**, não o plano).

## Evidência desta execução

R6: paridade de artefato Finder–Sugestão não confirmada; nenhum contato externo enviado. R7: Finder entrega resumo seguro; id+tipo+gravidade no HTTP Módulo→Core é obrigação 009-5 revisto/010-23 do Módulo; LLT-218 e leitura da Ocorrência por Notificação são externos. R8: nenhuma minimização ampla LLT-51, nem mudança ADR_019/ADR_167/SSID_007.1.

R12: precisão ≥95% por classe é gate da 006/007 sobre JSONL em datasets/labeled-logs/; a 010 não cria base nem mede classificação. P-010-27 e P-010-2 ainda bloqueiam o aceite formal. Constituição I e coordenação 006-10 A/LLT-202 permanecem fora desta implementação.
