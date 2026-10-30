package application

import (
	"context"
	"errors"
	"testing"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

func TestExecutarConcluiEDisponibilizaUmaVez(t *testing.T) {
	fila := novaFila(domain.Item{ID: "1", Codigo: "apontamento-a", Estado: domain.Livre})
	codigos := &memCodigos{}
	disp := &memDisp{}
	s := &Servico{Fila: fila, Codigos: codigos, Disp: disp, Reservas: domain.NovoQuadro(nil)}
	saida, err := s.Executar(context.Background(), in.Pedido{
		Instancia: "a", WorkerNoAr: true, InstanciaNoAr: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if saida.Conclusoes != 1 || len(saida.CodigosDisponiveis) != 1 || saida.CodigosDisponiveis[0] != "apontamento-a" {
		t.Fatalf("saída: %+v", saida)
	}
	if len(disp.codigos) != 1 || fila.concluidos["1"] != true {
		t.Fatalf("disp=%v concluidos=%v", disp.codigos, fila.concluidos)
	}
	if fila.livres["1"] {
		t.Fatal("item concluído voltou a livre")
	}
}

func TestExecutarSegundoCodigoNaoConclui(t *testing.T) {
	fila := novaFila(
		domain.Item{ID: "1", Codigo: "mesmo", Estado: domain.Livre},
		domain.Item{ID: "2", Codigo: "mesmo", Estado: domain.Livre},
	)
	codigos := &memCodigos{}
	disp := &memDisp{}
	s := &Servico{Fila: fila, Codigos: codigos, Disp: disp, Reservas: domain.NovoQuadro(nil)}
	pedido := in.Pedido{Instancia: "a", WorkerNoAr: true, InstanciaNoAr: true}
	if _, err := s.Executar(context.Background(), pedido); err != nil {
		t.Fatal(err)
	}
	segunda, err := s.Executar(context.Background(), pedido)
	if err != nil {
		t.Fatal(err)
	}
	if segunda.Conclusoes != 0 || len(disp.codigos) != 1 {
		t.Fatalf("segunda=%+v disp=%v", segunda, disp.codigos)
	}
	if fila.reservados["2"] || fila.livres["2"] || fila.concluidos["2"] {
		t.Fatalf("estado do segundo: livres=%v reservados=%v concluidos=%v encerrados=%v",
			fila.livres, fila.reservados, fila.concluidos, fila.encerrados)
	}
	if !fila.encerrados["2"] {
		t.Fatal("segundo item não foi encerrado sem conclusão")
	}
}

func TestExecutarParadaDoWorkerLibera(t *testing.T) {
	fila := novaFila(domain.Item{ID: "1", Codigo: "apontamento-a", Estado: domain.Livre})
	disp := &memDisp{}
	s := &Servico{Fila: fila, Codigos: &memCodigos{}, Disp: disp, Reservas: domain.NovoQuadro(nil)}
	saida, err := s.Executar(context.Background(), in.Pedido{
		Instancia: "a", WorkerNoAr: false, InstanciaNoAr: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if saida.Conclusoes != 0 || len(disp.codigos) != 0 || !fila.livres["1"] {
		t.Fatalf("saída=%+v livres=%v", saida, fila.livres)
	}
}

func TestExecutarFilaVaziaEFalhaDeLeitura(t *testing.T) {
	vazia := novaFila()
	s := &Servico{Fila: vazia, Codigos: &memCodigos{}, Disp: &memDisp{}, Reservas: domain.NovoQuadro(nil)}
	saida, err := s.Executar(context.Background(), in.Pedido{Instancia: "a", WorkerNoAr: true, InstanciaNoAr: true})
	if err != nil || saida.Conclusoes != 0 {
		t.Fatalf("vazia: %+v err=%v", saida, err)
	}
	falha := &filaFake{err: errors.New("falha de leitura")}
	s.Fila = falha
	saida, err = s.Executar(context.Background(), in.Pedido{Instancia: "a", WorkerNoAr: true, InstanciaNoAr: true})
	if err == nil || !saida.FalhaLeitura || saida.Conclusoes != 0 {
		t.Fatalf("falha: %+v err=%v", saida, err)
	}
}

func TestExecutarConteudoExtraNaoConclui(t *testing.T) {
	fila := novaFila(domain.Item{ID: "1", Codigo: "apontamento-a", Estado: domain.Livre, ConteudoExtra: true})
	disp := &memDisp{}
	s := &Servico{Fila: fila, Codigos: &memCodigos{}, Disp: disp, Reservas: domain.NovoQuadro(nil)}
	saida, err := s.Executar(context.Background(), in.Pedido{Instancia: "a", WorkerNoAr: true, InstanciaNoAr: true})
	if err != nil || !saida.ForaDeContrato || saida.Conclusoes != 0 || len(disp.codigos) != 0 {
		t.Fatalf("saída=%+v disp=%v err=%v", saida, disp.codigos, err)
	}
}

type filaFake struct {
	espera     []domain.Item
	livres     map[domain.ID]bool
	reservados map[domain.ID]bool
	concluidos map[domain.ID]bool
	encerrados map[domain.ID]bool
	err        error
}

func novaFila(itens ...domain.Item) *filaFake {
	f := &filaFake{
		livres:     map[domain.ID]bool{},
		reservados: map[domain.ID]bool{},
		concluidos: map[domain.ID]bool{},
		encerrados: map[domain.ID]bool{},
	}
	for _, item := range itens {
		f.espera = append(f.espera, item)
		f.livres[item.ID] = true
	}
	return f
}

func (f *filaFake) TomarLivre(context.Context) (domain.Item, bool, error) {
	if f.err != nil {
		return domain.Item{}, false, f.err
	}
	if len(f.espera) == 0 {
		return domain.Item{}, false, nil
	}
	item := f.espera[0]
	f.espera = f.espera[1:]
	delete(f.livres, item.ID)
	f.reservados[item.ID] = true
	return item, true, nil
}

func (f *filaFake) Concluir(_ context.Context, id domain.ID) error {
	return f.mover(id, func() { f.concluidos[id] = true })
}

func (f *filaFake) Liberar(_ context.Context, id domain.ID) error {
	return f.mover(id, func() { f.livres[id] = true })
}

func (f *filaFake) EncerrarSemConclusao(_ context.Context, id domain.ID) error {
	return f.mover(id, func() { f.encerrados[id] = true })
}

func (f *filaFake) mover(id domain.ID, marcar func()) error {
	if !f.reservados[id] {
		return errors.New("item não reservado")
	}
	delete(f.reservados, id)
	marcar()
	return nil
}

type memCodigos struct {
	vistos map[domain.Codigo]bool
}

func (m *memCodigos) RegistrarPrimeira(_ context.Context, codigo domain.Codigo) (bool, error) {
	if m.vistos == nil {
		m.vistos = map[domain.Codigo]bool{}
	}
	if m.vistos[codigo] {
		return false, nil
	}
	m.vistos[codigo] = true
	return true, nil
}

func (m *memCodigos) JaConcluido(_ context.Context, codigo domain.Codigo) (bool, error) {
	return m.vistos[codigo], nil
}

type memDisp struct {
	codigos []domain.Codigo
}

func (m *memDisp) Disponibilizar(_ context.Context, codigo domain.Codigo) error {
	m.codigos = append(m.codigos, codigo)
	return nil
}
