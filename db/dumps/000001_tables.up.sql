BEGIN;

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    amount BIGINT NOT NULL CHECK (amount>=0),
    type TEXT NOT NULL CHECK (type IN ('income', 'expense')),
    category TEXT NOT NULL,
    event_date TIMESTAMP NOT NULL

);

COMMIT;
