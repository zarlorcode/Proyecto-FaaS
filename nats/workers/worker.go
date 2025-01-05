package main
import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
    "os"
    "os/signal"
	"github.com/nats-io/nats.go"
    //"github.com/google/uuid" 
    "context"
    "time"
    "github.com/nats-io/nats.go/jetstream"
)

func main() {
    // Configurar el nombre del Worker
    idWorker:= os.Args[1]
	workerName := fmt.Sprintf("worker_%s",idWorker)
    
    // Conectarse a NATS
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Error conectando a NATS: %v", err)
	}
	defer nc.Close()
    
    // Iniciar contexto JetStream
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}
	
	ctx:=context.Background()
	// Crear consumidor (Worker Consumer)
	consumer, err := js.CreateOrUpdateConsumer(ctx, "activations", jetstream.ConsumerConfig{
		Name:          workerName,
		Durable:       workerName,
		MaxDeliver:    5,
		BackOff:       []time.Duration{5 * time.Second, 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("Error creando el consumidor: %v", err)
	}
	
	// Suscribirse al stream y procesar mensajes
	sub, err := consumer.Consume(func(msg jetstream.Msg) {
		// Procesar el mensaje
		log.Printf("[%s] Procesando mensaje: %s\n", workerName, string(msg.Data()))
        data := strings.Split(string(msg.Data()), "|")
        subject := msg.Subject()
		requestId := strings.Split(subject, ".")

		username, functionName, parameter := data[0], data[1], data[2]
		log.Printf("Procesando activación: Usuario=%s, Función=%s, Parámetro=%s\n", username, functionName, parameter)

		// Procesar la función
		result, err := processFunction(functionName, parameter)
		if err != nil {
			log.Printf("Error ejecutando función para %s: %s\n", username, err.Error())
			js.Publish(ctx,"results."+requestId[1], []byte(fmt.Sprintf("Error para %s: %s", username, err.Error())))
			return
		}
		msg.Ack()  // Confirmar mensaje procesado
        
        // Publicar el resultado en "results"
		js.Publish(ctx,"results."+requestId[1], []byte(result))
		log.Printf("Resultado enviado para %s\n", username)
        
	})
	if err != nil {
		log.Fatalf("Error suscribiéndose al stream: %v", err)
	}
    
    // Manejar cierre limpio
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	sub.Stop()
}

// createStreams se asegura de que los streams "activations" y "results" estén configurados
func createStreams(js nats.JetStreamContext) {
	// Configurar el stream "activations"
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "activations",         // Nombre del stream
		Subjects: []string{"activations.*"}, // Temas relacionados
        Retention: nats.WorkQueuePolicy, // Elimina mensajes después de ser procesados
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
