package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"billing-service/internal/repository"
)

type BillingService struct {
	repo              repository.BillingRepository
	appointmentSvcURL string
	userSvcURL        string
	notifURL          string
	client            *http.Client
}

func NewBillingService(r repository.BillingRepository, apptURL, userURL, notifURL string) *BillingService {
	return &BillingService{
		repo:              r,
		appointmentSvcURL: apptURL,
		userSvcURL:        userURL,
		notifURL:          notifURL,
		client:            &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *BillingService) CreateBill(appointmentID string) (*repository.Bill, error) {
	url := fmt.Sprintf("%s/v1/appointments/%s", s.appointmentSvcURL, appointmentID)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("appointment fetch failed: %s", resp.Status)
	}
	var appt struct {
		AppointmentID string `json:"appointment_id"`
		PatientID     string `json:"patient_id"`
		DoctorID      string `json:"doctor_id"`
		Status        string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&appt); err != nil {
		return nil, err
	}
	if appt.Status != "COMPLETED" {
		return nil, errors.New("appointment not completed yet")
	}
	b := &repository.Bill{
		AppointmentID: appointmentID,
		PatientID:     appt.PatientID,
		DoctorID:      appt.DoctorID,
		Amount:        500.0,
	}
	if err := s.repo.Create(b); err != nil {
		return nil, err
	}
	// notify patient asynchronously
	go func() {
		payload := map[string]string{"to": appt.PatientID, "message": fmt.Sprintf("Bill %s created for appointment %s", b.BillID, appointmentID)}
		body, _ := json.Marshal(payload)
		s.client.Post(s.notifURL+"/v1/notify", "application/json", bytes.NewReader(body))
	}()
	return b, nil
}

func (s *BillingService) GetBill(id string) (*repository.Bill, error) {
	return s.repo.GetByID(id)
}

func (s *BillingService) ListBills() ([]*repository.Bill, error) {
	return s.repo.List()
}

func (s *BillingService) PayBill(id string) error {
	b, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if b == nil {
		return errors.New("bill not found")
	}
	if b.Status == "PAID" {
		return errors.New("already paid")
	}
	paidAt := time.Now().UTC()
	if err := s.repo.MarkPaid(id, paidAt); err != nil {
		return err
	}
	go func() {
		payload := map[string]string{"to": b.PatientID, "message": fmt.Sprintf("Bill %s paid", id)}
		body, _ := json.Marshal(payload)
		s.client.Post(s.notifURL+"/v1/notify", "application/json", bytes.NewReader(body))
	}()
	return nil
}
