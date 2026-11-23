package domain

import "testing"

func TestClassificarLogContrato(t *testing.T) {
	casos := []struct {
		nivel, msg string
		classe     ClasseLog
	}{
		{"error", "mensagem", ERROR}, {"erro", "mensagem", ERROR}, {"fatal", "info", ERROR},
		{"warning", "mensagem", WARNING}, {"warn", "mensagem", WARNING}, {"alerta", "mensagem", WARNING},
		{"info", "mensagem", INFO}, {"debug", "mensagem", INFO}, {"trace", "mensagem", INFO},
		{"info", "ERROR ao conectar", ERROR}, {"info", "errorless", INFO}, {"", "traceback", INFO},
		{"", "(warning)", WARNING}, {"", "pre_error_pos", ERROR}, {"", "terror", INFO}, {"", "error2", INFO},
		{"", "erro! alerta", ERROR}, {"", "informação", INFO},
	}
	for _, c := range casos {
		r := ClassificarLog(c.nivel, c.msg, nil)
		if r.Classe != c.classe || r.CatalogoUsado {
			t.Errorf("%q %q: %+v", c.nivel, c.msg, r)
		}
	}
}
func TestCatalogoOpacoSomenteSemNivel(t *testing.T) {
	chamadas := 0
	regra := func(string) bool { chamadas++; return true }
	cat := []RegraClassificacao{{Classe: "SPAM", Corresponde: regra}, {Classe: WARNING, Corresponde: regra}, {Classe: ERROR, Corresponde: regra}}
	r := ClassificarLog("info", "x", cat)
	if r.Classe != INFO || chamadas != 0 {
		t.Fatal("catalogo venceu nivel")
	}
	r = ClassificarLog("", "x", cat)
	if r.Classe != WARNING || !r.CatalogoUsado || chamadas != 1 {
		t.Fatal("primeira valida nao decidiu")
	}
}
func TestClasseValidaEOrdem(t *testing.T) {
	if !ERROR.Valida() || !WARNING.Valida() || !INFO.Valida() || ClasseLog("SPAM").Valida() {
		t.Fatal("taxonomia")
	}
	if MaisGrave(INFO, WARNING) != WARNING || MaisGrave(ERROR, WARNING) != ERROR {
		t.Fatal("ordem")
	}
}
