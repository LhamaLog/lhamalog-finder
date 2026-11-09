# Data Model: Modelo de classificação de logs

Modelo do caso de uso. URI, banco, coleção e o catálogo não estão aqui como valores. Vêm da execução.

## Apontamento de entrada

Conjunto já capturado. Esta feature não o busca na fila nem no OpenObserve.

| Campo | Regra |
|-------|--------|
| Identidade | Identifica o apontamento. Não é o texto dos logs. Duas gravações da mesma identidade inserem duas ocorrências. |
| Tenant | Autoridade do tenant. Vazio: não grava. Não é lido na mensagem. |
| Produto | Autoridade do produto. Mesma regra. |
| Armazenamento inacessível | Quando verdadeiro, não grava. O motivo é armazenamento inacessível, mesmo que a lista de logs venha preenchida. |
| Logs | Lista do apontamento. Pode misturar presente e ausente. |

## Log do apontamento

| Campo | Regra |
|-------|--------|
| Presente | Falso: não entra na classe e o texto, se vier, é ignorado. |
| Campo de nível | Texto já separado. Espaços nas bordas não mudam o reconhecimento. O valor inteiro é comparado à lista. |
| Mensagem | Texto já separado. Só a palavra inteira da lista conta. Não é gravada. |

Lista do campo e da palavra, sem distinção de maiúsculas: erro, error e fatal são erro; alerta, warn e warning são alerta; informação, info, debug e trace são informação.

Palavra inteira: o trecho igual ao termo, e o caractere ao lado não é letra nem dígito. Borda do texto vale.

## Catálogo da execução

Opcional. Esta feature não cria entrada.

| Campo | Regra |
|-------|--------|
| Identificação | Se vier, é guardada na ocorrência. Se não vier, a ocorrência informa que só as regras de nível desta feature foram usadas. |
| Entrada | Classe mais texto, na ordem fornecida. Classe fora das três é ignorada. A primeira palavra inteira na mensagem decide. Sem acerto, a classe do log é informação. |

O catálogo só é consultado quando o campo não está na lista e a mensagem não tem palavra inteira da lista.

## Classe

Três valores: informação, alerta e erro. Erro vence alerta. Alerta vence informação.

| Origem | Classe |
|--------|--------|
| Campo e palavras da mensagem | A mais grave entre eles. O catálogo não entra. |
| Só o catálogo | A primeira entrada válida que casa. |
| Nenhum sinal | Informação. |

A classe da ocorrência é a mais grave entre os logs presentes. Log ausente não entra.

## Ocorrência gravada

Um documento inserido. Não substitui outro da mesma identidade.

| Campo | Regra |
|-------|--------|
| Classe | A da ocorrência. |
| Identificação das regras | A identificação do catálogo, ou a informação de que só as regras de nível desta feature foram usadas. |
| Catálogo entrou | Verdadeiro só se a classe de algum log presente saiu do catálogo. |
| Identidade do apontamento | A da entrada. |
| Tenant | O da entrada. |
| Produto | O da entrada. |

Não há campo de mensagem, de campo de nível, de texto de log, de gravidade, de sugestão, de evento nem de notificação.

## Resultado da chamada

Uma ocorrência pronta ou um motivo sem ocorrência.

| Situação | Ocorrência pronta | Motivo |
|----------|-------------------|--------|
| Ao menos um log presente, identidade completa, gravação concluída | Sim | — |
| Nenhum log presente | Não | Log ausente |
| Armazenamento inacessível | Não | Armazenamento inacessível |
| Gravação não concluída | Não | Falha de gravação |
| URI, banco ou coleção ausente | Não | Configuração ausente |
| Tenant ou produto vazio | Não | Identidade incompleta |

Log ausente no meio de um conjunto com outro log presente não é a linha "nenhum log presente".
