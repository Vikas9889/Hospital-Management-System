package main

import (
    "log"

    "appointment-service/internal/config"
    "appointment-service/internal/handler"
    "appointment-service/internal/repository"
    "appointment-service/internal/service"

    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()

    db := repository.ConnectDB(cfg.DatabaseURL)
    repo := repository.NewAppointmentRepository(db)
    svc := service.NewAppointmentService(repo, cfg.UserServiceURL)
    h := handler.NewAppointmentHandler(svc)

    r := gin.Default()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    api := r.Group("/v1")
    {
        api.POST("/appointments", h.CreateAppointment)
        api.GET("/appointments", h.ListAppointments)
        api.GET("/appointments/:id", h.GetAppointment)
        api.POST("/appointments/:id/reschedule", h.RescheduleAppointment)
        api.DELETE("/appointments/:id", h.CancelAppointment)
        api.POST("/appointments/:id/complete", h.CompleteAppointment)
    }

    log.Printf("Appointment Service running on port %s", cfg.Port)
    r.Run(":" + cfg.Port)
}
