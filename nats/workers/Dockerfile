# Dockerfile para el Worker
FROM golang:1.20

WORKDIR /app

# Copiar los archivos del Worker
COPY . .

# Compilar el Worker
RUN go mod tidy && go build -o worker

# Comando para iniciar el Worker
CMD ["./worker"]
