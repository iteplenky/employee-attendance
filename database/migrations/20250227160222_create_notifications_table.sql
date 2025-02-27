-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN notifications_enabled BOOLEAN DEFAULT FALSE;
DROP TABLE IF EXISTS schedules;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS schedules (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL
);
-- +goose StatementEnd
