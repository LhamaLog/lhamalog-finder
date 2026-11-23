# Auditoria das 31 tarefas previamente marcadas

As 31 tarefas foram conferidas antes do ajuste. A revisão cruzada posterior acrescentou regressões e demonstrou uma lacuna em FR-010/T012/T015: JSON com sufixo inválido era aceito. Correção mínima exige EOF após o objeto; nenhuma reimplementação. Marcadores de checklist preservados.

| Tarefa | Código/prova conferidos | Resultado |
|---|---|---|
| T001 | go.mod: go-redis v9.23.0; sem dependências novas | Validada no escopo da tarefa |
| T002 | docker-compose.yml e .env.example: stream e nomes externos | Validada no escopo da tarefa |
| T003 | internal/config/config.go: Carregar e ValidarConsumo | Validada no escopo da tarefa |
| T004 | TestCarregarNaoCriaParametrosDeConsumo / TestCarregarIgnoraArquivoEnvComoPadrao | Validada no escopo da tarefa |
| T005 | provaLote/nomesRedis: opt-in e falha com configuração ausente | Validada no escopo da tarefa |
| T006 | TestEnvelopeValidoComESemTrace / TestEnvelopeOpcionalAusenteContinuaValido | Validada no escopo da tarefa |
| T007 | TestInterpretacaoEntregaSemConfirmar | Validada no escopo da tarefa |
| T008 | internal/domain/envelope.go: InterpretarEnvelope | Validada no escopo da tarefa |
| T009 | internal/application/port/out/efeito.go: Efeito | Validada no escopo da tarefa |
| T010 | internal/adapter/in/queue/redisStream.go: Lote.Ler | Validada no escopo da tarefa |
| T011 | internal/application/consumir.go: ConsumoReferencia.Executar | Validada no escopo da tarefa |
| T012 | TestEnvelopeInvalido + TestEnvelopeRejeitaJSONComSufixoInvalido (T032) | Lacuna de sufixo JSON reproduzida e corrigida após auditoria inicial |
| T013 | TestEntradaInvalidaConfirmaUmaVezSemConcluir | Validada no escopo da tarefa |
| T014 | TestInvalidaSaiDaPendenciaSemChaveConcluida | Validada no escopo da tarefa |
| T015 | InterpretarEnvelope/ErroEnvelope + teste Redis JSON inválido (T032) | Motivo seguro; fim de JSON reforçado após regressão |
| T016 | Lote.EncerrarInvalida; prova Redis | Validada no escopo da tarefa |
| T017 | ConsumoReferencia.encerrarInvalida; teste de registro sem bruto | Validada no escopo da tarefa |
| T018 | TestChaveReservaConcluiENaoRenova | Validada no escopo da tarefa |
| T019 | TestSegundaEntradaEsperaEEfeitoConcluidoNaoRepete / TestQuedaAntesEDepoisDoEfeito | Validada no escopo da tarefa |
| T020 | TestChaveEmProcessamentoSemPrazo / TestChaveConcluidaDuraVinteCincoHorasSemRenovar | Validada no escopo da tarefa |
| T021 | internal/domain/chave.go: PrazoConcluido/Reservar/Concluir/RemoverReserva | Validada no escopo da tarefa |
| T022 | internal/application/port/out/chave.go | Validada no escopo da tarefa |
| T023 | TestSegundaEntradaNaoRepeteEfeito; chave sob prefixo de teste | Validada no escopo da tarefa |
| T024 | ConsumoReferencia.executarEfeito/retomar; testes de queda | Validada no escopo da tarefa |
| T025 | TestValidarRecusaCadaParametroNumerico | Validada no escopo da tarefa |
| T026 | TestGrupoLeEntradaAnteriorSemDolar / TestDoisGruposObservamAMesmaEntrada / TestRedisFechadoNaoConfirma / TestReivindicacaoRespeitaQuantidade | Validada no escopo da tarefa |
| T027 | Lote: grupo 0, COUNT/BLOCK/ociosidade externos; provas Redis | Validada no escopo da tarefa |
| T028 | cmd/finder/main.go: ValidarConsumo antes de NovaLote; senha externa | Validada no escopo da tarefa |
| T029 | internal/adapter/out/openobserve/doc.go: sem cliente na branch 004 | Validada no escopo da tarefa |
| T030 | quickstart comparado com baseline/race/redis-completo | Validada no escopo da tarefa |
| T031 | TestFeedbackProcessamentoNoTerminal; slog.SetDefault JSON em main | Validada no escopo da tarefa |

Comandos e saídas completos: `evidencias/baseline.txt`, `evidencias/race.txt`, `evidencias/redis-completo.txt` (exit 0). A primeira execução `redis.txt` falhou por ausência de FINDER_REDIS_CODIGOS exigido pelos testes legados; preservada como evidência.

SC-001/002: referências válidas/sem trace e lotes unitário/múltiplo; SC-003: inválida confirmada; SC-004: reserva, TTL e reentrega; SC-005: início 0 e grupos distintos; SC-006: Redis fechado sem ACK; SC-007: reivindicação; SC-008: registro e identidade; SC-009: recusa dos parâmetros ausentes. Cobertura nos testes mapeados acima. A prova usa Redis 8.2.10 descartável em 127.0.0.1:16379; valores numéricos são dados de teste. Sem homologação de produção.

T032: json-envelope-red.txt e json-envelope-redis-red.txt falharam antes do reparo; json-envelope-green-suite.txt passou com race e Redis real (exit 0). As 31 marcas anteriores foram auditadas antes deste ajuste adicional.
