-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE internal_object ADD IF NOT EXISTS description varchar not null default '';
ALTER TABLE internal_object ADD IF NOT EXISTS created_at timestamp not null default now();
ALTER TABLE internal_object ADD IF NOT EXISTS is_deleted timestamp null;

ALTER TABLE internal_action ADD IF NOT EXISTS description varchar not null default '';
ALTER TABLE internal_action ADD IF NOT EXISTS created_at timestamp not null default now();
ALTER TABLE internal_action ADD IF NOT EXISTS is_deleted timestamp null;

ALTER TABLE internal_role ADD IF NOT EXISTS description varchar not null default '';
ALTER TABLE internal_role ADD IF NOT EXISTS created_at timestamp not null default now();
ALTER TABLE internal_role ADD IF NOT EXISTS is_deleted timestamp null;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE internal_object DROP COLUMN description;
ALTER TABLE internal_object DROP COLUMN is_deleted;

ALTER TABLE internal_action DROP COLUMN description;
ALTER TABLE internal_action DROP COLUMN is_deleted;

ALTER TABLE internal_role DROP COLUMN description;
ALTER TABLE internal_role DROP COLUMN is_deleted;
-- +goose StatementEnd
