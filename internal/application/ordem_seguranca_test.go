package application

import (
	"context"
	entrada "github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"testing"
)

func TestFronteira010AntesDaClassificacaoEAssinaturaCanonica(t *testing.T) {
	chamadas := 0
	scanner := Seguranca{Catalogo: domain.CatalogoSeguranca{Versao: "catalogo-teste", Transformar: func(string) (domain.LinhaLogMascarada, int, error) {
		chamadas++
		return domain.LinhaLogMascarada{Texto: "seguro WARNING"}, 2, nil
	}}}
	s := ServicoAnalise{Seguranca: scanner}
	p := entrada.PedidoAnalise{Referencia: domain.Referencia{BatchID: "b", TenantID: "t", ProductID: "p", RecordCount: 1}, ProdutoBody: "p", SchemaVersion: "1", Completo: true, VersaoRegras: "v1", VersaoCatalogo: "catalogo-teste", Logs: []entrada.LogAnalise{{Apontamento: "r", NivelCampo: "info", Mensagem: "BRUTO ERROR SENTINELA", ServiceName: "svc", Host: "host"}}}
	r, err := s.Montar(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
	a := r.Achados[0]
	if chamadas != 1 || a.Tipo != domain.WARNING || a.Assinatura != "WARNING|seguro warning|svc|host" || a.Linhas[0] != "seguro WARNING" || r.Contadores.DadosSensiveisMascarados != 2 || r.VersaoCatalogo != "catalogo-teste" {
		t.Fatal("ordem/saida/contador/versao invalidos")
	}
	// Este é um teste de mecanismo com transformador opaco, não caso oficial P-010-27.
}
