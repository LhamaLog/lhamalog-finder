package labeleddataset

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestProvenienciaImpedeAlterarOuroOuInventarMensagem(t *testing.T) {
	original := filepath.Join("..", "..", "datasets", "labeled-logs")
	for _, caso := range []string{"ouro trocado mantendo balanceamento", "mensagem fabricada"} {
		t.Run(caso, func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, "datasets", "labeled-logs")
			for _, name := range []string{"manifest.json", "logs.jsonl", "apontamentos.jsonl", "proveniencia/captura.jsonl.gz"} {
				data, err := os.ReadFile(filepath.Join(original, name))
				if err != nil {
					t.Fatal(err)
				}
				path := filepath.Join(dir, name)
				if os.MkdirAll(filepath.Dir(path), 0700) != nil {
					t.Fatal("diretório")
				}
				if os.WriteFile(path, data, 0600) != nil {
					t.Fatal("cópia")
				}
			}
			code, err := os.ReadFile(filepath.Join("..", "..", "cmd", "capturar-base", "main.go"))
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(root, "cmd", "capturar-base", "main.go")
			os.MkdirAll(filepath.Dir(path), 0700)
			os.WriteFile(path, code, 0600)
			m, _, err := CarregarBasePublicada(dir)
			if err != nil {
				t.Fatalf("controle antes da alteração: %v", err)
			}
			f, err := os.Open(filepath.Join(dir, "logs.jsonl"))
			if err != nil {
				t.Fatal(err)
			}
			logs, err := LoadConjunto(f)
			f.Close()
			if err != nil {
				t.Fatal(err)
			}
			if caso == "ouro trocado mantendo balanceamento" {
				logs[0].ClassLabel, logs[1000].ClassLabel = logs[1000].ClassLabel, logs[0].ClassLabel
			} else {
				logs[0].MessageText = `{"level":"ERROR","msg":"mensagem fabricada"}`
			}
			var data []byte
			for _, r := range logs {
				line, _ := json.Marshal(r)
				data = append(data, line...)
				data = append(data, '\n')
			}
			os.WriteFile(filepath.Join(dir, "logs.jsonl"), data, 0600)
			for i := range m.Arquivos {
				if m.Arquivos[i].Nome == "logs.jsonl" {
					m.Arquivos[i].SHA256 = digest(data)
				}
			}
			raw, _ := json.Marshal(m)
			os.WriteFile(filepath.Join(dir, "manifest.json"), raw, 0600)
			if _, _, err := CarregarBasePublicada(dir); err == nil {
				t.Fatal("aceitou dados que não correspondem à captura, mesmo com hashes atualizados")
			}
		})
	}
}
