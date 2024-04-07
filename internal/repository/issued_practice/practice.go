package issued_practice

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	l  *zap.Logger
	db *pgxpool.Pool
}

func NewIssuedPracticeRepo(db *pgxpool.Pool, logger *zap.Logger) Repository {
	return Repository{
		l:  logger,
		db: db,
	}
}
