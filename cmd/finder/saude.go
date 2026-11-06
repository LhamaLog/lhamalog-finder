package main

import (
	"net/http"

	httpin "github.com/LhamaLog/lhamalog-finder/internal/adapter/in/http"
)

// servidor atende a consulta de saúde já existente. Não cria outra consulta.
func servidor(addr string) *http.Server {
	return &http.Server{Addr: addr, Handler: httpin.NovoRoteador()}
}
