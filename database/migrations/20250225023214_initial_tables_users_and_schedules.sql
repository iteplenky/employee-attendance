-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    tg_id BIGINT UNIQUE NOT NULL,
    iin VARCHAR(12) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS schedules (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
