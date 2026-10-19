// Package http expõe endpoints internos (health).
package http

import "net/http"

// NovoRoteador cria o roteador com /health.
func NovoRoteador() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"UP"}`))
	})
	return mux
}
