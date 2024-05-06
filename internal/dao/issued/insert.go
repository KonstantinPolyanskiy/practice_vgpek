package issued

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

func (dao DAO) Save(ctx context.Context, data dto.NewIssuedPractice) (entity.IssuedPractice, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SaveIssuedPracticeDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	insertQuery := `INSERT INTO 
						issued_practice (account_id, target_groups, title, theme, major, practice_path, upload_at) 
					VALUES 
					    (@AccountId, @TargetGroups, @Title, @Theme, @Major, @PracticePath, @UploadAt)
					RETURNING issued_practice_id`

	args := pgx.NamedArgs{
		"AccountId":    data.AccountId,
		"TargetGroups": data.TargetGroups,
		"Title":        data.Title,
		"Theme":        data.Theme,
		"Major":        data.Major,
		"PracticePath": data.Path,
		"UploadAt":     data.UploadAt,
	}

	l.Debug("аргументы запроса",
		zap.Int("id аккаунта", args["AccountId"].(int)),
		zap.Strings("целевые группы", args["TargetGroups"].([]string)),
		zap.String("название", args["Title"].(string)),
		zap.String("тема", args["Theme"].(string)),
		zap.String("специальность", args["Major"].(string)),
		zap.String("путь к практике", args["PracticePath"].(string)),
		zap.Time("дата загрузки", args["UploadAt"].(time.Time)),
	)

	var issuedPracticeId int

	now := time.Now()
	err := dao.db.QueryRow(ctx, insertQuery, args).Scan(&issuedPracticeId)
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.IssuedPractice{}, err
	}

	l.Debug(operation.Insert, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	selectQuery := `SELECT * FROM issued_practice WHERE issued_practice_id=@IssuedPracticeId`

	args = pgx.NamedArgs{
		"IssuedPracticeId": issuedPracticeId,
	}

	now = time.Now()
	rows, err := dao.db.Query(ctx, selectQuery, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.IssuedPractice{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	practice, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.IssuedPractice])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.IssuedPractice{}, err
	}

	return practice, nil
}
