package repository

import (
    "database/sql"
    "errors"
    "time"
)

type User struct {
    UserID        string    `json:"user_id"`
    Name          string    `json:"name"`
    Email         string    `json:"email"`
    Phone         string    `json:"phone"`
    Role          string    `json:"role"`
    Department    string    `json:"department,omitempty"`
    Specialization string   `json:"specialization,omitempty"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type UserRepository interface {
    Create(u *User) error
    GetAll() ([]*User, error)
    GetByID(id string) (*User, error)
    Update(id string, u *User) error
    Delete(id string) error
}

type userRepo struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepo{db: db}
}

func (r *userRepo) Create(u *User) error {
    query := `INSERT INTO users(user_id,name,email,phone,role,department,specialization,created_at,updated_at) VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,NOW(),NOW()) RETURNING user_id,created_at,updated_at`
    row := r.db.QueryRow(query, u.Name, u.Email, u.Phone, u.Role, u.Department, u.Specialization)
    var id string
    var created, updated time.Time
    if err := row.Scan(&id, &created, &updated); err != nil {
        return err
    }
    u.UserID = id
    u.CreatedAt = created
    u.UpdatedAt = updated
    return nil
}

func (r *userRepo) GetAll() ([]*User, error) {
    rows, err := r.db.Query(`SELECT user_id,name,email,phone,role,department,specialization,created_at,updated_at FROM users ORDER BY created_at DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var res []*User
    for rows.Next() {
        u := &User{}
        if err := rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Phone, &u.Role, &u.Department, &u.Specialization, &u.CreatedAt, &u.UpdatedAt); err != nil {
            return nil, err
        }
        res = append(res, u)
    }
    return res, nil
}

func (r *userRepo) GetByID(id string) (*User, error) {
    u := &User{}
    row := r.db.QueryRow(`SELECT user_id,name,email,phone,role,department,specialization,created_at,updated_at FROM users WHERE user_id=$1`, id)
    if err := row.Scan(&u.UserID, &u.Name, &u.Email, &u.Phone, &u.Role, &u.Department, &u.Specialization, &u.CreatedAt, &u.UpdatedAt); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return u, nil
}

func (r *userRepo) Update(id string, u *User) error {
    res, err := r.db.Exec(`UPDATE users SET name=$1,email=$2,phone=$3,role=$4,department=$5,specialization=$6,updated_at=NOW() WHERE user_id=$7`, u.Name, u.Email, u.Phone, u.Role, u.Department, u.Specialization, id)
    if err != nil {
        return err
    }
    cnt, _ := res.RowsAffected()
    if cnt == 0 {
        return errors.New("no rows updated")
    }
    return nil
}

func (r *userRepo) Delete(id string) error {
    res, err := r.db.Exec(`DELETE FROM users WHERE user_id=$1`, id)
    if err != nil {
        return err
    }
    cnt, _ := res.RowsAffected()
    if cnt == 0 {
        return errors.New("no rows deleted")
    }
    return nil
}
