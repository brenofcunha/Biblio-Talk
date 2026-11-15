.PHONY: run build test vet dev-db-up dev-db-down migrate-up migrate-down format lint clean

# ── Execução ────────────────────────────────────────────────
run:
	go run ./cmd/api

build:
	go build -o bin/lume-api ./cmd/api

# ── Qualidade ───────────────────────────────────────────────
test:
	go test ./...

vet:
	go vet ./...

format:
	gofmt -w .

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run; else echo "golangci-lint não instalado — pulando"; fi

# ── Banco de desenvolvimento ───────────────────────────────
dev-db-up:
	docker compose up -d postgres

dev-db-down:
	docker compose down

# ── Migrações ───────────────────────────────────────────────
migrate-up:
	go run ./cmd/migrate -direction=up

migrate-down:
	go run ./cmd/migrate -direction=down

clean:
	rm -rf bin
