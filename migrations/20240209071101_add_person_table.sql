-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE person (
    person_id serial PRIMARY KEY NOT NULL,
    account_id integer REFERENCES account (account_id),
    first_name varchar NOT NULL,
    middle_name varchar NOT NULL DEFAULT '',
    last_name varchar NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE person;
-- +goose StatementEnd
