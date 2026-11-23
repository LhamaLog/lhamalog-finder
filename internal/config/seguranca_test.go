package config

import (
	"context"
	"testing"
)

func TestInicializacaoSegurancaAntesDaProntidao(t *testing.T) {
	t.Setenv("FINDER_CATALOGO_SEGURANCA", "../../config/seguranca/catalogo.json")
	for _, v := range []string{"", "outra"} {
		t.Setenv("FINDER_VERSAO_REGRAS", v)
		if _, e := InicializarSeguranca(context.Background()); e == nil {
			t.Fatal("inicializou versão ausente/divergente")
		}
	}
	t.Setenv("FINDER_VERSAO_REGRAS", "finder-2026-11-22.1")
	s, e := InicializarSeguranca(context.Background())
	if e != nil || s.Catalogo.Transformar == nil {
		t.Fatalf("catálogo não executável: %v", e)
	}
	t.Setenv("FINDER_CATALOGO_SEGURANCA", "ausente")
	if _, e := InicializarSeguranca(context.Background()); e == nil {
		t.Fatal("arquivo ausente aceito")
	}
}
