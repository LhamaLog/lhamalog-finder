package recorrencia

import "github.com/LhamaLog/lhamalog-finder/internal/domain"

const JanelaLotes = 500

// Avaliar segue quickstart 008 §7: >=3 lotes distintos, mesma máquina/chave.
// A porta entrega ordem recente→antiga; o adaptador operacional ainda está bloqueado.
func Avaliar(historico []domain.EvidenciaRecorrencia, chave domain.IdentidadeAchado, host string) bool {
	janela := map[string]bool{}
	equivalentes := map[string]bool{}
	for _, e := range historico {
		if e.BatchID == "" {
			continue
		}
		if !janela[e.BatchID] {
			if len(janela) == JanelaLotes {
				break
			}
			janela[e.BatchID] = true
		}
		if e.Host == host && e.Chave == chave {
			equivalentes[e.BatchID] = true
		}
	}
	return len(equivalentes) >= 3
}
