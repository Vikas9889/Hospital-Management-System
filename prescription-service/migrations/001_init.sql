-- prescriptions table scaffold
CREATE TABLE IF NOT EXISTS prescriptions (
  id SERIAL PRIMARY KEY,
  appointment_id INT NOT NULL,
  content TEXT,
  created_at TIMESTAMP DEFAULT now()
);
