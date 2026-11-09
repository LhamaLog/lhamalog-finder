package application

import (
	"context"
	"errors"
	"log/slog"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// Servico executa o consumo de um item livre.
// Não há porta para o OpenObserve.
// Codigos e Disp continuam na montagem. A conclusão grava o código e tira a pendência numa operação só da fila. Essa operação é a disponibilidade.
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
		if err := s.Fila.EncerrarSemConclusao(ctx, item.ID); err != nil {
			return in.Saida{ForaDeContrato: true}, err
		}
		slog.Info("item descartado", "id", string(item.ID))
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
	nova, err := s.Fila.ConcluirPrimeira(ctx, item.ID, item.Codigo)
	if err != nil {
		return in.Saida{}, err
	}
	if !nova {
		_ = domain.Fim(emPasso, true, true, true)
		if s.Reservas != nil {
			s.Reservas.Encerrar(item.ID)
		}
		return in.Saida{}, nil
	}
	_ = domain.Fim(emPasso, true, true, false)
	if s.Reservas != nil {
		s.Reservas.ConcluirItem(item.ID)
	}
	return in.Saida{Conclusoes: 1, CodigosDisponiveis: []domain.Codigo{item.Codigo}}, nil
}

// ConsumoReferencia interpreta log.ingested e entrega a referência à porta de efeito.
// A entrada válida só é confirmada depois do efeito concluído.
// Não há consulta ao OpenObserve nem uso de correlationId ou do id Redis como idempotência.
type ConsumoReferencia struct {
	Fila   out.Entradas
	Chaves out.ChaveLote
	Efeito out.Efeito
	dono   map[string]string
}

// SaidaLote é o efeito observável de uma entrada do stream de logs.
type SaidaLote struct {
	Confirmacoes int
	Efeitos      int
	Invalida     bool
	FalhaLeitura bool
}

// Executar lê uma entrada. Inválida é confirmada uma vez. Válida espera o efeito.
func (s *ConsumoReferencia) Executar(ctx context.Context) (SaidaLote, error) {
	if err := ctx.Err(); err != nil {
		return SaidaLote{}, err
	}
	entrada, ok, err := s.Fila.Ler(ctx)
	if err != nil {
		return SaidaLote{FalhaLeitura: true}, err
	}
	if !ok {
		return SaidaLote{}, nil
	}
	ref, err := domain.InterpretarEnvelope(entrada.Campos)
	if err != nil {
		return s.encerrarInvalida(ctx, entrada, err)
	}
	ref.EntradaID = entrada.ID
	return s.tratarValida(ctx, entrada, ref)
}

func (s *ConsumoReferencia) encerrarInvalida(ctx context.Context, entrada domain.Entrada, err error) (SaidaLote, error) {
	var erro *domain.ErroEnvelope
	idem := ""
	motivo := err.Error()
	if errors.As(err, &erro) {
		idem = erro.IdempotencyID
		motivo = erro.Motivo
	}
	if idem != "" {
		delete(s.dono, idem)
	}
	args := []any{"id", entrada.ID, "motivo", motivo}
	if idem != "" {
		args = append(args, "idempotencyId", idem)
	}
	slog.Info("entrada inválida", args...)
	if err := s.Fila.EncerrarInvalida(ctx, entrada.ID, idem); err != nil {
		return SaidaLote{}, err
	}
	return SaidaLote{Confirmacoes: 1, Invalida: true}, nil
}

func (s *ConsumoReferencia) tratarValida(ctx context.Context, entrada domain.Entrada, ref domain.Referencia) (SaidaLote, error) {
	estado, err := s.Chaves.Consultar(ctx, ref.IdempotencyID)
	if err != nil {
		return SaidaLote{}, err
	}
	if estado == domain.EstadoConcluido {
		return s.confirmarSemEfeito(ctx, entrada.ID, ref.IdempotencyID)
	}
	if estado == domain.EstadoEmProcessamento {
		return s.retomar(ctx, entrada, ref)
	}
	criou, err := s.Chaves.Reservar(ctx, ref.IdempotencyID)
	if err != nil {
		return SaidaLote{}, err
	}
	if !criou {
		return s.retomar(ctx, entrada, ref)
	}
	s.definirDono(ref.IdempotencyID, entrada.ID)
	return s.executarEfeito(ctx, entrada, ref)
}

func (s *ConsumoReferencia) retomar(ctx context.Context, entrada domain.Entrada, ref domain.Referencia) (SaidaLote, error) {
	estado, err := s.Chaves.Consultar(ctx, ref.IdempotencyID)
	if err != nil {
		return SaidaLote{}, err
	}
	if estado == domain.EstadoConcluido {
		return s.confirmarSemEfeito(ctx, entrada.ID, ref.IdempotencyID)
	}
	dono := s.dono[ref.IdempotencyID]
	if (dono == entrada.ID || dono == "") && (entrada.Ociosa || entrada.Reivindicada) {
		s.definirDono(ref.IdempotencyID, entrada.ID)
		return s.executarEfeito(ctx, entrada, ref)
	}
	if err := s.Fila.Adiar(ctx, entrada.ID); err != nil {
		return SaidaLote{}, err
	}
	return SaidaLote{}, nil
}

func (s *ConsumoReferencia) executarEfeito(ctx context.Context, entrada domain.Entrada, ref domain.Referencia) (SaidaLote, error) {
	if s.Efeito == nil {
		return SaidaLote{}, nil
	}
	concluido, err := s.Efeito.Executar(ctx, ref)
	if err != nil || !concluido {
		return SaidaLote{}, err
	}
	if err := s.Chaves.Concluir(ctx, ref.IdempotencyID); err != nil {
		return SaidaLote{Efeitos: 1}, err
	}
	delete(s.dono, ref.IdempotencyID)
	if err := s.Fila.Confirmar(ctx, entrada.ID); err != nil {
		return SaidaLote{Efeitos: 1}, err
	}
	return SaidaLote{Confirmacoes: 1, Efeitos: 1}, nil
}

func (s *ConsumoReferencia) confirmarSemEfeito(ctx context.Context, id, idem string) (SaidaLote, error) {
	if err := s.Fila.Confirmar(ctx, id); err != nil {
		return SaidaLote{}, err
	}
	delete(s.dono, idem)
	return SaidaLote{Confirmacoes: 1}, nil
}

func (s *ConsumoReferencia) definirDono(idem, id string) {
	if s.dono == nil {
		s.dono = map[string]string{}
	}
	s.dono[idem] = id
}
