# Build stage
FROM golang:1.23.2-alpine3.20 AS builder

# Installing build dependencies
RUN apk --no-cache add build-base

# Setting up the working directory and copying dependencies
WORKDIR /lp
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .

# Building the application
RUN go build -o lp ./cmd/main.go

# Final stage
FROM alpine:3.20

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /lp

# Copying the binary
COPY --from=builder /lp/lp .
COPY ./migrate.sh /lp/migrate.sh

# Setting permissions
RUN chmod +x ./lp \
    && chown -R appuser:appgroup /lp

# Specifying a port for a container
EXPOSE 8002

# Switch to non-root user
USER appuser

# Launch command
ENTRYPOINT ["./lp", "serve", "--config=/lp/config.yml"]