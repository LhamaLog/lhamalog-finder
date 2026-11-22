package domain

import "time"

type Achado struct {
	Assinatura                           string
	Classe                               ClasseLog
	Logs                                 []Codigo
	Contagem                             int64
	PrimeiraOcorrencia, UltimaOcorrencia time.Time
}
