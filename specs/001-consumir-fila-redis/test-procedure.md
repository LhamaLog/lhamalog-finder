# Procedimento de Teste — 001-consumir-fila-redis

Feature: consumo da fila (reserva exclusiva e conclusão no mesmo passo). Spec: `spec.md`, contrato: `contracts/consumo-fila.md`, prova: `quickstart.md`, tarefas: `tasks.md`.

## 1. Pré-requisitos e ambiente
- Go 1.27.0 (local) **ou** Docker + Docker Compose (container `gotools`, Go 1.27.0).
- `make` opcional (sem make, usar os comandos `docker compose` equivalentes).
- `-race` exige CGO (gcc) no Windows; sem gcc, rodar `go test` sem `-race` localmente e `make test` (Linux/container) para a corrida.
- Gravação de vídeo (em caso de erro): ffmpeg (`gdigrab` no Windows) ou gravador de tela equivalente (Xbox Game Bar / OBS).

## 2. Como subir o projeto
- Testes: `make test` ou `docker compose run --rm --no-deps gotools test -race ./...` (local: `go test -v -count=1 ./...`).
- Lint: `make lint` (`go vet ./...`).
- App: `make run` (sobe mongo, redis, openobserve e finder). **Observação:** o adaptador Redis não está conectado ao caso de uso nesta feature (T014 bloqueada), então não há prova caixa-preta com Redis real.

## 3. Dados de teste
- Itens com conteúdo = somente o código de apontamento (ex.: `cod-1`, `cod-2`).
- Dois itens com o mesmo código (`cod-1` duas vezes).
- Item fora do contrato (sem código / com conteúdo extra).
- Instâncias: `instancia-a`, `instancia-b`. Substitutos (fakes) das portas fila, códigos concluídos e disponibilidade.
- **Não** usar `fila:logs`, campo `id` nem id `110` (FR-007).

## 4. Casos de teste
| ID | Passos | Resultado esperado | Requisitos | Teste automatizado |
|----|--------|--------------------|------------|--------------------|
| CT01 | Item com só o código; worker e instância no ar; executar consumo | 1 conclusão, item fora da pendência, código disponível 1 vez, OpenObserve não aberto | FR-002, FR-003, FR-005, SC-002, SC-003 | `TestExecutarConcluiEDisponibilizaUmaVez`, `TestFimConcluiNoMesmoPasso` |
| CT02 | Passo em andamento; outra instância tenta tomar item | Outra instância não recebe o mesmo item (outro item ou nenhum) | FR-004, FR-006, SC-001 | `TestReservaExclusivaELiberacaoSemPrazo`, `TestSegundoTomadorRecebeOutroItem` |
| CT03 | Worker para no meio do passo (instância no ar) | Não concluído, item livre de novo, sem prazo numérico | FR-006, SC-004 | `TestExecutarParadaDoWorkerLibera`, `TestFimParadaDoWorkerNaoConclui` |
| CT04 | Instância para no meio do passo; fato informado ao domínio | Item livre para outra instância, sem prazo | FR-006, SC-004 | `TestReservaExclusivaELiberacaoSemPrazo` |
| CT05 | Dois itens com o mesmo código, em sequência | 1 conclusão; o segundo nem concluído nem reservado | FR-009, SC-005 | `TestExecutarSegundoCodigoNaoConclui`, `TestFimNaoConcluiSeCodigoJaConcluido` |
| CT06 | Dois itens com mesmo código terminam juntos | 1 conclusão, 0 reservas restantes (vencedor não verificado) | FR-009, SC-005, CHK011, CHK024 | `TestConcluirJuntosUmaConclusao` |
| CT07 | Fila vazia / falha de leitura | 0 conclusões na tentativa | Edge cases, CHK015, CHK016, CHK020 | `TestExecutarFilaVaziaEFalhaDeLeitura`, `TestSemConsumoNaoConclui` |
| CT08 | Item sem código / conteúdo extra | Nenhuma conclusão por regra nova | Edge cases, CHK019 | `TestExecutarConteudoExtraNaoConclui`, `TestItemForaDeContrato` |
| CT09 | Concorrência com `-race` | Sem data race | FR-004, SC-001 | `make test` |
| CT10 | Escopo: `git diff --stat` em `internal/adapter/`, `cmd/finder/main.go`; grep por `fila:logs`, `time.` no domínio | Sem alterações fora do escopo; sem prazo/nomes inventados | FR-007, FR-008, T014, T016 | revisão manual |
| CT11 | `GET /health` | Inalterado, responde OK | Contrato (sem endpoint novo) | `TestHealth` |
| BLOQ | Duas instâncias reais na mesma fila Redis | — | FR-001 | Bloqueado até informarem stream/grupo/campos |

## 5. Registro de evidências
- Salvar saídas de comando (`go test -v`, `git diff --stat`) em pasta temporária fora do repo (ex.: `%TEMP%\qa-lhamalog\`).
- **Em caso de erro:** gravar vídeo da reprodução: `ffmpeg -f gdigrab -framerate 15 -i desktop %TEMP%\qa-lhamalog\CTxx.mp4` enquanto executa o caso; anotar ID do caso, commit (`git rev-parse HEAD`), data/hora (America/Sao_Paulo) e saída.
- Nunca salvar evidências dentro do repositório.

## 6. Critérios de aprovação
- CT01–CT08 e CT11 passam; CT09 passa no container (`make test`).
- CT10 sem divergência de escopo (ou divergência justificada e aceita pelo solicitante).
- Nenhuma task marcada `[X]` sem implementação correspondente.
- BLOQ registrado como pendência, não como falha.
