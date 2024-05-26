package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/bus"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/client"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/repositories/command"
	"github.com/google/uuid"
	"log"
	"os"
	"strings"
	"time"
)

func processCommand(
	ctx context.Context,
	commandRepo commandRepos.CommandOutboxRepositoryInterface,
	producer sarama.SyncProducer,
) error {
	outboxes, err := commandRepo.GetCommandOutboxesByStatus(ctx, "")
	if err != nil {
		return err
	}
	eventUuids := make([]uuid.UUID, 0, len(outboxes))
	messages := make([]*sarama.ProducerMessage, 0, len(outboxes))
	for _, outbox := range outboxes {
		eventUuids = append(eventUuids, outbox.EventUuid)
		messages = append(
			messages,
			&sarama.ProducerMessage{
				Topic: "user_approve_events",
				Value: sarama.StringEncoder(
					fmt.Sprintf(
						"{\"event_uuid\": \"%s\", \"user_id\": %d, \"amount\": %d}",
						outbox.EventUuid,
						outbox.UserId,
						outbox.Amount,
					),
				),
			},
		)
	}
	err = producer.SendMessages(messages)
	if err != nil {
		return err
	}
	err = commandRepo.SetCommandOutboxesStatus(ctx, eventUuids, "processing")
	return nil
}

func processSaga(
	ctx context.Context,
	commandRepo commandRepos.CommandOutboxRepositoryInterface,
	producer sarama.SyncProducer,
) error {
	err := processCommand(ctx, commandRepo, producer)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	pool, err := client.New(ctx, os.Getenv("db_uri"))
	if err != nil {
		panic("cant connect to db")
	}
	repo := commandRepos.NewCommandRepository(pool)
	kafkaServers := strings.Split(strings.Trim(os.Getenv("KAFKA_SERVERS"), "[]"), ", ")
	producer, err := bus.NewProducer(kafkaServers)
	if err != nil {
		panic(fmt.Sprintf("failed to start producer: %s", err))
	}

	for {
		err = processSaga(ctx, repo, producer)
		if err != nil {
			log.Println("process saga err: ", err)
		}
		time.Sleep(5 * time.Second)
	}
}
