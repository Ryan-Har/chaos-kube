package common

import (
	"context"
	"encoding/json"
	"github.com/Ryan-Har/chaos-kube/pkg/config"
	"github.com/Ryan-Har/chaos-kube/pkg/message"
	"github.com/go-redis/redis/v8"
	"log/slog"
	"time"
)

type RedisClientWrapper struct {
	*redis.Client
}

// Creates redis client from config.
func NewRedisClient(cfg config.RedisConfig) *RedisClientWrapper {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &RedisClientWrapper{client}
}

// Retrieves messages from redis stream and adds them to messageChan.
func (c *RedisClientWrapper) ReadStreamToChan(rExArgs *redis.XReadGroupArgs, messageChan chan<- message.Message) {
	ctx := context.Background()
	for {
		messages, err := c.XReadGroup(ctx, rExArgs).Result()

		if err != nil {
			clientName, _ := c.ClientGetName(ctx).Result()
			slog.Error("Error reading from redis stream", "client", clientName, "stream", rExArgs.Streams[0])
			time.Sleep(3 * time.Second)
			continue
		}

		for _, msg := range messages {
			for _, xMessage := range msg.Messages {
				dataStr, ok := xMessage.Values["data"].(string)
				if !ok {
					slog.Error("Error: No 'data' field found or type assertion failed")
					continue
				}

				var data message.Message
				err := json.Unmarshal([]byte(dataStr), &data)
				if err != nil {
					slog.Error("Error unmarshalling JSON data to message.Message", "error", err, "data", data)
					continue
				}
				slog.Debug("received message", "stream", rExArgs.Streams[0], "consumer group", rExArgs.Group, "message", data)
				messageChan <- data
			}
		}
	}
}

// Create consumer group for a stream if it doesn't already exist
func (c *RedisClientWrapper) CreateConsumerGroup(stream string, consumerGroup string) {
	ctx := context.Background()
	c.XGroupCreateMkStream(ctx, stream, consumerGroup, "$")
}

// Send Message to Specific stream
func (c *RedisClientWrapper) SendMessage(msg *message.Message, stream string) error {
	if err := msg.Validate(); err != nil {
		return err
	}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	ctx := context.Background()
	_, err = c.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{
			"data": jsonMsg,
		},
	}).Result()
	return err
}
