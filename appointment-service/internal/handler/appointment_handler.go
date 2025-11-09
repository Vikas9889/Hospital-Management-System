package handler

import (
	"appointment-service/internal/config"
	"appointment-service/internal/repository"
	"appointment-service/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	service *service.AppointmentService
	config  *config.Config
}

func NewAppointmentHandler(s *service.AppointmentService, cfg *config.Config) *AppointmentHandler {
	return &AppointmentHandler{
		service: s,
		config:  cfg,
	}
}

type createReq struct {
	PatientID string `json:"patient_id" binding:"required"`
	DoctorID  string `json:"doctor_id" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time"})
		return
	}
	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time"})
		return
	}

	// ✅ Validate Patient via User Service
	validPatient, err := service.ValidatePatient(h.config.UserServiceURL, req.PatientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user-service unreachable"})
		return
	}
	if !validPatient {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	// ✅ Validate Doctor via Doctor Service
	validDoctor, err := service.ValidateDoctor(h.config.DoctorServiceURL, req.DoctorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "doctor-service unreachable"})
		return
	}
	if !validDoctor {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor_id"})
		return
	}

	a := &repository.Appointment{
		PatientID: req.PatientID,
		DoctorID:  req.DoctorID,
		StartTime: start.UTC(),
		EndTime:   end.UTC(),
	}

	if err := h.service.CreateAppointment(a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, a)
}

func (h *AppointmentHandler) ListAppointments(c *gin.Context) {
	res, err := h.service.ListAppointments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *AppointmentHandler) GetAppointment(c *gin.Context) {
	id := c.Param("id")
	a, err := h.service.GetAppointment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if a == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}

type rescheduleReq struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

func (h *AppointmentHandler) RescheduleAppointment(c *gin.Context) {
	id := c.Param("id")
	var req rescheduleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time"})
		return
	}
	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time"})
		return
	}
	if err := h.service.Reschedule(id, start.UTC(), end.UTC()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "rescheduled"})
}

func (h *AppointmentHandler) CancelAppointment(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Cancel(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}

func (h *AppointmentHandler) CompleteAppointment(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Complete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "completed"})
}
