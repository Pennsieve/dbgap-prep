FROM golang:1.23-alpine

WORKDIR /app

RUN mkdir -p data

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN ["go", "mod", "download"]

COPY cmd cmd
COPY internal internal

RUN ["go", "build", "-o", "main", "cmd/app/main.go"]

ENTRYPOINT [ "/app/main" ]