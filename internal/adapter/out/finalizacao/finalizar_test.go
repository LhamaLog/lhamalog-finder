package finalizacao

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

type chavesTeste struct{ conclusoes int }

func (c *chavesTeste) Reservar(context.Context, string) (bool, error)    { return true, nil }
func (c *chavesTeste) Concluir(context.Context, string) error            { c.conclusoes++; return nil }
func (c *chavesTeste) Consultar(context.Context, string) (string, error) { return "", nil }
func (c *chavesTeste) RemoverReserva(context.Context, string) error      { return nil }

func TestFinalizacaoV1LocalESemExecucaoInventada(t *testing.T) {
	c := &chavesTeste{}
	ref := domain.Referencia{BatchID: "b", IdempotencyID: "b"}
	f := Finalizador{Chaves: c, Ref: ref}
	if err := f.ConcluirAnalise(context.Background(), domain.ChaveResumo{BatchID: "b"}); err != nil || c.conclusoes != 1 {
		t.Fatal("não concluiu chave local")
	}
	if f.ConcluirAnalise(context.Background(), domain.ChaveResumo{BatchID: "outro"}) == nil || c.conclusoes != 1 {
		t.Fatal("aceitou outra identidade")
	}
	if f.ConcluirAnalise(context.Background(), domain.ChaveResumo{BatchID: "b", ExecID: "V00000000-0000-0000-0000-000000000001"}) == nil || c.conclusoes != 1 {
		t.Fatal("inventou execução ou finalização sem API")
	}
}

func TestVarreduraExigeReciboEFalseConfirmaSomenteEsteItem(t *testing.T) {
	exec := "V00000000-0000-0000-0000-000000000001"
	for _, caso := range []string{"pendem outros", "toda finalizada", "criado", "sem flag", "outra identidade", "assincrono", "alias", "duplicado"} {
		t.Run(caso, func(t *testing.T) {
			s := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/internal/v1/execucoes/conclusoes" || r.Header.Get("Authorization") != "Bearer teste" {
					t.Error("rota/credencial")
				}
				var p pedido
				if json.NewDecoder(r.Body).Decode(&p) != nil || p.ExecID != exec || p.TenantID != "t" {
					t.Error("payload")
				}
				ack := map[string]any{"aceito": true, "tenantId": "t", "productId": "p", "batchId": "b", "execId": exec, "execucaoFinalizada": caso == "toda finalizada"}
				switch caso {
				case "criado":
					w.WriteHeader(201)
				case "assincrono":
					w.WriteHeader(202)
				case "sem flag":
					delete(ack, "execucaoFinalizada")
				case "outra identidade":
					ack["batchId"] = "outro"
				case "alias":
					delete(ack, "aceito")
					ack["Aceito"] = true
				case "duplicado":
					w.Write([]byte(`{"aceito":false,"aceito":true}`))
					return
				}
				json.NewEncoder(w).Encode(ack)
			}))
			defer s.Close()
			ref := domain.Referencia{TenantID: "t", ProductID: "p", BatchID: "b", IdempotencyID: "b"}
			v, err := NovaVarredura(s.URL, "teste", time.Second, ref, s.Client().Transport)
			if err != nil {
				t.Fatal(err)
			}
			c := &chavesTeste{}
			f := Finalizador{Chaves: c, Ref: ref, Varredura: v}
			err = f.ConcluirAnalise(context.Background(), domain.ChaveResumo{BatchID: "b", ExecID: exec})
			ok := caso == "pendem outros" || caso == "toda finalizada" || caso == "criado"
			if (err == nil) != ok || (c.conclusoes == 1) != ok {
				t.Fatalf("erro=%v conclusões=%d", err, c.conclusoes)
			}
		})
	}
}
