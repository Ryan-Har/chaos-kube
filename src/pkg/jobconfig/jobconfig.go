package jobconfig

import (
	"errors"
	"time"

	"github.com/Ryan-Har/chaos-kube/pkg/tasks"
	"github.com/google/uuid"
)

type JobConfig struct {
	ID        uuid.UUID    `json:"id,omitempty"`
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

func (jc *JobConfig) Validate() error {
	// Check if ID is a zero UUID
	if jc.ID == uuid.Nil {
		return errors.New("job ID cannot be empty")
	}

	// Validate Name
	if jc.Name == "" {
		return errors.New("job name cannot be empty")
	}

	// Validate Tasks
	if len(jc.Tasks) == 0 {
		return errors.New("job must have at least one task")
	}
	for _, task := range jc.Tasks {
		if err := task.Validate(); err != nil {
			return errors.New("task validation failed: " + err.Error())
		}
	}

	// Validate Start and End Times
	if jc.StartTime.IsZero() {
		return errors.New("start time cannot be empty")
	}
	if !jc.EndTime.IsZero() && jc.EndTime.Before(jc.StartTime) {
		return errors.New("end time cannot be before start time")
	}

	return nil
}

func (tc *TaskConfig) Validate() error {
	// Validate TaskType (assuming TaskType has some constraints)
	if tc.TaskType == tasks.TaskUnknown {
		return errors.New("invalid task type")
	}

	// Validate Selectors map
	if len(tc.Selectors) == 0 {
		return errors.New("selectors cannot be empty")
	}
	for key, value := range tc.Selectors {
		if key == "" || value == "" {
			return errors.New("selector keys and values cannot be empty")
		}

	}

	return nil
}
