package issued_practice

import (
	"context"
	"go.uber.org/zap"
	"mime/multipart"
	"practice_vgpek/internal/model/practice/issued"
)

const (
	AddActionName    = "ADD"
	GetActionName    = "GET"
	IssuedObjectName = "ISSUED PRACTICE"
)

type IssuedPracticeRepository interface {
	Save(ctx context.Context, dto issued.DTO) (issued.Entity, error)
	ById(ctx context.Context, id int) (issued.Entity, error)
}

type PracticeMediator interface {
	// IssuedGroupMatch Проверяет, совпадает ли группа студента с одной из целевых груп практического задания
	IssuedGroupMatch(ctx context.Context, accountId, practiceId int) (bool, error)
}

type PracticeFileStorage interface {
	// SaveFile возвращает путь, по которому был сохранен файл
	SaveFile(ctx context.Context, file *multipart.File, root, ext, name string) (string, error)
}

type AccountMediator interface {
	HasAccess(ctx context.Context, roleId int, objectName, actionName string) (bool, error)
}

type Service struct {
	l  *zap.Logger
	r  IssuedPracticeRepository
	fs PracticeFileStorage
	am AccountMediator
	pm PracticeMediator
}

func NewIssuedPracticeService(issuedRepo IssuedPracticeRepository, fileStorage PracticeFileStorage,
	accountMediator AccountMediator, practiceMediator PracticeMediator, logger *zap.Logger) Service {
	return Service{
		l:  logger,
		r:  issuedRepo,
		fs: fileStorage,
		am: accountMediator,
		pm: practiceMediator,
	}
}
