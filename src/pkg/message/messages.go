package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
	Contents   Contents    `json:"contents,omitempty"`
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
func WithContents(contents Contents) MessageOption {
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
// TODO: Add vailidation for contents based on Type, certain types cannot have nil contents
func (m *Message) Validate() error {
	if m.ID == uuid.Nil {
		return errors.New("ID is required")
	}
	if m.Type == 0 {
		return errors.New("type is required")
	}
	if m.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}
	if m.Source == "" {
		return errors.New("source is required")
	}
	// Validate the Contents field, if it exists
	if err := m.Contents.Validate(); err != nil {
		return fmt.Errorf("contents validation failed: %w", err)
	}
	return nil
}

// Custom UnmarshalJSON for Message to handle the dynamic unmarshalling of Contents.Data based on Type
func (m *Message) UnmarshalJSON(data []byte) error {
	// First, unmarshal the message without the Data field
	type Alias Message
	aux := &struct {
		*Alias
		Contents *struct { // Use pointer to allow nil contents
			Status Status          `json:"status"`
			Error  error           `json:"error,omitempty"`
			Data   json.RawMessage `json:"data,omitempty"` // Use RawMessage for deferred unmarshalling
		} `json:"contents,omitempty"` // Handle missing contents gracefully
	}{
		Alias: (*Alias)(m),
	}

	// Unmarshal the JSON into the auxiliary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Contents != nil {
		m.Contents.Status = aux.Contents.Status
		m.Contents.Error = aux.Contents.Error
	}

	switch m.Type {
	// In these cases we expect no content
	case ExperimentStartRequest:
		m.Contents = Contents{}
		return nil
	case ExperimentStart:
		if aux.Contents != nil {
			var data ExperimentStartContentData
			if err := json.Unmarshal(aux.Contents.Data, &data); err != nil {
				return fmt.Errorf("failed to unmarshal ExperimentStartContentData: %w", err)
			}
			m.Contents.Data = data
		}
	case ExperimentStopRequest:
		if aux.Contents != nil {
			var data ExperimentStopRequestContentData
			if err := json.Unmarshal(aux.Contents.Data, &data); err != nil {
				return fmt.Errorf("failed to unmarshal ExperimentStopRequestContentData: %w", err)
			}
			m.Contents.Data = data
		}
	default:
		return fmt.Errorf("unknown message type: %v", m.Type)
	}
	return nil
}
