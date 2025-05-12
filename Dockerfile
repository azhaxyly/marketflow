FROM golang:1.24-bullseye AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /marketflow main.go

FROM gcr.io/distroless/base-debian12
COPY --from=builder /marketflow /app/marketflow
ENTRYPOINT ["/app/marketflow"]
