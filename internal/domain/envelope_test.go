package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEnvelopeValidoComESemTrace(t *testing.T) {
	um := envelope(t, 1, "rec-1", "trace-1", true)
	ref, err := InterpretarEnvelope(um)
	if err != nil {
		t.Fatal(err)
	}
	conferirReferencia(t, ref, 1, "rec-1", "trace-1")
	if ref.ReceivedAt == nil || *ref.ReceivedAt != "2026-11-08T12:00:00Z" || ref.Bytes == nil || *ref.Bytes != 32 {
		t.Fatalf("opcionais: %+v", ref)
	}
	if len(ref.ServiceNames) != 1 || ref.ServiceNames[0] != "api" || ref.Priority == nil || *ref.Priority != "qualquer" {
		t.Fatalf("agregados: %+v", ref)
	}
	if ref.SeverityCounts["info"] != 2 || ref.MaxSeverityNumber == nil || *ref.MaxSeverityNumber != 9 {
		t.Fatalf("severidade: %+v", ref)
	}

	semTrace := envelope(t, 1, "rec-1", "", false)
	ref, err = InterpretarEnvelope(semTrace)
	if err != nil {
		t.Fatal(err)
	}
	if ref.TraceID != nil || ref.IdempotencyID != "batch-1" {
		t.Fatalf("sem trace: %+v", ref)
	}

	varios := envelope(t, 4, nil, "", false)
	ref, err = InterpretarEnvelope(varios)
	if err != nil {
		t.Fatal(err)
	}
	if ref.RecordID != nil || ref.RecordCount != 4 {
		t.Fatalf("lote múltiplo: %+v", ref)
	}

	outroTrace := envelope(t, 1, "rec-1", "trace-2", true)
	segunda, err := InterpretarEnvelope(outroTrace)
	if err != nil {
		t.Fatal(err)
	}
	if ref.IdempotencyID != segunda.IdempotencyID || segunda.TraceID == nil || *segunda.TraceID != "trace-2" {
		t.Fatal("traceId passou a deduplicar")
	}
}

func TestEnvelopeOpcionalAusenteContinuaValido(t *testing.T) {
	campos := envelope(t, 1, "rec-1", "", false)
	body := map[string]any{}
	if err := json.Unmarshal([]byte(campos["body"]), &body); err != nil {
		t.Fatal(err)
	}
	for _, campo := range []string{"receivedAt", "bytes", "serviceNames", "priority", "criticality", "maxSeverityNumber", "severityCounts"} {
		delete(body, campo)
	}
	bruto, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	campos["body"] = string(bruto)
	ref, err := InterpretarEnvelope(campos)
	if err != nil {
		t.Fatal(err)
	}
	if ref.ReceivedAt != nil || ref.Bytes != nil || ref.ServiceNames != nil || ref.Priority != nil || ref.Criticality != nil || ref.MaxSeverityNumber != nil || ref.SeverityCounts != nil {
		t.Fatalf("ausência não foi só omissão: %+v", ref)
	}
	if ref.TenantID != "tenant-1" || ref.ProductID != "prod-1" {
		t.Fatalf("identidade: %+v", ref)
	}
}

func TestEnvelopeInvalido(t *testing.T) {
	casos := []struct {
		nome   string
		mutar  func(map[string]string)
		motivo string
	}{
		{nome: "json metadata", mutar: func(c map[string]string) { c["metadata"] = "CONTEUDO-DE-LOG-SECRETO" }, motivo: "metadata não é JSON válido"},
		{nome: "json body", mutar: func(c map[string]string) { c["body"] = `senha=segredo-redis CONTEUDO-DE-LOG-SECRETO` }, motivo: "body não é JSON válido"},
		{nome: "tenant ausente", mutar: func(c map[string]string) { c["metadata"] = trocarJSON(t, c["metadata"], "tenantId", nil) }, motivo: "tenantId ausente"},
		{nome: "event", mutar: func(c map[string]string) { c["event"] = "trace.ingested" }, motivo: "event diferente de log.ingested"},
		{nome: "origin", mutar: func(c map[string]string) { c["origin"] = "api" }, motivo: "origin diferente de collector"},
		{nome: "schema 2", mutar: func(c map[string]string) { c["schemaVersion"] = "2" }, motivo: "schemaVersion diferente de 1"},
		{nome: "signal", mutar: func(c map[string]string) { c["body"] = trocarJSON(t, c["body"], "signal", "TRACES") }, motivo: "signal diferente de LOGS"},
		{nome: "productId divergente", mutar: func(c map[string]string) { c["body"] = trocarJSON(t, c["body"], "productId", "outro") }, motivo: "productId divergente"},
		{nome: "batch diferente", mutar: func(c map[string]string) {
			c["body"] = trocarAberto(t, c["body"], "batchId", "outro-batch")
		}, motivo: "idempotencyId diferente de batchId"},
		{nome: "faixa", mutar: func(c map[string]string) {
			c["body"] = trocarAberto(t, c["body"], "minTimestampMicros", 500)
		}, motivo: "mínimo da faixa maior que o máximo"},
		{nome: "recordCount zero", mutar: func(c map[string]string) { c["body"] = trocarJSON(t, c["body"], "recordCount", 0) }, motivo: "recordCount não positivo"},
		{nome: "recordCount negativo", mutar: func(c map[string]string) { c["body"] = trocarJSON(t, c["body"], "recordCount", -2) }, motivo: "recordCount não positivo"},
		{nome: "um sem recordId", mutar: func(c map[string]string) { c["body"] = trocarAberto(t, c["body"], "recordId", nil) }, motivo: "lote de um registro sem recordId"},
		{nome: "varios com recordId", mutar: func(c map[string]string) {
			c["body"] = trocarJSON(t, c["body"], "recordCount", 3)
			c["body"] = trocarAberto(t, c["body"], "recordId", "rec-1")
		}, motivo: "lote de vários registros com recordId"},
		{nome: "publishedAt", mutar: func(c map[string]string) { c["publishedAt"] = "2026-11-08T12:00:00-03:00" }, motivo: "publishedAt não é instante ISO UTC"},
		{nome: "publishedAt sem zona", mutar: func(c map[string]string) { c["publishedAt"] = "2026-11-08T12:00:00" }, motivo: "publishedAt não é instante ISO UTC"},
	}
	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			campos := envelope(t, 1, "rec-1", "trace-1", true)
			caso.mutar(campos)
			_, err := InterpretarEnvelope(campos)
			if err == nil {
				t.Fatal("aceitou referência inválida")
			}
			erro, ok := err.(*ErroEnvelope)
			if !ok {
				t.Fatalf("erro = %T", err)
			}
			if erro.Motivo != caso.motivo {
				t.Fatalf("motivo = %q", erro.Motivo)
			}
			if strings.Contains(erro.Motivo, "CONTEUDO-DE-LOG-SECRETO") || strings.Contains(erro.Motivo, "segredo-redis") || strings.Contains(erro.Motivo, campos["metadata"]) || strings.Contains(erro.Motivo, campos["body"]) {
				t.Fatalf("motivo contém texto bruto: %s", erro.Motivo)
			}
			if caso.nome != "json metadata" && caso.nome != "event" && caso.nome != "origin" && caso.nome != "schema 2" && caso.nome != "publishedAt" && caso.nome != "publishedAt sem zona" {
				if erro.IdempotencyID != "batch-1" {
					t.Fatalf("idempotencyId = %q", erro.IdempotencyID)
				}
			}
		})
	}
}

func conferirReferencia(t *testing.T, ref Referencia, contagem int64, recordID, trace string) {
	t.Helper()
	if ref.TenantID != "tenant-1" || ref.ProductID != "prod-1" || ref.Organizacao != "org-1" || ref.Stream != "logs" {
		t.Fatalf("consulta: %+v", ref)
	}
	if ref.BatchID != "batch-1" || ref.IdempotencyID != "batch-1" || ref.CorrelationID != "corr-1" || ref.RecordCount != contagem {
		t.Fatalf("lote: %+v", ref)
	}
	if ref.MinTimestampMicros != 100 || ref.MaxTimestampMicros != 200 {
		t.Fatalf("faixa: %+v", ref)
	}
	if ref.RecordID == nil || *ref.RecordID != recordID {
		t.Fatalf("recordId: %+v", ref.RecordID)
	}
	if ref.TraceID == nil || *ref.TraceID != trace {
		t.Fatalf("trace: %+v", ref.TraceID)
	}
}

func envelope(t *testing.T, contagem int, recordID any, trace string, comTrace bool) map[string]string {
	t.Helper()
	meta := map[string]any{
		"tenantId":      "tenant-1",
		"productId":     "prod-1",
		"correlationId": "corr-1",
	}
	if comTrace {
		meta["traceId"] = trace
	}
	body := map[string]any{
		"signal":    "LOGS",
		"productId": "prod-1",
		"openobserve": map[string]any{
			"organization":       "org-1",
			"stream":             "logs",
			"batchId":            "batch-1",
			"recordId":           recordID,
			"minTimestampMicros": 100,
			"maxTimestampMicros": 200,
		},
		"recordCount":       contagem,
		"receivedAt":        "2026-11-08T12:00:00Z",
		"bytes":             32,
		"serviceNames":      []string{"api"},
		"priority":          "qualquer",
		"criticality":       "alta",
		"maxSeverityNumber": 9,
		"severityCounts":    map[string]int{"info": 2},
	}
	return map[string]string{
		"idempotencyId": "batch-1",
		"event":         "log.ingested",
		"origin":        "collector",
		"schemaVersion": "1",
		"publishedAt":   "2026-11-08T12:00:00.123Z",
		"metadata":      string(deveJSON(t, meta)),
		"body":          string(deveJSON(t, body)),
	}
}

func trocarJSON(t *testing.T, bruto, campo string, valor any) string {
	t.Helper()
	var objeto map[string]any
	if err := json.Unmarshal([]byte(bruto), &objeto); err != nil {
		t.Fatal(err)
	}
	if valor == nil {
		delete(objeto, campo)
	} else {
		objeto[campo] = valor
	}
	return string(deveJSON(t, objeto))
}

func trocarAberto(t *testing.T, bruto, campo string, valor any) string {
	t.Helper()
	var objeto map[string]any
	if err := json.Unmarshal([]byte(bruto), &objeto); err != nil {
		t.Fatal(err)
	}
	aberto := objeto["openobserve"].(map[string]any)
	if valor == nil {
		delete(aberto, campo)
	} else {
		aberto[campo] = valor
	}
	return string(deveJSON(t, objeto))
}

func deveJSON(t *testing.T, valor any) []byte {
	t.Helper()
	bruto, err := json.Marshal(valor)
	if err != nil {
		t.Fatal(err)
	}
	return bruto
}
