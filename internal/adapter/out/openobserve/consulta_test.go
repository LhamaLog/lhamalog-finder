package openobserve

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

func TestMontarCorpoDaConsulta(t *testing.T) {
	cfg := domain.ConfiguracaoAcesso{
		Ordem:        "campo_ordem",
		CampoTenant:  "campo_tenant",
		CampoProduto: "campo_produto",
		Timeout:      time.Second,
	}
	ref := domain.Referencia{
		TenantID: "tenant-a", ProductID: "produto-a", Organizacao: "org-a", Stream: "stream_sintetico",
		BatchID: "batch-a", MinTimestampMicros: 100, MaxTimestampMicros: 200, RecordCount: 9,
	}
	corpo, err := Montar(ref, cfg, 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	decodificado := decodificar(t, corpo)
	if decodificado.Query.From != 0 || decodificado.Query.Size != 2 {
		t.Fatalf("from=%d size=%d", decodificado.Query.From, decodificado.Query.Size)
	}
	if decodificado.Query.Size == int(ref.RecordCount) {
		t.Fatal("size usou recordCount")
	}
	if decodificado.Query.StartTime != 100 || decodificado.Query.EndTime != 201 {
		t.Fatalf("faixa %d..%d", decodificado.Query.StartTime, decodificado.Query.EndTime)
	}
	esperado := `SELECT * FROM "stream_sintetico" WHERE "lhamalog_batch_id" = 'batch-a' AND "campo_tenant" = 'tenant-a' AND "campo_produto" = 'produto-a' ORDER BY campo_ordem`
	if decodificado.Query.SQL != esperado {
		t.Fatalf("sql=%s", decodificado.Query.SQL)
	}
	if strings.Contains(string(corpo), "timeout") || strings.Contains(decodificado.Query.SQL, "_timestamp") {
		t.Fatal("corpo gravou timeout ou ordem inventada")
	}
	var raiz map[string]any
	if err := json.Unmarshal(corpo, &raiz); err != nil {
		t.Fatal(err)
	}
	query := raiz["query"].(map[string]any)
	if _, ok := query["timeout"]; ok {
		t.Fatal("campo timeout presente")
	}

	id := "rec-1"
	ref.RecordID = &id
	corpo, err = Montar(ref, cfg, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	comID := decodificar(t, corpo)
	if comID.Query.From != 2 || !strings.Contains(comID.Query.SQL, `"lhamalog_record_id" = 'rec-1'`) {
		t.Fatalf("from=%d sql=%s", comID.Query.From, comID.Query.SQL)
	}

	ref.RecordID = nil
	ref.BatchID = "42"
	corpo, err = Montar(ref, cfg, 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	inteiro := decodificar(t, corpo)
	if !strings.Contains(inteiro.Query.SQL, `"lhamalog_batch_id" = 42`) || strings.Contains(inteiro.Query.SQL, `'42'`) {
		t.Fatalf("inteiro: %s", inteiro.Query.SQL)
	}

	ref.BatchID = "a'b"
	corpo, err = Montar(ref, cfg, 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(decodificar(t, corpo).Query.SQL, `'a''b'`) {
		t.Fatal("aspas simples não foi repetida")
	}

	ref.BatchID = "batch-a"
	ref.Stream = `logs" quebrado`
	if _, err = Montar(ref, cfg, 0, 2); err == nil {
		t.Fatal("valor inaplicável montou corpo")
	}
}

func decodificar(t *testing.T, corpo []byte) corpoConsulta {
	t.Helper()
	var saida corpoConsulta
	if err := json.Unmarshal(corpo, &saida); err != nil {
		t.Fatal(err)
	}
	return saida
}
