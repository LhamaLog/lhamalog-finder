# Execução da spec 005

Referência: Development `6e8b0fd`. Branch: `LLT-270-implementar-recuperacao-dos-logs-no-open-observe-pela-referencia-do-lote`.
Estado publicado: PR #8, aberta para `Development`, com HEAD inicial desta revisão em `735e556d7b8f6981362cf93ffcc60b6f7856efdc`. Esse commit de merge registra a implementação e a sincronização anterior com `Development` (`6e8b0fd`). O plano abaixo descreve a execução anterior à publicação.

## Plano anterior à implementação

Preservar LLT-270 e incorporar Development sem commit; T031 remove teto e T032 corrige falhas de transporte/status, testes primeiro.

Entregas: código/testes somente nesta branch; contratos e pendências nesta pasta; evidências com comando, saída observada e exit code.
Testes/aceite: Unitários, race, regressões HTTP e prova OpenObserve; SC-001–SC-010.
Bloqueios: Merge permanece sem commit; nenhuma homologação de produção.

Ordem compartilhada: auditoria 004 → recuperação 005 → fundação/normalizador 008 → segurança 010 → classificação/resumo 006 → base 007 → agrupamento 008 → handoff 009.
Os patches necessários de outras specs serão aplicados somente em uma cópia de validação composta, identificados nas evidências; não serão copiados como implementação para esta branch.
Dono do normalizador/compositor: 008; contexto seguro: 010; tipos de classificação e resumo comuns: 006. Consumidores usam portas/código canônico.
Tarefas bloqueadas continuam desmarcadas; conclusão parcial não equivale a aceite E2E.

## Checklists (estado inicial, somente leitura)

| Checklist | Total | Marcados | Abertos |
|---|---:|---:|---:|
| requirements.md | 16 | 16 | 0 |

Autorização explícita do usuário para continuar tarefas desbloqueadas 006–010 apesar dos itens abertos. Marcadores dos checklists não serão alterados.

## Pré-requisitos da skill

Skill: `.agents/skills/speckit-implement/SKILL.md`; regras Git: `.agents/rules/git-autorizacao.md`; constituição 1.0.0.
`.specify/extensions.yml` ausente: sem hooks before/after implement.
Comando solicitado: `pwsh -File .specify/scripts/powershell/check-prerequisites.ps1 -Json -RequireTasks -IncludeTasks`.
Saída observada: `pwsh: command not found` (exit 127). Conferência manual de spec, plan, tasks, research, data-model, contratos e quickstart; caminhos absolutos nesta worktree.
Ignore files existentes cobrem artefatos Go, IDE, segredos .env e contexto Docker; não exigem dependência nova.

## Resultados

Branch existente preservada, comparada com Development (2 commits exclusivos Development / 1 da LLT-270). Merge sem commit e sem conflitos incorporou a referência. T031 remove o teto de 32 MiB; T032 separa falha de leitura/status de JSON de sucesso inválido. Corrigidos sufixo JSON e problemas observados nos testes de integração (assert de filtro e acesso concorrente ao coletor).

Tarefas: **32/32 verificadas**, 0 parciais e 0 bloqueadas. IDs, critérios, arquivos e testes em [rastreabilidade.md](rastreabilidade.md). Os checklists continuam com o mesmo estado inicial.

`go test -race -count=1 ./...` na própria branch com Redis/OpenObserve reais: exit 0, validacao-final-race-integracoes.txt. unit.txt e race.txt também passaram na execução registrada nesses arquivos. A sincronização foi posteriormente publicada no commit `735e556`, com pais `5f82490` e `6e8b0fd`. Não há operação de merge em andamento na revisão de 21/11/2026.

## Cobertura do aceite e pendências

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

## Revisão e reprodução

Ver [rastreabilidade.md](rastreabilidade.md), contratos, quickstart e evidencias. A PR publicada é [#8](https://github.com/LhamaLog/lhamalog-finder/pull/8). As evidências antigas descrevem suas próprias execuções e não atestam automaticamente commits posteriores.

## Revisão de 21/11/2026

Requisitos conferidos na LLT-270 do Jira e no [Contrato do Consumidor v1 no Drive](https://docs.google.com/document/d/14yLd8WA7pRm6es4okYiOnUjHfZp0T_4uyUnVkFk2tQk/edit), seção 7. A [ADR_168](https://docs.google.com/document/d/1cq0KUW2gWsIq6Mv54SZhXcWIoHlKCqwqUG5XHD3rVFs/edit) mantém o isolamento por tenant no OpenObserve. A revisão corrigiu os registros desatualizados de publicação e a lista de arquivos do diff.

| Verificação executada | Resultado |
|---|---|
| `go test -count=1 ./...`, Go 1.27.0 Windows/amd64 | Exit 0; integrações opt-in desativadas nessa execução |
| `go vet ./...` e `go build ./...`, Go 1.27.0 Windows/amd64 | Exit 0 |
| `go test -race -count=1 -v ./...`, imagem `golang:1.27.0` Linux/amd64 | Exit 0, sem testes ignorados; Redis 8.2.10 e OpenObserve 1.1.2 reais em contêineres exclusivos de teste |
| `gofmt` nos arquivos Go alterados pela PR | Sem diferença de conteúdo no Git |
| `git diff --check` | Exit 0 |

Saída da execução com integrações em [evidencias/revisao-2026-11-21-race.txt](evidencias/revisao-2026-11-21-race.txt). A primeira tentativa de race nativo falhou pela ausência de suporte 64 bits no GCC instalado. A primeira tentativa no contêiner deixou de fornecer parâmetros obrigatórios do Redis; a execução final forneceu consumidor, prefixo, count, block e ociosidade e passou integralmente.

A prova usa a própria branch da PR #8. O HEAD `3b26753` da PR #7 não é ancestral desta branch e suas correções não foram importadas. A validação exercita a recuperação desta PR sobre a base existente e não atesta a composição futura das duas entregas.

Na consulta anterior ao push, a PR estava aberta, integrável, sem reviews ou comentários no GitHub, sem check runs e sem regras de proteção ou rulesets. A revisão pelo Discord será solicitada para o novo HEAD após a confirmação do push.

O parecer do Discord de 21/11/2026 às 22h39, em nove partes para `b7eb540`, concluiu que a PR ainda não está apta para merge. Os ajustes locais posteriores corrigem T029/T030 nesta rastreabilidade, documentam FINDER_OO_* e o destino atual dos registros no README, e mantêm .specify/feature.json no mesmo valor de Development. A descrição da PR passa a indicar apenas a validação reproduzida nesta revisão.

Permanece o bloqueio de documentação externa em `lhamalog-docs/planejamento/LLT-270-recuperar-logs-lote/`. O revisor pediu o registro do contrato de busca, das regras de efeito/confirmação e das decisões de paginação e memória nesse repositório. A autorização desta execução abrange somente as branches do Finder das PRs #8–#13. A publicação externa depende de ampliação expressa desse escopo. A revisão também pede decisão registrada sobre o fluxo que confirma a entrada após recuperar os registros, ainda sem encaminhá-los para a próxima etapa; esse comportamento já está descrito no contrato da spec 005 e no README.
