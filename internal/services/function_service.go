package services

import (
    //"bytes"
	//"context"
	"errors"
	"fmt"
	//"os/exec"
	"time"

	"github.com/nats-io/nats.go"
    "github.com/google/uuid" 
)

// Límite de concurrencia
var maxConcurrentActivations = 5

// Canal para manejar concurrencia
var sem = make(chan struct{}, maxConcurrentActivations) 

// FunctionService es el servicio que maneja las funciones
type FunctionService struct {
	KV        nats.KeyValue
	JetStream nats.JetStreamContext
}

// NewFunctionService crea una nueva instancia de FunctionService
func NewFunctionService(kv nats.KeyValue, js nats.JetStreamContext) *FunctionService {
	return &FunctionService{
		KV:        kv,
		JetStream: js,
	}
}

// Registrar una función
func (s *FunctionService) RegisterFunction(username, functionName, dockerImage string) error {
    // Usar un formato de clave válido
    key := fmt.Sprintf("%s/%s", username, functionName)
    fmt.Println("El username ", username ," ha REGISTRADO la función ", functionName )
    // Verificar si la función ya existe
    _, err := s.KV.Get(key)
    if err == nil {
        return errors.New("la función ya está registrada")
    }

    // Registrar la función
    _, err = s.KV.PutString(key,dockerImage)
    if err != nil {
        return err
    }

    return nil
}

func (s *FunctionService) DeleteFunction(username, functionName string) error {
	// Construir la clave de la función
	key := fmt.Sprintf("%s/%s", username, functionName)
    fmt.Println("El username ", username ," ha ELIMINADO la función ", functionName )
    
	// Intentar obtener la función
	value, err := s.KV.Get(key)
	if err != nil {
		if err == nats.ErrKeyNotFound {
			return errors.New("función no encontrada")
		}
		return fmt.Errorf("error al buscar la función: %s", err)
	}

	// Verificar que la clave pertenece al usuario correcto
	if string(value.Value()) == "" {
		return errors.New("la función no pertenece a este usuario")
	}

	// Intentar eliminar la clave
	err = s.KV.Delete(key)
	if err != nil {
		return fmt.Errorf("error eliminando la función: %s", err)
	}

	return nil
}


// Activar una función con control de concurrencia
func (s *FunctionService) ActivateFunction(username, functionName, parameter string) (string, error) {
    // Construir la clave de la función
    key := fmt.Sprintf("%s/%s", username, functionName)
    
    // Verificar si la función existe
    value, err := s.KV.Get(key)
    if err != nil {
        if err == nats.ErrKeyNotFound {
            return "", errors.New("función no encontrada")
        }
        return "", fmt.Errorf("error al buscar la función: %s", err)
    }

    // Verificar que la función pertenece al usuario
    if string(value.Value()) == "" {
        return "", errors.New("la función no pertenece a este usuario")
    }

    // Generar un ID único para esta petición
    requestID := uuid.NewString()

    // Construir el mensaje
    message := fmt.Sprintf("%s|%s|%s", username, value.Value(), parameter)
    
    // Publicar el mensaje en el stream "activations"
    _, err = s.JetStream.Publish("activations."+requestID, []byte(message))
    if err != nil {
        return "", fmt.Errorf("error enviando activación: %s", err)
    }
    
    // Configurar un suscriptor para escuchar el resultado en "results.<username>"
    subject := "results." + requestID
    sub, err := s.JetStream.SubscribeSync(subject)
    if err != nil {
        return "", fmt.Errorf("error suscribiéndose al stream de resultados: %s", err)
    }
    defer sub.Unsubscribe() // Asegurarse de cerrar la suscripción al final
    
    // Esperar el mensaje de resultado
    timeout := 40 * time.Second // Tiempo máximo de espera para el resultado
    msg, err := sub.NextMsg(timeout)
    if err != nil {
        return "", fmt.Errorf("error esperando el resultado: %s", err)
    }

    // Procesar el mensaje recibido
    result := string(msg.Data) // Convertir los datos a string para el retorno
    return result, nil
}

