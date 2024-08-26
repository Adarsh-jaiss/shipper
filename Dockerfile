# Build stage
FROM golang:1.22-alpine3.18 AS build

# Set the working directory to /app
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and install any required Go dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shipper .

# Final stage
FROM alpine:3.18

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN adduser -D -g '' appuser

# Set the working directory to /app
WORKDIR /app

# Copy the .env file from the build context to the final image
COPY --from=build /app/.env ./.env

# Copy the binary from the build stage
COPY --from=build /app/shipper .

# Create the .kube directory and set ownership
RUN mkdir -p /home/appuser/.kube && chown -R appuser:appuser /home/appuser/.kube

# Set the user to appuser
USER appuser

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./shipper"]