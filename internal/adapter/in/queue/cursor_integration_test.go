package queue

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestReivindicacaoAvancaPeloCursorEReiniciaNovaPassagem(t *testing.T) {
	for _, quantidade := range []int64{1, 2} {
		t.Run(fmt.Sprintf("count-%d", quantidade), func(t *testing.T) {
			op := provaLote(t)
			op.Count = quantidade
			op.Ociosidade = 500 * time.Millisecond
			op.Block = time.Millisecond
			lote := abrirLote(t, op)
			rdb := clienteSenha(t, op.Addr, op.Senha)
			ctx, cancelar := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancelar()
			prefixo := int(10 * quantidade)
			ids := make([]string, prefixo+1)
			for i := range ids {
				ids[i] = publicarID(t, op, camposIntegracao("batch-cursor"))
			}
			for _, consumidor := range []struct {
				nome  string
				count int64
			}{{"ativo", int64(prefixo)}, {"abandonado", 1}} {
				_, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{Group: op.Group, Consumer: consumidor.nome, Streams: []string{op.Stream, ">"}, Count: consumidor.count}).Result()
				if err != nil {
					t.Fatal(err)
				}
			}
			time.Sleep(2 * op.Ociosidade)
			recuperou := false
			for tentativa := 0; tentativa < 3; tentativa++ {
				// Outro consumidor mantém o prefixo ativo durante a varredura.
				if _, err := rdb.XClaim(ctx, &redis.XClaimArgs{Stream: op.Stream, Group: op.Group, Consumer: "ativo", Messages: ids[:prefixo]}).Result(); err != nil {
					t.Fatal(err)
				}
				entrada, ok, err := lote.Ler(ctx)
				if err != nil {
					t.Fatal(err)
				}
				if !ok {
					continue
				}
				if entrada.ID != ids[prefixo] || !entrada.Ociosa || !entrada.Reivindicada {
					t.Fatalf("reivindicacao incorreta: %+v", entrada)
				}
				if err := lote.Confirmar(ctx, entrada.ID); err != nil {
					t.Fatal(err)
				}
				recuperou = true
				break
			}
			if !recuperou {
				t.Fatal("a varredura nao alcancou a pendencia apos o prefixo ativo")
			}
			// Ao terminar a passagem, uma entrada anterior que ficou ociosa
			// precisa voltar a ser elegível na próxima passagem.
			time.Sleep(2 * op.Ociosidade)
			entrada, ok, err := lote.Ler(ctx)
			if err != nil || !ok || entrada.ID != ids[0] || !entrada.Reivindicada {
				t.Fatalf("nova passagem nao voltou ao inicio: entrada=%+v ok=%v erro=%v", entrada, ok, err)
			}
		})
	}
}
