package labeleddataset

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var ErrBaseIndisponivel = errors.New("base_real_rotulada_indisponivel")

// LoadConjunto valida cada registro JSONL. ReadString evita teto artificial de 64KiB.
// Erros de registro trazem número da linha, nunca mensagem/PII da entrada.
// Conflitos de rótulos são verificados no conjunto completo.
func LoadConjunto(reader io.Reader, documentacao ...DocumentacaoConjunto) ([]Registro, error) {
	registros, err := decodificarRegistros(reader)
	if err != nil {
		return nil, err
	}
	if err := ValidarConjunto(registros, documentacao...); err != nil {
		return nil, err
	}
	return registros, nil
}

// A validação de conflitos pertence à coleção completa, não a cada arquivo.
func decodificarRegistros(reader io.Reader) ([]Registro, error) {
	br := bufio.NewReader(reader)
	var registros []Registro
	for linha := 1; ; linha++ {
		texto, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("linha %d: leitura falhou", linha)
		}
		if strings.TrimSpace(texto) != "" {
			if err := validarChavesJSON(texto); err != nil {
				return nil, fmt.Errorf("linha %d: %w", linha, ErrRegistro)
			}
			var r Registro
			dec := json.NewDecoder(strings.NewReader(texto))
			dec.DisallowUnknownFields()
			if e := dec.Decode(&r); e != nil {
				return nil, fmt.Errorf("linha %d: %w", linha, ErrRegistro)
			}
			var extra any
			if e := dec.Decode(&extra); e != io.EOF {
				return nil, fmt.Errorf("linha %d: %w", linha, ErrRegistro)
			}
			if e := ValidateRegistro(r); e != nil {
				return nil, fmt.Errorf("linha %d: %w", linha, e)
			}
			registros = append(registros, r)
		}
		if err == io.EOF {
			break
		}
	}
	return registros, nil
}
func LoadDiretorio(dir string, documentacao ...DocumentacaoConjunto) ([]Registro, error) {
	nomes, err := filepath.Glob(filepath.Join(dir, "*.jsonl"))
	if err != nil || len(nomes) == 0 {
		return nil, ErrBaseIndisponivel
	}
	var registros []Registro
	for _, nome := range nomes {
		f, e := os.Open(nome)
		if e != nil {
			return nil, ErrBaseIndisponivel
		}
		r, e := decodificarRegistros(f)
		f.Close()
		if e != nil {
			return nil, e
		}
		registros = append(registros, r...)
	}
	if len(registros) == 0 {
		return nil, ErrBaseIndisponivel
	}
	if err := ValidarConjunto(registros, documentacao...); err != nil {
		return nil, err
	}
	return registros, nil
}
func IterRegistros(r []Registro, visitar func(Registro) error, documentacao ...DocumentacaoConjunto) error {
	if err := ValidarConjunto(r, documentacao...); err != nil {
		return err
	}
	for _, item := range r {
		if err := visitar(item); err != nil {
			return err
		}
	}
	return nil
}

// A representação física deve ter um único gold; encoding/json sozinho aceitaria
// nomes repetidos escolhendo o último. A verificação também cobre objetos members.
func validarChavesJSON(texto string) error {
	dec := json.NewDecoder(strings.NewReader(texto))
	var valor func() error
	valor = func() error {
		token, err := dec.Token()
		if err != nil {
			return err
		}
		delim, ok := token.(json.Delim)
		if !ok {
			return nil
		}
		switch delim {
		case '{':
			vistos := map[string]bool{}
			for dec.More() {
				chave, err := dec.Token()
				if err != nil {
					return err
				}
				nome, ok := chave.(string)
				if !ok || vistos[nome] || !chaveCanonica(nome) {
					return ErrRegistro
				}
				vistos[nome] = true
				if err := valor(); err != nil {
					return err
				}
			}
			fim, err := dec.Token()
			if err != nil || fim != json.Delim('}') {
				return ErrRegistro
			}
		case '[':
			for dec.More() {
				if err := valor(); err != nil {
					return err
				}
			}
			fim, err := dec.Token()
			if err != nil || fim != json.Delim(']') {
				return ErrRegistro
			}
		default:
			return ErrRegistro
		}
		return nil
	}
	if err := valor(); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF {
		return ErrRegistro
	}
	return nil
}

// encoding/json aceita aliases sem distinguir caixa; o contrato usa nomes exatos.
// DisallowUnknownFields continua responsável por rejeitar campos fora do contexto,
// como class_label em members, mesmo quando o valor é null.
func chaveCanonica(nome string) bool {
	switch nome {
	case "itemType", "level_text", "message_text", "lhamalog_batch_id", "lhamalog_record_id", "tenantId", "productId", "split", "notas", "class_label", "members", "documentacao":
		return true
	}
	return false
}
