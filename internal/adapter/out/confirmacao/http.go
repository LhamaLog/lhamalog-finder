// Package confirmacao implementa POST interno HTTPS com JWT e recibo JSON estrito.
package confirmacao

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/strictjson"
)

var ErrConfirmacao = errors.New("confirmacao_interna_indisponivel")

type Cliente struct {
	endpoint, token string
	http            *http.Client
}

// Novo não permite redirecionar o JWT nem incluir credenciais na URL.
// transporte nil usa a implementação padrão; substituição permite testes TLS.
func Novo(base, rota, token string, timeout time.Duration, transporte http.RoundTripper) (*Cliente, error) {
	u, err := url.Parse(base)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") || strings.TrimSpace(token) == "" || strings.ContainsAny(token, "\r\n") || timeout <= 0 || !strings.HasPrefix(rota, "/") {
		return nil, ErrConfirmacao
	}
	u.Path = rota
	return &Cliente{endpoint: u.String(), token: token, http: &http.Client{Timeout: timeout, Transport: transporte, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}, nil
}

// Enviar não retorna corpo, token, erro de rede ou URL em diagnósticos.
func (c *Cliente) Enviar(ctx context.Context, pedido, recibo any, nomes ...string) error {
	if c == nil {
		return ErrConfirmacao
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	data, err := json.Marshal(pedido)
	if err != nil {
		return ErrConfirmacao
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(data))
	if err != nil {
		return ErrConfirmacao
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.http.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return ErrConfirmacao
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return ErrConfirmacao
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024+1))
	if err != nil || len(body) > 64*1024 {
		return ErrConfirmacao
	}
	if strictjson.Decode(body, recibo, nomes...) != nil {
		return ErrConfirmacao
	}
	return nil
}
