package agrupar

import (
	"errors"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"math"
)

func Fundir(a domain.AchadoAnalise, estados []domain.EstadoOcorrencia) (domain.AchadoAnalise, error) {
	if len(estados) > 1 {
		return domain.AchadoAnalise{}, errors.New("multiplas_ocorrencias_para_mesma_chave")
	}
	if len(estados) == 0 {
		return a, nil
	}
	e := estados[0]
	if e.ContagemAcumulada < 0 || a.Contagem < 0 || a.Contagem > math.MaxInt64-e.ContagemAcumulada {
		return domain.AchadoAnalise{}, errors.New("estado_ocorrencia_invalido")
	}
	a.Contagem += e.ContagemAcumulada
	if e.PrimeiraGlobal.Before(a.PrimeiraOcorrencia) {
		a.PrimeiraOcorrencia = e.PrimeiraGlobal
	}
	if e.UltimaGlobal.After(a.UltimaOcorrencia) {
		a.UltimaOcorrencia = e.UltimaGlobal
	}
	a.StatusAnterior = e.Status // Nunca transforma FECHADA em ABERTA: responsabilidade do Módulo.
	return a, nil
}
