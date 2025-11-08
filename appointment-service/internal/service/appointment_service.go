package service

import (
	"errors"
	"net/http"
	"time"

	"appointment-service/internal/repository"
)

type AppointmentService struct {
	repo       repository.AppointmentRepository
	userSvcURL string
	client     *http.Client
}

func NewAppointmentService(r repository.AppointmentRepository, userSvcURL string) *AppointmentService {
	return &AppointmentService{repo: r, userSvcURL: userSvcURL, client: &http.Client{Timeout: 5 * time.Second}}
}

func (s *AppointmentService) CreateAppointment(a *repository.Appointment) error {
	// Basic creation - status set to SCHEDULED
	a.Status = "SCHEDULED"
	return s.repo.CreateAppointment(a)
}

func (s *AppointmentService) ListAppointments() ([]repository.Appointment, error) {
	return s.repo.ListAppointments()
}

func (s *AppointmentService) GetAppointment(id string) (*repository.Appointment, error) {
	return s.repo.GetAppointment(id)
}

func (s *AppointmentService) Reschedule(id string, start, end time.Time) error {
	a, err := s.repo.GetAppointment(id)
	if err != nil {
		return err
	}
	if a == nil {
		return errors.New("appointment not found")
	}
	if a.Status == "CANCELLED" || a.Status == "COMPLETED" {
		return errors.New("cannot reschedule appointment in its current state")
	}
	a.StartTime = start
	a.EndTime = end
	a.RescheduleCount += 1
	return s.repo.UpdateAppointment(a)
}

func (s *AppointmentService) Cancel(id string) error {
	a, err := s.repo.GetAppointment(id)
	if err != nil {
		return err
	}
	if a == nil {
		return errors.New("appointment not found")
	}
	if a.Status == "CANCELLED" {
		return errors.New("appointment already cancelled")
	}
	a.Status = "CANCELLED"
	return s.repo.UpdateAppointment(a)
}

func (s *AppointmentService) Complete(id string) error {
	a, err := s.repo.GetAppointment(id)
	if err != nil {
		return err
	}
	if a == nil {
		return errors.New("appointment not found")
	}
	if a.Status == "COMPLETED" {
		return errors.New("appointment already completed")
	}
	a.Status = "COMPLETED"
	return s.repo.UpdateAppointment(a)
}
