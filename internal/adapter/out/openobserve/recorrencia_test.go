package openobserve

import (
	"context"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"testing"
)

type historicoTeste struct {
	chamadas   int
	temporario bool
}

func (h *historicoTeste) ConsultarHistorico(context.Context, out.PedidoRecorrencia) ([]domain.EvidenciaRecorrencia, error) {
	h.chamadas++
	return nil, out.FalhaRecorrencia{Temporaria: h.temporario}
}

type avisoTeste struct {
	chamadas int
	batch    string
}

func (a *avisoTeste) AvisarFalhaOpenObserve(ctx context.Context, batch, motivo string) error {
	a.chamadas++
	a.batch = batch
	return nil
}
func TestRecorrenciaTresTentativasSemFalsoNegativo(t *testing.T) {
	h := &historicoTeste{temporario: true}
	a := &avisoTeste{}
	s := Recorrencia{Historico: h, Aviso: a}
	marca, err := s.ConsultarRecorrencia(context.Background(), out.PedidoRecorrencia{BatchID: "batch", Host: "h"})
	if err == nil || marca.Valor != nil || h.chamadas != 3 || a.chamadas != 1 || a.batch != "batch" {
		t.Fatal("falha virou recorrente=false ou nao avisou porta")
	}
	h = &historicoTeste{temporario: false}
	s.Historico = h
	_, err = s.ConsultarRecorrencia(context.Background(), out.PedidoRecorrencia{BatchID: "batch", Host: "h"})
	if err == nil || h.chamadas != 1 {
		t.Fatal("erro critico repetido")
	}
}

type historicoPonteiro struct{ chamadas int }

func (h *historicoPonteiro) ConsultarHistorico(context.Context, out.PedidoRecorrencia) ([]domain.EvidenciaRecorrencia, error) {
	h.chamadas++
	return nil, &out.FalhaRecorrencia{Temporaria: true}
}
func TestErroTemporarioPonteiroTambemRepeteTres(t *testing.T) {
	h := &historicoPonteiro{}
	a := &avisoTeste{}
	s := Recorrencia{Historico: h, Aviso: a}
	_, err := s.ConsultarRecorrencia(context.Background(), out.PedidoRecorrencia{BatchID: "b"})
	if err == nil || h.chamadas != 3 || a.chamadas != 1 {
		t.Fatal("falha temporaria por ponteiro tratada como critica")
	}
}
