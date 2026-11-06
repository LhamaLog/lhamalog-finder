package domain

import "testing"

func itemLivre(id, codigo string) Item {
	return Item{ID: ID(id), Codigo: Codigo(codigo), Estado: Livre}
}

func TestFimConcluiNoMesmoPasso(t *testing.T) {
	item, ok := Iniciar(itemLivre("1", "apontamento-a"))
	if !ok {
		t.Fatal("iniciar")
	}
	res := Fim(item, true, true, false)
	if !res.Concluido || !res.Disponibilizar || res.Item.Estado != Concluido {
		t.Fatalf("conclusão: %+v", res)
	}
}

func TestFimNaoConcluiSeCodigoJaConcluido(t *testing.T) {
	item, ok := Iniciar(itemLivre("2", "apontamento-a"))
	if !ok {
		t.Fatal("iniciar")
	}
	res := Fim(item, true, true, true)
	if res.Concluido || res.Disponibilizar || res.Item.Estado != NaoConcluido {
		t.Fatalf("segunda conclusão: %+v", res)
	}
	if res.LivreDeNovo {
		t.Fatal("item repetido voltou a livre")
	}
}

func TestFimParadaDoWorkerNaoConclui(t *testing.T) {
	item, ok := Iniciar(itemLivre("3", "apontamento-b"))
	if !ok {
		t.Fatal("iniciar")
	}
	res := Fim(item, false, true, false)
	if res.Concluido || !res.LivreDeNovo || res.Item.Estado != Livre {
		t.Fatalf("parada do worker: %+v", res)
	}
}

func TestConcluirJuntosUmaConclusao(t *testing.T) {
	saida := ConcluirJuntos([]Item{
		itemLivre("4", "mesmo"),
		itemLivre("5", "mesmo"),
	})
	conclusoes := 0
	for _, res := range saida {
		if res.Concluido {
			conclusoes++
		}
		if res.Item.Estado == EmPasso {
			t.Fatalf("reserva restante: %+v", res)
		}
	}
	if conclusoes != 1 {
		t.Fatalf("conclusões = %d", conclusoes)
	}
}

func TestItemForaDeContrato(t *testing.T) {
	semCodigo := Item{ID: "6", Estado: Livre}
	if semCodigo.Aceito() {
		t.Fatal("código vazio foi aceito")
	}
	extra := Item{ID: "7", Codigo: "apontamento-c", Estado: Livre, ConteudoExtra: true}
	res := Fim(extra, true, true, false)
	if res.Concluido || !res.ForaDeContrato {
		t.Fatalf("conteúdo extra: %+v", res)
	}
}

func TestSemConsumoNaoConclui(t *testing.T) {
	vazio := SemConsumo(false)
	falha := SemConsumo(true)
	if vazio.Concluido || falha.Concluido || !falha.FalhaLeitura {
		t.Fatalf("vazio=%+v falha=%+v", vazio, falha)
	}
}
