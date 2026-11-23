package core

import (
	"context"
	"net/http"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/adapter/out/confirmacao"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

type aviso struct {
	TenantID  string `json:"tenantId"`
	ProductID string `json:"productId"`
	BatchID   string `json:"batchId"`
	Motivo    string `json:"motivo"`
}
type recibo struct {
	Aceito bool `json:"aceito"`
	aviso
}
type Aviso struct {
	cliente    *confirmacao.Cliente
	identidade aviso
}

func NovoAviso(base, token string, timeout time.Duration, ref domain.Referencia, transporte http.RoundTripper) (*Aviso, error) {
	if ref.TenantID == "" || ref.ProductID == "" || ref.BatchID == "" {
		return nil, confirmacao.ErrConfirmacao
	}
	c, err := confirmacao.Novo(base, "/internal/v1/alertas/finder", token, timeout, transporte)
	if err != nil {
		return nil, err
	}
	return &Aviso{cliente: c, identidade: aviso{ref.TenantID, ref.ProductID, ref.BatchID, "consulta_openobserve_indisponivel"}}, nil
}

func (a *Aviso) AvisarFalhaOpenObserve(ctx context.Context, batch, motivo string) error {
	if a == nil || batch != a.identidade.BatchID || motivo != a.identidade.Motivo {
		return confirmacao.ErrConfirmacao
	}
	var r recibo
	if err := a.cliente.Enviar(ctx, a.identidade, &r, "aceito", "tenantId", "productId", "batchId", "motivo"); err != nil {
		return err
	}
	if !r.Aceito || r.aviso != a.identidade {
		return confirmacao.ErrConfirmacao
	}
	return nil
}
