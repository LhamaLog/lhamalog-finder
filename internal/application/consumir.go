package application

import (
	"context"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// Servico executa o consumo de um item livre.
// Não há porta para o OpenObserve.
type Servico struct {
	Fila     out.Fila
	Codigos  out.Codigos
	Disp     out.Disponibilidade
	Reservas *domain.Quadro
}

var _ in.Consumo = (*Servico)(nil)

// Executar toma um item e, com worker e instância no ar, conclui no mesmo passo.
func (s *Servico) Executar(ctx context.Context, pedido in.Pedido) (in.Saida, error) {
	if err := ctx.Err(); err != nil {
		return in.Saida{}, err
	}
	item, ok, err := s.Fila.TomarLivre(ctx)
	if err != nil {
		return in.Saida{FalhaLeitura: true}, err
	}
	if !ok {
		return in.Saida{}, nil
	}
	if !item.Aceito() {
		return in.Saida{ForaDeContrato: true}, nil
	}
	if s.Reservas != nil && !s.Reservas.Marcar(pedido.Instancia, item) {
		if err := s.Fila.EncerrarSemConclusao(ctx, item.ID); err != nil {
			return in.Saida{}, err
		}
		return in.Saida{}, nil
	}
	if item.Estado == "" {
		item.Estado = domain.Livre
	}
	emPasso, okInicio := domain.Iniciar(item)
	if !okInicio {
		emPasso = item
		emPasso.Estado = domain.EmPasso
	}
	if !pedido.WorkerNoAr || !pedido.InstanciaNoAr {
		_ = domain.Fim(emPasso, pedido.WorkerNoAr, pedido.InstanciaNoAr, false)
		if s.Reservas != nil {
			if pedido.InstanciaNoAr {
				s.Reservas.LiberarItem(item.ID)
			} else {
				s.Reservas.SoltarInstancia(pedido.Instancia)
			}
		}
		if err := s.Fila.Liberar(ctx, item.ID); err != nil {
			return in.Saida{}, err
		}
		return in.Saida{}, nil
	}
	aceito, err := s.Codigos.RegistrarPrimeira(ctx, item.Codigo)
	if err != nil {
		return in.Saida{}, err
	}
	if !aceito {
		_ = domain.Fim(emPasso, true, true, true)
		if s.Reservas != nil {
			s.Reservas.Encerrar(item.ID)
		}
		if err := s.Fila.EncerrarSemConclusao(ctx, item.ID); err != nil {
			return in.Saida{}, err
		}
		return in.Saida{}, nil
	}
	_ = domain.Fim(emPasso, true, true, false)
	if s.Reservas != nil {
		s.Reservas.ConcluirItem(item.ID)
	}
	if err := s.Fila.Concluir(ctx, item.ID); err != nil {
		return in.Saida{}, err
	}
	if err := s.Disp.Disponibilizar(ctx, item.Codigo); err != nil {
		return in.Saida{}, err
	}
	return in.Saida{Conclusoes: 1, CodigosDisponiveis: []domain.Codigo{item.Codigo}}, nil
}
