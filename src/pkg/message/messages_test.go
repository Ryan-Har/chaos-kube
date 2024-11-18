package message_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Ryan-Har/chaos-kube/pkg/message"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Used to set specific id on message after it is generated since we cannot generate with a specific id
func setID(id uuid.UUID, msg *message.Message) *message.Message {
	msg.ID = id
	return msg
}

// TestMessageTypeString ensures the correct string representation for MessageType values.
func TestMessageTypeString(t *testing.T) {
	tests := []struct {
		msgType  message.MessageType
		expected string
	}{
		{message.Unknown, "Unknown"},
		{message.JobStartRequest, "JobStartRequest"},
		{message.JobStart, "JobStart"},
		{message.JobStopRequest, "JobStopRequest"},
		{message.JobStop, "JobStop"},
		{message.ExperimentStartRequest, "ExperimentStartRequest"},
		{message.ExperimentStart, "ExperimentStart"},
		{message.ExperimentStopRequest, "ExperimentStopRequest"},
		{message.ExperimentStop, "ExperimentStop"},
		{message.MessageType(-1), "Unknown"},
		{message.MessageType(999), "Unknown"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.msgType.String(), "Unexpected string for MessageType %v", tt.msgType)
	}
}

// TestNewMessageWithDefaults verifies the creation of a new message with default values.
func TestNewMessageWithDefaults(t *testing.T) {
	msg := message.New()

	assert.NotEqual(t, uuid.Nil, msg.ID, "Expected a non-nil UUID for default message ID")
	assert.False(t, msg.Timestamp.IsZero(), "Expected a non-zero timestamp for default message")
	assert.Equal(t, message.Unknown, msg.Type, "Expected default message type to be Unknown")
	assert.Empty(t, msg.Source, "Expected default source to be an empty string")
	assert.Empty(t, msg.Contents, "Expected default contents to be empty")
}

// TestNewMessageWithOptions verifies the creation of a new message with specified options.
func TestNewMessageWithOptions(t *testing.T) {
	msgID := uuid.New()
	timestamp := time.Now()
	contents := message.Contents{}
	msg := message.New(
		message.WithType(message.JobStartRequest),
		message.WithTimestamp(timestamp),
		message.WithSource("test-source"),
		message.WithContents(contents),
	)
	msg.ID = msgID

	assert.Equal(t, msgID, msg.ID, "Expected message ID to be set by option")
	assert.Equal(t, message.JobStartRequest, msg.Type, "Expected message type to be set by option")
	assert.Equal(t, timestamp, msg.Timestamp, "Expected message timestamp to be set by option")
	assert.Equal(t, "test-source", msg.Source, "Expected message source to be set by option")
	assert.Equal(t, contents, msg.Contents, "Expected message contents to be set by option")
}

// TestMessageValidate ensures message validation rules are enforced correctly.
func TestMessageValidate(t *testing.T) {
	validID := uuid.New()
	validTime := time.Now()

	tests := []struct {
		name      string
		msg       *message.Message
		expectErr error
	}{
		{
			name: "Valid Message",
			msg: setID(validID, message.New(
				message.WithType(message.JobStart),
				message.WithTimestamp(validTime),
				message.WithSource("test-source"),
			)),
			expectErr: nil,
		},
		{
			name: "Missing ID",
			msg: setID(uuid.Nil, message.New(
				message.WithType(message.JobStart),
				message.WithTimestamp(validTime),
				message.WithSource("test-source"),
			)),
			expectErr: errors.New("ID is required"),
		},
		{
			name: "Missing Type",
			msg: setID(validID, message.New(
				message.WithType(0), // Unknown type, should trigger validation error
				message.WithTimestamp(validTime),
				message.WithSource("test-source"),
			)),
			expectErr: errors.New("type is required"),
		},
		{
			name: "Missing Timestamp",
			msg: setID(validID, message.New(
				message.WithType(message.JobStart),
				message.WithSource("test-source"),
				message.WithTimestamp(time.Time{}),
			)),
			expectErr: errors.New("timestamp is required"),
		},
		{
			name: "Missing Source",
			msg: setID(validID, message.New(
				message.WithType(message.JobStart),
				message.WithTimestamp(validTime),
				message.WithSource(""),
			)),
			expectErr: errors.New("source is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if tt.expectErr != nil {
				assert.EqualError(t, err, tt.expectErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestSendToRedis verifies that a message can be sent to Redis as expected.
func TestSendToRedis(t *testing.T) {
	rdb, mock := redismock.NewClientMock()

	// Create a valid message to send
	msg := message.New(
		message.WithType(message.JobStart),
		message.WithTimestamp(time.Now()),
		message.WithSource("test-source"),
	)

	// Marshal the message to JSON to set up the expected arguments
	jsonMsg, err := json.Marshal(msg)
	assert.NoError(t, err, "Expected no error when marshalling message")

	// Set up the expectation for XAdd with the stream name and values map
	mock.ExpectXAdd(&redis.XAddArgs{
		Stream: "test-stream",
		Values: map[string]interface{}{
			"data": jsonMsg,
		},
	}).SetVal("message-id")

	// Attempt to send the message to Redis
	err = msg.SendToRedis(rdb, "test-stream")
	assert.NoError(t, err, "Expected no error when sending valid message")

	// Ensure all expectations are met
	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations should be met")
}
