package service

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/dao"
	"practice_vgpek/internal/mediator/account"
	"practice_vgpek/internal/mediator/practice"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/service/issued_practice"
	"practice_vgpek/internal/service/key"
	"practice_vgpek/internal/service/person"
	"practice_vgpek/internal/service/rbac"
	"practice_vgpek/internal/service/solved_practice"
	"practice_vgpek/internal/service/token"
	"practice_vgpek/internal/storage"
)

type AuthnService interface {
	NewUser(ctx context.Context, registration dto.RegistrationReq) (domain.Person, error)
}

type TokenService interface {
	ParseToken(ctx context.Context, token string) (int, error)
	CreateToken(ctx context.Context, cred dto.Credentials) (string, error)
}

type PersonService interface {
	NewUser(ctx context.Context, registration dto.RegistrationReq) (domain.Person, error)

	EntityAccountById(ctx context.Context, req dto.EntityId) (entity.Account, error)
	AccountById(ctx context.Context, req dto.EntityId) (domain.Account, error)

	EntityAccountByParam(ctx context.Context, p params.State) ([]entity.Account, error)
}

type RBACService interface {
	NewAction(ctx context.Context, req dto.NewRBACReq) (domain.Action, error)
	ActionById(ctx context.Context, req dto.EntityId) (domain.Action, error)
	ActionsByParams(ctx context.Context, params params.State) ([]domain.Action, error)

	NewObject(ctx context.Context, req dto.NewRBACReq) (domain.Object, error)
	ObjectById(ctx context.Context, req dto.EntityId) (domain.Object, error)
	ObjectsByParams(ctx context.Context, params params.State) ([]domain.Object, error)

	NewRole(ctx context.Context, req dto.NewRBACReq) (domain.Role, error)
	RoleById(ctx context.Context, req dto.EntityId) (domain.Role, error)
	RolesByParams(ctx context.Context, params params.State) ([]domain.Role, error)

	NewPermission(ctx context.Context, req dto.SetPermissionReq) error
	ByRoleId(ctx context.Context, req dto.EntityId) ([]domain.Permissions, error)
}

type KeyService interface {
	NewKey(ctx context.Context, req dto.NewKeyReq) (domain.Key, error)
	ById(ctx context.Context, req dto.EntityId) (domain.Key, error)
	InvalidateKey(ctx context.Context, id int) (domain.InvalidatedKey, error)
	KeysByParams(ctx context.Context, keyParams params.State) ([]domain.Key, error)
}

type IssuedPracticeService interface {
	Save(ctx context.Context, req dto.NewIssuedPracticeReq) (domain.IssuedPractice, error)
	ById(ctx context.Context, req dto.EntityId) (domain.IssuedPractice, error)
}

type SolvedPracticeService interface {
	Save(ctx context.Context, req dto.NewSolvedPracticeReq) (domain.SolvedPractice, error)
	ById(ctx context.Context, req dto.EntityId) (domain.SolvedPractice, error)

	SetMark(ctx context.Context, req dto.MarkPracticeReq) (domain.SolvedPractice, error)
}

type Service struct {
	PersonService
	TokenService
	KeyService
	RBACService
	IssuedPracticeService
	SolvedPracticeService
}

func New(daoAggregator dao.Aggregator, logger *zap.Logger) Service {
	issuedMediator := practice.NewIssuedPracticeMediator(daoAggregator.AccountDAO, daoAggregator.IssuedDAO, daoAggregator.KeyDAO)
	fileStorage := storage.NewFileStorage()
	rbacService := rbac.New(daoAggregator.ActionDAO, daoAggregator.ObjectDAO, daoAggregator.RoleDAO, daoAggregator.PermissionDAO, logger)

	keyService := key.New(daoAggregator.KeyDAO, daoAggregator.RoleDAO, logger)

	personService := person.New(rbacService, daoAggregator.PermissionDAO, daoAggregator.KeyDAO, daoAggregator.PersonDAO, daoAggregator.AccountDAO, daoAggregator.RoleDAO, keyService, logger)

	accountMediator := account.NewAccountMediator(personService, keyService, rbacService, rbacService)

	tokenService := token.New(daoAggregator.AccountDAO, "ioj9t3r89ug489h", logger)
	issuedService := issued_practice.New(daoAggregator.IssuedDAO, daoAggregator.PersonDAO, fileStorage, accountMediator, issuedMediator, logger)
	solvedService := solved_practice.New(accountMediator, issuedMediator, fileStorage, daoAggregator.SolvedDAO, daoAggregator.IssuedDAO, daoAggregator.PersonDAO, daoAggregator.AccountDAO, logger)

	return Service{
		PersonService:         personService,
		TokenService:          tokenService,
		KeyService:            keyService,
		RBACService:           rbacService,
		IssuedPracticeService: issuedService,
		SolvedPracticeService: solvedService,
	}
}
