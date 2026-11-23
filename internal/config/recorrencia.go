package config

import "os"

const (
	EnvHistoricoOrganizacao = "FINDER_HISTORICO_OPENOBSERVE_ORGANIZACAO"
	EnvHistoricoStream      = "FINDER_HISTORICO_OPENOBSERVE_STREAM"
	EnvHistoricoUsuario     = "FINDER_HISTORICO_OPENOBSERVE_USUARIO"
	EnvHistoricoSenha       = "FINDER_HISTORICO_OPENOBSERVE_SENHA"
)

type Historico struct{ Organizacao, Stream, Usuario, Senha string }

// CarregarHistorico fornece apenas parametros explicitos; nao define query,
// indice, stream ou credencial operacional. O adapter HTTP historico esta pendente.
func CarregarHistorico() Historico {
	return Historico{Organizacao: os.Getenv(EnvHistoricoOrganizacao), Stream: os.Getenv(EnvHistoricoStream), Usuario: os.Getenv(EnvHistoricoUsuario), Senha: os.Getenv(EnvHistoricoSenha)}
}
