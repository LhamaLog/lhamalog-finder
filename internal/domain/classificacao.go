package domain

import (
	"strings"
	"unicode"
)

// ClassificarLog é pura: nenhuma porta/rede/LLM. Nível e palavras inteiras
// na mensagem pós-segurança decidem juntos; catálogo somente sem sinal reconhecido.
func ClassificarLog(nivel, mensagem string, catalogo []RegraClassificacao) ResultadoClassificacao {
	classe := ClasseLog("")
	termos := strings.FieldsFunc(strings.ToLower(nivel+" "+mensagem), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r) })
	for _, termo := range termos {
		var c ClasseLog
		switch termo {
		case "error", "erro", "fatal":
			c = ERROR
		case "warn", "warning", "alerta":
			c = WARNING
		case "info", "informação", "informacao", "debug", "trace":
			c = INFO
		}
		classe = MaisGrave(classe, c)
	}
	if classe.Valida() {
		return ResultadoClassificacao{Classe: classe}
	}
	for _, r := range catalogo {
		if r.Classe.Valida() && r.Corresponde != nil && r.Corresponde(mensagem) {
			return ResultadoClassificacao{Classe: r.Classe, CatalogoUsado: true}
		}
	}
	return ResultadoClassificacao{Classe: INFO}
}
