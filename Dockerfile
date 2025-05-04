FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /marketflow main.go

FROM debian:stable-slim
COPY --from=builder /marketflow /app/marketflow
ENTRYPOINT ["/app/marketflow"]
