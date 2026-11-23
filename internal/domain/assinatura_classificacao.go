package domain

import canonica "github.com/LhamaLog/lhamalog-finder/internal/application/assinatura"

// AssinarClassificado delega ao pacote canônico 008 (T051); nenhuma cópia do algoritmo.
func AssinarClassificado(log LogClassificado) (string, error) {
	return canonica.Calcular(log.NivelCampo, log.Mensagem, log.ServiceName, log.Host)
}
