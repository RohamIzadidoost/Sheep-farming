# Build stage
FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o sheep-farm ./cmd/api

# Runtime stage
FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/sheep-farm ./sheep-farm
EXPOSE 8080
CMD ["./sheep-farm"]
