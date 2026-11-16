package application

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/LhamaLog/lhamalog-finder/internal/application/port/out"
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

func TestPaginaCurtaEntregaLoteENaoConfirmaRedis(t *testing.T) {
	ref := referenciaTeste()
	ref.RecordCount = 1
	acesso := acessoTeste()
	doc := documentoTeste("a", "linha-sintetica")
	busca := &buscaFake{paginas: [][]domain.DocumentoFisico{{doc}}}
	svc := &ServicoRecuperacao{Busca: busca, Acesso: acesso}

	resultado, err := svc.Recuperar(context.Background(), ref)
	if err != nil {
		t.Fatal(err)
	}
	if resultado.Classificacao != domain.ClasseLoteCompleto || !resultado.EfeitoConcluido {
		t.Fatalf("classe=%s efeito=%v", resultado.Classificacao, resultado.EfeitoConcluido)
	}
	if resultado.TenantID != ref.TenantID || resultado.ProductID != ref.ProductID || resultado.Organizacao != ref.Organizacao || resultado.Stream != ref.Stream || resultado.BatchID != ref.BatchID || resultado.IdempotencyID != ref.IdempotencyID {
		t.Fatal("contexto do lote incompleto")
	}
	if resultado.RecordCount != 1 || resultado.Encontrado != 1 || len(resultado.Registros) != 1 {
		t.Fatalf("contagem record=%d encontrado=%d n=%d", resultado.RecordCount, resultado.Encontrado, len(resultado.Registros))
	}
	if resultado.Registros[0]["mensagem"] != "linha-sintetica" {
		t.Fatal("o conteúdo não veio da busca")
	}
	if len(busca.chamadas) != 1 {
		t.Fatalf("chamadas=%d", len(busca.chamadas))
	}
	pedido := busca.chamadas[0]
	if pedido.From != 0 || pedido.Size != acesso.TamanhoPagina || pedido.Size == int(ref.RecordCount) {
		t.Fatalf("página from=%d size=%d", pedido.From, pedido.Size)
	}
	if pedido.Inicio != ref.MinTimestampMicros || pedido.Fim != ref.MaxTimestampMicros+1 {
		t.Fatalf("faixa %d..%d", pedido.Inicio, pedido.Fim)
	}
	if !strings.Contains(pedido.SQL, `"stream_sintetico"`) || strings.Contains(pedido.SQL, "lhamalog:stream:logs") {
		t.Fatal("a consulta não usou o stream da referência")
	}
	if strings.Contains(pedido.SQL, "lhamalog_record_id") {
		t.Fatal("lote de vários filtrou recordId")
	}
	if _, confirma := any(svc).(interface {
		Confirmar(context.Context, string) error
	}); confirma {
		t.Fatal("recuperação confirma a entrada")
	}
	outro := &buscaFake{paginas: [][]domain.DocumentoFisico{{doc}}}
	concluido, err := (&ServicoRecuperacao{Busca: outro, Acesso: acesso}).Executar(context.Background(), ref)
	if err != nil || !concluido {
		t.Fatalf("efeito=%v err=%v", concluido, err)
	}
}

func TestLoteUnitarioFiltraRecordID(t *testing.T) {
	ref := referenciaTeste()
	id := "rec-1"
	ref.RecordID = &id
	ref.RecordCount = 1
	doc := documentoTeste("rec-1", "um")
	doc["lhamalog_record_id"] = "rec-1"
	busca := &buscaFake{paginas: [][]domain.DocumentoFisico{{doc}}}
	svc := &ServicoRecuperacao{Busca: busca, Acesso: acessoTeste()}
	resultado, err := svc.Recuperar(context.Background(), ref)
	if err != nil || resultado.Classificacao != domain.ClasseLoteCompleto {
		t.Fatalf("%s err=%v", resultado.Classificacao, err)
	}
	if !strings.Contains(busca.chamadas[0].SQL, `"lhamalog_record_id" = 'rec-1'`) {
		t.Fatal("lote unitário não filtrou recordId")
	}
}

func TestConfiguracaoAusenteNaoConsulta(t *testing.T) {
	base := acessoTeste()
	casos := []struct {
		nome string
		alt  func(*domain.ConfiguracaoAcesso)
	}{
		{nome: "url", alt: func(c *domain.ConfiguracaoAcesso) { c.URL = "" }},
		{nome: "usuario", alt: func(c *domain.ConfiguracaoAcesso) { c.Usuario = "" }},
		{nome: "senha", alt: func(c *domain.ConfiguracaoAcesso) { c.Senha = "" }},
		{nome: "timeout", alt: func(c *domain.ConfiguracaoAcesso) { c.Timeout = 0 }},
		{nome: "pagina", alt: func(c *domain.ConfiguracaoAcesso) { c.TamanhoPagina = 0 }},
		{nome: "lista", alt: func(c *domain.ConfiguracaoAcesso) { c.CamposCopia = nil }},
		{nome: "ordem", alt: func(c *domain.ConfiguracaoAcesso) { c.Ordem = "" }},
	}
	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			acesso := base
			caso.alt(&acesso)
			busca := &buscaFake{}
			resultado, err := (&ServicoRecuperacao{Busca: busca, Acesso: acesso}).Recuperar(context.Background(), referenciaTeste())
			if err != nil {
				t.Fatal(err)
			}
			if resultado.Classificacao != domain.ClasseConfigAusente || resultado.EfeitoConcluido || resultado.ContagemCalculada || resultado.Registros != nil {
				t.Fatalf("classe=%s efeito=%v", resultado.Classificacao, resultado.EfeitoConcluido)
			}
			if len(busca.chamadas) != 0 {
				t.Fatal("consulta foi enviada sem configuração")
			}
			if _, confirma := any(&ServicoRecuperacao{}).(interface {
				Confirmar(context.Context, string) error
			}); confirma {
				t.Fatal("recuperação confirma a entrada")
			}
		})
	}
}

func TestListaVaziaTornaDiferencaDivergente(t *testing.T) {
	vazia := []string{}
	acesso := acessoTeste()
	acesso.CamposCopia = &vazia
	a := documentoTeste("a", "um")
	b := documentoTeste("a", "dois")
	busca := &buscaFake{paginas: [][]domain.DocumentoFisico{{a, b}}}
	resultado, err := (&ServicoRecuperacao{Busca: busca, Acesso: acesso}).Recuperar(context.Background(), referenciaComContagem(1))
	if err != nil {
		t.Fatal(err)
	}
	if len(busca.chamadas) == 0 || resultado.Classificacao != domain.ClasseDivergente || !resultado.EfeitoConcluido || resultado.Registros != nil {
		t.Fatalf("classe=%s chamadas=%d", resultado.Classificacao, len(busca.chamadas))
	}
}

func TestNomesAusentesNaoConsultamEIsolam(t *testing.T) {
	acesso := acessoTeste()
	acesso.CampoProduto = ""
	busca := &buscaFake{paginas: [][]domain.DocumentoFisico{{documentoTeste("a", "um")}}}
	resultado, err := (&ServicoRecuperacao{Busca: busca, Acesso: acesso}).Recuperar(context.Background(), referenciaTeste())
	if err != nil {
		t.Fatal(err)
	}
	if len(busca.chamadas) != 0 || resultado.Classificacao != domain.ClasseIsolamento || !resultado.EfeitoConcluido || resultado.Registros != nil {
		t.Fatalf("classe=%s chamadas=%d efeito=%v", resultado.Classificacao, len(busca.chamadas), resultado.EfeitoConcluido)
	}
	if strings.Contains(resultadoSemSegredo(resultado), acesso.Senha) {
		t.Fatal("resultado expõe a senha")
	}
}

func TestPaginacaoDedupEContagem(t *testing.T) {
	t.Run("varias paginas e multiplo exato", func(t *testing.T) {
		acesso := acessoTeste()
		acesso.TamanhoPagina = 2
		ref := referenciaComContagem(4)
		busca := &buscaFake{paginas: [][]domain.DocumentoFisico{
			{documentoTeste("a", "a"), documentoTeste("b", "b")},
			{documentoTeste("c", "c"), documentoTeste("d", "d")},
			{},
		}}
		resultado, err := (&ServicoRecuperacao{Busca: busca, Acesso: acesso}).Recuperar(context.Background(), ref)
		if err != nil || resultado.Classificacao != domain.ClasseLoteCompleto || resultado.Encontrado != 4 {
			t.Fatalf("%s encontrado=%d err=%v", resultado.Classificacao, resultado.Encontrado, err)
		}
		if len(busca.chamadas) != 3 || busca.chamadas[0].From != 0 || busca.chamadas[1].From != 2 || busca.chamadas[2].From != 4 {
			t.Fatalf("páginas=%v", froms(busca))
		}
		for _, pedido := range busca.chamadas {
			if pedido.Size != 2 {
				t.Fatal("size foi trocado por recordCount")
			}
		}
	})

	t.Run("copias e identificador depois de recordCount", func(t *testing.T) {
		acesso := acessoTeste()
		acesso.TamanhoPagina = 2
		ref := referenciaComContagem(2)
		a1 := documentoTeste("a", "mesmo")
		a1["copia"] = "1"
		a2 := documentoTeste("a", "mesmo")
		a2["copia"] = "2"
		b := documentoTeste("b", "novo")
		busca := &buscaFake{paginas: [][]domain.DocumentoFisico{{a1, a2}, {b}}}
		resultado, err := (&ServicoRecuperacao{Busca: busca, Acesso: acesso}).Recuperar(context.Background(), ref)
		if err != nil || resultado.Classificacao != domain.ClasseLoteCompleto || resultado.Encontrado != 2 || len(resultado.Registros) != 2 {
			t.Fatalf("%s n=%d err=%v", resultado.Classificacao, len(resultado.Registros), err)
		}
		ids := map[string]int{}
		for _, doc := range resultado.Registros {
			id, _ := domain.IdentificadorRegistro(doc)
			ids[id]++
		}
		if ids["a"] != 1 || ids["b"] != 1 {
			t.Fatalf("ids=%v", ids)
		}
		if len(busca.chamadas) < 2 {
			t.Fatal("a leitura parou em recordCount")
		}
	})

	t.Run("sem teto de paginas", func(t *testing.T) {
		acesso := acessoTeste()
		acesso.TamanhoPagina = 1
		ref := referenciaComContagem(4)
		busca := &buscaFake{paginas: [][]domain.DocumentoFisico{
			{documentoTeste("a", "a")},
			{documentoTeste("b", "b")},
			{documentoTeste("c", "c")},
			{documentoTeste("d", "d")},
			{},
		}}
		resultado, err := (&ServicoRecuperacao{Busca: busca, Acesso: acesso}).Recuperar(context.Background(), ref)
		if err != nil || resultado.Classificacao != domain.ClasseLoteCompleto || len(busca.chamadas) != 5 {
			t.Fatalf("chamadas=%d classe=%s err=%v", len(busca.chamadas), resultado.Classificacao, err)
		}
	})
}

func TestCompletudeNaoAfirmaCausa(t *testing.T) {
	acesso := acessoTeste()
	acesso.TamanhoPagina = 5
	casos := []struct {
		nome       string
		esperado   int64
		docs       []domain.DocumentoFisico
		classe     domain.Classificacao
		efeito     bool
		encontrado int64
	}{
		{nome: "igual", esperado: 1, docs: []domain.DocumentoFisico{documentoTeste("a", "a")}, classe: domain.ClasseLoteCompleto, efeito: true, encontrado: 1},
		{nome: "menor", esperado: 2, docs: []domain.DocumentoFisico{documentoTeste("a", "a")}, classe: domain.ClasseIncompleto, efeito: false, encontrado: 1},
		{nome: "maior", esperado: 1, docs: []domain.DocumentoFisico{documentoTeste("a", "a"), documentoTeste("b", "b")}, classe: domain.ClasseExcesso, efeito: true, encontrado: 2},
		{nome: "zero", esperado: 1, docs: nil, classe: domain.ClasseIncompleto, efeito: false, encontrado: 0},
	}
	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			busca := &buscaFake{paginas: [][]domain.DocumentoFisico{caso.docs}}
			resultado, err := (&ServicoRecuperacao{Busca: busca, Acesso: acesso}).Recuperar(context.Background(), referenciaComContagem(caso.esperado))
			if err != nil {
				t.Fatal(err)
			}
			if resultado.Classificacao != caso.classe || resultado.EfeitoConcluido != caso.efeito || resultado.Encontrado != caso.encontrado || resultado.RecordCount != caso.esperado {
				t.Fatalf("classe=%s efeito=%v encontrado=%d", resultado.Classificacao, resultado.EfeitoConcluido, resultado.Encontrado)
			}
			if caso.classe != domain.ClasseLoteCompleto && resultado.Registros != nil {
				t.Fatal("contagem divergente entregou registros")
			}
			if strings.Contains(string(resultado.Classificacao), "atraso") {
				t.Fatal("causa afirmada")
			}
		})
	}

	t.Run("lote misto usa só os aceitos", func(t *testing.T) {
		ref := referenciaComContagem(1)
		busca := &buscaFake{paginas: [][]domain.DocumentoFisico{{documentoTeste("aceito", "ok")}}}
		resultado, err := (&ServicoRecuperacao{Busca: busca, Acesso: acesso}).Recuperar(context.Background(), ref)
		if err != nil || resultado.Classificacao != domain.ClasseLoteCompleto || len(busca.chamadas) != 1 {
			t.Fatalf("%s chamadas=%d err=%v", resultado.Classificacao, len(busca.chamadas), err)
		}
	})
}

func TestIsolamentoDeDocumentoECredencial(t *testing.T) {
	acesso := acessoTeste()
	acesso.TamanhoPagina = 5
	outro := documentoTeste("b", "alheio")
	outro["campo_tenant"] = "tenant-b"
	busca := &buscaFake{paginas: [][]domain.DocumentoFisico{{documentoTeste("a", "ok"), outro}}}
	resultado, err := (&ServicoRecuperacao{Busca: busca, Acesso: acesso}).Recuperar(context.Background(), referenciaComContagem(1))
	if err != nil || resultado.Classificacao != domain.ClasseIsolamento || !resultado.EfeitoConcluido || resultado.Registros != nil {
		t.Fatalf("%s efeito=%v err=%v", resultado.Classificacao, resultado.EfeitoConcluido, err)
	}
	if !strings.Contains(busca.chamadas[0].SQL, `"campo_tenant" = 'tenant-a'`) || !strings.Contains(busca.chamadas[0].SQL, `"campo_produto" = 'produto-a'`) {
		t.Fatal("a consulta não filtrou tenant e produto")
	}
}

func TestFalhaClassificadaSemParcial(t *testing.T) {
	acesso := acessoTeste()
	acesso.TamanhoPagina = 1
	senha := acesso.Senha
	primeira := []domain.DocumentoFisico{documentoTeste("a", "parcial-nao-entregar")}

	t.Run("timeout na primeira", func(t *testing.T) {
		busca := &buscaFake{erros: []error{&out.ErroBusca{Tipo: out.FalhaTimeout}}}
		resultado := recuperar(t, acesso, referenciaComContagem(1), busca)
		assertarFalha(t, resultado, domain.ClasseIndisponibilidade, false, senha, "parcial-nao-entregar")
	})
	t.Run("credencial", func(t *testing.T) {
		busca := &buscaFake{erros: []error{&out.ErroBusca{Tipo: out.FalhaCredencial}}}
		resultado := recuperar(t, acesso, referenciaComContagem(1), busca)
		assertarFalha(t, resultado, domain.ClasseAutenticacao, true, senha, "")
		if resultado.Classificacao == domain.ClasseIndisponibilidade {
			t.Fatal("credencial foi tratada como falha temporária")
		}
	})
	t.Run("resposta invalida", func(t *testing.T) {
		busca := &buscaFake{erros: []error{&out.ErroBusca{Tipo: out.FalhaResposta}}}
		resultado := recuperar(t, acesso, referenciaComContagem(1), busca)
		assertarFalha(t, resultado, domain.ClasseRespostaInvalida, true, senha, "")
	})
	t.Run("falha na segunda pagina", func(t *testing.T) {
		busca := &buscaFake{
			paginas: [][]domain.DocumentoFisico{primeira, nil},
			erros:   []error{nil, &out.ErroBusca{Tipo: out.FalhaPagina}},
		}
		resultado := recuperar(t, acesso, referenciaComContagem(2), busca)
		assertarFalha(t, resultado, domain.ClasseFalhaPagina, false, senha, "parcial-nao-entregar")
		if len(busca.chamadas) != 2 {
			t.Fatalf("chamadas=%d", len(busca.chamadas))
		}
	})
	t.Run("timeout na segunda", func(t *testing.T) {
		busca := &buscaFake{
			paginas: [][]domain.DocumentoFisico{primeira, nil},
			erros:   []error{nil, &out.ErroBusca{Tipo: out.FalhaTimeout}},
		}
		resultado := recuperar(t, acesso, referenciaComContagem(2), busca)
		assertarFalha(t, resultado, domain.ClasseIndisponibilidade, false, senha, "parcial-nao-entregar")
	})
	t.Run("isolamento precede pagina", func(t *testing.T) {
		alheio := documentoTeste("a", "alheio")
		alheio["campo_tenant"] = "tenant-b"
		busca := &buscaFake{
			paginas: [][]domain.DocumentoFisico{{alheio}, nil},
			erros:   []error{nil, &out.ErroBusca{Tipo: out.FalhaIndisponivel}},
		}
		resultado := recuperar(t, acesso, referenciaComContagem(1), busca)
		assertarFalha(t, resultado, domain.ClasseIsolamento, true, senha, "alheio")
	})
	t.Run("autenticacao precede isolamento", func(t *testing.T) {
		alheio := documentoTeste("a", "alheio")
		alheio["campo_tenant"] = "tenant-b"
		busca := &buscaFake{
			paginas: [][]domain.DocumentoFisico{{alheio}, nil},
			erros:   []error{nil, &out.ErroBusca{Tipo: out.FalhaCredencial}},
		}
		resultado := recuperar(t, acesso, referenciaComContagem(1), busca)
		assertarFalha(t, resultado, domain.ClasseAutenticacao, true, senha, "alheio")
	})
	t.Run("valor que nao pode ser filtro", func(t *testing.T) {
		ref := referenciaComContagem(1)
		ref.Stream = `stream" OR "lhamalog_batch_id" = "outro`
		busca := &buscaFake{paginas: [][]domain.DocumentoFisico{{documentoTeste("a", "nao")}}}
		resultado := recuperar(t, acesso, ref, busca)
		assertarFalha(t, resultado, domain.ClasseRespostaInvalida, true, senha, "nao")
		if len(busca.chamadas) != 0 || resultado.BatchID != ref.BatchID || resultado.Stream != ref.Stream || resultado.Organizacao != ref.Organizacao {
			t.Fatal("a consulta trocou lote, organização, stream ou faixa")
		}
	})
	t.Run("documento sem identificador precede contagem", func(t *testing.T) {
		sem := documentoTeste("a", "conteudo-oculto")
		delete(sem, "lhamalog_record_id")
		busca := &buscaFake{paginas: [][]domain.DocumentoFisico{{documentoTeste("b", "b"), sem}}}
		acessoPagina := acesso
		acessoPagina.TamanhoPagina = 5
		resultado := recuperar(t, acessoPagina, referenciaComContagem(1), busca)
		assertarFalha(t, resultado, domain.ClasseSemIdentificador, true, senha, "conteudo-oculto")
		if !resultado.ContagemCalculada || resultado.Encontrado != 1 || resultado.RecordCount != 1 {
			t.Fatal("a contagem calculada não foi informada")
		}
	})
}

type buscaFake struct {
	paginas  [][]domain.DocumentoFisico
	erros    []error
	chamadas []out.PedidoPagina
}

func (b *buscaFake) Consultar(_ context.Context, pedido out.PedidoPagina) ([]domain.DocumentoFisico, error) {
	b.chamadas = append(b.chamadas, pedido)
	i := len(b.chamadas) - 1
	if i < len(b.erros) && b.erros[i] != nil {
		return nil, b.erros[i]
	}
	if i < len(b.paginas) {
		return b.paginas[i], nil
	}
	return nil, nil
}

func recuperar(t *testing.T, acesso domain.ConfiguracaoAcesso, ref domain.Referencia, busca out.Busca) domain.ResultadoRecuperacao {
	t.Helper()
	resultado, err := (&ServicoRecuperacao{Busca: busca, Acesso: acesso}).Recuperar(context.Background(), ref)
	if err != nil {
		t.Fatal(err)
	}
	return resultado
}

func assertarFalha(t *testing.T, resultado domain.ResultadoRecuperacao, classe domain.Classificacao, efeito bool, senha, conteudo string) {
	t.Helper()
	if resultado.Classificacao != classe || resultado.EfeitoConcluido != efeito || resultado.Registros != nil {
		t.Fatalf("classe=%s efeito=%v registros=%d", resultado.Classificacao, resultado.EfeitoConcluido, len(resultado.Registros))
	}
	texto := resultadoSemSegredo(resultado)
	if senha != "" && strings.Contains(texto, senha) {
		t.Fatal("resultado expõe a senha")
	}
	if conteudo != "" && strings.Contains(texto, conteudo) {
		t.Fatal("resultado expõe conteúdo de log")
	}
}

func resultadoSemSegredo(resultado domain.ResultadoRecuperacao) string {
	return strings.Join([]string{
		string(resultado.Classificacao),
		resultado.TenantID,
		resultado.ProductID,
		resultado.Organizacao,
		resultado.Stream,
		resultado.BatchID,
		resultado.IdempotencyID,
		resultado.CorrelationID,
		resultado.EntradaID,
	}, " ")
}

func froms(b *buscaFake) []int {
	saida := make([]int, len(b.chamadas))
	for i, pedido := range b.chamadas {
		saida[i] = pedido.From
	}
	return saida
}

func referenciaTeste() domain.Referencia {
	return referenciaComContagem(2)
}

func referenciaComContagem(n int64) domain.Referencia {
	return domain.Referencia{
		TenantID:           "tenant-a",
		ProductID:          "produto-a",
		Organizacao:        "org-a",
		Stream:             "stream_sintetico",
		BatchID:            "batch-a",
		MinTimestampMicros: 100,
		MaxTimestampMicros: 200,
		RecordCount:        n,
		IdempotencyID:      "batch-a",
		CorrelationID:      "corr-1",
		EntradaID:          "1-0",
	}
}

func acessoTeste() domain.ConfiguracaoAcesso {
	copia := []string{"copia"}
	return domain.ConfiguracaoAcesso{
		URL:           "http://oo.teste",
		Usuario:       "usuario",
		Senha:         "segredo-nao-logar",
		Timeout:       time.Second,
		TamanhoPagina: 2,
		CamposCopia:   &copia,
		Ordem:         "campo_ordem",
		CampoTenant:   "campo_tenant",
		CampoProduto:  "campo_produto",
	}
}

func documentoTeste(id, mensagem string) domain.DocumentoFisico {
	return domain.DocumentoFisico{
		"lhamalog_record_id": id,
		"campo_tenant":       "tenant-a",
		"campo_produto":      "produto-a",
		"mensagem":           mensagem,
	}
}
