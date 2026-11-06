# Plano de teste automatizado — Corrigir entrega da fila

Executor: agente. Não pergunta, não altera código de produção, não cria commit e não faz merge. Trabalha na raiz do repositório. Cada caso termina em `PASSOU` ou `FALHOU`. Um `FALHOU` não cancela os casos seguintes, exceto quando a preparação do ambiente falha.

Feature: `specs/002-corrigir-entrega-fila`. Spec: `spec.md`. Contrato: `contracts/entrega-fila.md`. Os casos da feature 001 que o contrato mantém entram neste plano.

## 1. Ambiente

| Item | Valor |
|------|--------|
| Raiz | diretório do repositório, o mesmo de `go.mod` |
| Shell | PowerShell |
| Go da prova com corrida | imagem `golang:1.27.0` |
| Redis | serviço `redis` do `docker-compose.yml`, contêiner `lhamalog-finder-redis`, porta do host `6379` |
| Consulta de saúde | `http://127.0.0.1:8085/health` |
| Fuso do relatório | America/Sao_Paulo |
| Evidências | `%TEMP%\qa-lhamalog-002\`. Não gravar evidência dentro do repositório |

Nomes externos, só no ambiente de teste. O código de produção não tem esses padrões.

```text
FINDER_REDIS_ADDR=127.0.0.1:6379
FINDER_REDIS_STREAM=fila:logs
FINDER_REDIS_GROUP=finder
FINDER_REDIS_CODIGOS=finder:codigos
FINDER_HTTP_ADDR=:8085
FINDER_REDIS_INTEGRATION=1
```

`docker compose run` nesta instalação não aceita `--network`. A prova com Redis usa `docker run --network container:lhamalog-finder-redis`.

## 2. Preparação

Executar uma vez, nesta ordem. Se um passo falhar, parar e marcar o plano inteiro como `BLOQUEADO`.

1. `docker info`. Se o daemon estiver parado, iniciar o Docker Desktop e repetir até `docker info` sair com código 0. Limite: 3 minutos.
2. `docker compose up -d redis`. Esperar o healthcheck do serviço `redis` ficar saudável. Limite: 2 minutos.
3. Confirmar que o contêiner `lhamalog-finder-redis` está em execução.
4. Criar `%TEMP%\qa-lhamalog-002`.
5. Registrar `git rev-parse HEAD` e a data local nesse diretório, arquivo `ambiente.txt`.

## 3. Regras de julgamento

- Código de saída 0 e ausência de `FAIL` e `panic` no log: `PASSOU`.
- Código diferente de 0, `FAIL`, `panic` ou corpo diferente do esperado: `FALHOU`.
- Com `FINDER_REDIS_INTEGRATION=1`, a linha `--- SKIP` num teste Redis deste plano é `FALHOU`. Sem essa variável, o mesmo skip é o comportamento esperado da suíte geral.
- Não há prazo numérico de produto. Espera de healthcheck, de Docker e de a porta abrir são só da infraestrutura de teste.
- Não abrir o OpenObserve. Não criar fila de erro. Não escolher nível.

## 4. Comandos reutilizáveis

Suíte sem Redis, na raiz:

```text
docker compose run --rm --no-deps gotools test -race -count=1 ./...
```

Um teste sem Redis. Substituir `PACOTE` e `TESTE`:

```text
docker compose run --rm --no-deps gotools test -race -count=1 -run TESTE PACOTE
```

Um teste com Redis. Substituir `PACOTE` e `TESTE`. O volume é a raiz do repositório:

```text
docker run --rm --network container:lhamalog-finder-redis -v "RAIZ:/src" -w /src -e FINDER_REDIS_INTEGRATION=1 -e FINDER_REDIS_ADDR=127.0.0.1:6379 -e FINDER_REDIS_STREAM=fila:logs -e FINDER_REDIS_GROUP=finder -e FINDER_REDIS_CODIGOS=finder:codigos golang:1.27.0 go test -race -count=1 -run TESTE PACOTE
```

`RAIZ` é o caminho absoluto da raiz. No PowerShell, `%CD%` não expande; usar `(Get-Location).Path`.

## 5. Casos

### Sem Redis

| ID | Comando | Resultado esperado |
|----|---------|--------------------|
| CT01 | `-run TestExecutarConcluiEDisponibilizaUmaVez ./internal/application/` | Um item aceito produz uma conclusão, o código fica disponível uma vez e o item não volta a livre |
| CT02 | `-run 'TestReservaExclusivaELiberacaoSemPrazo\|TestSegundoTomadorRecebeOutroItem' ./internal/domain/` | Outra instância não recebe o item em passo. A parada da instância devolve o item sem prazo |
| CT03 | `-run 'TestExecutarParadaDoWorkerLibera\|TestFimParadaDoWorkerNaoConclui' ./internal/application/ ./internal/domain/` | Worker parado no meio não conclui e o item volta a livre |
| CT04 | `-run TestFimConcluiNoMesmoPasso ./internal/domain/` | Com worker e instância no ar, o passo termina em conclusão |
| CT05 | `-run 'TestExecutarSegundoCodigoNaoConclui\|TestFimNaoConcluiSeCodigoJaConcluido' ./internal/application/ ./internal/domain/` | O segundo item do mesmo código não é concluído e não permanece reservado |
| CT06 | `-run TestConcluirJuntosUmaConclusao ./internal/domain/` | Dois itens do mesmo código que terminam juntos produzem uma conclusão e nenhuma reserva restante |
| CT07 | `-run 'TestExecutarFilaVaziaEFalhaDeLeitura\|TestSemConsumoNaoConclui' ./internal/application/ ./internal/domain/` | Fila vazia ou falha de leitura não conclui |
| CT08 | `-run 'TestExecutarConteudoExtraNaoConclui\|TestItemForaDeContrato\|TestItemDaMensagem' ./internal/application/ ./internal/domain/ ./internal/adapter/in/queue/` | Item sem código ou com conteúdo extra não é aceito |
| CT09 | `-run TestExecutarParcialDepoisCompletaSemSegundaContagem ./internal/application/` | A falha deixa o item parcial. A tentativa seguinte não aumenta a contagem e tira a reserva |
| CT10 | `-run TestExecutarDescartaSemPrender ./internal/application/` | Dois itens fora do contrato são descartados, o registro contém `item descartado` e não contém `item concluído`, e o item aceito seguinte conclui |
| CT11 | `-run TestDescarteMapaVazioNaoFicaNoMapa ./internal/adapter/in/queue/` | Mensagem sem valores sai do mapa `porID`. A mensagem aceita seguinte conclui. Este teste não exige Redis |
| CT12 | `-run TestLiberarDevolveOMesmoItem ./internal/adapter/in/queue/` | Liberar devolve o mesmo item sem nova leitura |
| CT13 | `-run TestNovaSemNomes ./internal/adapter/in/queue/` | Sem stream, grupo ou conjunto, a abertura devolve `ErrNomesAusentes` |
| CT14 | `-run 'TestHealth\|TestServidorAtendeHealth' ./internal/adapter/in/http/ ./cmd/finder/` | `GET /health` responde 200 e o corpo é exatamente `{"status":"UP"}` |
| CT15 | `-run 'TestCarregarPadrao\|TestCarregarWorkers\|TestAplicarArquivoEnv\|TestAplicarArquivoEnvAusente' ./internal/config/` | A configuração lê o ambiente. O endereço HTTP padrão é `:8085`. Workers inválido não substitui o padrão |
| CT16 | Suíte sem Redis, comando da seção 4, pacote `./...` | Código 0. Testes Redis deste plano podem aparecer como `SKIP`. Os demais pacotes com teste passam |

CT02, CT03, CT05 e CT08 usam dois pacotes. Rodar o `-run` uma vez em cada pacote. Os dois precisam passar.

### Com Redis

Subir o Redis na preparação. Usar o comando com `FINDER_REDIS_INTEGRATION=1`. `SKIP` é falha.

| ID | `-run` e pacote | Resultado esperado |
|----|-----------------|--------------------|
| CT17 | `TestConclusaoParcialEFalhaNoRedis ./internal/adapter/in/queue/` | Código já no conjunto e mensagem pendente: a operação deixa um membro e zera a pendência. Contexto cancelado não tira a pendência. A chamada seguinte conclui |
| CT18 | `TestInstanciaParadaDevolveNoLugarOriginal ./internal/adapter/in/queue/` | A instância que para antes de qualquer efeito devolve o item da frente, na mesma identidade, antes do item que já estava atrás |
| CT19 | `TestConsumidorVivoNaoPerdeOItem ./internal/adapter/in/queue/` | Consumidor ainda conectado não entrega o item dele à outra instância |
| CT20 | `TestOrdemDeEntradaComUmaInstancia ./internal/adapter/in/queue/` | Uma instância reserva `c`, depois `a`, depois `b`, que é a ordem de entrada, não a ordem alfabética |
| CT21 | `TestDuasInstanciasReservamItensDiferentes ./internal/adapter/in/queue/` | Duas instâncias reservam itens diferentes ao mesmo tempo |
| CT22 | `TestMesmoCodigoUmaConclusao ./internal/adapter/in/queue/` | Dois itens com o código `mesmo` deixam um membro no conjunto, o texto `mesmo` sem sufixo, e nenhuma pendência |
| CT23 | `TestDescarteInvalidoEItemAceito ./internal/adapter/in/queue/` | Código vazio e mensagem com dois campos são descartados. O mapa fica vazio. O item `valido` conclui. O conjunto tem um membro e a pendência fica em zero |
| CT24 | Os sete `-run` de CT17 a CT23 num único `go test` do pacote `./internal/adapter/in/queue/` | Código 0 e nenhum `SKIP` |

### Processo no ar

| ID | Passos | Resultado esperado |
|----|--------|--------------------|
| CT25 | `docker build -t lhamalog-finder .`. Remover `finder-plano-teste` se existir. Subir `docker run -d --name finder-plano-teste --network container:lhamalog-finder-redis` com `FINDER_REDIS_ADDR=127.0.0.1:6379`, `FINDER_REDIS_STREAM=fila:logs`, `FINDER_REDIS_GROUP=finder`, `FINDER_REDIS_CODIGOS=finder:codigos` e `FINDER_HTTP_ADDR=:8085`. Esperar até 15 segundos e consultar `http://127.0.0.1:8085/health` | Status 200 e corpo `{"status":"UP"}`. O log contém `[LhamaLog] lhamalog-finder iniciado` e `versao` `dev` |
| CT26 | `docker rm -f finder-plano-teste`. Consultar de novo `http://127.0.0.1:8085/health` com timeout de 3 segundos | A consulta não é atendida |

Se CT25 falhar, executar CT26 mesmo assim para não deixar o processo no ar.

### Conferências estáticas

Rodar na raiz. A busca é em `cmd\` e `internal\`, excluindo `_test.go`.

| ID | O que procurar | Resultado esperado |
|----|----------------|--------------------|
| CT27 | `time.Sleep`, `XAUTOCLAIM` e `MinIdle:` com valor diferente de zero em `internal\adapter\in\queue\redisStream.go` | Nenhuma ocorrência. O único `MinIdle` é `0` |
| CT28 | `Very High`, `Very Low`, `organizacao` e `tenant` em `internal\domain` e `internal\adapter\in\queue` | Nenhuma ocorrência no código de produção |
| CT29 | `.env.example` | Seis linhas `CHAVE=valor`, sem texto depois de `:`. As chaves são `FINDER_HTTP_ADDR`, `FINDER_WORKERS`, `FINDER_REDIS_ADDR`, `FINDER_REDIS_STREAM`, `FINDER_REDIS_GROUP` e `FINDER_REDIS_CODIGOS` |
| CT30 | `README.md` | Promete `GET /health` com `{"status":"UP"}` e o log de arranque sem `porta`, `go` e `workers` |

## 6. Fora deste plano

Leitura de log no OpenObserve, testes de `internal/adapter/out/openobserve` e `internal/application/capturar_test.go`, ordenação por nível dentro do Finder, fila de erro e prazo numérico de produto.

## 7. Relatório

Gravar `%TEMP%\qa-lhamalog-002\relatorio.md` com uma linha por caso:

```text
| ID | Resultado | Código de saída | Arquivo do log |
```

No fim: total de `PASSOU`, total de `FALHOU` e o commit. O plano é aprovado só quando CT01–CT30 estão `PASSOU`.

Em falha, guardar o log completo do comando no diretório de evidências com o nome do caso. Não gravar vídeo salvo se o comando não puder ser reproduzido por log.
