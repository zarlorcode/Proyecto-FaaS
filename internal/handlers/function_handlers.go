package handlers

import (
	"net/http"
	"faas-api/internal/services"
	"github.com/gin-gonic/gin"
)

type FunctionHandler struct {
	Service *services.FunctionService
}

// Nueva instancia del controlador
func NewFunctionHandler(service *services.FunctionService) *FunctionHandler {
	return &FunctionHandler{Service: service}
}

// Endpoint para registrar funciones
func (h *FunctionHandler) RegisterFunction(c *gin.Context) {
	var req struct {
		FunctionName string `json:"functionName"`
		DockerImage  string `json:"dockerImage"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Datos inválidos"})
		return
	}

	// Obtener el usuario desde el contexto
	username, _ := c.Get("username")

	// Registrar la función
	err := h.Service.RegisterFunction(username.(string), req.FunctionName, req.DockerImage)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Función registrada exitosamente"})
}

//Endpoint para eliminar funciones
func (h *FunctionHandler) DeregisterFunction(c *gin.Context) {
	var req struct {
		FunctionName string `json:"functionName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Datos inválidos"})
		return
	}

	// Obtener el usuario desde el contexto
	username, _ := c.Get("username")

	// Intentar eliminar la función
	err := h.Service.DeleteFunction(username.(string), req.FunctionName)
	if err != nil {
		if err.Error() == "función no encontrada" || err.Error() == "la función no pertenece a este usuario" {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Función eliminada exitosamente"})
}

// Endpoint para activar funciones
func (h *FunctionHandler) ActivateFunction(c *gin.Context) {
	var req struct {
		FunctionName string `json:"functionName"`
		Parameter    string `json:"parameter"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Datos inválidos"})
		return
	}

	username, _ := c.Get("username")

	// Intentar activar la función
	result, err := h.Service.ActivateFunction(username.(string), req.FunctionName, req.Parameter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "result": result})
}

// Endpoint para activar funciones SIN CONCURRENCIA
/*func (h *FunctionHandler) ActivateFunction(c *gin.Context) {
	var req struct {
		FunctionName string `json:"functionName"`
		Parameter    string `json:"parameter"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Datos inválidos"})
		return
	}

	// Obtener el usuario desde el contexto
	username, _ := c.Get("username")

	// Activar la función
	result, err := h.Service.ActivateFunction(username.(string), req.FunctionName, req.Parameter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "result": result})
}*/
