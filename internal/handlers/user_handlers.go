package handlers

import (
	"faas-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

// Endpoint para registrar usuarios
func (h *UserHandler) RegisterUser(c *gin.Context) {
    // Obtén las credenciales del encabezado de Basic Auth
    username, password, ok := c.Request.BasicAuth()
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Credenciales no proporcionadas o inválidas"})
        return
    }

    // Llama al servicio para registrar al usuario
    if err := h.Service.RegisterUser(username, password); err != nil {
        c.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Usuario registrado exitosamente"})
}


