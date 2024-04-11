package issued_practice

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dberr"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/practice/issued"
)

func (r Repository) ById(ctx context.Context, id int) (issued.Entity, error) {
	l := r.l.With(
		zap.String("операция", operation.GetIssuedPracticeInfoById),
		zap.String("слой", "репозиторий"),
	)

	getPracticeQuery := `SELECT * FROM issued_practice WHERE issued_practice_id=$1`

	row, err := r.db.Query(ctx, getPracticeQuery, id)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return issued.Entity{}, dberr.ErrNotFound
		}

		return issued.Entity{}, errors.New("unknown error")
	}

	practice, err := pgx.CollectOneRow(row, pgx.RowToStructByName[issued.Entity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return issued.Entity{}, errors.New("unknown error")
	}

	return practice, nil
}
