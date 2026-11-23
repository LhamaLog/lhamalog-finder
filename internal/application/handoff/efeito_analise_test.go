package handoff_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/LhamaLog/lhamalog-finder/internal/adapter/out/openobserve"
	"github.com/LhamaLog/lhamalog-finder/internal/application/handoff"
	entrada "github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

type recuperacaoEfeito struct {
	resultado domain.ResultadoRecuperacao
}

func (r recuperacaoEfeito) Recuperar(context.Context, domain.Referencia) (domain.ResultadoRecuperacao, error) {
	return r.resultado, nil
}

type preparadorEfeito struct {
	chamadas int
	pedido   entrada.PedidoAnalise
}

func (p *preparadorEfeito) Montar(_ context.Context, pedido entrada.PedidoAnalise) (domain.ResumoAnalise, error) {
	p.chamadas++
	p.pedido = pedido
	return domain.ResumoAnalise{
		BatchID: pedido.Referencia.BatchID, TenantID: pedido.Referencia.TenantID, ProductID: pedido.Referencia.ProductID,
		VersaoRegras: pedido.VersaoRegras,
		Contadores:   domain.ContadoresResumo{LinhasProcessadas: 1, Warnings: 1},
		Achados:      []domain.AchadoAnalise{{Tipo: domain.WARNING}},
	}, nil
}

type moduloEfeito struct {
	aceita             bool
	envios, conclusoes int
}

func (m *moduloEfeito) Enviar(context.Context, domain.ResumoAnalise) (bool, error) {
	m.envios++
	return m.aceita, nil
}

func (m *moduloEfeito) ConcluirAnalise(context.Context, domain.ChaveResumo) error {
	m.conclusoes++
	return nil
}

func TestEfeitoAnaliseMantemConfirmacaoAposEntrega(t *testing.T) {
	ref := domain.Referencia{BatchID: "b", TenantID: "t", ProductID: "p", RecordCount: 1}
	doc := domain.DocumentoFisico{"lhamalog_record_id": "r", "_timestamp": json.Number("123"), "severitytext": "WARN", "servicename": "s", "resourceattributes": `{"host.name":"h"}`}
	for _, c := range []struct {
		nome                               string
		classe                             domain.Classificacao
		consultaConcluida, aceita, liberar bool
		mapeados, enviados, concluidos     int
	}{
		{"incompleto", domain.ClasseIncompleto, false, false, false, 0, 0, 0},
		{"recuperacao rejeitada definitivamente", domain.ClasseIsolamento, true, false, true, 0, 0, 0},
		{"lote completo sem confirmacao", domain.ClasseLoteCompleto, true, false, false, 1, 1, 0},
		{"lote completo confirmado", domain.ClasseLoteCompleto, true, true, true, 1, 1, 1},
	} {
		t.Run(c.nome, func(t *testing.T) {
			p := &preparadorEfeito{}
			m := &moduloEfeito{aceita: c.aceita}
			mapeados := 0
			e := handoff.EfeitoAnalise{
				Recuperacao: recuperacaoEfeito{domain.ResultadoRecuperacao{Classificacao: c.classe, EfeitoConcluido: c.consultaConcluida, Registros: []domain.DocumentoFisico{doc}}},
				Mapear: func(d domain.DocumentoFisico) (entrada.LogAnalise, error) {
					mapeados++
					return openobserve.MapearDocumentoAnalise(d)
				},
				Entrega:      handoff.OrquestradorEntrega{Preparador: p, Entrega: handoff.Entrega{Envio: m, Finalizar: m}},
				VersaoRegras: "regras-v1", VersaoCatalogo: "catalogo-v1",
			}
			ok, err := e.Executar(context.Background(), ref)
			querErro := c.classe == domain.ClasseLoteCompleto && !c.aceita
			if (err != nil) != querErro || ok != c.liberar || mapeados != c.mapeados || p.chamadas != c.mapeados || m.envios != c.enviados || m.conclusoes != c.concluidos {
				t.Fatalf("resultado=%v erro=%v mapeados=%d preparados=%d envios=%d conclusoes=%d", ok, err, mapeados, p.chamadas, m.envios, m.conclusoes)
			}
			if p.chamadas == 1 && (p.pedido.ExecID != "" || p.pedido.Logs[0].Apontamento != "r" || p.pedido.ProdutoBody != ref.ProductID || p.pedido.VersaoRegras != "regras-v1" || p.pedido.VersaoCatalogo != "catalogo-v1") {
				t.Fatalf("pedido inconsistente: %+v", p.pedido)
			}
		})
	}
}

func TestEfeitoAnaliseInterrompeAntesDaEntregaEmErroDeMapeamento(t *testing.T) {
	falha := errors.New("documento invalido")
	p := &preparadorEfeito{}
	e := handoff.EfeitoAnalise{
		Recuperacao: recuperacaoEfeito{domain.ResultadoRecuperacao{Classificacao: domain.ClasseLoteCompleto, EfeitoConcluido: true, Registros: []domain.DocumentoFisico{{}}}},
		Mapear:      func(domain.DocumentoFisico) (entrada.LogAnalise, error) { return entrada.LogAnalise{}, falha },
		Entrega:     handoff.OrquestradorEntrega{Preparador: p},
	}
	ok, err := e.Executar(context.Background(), domain.Referencia{RecordCount: 1})
	if ok || !errors.Is(err, falha) || p.chamadas != 0 {
		t.Fatalf("efeito=%v erro=%v, chamadas=%d", ok, err, p.chamadas)
	}
}

func (p *preparadorEfeito) Agrupar(_ context.Context, r domain.ResumoAnalise) (domain.ResumoAnalise, error) {
	return r, nil
}
