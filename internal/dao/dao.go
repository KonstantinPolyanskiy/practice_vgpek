package dao

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"practice_vgpek/internal/dao/action"
)

type Aggregator struct {
	actionDAO action.DAO
}

func New(db *pgxpool.Pool, logger *zap.Logger) Aggregator {
	actionDAO := action.New(db, logger)

	return Aggregator{
		actionDAO: actionDAO,
	}
}
