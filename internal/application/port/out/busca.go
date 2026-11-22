package out

import (
	"context"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// TipoFalhaBusca é o que a consulta de uma página informa, sem credencial e sem hits.
type TipoFalhaBusca int

const (
	FalhaCredencial TipoFalhaBusca = iota + 1
	FalhaResposta
	FalhaTimeout
	FalhaIndisponivel
	FalhaPagina
)

// ErroBusca é a falha de uma página. O texto não inclui senha nem conteúdo de log.
type ErroBusca struct {
	Tipo TipoFalhaBusca
}

func (e *ErroBusca) Error() string {
	if e == nil {
		return ""
	}
	switch e.Tipo {
	case FalhaCredencial:
		return "credencial recusada"
	case FalhaResposta:
		return "resposta inválida"
	case FalhaTimeout:
		return "timeout"
	case FalhaIndisponivel:
		return "indisponibilidade"
	case FalhaPagina:
		return "falha de página"
	default:
		return "falha de busca"
	}
}

// PedidoPagina é uma página já montada. A porta não escolhe lote, faixa nem versão.
type PedidoPagina struct {
	URL         string
	Usuario     string
	Senha       string
	Timeout     time.Duration
	Organizacao string
	SQL         string
	Inicio      int64
	Fim         int64
	From        int
	Size        int
}

// Busca consulta uma página e devolve os hits.
// Não confirma a entrada Redis, não oculta campo e não escolhe versão divergente.
type Busca interface {
	Consultar(ctx context.Context, pedido PedidoPagina) ([]domain.DocumentoFisico, error)
}
