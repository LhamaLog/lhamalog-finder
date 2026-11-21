package domain

import "time"

// ResumoAnalise é o contrato lógico em memória. O adaptador 009 define o JSON.
// Não contém documento bruto OpenObserve, gravidade nem entidade de ocorrência.
type ResumoAnalise struct {
	BatchID             string              `json:"batchId"`
	ExecID              string              `json:"execId,omitempty"`
	TenantID            string              `json:"tenantId"`
	ProductID           string              `json:"productId"`
	VersaoRegras        string              `json:"versaoRegras"`
	IdentificacaoRegras IdentificacaoRegras `json:"identificacaoRegras"`
	Contadores          ContadoresResumo    `json:"contadores"`
	Achados             []AchadoAnalise     `json:"achados"`
	VersaoCatalogo      string              `json:"versaoCatalogo,omitempty"`
}
type ContadoresResumo struct {
	LinhasProcessadas        int64 `json:"linhasProcessadas"`
	Erros                    int64 `json:"erros"`
	Warnings                 int64 `json:"warnings"`
	Info                     int64 `json:"info"`
	DadosSensiveisMascarados int64 `json:"dadosSensiveisMascarados"`
}

// Campos de agrupamento são serializados pelo adaptador 009.
// Titulo/Resumo/NumeroLinha/CamposExtraido não recebem conteúdo inventado.
type AchadoAnalise struct {
	Tipo               ClasseLog         `json:"tipo"`
	Titulo             string            `json:"titulo"`
	Resumo             string            `json:"resumo"`
	NumeroLinha        int               `json:"numeroLinha"`
	CamposExtraido     map[string]string `json:"camposExtraido"`
	DadosSensiveis     *bool             `json:"dadosSensiveis"`
	TraceID            string            `json:"traceId,omitempty"`
	CorrelationID      string            `json:"correlationId,omitempty"`
	HostName           string            `json:"hostName,omitempty"`
	ModuloID           string            `json:"moduloId,omitempty"`
	Linhas             []string          `json:"linhas"`
	Recorrente         bool              `json:"recorrente"`
	Assinatura         string            `json:"-"`
	Contagem           int64             `json:"-"`
	PrimeiraOcorrencia time.Time         `json:"-"`
	UltimaOcorrencia   time.Time         `json:"-"`
	VersaoRegras       string            `json:"-"`
	StatusAnterior     string            `json:"-"`
	RevisaoConsultada  *string           `json:"-"`
}
