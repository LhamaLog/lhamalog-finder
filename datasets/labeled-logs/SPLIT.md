# Split qualitativo

Uso no MVP: somente avaliação das regras; nenhum treino ou LLM.

- `validacao`: registros com gold revisado usados para medir precisão, recall e F1 por classe; itens apontamento contam uma vez, pela previsão mais grave dos membros.
- `teste`: registros reservados para verificação adicional sem ajustar regras a partir de seus resultados.
- `treinamento`: reserva para pós-MVP; nenhum uso de treino autorizado no MVP.
- `so_catalogo`: apenas itens sem nível reconhecível em campo/mensagem cuja classe dependa exclusivamente do catálogo da 006. Conteúdo bloqueado por 006-7 D; nenhum registro inventado.

Sem proporções, seed ou distribuição numérica escolhidas. A atribuição dos registros reais deverá ser documentada na entrega. Os conjuntos de avaliação devem ser identificados sem apresentar fixtures como evidência de qualidade.
