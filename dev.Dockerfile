FROM golang:1.25-alpine AS development

WORKDIR /app

RUN apk add --no-cache \
  make \
  bash

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 5000

CMD ["make", "dev"]
