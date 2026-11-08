package application

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/in"
	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

func TestExecutarDevolveRegistroInteiro(t *testing.T) {
	conteudo := domain.Registro(`{"senha":"segredo"}`)
	fonte := &fonteFixa{lista: []domain.Codigo{"apontamento-a"}}
	leitura := &leituraFixa{regs: map[domain.Codigo]domain.Registro{"apontamento-a": conteudo}}
	prontos := novoMemProntos()
	s := &ServicoCaptura{Fonte: fonte, Leitura: leitura, Prontos: prontos}
	saida, err := s.Executar(context.Background(), in.PedidoCaptura{InstanciaNoAr: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(saida.Prontos) != 1 || len(saida.Falhas) != 0 {
		t.Fatalf("saída: %+v", saida)
	}
	if !bytes.Equal(saida.Prontos[0].Conteudo, conteudo) {
		t.Fatalf("conteúdo: %q", saida.Prontos[0].Conteudo)
	}
	lido, ok, err := prontos.Ler(context.Background(), "apontamento-a", true)
	if err != nil || !ok || !bytes.Equal(lido, conteudo) {
		t.Fatalf("segunda leitura: %q ok=%v err=%v", lido, ok, err)
	}
}

func TestExecutarSemApontamento(t *testing.T) {
	leitura := &leituraFixa{}
	s := &ServicoCaptura{Fonte: &fonteFixa{}, Leitura: leitura, Prontos: novoMemProntos()}
	saida, err := s.Executar(context.Background(), in.PedidoCaptura{InstanciaNoAr: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(saida.Prontos) != 0 || len(saida.Falhas) != 0 || leitura.chamadas != 0 {
		t.Fatalf("saída=%+v chamadas=%d", saida, leitura.chamadas)
	}
}

func TestExecutarParadaNoMeioNaoDeixaPronto(t *testing.T) {
	leitura := &leituraFixa{regs: map[domain.Codigo]domain.Registro{"apontamento-a": domain.Registro("log")}}
	prontos := novoMemProntos()
	s := &ServicoCaptura{Fonte: &fonteFixa{lista: []domain.Codigo{"apontamento-a"}}, Leitura: leitura, Prontos: prontos}
	saida, err := s.Executar(context.Background(), in.PedidoCaptura{InstanciaNoAr: false})
	if err != nil {
		t.Fatal(err)
	}
	if len(saida.Prontos) != 0 || len(saida.Falhas) != 0 || leitura.chamadas != 0 {
		t.Fatalf("saída=%+v chamadas=%d", saida, leitura.chamadas)
	}
	if _, ok, _ := prontos.Ler(context.Background(), "apontamento-a", true); ok {
		t.Fatal("log ficou pronto")
	}
}

func TestExecutarApontamentoVazioNaoConsulta(t *testing.T) {
	leitura := &leituraFixa{}
	s := &ServicoCaptura{Fonte: &fonteFixa{lista: []domain.Codigo{""}}, Leitura: leitura, Prontos: novoMemProntos()}
	saida, err := s.Executar(context.Background(), in.PedidoCaptura{InstanciaNoAr: true})
	if err != nil {
		t.Fatal(err)
	}
	if leitura.chamadas != 0 || len(saida.Prontos) != 0 || len(saida.Falhas) != 1 {
		t.Fatalf("saída=%+v chamadas=%d", saida, leitura.chamadas)
	}
	if saida.Falhas[0].Motivo != domain.MotivoInutilizavel || saida.Falhas[0].Apontamento != "" {
		t.Fatalf("falha: %+v", saida.Falhas[0])
	}
}

func TestExecutarLogAusente(t *testing.T) {
	leitura := &leituraFixa{estados: map[domain.Codigo]out.EstadoLeitura{"apontamento-b": out.LeituraAusente}}
	s := &ServicoCaptura{Fonte: &fonteFixa{lista: []domain.Codigo{"apontamento-b"}}, Leitura: leitura, Prontos: novoMemProntos()}
	saida, err := s.Executar(context.Background(), in.PedidoCaptura{InstanciaNoAr: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(saida.Prontos) != 0 || len(saida.Falhas) != 1 || saida.Falhas[0].Motivo != domain.MotivoAusente {
		t.Fatalf("saída: %+v", saida)
	}
	if saida.Falhas[0].Apontamento != "apontamento-b" {
		t.Fatalf("apontamento: %+v", saida.Falhas[0])
	}
}

func TestExecutarInacessivel(t *testing.T) {
	leitura := &leituraFixa{estados: map[domain.Codigo]out.EstadoLeitura{"apontamento-c": out.LeituraInacessivel}}
	s := &ServicoCaptura{Fonte: &fonteFixa{lista: []domain.Codigo{"apontamento-c"}}, Leitura: leitura, Prontos: novoMemProntos()}
	saida, err := s.Executar(context.Background(), in.PedidoCaptura{InstanciaNoAr: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(saida.Prontos) != 0 || len(saida.Falhas) != 1 || saida.Falhas[0].Motivo != domain.MotivoInacessivel {
		t.Fatalf("saída: %+v", saida)
	}
	if saida.Falhas[0].Motivo == domain.MotivoAusente {
		t.Fatal("inacessível foi tratado como ausente")
	}
}

func TestExecutarSobrepostas(t *testing.T) {
	leitura := &leituraSobreposta{
		regs:    map[domain.Codigo]domain.Registro{"ok": domain.Registro("log-ok")},
		estados: map[domain.Codigo]out.EstadoLeitura{"falta": out.LeituraAusente},
	}
	prontos := novoMemProntos()
	s := &ServicoCaptura{
		Fonte:   &fonteFixa{lista: []domain.Codigo{"ok", "falta"}},
		Leitura: leitura,
		Prontos: prontos,
	}
	done := make(chan struct{})
	var saida in.SaidaCaptura
	var err error
	go func() {
		saida, err = s.Executar(context.Background(), in.PedidoCaptura{InstanciaNoAr: true})
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("leituras não se sobrepuseram")
	}
	if err != nil {
		t.Fatal(err)
	}
	if !leitura.sobrepos {
		t.Fatal("as duas leituras não estiveram em voo juntas")
	}
	if len(saida.Prontos) != 1 || saida.Prontos[0].Apontamento != "ok" || !bytes.Equal(saida.Prontos[0].Conteudo, domain.Registro("log-ok")) {
		t.Fatalf("prontos: %+v", saida.Prontos)
	}
	if len(saida.Falhas) != 1 || saida.Falhas[0].Motivo != domain.MotivoAusente || saida.Falhas[0].Apontamento != "falta" {
		t.Fatalf("falhas: %+v", saida.Falhas)
	}
	if _, ok, _ := prontos.Ler(context.Background(), "ok", true); !ok {
		t.Fatal("a falha apagou o log pronto")
	}
}

func TestExecutarApontamentoRepetido(t *testing.T) {
	leitura := &leituraFixa{regs: map[domain.Codigo]domain.Registro{"mesmo": domain.Registro("um")}}
	s := &ServicoCaptura{
		Fonte:   &fonteFixa{lista: []domain.Codigo{"mesmo", "mesmo"}},
		Leitura: leitura,
		Prontos: novoMemProntos(),
	}
	saida, err := s.Executar(context.Background(), in.PedidoCaptura{InstanciaNoAr: true})
	if err != nil {
		t.Fatal(err)
	}
	if leitura.chamadas != 1 || len(saida.Prontos) != 1 || saida.Prontos[0].Apontamento != "mesmo" {
		t.Fatalf("saída=%+v chamadas=%d", saida, leitura.chamadas)
	}
}

func TestExecutarErroComEstadoProntoNaoGuarda(t *testing.T) {
	leitura := &leituraFixa{err: errors.New("falha de leitura")}
	prontos := novoMemProntos()
	s := &ServicoCaptura{
		Fonte:   &fonteFixa{lista: []domain.Codigo{"apontamento-a"}},
		Leitura: leitura,
		Prontos: prontos,
	}
	saida, err := s.Executar(context.Background(), in.PedidoCaptura{InstanciaNoAr: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(saida.Prontos) != 0 || len(saida.Falhas) != 1 || saida.Falhas[0].Motivo != domain.MotivoInacessivel {
		t.Fatalf("saída: %+v", saida)
	}
	if saida.Falhas[0].Apontamento != "apontamento-a" {
		t.Fatalf("apontamento: %+v", saida.Falhas[0])
	}
	if _, ok, _ := prontos.Ler(context.Background(), "apontamento-a", true); ok {
		t.Fatal("erro com estado pronto guardou log")
	}
}

func TestExecutarRegistroVazioNaoFicaPronto(t *testing.T) {
	leitura := &leituraFixa{regs: map[domain.Codigo]domain.Registro{
		"nulo":  nil,
		"vazio": {},
	}}
	prontos := novoMemProntos()
	s := &ServicoCaptura{
		Fonte:   &fonteFixa{lista: []domain.Codigo{"nulo", "vazio"}},
		Leitura: leitura,
		Prontos: prontos,
	}
	saida, err := s.Executar(context.Background(), in.PedidoCaptura{InstanciaNoAr: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(saida.Prontos) != 0 || len(saida.Falhas) != 2 {
		t.Fatalf("saída: %+v", saida)
	}
	for _, falha := range saida.Falhas {
		if falha.Motivo != domain.MotivoAusente {
			t.Fatalf("falha: %+v", falha)
		}
	}
	if _, ok, _ := prontos.Ler(context.Background(), "nulo", true); ok {
		t.Fatal("registro nulo ficou pronto")
	}
	if _, ok, _ := prontos.Ler(context.Background(), "vazio", true); ok {
		t.Fatal("registro vazio ficou pronto")
	}
}

func TestExecutarCancelamentoNoMeioDevolveErro(t *testing.T) {
	leitura := &leituraEspera{entrou: make(chan struct{}), solta: make(chan struct{}), reg: domain.Registro("log")}
	prontos := novoMemProntos()
	s := &ServicoCaptura{
		Fonte:   &fonteFixa{lista: []domain.Codigo{"apontamento-a"}},
		Leitura: leitura,
		Prontos: prontos,
	}
	ctx, cancelar := context.WithCancel(context.Background())
	defer cancelar()
	done := make(chan struct{})
	var saida in.SaidaCaptura
	var err error
	go func() {
		saida, err = s.Executar(ctx, in.PedidoCaptura{InstanciaNoAr: true})
		close(done)
	}()
	select {
	case <-leitura.entrou:
	case <-time.After(2 * time.Second):
		t.Fatal("a leitura não começou")
	}
	cancelar()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("a captura não terminou depois do cancelamento")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
	if len(saida.Prontos) != 0 {
		t.Fatalf("prontos: %+v", saida.Prontos)
	}
	if _, ok, _ := prontos.Ler(context.Background(), "apontamento-a", true); ok {
		t.Fatal("cancelamento guardou log")
	}
}

func TestExecutarParadaDuranteLeituraNaoGuarda(t *testing.T) {
	var noAr atomic.Bool
	noAr.Store(true)
	leitura := &leituraEspera{entrou: make(chan struct{}), solta: make(chan struct{}), reg: domain.Registro("log")}
	prontos := novoMemProntos()
	s := &ServicoCaptura{
		Fonte:   &fonteFixa{lista: []domain.Codigo{"apontamento-a"}},
		Leitura: leitura,
		Prontos: prontos,
	}
	done := make(chan struct{})
	var saida in.SaidaCaptura
	var err error
	go func() {
		saida, err = s.Executar(context.Background(), in.PedidoCaptura{
			InstanciaNoAr: true,
			SegueNoAr:     noAr.Load,
		})
		close(done)
	}()
	select {
	case <-leitura.entrou:
	case <-time.After(2 * time.Second):
		t.Fatal("a leitura não começou")
	}
	noAr.Store(false)
	close(leitura.solta)
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("a captura não terminou depois da parada")
	}
	if err != nil {
		t.Fatal(err)
	}
	if len(saida.Prontos) != 0 || len(saida.Falhas) != 0 {
		t.Fatalf("saída: %+v", saida)
	}
	if _, ok, _ := prontos.Ler(context.Background(), "apontamento-a", true); ok {
		t.Fatal("parada durante a leitura guardou log")
	}
}

type fonteFixa struct {
	lista []domain.Codigo
}

func (f *fonteFixa) Obter(context.Context) ([]domain.Codigo, error) {
	return f.lista, nil
}

type leituraFixa struct {
	mu       sync.Mutex
	chamadas int
	regs     map[domain.Codigo]domain.Registro
	estados  map[domain.Codigo]out.EstadoLeitura
	err      error
}

func (f *leituraFixa) Ler(_ context.Context, codigo domain.Codigo) (domain.Registro, out.EstadoLeitura, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.chamadas++
	if f.err != nil {
		return nil, out.LeituraPronta, f.err
	}
	if f.estados != nil {
		if est, ok := f.estados[codigo]; ok {
			return nil, est, nil
		}
	}
	return f.regs[codigo], out.LeituraPronta, nil
}

type leituraEspera struct {
	uma    sync.Once
	entrou chan struct{}
	solta  chan struct{}
	reg    domain.Registro
}

func (f *leituraEspera) Ler(ctx context.Context, _ domain.Codigo) (domain.Registro, out.EstadoLeitura, error) {
	f.uma.Do(func() { close(f.entrou) })
	select {
	case <-f.solta:
		return f.reg, out.LeituraPronta, nil
	case <-ctx.Done():
		return f.reg, out.LeituraPronta, nil
	}
}

type leituraSobreposta struct {
	mu       sync.Mutex
	emVoo    int
	release  chan struct{}
	sobrepos bool
	regs     map[domain.Codigo]domain.Registro
	estados  map[domain.Codigo]out.EstadoLeitura
}

func (f *leituraSobreposta) Ler(_ context.Context, codigo domain.Codigo) (domain.Registro, out.EstadoLeitura, error) {
	f.mu.Lock()
	if f.release == nil {
		f.release = make(chan struct{})
	}
	f.emVoo++
	if f.emVoo == 2 {
		f.sobrepos = true
		close(f.release)
	}
	release := f.release
	f.mu.Unlock()
	<-release
	if est, ok := f.estados[codigo]; ok {
		return nil, est, nil
	}
	return f.regs[codigo], out.LeituraPronta, nil
}

type memProntos struct {
	mu sync.Mutex
	q  *domain.QuadroCaptura
}

func novoMemProntos() *memProntos {
	return &memProntos{q: domain.NovoQuadroCaptura()}
}

func (m *memProntos) Guardar(_ context.Context, codigo domain.Codigo, conteudo domain.Registro) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.q.Concluir(codigo, conteudo, true)
	return nil
}

func (m *memProntos) Ler(_ context.Context, codigo domain.Codigo, instanciaNoAr bool) (domain.Registro, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	r, ok := m.q.Ler(codigo, instanciaNoAr)
	return r, ok, nil
}
