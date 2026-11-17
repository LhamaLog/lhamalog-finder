package labeleddataset

import "github.com/LhamaLog/lhamalog-finder/internal/domain"

type ClassLabel = domain.ClasseLog
type CamposEntrada struct {
	ItemType    string `json:"itemType"`
	LevelText   string `json:"level_text"`
	MessageText string `json:"message_text"`
	BatchID     string `json:"lhamalog_batch_id"`
	RecordID    string `json:"lhamalog_record_id"`
	TenantID    string `json:"tenantId"`
	ProductID   string `json:"productId"`
	Split       string `json:"split,omitempty"`
	Notas       string `json:"notas,omitempty"`
}

// Members não possuem campo gold; o decoder estrito rejeita até class_label:null.
type Registro struct {
	CamposEntrada
	ClassLabel   ClassLabel      `json:"class_label"`
	Members      []CamposEntrada `json:"members,omitempty"`
	Documentacao string          `json:"documentacao,omitempty"`
}
type RegistroLogRotulado = Registro
type RegistroApontamentoRotulado = Registro
