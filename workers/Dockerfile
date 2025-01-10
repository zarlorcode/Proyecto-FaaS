# Etapa de construcción
FROM golang:1.22 as builder 

# Instala Docker CLI
RUN apt-get update && apt-get install -y docker.io && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copiar los archivos necesarios
COPY go.mod go.mod
COPY go.sum go.sum
COPY worker.go worker.go

RUN go mod download
RUN go mod tidy
RUN go build -o main /app/worker.go 
ENTRYPOINT ["/app/main"]


