package in

import (
	"context"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// PedidoCaptura informa se a instância segue no ar durante a captura.
// SegueNoAr, quando informado, é lido de novo depois da leitura.
// Sem ele, vale InstanciaNoAr.
type PedidoCaptura struct {
	InstanciaNoAr bool
	SegueNoAr     func() bool
}

// SaidaCaptura é o efeito observável de uma execução.
type SaidaCaptura struct {
	Prontos []domain.LogPronto
	Falhas  []domain.FalhaCaptura
}

// Captura é a porta de entrada. A fila não é entrada.
type Captura interface {
	Executar(ctx context.Context, pedido PedidoCaptura) (SaidaCaptura, error)
}
