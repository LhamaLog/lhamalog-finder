# Quickstart: Recuperar logs do lote

Prova do comportamento de [contracts/recuperar-lote.md](./contracts/recuperar-lote.md) e [data-model.md](./data-model.md). Não substitui o código de produção. URL, usuário, senha, timeout, tamanho de página, nomes de campo e a expressão de ordenação abaixo são dados desta prova. O processo não os traz como padrão.

Esta sessão de planejamento não executou os comandos. Nenhum resultado de teste é afirmado aqui. A prova não usa log de produção e não declara homologação dos consumidores.

## Pré-requisitos

- Docker Compose na raiz do repositório.
- OpenObserve do Compose, serviço `openobserve`, imagem `openobserve/openobserve:v1.1.2`, publicado em `localhost:5080`.
- Dados desta prova, fornecidos na execução:

```text
FINDER_OO_URL=http://127.0.0.1:5080
FINDER_OO_USUARIO=<usuário da imagem local>
FINDER_OO_SENHA=<senha da imagem local>
FINDER_OO_TIMEOUT=2s
FINDER_OO_PAGINA=2
FINDER_OO_CAMPOS_COPIA=<campo que a prova usa só para distinguir a cópia>
FINDER_OO_ORDEM=<expressão que mantém o identificador lógico estável entre páginas>
FINDER_OO_CAMPO_TENANT=<nome persistido do tenant>
FINDER_OO_CAMPO_PRODUTO=<nome persistido do produto>
FINDER_OPENOBSERVE_INTEGRATION=1
```

A senha não entra no repositório. Organização e stream vêm da referência sintética, não destas variáveis. `lhamalog:stream:logs` não é o stream consultado.

## 1. Prova sem OpenObserve

Na raiz:

```text
make test
```

Resultado esperado: `go test -race ./...` termina sem falha. Testes que exigem OpenObserve não rodam neste alvo. A prova de domínio cobre montagem da consulta, lote unitário e lote de vários registros, faixa com fim no máximo mais 1, uma página, várias páginas e múltiplo exato do tamanho de página, cópia na mesma página e entre páginas, contagem depois da deduplicação, busca vazia, lote incompleto, excesso, documento sem identificador, conteúdo divergente, documento de outro tenant ou produto, timeout, credencial recusada, falha no meio da paginação e a precedência de uma só classificação.

## 2. Arranque da recuperação sem parâmetro

Repetir a recuperação sem cada um destes, um por vez: URL, credencial, timeout, tamanho de página, lista de cópia e expressão de ordenação.

Resultado esperado: a consulta não é enviada. Nenhum valor padrão é criado. O resultado não marca efeito concluído. A entrada Redis não é confirmada.

Lista de cópia presente e vazia é outro caso: a recuperação pode consultar, e qualquer diferença no mesmo identificador é conteúdo divergente.

## 3. Nomes de tenant e produto ausentes

Chamar a recuperação com o restante da configuração e sem um dos dois nomes.

Resultado esperado: a consulta não é enviada. A classificação é isolamento. O resultado marca efeito concluído. Nenhum documento é entregue. A entrada Redis não é confirmada por esta etapa.

## 4. Lote sintético

Gravar no OpenObserve de teste um lote aceito, com organização e stream da referência, `lhamalog_batch_id`, a faixa e os campos de tenant e produto da prova. Incluir uma cópia física do mesmo `lhamalog_record_id`, diferente só no campo da lista de cópia, e um identificador lógico que só aparece depois dessas cópias.

Resultado esperado: cada identificador lógico aparece uma vez. O identificador que ficaria de fora, se a leitura parasse em `recordCount`, está no conjunto. Se a quantidade distinta fecha e não há outra falha, o lote é completo e o resultado marca efeito concluído. A entrada Redis segue sem confirmação desta etapa.

Repetir com lote de um registro e conferir o filtro de `lhamalog_record_id`. Repetir com lote misto e conferir que a completude usa só a quantidade aceita.

## 5. Dois tenants

Gravar documento do mesmo `batchId` com outro tenant ou outro produto.

Resultado esperado: a consulta filtra os campos configurados. Nenhum documento sai como resultado válido. A classificação é isolamento, mesmo que o lote também tenha documento sem identificador. O resultado marca efeito concluído.

## 6. Falha de armazenamento

Com a configuração presente, tornar o OpenObserve inalcançável no meio da segunda página, sem falha que pede correção já encontrada.

Resultado esperado: a classificação é indisponibilidade ou timeout. O conjunto da primeira página não é lote pronto. O resultado não marca efeito concluído.

Se a primeira página já trouxe documento de outro tenant, a classificação da segunda página indisponível é isolamento e o resultado marca efeito concluído. O parcial continua de fora.
