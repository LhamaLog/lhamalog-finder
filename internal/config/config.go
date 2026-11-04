// Package config carrega a configuração por variáveis de ambiente.
package config

import (
	"os"
	"strconv"
)

// Config reúne as configurações do Finder.
type Config struct {
	HTTPAddr            string
	Workers             int
	RedisAddr           string
	RedisStream         string
	RedisGroup          string
	RedisCodigos        string
	OpenObserveURL      string
	OpenObserveOrg      string
	OpenObserveUser     string
	OpenObservePassword string
}

// Carregar lê as variáveis de ambiente com valores padrão.
// RedisAddr usa localhost:6379, o endereço já usado pelo processo e pelo Compose.
// Stream, grupo e o conjunto dos códigos concluídos não têm padrão no código.
// O processo nativo recebe esses nomes de um .env, aplicado antes desta leitura.
func Carregar() Config {
	c := Config{
		HTTPAddr:            valor("FINDER_HTTP_ADDR", ":8085"),
		Workers:             4,
		RedisAddr:           valor("FINDER_REDIS_ADDR", "localhost:6379"),
		RedisStream:         os.Getenv("FINDER_REDIS_STREAM"),
		RedisGroup:          os.Getenv("FINDER_REDIS_GROUP"),
		RedisCodigos:        os.Getenv("FINDER_REDIS_CODIGOS"),
		OpenObserveURL:      valor("FINDER_OPENOBSERVE_URL", "http://127.0.0.1:5080"),
		OpenObserveOrg:      valor("FINDER_OPENOBSERVE_ORG", "default"),
		OpenObserveUser:     valor("FINDER_OPENOBSERVE_USER", "root@example.com"),
		OpenObservePassword: valor("FINDER_OPENOBSERVE_PASSWORD", "Complexpass#123"),
	}
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
