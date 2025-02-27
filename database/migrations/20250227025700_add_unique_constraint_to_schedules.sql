-- +goose Up
-- +goose StatementBegin
ALTER TABLE schedules ADD CONSTRAINT schedules_user_id_unique UNIQUE (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE schedules DROP CONSTRAINT schedules_user_id_unique;
-- +goose StatementEnd