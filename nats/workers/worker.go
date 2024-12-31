package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/nats-io/nats.go"
)

func main() {
	// Conectarse a NATS
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Error conectando a NATS: %v", err)
	}
	defer nc.Close()

	// Iniciar contexto JetStream
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Error inicializando JetStream: %v", err)
	}

	// Crear Streams (activations y results)
	createStreams(js)

	// Suscribirse al tema "activations.*"
	sub, err := js.Subscribe("activations.*", func(msg *nats.Msg) {
		data := strings.Split(string(msg.Data), "|")
		if len(data) != 3 {
			log.Println("Error: Formato de mensaje inválido")
			js.Publish("results", []byte("Error: Formato de mensaje inválido"))
			return
		}
		subject := msg.Subject
		requestId := strings.Split(subject, ".")

		username, functionName, parameter := data[0], data[1], data[2]
		log.Printf("Procesando activación: Usuario=%s, Función=%s, Parámetro=%s\n", username, functionName, parameter)

		// Procesar la función
		result, err := processFunction(functionName, parameter)
		if err != nil {
			log.Printf("Error ejecutando función para %s: %s\n", username, err.Error())
			js.Publish("results."+requestId[1], []byte(fmt.Sprintf("Error para %s: %s", username, err.Error())))
			return
		}

		// Publicar el resultado en "results"
		js.Publish("results."+requestId[1], []byte(result))
		log.Printf("Resultado enviado para %s\n", username)
	})

	if err != nil {
		log.Fatalf("Error suscribiéndose a activations: %v", err)
	}

	defer sub.Unsubscribe()
	select {} // Mantener el Worker activo
}

// createStreams se asegura de que los streams "activations" y "results" estén configurados
func createStreams(js nats.JetStreamContext) {
	// Configurar el stream "activations"
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "activations",         // Nombre del stream
		Subjects: []string{"activations.*"}, // Temas relacionados
		Storage:  nats.FileStorage,     // Almacenamiento en disco
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Printf("Error creando stream 'activations': %v", err)
	} else {
		log.Println("Stream 'activations' configurado")
	}

	// Configurar el stream "results"
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "results",         // Nombre del stream
		Subjects: []string{"results.*"}, // Temas relacionados
		Storage:  nats.FileStorage, // Almacenamiento en disco
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Printf("Error creando stream 'results': %v", err)
	} else {
		log.Println("Stream 'results' configurado")
	}
}

// processFunction ejecuta una función utilizando Docker
func processFunction(functionName, parameter string) (string, error) {
	// Simulación de ejecución con Docker
    log.Println("Worker PROCESANDO la función ", functionName)
	cmd := exec.Command("docker", "run", "--rm", functionName, parameter)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error ejecutando la función: %s", stderr.String())
	}
	return out.String(), nil
}
