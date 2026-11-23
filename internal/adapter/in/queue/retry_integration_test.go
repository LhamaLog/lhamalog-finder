package queue

import (
	"context"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/application"
)

func TestPendenciaPropriaPreservaElegibilidadeDeRetomada(t *testing.T) {
	op := provaLote(t)
	op.Count = 1
	op.Ociosidade = 80 * time.Millisecond
	lote := abrirLote(t, op)
	id := publicarID(t, op, camposIntegracao("batch-1"))
	ctx, cancelar := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelar()
	if entrada, ok, err := lote.Ler(ctx); err != nil || !ok || entrada.ID != id {
		t.Fatalf("primeira leitura: entrada=%+v ok=%v erro=%v", entrada, ok, err)
	}
	time.Sleep(2 * op.Ociosidade)
	entrada, ok, err := lote.Ler(ctx)
	if err != nil || !ok || entrada.ID != id || !entrada.Ociosa {
		t.Fatalf("retomada perdeu ociosidade: entrada=%+v ok=%v erro=%v", entrada, ok, err)
	}
}

func TestFalhaTemporariaRetomaNoMesmoConsumidorAntesDeConfirmar(t *testing.T) {
	op := provaLote(t)
	op.Count = 1
	op.Ociosidade = 80 * time.Millisecond
	op.Block = 20 * time.Millisecond
	lote := abrirLote(t, op)
	id := publicarID(t, op, camposIntegracao("batch-1"))
	efeito := &efeitoIntegracao{}
	consumo := &application.ConsumoReferencia{Fila: lote, Chaves: lote, Efeito: efeito}
	ctx, cancelar := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelar()
	saida, err := consumo.Executar(ctx)
	if err != nil || saida.Confirmacoes != 0 || len(efeito.refs) != 1 {
		t.Fatalf("falha temporária: saída=%+v erro=%v efeitos=%d", saida, err, len(efeito.refs))
	}
	if pendencias(t, op.Addr, op.Stream, op.Group) != 1 {
		t.Fatal("confirmou antes de concluir o efeito")
	}
	efeito.concluido = true
	for saida.Confirmacoes == 0 && ctx.Err() == nil {
		saida, err = consumo.Executar(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}
	if saida.Confirmacoes != 1 || len(efeito.refs) != 2 || efeito.refs[1].EntradaID != id {
		t.Fatalf("retomada: saída=%+v efeitos=%d", saida, len(efeito.refs))
	}
	if pendencias(t, op.Addr, op.Stream, op.Group) != 0 {
		t.Fatal("efeito concluído deixou a entrada pendente")
	}
}
