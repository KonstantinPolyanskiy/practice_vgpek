package solved_practice

import (
	"context"
	"go.uber.org/zap"
	"mime/multipart"
	"practice_vgpek/internal/model/practice/issued"
	"practice_vgpek/internal/model/practice/solved"
)

const (
	GetActionName    = "GET"
	AddActionName    = "ADD"
	SolvedObjectName = "SOLVED PRACTICE"
	MarkObjectName   = "MARK"
)

type SolvedPracticeRepository interface {
	Save(ctx context.Context, dto solved.DTO) (solved.Entity, error)
	ById(ctx context.Context, id int) (solved.Entity, error)
	Update(ctx context.Context, practice solved.Entity) (solved.Entity, error)
}

type IssuedPracticeRepository interface {
	ById(ctx context.Context, id int) (issued.Entity, error)
}

type AccountMediator interface {
	HasAccess(ctx context.Context, accountId int, objectName, actionName string) (bool, error)
}

type IssuedPracticeMediator interface {
	IssuedGroupMatch(ctx context.Context, accountId, practiceId int) (bool, error)
}

type PracticeFileStorage interface {
	// SaveFile возвращает путь, по которому был сохранен файл
	SaveFile(ctx context.Context, file *multipart.File, root, ext, name string) (string, error)
}

type Service struct {
	l   *zap.Logger
	am  AccountMediator
	ipm IssuedPracticeMediator
	fs  PracticeFileStorage
	spr SolvedPracticeRepository
	ipr IssuedPracticeRepository
}

func NewSolvedPracticeService(am AccountMediator, ipm IssuedPracticeMediator, fs PracticeFileStorage, spr SolvedPracticeRepository, ipr IssuedPracticeRepository, logger *zap.Logger) Service {
	return Service{
		am:  am,
		ipm: ipm,
		fs:  fs,
		spr: spr,
		ipr: ipr,
		l:   logger,
	}
}
