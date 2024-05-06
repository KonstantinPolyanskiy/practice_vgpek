package solved_practice

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"mime/multipart"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
)

type SolvedPracticeDAO interface {
	Save(ctx context.Context, data dto.NewSolvedPractice) (entity.SolvedPractice, error)
	ById(ctx context.Context, id int) (entity.SolvedPractice, error)
	Update(ctx context.Context, old entity.SolvedPracticeUpdate) (entity.SolvedPractice, error)
}

type AccountDAO interface {
	ById(ctx context.Context, id int) (entity.Account, error)
}

type PersonDAO interface {
	ByAccountId(ctx context.Context, id int) (entity.Person, error)
}

type IssuedPracticeDAO interface {
	ById(ctx context.Context, id int) (entity.IssuedPractice, error)
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
	logger *zap.Logger

	accountDAO AccountDAO
	personDAO  PersonDAO

	accountMediator        AccountMediator
	issuedPracticeMediator IssuedPracticeMediator

	fileStorage PracticeFileStorage

	solvedPracticeDAO SolvedPracticeDAO
	issuedPracticeDAO IssuedPracticeDAO
}

func New(
	accountMediator AccountMediator, issuedPracticeMediator IssuedPracticeMediator,
	fileStorage PracticeFileStorage, solvedPracticeDAO SolvedPracticeDAO, issuedPracticeDAO IssuedPracticeDAO,
	personDAO PersonDAO, accountDAO AccountDAO, logger *zap.Logger) Service {
	return Service{
		accountDAO: accountDAO,
		personDAO:  personDAO,

		accountMediator:        accountMediator,
		issuedPracticeMediator: issuedPracticeMediator,

		fileStorage: fileStorage,

		solvedPracticeDAO: solvedPracticeDAO,
		issuedPracticeDAO: issuedPracticeDAO,

		logger: logger,
	}
}

func (s Service) EntityToDomain(ctx context.Context, accountId int, entity entity.SolvedPractice) (domain.SolvedPractice, error) {
	issuedPracticeEntity, err := s.issuedPracticeDAO.ById(ctx, entity.IssuedPracticeId)
	if err != nil {
		return domain.SolvedPractice{}, err
	}

	teacherAccount, err := s.accountDAO.ById(ctx, issuedPracticeEntity.AccountId)
	if err != nil {
		return domain.SolvedPractice{}, err
	}

	teacher, err := s.personDAO.ByAccountId(ctx, teacherAccount.Id)
	if err != nil {
		return domain.SolvedPractice{}, err
	}

	student, err := s.personDAO.ByAccountId(ctx, accountId)
	if err != nil {
		return domain.SolvedPractice{}, err
	}

	var isDeleted bool

	if entity.IsDeleted != nil {
		isDeleted = true
	}

	practice := domain.SolvedPractice{
		Id:               entity.Id,
		IssuedPracticeId: entity.IssuedPracticeId,
		IssuerName:       fmt.Sprintf("%s %s %s", teacher.FirstName, teacher.MiddleName, teacher.LastName),
		AuthorName:       fmt.Sprintf("%s %s %s", student.FirstName, student.MiddleName, student.LastName),
		AuthorId:         student.AccountId,
		Mark:             entity.Mark,
		MarkTime:         entity.MarkTime,
		SolvedTime:       *entity.SolvedTime,
		IsDeleted:        isDeleted,
		DeletedAt:        entity.IsDeleted,
	}

	return practice, nil
}
