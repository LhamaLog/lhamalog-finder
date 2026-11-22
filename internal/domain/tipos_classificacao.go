package domain

import "time"

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
