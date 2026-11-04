<!--
Sync Impact Report
- Version change: template sem versão → 1.0.0
- Modified principles:
  - [PRINCIPLE_1_NAME] → I. Finalidade: logs em ocorrências
  - [PRINCIPLE_2_NAME] → II. Origem dos logs no coletor
  - [PRINCIPLE_3_NAME] → III. Categorias de ocorrência
  - [PRINCIPLE_4_NAME] → IV. Fila Redis e apontamento no OpenObserve
  - [PRINCIPLE_5_NAME] → V. Incerteza sem invenção (NÃO NEGOCIÁVEL)
- Added sections:
  - Restrições Técnicas
  - Evidência e Consulta
  - Governance (regras concretas)
- Removed sections: nenhuma (placeholders do template substituídos)
- Follow-up TODOs: nenhum
- Nota: persistência da ocorrência e nomes de fila, stream ou índice
  não foram informados e não foram definidos aqui.
-->

# LhamaLog-Finder Constitution

## Core Principles

### I. Finalidade: logs em ocorrências

O LhamaLog-Finder é um módulo do produto LhamaLog. Ele MUST coletar logs
e agregá-los, transformando o resultado em ocorrência.

Rationale: o módulo existe para esse fluxo. Trabalho que não coleta,
analisa ou agrega logs em ocorrência está fora do propósito deste
repositório.

### II. Origem dos logs no coletor

O Finder MUST analisar logs coletados pelo módulo LhamaLog-collector via
OpenTelemetry. A análise MUST incluir a agregação desses logs em
ocorrências.

Rationale: a coleta de telemetria pertence ao LhamaLog-collector. O
Finder parte dos logs que esse módulo já coletou.

### III. Categorias de ocorrência

Toda ocorrência MUST ser categorizada como erro, alerta ou informação.
Essas são as únicas classes reconhecidas por esta constituição.

Rationale: a classificação da ocorrência é parte do resultado do Finder.
Uma classe fora dessa lista não pode ser tratada como requisito válido
sem emenda desta constituição.

### IV. Fila Redis e apontamento no OpenObserve

O Finder MUST acessar uma fila Redis com códigos de apontamento, no
OpenObserve, dos logs coletados via OpenTelemetry. Após obter a posição,
o Finder MUST agregar e processar o log coletado correspondente. O código
de apontamento MUST NOT ser tratado como o conteúdo do log.

Rationale: a fila entrega a posição do log. A agregação e o processamento
ocorrem sobre o log apontado.

### V. Incerteza sem invenção (NÃO NEGOCIÁVEL)

Quando a informação estiver incompleta ou incerta, o trabalho MUST parar
e consultar quem solicitou a mudança, descrevendo o problema identificado,
antes de criar dados, regras ou resultados. Resultados MUST NOT ser
inventados. Código, teste ou comando MUST NOT ser descrito como executado
sem execução real e observação da saída.

Rationale: especificação, implementação e relato precisam apoiar-se no
que foi informado ou no que foi de fato executado.

## Restrições Técnicas

O fluxo normativo é: ler o código de apontamento na fila Redis, resolver
a posição no OpenObserve e, só então, agregar e processar o log coletado.

O README deste repositório documenta o estado atual do módulo: Go,
arquitetura hexagonal leve, consumo de Redis Streams, análise em paralelo
com pool de workers e mascaramento de dados sensíveis. Especificações
MUST NOT contradizer essa descrição sem emenda desta constituição ou
atualização explícita do README.

Dependências listadas no README como candidatas e ainda não adicionadas
MUST NOT ser tratadas como já presentes no módulo.

Persistência da ocorrência e os nomes concretos de fila, stream ou índice
não estão definidos nesta constituição.

## Evidência e Consulta

- Especificações, planos e tarefas MUST separar fato, premissa e pergunta
  em aberto. Premissa sem fonte no pedido, no repositório ou em execução
  observada MUST ser apresentada como consulta, não como dado criado.
- Relatos de verificação MUST indicar o que foi lido, o que foi executado
  e o que permanece em aberto.
- Uma afirmação de compilação, teste ou execução MUST corresponder a um
  comando executado e à saída observada.

## Governance

Esta constituição prevalece sobre práticas conflitantes em especificações,
planos, tarefas e implementação deste repositório.

Emendas:

- Mudança de finalidade, origem dos logs, categoria de ocorrência, fluxo
  da fila Redis até o OpenObserve, ou da regra de incerteza MUST ser
  registrada neste arquivo antes de orientar especificação ou código.
- A versão segue SemVer. MAJOR: remoção ou redefinição incompatível de
  princípio. MINOR: princípio ou seção nova, ou expansão material de
  regra. PATCH: clareza, redação ou correção sem mudança de regra.
- A data de ratificação é a da adoção inicial e não muda. A data da
  última emenda MUST ser a data da mudança.

Conformidade:

- Revisão de spec, plano, tarefas e código MUST checar a finalidade, as
  três categorias, o fluxo fila Redis → posição no OpenObserve →
  agregação do log, e a regra de incerteza.
- Complexidade além desse fluxo MUST ser justificada na especificação
  antes de ser tratada como requisito.
- O comentário HTML de Sync Impact Report no topo deste arquivo é rascunho
  de revisão da emenda e MUST ser removido antes do commit da constituição.

**Version**: 1.0.0 | **Ratified**: 2026-10-29 | **Last Amended**: 2026-10-29
