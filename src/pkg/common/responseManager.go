package common

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/Ryan-Har/chaos-kube/pkg/message"
	"github.com/google/uuid"
)

// Used for managing resonse messages
type RedisResponseManager struct {
	WaitingResponses    map[uuid.UUID]chan<- message.Message
	waitingResponsesMux sync.RWMutex
}

func NewResponseManager() *RedisResponseManager {
	return &RedisResponseManager{
		WaitingResponses: make(map[uuid.UUID]chan<- message.Message),
	}
}

// Method adds an entry to the manager
// THIS WILL NOT WORK IF RESPONSEMANAGER IS DISTRIBUTED, FIND ANOTHER WAY
func (r *RedisResponseManager) Add(id uuid.UUID, c chan<- message.Message) {
	r.waitingResponsesMux.Lock()
	defer r.waitingResponsesMux.Unlock()
	r.WaitingResponses[id] = c
	slog.Info("Item added to response manager", "ID", id)
}

// Define a custom error type for when a channel doesn't exist for a message ID
type ChannelNotFoundError struct {
	ID uuid.UUID
}

// Implement the Error() method for ChannelNotFoundError
func (e *ChannelNotFoundError) Error() string {
	return fmt.Sprintf("unable to send message with ID %v to response manager: no channel exists for this ID", e.ID)
}

// Implement the Is method for ChannelNotFoundError
func (e *ChannelNotFoundError) Is(err error) bool {
	var target *ChannelNotFoundError
	return errors.As(err, &target)
}

// Method sends message to the waiting channel, if it exists
func (r *RedisResponseManager) Send(m *message.Message) error {
	if m.ResponseID == uuid.Nil {
		return &ChannelNotFoundError{ID: uuid.Nil}
	}
	r.waitingResponsesMux.RLock()
	channel, ok := r.WaitingResponses[m.ResponseID]
	r.waitingResponsesMux.RUnlock()
	if !ok {
		return &ChannelNotFoundError{ID: m.ResponseID}
	}
	channel <- *m
	close(channel)
	r.waitingResponsesMux.Lock()
	delete(r.WaitingResponses, m.ResponseID)
	r.waitingResponsesMux.Unlock()
	slog.Debug("Item removed from response manager", "ID", m.ResponseID)
	return nil
}
