# Use the official Golang image to create a build artifact.
# This is a multi-stage build to keep the final image small.
FROM golang:1.21-alpine as builder

WORKDIR /app

# Copy the Go module files and download dependencies.
# This leverages Docker layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code.
COPY . .

# Build the Go app.
# -o /versify places the output binary in the root of the container.
# CGO_ENABLED=0 is important for creating a static binary that can run in a minimal image.
RUN CGO_ENABLED=0 go build -o /versify .

# ---

# Start a new, minimal image.
FROM alpine:latest

# Copy the built binary from the builder stage.
COPY --from=builder /versify /versify

# The entrypoint script will execute the binary with the correct flags.
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
