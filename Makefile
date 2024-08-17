CMD_DIR=cmd
INTERNAL_DIR=internal
BINARY_NAME=build/main
POSTGRES_DSN_DEV=postgres://postgres:postgres@localhost:5432/app

define run_migration
	psql $(POSTGRES_DSN_DEV) -v "ON_ERROR_STOP=1" -f $(1) || exit $?
endef


build:
	go build -o $(BINARY_NAME) $(CMD_DIR)/main.go

run: build
	./$(BINARY_NAME)

dev:
	go run $(CMD_DIR)/main.go

test:
	go test ./$(INTERNAL_DIR)/...

clean:
	go clean
	rm -f $(BINARY_NAME)

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

all: fmt vet test build run

migrate-dev:
	$(call run_migration,./migrations/init.psql)
	$(call run_migration,./migrations/test_data.psql)

.PHONY: build run dev test clean fmt vet lint all migrate-dev
