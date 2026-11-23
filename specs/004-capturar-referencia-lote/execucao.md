# Execução da spec 004

Referência: Development `6e8b0fd`. Branch: `LLT-269-004-capturar-referencia-lote`.
Sem commit, push ou PR; nenhuma mudança nas specs 001–003.

## Plano anterior à implementação

Auditar T001–T031 já marcadas; corrigir somente defeito demonstrado.

Entregas: código/testes somente nesta branch; contratos e pendências nesta pasta; evidências com comando, saída observada e exit code.
Testes/aceite: Unitários, race e prova Redis; SC-001–SC-009.
Bloqueios: Nenhum contrato novo; integração depende do ambiente.

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

31 tarefas pré-existentes auditadas antes de qualquer ajuste. Regressão posterior comprovou aceitação de sufixo JSON inválido; correção mínima exige EOF após o objeto. T032 adicionada com testes de domínio e Redis escritos antes da correção. Nenhuma reimplementação.

Tarefas: **32/32 verificadas**, 0 parciais e 0 bloqueadas. IDs, critérios, arquivos e testes em [rastreabilidade.md](rastreabilidade.md). Os checklists continuam com o mesmo estado inicial.

A auditoria inicial sem ajuste foi revisada após a regressão JSON: a versão final contém apenas o reparo de fim de JSON e seus testes. [Auditoria das 31 tarefas](auditoria-31.md). `go test -race -count=1 ./...` com Redis real: exit 0, json-envelope-green-suite.txt.

## Cobertura do aceite e pendências

| Critério | Evidência e alcance |
|---|---|
| SC-001/002 | Referências unitárias/múltiplas, campos corretos, sem conteúdo de log/ACK na interpretação |
| SC-003 | Redis real: inválidas confirmadas, sem efeito/estado concluído; regressão JSON adicional |
| SC-004 | Reserva/reentrega/TTL de 25h e queda; sem repetição de efeito |
| SC-005 | Grupo com posição 0 e dois grupos distintos |
| SC-006/007 | Redis fechado sem ACK; reivindicação de pendências |
| SC-008/009 | Logs operacionais sem bruto/credencial; parâmetros ausentes impedem leitura |

Cobertura nos testes mapeados em auditoria-31.md. Prova sintética local, sem homologação de produção.

## Revisão e reprodução

Ver [rastreabilidade.md](rastreabilidade.md), contratos, quickstart e evidencias. Desenvolvimento sem commit/push/PR. Relatório central em `/workspace/lhamalog-specs/README-execucao.md`.
