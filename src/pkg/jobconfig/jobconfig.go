package jobconfig

import (
	"time"

	"github.com/Ryan-Har/chaos-kube/pkg/tasks"
	"github.com/google/uuid"
)

type JobConfig struct {
	ID        uuid.UUID    `json:"id"`
	Name      string       `json:"name"`
	Tasks     []TaskConfig `json:"tasks"`
	StartTime time.Time    `json:"starttime"`
	EndTime   time.Time    `json:"endtime,omitempty"`
	Runoption RunOption    `json:"runoption"`
}

type RunOption string

const (
	RunOptionSpaced RunOption = "Spaced"
	RunOptionRandom RunOption = "Random"
	RunOptionNow    RunOption = "Now"
)

type TaskConfig struct {
	TaskType  tasks.TaskType          `json:"tasktype"`
	Selectors map[SelectorType]string `json:"selectors"`
}

type SelectorType string

const (
	SelectorTypeNamespace  SelectorType = "Namespace"
	SelectorTypeDeployment SelectorType = "Deployment"
	SelectorTypePodName    SelectorType = "Pod"
)
