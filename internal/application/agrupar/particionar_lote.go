package agrupar

import (
	"errors"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

func ParticionarLote(logs []domain.LogClassificado) ([]domain.AchadoLote, error) {
	return ParticionarComAssinador(logs, domain.AssinarClassificado)
}

func ParticionarComAssinador(logs []domain.LogClassificado, assinar func(domain.LogClassificado) (string, error)) ([]domain.AchadoLote, error) {
	indices := map[domain.IdentidadeAchado]int{}
	var achados []domain.AchadoLote
	for _, l := range logs {
		if !l.Classe.Valida() || l.TenantID == "" || l.ProductID == "" || l.Apontamento == "" {
			return nil, errors.New("identidade_incompleta")
		}
		assinatura, err := assinar(l)
		if err != nil {
			return nil, err
		}
		chave := domain.IdentidadeAchado{Assinatura: assinatura, Classe: l.Classe, TenantID: l.TenantID, ProductID: l.ProductID}
		i, ok := indices[chave]
		if !ok {
			i = len(achados)
			indices[chave] = i
			achados = append(achados, domain.AchadoLote{Chave: chave, PrimeiraNoLote: l.Timestamp, UltimaNoLote: l.Timestamp})
		}
		a := &achados[i]
		a.Logs = append(a.Logs, l)
		a.ContagemLote++
		if l.Timestamp.Before(a.PrimeiraNoLote) {
			a.PrimeiraNoLote = l.Timestamp
		}
		if l.Timestamp.After(a.UltimaNoLote) {
			a.UltimaNoLote = l.Timestamp
		}
	}
	return achados, nil
}
