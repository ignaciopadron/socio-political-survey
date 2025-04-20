# Etapa 1: Construcción
FROM golang:1.22-alpine AS builder

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de módulos y descargar dependencias
# Copiar primero estos archivos aprovecha el cache de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar la aplicación
# -o /app/main especifica el nombre del archivo de salida
# ./cmd/main.go es la ruta a tu archivo principal
# CGO_ENABLED=0 crea un binario estático (recomendado para contenedores)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/main.go

# Etapa 2: Ejecución
FROM alpine:latest

# Instalar certificados CA (necesario para peticiones HTTPS si las hubiera)
RUN apk --no-cache add ca-certificates

# Establecer directorio de trabajo
WORKDIR /app

# Copiar el binario compilado desde la etapa de construcción
COPY --from=builder /app/main .

# Copiar directorios estáticos y de plantillas
# Asegúrate de que estas rutas coincidan con la estructura de tu proyecto
COPY static ./static
COPY templates ./templates

# Exponer el puerto en el que la aplicación escucha
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]
