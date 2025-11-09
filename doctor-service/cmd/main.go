package main

import (
	"doctor-service/internal/config"
	"doctor-service/internal/handler"
	"doctor-service/internal/repository"
	"doctor-service/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db := repository.ConnectDB(cfg.DatabaseURL)
	repo := repository.NewDoctorRepository(db)
	svc := service.NewDoctorService(repo)
	h := handler.NewDoctorHandler(svc)

	r := gin.Default()
    // Add this health endpoint for readiness/liveness probes
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
	api := r.Group("/v1")
	{
		api.POST("/doctors", h.CreateDoctor)
		api.GET("/doctors", h.ListDoctors)
		api.GET("/doctors/:id", h.GetDoctor)
		api.PUT("/doctors/:id", h.UpdateDoctor)
		api.DELETE("/doctors/:id", h.DeleteDoctor)
	}

	log.Printf("Doctor Service running on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
