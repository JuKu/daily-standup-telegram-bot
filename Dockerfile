# Base image
FROM golang:1.20-alpine AS builder

# ADDED
ENV GO111MODULE=on

# Set working directory
WORKDIR /app

# Copy the Go modules and dependencies files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bot .
#RUN go build -o bot .

# Use a minimal image for the final build
FROM alpine:3.18

# Set timezone if needed (optional)
RUN apk add --no-cache tzdata

# Copy the built binary from the builder stage
COPY --from=builder /app/bot /usr/local/bin/bot

# Set environment variable for the Bot API key
ENV TELEGRAM_BOT_API_KEY=""

# Run the bot
ENTRYPOINT ["bot"]
