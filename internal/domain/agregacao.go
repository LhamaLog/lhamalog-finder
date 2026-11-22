package domain

// Valor nil significa desconhecido; falha de consulta nunca produz false conhecido.
type MarcaRecorrente struct{ Valor *bool }
type EvidenciaRecorrencia struct {
	BatchID, Host string
	Chave         IdentidadeAchado
}
