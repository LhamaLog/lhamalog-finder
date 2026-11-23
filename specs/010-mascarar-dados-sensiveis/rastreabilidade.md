# Rastreabilidade e estado das tarefas

Tipos e contexto seguro, validação única por lote com transformador opaco, política de falha sem trechos, propagação de versão/contagem e ordem antes da classificação/assinatura. Sem regex, placeholder de e-mail, caminho ou formato de artefato inventados.

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
| T013 | Bloqueada | P-010-2/DN-010-4 C: caminho/formato/artefato não publicados |
| T014 | Parcial | Mecanismo opaco implementado/testado; regras interim concretas dependem padrões e placeholder acordados |
| T015 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T016 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T017 | Bloqueada | P-010-2: padrões oficiais não publicados |
| T018 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T019 | Bloqueada | Placeholder canônico e padrão de e-mail não definidos; não escolher valor |
| T020 | Parcial | Marcador fechado e mecanismo testados; padrão real senha/token/secret depende catálogo |
| T021 | Parcial | Transformador opaco disponível; casos reais fora do catálogo dependem artefato |
| T022 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T023 | Bloqueada | P-010-27: aceite formal SC-001/003 |
| T024 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T025 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T026 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T027 | Parcial | Handoff lógico seguro validado; contrato N3 A/HTTP final não fechado |
| T028 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T029 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T030 | Bloqueada | P-010-27: tela/e-mail/Sugestão fora deste ambiente |
| T041 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T043 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T046 | Parcial | Factory rejeita versão divergente; bootstrap/readiness real depende artefato P-010-2 |
| T031 | Parcial | Contador lógico mapeado; resumo externo do arquivo e chave final N3/ADR_102 pendentes |
| T032 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T033 | Bloqueada | P-010-27: contagens de aceite SC-004 |
| T034 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T035 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T036 | Parcial | Ordem/falha/ausência de bruto testadas; paridade real, regras de catálogo e cenário externo não executados |
| T037 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T038 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T039 | Bloqueada | Artefato ausente e confirmação do time Sugestão não recebida; nenhuma mensagem externa enviada |
| T040 | Bloqueada | P-010-2/P-010-27 e cadeia 008→009 impedem aceite formal |
| T042 | Parcial | Cliente 009 verificado na composição; handoff mascarado homologado depende de catálogo, casos oficiais e Módulo real. |
| T044 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T045 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |

## Cobertura dos critérios de aceite

| Critério | Estado |
|---|---|
| SC-001/002/003/004 | Bloqueados por catálogo/casos formais; fixtures verificam apenas mecanismo/ordem/contadores |
| SC-005 | Aceite formal permanece bloqueado explicitamente |
| SC-006 | Factory e propagação de versão verificadas; arquivo, bootstrap/readiness e paridade Sugestão pendentes |
| SC-007 | Referência ao gate de classificação 006/007; sem base real, sem alegação ≥95% |

## Decisões e coordenação

006-10 A permanece ação externa: alinhar LLT-197/199 (classificação/métricas), LLT-198 (base), LLT-200 (agrupamento), LLT-201 (handoff em vez de publicação pelo Finder) e LLT-202 (contexto seguro). T1 revisto prevalece: o Módulo cria/persiste a ocorrência. A constituição não foi emendada.

Avaliação offline: carregar JSONL 007 com origem verificada, validar todos os registros, selecionar validacao e comparar gold×predição. Para apontamento, usar maior classe prevista dos membros e contabilizar uma unidade. Calcular TP/FP/FN, precisão, recall e F1 por ERROR/WARNING/INFO; assert ≥0,95 para as três métricas em cada classe e rejeitar ausência de suporte. Sem conjunto real, o gate não pode ser aprovado.

## Arquivos de código/dados alterados contra Development

- `internal/application/ordem_seguranca_test.go`
- `internal/application/port/out/catalogo_padroes.go`
- `internal/application/seguranca.go`
- `internal/application/seguranca_test.go`
- `internal/domain/seguranca.go`

## Testes existentes/novos associados aos arquivos alterados

- `internal/application/ordem_seguranca_test.go`: `TestFronteira010AntesDaClassificacaoEAssinaturaCanonica`
- `internal/application/seguranca_test.go`: `TestBootstrapSegurancaExigeVersaoCombinada`
- `internal/application/seguranca_test.go`: `TestMontarContextoUmaValidacaoSemReusarBruto`
- `internal/application/seguranca_test.go`: `TestValidarSegurancaDelegaCatalogoSemInventarRegex`

Comandos e saídas observadas estão em [evidencias](evidencias/); resultados compostos incluem manifesto de origem e SHA-256. Ausência de teste no pacote é indicada pelo próprio Go; não equivale a aceite externo.
