# syntax=docker/dockerfile:1@sha256:a57df69d0ea827fb7266491f2813635de6f17269be881f696fbfdf2d83dda33e
FROM golang:1.22-bullseye@sha256:0742798b94c656e75ba5409583c2485f18b52de5c9b7cf12859c945bb7414718 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build

FROM debian:bullseye-slim@sha256:0e75382930ceb533e2f438071307708e79dc86d9b8e433cc6dd1a96872f2651d
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
