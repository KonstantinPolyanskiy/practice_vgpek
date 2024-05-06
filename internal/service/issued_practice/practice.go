package issued_practice

import (
	"context"
	"go.uber.org/zap"
	"mime/multipart"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
)

type IssuedPracticeDAO interface {
	Save(ctx context.Context, data dto.NewIssuedPractice) (entity.IssuedPractice, error)
	ById(ctx context.Context, id int) (entity.IssuedPractice, error)
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

type PersonDAO interface {
	ByAccountId(ctx context.Context, accountId int) (entity.Person, error)
}

type Service struct {
	logger *zap.Logger

	issuedPracticeDAO IssuedPracticeDAO

	personDAO PersonDAO

	fileStorage PracticeFileStorage

	accountMediator AccountMediator
	mediator        PracticeMediator
}

func New(issuedPracticeDAO IssuedPracticeDAO, personDAO PersonDAO, fileStorage PracticeFileStorage,
	accountMediator AccountMediator, practiceMediator PracticeMediator, logger *zap.Logger) Service {
	return Service{
		logger:            logger,
		issuedPracticeDAO: issuedPracticeDAO,
		fileStorage:       fileStorage,
		personDAO:         personDAO,
		accountMediator:   accountMediator,
		mediator:          practiceMediator,
	}
}
