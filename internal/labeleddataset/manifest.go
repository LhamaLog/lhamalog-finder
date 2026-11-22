package labeleddataset

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"github.com/LhamaLog/lhamalog-finder/internal/strictjson"
)

var ErrManifesto = errors.New("manifesto_base_invalido")

type ArquivoBase struct {
	Nome      string `json:"nome"`
	SHA256    string `json:"sha256"`
	Registros int    `json:"registros"`
}
type CapturaBase struct {
	Arquivo        string `json:"arquivo"`
	SHA256         string `json:"sha256"`
	CodigoSHA256LF string `json:"codigoSha256LF"`
	Eventos        int    `json:"eventos"`
	PorClasse      int    `json:"porClasse"`
	Inicio         string `json:"inicio"`
	Fim            string `json:"fim"`
	Go             string `json:"goVersion"`
	Redis          string `json:"redisVersion"`
}
type ParticaoBase struct {
	Nome     string `json:"nome"`
	Estado   string `json:"estado"`
	Criterio string `json:"criterio"`
}
type Manifesto struct {
	Schema      string         `json:"schemaVersion"`
	Versao      string         `json:"versao"`
	OrigemDrive string         `json:"origemDrive"`
	TipoOrigem  string         `json:"tipoOrigem"`
	Autoria     string         `json:"autoria"`
	Metodo      string         `json:"metodoRotulagem"`
	UsoMVP      string         `json:"usoMVP"`
	Minimo      int            `json:"minimoPorClasse"`
	Registros   int            `json:"registros"`
	Contagens   map[string]int `json:"contagens"`
	Arquivos    []ArquivoBase  `json:"arquivos"`
	Captura     CapturaBase    `json:"captura"`
	Particoes   []ParticaoBase `json:"particoes"`
}

// CarregarBasePublicada valida proveniência, hashes e cobertura de cada evento.
// Não usa o classificador para construir ou verificar o rótulo de origem.
func CarregarBasePublicada(dir string) (Manifesto, []Registro, error) {
	fail := func() (Manifesto, []Registro, error) { return Manifesto{}, nil, ErrManifesto }
	raw, e := os.ReadFile(filepath.Join(dir, "manifest.json"))
	if e != nil || len(raw) > 1024*1024 || !utf8.Valid(raw) {
		return fail()
	}
	var m Manifesto
	if strictjson.Decode(raw, &m, "schemaVersion", "versao", "origemDrive", "tipoOrigem", "autoria", "metodoRotulagem", "usoMVP", "minimoPorClasse", "registros", "contagens", "ERROR", "WARNING", "INFO", "arquivos", "nome", "sha256", "captura", "arquivo", "codigoSha256LF", "eventos", "porClasse", "inicio", "fim", "goVersion", "redisVersion", "particoes", "estado", "criterio") != nil {
		return fail()
	}
	if m.Schema != "1" || m.Versao != "finder-2026-11-22.1" || !strings.HasPrefix(m.OrigemDrive, "https://docs.google.com/document/d/") || m.TipoOrigem != "captura-controlada-redis" || m.Autoria != "engenharia-finder-codex" || m.Metodo != "severidade-emitida-na-origem" || m.UsoMVP != "avaliacao" || m.Registros < 3000 || m.Minimo < 1000 || len(m.Contagens) != 3 {
		return fail()
	}
	if len(m.Particoes) != 3 || m.Particoes[0] != (ParticaoBase{"validacao", "ativa", "sessao-controlada-unica"}) || m.Particoes[1] != (ParticaoBase{"treino", "reservada", "novas-origens-disjuntas"}) || m.Particoes[2] != (ParticaoBase{"teste", "reservada", "novas-sessoes-independentes"}) {
		return fail()
	}
	nomes, _ := filepath.Glob(filepath.Join(dir, "*.jsonl"))
	if len(nomes) != len(m.Arquivos) || len(nomes) == 0 {
		return fail()
	}
	seen := map[string]bool{}
	for _, a := range m.Arquivos {
		if a.Nome != filepath.Base(a.Nome) || !strings.HasSuffix(a.Nome, ".jsonl") || seen[a.Nome] || a.Registros <= 0 {
			return fail()
		}
		seen[a.Nome] = true
		data, e := os.ReadFile(filepath.Join(dir, a.Nome))
		if e != nil || digest(data) != a.SHA256 {
			return fail()
		}
		rs, e := decodificarRegistros(bytes.NewReader(data))
		if e != nil || len(rs) != a.Registros {
			return fail()
		}
	}
	registros, e := LoadDiretorio(dir)
	if e != nil || len(registros) != m.Registros {
		return fail()
	}
	counts := map[string]int{}
	for _, r := range registros {
		if r.Split != "validacao" {
			return fail()
		}
		counts[string(r.ClassLabel)]++
	}
	for _, c := range []string{"ERROR", "WARNING", "INFO"} {
		if counts[c] != m.Contagens[c] || counts[c] != m.Minimo {
			return fail()
		}
	}
	if err := validarCaptura(dir, m, registros); err != nil {
		return fail()
	}
	return m, registros, nil
}

type eventoCapturado struct {
	RecordID string          `json:"recordId"`
	Classe   string          `json:"classeOrigem"`
	Log      json.RawMessage `json:"evento"`
}
type origemEvento struct {
	classe domain.ClasseLog
	nivel  string
	log    map[string]any
}

func validarCaptura(dir string, m Manifesto, registros []Registro) error {
	c := m.Captura
	inicio, e := time.Parse(time.RFC3339Nano, c.Inicio)
	if e != nil {
		return ErrManifesto
	}
	fim, e := time.Parse(time.RFC3339Nano, c.Fim)
	if e != nil || fim.Before(inicio) || c.Arquivo != "proveniencia/captura.jsonl.gz" || c.Eventos <= 0 || c.PorClasse <= 0 || c.Go == "" || c.Redis == "" {
		return ErrManifesto
	}
	source, e := os.ReadFile(filepath.Join(dir, c.Arquivo))
	if e != nil || digest(source) != c.SHA256 {
		return ErrManifesto
	}
	code, e := os.ReadFile(filepath.Join(dir, "..", "..", "cmd", "capturar-base", "main.go"))
	if e != nil || digest(bytes.ReplaceAll(code, []byte("\r\n"), []byte("\n"))) != c.CodigoSHA256LF {
		return ErrManifesto
	}
	gz, e := gzip.NewReader(bytes.NewReader(source))
	if e != nil {
		return ErrManifesto
	}
	defer gz.Close()
	data, e := io.ReadAll(io.LimitReader(gz, 64*1024*1024+1))
	if e != nil || len(data) > 64*1024*1024 {
		return ErrManifesto
	}
	origens := map[string]origemEvento{}
	classes := map[string]int{}
	br := bufio.NewScanner(bytes.NewReader(data))
	br.Buffer(make([]byte, 4096), 1024*1024)
	for br.Scan() {
		var row eventoCapturado
		if strictjson.Decode(br.Bytes(), &row, "recordId", "classeOrigem", "evento", "time", "level", "msg", "operation", "sequence", "elapsed_us", "diagnostic") != nil || row.RecordID == "" {
			return ErrManifesto
		}
		if _, ok := origens[row.RecordID]; ok {
			return ErrManifesto
		}
		var log map[string]any
		if json.Unmarshal(row.Log, &log) != nil {
			return ErrManifesto
		}
		level, _ := log["level"].(string)
		timestamp, _ := log["time"].(string)
		instante, err := time.Parse(time.RFC3339Nano, timestamp)
		if err != nil || instante.Before(inicio) || instante.After(fim) {
			return ErrManifesto
		}
		var classe domain.ClasseLog
		switch level {
		case "DEBUG", "INFO":
			classe = domain.INFO
		case "WARN":
			classe = domain.WARNING
		case "ERROR":
			classe = domain.ERROR
		default:
			return ErrManifesto
		}
		if string(classe) != row.Classe {
			return ErrManifesto
		}
		origens[row.RecordID] = origemEvento{classe, level, log}
		classes[string(classe)]++
	}
	if br.Err() != nil || len(origens) != c.Eventos {
		return ErrManifesto
	}
	for _, classe := range []string{"ERROR", "WARNING", "INFO"} {
		if classes[classe] != c.PorClasse {
			return ErrManifesto
		}
	}
	usados := map[string]bool{}
	conferir := func(campos CamposEntrada) (domain.ClasseLog, error) {
		origem, ok := origens[campos.RecordID]
		if !ok || usados[campos.RecordID] || campos.LevelText != origem.nivel {
			return "", ErrManifesto
		}
		var log map[string]any
		if json.Unmarshal([]byte(campos.MessageText), &log) != nil || !reflect.DeepEqual(log, origem.log) {
			return "", ErrManifesto
		}
		usados[campos.RecordID] = true
		return origem.classe, nil
	}
	for _, r := range registros {
		var gold domain.ClasseLog
		if r.ItemType == "log" {
			classe, err := conferir(r.CamposEntrada)
			if err != nil {
				return ErrManifesto
			}
			gold = classe
		} else {
			if r.Documentacao != m.OrigemDrive {
				return ErrManifesto
			}
			for _, member := range r.Members {
				classe, err := conferir(member)
				if err != nil {
					return ErrManifesto
				}
				gold = domain.MaisGrave(gold, classe)
			}
		}
		if r.ClassLabel != gold {
			return ErrManifesto
		}
	}
	if len(usados) != c.Eventos {
		return ErrManifesto
	}
	return nil
}

func digest(data []byte) string { hash := sha256.Sum256(data); return hex.EncodeToString(hash[:]) }
