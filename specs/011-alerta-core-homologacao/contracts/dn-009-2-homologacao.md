# Contrato referenciado: DN-009-2

Esta feature não publica outra versão do contrato. O texto normativo continua em:

- `specs/009-publicar-evento-ocorrencia/contracts/proposta-dn-009-2.md`
- `specs/009-publicar-evento-ocorrencia/contracts/proposta-dn-009-2.openapi.yaml`

(fonte: spec.md, abertura e FR-006)

## Rotas que o Finder já chama

| Papel | Método e caminho | Variáveis |
| --- | --- | --- |
| Consulta do achado | `POST /internal/v1/ocorrencias/consultas` | `FINDER_MODULO_OCORRENCIA_URL`, `FINDER_MODULO_OCORRENCIA_JWT` |
| Entrega do resumo | `POST /internal/v1/resumos-analise` | as mesmas |

Base `https://`, mesmas restrições de URL do alerta. No teste local, a URL é a do Caddy no compose do Core, o mesmo host de `FINDER_CORE_URL`. Prazo padrão no cliente: 2 s na consulta e 5 s na entrega. Bearer com o JWT do ambiente. O cliente não emite o token. (fonte: contexto-llt201.md, §3; `internal/adapter/out/modocorrencia/cliente.go`; spec.md, FR-003, FR-004 e FR-005; decisoes-joao-llt201.md, pendência 13, formulário 2)

Audiência exigida nestas rotas: `lhamalog-modulo-ocorrencia`. Não é trocada por `lhamalog-core`. O JWT pode ser o mesmo de `FINDER_CORE_JWT`, com as duas audiências no `aud`. A rota de alerta, à parte, exige `lhamalog-core` via `core-audience`. (fonte: spec.md, FR-003; decisoes-joao-llt201.md, pendências 6 e 20, formulário 1; `application.properties` do #20, via `git show`)

Fato de código: a cadeia `/internal/v1/**` do #20 confere `aud` por "contém" e, hoje, exige a audiência do Módulo em toda a cadeia. Um token só com `lhamalog-core` recebe 401 sem o ajuste por rota. `scripts/InternalJwt.java` gera uma audiência por `--aud`. Este contrato não altera o script. (fonte: respostas-pendentes-llt201.md, resumo, fato de código)

## Lado servidor assumido mergeado

Recibo idempotente por (tenant, produto, batchId, execId): 201 na primeira entrega e 200 com `repeticao: true` no reenvio do resumo. `409 REVISAO_DESATUALIZADA` e `409 ENTREGA_EM_ANDAMENTO`. JWT interno HS256 com segredo, `iss`, `aud` e `sub`. FECHADA reabre como ABERTA. Pacote `br.com.lhamalog.core.occurrence`. O campo `repeticao` deste recibo não entra no corpo do alerta. (fonte: contexto-llt201.md, §1; decisoes-joao-llt201.md, base da branch; respostas-pendentes-llt201.md, pendência 8)

## Cenários que a homologação percorre

Os cinco já escritos na seção "Cenários de homologação com o Módulo" de `proposta-dn-009-2.md`, e repetidos na User Story 3 da spec 011:

1. A consulta distingue ocorrência ausente de falha e rejeita acesso a outro tenant ou produto.
2. A primeira entrega grava exatamente um efeito e retorna 201. Reenvio depois de perda de resposta, inclusive depois de reinício do Finder, retorna 200 e o mesmo `reciboId`, sem duplicar.
3. Dois lotes que consultam a mesma revisão preservam todos os incrementos. Revisão desatualizada em um achado rejeita a entrega inteira antes de qualquer escrita.
4. Confirmação inválida ou inesperada mantém o lote pendente no Finder. Um recibo válido antecede a conclusão do lote e o reconhecimento da fila.
5. O JSON mantém as três classes, carrega somente trechos protegidos e conserva `linhas` vazias quando for o caso. A API rejeita tipos e contadores inválidos.

(fonte: spec.md, User Story 3 cenário 2; `proposta-dn-009-2.md`)

## O que a conferência no Core não pede a este contrato

Ocorrência nova, reaberta e "virou recorrente", gravidade e o canal `lhamalog:occurrence-events` não são campos novos do `ResumoAnalise`. O Finder só entrega o resumo. A conferência é no Core, spec 016, e segue o FR-007 da spec 011: precedência, gatilho (`before.recurring` não é TRUE e `after.recurring` é TRUE; ausente ou `null` conta como `false`) e literal `RECURRENT`. A gravidade é a classe no campo `gravidade`. O aviso "virou recorrente" sai em até 3 s depois do 201. Esse teto, com o de 2 s da rota de alerta (`FINDER_CORE_TIMEOUT`), é o SC-009. Os dois valem em alvos diferentes. Não usar `RECURRING`. (fonte: spec.md, FR-007 e SC-009; decisoes-joao-llt201.md, pendências 3 e 5, formulário 1, e Rodada 3; respostas-pendentes-llt201.md, pendências 2, 11 e 16; origem: spec 015, SC-001 + DF-FINDER §3 e §8)

Não há envelope novo neste diretório. O alerta aceito, depois do 201, grava em `finder_alerts` e avisa a tela no canal `/ws/occurrences`, tipo `FINDER_ALERT`, corpo `{tenantId, productId, batchId, motivo}`, sem publicar em `lhamalog:occurrence-events` nem em `lhamalog:notifications:events`. Esse contrato está em [alerta-finder-cliente.md](alerta-finder-cliente.md), não neste. (fonte: decisoes-joao-llt201.md, evento para a Notificação, pendência 1, formulário 1, e Rodada 3)
