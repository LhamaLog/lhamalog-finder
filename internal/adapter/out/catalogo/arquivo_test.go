package catalogo

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type caso struct {
	Entrada  string `json:"entrada"`
	Saida    string `json:"saida"`
	Total    int    `json:"total"`
	Removida bool   `json:"removida"`
}

func TestCatalogoDistribuidoEConformidade(t *testing.T) {
	cat, err := (Arquivo{Caminho: "../../../../config/seguranca/catalogo.json"}).Carregar(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile("../../../../testdata/seguranca/conformidade.json")
	if err != nil {
		t.Fatal(err)
	}
	var casos []caso
	if json.Unmarshal(raw, &casos) != nil {
		t.Fatal("casos inválidos")
	}
	for _, c := range casos {
		t.Run(c.Entrada, func(t *testing.T) {
			r, n, e := cat.Transformar(c.Entrada)
			if e != nil || r.Texto != c.Saida || n != c.Total || r.LinhaRemovida != c.Removida {
				t.Fatalf("resultado divergente: total=%d removida=%v", n, r.LinhaRemovida)
			}
		})
	}
	if _, _, err := cat.Transformar(string([]byte{0xff})); err == nil {
		t.Fatal("aceitou UTF-8 inválido")
	}
}

func TestArquivoImutavelECatalogoEstrito(t *testing.T) {
	good, err := os.ReadFile("../../../../config/seguranca/catalogo.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, data := range []string{
		strings.Replace(string(good), `"versao"`, `"Versao"`, 1),
		strings.Replace(string(good), `"schemaVersion": "1"`, `"schemaVersion": "1", "schemaVersion": "1"`, 1),
		strings.Replace(string(good), `"id": "linha-com-segredo"`, `"id": "email"`, 1),
		strings.Replace(string(good), `"OMITIR_LINHA"`, `"IGNORAR"`, 1),
		strings.Replace(string(good), `"schemaVersion": "1"`, `"schemaVersion": "2"`, 1),
		strings.Replace(string(good), `"regex":`, `"Regex":`, 1),
	} {
		if _, e := compilar([]byte(data)); e == nil {
			t.Fatal("aceitou contrato inválido")
		}
	}
	path := filepath.Join(t.TempDir(), "catalogo.json")
	os.WriteFile(path, append(good, ' '), 0600)
	if _, e := (Arquivo{Caminho: path}).Carregar(context.Background()); e == nil {
		t.Fatal("mudou bytes da versão publicada")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, e := (Arquivo{Caminho: path}).Carregar(ctx); e == nil {
		t.Fatal("ignorou cancelamento")
	}
}
