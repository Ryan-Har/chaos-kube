package tasks

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID          uuid.UUID              `json:"id"`
	Type        TaskType               `json:"type,omitempty"`
	NameSpace   string                 `json:"namespace,omitempty"`
	Target      string                 `json:"target,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	ScheduledAt time.Time              `json:"scheduled_at,omitempty"`
	Status      TaskStatus             `json:"status,omitempty"`
}

type TaskStatus string

const (
	StatusUnknown   TaskStatus = "Unknown"
	StatusPending   TaskStatus = "Pending"
	StatusScheduled TaskStatus = "Scheduled"
	StatusRunning   TaskStatus = "Running"
	StatusCompleted TaskStatus = "Completed"
	StatusFailed    TaskStatus = "Failed"
	StatusCanceled  TaskStatus = "Canceled"
	StatusTimedOut  TaskStatus = "TimedOut"
	StatusRetrying  TaskStatus = "Retrying"
	StatusSkipped   TaskStatus = "Skipped"
)

// Task type represents the available interactions that can be done on the cluster
type TaskType string

const (
	TaskUnknown                  TaskType = "Unknown"
	TaskDeletePod                TaskType = "DeletePod"
	TaskEvictPod                 TaskType = "EvictPod"
	TaskTerminatePod             TaskType = "TerminatePod"
	TaskPodCrashLoop             TaskType = "PodCrashLoop"
	TaskScaleDeployment          TaskType = "ScaleDeployment"
	TaskUpdateDeployment         TaskType = "UpdateDeployment"
	TaskRollbackDeployment       TaskType = "RollbackDeployment"
	TaskDrainNode                TaskType = "DrainNode"
	TaskTerminateNode            TaskType = "TerminateNode"
	TaskSimulateLatency          TaskType = "SimulateLatency"
	TaskSimulateNetworkLoss      TaskType = "SimulateNetworkLoss"
	TaskSimulateNetworkPartition TaskType = "SimulateNetworkPartition"
	TaskOomKillPod               TaskType = "OomKillPod"
	TaskCPUThrottling            TaskType = "CPUThrottling"
	TaskMemoryPressure           TaskType = "MemoryPressure"
	TaskSimulateAppCrash         TaskType = "SimulateAppCrash"
	TaskSimulateDiskFailure      TaskType = "SimulateDiskFailure"
	TaskRotateSecrets            TaskType = "RotateSecrets"
)

func (t *Task) Validate() ([]string, bool) {
	var failReasons []string

	if t.ID == uuid.Nil {
		failReasons = append(failReasons, "task UUID: ID must not be the zero value")
	}

	if t.Type == TaskUnknown {
		failReasons = append(failReasons, "task type: must not be unknown")
	}

	if t.Status == StatusUnknown {
		failReasons = append(failReasons, "task status: must not be unknown")
	}

	if t.NameSpace == "" {
		failReasons = append(failReasons, "task namespace: must not be empty")
	}

	if t.Target == "" {
		failReasons = append(failReasons, "task target: must not be empty")
	}

	if len(failReasons) > 0 {
		return failReasons, false
	}

	return failReasons, true
}
