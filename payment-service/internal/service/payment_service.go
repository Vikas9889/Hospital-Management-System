package service

import "payment-service/internal/repository"

type PaymentService struct{ repo *repository.PaymentRepository }

func NewPaymentService(r *repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: r}
}

// Implement idempotent charge/refund logic here
