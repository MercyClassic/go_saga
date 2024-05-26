package bus

import (
	"github.com/IBM/sarama"
)

func NewConsumer(consumerAddrs []string) (sarama.Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true

	consumer, err := sarama.NewConsumer(consumerAddrs, cfg)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
