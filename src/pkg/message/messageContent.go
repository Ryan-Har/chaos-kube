package message

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Contents represents the main contents of a Message
type Contents struct {
	Status Status      `json:"status"`
	Error  error       `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"` //data depends on the source system and what it produces
}

type Status int

// Const representing possible status types
const (
	Success Status = iota
	Warn
	Fail
	Cancel
)

func (s Status) String() string {
	statuses := [...]string{"Success", "Warning", "Fail", "Cancel"}
	if int(s) < 0 || int(s) >= len(statuses) {
		return "Unknown"
	}
	return statuses[s]
}

// Define a custom error type for when a message cannot be processed
type ContentNotValidError struct {
	ContentType string
	Reasons     []string
}

func (e *ContentNotValidError) Error() string {
	return fmt.Sprintf("content not valid for type %v. Reason: %s", e.ContentType, strings.Join(e.Reasons, ", "))
}

func (e *ContentNotValidError) Is(err error) bool {
	var target *ContentNotValidError
	return errors.As(err, &target)
}

type ExperimentStartContentData struct {
	ExperimentID uuid.UUID
}

type ExperimentStopContentData struct {
	ExperimentID uuid.UUID
	Progress     string
}

type ExperimentStopRequestContentData struct {
	ExperimentID uuid.UUID
}

// Validate checks the Contents fields for validity
func (c *Contents) Validate() error {
	// Check if Status is within a valid range of values (Success, Warn, Fail, Cancel)

	var dataType string
	var failReasons []string

	switch c.Status {
	case Success, Warn, Fail, Cancel:
	default:
		return errors.New("invalid status in contents")
	}

	//not always required
	if c.Data == nil {
		return nil
	}

	switch data := c.Data.(type) {
	case ExperimentStartContentData:
		dataType = "ExperimentStartContentData"
		if data.ExperimentID == uuid.Nil {
			failReasons = append(failReasons, "Experiment Id is nil")
		}
	case ExperimentStopContentData:
		dataType = "ExperimentStopContentData"
		if data.ExperimentID == uuid.Nil {
			failReasons = append(failReasons, "Experiment Id is nil")
		}
	case ExperimentStopRequestContentData:
		dataType = "ExperimentStopRequestContentData"
		if data.ExperimentID == uuid.Nil {
			failReasons = append(failReasons, "Experiment Id is nil")
		}
	default:
		dataType = "Unknown"
		failReasons = append(failReasons, "Unknown Contents Data Type")
	}

	if len(failReasons) > 0 {
		return &ContentNotValidError{
			ContentType: dataType,
			Reasons:     failReasons,
		}
	}
	return nil
}
