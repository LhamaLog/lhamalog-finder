package labeleddataset

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBaseRealRotulada(t *testing.T) {
	if os.Getenv("FINDER_DATASET_INTEGRATION") != "1" {
		t.Skip("base real anonimizada/origem/autoria indisponiveis; fixtures nao fecham SC-001/003/005")
	}
	dir := filepath.Join("..", "..", "datasets", "labeled-logs")
	registros, err := LoadDiretorio(dir)
	if err != nil {
		t.Fatal(err)
	}
	origem, err := os.ReadFile(filepath.Join(dir, "ORIGEM.md"))
	if err != nil || strings.TrimSpace(string(origem)) == "" {
		t.Fatal("declaracao de origem pendente")
	}
	// Proveniência/privacidade e cobertura de milhares requerem revisão da
	// declaração e manifest real; nenhum volume mínimo sem fonte é inventado aqui.
	m, err := Avaliar(registros)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.ValidarGate(); err != nil {
		t.Fatal(err)
	}
}
