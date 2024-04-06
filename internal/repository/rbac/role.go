package rbac

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dberr"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/params"
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

func (rr RoleRepository) SaveRole(ctx context.Context, savingRole permissions.RoleDTO) (permissions.RoleEntity, error) {
	l := rr.l.With(
		zap.String("операция", operation.AddRoleOperation),
		zap.String("слой", "репозиторий"),
	)

	var insertedRoleId int

	insertedRoleQuery := `
	INSERT INTO internal_role (role_name)
	VALUES (@RoleName)
	RETURNING internal_role_id
`

	args := pgx.NamedArgs{
		"RoleName": savingRole.Name,
	}

	// Вставляем объект в БД
	err := rr.db.QueryRow(ctx, insertedRoleQuery, args).Scan(&insertedRoleId)
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.RoleEntity{}, dberr.ErrNotFound
		}
		return permissions.RoleEntity{}, err
	}

	getRoleQuery := `
	SELECT * FROM internal_role WHERE internal_role_id=$1
`

	row, err := rr.db.Query(ctx, getRoleQuery, insertedRoleId)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.RoleEntity{}, dberr.ErrNotFound
		}
		return permissions.RoleEntity{}, err
	}

	savedRole, err := pgx.CollectOneRow(row, pgx.RowToStructByName[permissions.RoleEntity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return permissions.RoleEntity{}, err
	}

	return savedRole, nil
}

func (rr RoleRepository) RoleById(ctx context.Context, id int) (permissions.RoleEntity, error) {
	l := rr.l.With(
		zap.String("операция", operation.GetRoleOperation),
		zap.String("слой", "репозиторий"),
	)

	getRoleQuery := `SELECT * FROM internal_role WHERE internal_role_id = $1`

	row, err := rr.db.Query(ctx, getRoleQuery, id)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.RoleEntity{}, dberr.ErrNotFound
		}

		return permissions.RoleEntity{}, err
	}

	role, err := pgx.CollectOneRow(row, pgx.RowToStructByName[permissions.RoleEntity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return permissions.RoleEntity{}, err
	}

	return role, nil
}

func (rr RoleRepository) RolesByParams(ctx context.Context, params params.Default) ([]permissions.RoleEntity, error) {
	l := rr.l.With(
		zap.String("операция", operation.GetRolesOperation),
		zap.String("слой", "репозиторий"),
	)

	getRolesQuery := squirrel.Select("*").From("internal_role").
		Limit(uint64(params.Limit)).
		Offset(uint64(params.Offset)).
		PlaceholderFormat(squirrel.Dollar)

	q, args, err := getRolesQuery.ToSql()
	if err != nil {
		l.Warn("ошибка подготовки запроса", zap.Error(err))

		return nil, err
	}

	row, err := rr.db.Query(ctx, q, args...)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		return nil, err
	}

	roles, err := pgx.CollectRows(row, pgx.RowToStructByName[permissions.RoleEntity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return nil, err
	}

	return roles, nil
}
