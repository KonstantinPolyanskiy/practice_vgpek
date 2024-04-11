-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS solved_practice (
    solved_practice_id serial primary key not null,
    performed_account_id integer not null  references account(account_id),
    issued_practice_id integer not null references issued_practice(issued_practice_id),
    mark int not null default 0,
    mark_time timestamp default null,
    solved_time timestamp default null,
    path varchar not null,
    is_deleted timestamp default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS solved_practice;
-- +goose StatementEnd
