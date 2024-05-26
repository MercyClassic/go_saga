package commandRepos

import (
	"context"
	"fmt"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/client"
	"github.com/google/uuid"
)

type CommandOutbox struct {
	EventUuid uuid.UUID
	UserId    int
	Amount    int
	Status    string
}

type CommandOutboxRepositoryInterface interface {
	GetCommandOutboxesByStatus(ctx context.Context, status string) ([]*CommandOutbox, error)
	SetCommandOutboxesStatus(ctx context.Context, outboxes []uuid.UUID, status string) error
	RollbackCommandTransaction(ctx context.Context, eventUuid uuid.UUID) error
}

type CommandRepository struct {
	client client.Client
}

func NewCommandRepository(client client.Client) *CommandRepository {
	return &CommandRepository{client: client}
}

func (c *CommandRepository) GetCommandOutboxesByStatus(ctx context.Context, status string) ([]*CommandOutbox, error) {
	q := `
		select event_uuid, user_id, status, amount
		from command_outboxes
		join command on command_outboxes.command_id = command.id
		where $1
		limit 10;
	`

	var statusExpr string
	if status != "" {
		statusExpr = fmt.Sprintf("status='%s'", status)
	} else {
		statusExpr = "status is null"
	}

	rows, err := c.client.Query(ctx, q, statusExpr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outboxes := make([]*CommandOutbox, 0, 100)
	for rows.Next() {
		var outbox CommandOutbox
		if err = rows.Scan(
			&outbox.EventUuid,
			&outbox.UserId,
			&outbox.Status,
			&outbox.Amount,
		); err != nil {
			return nil, err
		}
		outboxes = append(outboxes, &outbox)
	}
	return outboxes, nil
}

func (c *CommandRepository) SetCommandOutboxesStatus(
	ctx context.Context,
	outboxes []uuid.UUID,
	status string,
) error {
	tx, err := c.client.Begin(ctx)
	if err != nil {
		return err
	}
	q := `
		update command_outboxes 
		set status = $1
		where event_uuid = any($2);
	`
	_, err = tx.Exec(ctx, q, status, outboxes)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommandRepository) RollbackCommandTransaction(ctx context.Context, eventUuid uuid.UUID) error {
	tx, err := c.client.Begin(ctx)
	if err != nil {
		return err
	}
	q := `
		update command_outboxes 
		set status = 'failed'
		where event_uuid = $1;
	`
	_, err = tx.Exec(ctx, q, eventUuid)
	if err != nil {
		return err
	}
	q = `
		delete from command
		join command_outbox on command_outbox.command_id = command.id
		where event_uuid = $1;
	`
	_, err = tx.Exec(ctx, q, eventUuid)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
