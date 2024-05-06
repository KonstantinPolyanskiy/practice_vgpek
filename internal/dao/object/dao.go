package object

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DAO struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func New(db *pgxpool.Pool, logger *zap.Logger) DAO {
	return DAO{
		db:     db,
		logger: logger,
	}
}
