package application

import (
	"context"
	"errors"
	"testing"
)

func TestBootstrapRecusaCatalogoSemTransformacao(t *testing.T) {
	_, err := InicializarSeguranca(context.Background(), catalogoTeste{versao: "combinada"}, "combinada")
	if !errors.Is(err, ErrCatalogoSeguranca) {
		t.Fatalf("catalogo sem transformacao declarado pronto: %v", err)
	}
}
