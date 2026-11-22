// Package bootstrap monta o único efeito usado pelo consumidor de produção.
package bootstrap

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/adapter/out/core"
	"github.com/LhamaLog/lhamalog-finder/internal/adapter/out/finalizacao"
	"github.com/LhamaLog/lhamalog-finder/internal/adapter/out/modocorrencia"
	"github.com/LhamaLog/lhamalog-finder/internal/adapter/out/openobserve"
	"github.com/LhamaLog/lhamalog-finder/internal/application"
	"github.com/LhamaLog/lhamalog-finder/internal/application/agrupar"
	"github.com/LhamaLog/lhamalog-finder/internal/application/handoff"
	entrada "github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/config"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

var ErrConfiguracao = errors.New("configuracao_pipeline_invalida")

// Transporte é opcional para testes TLS. A produção usa a cadeia de confiança do sistema.
type Opcoes struct{ Transporte http.RoundTripper }
type Pipeline struct {
	acesso                                       domain.ConfiguracaoAcesso
	seguranca                                    application.Seguranca
	modulo                                       *modocorrencia.Cliente
	chaves                                       out.ChaveLote
	coreURL, coreJWT, varreduraURL, varreduraJWT string
	coreTimeout, varreduraTimeout                time.Duration
	transporte                                   http.RoundTripper
}

var _ out.Efeito = (*Pipeline)(nil)

func Novo(ctx context.Context, cfg config.Config, chaves out.ChaveLote, o Opcoes) (*Pipeline, error) {
	if chaves == nil || cfg.ValidarConsumo() != nil {
		return nil, ErrConfiguracao
	}
	acesso := cfg.Acesso()
	if !acesso.Pronta() || acesso.CampoTenant == "" || acesso.CampoProduto == "" {
		return nil, ErrConfiguracao
	}
	seguranca, err := config.InicializarSeguranca(ctx)
	if err != nil {
		return nil, ErrConfiguracao
	}
	analise := config.CarregarAnalise()
	modulo, err := modocorrencia.NovoCliente(modocorrencia.OpcoesCliente{URLBase: analise.URLModulo, JWT: analise.JWTModulo, Transporte: o.Transporte})
	if err != nil {
		return nil, ErrConfiguracao
	}
	coreTimeout, err := prazo("FINDER_CORE_TIMEOUT")
	if err != nil {
		return nil, err
	}
	varreduraTimeout, err := prazo("FINDER_VARREDURA_TIMEOUT")
	if err != nil {
		return nil, err
	}
	p := &Pipeline{acesso: acesso, seguranca: seguranca, modulo: modulo, chaves: chaves, coreURL: os.Getenv("FINDER_CORE_URL"), coreJWT: os.Getenv("FINDER_CORE_JWT"), varreduraURL: os.Getenv("FINDER_VARREDURA_URL"), varreduraJWT: os.Getenv("FINDER_VARREDURA_JWT"), coreTimeout: coreTimeout, varreduraTimeout: varreduraTimeout, transporte: o.Transporte}
	// Validar formatos e contratos antes de abrir saúde/consumo. Sem chamadas externas.
	ref := domain.Referencia{TenantID: "validacao", ProductID: "validacao", Organizacao: "validacao", Stream: "validacao", BatchID: "validacao", RecordCount: 1, MaxTimestampMicros: 1}
	if _, err := openobserve.NovoHistorico(acesso, ref, func(context.Context, domain.Referencia, []domain.DocumentoFisico) (domain.ResumoAnalise, error) {
		return domain.ResumoAnalise{}, nil
	}); err != nil {
		return nil, ErrConfiguracao
	}
	if _, err := core.NovoAviso(p.coreURL, p.coreJWT, coreTimeout, ref, o.Transporte); err != nil {
		return nil, ErrConfiguracao
	}
	if p.varreduraURL != "" || p.varreduraJWT != "" {
		if _, err := finalizacao.NovaVarredura(p.varreduraURL, p.varreduraJWT, varreduraTimeout, ref, o.Transporte); err != nil {
			return nil, ErrConfiguracao
		}
	}
	return p, nil
}

func prazo(env string) (time.Duration, error) {
	s := os.Getenv(env)
	if s == "" {
		return 2 * time.Second, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return 0, ErrConfiguracao
	}
	return d, nil
}

type montadorFunc func(context.Context, entrada.PedidoAnalise) (domain.ResumoAnalise, error)

func (f montadorFunc) Montar(ctx context.Context, p entrada.PedidoAnalise) (domain.ResumoAnalise, error) {
	return f(ctx, p)
}

func (p *Pipeline) Executar(ctx context.Context, ref domain.Referencia) (bool, error) {
	if p == nil {
		return false, ErrConfiguracao
	}
	servico := application.ServicoAnalise{Seguranca: p.seguranca}
	var atual *domain.ResumoAnalise
	analisar := func(ctx context.Context, r domain.Referencia, docs []domain.DocumentoFisico) (domain.ResumoAnalise, error) {
		// O lote corrente já passou pelo scanner em Montar; não o executar duas vezes.
		if r.BatchID == ref.BatchID && atual != nil {
			return *atual, nil
		}
		logs := make([]entrada.LogAnalise, 0, len(docs))
		for _, doc := range docs {
			l, e := openobserve.MapearDocumentoAnalise(doc)
			if e != nil {
				return domain.ResumoAnalise{}, e
			}
			logs = append(logs, l)
		}
		return servico.Montar(ctx, entrada.PedidoAnalise{Referencia: r, ProdutoBody: r.ProductID, SchemaVersion: "1", Completo: true, VersaoRegras: p.seguranca.Catalogo.Versao, VersaoCatalogo: p.seguranca.Catalogo.Versao, Logs: logs})
	}
	historico, err := openobserve.NovoHistorico(p.acesso, ref, analisar)
	if err != nil {
		return false, err
	}
	aviso, err := core.NovoAviso(p.coreURL, p.coreJWT, p.coreTimeout, ref, p.transporte)
	if err != nil {
		return false, err
	}
	var varredura *finalizacao.Varredura
	if p.varreduraURL != "" {
		varredura, err = finalizacao.NovaVarredura(p.varreduraURL, p.varreduraJWT, p.varreduraTimeout, ref, p.transporte)
		if err != nil {
			return false, err
		}
	}
	preparar := agrupar.Orquestrador{Montador: montadorFunc(func(ctx context.Context, pedido entrada.PedidoAnalise) (domain.ResumoAnalise, error) {
		r, e := servico.Montar(ctx, pedido)
		if e == nil {
			atual = &r
		}
		return r, e
	}), Consulta: p.modulo, Recorrencia: openobserve.Recorrencia{Historico: historico, Aviso: aviso}}
	efeito := handoff.EfeitoAnalise{Recuperacao: &application.ServicoRecuperacao{Busca: openobserve.Cliente{}, Acesso: p.acesso}, Mapear: openobserve.MapearDocumentoAnalise, VersaoRegras: p.seguranca.Catalogo.Versao, VersaoCatalogo: p.seguranca.Catalogo.Versao, Entrega: handoff.OrquestradorEntrega{Preparador: preparar, Entrega: handoff.Entrega{Envio: p.modulo, Finalizar: finalizacao.Finalizador{Chaves: p.chaves, Ref: ref, Varredura: varredura}}}}
	return efeito.Executar(ctx, ref)
}
