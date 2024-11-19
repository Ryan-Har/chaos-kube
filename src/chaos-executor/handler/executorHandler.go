package handler

import (
	"github.com/Ryan-Har/chaos-kube/pkg/common"
	"github.com/Ryan-Har/chaos-kube/pkg/message"
	"github.com/Ryan-Har/chaos-kube/pkg/streams"
	"github.com/google/uuid"
)

type ExecutorHandler struct {
	Redis  common.RedisClient
	Source string
}

func NewExecutorHandler(redisClient common.RedisClient, consumerGroup string) *ExecutorHandler {
	return &ExecutorHandler{
		Redis:  redisClient,
		Source: consumerGroup,
	}
}

func (h *ExecutorHandler) Message(msg *message.Message) error {
	switch msg.Type {
	case message.ExperimentStartRequest:
		experimentId := uuid.New()
		// generate uuid for experiment
		// start experiment
		// send jobStart message
		returnMsg := message.New(
			message.WithType(message.ExperimentStart),
			message.WithResponseID(msg.ID),
			message.WithSource(h.Source),
			message.WithContents(message.Contents{
				Status: message.Success,
				Error:  nil,
				Data: &message.ExperimentStartContentData{
					ExperimentID: experimentId,
				},
			}),
		)
		err := h.Redis.SendMessageToStream(returnMsg, streams.ExperimentControl)
		return err
	case message.ExperimentStopRequest:
		experimentStopRequestData, ok := msg.Contents.Data.(message.ExperimentStopContentData)
		if !ok {
			return &message.MessageNotProcessedError{
				ID:     msg.ID,
				Type:   msg.Type,
				Reason: "Unable to extract Experiment Stop Request Data from message",
			}
		}

		// Stop experiment

		// Add jobstop message back to the ExperimentControl stream
		returnMsg := message.New(
			message.WithType(message.ExperimentStop),
			message.WithResponseID(msg.ID),
			message.WithSource(h.Source),
			message.WithContents(message.Contents{
				Status: message.Success,
				Error:  nil,
				Data: &message.ExperimentStopContentData{
					ExperimentID: experimentStopRequestData.ExperimentID,
					Progress:     "success",
				},
			}),
		)
		h.Redis.SendMessageToStream(returnMsg, streams.ExperimentControl)
		return nil
	default:
		return &message.MessageNotProcessedError{
			ID:     msg.ID,
			Type:   msg.Type,
			Reason: "ExecutorHandler not configured to handle type",
		}
	}
}
