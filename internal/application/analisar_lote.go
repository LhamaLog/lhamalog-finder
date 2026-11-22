package application

import (
	"context"
	"errors"
	"github.com/LhamaLog/lhamalog-finder/internal/application/agrupar"

	entrada "github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

type ServicoAnalise struct {
	Seguranca  ValidacaoSeguranca
	Consulta   out.ConsultarAchadoModulo
	Envio      out.EnviarResumo
	Assinatura out.Assinador
}

var _ entrada.Analise = ServicoAnalise{}

// Montar executa somente o ciclo local: segurança → classificação → assinatura → resumo.
// A medição SC-002 exclui consulta/envio; a medição formal exige VM-A real.
func (s ServicoAnalise) Montar(ctx context.Context, p entrada.PedidoAnalise) (domain.ResumoAnalise, error) {
	falhar := func(m string) (domain.ResumoAnalise, error) {
		return domain.ResumoAnalise{}, domain.FalhaAnalise{Motivo: m, SemRetry: m == "produto_ausente_ou_inconsistente"}
	}
	if err := ctx.Err(); err != nil {
		return domain.ResumoAnalise{}, err
	}
	if p.ArmazenamentoInacessivel {
		return falhar("armazenamento_inacessivel")
	}
	if p.SchemaVersion != "1" {
		return falhar("envelope_invalido")
	}
	if p.Referencia.ProductID == "" || p.ProdutoBody == "" || p.ProdutoBody != p.Referencia.ProductID {
		return falhar("produto_ausente_ou_inconsistente")
	}
	if p.Referencia.TenantID == "" || p.Referencia.BatchID == "" {
		return falhar("identidade_incompleta")
	}
	if len(p.Logs) == 0 {
		return falhar("log_ausente")
	}
	if !p.Completo || int64(len(p.Logs)) != p.Referencia.RecordCount {
		return falhar("lote_incompleto")
	}
	if p.VersaoRegras == "" {
		return falhar("versao_regras_ausente")
	}
	brutas := make([]string, len(p.Logs))
	for i, l := range p.Logs {
		brutas[i] = l.Mensagem
	}
	seguro := MontarContexto(ctx, p.Referencia.BatchID, brutas, p.VersaoCatalogo, s.Seguranca)
	r := domain.ResumoAnalise{BatchID: p.Referencia.BatchID, ExecID: p.ExecID, TenantID: p.Referencia.TenantID, ProductID: p.Referencia.ProductID, VersaoRegras: p.VersaoRegras, VersaoCatalogo: seguro.VersaoCatalogo, IdentificacaoRegras: domain.IdentificacaoRegras{VersaoUnica: p.VersaoRegras, SomenteNivel: true}, Achados: []domain.AchadoAnalise{}}
	r.Contadores.DadosSensiveisMascarados = int64(seguro.DadosSensiveisIdentificados)
	classificados := make([]domain.LogClassificado, 0, len(p.Logs))
	for i, l := range p.Logs {
		if err := ctx.Err(); err != nil {
			return domain.ResumoAnalise{}, err
		}
		linha := domain.LinhaLogMascarada{}
		if !seguro.FalhaSeguranca {
			linha = seguro.Linhas[i]
		}
		classificado, err := ClassificarLogCapturado(p.Referencia, l, linha, p.VersaoRegras, p.Catalogo)
		if err != nil {
			return domain.ResumoAnalise{}, err
		}
		classificados = append(classificados, classificado)
		if classificado.CatalogoUsado {
			r.IdentificacaoRegras.CatalogoEntrou = true
			r.IdentificacaoRegras.SomenteNivel = false
		}
		r.Contadores.LinhasProcessadas++
		switch classificado.Classe {
		case domain.ERROR:
			r.Contadores.Erros++
		case domain.WARNING:
			r.Contadores.Warnings++
		case domain.INFO:
			r.Contadores.Info++
		}
	}
	assinar := domain.AssinarClassificado
	if s.Assinatura != nil {
		assinar = s.Assinatura.Calcular
	}
	particoes, err := agrupar.ParticionarComAssinador(classificados, assinar)
	if err != nil {
		return domain.ResumoAnalise{}, err
	}
	for _, particao := range particoes {
		primeiro := particao.Logs[0]
		achado := domain.AchadoAnalise{Tipo: particao.Chave.Classe, Assinatura: particao.Chave.Assinatura, Contagem: particao.ContagemLote, HostName: primeiro.Host, CorrelationID: p.Referencia.CorrelationID, TraceID: primeiro.TraceID, ModuloID: primeiro.ModuloID, Linhas: []string{}, CamposExtraido: map[string]string{}, VersaoRegras: p.VersaoRegras, PrimeiraOcorrencia: particao.PrimeiraNoLote, UltimaOcorrencia: particao.UltimaNoLote}
		if !seguro.FalhaSeguranca {
			for _, l := range particao.Logs {
				achado.Linhas = append(achado.Linhas, l.Mensagem)
			}
		}
		r.Achados = append(r.Achados, achado)
	}

	return r, nil
}

// Analisar usa portas injetadas. Fake é só teste; nenhuma confirmação é presumida.
func (s ServicoAnalise) Analisar(ctx context.Context, p entrada.PedidoAnalise) (entrada.ResultadoAnalise, error) {
	r, err := s.Montar(ctx, p)
	if err != nil {
		return entrada.ResultadoAnalise{}, err
	}
	resultado := entrada.ResultadoAnalise{Resumo: r}
	if s.Consulta == nil {
		return resultado, errors.New("consulta_modulo_indisponivel")
	}
	for _, a := range r.Achados {
		_, err := s.Consulta.ConsultarAchado(ctx, domain.IdentidadeAchado{Assinatura: a.Assinatura, Classe: a.Tipo, TenantID: r.TenantID, ProductID: r.ProductID})
		if err != nil {
			return resultado, errors.New("consulta_modulo_falhou")
		}
	}
	if s.Envio == nil {
		return resultado, errors.New("envio_modulo_falhou")
	}
	aceito, err := s.Envio.Enviar(ctx, r)
	if err != nil || !aceito {
		return resultado, errors.New("envio_modulo_falhou")
	}
	resultado.Concluida = true
	return resultado, nil
}
