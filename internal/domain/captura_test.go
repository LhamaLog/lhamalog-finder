package domain

import (
	"bytes"
	"testing"
)

func TestQuadroSemApontamento(t *testing.T) {
	q := NovoQuadroCaptura()
	prontos, falhas := q.Contagens()
	if prontos != 0 || falhas != 0 {
		t.Fatalf("prontos=%d falhas=%d", prontos, falhas)
	}
}

func TestConcluirGuardaRegistroInteiro(t *testing.T) {
	q := NovoQuadroCaptura()
	conteudo := Registro(`{"senha":"segredo","msg":"inteiro"}`)
	if Inutilizavel("apontamento-a") {
		t.Fatal("apontamento utilizável foi recusado")
	}
	pronto, ok := q.Concluir("apontamento-a", conteudo, true)
	if !ok || !bytes.Equal(pronto.Conteudo, conteudo) {
		t.Fatalf("concluir: %+v", pronto)
	}
	lido, ok := q.Ler("apontamento-a", true)
	if !ok || !bytes.Equal(lido, conteudo) {
		t.Fatalf("leitura: %q ok=%v", lido, ok)
	}
	conteudo[0] = 'X'
	lido2, _ := q.Ler("apontamento-a", true)
	if bytes.Equal(lido2, conteudo) {
		t.Fatal("o registro guardado mudou com o original")
	}
}

func TestParadaNoMeioNaoDeixaPronto(t *testing.T) {
	q := NovoQuadroCaptura()
	if _, ok := q.Concluir("apontamento-a", Registro("log"), false); ok {
		t.Fatal("parada no meio guardou log")
	}
	if _, ok := q.Ler("apontamento-a", true); ok {
		t.Fatal("log pronto depois da parada no meio")
	}
}

func TestDepoisQueInstanciaParaNaoExigeLog(t *testing.T) {
	q := NovoQuadroCaptura()
	if _, ok := q.Concluir("apontamento-a", Registro("log"), true); !ok {
		t.Fatal("concluir")
	}
	if _, ok := q.Ler("apontamento-a", false); ok {
		t.Fatal("a leitura exigiu o log depois que a instância parou")
	}
	prontos, falhas := q.Contagens()
	if falhas != 0 || prontos != 1 {
		t.Fatalf("prontos=%d falhas=%d", prontos, falhas)
	}
}

func TestApontamentoVazioEInutilizavel(t *testing.T) {
	q := NovoQuadroCaptura()
	if !Inutilizavel("") {
		t.Fatal("vazio foi aceito")
	}
	f := q.Falhar("", MotivoInutilizavel)
	if f.Motivo != MotivoInutilizavel || f.Apontamento != "" {
		t.Fatalf("falha: %+v", f)
	}
	if _, ok := q.Concluir("", Registro("log"), true); ok || q.ProntosDe("") != 0 {
		t.Fatal("apontamento vazio ficou pronto")
	}
}

func TestAusenteNaoCriaRegistroVazio(t *testing.T) {
	q := NovoQuadroCaptura()
	f := q.Falhar("apontamento-b", MotivoAusente)
	if f.Motivo != MotivoAusente || q.ProntosDe("apontamento-b") != 0 {
		t.Fatalf("falha=%+v prontos=%d", f, q.ProntosDe("apontamento-b"))
	}
}

func TestInacessivelDistintoDeAusente(t *testing.T) {
	q := NovoQuadroCaptura()
	ausente := q.Falhar("apontamento-b", MotivoAusente)
	inacessivel := q.Falhar("apontamento-c", MotivoInacessivel)
	if ausente.Motivo == inacessivel.Motivo {
		t.Fatal("motivos iguais")
	}
	if q.ProntosDe("apontamento-b") != 0 || q.ProntosDe("apontamento-c") != 0 {
		t.Fatal("falha produziu log pronto")
	}
}

func TestMesmoApontamentoUmLogPronto(t *testing.T) {
	q := NovoQuadroCaptura()
	primeiro := Registro("primeiro")
	if _, ok := q.Concluir("mesmo", primeiro, true); !ok {
		t.Fatal("primeiro")
	}
	segundo, ok := q.Concluir("mesmo", Registro("segundo"), true)
	if !ok || q.ProntosDe("mesmo") != 1 {
		t.Fatalf("prontos=%d ok=%v", q.ProntosDe("mesmo"), ok)
	}
	if !bytes.Equal(segundo.Conteudo, primeiro) {
		t.Fatalf("segunda guarda trocou o registro: %q", segundo.Conteudo)
	}
	if len(Unicos([]Codigo{"mesmo", "mesmo", "outro"})) != 2 {
		t.Fatal("indicações repetidas não convergiram")
	}
}
