# Etapa de compilación
FROM golang:1.25-alpine AS compilacion

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/api ./cmd/api

# Etapa de ejecución
FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=compilacion /bin/api /app/api

EXPOSE 8080
CMD ["/app/api"]
