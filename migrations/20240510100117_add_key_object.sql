-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO internal_object
    (internal_object_name, description)
VALUES
    ('KEY', 'Объект для действия с ключами');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM internal_object WHERE internal_object_name = 'KEY';
-- +goose StatementEnd
