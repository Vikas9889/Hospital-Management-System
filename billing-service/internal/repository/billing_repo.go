package repository

import (
    "database/sql"
    "errors"
    "time"

    "github.com/google/uuid"
)

type Bill struct {
    BillID        string     `json:"bill_id"`
    AppointmentID string     `json:"appointment_id"`
    PatientID     string     `json:"patient_id"`
    Amount        float64    `json:"amount"`
    Status        string     `json:"status"`
    CreatedAt     time.Time  `json:"created_at"`
    PaidAt        *time.Time `json:"paid_at,omitempty"`
}

type BillingRepository interface {
    Create(b *Bill) error
    GetByID(id string) (*Bill, error)
    List() ([]*Bill, error)
    MarkPaid(id string, paidAt time.Time) error
}

type billingRepo struct {
    db *sql.DB
}

func NewBillingRepository(db *sql.DB) BillingRepository {
    return &billingRepo{db: db}
}

func (r *billingRepo) Create(b *Bill) error {
    id := uuid.New().String()
    query := `INSERT INTO bills(bill_id,appointment_id,patient_id,amount,status,created_at) VALUES($1,$2,$3,$4,$5,NOW())`
    _, err := r.db.Exec(query, id, b.AppointmentID, b.PatientID, b.Amount, "OPEN")
    if err != nil {
        return err
    }
    b.BillID = id
    return nil
}

func (r *billingRepo) GetByID(id string) (*Bill, error) {
    b := &Bill{}
    row := r.db.QueryRow(`SELECT bill_id,appointment_id,patient_id,amount,status,created_at,paid_at FROM bills WHERE bill_id=$1`, id)
    var paidAt sql.NullTime
    if err := row.Scan(&b.BillID, &b.AppointmentID, &b.PatientID, &b.Amount, &b.Status, &b.CreatedAt, &paidAt); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    if paidAt.Valid {
        b.PaidAt = &paidAt.Time
    }
    return b, nil
}

func (r *billingRepo) List() ([]*Bill, error) {
    rows, err := r.db.Query(`SELECT bill_id,appointment_id,patient_id,amount,status,created_at,paid_at FROM bills ORDER BY created_at DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var res []*Bill
    for rows.Next() {
        b := &Bill{}
        var paidAt sql.NullTime
        if err := rows.Scan(&b.BillID, &b.AppointmentID, &b.PatientID, &b.Amount, &b.Status, &b.CreatedAt, &paidAt); err != nil {
            return nil, err
        }
        if paidAt.Valid {
            b.PaidAt = &paidAt.Time
        }
        res = append(res, b)
    }
    return res, nil
}

func (r *billingRepo) MarkPaid(id string, paidAt time.Time) error {
    res, err := r.db.Exec(`UPDATE bills SET status='PAID', paid_at=$1 WHERE bill_id=$2`, paidAt, id)
    if err != nil {
        return err
    }
    cnt, _ := res.RowsAffected()
    if cnt == 0 {
        return errors.New("no rows updated")
    }
    return nil
}
