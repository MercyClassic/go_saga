package commandServices

import (
	"context"
	"github.com/MercyClassic/go_saga/src/app/domain/entities/command"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/repositories/command"
)

type CommandServiceInterface interface {
	GetCommand(ctx context.Context) (*commandEntities.Command, error)
	GetCommands(ctx context.Context) ([]*commandEntities.Command, error)
	CreateCommand(ctx context.Context) error
}

type CommandService struct {
	commandRepo commandRepos.CommandRepositoryInterface
}

func NewCommandService(commandRepo commandRepos.CommandRepositoryInterface) *CommandService {
	return &CommandService{commandRepo: commandRepo}
}

func (s *CommandService) GetUser(
	ctx context.Context,
	userId int,
) (*commandEntities.Command, error) {
	user, err := s.commandRepo.GetCommandById(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *CommandService) GetUsers(ctx context.Context) ([]*commandEntities.Command, error) {
	userList, err := s.commandRepo.GetCommands(ctx)
	if err != nil {
		return nil, err
	}
	return userList, nil
}

func (s *CommandService) CreateUser(
	ctx context.Context,
	description string,
	userId int,
	amount float32,
) (*commandEntities.Command, error) {
	command := commandEntities.NewCommand(description, userId, amount)
	err := s.commandRepo.SaveCommand(ctx, command)
	if err != nil {
		return nil, err
	}
	return command, nil
}
