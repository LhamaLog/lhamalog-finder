# Quickstart: Validação da base de logs rotulados

**Feature**: `007-base-logs-rotulados` | **Date**: 2026-11-16

Guia de **validação** (não implementação). Não executa comandos neste documento — descreve resultados esperados quando os artefatos e testes existirem.

## Pré-requisitos

- Go 1.27 via Docker, conforme [README.md](../../../README.md)
- [spec.md](../spec.md) (decisões 007-meta, DN-007-2 B, DN-006-2/007-3 B, DN-007-4 A, **007 H2 D**, 007-formato, 007-caminho, DN-007-1 C)
- Contratos: [registro-rotulado.md](./contracts/registro-rotulado.md), [entrada-classificador-desde-base.md](./contracts/entrada-classificador-desde-base.md)
- Modelo: [data-model.md](../data-model.md)
- Classificador: [spec 006 `classificar-log.md`](../006-modelo-classificacao-logs/contracts/classificar-log.md)

## 1. Declaração de origem e privacidade (SC-001)

**Objetivo**: User Story 1 — anonimização; zero PII.

1. Abrir o arquivo de declaração de origem em **`datasets/labeled-logs/`** na raiz do repositório (007-caminho C; registros em **JSONL**, 007-formato A).
2. Confirmar procedimento de anonimização e ausência de dados reais de clientes.
3. Amostrar ou varrer o conjunto em busca de PII (e-mail, senha, token) — esperado **0** ocorrências.

**Esperado**: 100% dos registros cobertos pela declaração (SC-001).

## 2. Taxonomia e balanceamento (SC-002, SC-003)

**Objetivo**: ERROR/WARNING/INFO; contagens iguais no mínimo; total em milhares.

1. Ler manifest de contagens publicado (spec § Cobertura alvo + `checklists/requirements.md`).
2. Conferir que nenhum `class_label` é spam ou "não classificado".
3. Conferir mínimo **igual** nas três classes e total na ordem de milhares.

**Esperado**: alinhado a 007-1 C, 007-2 B, 007-25 A.

## 3. Formato e mapeamento (SC-004)

**Objetivo**: 100% dos registros com campos FR-004.

1. Executar validador estrutural (teste automatizado ou script de CI quando existir).
2. Conferir tabela de mapeamento na spec com nomes físicos reais dos arquivos (007-6 D).

**Esperado**: 0 falhas em 100% dos registros (007-21 A).

**Comando** (quando existir): `make test` incluindo pacote de fixtures da 007 — alvo em `tasks.md`.

## 4. Entrada no classificador sem fila (SC-004 + 007-23 A)

**Objetivo**: User Story 2.

1. Para cada registro `itemType=log`, mapear para `Classificar log` conforme [entrada-classificador-desde-base.md](./contracts/entrada-classificador-desde-base.md).
2. Invocar classificador (implementação spec 006).
3. Comparar `class_label` (gold) com `classe` (predição).

**Esperado**: operação completa sem Redis e sem OpenObserve para obter texto.

## 5. Itens apontamento (User Story 3)

1. Listar registros `itemType=apontamento`.
2. Confirmar `documentacao` presente e `members` não vazio.
3. Conferir **exatamente um** `class_label` no apontamento e **ausência** de `class_label` em cada `member` (DN-007-4 A).
4. Para cada `member`, invocar `Classificar log` e reunir predições `classe`.
5. Derivar **classe prevista do apontamento** = **mais grave** entre predições (**ERROR** > **WARNING** > **INFO**, ADR_170); comparar ao `class_label` gold do agregado (igualdade exata). **Unidade** da métrica = apontamento (**007 H2 D**; FR-007c; [entrada-classificador-desde-base.md](./contracts/entrada-classificador-desde-base.md) § Registro de apontamento).
6. Executar cenário de teste documentado por item (sem inventar agrupamento 008).

**Esperado**: papéis distintos de log único documentados (007-7 C); gold **somente** no agregado (DN-007-4 A; FR-007b); avaliação conforme **007 H2 D** (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40); respostas-etapa3f.md §D).

## 6. Split e uso MVP (SC-006)

1. Abrir `SPLIT.md` / manifest de split (007-20 B; analyze M2).
2. Confirmar **nomes** de partições e critérios **qualitativos** (FR-005); MUST NOT exigir proporção/seed numéricos sem fonte.
3. Confirmar documentação: MVP **somente avaliação** (007-3 A).
4. Confirmar partição `validacao` usada para métricas de precisão (contrato E3).

**Esperado**: split nomeado; sem obrigação de treino no MVP.

## 7. Métricas por classe ≥ 95% (SC-005)

**Objetivo**: LLT-199 critério 2 / 007-meta / 006-2 B.

1. Calcular **precisão, recall e F1 por classe** (ERROR, WARNING, INFO) na partição de **validação** — cada métrica MUST atingir **≥ 95%** em cada classe (007-meta; DN-007-2 B; DN-006-2/007-3 B). Itens `log`: um par gold/predição por registro. Itens `apontamento`: contagem **por apontamento** após **007 H2 D** (prevista = mais grave entre predições dos `members`).
2. Registrar resultados no relatório de avaliação acordado na spec 006 quando aplicável.

**Esperado**: precisão, recall e F1 **≥ 95% em cada classe** (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h; respostas-etapa3c.md §1).

## 8. Subconjunto só-catálogo

1. Verificar manifest/critério de inclusão (FR-011).
2. **Não** exigir registros de exemplo até catálogo existir (DN-007-1 C).

**Esperado**: zero registros inventados na partição `so_catalogo` se catálogo ausente.

## 9. Substituição dos sintéticos 006 (007-19 B)

1. Desativar ou remover dependência de dados sintéticos legados da spec 006 nos testes de métricas.
2. Apontar testes para este conjunto.

## 10. Volume sem teto (FR-013 / SC implícito)

1. Confirmar na spec/plan que **não** há requisito de volume máximo (007-24 A; FR-013).
2. Crescimento do conjunto fica a critério de engenharia e revisão de PR — sem assert de teto em CI salvo política futura com fonte.

**Esperado**: documentação satisfaz FR-013 (analyze M6; [research.md](../research.md) R18).

## O que não validar nesta feature

- API do Módulo Ocorrência (DN-009-2; [plan.md](../plan.md) § API do Módulo)
- Normalização de mensagem para assinatura (DN-008-4 B decidido na etapa 3b; implementação na spec **008**, não na base 007)
- Conteúdo do subconjunto só-catálogo sem catálogo (DN-007-1 C)
- Consumo da fila `log.ingested` (007-23 A)
- Treino de modelo / LLM (007-3 A; ADR_167)

## Referências rápidas

| ID | Uso na validação |
|----|------------------|
| 007-meta / DN-006-2/007-3 B | Cenário 7 |
| DN-007-1 C | Cenário 8 |
| 007-19 B | Cenário 9 |
| 006-2 B | Métricas por classe |
| DN-007-4 A | Cenário 5 (apontamento) |
| 007 H2 D | Cenários 5 e 7 (avaliação apontamento) |
| 007-24 A / FR-013 | Cenário 10 |

(fonte: respostas-joao.md; detalhes em [research.md](../research.md))