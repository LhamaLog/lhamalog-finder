package openobserve

import (
	"os"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

func TestMapearDocumentoOpenObserveIntegracao(t *testing.T) {
	if _, habilitado := os.LookupEnv("FINDER_OPENOBSERVE_INTEGRATION"); !habilitado {
		t.Skip("FINDER_OPENOBSERVE_INTEGRATION não está definida")
	}
	prova := carregarProva(t)
	const org = "default"
	const stream = "finder_analise_sintetica"
	batch := "lote-" + unico()
	micros := time.Now().UnixMicro()
	doc := docSintetico(prova, batch, "registro", "tenant-a", "produto-a", "1", micros)
	// O Collector grava camelCase e o OpenObserve retorna colunas minúsculas.
	doc["severityText"] = "WARN"
	doc["body"] = "mensagem sintética"
	doc["serviceName"] = "checkout"
	doc["resourceAttributes"] = `{"host.name":"maquina","service.name":"checkout"}`
	doc["traceId"] = "0123456789abcdef0123456789abcdef"
	ingerir(t, prova, org, stream, []map[string]any{doc})
	ref := referenciaSintetica(org, stream, batch, micros, micros, 1)
	resultado := recuperarVia(t, prova, ref, &coletor{}, classificarPronto(domain.ClasseLoteCompleto))
	if resultado.Classificacao != domain.ClasseLoteCompleto || len(resultado.Registros) != 1 {
		t.Fatalf("recuperação não entregou o documento: %s", resultado.Classificacao)
	}
	log, err := MapearDocumentoAnalise(resultado.Registros[0])
	if err != nil {
		t.Fatal(err)
	}
	if log.Apontamento != "registro" || log.Timestamp.UnixMicro() != micros || log.NivelCampo != "WARN" || log.Host != "maquina" || log.ServiceName != "checkout" || log.Mensagem != doc["body"] || log.TraceID != doc["traceId"] {
		t.Fatal("o documento persistido não manteve os campos esperados pelo mapeamento")
	}
}
