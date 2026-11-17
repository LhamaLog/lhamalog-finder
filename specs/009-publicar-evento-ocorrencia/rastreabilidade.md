# Rastreabilidade e estado das tarefas

Entrega lógica com confirmação, chave/snapshot estáveis, reenvio e DTO provisório. Reutiliza resumo 006 e preparação 008; não repete consultas nem publica evento de ocorrência. Cópia profunda inclui campos opcionais, mapas e linhas.

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
| T016 | Bloqueada | N3 A: nomes JSON ampliados sem ADR |
| T048 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T051 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T049 | Parcial | 009 impede handoff após falha na preparação; política/porta na 008; consumir.go e HTTP Core dependem contratos |
| T017 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T018 | Bloqueada | DN-009-2: cliente de entrega HTTP/JWT |
| T019 | Bloqueada | DN-009-2: confirmação HTTP não definida |
| T020 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T021 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T022 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T023 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T024 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T025 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T026 | Bloqueada | 009-22 A: lógica 008 testada, mas aceite completo 008 não fechado |
| T027 | Bloqueada | P-010-27: casos formais de linhas mascaradas indisponíveis |
| T050 | Bloqueada | DN-009-2: cliente de consulta HTTP/JWT |
| T028 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T029 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T030 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T031 | Bloqueada | 009-22 A: aceite SC-002 continua bloqueado; estado mantido no quickstart |
| T032 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T033 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T034 | Bloqueada | DN-009-2: header versus corpo da chave não definidos |
| T035 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T036 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T037 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T038 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T039 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T040 | Bloqueada | SC-002 exige 008 concluída e ambiente produto |
| T041 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T042 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T043 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T044 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T045 | Parcial | Equivalentes lógicos dos cenários 1–5 testados; roteiro manual com Módulo real não executável |
| T046 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T047 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |

## Cobertura dos critérios de aceite

| Critério | Estado |
|---|---|
| SC-F001 | Um resumo após preparação, confirmação antes do callback; verificado com portas |
| SC-F002 | Snapshot/chave estáveis no reenvio, fake sem efeito duplicado; idempotência real Módulo não comprovada |
| SC-F003 | Rejeição de classes inválidas verificada |
| SC-001/003/004/006 | Pipeline Módulo→Core→Notificação apenas documentado; serviços externos não executados |
| SC-002/005 | Bloqueados até conclusão formal 008 e contratos do Módulo |

O teste de cadeia 010→006→008→009 usa serviços de aplicação reais e portas externas fake. Não é E2E de produto nem prova de notificação real.

## Decisões e coordenação

006-10 A permanece ação externa: alinhar LLT-197/199 (classificação/métricas), LLT-198 (base), LLT-200 (agrupamento), LLT-201 (handoff em vez de publicação pelo Finder) e LLT-202 (contexto seguro). T1 revisto prevalece: o Módulo cria/persiste a ocorrência. A constituição não foi emendada.

Avaliação offline: carregar JSONL 007 com origem verificada, validar todos os registros, selecionar validacao e comparar gold×predição. Para apontamento, usar maior classe prevista dos membros e contabilizar uma unidade. Calcular TP/FP/FN, precisão, recall e F1 por ERROR/WARNING/INFO; assert ≥0,95 para as três métricas em cada classe e rejeitar ausência de suporte. Sem conjunto real, o gate não pode ser aprovado.

## Matriz FR-F01–FR-F11

| Requisito | Tarefas / situação |
|---|---|
| FR-F01 | T017–T024: resumo lógico, aceitação e conclusão; HTTP T018/T019 bloqueado |
| FR-F02 | T025: nenhuma publicação de ocorrência pelo Finder |
| FR-F03 | T023/T015: execId opcional e chave lógica |
| FR-F04 | T020/T022: INFO incluído no resumo e no fluxo lógico |
| FR-F05 | T024/T036: confirmação antes de conclusão |
| FR-F06 | T035/T037: reenvio da mesma chave; idempotência real depende Módulo |
| FR-F07 | T049: preparo impede envio na falha; política na 008, HTTP Core pendente |
| FR-F08 | T012/T017: rejeição de classes inválidas |
| FR-F09 | T002/T022/T050: consulta por achado, cliente bloqueado |
| FR-F10 | Auditoria Git: ADR_019, ADR_167, SSID_007.1 inalterados; taxonomia ADR_170 |
| FR-F11 | T048/T051: reuso canônico 008, sem duplicar algoritmo |

Draft mantido; SC-002 bloqueado 009-22 A; DN-009-2 bloqueia adaptadores.

## Arquivos de código/dados alterados contra Development

- `internal/adapter/out/modocorrencia/serializar_resumo.go`
- `internal/adapter/out/modocorrencia/serializar_resumo_test.go`
- `internal/application/handoff/entrega.go`
- `internal/application/handoff/entrega_test.go`
- `internal/application/handoff/fluxo_test.go`
- `internal/application/handoff/orquestrar_entrega.go`
- `internal/application/port/out/finalizar_analise.go`
- `internal/application/port/out/modulo_ocorrencia.go`
- `internal/domain/confirmacao_handoff.go`
- `internal/domain/idempotencia_handoff.go`

## Testes existentes/novos associados aos arquivos alterados

- `internal/adapter/out/modocorrencia/serializar_resumo_test.go`: `TestSerializarResumoProvisorioSemCamposBloqueados`
- `internal/adapter/out/modocorrencia/serializar_resumo_test.go`: `TestFlagDeAchadoSemSemanticaNaoAfirmaAusenciaDeSensivel`
- `internal/application/handoff/entrega_test.go`: `TestHandoffEsperaConfirmacaoEReenviaMesmaChave`
- `internal/application/handoff/entrega_test.go`: `TestHandoffRejeitaTaxonomiaEExecIdInvalido`
- `internal/application/handoff/entrega_test.go`: `TestOrquestrarEntregaBloqueiaLoteIncompleto`
- `internal/application/handoff/entrega_test.go`: `TestSnapshotHandoffIsolaPonteirosDaOrigemEDoAdapter`
- `internal/application/handoff/fluxo_test.go`: `TestFluxoLogico010006008009`

Comandos e saídas observadas estão em [evidencias](evidencias/); resultados compostos incluem manifesto de origem e SHA-256. Ausência de teste no pacote é indicada pelo próprio Go; não equivale a aceite externo.
