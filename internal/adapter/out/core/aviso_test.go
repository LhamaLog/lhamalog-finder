package core

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

func TestAvisoHTTPSExigeReciboCorrespondente(t *testing.T) {
	good := `{"aceito":true,"tenantId":"t","productId":"p","batchId":"b","motivo":"consulta_openobserve_indisponivel"}`
	for _, c := range []struct {
		nome, body string
		status     int
		ok         bool
	}{
		{"confirmado", good, 200, true}, {"criado", good, 201, true}, {"assincrono", good, 202, false}, {"vazio", "", 204, false},
		{"outro tenant", `{"aceito":true,"tenantId":"outro","productId":"p","batchId":"b","motivo":"consulta_openobserve_indisponivel"}`, 200, false},
		{"alias", `{"Aceito":true,"tenantId":"t","productId":"p","batchId":"b","motivo":"consulta_openobserve_indisponivel"}`, 200, false},
		{"duplicada", `{"aceito":false,"aceito":true,"tenantId":"t","productId":"p","batchId":"b","motivo":"consulta_openobserve_indisponivel"}`, 200, false},
	} {
		t.Run(c.nome, func(t *testing.T) {
			calls := 0
			s := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != "POST" || r.URL.Path != "/internal/v1/alertas/finder" || r.Header.Get("Authorization") != "Bearer teste" {
					t.Error("transporte/rota inválidos")
				}
				var p map[string]any
				if json.NewDecoder(r.Body).Decode(&p) != nil || len(p) != 4 || p["tenantId"] != "t" || p["batchId"] != "b" {
					t.Error("payload inválido")
				}
				w.WriteHeader(c.status)
				w.Write([]byte(c.body))
			}))
			defer s.Close()
			a, err := NovoAviso(s.URL, "teste", time.Second, domain.Referencia{TenantID: "t", ProductID: "p", BatchID: "b"}, s.Client().Transport)
			if err != nil {
				t.Fatal(err)
			}
			err = a.AvisarFalhaOpenObserve(context.Background(), "b", "consulta_openobserve_indisponivel")
			if (err == nil) != c.ok {
				t.Fatalf("erro=%v", err)
			}
			if a.AvisarFalhaOpenObserve(context.Background(), "outro", "consulta_openobserve_indisponivel") == nil || calls != 1 {
				t.Fatal("aceitou troca da identidade")
			}
		})
	}
	if _, err := NovoAviso("http://localhost", "teste", time.Second, domain.Referencia{TenantID: "t", ProductID: "p", BatchID: "b"}, nil); err == nil {
		t.Fatal("aceitou HTTP")
	}
}
