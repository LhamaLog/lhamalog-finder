package queue

import (
	"context"
	"testing"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

func TestItemDaMensagem(t *testing.T) {
	um := itemDaMensagem("1-0", map[string]any{"qualquer": "apontamento-a"})
	if !um.Aceito() || um.Codigo != "apontamento-a" || um.ID != "1-0" {
		t.Fatalf("um valor: %+v", um)
	}
	vazio := itemDaMensagem("2-0", map[string]any{})
	if vazio.Aceito() || vazio.ConteudoExtra {
		t.Fatalf("sem valor: %+v", vazio)
	}
	extra := itemDaMensagem("3-0", map[string]any{"a": "1", "b": "2"})
	if extra.Aceito() || !extra.ConteudoExtra {
		t.Fatalf("dois valores: %+v", extra)
	}
}

func TestLiberarDevolveOMesmoItem(t *testing.T) {
	f := &Fila{
		fluxo: &fluxoParado{item: domain.Item{ID: "9-0", Codigo: "apontamento-a", Estado: domain.Livre}},
		porID: map[domain.ID]domain.Item{},
	}
	primeiro, ok, err := f.TomarLivre(context.Background())
	if err != nil || !ok {
		t.Fatalf("tomar: %+v ok=%v err=%v", primeiro, ok, err)
	}
	if err := f.Liberar(context.Background(), primeiro.ID); err != nil {
		t.Fatal(err)
	}
	segundo, ok, err := f.TomarLivre(context.Background())
	if err != nil || !ok || segundo.ID != primeiro.ID {
		t.Fatalf("liberado: %+v ok=%v err=%v", segundo, ok, err)
	}
	if f.fluxo.(*fluxoParado).leituras != 1 {
		t.Fatalf("leituras = %d", f.fluxo.(*fluxoParado).leituras)
	}
}

func TestNovaSemNomes(t *testing.T) {
	if _, err := Nova(Opcoes{Addr: "localhost:6379"}); err != ErrNomesAusentes {
		t.Fatalf("err = %v", err)
	}
}

type fluxoParado struct {
	item     domain.Item
	leituras int
}

func (f *fluxoParado) Ler(context.Context) (domain.Item, bool, error) {
	f.leituras++
	return f.item, true, nil
}

func (f *fluxoParado) Ack(context.Context, domain.ID) error { return nil }

func (f *fluxoParado) Marcar(context.Context, domain.Codigo) (bool, error) { return true, nil }

func (f *fluxoParado) Existe(context.Context, domain.Codigo) (bool, error) { return false, nil }

func (f *fluxoParado) Fechar() error { return nil }
