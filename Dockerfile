FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/server ./cmd/api/main.go

FROM golang:1.22-alpine

WORKDIR /app

# Install language runtimes for code execution
RUN apk add --no-cache \
    nodejs npm \
    python3 py3-pip \
    php \
    rust cargo \
    openjdk17-jdk \
    git \
    bash

COPY --from=builder /app/server .
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations

ENV config=docker

EXPOSE 8080 8081

CMD ["./server"]
