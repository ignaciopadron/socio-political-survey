# Etapa 1: Build
FROM golang:1.24.2-alpine AS builder

# Metadata y argumentos
ARG VERSION=1.0.0
LABEL org.opencontainers.image.version=$VERSION

# Crear user en build para permisos coherentes
RUN addgroup -g 10001 radar \
  && adduser -D -u 10001 -G radar radar

WORKDIR /app

# Dependencias (cache)
COPY go.mod ./
RUN apk add --no-cache ca-certificates \
  && go mod download

# Código fuente
COPY . .

# Compilar binario estático sin símbolos
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
      -ldflags="-s -w" \
      -o radar ./cmd/main.go

# Etapa 2: Runtime mínimo
FROM scratch AS runtime

# Copiar certificados para HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Mantener mismo usuario
USER 10001:10001
WORKDIR /app

# Copiar binario y assets
COPY --from=builder /app/radar .
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

# Metadatos finales
LABEL org.opencontainers.image.title="radar"
LABEL org.opencontainers.image.description="Aplicación de encuesta sociopolítica"

# Exponer y comprobar salud
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s \
  CMD ["./radar", "health"]

ENTRYPOINT ["./radar"]
