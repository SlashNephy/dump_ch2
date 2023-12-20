# syntax=docker/dockerfile:1@sha256:ac85f380a63b13dfcefa89046420e1781752bab202122f8f50032edf31be0021
FROM golang:1.21-bullseye@sha256:45e3d8f6c1f76b20da6595077571358d451af9b58270e9054504f7ebe0a8598d AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build

FROM debian:bullseye-slim@sha256:d3d0d14f49b49a4dd98a436711f5646dc39e1c99203ef223d1b6620061e2c0e5
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
