package handler

import (
	"billing-service/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BillingHandler struct {
	service *service.BillingService
}

func NewBillingHandler(s *service.BillingService) *BillingHandler {
	return &BillingHandler{service: s}
}

type createReq struct {
	// accept either a number or string for appointment_id
	AppointmentID json.RawMessage `json:"appointment_id" binding:"required"`
}

func (h *BillingHandler) CreateBill(c *gin.Context) {
	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// support both numeric and string appointment_id in the incoming JSON
	var apptIDStr string
	var tmpStr string
	if err := json.Unmarshal(req.AppointmentID, &tmpStr); err == nil {
		apptIDStr = tmpStr
	} else {
		var tmpInt int64
		if err := json.Unmarshal(req.AppointmentID, &tmpInt); err == nil {
			apptIDStr = strconv.FormatInt(tmpInt, 10)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "appointment_id must be a string or number"})
			return
		}
	}

	b, err := h.service.CreateBill(apptIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, b)
}

func (h *BillingHandler) ListBills(c *gin.Context) {
	res, err := h.service.ListBills()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *BillingHandler) GetBill(c *gin.Context) {
	id := c.Param("id")
	b, err := h.service.GetBill(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if b == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *BillingHandler) PayBill(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.PayBill(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "paid"})
}
