CREATE SCHEMA testing;

CREATE TABLE testing.todos (
    id INT GENERATED ALWAYS AS IDENTITY,
    task TEXT NOT NULL,
    complete BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP
);
