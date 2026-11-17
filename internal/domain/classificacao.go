package domain

import (
	"strings"
	"time"
	"unicode"
)

// RegraClassificacao é um predicado opaco fornecido pelo chamador.
// Não define sintaxe de padrão nem lê/cria catálogo (006-7 D permanece aberta).
type RegraClassificacao struct {
	Classe      ClasseLog
	Corresponde func(string) bool
}
type ResultadoClassificacao struct {
	Classe        ClasseLog
	CatalogoUsado bool
}
type LogClassificado struct {
	Apontamento                                           Codigo
	Classe                                                ClasseLog
	TenantID, ProductID                                   string
	CatalogoUsado                                         bool
	ServiceName, Host, NivelCampo, Mensagem, VersaoRegras string
	TraceID, ModuloID                                     string
	Timestamp                                             time.Time
}

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
