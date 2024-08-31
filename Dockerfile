# Build stage
FROM golang:1.22-alpine3.18 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shipper .

# Final stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates curl bash

# Install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
    && chmod +x kubectl \
    && mv kubectl /usr/local/bin/

# Create a non-root user
RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=build /app/shipper .
COPY --from=build /app/.env ./.env

RUN mkdir -p /home/appuser/.kube && chown -R appuser:appuser /home/appuser/.kube

USER appuser

EXPOSE 8080

CMD ["./shipper"]