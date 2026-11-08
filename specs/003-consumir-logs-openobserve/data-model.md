# Data Model: Consumir logs no OpenObserve

Organização, stream, índice, campos, credenciais e o conteúdo de `otel-logs.zip` não fazem parte deste modelo. Eles continuam entrada externa. O lugar do log pronto não tem nome.

## Apontamento mockado

| Campo | Regra |
|-------|--------|
| Valor | Seleciona um único log no OpenObserve. Não é o conteúdo do log. Formato, tamanho e alfabeto não estão definidos. |
| Vazio | Inutilizável. A captura não começa e a falha leva esse motivo. |
| Repetido | Várias indicações do mesmo valor na mesma execução contam como um apontamento. |
| Fora | Um valor que selecionaria vários logs, ou o exemplo inteiro, não entra neste modelo. |

O tipo já existente `domain.Codigo` é a posição do log e não é o log. Este plano reutiliza esse papel. Não cria um segundo formato.

## Registro armazenado

| Campo | Regra |
|-------|--------|
| Conteúdo | O que o OpenObserve devolve para o apontamento. Os campos internos não são descritos. |
| Aceite | No SC-001, o conteúdo é um registro de `otel-logs.zip`. O arquivo precisa ser entregue. Este modelo não lista esses registros. |
| Ausente | O apontamento é utilizável e o OpenObserve não tem o log. A falha leva o motivo log ausente. Nenhum registro vazio é criado no lugar. |

## Log pronto

| Campo | Regra |
|-------|--------|
| Apontamento | O valor que selecionou o registro. |
| Conteúdo | Igual ao registro armazenado. Nenhum campo é ocultado ou reescrito. |
| Leitura | Com a instância no ar, o mesmo apontamento devolve esse único conteúdo. |
| Parada | Depois que a instância para, o modelo não exige que o log continue. Parada no meio da captura não deixa log pronto. |
| Fora | Ocorrência, categoria erro, alerta ou informação, e lugar nomeado. |

## Falha de captura

| Campo | Regra |
|-------|--------|
| Apontamento | O valor usado na tentativa. |
| Motivo | Um destes três: inutilizável, ausente, inacessível. Os três são distinguíveis. |
| Efeito | O log desse apontamento não fica pronto. Uma falha não apaga o resultado de outro apontamento. |
| Fora | Nova tentativa, prazo numérico e fila de erro. |

## Estados da captura

| Transição | Quando | Efeito |
|-----------|--------|--------|
| Sem apontamento | A fonte mockada não fornece valor | Nenhum log pronto e nenhuma falha |
| Indicado → inutilizável | O valor é vazio | Falha com esse motivo. O OpenObserve não é consultado para esse valor |
| Indicado → em captura | O valor é utilizável e a instância está no ar | Uma captura por valor, mesmo que a indicação se repita |
| Em captura → pronto | O registro é obtido | O conteúdo fica igual ao armazenado e legível pelo apontamento |
| Em captura → ausente | O armazenamento responde e não tem o log | Falha com esse motivo |
| Em captura → inacessível | O armazenamento não pode ser alcançado | Falha com esse motivo, distinta da ausência |
| Em captura → interrompida | A instância para no meio | O log não fica pronto. Sem prazo e sem retomada |

## Relações

- Uma execução tem zero ou mais apontamentos.
- Um apontamento tem um resultado: log pronto ou falha. Apontamento repetido não cria o segundo resultado.
- Duas capturas de apontamentos distintos podem ocorrer ao mesmo tempo. Cada uma conserva o próprio resultado.
- A fila Redis não participa desta relação.
- Várias instâncias do Finder não compartilham o log pronto neste modelo.

## Validação que o domínio aplica

1. Apontamento vazio é inutilizável.
2. O conteúdo pronto é o conteúdo armazenado, sem campo oculto.
3. O mesmo apontamento tem no máximo um log pronto na execução.
4. Falha não produz log pronto.
5. Motivo ausente e motivo inacessível não se substituem.
