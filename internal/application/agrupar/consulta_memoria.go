package agrupar

import (
	"context"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"sync"
)

// ConsultaMemoria é prova explícita da semântica do Módulo. Não é fallback
// de produção e não é ligada ao bootstrap; HTTP permanece DN-009-2.
type ConsultaMemoria struct {
	mu       sync.Mutex
	estados  map[domain.IdentidadeAchado][]domain.EstadoOcorrencia
	chamadas map[domain.IdentidadeAchado]int
}

func NovaConsultaMemoria(estados map[domain.IdentidadeAchado][]domain.EstadoOcorrencia) *ConsultaMemoria {
	c := &ConsultaMemoria{estados: map[domain.IdentidadeAchado][]domain.EstadoOcorrencia{}, chamadas: map[domain.IdentidadeAchado]int{}}
	for k, v := range estados {
		c.estados[k] = append([]domain.EstadoOcorrencia(nil), v...)
	}
	return c
}
func (c *ConsultaMemoria) ConsultarAchado(ctx context.Context, k domain.IdentidadeAchado) ([]domain.EstadoOcorrencia, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.chamadas == nil {
		c.chamadas = map[domain.IdentidadeAchado]int{}
	}
	c.chamadas[k]++
	return append([]domain.EstadoOcorrencia(nil), c.estados[k]...), nil
}
func (c *ConsultaMemoria) Chamadas(k domain.IdentidadeAchado) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.chamadas[k]
}
