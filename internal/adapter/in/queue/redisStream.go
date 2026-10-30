package queue

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"

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
	rdb := redis.NewClient(&redis.Options{Addr: op.Addr})
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
	lista, err := r.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    r.group,
		Consumer: r.consumidor,
		Streams:  []string{r.stream, ">"},
		Count:    1,
		Block:    0,
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
