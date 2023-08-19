# Always prefer the latest stable Go version.
FROM golang:1.21-alpine AS builder

# Create and change to the app directory.
WORKDIR /app

# Copy go mod files
COPY ../go.mod ./
COPY ../go.sum ./

# Install the required dependencies.
RUN go mod download

# Copy source/configuration code.
COPY ../ ./

# Build the executable file.
RUN go build -v -o server ./cmd/main.go

# Use multi-image container for a lean production container.
FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates

RUN apk update && apk add postgresql

WORKDIR /bin

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /bin/server

# Run the web service on container startup.
ENTRYPOINT [ "/bin/server" ]
