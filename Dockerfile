FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod tidy
COPY . .
RUN go build -o gyanpass ./cmd/main.go
 
 
FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/gyanpass .
COPY --from=builder /app/.env .env

EXPOSE 8080
ENTRYPOINT ["./gyanpass"]
