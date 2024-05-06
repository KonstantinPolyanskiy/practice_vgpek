package solved

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/timeutils"
	"time"
)

func (dao DAO) Update(ctx context.Context, practice entity.SolvedPracticeUpdate) (entity.SolvedPractice, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.UpdateSolvedPracticeDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	update := updateQ("solved_practice", practice)

	update = update.Where("solved_practice_id = $2", practice.Id)

	updateQuery, args, err := update.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		l.Error("ошибка сборки запроса", zap.Error(err))
		return entity.SolvedPractice{}, err
	}

	now := time.Now()
	_, err = dao.db.Exec(ctx, updateQuery, args...)
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.SolvedPractice{}, err
	}

	l.Debug(operation.Operation, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	getQuery := `SELECT * FROM solved_practice WHERE solved_practice_id=$1`

	now = time.Now()
	rows, err := dao.db.Query(ctx, getQuery, practice.Id)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.SolvedPractice{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.SolvedPractice])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.SolvedPractice{}, err
	}

	return updated, err
}

func updateQ(table string, newPractice entity.SolvedPracticeUpdate) squirrel.UpdateBuilder {
	updateBuilder := squirrel.Update(table)

	if newPractice.PerformedAccountId != nil {
		updateBuilder = updateBuilder.Set("performed_account_id", newPractice.PerformedAccountId)
	}
	if newPractice.IssuedPracticeId != nil {
		updateBuilder = updateBuilder.Set("issued_practice_id", newPractice.IssuedPracticeId)
	}
	if newPractice.Path != nil {
		updateBuilder = updateBuilder.Set("path", newPractice.Path)
	}
	if newPractice.SolvedTime != nil {
		updateBuilder = updateBuilder.Set("solved_time", newPractice.SolvedTime)
	}
	if newPractice.MarkTime != nil {
		updateBuilder = updateBuilder.Set("mark_time", newPractice.MarkTime)
	}
	if newPractice.IsDeleted != nil {
		updateBuilder = updateBuilder.Set("is_deleted", newPractice.IsDeleted)
	}
	if newPractice.Mark != nil {
		updateBuilder = updateBuilder.Set("mark", newPractice.Mark)
	}

	return updateBuilder
}
