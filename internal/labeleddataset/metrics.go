package labeleddataset

import (
	"errors"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

type MetricasClasse struct {
	TP, FP, FN, Suporte  int
	Precisao, Recall, F1 float64
}
type Relatorio struct {
	Unidades int
	Classes  map[ClassLabel]MetricasClasse
}

// Avaliar conta um par gold/predição por log OU por apontamento, só em validacao.
func Avaliar(registros []Registro, documentacao ...DocumentacaoConjunto) (Relatorio, error) {
	if err := ValidarConjunto(registros, documentacao...); err != nil {
		return Relatorio{}, err
	}
	r := Relatorio{Classes: map[ClassLabel]MetricasClasse{domain.ERROR: {}, domain.WARNING: {}, domain.INFO: {}}}
	for _, item := range registros {
		if item.Split != "validacao" {
			continue
		}
		pred, err := Predizer(item)
		if err != nil {
			return Relatorio{}, err
		}
		r.Unidades++
		for classe, m := range r.Classes {
			if item.ClassLabel == classe {
				m.Suporte++
				if pred == classe {
					m.TP++
				} else {
					m.FN++
				}
			} else if pred == classe {
				m.FP++
			}
			r.Classes[classe] = m
		}
	}
	if r.Unidades == 0 {
		return r, errors.New("avaliacao_vazia")
	}
	for c, m := range r.Classes {
		if m.TP+m.FP > 0 {
			m.Precisao = float64(m.TP) / float64(m.TP+m.FP)
		}
		if m.TP+m.FN > 0 {
			m.Recall = float64(m.TP) / float64(m.TP+m.FN)
		}
		if m.Precisao+m.Recall > 0 {
			m.F1 = 2 * m.Precisao * m.Recall / (m.Precisao + m.Recall)
		}
		r.Classes[c] = m
	}
	return r, nil
}
func (r Relatorio) ValidarGate() error {
	if r.Unidades == 0 {
		return errors.New("avaliacao_vazia")
	}
	for _, c := range []ClassLabel{domain.ERROR, domain.WARNING, domain.INFO} {
		m, ok := r.Classes[c]
		if !ok || m.Suporte == 0 || m.Precisao < .95 || m.Recall < .95 || m.F1 < .95 {
			return errors.New("gate_95_por_classe_nao_atendido")
		}
	}
	return nil
}
