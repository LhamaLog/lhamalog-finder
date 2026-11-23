package out

import (
	"context"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// EnviarResumo aguarda confirmação explícita; reenvio preserva ResumoAnalise.Chave.
// HTTP/JWT: DN-009-2. A interface não escolhe cabeçalho, rota ou código HTTP.
type EnviarResumo interface {
	Enviar(context.Context, domain.ResumoAnalise) (aceito bool, err error)
}
