package openobserve

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

func TestMapearDocumentoAnaliseDoColetor(t *testing.T) {
	doc := domain.DocumentoFisico{
		"lhamalog_record_id": "registro", "_timestamp": json.Number("1795104000123456"),
		"severitytext": "WARN", "body": "mensagem", "servicename": "checkout", "traceid": "trace",
		"resourceattributes": `{"host.name":"maquina","service.name":"checkout"}`,
	}
	got, err := MapearDocumentoAnalise(doc)
	if err != nil {
		t.Fatal(err)
	}
	if got.Apontamento != "registro" || !got.Timestamp.Equal(time.UnixMicro(1795104000123456)) || got.NivelCampo != "WARN" || got.Mensagem != "mensagem" || got.ServiceName != "checkout" || got.Host != "maquina" || got.TraceID != "trace" {
		t.Fatalf("campos do coletor mapeados incorretamente: %+v", got)
	}
	delete(doc, "severitytext")
	delete(doc, "body")
	delete(doc, "traceid")
	got, err = MapearDocumentoAnalise(doc)
	if err != nil || got.NivelCampo != "" || got.Mensagem != "" || got.TraceID != "" {
		t.Fatalf("opcionais ausentes: %+v %v", got, err)
	}
}

func TestMapearDocumentoAnaliseRejeitaCamposInvalidos(t *testing.T) {
	for _, c := range []struct {
		campo string
		valor any
		quer  error
	}{
		{"lhamalog_record_id", nil, errDocumentoAnalise},
		{"_timestamp", json.Number("123.4"), errDocumentoAnalise},
		{"resourceattributes", `{"host.name":""}`, errIdentidadeAssinatura},
		{"resourceattributes", `{"host.name":"x"} }`, errDocumentoAnalise},
		{"servicename", " ", errIdentidadeAssinatura},
		{"body", map[string]string{"privado": "nao copiar"}, errDocumentoAnalise},
	} {
		t.Run(c.campo, func(t *testing.T) {
			doc := domain.DocumentoFisico{"lhamalog_record_id": "x", "_timestamp": json.Number("123"), "servicename": "svc", "resourceattributes": `{"host.name":"host"}`}
			doc[c.campo] = c.valor
			got, err := MapearDocumentoAnalise(doc)
			if !errors.Is(err, c.quer) || got.Apontamento != "" || got.Mensagem != "" {
				t.Fatalf("falha esperada %v, resultado %+v e erro %v", c.quer, got, err)
			}
		})
	}
}
