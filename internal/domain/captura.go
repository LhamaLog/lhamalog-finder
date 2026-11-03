package domain

// Registro é o log armazenado, opaco. Os campos internos não estão definidos.
type Registro []byte

// Motivo distingue as três falhas de captura.
type Motivo string

const (
	MotivoInutilizavel Motivo = "inutilizavel"
	MotivoAusente      Motivo = "ausente"
	MotivoInacessivel  Motivo = "inacessivel"
)

// FalhaCaptura é o apontamento e um motivo. Não é log pronto.
type FalhaCaptura struct {
	Apontamento Codigo
	Motivo      Motivo
}

// LogPronto é o registro inteiro associado ao apontamento.
type LogPronto struct {
	Apontamento Codigo
	Conteudo    Registro
}

// QuadroCaptura guarda logs prontos de uma execução. O lugar não tem nome.
type QuadroCaptura struct {
	prontos map[Codigo]Registro
	falhas  []FalhaCaptura
}

// NovoQuadroCaptura começa sem log pronto e sem falha.
func NovoQuadroCaptura() *QuadroCaptura {
	return &QuadroCaptura{prontos: make(map[Codigo]Registro)}
}

// Contagens devolve quantos logs prontos e quantas falhas existem.
func (q *QuadroCaptura) Contagens() (prontos, falhas int) {
	if q == nil {
		return 0, 0
	}
	return len(q.prontos), len(q.falhas)
}

// Inutilizavel é verdadeiro quando o apontamento é vazio.
// Formato, tamanho e alfabeto de um valor não vazio não estão definidos.
func Inutilizavel(c Codigo) bool {
	return c == ""
}

// Unicos reduz indicações repetidas do mesmo valor a um apontamento.
// A ordem da primeira aparição permanece.
func Unicos(lista []Codigo) []Codigo {
	vistos := make(map[Codigo]struct{}, len(lista))
	saida := make([]Codigo, 0, len(lista))
	for _, c := range lista {
		if _, ok := vistos[c]; ok {
			continue
		}
		vistos[c] = struct{}{}
		saida = append(saida, c)
	}
	return saida
}

// Falhar registra o motivo e não cria log pronto.
func (q *QuadroCaptura) Falhar(c Codigo, motivo Motivo) FalhaCaptura {
	f := FalhaCaptura{Apontamento: c, Motivo: motivo}
	q.falhas = append(q.falhas, f)
	return f
}

// Concluir guarda o registro inteiro quando a instância está no ar.
// Uma segunda guarda do mesmo apontamento não cria outro log.
// Com a instância parada, não guarda.
func (q *QuadroCaptura) Concluir(c Codigo, conteudo Registro, instanciaNoAr bool) (LogPronto, bool) {
	if q == nil || !instanciaNoAr || Inutilizavel(c) {
		return LogPronto{}, false
	}
	if existente, ok := q.prontos[c]; ok {
		return LogPronto{Apontamento: c, Conteudo: copiarRegistro(existente)}, true
	}
	q.prontos[c] = copiarRegistro(conteudo)
	return LogPronto{Apontamento: c, Conteudo: copiarRegistro(q.prontos[c])}, true
}

// Ler devolve o registro inteiro enquanto a instância está no ar.
// Com a instância parada, a leitura não é exigida.
func (q *QuadroCaptura) Ler(c Codigo, instanciaNoAr bool) (Registro, bool) {
	if q == nil || !instanciaNoAr {
		return nil, false
	}
	r, ok := q.prontos[c]
	if !ok {
		return nil, false
	}
	return copiarRegistro(r), true
}

// ProntosDe conta logs prontos daquele apontamento. O máximo é um.
func (q *QuadroCaptura) ProntosDe(c Codigo) int {
	if q == nil {
		return 0
	}
	if _, ok := q.prontos[c]; ok {
		return 1
	}
	return 0
}

func copiarRegistro(r Registro) Registro {
	if r == nil {
		return nil
	}
	copia := make(Registro, len(r))
	copy(copia, r)
	return copia
}
