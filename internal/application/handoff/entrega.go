package handoff

import (
	"context"
	"errors"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

type Pendente struct {
	resumo domain.ResumoAnalise
	chave  domain.ChaveResumo
}

func (p *Pendente) Chave() domain.ChaveResumo {
	if p == nil {
		return domain.ChaveResumo{}
	}
	return p.chave
}
func copiarResumo(r domain.ResumoAnalise) domain.ResumoAnalise {
	r.Achados = append([]domain.AchadoAnalise(nil), r.Achados...)
	for i := range r.Achados {
		a := &r.Achados[i]
		if a.RevisaoConsultada != nil {
			revisao := *a.RevisaoConsultada
			a.RevisaoConsultada = &revisao
		}
		if a.DadosSensiveis != nil {
			valor := *a.DadosSensiveis
			a.DadosSensiveis = &valor
		}
		a.Linhas = append([]string(nil), a.Linhas...)
		if a.CamposExtraido != nil {
			c := make(map[string]string, len(a.CamposExtraido))
			for k, v := range a.CamposExtraido {
				c[k] = v
			}
			a.CamposExtraido = c
		}
	}
	return r
}
func NovaPendente(r domain.ResumoAnalise) (*Pendente, error) {
	if err := domain.ValidarResumoHandoff(r); err != nil {
		return nil, err
	}
	return &Pendente{resumo: copiarResumo(r), chave: r.Chave()}, nil
}

type ResultadoEntrega struct {
	Pendente              *Pendente
	Confirmada, Concluida bool
}
type Entrega struct {
	Envio     out.EnviarResumo
	Finalizar out.FinalizarAnalise
}

func (s Entrega) Enviar(ctx context.Context, p *Pendente) (ResultadoEntrega, error) {
	r := ResultadoEntrega{Pendente: p}
	if err := ctx.Err(); err != nil {
		return r, err
	}
	if p == nil || p.chave.BatchID == "" {
		return r, errors.New("resumo_analise_invalido")
	}
	if s.Envio == nil {
		return r, errors.New("envio_modulo_indisponivel")
	}
	aceito, err := s.Envio.Enviar(ctx, copiarResumo(p.resumo))
	if err != nil {
		for _, tipo := range []error{out.ErrRevisaoDesatualizada, out.ErrModuloTemporario, out.ErrModuloRequisicao, out.ErrModuloProtocolo, context.Canceled, context.DeadlineExceeded} {
			if errors.Is(err, tipo) {
				return r, tipo
			}
		}
	}
	if err != nil || !aceito {
		return r, errors.New("envio_modulo_falhou")
	}
	r.Confirmada = true
	if s.Finalizar == nil {
		return r, errors.New("conclusao_varredura_indisponivel")
	}
	if err := ctx.Err(); err != nil {
		return r, err
	}
	if err := s.Finalizar.ConcluirAnalise(ctx, p.chave); err != nil {
		return r, errors.New("conclusao_varredura_falhou")
	}
	r.Concluida = true
	return r, nil
}

// Reenviar preserva o snapshot e a mesma tupla. Em conflito de revisão,
// OrquestradorEntrega consulta o estado novamente antes de um novo envio.
func (s Entrega) Reenviar(ctx context.Context, p *Pendente) (ResultadoEntrega, error) {
	return s.Enviar(ctx, p)
}
