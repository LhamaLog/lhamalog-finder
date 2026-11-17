package application

import (
	"context"
	"errors"
	"testing"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

type catalogoTeste struct {
	versao string
	err    error
}

func (c catalogoTeste) Carregar(context.Context) (domain.CatalogoSeguranca, error) {
	return domain.CatalogoSeguranca{Versao: c.versao}, c.err
}
func TestBootstrapSegurancaExigeVersaoCombinada(t *testing.T) {
	for _, v := range []string{"", "outra"} {
		if _, err := InicializarSeguranca(context.Background(), catalogoTeste{versao: v}, "combinada"); err == nil {
			t.Fatal("versao divergente aceita")
		}
	}
	if _, err := InicializarSeguranca(context.Background(), catalogoTeste{versao: "combinada"}, "combinada"); err != nil {
		t.Fatal(err)
	}
	if _, err := InicializarSeguranca(context.Background(), catalogoTeste{err: errors.New("privado")}, "combinada"); err == nil || err.Error() == "privado" {
		t.Fatal("falha ausente ou insegura")
	}
}
func TestMontarContextoUmaValidacaoSemReusarBruto(t *testing.T) {
	chamadas := 0
	v := ValidacaoSegurancaFunc(func(ctx context.Context, linhas []string) domain.ResultadoValidarSeguranca {
		chamadas++
		return domain.ResultadoValidarSeguranca{LinhasMascaradas: []domain.LinhaLogMascarada{{Texto: "seguro", IndiceLinha: 0}}, TotalSensivel: 2}
	})
	c := MontarContexto(context.Background(), "batch", []string{"bruto"}, "versao", v)
	if chamadas != 1 || c.Linhas[0].Texto != "seguro" || c.DadosSensiveisIdentificados != 2 || c.VersaoCatalogo != "versao" {
		t.Fatalf("contexto invalido: chamadas=%d", chamadas)
	}
	falho := ValidacaoSegurancaFunc(func(context.Context, []string) domain.ResultadoValidarSeguranca {
		return domain.ResultadoValidarSeguranca{LinhasMascaradas: []domain.LinhaLogMascarada{{Texto: "bruto"}}, Erro: errors.New("falha")}
	})
	c = MontarContexto(context.Background(), "batch", []string{"bruto"}, "versao", falho)
	if !c.FalhaSeguranca || len(c.Linhas) != 0 {
		t.Fatal("falha reutilizou linhas")
	}
}
func TestValidarSegurancaDelegaCatalogoSemInventarRegex(t *testing.T) {
	v := Seguranca{Catalogo: domain.CatalogoSeguranca{Versao: "v", Transformar: func(s string) (domain.LinhaLogMascarada, int, error) {
		return domain.LinhaLogMascarada{Texto: domain.MarcadorOmisso, LinhaRemovida: true}, 1, nil
	}}}
	r := v.ValidarSeguranca(context.Background(), []string{"entrada-opaca"})
	if r.Erro != nil || r.TotalSensivel != 1 || r.LinhasMascaradas[0].Texto != "log com senha omitido" {
		t.Fatal("contrato de transformacao")
	}
	r = (Seguranca{}).ValidarSeguranca(context.Background(), []string{"entrada-opaca"})
	if r.Erro == nil || len(r.LinhasMascaradas) != 0 {
		t.Fatal("catalogo ausente fez bypass")
	}
}
