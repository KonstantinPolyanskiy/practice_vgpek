-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO internal_role (role_name, description) VALUES
                                                       ('STUDENT', 'Роль студента, может получить свои практические задания и загрузить работу'),
                                                       ('TEACHER', 'Роль учителя, пока права равны администраторским');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM internal_role WHERE role_name = 'STUDENT' AND role_name = 'TEACHER'
-- +goose StatementEnd
