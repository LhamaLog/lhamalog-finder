---
trigger: always_on
description: Exige autorização antes de branch, commit, push e pull request
---

# Autorização de Git

Não crie branch, commit, push nem pull request sem autorização explícita nesta conversa.

## O que pedir antes

- Branch: confirme o nome e espere um sim.
- Commit: confirme que pode gravar o commit e peça o título.
- Push: confirme que pode enviar ao remoto.
- Pull request: confirme que pode abrir e peça o título.

Sem o título, não faça o commit nem abra a pull request. Não invente o título.

## Exemplo

Pedido: "sobe isso".

```text
Posso criar o commit e o push? Qual é o título do commit?
```

Só depois do sim e do título:

```text
título: Registra a integração cursor-agent do Spec Kit
```

Pedido: "abre a PR".

```text
Posso abrir a pull request? Qual é o título?
```
