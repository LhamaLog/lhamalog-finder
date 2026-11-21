package out

import "errors"

// Erros da porta têm mensagens controladas, apropriadas para o log do consumidor.
var (
	ErrRevisaoDesatualizada = errors.New("revisao_ocorrencia_desatualizada")
	ErrModuloTemporario     = errors.New("modulo_ocorrencia_temporariamente_indisponivel")
	ErrModuloRequisicao     = errors.New("modulo_ocorrencia_recusou_requisicao")
	ErrModuloProtocolo      = errors.New("resposta_modulo_ocorrencia_invalida")
)
