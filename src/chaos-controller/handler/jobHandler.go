package handler

import (
	"github.com/Ryan-Har/chaos-kube/pkg/common"
	"github.com/google/uuid"
)

type JobHandler interface {
	Begin(redisClient common.RedisClient)
	StartJob(id uuid.UUID, j *JobManifest)
}

type ChaosJobHandler struct {
	//db con TODO
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

func (c *ChaosJobHandler) StartJob(id uuid.UUID, j *JobManifest) {
	// handle startjob steps
	// validate JobManifest
	// update redis data structure for jobstart
	// update db with job start info
	//
}
