package queue

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

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

func TestNovaLoteRecusaParametroAusente(t *testing.T) {
	_, err := NovaLote(OpcoesLote{Senha: "segredo-nao-logar"})
	if !errors.Is(err, ErrParametrosConsumo) {
		t.Fatalf("err = %v", err)
	}
	if err != nil && strings.Contains(err.Error(), "segredo-nao-logar") {
		t.Fatal("erro contém senha")
	}
	_, err = NovaLote(OpcoesLote{
		Stream:       "lhamalog:stream:logs",
		Group:        "g",
		Consumidor:   "c",
		Idempotencia: "p:",
		Count:        0,
		Block:        time.Second,
		Ociosidade:   time.Second,
	})
	if !errors.Is(err, ErrParametrosConsumo) {
		t.Fatalf("quantidade zero: %v", err)
	}
}

func TestNomeDaChaveUsaPrefixoExterno(t *testing.T) {
	finder := nomeDaChave("finder:log-ingested:", "batch")
	varredura := nomeDaChave("varredura:", "batch")
	if finder == varredura {
		t.Fatal("prefixos diferentes produziram a mesma chave")
	}
	if finder != "finder:log-ingested:batch" {
		t.Fatalf("chave = %s", finder)
	}
	if strings.Contains(finder, "finder:codigos") {
		t.Fatal("a chave do lote usou o conjunto de códigos")
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

func (f *fluxoParado) ConcluirCodigo(context.Context, domain.ID, domain.Codigo) (bool, error) {
	return true, nil
}

func (f *fluxoParado) Marcar(context.Context, domain.Codigo) (bool, error) { return true, nil }

func (f *fluxoParado) Existe(context.Context, domain.Codigo) (bool, error) { return false, nil }

func (f *fluxoParado) Fechar() error { return nil }
