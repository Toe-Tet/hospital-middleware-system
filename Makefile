APP_NAME=hospital-middleware
CMD_DIR=./cmd/api
BIN_DIR=./bin
SRC_DIR=./src

.PHONY: all build run test clean migrate-up migrate-down seed swagger

all: build

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

run:
	go run $(CMD_DIR)/main.go

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	@rm -rf $(BIN_DIR) coverage.out coverage.html

install-tools:
	@mkdir -p bin
	@GOBIN=$(CURDIR)/bin go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1
	@GOBIN=$(CURDIR)/bin go install github.com/swaggo/swag/cmd/swag@v1.8.12

swagger:
	@./bin/swag init -g $(CMD_DIR)/main.go -o ./docs --parseDependency --parseInternal
  
migrate-up:
	./scripts/migrate.sh up

migrate-down:
	./scripts/rollback.sh last

migrate-down-all:
	./scripts/rollback.sh all

seed:
	./scripts/seed.sh

fmt:
	gofmt -s -w .
	goimports -w . || true

tidy:
	go mod tidy
