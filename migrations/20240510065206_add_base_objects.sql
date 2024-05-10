-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO internal_object (internal_object_name, description)
VALUES
            ('MARK', 'Объект для работы с оценками'), ('SOLVED_PRACTICE', 'Объект для работы с практическими работами'),
            ('ISSUED_PRACTICE', 'Объект для работы с практическими заданиями'), ('ACCOUNT', 'Объект для работы с аккаунтами'),
            ('PERSON', 'Объект для работы с пользователями'), ('RBAC', 'Объект для работы с пользователями');


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM internal_object WHERE internal_object_name = 'MARK' AND internal_object_name = 'SOLVED_PRACTICE'
                              AND internal_object_name = 'ISSUED_PRACTICE' AND internal_object_name = 'ACCOUNT'
                              AND internal_object_name = 'PERSON'
-- +goose StatementEnd
