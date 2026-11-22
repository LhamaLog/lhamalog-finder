package openobserve

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// AnalisarHistorico usa o scanner/classificador corrente, nunca mensagem bruta
// como assinatura. O callback deve ser o mesmo montador empregado no lote atual.
type AnalisarHistorico func(context.Context, domain.Referencia, []domain.DocumentoFisico) (domain.ResumoAnalise, error)

// HistoricoHTTP pertence a uma única referência/efeito. O snapshot derivado é
// reutilizado por todos os achados e pelo único recálculo de revisão do Módulo.
type HistoricoHTTP struct {
	cfg        domain.ConfiguracaoAcesso
	ref        domain.Referencia
	analisar   AnalisarHistorico
	mu         sync.Mutex
	pronto     bool
	evidencias []domain.EvidenciaRecorrencia
}

func NovoHistorico(cfg domain.ConfiguracaoAcesso, ref domain.Referencia, analisar AnalisarHistorico) (*HistoricoHTTP, error) {
	if !cfg.Pronta() || analisar == nil || ref.RecordCount <= 0 {
		return nil, out.FalhaRecorrencia{}
	}
	if _, _, err := domain.ConsultasHistorico(ref, cfg, nil); err != nil {
		return nil, out.FalhaRecorrencia{}
	}
	if _, err := urlBusca(cfg.URL, ref.Organizacao); err != nil {
		return nil, out.FalhaRecorrencia{}
	}
	copia := append([]string{}, (*cfg.CamposCopia)...)
	cfg.CamposCopia = &copia
	return &HistoricoHTTP{cfg: cfg, ref: ref, analisar: analisar}, nil
}

func (h *HistoricoHTTP) ConsultarHistorico(ctx context.Context, p out.PedidoRecorrencia) ([]domain.EvidenciaRecorrencia, error) {
	if h == nil {
		return nil, out.FalhaRecorrencia{}
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if p.BatchID != h.ref.BatchID || p.Chave.TenantID != h.ref.TenantID || p.Chave.ProductID != h.ref.ProductID {
		return nil, out.FalhaRecorrencia{}
	}
	if !h.pronto {
		ev, err := h.carregar(ctx)
		if err != nil {
			return nil, err
		}
		h.evidencias = ev
		h.pronto = true
	}
	return append([]domain.EvidenciaRecorrencia(nil), h.evidencias...), nil
}

type loteHistorico struct {
	id                 string
	quantidade, ultimo int64
}

func (h *HistoricoHTTP) carregar(ctx context.Context) ([]domain.EvidenciaRecorrencia, error) {
	sql, fim, err := domain.ConsultasHistorico(h.ref, h.cfg, nil)
	if err != nil {
		return nil, out.FalhaRecorrencia{}
	}
	grupos, err := h.pagina(ctx, sql, fim, 0, 499)
	if err != nil {
		return nil, err
	}
	if len(grupos) > 499 {
		return nil, out.FalhaRecorrencia{}
	}
	lotes := []loteHistorico{{h.ref.BatchID, h.ref.RecordCount, h.ref.MaxTimestampMicros}}
	esperados := map[string]int64{h.ref.BatchID: h.ref.RecordCount}
	for _, g := range grupos {
		id, ok := idHistorico(g["batch_id"])
		q, okq := inteiroHistorico(g["record_count"])
		ultimo, oku := inteiroHistorico(g["last_timestamp"])
		if !ok || !okq || !oku || q <= 0 || ultimo < 0 || ultimo > h.ref.MaxTimestampMicros || esperados[id] != 0 {
			return nil, out.FalhaRecorrencia{}
		}
		esperados[id] = q
		lotes = append(lotes, loteHistorico{id, q, ultimo})
	}
	sort.Slice(lotes[1:], func(i, j int) bool {
		a, b := lotes[i+1], lotes[j+1]
		if a.ultimo != b.ultimo {
			return a.ultimo > b.ultimo
		}
		return a.id < b.id
	})
	ids := make([]string, len(lotes))
	for i, l := range lotes {
		ids[i] = l.id
	}
	sql, fim, err = domain.ConsultasHistorico(h.ref, h.cfg, ids)
	if err != nil {
		return nil, out.FalhaRecorrencia{}
	}
	derivados := map[string][]domain.EvidenciaRecorrencia{}
	vistos := map[string]bool{}
	atual := ""
	var docs []domain.DocumentoFisico
	fechar := func() error {
		if atual == "" {
			return nil
		}
		unicos, divergente, semID := domain.Deduplicar(docs, *h.cfg.CamposCopia)
		if divergente || semID {
			return out.FalhaRecorrencia{}
		}
		if int64(len(unicos)) != esperados[atual] {
			return out.FalhaRecorrencia{Temporaria: true}
		}
		ref := h.ref
		ref.BatchID = atual
		ref.IdempotencyID = atual
		ref.RecordID = nil
		ref.RecordCount = int64(len(unicos))
		ref.MinTimestampMicros = 0
		resumo, e := h.analisar(ctx, ref, unicos)
		if e != nil {
			return out.FalhaRecorrencia{}
		}
		if resumo.BatchID != atual || resumo.TenantID != h.ref.TenantID || resumo.ProductID != h.ref.ProductID {
			return out.FalhaRecorrencia{}
		}
		for _, a := range resumo.Achados {
			if a.Assinatura == "" || a.HostName == "" {
				return out.FalhaRecorrencia{}
			}
			derivados[atual] = append(derivados[atual], domain.EvidenciaRecorrencia{BatchID: atual, Host: a.HostName, Chave: domain.IdentidadeAchado{Assinatura: a.Assinatura, Classe: a.Tipo, TenantID: h.ref.TenantID, ProductID: h.ref.ProductID}})
		}
		vistos[atual] = true
		docs = nil
		return nil
	}
	for from := 0; ; {
		pagina, e := h.pagina(ctx, sql, fim, from, h.cfg.TamanhoPagina)
		if e != nil {
			return nil, e
		}
		if len(pagina) > h.cfg.TamanhoPagina || domain.Isolado(pagina, h.cfg.CampoTenant, h.cfg.CampoProduto, h.ref.TenantID, h.ref.ProductID) {
			return nil, out.FalhaRecorrencia{}
		}
		for _, doc := range pagina {
			id, ok := idHistorico(doc["lhamalog_batch_id"])
			ts, okt := inteiroHistorico(doc["_timestamp"])
			if !ok || esperados[id] == 0 || !okt || ts < 0 || ts > h.ref.MaxTimestampMicros {
				return nil, out.FalhaRecorrencia{}
			}
			if atual != id {
				if err := fechar(); err != nil {
					return nil, err
				}
				if vistos[id] {
					return nil, out.FalhaRecorrencia{}
				}
				atual = id
			}
			docs = append(docs, doc)
		}
		if len(pagina) < h.cfg.TamanhoPagina {
			break
		}
		if from > math.MaxInt-h.cfg.TamanhoPagina {
			return nil, out.FalhaRecorrencia{}
		}
		from += h.cfg.TamanhoPagina
	}
	if err := fechar(); err != nil {
		return nil, err
	}
	var ev []domain.EvidenciaRecorrencia
	for _, l := range lotes {
		if !vistos[l.id] {
			return nil, out.FalhaRecorrencia{Temporaria: true}
		}
		// O marcador inclui também os lotes sem achados equivalentes na janela.
		ev = append(ev, domain.EvidenciaRecorrencia{BatchID: l.id})
		ev = append(ev, derivados[l.id]...)
	}
	return ev, nil
}

func (h *HistoricoHTTP) pagina(ctx context.Context, sql string, fim int64, from, size int) ([]domain.DocumentoFisico, error) {
	endpoint, _ := urlBusca(h.cfg.URL, h.ref.Organizacao)
	data, err := CorpoJSON(sql, 1, fim, from, size)
	if err != nil {
		return nil, out.FalhaRecorrencia{}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, out.FalhaRecorrencia{}
	}
	req.SetBasicAuth(h.cfg.Usuario, h.cfg.Senha)
	req.Header.Set("Content-Type", "application/json")
	cliente := http.Client{Timeout: h.cfg.Timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := cliente.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, out.FalhaRecorrencia{Temporaria: true}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 || resp.StatusCode >= 500 || resp.StatusCode == 408 {
		return nil, out.FalhaRecorrencia{Temporaria: true}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, out.FalhaRecorrencia{}
	}
	// Limite por página; uma resposta truncada nunca produz recorrente=false.
	data, err = io.ReadAll(io.LimitReader(resp.Body, 64*1024*1024+1))
	if err != nil {
		return nil, out.FalhaRecorrencia{Temporaria: true}
	}
	if len(data) > 64*1024*1024 {
		return nil, out.FalhaRecorrencia{}
	}
	var estado struct {
		Parcial bool   `json:"is_partial"`
		Erro    string `json:"function_error"`
	}
	if json.Unmarshal(data, &estado) != nil {
		return nil, out.FalhaRecorrencia{}
	}
	if estado.Parcial || estado.Erro != "" {
		return nil, out.FalhaRecorrencia{Temporaria: true}
	}
	docs, err := lerHits(bytes.NewReader(data))
	if err != nil {
		return nil, out.FalhaRecorrencia{}
	}
	return docs, nil
}

func idHistorico(v any) (string, bool) {
	switch x := v.(type) {
	case string:
		return x, x != ""
	case json.Number:
		_, err := strconv.ParseInt(x.String(), 10, 64)
		return x.String(), err == nil
	}
	return "", false
}
func inteiroHistorico(v any) (int64, bool) {
	switch x := v.(type) {
	case json.Number:
		n, e := x.Int64()
		return n, e == nil
	case int64:
		return x, true
	case int:
		return int64(x), true
	}
	return 0, false
}
