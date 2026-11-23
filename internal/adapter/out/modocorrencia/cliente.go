package modocorrencia

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

const (
	rotaConsulta = "/internal/v1/ocorrencias/consultas"
	rotaEntrega  = "/internal/v1/resumos-analise"
)

// OpcoesCliente recebe a URL e a credencial do ambiente. O cliente não emite JWT.
type OpcoesCliente struct {
	URLBase         string
	JWT             string
	TimeoutConsulta time.Duration
	TimeoutEntrega  time.Duration
	Transporte      http.RoundTripper
}

type Cliente struct {
	base            string
	jwt             string
	http            http.Client
	TimeoutConsulta time.Duration
	TimeoutEntrega  time.Duration
	consultaApos    atomic.Pointer[time.Time]
	entregaApos     atomic.Pointer[time.Time]
}

var _ out.ModuloOcorrencia = (*Cliente)(nil)

func NovoCliente(o OpcoesCliente) (*Cliente, error) {
	base, err := url.Parse(o.URLBase)
	if err != nil || base.Scheme != "https" || base.Hostname() == "" || base.User != nil || base.RawQuery != "" || base.ForceQuery || base.Fragment != "" || (base.Path != "" && base.Path != "/") {
		return nil, errors.New("url_modulo_ocorrencia_invalida")
	}
	if o.JWT == "" || strings.ContainsAny(o.JWT, " \t\r\n\x00") {
		return nil, errors.New("jwt_modulo_ocorrencia_invalido")
	}
	if o.TimeoutConsulta < 0 || o.TimeoutEntrega < 0 {
		return nil, errors.New("timeout_modulo_ocorrencia_invalido")
	}
	if o.TimeoutConsulta == 0 {
		o.TimeoutConsulta = 2 * time.Second
	}
	if o.TimeoutEntrega == 0 {
		o.TimeoutEntrega = 5 * time.Second
	}
	return &Cliente{
		base: "https://" + base.Host, jwt: o.JWT,
		TimeoutConsulta: o.TimeoutConsulta, TimeoutEntrega: o.TimeoutEntrega,
		http: http.Client{Transport: o.Transporte, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }},
	}, nil
}

func (c *Cliente) ConsultarAchado(ctx context.Context, chave domain.IdentidadeAchado) ([]domain.EstadoOcorrencia, error) {
	if !chave.Classe.Valida() || algumVazio(chave.TenantID, chave.ProductID, chave.Assinatura) {
		return nil, out.ErrModuloRequisicao
	}
	corpo, err := json.Marshal(struct {
		TenantID   string           `json:"tenantId"`
		ProductID  string           `json:"productId"`
		Tipo       domain.ClasseLog `json:"tipo"`
		Assinatura string           `json:"assinatura"`
	}{chave.TenantID, chave.ProductID, chave.Classe, chave.Assinatura})
	if err != nil {
		return nil, out.ErrModuloRequisicao
	}
	var estados []domain.EstadoOcorrencia
	err = c.post(ctx, rotaConsulta, c.TimeoutConsulta, &c.consultaApos, corpo, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return erroHTTP(resp, &c.consultaApos)
		}
		var dto struct {
			Estado json.RawMessage `json:"estado"`
		}
		if lerObjeto(resp.Body, &dto, "estado") != nil || len(dto.Estado) == 0 {
			return out.ErrModuloProtocolo
		}
		if bytes.Equal(dto.Estado, []byte("null")) {
			return nil
		}
		estado, err := lerEstado(dto.Estado)
		if err != nil {
			return err
		}
		estados = []domain.EstadoOcorrencia{estado}
		return nil
	})
	return estados, err
}

func lerEstado(raw json.RawMessage) (domain.EstadoOcorrencia, error) {
	var dto struct {
		Contagem int64  `json:"contagemAcumulada"`
		Primeira string `json:"primeiraOcorrencia"`
		Ultima   string `json:"ultimaOcorrencia"`
		Status   string `json:"status"`
		Revisao  string `json:"revisao"`
	}
	if lerObjeto(bytes.NewReader(raw), &dto, "contagemAcumulada", "primeiraOcorrencia", "ultimaOcorrencia", "status", "revisao") != nil || dto.Contagem < 1 || algumVazio(dto.Revisao) {
		return domain.EstadoOcorrencia{}, out.ErrModuloProtocolo
	}
	switch dto.Status {
	case "ABERTA", "EM_ANALISE", "FECHADA":
	default:
		return domain.EstadoOcorrencia{}, out.ErrModuloProtocolo
	}
	primeira, err1 := time.Parse(time.RFC3339Nano, dto.Primeira)
	ultima, err2 := time.Parse(time.RFC3339Nano, dto.Ultima)
	if err1 != nil || err2 != nil || !strings.HasSuffix(dto.Primeira, "Z") || !strings.HasSuffix(dto.Ultima, "Z") || ultima.Before(primeira) {
		return domain.EstadoOcorrencia{}, out.ErrModuloProtocolo
	}
	return domain.EstadoOcorrencia{ContagemAcumulada: dto.Contagem, PrimeiraGlobal: primeira, UltimaGlobal: ultima, Status: dto.Status, Revisao: dto.Revisao}, nil
}

func (c *Cliente) Enviar(ctx context.Context, resumo domain.ResumoAnalise) (bool, error) {
	corpo, err := SerializarResumo(resumo)
	if err != nil {
		return false, err
	}
	err = c.post(ctx, rotaEntrega, c.TimeoutEntrega, &c.entregaApos, corpo, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			return erroHTTP(resp, &c.entregaApos)
		}
		return validarConfirmacao(resp, resumo)
	})
	return err == nil, err
}

func validarConfirmacao(resp *http.Response, r domain.ResumoAnalise) error {
	var recibo struct {
		Aceito    bool            `json:"aceito"`
		ReciboID  string          `json:"reciboId"`
		TenantID  string          `json:"tenantId"`
		ProductID string          `json:"productId"`
		BatchID   string          `json:"batchId"`
		ExecID    json.RawMessage `json:"execId"`
		Repeticao *bool           `json:"repeticao"`
	}
	if lerObjeto(resp.Body, &recibo, "aceito", "reciboId", "tenantId", "productId", "batchId", "execId", "repeticao") != nil || !recibo.Aceito || algumVazio(recibo.ReciboID) || recibo.Repeticao == nil || *recibo.Repeticao != (resp.StatusCode == http.StatusOK) || recibo.TenantID != r.TenantID || recibo.ProductID != r.ProductID || recibo.BatchID != r.BatchID {
		return out.ErrModuloProtocolo
	}
	if r.ExecID == "" {
		if recibo.ExecID != nil {
			return out.ErrModuloProtocolo
		}
	} else {
		var exec string
		if json.Unmarshal(recibo.ExecID, &exec) != nil || exec != r.ExecID {
			return out.ErrModuloProtocolo
		}
	}
	return nil
}

func (c *Cliente) post(ctx context.Context, rota string, timeout time.Duration, apos *atomic.Pointer[time.Time], corpo []byte, ler func(*http.Response) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if limite := apos.Load(); limite != nil && time.Now().Before(*limite) {
		return out.ErrModuloTemporario
	}
	prazo, cancelar := context.WithTimeout(ctx, timeout)
	defer cancelar()
	req, err := http.NewRequestWithContext(prazo, http.MethodPost, c.base+rota, bytes.NewReader(corpo))
	if err != nil {
		return out.ErrModuloRequisicao
	}
	req.Header.Set("Authorization", "Bearer "+c.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return out.ErrModuloTemporario
	}
	defer resp.Body.Close()
	err = ler(resp)
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if prazo.Err() != nil {
		return out.ErrModuloTemporario
	}
	return err
}

func erroHTTP(resp *http.Response, apos *atomic.Pointer[time.Time]) error {
	switch resp.StatusCode {
	case http.StatusConflict:
		var dto struct {
			Codigo string `json:"codigo"`
		}
		if lerObjeto(resp.Body, &dto, "codigo") != nil {
			return out.ErrModuloProtocolo
		}
		switch dto.Codigo {
		case "REVISAO_DESATUALIZADA":
			return out.ErrRevisaoDesatualizada
		case "ENTREGA_EM_ANDAMENTO":
		default:
			return out.ErrModuloProtocolo
		}
	case http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
	case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusUnprocessableEntity:
		return out.ErrModuloRequisicao
	default:
		return out.ErrModuloProtocolo
	}
	guardarRetryAfter(resp.Header.Get("Retry-After"), apos)
	return out.ErrModuloTemporario
}

func guardarRetryAfter(valor string, apos *atomic.Pointer[time.Time]) {
	agora := time.Now()
	var limite time.Time
	segundos, err := strconv.ParseUint(strings.TrimSpace(valor), 10, 63)
	if err == nil && segundos <= uint64((1<<63-1)/time.Second) {
		limite = agora.Add(time.Duration(segundos) * time.Second)
	} else {
		limite, _ = http.ParseTime(valor)
	}
	if !limite.After(agora) {
		return
	}
	for {
		anterior := apos.Load()
		if anterior != nil && !limite.After(*anterior) {
			return
		}
		if apos.CompareAndSwap(anterior, &limite) {
			return
		}
	}
}

// lerObjeto exige um único objeto com nomes exatos e sem duplicatas.
func lerObjeto(r io.Reader, destino any, campos ...string) error {
	dec := json.NewDecoder(r)
	inicio, err := dec.Token()
	if err != nil || inicio != json.Delim('{') {
		return out.ErrModuloProtocolo
	}
	obj := make(map[string]json.RawMessage, len(campos))
	for dec.More() {
		chave, err := dec.Token()
		nome, ok := chave.(string)
		if err != nil || !ok || obj[nome] != nil {
			return out.ErrModuloProtocolo
		}
		permitido := false
		for _, campo := range campos {
			if campo == nome {
				permitido = true
				break
			}
		}
		if !permitido {
			return out.ErrModuloProtocolo
		}
		var valor json.RawMessage
		if dec.Decode(&valor) != nil {
			return out.ErrModuloProtocolo
		}
		obj[nome] = valor
	}
	fim, err := dec.Token()
	if err != nil || fim != json.Delim('}') {
		return out.ErrModuloProtocolo
	}
	if _, err := dec.Token(); err != io.EOF {
		return out.ErrModuloProtocolo
	}
	dados, err := json.Marshal(obj)
	if err != nil || json.Unmarshal(dados, destino) != nil {
		return out.ErrModuloProtocolo
	}
	return nil
}

func algumVazio(valores ...string) bool {
	for _, v := range valores {
		if strings.TrimSpace(v) == "" {
			return true
		}
	}
	return false
}
