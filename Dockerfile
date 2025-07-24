# syntax=docker/dockerfile:1@sha256:9857836c9ee4268391bb5b09f9f157f3c91bb15821bb77969642813b0d00518d
FROM golang:1.24-bullseye@sha256:62ba6b19de03e891f7fa1001326bd48411f2626ff35e7ba5b9d890711ce581d9 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build

FROM debian:bullseye-slim@sha256:c2c58af6e3ceeb3ed40adba85d24cfa62b7432091597ada9b76b56a51b62f4c6
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
