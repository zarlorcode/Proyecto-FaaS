package database

import (
    "github.com/nats-io/nats.go"
    "time"
    "log"
)

func ConnectNATS() (nats.KeyValue,nats.JetStreamContext, error) {
    
    // Esperar unos segundos para asegurar que NATS esté disponible
    time.Sleep(5*time.Second)
    
    // Conectar a NATS
    nc, err := nats.Connect("nats://nats:4222")
    if err != nil {
        return nil, nil, err
    }
    
    // Inicializar JetStream
    js, err := nc.JetStream()
    if err != nil {
        return nil, nil, err
    }   
    // Configurar el stream "activations"
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
