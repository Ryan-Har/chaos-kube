package services

import (
	"errors"
	"log/slog"
	"time"

	"github.com/Ryan-Har/chaos-kube/chaos-controller/handler"
	"github.com/Ryan-Har/chaos-kube/pkg/common"
	"github.com/Ryan-Har/chaos-kube/pkg/config"
	"github.com/Ryan-Har/chaos-kube/pkg/message"
	"github.com/go-redis/redis/v8"
)

type ControllerService struct {
	cfg     *config.Config
	Handler *handler.ChaosHandler
}

func NewControllerService(cfg *config.Config, handler *handler.ChaosHandler) *ControllerService {
	return &ControllerService{
		cfg:     cfg,
		Handler: handler,
	}
}

func (s *ControllerService) Start() {
	// Begin background JobHandler
	//go s.Handler.JobHandler.Begin(s.Handler.Redis)
	go s.readStreams()
	for i := 1; i <= 10; i++ {
		go s.Handler.SendNewExperimentStartRequest()
		time.Sleep(3 * time.Second)
	}

	select {}
}

func (s *ControllerService) readStreams() {
	for _, stream := range s.cfg.RedisStreams.ConsumerStreams {
		go s.Handler.Redis.CreateConsumerGroup(stream, s.cfg.RedisStreams.ConsumerGroup)

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

// Processes the message by sending it to the response manager. If nothing is waiting for this response
// then the handler takes it
func (s *ControllerService) ProcessMessage(msg *message.Message) {
	if msg.Source == s.Handler.Source {
		return
	}
	err := s.Handler.ResponseManager.Send(msg)
	if err == nil {
		return
	}
	if !errors.Is(err, &common.ChannelNotFoundError{}) {
		slog.Error("Unknown error when processing message through response manager", "error", err)
		return
	}
	err = s.Handler.Message(msg)
	if err != nil {
		slog.Error(err.Error())
	}
}
