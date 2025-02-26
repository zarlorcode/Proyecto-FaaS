package services

import (
	"github.com/nats-io/nats.go"
    "fmt"
    "bytes"
	"encoding/json"
	"net/http"
)

type UserService struct {
	KV nats.KeyValue
}

// Nueva instancia del servicio
func NewUserService(kv nats.KeyValue) *UserService {
	return &UserService{KV: kv}
}

// Registrar un usuario
func (s *UserService) RegisterUser(username, password string) error {
    // Define the API endpoint
	url := "http://apisix:9180/apisix/admin/consumers"
    adminKey := "edd1c9f034335f136f87ad84b625c8f1"
	// Create the request payload
	payload := map[string]interface{}{
		"username": username,
		"plugins": map[string]interface{}{
			"basic-auth": map[string]string{
				"username": username,
				"password": password,
			},
		},
	}

	// Serialize the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", adminKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("El usuario ya existe")
	}

	fmt.Println("Consumer created successfully")
	return nil
}

