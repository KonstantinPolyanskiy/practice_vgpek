package solved

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/timeutils"
	"time"
)

func (dao DAO) ById(ctx context.Context, id int) (entity.SolvedPractice, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.GetSolvedPracticeInfoByIdDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	selectQuery := `SELECT * FROM solved_practice WHERE solved_practice_id=@SolvedPracticeId`

	args := pgx.NamedArgs{
		"SolvedPracticeId": id,
	}

	now := time.Now()
	rows, err := dao.db.Query(ctx, selectQuery, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.SolvedPractice{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	solvedPractice, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.SolvedPractice])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.SolvedPractice{}, err
	}

	return solvedPractice, nil
}
