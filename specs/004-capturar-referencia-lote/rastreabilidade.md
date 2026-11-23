# Rastreabilidade e estado das tarefas

31 tarefas pré-existentes auditadas antes de qualquer ajuste. Regressão posterior comprovou aceitação de sufixo JSON inválido; correção mínima exige EOF após o objeto. T032 adicionada com testes de domínio e Redis escritos antes da correção. Nenhuma reimplementação.

[tasks.md](tasks.md) mantém a descrição e a sequência das tarefas. [X] significa trabalho verificado, não aceite de produto. Tarefas parciais/bloqueadas continuam [ ]. Destinos documentais originalmente em checklists são registrados neste arquivo, preservando os artefatos dos revisores. Nenhuma issue Jira foi alterada.

| Tarefa | Estado | Evidência / pendência |
|---|---|---|
| T001 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T002 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T003 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T004 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T005 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T006 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T007 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T008 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T009 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T010 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T011 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T012 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T013 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T014 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T015 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T016 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T017 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T018 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T019 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T020 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T021 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T022 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T023 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T024 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T025 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T026 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T027 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T028 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T029 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T030 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T031 | Validada pré-existente | Código existente conferido; sem reimplementação. |
| T032 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |

## Cobertura dos critérios de aceite

| Critério | Evidência e alcance |
|---|---|
| SC-001/002 | Referências unitárias/múltiplas, campos corretos, sem conteúdo de log/ACK na interpretação |
| SC-003 | Redis real: inválidas confirmadas, sem efeito/estado concluído; regressão JSON adicional |
| SC-004 | Reserva/reentrega/TTL de 25h e queda; sem repetição de efeito |
| SC-005 | Grupo com posição 0 e dois grupos distintos |
| SC-006/007 | Redis fechado sem ACK; reivindicação de pendências |
| SC-008/009 | Logs operacionais sem bruto/credencial; parâmetros ausentes impedem leitura |

Cobertura nos testes mapeados em auditoria-31.md. Prova sintética local, sem homologação de produção.

## Arquivos de código/dados alterados contra Development

- `internal/adapter/in/queue/redis_integration_test.go`
- `internal/domain/envelope.go`
- `internal/domain/envelope_test.go`

## Testes existentes/novos associados aos arquivos alterados

- `internal/adapter/in/queue/redis_integration_test.go`: `TestDescarteMapaVazioNaoFicaNoMapa`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestConclusaoParcialEFalhaNoRedis`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestInstanciaParadaDevolveNoLugarOriginal`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestConsumidorVivoNaoPerdeOItem`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestOrdemDeEntradaComUmaInstancia`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestDuasInstanciasReservamItensDiferentes`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestMesmoCodigoUmaConclusao`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestDescarteInvalidoEItemAceito`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestGrupoLeEntradaAnteriorSemDolar`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestDoisGruposObservamAMesmaEntrada`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestInvalidaSaiDaPendenciaSemChaveConcluida`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestChaveEmProcessamentoSemPrazo`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestChaveConcluidaDuraVinteCincoHorasSemRenovar`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestSegundaEntradaNaoRepeteEfeito`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestRedisFechadoNaoConfirma`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestReivindicacaoRespeitaQuantidade`
- `internal/adapter/in/queue/redis_integration_test.go`: `TestJSONInvalidoConfirmadoSemEfeito`
- `internal/domain/envelope_test.go`: `TestEnvelopeValidoComESemTrace`
- `internal/domain/envelope_test.go`: `TestEnvelopeOpcionalAusenteContinuaValido`
- `internal/domain/envelope_test.go`: `TestEnvelopeInvalido`
- `internal/domain/envelope_test.go`: `TestEnvelopeRejeitaJSONComSufixoInvalido`

Comandos e saídas observadas estão em [evidencias](evidencias/); resultados compostos incluem manifesto de origem e SHA-256. Ausência de teste no pacote é indicada pelo próprio Go; não equivale a aceite externo.
