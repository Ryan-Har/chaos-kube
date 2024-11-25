package message

import (
	"errors"
	"fmt"
	"github.com/Ryan-Har/chaos-kube/pkg/tasks"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Stringer interface {
	String() string
}

type Message struct {
	ID         uuid.UUID   `json:"id"`                   //message id, unique to each message
	ResponseID uuid.UUID   `json:"responseid,omitempty"` //message id that is being responded to, if applicaple
	Type       MessageType `json:"type"`                 //message type, enum list
	Timestamp  time.Time   `json:"timestamp"`            //timestamp message was created
	Source     string      `json:"source"`               //the service that created the message
	Contents   interface{} `json:"contents,omitempty"`
}

// holds both the message with ack and nack functions
type MessageWithRedisOperations struct {
	Message Message
	Ack     func() error // Acknowledge the message to remove it from the pending list
	Nack    func() error // Release the message so that another consumer can handle it
}

// MessageType represents the different types of messages that can be send through redis
type MessageType int

// Const representing possible request types
const (
	Unknown MessageType = iota
	JobStartRequest
	JobStart
	JobStopRequest
	JobStop
	ExperimentStartRequest
	ExperimentStart
	ExperimentStopRequest
	ExperimentStop
)

func (m MessageType) String() string {
	msgTypes := [...]string{
		"Unknown",
		"JobStartRequest",
		"JobStart",
		"JobStopRequest",
		"JobStop",
		"ExperimentStartRequest",
		"ExperimentStart",
		"ExperimentStopRequest",
		"ExperimentStop",
	}
	if int(m) < 0 || int(m) >= len(msgTypes) {
		return "Unknown"
	}
	return msgTypes[m]
}

// Define a custom error type for when a message cannot be processed
type MessageNotProcessedError struct {
	ID     uuid.UUID
	Type   MessageType
	Reason string
}

// Implement the Error() method for MessageNotProcessedError
func (e *MessageNotProcessedError) Error() string {
	return fmt.Sprintf("unable to process message with ID %v and type %v. Reason: %s", e.ID, e.Type, e.Reason)
}

// Implement the Is method for MessageNotProcessedError
func (e *MessageNotProcessedError) Is(err error) bool {
	var target *MessageNotProcessedError
	return errors.As(err, &target)
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

// Functional option type that modifies a Message.
type MessageOption func(*Message)

// WithID sets the ID of the message.
// If no ID it provided, it defaults to uuid.New()
func WithResponseID(id uuid.UUID) MessageOption {
	return func(m *Message) {
		m.ResponseID = id
	}
}

// WithType sets the Type of the message
func WithType(msgType MessageType) MessageOption {
	return func(m *Message) {
		m.Type = msgType
	}
}

// WithTimestamp sets the Timestamp of the message.
// If no timestamp is provided, it defaults to time.Now()
func WithTimestamp(timestamp time.Time) MessageOption {
	return func(m *Message) {
		m.Timestamp = timestamp
	}
}

// WithSource sets the Source of the message.
func WithSource(source string) MessageOption {
	return func(m *Message) {
		m.Source = source
	}
}

// WithContents sets the Contents of the message.
func WithContents(contents interface{}) MessageOption {
	return func(m *Message) {
		m.Contents = contents
	}
}

/*
Function to create a new message with options. Applies ID and Timestamp by default if not added.
Add Type message.NilData if not adding contents. Available options:
  - message.WithID // Optional
  - message.WithType
  - message.WithTimestamp // Optional
  - message.WithSource
  - message.WithContents // Optional, not required for all messages
*/
func New(opts ...MessageOption) *Message {
	m := &Message{
		ID:        uuid.New(),
		Timestamp: time.Now(),
	}

	// Apply all options
	for _, opt := range opts {
		opt(m)
	}

	return m
}

// method to validate message is ok to send
func (m *Message) Validate() error {
	var dataType string
	var failReasons []string

	if m.ID == uuid.Nil {
		failReasons = append(failReasons, "message ID is required")
	}
	if m.Type == 0 {
		failReasons = append(failReasons, "message type is unknown")
	}
	if m.Timestamp.IsZero() {
		failReasons = append(failReasons, "message timestamp is required")
	}
	if m.Source == "" {
		failReasons = append(failReasons, "source is required")
	}

	switch data := m.Contents.(type) {
	case tasks.Task:
		dataType = "Task"
		// only ID required unless it is a start request or the final stop
		if m.Type == ExperimentStartRequest || m.Type == ExperimentStop {
			errs, ok := data.Validate()
			if !ok {
				failReasons = append(failReasons, errs...)
			}
			break
		}
		if data.ID == uuid.Nil {
			failReasons = append(failReasons, "task UUID: ID must not be the zero value")
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
