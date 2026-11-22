package recorrencia

import (
	"fmt"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"testing"
)

func TestRecorrenciaDistintaMesmaMaquinaJanela500(t *testing.T) {
	chave := domain.IdentidadeAchado{Assinatura: "s", Classe: domain.ERROR, TenantID: "t", ProductID: "p"}
	r := []domain.EvidenciaRecorrencia{{BatchID: "a", Host: "h", Chave: chave}, {BatchID: "a", Host: "h", Chave: chave}, {BatchID: "b", Host: "h", Chave: chave}}
	if Avaliar(r, chave, "h") {
		t.Fatal("repeticao no mesmo lote contou")
	}
	r = append(r, domain.EvidenciaRecorrencia{BatchID: "c", Host: "h", Chave: chave})
	if !Avaliar(r, chave, "h") {
		t.Fatal("3 lotes nao contaram")
	}
	if Avaliar(r, chave, "outro") {
		t.Fatal("maquina diferente contou")
	}
}

func TestJanelaContaLotesNaoRegistros(t *testing.T) {
	chave := domain.IdentidadeAchado{Assinatura: "s", Classe: domain.INFO, TenantID: "t", ProductID: "p"}
	var h []domain.EvidenciaRecorrencia
	for i := 0; i < 501; i++ {
		c := domain.IdentidadeAchado{}
		if i == 0 || i == 1 || i == 500 {
			c = chave
		}
		h = append(h, domain.EvidenciaRecorrencia{BatchID: fmt.Sprint(i), Host: "h", Chave: c})
	}
	if Avaliar(h, chave, "h") {
		t.Fatal("lote alem da janela contou")
	}
	h[499].Chave = chave
	if !Avaliar(h, chave, "h") {
		t.Fatal("lote dentro da janela nao contou")
	}
}
