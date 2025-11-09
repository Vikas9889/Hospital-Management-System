package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	r := gin.Default()
    // Add this health endpoint for readiness/liveness probes
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
	// Group all APIs under /v1
	api := r.Group("/v1")
	{
		api.GET("/notify/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Notification service is live 🚀",
			})
		})

		api.POST("/notify", func(c *gin.Context) {
			var payload struct {
				Email   string `json:"email"`
				Subject string `json:"subject"`
				Body    string `json:"body"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			// Mock sending notification
			log.Printf("📧 Notification sent to %s | %s: %s", payload.Email, payload.Subject, payload.Body)
			c.JSON(http.StatusOK, gin.H{"status": "sent", "to": payload.Email})
		})
	}

	log.Printf("🚀 Notification Service running on port %s", port)
	r.Run(":" + port)
}
