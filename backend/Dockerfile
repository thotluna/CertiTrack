# Dockerfile para el Backend de Go (CertiTrack)

# Usa la imagen oficial de Go como base para la etapa de construcción
FROM golang:1.23-alpine AS builder

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos go.mod y go.sum para descargar las dependencias
COPY go.mod .
COPY go.sum .

# Descarga las dependencias del módulo Go
RUN go mod download

# Copia el resto del código fuente de la aplicación
COPY . .

# Construye la aplicación Go
# La ruta al paquete principal ahora es ./cmd/certitrack-backend
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o certitrack-backend ./cmd/certitrack-backend

# Usa una imagen base mínima para la etapa final (más pequeña y segura)
FROM alpine:latest

# Instala certificados SSL para que el backend pueda hacer llamadas HTTPS (ej. a SES)
RUN apk --no-cache add ca-certificates

# Establece el directorio de trabajo
WORKDIR /root/

# Copia el ejecutable construido desde la etapa 'builder'
COPY --from=builder /app/certitrack-backend .

# Expone el puerto en el que el backend escuchará
EXPOSE 8080

# Comando para ejecutar la aplicación cuando el contenedor se inicie
CMD ["./certitrack-backend"]