package modocorrencia

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

func TestSerializarResumoDoContratoAprovado(t *testing.T) {
	r := resumoHTTPTeste()
	revisao := "rev-8"
	r.Achados = []domain.AchadoAnalise{{Tipo: domain.INFO, Assinatura: "INFO|seguro|svc|host", HostName: "host", Contagem: 9, PrimeiraOcorrencia: time.Unix(1, 0), UltimaOcorrencia: time.Unix(2, 0), RevisaoConsultada: &revisao}}
	b, err := SerializarResumo(r)
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]any
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	if _, ok := v["execId"]; ok {
		t.Fatal("execId inventado")
	}
	if v["contadores"].(map[string]any)["informacao"] != float64(1) {
		t.Fatal("contador de INFO incorreto")
	}
	a := v["achados"].([]any)[0].(map[string]any)
	if a["assinatura"] != r.Achados[0].Assinatura || a["contagem"] != float64(9) || a["revisaoConsultada"] != revisao || a["primeiraOcorrencia"] != "1970-01-01T00:00:01Z" || a["ultimaOcorrencia"] != "1970-01-01T00:00:02Z" {
		t.Fatalf("campos N3 ou revisão incorretos: %v", a)
	}
	for _, nome := range []string{"gravidade", "criticidade", "conteudo", "raw", "versaoRegras"} {
		if _, ok := a[nome]; ok {
			t.Fatalf("campo fora do contrato: %s", nome)
		}
	}
	for _, nome := range []string{"titulo", "resumo", "numeroLinha", "camposExtraido", "dadosSensiveis"} {
		if valor, ok := a[nome]; !ok || valor != nil {
			t.Fatalf("campo desconhecido %s deveria ser null", nome)
		}
	}
	if linhas, ok := a["linhas"].([]any); !ok || len(linhas) != 0 {
		t.Fatal("linhas deve ser array vazio")
	}
}

func TestSerializarResumoNaoAplicaEstadoExistenteSemRevisao(t *testing.T) {
	r := resumoHTTPTeste()
	r.Achados = []domain.AchadoAnalise{{Tipo: domain.INFO, Assinatura: "a", HostName: "h", Contagem: 1, PrimeiraOcorrencia: time.Unix(1, 0), UltimaOcorrencia: time.Unix(1, 0), StatusAnterior: "ABERTA"}}
	if _, err := SerializarResumo(r); !errors.Is(err, out.ErrModuloRequisicao) {
		t.Fatalf("falta de revisão deveria impedir entrega: %v", err)
	}
	r.Achados[0].StatusAnterior = ""
	if _, err := SerializarResumo(r); err != nil {
		t.Fatalf("consulta sem ocorrência existente deve permitir entrega: %v", err)
	}
}
