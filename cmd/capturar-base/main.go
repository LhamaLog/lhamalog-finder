// Captura diagnósticos de operações Redis efetivamente executadas. Os estímulos
// são controlados; não representa tráfego de produção nem usa dados de clientes.
package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type evento struct {
	RecordID string          `json:"recordId"`
	Classe   string          `json:"classeOrigem"`
	Log      json.RawMessage `json:"evento"`
}

func executar() error {
	dir := flag.String("destino", "", "diretório novo da captura")
	count := flag.Int("por-classe", 1020, "eventos por classe; incluindo membros dos apontamentos")
	flag.Parse()
	if *dir == "" || *count <= 0 || os.Getenv("FINDER_REDIS_ADDR") == "" {
		return fmt.Errorf("configuracao_captura_ausente")
	}
	if err := os.Mkdir(*dir, 0700); err != nil {
		return fmt.Errorf("destino_captura_indisponivel")
	}
	f, err := os.OpenFile(filepath.Join(*dir, "captura.jsonl"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("arquivo_captura_indisponivel")
	}
	defer f.Close()
	client := redis.NewClient(&redis.Options{Addr: os.Getenv("FINDER_REDIS_ADDR"), Password: os.Getenv("FINDER_REDIS_PASSWORD"), MaxRetries: -1})
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	if client.Ping(ctx).Err() != nil {
		return fmt.Errorf("redis_captura_indisponivel")
	}
	var random [12]byte
	if _, err := rand.Read(random[:]); err != nil {
		return fmt.Errorf("id_captura_indisponivel")
	}
	prefix := "finder-capture:" + hex.EncodeToString(random[:]) + ":"
	chave := prefix + "value"
	lista := prefix + "list"
	hash := prefix + "hash"
	ausente := prefix + "absent"
	defer client.Del(context.Background(), chave, lista, hash, ausente)
	if client.Set(ctx, chave, "nao-inteiro", time.Minute).Err() != nil || client.RPush(ctx, lista, "item").Err() != nil || client.Expire(ctx, lista, time.Minute).Err() != nil || client.HSet(ctx, hash, "campo", "item").Err() != nil || client.Expire(ctx, hash, time.Minute).Err() != nil {
		return fmt.Errorf("preparo_captura_falhou")
	}
	info, err := client.Info(ctx, "server").Result()
	if err != nil {
		return fmt.Errorf("versao_redis_indisponivel")
	}
	serverVersion := ""
	for _, line := range strings.Split(info, "\n") {
		if strings.HasPrefix(line, "redis_version:") {
			serverVersion = strings.TrimSpace(strings.TrimPrefix(line, "redis_version:"))
		}
	}
	enc := json.NewEncoder(f)
	inicio := time.Now().UTC()
	total := 0
	for _, classe := range []string{"INFO", "WARNING", "ERROR"} {
		for i := 0; i < *count; i++ {
			started := time.Now()
			var result error
			operation := ""
			level := slog.LevelInfo
			message := "redis operation completed"
			switch classe {
			case "INFO":
				switch i % 6 {
				case 0:
					operation = "PING"
					result = client.Ping(ctx).Err()
				case 1:
					operation = "GET"
					result = client.Get(ctx, chave).Err()
				case 2:
					operation = "LLEN"
					result = client.LLen(ctx, lista).Err()
				case 3:
					operation = "HGET"
					result = client.HGet(ctx, hash, "campo").Err()
				case 4:
					operation = "EXISTS"
					result = client.Exists(ctx, chave).Err()
				case 5:
					operation = "TTL"
					result = client.TTL(ctx, chave).Err()
				}
				if result != nil {
					return fmt.Errorf("cenario_info_falhou")
				}
				if i%2 == 1 {
					level = slog.LevelDebug
				}
			case "WARNING":
				level = slog.LevelWarn
				message = "redis optional lookup unavailable"
				switch i % 3 {
				case 0:
					operation = "GET_ABSENT"
					result = client.Get(ctx, ausente).Err()
				case 1:
					operation = "HGET_ABSENT"
					result = client.HGet(ctx, hash, "ausente").Err()
				case 2:
					operation = "LINDEX_ABSENT"
					result = client.LIndex(ctx, lista, 10).Err()
				}
				if result != redis.Nil {
					return fmt.Errorf("cenario_warning_falhou")
				}
			case "ERROR":
				level = slog.LevelError
				message = "redis operation failed"
				switch i % 4 {
				case 0:
					operation = "GET_WRONGTYPE"
					result = client.Get(ctx, lista).Err()
				case 1:
					operation = "INCR_NOT_INTEGER"
					result = client.Incr(ctx, chave).Err()
				case 2:
					operation = "UNKNOWN_COMMAND"
					result = client.Do(ctx, "CAPTURE_CONTROLLED_UNKNOWN").Err()
				case 3:
					operation = "SET_MISSING_ARGUMENT"
					result = client.Do(ctx, "SET").Err()
				}
				if result == nil {
					return fmt.Errorf("cenario_error_nao_falhou")
				}
			}
			var buffer bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buffer, &slog.HandlerOptions{Level: slog.LevelDebug}))
			attrs := []any{"operation", operation, "sequence", i, "elapsed_us", time.Since(started).Microseconds()}
			if result != nil {
				attrs = append(attrs, "diagnostic", result.Error())
			}
			logger.Log(ctx, level, message, attrs...)
			total++
			if err := enc.Encode(evento{RecordID: fmt.Sprintf("r%06d", total), Classe: classe, Log: json.RawMessage(bytes.TrimSpace(buffer.Bytes()))}); err != nil {
				return fmt.Errorf("gravacao_captura_falhou")
			}
		}
	}
	metadata := map[string]any{"schemaVersion": "1", "tipo": "captura-controlada-redis", "goVersion": runtime.Version(), "redisVersion": serverVersion, "startedAt": inicio, "finishedAt": time.Now().UTC(), "eventosPorClasse": *count, "eventos": total, "instrumentacao": "cmd/capturar-base/main.go", "dadosClientes": false}
	raw, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(*dir, "sessao.json"), raw, 0600); err != nil {
		return fmt.Errorf("sessao_captura_nao_gravada")
	}
	fmt.Printf("captura concluida: %d eventos reais controlados; Go %s; Redis %s\n", total, runtime.Version(), serverVersion)
	return f.Sync()
}

func main() {
	if err := executar(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
