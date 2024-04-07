package issued_practice

import (
	"context"
	"go.uber.org/zap"
	"mime/multipart"
	"practice_vgpek/internal/model/practice/issued"
)

type IssuedPracticeRepository interface {
	Save(ctx context.Context, dto issued.DTO) (issued.)
}

type IssuedPracticeFileStorage interface {
	// SaveFile возвращает путь, по которому был сохранен файл
	SaveFile(ctx context.Context, file *multipart.File, root, ext, name string) (string, error)
}

type AccountMediator interface {
	HasAccess(ctx context.Context, roleId int, objectName, actionName string) (bool, error)
}

type Service struct {
	l  *zap.Logger
	r  IssuedPracticeRepository
	fs IssuedPracticeFileStorage
	am AccountMediator
}

func NewIssuedPracticeService(issuedRepo IssuedPracticeRepository, fileStorage IssuedPracticeFileStorage, accountMediator AccountMediator, logger *zap.Logger) Service {
	return Service{
		l:  logger,
		r:  issuedRepo,
		fs: fileStorage,
		am: accountMediator,
	}
}
