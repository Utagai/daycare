# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod ./
RUN go mod tidy
RUN go mod download

COPY . .

RUN go build -o daycare .

# Test environment stage
FROM alpine:3.19

RUN apk add --no-cache \
    git \
    bash \
    zsh \
    go \
    ca-certificates

COPY --from=builder /build/daycare /usr/local/bin/daycare

RUN mkdir -p /root/projects && \
    mkdir -p /root/.config/daycare && \
    git config --global user.email "test@example.com" && \
    git config --global user.name "Test User"
