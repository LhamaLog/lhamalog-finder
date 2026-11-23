package assinatura

import "testing"

func TestNivelOpcionalPreservaPosicaoNaAssinatura(t *testing.T) {
	valor, err := Calcular("", "processado", "svc", "h")
	if err != nil || valor != "|processado|svc|h" {
		t.Fatalf("assinatura=%q erro=%v", valor, err)
	}
}
