package main

import (
    "log"
    "notification-service/internal/config"
    "notification-service/internal/handler"

    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()
    r := gin.Default()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    h := handler.NewNotificationHandler()

    api := r.Group("/v1")
    {
        api.POST("/notify", h.Notify)
        api.POST("/notify/test", h.Test)
    }

    log.Printf("Notification Service running on port %s", cfg.Port)
    r.Run(":" + cfg.Port)
}
