package issued

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

func (dao DAO) ById(ctx context.Context, id int) (entity.IssuedPractice, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.GetIssuedPracticeInfoById),
		zap.String(layer.Layer, layer.DataLayer),
	)

	selectQuery := `SELECT * FROM issued_practice WHERE issued_practice_id=@IssuedPracticeId`

	args := pgx.NamedArgs{
		"IssuedPracticeId": id,
	}

	now := time.Now()
	rows, err := dao.db.Query(ctx, selectQuery, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.IssuedPractice{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	issuedPractice, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.IssuedPractice])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.IssuedPractice{}, err
	}

	return issuedPractice, nil
}
