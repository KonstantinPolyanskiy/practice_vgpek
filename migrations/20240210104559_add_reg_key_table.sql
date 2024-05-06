-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE  IF NOT EXISTS registration_key (
    reg_key_id serial PRIMARY KEY NOT NULL,
    internal_role_id integer NOT NULL REFERENCES internal_role(internal_role_id),
    body_key varchar NOT NULL,
    max_count_usages integer NOT NULL,
    current_count_usages integer NOT NULL,
    created_at timestamp NOT NULL,
    is_valid boolean NOT NULL DEFAULT true,
    invalidation_time timestamp DEFAULT NULL
);
ALTER TABLE account  ADD IF NOT EXISTS reg_key_id integer REFERENCES registration_key(reg_key_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE account DROP COLUMN IF EXISTS reg_key_id;
DROP TABLE IF EXISTS registration_key;
-- +goose StatementEnd
