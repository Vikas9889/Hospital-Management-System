package handler

import (
    "net/http"
    "fmt"

    "github.com/gin-gonic/gin"
)

type NotificationHandler struct{}

func NewNotificationHandler() *NotificationHandler { return &NotificationHandler{} }

type notifyReq struct {
    To      string `json:"to" binding:"required"`
    Message string `json:"message" binding:"required"`
}

func (h *NotificationHandler) Notify(c *gin.Context) {
    var req notifyReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    fmt.Printf("[NOTIFY] to=%s message=%s\n", req.To, req.Message)
    c.JSON(http.StatusOK, gin.H{"status": "sent"})
}

func (h *NotificationHandler) Test(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
