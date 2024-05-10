-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO registration_key
    (internal_role_id, body_key, max_count_usages, current_count_usages, created_at)
VALUES
    (2, 'example_teacher', 100, 0, now());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM registration_key WHERE body_key = 'example_teacher';
-- +goose StatementEnd
