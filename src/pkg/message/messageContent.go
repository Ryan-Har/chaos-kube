package message

// Generic Interface representing the Message Content Data
// Data depends on the source system and what it produces
// E.g Log Content Data will contain log data
// and Config Content Data will contain config data
type MessageContentData interface {
	IsValidMessageDataContent() bool
}

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
