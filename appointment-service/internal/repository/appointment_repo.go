package repository

import (
    "database/sql"
    //"log"
    "time"

    _ "github.com/lib/pq"
)

type Appointment struct {
    ID              string    `json:"id"`
    PatientID       string    `json:"patient_id"`
    DoctorID        string    `json:"doctor_id"`
    StartTime       time.Time `json:"start_time"`
    EndTime         time.Time `json:"end_time"`
    Status          string    `json:"status"`
    RescheduleCount int       `json:"reschedule_count"`
    CreatedAt       time.Time `json:"created_at"`
}

type AppointmentRepository interface {
    CreateAppointment(a *Appointment) error
    GetAppointment(id string) (*Appointment, error)
    ListAppointments() ([]Appointment, error)
    UpdateAppointment(a *Appointment) error
}

type repository struct {
    db *sql.DB
}

// func ConnectDB(url string) *sql.DB {
//     db, err := sql.Open("postgres", url)
//     if err != nil {
//         log.Fatalf("DB connect error: %v", err)
//     }
//     if err := db.Ping(); err != nil {
//         log.Fatalf("DB ping failed: %v", err)
//     }
//     log.Println("Connected to Appointment DB")
//     return db
// }

func NewAppointmentRepository(db *sql.DB) AppointmentRepository {
    return &repository{db: db}
}

func (r *repository) CreateAppointment(a *Appointment) error {
    query := `
        INSERT INTO appointments (patient_id, doctor_id, start_time, end_time, status, reschedule_count, created_at)
        VALUES ($1, $2, $3, $4, 'SCHEDULED', 0, NOW())
        RETURNING id, created_at
    `
    return r.db.QueryRow(query, a.PatientID, a.DoctorID, a.StartTime, a.EndTime).Scan(&a.ID, &a.CreatedAt)
}

func (r *repository) GetAppointment(id string) (*Appointment, error) {
    var a Appointment
    query := `SELECT id, patient_id, doctor_id, start_time, end_time, status, reschedule_count, created_at FROM appointments WHERE id=$1`
    err := r.db.QueryRow(query, id).Scan(
        &a.ID, &a.PatientID, &a.DoctorID, &a.StartTime, &a.EndTime, &a.Status, &a.RescheduleCount, &a.CreatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &a, err
}

func (r *repository) ListAppointments() ([]Appointment, error) {
    rows, err := r.db.Query(`SELECT id, patient_id, doctor_id, start_time, end_time, status, reschedule_count, created_at FROM appointments ORDER BY created_at DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var res []Appointment
    for rows.Next() {
        var a Appointment
        if err := rows.Scan(&a.ID, &a.PatientID, &a.DoctorID, &a.StartTime, &a.EndTime, &a.Status, &a.RescheduleCount, &a.CreatedAt); err != nil {
            return nil, err
        }
        res = append(res, a)
    }
    return res, nil
}

func (r *repository) UpdateAppointment(a *Appointment) error {
    _, err := r.db.Exec(`
        UPDATE appointments SET start_time=$1, end_time=$2, status=$3, reschedule_count=$4 WHERE id=$5
    `, a.StartTime, a.EndTime, a.Status, a.RescheduleCount, a.ID)
    return err
}
