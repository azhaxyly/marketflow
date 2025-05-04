BINARY = marketflow

.PHONY: fmt build run

fmt:
	go run mvdan.cc/gofumpt -w .

build:
	go build -o $(BINARY) .

run: build
	./$(BINARY) --help
