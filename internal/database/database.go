package database

import (
    "github.com/nats-io/nats.go"
    "time"
)

func ConnectNATS() (nats.KeyValue, error) {
    time.Sleep(5*time.Second)
    nc, err := nats.Connect("nats://nats:4222")
    if err != nil {
        return nil, err
    }

    js, err := nc.JetStream()
    if err != nil {
        return nil, err
    }

    kv, err := js.KeyValue("users")
    if err != nil {
        kv, err = js.CreateKeyValue(&nats.KeyValueConfig{
            Bucket: "users",
        })
        if err != nil {
            return nil, err
        }
    }

    return kv, nil
}
