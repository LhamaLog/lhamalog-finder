package http

import (
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	NovoRoteador().ServeHTTP(rec, httptest.NewRequest("GET", "/health", nil))
	if rec.Code != 200 {
		t.Fatalf("status = %d", rec.Code)
	}
}
