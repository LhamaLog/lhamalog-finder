package labeleddataset_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	ds "github.com/LhamaLog/lhamalog-finder/internal/labeleddataset"
)

func TestReviewParticaoDosArquivosNaoAlteraJustificativaDoConjunto(t *testing.T) {
	a := `{"itemType":"log","level_text":"ERROR","message_text":"falha de conexao","lhamalog_batch_id":"lote-teste","lhamalog_record_id":"registro-teste","tenantId":"tenant-teste","productId":"produto-teste","class_label":"ERROR","split":"validacao"}`
	b := strings.Replace(a, `"class_label":"ERROR"`, `"class_label":"WARNING"`, 1)
	c := strings.Replace(a, `"class_label":"ERROR"`, `"class_label":"ERROR","notas":"fixture de gold alternativo intencional para o mesmo registro"`, 1)
	registros, err := ds.LoadConjunto(strings.NewReader(a + "\n" + b + "\n" + c))
	if err != nil || len(registros) != 3 {
		t.Fatalf("controle em arquivo unico: %v", err)
	}
	if err := ds.ValidarConjunto(registros); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	for nome, conteudo := range map[string]string{"a.jsonl": a + "\n" + b, "b.jsonl": c} {
		if err := os.WriteFile(filepath.Join(dir, nome), []byte(conteudo), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := ds.LoadDiretorio(dir); err != nil {
		t.Fatalf("a mesma base aceita em um arquivo foi rejeitada ao dividir arquivos com justificativa no segundo: %v", err)
	}
}
