ARG GO_VERSION=1.24.4

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go install github.com/air-verse/air@latest

COPY . .

EXPOSE 5000
