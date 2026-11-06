package in

import (
	"context"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// Pedido descreve quem executa o passo e se worker e instância seguem no ar.
type Pedido struct {
	Instancia     string
	WorkerNoAr    bool
	InstanciaNoAr bool
}

// Saida é o efeito observável de uma tentativa de consumo.
type Saida struct {
	Conclusoes         int
	CodigosDisponiveis []domain.Codigo
	ForaDeContrato     bool
	FalhaLeitura       bool
}

// Consumo é a porta de entrada do caso de uso.
type Consumo interface {
	Executar(ctx context.Context, pedido Pedido) (Saida, error)
}
