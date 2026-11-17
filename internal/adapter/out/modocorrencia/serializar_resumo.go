package modocorrencia

import (
	"encoding/json"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// DTO provisório da 009: informacao (ADR_102 pendente), diferente do modelo
// lógico 006 info. Nenhum cliente HTTP usa este DTO até DN-009-2 fechar.
type contadoresDTO struct {
	LinhasProcessadas        int64 `json:"linhasProcessadas"`
	Erros                    int64 `json:"erros"`
	Warnings                 int64 `json:"warnings"`
	Informacao               int64 `json:"informacao"`
	DadosSensiveisMascarados int64 `json:"dadosSensiveisMascarados"`
}
type achadoDTO struct {
	Tipo           domain.ClasseLog  `json:"tipo"`
	Titulo         string            `json:"titulo"`
	Resumo         string            `json:"resumo"`
	NumeroLinha    int               `json:"numeroLinha"`
	CamposExtraido map[string]string `json:"camposExtraido"`
	DadosSensiveis *bool             `json:"dadosSensiveis"`
	TraceID        string            `json:"traceId,omitempty"`
	CorrelationID  string            `json:"correlationId,omitempty"`
	HostName       string            `json:"hostName,omitempty"`
	ModuloID       string            `json:"moduloId,omitempty"`
	Linhas         []string          `json:"linhas"`
	Recorrente     bool              `json:"recorrente"`
}
type resumoDTO struct {
	BatchID             string                     `json:"batchId"`
	ExecID              string                     `json:"execId,omitempty"`
	TenantID            string                     `json:"tenantId"`
	ProductID           string                     `json:"productId"`
	VersaoRegras        string                     `json:"versaoRegras"`
	VersaoCatalogo      string                     `json:"versaoCatalogo,omitempty"`
	IdentificacaoRegras domain.IdentificacaoRegras `json:"identificacaoRegras"`
	Contadores          contadoresDTO              `json:"contadores"`
	Achados             []achadoDTO                `json:"achados"`
}

func SerializarResumo(r domain.ResumoAnalise) ([]byte, error) {
	if err := domain.ValidarResumoHandoff(r); err != nil {
		return nil, err
	}
	c := r.Contadores
	d := resumoDTO{BatchID: r.BatchID, ExecID: r.ExecID, TenantID: r.TenantID, ProductID: r.ProductID, VersaoRegras: r.VersaoRegras, VersaoCatalogo: r.VersaoCatalogo, IdentificacaoRegras: r.IdentificacaoRegras, Contadores: contadoresDTO{c.LinhasProcessadas, c.Erros, c.Warnings, c.Info, c.DadosSensiveisMascarados}, Achados: []achadoDTO{}}
	for _, a := range r.Achados {
		linhas := append([]string{}, a.Linhas...)
		d.Achados = append(d.Achados, achadoDTO{Tipo: a.Tipo, Titulo: a.Titulo, Resumo: a.Resumo, NumeroLinha: a.NumeroLinha, CamposExtraido: a.CamposExtraido, DadosSensiveis: a.DadosSensiveis, TraceID: a.TraceID, CorrelationID: a.CorrelationID, HostName: a.HostName, ModuloID: a.ModuloID, Linhas: linhas, Recorrente: a.Recorrente})
	}
	return json.Marshal(d)
}
