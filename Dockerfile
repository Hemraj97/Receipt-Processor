# Use the official Golang image to build the Go app
FROM golang:1.22-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY ./go.mod ./go.sum ./


# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o receipt-processor .

# Start a new stage from a smaller image for the final image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary from the builder stage
COPY --from=builder /app/receipt-processor .

# Expose the port that your app will run on
EXPOSE 8080

# Command to run the application
CMD ["./receipt-processor"]
