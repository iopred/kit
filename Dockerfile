FROM rust:1.70 as builder

WORKDIR /usr/src/kit

# Copy the source code.
COPY . .

# Copy the .env file into the image.
COPY .env .env

RUN cargo build --release

# Now copy it into our base image.
# FROM gcr.io/distroless/cc-debian10

# Second stage: use a more feature-rich base image for debugging
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

COPY --from=builder /usr/src/kit/target/release/kit /usr/local/bin/kit

CMD ["kit"]

