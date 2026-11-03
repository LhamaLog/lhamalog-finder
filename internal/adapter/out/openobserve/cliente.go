// Package openobserve lê um log na API HTTP do OpenObserve.
// O apontamento é o nome do stream. A organização padrão da ferramenta é "default".
// O registro devolvido é o primeiro hit, sem campo removido.
package openobserve

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// Cliente implementa a leitura de um stream no OpenObserve.
type Cliente struct {
	BaseURL string
	Org     string
	Usuario string
	Senha   string
	HTTP    *http.Client
	Agora   func() time.Time
}

// Novo monta o cliente com a biblioteca HTTP padrão.
func Novo(baseURL, org, usuario, senha string) *Cliente {
	return &Cliente{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Org:     org,
		Usuario: usuario,
		Senha:   senha,
		HTTP:    &http.Client{},
		Agora:   time.Now,
	}
}

var _ out.LeituraLog = (*Cliente)(nil)

type consulta struct {
	Query struct {
		SQL       string `json:"sql"`
		StartTime int64  `json:"start_time"`
		EndTime   int64  `json:"end_time"`
		From      int    `json:"from"`
		Size      int    `json:"size"`
	} `json:"query"`
}

type resposta struct {
	Hits  []json.RawMessage `json:"hits"`
	Total int               `json:"total"`
}

// Ler busca um registro no stream indicado pelo apontamento.
// Falha de rede ou resposta 5xx é inacessível. Resposta sem hit é ausência.
func (c *Cliente) Ler(ctx context.Context, codigo domain.Codigo) (domain.Registro, out.EstadoLeitura, error) {
	if c == nil || c.BaseURL == "" {
		return nil, out.LeituraInacessivel, fmt.Errorf("openobserve sem url")
	}
	stream := string(codigo)
	if !streamAceito(stream) {
		return nil, out.LeituraInacessivel, fmt.Errorf("apontamento %q não é um nome de stream aceito", codigo)
	}
	corpo, err := json.Marshal(novaConsulta(stream, c.agora()))
	if err != nil {
		return nil, out.LeituraInacessivel, err
	}
	endpoint, err := url.JoinPath(c.BaseURL, "api", c.Org, "_search")
	if err != nil {
		return nil, out.LeituraInacessivel, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(corpo))
	if err != nil {
		return nil, out.LeituraInacessivel, err
	}
	req.SetBasicAuth(c.Usuario, c.Senha)
	req.Header.Set("Content-Type", "application/json")
	cliente := c.HTTP
	if cliente == nil {
		cliente = http.DefaultClient
	}
	resp, err := cliente.Do(req)
	if err != nil {
		return nil, out.LeituraInacessivel, fmt.Errorf("openobserve inacessível em %s: %w", endpoint, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, out.LeituraInacessivel, fmt.Errorf("openobserve leitura da resposta %s: %w", endpoint, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, out.LeituraAusente, fmt.Errorf("openobserve %s: log ausente", endpoint)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, out.LeituraInacessivel, fmt.Errorf("openobserve %s: status %d: %s", endpoint, resp.StatusCode, bytes.TrimSpace(body))
	}
	var parsed resposta
	if err := json.Unmarshal(body, &parsed); err != nil {
		return domain.Registro(body), out.LeituraPronta, nil
	}
	if len(parsed.Hits) == 0 {
		return nil, out.LeituraAusente, fmt.Errorf("openobserve %s: nenhum hit para %s", endpoint, codigo)
	}
	return domain.Registro(parsed.Hits[0]), out.LeituraPronta, nil
}

func (c *Cliente) agora() time.Time {
	if c.Agora != nil {
		return c.Agora()
	}
	return time.Now()
}

func novaConsulta(stream string, agora time.Time) consulta {
	var q consulta
	q.Query.SQL = `SELECT * FROM "` + stream + `"`
	q.Query.StartTime = 0
	q.Query.EndTime = agora.UnixMicro()
	q.Query.From = 0
	q.Query.Size = 1
	return q
}

func streamAceito(stream string) bool {
	if stream == "" {
		return false
	}
	for _, r := range stream {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '.' {
			continue
		}
		return false
	}
	return true
}
