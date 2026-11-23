package openobserve

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	entrada "github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

var (
	errDocumentoAnalise     = errors.New("documento_openobserve_invalido")
	errIdentidadeAssinatura = errors.New("assinatura_incompleta")
)

// MapearDocumentoAnalise usa os nomes de coluna do documento de logs do Collector.
// O chamador entrega apenas documentos de um lote completo, já isolado e deduplicado.
func MapearDocumentoAnalise(doc domain.DocumentoFisico) (entrada.LogAnalise, error) {
	id, ok := domain.IdentificadorRegistro(doc)
	if !ok || strings.TrimSpace(id) == "" {
		return entrada.LogAnalise{}, errDocumentoAnalise
	}
	// O cliente OpenObserve usa Decoder.UseNumber para preservar os microssegundos.
	numero, ok := doc["_timestamp"].(json.Number)
	if !ok {
		return entrada.LogAnalise{}, errDocumentoAnalise
	}
	micros, err := numero.Int64()
	if err != nil {
		return entrada.LogAnalise{}, errDocumentoAnalise
	}

	servico, ok := textoColuna(doc, "servicename")
	if !ok || strings.TrimSpace(servico) == "" {
		return entrada.LogAnalise{}, errIdentidadeAssinatura
	}
	atributos, ok := textoColuna(doc, "resourceattributes")
	if !ok {
		return entrada.LogAnalise{}, errDocumentoAnalise
	}
	if atributos == "" {
		return entrada.LogAnalise{}, errIdentidadeAssinatura
	}
	var recurso map[string]json.RawMessage
	if err := json.Unmarshal([]byte(atributos), &recurso); err != nil || recurso == nil {
		return entrada.LogAnalise{}, errDocumentoAnalise
	}
	var host string
	if err := json.Unmarshal(recurso["host.name"], &host); err != nil || strings.TrimSpace(host) == "" {
		return entrada.LogAnalise{}, errIdentidadeAssinatura
	}
	mensagem, corpoValido := textoColuna(doc, "body")
	nivel, nivelValido := textoColuna(doc, "severitytext")
	trace, traceValido := textoColuna(doc, "traceid")
	if !corpoValido || !nivelValido || !traceValido {
		return entrada.LogAnalise{}, errDocumentoAnalise
	}
	return entrada.LogAnalise{
		Apontamento: domain.Codigo(id),
		Timestamp:   time.UnixMicro(micros).UTC(),
		NivelCampo:  nivel,
		Mensagem:    mensagem,
		ServiceName: servico,
		Host:        host,
		TraceID:     trace,
	}, nil
}

func textoColuna(doc domain.DocumentoFisico, nome string) (string, bool) {
	if doc[nome] == nil {
		return "", true
	}
	valor, ok := doc[nome].(string)
	return valor, ok
}
