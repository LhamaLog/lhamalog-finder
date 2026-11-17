# Contract: Consultar ocorrências existentes (Finder → Módulo Ocorrência)

Interface **externa** para agregação **entre lotes** (**N2 B**). Transporte alinhado a **N1 A** (HTTP interno + JWT). **BLOQUEADO** para implementação até **DN-009-2** (respostas-joao.md §Etapa 3; respostas-etapa3.md §10).

## Propósito

Antes de enviar o `ResumoAnalise`, o Finder MUST obter estado agregado de ocorrências que compartilham a **chave de unicidade** com cada achado do lote — **uma chamada por achado** (**DN-008-7 B**, etapa 3e; spec.md FR-009).

## Chave de consulta

| Campo | Descrição |
|-------|-----------|
| `assinatura` | Valor calculado (DN-008-1 B) |
| `classe` | ERROR \| WARNING \| INFO |
| `productId` | Produto |
| `tenantId` | Tenant |

(fonte: respostas-joao.md, 008-11 rev.; spec.md FR-005)

## Semântica esperada (lógica — independente do transporte)

| Cenário | Comportamento no Finder |
|---------|-------------------------|
| Nenhuma ocorrência | `contagem` do achado = `contagemLote`; extremos = do lote |
| Ocorrência existente ABERTA/EM_ANALISE | Somar contador; min/max timestamps globais (**008-12**, **008-14**) |
| Ocorrência **FECHADA** | Finder ainda envia achado; **reabertura** FECHADA→ABERTA é do **Módulo** (**008-4 B**) |
| Status após incremento | **Não** muda por contador (**008-10 A**) |

## Request / Response — **DN-009-2 [PENDENTE]**

| Item | Valor |
|------|-------|
| Método / path | **PENDENTE** — Drive não documenta rota Finder→Módulo (respostas-etapa3.md §10) |
| Autenticação | JWT interno (**N1 A**); algoritmo/claims **ADR_028 [PENDENTE]** |
| Granularidade consulta | **Uma chamada por achado** (**DN-008-7 B**) — (fonte: respostas-joao.md §Etapa 3c (João, 16/11 11h02); origem: N2 B, respostas-etapa3c.md §6). Envio de `ResumoAnalise`: **um por lote** (**DN-008-7**, ADR_062) — contrato [resumo-analise-agrupado.md](./resumo-analise-agrupado.md) |
| Corpo resposta mínimo lógico | `existe`, `contagemAcumulada`, `primeiraOcorrencia`, `ultimaOcorrencia`, `status` — nomes JSON **PENDENTE** |

## Erros — **DN-009-2 [PENDENTE]**

| Código | Comportamento Finder |
|--------|---------------------|
| 4xx | Falha de contrato — não concluir handoff do lote (inferência lógica ADR_014) |
| 5xx / timeout | **PENDENTE** — opções 10d respostas-etapa3.md sem decisão; ADR_118 não copiada automaticamente |

Até contrato publicado, testes MUST usar **mock** da semântica acima, não URL real.

## O que o Finder não faz

- Criar ou atualizar ocorrência persistida (**T1 rev.**)
- Aplicar reabertura de status (**008-4** no Módulo)

## Fontes

- (fonte: respostas-joao.md, N2 B / N1 A / 008-11 rev. / Etapa 3 API Módulo; origem: ADR_062, LLT-200, respostas-etapa3.md §10)
- (fonte: revisao-T1-drive.md, §2)