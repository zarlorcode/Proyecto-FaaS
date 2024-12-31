package handlers

import (
	"faas-api/internal/config"
	"faas-api/internal/models"
	"faas-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
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

// Endpoint para loguearse con usuarios
func (h *UserHandler) LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Datos inválidos"})
		return
	}

	err := h.Service.AuthenticateUser(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": err.Error()})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "No se pudo generar el token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "token": tokenString})
}

