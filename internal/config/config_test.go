package config

import "testing"

func TestCarregarPadrao(t *testing.T) {
	t.Setenv("FINDER_HTTP_ADDR", "")
	t.Setenv("FINDER_WORKERS", "")
	t.Setenv("FINDER_REDIS_ADDR", "")
	t.Setenv("FINDER_REDIS_STREAM", "")
	t.Setenv("FINDER_REDIS_GROUP", "")
	t.Setenv("FINDER_REDIS_CODIGOS", "")
	t.Setenv("FINDER_OPENOBSERVE_URL", "")
	t.Setenv("FINDER_OPENOBSERVE_ORG", "")
	t.Setenv("FINDER_OPENOBSERVE_USER", "")
	t.Setenv("FINDER_OPENOBSERVE_PASSWORD", "")
	c := Carregar()
	if c.HTTPAddr != ":8085" || c.Workers != 4 || c.RedisAddr != "localhost:6379" {
		t.Fatalf("config inesperada: %+v", c)
	}
	if c.OpenObserveURL != "http://127.0.0.1:5080" || c.OpenObserveOrg != "default" || c.OpenObserveUser != "root@example.com" {
		t.Fatalf("openobserve: %+v", c)
	}
	if c.RedisStream != "" || c.RedisGroup != "" || c.RedisCodigos != "" {
		t.Fatalf("nomes da fila não podem ter padrão: %+v", c)
	}
}

func TestCarregarWorkers(t *testing.T) {
	t.Setenv("FINDER_WORKERS", "8")
	if c := Carregar(); c.Workers != 8 {
		t.Fatalf("workers = %d", c.Workers)
	}
}
