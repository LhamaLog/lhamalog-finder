# Contract: ResumoAnalise com achados agrupados

Estende o handoff da spec **006** com campos de agrupamento desta feature. Base comum: [resumo-analise-handoff.md](../../006-modelo-classificacao-logs/contracts/resumo-analise-handoff.md).

## Escopo 008

- O Finder MUST enviar **exatamente um** envelope `ResumoAnalise` **por lote** processado com achados (**DN-008-7**, ADR_062).
- Antes do envio, MUST **consultar** o Módulo **uma vez por achado** para fusão entre lotes (**N2 B**, **DN-008-7 B**).
- Achados MUST refletir agregação **no lote** e valores **após** consulta (contagem e extremos globais quando aplicável).
- O Módulo Ocorrência MUST **criar ou atualizar** pela chave **(assinatura, classe, produto, tenant)** sem duplicata (FR-004, FR-005).

(fonte: respostas-joao.md §Etapa 3c 11h, DN-008-7; §Etapa 3c (João, 16/11 11h02), DN-008-7 B; origem: ADR_062, N2 B)

## Campos adicionais por achado (N3 A — [PENDENTE ADR])

| Campo | Obrigatório | Descrição |
|-------|-------------|-----------|
| `assinatura` | sim* | *Após ADR N3 |
| `contagem` | sim* | Total após fusão (**008-12**) |
| `primeiraOcorrencia` | sim* | Timestamp log (**008-14**) |
| `ultimaOcorrencia` | sim* | Timestamp log |
| `recorrente` | sim | Preenchido via OpenObserve (**008-5**); falha consulta segue **DN-008-3** (rev. 3h: 3 tentativas, **sem DLQ**; aviso Core **batchId** + **motivo**) |

Até a ADR vigente, implementação MUST tratar estes campos como **planejados** (SC-005).

## Idempotência

`batchId` + `execId` (**N4 A**). Origem de `execId`: **DN-006-1** (spec 006).

## Transporte — **DN-009-2 [PENDENTE]**

HTTP interno + JWT (**N1 A**); rota, respostas 2xx/4xx/5xx e timeout **sem fonte** (respostas-etapa3.md §10). Envio real **BLOQUEADO**; montagem em memória **desbloqueada**.

## Pós-condições (Módulo — fora do Finder)

| Efeito | Fonte |
|--------|-------|
| Upsert ocorrência pela chave | 008-11, FR-004 |
| FECHADA + novo log → ABERTA | 008-4 B |
| Contador sem mudança de status | 008-10 A |
| Aceitar achados INFO | T7 A |

## O que permanece fora do escopo

- Política ADR_035 / notificação (**008-17 C**)
- Publicação evento Core (spec 009)
- Gravidade no payload Finder (ADR_062)

## Fontes

- (fonte: respostas-joao.md, N2 B / N3 A / N4 A / N1 A / T7 A / 008-4 B / 008-10 A / 008-11 rev.; origem: ADR_062, ADR_040)
- (fonte: spec.md, FR-001, FR-009, DN-008-3 rev. 3h / DLQ B; respostas-joao.md §Etapa 3g (João, 16/11 13h22))

## Evidência desta execução

N3 A **[PENDENTE ADR]**: assinatura, contagem, primeira/última ocorrência e versão existem no modelo lógico comum 006 e são excluídos do DTO JSON provisório 009. Sem fechamento de nomes externos. Um resumo por lote, batchId com execId opcional; sem geração de execId neste caminho.
