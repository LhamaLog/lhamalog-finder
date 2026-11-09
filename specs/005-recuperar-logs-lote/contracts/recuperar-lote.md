# Contract: Recuperar lote

Contrato entre o caso de uso da recuperação e as portas. Não é um endpoint novo do Finder. A confirmação da entrada Redis continua na primeira tarefa.

## Porta de entrada

`Recuperar logs do lote`

| Entrada | Saída |
|---------|--------|
| Referência válida, configuração completa, nomes de tenant e produto, lote cuja contagem distinta fecha e sem falha | Efeito concluído. Registros lógicos, um por `lhamalog_record_id`, e a correlação |
| Mesma referência com `recordId` | A consulta também filtra `lhamalog_record_id` |
| Referência sem `recordId` | A consulta não filtra identificador de registro |
| URL, credencial, timeout, tamanho de página, lista de cópia ou expressão de ordenação ausente | Não consulta. Sem efeito concluído. Nenhum padrão criado |
| Lista de cópia presente e vazia, dois documentos do mesmo identificador com qualquer diferença | Conteúdo divergente. Efeito concluído. Nenhuma versão escolhida |
| Nomes de tenant ou produto ausentes | Não consulta. Isolamento. Efeito concluído. Nenhum documento válido |
| Documento de outro tenant ou produto | Isolamento. Efeito concluído. Nenhum documento válido |
| Documento sem `lhamalog_record_id` e também outra falha de menor precedência | Documento sem identificador. Efeito concluído |
| Isolamento junto com documento sem identificador, divergência ou contagem diferente | Isolamento. Efeito concluído |
| Contagem distinta menor, sem falha de maior precedência | Lote incompleto. Sem efeito concluído. Quantidades esperada e encontrada |
| Contagem distinta maior, sem falha de maior precedência | Excesso. Efeito concluído. Quantidades esperada e encontrada |
| Busca vazia com `recordCount` positivo | Lote incompleto. Sem efeito concluído. Encontrado igual a zero |
| Credencial recusada | Falha de autenticação ou autorização. Efeito concluído. Sem credencial no resultado |
| Armazenamento indisponível ou timeout, sem falha que pede correção | Sem efeito concluído. Parcial não é lote pronto |
| Página posterior falha depois de falha que pede correção | A classificação é a falha que pede correção. Efeito concluído. Parcial não é lote pronto |
| Página cheia | A seguinte é consultada. Não há máximo de páginas |
| Valor que não pode ser dado do filtro | Resposta inválida. Efeito concluído. A consulta não troca lote, organização, stream nem faixa |

A entrada Redis não é confirmada em nenhuma linha.

## Porta de saída: busca

| Operação | Obrigação |
|----------|-----------|
| Consultar uma página | `POST /api/{organization}/_search` com `query.sql`, `query.start_time`, `query.end_time`, `query.from` e `query.size` |
| Faixa | `start_time = minTimestampMicros` e `end_time = maxTimestampMicros + 1` |
| SQL | Stream da referência, `lhamalog_batch_id`, `lhamalog_record_id` quando houver, tenant e produto quando os nomes existem, e a expressão de ordenação configurada |
| Página | `from` começa em 0. A página seguinte soma `size`. `size` não é `recordCount` |
| Resposta | Devolve `hits` ou informa recusa de credencial, resposta inválida, timeout ou indisponibilidade |
| Fora desta porta | Confirmação Redis, classificação de ocorrência, ocultação de campo e escolha de versão divergente |

O adaptador não registra a senha nem o corpo dos hits em falha.

## Porta de saída: efeito para a primeira tarefa

| Operação | Obrigação |
|----------|-----------|
| Informar | Entrega a classificação, a marca de efeito concluído, o contexto do lote e as quantidades quando a contagem foi calculada |
| Sucesso | Entrega também os registros lógicos |
| Fora desta porta | `XACK`, chave Redis e leitura do stream |

Sem efeito concluído, a primeira tarefa não confirma a entrada válida. Esta obrigação já está em `specs/004-capturar-referencia-lote/contracts/referencia-lote.md`.

## Registro operacional

Classificação, `idempotencyId` e identificador da entrada. Nas falhas de contagem, as duas quantidades. Sem senha, sem conteúdo de log e sem documento de outro tenant.

## O que o Finder não expõe nesta feature

Nenhum endpoint HTTP novo. `GET /health` permanece como está.
