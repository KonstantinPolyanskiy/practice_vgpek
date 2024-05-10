-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO internal_action (internal_action_name, description)
VALUES
            ('ADD', 'Действие для добавления'),
            ('GET', 'Действие для получения'),
            ('EDIT', 'Действие для изменения'),
            ('DELETE', 'Действие для удаления');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM internal_action
WHERE
internal_action_name = 'DELETE' AND internal_action_name = 'EDIT' AND internal_action_name = 'GET' AND internal_action_name = 'ADD';
-- +goose StatementEnd
