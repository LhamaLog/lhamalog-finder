package domain

// ID é a identidade do item da fila. Ela é distinta do código de apontamento.
type ID string

// Codigo é o apontamento do log no OpenObserve. Não é o conteúdo do log.
// Formato, tamanho e alfabeto não estão definidos.
type Codigo string

// Estado do item no consumo.
// Os estados gravados são só estes quatro.
// A condição parcial não é um estado: é o código já no conjunto com o item ainda pendente.
const (
	Livre        = "livre"
	EmPasso      = "em_passo"
	Concluido    = "concluido"
	NaoConcluido = "nao_concluido"
)

// Item é um registro da fila. ConteudoExtra marca carga além do código.
type Item struct {
	ID            ID
	Codigo        Codigo
	Estado        string
	ConteudoExtra bool
}

// Aceito é verdadeiro só quando o conteúdo é o código de apontamento.
// Item sem código, ou com conteúdo além do código, está fora do contrato.
func (i Item) Aceito() bool {
	return i.Codigo != "" && !i.ConteudoExtra
}
