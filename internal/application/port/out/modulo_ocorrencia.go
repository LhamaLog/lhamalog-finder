package out

// ModuloOcorrencia reutiliza as portas comuns 006/008; O adaptador HTTP/JWT segue DN-009-2 aprovado.
// Consulta é por achado, envio do resumo é por lote.
type ModuloOcorrencia interface {
	EnviarResumo
	ConsultarAchadoModulo
}
