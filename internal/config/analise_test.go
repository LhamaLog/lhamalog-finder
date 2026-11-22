package config

import "testing"

func TestAnaliseSemPadroesInventados(t *testing.T) {
	for _, k := range []string{"FINDER_MODULO_OCORRENCIA_URL", "FINDER_MODULO_OCORRENCIA_JWT", "FINDER_VERSAO_REGRAS"} {
		t.Setenv(k, "")
	}
	c := CarregarAnalise()
	if c.URLModulo != "" || c.JWTModulo != "" || c.VersaoRegras != "" {
		t.Fatal("config criou padrao")
	}
	t.Setenv("FINDER_VERSAO_REGRAS", "regra-v1")
	if CarregarAnalise().VersaoRegras != "regra-v1" {
		t.Fatal("versao ignorada")
	}
}
