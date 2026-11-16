// Package config carrega a configuração por variáveis de ambiente.
package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

const streamContrato = "lhamalog:stream:logs"

// Config reúne as configurações do Finder.
// Quantidade, bloqueio, ociosidade, grupo, consumidor e prefixo só existem quando o ambiente os informa.
// Workers continua com o padrão 4 e não é o grau de paralelismo deste consumo.
// RedisSenha vazia permanece vazia. O valor não entra em registro.
type Config struct {
	HTTPAddr             string
	Workers              int
	RedisAddr            string
	RedisStream          string
	RedisGroup           string
	RedisCodigos         string
	RedisConsumidor      string
	RedisIdempotencia    string
	RedisCount           int
	RedisBlock           time.Duration
	RedisPendenciaOciosa time.Duration
	RedisSenha           string
	// Parâmetros do OpenObserve. Nenhum recebe URL, número, lista ou nome de campo nesta leitura.
	// OOSenha vazia permanece vazia e não entra em registro.
	// OOCamposCopiaPresente distingue lista ausente de lista presente vazia.
	OOURL                 string
	OOUsuario             string
	OOSenha               string
	OOTimeout             time.Duration
	OOPagina              int
	OOCamposCopiaPresente bool
	OOCamposCopia         []string
	OOOrdem               string
	OOCampoTenant         string
	OOCampoProduto        string
}

// ErroConsumo lista o que impede abrir o consumo. Não inclui senha nem o valor informado.
type ErroConsumo struct {
	Pendencias []string
}

func (e ErroConsumo) Error() string {
	return strings.Join(e.Pendencias, "; ")
}

// Carregar lê as variáveis de ambiente com valores padrão.
// RedisAddr usa localhost:6379, o endereço já usado pelo processo e pelo Compose.
// Stream, grupo, consumidor, prefixo, quantidade, bloqueio e ociosidade não têm padrão no código.
// A senha vem de FINDER_REDIS_PASSWORD. Vazia continua vazia.
// URL, usuário, senha, timeout, página, lista de cópia, ordenação e nomes de tenant e produto
// do OpenObserve também não têm padrão. A senha do OpenObserve não é registrada.
// O processo nativo recebe esses nomes de um .env, aplicado antes desta leitura.
func Carregar() Config {
	copiaPresente, copia := listaCopia()
	c := Config{
		HTTPAddr:              valor("FINDER_HTTP_ADDR", ":8085"),
		Workers:               4,
		RedisAddr:             valor("FINDER_REDIS_ADDR", "localhost:6379"),
		RedisStream:           os.Getenv("FINDER_REDIS_STREAM"),
		RedisGroup:            os.Getenv("FINDER_REDIS_GROUP"),
		RedisCodigos:          os.Getenv("FINDER_REDIS_CODIGOS"),
		RedisConsumidor:       os.Getenv("FINDER_REDIS_CONSUMIDOR"),
		RedisIdempotencia:     os.Getenv("FINDER_REDIS_IDEMPOTENCIA"),
		RedisCount:            inteiroPositivo(os.Getenv("FINDER_REDIS_COUNT")),
		RedisBlock:            duracaoPositiva(os.Getenv("FINDER_REDIS_BLOCK")),
		RedisPendenciaOciosa:  duracaoPositiva(os.Getenv("FINDER_REDIS_PENDENCIA_OCIOSA")),
		RedisSenha:            os.Getenv("FINDER_REDIS_PASSWORD"),
		OOURL:                 os.Getenv("FINDER_OO_URL"),
		OOUsuario:             os.Getenv("FINDER_OO_USUARIO"),
		OOSenha:               os.Getenv("FINDER_OO_SENHA"),
		OOTimeout:             duracaoPositiva(os.Getenv("FINDER_OO_TIMEOUT")),
		OOPagina:              inteiroPositivo(os.Getenv("FINDER_OO_PAGINA")),
		OOCamposCopiaPresente: copiaPresente,
		OOCamposCopia:         copia,
		OOOrdem:               os.Getenv("FINDER_OO_ORDEM"),
		OOCampoTenant:         os.Getenv("FINDER_OO_CAMPO_TENANT"),
		OOCampoProduto:        os.Getenv("FINDER_OO_CAMPO_PRODUTO"),
	}
	if n, err := strconv.Atoi(os.Getenv("FINDER_WORKERS")); err == nil && n > 0 {
		c.Workers = n
	}
	return c
}

// ValidarConsumo recusa stream vazio ou diferente do contrato, e recusa grupo,
// consumidor, prefixo, quantidade, bloqueio ou ociosidade ausentes ou não positivos.
func (c Config) ValidarConsumo() error {
	var pendencias []string
	if c.RedisStream == "" || c.RedisStream != streamContrato {
		pendencias = append(pendencias, "stream ausente ou diferente de lhamalog:stream:logs")
	}
	if c.RedisGroup == "" {
		pendencias = append(pendencias, "grupo ausente")
	}
	if c.RedisConsumidor == "" {
		pendencias = append(pendencias, "consumidor ausente")
	}
	if c.RedisIdempotencia == "" {
		pendencias = append(pendencias, "prefixo ausente")
	}
	if c.RedisCount <= 0 {
		pendencias = append(pendencias, "quantidade ausente ou não positiva")
	}
	if c.RedisBlock <= 0 {
		pendencias = append(pendencias, "bloqueio ausente ou não positivo")
	}
	if c.RedisPendenciaOciosa <= 0 {
		pendencias = append(pendencias, "ociosidade ausente ou não positiva")
	}
	if len(pendencias) == 0 {
		return nil
	}
	return ErroConsumo{Pendencias: pendencias}
}

func inteiroPositivo(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return 0
	}
	return n
}

func duracaoPositiva(s string) time.Duration {
	if s == "" {
		return 0
	}
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return 0
	}
	return d
}

// Acesso entrega a configuração lida, sem preencher o que o ambiente omitiu.
// A senha segue no valor de retorno e não é registrada aqui.
func (c Config) Acesso() domain.ConfiguracaoAcesso {
	acesso := domain.ConfiguracaoAcesso{
		URL:           c.OOURL,
		Usuario:       c.OOUsuario,
		Senha:         c.OOSenha,
		Timeout:       c.OOTimeout,
		TamanhoPagina: c.OOPagina,
		Ordem:         c.OOOrdem,
		CampoTenant:   c.OOCampoTenant,
		CampoProduto:  c.OOCampoProduto,
	}
	if c.OOCamposCopiaPresente {
		copia := append([]string(nil), c.OOCamposCopia...)
		acesso.CamposCopia = &copia
	}
	return acesso
}

// listaCopia distingue variável ausente de variável presente e vazia.
func listaCopia() (bool, []string) {
	valor, ok := os.LookupEnv("FINDER_OO_CAMPOS_COPIA")
	if !ok {
		return false, nil
	}
	valor = strings.TrimSpace(valor)
	if valor == "" {
		return true, []string{}
	}
	partes := strings.Split(valor, ",")
	saida := make([]string, 0, len(partes))
	for _, parte := range partes {
		parte = strings.TrimSpace(parte)
		if parte != "" {
			saida = append(saida, parte)
		}
	}
	return true, saida
}

func valor(chave, padrao string) string {
	if v := os.Getenv(chave); v != "" {
		return v
	}
	return padrao
}
