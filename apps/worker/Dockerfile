FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY . .

# Compilation optimisée (désactive le CGO pour une image ultra légère)
RUN CGO_ENABLED=0 GOOS=linux go build -o worker-service ./cmd/worker

# Étape finale (image ultra légère)
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/worker-service .

CMD ["./worker-service"]
