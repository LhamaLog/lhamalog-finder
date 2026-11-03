IMAGE := lhamalog-finder

# Portabilidade: no Windows (cmd/PowerShell) a remoção usa comandos do cmd.
ifeq ($(OS),Windows_NT)
RM := if exist bin rmdir /s /q bin
else
RM := rm -rf bin
endif

.PHONY: run build test lint docker infra-up infra-down infra-config clean

run:
	docker compose up -d mongo redis openobserve
	docker compose up -d --wait redis
	docker compose up --build --no-deps finder

build:
	docker build -t $(IMAGE) .

test:
	docker compose run --rm --no-deps gotools test -race ./...

lint:
	docker compose run --rm --no-deps gotools vet ./...

docker: build

infra-up:
	docker compose up -d mongo redis openobserve

infra-down:
	docker compose down

infra-config:
	docker compose config

clean:
	$(RM)
