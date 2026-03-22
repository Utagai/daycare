FROM golang:1.25-alpine

RUN apk add --no-cache \
    git \
    bash \
    zsh 

WORKDIR /build

RUN git config --global user.email "test@example.com" && \
    git config --global user.name "Test User"
COPY . .
RUN go build -o /usr/local/bin/daycare .

RUN mkdir -p /root/projects && \
    mkdir -p /root/.config/daycare

WORKDIR /build/e2e
CMD ["go", "test", "-tags=e2e", "-v", "."]
