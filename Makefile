IMAGE := lhamalog-finder
COMPOSE := docker compose -f lhamalog-finder/docker-compose.yml

# Portabilidade: no Windows (cmd/PowerShell) o binário é .exe, a remoção usa comandos do cmd
# e o -race fica de fora (exige cgo/gcc no Windows).
ifeq ($(OS),Windows_NT)
BIN := bin\finder.exe
RM := if exist bin rmdir /s /q bin
RACE :=
LINT_EXTRA := where golangci-lint >nul 2>&1 && golangci-lint run || echo golangci-lint nao instalado (opcional)
else
BIN := bin/finder
RM := rm -rf bin
RACE := -race
LINT_EXTRA := command -v golangci-lint >/dev/null && golangci-lint run || echo "golangci-lint não instalado (opcional)"
endif

.PHONY: run build test lint docker infra-up infra-down clean

run:
	go run ./cmd/finder

build:
	go build -o $(BIN) ./cmd/finder

test:
	go test $(RACE) ./...

lint:
	go vet ./...
	@$(LINT_EXTRA)

docker:
	docker build -t $(IMAGE) .

infra-up:
	docker compose up -d

infra-down:
	docker compose down

infra-config:
	docker compose config

clean:
	$(RM)
