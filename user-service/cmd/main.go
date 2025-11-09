package main

import (
    "log"
    "user-service/internal/config"
    "user-service/internal/handler"
    "user-service/internal/repository"
    "user-service/internal/service"

    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()
    db := repository.ConnectDB(cfg.DatabaseURL)
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

    r := gin.Default()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    // Add this health endpoint for readiness/liveness probes
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    api := r.Group("/v1")
    {
        api.GET("/users", userHandler.GetAll)
        api.GET("/users/:id", userHandler.GetByID)
        api.POST("/users", userHandler.Create)
        api.PUT("/users/:id", userHandler.Update)
        api.DELETE("/users/:id", userHandler.Delete)
    }

    log.Printf("User Service running on port %s", cfg.Port)
    r.Run(":" + cfg.Port)
}
