package in

import (
	"context"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"time"
)

// LogAnalise contém campos já extraídos pelo adaptador; Mensagem só entra no scanner.
// Nenhum mapeamento físico OTel adicional é adivinhado nesta porta lógica.
type LogAnalise struct {
	Apontamento                                                domain.Codigo
	NivelCampo, Mensagem, ServiceName, Host, TraceID, ModuloID string
	Timestamp                                                  time.Time
}
type PedidoAnalise struct {
	Referencia                                                       domain.Referencia
	ProdutoBody, SchemaVersion, ExecID, VersaoRegras, VersaoCatalogo string
	Logs                                                             []LogAnalise
	Completo, ArmazenamentoInacessivel                               bool
	Catalogo                                                         []domain.RegraClassificacao
}
type ResultadoAnalise struct {
	Resumo    domain.ResumoAnalise
	Concluida bool
}
type Analise interface {
	Analisar(context.Context, PedidoAnalise) (ResultadoAnalise, error)
}
