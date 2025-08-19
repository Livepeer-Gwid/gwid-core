FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app ./cmd

FROM alpine:3.21 AS final

COPY --from=builder /app /bin/app

EXPOSE 5000

CMD ["bin/app"]
