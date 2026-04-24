CREATE TYPE status_enum AS ENUM ('open', 'closed', 'in_active');

CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    text VARCHAR(255) NOT NULL,
    status status_enum NOT NULL DEFAULT 'open',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT NOT NULL REFERENCES users(id),
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX todos_status_idx ON todos (status);