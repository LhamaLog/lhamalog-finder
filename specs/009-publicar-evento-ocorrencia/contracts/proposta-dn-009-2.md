# Proposta de contrato Finder → Módulo de Ocorrência

**Status** · Aprovado em 20/11/2026 · Spec 009 · DN-009-2 / N3 A · LLT-201

O Finder consulta o estado de cada achado, calcula os totais e entrega um `ResumoAnalise` por lote. O Módulo de Ocorrência persiste as ocorrências, determina a gravidade e envia os avisos ao Core. O Finder só conclui o processamento após confirmação explícita do Módulo. As rotas, os campos JSON e as regras de transporte abaixo definem a interface entre os serviços. [ADR 062](https://docs.google.com/document/d/16XOz3bMdXwRTDHIkY0zuGx6faDSi8COtD84e-zecpR8/edit) · [Varredura e Finder v6](https://docs.google.com/document/d/1XdmrVlWxS5P82wAzETQwZTyTqtXP1HDe/edit) · [LLT-201](https://logviewer2026.atlassian.net/browse/LLT-201)

[OpenAPI 3.1 com esquemas e exemplos](proposta-dn-009-2.openapi.yaml)

## Rotas

| Operação | Rota | Resposta |
|---|---|---|
| Consultar um achado | `POST /internal/v1/ocorrencias/consultas` | `200` com `estado: null` quando não existe ocorrência, ou `estado` com contagem acumulada, extremos temporais, status e revisão. |
| Entregar um resumo | `POST /internal/v1/resumos-analise` | `201` após o primeiro registro confirmado e `200` em reenvios de chave já aceita. Ambos retornam um recibo explícito. |

A consulta é somente leitura e recebe a tupla `(tenantId, productId, assinatura, tipo)` em JSON. O `POST` mantém a assinatura fora da URL, que costuma aparecer em logs de acesso. Cliente e servidor devem excluir assinatura, trechos de log e credenciais dos registros operacionais. Cada consulta retorna zero ou uma ocorrência para a tupla. Se o armazenamento tiver mais de uma, o Módulo devolve erro e o Finder mantém o lote sem confirmação. [Contrato lógico da 008](../../008-agrupar-logs-ocorrencia/contracts/consulta-ocorrencias-modulo.md)

## Identidade, autenticação e confirmação

- **Autenticação.** HTTPS interno em rede privada e `Authorization: Bearer <JWT interno>`. Propõe-se validar assinatura, `iss`, `aud`, `sub` e `exp`, além de `nbf` quando presente. O emissor confiável, a audiência, a identidade do Finder, o algoritmo permitido, as chaves e a rotação são configuração do ambiente. O Módulo valida a autorização do serviço para a operação e para o par tenant/produto, além de conferir que o produto pertence ao tenant. [ADR 028](https://docs.google.com/document/d/1xJSvxvHGgRrC5wxJTDISyxN5onMbNnW8_6GCe7KRVCw/edit)
- **Chave no corpo.** O Módulo usa `batchId` para `log.ingested` v1 e a tupla `batchId + execId` quando o Finder recebe um ID real de Varredura (`V<uuid>` ou `M<uuid>`). A reserva fica no escopo de `(tenantId, productId)` e usa valores estruturados. `execId` ausente é omitido. [Contrato do Coletor](https://docs.google.com/document/d/14yLd8WA7pRm6es4okYiOnUjHfZp0T_4uyUnVkFk2tQk/edit) · [ADR 070](https://docs.google.com/document/d/1dJvPaQHU32ThNt_epI13JcdxdLcPe51DiWsssH3LtYE/edit)
- **Recibo.** O Módulo responde `{ "aceito": true, "reciboId": "…", "tenantId": "…", "productId": "…", "batchId": "…", "repeticao": false }` e inclui `execId` quando aplicável. O `reciboId` é opaco e estável. O Finder valida status, corpo e a identidade devolvida antes da conclusão e da confirmação Redis. Recebimento assíncrono (`202`), resposta sem corpo (`204`), recibo inválido ou divergência de identidade mantêm o lote sem confirmação.
- **Reenvio.** A primeira entrega confirmada para a chave é definitiva. O Módulo retorna o recibo original com `repeticao: true` antes de revalidar as revisões consultadas. Um resumo recalculado após reinício pode ter contagens diferentes e recebe o mesmo recibo sem reaplicar alterações. Um processamento intencionalmente novo exige um `execId` distinto vindo da Varredura. Durante o mesmo envio, o Finder reutiliza o conteúdo que já montou. A retenção do recibo deve cobrir o período em que a mesma chave pode ser reentregue. O Módulo precisa definir essa retenção com a equipe que opera o Redis. [ADR 014](https://docs.google.com/document/d/1rcWvAzhaatapzJ3hQch6NiyYaCTk6qVGzWmTZnlrdZ0/edit)

## Campos do resumo e escrita concorrente

O contrato mantém `contadores.informacao` do serializador atual como nome de transporte para a classe `INFO`. `linhasProcessadas = erros + warnings + informacao`, com contagens de **logs** do lote. O contador `dadosSensiveisMascarados` contém o total do scanner e é independente da soma das classes. A versão das regras fica no resumo e vale para todos os achados daquele lote. A versão do catálogo de mascaramento, quando conhecida, segue em `versaoCatalogo`. [ADR 102](https://docs.google.com/document/d/19IY82QrXHf-OewsdQ1HqPvw0jOnmDu6iXSvR-Yr6IhY/edit) · [ADR 170](https://docs.google.com/document/d/12SGFJs3Ue-XG9eRCBrRqVnq9IIYs2uacrCL0tH0dtTU/edit)

Cada achado tem `tipo`, `assinatura`, `hostName`, `contagem`, `primeiraOcorrencia`, `ultimaOcorrencia`, `recorrente`, `linhas` e `revisaoConsultada`. `contagem` já inclui a contagem acumulada consultada mais o incremento do lote. Os horários representam o primeiro e o último log, em UTC. A assinatura usa o formato da 008. `linhas` contém somente conteúdo que passou pelo scanner e pode ser vazio. O contrato inclui `titulo`, `resumo`, `numeroLinha`, `camposExtraido` e `dadosSensiveis` com `null` para informação indisponível, `numeroLinha` baseado em 1 quando a linha de origem é conhecida e `dadosSensiveis` booleano quando o resultado é conhecido para o achado. Qualquer texto derivado do log nesses campos de apresentação também vem da saída do scanner. `traceId`, `correlationId` e `moduloId` são opcionais. O Módulo é responsável pelos IDs de ocorrência, gravidade e criticidade. [ADR 040](https://docs.google.com/document/d/1j35UoBKtjqS82ihh7wHuHDj9iTzBe3zBS3Qz6tNJU34/edit) · [ADR 064](https://docs.google.com/document/d/10SmZvV1pHH7HwLkpfVpEzZ4DaNTcieGIlzrqwPpphr0/edit)

**Proposta para concorrência.** O Módulo gera uma revisão opaca que muda quando o estado retornado pela consulta muda. O Finder inclui essa revisão no achado, ou `null` quando a consulta informou ausência. Para uma chave nova, o Módulo compara todas as revisões e aplica os valores finais uma única vez, com a gravação durável do recibo. Por exemplo, se a consulta devolveu `10` e o lote traz `3`, o Finder envia `13` e o Módulo grava `13`. Se outro lote mudar a contagem para `12` antes da entrega, a primeira tentativa recebe `409 REVISAO_DESATUALIZADA` sem efeito parcial. O Finder pode consultar novamente e enviar `15` sob a mesma chave porque ainda não existe recibo. Uma chave já confirmada segue a regra de reenvio. O Módulo aplica a reabertura conforme a regra da 008. [Contrato lógico da 008](../../008-agrupar-logs-ocorrencia/contracts/consulta-ocorrencias-modulo.md)

## Falhas e tempo de espera

| Resultado | Ação do Finder |
|---|---|
| `400`, `401`, `403`, `422` | Registrar código sanitizado, manter o lote sem confirmação e encaminhar para correção. Corrigir credencial, autorização ou formato antes de tentar novamente. |
| `409 REVISAO_DESATUALIZADA` | Consultar os estados novamente e recalcular o resumo sob a mesma chave. Uma nova tentativa nesse ciclo antes de deixar o lote pendente. |
| `409 ENTREGA_EM_ANDAMENTO`, `429`, `500`, `502`, `503`, `504` ou timeout | Reenviar com a mesma chave quando o agendamento de retry permitir. Respeitar `Retry-After` se fornecido. Manter o lote pendente até receber um recibo válido. |
| Recibo ausente ou inválido | Registrar falha de protocolo. Uma nova tentativa deve preservar a chave. |

O cliente usa **2 s por consulta** e **5 s por entrega** por padrão, configuráveis na construção do adaptador, com cancelamento por contexto. A regra de três tentativas do OpenObserve e o alerta operacional ao Core pertencem à 008/ADR 051. [ADR 051](https://docs.google.com/document/d/1tmtv9FPzSLWJmBfP1vzVXp1re6ogibfqcIQshrV8u14/edit)

## Cenários de homologação com o Módulo

1. A consulta distingue ocorrência ausente de falha e rejeita acesso a outro tenant/produto.
2. A primeira entrega grava exatamente um efeito e retorna `201`. Reenvio após perda de resposta, inclusive após reinício do Finder, retorna `200` e o mesmo `reciboId` sem duplicação.
3. Dois lotes que consultam a mesma revisão preservam todos os incrementos. Revisão desatualizada em um achado rejeita a entrega inteira antes de qualquer escrita.
4. Confirmação inválida ou inesperada mantém o lote pendente. Um recibo válido antecede a conclusão e o `XACK`.
5. JSON mantém as três classes, carrega somente trechos protegidos e conserva `linhas: []` quando aplicável. A API rejeita tipos e contadores inválidos.

O Finder implementa os cenários de cliente. A homologação requer o Módulo implementado e configurado com emissor, audiência e rotação JWT, além da retenção operacional de recibos. O agendamento em produção, incluindo o encaminhamento de rejeições permanentes para correção, faz parte da ativação do fluxo. [Evidência do cliente e do Redis](../quickstart.md#verificação-do-contrato-dn-009-2)
