# Build Stage
FROM golang:1.22.2-alpine3.18 AS builder

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o main cmd/app/main.go

# Final Stage
FROM debian:bookworm-slim

WORKDIR /app

# Download and install Google Chrome
RUN apt-get update && apt-get install -y wget gnupg \
    && wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list' \
    && apt-get update && apt-get install -y google-chrome-stable \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Copy the Go app binary from the builder stage
COPY --from=builder /app .

CMD ["/app/main"]
