package openobserve

import (
	"encoding/json"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

type corpoConsulta struct {
	Query queryConsulta `json:"query"`
}

type queryConsulta struct {
	SQL       string `json:"sql"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	From      int    `json:"from"`
	Size      int    `json:"size"`
}

// Montar junta o corpo do POST. O prazo do cliente não entra neste JSON.
func Montar(ref domain.Referencia, cfg domain.ConfiguracaoAcesso, from, size int) ([]byte, error) {
	sql, inicio, fim, err := domain.ConsultaSQL(ref, cfg)
	if err != nil {
		return nil, err
	}
	if from < 0 || size <= 0 {
		return nil, domain.ErrFiltro
	}
	return CorpoJSON(sql, inicio, fim, from, size)
}

// CorpoJSON serializa a página. Não grava timeout nem credencial.
func CorpoJSON(sql string, inicio, fim int64, from, size int) ([]byte, error) {
	corpo := corpoConsulta{
		Query: queryConsulta{
			SQL:       sql,
			StartTime: inicio,
			EndTime:   fim,
			From:      from,
			Size:      size,
		},
	}
	return json.Marshal(corpo)
}
