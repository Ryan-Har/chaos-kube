package handler

import (
	"log/slog"

	"github.com/Ryan-Har/chaos-kube/pkg/common"
	"github.com/Ryan-Har/chaos-kube/pkg/message"
	"github.com/Ryan-Har/chaos-kube/pkg/streams"
)

type ChaosHandler struct {
	Redis           common.RedisClient
	ResponseManager common.ResponseManager
	JobHandler      JobHandler
	Source          string
}

func NewChaosHandler(redisClient common.RedisClient, responseManager common.ResponseManager, ConsumerGroup string) *ChaosHandler {
	return &ChaosHandler{
		Redis:           redisClient,
		ResponseManager: responseManager,
		JobHandler:      NewJobHandler(),
		Source:          ConsumerGroup,
	}
}

// Message handler to route any message that is configured
func (c *ChaosHandler) Message(msg *message.Message) error {
	switch msg.Type {
	case message.ExperimentStartRequest:
		//add Job ExperimentStart Info to database
		return nil
	case message.ExperimentStopRequest:
		//add Job ExperimentStop Info to database
		//Update Job progress (details tbd)
		return nil
	case message.JobStartRequest:
		// Creates a UUID for the job
		// Call JobHandler passing the job UUID and . JobHandler will keep track of the job, updating the db of progress
		// and keeping track within a redis data structure. This means different chaos-controller instances also have the progress

		//generate JobStart message
		return nil
	case message.JobStopRequest:
		// Jobstop request could acompany the jobstartuuid
		// calls JobHandler method to cancel the job. Still need to figure out how this would work so that if multiple controllers exist it still works.
		// should receive a response which generates either a error message or success message. Regardless, it should create a redis message response.
		return nil
	default:
		return &message.MessageNotProcessedError{
			ID:   msg.ID,
			Type: msg.Type,
		}
	}
}

// Used for testing
func (c *ChaosHandler) SendNewExperimentStartRequest() {
	slog.Info("sending ExperimentStartRequest")
	msg := message.New(
		message.WithType(message.ExperimentStartRequest),
		message.WithSource(c.Source),
	)
	returnChan := make(chan message.Message, 1)
	c.ResponseManager.Add(msg.ID, returnChan)
	if err := c.Redis.SendMessageToStream(msg, streams.ExperimentControl); err != nil {
		slog.Error("sending message to redis", "err", err)
	}
	response, ok := <-returnChan
	if ok {
		slog.Info("response received", "response", response)
	} else {
		slog.Info("attempted to receive from channel which is not ok")
	}
}
