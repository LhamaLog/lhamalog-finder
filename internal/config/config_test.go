package config

import "testing"

func TestCarregarPadrao(t *testing.T) {
	t.Setenv("FINDER_HTTP_ADDR", "")
	t.Setenv("FINDER_WORKERS", "")
	c := Carregar()
	if c.HTTPAddr != ":8085" || c.Workers != 4 {
		t.Fatalf("config inesperada: %+v", c)
	}
}

func TestCarregarWorkers(t *testing.T) {
	t.Setenv("FINDER_WORKERS", "8")
	if c := Carregar(); c.Workers != 8 {
		t.Fatalf("workers = %d", c.Workers)
	}
}
