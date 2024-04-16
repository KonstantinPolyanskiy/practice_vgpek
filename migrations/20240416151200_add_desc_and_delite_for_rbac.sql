-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE internal_object ADD description varchar not null default '';
ALTER TABLE internal_object ADD created_at timestamp not null default now();
ALTER TABLE internal_object ADD is_deleted timestamp null;

ALTER TABLE internal_action ADD description varchar not null default '';
ALTER TABLE internal_action ADD created_at timestamp not null default now();
ALTER TABLE internal_action ADD is_deleted timestamp null;

ALTER TABLE internal_role ADD description varchar not null default '';
ALTER TABLE internal_role ADD created_at timestamp not null default now();
ALTER TABLE internal_role ADD is_deleted timestamp null;

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
