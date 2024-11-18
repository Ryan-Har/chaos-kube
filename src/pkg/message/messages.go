package message

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
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
	return nil
}

// Sends the message to a specific stream using the redis client
func (m *Message) SendToRedis(rdb *redis.Client, stream string) error {
	if err := m.Validate(); err != nil {
		return err
	}
	jsonMsg, err := json.Marshal(m)
	if err != nil {
		return err
	}
	ctx := context.Background()
	_, err = rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{
			"data": jsonMsg,
		},
	}).Result()
	return err
}
