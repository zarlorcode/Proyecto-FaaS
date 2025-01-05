package database

import (
    "github.com/nats-io/nats.go"
    "time"
    "log"
    "errors"
)

func ConnectNATS() (nats.KeyValue,nats.JetStreamContext, error) {
    
    //Logica de reintento para esperar a que nats esté disponible
    var nc *nats.Conn
    var err error

    maxRetries := 10              // Número máximo de reintentos
    retryInterval := 2 * time.Second // Intervalo entre reintentos

    for i := 1; i <= maxRetries; i++ {
        nc, err = nats.Connect("nats://nats:4222")
        if err == nil {
            break // Conexión exitosa
        }

        log.Printf("Intento %d de %d: no se pudo conectar a NATS. Error: %v", i, maxRetries, err)
        time.Sleep(retryInterval) // Esperar antes del próximo intento
    }

    if err != nil {
        return nil, nil, errors.New("no se pudo conectar a NATS después de múltiples intentos")
    }

    log.Println("Conexión exitosa a NATS")
    //FIN DE LA CONEXIÓN CON NATS
    
    // Inicializar JetStream
    js, err := nc.JetStream()
    if err != nil {
        return nil, nil, err
    }   
    // Configurar el stream "activations" (Work Queue Stream)
    _, err = js.AddStream(&nats.StreamConfig{
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
 
    // Configurar el Key-Value "users"
    kv, err := js.KeyValue("users")
    if err != nil {
        kv, err = js.CreateKeyValue(&nats.KeyValueConfig{
            Bucket: "users",
        })
        if err != nil {
            return nil, nil, err
        }
    }
    
    // Devolver Key-Value y JetStream
    return kv,js, nil
}
