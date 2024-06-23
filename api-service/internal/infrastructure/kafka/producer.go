package kafka

import (
	"context"
	"encoding/json"

	"github.com/dostonshernazarov/resume_maker/api/models"

	// "github.com/dostonshernazarov/resume_maker/internal/entity"
	configpkg "github.com/dostonshernazarov/resume_maker/internal/pkg/config"

	"github.com/segmentio/kafka-go"

	// otlp_pkg "github.com/dostonshernazarov/resume_maker/internal/pkg/otlp"
	"go.uber.org/zap"
)

type Producer struct {
	logger     *zap.Logger
	userCreate *kafka.Writer
}

func NewProducer(config *configpkg.Config, logger *zap.Logger) *Producer {
	return &Producer{
		logger: logger,
		userCreate: &kafka.Writer{
			Addr:                   kafka.TCP(config.Kafka.Address...),
			Topic:                  config.Kafka.Topic.UserCreateTopic,
			Balancer:               &kafka.Hash{},
			RequiredAcks:           kafka.RequireAll,
			AllowAutoTopicCreation: true,
			Async:                  true,
			Completion: func(messages []kafka.Message, err error) {
				if err != nil {
					logger.Error("kafka userCreate", zap.Error(err))
				}
				for _, message := range messages {
					logger.Sugar().Info(
						"kafka UserC reatemessage",
						zap.Int("partition", message.Partition),
						zap.Int64("offset", message.Offset),
						zap.String("key", string(message.Key)),
						zap.String("value", string(message.Value)),
					)
				}
			},
		},
	}
}

func (p *Producer) ProduceUserToCreate(ctx context.Context, key string, value *models.UserRes) error {
	byteValue, err := json.Marshal(&value)
	if err != nil {
		return err
	}

	message := p.buildMessage(key, byteValue)

	return p.userCreate.WriteMessages(ctx, message)

}

func (p *Producer) buildMessage(key string, value []byte) kafka.Message {
	return kafka.Message{
		Key:   []byte(key),
		Value: value,
	}
}

func (p *Producer) Close() {
	if err := p.userCreate.Close(); err != nil {
		p.logger.Error("error during close writer investmentCreated", zap.Error(err))
	}
}
