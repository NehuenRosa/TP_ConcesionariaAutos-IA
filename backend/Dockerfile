# Etapa de compilación
FROM golang:1.25-alpine AS compilacion

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/api ./cmd/api

# Etapa de ejecución
FROM alpine:3.22

# curl para el healthcheck; usuario no-root de menor privilegio.
RUN apk add --no-cache ca-certificates curl && \
    addgroup -S -g 10001 app && \
    adduser -S -G app -u 10001 app

WORKDIR /app
COPY --from=compilacion /bin/api /app/api

EXPOSE 8080
USER app
CMD ["/app/api"]
