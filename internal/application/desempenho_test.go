package application

import (
	"context"
	"fmt"
	entrada "github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"os"
	"testing"
	"time"
)

func TestAceiteDesempenhoVMA(t *testing.T) {
	if os.Getenv("FINDER_ACEITE_VMA") != "1" {
		t.Skip("SC-002: requer VM-A real, todos os servicos e lote de 1000; nao executado neste ambiente")
	}
	// Operador só habilita após comprovar o perfil do quickstart §7b. A flag não
	// comprova hardware ou serviços; evidência operacional deve acompanhar o teste.
	p := pedidoAnaliseTeste()
	p.Logs = make([]entrada.LogAnalise, 1000)
	p.Referencia.RecordCount = 1000
	for i := range p.Logs {
		p.Logs[i] = entrada.LogAnalise{Apontamento: domain.Codigo(fmt.Sprint(i)), NivelCampo: "INFO", ServiceName: "svc", Host: "h", Timestamp: time.Unix(10, 0)}
	}
	n := 0
	s := ServicoAnalise{Seguranca: scannerAnaliseTeste(&n)}
	inicio := time.Now()
	_, err := s.Montar(context.Background(), p)
	duracao := time.Since(inicio)
	if err != nil || duracao > 2*time.Second {
		t.Fatalf("ciclo local: %s erro=%v", duracao, err)
	}
}
