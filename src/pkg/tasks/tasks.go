package tasks

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID              `json:"id"`
	Type        TaskType               `json:"type,omitempty"`
	NameSpace   string                 `json:"namespace,omitempty"`
	Target      string                 `json:"target,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	ScheduledAt time.Time              `json:"scheduled_at,omitempty"`
	Status      TaskStatus             `json:"status,omitempty"`
	Timeout     int                    `json:"timeout,omitempty"` //timeout, in seconds
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

// Define a custom error type for when a Task has an error
type TaskError struct {
	ID     uuid.UUID
	Type   TaskType
	Reason string
}

// Implement the Error() method for TaskError
func (e *TaskError) Error() string {
	return fmt.Sprintf("unable to process task with ID %v and type %v. Reason: %s", e.ID, e.Type, e.Reason)
}

// Implement the Is method for TaskError
func (e *TaskError) Is(err error) bool {
	var target *TaskError
	return errors.As(err, &target)
}

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

// converts task details into TaskOptions interface
func (t *Task) Options() TaskOptions {
	switch t.Type {
	case TaskDeletePod:
		if len(t.Details) == 0 {
			return defaultDeletePodOptions()
		}
		var gracePeriod int64
		gpVal, ok := t.Details["GracePeriodSeconds"]
		if ok {
			gp, err := xToInt64(gpVal)
			if err != nil {
				slog.Warn("build task options failure, using defaults", "convert error", err.Error())
			} else {
				gracePeriod = gp
			}
		}
		return &DeletePodOptions{
			GracePeriodSeconds: &gracePeriod,
		}
	default:
		return nil
	}
}

// accepts any of int, int32, int64 or string and attempts to convert to int64
func xToInt64(x interface{}) (int64, error) {
	switch v := x.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case string:
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("unable to convert string %s to int64", v)
		}
		return num, nil
	default:
		return 0, errors.ErrUnsupported
	}
}
