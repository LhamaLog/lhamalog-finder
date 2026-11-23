package application

import (
	"context"
	"testing"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

func TestAnaliseMantemNivelOriginalMesmoQuandoMensagemElevaClasse(t *testing.T) {
	p := pedidoAnaliseTeste()
	p.Logs[0].NivelCampo = "INFO"
	p.Logs[1].NivelCampo = "WARN"
	p.Logs[2].NivelCampo = ""
	n := 0
	s := ServicoAnalise{Seguranca: scannerAnaliseTeste(&n)}
	r, err := s.Montar(context.Background(), p)
	if err != nil || len(r.Achados) != 3 || n != 1 {
		t.Fatalf("achados=%v validacoes=%d erro=%v", r.Achados, n, err)
	}
	for i, nivel := range []string{"INFO", "WARN", ""} {
		a := r.Achados[i]
		if a.Tipo != domain.WARNING || a.Assinatura != nivel+"|seguro warning|svc|host" || a.Contagem != 1 {
			t.Fatalf("nivel original perdido: %+v", a)
		}
	}
}
