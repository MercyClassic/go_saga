package userServices

import (
	"context"
	"github.com/MercyClassic/go_saga/src/app/domain/entities/user"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/repositories/user"
)

type UserServiceInterface interface {
	GetUser(ctx context.Context) (*userEntities.User, error)
	GetUsers(ctx context.Context) ([]*userEntities.User, error)
	CreateUser(ctx context.Context) error
}

type UserService struct {
	userRepo userRepos.UserRepositoryInterface
}

func NewUserService(userRepo userRepos.UserRepositoryInterface) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUser(
	ctx context.Context,
	userId int,
) (*userEntities.User, error) {
	user, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]*userEntities.User, error) {
	userList, err := s.userRepo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return userList, nil
}

func (s *UserService) CreateUser(
	ctx context.Context,
	name, username string,
) (*userEntities.User, error) {
	user := userEntities.NewUser(name, username)
	err := s.userRepo.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
