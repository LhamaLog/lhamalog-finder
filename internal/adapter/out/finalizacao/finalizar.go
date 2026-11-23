package finalizacao

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/adapter/out/confirmacao"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

var ErrConclusao = errors.New("conclusao_analise_indisponivel")

type pedido struct {
	TenantID  string `json:"tenantId"`
	ProductID string `json:"productId"`
	BatchID   string `json:"batchId"`
	ExecID    string `json:"execId"`
}
type recibo struct {
	Aceito bool `json:"aceito"`
	pedido
	Finalizada *bool `json:"execucaoFinalizada"`
}
type Varredura struct {
	cliente *confirmacao.Cliente
	ref     domain.Referencia
}

func NovaVarredura(base, token string, timeout time.Duration, ref domain.Referencia, transporte http.RoundTripper) (*Varredura, error) {
	if ref.TenantID == "" || ref.ProductID == "" || ref.BatchID == "" {
		return nil, ErrConclusao
	}
	c, err := confirmacao.Novo(base, "/internal/v1/execucoes/conclusoes", token, timeout, transporte)
	if err != nil {
		return nil, ErrConclusao
	}
	return &Varredura{c, ref}, nil
}

func (v *Varredura) Concluir(ctx context.Context, exec string) error {
	if v == nil || exec == "" {
		return ErrConclusao
	}
	p := pedido{v.ref.TenantID, v.ref.ProductID, v.ref.BatchID, exec}
	var r recibo
	if err := v.cliente.Enviar(ctx, p, &r, "aceito", "tenantId", "productId", "batchId", "execId", "execucaoFinalizada"); err != nil {
		return ErrConclusao
	}
	// false indica outros arquivos pendentes, mas confirma a conclusão deste item.
	if !r.Aceito || r.pedido != p || r.Finalizada == nil {
		return ErrConclusao
	}
	return nil
}

type Finalizador struct {
	Chaves    out.ChaveLote
	Ref       domain.Referencia
	Varredura *Varredura
}

func (f Finalizador) ConcluirAnalise(ctx context.Context, k domain.ChaveResumo) error {
	if f.Chaves == nil || k.BatchID == "" || k.BatchID != f.Ref.BatchID || f.Ref.IdempotencyID != k.BatchID {
		return ErrConclusao
	}
	if k.ExecID != "" {
		if err := f.Varredura.Concluir(ctx, k.ExecID); err != nil {
			return err
		}
	}
	// log.ingested v1 não traz execId/arquivoId: a marca local é a finalização real.
	return f.Chaves.Concluir(ctx, f.Ref.IdempotencyID)
}
