package commandRepos

import (
	"context"
	"github.com/MercyClassic/go_saga/src/app/domain/entities/command"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/client"
)

type CommandRepositoryInterface interface {
	GetCommandById(ctx context.Context, userId int) (*commandEntities.Command, error)
	GetCommands(ctx context.Context) ([]*commandEntities.Command, error)
	SaveCommand(ctx context.Context, user *commandEntities.Command) error
}

type CommandRepository struct {
	client client.Client
}

func NewCommandRepository(client client.Client) *CommandRepository {
	return &CommandRepository{client: client}
}

func (r *CommandRepository) GetCommandById(
	ctx context.Context,
	commandId int,
) (*commandEntities.Command, error) {
	q := `
		select 
		    id, description, user_id, amount, created_at
		from command
		where id=$1;
	`
	command := commandEntities.Command{}
	err := r.client.QueryRow(ctx, q, commandId).Scan(
		&command.Id,
		&command.Description,
		&command.UserId,
		&command.Amount,
		&command.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &command, nil
}

func (r *CommandRepository) GetCommands(ctx context.Context) ([]*commandEntities.Command, error) {
	q := `
		select 
		    id, description, user_id, amount, created_at
		from command;
	`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	commands := make([]*commandEntities.Command, 0, 100)
	for rows.Next() {
		var command commandEntities.Command
		if err := rows.Scan(
			&command.Id,
			&command.Description,
			&command.UserId,
			&command.Amount,
			&command.CreatedAt,
		); err != nil {
			return nil, err
		}
		commands = append(commands, &command)
	}
	return commands, nil
}

func (r *CommandRepository) SaveCommand(
	ctx context.Context,
	command *commandEntities.Command,
) error {
	q := `
		insert into command
		    (description, user_id, amount)
		values ($1, $2, $3)
		returning id, created_at;
	`
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return err
	}

	err = tx.QueryRow(
		ctx,
		q,
		command.Description,
		command.UserId,
		command.Amount,
	).Scan(&command.Id, &command.CreatedAt)

	if err != nil {
		return err
	}

	outbox := commandEntities.NewCommandOutbox(command.Id)
	q = `
		insert into command_outbox
		    (event_uuid, command_id)
		values ($1, $2)
	`
	_, err = tx.Exec(ctx, q, outbox.EventUuid, outbox.CommandId)

	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
