# First stage: Build the Rust application
FROM rust:1.83 as kit_builder

WORKDIR /usr/src/kit

# Copy the source code and .env file.
COPY . .
COPY .env .env

# Build the Rust application in release mode.
RUN cargo build --release

# Second stage: Build the Go application
FROM golang:1.22 as qr_kit_builder

WORKDIR /qr.kit

# Copy go.mod and go.sum files, then download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the Go source files
COPY *.go ./

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /qr.kit

# Final stage: Create a lightweight image with both Rust and Go binaries
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the Rust binary from the first stage
COPY --from=kit_builder /usr/src/kit/target/release/kit /usr/local/bin/kit

# Copy the Go binary from the second stage
COPY --from=qr_kit_builder /qr.kit /usr/local/bin/qr.kit

# Expose a port for the Go application (if needed)
EXPOSE 3242

# Ensure the binaries are executable
RUN chmod +x /usr/local/bin/kit /usr/local/bin/qr.kit

# Default command (you might need to adjust this to suit your needs)
CMD ["/bin/bash"]

# Optionally, if you want to run one of the applications by default, uncomment one of the following lines:
# CMD ["kit"]
# CMD ["qr.kit"]
