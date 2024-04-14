package solved_practice

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dberr"
	"practice_vgpek/internal/model/practice/solved"
)

func (r Repository) Update(ctx context.Context, practice solved.Entity) (solved.Entity, error) {
	l := r.l.With(
		zap.String("операция", "обновление практической работы"),
		zap.String("слой", "репозиторий"),
	)

	updatePracticeQuery := `UPDATE solved_practice
							SET performed_account_id = @PerformedAccountId,
								issued_practice_id = @IssuedPracticeId,
								mark = @Mark,
								mark_time = @MarkTime,
								solved_time = @SolvedTime,
								path = @Path,
								is_deleted = @IsDeleted
							WHERE solved_practice_id = @SolvedPracticeId`

	args := pgx.NamedArgs{
		"PerformedAccountId": practice.PerformedAccountId,
		"IssuedPracticeId":   practice.IssuedPracticeId,
		"Mark":               practice.Mark,
		"MarkTime":           practice.MarkTime,
		"SolvedTime":         practice.SolvedTime,
		"Path":               practice.Path,
		"IsDeleted":          practice.IsDeleted,
		"SolvedPracticeId":   practice.SolvedPracticeId,
	}

	_, err := r.db.Exec(ctx, updatePracticeQuery, args)
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		return solved.Entity{}, errors.New("unknown error")
	}

	getPracticeQuery := `SELECT * FROM solved_practice WHERE solved_practice_id=$1`

	row, err := r.db.Query(ctx, getPracticeQuery, practice.SolvedPracticeId)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return solved.Entity{}, dberr.ErrNotFound
		}

		return solved.Entity{}, errors.New("unknown error")
	}

	updatedPractice, err := pgx.CollectOneRow(row, pgx.RowToStructByName[solved.Entity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return solved.Entity{}, errors.New("unknown error")
	}

	return updatedPractice, nil
}
