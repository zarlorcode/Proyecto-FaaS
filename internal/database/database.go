package database

import (
    "github.com/nats-io/nats.go"
    "time"
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
