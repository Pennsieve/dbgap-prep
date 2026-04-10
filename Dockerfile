FROM golang:1.26-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN ["go", "mod", "download"]

COPY cmd cmd
COPY internal internal

RUN ["go", "build", "-o", "main", "cmd/app/main.go"]

FROM alpine:3.21

WORKDIR /app

RUN mkdir -p data

COPY --from=builder /app/main .

ENTRYPOINT [ "/app/main" ]
