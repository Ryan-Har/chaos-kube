package message

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock implementation of MessageContentData to test IsValidMessageDataContent method.
type MockContentData struct {
	Valid bool
}

func (m MockContentData) IsValidMessageDataContent() bool {
	return m.Valid
}

func TestStatusString(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{Success, "Success"},
		{Warn, "Warning"},
		{Fail, "Fail"},
		{Cancel, "Cancel"},
		// Test for a case outside of defined constants (to verify array bounds)
		{Status(999), "Unknown"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.status.String(), "Expected status string representation to match")
	}
}

func TestContentsStatus(t *testing.T) {
	// Test contents with different statuses
	tests := []struct {
		contents Contents
		expected Status
	}{
		{Contents{Status: Success}, Success},
		{Contents{Status: Warn}, Warn},
		{Contents{Status: Fail}, Fail},
		{Contents{Status: Cancel}, Cancel},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.contents.Status, "Expected contents status to match")
	}
}

func TestContentsError(t *testing.T) {
	errMsg := errors.New("sample error")
	contents := Contents{
		Error: errMsg,
	}

	assert.Equal(t, errMsg, contents.Error, "Expected contents error to match the set error")
}

func TestContentsData(t *testing.T) {
	// Test Contents with different data types
	mockData := MockContentData{Valid: true}
	contents := Contents{
		Data: mockData,
	}

	assert.Equal(t, mockData, contents.Data, "Expected contents data to match the set data")
}

func TestMockContentDataIsValid(t *testing.T) {
	validContent := MockContentData{Valid: true}
	invalidContent := MockContentData{Valid: false}

	assert.True(t, validContent.IsValidMessageDataContent(), "Expected validContent to be valid")
	assert.False(t, invalidContent.IsValidMessageDataContent(), "Expected invalidContent to be invalid")
}

func TestContentsWithMockData(t *testing.T) {
	// Test valid mock data
	validData := MockContentData{Valid: true}
	contentsWithValidData := Contents{
		Data: validData,
	}

	assert.IsType(t, MockContentData{}, contentsWithValidData.Data, "Expected contents data type to be MockContentData")
	assert.True(t, contentsWithValidData.Data.(MockContentData).IsValidMessageDataContent(), "Expected mock data to be valid")

	// Test invalid mock data
	invalidData := MockContentData{Valid: false}
	contentsWithInvalidData := Contents{
		Data: invalidData,
	}

	assert.IsType(t, MockContentData{}, contentsWithInvalidData.Data, "Expected contents data type to be MockContentData")
	assert.False(t, contentsWithInvalidData.Data.(MockContentData).IsValidMessageDataContent(), "Expected mock data to be invalid")
}
