package solved_practice

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/practice/solved"
	"time"
)

func (r Repository) Save(ctx context.Context, dto solved.DTO) (solved.Entity, error) {
	l := r.l.With(
		zap.String("операция", operation.UploadSolvedPracticeOperation),
		zap.String("слой", "репозиторий"),
	)

	var insertedPracticeId int

	insertPracticeQuery := `
	INSERT INTO solved_practice
		(performed_account_id, issued_practice_id, solved_time, path) 
	VALUES 
		(@PerformedAccountId, @IssuedPracticeId, @SolvedTime, @Path)
	RETURNING solved_practice_id
`
	args := pgx.NamedArgs{
		"PerformedAccountId": dto.PerformedAccountId,
		"IssuedPracticeId":   dto.IssuedPracticeId,
		"SolvedTime":         time.Now(),
		"Path":               dto.Path,
	}

	err := r.db.QueryRow(ctx, insertPracticeQuery, args).Scan(&insertedPracticeId)
	if err != nil {
		l.Warn("ошибка при выполнении запроса", zap.Error(err))

		return solved.Entity{}, err
	}

	getPracticeQuery := `SELECT * FROM solved_practice WHERE solved_practice_id=$1`

	row, err := r.db.Query(ctx, getPracticeQuery, insertedPracticeId)
	if err != nil {
		l.Warn("ошибка при выполнении запроса", zap.Error(err))

		return solved.Entity{}, err
	}
	defer row.Close()

	savedPractice, err := pgx.CollectOneRow(row, pgx.RowToStructByName[solved.Entity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return solved.Entity{}, err
	}

	return savedPractice, nil
}
