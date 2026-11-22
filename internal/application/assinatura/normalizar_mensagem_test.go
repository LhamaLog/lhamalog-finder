package assinatura

import "testing"

func TestNormalizarMarcadoresCanonicos(t *testing.T) {
	casos := []struct{ entrada, saida string }{
		{"  FALHA 123  ", "falha <NUM>"},
		{"id 550e8400-e29b-41d4-a716-446655440000", "id <UUID>"},
		{"falha ao conectar 192.0.2.10", "falha ao conectar <IP>"},
		{"falha em 2026-11-16", "falha em <DATE>"},
		{"valor 0xDEADBEEF", "valor <HEX>"},
		{"a <NUM> <UUID> <IP> <DATE> <HEX>", "a <NUM> <UUID> <IP> <DATE> <HEX>"},
	}
	for _, c := range casos {
		if got := NormalizarMensagem(c.entrada); got != c.saida {
			t.Errorf("normalizar: %q != %q", got, c.saida)
		}
	}
}
func TestComposicaoCanonicaEEscape(t *testing.T) {
	got, err := Calcular("ERROR", "  FALHA 42 | rede ", "orders|api", "host-a")
	if err != nil || got != `ERROR|falha <NUM> \| rede|orders\|api|host-a` {
		t.Fatalf("assinatura=%q erro=%v", got, err)
	}
	for _, v := range [][4]string{{"ERROR", "x", "", "h"}, {"ERROR", "x", "s", ""}} {
		if _, err := Calcular(v[0], v[1], v[2], v[3]); err == nil {
			t.Fatal("assinatura incompleta aceita")
		}
	}
}
func TestNormalizacaoDeterministicaIdempotente(t *testing.T) {
	a := NormalizarMensagem(" FALHA 12 em 192.0.2.1 ")
	b := NormalizarMensagem("falha 99 em 198.51.100.5")
	if a != b || NormalizarMensagem(a) != a {
		t.Fatalf("nao canonico: %q %q", a, b)
	}
}
