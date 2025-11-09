package handler

import (
	"doctor-service/internal/repository"
	"doctor-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	Service *service.DoctorService
}

func NewDoctorHandler(s *service.DoctorService) *DoctorHandler {
	return &DoctorHandler{Service: s}
}

func (h *DoctorHandler) CreateDoctor(c *gin.Context) {
	var d repository.Doctor
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.CreateDoctor(&d); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *DoctorHandler) ListDoctors(c *gin.Context) {
	dept := c.Query("department")
	doctors, err := h.Service.ListDoctors(dept)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doctors)
}

func (h *DoctorHandler) GetDoctor(c *gin.Context) {
	id := c.Param("id")
	d, err := h.Service.GetDoctor(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if d == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "doctor not found"})
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *DoctorHandler) UpdateDoctor(c *gin.Context) {
	id := c.Param("id")
	var d repository.Doctor
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.UpdateDoctor(id, &d); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *DoctorHandler) DeleteDoctor(c *gin.Context) {
	id := c.Param("id")
	if err := h.Service.DeleteDoctor(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
