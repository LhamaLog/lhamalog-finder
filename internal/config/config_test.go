package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCarregarPadrao(t *testing.T) {
	t.Setenv("FINDER_HTTP_ADDR", "")
	t.Setenv("FINDER_WORKERS", "")
	t.Setenv("FINDER_REDIS_ADDR", "")
	t.Setenv("FINDER_REDIS_STREAM", "")
	t.Setenv("FINDER_REDIS_GROUP", "")
	t.Setenv("FINDER_REDIS_CODIGOS", "")
	c := Carregar()
	if c.HTTPAddr != ":8085" || c.Workers != 4 || c.RedisAddr != "localhost:6379" {
		t.Fatalf("config inesperada: %+v", c)
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

func TestAplicarArquivoEnv(t *testing.T) {
	dir := t.TempDir()
	conteudo := "" +
		"# comentario\n" +
		"\n" +
		"export FINDER_REDIS_STREAM=\"fila:logs\"\n" +
		"FINDER_REDIS_GROUP=finder\n" +
		"FINDER_REDIS_CODIGOS='finder:codigos'\n" +
		"FINDER_WORKERS=9\n"
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)

	restaurar := ausentar(t, "FINDER_REDIS_STREAM", "FINDER_REDIS_GROUP", "FINDER_REDIS_CODIGOS")
	defer restaurar()
	t.Setenv("FINDER_WORKERS", "8")

	if err := AplicarArquivoEnv(); err != nil {
		t.Fatal(err)
	}
	if os.Getenv("FINDER_REDIS_STREAM") != "fila:logs" || os.Getenv("FINDER_REDIS_GROUP") != "finder" || os.Getenv("FINDER_REDIS_CODIGOS") != "finder:codigos" {
		t.Fatalf("stream=%q grupo=%q codigos=%q", os.Getenv("FINDER_REDIS_STREAM"), os.Getenv("FINDER_REDIS_GROUP"), os.Getenv("FINDER_REDIS_CODIGOS"))
	}
	if os.Getenv("FINDER_WORKERS") != "8" {
		t.Fatalf("workers substituído: %s", os.Getenv("FINDER_WORKERS"))
	}
}

func TestAplicarArquivoEnvAusente(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := AplicarArquivoEnv(); err != nil {
		t.Fatal(err)
	}
}

func ausentar(t *testing.T, chaves ...string) func() {
	t.Helper()
	type anterior struct {
		chave   string
		valor   string
		existia bool
	}
	salvos := make([]anterior, 0, len(chaves))
	for _, chave := range chaves {
		valor, existia := os.LookupEnv(chave)
		salvos = append(salvos, anterior{chave, valor, existia})
		if err := os.Unsetenv(chave); err != nil {
			t.Fatal(err)
		}
	}
	return func() {
		for _, item := range salvos {
			if item.existia {
				os.Setenv(item.chave, item.valor)
			} else {
				os.Unsetenv(item.chave)
			}
		}
	}
}
