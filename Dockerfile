FROM golang:1.25-alpine

RUN apk add --no-cache \
    git \
    bash \
    zsh \
    ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /usr/local/bin/daycare .

RUN mkdir -p /root/projects && \
    mkdir -p /root/.config/daycare && \
    git config --global user.email "test@example.com" && \
    git config --global user.name "Test User"

WORKDIR /build/e2e
COPY go.mod go.sum ./
COPY e2e/ ./
CMD ["go", "test", "-v", "."]
