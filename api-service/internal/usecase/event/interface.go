package event

import (
	"context"

	"github.com/dostonshernazarov/resume_maker/api-service/api/models"
)

type ConsumerConfig interface {
	GetBrokers() []string
	GetTopic() string
	GetGroupID() string
	GetHandler() ConsumerHandler
}

type ConsumerHandler interface {
	Handle(ctx context.Context, key, value []byte) error
}

type BrokerConsumer interface {
	Run() error
	RegisterConsumer(cfg ConsumerConfig)
	Close()
}

type BrokerProducer interface {
	ProduceUserToCreate(ctx context.Context, key string, value *models.UserRes) error
	Close()
}
