package agrupar

import (
	"context"
	"errors"
	entrada "github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

type PedidoAnalise = entrada.PedidoAnalise
type PedidoRecorrencia = out.PedidoRecorrencia
type MontadorResumo interface {
	Montar(context.Context, entrada.PedidoAnalise) (domain.ResumoAnalise, error)
}

var ErrRecorrenciaIndisponivel = errors.New("recorrencia_indisponivel")

type Orquestrador struct {
	Montador    MontadorResumo
	Consulta    out.ConsultarAchadoModulo
	Recorrencia out.ConsultaRecorrencia
	Envio       out.EnviarResumo
}

// Preparar monta a parte local 006/008 e faz a única consulta ao Módulo por achado.
// A porta de entrega da 009 reutiliza o resumo já preparado, sem repetir consultas.
func (o Orquestrador) Preparar(ctx context.Context, p entrada.PedidoAnalise) (domain.ResumoAnalise, error) {
	if !p.Completo {
		return domain.ResumoAnalise{}, errors.New("lote_incompleto")
	}
	if o.Montador == nil || o.Consulta == nil {
		return domain.ResumoAnalise{}, errors.New("consulta_modulo_indisponivel")
	}
	r, err := o.Montador.Montar(ctx, p)
	if err != nil {
		return domain.ResumoAnalise{}, err
	}
	r.Achados = append([]domain.AchadoAnalise(nil), r.Achados...)
	for i, a := range r.Achados {
		chave := domain.IdentidadeAchado{Assinatura: a.Assinatura, Classe: a.Tipo, TenantID: r.TenantID, ProductID: r.ProductID}
		estados, err := o.Consulta.ConsultarAchado(ctx, chave)
		if err != nil {
			return domain.ResumoAnalise{}, errors.New("consulta_modulo_falhou")
		}
		a, err = Fundir(a, estados)
		if err != nil {
			return domain.ResumoAnalise{}, err
		}
		if o.Recorrencia == nil {
			return domain.ResumoAnalise{}, ErrRecorrenciaIndisponivel
		}
		marca, err := o.Recorrencia.ConsultarRecorrencia(ctx, out.PedidoRecorrencia{BatchID: r.BatchID, Host: a.HostName, Chave: chave, LimiteLotes: 500})
		if err != nil {
			return domain.ResumoAnalise{}, err
		}
		if marca.Valor == nil {
			return domain.ResumoAnalise{}, ErrRecorrenciaIndisponivel
		}
		a.Recorrente = *marca.Valor
		r.Achados[i] = a
	}
	return r, nil
}
func (o Orquestrador) Executar(ctx context.Context, p entrada.PedidoAnalise) (entrada.ResultadoAnalise, error) {
	r, err := o.Preparar(ctx, p)
	if err != nil {
		return entrada.ResultadoAnalise{}, err
	}
	saida := entrada.ResultadoAnalise{Resumo: r}
	if o.Envio == nil {
		return saida, errors.New("envio_modulo_falhou")
	}
	aceito, err := o.Envio.Enviar(ctx, r)
	if err != nil || !aceito {
		return saida, errors.New("envio_modulo_falhou")
	}
	saida.Concluida = true
	return saida, nil
}
