package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// AplicarArquivoEnv lê o .env do diretório atual com o GoDotEnv.
// Arquivo ausente não é erro. Variável já definida no terminal permanece.
func AplicarArquivoEnv() error {
	err := godotenv.Load()
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
