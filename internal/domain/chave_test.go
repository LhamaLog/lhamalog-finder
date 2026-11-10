package domain

import (
	"testing"
	"time"
)

func TestChaveReservaConcluiENaoRenova(t *testing.T) {
	agora := time.Date(2026, 11, 8, 12, 0, 0, 0, time.UTC)
	ausente := Chave{}
	if !ausente.Ausente() || ausente.Suprime(agora) {
		t.Fatalf("ausente: %+v", ausente)
	}
	reservada, ok := ausente.Reservar()
	if !ok || reservada.Estado != EstadoEmProcessamento || !reservada.ExpiraEm.IsZero() {
		t.Fatalf("reserva: %+v ok=%v", reservada, ok)
	}
	if _, ok := reservada.Reservar(); ok {
		t.Fatal("reserva repetida")
	}
	concluida := reservada.Concluir(agora)
	if concluida.Estado != EstadoConcluido || !concluida.ExpiraEm.Equal(agora.Add(PrazoConcluido)) {
		t.Fatalf("conclusão: %+v", concluida)
	}
	if !concluida.Suprime(agora.Add(PrazoConcluido - time.Nanosecond)) {
		t.Fatal("concluída não suprime dentro de 25 horas")
	}
	if concluida.Suprime(agora.Add(PrazoConcluido)) {
		t.Fatal("concluída ainda suprime depois de 25 horas")
	}
	renovada := concluida.Concluir(agora.Add(time.Hour))
	if !renovada.ExpiraEm.Equal(concluida.ExpiraEm) || renovada.Estado != EstadoConcluido {
		t.Fatalf("prazo renovado: %+v", renovada)
	}
	removida := reservada.RemoverReserva()
	if removida.Estado != EstadoAusente || removida.Estado == EstadoConcluido || !removida.ExpiraEm.IsZero() {
		t.Fatalf("remoção: %+v", removida)
	}
	if concluida.RemoverReserva().Estado != EstadoConcluido {
		t.Fatal("remover a reserva gravou outro estado na concluída")
	}
}
