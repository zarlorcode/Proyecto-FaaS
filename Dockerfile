# Etapa de construcción
FROM golang:1.20 as builder 

WORKDIR /app

# Copiar los archivos necesarios
COPY go.mod go.sum ./
COPY cmd/ ./cmd/
COPY functions/ ./functions/
COPY internal/ ./internal/

RUN go mod download
RUN go mod tidy
RUN go build -o main ./cmd/api/main.go 
ENTRYPOINT ["/app/main"]




