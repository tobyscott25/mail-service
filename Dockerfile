# Use the official Go image to create a build artifact.
FROM golang:1.21 AS builder

# Copy local code to the container image.
WORKDIR /app
COPY . .

# Build the command inside the container.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server

# Use a minimal alpine image for production.
FROM alpine:3.14

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /server

# Run the web service on container startup.
CMD ["/server"]
