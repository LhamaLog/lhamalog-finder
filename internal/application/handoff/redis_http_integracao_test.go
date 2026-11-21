package handoff_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/adapter/in/queue"
	"github.com/LhamaLog/lhamalog-finder/internal/adapter/out/modocorrencia"
	"github.com/LhamaLog/lhamalog-finder/internal/adapter/out/openobserve"
	"github.com/LhamaLog/lhamalog-finder/internal/application"
	"github.com/LhamaLog/lhamalog-finder/internal/application/agrupar"
	"github.com/LhamaLog/lhamalog-finder/internal/application/handoff"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"github.com/redis/go-redis/v9"
)

type finalizadorHTTPTeste struct{ chamadas int }

func (f *finalizadorHTTPTeste) ConcluirAnalise(context.Context, domain.ChaveResumo) error {
	f.chamadas++
	return nil
}

// Exercita a confirmação no Redis real e o cliente HTTPS. O servidor de teste
// simula os resultados acordados do Módulo, sem afirmar persistência nesse serviço.
func TestHandoffHTTPComRedisMantemPendenteAteRecibo(t *testing.T) {
	if os.Getenv("FINDER_REDIS_INTEGRATION") != "1" {
		t.Skip("Redis de integração não habilitado")
	}
	ctx, cancelar := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancelar()
	sufixo := fmt.Sprintf(":handoff:%d", time.Now().UnixNano())
	stream := os.Getenv("FINDER_REDIS_STREAM") + sufixo
	grupo := os.Getenv("FINDER_REDIS_GROUP") + sufixo
	prefixo := os.Getenv("FINDER_REDIS_IDEMPOTENCIA") + sufixo + ":"
	op := queue.OpcoesLote{Addr: os.Getenv("FINDER_REDIS_ADDR"), Senha: os.Getenv("FINDER_REDIS_SENHA"), Stream: stream, Group: grupo, Consumidor: "finder-handoff", Idempotencia: prefixo, Count: 1, Block: 20 * time.Millisecond, Ociosidade: 10 * time.Millisecond}
	if op.Addr == "" {
		t.Fatal("FINDER_REDIS_ADDR ausente")
	}
	lote, err := queue.NovaLote(op)
	if err != nil {
		t.Fatal(err)
	}
	defer lote.Fechar()
	rdb := redis.NewClient(&redis.Options{Addr: op.Addr, Password: op.Senha})
	defer rdb.Close()
	t.Cleanup(func() {
		limpeza := redis.NewClient(&redis.Options{Addr: op.Addr, Password: op.Senha})
		defer limpeza.Close()
		limpeza.Del(context.Background(), stream, prefixo+"lote")
	})
	entradaRedis := map[string]any{
		"event": "log.ingested", "origin": "collector", "schemaVersion": "1", "publishedAt": "2026-11-20T13:00:00Z", "idempotencyId": "lote",
		"metadata": `{"tenantId":"t","productId":"p","correlationId":"cx"}`,
		"body":     `{"signal":"LOGS","productId":"p","recordCount":1,"openobserve":{"organization":"default","stream":"logs","batchId":"lote","minTimestampMicros":1795180000000000,"maxTimestampMicros":1795180000000000,"recordId":"r"}}`,
	}
	if err := rdb.XAdd(ctx, &redis.XAddArgs{Stream: stream, Values: entradaRedis}).Err(); err != nil {
		t.Fatal(err)
	}
	var mutex sync.Mutex
	contagem, consultas, envios, escritas := 10, 0, 0, 0
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") != "Bearer jwt-de-teste" {
			w.WriteHeader(401)
			return
		}
		switch r.URL.Path {
		case "/internal/v1/ocorrencias/consultas":
			consultas++
			fmt.Fprintf(w, `{"estado":{"contagemAcumulada":%d,"primeiraOcorrencia":"2026-11-18T13:00:00Z","ultimaOcorrencia":"2026-11-19T13:00:00Z","status":"ABERTA","revisao":"rev-%d"}}`, contagem, contagem)
		case "/internal/v1/resumos-analise":
			envios++
			var v struct {
				BatchID string `json:"batchId"`
				Achados []struct {
					Contagem int      `json:"contagem"`
					Revisao  string   `json:"revisaoConsultada"`
					Linhas   []string `json:"linhas"`
				} `json:"achados"`
			}
			if err := json.NewDecoder(r.Body).Decode(&v); err != nil || v.BatchID != "lote" || len(v.Achados) != 1 || strings.Join(v.Achados[0].Linhas, "") != "mensagem segura" {
				t.Errorf("payload de entrega inválido: %+v, %v", v, err)
				w.WriteHeader(422)
				return
			}
			a := v.Achados[0]
			if a.Contagem != contagem+1 || a.Revisao != fmt.Sprintf("rev-%d", contagem) {
				t.Errorf("revisão/contagem não propagadas: %+v", a)
				w.WriteHeader(422)
				return
			}
			switch envios {
			case 1:
				contagem = 12 // outro lote foi persistido entre consulta e escrita
				w.WriteHeader(409)
				fmt.Fprint(w, `{"codigo":"REVISAO_DESATUALIZADA"}`)
			case 2:
				contagem = a.Contagem
				escritas++
				w.WriteHeader(503) // o recibo não chega ao Finder
			default:
				fmt.Fprint(w, `{"aceito":true,"reciboId":"recibo-1","tenantId":"t","productId":"p","batchId":"lote","repeticao":true}`)
			}
		default:
			w.WriteHeader(404)
		}
	}))
	defer srv.Close()
	cliente, err := modocorrencia.NovoCliente(modocorrencia.OpcoesCliente{URLBase: srv.URL, JWT: "jwt-de-teste", Transporte: srv.Client().Transport})
	if err != nil {
		t.Fatal(err)
	}
	mascaramentos := 0
	seguranca := application.ValidacaoSegurancaFunc(func(_ context.Context, _ []string) domain.ResultadoValidarSeguranca {
		mascaramentos++
		return domain.ResultadoValidarSeguranca{LinhasMascaradas: []domain.LinhaLogMascarada{{Texto: "mensagem segura"}}}
	})
	eventos := []string{}
	preparar := agrupar.Orquestrador{Montador: application.ServicoAnalise{Seguranca: seguranca}, Consulta: cliente, Recorrencia: openobserve.Recorrencia{Historico: historicoFluxo{eventos: &eventos}}}
	finalizador := &finalizadorHTTPTeste{}
	efeito := handoff.EfeitoAnalise{Recuperacao: recuperacaoEfeito{domain.ResultadoRecuperacao{Classificacao: domain.ClasseLoteCompleto, EfeitoConcluido: true, Registros: []domain.DocumentoFisico{{"lhamalog_record_id": "r", "_timestamp": json.Number("1795180000000000"), "severitytext": "WARN", "body": "bruto", "servicename": "s", "resourceattributes": `{"host.name":"h"}`}}}},
		Mapear: openobserve.MapearDocumentoAnalise, VersaoRegras: "v1", VersaoCatalogo: "catalogo-v1", Entrega: handoff.OrquestradorEntrega{Preparador: preparar, Entrega: handoff.Entrega{Envio: cliente, Finalizar: finalizador}},
	}
	consumo := application.ConsumoReferencia{Fila: lote, Chaves: lote, Efeito: efeito}
	r, err := consumo.Executar(ctx)
	if !errors.Is(err, out.ErrModuloTemporario) || r.Confirmacoes != 0 || finalizador.chamadas != 0 || mascaramentos != 1 {
		t.Fatalf("confirmou sem recibo: %+v %v", r, err)
	}
	if p, err := rdb.XPending(ctx, stream, grupo).Result(); err != nil || p.Count != 1 {
		t.Fatalf("esperava uma entrada pendente: %v %v", p, err)
	}
	for r.Confirmacoes == 0 && ctx.Err() == nil {
		r, err = consumo.Executar(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}
	if r.Confirmacoes != 1 || finalizador.chamadas != 1 || mascaramentos != 2 {
		t.Fatalf("recibo não encerrou a entrada: %+v", r)
	}
	if p, err := rdb.XPending(ctx, stream, grupo).Result(); err != nil || p.Count != 0 {
		t.Fatalf("entrada ainda pendente: %v %v", p, err)
	}
	mutex.Lock()
	defer mutex.Unlock()
	if contagem != 13 || consultas != 3 || envios != 3 || escritas != 1 {
		t.Fatalf("contagem=%d consultas=%d envios=%d escritas=%d", contagem, consultas, envios, escritas)
	}
}
