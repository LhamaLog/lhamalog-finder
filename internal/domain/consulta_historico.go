package domain

import (
	"fmt"
	"strings"
)

// ConsultasHistorico implementa DF-FINDER-2026-11-22 sem interpolar SQL recebido.
// A lista vazia seleciona a janela; uma lista presente carrega seus registros.
func ConsultasHistorico(ref Referencia, cfg ConfiguracaoAcesso, lotes []string) (string, int64, error) {
	_, _, fim, err := ConsultaSQL(ref, cfg)
	if err != nil || ref.MaxTimestampMicros < 0 {
		return "", 0, ErrFiltro
	}
	stream, _ := identificar(ref.Stream)
	tenant, _ := identificar(cfg.CampoTenant)
	produto, _ := identificar(cfg.CampoProduto)
	t, _ := literalObrigatorio(ref.TenantID)
	p, _ := literalObrigatorio(ref.ProductID)
	b, _ := literalObrigatorio(ref.BatchID)
	filtro := tenant + " = " + t + " AND " + produto + " = " + p
	if lotes == nil {
		return fmt.Sprintf(`SELECT "lhamalog_batch_id" AS batch_id, MAX("_timestamp") AS last_timestamp, COUNT(DISTINCT "lhamalog_record_id") AS record_count FROM %s WHERE %s AND "lhamalog_batch_id" <> %s GROUP BY "lhamalog_batch_id" ORDER BY last_timestamp DESC, batch_id ASC LIMIT 499`, stream, filtro, b), fim, nil
	}
	if len(lotes) == 0 || len(lotes) > 500 {
		return "", 0, ErrFiltro
	}
	valores := make([]string, len(lotes))
	for i, id := range lotes {
		v, e := literalObrigatorio(id)
		if e != nil {
			return "", 0, e
		}
		valores[i] = v
	}
	return fmt.Sprintf(`SELECT * FROM %s WHERE %s AND "lhamalog_batch_id" IN (%s) ORDER BY "lhamalog_batch_id" ASC, "lhamalog_record_id" ASC, "_timestamp" DESC`, stream, filtro, strings.Join(valores, ",")), fim, nil
}
