package handler

import (
	"net/http"
	"prescription-service/internal/service"

	"github.com/gin-gonic/gin"
)

type PrescriptionHandler struct{ svc *service.PrescriptionService }

func NewPrescriptionHandler(s *service.PrescriptionService) *PrescriptionHandler {
	return &PrescriptionHandler{svc: s}
}

func (h *PrescriptionHandler) CreatePrescription(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "create prescription - scaffold"})
}

func (h *PrescriptionHandler) ListPrescriptions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "list prescriptions - scaffold"})
}

func (h *PrescriptionHandler) GetPrescription(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get prescription - scaffold"})
}
