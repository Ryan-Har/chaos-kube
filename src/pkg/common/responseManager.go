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
func (r *RedisResponseManager) Add(id uuid.UUID, c chan<- message.Message) {
	r.waitingResponsesMux.Lock()
	defer r.waitingResponsesMux.Unlock()
	r.WaitingResponses[id] = c
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
	slog.Info("response manager handling message", "msg id", m.ID)
	r.waitingResponsesMux.RLock()
	channel, ok := r.WaitingResponses[m.ID]
	r.waitingResponsesMux.RUnlock()
	if !ok {
		return &ChannelNotFoundError{ID: m.ID}
	}
	channel <- *m
	r.waitingResponsesMux.Lock()
	delete(r.WaitingResponses, m.ID)
	r.waitingResponsesMux.Unlock()
	return nil
}
