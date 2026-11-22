// Package catalogo carrega o artefato de segurança DF-FINDER-2026-11-22.
package catalogo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"github.com/LhamaLog/lhamalog-finder/internal/strictjson"
)

const VersaoDistribuida = "finder-2026-11-22.1"

// A mesma versão é imutável, inclusive entre as distribuições Finder/Sugestão.
const SHA256Distribuido = "e5ef721a8337390c12745b815e36f7dfb3f353a289cf41bc13d2d951b87dc2ad"

var ErrArquivo = errors.New("catalogo_seguranca_invalido")

type Arquivo struct{ Caminho string }
type especificacao struct {
	Schema   string   `json:"schemaVersion"`
	Versao   string   `json:"versao"`
	Email    string   `json:"placeholderEmail"`
	Marcador string   `json:"marcadorLinhaSegredo"`
	Padroes  []padrao `json:"padroes"`
}
type padrao struct {
	ID    string `json:"id"`
	Acao  string `json:"acao"`
	Regex string `json:"regex"`
}

func (a Arquivo) Carregar(ctx context.Context) (domain.CatalogoSeguranca, error) {
	if err := ctx.Err(); err != nil {
		return domain.CatalogoSeguranca{}, err
	}
	f, err := os.Open(a.Caminho)
	if err != nil {
		return domain.CatalogoSeguranca{}, ErrArquivo
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, 1024*1024+1))
	if err != nil || len(data) > 1024*1024 {
		return domain.CatalogoSeguranca{}, ErrArquivo
	}
	hash := sha256.Sum256(data)
	if hex.EncodeToString(hash[:]) != SHA256Distribuido {
		return domain.CatalogoSeguranca{}, ErrArquivo
	}
	return compilar(data)
}

func compilar(data []byte) (domain.CatalogoSeguranca, error) {
	falhar := func() (domain.CatalogoSeguranca, error) { return domain.CatalogoSeguranca{}, ErrArquivo }
	if !utf8.Valid(data) {
		return falhar()
	}
	var spec especificacao
	if strictjson.Decode(data, &spec, "schemaVersion", "versao", "placeholderEmail", "marcadorLinhaSegredo", "padroes", "id", "acao", "regex") != nil {
		return falhar()
	}
	if spec.Schema != "1" || spec.Versao != VersaoDistribuida || spec.Email != "[email omitido]" || spec.Marcador != domain.MarcadorOmisso || len(spec.Padroes) == 0 {
		return falhar()
	}
	var emails, segredos []*regexp.Regexp
	ids := map[string]bool{}
	for _, p := range spec.Padroes {
		if p.ID == "" || ids[p.ID] || p.Regex == "" {
			return falhar()
		}
		ids[p.ID] = true
		re, err := regexp.Compile(p.Regex)
		if err != nil || re.MatchString("") {
			return falhar()
		}
		switch p.Acao {
		case "SUBSTITUIR_EMAIL":
			emails = append(emails, re)
		case "OMITIR_LINHA":
			segredos = append(segredos, re)
		default:
			return falhar()
		}
	}
	if len(emails) == 0 || len(segredos) == 0 {
		return falhar()
	}
	return domain.CatalogoSeguranca{Versao: spec.Versao, Transformar: func(s string) (domain.LinhaLogMascarada, int, error) {
		if !utf8.ValidString(s) {
			return domain.LinhaLogMascarada{}, 0, ErrArquivo
		}
		matches := intervalos(emails, s)
		var out strings.Builder
		pos := 0
		for _, m := range matches {
			out.WriteString(s[pos:m[0]])
			out.WriteString(spec.Email)
			pos = m[1]
		}
		out.WriteString(s[pos:])
		texto := out.String()
		secrets := intervalos(segredos, texto)
		n := len(matches) + len(secrets)
		if len(secrets) > 0 {
			return domain.LinhaLogMascarada{Texto: spec.Marcador, LinhaRemovida: true}, n, nil
		}
		return domain.LinhaLogMascarada{Texto: texto}, n, nil
	}}, nil
}

// Une ocorrências sobrepostas de vários padrões, sem contar ou substituir duas vezes.
func intervalos(padroes []*regexp.Regexp, s string) [][2]int {
	var encontrados [][2]int
	for _, p := range padroes {
		for _, m := range p.FindAllStringIndex(s, -1) {
			encontrados = append(encontrados, [2]int{m[0], m[1]})
		}
	}
	sort.Slice(encontrados, func(i, j int) bool { return encontrados[i][0] < encontrados[j][0] })
	var unicos [][2]int
	for _, m := range encontrados {
		n := len(unicos)
		if n > 0 && m[0] < unicos[n-1][1] {
			if m[1] > unicos[n-1][1] {
				unicos[n-1][1] = m[1]
			}
		} else {
			unicos = append(unicos, m)
		}
	}
	return unicos
}
