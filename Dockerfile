# Etapa de construcción
FROM golang:1.20 as builder 

WORKDIR /app
COPY . . 
RUN go mod download
RUN go build -o main ./cmd/api/main.go 
ENTRYPOINT ["/app/main"]



