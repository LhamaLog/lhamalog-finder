package domain

// IdentificacaoRegras cobre regras+catálogo+base com a versão única fornecida.
// Cada resultado captura sua versão; mudar configuração não altera resultados antigos.
type IdentificacaoRegras struct {
	VersaoUnica    string `json:"versaoUnica"`
	SomenteNivel   bool   `json:"somenteNivel"`
	CatalogoEntrou bool   `json:"catalogoEntrou"`
}
