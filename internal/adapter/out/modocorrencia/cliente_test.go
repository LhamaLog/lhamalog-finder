package modocorrencia

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

const reciboTeste = `{"aceito":true,"reciboId":"recibo","tenantId":"t","productId":"p","batchId":"b","repeticao":false}`

func clienteTeste(t *testing.T, h http.HandlerFunc) *Cliente {
	t.Helper()
	srv := httptest.NewTLSServer(h)
	t.Cleanup(srv.Close)
	c, err := NovoCliente(OpcoesCliente{URLBase: srv.URL, JWT: "jwt-de-teste", Transporte: srv.Client().Transport})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestClienteConsultaAchado(t *testing.T) {
	for _, caso := range []struct {
		nome, corpo string
		total       int
	}{
		{"existente", `{"estado":{"contagemAcumulada":10,"primeiraOcorrencia":"2026-11-19T09:00:00Z","ultimaOcorrencia":"2026-11-19T10:00:00.123456Z","status":"FECHADA","revisao":"rev-10"}}`, 1},
		{"ausente", `{"estado":null}`, 0},
	} {
		t.Run(caso.nome, func(t *testing.T) {
			c := clienteTeste(t, func(w http.ResponseWriter, req *http.Request) {
				if req.Method != "POST" || req.URL.Path != "/internal/v1/ocorrencias/consultas" || req.URL.RawQuery != "" || req.Header.Get("Authorization") != "Bearer jwt-de-teste" {
					t.Errorf("requisição fora do contrato")
				}
				var pedido map[string]string
				if err := json.NewDecoder(req.Body).Decode(&pedido); err != nil || len(pedido) != 4 || pedido["tipo"] != "INFO" || pedido["tenantId"] != "t" || pedido["productId"] != "p" || pedido["assinatura"] != "INFO|seguro|svc|host" {
					t.Errorf("identidade da consulta inválida: %v %v", pedido, err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(caso.corpo))
			})
			estados, err := c.ConsultarAchado(context.Background(), domain.IdentidadeAchado{TenantID: "t", ProductID: "p", Classe: domain.INFO, Assinatura: "INFO|seguro|svc|host"})
			if err != nil || len(estados) != caso.total {
				t.Fatalf("consulta: %v %v", estados, err)
			}
			if len(estados) == 1 && (estados[0].ContagemAcumulada != 10 || estados[0].Revisao != "rev-10" || estados[0].Status != "FECHADA" || estados[0].UltimaGlobal.Nanosecond() != 123456000) {
				t.Fatalf("estado incompleto: %+v", estados)
			}
		})
	}
}

func TestClienteRejeitaConsultaSemEstado(t *testing.T) {
	for _, corpo := range []string{`{}`, `{"estado":{}}`, `{"estado":null} {}`, `{"estado":{"contagemAcumulada":1,"primeiraOcorrencia":"2026-11-19T09:00:00+00:00","ultimaOcorrencia":"2026-11-19T10:00:00Z","status":"ABERTA","revisao":"rev"}}`} {
		c := clienteTeste(t, func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte(corpo)) })
		if _, err := c.ConsultarAchado(context.Background(), domain.IdentidadeAchado{TenantID: "t", ProductID: "p", Classe: domain.INFO, Assinatura: "a"}); err == nil {
			t.Fatalf("resposta inválida aceita: %s", corpo)
		}
	}
}

func resumoHTTPTeste() domain.ResumoAnalise {
	return domain.ResumoAnalise{BatchID: "b", TenantID: "t", ProductID: "p", VersaoRegras: "v", Contadores: domain.ContadoresResumo{Info: 1, LinhasProcessadas: 1}}
}

func TestClienteSoConfirmaReciboCorrespondente(t *testing.T) {
	for _, caso := range []struct {
		nome   string
		status int
		corpo  string
		aceito bool
	}{
		{"primeira", 201, reciboTeste, true},
		{"repetida", 200, strings.Replace(reciboTeste, "false", "true", 1), true},
		{"status divergente", 200, reciboTeste, false},
		{"recebido sem persistência", 202, reciboTeste, false},
		{"sem recibo", 204, "", false},
		{"outro tenant", 201, strings.Replace(reciboTeste, `"t"`, `"outro"`, 1), false},
		{"execId não enviado", 201, strings.Replace(reciboTeste, "}", `,"execId":"mais"}`, 1), false},
		{"recibo vazio", 201, strings.Replace(reciboTeste, `"recibo"`, `""`, 1), false},
		{"repetição ausente", 201, strings.Replace(reciboTeste, `,"repeticao":false`, "", 1), false},
		{"chave JSON duplicada", 201, strings.Replace(reciboTeste, "}", `,"aceito":true}`, 1), false},
	} {
		t.Run(caso.nome, func(t *testing.T) {
			c := clienteTeste(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/internal/v1/resumos-analise" || r.Header.Get("Authorization") != "Bearer jwt-de-teste" || r.Header.Get("Idempotency-Key") != "" {
					t.Error("transporte fora do contrato")
				}
				var pedido map[string]any
				if err := json.NewDecoder(r.Body).Decode(&pedido); err != nil || pedido["batchId"] != "b" || pedido["execId"] != nil {
					t.Errorf("chave não corresponde: %v %v", pedido, err)
				}
				w.WriteHeader(caso.status)
				w.Write([]byte(caso.corpo))
			})
			aceito, err := c.Enviar(context.Background(), resumoHTTPTeste())
			if aceito != caso.aceito || (err == nil) != caso.aceito {
				t.Fatalf("aceito=%v err=%v", aceito, err)
			}
		})
	}
}

func TestClientePropagaConflitoERespeitaRetryAfter(t *testing.T) {
	var n atomic.Int32
	c := clienteTeste(t, func(w http.ResponseWriter, _ *http.Request) {
		if n.Add(1) == 1 {
			w.WriteHeader(409)
			w.Write([]byte(`{"codigo":"REVISAO_DESATUALIZADA"}`))
			return
		}
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(429)
	})
	if _, err := c.Enviar(context.Background(), resumoHTTPTeste()); !errors.Is(err, out.ErrRevisaoDesatualizada) {
		t.Fatalf("conflito: %v", err)
	}
	for i := 0; i < 2; i++ {
		if _, err := c.Enviar(context.Background(), resumoHTTPTeste()); !errors.Is(err, out.ErrModuloTemporario) {
			t.Fatalf("429: %v", err)
		}
	}
	if n.Load() != 2 {
		t.Fatalf("reenviou antes de Retry-After: %d chamadas", n.Load())
	}
}

func TestClienteSanitizaErrosENaoSegueRedirect(t *testing.T) {
	var redirects atomic.Int32
	destino := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { redirects.Add(1) }))
	defer destino.Close()
	c := clienteTeste(t, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, destino.URL, http.StatusTemporaryRedirect)
	})
	if _, err := c.Enviar(context.Background(), resumoHTTPTeste()); err == nil || strings.Contains(err.Error(), destino.URL) || redirects.Load() != 0 {
		t.Fatalf("redirect inesperado: %v", err)
	}
	c.TimeoutEntrega = 20 * time.Millisecond
	c.http.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		<-req.Context().Done()
		return nil, errors.New("segredo de transporte")
	})
	if _, err := c.Enviar(context.Background(), resumoHTTPTeste()); !errors.Is(err, out.ErrModuloTemporario) || strings.Contains(err.Error(), "segredo") {
		t.Fatalf("timeout: %v", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
