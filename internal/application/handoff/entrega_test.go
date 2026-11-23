package handoff

import (
	"context"
	"errors"
	entrada "github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"testing"
)

type moduloTeste struct {
	envios  int
	aceitar bool
	chaves  []domain.ChaveResumo
}

func (m *moduloTeste) Enviar(ctx context.Context, r domain.ResumoAnalise) (bool, error) {
	m.envios++
	m.chaves = append(m.chaves, r.Chave())
	if !m.aceitar {
		return false, errors.New("falha privada")
	}
	return true, nil
}

type fimTeste struct {
	n int
	m *moduloTeste
}

func (f *fimTeste) ConcluirAnalise(context.Context, domain.ChaveResumo) error {
	if !f.m.aceitar {
		return errors.New("concluiu antes da confirmacao")
	}
	f.n++
	return nil
}
func resumoTeste() domain.ResumoAnalise {
	return domain.ResumoAnalise{BatchID: "batch", TenantID: "t", ProductID: "p", VersaoRegras: "v1", Contadores: domain.ContadoresResumo{LinhasProcessadas: 1, Info: 1}, Achados: []domain.AchadoAnalise{{Tipo: domain.INFO, Linhas: []string{"seguro"}}}}
}
func TestHandoffEsperaConfirmacaoEReenviaMesmaChave(t *testing.T) {
	m := &moduloTeste{}
	f := &fimTeste{m: m}
	s := Entrega{Envio: m, Finalizar: f}
	original := resumoTeste()
	p, err := NovaPendente(original)
	if err != nil {
		t.Fatal(err)
	}
	original.BatchID = "mudou"
	original.Achados[0].Linhas[0] = "mudou"
	r, err := s.Enviar(context.Background(), p)
	if err == nil || r.Concluida || f.n != 0 {
		t.Fatal("falha concluiu varredura")
	}
	m.aceitar = true
	r, err = s.Reenviar(context.Background(), p)
	if err != nil || !r.Concluida || f.n != 1 || m.chaves[0] != m.chaves[1] || m.chaves[0].BatchID != "batch" {
		t.Fatal("reenvio mudou chave ou ordem")
	}
}
func TestHandoffRejeitaTaxonomiaEExecIdInvalido(t *testing.T) {
	r := resumoTeste()
	r.Achados[0].Tipo = "SPAM"
	if _, err := NovaPendente(r); err == nil {
		t.Fatal("taxonomia invalida")
	}
	r = resumoTeste()
	r.ExecID = "inventado"
	if _, err := NovaPendente(r); err == nil {
		t.Fatal("execId sem origem/forma")
	}
	r.ExecID = "V550e8400-e29b-41d4-a716-446655440000"
	p, err := NovaPendente(r)
	if err != nil || p.Chave().ExecID == "" {
		t.Fatal("execId valido nao entrou")
	}
}

type preparadorTeste struct{ n int }

func (p *preparadorTeste) Montar(context.Context, entrada.PedidoAnalise) (domain.ResumoAnalise, error) {
	p.n++
	return resumoTeste(), nil
}
func TestOrquestrarEntregaBloqueiaLoteIncompleto(t *testing.T) {
	p := &preparadorTeste{}
	o := OrquestradorEntrega{Preparador: p}
	r, err := o.Executar(context.Background(), entrada.PedidoAnalise{Completo: false})
	if err == nil || r.Concluida || p.n != 0 {
		t.Fatal("incompleto analisado/enviado")
	}
}

type moduloMutadorTeste struct {
	recebidos []bool
	linhas    []string
	campos    []string
}

func (m *moduloMutadorTeste) Enviar(_ context.Context, r domain.ResumoAnalise) (bool, error) {
	a := r.Achados[0]
	m.recebidos = append(m.recebidos, *a.DadosSensiveis)
	m.linhas = append(m.linhas, a.Linhas[0])
	m.campos = append(m.campos, a.CamposExtraido["chave"])
	*a.DadosSensiveis = false
	a.Linhas[0] = "mutou"
	a.CamposExtraido["chave"] = "mutou"
	return false, errors.New("sem confirmacao")
}
func TestSnapshotHandoffIsolaPonteirosDaOrigemEDoAdapter(t *testing.T) {
	r := resumoTeste()
	flag := true
	r.Achados[0].DadosSensiveis = &flag
	r.Achados[0].CamposExtraido = map[string]string{"chave": "seguro"}
	p, err := NovaPendente(r)
	if err != nil {
		t.Fatal(err)
	}
	flag = false
	r.Achados[0].Linhas[0] = "origem"
	r.Achados[0].CamposExtraido["chave"] = "origem"
	m := &moduloMutadorTeste{}
	s := Entrega{Envio: m}
	s.Enviar(context.Background(), p)
	s.Reenviar(context.Background(), p)
	for i := 0; i < 2; i++ {
		if !m.recebidos[i] || m.linhas[i] != "seguro" || m.campos[i] != "seguro" {
			t.Fatalf("snapshot alterado: flags=%v linhas=%v campos=%v", m.recebidos, m.linhas, m.campos)
		}
	}
}

func (p *preparadorTeste) Agrupar(_ context.Context, r domain.ResumoAnalise) (domain.ResumoAnalise, error) {
	return r, nil
}
