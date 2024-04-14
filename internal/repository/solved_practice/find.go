package solved_practice

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dberr"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/practice/solved"
)

func (r Repository) ById(ctx context.Context, id int) (solved.Entity, error) {
	l := r.l.With(
		zap.String("операция", operation.GetSolvedPracticeInfoById),
		zap.String("слой", "репозиторий"),
	)

	getPracticeQuery := `SELECT * FROM solved_practice WHERE solved_practice_id=$1`

	row, err := r.db.Query(ctx, getPracticeQuery, id)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return solved.Entity{}, dberr.ErrNotFound
		}

		return solved.Entity{}, errors.New("unknown error")
	}

	practice, err := pgx.CollectOneRow(row, pgx.RowToStructByName[solved.Entity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return solved.Entity{}, errors.New("unknown error")
	}

	return practice, nil
}
