package labeleddataset

import (
	"errors"
	"regexp"
	"strings"
)

var ErrRegistro = errors.New("registro_rotulado_invalido")

// Heurísticas de publicação da base 007, não regex oficiais de mascaramento 010.
// Não certificam ausência completa de PII; a declaração/revisão de origem é obrigatória.
var emailHeuristico = regexp.MustCompile(`[^\s@]+@[^\s@]+\.[^\s@]+`)
var segredoHeuristico = regexp.MustCompile(`(?i)(?:senha|password|token|secret)\s*[:=]\s*\S+`)

func camposValidos(c CamposEntrada) bool {
	for _, s := range []string{c.LevelText, c.MessageText, c.BatchID, c.RecordID, c.TenantID, c.ProductID} {
		if strings.TrimSpace(s) == "" {
			return false
		}
	}
	return !emailHeuristico.MatchString(c.MessageText) && !segredoHeuristico.MatchString(c.MessageText)
}
func ValidateRegistro(r Registro) error {
	if !r.ClassLabel.Valida() {
		return ErrRegistro
	}
	switch r.ItemType {
	case "log":
		if !camposValidos(r.CamposEntrada) || len(r.Members) != 0 {
			return ErrRegistro
		}
	case "apontamento":
		if len(r.Members) == 0 || strings.TrimSpace(r.Documentacao) == "" || strings.TrimSpace(r.BatchID) == "" || strings.TrimSpace(r.TenantID) == "" || strings.TrimSpace(r.ProductID) == "" {
			return ErrRegistro
		}
		tenant, produto := r.TenantID, r.ProductID
		for _, m := range r.Members {
			if m.ItemType != "log" || !camposValidos(m) || m.TenantID != tenant || m.ProductID != produto {
				return ErrRegistro
			}
		}
	default:
		return ErrRegistro
	}
	return nil
}
