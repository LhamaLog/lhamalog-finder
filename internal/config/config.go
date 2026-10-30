// Package config carrega a configuração por variáveis de ambiente.
package config

import (
	"os"
	"strconv"
)

// Config reúne as configurações do Finder.
type Config struct {
	HTTPAddr     string
	Workers      int
	RedisAddr    string
	RedisStream  string
	RedisGroup   string
	RedisCodigos string
}

// Carregar lê as variáveis de ambiente com valores padrão.
// RedisAddr usa localhost:6379, o endereço já usado pelo processo e pelo Compose.
// Stream, grupo e o conjunto dos códigos concluídos não têm padrão: não há nome definido para eles.
func Carregar() Config {
	c := Config{
		HTTPAddr:     valor("FINDER_HTTP_ADDR", ":8085"),
		Workers:      4,
		RedisAddr:    valor("FINDER_REDIS_ADDR", "localhost:6379"),
		RedisStream:  os.Getenv("FINDER_REDIS_STREAM"),
		RedisGroup:   os.Getenv("FINDER_REDIS_GROUP"),
		RedisCodigos: os.Getenv("FINDER_REDIS_CODIGOS"),
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
