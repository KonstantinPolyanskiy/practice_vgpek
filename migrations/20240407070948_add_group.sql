-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE registration_key ADD IF NOT EXISTS group_name varchar DEFAULT 'unknown' NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE registration_key DROP COLUMN group_name;
-- +goose StatementEnd
