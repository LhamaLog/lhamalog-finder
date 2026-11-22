// Package assinatura é o algoritmo canônico da spec 008, compartilhado por 006/009/010.
// Recebe exclusivamente mensagem pós-segurança. Casos formais adicionais: P-010-27.
package assinatura

import (
	"net/netip"
	"regexp"
	"strings"
)

const (
	MarcadorNumero = "<NUM>"
	MarcadorUUID   = "<UUID>"
	MarcadorIP     = "<IP>"
	MarcadorData   = "<DATE>"
	MarcadorHex    = "<HEX>"
)

// Uma substituição única preserva categorias compostas e os próprios marcadores:
// UUID/data/IP/hex são reconhecidos antes de números internos a esses valores.
var categorias = regexp.MustCompile(`(?i)<(?:num|uuid|ip|date|hex)>|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}|[0-9]{4}-[0-9]{2}-[0-9]{2}(?:t[0-9]{2}:[0-9]{2}:[0-9]{2}(?:\.[0-9]+)?(?:z|[+-][0-9]{2}:[0-9]{2})?)?|(?:[0-9]{1,3}\.){3}[0-9]{1,3}|(?:[0-9a-f]{0,4}:){2,}[0-9a-f:.]*|0x[0-9a-f]+|[-+]?[0-9]+(?:\.[0-9]+)?(?:e[-+]?[0-9]+)?`)

func NormalizarMensagem(mensagem string) string {
	return categorias.ReplaceAllStringFunc(strings.ToLower(strings.TrimSpace(mensagem)), func(v string) string {
		switch {
		case strings.HasPrefix(v, "<"):
			return strings.ToUpper(v)
		case len(v) == 36 && v[8] == '-' && v[13] == '-':
			return MarcadorUUID
		case len(v) >= 10 && v[4] == '-' && v[7] == '-':
			return MarcadorData
		case strings.HasPrefix(v, "0x"):
			return MarcadorHex
		default:
			if _, err := netip.ParseAddr(v); err == nil {
				return MarcadorIP
			}
			if strings.Count(v, ".") > 1 || strings.Contains(v, ":") {
				return v
			}
			return MarcadorNumero
		}
	})
}
