package labeleddataset

import (
	"strings"
	"testing"
)

const logValido = `{"itemType":"log","level_text":"ERROR","message_text":"falha de conexao","lhamalog_batch_id":"teste-lote","lhamalog_record_id":"teste-registro","tenantId":"teste-tenant","productId":"teste-produto","class_label":"ERROR","split":"validacao"}`

func TestLoadValidaCemPorCentoERejeitaCampos(t *testing.T) {
	r, err := LoadConjunto(strings.NewReader(logValido + "\n" + logValido))
	if err != nil || len(r) != 2 {
		t.Fatalf("carga: %v", err)
	}
	for _, entrada := range []string{
		strings.Replace(logValido, `"ERROR","split"`, `"SPAM","split"`, 1),
		strings.Replace(logValido, `"level_text":"ERROR"`, `"level_text":" "`, 1),
		strings.Replace(logValido, `"tenantId":"teste-tenant"`, `"tenantId":""`, 1),
		`null`, `{"itemType":"outro"}`, logValido + ` {}`,
	} {
		if _, err := LoadConjunto(strings.NewReader(entrada)); err == nil {
			t.Fatal("registro invalido aceito")
		}
	}
	grande := strings.Replace(logValido, "falha de conexao", strings.Repeat("x", 100000), 1)
	if _, err := LoadConjunto(strings.NewReader(grande)); err != nil {
		t.Fatal("loader criou teto de linha")
	}
}
func TestMembroNuncaTemGoldMesmoNulo(t *testing.T) {
	member := strings.Replace(logValido, `,"class_label":"ERROR"`, "", 1)
	apontamento := `{"itemType":"apontamento","class_label":"ERROR","lhamalog_batch_id":"teste-lote","tenantId":"teste-tenant","productId":"teste-produto","documentacao":"cenario agregado","members":[` + member + `]}`
	r, err := LoadConjunto(strings.NewReader(apontamento))
	if err != nil || len(r) != 1 {
		t.Fatalf("apontamento: %v", err)
	}
	for _, gold := range []string{`"ERROR"`, `""`, `null`} {
		bad := strings.Replace(apontamento, `"members":[{`, `"members":[{"class_label":`+gold+`,`, 1)
		if _, err := LoadConjunto(strings.NewReader(bad)); err == nil {
			t.Fatal("gold em member aceito")
		}
	}
}
func TestPIIHeuristicaRejeitaSemPublicarBase(t *testing.T) {
	for _, mensagem := range []string{"usuario@" + "example.test", "token=" + "valor-teste"} {
		if _, err := LoadConjunto(strings.NewReader(strings.Replace(logValido, "falha de conexao", mensagem, 1))); err == nil {
			t.Fatal("heuristica de privacidade nao rejeitou")
		}
	}
}

func TestGoldDuplicadoNaoPodeEscolherUltimoValor(t *testing.T) {
	entrada := strings.Replace(logValido, `"class_label":"ERROR"`, `"class_label":"WARNING","class_label":"ERROR"`, 1)
	if _, err := LoadConjunto(strings.NewReader(entrada)); err == nil {
		t.Fatal("gold duplicado aceito")
	}
}
