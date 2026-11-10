package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

func TestCarregarNaoCriaParametrosDeConsumo(t *testing.T) {
	ausentarConsumo(t)
	c := Carregar()
	if c.Workers != 4 {
		t.Fatalf("workers = %d", c.Workers)
	}
	if c.RedisCount != 0 || c.RedisBlock != 0 || c.RedisPendenciaOciosa != 0 {
		t.Fatal("quantidade, bloqueio ou ociosidade receberam padrão")
	}
	if c.RedisGroup != "" || c.RedisConsumidor != "" || c.RedisIdempotencia != "" {
		t.Fatal("grupo, consumidor ou prefixo receberam padrão")
	}
	if c.RedisSenha != "" {
		t.Fatal("senha vazia não permaneceu vazia")
	}
	err := c.ValidarConsumo()
	if err == nil {
		t.Fatal("validação aceitou configuração vazia")
	}
	texto := err.Error()
	for _, parte := range []string{"stream", "grupo", "consumidor", "prefixo", "quantidade", "bloqueio", "ociosidade"} {
		if !strings.Contains(texto, parte) {
			t.Fatalf("validação não aponta %s: %s", parte, texto)
		}
	}
}

func TestCarregarIgnoraArquivoEnvComoPadrao(t *testing.T) {
	dir := t.TempDir()
	conteudo := "" +
		"FINDER_REDIS_STREAM=lhamalog:stream:logs\n" +
		"FINDER_REDIS_GROUP=grupo-exemplo\n" +
		"FINDER_REDIS_CONSUMIDOR=consumidor-exemplo\n" +
		"FINDER_REDIS_IDEMPOTENCIA=prefixo-exemplo:\n" +
		"FINDER_REDIS_COUNT=2\n" +
		"FINDER_REDIS_BLOCK=200ms\n" +
		"FINDER_REDIS_PENDENCIA_OCIOSA=1s\n" +
		"FINDER_REDIS_PASSWORD=segredo-exemplo\n"
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
	ausentarConsumo(t)
	c := Carregar()
	if c.RedisGroup != "" || c.RedisConsumidor != "" || c.RedisIdempotencia != "" || c.RedisCount != 0 || c.RedisBlock != 0 || c.RedisPendenciaOciosa != 0 || c.RedisSenha != "" {
		t.Fatal("nomes de exemplo viraram padrão de Carregar")
	}
}

func TestValidarRecusaCadaParametroNumerico(t *testing.T) {
	casos := []struct {
		nome  string
		chave string
		parte string
	}{
		{nome: "sem quantidade", chave: "FINDER_REDIS_COUNT", parte: "quantidade"},
		{nome: "sem bloqueio", chave: "FINDER_REDIS_BLOCK", parte: "bloqueio"},
		{nome: "sem ociosidade", chave: "FINDER_REDIS_PENDENCIA_OCIOSA", parte: "ociosidade"},
	}
	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			definirConsumoCompleto(t)
			t.Setenv(caso.chave, "")
			c := Carregar()
			if c.RedisCount < 0 || c.RedisBlock < 0 || c.RedisPendenciaOciosa < 0 {
				t.Fatal("a validação criou número")
			}
			switch caso.chave {
			case "FINDER_REDIS_COUNT":
				if c.RedisCount != 0 {
					t.Fatal("quantidade ausente ganhou valor")
				}
			case "FINDER_REDIS_BLOCK":
				if c.RedisBlock != 0 {
					t.Fatal("bloqueio ausente ganhou valor")
				}
			case "FINDER_REDIS_PENDENCIA_OCIOSA":
				if c.RedisPendenciaOciosa != 0 {
					t.Fatal("ociosidade ausente ganhou valor")
				}
			}
			err := c.ValidarConsumo()
			if err == nil || !strings.Contains(err.Error(), caso.parte) {
				t.Fatalf("err = %v", err)
			}
		})
	}
}

func TestValidarRecusaNaoPositivoEStreamForaDoContrato(t *testing.T) {
	definirConsumoCompleto(t)
	t.Setenv("FINDER_REDIS_STREAM", "fila:logs")
	t.Setenv("FINDER_REDIS_COUNT", "0")
	t.Setenv("FINDER_REDIS_BLOCK", "0s")
	t.Setenv("FINDER_REDIS_PENDENCIA_OCIOSA", "-1ms")
	t.Setenv("FINDER_REDIS_PASSWORD", "segredo-nao-logar")
	c := Carregar()
	if c.RedisCount != 0 || c.RedisBlock != 0 || c.RedisPendenciaOciosa != 0 {
		t.Fatal("valor não positivo foi aceito como configuração")
	}
	err := c.ValidarConsumo()
	if err == nil {
		t.Fatal("validação aceitou stream e números inválidos")
	}
	texto := err.Error()
	if strings.Contains(texto, "segredo-nao-logar") {
		t.Fatal("validação registrou a senha")
	}
	for _, parte := range []string{"stream", "quantidade", "bloqueio", "ociosidade"} {
		if !strings.Contains(texto, parte) {
			t.Fatalf("validação não aponta %s: %s", parte, texto)
		}
	}
}

func TestValidarAceitaConsumoInformado(t *testing.T) {
	definirConsumoCompleto(t)
	t.Setenv("FINDER_REDIS_PASSWORD", "segredo-nao-logar")
	c := Carregar()
	if err := c.ValidarConsumo(); err != nil {
		t.Fatal(err)
	}
	if c.RedisStream != "lhamalog:stream:logs" || c.RedisGroup != "grupo-exemplo" || c.RedisConsumidor != "consumidor-exemplo" || c.RedisIdempotencia != "prefixo-exemplo:" {
		t.Fatal("nomes informados não foram lidos")
	}
	if c.RedisCount != 2 || c.RedisBlock != 200*time.Millisecond || c.RedisPendenciaOciosa != time.Second {
		t.Fatal("números informados não foram lidos")
	}
	if c.RedisSenha != "segredo-nao-logar" {
		t.Fatal("senha informada não foi lida")
	}
	if c.Workers != 4 {
		t.Fatalf("workers = %d", c.Workers)
	}
}

func definirConsumoCompleto(t *testing.T) {
	t.Helper()
	t.Setenv("FINDER_REDIS_STREAM", "lhamalog:stream:logs")
	t.Setenv("FINDER_REDIS_GROUP", "grupo-exemplo")
	t.Setenv("FINDER_REDIS_CONSUMIDOR", "consumidor-exemplo")
	t.Setenv("FINDER_REDIS_IDEMPOTENCIA", "prefixo-exemplo:")
	t.Setenv("FINDER_REDIS_COUNT", "2")
	t.Setenv("FINDER_REDIS_BLOCK", "200ms")
	t.Setenv("FINDER_REDIS_PENDENCIA_OCIOSA", "1s")
	t.Setenv("FINDER_REDIS_PASSWORD", "")
	t.Setenv("FINDER_WORKERS", "")
}

func ausentarConsumo(t *testing.T) {
	t.Helper()
	for _, chave := range []string{
		"FINDER_REDIS_STREAM",
		"FINDER_REDIS_GROUP",
		"FINDER_REDIS_CONSUMIDOR",
		"FINDER_REDIS_IDEMPOTENCIA",
		"FINDER_REDIS_COUNT",
		"FINDER_REDIS_BLOCK",
		"FINDER_REDIS_PENDENCIA_OCIOSA",
		"FINDER_REDIS_PASSWORD",
		"FINDER_WORKERS",
	} {
		t.Setenv(chave, "")
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
