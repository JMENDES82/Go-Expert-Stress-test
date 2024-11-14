# Etapa de construção
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o stresstest

# Etapa final
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/stresstest .
ENTRYPOINT ["./stresstest"]
