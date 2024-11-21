package common

import (
	"context"
	"encoding/json"
	"github.com/Ryan-Har/chaos-kube/pkg/config"
	"github.com/Ryan-Har/chaos-kube/pkg/message"
	"github.com/Ryan-Har/chaos-kube/pkg/streams"
	"github.com/go-redis/redis/v8"
	"log/slog"
	"os"
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

// Retrieves messages from one redis stream and adds them to messageChan.
func (c *RedisClientWrapper) ReadStreamToChan(rExArgs *redis.XReadGroupArgs, messageChan *chan message.Message) {
	ctx := context.Background()
	slog.Info("bagan reading redis stream", "stream", rExArgs.Streams[0])
	for {
		retryOnError("read stream", func() error {
			messages, err := c.XReadGroup(ctx, rExArgs).Result()
			if err != nil {
				return err
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
					slog.Info("received message", "stream", rExArgs.Streams[0], "consumer group", rExArgs.Group, "message", data)
					*messageChan <- data
					// ack message now that it's being handled
					c.XAck(ctx, rExArgs.Streams[0], rExArgs.Consumer, xMessage.ID)
				}
			}
			return nil
		})
	}
}

// Create consumer group for a stream if it doesn't already exist
func (c *RedisClientWrapper) CreateConsumerGroup(stream string, consumerGroup string) {
	ctx := context.Background()

	retryOnError("create consumer group", func() error {
		statusCmd := c.XGroupCreateMkStream(ctx, stream, consumerGroup, "$")
		if statusCmd.Err() == redis.Nil || statusCmd.Err().Error() == "BUSYGROUP Consumer Group name already exists" {
			return nil
		}
		return statusCmd.Err()
	})
}

// Send Message to Specific stream
func (c *RedisClientWrapper) SendMessageToStream(msg *message.Message, stream streams.RedisStreams) error {
	if err := msg.Validate(); err != nil {
		return err
	}
	slog.Info("sending message", "stream", stream, "message", msg)
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	ctx := context.Background()
	_, err = c.XAdd(ctx, &redis.XAddArgs{
		Stream: stream.String(),
		Values: map[string]interface{}{
			"data": jsonMsg,
		},
	}).Result()
	return err
}

// Helper function to detect if the error is related to a connection issue
func isConnectionError(err error) bool {
	if err != nil {
		if err.Error() == "connection refused" || err.Error() == "no reachable nodes" {
			return true
		}
	}
	return false
}

// helper function that adds retries
func retryOnError(operationName string, operation func() error) {
	maxRetries := 5
	retryDelay := 3 * time.Second
	retries := 0

	// Keep trying until maxRetries is reached
	for {
		err := operation()
		if err != nil {
			retries++

			if retries >= maxRetries {
				// Log after reaching max retries with the operation name
				slog.Error("failed after consecutive retries", "operation", operationName, "retries", retries, "error", err)
				os.Exit(2)
			}
			if isConnectionError(err) {
				slog.Error("operation failed connecting to redis", "operation", operationName, "retries", retries, "error", err)
			} else {
				slog.Error("operation failed, retrying", "operation", operationName, "retries", retries, "error", err)
			}
			time.Sleep(retryDelay)
			continue
		}
		// Reset retries on success
		retries = 0
	}
}
