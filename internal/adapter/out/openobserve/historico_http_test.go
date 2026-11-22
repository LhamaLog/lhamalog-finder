package openobserve

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

func configHistorico(base string) domain.ConfiguracaoAcesso {
	c := []string{"copia"}
	return domain.ConfiguracaoAcesso{URL: base, Usuario: "teste", Senha: "teste", Timeout: time.Second, TamanhoPagina: 2, CamposCopia: &c, Ordem: "lhamalog_record_id", CampoTenant: "tenantid", CampoProduto: "productid"}
}
func referenciaHistorico() domain.Referencia {
	return domain.Referencia{TenantID: "t", ProductID: "p", Organizacao: "default", Stream: "logs", BatchID: "atual", RecordCount: 1, MaxTimestampMicros: 100}
}
func pedidoHistorico() out.PedidoRecorrencia {
	return out.PedidoRecorrencia{BatchID: "atual", Host: "h", Chave: domain.IdentidadeAchado{Assinatura: "segura", Classe: domain.ERROR, TenantID: "t", ProductID: "p"}}
}
func documentoHistorico(batch, id, copia string) domain.DocumentoFisico {
	return domain.DocumentoFisico{"lhamalog_batch_id": batch, "lhamalog_record_id": id, "tenantid": "t", "productid": "p", "_timestamp": int64(90), "body": "não derivar assinatura sem scanner", "copia": copia}
}
func montadorHistorico(ctx context.Context, ref domain.Referencia, docs []domain.DocumentoFisico) (domain.ResumoAnalise, error) {
	return domain.ResumoAnalise{TenantID: ref.TenantID, ProductID: ref.ProductID, BatchID: ref.BatchID, Achados: []domain.AchadoAnalise{{Assinatura: "segura", HostName: "h", Tipo: domain.ERROR}}}, nil
}

func TestHistoricoPaginasDuplicatasSnapshotEIsolamento(t *testing.T) {
	docs := []domain.DocumentoFisico{documentoHistorico("antigo", "r", "1"), documentoHistorico("antigo", "r", "2"), documentoHistorico("atual", "r", "1"), documentoHistorico("outro", "r", "1")}
	calls := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		var q corpoConsulta
		if json.NewDecoder(r.Body).Decode(&q) != nil {
			t.Error("query inválida")
		}
		if q.Query.StartTime != 1 || q.Query.EndTime != 101 || !strings.Contains(q.Query.SQL, `"tenantid" = 't' AND "productid" = 'p'`) {
			t.Error("filtro/limite ausentes")
		}
		var hits any
		if strings.Contains(q.Query.SQL, "GROUP BY") {
			if !strings.Contains(q.Query.SQL, "LIMIT 499") || q.Query.Size != 499 {
				t.Error("janela errada")
			}
			hits = []any{map[string]any{"batch_id": "outro", "record_count": 1, "last_timestamp": 95}, map[string]any{"batch_id": "antigo", "record_count": 1, "last_timestamp": 90}}
		} else {
			if !strings.Contains(q.Query.SQL, `ORDER BY "lhamalog_batch_id" ASC, "lhamalog_record_id" ASC, "_timestamp" DESC`) {
				t.Error("ordem instável")
			}
			a := q.Query.From
			b := a + q.Query.Size
			if b > len(docs) {
				b = len(docs)
			}
			hits = docs[a:b]
		}
		json.NewEncoder(w).Encode(map[string]any{"hits": hits})
	}))
	defer s.Close()
	analisados := 0
	h, err := NovoHistorico(configHistorico(s.URL), referenciaHistorico(), func(c context.Context, r domain.Referencia, d []domain.DocumentoFisico) (domain.ResumoAnalise, error) {
		analisados++
		if len(d) != 1 {
			t.Error("duplicata chegou ao scanner")
		}
		return montadorHistorico(c, r, d)
	})
	if err != nil {
		t.Fatal(err)
	}
	svc := Recorrencia{Historico: h}
	for i := 0; i < 2; i++ {
		m, e := svc.ConsultarRecorrencia(context.Background(), pedidoHistorico())
		if e != nil || m.Valor == nil || !*m.Valor {
			t.Fatalf("marca=%v erro=%v", m, e)
		}
	}
	if calls != 4 || analisados != 3 {
		t.Fatalf("snapshot não reutilizado: http=%d análise=%d", calls, analisados)
	}
	ev, _ := h.ConsultarHistorico(context.Background(), pedidoHistorico())
	if ev[0].BatchID != "atual" || ev[2].BatchID != "outro" || ev[4].BatchID != "antigo" {
		t.Fatal("ordem da janela")
	}
	p := pedidoHistorico()
	p.Chave.TenantID = "intruso"
	if _, e := h.ConsultarHistorico(context.Background(), p); e == nil || calls != 4 {
		t.Fatal("permitiu consultar outra identidade")
	}
}

func TestHistoricoRejeitaDadosInvalidosENaoConfirmaParcial(t *testing.T) {
	for _, caso := range []string{"isolamento", "sem id", "divergencia", "mudou contagem", "parcial", "corpo invalido", "autenticacao", "query invalida", "temporario"} {
		t.Run(caso, func(t *testing.T) {
			calls := 0
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				var q corpoConsulta
				json.NewDecoder(r.Body).Decode(&q)
				switch caso {
				case "autenticacao":
					w.WriteHeader(401)
					return
				case "query invalida":
					w.WriteHeader(400)
					return
				case "temporario":
					w.WriteHeader(503)
					return
				case "corpo invalido":
					fmt.Fprint(w, `{"hits":null}`)
					return
				case "parcial":
					fmt.Fprint(w, `{"hits":[],"is_partial":true}`)
					return
				}
				if strings.Contains(q.Query.SQL, "GROUP BY") {
					fmt.Fprint(w, `{"hits":[]}`)
					return
				}
				d := documentoHistorico("atual", "r", "1")
				docs := []domain.DocumentoFisico{d}
				switch caso {
				case "isolamento":
					d["tenantid"] = "intruso"
				case "sem id":
					delete(d, "lhamalog_record_id")
				case "divergencia":
					e := documentoHistorico("atual", "r", "2")
					e["body"] = "diferente"
					docs = append(docs, e)
				case "mudou contagem":
					docs = []domain.DocumentoFisico{}
				}
				if q.Query.From > 0 {
					docs = nil
				}
				json.NewEncoder(w).Encode(map[string]any{"hits": docs})
			}))
			defer s.Close()
			h, _ := NovoHistorico(configHistorico(s.URL), referenciaHistorico(), montadorHistorico)
			a := &avisoTeste{}
			marca, err := (Recorrencia{Historico: h, Aviso: a}).ConsultarRecorrencia(context.Background(), pedidoHistorico())
			if err == nil || marca.Valor != nil {
				t.Fatal("falha produziu recorrente=false")
			}
			temporaria := caso == "temporario" || caso == "parcial" || caso == "mudou contagem"
			if temporaria && a.chamadas != 1 {
				t.Fatal("não avisou após três tentativas")
			}
			if !temporaria && a.chamadas != 0 {
				t.Fatal("repetiu falha crítica")
			}
			if (caso == "autenticacao" || caso == "query invalida" || caso == "corpo invalido") && calls != 1 {
				t.Fatalf("retry indevido: %d", calls)
			}
		})
	}
}

// Consulta a um OpenObserve real de teste: comprova os aliases, GROUP BY,
// DISTINCT, limite, corte temporal e paginação da query definida no Drive.
func TestHistoricoOpenObserveReal(t *testing.T) {
	if os.Getenv("FINDER_OPENOBSERVE_INTEGRATION") == "" {
		t.Skip("OpenObserve de teste opt-in")
	}
	cfg := configHistorico(os.Getenv("FINDER_OO_URL"))
	cfg.Usuario = os.Getenv("FINDER_OO_USUARIO")
	cfg.Senha = os.Getenv("FINDER_OO_SENHA")
	cfg.Timeout = 5 * time.Second
	ref := referenciaHistorico()
	ref.Stream = fmt.Sprintf("historico_%d", time.Now().UnixNano())
	ref.MaxTimestampMicros = time.Now().UnixMicro()
	ref.RecordCount = 1
	var docs []domain.DocumentoFisico
	for i := 0; i < 502; i++ {
		id := fmt.Sprintf("b%03d", i)
		if i == 0 {
			id = "atual"
		}
		d := documentoHistorico(id, "r", "1")
		d["_timestamp"] = ref.MaxTimestampMicros - int64(i)
		docs = append(docs, d)
	}
	// Outro tenant e evento posterior jamais podem ocupar uma posição na janela.
	intruso := documentoHistorico("intruso", "r", "1")
	intruso["tenantid"] = "outro"
	intruso["_timestamp"] = ref.MaxTimestampMicros
	docs = append(docs, intruso)
	futuro := documentoHistorico("futuro", "r", "1")
	futuro["_timestamp"] = ref.MaxTimestampMicros + 1000000
	docs = append(docs, futuro)
	data, _ := json.Marshal(docs)
	req, _ := http.NewRequest("POST", cfg.URL+"/api/default/"+ref.Stream+"/_json", strings.NewReader(string(data)))
	req.SetBasicAuth(cfg.Usuario, cfg.Senha)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("ingestão indisponível")
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("ingestão status=%d", resp.StatusCode)
	}
	cfg.TamanhoPagina = 73
	var ev []domain.EvidenciaRecorrencia
	deadline := time.Now().Add(20 * time.Second)
	for {
		h, e := NovoHistorico(cfg, ref, montadorHistorico)
		if e != nil {
			t.Fatal(e)
		}
		ev, err = h.ConsultarHistorico(context.Background(), pedidoHistorico())
		if err == nil && len(ev) == 1000 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("janela: evidências=%d erro=%v", len(ev), err)
		}
		time.Sleep(200 * time.Millisecond)
	}
	ids := map[string]bool{}
	for _, e := range ev {
		ids[e.BatchID] = true
	}
	if len(ids) != 500 || !ids["atual"] || !ids["b499"] || ids["b500"] || ids["intruso"] || ids["futuro"] {
		t.Fatalf("janela incorreta: %d", len(ids))
	}
	sorted := []string{}
	for id := range ids {
		sorted = append(sorted, id)
	}
	sort.Strings(sorted)
	t.Logf("OpenObserve real: %d lotes distintos; paginação e isolamento confirmados", len(sorted))
}
