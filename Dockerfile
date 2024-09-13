FROM golang:1.23.0-alpine AS builder
WORKDIR /data

COPY go.mod go.mod
RUN go mod download

COPY . .
COPY .env .env

RUN go build \
    -a \
    -o avito \
    ./cmd/main.go

FROM alpine:latest

WORKDIR /bin

COPY --from=builder /data /bin
#COPY --from=builder /data/avito /bin
#COPY --from=builder /data/.env /bin

ENTRYPOINT ["/bin/avito"]
