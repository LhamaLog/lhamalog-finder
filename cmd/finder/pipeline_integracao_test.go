package main

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// Executa o binário real com Redis/OpenObserve reais de teste. Somente o Módulo
// e o Core são servidores de contrato controlados; não é homologação externa.
func TestBinarioExecutaPipelineAntesDoACK(t *testing.T) {
	if runtime.GOOS != "linux" || os.Getenv("FINDER_OPENOBSERVE_INTEGRATION") == "" || os.Getenv("FINDER_REDIS_INTEGRATION") != "1" {
		t.Skip("binário Linux com Redis/OpenObserve opt-in")
	}
	binary := filepath.Join(t.TempDir(), "finder")
	build := exec.Command("go", "build", "-o", binary, ".")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build do binário: %v %s", err, output)
	}
	for _, modo := range []string{"sucesso", "modulo temporario", "historico indisponivel", "recibo invalido"} {
		t.Run(modo, func(t *testing.T) { executarBinario(t, binary, modo) })
	}
}

func executarBinario(t *testing.T, binary, modo string) {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	tag := fmt.Sprintf("%d", time.Now().UnixNano())
	streamOO := "pipeline_" + tag
	batch := "batch_" + tag
	tenant := "tenant_" + tag
	produto := "produto_pipeline"
	grupo := "pipeline_" + tag
	prefixo := "pipeline_" + tag + ":"
	rdb := redis.NewClient(&redis.Options{Addr: os.Getenv("FINDER_REDIS_ADDR"), Password: os.Getenv("FINDER_REDIS_PASSWORD")})
	defer rdb.Close()
	stream := "lhamalog:stream:logs"
	if err := rdb.XGroupCreateMkStream(ctx, stream, grupo, "$").Err(); err != nil {
		t.Fatal(err)
	}
	defer rdb.XGroupDestroy(context.Background(), stream, grupo)
	defer rdb.Del(context.Background(), prefixo+batch)
	var entradaIDs []string
	defer func() {
		if len(entradaIDs) > 0 {
			rdb.XDel(context.Background(), stream, entradaIDs...)
		}
	}()
	var aceitar atomic.Bool
	aceitar.Store(modo == "sucesso")
	var mutex sync.Mutex
	consultas, envios, alertas := 0, 0, 0
	modulo := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()
		if r.Header.Get("Authorization") != "Bearer teste-pipeline" {
			t.Error("JWT ausente")
			w.WriteHeader(401)
			return
		}
		data, _ := io.ReadAll(r.Body)
		if strings.Contains(string(data), "example.invalid") || strings.Contains(string(data), "valor-sigiloso") {
			t.Error("dados brutos escaparam do scanner")
			w.WriteHeader(422)
			return
		}
		switch r.URL.Path {
		case "/internal/v1/ocorrencias/consultas":
			consultas++
			w.Write([]byte(`{"estado":null}`))
		case "/internal/v1/resumos-analise":
			envios++
			var resumo struct {
				Batch      string `json:"batchId"`
				Exec       string `json:"execId"`
				Contadores struct {
					Linhas     int `json:"linhasProcessadas"`
					Mascarados int `json:"dadosSensiveisMascarados"`
				} `json:"contadores"`
				Achados []struct {
					Recorrente bool     `json:"recorrente"`
					Contagem   int      `json:"contagem"`
					Linhas     []string `json:"linhas"`
				} `json:"achados"`
			}
			if json.Unmarshal(data, &resumo) != nil || resumo.Batch != batch || resumo.Exec != "" || resumo.Contadores.Linhas != 2 || resumo.Contadores.Mascarados != 4 || len(resumo.Achados) != 1 || resumo.Achados[0].Contagem != 2 || !resumo.Achados[0].Recorrente || strings.Join(resumo.Achados[0].Linhas, ",") != "log com senha omitido,log com senha omitido" {
				t.Error("resumo/recorrência/catálogo incorretos")
				w.WriteHeader(422)
				return
			}
			if modo == "recibo invalido" {
				w.Write([]byte(`{"aceito":true}`))
				return
			}
			if !aceitar.Load() {
				w.WriteHeader(503)
				return
			}
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]any{"aceito": true, "reciboId": "recibo-pipeline", "tenantId": tenant, "productId": produto, "batchId": batch, "repeticao": false})
		case "/internal/v1/alertas/finder":
			alertas++
			var p map[string]any
			if json.Unmarshal(data, &p) != nil || len(p) != 4 || p["tenantId"] != tenant || p["productId"] != produto || p["batchId"] != batch || p["motivo"] != "consulta_openobserve_indisponivel" {
				t.Error("aviso Core incorreto")
				w.WriteHeader(422)
				return
			}
			p["aceito"] = true
			json.NewEncoder(w).Encode(p)
		default:
			t.Error("rota externa inesperada")
			w.WriteHeader(404)
		}
	}))
	defer modulo.Close()
	cert := filepath.Join(t.TempDir(), "root.pem")
	if err := os.WriteFile(cert, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: modulo.Certificate().Raw}), 0600); err != nil {
		t.Fatal(err)
	}
	// O proxy preserva a consulta real e permite interromper somente o histórico.
	realOO, err := url.Parse(os.Getenv("FINDER_OO_URL"))
	if err != nil {
		t.Fatal(err)
	}
	reverse := httputil.NewSingleHostReverseProxy(realOO)
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(raw))
		if modo == "historico indisponivel" && !aceitar.Load() && bytes.Contains(raw, []byte("GROUP BY")) {
			w.WriteHeader(503)
			return
		}
		reverse.ServeHTTP(w, r)
	}))
	defer proxy.Close()
	agora := time.Now().UnixMicro()
	var docs []map[string]any
	for i, id := range []string{batch + "_anterior_a", batch + "_anterior_b", batch, batch} {
		docs = append(docs, map[string]any{"lhamalog_batch_id": id, "lhamalog_record_id": fmt.Sprintf("r%d", i), "_timestamp": agora + int64(i), "tenantid": tenant, "productid": produto, "severitytext": "ERROR", "body": "a+teste@example.invalid senha=valor-sigiloso", "servicename": "servico-pipeline", "resourceattributes": `{"host.name":"host-pipeline"}`})
	}
	payload, _ := json.Marshal(docs)
	req, _ := http.NewRequestWithContext(ctx, "POST", realOO.String()+"/api/default/"+streamOO+"/_json", bytes.NewReader(payload))
	req.SetBasicAuth(os.Getenv("FINDER_OO_USUARIO"), os.Getenv("FINDER_OO_SENHA"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		t.Fatal("OpenObserve de teste indisponível")
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("ingestão status=%d", resp.StatusCode)
	}
	// Aguardamos a indexação antes de publicar a referência ao processo.
	sql := fmt.Sprintf(`SELECT * FROM "%s"`, streamOO)
	data, _ := json.Marshal(map[string]any{"query": map[string]any{"sql": sql, "start_time": agora, "end_time": agora + 1000, "from": 0, "size": 10}})
	aguardarPipeline(t, ctx, func() bool {
		q, _ := http.NewRequestWithContext(ctx, "POST", realOO.String()+"/api/default/_search", bytes.NewReader(data))
		q.SetBasicAuth(os.Getenv("FINDER_OO_USUARIO"), os.Getenv("FINDER_OO_SENHA"))
		q.Header.Set("Content-Type", "application/json")
		r, e := http.DefaultClient.Do(q)
		if e != nil {
			return false
		}
		defer r.Body.Close()
		var v struct {
			Hits []any `json:"hits"`
		}
		return json.NewDecoder(r.Body).Decode(&v) == nil && len(v.Hits) == 4
	}, "indexação dos quatro documentos")
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	ln.Close()
	values := map[string]string{"FINDER_HTTP_ADDR": addr, "FINDER_REDIS_ADDR": os.Getenv("FINDER_REDIS_ADDR"), "FINDER_REDIS_PASSWORD": os.Getenv("FINDER_REDIS_PASSWORD"), "FINDER_REDIS_STREAM": stream, "FINDER_REDIS_GROUP": grupo, "FINDER_REDIS_CONSUMIDOR": "binario-pipeline", "FINDER_REDIS_IDEMPOTENCIA": prefixo, "FINDER_REDIS_COUNT": "1", "FINDER_REDIS_BLOCK": "100ms", "FINDER_REDIS_PENDENCIA_OCIOSA": "100ms", "FINDER_OO_URL": proxy.URL, "FINDER_OO_USUARIO": os.Getenv("FINDER_OO_USUARIO"), "FINDER_OO_SENHA": os.Getenv("FINDER_OO_SENHA"), "FINDER_OO_TIMEOUT": "3s", "FINDER_OO_PAGINA": "2", "FINDER_OO_CAMPOS_COPIA": "", "FINDER_OO_ORDEM": "lhamalog_record_id", "FINDER_OO_CAMPO_TENANT": "tenantid", "FINDER_OO_CAMPO_PRODUTO": "productid", "FINDER_CATALOGO_SEGURANCA": "../../config/seguranca/catalogo.json", "FINDER_VERSAO_REGRAS": "finder-2026-11-22.1", "FINDER_MODULO_OCORRENCIA_URL": modulo.URL, "FINDER_MODULO_OCORRENCIA_JWT": "teste-pipeline", "FINDER_CORE_URL": modulo.URL, "FINDER_CORE_JWT": "teste-pipeline", "SSL_CERT_FILE": cert}
	cmd := exec.CommandContext(ctx, binary)
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "FINDER_") && !strings.HasPrefix(e, "SSL_CERT_FILE=") {
			cmd.Env = append(cmd.Env, e)
		}
	}
	for k, v := range values {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	defer func() {
		cancel()
		<-done
		if bytes.Contains(output.Bytes(), []byte("valor-sigiloso")) || bytes.Contains(output.Bytes(), []byte("example.invalid")) {
			t.Error("dados brutos no log do processo")
		}
		if t.Failed() {
			t.Logf("diagnóstico do binário: %s", output.String())
		}
	}()
	aguardarPipeline(t, ctx, func() bool {
		r, e := (&http.Client{Timeout: time.Second}).Get("http://" + addr + "/health")
		if e != nil {
			return false
		}
		defer r.Body.Close()
		return r.StatusCode == 200
	}, "prontidão do binário")
	body, _ := json.Marshal(map[string]any{"signal": "LOGS", "productId": produto, "recordCount": 2, "openobserve": map[string]any{"organization": "default", "stream": streamOO, "batchId": batch, "minTimestampMicros": agora, "maxTimestampMicros": agora + 1000}})
	meta, _ := json.Marshal(map[string]string{"tenantId": tenant, "productId": produto, "correlationId": "correlacao_" + tag})
	entrada := map[string]any{"event": "log.ingested", "origin": "collector", "schemaVersion": "1", "publishedAt": time.Now().UTC().Format(time.RFC3339Nano), "idempotencyId": batch, "metadata": string(meta), "body": string(body)}
	id, err := rdb.XAdd(ctx, &redis.XAddArgs{Stream: stream, Values: entrada}).Result()
	if err != nil {
		t.Fatal(err)
	}
	entradaIDs = append(entradaIDs, id)
	if modo != "sucesso" {
		aguardarPipeline(t, ctx, func() bool {
			mutex.Lock()
			defer mutex.Unlock()
			if modo == "historico indisponivel" {
				return alertas > 0
			}
			return envios > 0
		}, "falha controlada observada")
		pendentes, e := rdb.XPending(ctx, stream, grupo).Result()
		estado, _ := rdb.Get(ctx, prefixo+batch).Result()
		if e != nil || pendentes.Count != 1 || estado == "concluido" {
			t.Fatal("confirmou entrada antes do recibo válido")
		}
		if modo == "recibo invalido" {
			t.Log("recibo inválido conservou pendência; nenhum ACK")
			return
		}
		aceitar.Store(true)
	}
	aguardarPipeline(t, ctx, func() bool {
		p, e := rdb.XPending(ctx, stream, grupo).Result()
		estado, _ := rdb.Get(ctx, prefixo+batch).Result()
		return e == nil && p.Count == 0 && estado == "concluido"
	}, "recibo, marca de conclusão e ACK")
	ttl, _ := rdb.TTL(ctx, prefixo+batch).Result()
	if ttl < 24*time.Hour || ttl > 25*time.Hour {
		t.Fatalf("retenção incorreta: %v", ttl)
	}
	mutex.Lock()
	enviosAntes := envios
	consultasAntes := consultas
	mutex.Unlock()
	replay, err := rdb.XAdd(ctx, &redis.XAddArgs{Stream: stream, Values: entrada}).Result()
	if err != nil {
		t.Fatal(err)
	}
	entradaIDs = append(entradaIDs, replay)
	aguardarPipeline(t, ctx, func() bool {
		groups, e := rdb.XInfoGroups(ctx, stream).Result()
		if e != nil {
			return false
		}
		for _, g := range groups {
			if g.Name == grupo && g.LastDeliveredID == replay && g.Pending == 0 {
				return true
			}
		}
		return false
	}, "reentrega idempotente")
	mutex.Lock()
	defer mutex.Unlock()
	if envios != enviosAntes || consultas != consultasAntes {
		t.Fatal("reentrega executou novamente o efeito")
	}
	t.Logf("binário real: modo=%s; resumo mascarado; recorrência de três lotes; ACK após recibo; TTL/idempotência confirmados", modo)
}

func aguardarPipeline(t *testing.T, ctx context.Context, check func() bool, etapa string) {
	t.Helper()
	for !check() {
		select {
		case <-ctx.Done():
			t.Fatalf("prazo aguardando %s", etapa)
		case <-time.After(40 * time.Millisecond):
		}
	}
}
