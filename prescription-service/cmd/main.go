package main

import (
	"log"
	"prescription-service/internal/config"
	"prescription-service/internal/handler"
	"prescription-service/internal/repository"
	"prescription-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db := repository.ConnectDB(cfg.DatabaseURL)
	repo := repository.NewPrescriptionRepository(db)
	svc := service.NewPrescriptionService(repo)
	h := handler.NewPrescriptionHandler(svc)

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	api := r.Group("/v1")
	{
		api.POST("/prescriptions", h.CreatePrescription)
		api.GET("/prescriptions", h.ListPrescriptions)
		api.GET("/prescriptions/:id", h.GetPrescription)
	}

	log.Printf("Prescription Service running on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
