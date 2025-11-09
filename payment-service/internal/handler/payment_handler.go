package handler

import (
	"net/http"
	"payment-service/internal/service"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct{ svc *service.PaymentService }

func NewPaymentHandler(s *service.PaymentService) *PaymentHandler { return &PaymentHandler{svc: s} }

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "create payment - scaffold"})
}

func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "refund payment - scaffold"})
}
