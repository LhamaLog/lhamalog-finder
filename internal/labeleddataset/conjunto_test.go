package labeleddataset

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNomesJSONExatosEmRaizEMembros(t *testing.T) {
	for _, campo := range []string{"itemType", "level_text", "tenantId", "class_label"} {
		entrada := strings.Replace(logValido, `"`+campo+`":`, `"`+strings.ToUpper(campo)+`":`, 1)
		if _, err := LoadConjunto(strings.NewReader(entrada)); err == nil {
			t.Fatalf("alias de %s aceito", campo)
		}
	}
	member := strings.Replace(logValido, `,"class_label":"ERROR"`, "", 1)
	member = strings.Replace(member, `"tenantId"`, `"TENANTID"`, 1)
	ap := `{"itemType":"apontamento","class_label":"ERROR","lhamalog_batch_id":"teste-lote","tenantId":"teste-tenant","productId":"teste-produto","documentacao":"cenario agregado","members":[` + member + `]}`
	if _, err := LoadConjunto(strings.NewReader(ap)); err == nil {
		t.Fatal("alias em member aceito")
	}
}

func TestConflitosEntreArquivosExigemJustificativa(t *testing.T) {
	dir := t.TempDir()
	outro := strings.Replace(logValido, `"class_label":"ERROR"`, `"class_label":"WARNING"`, 1)
	for nome, valor := range map[string]string{"a.jsonl": logValido, "b.jsonl": outro} {
		if err := os.WriteFile(filepath.Join(dir, nome), []byte(valor), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := LoadDiretorio(dir); err == nil {
		t.Fatal("conflito entre arquivos aceito")
	}
	if _, err := LoadDiretorio(dir, DocumentacaoConjunto{Conflitos: " "}); err == nil {
		t.Fatal("justificativa vazia aceita")
	}
	doc := DocumentacaoConjunto{Conflitos: "fixture: mesma entrada com gold alternativo para validar a matriz de confusao"}
	rs, err := LoadDiretorio(dir, doc)
	if err != nil || len(rs) != 2 {
		t.Fatalf("conflito documentado rejeitado: %v", err)
	}
	if _, err := Avaliar(rs); err == nil {
		t.Fatal("avaliacao direta ignorou conflito")
	}
	if _, err := Avaliar(rs, doc); err != nil {
		t.Fatal(err)
	}
}

func TestNotasJustificamConflitoEmQualquerOrdem(t *testing.T) {
	outro := strings.Replace(logValido, `"class_label":"ERROR"`, `"class_label":"WARNING","notas":"fixture de divergencia intencional"`, 1)
	for _, entrada := range []string{logValido + "\n" + outro, outro + "\n" + logValido} {
		if _, err := LoadConjunto(strings.NewReader(entrada)); err != nil {
			t.Fatal(err)
		}
	}
}

func TestConflitosNaoMisturamIdentidadesOuConteudos(t *testing.T) {
	outro := strings.Replace(logValido, `"class_label":"ERROR"`, `"class_label":"WARNING"`, 1)
	for _, entrada := range []string{
		strings.Replace(outro, "teste-tenant", "outro-tenant", 1),
		strings.Replace(outro, "teste-produto", "outro-produto", 1),
		strings.Replace(outro, "falha de conexao", "outra mensagem", 1),
	} {
		if _, err := LoadConjunto(strings.NewReader(logValido + "\n" + entrada)); err != nil {
			t.Fatal(err)
		}
	}
}
