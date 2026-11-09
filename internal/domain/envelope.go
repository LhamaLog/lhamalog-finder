package domain

import (
	"encoding/json"
	"strings"
	"time"
)

// Entrada é uma publicação do stream. O id confirma a entrada e não deduplica o lote.
// Ociosa indica que a ociosidade configurada já autoriza reassunção.
// Reivindicada indica que a entrada veio de outra entrega abandonada.
type Entrada struct {
	ID           string
	Campos       map[string]string
	Ociosa       bool
	Reivindicada bool
}

// Referencia é o recorte entregue à etapa de recuperação.
// recordCount, bytes e os agregados descrevem os registros aceitos pelo Collector.
// A entrega não confirma a entrada e não traz o conteúdo do log.
type Referencia struct {
	TenantID           string
	ProductID          string
	Organizacao        string
	Stream             string
	BatchID            string
	RecordID           *string
	MinTimestampMicros int64
	MaxTimestampMicros int64
	RecordCount        int64
	IdempotencyID      string
	CorrelationID      string
	EntradaID          string
	TraceID            *string
	ReceivedAt         *string
	Bytes              *int64
	ServiceNames       []string
	Priority           *string
	Criticality        *string
	MaxSeverityNumber  *int64
	SeverityCounts     map[string]int64
}

// ErroEnvelope é o motivo de uma referência rejeitada.
// O texto não inclui metadata bruta, body bruto nem conteúdo de log.
type ErroEnvelope struct {
	Motivo        string
	IdempotencyID string
}

func (e *ErroEnvelope) Error() string {
	if e == nil {
		return ""
	}
	return e.Motivo
}

// InterpretarEnvelope lê os campos de texto do envelope log.ingested v1.
// tenantId e productId da referência são os da metadata.
// correlationId e traceId não deduplicam.
func InterpretarEnvelope(campos map[string]string) (Referencia, error) {
	idem := campos["idempotencyId"]
	recusar := func(motivo string) (Referencia, error) {
		return Referencia{}, &ErroEnvelope{Motivo: motivo, IdempotencyID: idem}
	}
	if campos["event"] != "log.ingested" {
		return recusar("event diferente de log.ingested")
	}
	if campos["origin"] != "collector" {
		return recusar("origin diferente de collector")
	}
	if campos["schemaVersion"] != "1" {
		return recusar("schemaVersion diferente de 1")
	}
	if !instanteUTC(campos["publishedAt"]) {
		return recusar("publishedAt não é instante ISO UTC")
	}
	if idem == "" {
		return recusar("idempotencyId ausente")
	}
	meta, err := objetoJSON(campos["metadata"])
	if err != nil {
		return recusar("metadata não é JSON válido")
	}
	body, err := objetoJSON(campos["body"])
	if err != nil {
		return recusar("body não é JSON válido")
	}
	tenant, ok := textoObrigatorio(meta, "tenantId")
	if !ok {
		return recusar("tenantId ausente")
	}
	produtoMeta, ok := textoObrigatorio(meta, "productId")
	if !ok {
		return recusar("productId ausente")
	}
	correlacao, ok := textoObrigatorio(meta, "correlationId")
	if !ok {
		return recusar("correlationId ausente")
	}
	if body["signal"] != "LOGS" {
		return recusar("signal diferente de LOGS")
	}
	produtoBody, ok := textoObrigatorio(body, "productId")
	if !ok {
		return recusar("productId ausente")
	}
	if produtoMeta != produtoBody {
		return recusar("productId divergente")
	}
	aberto, ok := body["openobserve"].(map[string]any)
	if !ok {
		return recusar("openobserve ausente")
	}
	organizacao, ok := textoObrigatorio(aberto, "organization")
	if !ok {
		return recusar("organization ausente")
	}
	stream, ok := textoObrigatorio(aberto, "stream")
	if !ok {
		return recusar("stream ausente")
	}
	batchID, ok := textoObrigatorio(aberto, "batchId")
	if !ok {
		return recusar("batchId ausente")
	}
	if idem != batchID {
		return recusar("idempotencyId diferente de batchId")
	}
	minimo, ok := inteiro(aberto["minTimestampMicros"])
	if !ok {
		return recusar("minTimestampMicros ausente")
	}
	maximo, ok := inteiro(aberto["maxTimestampMicros"])
	if !ok {
		return recusar("maxTimestampMicros ausente")
	}
	if minimo > maximo {
		return recusar("mínimo da faixa maior que o máximo")
	}
	contagem, ok := inteiro(body["recordCount"])
	if !ok {
		return recusar("recordCount ausente")
	}
	if contagem <= 0 {
		return recusar("recordCount não positivo")
	}
	registro, nulo, tipoOK := lerRecordID(aberto)
	if !tipoOK {
		return recusar("recordId ilegível")
	}
	if contagem == 1 && (nulo || registro == "") {
		return recusar("lote de um registro sem recordId")
	}
	if contagem > 1 && !nulo {
		return recusar("lote de vários registros com recordId")
	}
	ref := Referencia{
		TenantID:           tenant,
		ProductID:          produtoMeta,
		Organizacao:        organizacao,
		Stream:             stream,
		BatchID:            batchID,
		MinTimestampMicros: minimo,
		MaxTimestampMicros: maximo,
		RecordCount:        contagem,
		IdempotencyID:      idem,
		CorrelationID:      correlacao,
		TraceID:            textoOpcional(meta, "traceId"),
		ReceivedAt:         textoOpcional(body, "receivedAt"),
		Bytes:              inteiroOpcional(body, "bytes"),
		ServiceNames:       listaTexto(body["serviceNames"]),
		Priority:           textoOpcional(body, "priority"),
		Criticality:        textoOpcional(body, "criticality"),
		MaxSeverityNumber:  inteiroOpcional(body, "maxSeverityNumber"),
		SeverityCounts:     mapaInteiro(body["severityCounts"]),
	}
	if contagem == 1 {
		ref.RecordID = &registro
	}
	return ref, nil
}

func instanteUTC(valor string) bool {
	instante, err := time.Parse(time.RFC3339Nano, valor)
	if err != nil {
		return false
	}
	_, deslocamento := instante.Zone()
	return deslocamento == 0
}

func objetoJSON(texto string) (map[string]any, error) {
	if strings.TrimSpace(texto) == "" {
		return nil, errJSON
	}
	dec := json.NewDecoder(strings.NewReader(texto))
	dec.UseNumber()
	var valor any
	if err := dec.Decode(&valor); err != nil {
		return nil, errJSON
	}
	if dec.More() {
		return nil, errJSON
	}
	objeto, ok := valor.(map[string]any)
	if !ok {
		return nil, errJSON
	}
	return objeto, nil
}

var errJSON = &ErroEnvelope{Motivo: "json inválido"}

func textoObrigatorio(objeto map[string]any, campo string) (string, bool) {
	valor, ok := objeto[campo]
	if !ok || valor == nil {
		return "", false
	}
	texto, ok := valor.(string)
	if !ok || texto == "" {
		return "", false
	}
	return texto, true
}

func textoOpcional(objeto map[string]any, campo string) *string {
	texto, ok := textoObrigatorio(objeto, campo)
	if !ok {
		return nil
	}
	return &texto
}

func inteiro(valor any) (int64, bool) {
	numero, ok := valor.(json.Number)
	if !ok {
		return 0, false
	}
	n, err := numero.Int64()
	if err != nil {
		return 0, false
	}
	return n, true
}

func inteiroOpcional(objeto map[string]any, campo string) *int64 {
	valor, existe := objeto[campo]
	if !existe || valor == nil {
		return nil
	}
	n, ok := inteiro(valor)
	if !ok {
		return nil
	}
	return &n
}

func listaTexto(valor any) []string {
	lista, ok := valor.([]any)
	if !ok {
		return nil
	}
	saida := make([]string, 0, len(lista))
	for _, item := range lista {
		texto, ok := item.(string)
		if !ok {
			return nil
		}
		saida = append(saida, texto)
	}
	return saida
}

func mapaInteiro(valor any) map[string]int64 {
	objeto, ok := valor.(map[string]any)
	if !ok {
		return nil
	}
	saida := make(map[string]int64, len(objeto))
	for chave, item := range objeto {
		n, ok := inteiro(item)
		if !ok {
			return nil
		}
		saida[chave] = n
	}
	return saida
}

// lerRecordID devolve o texto, se o campo é nulo ou ausente, e se o tipo é aceitável.
// String vazia não conta como preenchida nem como nulo.
func lerRecordID(aberto map[string]any) (string, bool, bool) {
	valor, existe := aberto["recordId"]
	if !existe || valor == nil {
		return "", true, true
	}
	texto, ok := valor.(string)
	if !ok || texto == "" {
		return "", false, ok && texto == ""
	}
	return texto, false, true
}
