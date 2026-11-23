package modocorrencia

import (
	"encoding/json"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

type contadoresDTO struct {
	LinhasProcessadas        int64 `json:"linhasProcessadas"`
	Erros                    int64 `json:"erros"`
	Warnings                 int64 `json:"warnings"`
	Informacao               int64 `json:"informacao"`
	DadosSensiveisMascarados int64 `json:"dadosSensiveisMascarados"`
}

type achadoDTO struct {
	Tipo               domain.ClasseLog  `json:"tipo"`
	Assinatura         string            `json:"assinatura"`
	HostName           string            `json:"hostName"`
	Contagem           int64             `json:"contagem"`
	PrimeiraOcorrencia time.Time         `json:"primeiraOcorrencia"`
	UltimaOcorrencia   time.Time         `json:"ultimaOcorrencia"`
	Recorrente         bool              `json:"recorrente"`
	Linhas             []string          `json:"linhas"`
	RevisaoConsultada  *string           `json:"revisaoConsultada"`
	Titulo             *string           `json:"titulo"`
	Resumo             *string           `json:"resumo"`
	NumeroLinha        *int              `json:"numeroLinha"`
	CamposExtraido     map[string]string `json:"camposExtraido"`
	DadosSensiveis     *bool             `json:"dadosSensiveis"`
	TraceID            string            `json:"traceId,omitempty"`
	CorrelationID      string            `json:"correlationId,omitempty"`
	ModuloID           string            `json:"moduloId,omitempty"`
}

type resumoDTO struct {
	BatchID             string                      `json:"batchId"`
	ExecID              string                      `json:"execId,omitempty"`
	TenantID            string                      `json:"tenantId"`
	ProductID           string                      `json:"productId"`
	VersaoRegras        string                      `json:"versaoRegras"`
	VersaoCatalogo      string                      `json:"versaoCatalogo,omitempty"`
	IdentificacaoRegras *domain.IdentificacaoRegras `json:"identificacaoRegras,omitempty"`
	Contadores          contadoresDTO               `json:"contadores"`
	Achados             []achadoDTO                 `json:"achados"`
}

func SerializarResumo(r domain.ResumoAnalise) ([]byte, error) {
	if err := domain.ValidarResumoHandoff(r); err != nil {
		return nil, out.ErrModuloRequisicao
	}
	c := r.Contadores
	dto := resumoDTO{
		BatchID: r.BatchID, ExecID: r.ExecID, TenantID: r.TenantID, ProductID: r.ProductID,
		VersaoRegras: r.VersaoRegras, VersaoCatalogo: r.VersaoCatalogo,
		Contadores: contadoresDTO{c.LinhasProcessadas, c.Erros, c.Warnings, c.Info, c.DadosSensiveisMascarados},
		Achados:    make([]achadoDTO, 0, len(r.Achados)),
	}
	if r.IdentificacaoRegras != (domain.IdentificacaoRegras{}) {
		if r.IdentificacaoRegras.VersaoUnica != r.VersaoRegras {
			return nil, out.ErrModuloRequisicao
		}
		dto.IdentificacaoRegras = &r.IdentificacaoRegras
	}
	chaves := make(map[domain.IdentidadeAchado]struct{}, len(r.Achados))
	for _, a := range r.Achados {
		if algumVazio(a.Assinatura, a.HostName) || a.Contagem < 1 || a.PrimeiraOcorrencia.IsZero() || a.UltimaOcorrencia.Before(a.PrimeiraOcorrencia) || a.NumeroLinha < 0 ||
			(a.RevisaoConsultada == nil && a.StatusAnterior != "") || (a.RevisaoConsultada != nil && algumVazio(*a.RevisaoConsultada)) {
			return nil, out.ErrModuloRequisicao
		}
		chave := domain.IdentidadeAchado{Assinatura: a.Assinatura, Classe: a.Tipo, TenantID: r.TenantID, ProductID: r.ProductID}
		if _, existe := chaves[chave]; existe {
			return nil, out.ErrModuloRequisicao
		}
		chaves[chave] = struct{}{}
		achado := achadoDTO{
			Tipo: a.Tipo, Assinatura: a.Assinatura, HostName: a.HostName, Contagem: a.Contagem,
			PrimeiraOcorrencia: a.PrimeiraOcorrencia.UTC(), UltimaOcorrencia: a.UltimaOcorrencia.UTC(),
			Recorrente: a.Recorrente, Linhas: append([]string{}, a.Linhas...), RevisaoConsultada: a.RevisaoConsultada,
			CamposExtraido: a.CamposExtraido, DadosSensiveis: a.DadosSensiveis, TraceID: a.TraceID,
			CorrelationID: a.CorrelationID, ModuloID: a.ModuloID,
		}
		if a.Titulo != "" {
			achado.Titulo = &a.Titulo
		}
		if a.Resumo != "" {
			achado.Resumo = &a.Resumo
		}
		if a.NumeroLinha > 0 {
			achado.NumeroLinha = &a.NumeroLinha
		}
		dto.Achados = append(dto.Achados, achado)
	}
	dados, err := json.Marshal(dto)
	if err != nil {
		return nil, out.ErrModuloRequisicao
	}
	return dados, nil
}
