# syntax=docker/dockerfile:1@sha256:e63addfe27b10e394a5f9f1e866961adc70d08573f1cb80f8d1a0999347b3553
FROM golang:1.24-bullseye@sha256:dfd72198d14bc22f270c9e000c304a2ffd19f5a5f693fad82643311afdc6b568 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build

FROM debian:bullseye-slim@sha256:779034981fec838da124ff6ab9211499ba5d4e769dabdfd6c42c6ae2553b9a3b
WORKDIR /app

RUN groupadd -g 1000 app && useradd -u 1000 -g app app

RUN <<EOT
  apt-get update
  apt-get install -yqq --no-install-recommends ca-certificates
  rm -rf /var/lib/apt/lists/*
EOT

USER app
COPY --from=build /app/dump_ch2 ./
CMD ["./dump_ch2"]
