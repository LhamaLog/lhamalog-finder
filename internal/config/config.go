// Package config carrega a configuração por variáveis de ambiente.
package config

import (
	"os"
	"strconv"
)

// Config reúne as configurações do Finder.
type Config struct {
	HTTPAddr string
	Workers  int
}

// Carregar lê as variáveis de ambiente com valores padrão.
func Carregar() Config {
	c := Config{HTTPAddr: valor("FINDER_HTTP_ADDR", ":8085"), Workers: 4}
	if n, err := strconv.Atoi(os.Getenv("FINDER_WORKERS")); err == nil && n > 0 {
		c.Workers = n
	}
	return c
}

func valor(chave, padrao string) string {
	if v := os.Getenv(chave); v != "" {
		return v
	}
	return padrao
}
