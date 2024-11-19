package handler

import (
	"github.com/Ryan-Har/chaos-kube/pkg/common"
	"github.com/google/uuid"
)

type JobHandler interface {
	Begin(redisClient common.RedisClient)
	StartJob(id uuid.UUID, j *JobManifest) error
	StopJob(id uuid.UUID) error
}

type ChaosJobHandler struct {
	// redisclient
	// channel used to talk with backgroundHandler maybe
}

func NewJobHandler() *ChaosJobHandler {
	return &ChaosJobHandler{}
}

// think of and replace properly later
type JobManifest struct {
	Steps map[string]string
}

// Entrypoint, used to begin all background jobs required
func (c *ChaosJobHandler) Begin(redisClient common.RedisClient) {

}

func (c *ChaosJobHandler) StartJob(id uuid.UUID, j *JobManifest) error {
	// handle startjob steps
	// validate JobManifest
	// update redis data structure for jobstart
	//
	return nil
}

func (c *ChaosJobHandler) StopJob(id uuid.UUID) error {
	// Handle stop job steps
	// add message to channel to notify the background handler that no further steps are to be done
	// remove data structure from ongoing jobs redis keystore / datastructure
	//
	return nil
}
