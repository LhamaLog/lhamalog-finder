package queue

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"github.com/redis/go-redis/v9"
)

// ErrNomesAusentes indica que stream, grupo ou o conjunto dos códigos não foram informados.
var ErrNomesAusentes = errors.New("nomes de stream, grupo ou conjunto de códigos ausentes")

// Opcoes liga o adaptador ao Redis. Stream, Group e Codigos vêm de fora. Não há nome padrão.
type Opcoes struct {
	Addr       string
	Stream     string
	Group      string
	Codigos    string
	Consumidor string
}

type fluxo interface {
	Ler(ctx context.Context) (domain.Item, bool, error)
	Ack(ctx context.Context, id domain.ID) error
	ConcluirCodigo(ctx context.Context, id domain.ID, codigo domain.Codigo) (bool, error)
	Marcar(ctx context.Context, codigo domain.Codigo) (bool, error)
	Existe(ctx context.Context, codigo domain.Codigo) (bool, error)
	Fechar() error
}

// Fila implementa as portas da fila, dos códigos concluídos e da disponibilidade.
type Fila struct {
	fluxo fluxo
	mu    sync.Mutex
	porID map[domain.ID]domain.Item
	volta []domain.ID
}

// Nova abre o Redis e prepara o grupo de consumo. Não lê nem apaga mensagem na abertura.
func Nova(op Opcoes) (*Fila, error) {
	if op.Stream == "" || op.Group == "" || op.Codigos == "" {
		return nil, ErrNomesAusentes
	}
	if op.Addr == "" {
		op.Addr = "localhost:6379"
	}
	if op.Consumidor == "" {
		op.Consumidor = consumidorLocal()
	}
	rdb := redis.NewClient(&redis.Options{Addr: op.Addr, ClientName: op.Consumidor})
	fx := &redisFluxo{
		rdb:        rdb,
		stream:     op.Stream,
		group:      op.Group,
		consumidor: op.Consumidor,
		codigos:    op.Codigos,
	}
	if err := fx.garantirGrupo(context.Background()); err != nil {
		_ = rdb.Close()
		return nil, err
	}
	return &Fila{fluxo: fx, porID: map[domain.ID]domain.Item{}}, nil
}

// Fechar encerra o cliente Redis.
func (f *Fila) Fechar() error {
	if f == nil || f.fluxo == nil {
		return nil
	}
	return f.fluxo.Fechar()
}

// Instancia identifica este processo no grupo. Cada processo tem o seu.
func (f *Fila) Instancia() string {
	if r, ok := f.fluxo.(*redisFluxo); ok {
		return r.consumidor
	}
	return ""
}

// TomarLivre devolve um item liberado neste processo ou a próxima mensagem ainda não reservada.
func (f *Fila) TomarLivre(ctx context.Context) (domain.Item, bool, error) {
	f.mu.Lock()
	if len(f.volta) > 0 {
		id := f.volta[0]
		f.volta = f.volta[1:]
		item := f.porID[id]
		f.mu.Unlock()
		return item, true, nil
	}
	f.mu.Unlock()

	item, ok, err := f.fluxo.Ler(ctx)
	if err != nil || !ok {
		return domain.Item{}, false, err
	}
	f.mu.Lock()
	f.porID[item.ID] = item
	f.mu.Unlock()
	return item, true, nil
}

// Concluir confirma a mensagem. Ela deixa de ser entregue como item livre.
func (f *Fila) Concluir(ctx context.Context, id domain.ID) error {
	if err := f.fluxo.Ack(ctx, id); err != nil {
		return err
	}
	f.esquecer(id)
	return nil
}

// ConcluirPrimeira registra o código uma vez e tira o item da pendência na mesma operação.
// Se o código já está no conjunto, só tira o item da pendência.
func (f *Fila) ConcluirPrimeira(ctx context.Context, id domain.ID, codigo domain.Codigo) (bool, error) {
	nova, err := f.fluxo.ConcluirCodigo(ctx, id, codigo)
	if err != nil {
		return false, err
	}
	f.esquecer(id)
	return nova, nil
}

// Liberar devolve o item ao processo, sem conclusão e sem prazo numérico.
func (f *Fila) Liberar(_ context.Context, id domain.ID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.porID[id]; !ok {
		return errors.New("item não reservado")
	}
	f.volta = append(f.volta, id)
	return nil
}

// EncerrarSemConclusao tira a mensagem da reserva sem devolvê-la como livre.
func (f *Fila) EncerrarSemConclusao(ctx context.Context, id domain.ID) error {
	if err := f.fluxo.Ack(ctx, id); err != nil {
		return err
	}
	f.esquecer(id)
	return nil
}

// RegistrarPrimeira aceita o código uma vez.
func (f *Fila) RegistrarPrimeira(ctx context.Context, codigo domain.Codigo) (bool, error) {
	return f.fluxo.Marcar(ctx, codigo)
}

// JaConcluido informa se o código já foi aceito.
func (f *Fila) JaConcluido(ctx context.Context, codigo domain.Codigo) (bool, error) {
	return f.fluxo.Existe(ctx, codigo)
}

// Disponibilizar guarda o código para a etapa futura. Não abre o OpenObserve.
func (f *Fila) Disponibilizar(ctx context.Context, codigo domain.Codigo) error {
	_, err := f.fluxo.Marcar(ctx, codigo)
	return err
}

func (f *Fila) esquecer(id domain.ID) {
	f.mu.Lock()
	delete(f.porID, id)
	f.mu.Unlock()
}

// itemDaMensagem usa o único valor da mensagem como código.
// Zero valores ou mais de um valor não viram um código nomeado.
func itemDaMensagem(id string, values map[string]any) domain.Item {
	item := domain.Item{ID: domain.ID(id), Estado: domain.Livre}
	if len(values) != 1 {
		item.ConteudoExtra = len(values) > 1
		return item
	}
	for _, valor := range values {
		texto, ok := valor.(string)
		if !ok {
			item.ConteudoExtra = true
			return item
		}
		item.Codigo = domain.Codigo(texto)
		return item
	}
	return item
}

type redisFluxo struct {
	rdb        *redis.Client
	stream     string
	group      string
	consumidor string
	codigos    string
}

func (r *redisFluxo) Fechar() error { return r.rdb.Close() }

func (r *redisFluxo) garantirGrupo(ctx context.Context) error {
	if err := r.rdb.Ping(ctx).Err(); err != nil {
		return err
	}
	err := r.rdb.XGroupCreateMkStream(ctx, r.stream, r.group, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return err
	}
	return nil
}

func (r *redisFluxo) Ler(ctx context.Context) (domain.Item, bool, error) {
	item, ok, err := r.lerID(ctx, "0", -1)
	if err != nil || ok {
		return item, ok, err
	}
	item, ok, err = r.reivindicarMaisAntiga(ctx)
	if err != nil || ok {
		return item, ok, err
	}
	return r.lerID(ctx, ">", 0)
}

func (r *redisFluxo) lerID(ctx context.Context, id string, block time.Duration) (domain.Item, bool, error) {
	lista, err := r.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    r.group,
		Consumer: r.consumidor,
		Streams:  []string{r.stream, id},
		Count:    1,
		Block:    block,
	}).Result()
	if err == redis.Nil {
		return domain.Item{}, false, nil
	}
	if err != nil {
		return domain.Item{}, false, err
	}
	if len(lista) == 0 || len(lista[0].Messages) == 0 {
		return domain.Item{}, false, nil
	}
	msg := lista[0].Messages[0]
	return itemDaMensagem(msg.ID, msg.Values), true, nil
}

// reivindicarMaisAntiga reserva a pendência mais antiga de um consumidor sem conexão viva.
// O tempo de ociosidade zero é só o argumento da operação. A regra é a conexão ausente.
func (r *redisFluxo) reivindicarMaisAntiga(ctx context.Context) (domain.Item, bool, error) {
	vivos, err := r.consumidoresVivos(ctx)
	if err != nil {
		return domain.Item{}, false, err
	}
	inicio := "-"
	for {
		pagina, err := r.rdb.XPendingExt(ctx, &redis.XPendingExtArgs{
			Stream: r.stream,
			Group:  r.group,
			Start:  inicio,
			End:    "+",
			Count:  100,
		}).Result()
		if err == redis.Nil || len(pagina) == 0 {
			return domain.Item{}, false, nil
		}
		if err != nil {
			return domain.Item{}, false, err
		}
		for _, p := range pagina {
			if p.ID == inicio || p.Consumer == "" || p.Consumer == r.consumidor || vivos[p.Consumer] {
				continue
			}
			msgs, err := r.rdb.XClaim(ctx, &redis.XClaimArgs{
				Stream:   r.stream,
				Group:    r.group,
				Consumer: r.consumidor,
				MinIdle:  0,
				Messages: []string{p.ID},
			}).Result()
			if err != nil {
				return domain.Item{}, false, err
			}
			if len(msgs) == 0 {
				continue
			}
			return itemDaMensagem(msgs[0].ID, msgs[0].Values), true, nil
		}
		if len(pagina) < 100 {
			return domain.Item{}, false, nil
		}
		inicio = pagina[len(pagina)-1].ID
	}
}

func (r *redisFluxo) consumidoresVivos(ctx context.Context) (map[string]bool, error) {
	texto, err := r.rdb.ClientList(ctx).Result()
	if err != nil {
		return nil, err
	}
	vivos := map[string]bool{}
	for _, linha := range strings.Split(texto, "\n") {
		for _, campo := range strings.Fields(linha) {
			nome, ok := strings.CutPrefix(campo, "name=")
			if ok && nome != "" {
				vivos[nome] = true
			}
		}
	}
	return vivos, nil
}

const scriptConcluir = `
local novo = redis.call('SADD', KEYS[1], ARGV[1])
redis.call('XACK', KEYS[2], ARGV[2], ARGV[3])
return novo
`

func (r *redisFluxo) ConcluirCodigo(ctx context.Context, id domain.ID, codigo domain.Codigo) (bool, error) {
	n, err := r.rdb.Eval(ctx, scriptConcluir, []string{r.codigos, r.stream}, string(codigo), r.group, string(id)).Int64()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

func (r *redisFluxo) Ack(ctx context.Context, id domain.ID) error {
	return r.rdb.XAck(ctx, r.stream, r.group, string(id)).Err()
}

func (r *redisFluxo) Marcar(ctx context.Context, codigo domain.Codigo) (bool, error) {
	n, err := r.rdb.SAdd(ctx, r.codigos, string(codigo)).Result()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

func (r *redisFluxo) Existe(ctx context.Context, codigo domain.Codigo) (bool, error) {
	return r.rdb.SIsMember(ctx, r.codigos, string(codigo)).Result()
}

func consumidorLocal() string {
	nome, err := os.Hostname()
	if err != nil || nome == "" {
		nome = "finder"
	}
	return nome + "-" + strconv.Itoa(os.Getpid())
}

// ErrParametrosConsumo indica que o consumo de log.ingested não recebeu os parâmetros externos.
var ErrParametrosConsumo = errors.New("parâmetros de consumo ausentes ou não positivos")

// posicaoInicialGrupo é o início do que o stream ainda guarda. A posição "$" não é usada.
const posicaoInicialGrupo = "0"

// OpcoesLote liga a leitura de log.ingested. Nenhum número, grupo, consumidor ou prefixo é preenchido aqui.
type OpcoesLote struct {
	Addr         string
	Senha        string
	Stream       string
	Group        string
	Consumidor   string
	Idempotencia string
	Count        int64
	Block        time.Duration
	Ociosidade   time.Duration
}

// Lote lê o envelope, guarda a chave do lote e confirma pelo id Redis.
type Lote struct {
	rdb        *redis.Client
	stream     string
	group      string
	consumidor string
	prefixo    string
	count      int64
	block      time.Duration
	ociosidade time.Duration
	mu         sync.Mutex
	buffer     []domain.Entrada
	adiados    map[string]struct{}
}

var _ out.Entradas = (*Lote)(nil)
var _ out.ChaveLote = (*Lote)(nil)

// NovaLote abre o Redis e cria o grupo em 0 se ele ainda não existe.
// Não lê nem confirma entrada na abertura. Não gera consumidor nem número.
func NovaLote(op OpcoesLote) (*Lote, error) {
	if op.Stream == "" || op.Group == "" || op.Consumidor == "" || op.Idempotencia == "" || op.Count <= 0 || op.Block <= 0 || op.Ociosidade <= 0 {
		return nil, ErrParametrosConsumo
	}
	if op.Addr == "" {
		op.Addr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:       op.Addr,
		Password:   op.Senha,
		ClientName: op.Consumidor,
	})
	lote := &Lote{
		rdb:        rdb,
		stream:     op.Stream,
		group:      op.Group,
		consumidor: op.Consumidor,
		prefixo:    op.Idempotencia,
		count:      op.Count,
		block:      op.Block,
		ociosidade: op.Ociosidade,
		adiados:    map[string]struct{}{},
	}
	if err := lote.garantirGrupo(context.Background()); err != nil {
		_ = rdb.Close()
		return nil, err
	}
	return lote, nil
}

// Fechar encerra o cliente Redis.
func (l *Lote) Fechar() error {
	if l == nil || l.rdb == nil {
		return nil
	}
	return l.rdb.Close()
}

// Ler entrega uma entrada deste grupo. A pendência já deste consumidor não espera o bloqueio.
// Entrada adiada nesta passagem não volta antes da leitura bloqueante seguinte.
func (l *Lote) Ler(ctx context.Context) (domain.Entrada, bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if entrada, ok := l.daFila(); ok {
		return entrada, true, nil
	}
	if err := l.puxarPendentes(ctx); err != nil {
		return domain.Entrada{}, false, err
	}
	if entrada, ok := l.daFila(); ok {
		return entrada, true, nil
	}
	if err := l.reivindicar(ctx); err != nil {
		return domain.Entrada{}, false, err
	}
	if entrada, ok := l.daFila(); ok {
		return entrada, true, nil
	}
	if err := l.puxarNovas(ctx); err != nil {
		return domain.Entrada{}, false, err
	}
	l.adiados = map[string]struct{}{}
	if entrada, ok := l.daFila(); ok {
		return entrada, true, nil
	}
	return domain.Entrada{}, false, nil
}

// Confirmar tira a pendência desse id. Não usa idempotencyId.
func (l *Lote) Confirmar(ctx context.Context, id string) error {
	if err := l.rdb.XAck(ctx, l.stream, l.group, id).Err(); err != nil {
		return err
	}
	l.mu.Lock()
	delete(l.adiados, id)
	l.mu.Unlock()
	return nil
}

// Adiar mantém a entrada pendente e fora do restante desta passagem.
func (l *Lote) Adiar(_ context.Context, id string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.adiados == nil {
		l.adiados = map[string]struct{}{}
	}
	l.adiados[id] = struct{}{}
	destino := l.buffer[:0]
	for _, entrada := range l.buffer {
		if entrada.ID != id {
			destino = append(destino, entrada)
		}
	}
	l.buffer = destino
	return nil
}

// EncerrarInvalida confirma a entrada e remove a reserva em processamento.
// Não grava o estado concluído.
func (l *Lote) EncerrarInvalida(ctx context.Context, id, idempotencyID string) error {
	if idempotencyID != "" {
		if err := l.rdb.Eval(ctx, scriptRemoverReserva, []string{nomeDaChave(l.prefixo, idempotencyID)}, domain.EstadoEmProcessamento).Err(); err != nil {
			return err
		}
	}
	return l.Confirmar(ctx, id)
}

// Reservar cria em_processamento só se a chave não existe. Sem prazo.
func (l *Lote) Reservar(ctx context.Context, idempotencyID string) (bool, error) {
	return l.rdb.SetNX(ctx, nomeDaChave(l.prefixo, idempotencyID), domain.EstadoEmProcessamento, 0).Result()
}

// Concluir grava concluido e 25 horas. Não renova o prazo se a chave já está concluída.
func (l *Lote) Concluir(ctx context.Context, idempotencyID string) error {
	segundos := int64(domain.PrazoConcluido / time.Second)
	return l.rdb.Eval(ctx, scriptConcluirChave, []string{nomeDaChave(l.prefixo, idempotencyID)}, domain.EstadoConcluido, segundos).Err()
}

// Consultar devolve ausente, em_processamento ou concluido.
func (l *Lote) Consultar(ctx context.Context, idempotencyID string) (string, error) {
	valor, err := l.rdb.Get(ctx, nomeDaChave(l.prefixo, idempotencyID)).Result()
	if err == redis.Nil {
		return domain.EstadoAusente, nil
	}
	if err != nil {
		return "", err
	}
	switch valor {
	case domain.EstadoEmProcessamento, domain.EstadoConcluido:
		return valor, nil
	default:
		return domain.EstadoAusente, nil
	}
}

// RemoverReserva apaga a chave só quando ela está em processamento.
func (l *Lote) RemoverReserva(ctx context.Context, idempotencyID string) error {
	return l.rdb.Eval(ctx, scriptRemoverReserva, []string{nomeDaChave(l.prefixo, idempotencyID)}, domain.EstadoEmProcessamento).Err()
}

func nomeDaChave(prefixo, idempotencyID string) string {
	return prefixo + idempotencyID
}

func (l *Lote) garantirGrupo(ctx context.Context) error {
	if err := l.rdb.Ping(ctx).Err(); err != nil {
		return err
	}
	err := l.rdb.XGroupCreateMkStream(ctx, l.stream, l.group, posicaoInicialGrupo).Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return err
	}
	return nil
}

func (l *Lote) daFila() (domain.Entrada, bool) {
	if len(l.buffer) == 0 {
		return domain.Entrada{}, false
	}
	entrada := l.buffer[0]
	l.buffer = l.buffer[1:]
	return entrada, true
}

func (l *Lote) puxarPendentes(ctx context.Context) error {
	msgs, err := l.lerGrupo(ctx, "0", -1)
	if err != nil || len(msgs) == 0 {
		return err
	}
	pendentes, err := l.rdb.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream:   l.stream,
		Group:    l.group,
		Start:    "-",
		End:      "+",
		Count:    l.count,
		Consumer: l.consumidor,
	}).Result()
	if err == redis.Nil {
		pendentes = nil
		err = nil
	}
	if err != nil {
		return err
	}
	ociosidade := map[string]time.Duration{}
	for _, pendente := range pendentes {
		ociosidade[pendente.ID] = pendente.Idle
	}
	for _, msg := range msgs {
		if _, adiada := l.adiados[msg.ID]; adiada {
			continue
		}
		l.buffer = append(l.buffer, domain.Entrada{
			ID:     msg.ID,
			Campos: camposTexto(msg.Values),
			Ociosa: ociosidade[msg.ID] >= l.ociosidade,
		})
	}
	return nil
}

func (l *Lote) reivindicar(ctx context.Context) error {
	msgs, _, err := l.rdb.XAutoClaim(ctx, &redis.XAutoClaimArgs{
		Stream:   l.stream,
		Group:    l.group,
		Consumer: l.consumidor,
		MinIdle:  l.ociosidade,
		Start:    "0-0",
		Count:    l.count,
	}).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}
	for _, msg := range msgs {
		if _, adiada := l.adiados[msg.ID]; adiada {
			continue
		}
		l.buffer = append(l.buffer, domain.Entrada{
			ID:           msg.ID,
			Campos:       camposTexto(msg.Values),
			Ociosa:       true,
			Reivindicada: true,
		})
	}
	return nil
}

func (l *Lote) puxarNovas(ctx context.Context) error {
	msgs, err := l.lerGrupo(ctx, ">", l.block)
	if err != nil || len(msgs) == 0 {
		return err
	}
	for _, msg := range msgs {
		if _, adiada := l.adiados[msg.ID]; adiada {
			continue
		}
		l.buffer = append(l.buffer, domain.Entrada{
			ID:     msg.ID,
			Campos: camposTexto(msg.Values),
		})
	}
	return nil
}

func (l *Lote) lerGrupo(ctx context.Context, id string, block time.Duration) ([]redis.XMessage, error) {
	lista, err := l.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    l.group,
		Consumer: l.consumidor,
		Streams:  []string{l.stream, id},
		Count:    l.count,
		Block:    block,
	}).Result()
	if err == redis.Nil || len(lista) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return lista[0].Messages, nil
}

func camposTexto(valores map[string]any) map[string]string {
	campos := make(map[string]string, len(valores))
	for chave, valor := range valores {
		texto, ok := valor.(string)
		if !ok {
			continue
		}
		campos[chave] = texto
	}
	return campos
}

const scriptConcluirChave = `
local atual = redis.call('GET', KEYS[1])
if atual == ARGV[1] then
  return 0
end
redis.call('SET', KEYS[1], ARGV[1], 'EX', ARGV[2])
return 1
`

const scriptRemoverReserva = `
if redis.call('GET', KEYS[1]) == ARGV[1] then
  return redis.call('DEL', KEYS[1])
end
return 0
`
