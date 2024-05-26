package bus

import (
	"github.com/IBM/sarama"
)

func NewProducer(producerAddrs []string) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(producerAddrs, cfg)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
