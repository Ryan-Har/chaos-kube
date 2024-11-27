// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sql

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type JobStatus string

const (
	JobStatusUnknown   JobStatus = "Unknown"
	JobStatusPending   JobStatus = "Pending"
	JobStatusRunning   JobStatus = "Running"
	JobStatusCompleted JobStatus = "Completed"
)

func (e *JobStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = JobStatus(s)
	case string:
		*e = JobStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for JobStatus: %T", src)
	}
	return nil
}

type NullJobStatus struct {
	JobStatus JobStatus
	Valid     bool // Valid is true if JobStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullJobStatus) Scan(value interface{}) error {
	if value == nil {
		ns.JobStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.JobStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullJobStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.JobStatus), nil
}

type MessageType string

const (
	MessageTypeUnknown                MessageType = "Unknown"
	MessageTypeJobStartRequest        MessageType = "JobStartRequest"
	MessageTypeJobStart               MessageType = "JobStart"
	MessageTypeJobStopRequest         MessageType = "JobStopRequest"
	MessageTypeJobStop                MessageType = "JobStop"
	MessageTypeExperimentStartRequest MessageType = "ExperimentStartRequest"
	MessageTypeExperimentStart        MessageType = "ExperimentStart"
	MessageTypeExperimentStopRequest  MessageType = "ExperimentStopRequest"
	MessageTypeExperimentStop         MessageType = "ExperimentStop"
)

func (e *MessageType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MessageType(s)
	case string:
		*e = MessageType(s)
	default:
		return fmt.Errorf("unsupported scan type for MessageType: %T", src)
	}
	return nil
}

type NullMessageType struct {
	MessageType MessageType
	Valid       bool // Valid is true if MessageType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMessageType) Scan(value interface{}) error {
	if value == nil {
		ns.MessageType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MessageType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMessageType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MessageType), nil
}

type TaskStatus string

const (
	TaskStatusUnknown   TaskStatus = "Unknown"
	TaskStatusPending   TaskStatus = "Pending"
	TaskStatusScheduled TaskStatus = "Scheduled"
	TaskStatusRunning   TaskStatus = "Running"
	TaskStatusCompleted TaskStatus = "Completed"
	TaskStatusFailed    TaskStatus = "Failed"
	TaskStatusCancelled TaskStatus = "Cancelled"
	TaskStatusTimedOut  TaskStatus = "TimedOut"
	TaskStatusRetrying  TaskStatus = "Retrying"
	TaskStatusSkipped   TaskStatus = "Skipped"
)

func (e *TaskStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TaskStatus(s)
	case string:
		*e = TaskStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for TaskStatus: %T", src)
	}
	return nil
}

type NullTaskStatus struct {
	TaskStatus TaskStatus
	Valid      bool // Valid is true if TaskStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTaskStatus) Scan(value interface{}) error {
	if value == nil {
		ns.TaskStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TaskStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTaskStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TaskStatus), nil
}

type TaskType string

const (
	TaskTypeUnknown                  TaskType = "Unknown"
	TaskTypeDeletePod                TaskType = "DeletePod"
	TaskTypeEvictPod                 TaskType = "EvictPod"
	TaskTypeTerminatePod             TaskType = "TerminatePod"
	TaskTypePodCrashLoop             TaskType = "PodCrashLoop"
	TaskTypeScaleDeployment          TaskType = "ScaleDeployment"
	TaskTypeUpdateDeployment         TaskType = "UpdateDeployment"
	TaskTypeRollbackDeployment       TaskType = "RollbackDeployment"
	TaskTypeDrainNode                TaskType = "DrainNode"
	TaskTypeTerminateNode            TaskType = "TerminateNode"
	TaskTypeSimulateLatency          TaskType = "SimulateLatency"
	TaskTypeSimulateNetworkLoss      TaskType = "SimulateNetworkLoss"
	TaskTypeSimulateNetworkPartition TaskType = "SimulateNetworkPartition"
	TaskTypeOomKillPod               TaskType = "OomKillPod"
	TaskTypeCPUThrottling            TaskType = "CPUThrottling"
	TaskTypeMemoryPressure           TaskType = "MemoryPressure"
	TaskTypeSimulateAppCrash         TaskType = "SimulateAppCrash"
	TaskTypeSimulateDiskFailure      TaskType = "SimulateDiskFailure"
	TaskTypeRotateSecrets            TaskType = "RotateSecrets"
)

func (e *TaskType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TaskType(s)
	case string:
		*e = TaskType(s)
	default:
		return fmt.Errorf("unsupported scan type for TaskType: %T", src)
	}
	return nil
}

type NullTaskType struct {
	TaskType TaskType
	Valid    bool // Valid is true if TaskType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTaskType) Scan(value interface{}) error {
	if value == nil {
		ns.TaskType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TaskType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTaskType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TaskType), nil
}

type Configuration struct {
	ID        pgtype.UUID
	Name      string
	Key       string
	Value     []byte
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type Job struct {
	ID              pgtype.UUID
	ConfigurationID pgtype.UUID
	Name            string
	Description     pgtype.Text
	StartTime       pgtype.Timestamptz
	EndTime         pgtype.Timestamptz
	Status          NullJobStatus
	CreatedAt       pgtype.Timestamptz
	UpdatedAt       pgtype.Timestamptz
}

type Log struct {
	ID         pgtype.UUID
	JobID      pgtype.UUID
	Timestamp  pgtype.Timestamptz
	LogMessage pgtype.Text
}

type Message struct {
	ID         pgtype.UUID
	ResponseID pgtype.UUID
	Type       NullMessageType
	Timestamp  pgtype.Timestamptz
	Source     string
	Contents   []byte
}

type Task struct {
	ID          pgtype.UUID
	JobID       pgtype.UUID
	Type        TaskType
	Status      NullTaskStatus
	ScheduledAt pgtype.Timestamptz
	Timeout     pgtype.Int4
	Details     []byte
	Results     []byte
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}
