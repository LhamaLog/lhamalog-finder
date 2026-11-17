package out

import (
	"context"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// ConsultarAchadoModulo é invocada uma vez por achado antes do envio por lote.
// Estado da ocorrência pode ser usado pela 008; a 006 mínima não faz fusão histórica.
type ConsultarAchadoModulo interface {
	ConsultarAchado(context.Context, domain.IdentidadeAchado) ([]domain.EstadoOcorrencia, error)
}
