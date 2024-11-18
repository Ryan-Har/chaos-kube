package handler

import (
	"github.com/Ryan-Har/chaos-kube/pkg/common"
	"github.com/Ryan-Har/chaos-kube/pkg/message"
)

type ChaosHandler struct {
	Redis           common.RedisClient
	ResponseManager common.ResponseManager
	JobHandler      JobHandler
}

func NewChaosHandler(redisClient common.RedisClient, responseManager common.ResponseManager) *ChaosHandler {
	return &ChaosHandler{
		redisClient,
		responseManager,
		NewJobHandler(),
	}
}

// Message handler to route any message that is configured
func (c *ChaosHandler) Message(msg *message.Message) error {
	switch msg.Type {
	case message.ExperimentStart:
		//add Job ExperimentStart Info to database
	case message.ExperimentStop:
		//add Job ExperimentStop Info to database
		//Update Job progress (details tbd)
	case message.JobStartRequest:
		// Creates a UUID for the job
		// Call JobHandler passing the job UUID and . JobHandler will keep track of the job, updating the db of progress
		// and keeping track within a redis data structure. This means different chaos-controller instances also have the progress

		//generate JobStart message
	case message.JobStopRequest:
		// Jobstop request could acompany the jobstartuuid
		// calls JobHandler method to cancel the job. Still need to figure out how this would work so that if multiple controllers exist it still works.
		// should receive a response which generates either a error message or success message. Regardless, it should create a redis message response.
	default:
		//respond with error, shouldn't handle this type of message
	}
	return nil
}

// func (c *ChaosHandler) handleExperimentStartMessage(msg *message.Message) error {
// 	return nil
// }

// func (c *ChaosHandler) handleExperimentStopMessage(msg *message.Message) error {
// 	return nil
// }

// func (c *ChaosHandler) handleJobStartRequestMessage(msg *message.Message) error {
// 	return nil
// }

// func (c *ChaosHandler) handleJobStopRequestMessage(msg *message.Message) error {
// 	return nil
// }
