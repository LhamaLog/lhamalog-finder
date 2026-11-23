# Rastreabilidade e estado das tarefas

Loader JSONL estrito, validação de gold único e chaves duplicadas, ponte ao classificador real 006, métricas por classe e gate. Sem corpus real, autoria ou manifesto; ORIGEM.md é formulário pendente. Fixtures não são base publicada nem prova de qualidade do produto.

[tasks.md](tasks.md) mantém a descrição e a sequência das tarefas. [X] significa trabalho verificado, não aceite de produto. Tarefas parciais/bloqueadas continuam [ ]. Destinos documentais originalmente em checklists são registrados neste arquivo, preservando os artefatos dos revisores. Nenhuma issue Jira foi alterada.

| Tarefa | Estado | Evidência / pendência |
|---|---|---|
| T001 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T002 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T003 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T004 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T005 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T006 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T007 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T008 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T056 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T059 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T009 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T010 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T011 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T012 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T013 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T014 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T015 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T016 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T017 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T018 | Bloqueada | Sem corpus, manifesto e mínimo acordado verificável; não inventar contagens |
| T019 | Bloqueada | Origem/anonimização/responsável não fornecidos; ORIGEM.md é formulário pendente |
| T020 | Bloqueada | Sem logs reais anonimizados disponibilizados |
| T021 | Bloqueada | Sem T020; não há contagens reais para publicar |
| T022 | Bloqueada | Sem manifesto/contagens reais |
| T023 | Bloqueada | Método/autoria da rotulação não fornecidos |
| T024 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T025 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T026 | Bloqueada | Validador implementado; inexistem registros reais para atestar cobertura |
| T027 | Bloqueada | Esquema lógico documentado; falta artefato real T020 |
| T028 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T029 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T030 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T031 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T032 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T033 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T034 | Bloqueada | 006-7 D: catálogo e sintaxe indisponíveis |
| T035 | Bloqueada | Sem exemplos reais anonimizados e rotulados; fixtures não publicadas como base |
| T036 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T037 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T038 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T039 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T040 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T057 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T058 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T041 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T042 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T043 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T044 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T045 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T046 | Parcial | Loader/mapeamento/métricas verificados em fixtures; origem e contagens reais não podem ser verificadas |
| T047 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T048 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T049 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T050 | Parcial | Procedimento de substituição documentado; migração real depende T020 e coordenação 006 |
| T051 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T054 | Bloqueada | 006-7 D: injeção de catálogo real não definida |
| T052 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T053 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T055 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |

## Cobertura dos critérios de aceite

| Critério | Estado |
|---|---|
| SC-001 | Bloqueado: nenhuma declaração de origem/autoria validada |
| SC-002 | Rejeições mecânicas testadas; cobertura do corpus real não aferida |
| SC-003 | Bloqueado: nenhum corpus/manifesto real; sem contagens inventadas |
| SC-004 | Loader/mapeamento validam 100% das fixtures; corpus real ausente |
| SC-005 | Cálculo e gate por classe testados, apontamento = maior predição dos membros; aceite real bloqueado |
| SC-006 | MVP somente avaliação e split qualitativo documentados, sem treino/seeds/proporções inventados |

## Decisões e coordenação

006-10 A permanece ação externa: alinhar LLT-197/199 (classificação/métricas), LLT-198 (base), LLT-200 (agrupamento), LLT-201 (handoff em vez de publicação pelo Finder) e LLT-202 (contexto seguro). T1 revisto prevalece: o Módulo cria/persiste a ocorrência. A constituição não foi emendada.

Avaliação offline: carregar JSONL 007 com origem verificada, validar todos os registros, selecionar validacao e comparar gold×predição. Para apontamento, usar maior classe prevista dos membros e contabilizar uma unidade. Calcular TP/FP/FN, precisão, recall e F1 por ERROR/WARNING/INFO; assert ≥0,95 para as três métricas em cada classe e rejeitar ausência de suporte. Sem conjunto real, o gate não pode ser aprovado.

## Arquivos de código/dados alterados contra Development

- `datasets/labeled-logs/ORIGEM.md`
- `datasets/labeled-logs/README.md`
- `datasets/labeled-logs/SPLIT.md`
- `internal/labeleddataset/base_real_test.go`
- `internal/labeleddataset/doc.go`
- `internal/labeleddataset/load.go`
- `internal/labeleddataset/map_to_classifier.go`
- `internal/labeleddataset/metrics.go`
- `internal/labeleddataset/metrics_test.go`
- `internal/labeleddataset/types.go`
- `internal/labeleddataset/validate.go`
- `internal/labeleddataset/validate_test.go`

## Testes existentes/novos associados aos arquivos alterados

- `internal/labeleddataset/base_real_test.go`: `TestBaseRealRotulada`
- `internal/labeleddataset/metrics_test.go`: `TestMetricasPorApontamentoUsamMaisGrave`
- `internal/labeleddataset/metrics_test.go`: `TestGateRejeitaVazioOuClasseSemSuporte`
- `internal/labeleddataset/metrics_test.go`: `TestGateTresClassesEmFixtureNaoEhAceiteDaBase`
- `internal/labeleddataset/validate_test.go`: `TestLoadValidaCemPorCentoERejeitaCampos`
- `internal/labeleddataset/validate_test.go`: `TestMembroNuncaTemGoldMesmoNulo`
- `internal/labeleddataset/validate_test.go`: `TestPIIHeuristicaRejeitaSemPublicarBase`
- `internal/labeleddataset/validate_test.go`: `TestGoldDuplicadoNaoPodeEscolherUltimoValor`

Comandos e saídas observadas estão em [evidencias](evidencias/); resultados compostos incluem manifesto de origem e SHA-256. Ausência de teste no pacote é indicada pelo próprio Go; não equivale a aceite externo.
