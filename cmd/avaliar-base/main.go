package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/LhamaLog/lhamalog-finder/internal/labeleddataset"
)

func main() {
	dir := flag.String("diretorio", "datasets/labeled-logs", "base versionada com manifesto e proveniência")
	flag.Parse()
	m, registros, err := labeleddataset.CarregarBasePublicada(*dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	resultado, err := labeleddataset.Avaliar(registros)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := resultado.ValidarGate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_ = json.NewEncoder(os.Stdout).Encode(struct {
		Versao    string
		Populacao string
		Relatorio labeleddataset.Relatorio
	}{m.Versao, m.TipoOrigem, resultado})
}
