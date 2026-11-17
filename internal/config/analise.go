package config

import "os"

const (
	EnvModuloOcorrenciaURL = "FINDER_MODULO_OCORRENCIA_URL"
	EnvModuloOcorrenciaJWT = "FINDER_MODULO_OCORRENCIA_JWT"
	EnvVersaoRegras        = "FINDER_VERSAO_REGRAS"
)

type Analise struct{ URLModulo, JWTModulo, VersaoRegras string }

// CarregarAnalise não inventa host, rota, claims JWT nem versão.
func CarregarAnalise() Analise {
	return Analise{URLModulo: os.Getenv(EnvModuloOcorrenciaURL), JWTModulo: os.Getenv(EnvModuloOcorrenciaJWT), VersaoRegras: os.Getenv(EnvVersaoRegras)}
}
