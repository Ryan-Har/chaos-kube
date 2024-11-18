package services

import (
	"github.com/Ryan-Har/chaos-kube/pkg/common"
	"github.com/Ryan-Har/chaos-kube/pkg/config"
	"github.com/Ryan-Har/chaos-kube/pkg/message"
	"github.com/go-redis/redis/v8"
)

type ExecutorService struct {
	cfg         *config.Config
	redisClient common.RedisClient
}

func NewExecutorService(cfg *config.Config, redisClient common.RedisClient) *ExecutorService {
	return &ExecutorService{
		cfg:         cfg,
		redisClient: redisClient,
	}
}

func (s *ExecutorService) Start() {
	go s.readStreams()
	select {}
}

func (s *ExecutorService) readStreams() {
	for _, stream := range s.cfg.RedisStreams.ConsumerStreams {
		s.redisClient.CreateConsumerGroup(stream, s.cfg.RedisStreams.ConsumerGroup)

		go func() {
			rExArgs := &redis.XReadGroupArgs{
				Group:    s.cfg.RedisStreams.ConsumerGroup,
				Consumer: s.cfg.RedisStreams.ConsumerGroup + s.cfg.Hostname,
				Streams:  []string{stream, ">"},
				Block:    0,
				Count:    1, //messages at a time
			}

			receive := make(chan message.Message)
			go s.redisClient.ReadStreamToChan(rExArgs, receive)
			for msg := range receive {
				go s.ProcessMessage(&msg)
			}
		}()
	}
}

func (s *ExecutorService) ProcessMessage(msg *message.Message) {

}
