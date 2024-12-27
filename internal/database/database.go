package database

import (
    "github.com/nats-io/nats.go"
)

func ConnectNATS() (nats.KeyValue, error) {
    nc, err := nats.Connect(nats.DefaultURL)
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
