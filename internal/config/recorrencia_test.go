package config

import "testing"

func TestConfiguracaoHistoricoSemDefaults(t *testing.T) {
	for _, nome := range []string{EnvHistoricoOrganizacao, EnvHistoricoStream, EnvHistoricoUsuario, EnvHistoricoSenha} {
		t.Setenv(nome, "")
	}
	if c := CarregarHistorico(); c != (Historico{}) {
		t.Fatal("configuracao historica inventada")
	}
	t.Setenv(EnvHistoricoOrganizacao, "org-teste")
	t.Setenv(EnvHistoricoStream, "stream-teste")
	t.Setenv(EnvHistoricoUsuario, "usuario-teste")
	t.Setenv(EnvHistoricoSenha, "credencial-teste")
	c := CarregarHistorico()
	if c.Organizacao != "org-teste" || c.Stream != "stream-teste" || c.Usuario != "usuario-teste" || c.Senha != "credencial-teste" {
		t.Fatal("nao leu parametros explicitos")
	}
}
