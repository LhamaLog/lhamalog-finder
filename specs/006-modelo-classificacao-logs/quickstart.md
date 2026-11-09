# Quickstart: Modelo de classificação de logs

Prova do comportamento de [contracts/classificar-apontamento.md](./contracts/classificar-apontamento.md) e [data-model.md](./data-model.md). Não substitui o código de produção. URI, banco e coleção abaixo são dados desta prova. O processo não os traz como padrão.

Esta sessão de planejamento não executou os comandos. Nenhum resultado de teste é afirmado aqui. A prova não usa log de produção.

## Pré-requisitos

- Docker Compose na raiz do repositório.
- MongoDB do Compose, serviço `mongo`, imagem `mongo:8.0.34`.
- Dados desta prova, fornecidos na execução:

```text
FINDER_MONGO_URI=<uri do MongoDB local, sem gravar a senha no repositório>
FINDER_MONGO_BANCO=<banco desta prova>
FINDER_MONGO_COLECAO=<coleção desta prova>
FINDER_MONGO_INTEGRATION=1
```

A senha não entra no repositório. O catálogo, quando a prova usar um, também é dado da execução e não fica no código.

## 1. Prova sem MongoDB

Na raiz:

```text
make test
```

Resultado esperado: `go test -race ./...` termina sem falha. Testes que exigem MongoDB não rodam neste alvo. A prova de domínio cobre campo `error` sem palavra na mensagem, campo `informação` com a palavra `error`, campo `fatal` com a palavra `info`, `errorless` que não é erro, `traceback` que não é informação, catálogo ignorado quando o nível já decidiu, primeira entrada válida, classe informação sem sinal, uma ocorrência para três logs em erro, alerta e informação, log ausente no meio do conjunto, nenhum log presente, armazenamento inacessível e segunda gravação que não apaga a primeira.

## 2. Arranque sem parâmetro de gravação

Repetir a chamada sem cada um destes, um por vez: URI, banco e coleção.

Resultado esperado: não há inserção. O motivo é configuração ausente. Nenhuma URI, banco ou coleção é escolhido no lugar do que faltou.

## 3. Ocorrência sem o texto

Com os três parâmetros presentes e `FINDER_MONGO_INTEGRATION=1`, classificar um apontamento de um log cujo campo é `error` e cuja mensagem contém a palavra `fatal`.

Resultado esperado: a consulta pela identidade devolve uma ocorrência de classe erro, com tenant, produto e identificação das regras. O documento não contém a mensagem nem o campo de nível. A contagem de envios do texto a um serviço externo de linguagem é zero.

## 4. Prazo observado

Na mesma execução, medir do pedido até a consulta da ocorrência.

Resultado esperado: a consulta devolve a ocorrência em até 2 segundos. O processo não usa esse número como timeout interno.

## 5. O que não muda

`GET /health` segue respondendo `{"status":"UP"}`. O consumo da fila não passa a chamar esta classificação e não confirma entrada por causa dela.
