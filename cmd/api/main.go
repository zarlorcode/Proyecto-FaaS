package main

import (
    "log"
    "faas-api/internal/database"
    "faas-api/internal/handlers"
    "faas-api/internal/services"
    "github.com/gin-gonic/gin"
)

func main() {
    // Configurar conexión con NATS
    kv, err := database.ConnectNATS()
    if err != nil {
        log.Fatalf("Error conectando a NATS: %v", err)
    }

    // Inicializar servicio y controlador
    userService := services.NewUserService(kv)
    userHandler := handlers.NewUserHandler(userService)
    functionService := services.NewFunctionService(kv)
	functionHandler := handlers.NewFunctionHandler(functionService)

    // Configurar router
    router := gin.Default()
    
    // Rutas públicas
	router.POST("/register", userHandler.RegisterUser)
	router.POST("/login", userHandler.LoginUser)
    
    // Rutas protegidas
	auth := router.Group("/")
	auth.Use(handlers.AuthMiddleware())
	auth.POST("/functions/register", functionHandler.RegisterFunction)
    auth.POST("/functions/deregister", functionHandler.DeregisterFunction)
    auth.POST("/functions/activate", functionHandler.ActivateFunction)

	// Iniciar servidor
	log.Println("Servidor iniciado en http://localhost:8080")
	router.Run(":8080")
}
