# Contract: `validarSeguranca()` / ScannerSegurança (interno)

Interface **lógica** no Finder. Nomes conforme ADR_016 e ADR_064 (FR-007). Sem fixar pacote Go neste repositório.

## Posição no pipeline

| Ordem | Etapa | Fonte |
|-------|-------|-------|
| 1 | Recuperar linhas do lote (OpenObserve — specs 003–005) | spec.md FR-012 |
| 2 | **`montarContexto()`** chama **`validarSeguranca()`** uma vez | 010-3 A, 010-4 A, 010-15 A; ADR_016 |
| 3 | **MotorAnalise** (classificação, agrupamento — 006/008) | 010-5 A |

**Proibido**: classificar ou montar `AchadoAnalise.linhas` a partir de texto **não** passado por esta etapa (010-5 A, 010-21).

**Ordem vs. normalização de assinatura** (**DN-008-4 B**): normalização canônica na spec **008** (trim; case-insensitive; substituir números, UUIDs, IPs, datas e hex por marcadores distintos por categoria). Implementação: `internal/application/assinatura/normalizar_mensagem.go` (**T016**). A spec **010** MUST **invocar** esse artefato — **MUST NOT** definir algoritmo divergente (**010 H1**, etapa 3h). Esta etapa MUST receber mensagem/linha **já mascarada** quando derivada do lote processado (**FR-015**, **010-5 A**). Literais de normalização (**DN-008-5 A**): **`<NUM>`**, **`<UUID>`**, **`<IP>`**, **`<DATE>`**, **`<HEX>`** — definidos na spec **008**; **não** confundir com **`"log com senha omitido"`** (mascaramento — 010-10 A). (fonte: respostas-joao.md, 010 H1; origem: §Etapa 3g, João 16/11 13h22, DN-008-4 B)

## Entrada

| Parâmetro | Tipo | Obrig. | Descrição | Fonte |
|-----------|------|--------|-----------|-------|
| `linhasBrutas` | `string[]` | sim | Uma string por linha de log do lote | 010-14 C; ADR_011 |
| `versaoCatalogoEsperada` | string | sim | Versão do artefato carregado (paridade 010-19) | 010-19 A |
| `catalogo` | objeto ou referência | sim | Regras vindas do [catalogo-padroes-versionado.md](./catalogo-padroes-versionado.md) | FR-008 |

**Escopo (DN-010-2)**: entrada cobre **somente** corpo de **linhas de log** como em ADR_064. Atributos OTel sensíveis **não** são entrada deste contrato — responsabilidade do Cliente na origem (fonte: respostas-joao.md, DN-010-2; origem: respostas-etapa3.md §8, DPA Anexo IV item 6).

## Saída (`ResultadoValidarSeguranca`)

Ver [data-model.md](../data-model.md).

| Campo | Obrig. | Regra | Fonte |
|-------|--------|-------|-------|
| `linhasMascaradas` | sim | Mesmo cardinal que entrada, salvo linhas substituídas integralmente pelo marcador | ADR_064; FR-006 |
| `totalSensivel` | sim | Incrementa por cada aplicação de padrão catalogado (e-mail substituído, linha removida, etc.) | 010-17; ADR_064 |
| `erro` | condicional | Presente em falha de execução | 010-16 C |

## Regras de transformação (interim — até ADR P-010-2)

| Padrão (catálogo) | Efeito | Fonte |
|-------------------|--------|-------|
| E-mail | Substituir por placeholder fixo | 010-11 A; FR-006 |
| Linha com senha, token ou secret | `texto` = `"log com senha omitido"` (linha inteira) | 010-10 A; FR-006 |
| Fora do catálogo | Sem alteração | 010-13 A; FR-009 |

Regex concretas: **[PENDENTE ADR]** — não listadas neste contrato (010-2 A).

## Falha (`010-16 C`)

| Situação | Comportamento Finder | Fonte |
|----------|----------------------|-------|
| `validarSeguranca()` retorna erro | Continuar processamento do lote **sem** incluir trechos mascarados nos achados (`linhas` vazias ou omitidas conforme política de achado) | 010-16 C; FR-010 |
| Lote inteiro | **Não** fail-closed | 010-16 C (opção A/B rejeitadas) |

## Não-objetivos

- Revalidar mascarado no Módulo Ocorrência, Core ou Notificação (010-21).
- Incluir bruto em filas, handoff HTTP ou APIs Finder (010-24; FR-011).
- Mascarar atributos OTel do registro (DN-010-2).
- Aplicar escopo LLT-51 (010-26 A).
- Redefinir algoritmo de normalização (dona spec **008**, **DN-008-4 B**) — literais **DN-008-5 A** só por referência cruzada nesta feature.

## Dependências

- Carregamento do catálogo: [catalogo-padroes-versionado.md](./catalogo-padroes-versionado.md)
- Consumo da saída: classificação spec **006**, handoff [saida-mascarada-handoff.md](./saida-mascarada-handoff.md)

## Fontes consolidadas

- (fonte: spec.md FR-002, FR-004a, FR-006, FR-007, FR-009, FR-010)
- (fonte: respostas-documentadas.md, 010-3–010-5 / 010-14 / 010-15)
- (fonte: respostas-joao.md, 010-10 A / 010-16 C / DN-010-2)
- (origem: ADR_016, ADR_064, SSID_007.1)

## Evidência desta execução

Validação atua somente nas linhas, não em atributos OTel (DN-010-2). O mecanismo recebe transformador opaco do catálogo e preserva a política de exatamente uma passagem por lote. Falha, retorno parcial/cardinalidade inválida ou catálogo ausente elimina trechos do contexto; não existe fallback para linha bruta. Contador é o total retornado pelo artefato.

Regex oficiais, motor, caminho/formato e placeholder de e-mail não foram escolhidos. Logo, T014 tem mecanismo testado, mas as detecções/substituições reais T017/T019/T020/T021 dependem P-010-2. O literal MarcadorOmisso é o único texto de negócio fechado. Fixtures de transformador testam fluxo e ordem; não são novos casos formais de SC-001/002 do João.
