package message_test

import (
	"github.com/Ryan-Har/chaos-kube/pkg/message"
	"github.com/Ryan-Har/chaos-kube/pkg/tasks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Test New Message creation with options
func TestEmptyMessage(t *testing.T) {
	// Test without options
	msg := message.New(
		message.WithTimestamp(time.Time{}),
	)
	msg.ID = uuid.Nil

	assert.NotNil(t, msg.ID, "Message ID should not be nil")
	assert.NotNil(t, msg.Timestamp, "Message Timestamp should not be zero")
	assert.Equal(t, "", msg.Source, "Source should be empty initially")
	assert.Nil(t, msg.Contents, "Contents should be nil initially")

	// Test with options
	customID := uuid.New()
	msgWithID := message.New(message.WithSource("TestSource"), message.WithType(message.JobStartRequest))
	msgWithID.ID = customID // Set custom ID after message creation

	assert.Equal(t, customID, msgWithID.ID, "Message ID should be set via custom ID")
	assert.Equal(t, "TestSource", msgWithID.Source, "Source should be set via option")
	assert.Equal(t, message.JobStartRequest, msgWithID.Type, "Message type should be set via option")
}

// Test Message Validation for valid and invalid messages
func TestMessageValidation(t *testing.T) {
	// Valid message
	validMsg := message.New(message.WithSource("TestSource"), message.WithType(message.ExperimentStopRequest), message.WithContents(&tasks.Task{ID: uuid.New()}))
	err := validMsg.Validate()
	assert.Nil(t, err, "Validation should pass for valid message")

	// Invalid message: Missing required fields
	invalidMsg := message.New(message.WithType(message.JobStartRequest)) // Missing Source
	err = invalidMsg.Validate()
	assert.NotNil(t, err, "Validation should fail for invalid message")
	assert.IsType(t, &message.ContentNotValidError{}, err, "Expected ContentNotValidError")

	// Test with invalid task contents
	invalidTask := &tasks.Task{
		ID: uuid.New(), // This should be zero if not start or final stop
	}
	invalidMsgWithTask := message.New(message.WithSource("TestSource"), message.WithType(message.ExperimentStartRequest), message.WithContents(invalidTask))
	err = invalidMsgWithTask.Validate()
	assert.NotNil(t, err, "Validation should fail with invalid task contents")
	assert.IsType(t, &message.ContentNotValidError{}, err, "Expected ContentNotValidError")

	// Test with valid task contents
	validTask := &tasks.Task{
		ID: uuid.New(),
	}
	validMsgWithTask := message.New(message.WithSource("TestSource"), message.WithType(message.ExperimentStart), message.WithContents(validTask))
	err = validMsgWithTask.Validate()
	assert.Nil(t, err, "Validation should pass with valid task contents")
}

// Test Message Custom Error Handling
func TestMessageNotProcessedError(t *testing.T) {
	msgID := uuid.New()
	msgType := message.JobStartRequest
	reason := "Invalid data"

	err := &message.MessageNotProcessedError{
		ID:     msgID,
		Type:   msgType,
		Reason: reason,
	}

	// Ensure the error message is formatted correctly
	expectedErrorMessage := "unable to process message with ID " + msgID.String() + " and type " + msgType.String() + ". Reason: Invalid data"
	assert.Equal(t, expectedErrorMessage, err.Error(), "Error message should match the expected format")
}

// Test ContentNotValidError
func TestContentNotValidError(t *testing.T) {
	contentType := "Task"
	reasons := []string{"Invalid Task ID", "Missing field"}
	err := &message.ContentNotValidError{
		ContentType: contentType,
		Reasons:     reasons,
	}

	// Ensure the error message is formatted correctly
	expectedErrorMessage := "content not valid for type Task. Reason: Invalid Task ID, Missing field"
	assert.Equal(t, expectedErrorMessage, err.Error(), "Error message should match the expected format")
}
