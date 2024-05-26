package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/bus"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/client"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/repositories/user"
	"github.com/google/uuid"
	"log"
	"os"
	"strings"
)

func processCommand(
	ctx context.Context,
	userRepo userRepos.UserRepositoryInterface,
	userId int,
	amount float32,
) error {
	balance, err := userRepo.GetUserBalance(ctx, userId)
	if err != nil {
		return err
	}
	if balance-amount < 0 {
		return errors.New("insufficient balance")
	}
	err = userRepo.SetUserBalance(ctx, userId, amount)
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
	repo := userRepos.NewUserRepository(pool)

	kafkaServers := strings.Split(strings.Trim(os.Getenv("KAFKA_SERVERS"), "[]"), ", ")
	consumer, err := bus.NewConsumer(kafkaServers)
	if err != nil {
		panic(fmt.Sprintf("failed to start consumer: %s", err))
	}
	producer, err := bus.NewProducer(kafkaServers)
	if err != nil {
		panic(fmt.Sprintf("failed to start producer: %s", err))
	}
	partition, err := consumer.ConsumePartition("user_approve_events", 0, sarama.OffsetNewest)
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
					userId    int
					amount    float32
				}
				var cm commandMsg
				if err = json.Unmarshal(msg.Value, &cm); err != nil {
					return
				}

				if _, exists := alreadyProcessed[cm.eventUUID]; exists {
					return
				}
				alreadyProcessed[cm.eventUUID] = true

				err = processCommand(ctx, repo, cm.userId, cm.amount)
				var status string
				if err != nil {
					log.Printf("failed to consume command: %s", err)
					status = "failed"
				} else {
					status = "success"
				}
				_, _, err = producer.SendMessage(&sarama.ProducerMessage{
					Topic: "user_finished_events",
					Value: sarama.StringEncoder(
						fmt.Sprintf(
							"{\"event_uuid\": \"%s\", \"status\": %s}",
							cm.eventUUID,
							status,
						),
					),
				})
			}
		}
	}()
}
