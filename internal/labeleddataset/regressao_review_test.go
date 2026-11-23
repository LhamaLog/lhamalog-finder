package labeleddataset_test

import (
	ds "github.com/LhamaLog/lhamalog-finder/internal/labeleddataset"
	"strings"
	"testing"
)

const reviewLog = `{"itemType":"log","level_text":"ERROR","message_text":"falha","lhamalog_batch_id":"teste-b","lhamalog_record_id":"teste-r","tenantId":"teste-t","productId":"teste-p","class_label":"ERROR","split":"validacao"}`

func TestReviewGoldDuplicadoComCase(t *testing.T) {
	raw := strings.Replace(reviewLog, `"class_label":"ERROR"`, `"class_label":"WARNING","CLASS_LABEL":"ERROR"`, 1)
	r, err := ds.LoadConjunto(strings.NewReader(raw))
	if err == nil {
		t.Fatalf("gold duplicado aceito, resultado final %q", r[0].ClassLabel)
	}
}
func TestReviewIdentidadeObrigatoriaNoApontamento(t *testing.T) {
	member := strings.Replace(reviewLog, `,"class_label":"ERROR"`, "", 1)
	raw := `{"itemType":"apontamento","class_label":"ERROR","lhamalog_batch_id":"teste-b","documentacao":"cenario","members":[` + member + `]}`
	_, err := ds.LoadConjunto(strings.NewReader(raw))
	if err == nil {
		t.Fatal("apontamento sem tenantId e productId obrigatorios foi aceito")
	}
}
func TestReviewConflitoGoldSemDocumentacao(t *testing.T) {
	other := strings.Replace(reviewLog, `"class_label":"ERROR"`, `"class_label":"WARNING"`, 1)
	_, err := ds.LoadConjunto(strings.NewReader(reviewLog + "\n" + other))
	if err == nil {
		t.Fatal("mesmo registro com gold conflitante sem notas foi aceito")
	}
}
