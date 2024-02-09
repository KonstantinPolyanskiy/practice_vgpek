-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE account (
                         account_id serial PRIMARY KEY NOT NULL,
                         login varchar(16) NOT NULL UNIQUE,
                         password_hash varchar NOT NULL,
                         created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         is_active boolean NOT NULL DEFAULT TRUE,
                         deactivate_time timestamp DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE account;
-- +goose StatementEnd
