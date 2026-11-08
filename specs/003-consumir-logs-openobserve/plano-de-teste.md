# Plano de teste — LLT-196 captura de logs

Plano para execução automática. O bot executa os comandos, lê a saída e preenche o relatório. Este arquivo não é o relatório.

**Branch:** `LLT-196-consumir-logs-da-fila-por-prioridade`  
**Feature:** `specs/002-consumir-logs-openobserve/`  
**Spec:** `specs/002-consumir-logs-openobserve/spec.md`  
**Contrato:** `specs/002-consumir-logs-openobserve/contracts/captura-log.md`  
**Prova de referência:** `specs/002-consumir-logs-openobserve/quickstart.md`  
**Raiz do repositório:** `C:\Users\joaop\Documents\workarea\git\LhamaLog Workspace\lhamalog-finder`  
**Shell:** PowerShell  
**Data do plano:** 2026-11-02

O fluxo coberto é só a captura implementada nesta branch: apontamento mockado, log pronto em memória da execução, falhas distinguíveis e capturas sobrepostas. A fila Redis da feature 001 entra apenas como regressão: ela não pode quebrar e não é chamada pela captura.

## 1. O que foi desenvolvido e onde

| Papel | Caminho |
|-------|---------|
| Regras do log pronto e das falhas | `internal/domain/captura.go` |
| Testes de domínio | `internal/domain/captura_test.go` |
| Caso de uso `ServicoCaptura.Executar` | `internal/application/capturar.go` |
| Testes do caso de uso com substitutos | `internal/application/capturar_test.go` |
| Entrada `PedidoCaptura` / `SaidaCaptura` | `internal/application/port/in/capturar.go` |
| Fonte mockada dos apontamentos | `internal/application/port/out/apontamentos.go` |
| Leitura do registro (pronta, ausente, inacessível) | `internal/application/port/out/leitura.go` |
| Guarda e releitura do log pronto | `internal/application/port/out/prontos.go` |

O adaptador HTTP está em `internal/adapter/out/openobserve/`. `cmd/finder/main.go` continua no consumo da fila e não dispara a captura. A base local em `docker-compose.yml` e no `Makefile` sobe OpenObserve v1.1.2 no lugar do Elasticsearch. O arquivo `otel-logs.zip` não está no repositório.

O registro dos testes é opaco. O teste de domínio usa o texto `{"senha":"segredo","msg":"inteiro"}`. O teste do caso de uso usa `{"senha":"segredo"}`. Os dois exigem o texto de volta, inteiro. Isso não cumpre o aceite com um registro real de `otel-logs.zip`.

## 2. Ambiente

1. Confirmar a branch:

```powershell
Set-Location "C:\Users\joaop\Documents\workarea\git\LhamaLog Workspace\lhamalog-finder"
git branch --show-current
git rev-parse --short HEAD
go version
```

A branch esperada é `LLT-196-consumir-logs-da-fila-por-prioridade`. O `HEAD` commitado no momento deste plano é `117f9dd` (`LLT-132 Consumo da fila`). A captura está no working tree, ainda não commitada. Se a branch for outra, parar e registrar `AMBIENTE_INVALIDO`.

2. Go local. `go version` precisa responder. A imagem do projeto é Go 1.27.0. Se o Go local for outro, registrar a versão no relatório e seguir com os testes locais.

3. Prova com corrida, nesta ordem:

```powershell
docker info
```

- Se o Docker responder, a prova de corrida é:

```powershell
make test
```

Equivalente:

```powershell
docker compose run --rm --no-deps gotools test -race ./...
```

- Se o Docker não responder, não marcar a corrida como aprovada. Registrar `BLOQUEADO_AMBIENTE` com a saída do `docker info`. Em seguida rodar a suíte local sem `-race`:

```powershell
go test -count=1 -v ./internal/domain ./internal/application ./internal/adapter/in/http ./internal/adapter/in/queue ./internal/config
```

Na sessão de implementação de 2026-11-02, `make test` falhou porque o daemon Docker não estava no ar, e `go test -race` local falhou com `cc1.exe: sorry, unimplemented: 64-bit mode not compiled in`. O bot deve medir de novo. Não copiar esse resultado como se fosse desta execução.

4. `make run` e `make infra-up` sobem MongoDB, Redis e OpenObserve. Subir a base não prova a captura de um registro de `otel-logs.zip`.

## 3. Como executar cada caso

Diretório de trabalho: a raiz do repositório.

Cada caso roda sozinho. Sucesso do caso: código de saída `0` e a linha do teste contém `PASS`. Falha: código diferente de `0`, ou a saída contém `FAIL` ou `panic`.

Salvar a saída de cada comando em `%TEMP%\qa-lhamalog-llt196\`. Criar a pasta antes. Não gravar evidência dentro do repositório.

```powershell
New-Item -ItemType Directory -Force -Path "$env:TEMP\qa-lhamalog-llt196" | Out-Null
```

Modelo de um caso:

```powershell
go test -count=1 -v -run NOME_DO_TESTE ./PACOTE *> "$env:TEMP\qa-lhamalog-llt196\CTXX.txt"
echo $LASTEXITCODE >> "$env:TEMP\qa-lhamalog-llt196\CTXX.txt"
```

## 4. Casos positivos

| ID | O que prova | Comando | Esperado |
|----|-------------|---------|----------|
| CT01 | Um apontamento com a instância no ar deixa um log pronto, igual ao registro, e a leitura seguinte devolve o mesmo texto | `go test -count=1 -v -run TestExecutarDevolveRegistroInteiro ./internal/application` | `PASS`, saída 0. Um pronto, zero falhas, conteúdo `{"senha":"segredo"}` |
| CT02 | O domínio guarda o registro inteiro e a alteração do original não muda o guardado | `go test -count=1 -v -run TestConcluirGuardaRegistroInteiro ./internal/domain` | `PASS`, saída 0 |
| CT03 | Fonte sem apontamento: zero logs prontos e zero falhas, leitura não é chamada | `go test -count=1 -v -run TestExecutarSemApontamento ./internal/application` | `PASS`, saída 0 |
| CT04 | Quadro novo começa vazio | `go test -count=1 -v -run TestQuadroSemApontamento ./internal/domain` | `PASS`, saída 0 |
| CT05 | Dois apontamentos distintos ao mesmo tempo: o que existe fica pronto e o ausente conserva a própria falha | `go test -count=1 -v -run TestExecutarSobrepostas ./internal/application` | `PASS`, saída 0. O teste espera no máximo 2 segundos só para detectar sobreposição. Esse prazo é do teste, não do produto |
| CT06 | O mesmo apontamento duas vezes produz uma leitura e um log pronto | `go test -count=1 -v -run TestExecutarApontamentoRepetido ./internal/application` | `PASS`, saída 0 |
| CT07 | No domínio, a segunda guarda do mesmo apontamento não troca o primeiro registro | `go test -count=1 -v -run TestMesmoApontamentoUmLogPronto ./internal/domain` | `PASS`, saída 0 |
| CT08 | Regressão da fila, do health e da configuração | `go test -count=1 -v ./internal/adapter/in/http ./internal/adapter/in/queue ./internal/config ./internal/domain ./internal/application` | `PASS` em todos os pacotes, saída 0 |

## 5. Casos negativos

| ID | O que prova | Comando | Esperado |
|----|-------------|---------|----------|
| CT09 | Apontamento vazio é inutilizável, a leitura não é consultada e não há log pronto | `go test -count=1 -v -run TestExecutarApontamentoVazioNaoConsulta ./internal/application` | `PASS`, saída 0. Motivo `inutilizavel` |
| CT10 | O domínio recusa concluir apontamento vazio | `go test -count=1 -v -run TestApontamentoVazioEInutilizavel ./internal/domain` | `PASS`, saída 0 |
| CT11 | Log ausente: falha `ausente`, nenhum registro vazio no lugar | `go test -count=1 -v -run TestExecutarLogAusente ./internal/application` | `PASS`, saída 0 |
| CT12 | O domínio não cria log pronto na ausência | `go test -count=1 -v -run TestAusenteNaoCriaRegistroVazio ./internal/domain` | `PASS`, saída 0 |
| CT13 | Armazenamento inacessível: falha `inacessivel`, distinta de `ausente` | `go test -count=1 -v -run TestExecutarInacessivel ./internal/application` | `PASS`, saída 0 |
| CT14 | Os dois motivos permanecem diferentes e nenhum gera log pronto | `go test -count=1 -v -run TestInacessivelDistintoDeAusente ./internal/domain` | `PASS`, saída 0 |
| CT15 | Instância já parada: a captura não consulta a leitura e não deixa log pronto | `go test -count=1 -v -run TestExecutarParadaNoMeioNaoDeixaPronto ./internal/application` | `PASS`, saída 0 |
| CT16 | Conclusão com a instância parada não guarda log | `go test -count=1 -v -run TestParadaNoMeioNaoDeixaPronto ./internal/domain` | `PASS`, saída 0 |
| CT17 | Depois que a instância para, a leitura não é exigida e isso não vira falha | `go test -count=1 -v -run TestDepoisQueInstanciaParaNaoExigeLog ./internal/domain` | `PASS`, saída 0 |

## 6. Casos de escopo

Rodar depois dos casos funcionais. Falha de escopo se algum caminho proibido aparecer.

```powershell
git status --short *> "$env:TEMP\qa-lhamalog-llt196\CT18-status.txt"
```

Caminhos permitidos nesta feature:

- `internal/domain/captura.go`
- `internal/domain/captura_test.go`
- `internal/application/capturar.go`
- `internal/application/capturar_test.go`
- `internal/application/port/in/capturar.go`
- `internal/application/port/out/apontamentos.go`
- `internal/application/port/out/leitura.go`
- `internal/application/port/out/prontos.go`
- `internal/adapter/out/openobserve/`
- `internal/config/config.go`
- `internal/config/config_test.go`
- `docker-compose.yml`
- `Makefile`
- `README.md`
- `specs/002-consumir-logs-openobserve/`
- remoção de `internal/adapter/out/elastic/`

`.gitignore` pode aparecer modificado com `.agents/` e `.specify/`. Isso é anterior a esta feature. Registrar como observação. Não reprovar CT18 só por essa linha.

Caminhos que reprovam CT18 se aparecerem em `git status --short`:

- `cmd/finder/main.go`
- `internal/application/consumir.go`
- `internal/adapter/in/`
- `internal/adapter/out/mongo/`
- `go.mod`
- `go.sum`

```powershell
Test-Path "internal/adapter/out/elastic"
Select-String -Path "docker-compose.yml","Makefile","README.md","internal/adapter/out/openobserve/cliente.go" -Pattern "elasticsearch" -SimpleMatch
Select-String -Path "internal/application/capturar.go" -Pattern "queue|consumir|golang.org/x/sync|FINDER_WORKERS|otel-logs"
go vet ./internal/domain ./internal/application ./internal/application/port/in ./internal/application/port/out
```

Esperado:

- `Test-Path` do pacote `elastic` imprime `False`
- `Select-String` não imprime linha
- `go vet` termina com saída 0 e sem texto

| ID | Esperado |
|----|----------|
| CT18 | Status só com os caminhos permitidos, mais a observação do `.gitignore` se ela existir |
| CT19 | Pacote `elastic` ausente, Compose sem Elasticsearch, captura sem fila e sem biblioteca de sync, `go vet` limpo |

## 7. Casos bloqueados

Não executar como aprovação. Registrar `BLOQUEADO`, nunca `PASS` e nunca `FAIL` por ausência do ambiente.

| ID | Por que está bloqueado | O que não fazer |
|----|------------------------|-----------------|
| BL01 | `otel-logs.zip` não está no repositório. Organização, stream, índice, campos e credenciais do OpenObserve não foram informados | Não inventar JSON, URL, porta nem credencial. Não marcar CT01 como aceite SC-001 real |
| BL02 | A base sobe OpenObserve, mas o daemon Docker pode estar parado. Subir o contêiner não coloca um registro de `otel-logs.zip` nele | Não marcar a subida da base como captura de um log real |
| BL03 | `cmd/finder/main.go` não liga `ServicoCaptura` | Não exigir log do processo Finder como prova da captura |

## 8. Ordem do bot

1. Seção 2. Se a branch divergir, encerrar com `AMBIENTE_INVALIDO`.
2. CT01 a CT17, um comando por caso, saída em arquivo próprio.
3. CT08 por último entre os funcionais, porque repete os pacotes.
4. CT18 e CT19.
5. Tentar `make test`. Se o Docker não responder, `BLOQUEADO_AMBIENTE` e a suíte local sem `-race` da seção 2.
6. Preencher o relatório. Não alterar código para fazer um caso passar.

## 9. Relatório

Gravar em `%TEMP%\qa-lhamalog-llt196\relatorio.md`. Uma linha por caso.

```markdown
# Relatório LLT-196

- Branch:
- Commit HEAD:
- Go:
- Docker:
- Data e hora (America/Sao_Paulo):

| Caso | Resultado | Saída |
|------|-----------|-------|
| CT01 | PASS ou FAIL | caminho do txt |
| ... | | |
| CT19 | PASS ou FAIL | |
| BL01 | BLOQUEADO | |
| BL02 | BLOQUEADO | |
| BL03 | BLOQUEADO | |
| CORRIDA | PASS ou BLOQUEADO_AMBIENTE | |

## Veredito

APROVADO só se CT01–CT19 estão PASS e BL01–BL03 estão BLOQUEADO.
REPROVADO se qualquer CT está FAIL.
INCOMPLETO se a corrida ficou BLOQUEADO_AMBIENTE e o restante passou. A suíte local sem -race não substitui make test.
```

Cola da saída fica nos txt da pasta temporária. O relatório só aponta o arquivo.

## 10. Aprovação

- CT01–CT19 com `PASS`.
- BL01–BL03 com `BLOQUEADO`.
- Nenhum arquivo de produto fora da lista da seção 6.
- `make test` com `PASS` quando o Docker estiver no ar. Sem Docker, o veredito máximo é `INCOMPLETO`.
