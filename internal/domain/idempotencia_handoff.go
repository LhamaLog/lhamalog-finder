package domain

import (
	"errors"
	"regexp"
	"strings"
)

type ChaveIdempotenciaHandoff = ChaveResumo

var execVarredura = regexp.MustCompile(`^[VM][0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// A chave permanece tupla, sem inventar header, separador ou formato de transporte.
func ValidarResumoHandoff(r ResumoAnalise) error {
	invalido := errors.New("resumo_analise_invalido")
	for _, s := range []string{r.BatchID, r.TenantID, r.ProductID, r.VersaoRegras} {
		if strings.TrimSpace(s) == "" {
			return invalido
		}
	}
	if r.ExecID != "" && !execVarredura.MatchString(r.ExecID) {
		return invalido
	}
	c := r.Contadores
	if c.LinhasProcessadas < 0 || c.DadosSensiveisMascarados < 0 {
		return invalido
	}
	restantes := c.LinhasProcessadas
	for _, n := range []int64{c.Erros, c.Warnings, c.Info} {
		if n < 0 || n > restantes {
			return invalido
		}
		restantes -= n
	}
	if restantes != 0 {
		return invalido
	}
	for _, a := range r.Achados {
		if !a.Tipo.Valida() {
			return invalido
		}
	}
	return nil
}
