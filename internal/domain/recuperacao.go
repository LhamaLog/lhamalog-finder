package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const streamRedis = "lhamalog:stream:logs"

// ErrFiltro indica valor que não pode entrar só como dado do filtro.
// O texto não traz o valor recusado.
var ErrFiltro = errors.New("filtro inválido")

// Classificacao é o desfecho único da recuperação.
type Classificacao string

const (
	ClasseLoteCompleto      Classificacao = "lote completo"
	ClasseAutenticacao      Classificacao = "falha de autenticação ou autorização"
	ClasseRespostaInvalida  Classificacao = "resposta inválida"
	ClasseIsolamento        Classificacao = "inconsistência de identidade ou isolamento"
	ClasseSemIdentificador  Classificacao = "documento sem identificador"
	ClasseDivergente        Classificacao = "conteúdo divergente"
	ClasseExcesso           Classificacao = "excesso de registros distintos"
	ClasseIncompleto        Classificacao = "lote incompleto"
	ClasseIndisponibilidade Classificacao = "indisponibilidade ou timeout"
	ClasseFalhaPagina       Classificacao = "falha no meio da paginação"
	ClasseConfigAusente     Classificacao = "configuração de acesso ausente"
)

// FalhaAcesso é o problema da consulta, quando houve um.
type FalhaAcesso int

const (
	AcessoNenhum FalhaAcesso = iota
	AcessoCredencial
	AcessoResposta
	AcessoTimeout
	AcessoIndisponivel
	AcessoPagina
)

// ConfiguracaoAcesso é o que a execução informa. Nada aqui nasce com padrão no código.
// CamposCopia nil é lista ausente. Ponteiro para lista vazia é lista presente.
type ConfiguracaoAcesso struct {
	URL           string
	Usuario       string
	Senha         string
	Timeout       time.Duration
	TamanhoPagina int
	CamposCopia   *[]string
	Ordem         string
	CampoTenant   string
	CampoProduto  string
}

// Pronta informa se a recuperação pode consultar.
// Nomes de tenant e produto ficam de fora: a falta deles é isolamento, não configuração ausente.
func (c ConfiguracaoAcesso) Pronta() bool {
	if strings.TrimSpace(c.URL) == "" || c.Usuario == "" || c.Senha == "" {
		return false
	}
	if c.Timeout <= 0 || c.TamanhoPagina <= 0 || c.CamposCopia == nil {
		return false
	}
	return strings.TrimSpace(c.Ordem) != ""
}

// DocumentoFisico é um item de hits. O conteúdo permanece inteiro nesta etapa.
type DocumentoFisico map[string]any

// ResultadoRecuperacao é a classificação, a marca de efeito e o contexto do lote.
// Falha não carrega registros. A senha não tem campo aqui.
type ResultadoRecuperacao struct {
	Classificacao     Classificacao
	EfeitoConcluido   bool
	TenantID          string
	ProductID         string
	Organizacao       string
	Stream            string
	BatchID           string
	IdempotencyID     string
	CorrelationID     string
	EntradaID         string
	RecordCount       int64
	Encontrado        int64
	ContagemCalculada bool
	Registros         []DocumentoFisico
}

// Avaliacao reúne o que a leitura observou antes de escolher uma classificação.
type Avaliacao struct {
	Encontrado        int64
	ContagemCalculada bool
	SemIdentificador  bool
	Divergente        bool
	Isolamento        bool
	Acesso            FalhaAcesso
	PaginasLidas      int
}

// NovoResultadoRecuperacao copia o contexto da referência. Não copia log nem senha.
func NovoResultadoRecuperacao(ref Referencia) ResultadoRecuperacao {
	return ResultadoRecuperacao{
		TenantID:      ref.TenantID,
		ProductID:     ref.ProductID,
		Organizacao:   ref.Organizacao,
		Stream:        ref.Stream,
		BatchID:       ref.BatchID,
		IdempotencyID: ref.IdempotencyID,
		CorrelationID: ref.CorrelationID,
		EntradaID:     ref.EntradaID,
		RecordCount:   ref.RecordCount,
	}
}

// Fechar escolhe uma classificação e só entrega registros no lote completo.
func Fechar(ref Referencia, aval Avaliacao, registros []DocumentoFisico) ResultadoRecuperacao {
	classe, efeito := classificar(ref.RecordCount, aval)
	saida := NovoResultadoRecuperacao(ref)
	saida.Classificacao = classe
	saida.EfeitoConcluido = efeito
	if aval.ContagemCalculada {
		saida.ContagemCalculada = true
		saida.Encontrado = aval.Encontrado
	}
	if classe == ClasseLoteCompleto {
		saida.Registros = registros
	}
	return saida
}

func classificar(esperado int64, aval Avaliacao) (Classificacao, bool) {
	switch {
	case aval.Acesso == AcessoCredencial:
		return ClasseAutenticacao, true
	case aval.Acesso == AcessoResposta:
		return ClasseRespostaInvalida, true
	case aval.Isolamento:
		return ClasseIsolamento, true
	case aval.SemIdentificador:
		return ClasseSemIdentificador, true
	case aval.Divergente:
		return ClasseDivergente, true
	case aval.ContagemCalculada && aval.Encontrado > esperado:
		return ClasseExcesso, true
	case aval.ContagemCalculada && aval.Encontrado < esperado:
		return ClasseIncompleto, false
	case aval.Acesso == AcessoPagina && aval.PaginasLidas > 0:
		return ClasseFalhaPagina, false
	case aval.Acesso == AcessoTimeout || aval.Acesso == AcessoIndisponivel || aval.Acesso == AcessoPagina:
		return ClasseIndisponibilidade, false
	case aval.ContagemCalculada && aval.Encontrado == esperado:
		return ClasseLoteCompleto, true
	default:
		return ClasseIndisponibilidade, false
	}
}

// ConsultaSQL monta o texto da busca e a faixa. Não escolhe coluna de ordem.
// start é minTimestampMicros. fim é maxTimestampMicros + 1.
func ConsultaSQL(ref Referencia, cfg ConfiguracaoAcesso) (string, int64, int64, error) {
	if ref.Stream == "" || ref.Stream == streamRedis {
		return "", 0, 0, ErrFiltro
	}
	stream, err := identificar(ref.Stream)
	if err != nil {
		return "", 0, 0, err
	}
	if !organizacaoAplicavel(ref.Organizacao) {
		return "", 0, 0, ErrFiltro
	}
	if ref.MaxTimestampMicros == math.MaxInt64 {
		return "", 0, 0, ErrFiltro
	}
	lote, err := literalObrigatorio(ref.BatchID)
	if err != nil {
		return "", 0, 0, err
	}
	partes := []string{aspas("lhamalog_batch_id") + " = " + lote}
	if ref.RecordID != nil {
		registro, err := literalObrigatorio(*ref.RecordID)
		if err != nil {
			return "", 0, 0, err
		}
		partes = append(partes, aspas("lhamalog_record_id")+" = "+registro)
	}
	if cfg.CampoTenant == "" || cfg.CampoProduto == "" {
		return "", 0, 0, ErrFiltro
	}
	campoTenant, err := identificar(cfg.CampoTenant)
	if err != nil {
		return "", 0, 0, err
	}
	campoProduto, err := identificar(cfg.CampoProduto)
	if err != nil {
		return "", 0, 0, err
	}
	tenant, err := literalObrigatorio(ref.TenantID)
	if err != nil {
		return "", 0, 0, err
	}
	produto, err := literalObrigatorio(ref.ProductID)
	if err != nil {
		return "", 0, 0, err
	}
	partes = append(partes, campoTenant+" = "+tenant, campoProduto+" = "+produto)
	ordem := strings.TrimSpace(cfg.Ordem)
	if ordem == "" || strings.ContainsAny(ordem, ";\r\n\x00") {
		return "", 0, 0, ErrFiltro
	}
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s ORDER BY %s", stream, strings.Join(partes, " AND "), ordem)
	return sql, ref.MinTimestampMicros, ref.MaxTimestampMicros + 1, nil
}

func aspas(nome string) string {
	return `"` + nome + `"`
}

func identificar(nome string) (string, error) {
	if nome == "" || strings.ContainsAny(nome, "\"\x00\r\n;") {
		return "", ErrFiltro
	}
	return aspas(nome), nil
}

func literalObrigatorio(valor string) (string, error) {
	if valor == "" || strings.IndexByte(valor, 0) >= 0 {
		return "", ErrFiltro
	}
	if inteiroSQL(valor) {
		return valor, nil
	}
	return "'" + strings.ReplaceAll(valor, "'", "''") + "'", nil
}

func inteiroSQL(valor string) bool {
	if valor == "0" {
		return true
	}
	corpo := valor
	if strings.HasPrefix(valor, "-") {
		corpo = valor[1:]
	}
	if corpo == "" || corpo[0] == '0' {
		return false
	}
	for _, r := range corpo {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func organizacaoAplicavel(org string) bool {
	if org == "" || org == "." || org == ".." {
		return false
	}
	return !strings.ContainsAny(org, "/\\?#%\r\n\x00")
}

// Deduplicar agrupa por lhamalog_record_id.
// Diferença só nos campos da lista entrega uma das cópias.
// Diferença fora da lista não escolhe versão. Documento sem identificador fica de fora da contagem.
func Deduplicar(docs []DocumentoFisico, copia []string) (registros []DocumentoFisico, divergente, semIdentificador bool) {
	ignorar := make(map[string]struct{}, len(copia))
	for _, nome := range copia {
		ignorar[nome] = struct{}{}
	}
	ordem := make([]string, 0)
	grupos := map[string][]DocumentoFisico{}
	for _, doc := range docs {
		id, ok := IdentificadorRegistro(doc)
		if !ok {
			semIdentificador = true
			continue
		}
		if _, existe := grupos[id]; !existe {
			ordem = append(ordem, id)
		}
		grupos[id] = append(grupos[id], doc)
	}
	registros = make([]DocumentoFisico, 0, len(ordem))
	for _, id := range ordem {
		grupo := grupos[id]
		if !copiasCompativeis(grupo, ignorar) {
			divergente = true
			continue
		}
		registros = append(registros, grupo[0])
	}
	return registros, divergente, semIdentificador
}

func copiasCompativeis(grupo []DocumentoFisico, ignorar map[string]struct{}) bool {
	if len(grupo) <= 1 {
		return true
	}
	base := semCampos(grupo[0], ignorar)
	for _, doc := range grupo[1:] {
		if !reflect.DeepEqual(base, semCampos(doc, ignorar)) {
			return false
		}
	}
	return true
}

func semCampos(doc DocumentoFisico, ignorar map[string]struct{}) DocumentoFisico {
	saida := make(DocumentoFisico, len(doc))
	for chave, valor := range doc {
		if _, ok := ignorar[chave]; ok {
			continue
		}
		saida[chave] = valor
	}
	return saida
}

// IdentificadorRegistro lê lhamalog_record_id. Vazio ou ausente não conta.
func IdentificadorRegistro(doc DocumentoFisico) (string, bool) {
	valor, ok := doc["lhamalog_record_id"]
	if !ok || valor == nil {
		return "", false
	}
	switch item := valor.(type) {
	case string:
		if item == "" {
			return "", false
		}
		return item, true
	case json.Number:
		if item.String() == "" {
			return "", false
		}
		return item.String(), true
	case int:
		return strconv.Itoa(item), true
	case int64:
		return strconv.FormatInt(item, 10), true
	case float64:
		if item != float64(int64(item)) {
			return "", false
		}
		return strconv.FormatInt(int64(item), 10), true
	default:
		return "", false
	}
}

// Isolado informa se algum documento não confere com o tenant e o produto da referência.
// A conferência usa só os nomes informados, no topo do documento.
func Isolado(docs []DocumentoFisico, campoTenant, campoProduto, tenant, produto string) bool {
	if campoTenant == "" || campoProduto == "" {
		return true
	}
	for _, doc := range docs {
		if !pertence(doc, campoTenant, campoProduto, tenant, produto) {
			return true
		}
	}
	return false
}

// Proprios devolve os documentos que conferem com a referência.
func Proprios(docs []DocumentoFisico, campoTenant, campoProduto, tenant, produto string) []DocumentoFisico {
	saida := make([]DocumentoFisico, 0, len(docs))
	for _, doc := range docs {
		if pertence(doc, campoTenant, campoProduto, tenant, produto) {
			saida = append(saida, doc)
		}
	}
	return saida
}

func pertence(doc DocumentoFisico, campoTenant, campoProduto, tenant, produto string) bool {
	valorTenant, okTenant := doc[campoTenant]
	valorProduto, okProduto := doc[campoProduto]
	if !okTenant || !okProduto {
		return false
	}
	return mesmoTexto(tenant, valorTenant) && mesmoTexto(produto, valorProduto)
}

func mesmoTexto(esperado string, valor any) bool {
	switch item := valor.(type) {
	case string:
		return item == esperado
	case json.Number:
		return item.String() == esperado
	case int:
		return strconv.Itoa(item) == esperado
	case int64:
		return strconv.FormatInt(item, 10) == esperado
	case float64:
		return strconv.FormatInt(int64(item), 10) == esperado && item == float64(int64(item))
	default:
		return false
	}
}
