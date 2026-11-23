package bootstrap

import (
	"context"
	"testing"

	"github.com/LhamaLog/lhamalog-finder/internal/config"
)

type chavesTeste struct{}

func (chavesTeste) Reservar(context.Context, string) (bool, error)    { return true, nil }
func (chavesTeste) Concluir(context.Context, string) error            { return nil }
func (chavesTeste) Consultar(context.Context, string) (string, error) { return "", nil }
func (chavesTeste) RemoverReserva(context.Context, string) error      { return nil }

func TestBootstrapValidaTudoAntesDeConsumir(t *testing.T) {
	env := map[string]string{"FINDER_REDIS_STREAM": "lhamalog:stream:logs", "FINDER_REDIS_GROUP": "teste", "FINDER_REDIS_CONSUMIDOR": "teste", "FINDER_REDIS_IDEMPOTENCIA": "teste:", "FINDER_REDIS_COUNT": "1", "FINDER_REDIS_BLOCK": "100ms", "FINDER_REDIS_PENDENCIA_OCIOSA": "100ms", "FINDER_OO_URL": "http://localhost:15080", "FINDER_OO_USUARIO": "teste", "FINDER_OO_SENHA": "teste", "FINDER_OO_TIMEOUT": "1s", "FINDER_OO_PAGINA": "2", "FINDER_OO_CAMPOS_COPIA": "", "FINDER_OO_ORDEM": "lhamalog_record_id", "FINDER_OO_CAMPO_TENANT": "tenantid", "FINDER_OO_CAMPO_PRODUTO": "productid", "FINDER_CATALOGO_SEGURANCA": "../../config/seguranca/catalogo.json", "FINDER_VERSAO_REGRAS": "finder-2026-11-22.1", "FINDER_MODULO_OCORRENCIA_URL": "https://modulo.invalid", "FINDER_MODULO_OCORRENCIA_JWT": "teste", "FINDER_CORE_URL": "https://core.invalid", "FINDER_CORE_JWT": "teste", "FINDER_CORE_TIMEOUT": "2s", "FINDER_VARREDURA_URL": "", "FINDER_VARREDURA_JWT": "", "FINDER_VARREDURA_TIMEOUT": "2s"}
	for k, v := range env {
		t.Setenv(k, v)
	}
	if _, err := Novo(context.Background(), config.Carregar(), chavesTeste{}, Opcoes{}); err != nil {
		t.Fatal(err)
	}
	for k, bad := range map[string]string{"FINDER_VERSAO_REGRAS": "outra", "FINDER_CATALOGO_SEGURANCA": "ausente", "FINDER_MODULO_OCORRENCIA_URL": "http://modulo.invalid", "FINDER_MODULO_OCORRENCIA_JWT": "", "FINDER_CORE_URL": "http://core.invalid", "FINDER_CORE_JWT": "", "FINDER_CORE_TIMEOUT": "-1s", "FINDER_VARREDURA_URL": "https://scan.invalid", "FINDER_OO_CAMPO_TENANT": "", "FINDER_OO_CAMPOS_COPIA": ""} {
		if k == "FINDER_OO_CAMPOS_COPIA" {
			continue
		} // vazio presente é configuração válida.
		t.Run(k, func(t *testing.T) {
			t.Setenv(k, bad)
			if _, err := Novo(context.Background(), config.Carregar(), chavesTeste{}, Opcoes{}); err == nil {
				t.Fatal("inicializou configuração inválida")
			}
		})
	}
	if _, err := Novo(context.Background(), config.Carregar(), nil, Opcoes{}); err == nil {
		t.Fatal("inicializou sem armazenamento idempotente")
	}
}
