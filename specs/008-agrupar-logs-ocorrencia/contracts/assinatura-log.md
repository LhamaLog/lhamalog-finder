# Contract: Calcular assinatura do log

Operação **interna** do Finder após classificação (spec 006) e antes do agrupamento.

## Entrada

| Campo | Obrigatório | Descrição |
|-------|-------------|-----------|
| Log recuperado | sim | Conteúdo OTel no OpenObserve, referenciado pelo apontamento |
| `tenantId` / `productId` | sim | Do envelope v1 validado |
| `classe` | sim | ERROR \| WARNING \| INFO |

**Pré-condição**: lote **completo** (spec 005); `productId` coerente — senão falha crítica (**008-16 A**, ADR_049).

## Saída

| Campo | Descrição |
|-------|-----------|
| `assinatura` | Identificador canônico derivado dos componentes abaixo |

## Componentes (decididos — DN-008-1 B)

A assinatura MUST ser função determinística de:

| Componente | Origem OTel | Fonte |
|------------|-------------|-------|
| `nivel` | Nível/severity do log | DN-008-1 B |
| `mensagemNormalizada` | Corpo da mensagem após normalização | DN-008-1 B |
| `serviceName` | `service.name` | DN-008-1 B |
| `host` | Identificador de máquina/host (ex. `host.name`) | DN-008-1 B, ADR_013 |

## Regras (decididas)

1. **`traceId`** e **`correlationId`** MUST NOT integrar a assinatura (**008-2 rev.**; substitui 008-2 A).
2. **Classe** participa da chave de ocorrência **(assinatura, classe, produto, tenant)** (**008-15 A**); o `nivel` na assinatura não remove a dimensão `classe` na chave (**008-11 rev.**).

## Normalização da mensagem (**DN-008-4 B**, etapa 3c)

A mensagem MUST ser transformada **antes** de compor a assinatura:

1. **Aparar** espaços em branco nas extremidades.
2. **Ignorar** distinção entre maiúsculas e minúsculas na comparação (case-insensitive).
3. **Substituir** cada ocorrência reconhecível de **número**, **UUID**, **endereço IP**, **data** e **valor hexadecimal** pelos marcadores literais (**DN-008-5 A**, etapa 3e):

| Categoria | Literal |
|-----------|---------|
| Número | `<NUM>` |
| UUID | `<UUID>` |
| IP | `<IP>` |
| Data | `<DATE>` |
| Hexadecimal | `<HEX>` |

(fonte: respostas-joao.md §Etapa 3c 11h01/11h; origem: DN-008-5 A; respostas-etapa3c.md §4)

**Dono:** spec **008**. Specs **006**, **009** e **010** MUST reutilizar o mesmo algoritmo ao derivar assinatura a partir da mensagem.

(fonte: respostas-joao.md §Etapa 3b, João 16/11 09h02; origem: DN-008-4 B)

## Serialização de `valor` (**008-assinatura A**, etapa 3e; separador **DN-008-8 A**, etapa 3g)

1. Após normalizar e reunir os quatro componentes (**DN-008-1 B**), o Finder MUST produzir **`valor`** como **texto** concatenando, **nesta ordem**, **nível**, **mensagem normalizada**, **`service.name`** e **host**.
2. Entre cada par de componentes adjacentes, MUST usar o separador fixo **`|`** (barra vertical U+007C) (**DN-008-8 A**).
3. **`|`** literal em qualquer componente MUST ser escapado como **`\|`** (barra invertida + barra vertical) **antes** da concatenação (**DN-008-9 A**, etapa 3h). Escape de **barra invertida** literal nos componentes: **pendente** (sem fonte; **MUST NOT** inventar).

**Exemplo ilustrativo** (componentes sem `|` interno): `ERROR|falha ao conectar <IP>|orders-api|vm-a-01`

(fonte: respostas-joao.md, DN-008-8 A; origem: §Etapa 3e (João, 16/11 12h40); 008-assinatura A; DN-008-1 B)

## O que não fazer

- Inventar regras de normalização **divergentes** de **DN-008-4 B** ou literais **divergentes** de **DN-008-5 A**.
- Inventar escape divergente de **DN-008-9 A** (**`|`** → **`\|`**) ou regra para **`\`** literal sem fonte; o separador **`|`** está fechado (**DN-008-8 A**).
- Incluir `correlationId`/`traceId` na assinatura após **DN-008-1 B**.

## Fontes

- (fonte: respostas-joao.md, DN-008-1 B / 008-2 rev. / 008-15 A; origem: ADR_013, contrato-log-ingested-v1.txt §4)
- (fonte: respostas-joao.md §Etapa 3b, DN-008-4 B; spec.md, FR-002, FR-002b)

## Evidência desta execução

Implementação canônica em `internal/application/assinatura/`: cinco marcadores fechados, trim/case-insensitive; UUID/data/IP/hex reconhecidos antes de números internos. Cálculo usa classe, mensagem pós-segurança normalizada, service.name e host nessa ordem, com `|` e escape literal `|`→`\|`. `traceId`/`correlationId` ficam fora.

Um componente contendo barra invertida literal retorna `ErrEscapePendente`; não foi inventado um segundo protocolo de escape. A matriz formal P-010-27 A segue aberta; os testes mecânicos não a substituem.
