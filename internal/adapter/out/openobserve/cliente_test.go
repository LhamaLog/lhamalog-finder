package openobserve

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
)

func TestClienteEnviaBasicENaoGravaTimeout(t *testing.T) {
	senha := "segredo-nao-logar"
	var corpo []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/org-a/_search" {
			t.Errorf("pedido %s %s", r.Method, r.URL.Path)
		}
		esperado := "Basic " + base64.StdEncoding.EncodeToString([]byte("usuario:"+senha))
		if r.Header.Get("Authorization") != esperado {
			t.Error("autorização básica ausente")
		}
		lido, _ := io.ReadAll(r.Body)
		corpo = lido
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"hits":[{"lhamalog_record_id":"a","mensagem":"linha"}]}`)
	}))
	defer srv.Close()

	docs, err := (Cliente{}).Consultar(context.Background(), out.PedidoPagina{
		URL: srv.URL, Usuario: "usuario", Senha: senha, Timeout: time.Second,
		Organizacao: "org-a", SQL: `SELECT * FROM "stream_sintetico"`, Inicio: 100, Fim: 201, From: 0, Size: 2,
	})
	if err != nil || len(docs) != 1 || docs[0]["mensagem"] != "linha" {
		t.Fatalf("docs=%v err=%v", len(docs), err)
	}
	if strings.Contains(string(corpo), "timeout") || strings.Contains(string(corpo), senha) {
		t.Fatal("corpo levou timeout ou senha")
	}
}

func TestClienteClassificaFalhaSemConteudo(t *testing.T) {
	senha := "segredo-nao-logar"
	casos := []struct {
		nome   string
		status int
		corpo  string
		tipo   out.TipoFalhaBusca
	}{
		{nome: "credencial", status: http.StatusUnauthorized, corpo: `{"erro":"` + senha + `","hits":[{"mensagem":"linha-secreta"}]}`, tipo: out.FalhaCredencial},
		{nome: "invalida", status: http.StatusOK, corpo: `nao-json`, tipo: out.FalhaResposta},
		{nome: "pagina", status: http.StatusInternalServerError, corpo: `{"hits":[{"mensagem":"linha-secreta"}]}`, tipo: out.FalhaPagina},
		{nome: "indisponivel", status: http.StatusServiceUnavailable, corpo: `fora`, tipo: out.FalhaIndisponivel},
	}
	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(caso.status)
				_, _ = io.WriteString(w, caso.corpo)
			}))
			defer srv.Close()
			docs, err := (Cliente{}).Consultar(context.Background(), out.PedidoPagina{
				URL: srv.URL, Usuario: "usuario", Senha: senha, Timeout: time.Second,
				Organizacao: "org-a", SQL: "SELECT * FROM \"s\"", From: 0, Size: 1,
			})
			falha, ok := err.(*out.ErroBusca)
			if !ok || falha.Tipo != caso.tipo || docs != nil {
				t.Fatalf("tipo=%v err=%v", caso.tipo, err)
			}
			if strings.Contains(err.Error(), senha) || strings.Contains(err.Error(), "linha-secreta") {
				t.Fatal("erro expõe senha ou conteúdo")
			}
		})
	}
}

func TestClienteTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_, _ = io.WriteString(w, `{"hits":[]}`)
	}))
	defer srv.Close()
	_, err := (Cliente{}).Consultar(context.Background(), out.PedidoPagina{
		URL: srv.URL, Usuario: "usuario", Senha: "segredo-nao-logar", Timeout: 30 * time.Millisecond,
		Organizacao: "org-a", SQL: "SELECT * FROM \"s\"", From: 0, Size: 1,
	})
	falha, ok := err.(*out.ErroBusca)
	if !ok || falha.Tipo != out.FalhaTimeout {
		t.Fatalf("err=%v", err)
	}
	if strings.Contains(err.Error(), "segredo-nao-logar") {
		t.Fatal("timeout expõe a senha")
	}
}

func TestClienteIndisponivel(t *testing.T) {
	_, err := (Cliente{}).Consultar(context.Background(), out.PedidoPagina{
		URL: "http://127.0.0.1:1", Usuario: "usuario", Senha: "segredo-nao-logar", Timeout: 200 * time.Millisecond,
		Organizacao: "org-a", SQL: "SELECT * FROM \"s\"", From: 0, Size: 1,
	})
	falha, ok := err.(*out.ErroBusca)
	if !ok || (falha.Tipo != out.FalhaIndisponivel && falha.Tipo != out.FalhaTimeout) {
		t.Fatalf("err=%v", err)
	}
}

func TestHitsVaziosSaoPaginaCurta(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"hits":[]}`)
	}))
	defer srv.Close()
	docs, err := (Cliente{}).Consultar(context.Background(), out.PedidoPagina{
		URL: srv.URL, Usuario: "u", Senha: "s", Timeout: time.Second,
		Organizacao: "org-a", SQL: "SELECT * FROM \"s\"", From: 4, Size: 2,
	})
	if err != nil || docs == nil || len(docs) != 0 {
		t.Fatalf("docs=%v err=%v", docs, err)
	}
}
