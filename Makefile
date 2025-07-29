BIN_DIR := ./bin
SERVER_BIN := $(BIN_DIR)/server
CLIENT_BIN := $(BIN_DIR)/client

.PHONY: start
start:
	docker stop $$(docker ps -aq) \
	&& docker compose start

.PHONY: gen
gen:
	go generate ./...

.PHONY: test
test:
	./coverage.sh

.PHONY: build-server
build-server:
	echo "Сборка сервера..."
	go build -o $(SERVER_BIN) ./cmd/server

.PHONY: build-client
build-client:
	echo "Сборка клиента..."
	go build $(LDFLAGS) -o $(CLIENT_BIN) -ldflags "-X 'main.version=1.0.0' -X 'main.buildDate=$$(date +'%Y/%m/%d %H:%M:%S')'" ./cmd/client

.PHONY: build
build: build-server build-client