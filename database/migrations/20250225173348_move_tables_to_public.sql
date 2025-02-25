-- +goose Up
-- +goose StatementBegin
ALTER TABLE users SET SCHEMA public;
ALTER TABLE schedules SET SCHEMA public;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Откат не нужен, т.к public схема по умолчанию
-- +goose StatementEnd