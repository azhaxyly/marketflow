BINARY = marketflow

.PHONY: fmt build run

fmt:
	go run mvdan.cc/gofumpt -w .

build:
	go build -o $(BINARY) ./cmd/marketflow

run: build
	./$(BINARY) --help
