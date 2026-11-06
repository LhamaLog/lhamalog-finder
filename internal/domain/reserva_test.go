package domain

import "testing"

func TestReservaExclusivaELiberacaoSemPrazo(t *testing.T) {
	item := Item{ID: "1", Codigo: "apontamento-a", Estado: Livre}
	q := NovoQuadro([]Item{item})

	primeiro, ok := q.Tomar("instancia-a")
	if !ok {
		t.Fatal("primeira reserva")
	}
	if _, ok := q.Tomar("instancia-b"); ok {
		t.Fatal("o mesmo item foi reservado por duas instâncias")
	}
	dono, ok := q.Dono(primeiro.ID)
	if !ok || dono != "instancia-a" {
		t.Fatalf("dono = %q ok=%v", dono, ok)
	}

	soltos := q.SoltarInstancia("instancia-a")
	if len(soltos) != 1 || soltos[0].Estado != Livre {
		t.Fatalf("soltos: %+v", soltos)
	}
	if q.Estado(primeiro.ID) == Concluido {
		t.Fatal("parada concluiu o item")
	}
	segundo, ok := q.Tomar("instancia-b")
	if !ok || segundo.ID != primeiro.ID {
		t.Fatalf("depois da parada: %+v ok=%v", segundo, ok)
	}
}

func TestSegundoTomadorRecebeOutroItem(t *testing.T) {
	q := NovoQuadro([]Item{
		{ID: "1", Codigo: "a", Estado: Livre},
		{ID: "2", Codigo: "b", Estado: Livre},
	})
	a, ok := q.Tomar("instancia-a")
	if !ok {
		t.Fatal("tomada a")
	}
	b, ok := q.Tomar("instancia-b")
	if !ok {
		t.Fatal("tomada b")
	}
	if a.ID == b.ID {
		t.Fatal("as duas instâncias receberam o mesmo item")
	}
	if _, ok := q.Dono(a.ID); !ok {
		t.Fatal("item em andamento sem dono")
	}
	if _, ok := q.Dono(b.ID); !ok {
		t.Fatal("segundo item sem dono")
	}
}
