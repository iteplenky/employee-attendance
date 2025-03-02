-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS telegram_bot;

CREATE TABLE IF NOT EXISTS telegram_bot.users (
    id SERIAL PRIMARY KEY,
    tg_id BIGINT UNIQUE NOT NULL,
    iin TEXT UNIQUE NOT NULL,
    notifications_enabled BOOLEAN DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS telegram_bot CASCADE;
-- +goose StatementEnd