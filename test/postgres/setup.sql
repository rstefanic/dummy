CREATE SCHEMA testing;

CREATE TABLE testing.todos (
    id SERIAL PRIMARY KEY,
    task TEXT NOT NULL,
    complete BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP
);
