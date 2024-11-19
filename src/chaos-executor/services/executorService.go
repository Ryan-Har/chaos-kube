package services

import (
	"github.com/Ryan-Har/chaos-kube/chaos-executor/handler"
	"github.com/Ryan-Har/chaos-kube/pkg/config"
	"github.com/Ryan-Har/chaos-kube/pkg/message"
	"github.com/go-redis/redis/v8"
	"log/slog"
)

type ExecutorService struct {
	cfg     *config.Config
	Handler *handler.ExecutorHandler
}

func NewExecutorService(cfg *config.Config, handler *handler.ExecutorHandler) *ExecutorService {
	return &ExecutorService{
		cfg,
		handler,
	}
}

func (s *ExecutorService) Start() {
	go s.readStreams()
	select {}
}

func (s *ExecutorService) readStreams() {
	for _, stream := range s.cfg.RedisStreams.ConsumerStreams {
		s.Handler.Redis.CreateConsumerGroup(stream, s.cfg.RedisStreams.ConsumerGroup)

		go func() {
			rExArgs := &redis.XReadGroupArgs{
				Group:    s.cfg.RedisStreams.ConsumerGroup,
				Consumer: s.cfg.RedisStreams.ConsumerGroup + s.cfg.Hostname,
				Streams:  []string{stream, ">"},
				Block:    0,
				Count:    1, //messages at a time
			}

			receive := make(chan message.Message)
			go s.Handler.Redis.ReadStreamToChan(rExArgs, &receive)
			for msg := range receive {
				go s.ProcessMessage(&msg)
			}
		}()
	}
}

func (s *ExecutorService) ProcessMessage(msg *message.Message) {
	if msg.Source == s.Handler.Source {
		return
	}
	err := s.Handler.Message(msg)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
