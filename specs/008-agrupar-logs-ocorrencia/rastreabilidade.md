# Rastreabilidade e estado das tarefas

Normalizador/compositor canônicos, partição isolada por chave, métricas de lote, fusão de estado e recorrência na janela de 500 lotes. Política de três tentativas e porta Core testadas. Query histórica é stub explícito; nenhum resultado falso de recorrência por degradação.

[tasks.md](tasks.md) mantém a descrição e a sequência das tarefas. [X] significa trabalho verificado, não aceite de produto. Tarefas parciais/bloqueadas continuam [ ]. Destinos documentais originalmente em checklists são registrados neste arquivo, preservando os artefatos dos revisores. Nenhuma issue Jira foi alterada.

| Tarefa | Estado | Evidência / pendência |
|---|---|---|
| T001 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T002 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T003 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T004 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T005 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T006 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T046 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T047 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T007 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T008 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T009 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T010 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T011 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T012 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T013 | Parcial | Gate completo/incompleto nos serviços lógicos; ligação física de consumir.go ao pipeline aguarda extração/contratos |
| T014 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T016 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T015 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T048 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T017 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T018 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T019 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T020 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T049 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T021 | Concluída | Usar a consulta HTTPS/JWT da spec 009 por achado, preservando a revisão opaca para a entrega. Verificado na composição das branches (009 T050). |
| T022 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T023 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T050 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T024 | Concluída | Usar o POST de ResumoAnalise DN-009-2 da spec 009 e permitir recálculo após revisão desatualizada. Verificado na composição das branches (009 T018/T019). |
| T025 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T026 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T027 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T051 | Bloqueada | P-010-27 A: matriz formal indisponível |
| T028 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T029 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T030 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T031 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T032 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T033 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T034 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T035 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T036 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T037 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T038 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T039 | Parcial | Três tentativas e porta de aviso verificadas; envio real Core depende ADR_051 |
| T040 | Parcial | Orquestração lógica testada; cadeia operacional/Core real pendente |
| T041 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T042 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T043 | Parcial | Partições, fusão e recorrência com portas testadas; reabertura/Módulo e histórico OpenObserve reais indisponíveis |
| T044 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T045 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T052 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T053 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T054 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T055 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |

## Cobertura dos critérios de aceite

| Critério | Estado |
|---|---|
| SC-001 | Partição lógica única por chave testada; ocorrência persistida/matriz formal exigem contratos |
| SC-002 | Soma e min/max de timestamp do log verificados em memória |
| SC-003 | FECHADA propagado sem alteração; reabertura é do Módulo e não foi testada em produto |
| SC-004 | Três classes exclusivas verificadas |
| SC-005 | Bloqueado: N3 A, HTTP Módulo/Core, query histórica e escape de barra invertida literal |
| SC-006 | Dependência ADR_035/LLT-216 documentada, sem alegar resultado externo |
| SC-007 | Política de três falhas temporárias/aviso por porta e marca desconhecida verificada; aviso HTTP real bloqueado |

## Decisões e coordenação

006-10 A permanece ação externa: alinhar LLT-197/199 (classificação/métricas), LLT-198 (base), LLT-200 (agrupamento), LLT-201 (handoff em vez de publicação pelo Finder) e LLT-202 (contexto seguro). T1 revisto prevalece: o Módulo cria/persiste a ocorrência. A constituição não foi emendada.

Avaliação offline: carregar JSONL 007 com origem verificada, validar todos os registros, selecionar validacao e comparar gold×predição. Para apontamento, usar maior classe prevista dos membros e contabilizar uma unidade. Calcular TP/FP/FN, precisão, recall e F1 por ERROR/WARNING/INFO; assert ≥0,95 para as três métricas em cada classe e rejeitar ausência de suporte. Sem conjunto real, o gate não pode ser aprovado.

## Arquivos de código/dados alterados contra Development

- `internal/adapter/out/openobserve/recorrencia.go`
- `internal/adapter/out/openobserve/recorrencia_test.go`
- `internal/application/agrupar/agrupar_test.go`
- `internal/application/agrupar/consulta_memoria.go`
- `internal/application/agrupar/fusao_entre_lotes.go`
- `internal/application/agrupar/orquestrar.go`
- `internal/application/agrupar/particionar_lote.go`
- `internal/application/assinatura/calcular.go`
- `internal/application/assinatura/normalizar_mensagem.go`
- `internal/application/assinatura/normalizar_mensagem_test.go`
- `internal/application/port/out/consulta_ocorrencia.go`
- `internal/application/port/out/consulta_recorrencia.go`
- `internal/application/port/out/resumo_analise.go`
- `internal/application/recorrencia/avaliar.go`
- `internal/application/recorrencia/avaliar_test.go`
- `internal/config/recorrencia.go`
- `internal/config/recorrencia_test.go`
- `internal/domain/achado_lote.go`
- `internal/domain/agregacao.go`
- `internal/domain/assinatura.go`

## Testes existentes/novos associados aos arquivos alterados

- `internal/adapter/out/openobserve/recorrencia_test.go`: `TestRecorrenciaTresTentativasSemFalsoNegativo`
- `internal/adapter/out/openobserve/recorrencia_test.go`: `TestErroTemporarioPonteiroTambemRepeteTres`
- `internal/application/agrupar/agrupar_test.go`: `TestParticionarIsolaChaveEContaTimestampLog`
- `internal/application/agrupar/agrupar_test.go`: `TestFusaoEntreLotesNaoAlteraStatus`
- `internal/application/agrupar/agrupar_test.go`: `TestOrquestrarConsultaPorAchadoEnviaUmResumo`
- `internal/application/agrupar/agrupar_test.go`: `TestConsultaMemoriaCopiaEstadoEContaChamadas`
- `internal/application/assinatura/normalizar_mensagem_test.go`: `TestNormalizarMarcadoresCanonicos`
- `internal/application/assinatura/normalizar_mensagem_test.go`: `TestComposicaoCanonicaEEscape`
- `internal/application/assinatura/normalizar_mensagem_test.go`: `TestNormalizacaoDeterministicaIdempotente`
- `internal/application/recorrencia/avaliar_test.go`: `TestRecorrenciaDistintaMesmaMaquinaJanela500`
- `internal/application/recorrencia/avaliar_test.go`: `TestJanelaContaLotesNaoRegistros`
- `internal/config/recorrencia_test.go`: `TestConfiguracaoHistoricoSemDefaults`

Comandos e saídas observadas estão em [evidencias](evidencias/); resultados compostos incluem manifesto de origem e SHA-256. Ausência de teste no pacote é indicada pelo próprio Go; não equivale a aceite externo.
