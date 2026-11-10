package domain

import "time"

const (
	// EstadoAusente é a chave que não reserva nem suprime.
	EstadoAusente = "ausente"
	// EstadoEmProcessamento é a reserva. Não tem prazo e não é conclusão.
	EstadoEmProcessamento = "em_processamento"
	// EstadoConcluido é o efeito da recuperação terminado.
	EstadoConcluido = "concluido"
	// PrazoConcluido é contado uma vez, na passagem para concluído.
	PrazoConcluido = 25 * time.Hour
)

// Chave é o estado de um idempotencyId no prefixo do Finder.
// ExpiraEm só existe no estado concluído.
type Chave struct {
	Estado   string
	ExpiraEm time.Time
}

// Ausente informa que não há reserva nem conclusão vigente.
func (c Chave) Ausente() bool {
	return c.Estado == "" || c.Estado == EstadoAusente
}

// Reservar cria em_processamento só quando a chave está ausente. Não define prazo.
func (c Chave) Reservar() (Chave, bool) {
	if !c.Ausente() {
		return c, false
	}
	return Chave{Estado: EstadoEmProcessamento}, true
}

// Concluir grava o estado concluído e o prazo de 25 horas.
// Se a chave já está concluída, o prazo permanece o da primeira passagem.
func (c Chave) Concluir(agora time.Time) Chave {
	if c.Estado == EstadoConcluido {
		return c
	}
	return Chave{Estado: EstadoConcluido, ExpiraEm: agora.Add(PrazoConcluido)}
}

// Suprime é verdadeiro só enquanto o estado concluído não completou 25 horas.
func (c Chave) Suprime(agora time.Time) bool {
	if c.Estado != EstadoConcluido || c.ExpiraEm.IsZero() {
		return false
	}
	return agora.Before(c.ExpiraEm)
}

// RemoverReserva apaga a reserva. Não grava concluído.
func (c Chave) RemoverReserva() Chave {
	if c.Estado != EstadoEmProcessamento {
		return c
	}
	return Chave{Estado: EstadoAusente}
}
