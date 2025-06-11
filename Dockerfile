# syntax=docker/dockerfile:1@sha256:9857836c9ee4268391bb5b09f9f157f3c91bb15821bb77969642813b0d00518d
FROM golang:1.24-bullseye@sha256:12152217ef79f60fd71a64375d818b7be68c69c7611ca471c2e2c28324cbf4cd AS build
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
