package assinatura

import (
	"errors"
	"strings"
)

var ErrIncompleta = errors.New("assinatura_incompleta")
var ErrEscapePendente = errors.New("escape_barra_invertida_pendente")

// Calcular preserva a ordem nível|mensagem|serviço|host (DN-008-8 A).
// Sem regra para barra invertida literal, esse caminho permanece bloqueado;
// não tenta parse simétrico nem escolhe um escape sem fonte.
func Calcular(nivel, mensagem, servico, host string) (string, error) {
	// O nivel e opcional na entrada (spec 006); sua ausencia preserva o componente vazio.
	if strings.TrimSpace(servico) == "" || strings.TrimSpace(host) == "" {
		return "", ErrIncompleta
	}
	partes := []string{nivel, NormalizarMensagem(mensagem), servico, host}
	for i, v := range partes {
		if strings.Contains(v, `\`) {
			return "", ErrEscapePendente
		}
		partes[i] = strings.ReplaceAll(v, "|", `\|`)
	}
	return strings.Join(partes, "|"), nil
}
