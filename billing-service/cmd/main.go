package main

import (
    "log"
    "billing-service/internal/config"
    "billing-service/internal/handler"
    "billing-service/internal/repository"
    "billing-service/internal/service"

    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()
    db := repository.ConnectDB(cfg.DatabaseURL)
    repo := repository.NewBillingRepository(db)
    svc := service.NewBillingService(repo, cfg.AppointmentServiceURL, cfg.UserServiceURL, cfg.NotificationURL)
    h := handler.NewBillingHandler(svc)

    r := gin.Default()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    api := r.Group("/v1")
    {
        api.POST("/bills", h.CreateBill)
        api.GET("/bills", h.ListBills)
        api.GET("/bills/:id", h.GetBill)
        api.POST("/bills/:id/pay", h.PayBill)
    }

    log.Printf("Billing Service running on port %s", cfg.Port)
    r.Run(":" + cfg.Port)
}
