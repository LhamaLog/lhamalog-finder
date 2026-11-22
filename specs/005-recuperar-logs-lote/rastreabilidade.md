# Rastreabilidade e estado das tarefas

A implementação e a sincronização anterior com `Development` foram publicadas no commit `735e556`, da PR #8. Seus pais são `5f82490` (implementação inicial da LLT-270) e `6e8b0fd` (Development). T031 remove o teto de 32 MiB; T032 separa falha de leitura/status de JSON de sucesso inválido. Corrigidos sufixo JSON e problemas observados nos testes de integração (assert de filtro e acesso concorrente ao coletor).

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
| T031 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |
| T032 | Concluída | Implementação/registro verificado; provas e arquivos abaixo. |

## Cobertura dos critérios de aceite

| Critério | Evidência e alcance |
|---|---|
| SC-001/002 | Consulta da referência e filtro recordId somente em lote unitário; teste real OO |
| SC-003/004 | Paginação até fim, cópias físicas, deduplicação e divergência |
| SC-005 | Vazio, abaixo/acima do recordCount sem inventar causa |
| SC-006 | Dois tenants sintéticos; filtros e conferência de identidade; isolamento distinguível |
| SC-007 | Timeout/status/autorização/corpo interrompido/JSON inválido; sem resultado parcial completo |
| SC-008/009 | Lote misto e configuração faltante; nenhuma regra padrão criada |
| SC-010 | EfeitoConcluido distingue falhas permanentes/transitórias; sem ACK pelo recuperador |

OpenObserve v1.1.2 real para os cenários de recuperação; falhas de transporte/status induzidas em servidor/proxy de teste. Nenhum teste depende de dados de produção.

## Arquivos de código/dados alterados contra Development

- `.env.example`
- `cmd/finder/main.go`
- `internal/adapter/out/openobserve/cliente.go`
- `internal/adapter/out/openobserve/cliente_regressao_test.go`
- `internal/adapter/out/openobserve/cliente_test.go`
- `internal/adapter/out/openobserve/consulta.go`
- `internal/adapter/out/openobserve/consulta_test.go`
- `internal/adapter/out/openobserve/doc.go`
- `internal/adapter/out/openobserve/integracao_test.go`
- `internal/application/port/out/busca.go`
- `internal/application/port/out/efeito.go`
- `internal/application/recuperar.go`
- `internal/application/recuperar_test.go`
- `internal/config/config.go`
- `internal/config/config_test.go`
- `internal/domain/recuperacao.go`
- `internal/domain/recuperacao_test.go`

## Testes existentes/novos associados aos arquivos alterados

- `internal/adapter/out/openobserve/cliente_regressao_test.go`: `TestRespostaAcimaDe32MiBNaoTemTetoInventado`
- `internal/adapter/out/openobserve/cliente_regressao_test.go`: `TestErroDeLeituraNaoEhRespostaInvalida`
- `internal/adapter/out/openobserve/cliente_regressao_test.go`: `TestStatusNaoPrevistoMantemPendencia`
- `internal/adapter/out/openobserve/cliente_regressao_test.go`: `TestCorpoHTTPTruncadoMantemPendencia`
- `internal/adapter/out/openobserve/cliente_regressao_test.go`: `TestJSONDeSucessoExigeFimValido`
- `internal/adapter/out/openobserve/cliente_test.go`: `TestClienteEnviaBasicENaoGravaTimeout`
- `internal/adapter/out/openobserve/cliente_test.go`: `TestClienteClassificaFalhaSemConteudo`
- `internal/adapter/out/openobserve/cliente_test.go`: `TestClienteTimeout`
- `internal/adapter/out/openobserve/cliente_test.go`: `TestClienteIndisponivel`
- `internal/adapter/out/openobserve/cliente_test.go`: `TestHitsVaziosSaoPaginaCurta`
- `internal/adapter/out/openobserve/consulta_test.go`: `TestMontarCorpoDaConsulta`
- `internal/adapter/out/openobserve/integracao_test.go`: `TestOpenObserveIntegracao`
- `internal/application/recuperar_test.go`: `TestPaginaCurtaEntregaLoteENaoConfirmaRedis`
- `internal/application/recuperar_test.go`: `TestLoteUnitarioFiltraRecordID`
- `internal/application/recuperar_test.go`: `TestConfiguracaoAusenteNaoConsulta`
- `internal/application/recuperar_test.go`: `TestListaVaziaTornaDiferencaDivergente`
- `internal/application/recuperar_test.go`: `TestNomesAusentesNaoConsultamEIsolam`
- `internal/application/recuperar_test.go`: `TestPaginacaoDedupEContagem`
- `internal/application/recuperar_test.go`: `TestCompletudeNaoAfirmaCausa`
- `internal/application/recuperar_test.go`: `TestIsolamentoDeDocumentoECredencial`
- `internal/application/recuperar_test.go`: `TestFalhaClassificadaSemParcial`
- `internal/config/config_test.go`: `TestCarregarPadrao`
- `internal/config/config_test.go`: `TestCarregarWorkers`
- `internal/config/config_test.go`: `TestAplicarArquivoEnv`
- `internal/config/config_test.go`: `TestAplicarArquivoEnvAusente`
- `internal/config/config_test.go`: `TestCarregarNaoCriaParametrosDeConsumo`
- `internal/config/config_test.go`: `TestCarregarIgnoraArquivoEnvComoPadrao`
- `internal/config/config_test.go`: `TestValidarRecusaCadaParametroNumerico`
- `internal/config/config_test.go`: `TestValidarRecusaNaoPositivoEStreamForaDoContrato`
- `internal/config/config_test.go`: `TestValidarAceitaConsumoInformado`
- `internal/config/config_test.go`: `TestCarregarOpenObserveNaoCriaPadrao`
- `internal/config/config_test.go`: `TestCarregarListaCopiaVaziaEstaPresente`
- `internal/config/config_test.go`: `TestCarregarOpenObserveInformadoNaoRegistraSenha`
- `internal/config/config_test.go`: `TestCarregarOpenObserveRecusaNaoPositivo`
- `internal/domain/recuperacao_test.go`: `TestDeduplicarCopiaEDivergencia`
- `internal/domain/recuperacao_test.go`: `TestDeduplicarIgnoraDocumentoSemIdentificadorENaoUsaCorrelacao`
- `internal/domain/recuperacao_test.go`: `TestFecharCompletudeEPrecedencia`
- `internal/domain/recuperacao_test.go`: `TestConsultaSQLMontaFiltroERecusaValorInaplicavel`
- `internal/domain/recuperacao_test.go`: `TestIsoladoNaoUsaAtributoInterno`

Comandos e saídas observadas estão em [evidencias](evidencias/); resultados compostos incluem manifesto de origem e SHA-256. Ausência de teste no pacote é indicada pelo próprio Go; não equivale a aceite externo.
