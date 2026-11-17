package domain

import "time"

// IdentidadeAchado é uma tupla; não inventa codificação de chave de transporte.
type IdentidadeAchado struct {
	Assinatura          string
	Classe              ClasseLog
	TenantID, ProductID string
}
type EstadoOcorrencia struct {
	ContagemAcumulada            int64
	PrimeiraGlobal, UltimaGlobal time.Time
	Status                       string
}
type ChaveResumo struct{ BatchID, ExecID string }

func (r ResumoAnalise) Chave() ChaveResumo { return ChaveResumo{BatchID: r.BatchID, ExecID: r.ExecID} }
