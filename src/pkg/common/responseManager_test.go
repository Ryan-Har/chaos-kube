package common

import (
	"errors"
	"sync"
	"testing"

	"github.com/Ryan-Har/chaos-kube/pkg/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResponseManager(t *testing.T) {
	// Test that NewResponseManager returns a properly initialized RedisResponseManager
	manager := NewResponseManager()
	require.NotNil(t, manager)
	assert.NotNil(t, manager.WaitingResponses)
	assert.Empty(t, manager.WaitingResponses)
}

func TestAdd(t *testing.T) {
	manager := NewResponseManager()

	// Prepare a test message and channel
	id := uuid.New()
	ch := make(chan message.Message, 1)

	// Add the channel to the manager
	manager.Add(id, ch)

	// Check that the channel is correctly added
	manager.waitingResponsesMux.RLock()
	defer manager.waitingResponsesMux.RUnlock()
	_, exists := manager.WaitingResponses[id]
	assert.True(t, exists, "Expected channel to be added to the manager")
}

func TestSend_Success(t *testing.T) {
	manager := NewResponseManager()

	// Prepare a test message and channel
	id := uuid.New()
	ch := make(chan message.Message, 1)

	// Add the channel to the manager
	manager.Add(id, ch)

	// Prepare a test message to send
	msg := &message.Message{ID: id}

	// Send the message
	err := manager.Send(msg)

	// Assert no error occurred
	assert.NoError(t, err)

	// Check if the message was successfully sent to the channel
	select {
	case m := <-ch:
		assert.Equal(t, *msg, m, "Expected message sent to the channel")
	default:
		t.Error("Expected message to be received in the channel")
	}
}

func TestSend_NoChannel(t *testing.T) {
	manager := NewResponseManager()

	// Prepare a test message with an ID that doesn't have a channel
	id := uuid.New()
	msg := message.New(
		message.WithType(message.ExperimentStart),
		message.WithSource("test"),
	)
	msg.ID = id

	// Try sending the message (channel does not exist)
	err := manager.Send(msg)

	// Assert that an error occurred
	assert.Error(t, err)

	// Use errors.Is to check that the error is of type ChannelNotFoundError
	if assert.True(t, errors.Is(err, &ChannelNotFoundError{}), "Expected error to be of type ChannelNotFoundError") {
		var channelErr *ChannelNotFoundError
		if ok := errors.As(err, &channelErr); ok {
			assert.Equal(t, id, channelErr.ID)
		} else {
			// Handle unexpected error type
			t.Fatal("Unexpected error type:", err)
		}
	}
}

func TestAdd_Concurrency(t *testing.T) {
	manager := NewResponseManager()

	// Create multiple channels and add them concurrently
	numChannels := 100
	channels := make([]chan message.Message, numChannels)
	var wg sync.WaitGroup
	for i := 0; i < numChannels; i++ {
		wg.Add(1)
		channels[i] = make(chan message.Message, 1)
		go func(i int) {
			defer wg.Done()
			id := uuid.New()
			manager.Add(id, channels[i])

			// After adding, send a message
			msg := message.New(
				message.WithType(message.ExperimentStart),
				message.WithSource("test"),
			)
			msg.ID = id
			err := manager.Send(msg)

			// Assert no error occurred
			assert.NoError(t, err)

			// Verify the message was sent to the correct channel
			m := <-channels[i]
			assert.Equal(t, *msg, m, "Expected message sent to the correct channel")
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
