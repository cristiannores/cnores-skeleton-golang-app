# Paso 1: Imagen base con el entorno Go
FROM golang:1.20-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Compila la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .

# Paso 2: Construye una imagen limpia y pequeña desde 'scratch'
FROM scratch

# Copia el ejecutable a la imagen limpia
COPY --from=builder /app/myapp .
COPY app/shared/utils/config/config.json .

# El puerto se configura en tiempo de ejecución, basado en el archivo de configuración
# No es necesario exponer un puerto aquí

# Comando para ejecutar la aplicación
CMD ["./myapp"]
