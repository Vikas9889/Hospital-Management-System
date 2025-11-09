package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Doctor struct {
	DoctorID string `json:"doctor_id"`
	Name     string `json:"name"`
}

type Patient struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
}

// ValidateDoctor checks if doctor exists in doctor-service
func ValidateDoctor(baseURL, doctorID string) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v1/doctors/%s", baseURL, doctorID))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("doctor-service returned status %d", resp.StatusCode)
	}

	var d Doctor
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return false, err
	}
	return d.DoctorID != "", nil
}

// ValidatePatient checks if patient exists in user-service
func ValidatePatient(baseURL, patientID string) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v1/users/%s", baseURL, patientID))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("user-service returned status %d", resp.StatusCode)
	}

	var p Patient
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return false, err
	}
	return p.Role == "PATIENT", nil
}
