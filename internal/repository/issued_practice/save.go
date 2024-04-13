package issued_practice

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/practice/issued"
)

func (r Repository) Save(ctx context.Context, dto issued.DTO) (issued.Entity, error) {
	l := r.l.With(
		zap.String("операция", operation.UploadIssuedPracticeOperation),
		zap.String("слой", "репозиторий"),
	)

	var insertedPracticeId int

	insertPracticeQuery := `
	INSERT INTO issued_practice
	    (account_id, target_groups, title, theme, major, practice_path, upload_at) 
	VALUES 
		(@AccountId, @TargetGroups, @Title, @Theme, @Major, @PracticePath, @UploadAt)
	RETURNING issued_practice_id
`
	args := pgx.NamedArgs{
		"AccountId":    dto.AccountId,
		"TargetGroups": dto.TargetGroups,
		"Title":        dto.Title,
		"Theme":        dto.Theme,
		"Major":        dto.Major,
		"PracticePath": dto.Path,
		"UploadAt":     dto.UploadAt,
	}

	err := r.db.QueryRow(ctx, insertPracticeQuery, args).Scan(&insertedPracticeId)
	if err != nil {
		l.Warn("ошибка при выполнении запроса", zap.Error(err))

		return issued.Entity{}, err
	}

	getPracticeQuery := `SELECT * FROM issued_practice WHERE issued_practice_id=$1`

	row, err := r.db.Query(ctx, getPracticeQuery, insertedPracticeId)
	if err != nil {
		l.Warn("ошибка при выполнении запроса", zap.Error(err))

		return issued.Entity{}, err
	}
	defer row.Close()

	savedPractice, err := pgx.CollectOneRow(row, pgx.RowToStructByName[issued.Entity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return issued.Entity{}, err
	}

	return savedPractice, nil
}
