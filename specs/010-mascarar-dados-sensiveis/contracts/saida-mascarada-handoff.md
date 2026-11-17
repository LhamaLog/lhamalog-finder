# Contract: Saída mascarada no handoff (Finder → Módulo Ocorrência)

Complementa o handoff `ResumoAnalise` (specs **006** / **009**). Esta feature define **obrigações de conteúdo mascarado**. Transporte HTTP+JWT: **N1 A** — implementação **bloqueada** por **DN-009-2** (detalhes em spec **009** `resumo-analise-handoff.md` e [plan.md](../plan.md) § API do Módulo Ocorrência).

## Escopo no payload

| Área | Regra | Fonte |
|------|-------|-------|
| `achados[].linhas` | MUST ser texto **já mascarado**; cada elemento = exemplo de log | 010-6; 010-21; FR-003 |
| `achados[].linhas` em falha 010-16 | Pode ser `[]` ou ausência de trechos — MUST NOT conter bruto catalogado | 010-16 C; FR-010 |
| `contadores` (chave `[INFERIDO]` ex. `dadosSensiveisMascarados`) | MUST refletir `totalSensivel` / **`dadosSensiveisIdentificados`** do lote; rótulo humano “dados sensíveis mascarados” | 010-17; FR-005; glossário **I2** spec.md |
| Corpo inteiro | MUST NOT reexpõer log bruto do OpenObserve | 010-24 C; FR-011 |
| `gravidade`, ids de ocorrência | Fora do Finder (ADR_062) — inalterado | T1 revisto |

**NEEDS CLARIFICATION**: nomes canônicos exatos dos contadores na ADR_102 — semântica fixada; rótulos [INFERIDO] no Drive (revisao-T1-drive.md §2).

## Campos de classificação alimentados por texto mascarado

| Campo interno | Obrigação | Fonte |
|---------------|-----------|-------|
| `mensagem` (classificação) | Origem = linhas pós-`validarSeguranca()` | 010-5 A; spec 006 `classificar-log` |
| `assinatura` / normalização | Algoritmo **DN-008-4 B** (spec **008**, `normalizar_mensagem.go` **T016**): trim, case-insensitive, marcadores número/UUID/IP/data/hex; MUST operar sobre mensagem **mascarada** se derivada de linhas (**FR-015**, **010-5 A**). **010 H1**: a **010** **reaproveita** o normalizador da **008** — MUST NOT duplicar o módulo. Literais **DN-008-5 A**: `<NUM>` `<UUID>` `<IP>` `<DATE>` `<HEX>` | **010 H1**; DN-008-4 B; DN-008-5 A; 010-5 A; (fonte: respostas-joao.md, 010 H1; origem: §Etapa 3g, João 16/11 13h22) |
| Identificadores OTel na Ocorrência | **Não** passam por `validarSeguranca()` — ADR_040 | DN-010-2; FR-004a |

## Transporte HTTP (pendência)

| Aspecto | Status | Fonte |
|---------|--------|-------|
| Payload mascarado | **Definido** neste contrato | FR-003–FR-005 |
| Rota, JWT, 5xx, timeout | **DN-009-2** — pendência João | respostas-joao.md §Etapa 3; respostas-etapa3.md §10 |
| Adaptador `internal/adapter/out/modocorrencia/` | **Bloqueado** até contrato | spec 009 |

## Downstream (referência — fora da implementação Finder)

| Consumidor | Expectativa | Fonte |
|------------|-------------|-------|
| Módulo Ocorrência | Persiste `linhas` em `conteudo` | 010-25; ADR_040 |
| Telas E08.F1.3 | Exibem só `conteudo` mascarado | 010-20; LLT-205 |
| Sugestão | `conteudo` mascarado + reaplicação catálogo | 010-9; ADR_039 |
| Core (aviso HTTP) | **id + tipo + gravidade**, **sem** texto de log | 009-5 revisto; 010-23 |
| Notificação (e-mail LLT-218) | Lê Ocorrência **pelo id**; usa `conteudo` mascarado | DN-010-1 A; DN-009-1 A; respostas-etapa3.md §7 |

Finder garante mascarado **no que envia** (`ResumoAnalise`); **não** monta e-mail na Notificação (010-21). Templates NOT-03 [PENDENTE] — spec Notificação.

## Exemplo estrutural (ilustrativo — sem dados reais)

```json
{
  "contadores": {
    "dadosSensiveisMascarados": 2
  },
  "achados": [
    {
      "tipo": "ERROR",
      "linhas": [
        "2025-01-01 ERROR failed for [email]",
        "log com senha omitido"
      ]
    }
  ]
}
```

Marcador e placeholder MUST seguir FR-006 quando padrão estiver no catálogo. Valores acima são **forma**, não casos de teste (010-27 A).

## Cenários proibidos (Finder)

- Enviar senha/token/e-mail catalogado em formato original em `linhas` (SC-001 quando testes existirem).
- Incluir trechos de log no aviso Módulo → Core (responsabilidade 009; 010-23 / 009-5 revisto).

## Dependências

- [validar-seguranca.md](./validar-seguranca.md) — origem das linhas.
- Spec **009** `resumo-analise-handoff.md` — envelope completo, idempotência e bloqueio **DN-009-2**.

## Fontes consolidadas

- (fonte: spec.md FR-003–FR-005, FR-011, User Story 1–2)
- (fonte: revisao-T1-drive.md, 010-17 / 010-21 / 010-23 / 010-25)
- (fonte: respostas-documentadas.md, 010-9 / 010-20 / 010-22 / 010-24)
- (fonte: respostas-joao.md §Etapa 3, 009-5 revisto / DN-009-1 A)
- (origem: ADR_040, ADR_062, ADR_102, ADR_167, LLT-202)

## Evidência desta execução

Somente linhas do contexto seguro alimentam classificação, assinatura e handoff lógico 006→008→009. TotalSensivel alimenta contador lógico de dados sensíveis mascarados, sem fechar semântica JSON N3 A/ADR_102 nem resumo externo de arquivo.

Notificação lê Ocorrência por id, conteudo mascarado (DN-010-1 A); tela/e-mail/Sugestão ficam fora deste repo. O teste fluxo_test da 009 comprova que um sentinel bruto de fixture não chega ao DTO; não comprova catálogo real, tela, e-mail ou LLM.
