package out

// ModuloOcorrencia reutiliza as portas comuns 006/008; HTTP/JWT bloqueado DN-009-2.
// Consulta é por achado, envio do resumo é por lote.
type ModuloOcorrencia interface {
	EnviarResumo
	ConsultarAchadoModulo
}
