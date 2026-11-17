# Contract: Classificar log (regra pura)

Contrato **interno** entre o parser do log capturado e o agregador de achados. Não é endpoint HTTP. Não acessa fila, OpenObserve nem Módulo Ocorrência.

## Operação

`Classificar log`

### Entrada

| Campo | Obrigatório | Descrição |
|-------|-------------|-----------|
| `nivelCampo` | opcional | Texto/código de severidade OTel já extraído |
| `mensagem` | sim | Corpo da mensagem (esperado mascarado — spec 010) |
| `catalogo` | opcional | Lista ordenada de entradas `{padrao, classe}` — **sintaxe pendente** 006-7 D |

### Saída

| Campo | Descrição |
|-------|-----------|
| `classe` | `ERROR`, `WARNING` ou `INFO` |
| `catalogoUsado` | `true` somente se catálogo foi a única fonte decisória |

### Tabela de comportamento (espelho spec.md User Story 1)

| Cenário | Então `classe` | Catálogo |
|---------|----------------|----------|
| Campo nível = erro/error/fatal; mensagem sem termo mais grave | ERROR | não |
| Campo nível = warn/warning/alerta; sem termo mais grave na mensagem | WARNING | não |
| Campo nível = info/debug/trace; sem termo mais grave | INFO | não |
| Campo info + palavra inteira `error` na mensagem | ERROR | não |
| Campo fatal + palavra inteira `info` | ERROR | não |
| Campo info + `errorless` | INFO | não |
| Campo ausente + `traceback` sem outro termo; sem catálogo | INFO | não |
| Sem nível reconhecível; primeira entrada catálogo WARNING | WARNING | sim, primeira válida |
| Sem nível e sem catálogo | INFO | não |

### Obrigações transversais

- **FR-002**: o texto **não** é enviado a serviço externo de linguagem.
- **FR-003**: spam e “não classificado” **não** são valores de saída.
- Palavras de nível na mensagem contam só como **palavra inteira** (spec sessão 2026-11-09).

### Fontes

- (fonte: spec.md, User Story 1 e 3, FR-002–FR-006)
- (fonte: respostas-joao.md, 006-6 A / T7 A)