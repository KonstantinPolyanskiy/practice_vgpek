-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE issued_practice (
    issued_practice_id SERIAL PRIMARY KEY NOT NULL,
    account_id INTEGER NOT NULL REFERENCES account(account_id),
    target_groups VARCHAR[] NOT NULL,
    title VARCHAR NOT NULL DEFAULT 'unknown',
    theme VARCHAR NOT NULL DEFAULT 'unknown',
    major VARCHAR NOT NULL DEFAULT 'unknown',
    practice_path varchar NOT NULL DEFAULT 'unknown',
    upload_at timestamp NOT NULL ,
    deleted_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS issued_practice;
-- +goose StatementEnd
