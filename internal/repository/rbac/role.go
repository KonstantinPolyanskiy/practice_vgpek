package rbac

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/permissions"
)

var (
	ManyRoleErr  = errors.New("неоднозначный результат")
	RoleNotFound = errors.New("роль не найдена")
)

type RoleRepository struct {
	l  *zap.Logger
	db *pgxpool.Pool
}

func NewRoleRepo(db *pgxpool.Pool, logger *zap.Logger) RoleRepository {
	return RoleRepository{
		l:  logger,
		db: db,
	}
}

func (r RoleRepository) SaveRole(ctx context.Context, savingRole permissions.RoleDTO) (permissions.RoleEntity, error) {
	l := r.l.With(
		zap.String("executing query name", "save entity"),
		zap.String("layer", "repo"),
	)

	var insertedRoleId int

	insertedRoleQuery := `
	INSERT INTO internal_role (role_name)
	VALUES (@RoleName)
	RETURNING internal_role_id
`

	l.Debug("insert object", zap.String("query", insertedRoleQuery))

	args := pgx.NamedArgs{
		"RoleName": savingRole.Name,
	}

	l.Debug("args in insert role query", zap.Any("name role", args["RoleName"]))

	// Вставляем объект в БД
	err := r.db.QueryRow(ctx, insertedRoleQuery, args).Scan(&insertedRoleId)
	if err != nil {
		l.Warn("error insert role", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.RoleEntity{}, errors.New("сохраненная роль не найдена")
		}
		return permissions.RoleEntity{}, err
	}

	getRoleQuery := `
	SELECT * FROM internal_role WHERE internal_role_id=$1
`

	l.Debug("get inserted role", zap.String("query", getRoleQuery))

	row, err := r.db.Query(ctx, getRoleQuery, insertedRoleId)
	if err != nil {
		l.Warn("error get inserted role", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.RoleEntity{}, errors.New("сохраненная роль не найдена")
		}
		return permissions.RoleEntity{}, err
	}

	savedRole, err := pgx.CollectOneRow(row, pgx.RowToStructByName[permissions.RoleEntity])
	if err != nil {
		l.Warn("error collect role in struct", zap.Error(err))

		return permissions.RoleEntity{}, err
	}

	return savedRole, nil
}

func (r RoleRepository) RoleByName(ctx context.Context, name string) (permissions.RoleEntity, error) {
	l := r.l.With(
		zap.String("executing query name", "get role by name"),
		zap.String("layer", "repo"),
	)

	var role permissions.RoleEntity

	getRoleQuery := `SELECT * FROM internal_role WHERE role_name=$1`

	err := r.db.QueryRow(ctx, getRoleQuery, name).Scan(&role)
	if err != nil {
		l.Warn("error get role by name",
			zap.String("role name", name),
			zap.Error(err),
		)

		if errors.Is(err, pgx.ErrTooManyRows) {
			return permissions.RoleEntity{}, ManyRoleErr
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.RoleEntity{}, RoleNotFound
		}

		return permissions.RoleEntity{}, errors.New("unknown error")
	}

	return role, nil
}

func (r RoleRepository) RoleById(ctx context.Context, id int) (permissions.RoleEntity, error) {
	l := r.l.With(
		zap.String("executing query name", "get role by id"),
		zap.String("layer", "repo"),
	)

	getRoleQuery := `SELECT * FROM internal_role WHERE internal_role_id = $1`

	row, err := r.db.Query(ctx, getRoleQuery, id)
	if err != nil {
		l.Warn("error get role by id",
			zap.Int("Role id", id),
			zap.Error(err),
		)

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.RoleEntity{}, RoleNotFound
		}

		return permissions.RoleEntity{}, err
	}

	role, err := pgx.CollectOneRow(row, pgx.RowToStructByName[permissions.RoleEntity])
	if err != nil {
		return permissions.RoleEntity{}, err
	}

	return role, nil
}
