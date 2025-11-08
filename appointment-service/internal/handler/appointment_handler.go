package handler

import (
    "net/http"
    "time"

    "appointment-service/internal/repository"
    "appointment-service/internal/service"

    "github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
    service *service.AppointmentService
}

func NewAppointmentHandler(s *service.AppointmentService) *AppointmentHandler {
    return &AppointmentHandler{service: s}
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
