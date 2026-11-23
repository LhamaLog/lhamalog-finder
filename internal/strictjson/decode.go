// Package strictjson decodifica contratos fechados sem aliases de caixa ou chaves repetidas.
package strictjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

var ErrContrato = errors.New("json_fora_do_contrato")

// Decode aceita somente os nomes exatos informados. O decoder tipado também
// valida a posição de cada campo e rejeita campos desconhecidos e texto extra.
func Decode(data []byte, dest any, nomes ...string) error {
	allowed := map[string]bool{}
	for _, n := range nomes {
		allowed[n] = true
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	var value func() error
	value = func() error {
		t, err := d.Token()
		if err != nil {
			return ErrContrato
		}
		delim, ok := t.(json.Delim)
		if !ok {
			return nil
		}
		switch delim {
		case '{':
			seen := map[string]bool{}
			for d.More() {
				k, err := d.Token()
				if err != nil {
					return ErrContrato
				}
				key, ok := k.(string)
				if !ok || !allowed[key] || seen[key] {
					return ErrContrato
				}
				seen[key] = true
				if value() != nil {
					return ErrContrato
				}
			}
			t, err = d.Token()
			if err != nil || t != json.Delim('}') {
				return ErrContrato
			}
		case '[':
			for d.More() {
				if value() != nil {
					return ErrContrato
				}
			}
			t, err = d.Token()
			if err != nil || t != json.Delim(']') {
				return ErrContrato
			}
		default:
			return ErrContrato
		}
		return nil
	}
	if value() != nil {
		return ErrContrato
	}
	if _, err := d.Token(); err != io.EOF {
		return ErrContrato
	}
	d = json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(dest); err != nil {
		return ErrContrato
	}
	return nil
}
