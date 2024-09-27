FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o gyanpass ./cmd/main.go
 
 
FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/gyanpass .
EXPOSE 8080
ENTRYPOINT ["./gyanpass"]
