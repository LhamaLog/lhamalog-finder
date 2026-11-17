package openobserve

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
)

func TestRespostaAcimaDe32MiBNaoTemTetoInventado(t *testing.T) {
	texto := strings.Repeat("x", (32<<20)+1)
	docs, err := lerHits(strings.NewReader(`{"hits":[{"mensagem":"` + texto + `"}]}`))
	if err != nil || len(docs) != 1 || len(docs[0]["mensagem"].(string)) != len(texto) {
		t.Fatalf("resposta grande: documentos=%d erro=%v", len(docs), err)
	}
}

type leitorFalho struct{}

func (leitorFalho) Read([]byte) (int, error) {
	return 0, errors.New("falha de transporte com conteudo privado")
}

func TestErroDeLeituraNaoEhRespostaInvalida(t *testing.T) {
	_, err := lerHits(io.MultiReader(strings.NewReader(`{"hits":[`), leitorFalho{}))
	falha, ok := err.(*out.ErroBusca)
	if !ok || falha.Tipo != out.FalhaPagina {
		t.Fatalf("erro de leitura: %v", err)
	}
	if strings.Contains(err.Error(), "privado") {
		t.Fatal("erro expos conteudo")
	}
}

func TestStatusNaoPrevistoMantemPendencia(t *testing.T) {
	for _, status := range []int{302, 400, 404, 409, 418, 422} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(status); io.WriteString(w, `{"hits":[]}`) }))
			defer srv.Close()
			_, err := (Cliente{}).Consultar(context.Background(), out.PedidoPagina{URL: srv.URL, Organizacao: "org", SQL: `SELECT * FROM "s"`, Size: 1, Timeout: time.Second})
			falha, ok := err.(*out.ErroBusca)
			if !ok || falha.Tipo != out.FalhaPagina {
				t.Fatalf("status %d: %v", status, err)
			}
		})
	}
}

func TestCorpoHTTPTruncadoMantemPendencia(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		io.WriteString(w, `{"hits":[`)
	}))
	defer srv.Close()
	_, err := (Cliente{}).Consultar(context.Background(), out.PedidoPagina{URL: srv.URL, Organizacao: "org", SQL: `SELECT * FROM "s"`, Size: 1, Timeout: time.Second})
	falha, ok := err.(*out.ErroBusca)
	if !ok || falha.Tipo != out.FalhaPagina {
		t.Fatalf("corpo truncado: %v", err)
	}
}

func TestJSONDeSucessoExigeFimValido(t *testing.T) {
	for _, corpo := range []string{`{"hits":[]} }`, `{"hits":[]} ]`, `{"hits":[]} {"hits":[]}`} {
		if _, err := lerHits(strings.NewReader(corpo)); err == nil {
			t.Fatal("JSON com sufixo invalido aceito")
		}
	}
}
