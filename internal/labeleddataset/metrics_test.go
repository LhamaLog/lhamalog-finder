package labeleddataset

import (
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"strings"
	"testing"
)

func TestMetricasPorApontamentoUsamMaisGrave(t *testing.T) {
	member := strings.Replace(logValido, `,"class_label":"ERROR"`, "", 1)
	ap := `{"itemType":"apontamento","class_label":"WARNING","split":"validacao","lhamalog_batch_id":"teste-lote","tenantId":"teste-tenant","productId":"teste-produto","documentacao":"gold WARNING com membro ERROR deve errar","members":[` + member + `]}`
	rs, err := LoadConjunto(strings.NewReader(logValido + "\n" + ap))
	if err != nil {
		t.Fatal(err)
	}
	m, err := Avaliar(rs)
	if err != nil {
		t.Fatal(err)
	}
	if m.Unidades != 2 || m.Classes[domain.ERROR].TP != 1 || m.Classes[domain.ERROR].FP != 1 || m.Classes[domain.WARNING].FN != 1 {
		t.Fatalf("metricas: %+v", m)
	}
	if m.ValidarGate() == nil {
		t.Fatal("gate aceitou qualidade baixa")
	}
}
func TestGateRejeitaVazioOuClasseSemSuporte(t *testing.T) {
	if _, err := Avaliar(nil); err == nil {
		t.Fatal("avaliacao vazia aceita")
	}
	r, err := LoadConjunto(strings.NewReader(logValido))
	if err != nil {
		t.Fatal(err)
	}
	m, err := Avaliar(r)
	if err != nil {
		t.Fatal(err)
	}
	if m.ValidarGate() == nil {
		t.Fatal("gate aprovou classes sem suporte")
	}
}
func TestGateTresClassesEmFixtureNaoEhAceiteDaBase(t *testing.T) {
	var linhas []string
	for _, c := range []string{"ERROR", "WARNING", "INFO"} {
		linhas = append(linhas, strings.ReplaceAll(logValido, "ERROR", c))
	}
	rs, err := LoadConjunto(strings.NewReader(strings.Join(linhas, "\n")))
	if err != nil {
		t.Fatal(err)
	}
	m, err := Avaliar(rs)
	if err != nil || m.ValidarGate() != nil {
		t.Fatalf("calculo no fixture: %v", err)
	}
}
