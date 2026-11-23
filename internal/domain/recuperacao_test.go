package domain

import (
	"math"
	"strings"
	"testing"
	"time"
)

func TestDeduplicarCopiaEDivergencia(t *testing.T) {
	a1 := DocumentoFisico{"lhamalog_record_id": "a", "mensagem": "um", "copia": "1"}
	a2 := DocumentoFisico{"lhamalog_record_id": "a", "mensagem": "um", "copia": "2"}
	registros, divergente, semID := Deduplicar([]DocumentoFisico{a1, a2}, []string{"copia"})
	if divergente || semID || len(registros) != 1 {
		t.Fatalf("cópia: n=%d divergente=%v semID=%v", len(registros), divergente, semID)
	}
	if registros[0]["copia"] != "1" && registros[0]["copia"] != "2" {
		t.Fatal("a cópia entregue não é um dos documentos")
	}
	if registros[0]["mensagem"] != "um" {
		t.Fatal("o registro entregue ocultou campo")
	}

	outro := DocumentoFisico{"lhamalog_record_id": "a", "mensagem": "dois", "copia": "3"}
	registros, divergente, semID = Deduplicar([]DocumentoFisico{a1, outro}, []string{"copia"})
	if !divergente || semID || len(registros) != 0 {
		t.Fatalf("divergência: n=%d divergente=%v semID=%v", len(registros), divergente, semID)
	}

	registros, divergente, _ = Deduplicar([]DocumentoFisico{a1, a2}, []string{})
	if !divergente || len(registros) != 0 {
		t.Fatal("lista vazia aceitou diferença")
	}

	iguais := []DocumentoFisico{
		{"lhamalog_record_id": "a", "mensagem": "um"},
		{"lhamalog_record_id": "a", "mensagem": "um"},
	}
	registros, divergente, semID = Deduplicar(iguais, nil)
	if divergente || semID || len(registros) != 1 {
		t.Fatal("documentos iguais não viraram um registro")
	}
}

func TestDeduplicarIgnoraDocumentoSemIdentificadorENaoUsaCorrelacao(t *testing.T) {
	docs := []DocumentoFisico{
		{"mensagem": "sem-id", "correlationId": "c", "entrada": "1-0"},
		{"lhamalog_record_id": "", "mensagem": "vazio"},
		{"lhamalog_record_id": "a", "correlationId": "c", "entrada": "1-0", "mensagem": "a"},
		{"lhamalog_record_id": "b", "correlationId": "c", "entrada": "1-0", "mensagem": "b"},
	}
	registros, divergente, semID := Deduplicar(docs, nil)
	if !semID || divergente || len(registros) != 2 {
		t.Fatalf("n=%d divergente=%v semID=%v", len(registros), divergente, semID)
	}
	if registros[0]["lhamalog_record_id"] == registros[1]["lhamalog_record_id"] {
		t.Fatal("correlationId deduplicou registros diferentes")
	}
}

func TestFecharCompletudeEPrecedencia(t *testing.T) {
	ref := Referencia{RecordCount: 2, TenantID: "t", ProductID: "p", BatchID: "b", IdempotencyID: "b", CorrelationID: "c"}
	casos := []struct {
		nome    string
		aval    Avaliacao
		classe  Classificacao
		efeito  bool
		entrega bool
	}{
		{nome: "completo", aval: Avaliacao{Encontrado: 2, ContagemCalculada: true}, classe: ClasseLoteCompleto, efeito: true, entrega: true},
		{nome: "incompleto", aval: Avaliacao{Encontrado: 1, ContagemCalculada: true}, classe: ClasseIncompleto, efeito: false},
		{nome: "excesso", aval: Avaliacao{Encontrado: 3, ContagemCalculada: true}, classe: ClasseExcesso, efeito: true},
		{nome: "vazio", aval: Avaliacao{Encontrado: 0, ContagemCalculada: true}, classe: ClasseIncompleto, efeito: false},
		{nome: "sem identificador e contagem", aval: Avaliacao{Encontrado: 2, ContagemCalculada: true, SemIdentificador: true}, classe: ClasseSemIdentificador, efeito: true},
		{nome: "divergente e excesso", aval: Avaliacao{Encontrado: 3, ContagemCalculada: true, Divergente: true}, classe: ClasseDivergente, efeito: true},
		{nome: "isolamento e sem identificador", aval: Avaliacao{Encontrado: 0, ContagemCalculada: true, Isolamento: true, SemIdentificador: true}, classe: ClasseIsolamento, efeito: true},
		{nome: "credencial", aval: Avaliacao{Acesso: AcessoCredencial, Isolamento: true}, classe: ClasseAutenticacao, efeito: true},
		{nome: "resposta", aval: Avaliacao{Acesso: AcessoResposta, SemIdentificador: true}, classe: ClasseRespostaInvalida, efeito: true},
		{nome: "timeout", aval: Avaliacao{Acesso: AcessoTimeout}, classe: ClasseIndisponibilidade, efeito: false},
		{nome: "indisponivel", aval: Avaliacao{Acesso: AcessoIndisponivel, PaginasLidas: 1}, classe: ClasseIndisponibilidade, efeito: false},
		{nome: "pagina seguinte", aval: Avaliacao{Acesso: AcessoPagina, PaginasLidas: 1}, classe: ClasseFalhaPagina, efeito: false},
		{nome: "pagina seguinte depois de isolamento", aval: Avaliacao{Acesso: AcessoPagina, PaginasLidas: 1, Isolamento: true}, classe: ClasseIsolamento, efeito: true},
		{nome: "primeira pagina http", aval: Avaliacao{Acesso: AcessoPagina}, classe: ClasseIndisponibilidade, efeito: false},
	}
	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			doc := DocumentoFisico{"mensagem": "nao-entregar", "lhamalog_record_id": "a"}
			resultado := Fechar(ref, caso.aval, []DocumentoFisico{doc})
			if resultado.Classificacao != caso.classe || resultado.EfeitoConcluido != caso.efeito {
				t.Fatalf("classe=%s efeito=%v", resultado.Classificacao, resultado.EfeitoConcluido)
			}
			if caso.entrega && len(resultado.Registros) != 1 {
				t.Fatal("lote completo não entregou o registro")
			}
			if !caso.entrega && resultado.Registros != nil {
				t.Fatal("falha entregou conteúdo de log")
			}
			if caso.aval.ContagemCalculada && (resultado.RecordCount != 2 || resultado.Encontrado != caso.aval.Encontrado || !resultado.ContagemCalculada) {
				t.Fatal("quantidades não foram informadas")
			}
			if !caso.aval.ContagemCalculada && resultado.ContagemCalculada {
				t.Fatal("contagem foi marcada sem a leitura ter fechado")
			}
			if resultado.IdempotencyID != "b" || resultado.CorrelationID != "c" {
				t.Fatal("contexto do lote não foi preservado")
			}
			if strings.Contains(string(resultado.Classificacao), "atraso") || strings.Contains(string(resultado.Classificacao), "persistência") {
				t.Fatal("a contagem afirmou a causa da divergência")
			}
		})
	}
}

func TestConsultaSQLMontaFiltroERecusaValorInaplicavel(t *testing.T) {
	cfg := ConfiguracaoAcesso{Ordem: "campo_ordem", CampoTenant: "campo_tenant", CampoProduto: "campo_produto", Timeout: time.Second}
	ref := Referencia{
		TenantID: "tenant-a", ProductID: "produto-a", Organizacao: "org-a", Stream: "stream_sintetico",
		BatchID: "batch-a", MinTimestampMicros: 100, MaxTimestampMicros: 200, RecordCount: 9,
	}
	sql, inicio, fim, err := ConsultaSQL(ref, cfg)
	if err != nil {
		t.Fatal(err)
	}
	esperado := `SELECT * FROM "stream_sintetico" WHERE "lhamalog_batch_id" = 'batch-a' AND "campo_tenant" = 'tenant-a' AND "campo_produto" = 'produto-a' ORDER BY campo_ordem`
	if sql != esperado || inicio != 100 || fim != 201 {
		t.Fatalf("sql=%s inicio=%d fim=%d", sql, inicio, fim)
	}
	if strings.Contains(sql, "lhamalog:stream:logs") || strings.Contains(sql, "_timestamp") || strings.Contains(sql, "lhamalog_record_id") {
		t.Fatal("a consulta usou stream Redis, ordem ou recordId inventados")
	}

	id := "rec-1"
	ref.RecordID = &id
	sql, _, _, err = ConsultaSQL(ref, cfg)
	if err != nil || !strings.Contains(sql, `"lhamalog_record_id" = 'rec-1'`) {
		t.Fatalf("lote unitário: %s err=%v", sql, err)
	}

	ref.BatchID = "42"
	ref.RecordID = nil
	sql, _, _, err = ConsultaSQL(ref, cfg)
	if err != nil || !strings.Contains(sql, `"lhamalog_batch_id" = 42`) || strings.Contains(sql, `"lhamalog_batch_id" = '42'`) {
		t.Fatalf("inteiro: %s err=%v", sql, err)
	}

	ref.BatchID = "a'b"
	sql, _, _, err = ConsultaSQL(ref, cfg)
	if err != nil || !strings.Contains(sql, `"lhamalog_batch_id" = 'a''b'`) {
		t.Fatalf("aspas: %s err=%v", sql, err)
	}

	ref.BatchID = "batch-a"
	ref.Stream = `logs" OR "x" = "y`
	if _, _, _, err = ConsultaSQL(ref, cfg); err == nil {
		t.Fatal("stream com aspas foi aceito")
	}
	ref.Stream = streamRedis
	if _, _, _, err = ConsultaSQL(ref, cfg); err == nil {
		t.Fatal("stream Redis foi aceito")
	}
	ref.Stream = "stream_sintetico"
	ref.Organizacao = "org/outra"
	if _, _, _, err = ConsultaSQL(ref, cfg); err == nil {
		t.Fatal("organização que troca o caminho foi aceita")
	}
	ref.Organizacao = "org-a"
	ref.MaxTimestampMicros = math.MaxInt64
	if _, _, _, err = ConsultaSQL(ref, cfg); err == nil {
		t.Fatal("faixa que estoura foi aceita")
	}
}

func TestIsoladoNaoUsaAtributoInterno(t *testing.T) {
	interno := DocumentoFisico{
		"lhamalog_record_id": "a",
		"log":                map[string]any{"tenantId": "tenant-a", "productId": "produto-a"},
	}
	if !Isolado([]DocumentoFisico{interno}, "campo_tenant", "campo_produto", "tenant-a", "produto-a") {
		t.Fatal("atributo interno substituiu o campo configurado")
	}
	valido := DocumentoFisico{"lhamalog_record_id": "a", "campo_tenant": "tenant-a", "campo_produto": "produto-a"}
	outro := DocumentoFisico{"lhamalog_record_id": "b", "campo_tenant": "tenant-b", "campo_produto": "produto-a"}
	if !Isolado([]DocumentoFisico{valido, outro}, "campo_tenant", "campo_produto", "tenant-a", "produto-a") {
		t.Fatal("documento de outro tenant passou")
	}
	if Isolado([]DocumentoFisico{valido}, "campo_tenant", "campo_produto", "tenant-a", "produto-a") {
		t.Fatal("documento do tenant foi recusado")
	}
	if !Isolado(nil, "", "campo_produto", "tenant-a", "produto-a") {
		t.Fatal("nome ausente não é isolamento")
	}
}
