# Data Model: Recuperar logs do lote

Modelo do caso de uso. Os nomes de campo de tenant, produto e cópia física não estão aqui. Vêm da configuração da execução.

## Referência de entrada

Recorte já validado pela primeira tarefa. Esta feature não o reconstrói e não lê o log na entrada Redis.

| Campo | Regra |
|-------|--------|
| `tenantId` | Autoridade do tenant. Não é substituído por atributo do documento. |
| `productId` | Autoridade do produto. Mesma regra. |
| `organization` | Caminho do POST. Não é o nome do stream. |
| `stream` | Stream do OpenObserve na referência. Diferente de `lhamalog:stream:logs`. |
| `batchId` | Valor do filtro `lhamalog_batch_id`. |
| `recordId` | Presente no lote de um registro. Ausente no lote de vários. Quando presente, entra no filtro `lhamalog_record_id`. |
| `minTimestampMicros` | `query.start_time`. |
| `maxTimestampMicros` | `query.end_time` é este valor mais 1. |
| `recordCount` | Inteiro positivo. Conta os registros aceitos, não o lote bruto recusado. |
| `idempotencyId` | Correlação do resultado. Não deduplica documento. |
| `correlationId` | Correlação. Não deduplica documento. |
| identificador da entrada Redis | Permanece associado para a primeira tarefa. Esta feature não confirma a entrada. |

## Configuração da execução

Sem padrão no código. Ausência e lista vazia não são o mesmo caso.

| Dado | Ausente | Presente |
|------|---------|----------|
| URL base | A recuperação não começa. Efeito não concluído. | Base do POST. |
| Usuário e senha | A recuperação não começa. Efeito não concluído. | HTTP Basic. Não aparecem em log nem em falha. |
| Timeout de cada consulta | A recuperação não começa. Efeito não concluído. | Prazo do cliente HTTP. Não é prazo total da paginação. |
| Tamanho de página | A recuperação não começa. Efeito não concluído. | `query.size`. Inteiro positivo. |
| Lista de campos da cópia física | A recuperação não começa. Efeito não concluído. | Vazia: qualquer diferença é conteúdo divergente. Com nomes: só esses campos podem diferir numa cópia. |
| Expressão de ordenação | A recuperação não começa. Efeito não concluído. | `ORDER BY` da consulta. O código não escolhe a coluna. |
| Nome do campo de tenant | Consulta não é enviada. Classificação isolamento. Efeito concluído. Nenhum documento válido. | Filtro e conferência de cada documento. |
| Nome do campo de produto | Mesma regra do campo de tenant. Os dois precisam estar presentes para a consulta sair. | Filtro e conferência. |

Organização e stream não estão nesta tabela. Vêm da referência.

## Documento físico

Um item de `hits`. Vários podem compartilhar `lhamalog_record_id`.

| Situação | Efeito no conjunto lógico |
|----------|---------------------------|
| Sem `lhamalog_record_id` | Fora da contagem. Entra na falha documento sem identificador, se essa for a classificação emitida. |
| Mesmo identificador e diferença só na lista de cópia | Um registro lógico. O entregue é uma das cópias. |
| Mesmo identificador e diferença fora da lista | Nenhuma versão é escolhida. Conteúdo divergente, se essa for a classificação emitida. |
| Tenant ou produto diferente da referência | Nenhum documento do lote sai como válido. Isolamento, se essa for a classificação emitida. |

## Registro lógico

O que a contagem compara com `recordCount` e o que o sucesso devolve. Um por `lhamalog_record_id` válido, depois da deduplicação. O documento entregue é o armazenado, sem campo oculto nesta etapa.

## Resultado

Uma classificação e a marca de efeito concluído. Falha não leva conteúdo de log nem credencial. Quando a contagem foi calculada, leva a quantidade esperada e a encontrada.

| Classificação | Efeito concluído | Documentos no resultado |
|---------------|------------------|-------------------------|
| lote completo | sim | registros lógicos e a correlação |
| falha de autenticação ou autorização | sim | nenhum |
| resposta inválida | sim | nenhum |
| inconsistência de identidade ou isolamento | sim | nenhum |
| documento sem identificador | sim | nenhum |
| conteúdo divergente | sim | nenhum |
| excesso de registros distintos | sim | nenhum; com as duas quantidades |
| lote incompleto | não | nenhum; com as duas quantidades |
| indisponibilidade ou timeout | não | nenhum |
| falha no meio da paginação | não | nenhum |
| configuração de acesso ausente | não | nenhum; a consulta não começa |

Precedência quando há mais de uma falha. Sai uma só, a primeira que se aplicar:

1. falha de autenticação ou autorização
2. resposta inválida
3. inconsistência de identidade ou isolamento
4. documento sem identificador
5. conteúdo divergente
6. excesso de registros distintos
7. lote incompleto

Indisponibilidade, timeout e falha de página só são emitidas quando nenhuma das sete ocorreu. Configuração ausente não é uma dessas classes: a recuperação não começa.

Lote completo só quando a paginação terminou, a quantidade distinta válida é igual a `recordCount` e nenhuma das sete falhas está presente.

## Transição da leitura

```text
configuração de acesso ausente → não consulta → sem efeito concluído
nomes de tenant ou produto ausentes → não consulta → isolamento, efeito concluído
página com hits = size → from aumenta de size → outra consulta
página com hits < size → deduplica e compara com recordCount
autenticação recusada ou resposta inválida → encerra
falha de documento já vista e página seguinte indisponível → classificação da falha de documento, efeito concluído, sem lote parcial
```

Não há estado de máximo de páginas.
