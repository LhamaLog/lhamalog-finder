package openobserve

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// Cliente envia POST /api/{organization}/_search com net/http.
// Não escolhe URL, senha, timeout, página, lista, ordem nem nomes de campo.
type Cliente struct{}

var _ out.Busca = Cliente{}

// Consultar executa uma página. Falha não registra senha nem hits.
func (Cliente) Consultar(ctx context.Context, pedido out.PedidoPagina) ([]domain.DocumentoFisico, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	endereco, err := urlBusca(pedido.URL, pedido.Organizacao)
	if err != nil {
		return nil, &out.ErroBusca{Tipo: out.FalhaResposta}
	}
	corpo, err := CorpoJSON(pedido.SQL, pedido.Inicio, pedido.Fim, pedido.From, pedido.Size)
	if err != nil {
		return nil, &out.ErroBusca{Tipo: out.FalhaResposta}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endereco, bytes.NewReader(corpo))
	if err != nil {
		return nil, &out.ErroBusca{Tipo: out.FalhaResposta}
	}
	req.SetBasicAuth(pedido.Usuario, pedido.Senha)
	req.Header.Set("Content-Type", "application/json")

	cliente := &http.Client{
		Timeout: pedido.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := cliente.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if tempoEsgotado(err) {
			return nil, &out.ErroBusca{Tipo: out.FalhaTimeout}
		}
		return nil, &out.ErroBusca{Tipo: out.FalhaIndisponivel}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, &out.ErroBusca{Tipo: out.FalhaCredencial}
	}
	if resp.StatusCode == http.StatusRequestTimeout || resp.StatusCode == http.StatusGatewayTimeout {
		return nil, &out.ErroBusca{Tipo: out.FalhaTimeout}
	}
	if resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusServiceUnavailable {
		return nil, &out.ErroBusca{Tipo: out.FalhaIndisponivel}
	}
	if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusTooManyRequests {
		return nil, &out.ErroBusca{Tipo: out.FalhaPagina}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &out.ErroBusca{Tipo: out.FalhaPagina}
	}
	return lerHits(resp.Body)
}

func urlBusca(base, organizacao string) (string, error) {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	endereco, err := url.Parse(base)
	if err != nil || endereco.Scheme == "" || endereco.Host == "" {
		return "", errURL
	}
	if endereco.Scheme != "http" && endereco.Scheme != "https" {
		return "", errURL
	}
	endereco.User = nil
	endereco.RawQuery = ""
	endereco.Fragment = ""
	endereco.Path = strings.TrimRight(endereco.Path, "/") + "/api/" + url.PathEscape(organizacao) + "/_search"
	return endereco.String(), nil
}

var errURL = errors.New("url inválida")

func tempoEsgotado(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var rede net.Error
	return errors.As(err, &rede) && rede.Timeout()
}

func lerHits(corpo io.Reader) ([]domain.DocumentoFisico, error) {
	bruto, err := io.ReadAll(corpo)
	if err != nil {
		if tempoEsgotado(err) {
			return nil, &out.ErroBusca{Tipo: out.FalhaTimeout}
		}
		return nil, &out.ErroBusca{Tipo: out.FalhaPagina}
	}
	dec := json.NewDecoder(bytes.NewReader(bruto))
	dec.UseNumber()
	var raiz map[string]any
	if err := dec.Decode(&raiz); err != nil {
		return nil, &out.ErroBusca{Tipo: out.FalhaResposta}
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return nil, &out.ErroBusca{Tipo: out.FalhaResposta}
	}
	hits, ok := raiz["hits"]
	if !ok || hits == nil {
		return nil, &out.ErroBusca{Tipo: out.FalhaResposta}
	}
	lista, ok := hits.([]any)
	if !ok {
		return nil, &out.ErroBusca{Tipo: out.FalhaResposta}
	}
	docs := make([]domain.DocumentoFisico, 0, len(lista))
	for _, item := range lista {
		objeto, ok := item.(map[string]any)
		if !ok {
			return nil, &out.ErroBusca{Tipo: out.FalhaResposta}
		}
		docs = append(docs, domain.DocumentoFisico(objeto))
	}
	return docs, nil
}
