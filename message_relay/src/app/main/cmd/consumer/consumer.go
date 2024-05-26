package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/bus"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/client"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/repositories/command"
	"github.com/google/uuid"
	"log"
	"os"
	"strings"
)

func processCommand(
	ctx context.Context,
	repo commandRepos.CommandOutboxRepositoryInterface,
	eventUUID uuid.UUID,
	status string,
) error {
	if status == "failed" {
		err := repo.RollbackCommandTransaction(ctx, eventUUID)
		if err != nil {
			return err
		}
	} else if status == "success" {
		err := repo.SetCommandOutboxesStatus(ctx, []uuid.UUID{eventUUID}, status)
		if err != nil {
			return err
		}
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
	consumer, err := bus.NewConsumer(kafkaServers)
	if err != nil {
		panic(fmt.Sprintf("failed to start consumer: %s", err))
	}
	partition, err := consumer.ConsumePartition("user_finished_events", 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("failed to consume partition %d: %s", partition, err)
	}

	alreadyProcessed := make(map[uuid.UUID]bool) // stub
	go func() {
		for {
			select {
			case msg, ok := <-partition.Messages():
				if !ok {
					log.Println("Channel closed, exiting goroutine")
					return
				}
				type commandMsg struct {
					eventUUID uuid.UUID
					status    string
				}
				var cm commandMsg
				if err = json.Unmarshal(msg.Value, &cm); err != nil {
					return
				}

				if _, exists := alreadyProcessed[cm.eventUUID]; exists {
					return
				}
				alreadyProcessed[cm.eventUUID] = true

				err = processCommand(ctx, repo, cm.eventUUID, cm.status)
				if err != nil {
					log.Printf("failed to process command: %s", err)
				}
			}
		}
	}()
}
