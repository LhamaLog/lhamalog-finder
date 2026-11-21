# Execução da spec 009

Referência: Development `6e8b0fd`. Branch: `LLT-201-009-publicar-evento-ocorrencia`.
Sem commit, push ou PR; nenhuma mudança nas specs 001–003.

## Plano anterior à implementação

Reutilizar tipos 006/normalizador 008 → validar/serializar payload provisório → handoff por porta → confirmação/reenvio → documentação do produto.

Entregas: código/testes somente nesta branch; contratos e pendências nesta pasta; evidências com comando, saída observada e exit code.
Testes/aceite: Taxonomia, gate de lote, uma entrega por lote, mesma chave, confirmação explícita; race.
Bloqueios: DN-009-2; ADR_051; N3 A; P-010-27; aceite integrado 008. info vs informacao provisórios por contrato.

Ordem compartilhada: auditoria 004 → recuperação 005 → fundação/normalizador 008 → segurança 010 → classificação/resumo 006 → base 007 → agrupamento 008 → handoff 009.
Os patches necessários de outras specs serão aplicados somente em uma cópia de validação composta, identificados nas evidências; não serão copiados como implementação para esta branch.
Dono do normalizador/compositor: 008; contexto seguro: 010; tipos de classificação e resumo comuns: 006. Consumidores usam portas/código canônico.
Tarefas bloqueadas continuam desmarcadas; conclusão parcial não equivale a aceite E2E.

## Checklists (estado inicial, somente leitura)

| Checklist | Total | Marcados | Abertos |
|---|---:|---:|---:|
| requirements.md | 16 | 15 | 1 |
| revisao-etapa3c.md | 45 | 37 | 8 |

Autorização explícita do usuário para continuar tarefas desbloqueadas 006–010 apesar dos itens abertos. Marcadores dos checklists não serão alterados.

## Pré-requisitos da skill

Skill: `.agents/skills/speckit-implement/SKILL.md`; regras Git: `.agents/rules/git-autorizacao.md`; constituição 1.0.0.
`.specify/extensions.yml` ausente: sem hooks before/after implement.
Comando solicitado: `pwsh -File .specify/scripts/powershell/check-prerequisites.ps1 -Json -RequireTasks -IncludeTasks`.
Saída observada: `pwsh: command not found` (exit 127). Conferência manual de spec, plan, tasks, research, data-model, contratos e quickstart; caminhos absolutos nesta worktree.
Ignore files existentes cobrem artefatos Go, IDE, segredos .env e contexto Docker; não exigem dependência nova.

## Resultados

Entrega lógica com confirmação, chave/snapshot estáveis, reenvio e DTO provisório. Reutiliza resumo 006 e preparação 008; não repete consultas nem publica evento de ocorrência. Cópia profunda inclui campos opcionais, mapas e linhas.

Tarefas: **40/51 verificadas**, 2 parciais e 9 bloqueadas. IDs, critérios, arquivos e testes em [rastreabilidade.md](rastreabilidade.md). Os checklists continuam com o mesmo estado inicial.

Na composição identificada pelo manifesto: `go test -count=1 ./...`, `go test -race -count=1 ./...` com Redis/OpenObserve reais e `go vet ./...`: exit 0 (aceite-final-*.txt). Isso não afirma compilação isolada desta branch sem os patches pré-requisitos 006/008/010. Branches consumidoras preservam apenas o código da própria spec.

## Cobertura do aceite e pendências

| Critério | Estado |
|---|---|
| SC-F001 | Um resumo após preparação, confirmação antes do callback; verificado com portas |
| SC-F002 | Snapshot/chave estáveis no reenvio, fake sem efeito duplicado; idempotência real Módulo não comprovada |
| SC-F003 | Rejeição de classes inválidas verificada |
| SC-001/003/004/006 | Pipeline Módulo→Core→Notificação apenas documentado; serviços externos não executados |
| SC-002/005 | Bloqueados até conclusão formal 008 e contratos do Módulo |

O teste de cadeia 010→006→008→009 usa serviços de aplicação reais e portas externas fake. Não é E2E de produto nem prova de notificação real.

## Revisão e reprodução

Ver [rastreabilidade.md](rastreabilidade.md), contratos, quickstart e evidencias. Desenvolvimento sem commit/push/PR. Relatório central em `/workspace/lhamalog-specs/README-execucao.md`.

## Contrato DN-009-2 aprovado em 20/11/2026

| Spec | Tarefas verificadas | Alterações desta etapa | Pendências |
|---|---|---|---|
| 004 | T033 | Reivindicação de pendências com `XAUTOCLAIM`, ociosidade preservada e avanço de cursor. Duas regressões verificadas com Redis real. | — |
| 006 | T031 e T055 | Token opaco no modelo e uso do cliente comum da 009, verificado na composição. | Catálogo oficial, avaliação com a base real da 007 e desempenho em VM-A. |
| 008 | T021 e T024 | Propagação da revisão e separação da montagem local do recálculo de agrupamento. | Consulta histórica real, integração de aviso ao Core e aceite formal. |
| 009 | T016, T018, T019, T034 e T050 | DTO N3, cliente HTTPS/JWT, confirmação e reenvio. A integração exercitou também o efeito T052 até o `XACK`. | Servidor do Módulo, configuração de implantação, finalização de Varredura e homologação downstream. |
| 010 | T042 parcial | Reutilização do cliente 009 verificada na composição. | Catálogo, casos oficiais e integração real do handoff mascarado. |

O fluxo testado conserva a chave no caso de conflito de revisão e de resposta perdida. O cliente conclui após um recibo válido. O servidor de teste simula persistência, sem substituir a homologação do Módulo implantado.

[Contrato aprovado](contracts/proposta-dn-009-2.md) · [Rastreabilidade e critérios](rastreabilidade.md) · [Comandos e saídas](quickstart.md#verificação-do-contrato-dn-009-2)
