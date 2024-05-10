-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO role_permission
    (internal_role_id, internal_action_id, internal_object_id)
VALUES
    (2, 1, 7), -- Для роли TEACHER: добавить ключ
    (2, 2, 7), -- Для роли TEACHER: получить ключ
    (2, 3, 7), -- Для роли TEACHER: изменить ключ
    (2, 4, 7); -- Для роли TEACHER: удалить ключ

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
