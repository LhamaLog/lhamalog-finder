package assinatura

import "testing"

func TestEscapeDF20261122(t *testing.T) {
	for _, c := range []struct{ mensagem, servico, host, esperado string }{
		{`falha em c:\logs\app`, `svc`, `h`, `ERROR|falha em c:\\logs\\app|svc|h`},
		{`falha | detalhe`, `s\|v`, `h\`, `ERROR|falha \| detalhe|s\\\|v|h\\`},
		{`falha \| detalhe`, `s`, `h`, `ERROR|falha \\\| detalhe|s|h`},
	} {
		got, err := Calcular("ERROR", c.mensagem, c.servico, c.host)
		if err != nil || got != c.esperado {
			t.Fatalf("entrada=%q esperado=%q obtido=%q erro=%v", c.mensagem, c.esperado, got, err)
		}
	}
	a, _ := Calcular("ERROR", `falha \| detalhe`, "s", "h")
	b, _ := Calcular("ERROR", `falha | detalhe`, "s", "h")
	if a == b {
		t.Fatal("componentes distintos colidiram")
	}
}
