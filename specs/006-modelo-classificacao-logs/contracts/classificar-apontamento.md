# Contract: Analisar lote (classificar + achados + ResumoAnalise)

Contrato entre o caso de uso de análise e as portas de entrada/saída. **Não** é endpoint HTTP do Finder. A fila Redis **não** é lida dentro desta porta; a entrada são logs **já capturados** para o apontamento/lote.

## Porta de entrada

`Analisar lote` (classificação + agrupamento no lote)

### Entrada

| Campo | Obrigatório | Descrição |
|-------|-------------|-----------|
| `referencia` | sim | `batchId`, tenant, produto, metadados ingested v1 (`schemaVersion` "1") — T5 D |
| `execId` | opcional | Parte da chave N4 A quando disponível; **omitir** no fluxo só `log.ingested` v1 (DN-006-1b B) |
| `logs` | sim | Coleção de logs capturados (apontamento ≠ conteúdo) |
| `catalogo` | opcional | Fornecido externamente; feature **não cria** (006-7 D) |
| `versaoRegras` | sim | Versão única 006-3 C |

### Saída (sucesso de lote completo)

| Campo | Descrição |
|-------|-----------|
| `resumo` | `ResumoAnalise` pronto para envio (ver `resumo-analise-handoff.md`) |
| `achados` | Lista espelhada em `resumo.achados` |

### Saída (falha / não sucesso)

| Condição | Resultado | Motivo distinguível |
|----------|-----------|---------------------|
| Nenhum log presente | Sem resumo de sucesso | `log_ausente` |
| OpenObserve inacessível | Sem resumo de sucesso | `armazenamento_inacessivel` |
| Lote incompleto | Classificação dos presentes **pode** ocorrer; **não** entrega sucesso de lote | `lote_incompleto` (T4 A) |
| Tenant ou produto vazio | Rejeição | `identidade_incompleta` |
| Log contável sem `service.name` ou **host** (assinatura DN-008-1 B) | Rejeição / não agrupa em achado | `assinatura_incompleta` (analyze **E2**, **G1**; spec.md FR-007/FR-008) |
| `productId` ausente no evento ou **inconsistente** entre metadata validada e corpo `log.ingested` v1 | Rejeição; **sem retry** de processamento do lote (ADR_049) | `produto_ausente_ou_inconsistente` (**008-16 A**, analyze2 **E2**; spec.md FR-008, SC-003) |
| `schemaVersion` ≠ `"1"` | Rejeição | `envelope_invalido` |

### Regras de agrupamento

- Um **AchadoAnalise** por **assinatura** no lote (classe já embutida no componente `nivel` — 008-15 A).
- **Assinatura (DN-008-1 B)**: nível + mensagem normalizada + `service.name` + host.
- **Normalização da mensagem**: **DN-008-4 B** — trim, case-insensitive; marcadores **DN-008-5 A**: `<NUM>`, `<UUID>`, `<IP>`, `<DATE>`, `<HEX>`; **MUST** chamar `NormalizarMensagem` em **`internal/application/assinatura/normalizar_mensagem.go`** (spec **008**, T016); a 006 **MUST NOT** manter implementação paralela em `internal/domain/` (fonte: respostas-joao.md §Etapa 3b/3c; spec.md FR-013).
- **Valor serializado da assinatura**: **008 assinatura A** — texto concatenado com separador **`|`** (**DN-008-8 A**); **`|`** literal em componente → **`\|`** (**DN-008-9 A**; canônico na spec **008**). Escape de **`\`**: pendência (sem fonte). (fonte: respostas-joao.md §Etapa 3e (João, 16/11 12h40), DN-008-8 A; §Etapa 3g (João, 16/11 13h22), DN-008-9 A)
- Logs de assinaturas distintas → achados distintos (User Story 2 cenário 2).

### Pré-handoff: consulta ao Módulo (DN-008-7 B)

Após montar a lista de achados e **antes** de invocar `Enviar ResumoAnalise`:

| Regra | Detalhe | Fonte |
|-------|---------|-------|
| Grão da consulta | **Uma** chamada ao Módulo **por achado** | respostas-joao.md §Etapa 3c (João, 16/11 11h02), DN-008-7 B |
| Grão do envio | **Um** `ResumoAnalise` **por lote** | respostas-etapa3c.md §6; ADR_062 |
| Contrato HTTP | **BLOQUEADO (DN-009-2)** para consulta e envio reais | respostas-etapa3.md §10 |

### Cenários de aceite (amostra)

| Entrada | Então |
|---------|-------|
| Três logs mesma assinatura, ERROR | Um `AchadoAnalise` tipo ERROR no resumo |
| Vários achados no lote | Consulta Módulo **por achado**; **um** resumo ao enviar (DN-008-7 B) |
| Dois logs INFO, assinaturas distintas | Dois achados INFO (006-11 A, T7 A) |
| Campo error + catálogo diria INFO | Classe ERROR; catálogo não altera |
| Envio concluído | **Sem** publicação de evento; **sem** gravação de ocorrência pelo Finder |

## Porta de saída: consultar Módulo por achado

| Operação | Obrigação |
|----------|-----------|
| `Consultar ocorrência existente` (ou equivalente) | **Uma chamada por achado** antes do envio do lote (**DN-008-7 B**, N2 B) |
| Estado implementação | **BLOQUEADO (DN-009-2)** — rota e contrato do Módulo |

## Porta de saída: entregar ResumoAnalise

| Operação | Obrigação |
|----------|-----------|
| `Enviar` | HTTP interno + JWT ao Módulo Ocorrência (N1 A); corpo = **um** `ResumoAnalise` **por lote** (ADR_062) |
| Estado implementação | **BLOQUEADO (DN-009-2)** até contrato HTTP do Módulo |
| Idempotência | Chave `batchId` ou `batchId`+`execId` conforme montada (N4 A, DN-006-1b B) |
| Confirmação | Aguardar confirmação do Módulo antes de marcar varredura concluída — **bloqueado** com adaptador |
| Reenvio | Mesma chave; Módulo ignora duplicata (FA-02) |
| Fora do escopo | MongoDB ocorrência, gravidade, evento Core, sugestão, notificação |

## Porta de saída: ocorrência (removida do Finder)

**Não há** porta `Inserir ocorrência` nesta feature após T1 REVISTO. Upsert por achado é responsabilidade do **Módulo Ocorrência** (006-9 A).

## Registro operacional

Log operacional: `batchId`, `execId` (se houver), contagem de achados por classe, motivo de falha, versão das regras. **Sem** texto integral de log.

## Fontes

- (fonte: respostas-joao.md, T1 REVISTO / T2 C / T4 A / N1 A / N4 A / 006-9 A / Etapa 3 DN-008-1 B / DN-006-1b B)
- (fonte: spec.md, FR-001–FR-012, User Stories 1–2)
- (fonte: revisao-T1-drive.md, §2)