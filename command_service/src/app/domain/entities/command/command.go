package commandEntities

import (
	"github.com/google/uuid"
	"time"
)

type Command struct {
	Id          int
	Description string
	UserId      int
	Amount      float32
	CreatedAt   time.Time
}

type CommandOutbox struct {
	EventUuid uuid.UUID
	CommandId int
	Status    string
}

func NewCommand(
	description string,
	userId int,
	amount float32,
) *Command {
	return &Command{
		Description: description,
		UserId:      userId,
		Amount:      amount,
	}
}

func NewCommandOutbox(
	commandId int,
) *CommandOutbox {
	return &CommandOutbox{
		EventUuid: uuid.New(),
		CommandId: commandId,
	}
}
