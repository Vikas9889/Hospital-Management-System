package repository

import (
	"database/sql"
	"time"
)

type Doctor struct {
	DoctorID        string    `json:"doctor_id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	Department      string    `json:"department"`
	ExperienceYears int       `json:"experience_years"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type DoctorRepository struct {
	DB *sql.DB
}

func NewDoctorRepository(db *sql.DB) *DoctorRepository {
	return &DoctorRepository{DB: db}
}

func (r *DoctorRepository) CreateDoctor(d *Doctor) error {
	return r.DB.QueryRow(`
        INSERT INTO doctors (name, email, phone, department, experience_years)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING doctor_id, created_at, updated_at
    `, d.Name, d.Email, d.Phone, d.Department, d.ExperienceYears).Scan(&d.DoctorID, &d.CreatedAt, &d.UpdatedAt)
}

func (r *DoctorRepository) ListDoctors(dept string) ([]Doctor, error) {
	query := "SELECT doctor_id, name, email, phone, department, experience_years, created_at, updated_at FROM doctors"
	var rows *sql.Rows
	var err error

	if dept != "" {
		query += " WHERE department ILIKE $1"
		rows, err = r.DB.Query(query, "%"+dept+"%")
	} else {
		rows, err = r.DB.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []Doctor
	for rows.Next() {
		var d Doctor
		if err := rows.Scan(&d.DoctorID, &d.Name, &d.Email, &d.Phone, &d.Department, &d.ExperienceYears, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		doctors = append(doctors, d)
	}
	return doctors, nil
}

func (r *DoctorRepository) GetDoctor(id string) (*Doctor, error) {
	var d Doctor
	err := r.DB.QueryRow(`
        SELECT doctor_id, name, email, phone, department, experience_years, created_at, updated_at
        FROM doctors WHERE doctor_id=$1
    `, id).Scan(&d.DoctorID, &d.Name, &d.Email, &d.Phone, &d.Department, &d.ExperienceYears, &d.CreatedAt, &d.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &d, err
}

func (r *DoctorRepository) UpdateDoctor(id string, d *Doctor) error {
	_, err := r.DB.Exec(`
        UPDATE doctors
        SET name=$1, email=$2, phone=$3, department=$4, experience_years=$5, updated_at=NOW()
        WHERE doctor_id=$6
    `, d.Name, d.Email, d.Phone, d.Department, d.ExperienceYears, id)
	return err
}

func (r *DoctorRepository) DeleteDoctor(id string) error {
	_, err := r.DB.Exec(`DELETE FROM doctors WHERE doctor_id=$1`, id)
	return err
}
