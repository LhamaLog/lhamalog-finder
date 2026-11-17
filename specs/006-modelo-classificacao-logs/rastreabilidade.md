# Rastreabilidade e estado das tarefas

Classificação pura, tipos compartilhados, resumo seguro, contadores, catálogo opaco e portas de consulta/envio. Assinatura/partição reutilizam a 008; contexto seguro vem da 010. Falhas não concluem; snapshot/chave/versionamento testados.

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
| T009 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T010 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T011 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T012 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T013 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T014 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T015 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T016 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T017 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T018 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T019 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T020 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T021 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T022 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T023 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T024 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T025 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T026 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T051 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T027 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T053 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T028 | Parcial | Modelo e trechos seguros implementados; semântica de apresentação e N3 A não formalizadas |
| T029 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T030 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T031 | Bloqueada | DN-009-2: cliente HTTP/JWT não publicado |
| T032 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T033 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T054 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T055 | Bloqueada | DN-009-2: consulta HTTP por achado |
| T034 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T035 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T036 | Bloqueada | 006-7 D: sintaxe/artefato do catálogo não definidos |
| T037 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T038 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T039 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T040 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T041 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T042 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T043 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T044 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T045 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T046 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T047 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T048 | Parcial | Snapshot de versão testado; wiring de configuração no processo aguarda pipeline/contratos |
| T049 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T050 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T052 | Parcial | Cinco literais testados na 008; acréscimos formais 010-27 A ainda indisponíveis |
| T056 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |

## Cobertura dos critérios de aceite

| Critério | Estado |
|---|---|
| SC-001 | Verificado no núcleo lógico: três classes exclusivas, versão, operação sem rede/LLM |
| SC-002 | Resumo/contadores/consulta por achado em fakes verificados; HTTP real e ≤2s na VM-A bloqueados |
| SC-003 | Motivos e ausência de conclusão em falha verificados; produto é falha permanente |
| SC-004/005 | Default INFO, níveis e catálogo opaco com precedência verificados; catálogo concreto 006-7 D pendente |
| SC-006 | Auditoria do código novo: sem publicação de ocorrência |
| SC-007 | Avaliador/gate 007 implementados; não aferido na base real inexistente |

## Decisões e coordenação

006-10 A permanece ação externa: alinhar LLT-197/199 (classificação/métricas), LLT-198 (base), LLT-200 (agrupamento), LLT-201 (handoff em vez de publicação pelo Finder) e LLT-202 (contexto seguro). T1 revisto prevalece: o Módulo cria/persiste a ocorrência. A constituição não foi emendada.

Avaliação offline: carregar JSONL 007 com origem verificada, validar todos os registros, selecionar validacao e comparar gold×predição. Para apontamento, usar maior classe prevista dos membros e contabilizar uma unidade. Calcular TP/FP/FN, precisão, recall e F1 por ERROR/WARNING/INFO; assert ≥0,95 para as três métricas em cada classe e rejeitar ausência de suporte. Sem conjunto real, o gate não pode ser aprovado.

## Arquivos de código/dados alterados contra Development

- `.env.example`
- `internal/adapter/out/modocorrencia/doc.go`
- `internal/application/analisar_lote.go`
- `internal/application/analisar_lote_test.go`
- `internal/application/classificar_log.go`
- `internal/application/desempenho_test.go`
- `internal/application/port/in/analisar_lote.go`
- `internal/application/port/out/assinatura.go`
- `internal/application/port/out/consultar_achado_modulo.go`
- `internal/application/port/out/enviar_resumo.go`
- `internal/config/analise.go`
- `internal/config/analise_test.go`
- `internal/domain/achado.go`
- `internal/domain/assinatura_classificacao.go`
- `internal/domain/classe_log.go`
- `internal/domain/classificacao.go`
- `internal/domain/classificacao_test.go`
- `internal/domain/consulta_achado.go`
- `internal/domain/falha_analise.go`
- `internal/domain/resumo_analise.go`
- `internal/domain/versao_regras.go`

## Testes existentes/novos associados aos arquivos alterados

- `internal/application/analisar_lote_test.go`: `TestAnalisarMascaramentoAntesDaClasseResumoEConsulta`
- `internal/application/analisar_lote_test.go`: `TestAnaliseFalhasDistinguiveisSemConsultaOuEnvio`
- `internal/application/analisar_lote_test.go`: `TestAnaliseSemConfirmacaoReenviaMesmaChave`
- `internal/application/analisar_lote_test.go`: `TestAnaliseFalhaSegurancaContinuaSemTrechos`
- `internal/application/analisar_lote_test.go`: `TestProdutoInconsistenteEhPermanenteSemRetry`
- `internal/application/analisar_lote_test.go`: `TestAnaliseParticoesECatalogoSomenteDecisorio`
- `internal/application/analisar_lote_test.go`: `TestAnaliseFakeIdempotenteEVersionamentoSemRetroatividade`
- `internal/application/desempenho_test.go`: `TestAceiteDesempenhoVMA`
- `internal/config/analise_test.go`: `TestAnaliseSemPadroesInventados`
- `internal/domain/classificacao_test.go`: `TestClassificarLogContrato`
- `internal/domain/classificacao_test.go`: `TestCatalogoOpacoSomenteSemNivel`
- `internal/domain/classificacao_test.go`: `TestClasseValidaEOrdem`

Comandos e saídas observadas estão em [evidencias](evidencias/); resultados compostos incluem manifesto de origem e SHA-256. Ausência de teste no pacote é indicada pelo próprio Go; não equivale a aceite externo.
