package services

import (
	"errors"
	"github.com/nats-io/nats.go"
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
	// Verificar si el usuario ya existe
	_, err := s.KV.Get(username) // Ignoramos el primer valor (uint64)
	if err == nil {
		return errors.New("usuario ya registrado")
	}

	// Guardar el usuario
	_, err = s.KV.Put(username, []byte(password)) // Ignoramos el primer valor (uint64)
	if err != nil {
		return err
	}

	return nil
}

// Método para autenticar usuarios
func (s *UserService) AuthenticateUser(username, password string) error {
	// Verificar si el usuario existe
	value, err := s.KV.Get(username)
	if err != nil {
		return errors.New("usuario no encontrado")
	}

	// Verificar contraseña
	if string(value.Value()) != password {
		return errors.New("contraseña incorrecta")
	}

	return nil
}
