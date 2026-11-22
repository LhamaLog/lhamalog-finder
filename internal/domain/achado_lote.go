package domain

import "time"

type ChaveUnicidadeOcorrencia = IdentidadeAchado
type LogAgrupavel = LogClassificado
type AchadoLote struct {
	Chave                        IdentidadeAchado
	Logs                         []LogClassificado
	ContagemLote                 int64
	PrimeiraNoLote, UltimaNoLote time.Time
}
type EstadoAgregacaoEntreLotes = EstadoOcorrencia
type AchadoAnaliseAgrupado = AchadoAnalise
type ResumoAnaliseEntrega = ResumoAnalise
