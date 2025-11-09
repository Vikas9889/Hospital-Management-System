package service

import "doctor-service/internal/repository"

type DoctorService struct {
	Repo *repository.DoctorRepository
}

func NewDoctorService(r *repository.DoctorRepository) *DoctorService {
	return &DoctorService{Repo: r}
}

func (s *DoctorService) CreateDoctor(d *repository.Doctor) error {
	return s.Repo.CreateDoctor(d)
}

func (s *DoctorService) ListDoctors(dept string) ([]repository.Doctor, error) {
	return s.Repo.ListDoctors(dept)
}

func (s *DoctorService) GetDoctor(id string) (*repository.Doctor, error) {
	return s.Repo.GetDoctor(id)
}

func (s *DoctorService) UpdateDoctor(id string, d *repository.Doctor) error {
	return s.Repo.UpdateDoctor(id, d)
}

func (s *DoctorService) DeleteDoctor(id string) error {
	return s.Repo.DeleteDoctor(id)
}
