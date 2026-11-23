package labeleddataset

import (
	"path/filepath"
	"testing"
)

func TestBaseRealRotulada(t *testing.T) {
	dir := filepath.Join("..", "..", "datasets", "labeled-logs")
	manifesto, registros, err := CarregarBasePublicada(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("origem controlada=%s unidades=%d mínimo por classe=%d", manifesto.OrigemDrive, manifesto.Registros, manifesto.Minimo)
	m, err := Avaliar(registros)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.ValidarGate(); err != nil {
		t.Fatal(err)
	}
}
