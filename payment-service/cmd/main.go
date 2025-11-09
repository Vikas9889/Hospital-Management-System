package main

import (
	"log"
	"payment-service/internal/config"
	"payment-service/internal/handler"
	"payment-service/internal/repository"
	"payment-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db := repository.ConnectDB(cfg.DatabaseURL)
	repo := repository.NewPaymentRepository(db)
	svc := service.NewPaymentService(repo)
	h := handler.NewPaymentHandler(svc)

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	api := r.Group("/v1")
	{
		api.POST("/payments", h.CreatePayment)
		api.POST("/payments/:id/refund", h.RefundPayment)
	}

	log.Printf("Payment Service running on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
