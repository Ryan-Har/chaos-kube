package common

import (
	"github.com/Ryan-Har/chaos-kube/pkg/message"
	"github.com/Ryan-Har/chaos-kube/pkg/streams"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// RedisClient defines methods needed from the Redis client.
type RedisClient interface {
	ReadStreamToChan(rExArgs *redis.XReadGroupArgs, messageChan *chan message.Message)
	CreateConsumerGroup(stream string, consumerGroup string)
	SendMessageToStream(msg *message.Message, stream streams.RedisStreams) error
}

// ResponseSender defines methods for managing responses.
type ResponseManager interface {
	Send(m *message.Message) error
	Add(id uuid.UUID, c chan<- message.Message)
}
