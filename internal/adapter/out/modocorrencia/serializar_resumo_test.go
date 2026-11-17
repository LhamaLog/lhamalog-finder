package modocorrencia

import (
	"encoding/json"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"testing"
)

func TestSerializarResumoProvisorioSemCamposBloqueados(t *testing.T) {
	r := domain.ResumoAnalise{BatchID: "b", TenantID: "t", ProductID: "p", VersaoRegras: "v", Contadores: domain.ContadoresResumo{Info: 1, LinhasProcessadas: 1}, Achados: []domain.AchadoAnalise{{Tipo: domain.INFO, Assinatura: "planejada-N3", Contagem: 9, Linhas: []string{"seguro"}}}}
	b, err := SerializarResumo(r)
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]any
	if err = json.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	if _, ok := v["execId"]; ok {
		t.Fatal("execId inventado")
	}
	c := v["contadores"].(map[string]any)
	if c["informacao"] != float64(1) {
		t.Fatal("contador INFO provisório nao mapeado")
	}
	a := v["achados"].([]any)[0].(map[string]any)
	for _, nome := range []string{"assinatura", "contagem", "gravidade", "conteudo", "raw"} {
		if _, ok := a[nome]; ok {
			t.Fatalf("campo bloqueado: %s", nome)
		}
	}
}

func TestFlagDeAchadoSemSemanticaNaoAfirmaAusenciaDeSensivel(t *testing.T) {
	r := domain.ResumoAnalise{BatchID: "b", TenantID: "t", ProductID: "p", VersaoRegras: "v", Contadores: domain.ContadoresResumo{Info: 1, LinhasProcessadas: 1, DadosSensiveisMascarados: 2}, Achados: []domain.AchadoAnalise{{Tipo: domain.INFO, Linhas: []string{"seguro"}}}}
	b, err := SerializarResumo(r)
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]any
	json.Unmarshal(b, &v)
	if v["achados"].([]any)[0].(map[string]any)["dadosSensiveis"] != nil {
		t.Fatal("campo sem semantica formal afirmou resultado")
	}
}
