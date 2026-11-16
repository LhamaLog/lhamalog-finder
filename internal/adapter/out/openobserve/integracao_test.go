package openobserve

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/application"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// Prova com OpenObserve de teste e dados sintéticos.
// Não usa log de produção e não declara homologação.
func TestOpenObserveIntegracao(t *testing.T) {
	if _, definida := os.LookupEnv("FINDER_OPENOBSERVE_INTEGRATION"); !definida {
		t.Skip("FINDER_OPENOBSERVE_INTEGRATION não está definida")
	}
	prova := carregarProva(t)
	if err := responde(prova); err != nil {
		t.Fatalf("FINDER_OO_URL não responde: %v", err)
	}

	t.Run("lote sintetico", func(t *testing.T) {
		org := "default"
		stream := "finder_sintetico_lote"
		lote := "lote-" + unico()
		minimo := time.Now().UnixMicro()
		maximo := minimo + int64(time.Second/time.Microsecond)
		ingerir(t, prova, org, stream, []map[string]any{
			docSintetico(prova, lote, "r1", "tenant-a", "produto-a", "1", minimo+10),
		})
		id := "r1"
		ref := referenciaSintetica(org, stream, lote, minimo, maximo, 1)
		ref.RecordID = &id
		visto := &coletor{}
		resultado := recuperarVia(t, prova, ref, visto, classificarPronto(domain.ClasseLoteCompleto))
		if resultado.Classificacao != domain.ClasseLoteCompleto || !resultado.EfeitoConcluido || len(resultado.Registros) != 1 {
			t.Fatalf("classe=%s efeito=%v n=%d", resultado.Classificacao, resultado.EfeitoConcluido, len(resultado.Registros))
		}
		if len(visto.pedidos) == 0 {
			t.Fatal("consulta não foi observada")
		}
		pedido := visto.pedidos[0]
		if !strings.Contains(pedido.caminho, "/api/"+url.PathEscape(org)+"/_search") {
			t.Fatalf("caminho=%s", pedido.caminho)
		}
		if pedido.inicio != minimo || pedido.fim != maximo+1 {
			t.Fatalf("faixa %d..%d", pedido.inicio, pedido.fim)
		}
		if !strings.Contains(pedido.sql, `"`+stream+`"`) || strings.Contains(pedido.sql, "lhamalog:stream:logs") {
			t.Fatal("stream da consulta não é o da referência")
		}
		if !strings.Contains(pedido.sql, lote) || !strings.Contains(pedido.sql, "lhamalog_record_id") {
			t.Fatal("lote unitário não filtrou lote e recordId")
		}
		if _, confirma := any(&application.ServicoRecuperacao{}).(interface {
			Confirmar(context.Context, string) error
		}); confirma {
			t.Fatal("recuperação confirma a entrada")
		}
	})

	t.Run("lote de varios nao filtra recordId", func(t *testing.T) {
		org := "default"
		stream := "finder_sintetico_varios"
		lote := "lote-" + unico()
		minimo := time.Now().UnixMicro()
		maximo := minimo + int64(time.Second/time.Microsecond)
		ingerir(t, prova, org, stream, []map[string]any{
			docSintetico(prova, lote, "r1", "tenant-a", "produto-a", "1", minimo+10),
			docSintetico(prova, lote, "r2", "tenant-a", "produto-a", "1", minimo+11),
		})
		ref := referenciaSintetica(org, stream, lote, minimo, maximo, 2)
		visto := &coletor{}
		resultado := recuperarVia(t, prova, ref, visto, classificarPronto(domain.ClasseLoteCompleto))
		if resultado.Classificacao != domain.ClasseLoteCompleto || resultado.Encontrado != 2 {
			t.Fatalf("classe=%s encontrado=%d", resultado.Classificacao, resultado.Encontrado)
		}
		if strings.Contains(visto.pedidos[0].sql, "lhamalog_record_id") {
			t.Fatal("lote de vários filtrou recordId")
		}
	})

	t.Run("copias e identificador alem de recordCount", func(t *testing.T) {
		if prova.Pagina < 2 {
			t.Fatal("FINDER_OO_PAGINA não permite observar a página seguinte")
		}
		org := "default"
		stream := "finder_sintetico_copias"
		lote := "lote-" + unico()
		minimo := time.Now().UnixMicro()
		maximo := minimo + int64(time.Second/time.Microsecond)
		docs := []map[string]any{
			docSintetico(prova, lote, "r1", "tenant-a", "produto-a", "c1", minimo+10),
			docSintetico(prova, lote, "r1", "tenant-a", "produto-a", "c2", minimo+10),
		}
		for i := 0; i < prova.Pagina; i++ {
			docs = append(docs, docSintetico(prova, lote, "r2", "tenant-a", "produto-a", "c1", minimo+20))
		}
		ingerir(t, prova, org, stream, docs)
		ref := referenciaSintetica(org, stream, lote, minimo, maximo, 2)
		resultado := recuperarVia(t, prova, ref, &coletor{}, classificarPronto(domain.ClasseLoteCompleto))
		ids := map[string]int{}
		for _, doc := range resultado.Registros {
			id, ok := domain.IdentificadorRegistro(doc)
			if ok {
				ids[id]++
			}
		}
		if resultado.Classificacao != domain.ClasseLoteCompleto || ids["r1"] != 1 || ids["r2"] != 1 {
			t.Fatalf("classe=%s ids=%v", resultado.Classificacao, ids)
		}
	})

	t.Run("vazio menor maior e misto", func(t *testing.T) {
		org := "default"
		stream := "finder_sintetico_contagem"
		minimo := time.Now().UnixMicro()
		maximo := minimo + int64(time.Second/time.Microsecond)
		ingerir(t, prova, org, stream, []map[string]any{
			docSintetico(prova, "lote-outro-"+unico(), "rx", "tenant-a", "produto-a", "1", minimo+10),
		})
		vazio := referenciaSintetica(org, stream, "lote-vazio-"+unico(), minimo, maximo, 1)
		resultado := recuperarVia(t, prova, vazio, &coletor{}, func(domain.ResultadoRecuperacao) bool { return true })
		if resultado.Classificacao != domain.ClasseIncompleto || resultado.EfeitoConcluido || resultado.Encontrado != 0 {
			t.Fatalf("vazio: %s encontrado=%d efeito=%v", resultado.Classificacao, resultado.Encontrado, resultado.EfeitoConcluido)
		}

		lote := "lote-" + unico()
		ingerir(t, prova, org, stream, []map[string]any{
			docSintetico(prova, lote, "r1", "tenant-a", "produto-a", "1", minimo+10),
		})
		menor := referenciaSintetica(org, stream, lote, minimo, maximo, 2)
		resultado = recuperarVia(t, prova, menor, &coletor{}, func(r domain.ResultadoRecuperacao) bool {
			return r.Encontrado == 1 || r.Classificacao != domain.ClasseIncompleto
		})
		if resultado.Classificacao != domain.ClasseIncompleto || resultado.Encontrado != 1 || resultado.RecordCount != 2 {
			t.Fatalf("menor: %s encontrado=%d", resultado.Classificacao, resultado.Encontrado)
		}
		loteMaior := "lote-maior-" + unico()
		ingerir(t, prova, org, stream, []map[string]any{
			docSintetico(prova, loteMaior, "r1", "tenant-a", "produto-a", "1", minimo+10),
			docSintetico(prova, loteMaior, "r2", "tenant-a", "produto-a", "1", minimo+11),
		})
		excesso := referenciaSintetica(org, stream, loteMaior, minimo, maximo, 1)
		resultado = recuperarVia(t, prova, excesso, &coletor{}, classificarPronto(domain.ClasseExcesso))
		if resultado.Classificacao != domain.ClasseExcesso || !resultado.EfeitoConcluido || resultado.Encontrado != 2 || resultado.Registros != nil {
			t.Fatalf("maior: %s encontrado=%d", resultado.Classificacao, resultado.Encontrado)
		}
		misto := referenciaSintetica(org, stream, lote, minimo, maximo, 1)
		resultado = recuperarVia(t, prova, misto, &coletor{}, classificarPronto(domain.ClasseLoteCompleto))
		if resultado.Classificacao != domain.ClasseLoteCompleto || resultado.Encontrado != 1 {
			t.Fatalf("misto: %s encontrado=%d", resultado.Classificacao, resultado.Encontrado)
		}
	})

	t.Run("dois tenants", func(t *testing.T) {
		org := "default"
		stream := "finder_sintetico_tenant"
		lote := "lote-" + unico()
		minimo := time.Now().UnixMicro()
		maximo := minimo + int64(time.Second/time.Microsecond)
		ingerir(t, prova, org, stream, []map[string]any{
			docSintetico(prova, lote, "r1", "tenant-a", "produto-a", "1", minimo+10),
			docSintetico(prova, lote, "r2", "tenant-b", "produto-a", "1", minimo+11),
		})
		ref := referenciaSintetica(org, stream, lote, minimo, maximo, 1)
		visto := &coletor{}
		resultado := recuperarVia(t, prova, ref, visto, func(r domain.ResultadoRecuperacao) bool {
			return r.Classificacao == domain.ClasseLoteCompleto || r.Classificacao == domain.ClasseIsolamento
		})
		if len(visto.pedidos) == 0 || !strings.Contains(visto.pedidos[0].sql, prova.CampoTenant) || !strings.Contains(visto.pedidos[0].sql, "tenant-a") {
			t.Fatal("a consulta não filtrou o tenant da referência")
		}
		for _, doc := range resultado.Registros {
			if !mesmoCampo(doc, prova.CampoTenant, "tenant-a") || !mesmoCampo(doc, prova.CampoProduto, "produto-a") {
				t.Fatal("documento de outro tenant ou produto saiu como resultado válido")
			}
		}
		if resultado.Classificacao != domain.ClasseLoteCompleto && resultado.Classificacao != domain.ClasseIsolamento {
			t.Fatalf("classe=%s", resultado.Classificacao)
		}
		if resultado.Classificacao == domain.ClasseIsolamento && (resultado.Registros != nil || !resultado.EfeitoConcluido) {
			t.Fatal("isolamento entregou documento ou não concluiu o efeito")
		}
	})

	t.Run("segunda pagina inalcancavel", func(t *testing.T) {
		quebra := &proxyQuebra{pagina: prova.Pagina, tenant: prova.CampoTenant, produto: prova.CampoProduto}
		srv := httptest.NewServer(quebra)
		defer srv.Close()
		acesso := prova.acesso
		acesso.URL = srv.URL
		ref := referenciaSintetica("default", "stream_sintetico", "lote-quebra", 100, 200, int64(prova.Pagina))
		resultado, err := (&application.ServicoRecuperacao{Busca: Cliente{}, Acesso: acesso}).Recuperar(context.Background(), ref)
		if err != nil {
			t.Fatal(err)
		}
		if resultado.Classificacao != domain.ClasseIndisponibilidade || resultado.EfeitoConcluido || resultado.Registros != nil {
			t.Fatalf("classe=%s efeito=%v", resultado.Classificacao, resultado.EfeitoConcluido)
		}
		if quebra.buscas < 2 {
			t.Fatalf("buscas=%d", quebra.buscas)
		}
	})

	t.Run("isolamento precede pagina inalcancavel", func(t *testing.T) {
		quebra := &proxyQuebra{pagina: prova.Pagina, tenant: prova.CampoTenant, produto: prova.CampoProduto, outroTenant: true}
		srv := httptest.NewServer(quebra)
		defer srv.Close()
		acesso := prova.acesso
		acesso.URL = srv.URL
		ref := referenciaSintetica("default", "stream_sintetico", "lote-isolamento", 100, 200, 1)
		resultado, err := (&application.ServicoRecuperacao{Busca: Cliente{}, Acesso: acesso}).Recuperar(context.Background(), ref)
		if err != nil {
			t.Fatal(err)
		}
		if resultado.Classificacao != domain.ClasseIsolamento || !resultado.EfeitoConcluido || resultado.Registros != nil {
			t.Fatalf("classe=%s efeito=%v", resultado.Classificacao, resultado.EfeitoConcluido)
		}
	})
}

type provaOO struct {
	acesso       domain.ConfiguracaoAcesso
	URL          string
	Usuario      string
	Senha        string
	Timeout      time.Duration
	Pagina       int
	CampoTenant  string
	CampoProduto string
	Copia        string
}

func carregarProva(t *testing.T) provaOO {
	t.Helper()
	urlBase := exigir(t, "FINDER_OO_URL")
	usuario := exigir(t, "FINDER_OO_USUARIO")
	senha := exigir(t, "FINDER_OO_SENHA")
	textoTimeout := exigir(t, "FINDER_OO_TIMEOUT")
	textoPagina := exigir(t, "FINDER_OO_PAGINA")
	ordem := exigir(t, "FINDER_OO_ORDEM")
	campoTenant := exigir(t, "FINDER_OO_CAMPO_TENANT")
	campoProduto := exigir(t, "FINDER_OO_CAMPO_PRODUTO")
	copiaTexto, ok := os.LookupEnv("FINDER_OO_CAMPOS_COPIA")
	if !ok {
		t.Fatal("FINDER_OO_CAMPOS_COPIA ausente")
	}
	timeout, err := time.ParseDuration(textoTimeout)
	if err != nil || timeout <= 0 {
		t.Fatal("FINDER_OO_TIMEOUT inválido")
	}
	pagina := 0
	for _, r := range textoPagina {
		if r < '0' || r > '9' {
			t.Fatal("FINDER_OO_PAGINA inválido")
		}
		pagina = pagina*10 + int(r-'0')
	}
	if pagina <= 0 {
		t.Fatal("FINDER_OO_PAGINA inválido")
	}
	var copia []string
	if strings.TrimSpace(copiaTexto) != "" {
		for _, parte := range strings.Split(copiaTexto, ",") {
			parte = strings.TrimSpace(parte)
			if parte != "" {
				copia = append(copia, parte)
			}
		}
	} else {
		copia = []string{}
	}
	return provaOO{
		URL: urlBase, Usuario: usuario, Senha: senha, Timeout: timeout, Pagina: pagina,
		CampoTenant: campoTenant, CampoProduto: campoProduto,
		Copia: primeiro(copia),
		acesso: domain.ConfiguracaoAcesso{
			URL: urlBase, Usuario: usuario, Senha: senha, Timeout: timeout, TamanhoPagina: pagina,
			CamposCopia: &copia, Ordem: ordem, CampoTenant: campoTenant, CampoProduto: campoProduto,
		},
	}
}

func exigir(t *testing.T, chave string) string {
	t.Helper()
	valor, ok := os.LookupEnv(chave)
	if !ok || strings.TrimSpace(valor) == "" {
		t.Fatalf("%s ausente", chave)
	}
	return valor
}

func responde(prova provaOO) error {
	cliente := &http.Client{Timeout: prova.Timeout}
	req, err := http.NewRequest(http.MethodGet, prova.URL, nil)
	if err != nil {
		return err
	}
	resp, err := cliente.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
	return nil
}

func ingerir(t *testing.T, prova provaOO, org, stream string, docs []map[string]any) {
	t.Helper()
	corpo, err := json.Marshal(docs)
	if err != nil {
		t.Fatal(err)
	}
	endereco := strings.TrimRight(prova.URL, "/") + "/api/" + url.PathEscape(org) + "/" + url.PathEscape(stream) + "/_json"
	req, err := http.NewRequest(http.MethodPost, endereco, bytes.NewReader(corpo))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth(prova.Usuario, prova.Senha)
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{Timeout: prova.Timeout}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode >= 300 {
		t.Fatalf("ingestão sintética recusada: %d", resp.StatusCode)
	}
}

func recuperarVia(t *testing.T, prova provaOO, ref domain.Referencia, visto *coletor, pronto func(domain.ResultadoRecuperacao) bool) domain.ResultadoRecuperacao {
	t.Helper()
	alvo := prova.URL
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/_search") {
			corpo, _ := io.ReadAll(r.Body)
			visto.registrar(r.URL.Path, corpo)
			r.Body = io.NopCloser(bytes.NewReader(corpo))
		}
		encaminhar(w, r, alvo)
	}))
	defer srv.Close()
	acesso := prova.acesso
	acesso.URL = srv.URL
	var resultado domain.ResultadoRecuperacao
	prazo := time.Now().Add(8 * time.Second)
	for {
		saida, err := (&application.ServicoRecuperacao{Busca: Cliente{}, Acesso: acesso}).Recuperar(context.Background(), ref)
		if err != nil {
			t.Fatal(err)
		}
		resultado = saida
		if pronto(saida) || time.Now().After(prazo) {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	return resultado
}

func classificarPronto(classe domain.Classificacao) func(domain.ResultadoRecuperacao) bool {
	return func(r domain.ResultadoRecuperacao) bool {
		return r.Classificacao == classe || r.Classificacao == domain.ClasseAutenticacao || r.Classificacao == domain.ClasseRespostaInvalida
	}
}

func encaminhar(w http.ResponseWriter, r *http.Request, alvo string) {
	endereco := strings.TrimRight(alvo, "/") + r.URL.RequestURI()
	req, err := http.NewRequest(r.Method, endereco, r.Body)
	if err != nil {
		http.Error(w, "pedido inválido", http.StatusBadGateway)
		return
	}
	req.Header = r.Header.Clone()
	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		http.Error(w, "indisponível", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	for chave, valores := range resp.Header {
		for _, valor := range valores {
			w.Header().Add(chave, valor)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

type pedidoVisto struct {
	caminho string
	sql     string
	inicio  int64
	fim     int64
}

type coletor struct {
	mu      sync.Mutex
	pedidos []pedidoVisto
}

func (c *coletor) registrar(caminho string, corpo []byte) {
	var raiz struct {
		Query struct {
			SQL   string `json:"sql"`
			Start int64  `json:"start_time"`
			End   int64  `json:"end_time"`
		} `json:"query"`
	}
	_ = json.Unmarshal(corpo, &raiz)
	c.mu.Lock()
	c.pedidos = append(c.pedidos, pedidoVisto{caminho: caminho, sql: raiz.Query.SQL, inicio: raiz.Query.Start, fim: raiz.Query.End})
	c.mu.Unlock()
}

type proxyQuebra struct {
	pagina      int
	tenant      string
	produto     string
	outroTenant bool
	buscas      int
}

func (p *proxyQuebra) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !strings.Contains(r.URL.Path, "/_search") {
		http.NotFound(w, r)
		return
	}
	p.buscas++
	if p.buscas > 1 {
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "indisponível", http.StatusBadGateway)
			return
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			return
		}
		_ = conn.Close()
		return
	}
	hits := make([]map[string]any, p.pagina)
	for i := range hits {
		tenant := "tenant-a"
		if p.outroTenant && i == 0 {
			tenant = "tenant-b"
		}
		hits[i] = map[string]any{
			"lhamalog_record_id": "r" + unicoCurto(i),
			p.tenant:             tenant,
			p.produto:            "produto-a",
			"mensagem":           "sintetico",
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"hits": hits})
}

func docSintetico(prova provaOO, lote, id, tenant, produto, copia string, instante int64) map[string]any {
	doc := map[string]any{
		"_timestamp":         instante,
		"lhamalog_batch_id":  lote,
		"lhamalog_record_id": id,
		prova.CampoTenant:    tenant,
		prova.CampoProduto:   produto,
		"mensagem":           "sintetico",
	}
	if prova.Copia != "" {
		doc[prova.Copia] = copia
	}
	return doc
}

func referenciaSintetica(org, stream, lote string, minimo, maximo, contagem int64) domain.Referencia {
	return domain.Referencia{
		TenantID: "tenant-a", ProductID: "produto-a", Organizacao: org, Stream: stream,
		BatchID: lote, MinTimestampMicros: minimo, MaxTimestampMicros: maximo, RecordCount: contagem,
		IdempotencyID: lote, CorrelationID: "corr-sintetica", EntradaID: "sintetica",
	}
}

func mesmoCampo(doc domain.DocumentoFisico, campo, esperado string) bool {
	valor, ok := doc[campo]
	return ok && valor == esperado
}

func unico() string {
	return time.Now().Format("20060102150405.000000000")
}

func unicoCurto(i int) string {
	return strings.ReplaceAll(unico(), ".", "") + string(rune('a'+i))
}

func primeiro(lista []string) string {
	if len(lista) == 0 {
		return ""
	}
	return lista[0]
}
