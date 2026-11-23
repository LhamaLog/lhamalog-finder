package assinatura

import (
	"errors"
	"strings"
)

var ErrIncompleta = errors.New("assinatura_incompleta")
var ErrEscapePendente = errors.New("escape_barra_invertida_pendente")

// Calcular preserva a ordem nível|mensagem|serviço|host (DN-008-8 A).
// DF-FINDER-2026-11-22: escapar barra invertida antes do delimitador.
func Calcular(nivel, mensagem, servico, host string) (string, error) {
	// O nivel e opcional na entrada (spec 006); sua ausencia preserva o componente vazio.
	if strings.TrimSpace(servico) == "" || strings.TrimSpace(host) == "" {
		return "", ErrIncompleta
	}
	partes := []string{nivel, NormalizarMensagem(mensagem), servico, host}
	for i, v := range partes {
		v = strings.ReplaceAll(v, `\`, `\\`)
		partes[i] = strings.ReplaceAll(v, "|", `\|`)
	}
	return strings.Join(partes, "|"), nil
}
