package service

import "prescription-service/internal/repository"

type PrescriptionService struct {
	repo *repository.PrescriptionRepository
}

func NewPrescriptionService(r *repository.PrescriptionRepository) *PrescriptionService {
	return &PrescriptionService{repo: r}
}

// Add business logic: validate appointment exists, create/read prescriptions
