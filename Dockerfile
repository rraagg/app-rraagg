# syntax=docker/dockerfile:1
FROM golang:1.22.1-alpine3.19

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]