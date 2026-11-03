package openobserve

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

func TestLerDevolveHitInteiro(t *testing.T) {
	hit := []byte(`{"senha":"segredo","msg":"inteiro"}`)
	var recebido consulta
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/default/_search" {
			t.Errorf("caminho %s", r.URL.Path)
		}
		user, pass, ok := r.BasicAuth()
		if !ok || user != "root@example.com" || pass != "Complexpass#123" {
			t.Errorf("auth user=%s ok=%v", user, ok)
		}
		raw, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(raw, &recebido); err != nil {
			t.Errorf("corpo: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"total":1,"hits":[` + string(hit) + `]}`))
	}))
	defer srv.Close()

	c := Novo(srv.URL, "default", "root@example.com", "Complexpass#123")
	c.Agora = func() time.Time { return time.Unix(10, 0) }
	reg, estado, err := c.Ler(context.Background(), "otel_logs")
	if err != nil || estado != out.LeituraPronta || !bytes.Equal(reg, hit) {
		t.Fatalf("reg=%s estado=%v err=%v", reg, estado, err)
	}
	if recebido.Query.SQL != `SELECT * FROM "otel_logs"` || recebido.Query.Size != 1 {
		t.Fatalf("consulta: %+v", recebido.Query)
	}
}

func TestLerSemHitEAusente(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"total":0,"hits":[]}`))
	}))
	defer srv.Close()
	c := Novo(srv.URL, "default", "root@example.com", "x")
	reg, estado, err := c.Ler(context.Background(), "otel_logs")
	if estado != out.LeituraAusente || reg != nil || err == nil {
		t.Fatalf("reg=%s estado=%v err=%v", reg, estado, err)
	}
}

func TestLerRedeInacessivel(t *testing.T) {
	c := Novo("http://127.0.0.1:1", "default", "root@example.com", "x")
	reg, estado, err := c.Ler(context.Background(), "otel_logs")
	if estado != out.LeituraInacessivel || reg != nil || err == nil {
		t.Fatalf("reg=%s estado=%v err=%v", reg, estado, err)
	}
}

func TestLerStatus500Inacessivel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "falha interna", http.StatusInternalServerError)
	}))
	defer srv.Close()
	c := Novo(srv.URL, "default", "root@example.com", "x")
	_, estado, err := c.Ler(context.Background(), "otel_logs")
	if estado != out.LeituraInacessivel || err == nil || !bytes.Contains([]byte(err.Error()), []byte("500")) {
		t.Fatalf("estado=%v err=%v", estado, err)
	}
}

func TestLerApontamentoRecusado(t *testing.T) {
	c := Novo("http://127.0.0.1:5080", "default", "root@example.com", "x")
	_, estado, err := c.Ler(context.Background(), domain.Codigo(`nome'; drop`))
	if estado != out.LeituraInacessivel || err == nil {
		t.Fatalf("estado=%v err=%v", estado, err)
	}
}
