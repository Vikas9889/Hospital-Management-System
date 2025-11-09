-- payments table scaffold
CREATE TABLE IF NOT EXISTS payments (
  id SERIAL PRIMARY KEY,
  external_id TEXT,
  amount NUMERIC,
  status TEXT,
  created_at TIMESTAMP DEFAULT now()
);
