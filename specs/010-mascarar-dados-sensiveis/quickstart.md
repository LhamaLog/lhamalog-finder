# Quickstart: Validar mascaramento (feature 010)

Guia de **validação manual/contratual** após implementação. **Não** substitui casos de teste do João (**P-010-27**) nem regex oficiais (**P-010-2**).

**“Execuções acordadas”** (SC-001–SC-004): conjunto de casos definidos pelo João (**010-27 A**; analyze **A1**) — não inventar exemplos neste artefato.

## Pré-requisitos

- Finder com pipeline de lote completo (captura 003–005, mascaramento 010, classificação 006) — **hoje parcial** no código (`internal/application/consumir.go` sem `validarSeguranca`).
- Ambiente Docker (`make run` — README).
- Artefato de catálogo versionado publicado com `versao` alinhada à política 010-19 A e à versão **combinada** no deploy (**DN-010-3 B**). **Caminho** e **formato** conforme **ADR** (**DN-010-4 C**, **P-010-2**) — até a ADR publicar, cenários de leitura do arquivo permanecem **não executáveis** para aceite formal.
- **Bloqueio de aceite formal**: enquanto **P-010-27** ou **P-010-2** abertos, SC-001/SC-002 permanecem bloqueados (spec.md SC-005).
- Entrega HTTP ao Módulo: **DN-009-2** — quickstart de integração real permanece na spec **009**.

## Referências

- Regras internas: [validar-seguranca.md](./contracts/validar-seguranca.md)
- Handoff: [saida-mascarada-handoff.md](./contracts/saida-mascarada-handoff.md)
- Modelo: [data-model.md](./data-model.md)
- API Módulo (pendência): [plan.md](./plan.md) § API do Módulo Ocorrência

---

## Cenário 1 — Ordem do pipeline (P1)

**Objetivo**: Mascaramento ocorre antes da classificação (010-3, 010-5).

1. Ingerir lote de teste com linhas conhecidas (quando João publicar casos — **P-010-27**).
2. Instrumentar ou inspecionar logs internos: `validarSeguranca` MUST executar **antes** do passo de classificação.
3. **Esperado**: MotorAnalise/classificador nunca recebe linhas brutas com padrões **catalogados** (010-5 A).

**Fonte**: (fonte: spec.md FR-002; respostas-documentadas.md, 010-3 A / 010-5 A)

---

## Cenário 2 — Marcador e placeholder (interim ADR_064)

**Objetivo**: FR-006 com padrões no catálogo (quando ADR listar regex).

1. Processar linha contendo padrão de **senha/token/secret** coberto.
2. **Esperado**: string exata `"log com senha omitido"` (010-10 A).
3. Processar linha contendo **e-mail** coberto.
4. **Esperado**: e-mail substituído por placeholder fixo (010-11 A; valor exato **NEEDS CLARIFICATION** até ADR).

**Bloqueio**: sem casos João + ADR — marcar como **não executável para aceite** (SC-005).

**Fonte**: (fonte: respostas-joao.md, 010-10 A / 010-11 A; spec.md SC-003)

---

## Cenário 3 — Contador no ResumoAnalise

**Objetivo**: FR-005 / 010-17.

1. Lote com quantidade **conhecida** de ocorrências de padrões (após P-010-27).
2. Concluir lote e inspecionar `ResumoAnalise` / resumo do arquivo.
3. **Esperado**: contador com rótulo **“dados sensíveis mascarados”** = `dadosSensiveisIdentificados` / `totalSensivel` acordado (SC-004 quando testável; chave JSON handoff pode ser `[INFERIDO]` — glossário **I2**).

**Fonte**: (fonte: revisao-T1-drive.md, 010-17; origem: ADR_102)

---

## Cenário 4 — Handoff sem bruto

**Objetivo**: 010-24, FR-011.

1. Após processar lote, capturar payload `ResumoAnalise` (mock de porta de saída; HTTP real **bloqueado DN-009-2**).
2. **Esperado**: `achados[].linhas` só com texto mascarado; **não** há campo com conteúdo integral bruto do OpenObserve.

**Fonte**: (fonte: respostas-documentadas.md, 010-24 C; contrato-log-ingested-v1.txt §2)

---

## Cenário 5 — Falha de validarSeguranca (010-16 C)

**Objetivo**: FR-010.

1. Simular falha (catálogo ilegível, erro interno — forma **NEEDS CLARIFICATION** na implementação).
2. **Esperado**: lote **continua**; achados **sem** trechos mascarados indevidos; **não** bloqueio total do lote.

**Fonte**: (fonte: respostas-joao.md, 010-16 C)

---

## Cenário 6 — Paridade de versão do catálogo

**Objetivo**: FR-008 / SC-006 (analyze **E2**).

1. Conferir que o Finder **registra de forma auditável** o identificador `versao` do catálogo lido (startup ou por lote) — evidência mínima **neste repositório** para SC-006.
2. Paridade com Sugestão: conferir (fora deste repo) mesma `versao` no mesmo ambiente (010-19 A) — coordenação cross-módulo, não substitui evidência no Finder.
3. Configurar versão **combinada** distinta da `versao` embutida no artefato e reiniciar o Finder.
4. **Esperado**: processo **não** sobe / readiness falha com falha **observável** (log ou métrica — **DN-010-3 B**, **SC-006**).
5. Restaurar paridade e conferir subida normal.

**Fonte**: (fonte: respostas-joao.md §Etapa 3c, João 16/11 11h, DN-010-3 B; respostas-etapa3c.md §7; spec.md FR-008 / SC-006; 010-18 A / 010-19 A; Clarifications E2)

---

## Cenário 7 — Fora do catálogo (010-13)

**Objetivo**: FR-009.

1. Linha com PII **não** listada na ADR (ex. CPF — 010-12).
2. **Esperado**: texto **inalterado** no handoff; risco RIPD R03 documentado (não corrigido nesta história).

**Fonte**: (fonte: respostas-documentadas.md, 010-12 C / 010-13 A)

---

## Cenário 8 — Ordem mascaramento → normalização (FR-015, **010 H1**)

**Objetivo**: **DN-008-4 B** — assinatura/agrupamento na spec **008** não pode ver linha bruta com padrão catalogado; a **010** **reaproveita** o normalizador canônico da **008** (**010 H1**).

1. Processar lote com padrão catalogado na mensagem (após **P-010-2** / casos João).
2. Instrumentar ou inspecionar: `validarSeguranca()` MUST concluir **antes** de qualquer passo que normaliza `mensagem` para assinatura.
3. **Esperado**: normalização recebe texto já com placeholder/marcador FR-006; chamada MUST delegar a `internal/application/assinatura/normalizar_mensagem.go` (**008**, **T016**) — MUST NOT existir `normalizar_mensagem` duplicado na feature **010**. Literais **`<NUM>` `<UUID>` `<IP>` `<DATE>` `<HEX>`** (**DN-008-5 A**) — validação canônica na spec **008**, aqui só ordem pipeline e ponto de integração.

**Fonte**: (fonte: respostas-joao.md, 010 H1; origem: §Etapa 3g, João 16/11 13h22; §Etapa 3b, DN-008-4 B; §Etapa 3c, João 16/11 11h, DN-008-5 A; respostas-etapa3c.md §4; spec.md FR-015)

---

## Cenário 9 — Atributos OTel (DN-010-2)

**Objetivo**: FR-004a — escopo só linhas.

1. Registro com PII em **atributo** OTel (não na linha de mensagem).
2. **Esperado**: `validarSeguranca()` **não** altera o atributo; linha mascarada segue FR-006 se aplicável.

**Fonte**: (fonte: respostas-joao.md, DN-010-2; origem: respostas-etapa3.md §8)

---

## Não coberto neste quickstart

| Item | Motivo |
|------|--------|
| Casos LLT-202 crit. 1/2 completos | **P-010-27** |
| Regex oficiais | **P-010-2** |
| POST HTTP real ao Módulo | **DN-009-2** (spec 009) |
| Algoritmo completo de normalização / assinatura | Spec **008** (donas de **DN-008-4 B**; literais **DN-008-5 A**; **010 H1** = delegar a `normalizar_mensagem.go` **T016**) |
| Gate **≥95% precisão por classe** e base JSONL | Spec **007** (`datasets/labeled-logs/` documentado — **007 caminho C**; **DN-007-2 B**) |
| Meta de desempenho mascaramento | 010-28 A (sem meta) |
| Montagem template e-mail Notificação | NOT-03 [PENDENTE] — consumo via **DN-010-1 A** já definido |

## Comandos do repositório (quando implementado)

```bash
make test
```

Validação de aceite formal aguarda conjunto do João e ADR — não executar `go test` como substituto de SC-001/002 sem esses artefatos (010-27 A).

## Verificação desta execução

Cenários 1/4/5/8: testes de application/seguranca_test.go, ordem_seguranca_test.go e cadeia 009 verificam uma passagem, ausência de bruto após falha e ordem mascaramento → classificação → normalização canônica 008. Cenários 2/3: mecanismos de marcador/contagem exercitados com transformador opaco; não substituem catálogo/casos P-010-2/P-010-27. Cenário 6: factory rejeita versão divergente, porém bootstrap/readiness e paridade Java não verificados. Cenário 7: regras reais fora do catálogo dependem artefato. Cenário 9: o código novo atua somente em linhas, não modifica captura integral/atributos OTel; sem aceite cross-serviço. Não existe meta de desempenho própria (010-28 A), nem expansão da minimização ampla LLT-51.

Comandos realmente observados e saídas completas: `evidencias/aceite-final-unit.txt`, `aceite-final-race-integracoes.txt` e `aceite-final-vet.txt` (exit 0). Executados em composição temporária dos patches das specs, com manifesto SHA-256. Para reproduzir: compor os pré-requisitos conforme relatório central, usar Go 1.27.2 e fornecer explicitamente os parâmetros de teste de Redis/OpenObserve; credenciais não estão nos relatórios. `make test` falhou no contêiner por certificado ao baixar módulos; o comando Go local com race passou. Não executar `go test ./...` isoladamente numa branch consumidora e interpretar falta de dependências como falha da composição.
