package application

import (
	"context"
	"sync"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// ServicoCaptura obtém apontamentos da fonte mockada e deixa o log pronto.
// Não lê a fila e não nomeia armazenamento.
type ServicoCaptura struct {
	Fonte   out.FonteApontamentos
	Leitura out.LeituraLog
	Prontos out.LogsProntos
}

var _ in.Captura = (*ServicoCaptura)(nil)

// Executar captura cada apontamento distinto. A biblioteca padrão espera as capturas.
func (s *ServicoCaptura) Executar(ctx context.Context, pedido in.PedidoCaptura) (in.SaidaCaptura, error) {
	if err := ctx.Err(); err != nil {
		return in.SaidaCaptura{}, err
	}
	lista, err := s.Fonte.Obter(ctx)
	if err != nil {
		return in.SaidaCaptura{}, err
	}
	if len(lista) == 0 {
		return in.SaidaCaptura{}, nil
	}
	lista = domain.Unicos(lista)

	var (
		mu    sync.Mutex
		wg    sync.WaitGroup
		saida in.SaidaCaptura
		falha error
	)
	for _, codigo := range lista {
		wg.Add(1)
		go func(codigo domain.Codigo) {
			defer wg.Done()
			pronto, falhaCap, errUm := s.um(ctx, pedido, codigo)
			mu.Lock()
			defer mu.Unlock()
			if errUm != nil && falha == nil {
				falha = errUm
			}
			if falhaCap != nil {
				saida.Falhas = append(saida.Falhas, *falhaCap)
			}
			if pronto != nil {
				saida.Prontos = append(saida.Prontos, *pronto)
			}
		}(codigo)
	}
	wg.Wait()
	return saida, falha
}

func (s *ServicoCaptura) um(ctx context.Context, pedido in.PedidoCaptura, codigo domain.Codigo) (*domain.LogPronto, *domain.FalhaCaptura, error) {
	if domain.Inutilizavel(codigo) {
		f := domain.FalhaCaptura{Apontamento: codigo, Motivo: domain.MotivoInutilizavel}
		return nil, &f, nil
	}
	if !pedido.InstanciaNoAr || ctx.Err() != nil {
		return nil, nil, nil
	}
	reg, estado, err := s.Leitura.Ler(ctx, codigo)
	if ctx.Err() != nil || !pedido.InstanciaNoAr {
		return nil, nil, nil
	}
	switch estado {
	case out.LeituraAusente:
		f := domain.FalhaCaptura{Apontamento: codigo, Motivo: domain.MotivoAusente}
		return nil, &f, nil
	case out.LeituraInacessivel:
		f := domain.FalhaCaptura{Apontamento: codigo, Motivo: domain.MotivoInacessivel}
		return nil, &f, nil
	}
	if err != nil {
		return nil, nil, err
	}
	if err := s.Prontos.Guardar(ctx, codigo, reg); err != nil {
		return nil, nil, err
	}
	lido, ok, err := s.Prontos.Ler(ctx, codigo, true)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, nil
	}
	pronto := domain.LogPronto{Apontamento: codigo, Conteudo: lido}
	return &pronto, nil, nil
}
