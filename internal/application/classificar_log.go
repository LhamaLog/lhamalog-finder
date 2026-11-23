package application

import (
	"errors"
	entrada "github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// ClassificarLogCapturado recebe a linha segura separada; não lê Mensagem do log bruto.
func ClassificarLogCapturado(ref domain.Referencia, log entrada.LogAnalise, linha domain.LinhaLogMascarada, versao string, catalogo []domain.RegraClassificacao) (domain.LogClassificado, error) {
	if log.Apontamento == "" {
		return domain.LogClassificado{}, errors.New("log_ausente")
	}
	if ref.TenantID == "" || ref.ProductID == "" {
		return domain.LogClassificado{}, errors.New("identidade_incompleta")
	}
	if log.ServiceName == "" || log.Host == "" {
		return domain.LogClassificado{}, errors.New("assinatura_incompleta")
	}
	r := domain.ClassificarLog(log.NivelCampo, linha.Texto, catalogo)
	return domain.LogClassificado{Apontamento: log.Apontamento, Classe: r.Classe, TenantID: ref.TenantID, ProductID: ref.ProductID, CatalogoUsado: r.CatalogoUsado, ServiceName: log.ServiceName, Host: log.Host, NivelCampo: log.NivelCampo, Mensagem: linha.Texto, VersaoRegras: versao, Timestamp: log.Timestamp, TraceID: log.TraceID, ModuloID: log.ModuloID}, nil
}
