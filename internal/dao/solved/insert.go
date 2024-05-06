package solved

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/timeutils"
	"time"
)

func (dao DAO) Save(ctx context.Context, data dto.NewSolvedPractice) (entity.SolvedPractice, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SaveSolvedPracticeDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	insertQuery := `INSERT INTO 
						solved_practice (performed_account_id, issued_practice_id, solved_time, path) 
					VALUES 
					    (@PerformedAccountId, @IssuedPracticeId, @SolvedTime, @Path)
					RETURNING solved_practice_id`

	args := pgx.NamedArgs{
		"PerformedAccountId": data.PerformedAccountId,
		"IssuedPracticeId":   data.IssuedPracticeId,
		"SolvedTime":         data.SolvedTime,
		"Path":               data.Path,
	}

	l.Debug("аргументы запроса",
		zap.Int("id решившего аккаунта", args["PerformedAccountId"].(int)),
		zap.Int("id решенной практической", args["IssuedPracticeId"].(int)),
		zap.Time("время загрузки", args["SolvedTime"].(time.Time)),
		zap.String("путь к практике", args["Path"].(string)),
	)

	var solvedPracticeId int

	now := time.Now()
	err := dao.db.QueryRow(ctx, insertQuery, args).Scan(&solvedPracticeId)
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.SolvedPractice{}, err
	}

	l.Debug(operation.Insert, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	selectQuery := `SELECT * FROM solved_practice WHERE solved_practice_id=@SolvedPracticeId`

	args = pgx.NamedArgs{
		"SolvedPracticeId": solvedPracticeId,
	}

	now = time.Now()
	rows, err := dao.db.Query(ctx, selectQuery, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.SolvedPractice{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	practice, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.SolvedPractice])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.SolvedPractice{}, err
	}

	return practice, nil
}
