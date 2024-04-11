package solved_practice

import "go.uber.org/zap"

type Service struct {
	l *zap.Logger
}

func NewSolvedPracticeService(logger *zap.Logger) Service {
	return Service{
		l: logger,
	}
}
