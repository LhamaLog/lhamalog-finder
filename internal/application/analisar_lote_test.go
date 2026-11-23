package application

import (
	"context"
	"errors"
	"testing"
	"time"

	entrada "github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

type moduloAnaliseTeste struct {
	consultas, envios int
	aceita            bool
	falha             error
	chaves            []domain.ChaveResumo
}

func (m *moduloAnaliseTeste) ConsultarAchado(context.Context, domain.IdentidadeAchado) ([]domain.EstadoOcorrencia, error) {
	m.consultas++
	return nil, nil
}
func (m *moduloAnaliseTeste) Enviar(ctx context.Context, r domain.ResumoAnalise) (bool, error) {
	m.envios++
	m.chaves = append(m.chaves, r.Chave())
	return m.aceita, m.falha
}
func pedidoAnaliseTeste() entrada.PedidoAnalise {
	p := entrada.PedidoAnalise{Referencia: domain.Referencia{BatchID: "batch", TenantID: "tenant", ProductID: "produto", RecordCount: 3}, ProdutoBody: "produto", SchemaVersion: "1", Completo: true, VersaoRegras: "v1", VersaoCatalogo: "c1"}
	for i := 0; i < 3; i++ {
		p.Logs = append(p.Logs, entrada.LogAnalise{Apontamento: domain.Codigo(string(rune('a' + i))), NivelCampo: "info", Mensagem: "bruto ERROR", ServiceName: "svc", Host: "host", Timestamp: time.Unix(int64(i+1), 0)})
	}
	return p
}
func scannerAnaliseTeste(contador *int) ValidacaoSeguranca {
	return ValidacaoSegurancaFunc(func(ctx context.Context, s []string) domain.ResultadoValidarSeguranca {
		*contador++
		r := domain.ResultadoValidarSeguranca{TotalSensivel: 2}
		for i := range s {
			r.LinhasMascaradas = append(r.LinhasMascaradas, domain.LinhaLogMascarada{Texto: "seguro WARNING", IndiceLinha: i})
		}
		return r
	})
}
func TestAnalisarMascaramentoAntesDaClasseResumoEConsulta(t *testing.T) {
	n := 0
	m := &moduloAnaliseTeste{aceita: true}
	s := ServicoAnalise{Seguranca: scannerAnaliseTeste(&n), Consulta: m, Envio: m}
	r, err := s.Analisar(context.Background(), pedidoAnaliseTeste())
	if err != nil || !r.Concluida || n != 1 || m.consultas != 1 || m.envios != 1 {
		t.Fatalf("resultado: err=%v segura=%d consultas=%d envios=%d", err, n, m.consultas, m.envios)
	}
	a := r.Resumo.Achados[0]
	if len(r.Resumo.Achados) != 1 || a.Tipo != domain.WARNING || a.Contagem != 3 || a.Recorrente || a.HostName != "host" || a.Linhas[0] != "seguro WARNING" {
		t.Fatal("achado nao reflete somente saida segura")
	}
	if r.Resumo.Contadores.Warnings != 3 || r.Resumo.Contadores.DadosSensiveisMascarados != 2 || r.Resumo.VersaoRegras != "v1" {
		t.Fatal("contadores/versao")
	}
}
func TestAnaliseFalhasDistinguiveisSemConsultaOuEnvio(t *testing.T) {
	for _, c := range []struct {
		motivo string
		mudar  func(*entrada.PedidoAnalise)
	}{
		{"log_ausente", func(p *entrada.PedidoAnalise) { p.Logs = nil }},
		{"armazenamento_inacessivel", func(p *entrada.PedidoAnalise) { p.ArmazenamentoInacessivel = true }},
		{"lote_incompleto", func(p *entrada.PedidoAnalise) { p.Completo = false }},
		{"identidade_incompleta", func(p *entrada.PedidoAnalise) { p.Referencia.TenantID = "" }},
		{"produto_ausente_ou_inconsistente", func(p *entrada.PedidoAnalise) { p.ProdutoBody = "outro" }},
		{"assinatura_incompleta", func(p *entrada.PedidoAnalise) { p.Logs[0].Host = "" }},
		{"envelope_invalido", func(p *entrada.PedidoAnalise) { p.SchemaVersion = "2" }},
	} {
		t.Run(c.motivo, func(t *testing.T) {
			p := pedidoAnaliseTeste()
			c.mudar(&p)
			n := 0
			m := &moduloAnaliseTeste{aceita: true}
			s := ServicoAnalise{Seguranca: scannerAnaliseTeste(&n), Consulta: m, Envio: m}
			r, err := s.Analisar(context.Background(), p)
			if err == nil || err.Error() != c.motivo || r.Concluida || m.envios != 0 {
				t.Fatalf("motivo=%v concluida=%v envios=%d", err, r.Concluida, m.envios)
			}
		})
	}
}
func TestAnaliseSemConfirmacaoReenviaMesmaChave(t *testing.T) {
	n := 0
	m := &moduloAnaliseTeste{falha: errors.New("falha privada")}
	s := ServicoAnalise{Seguranca: scannerAnaliseTeste(&n), Consulta: m, Envio: m}
	p := pedidoAnaliseTeste()
	r, err := s.Analisar(context.Background(), p)
	if err == nil || err.Error() != "envio_modulo_falhou" || r.Concluida {
		t.Fatal("concluiu sem confirmacao")
	}
	m.falha = nil
	m.aceita = true
	r, err = s.Analisar(context.Background(), p)
	if err != nil || !r.Concluida || m.chaves[0] != m.chaves[1] {
		t.Fatal("reenvio mudou chave")
	}
	p.ExecID = "V550e8400-e29b-41d4-a716-446655440000"
	r, err = s.Analisar(context.Background(), p)
	if err != nil || r.Resumo.Chave() == m.chaves[0] {
		t.Fatal("execId nao entrou na chave")
	}
}
func TestAnaliseFalhaSegurancaContinuaSemTrechos(t *testing.T) {
	m := &moduloAnaliseTeste{aceita: true}
	s := ServicoAnalise{Consulta: m, Envio: m, Seguranca: ValidacaoSegurancaFunc(func(context.Context, []string) domain.ResultadoValidarSeguranca {
		return domain.ResultadoValidarSeguranca{Erro: errors.New("falha")}
	})}
	r, err := s.Analisar(context.Background(), pedidoAnaliseTeste())
	if err != nil || !r.Concluida {
		t.Fatalf("lote nao continuou: %v", err)
	}
	for _, a := range r.Resumo.Achados {
		if len(a.Linhas) != 0 || a.Tipo != domain.INFO {
			t.Fatal("bruto reutilizado")
		}
	}
}

func TestProdutoInconsistenteEhPermanenteSemRetry(t *testing.T) {
	p := pedidoAnaliseTeste()
	p.ProdutoBody = "outro"
	s := ServicoAnalise{}
	_, err := s.Montar(context.Background(), p)
	f, ok := err.(interface{ Permanente() bool })
	if !ok || !f.Permanente() {
		t.Fatal("falha de produto nao sinaliza sem retry")
	}
}

func TestAnaliseParticoesECatalogoSomenteDecisorio(t *testing.T) {
	p := pedidoAnaliseTeste()
	p.Logs = nil
	for i, mensagem := range []string{"ERROR falha 1", "ERROR falha 2", "ERROR falha 3", "INFO alfa", "INFO beta"} {
		p.Logs = append(p.Logs, entrada.LogAnalise{Apontamento: domain.Codigo(string(rune('a' + i))), NivelCampo: "info", Mensagem: mensagem, ServiceName: "svc", Host: "host", Timestamp: time.Unix(int64(i), 0)})
	}
	p.Referencia.RecordCount = 5
	p.Catalogo = []domain.RegraClassificacao{{Classe: domain.WARNING, Corresponde: func(string) bool { return true }}}
	scanner := ValidacaoSegurancaFunc(func(_ context.Context, linhas []string) domain.ResultadoValidarSeguranca {
		r := domain.ResultadoValidarSeguranca{}
		for i, linha := range linhas {
			r.LinhasMascaradas = append(r.LinhasMascaradas, domain.LinhaLogMascarada{Texto: linha, IndiceLinha: i})
		}
		return r
	})
	m := &moduloAnaliseTeste{aceita: true}
	s := ServicoAnalise{Seguranca: scanner, Consulta: m, Envio: m}
	r, err := s.Analisar(context.Background(), p)
	if err != nil || !r.Concluida || len(r.Resumo.Achados) != 3 || m.consultas != 3 || m.envios != 1 || r.Resumo.IdentificacaoRegras.CatalogoEntrou {
		t.Fatalf("particoes/catalogo: %+v err=%v consultas=%d", r, err, m.consultas)
	}
	if r.Resumo.Contadores.Erros != 3 || r.Resumo.Contadores.Info != 2 {
		t.Fatal("contadores por log")
	}
	for _, a := range r.Resumo.Achados {
		if a.Tipo == domain.ERROR && a.Contagem != 3 {
			t.Fatal("ERROR nao agrupado")
		}
	}
	p.Logs[0].NivelCampo = ""
	p.Logs[0].Mensagem = "desconhecido"
	r, err = s.Analisar(context.Background(), p)
	if err != nil || !r.Resumo.IdentificacaoRegras.CatalogoEntrou || r.Resumo.IdentificacaoRegras.SomenteNivel {
		t.Fatal("catalogo decisorio nao propagado")
	}
}

type moduloIdempotenteTeste struct {
	efeitos  map[domain.ChaveResumo]bool
	chamadas int
}

func (m *moduloIdempotenteTeste) ConsultarAchado(context.Context, domain.IdentidadeAchado) ([]domain.EstadoOcorrencia, error) {
	return nil, nil
}
func (m *moduloIdempotenteTeste) Enviar(_ context.Context, r domain.ResumoAnalise) (bool, error) {
	m.chamadas++
	m.efeitos[r.Chave()] = true
	return true, nil
}
func TestAnaliseFakeIdempotenteEVersionamentoSemRetroatividade(t *testing.T) {
	m := &moduloIdempotenteTeste{efeitos: map[domain.ChaveResumo]bool{}}
	n := 0
	s := ServicoAnalise{Seguranca: scannerAnaliseTeste(&n), Consulta: m, Envio: m}
	p := pedidoAnaliseTeste()
	primeiro, err := s.Analisar(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.Analisar(context.Background(), p); err != nil || len(m.efeitos) != 1 || m.chamadas != 2 {
		t.Fatal("duplicou efeito no fake")
	}
	p.VersaoRegras = "v2"
	p.Referencia.BatchID = "novo"
	segundo, err := s.Analisar(context.Background(), p)
	if err != nil || primeiro.Resumo.VersaoRegras != "v1" || segundo.Resumo.VersaoRegras != "v2" {
		t.Fatal("versao alterou resumo anterior")
	}
}
