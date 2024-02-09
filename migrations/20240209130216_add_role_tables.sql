-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE internal_role (
    internal_role_id serial PRIMARY KEY NOT NULL,
    role_name varchar NOT NULL UNIQUE
);
CREATE TABLE internal_action (
    internal_action_id serial PRIMARY KEY NOT NULL,
    internal_action_name varchar NOT NULL
);
CREATE TABLE internal_object (
    internal_object_id serial PRIMARY KEY NOT NULL,
    internal_object_name varchar NOT NULL
);

CREATE TABLE role_permission (
    role_perm_id serial PRIMARY KEY NOT NULL,
    internal_role_id integer REFERENCES internal_role(internal_role_id),
    internal_action_id integer REFERENCES internal_action(internal_action_id),
    internal_object_id integer REFERENCES internal_object(internal_object_id)
);

ALTER TABLE account ADD internal_role_id integer references internal_role(internal_role_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE account DROP COLUMN IF EXISTS internal_role_id;
DROP TABLE IF EXISTS role_permission;
DROP TABLE IF EXISTS internal_role;
DROP TABLE IF EXISTS internal_action;
DROP TABLE IF EXISTS internal_object;
-- +goose StatementEnd
