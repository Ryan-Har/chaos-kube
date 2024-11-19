package message

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidate_ValidContents(t *testing.T) {
	// Create valid ExperimentStartContentData with a valid UUID
	validExperimentID := uuid.New()
	contents := &Contents{
		Status: Success,
		Data:   &ExperimentStartContentData{ExperimentID: validExperimentID},
	}

	err := contents.Validate()

	// Assert that no error is returned for valid contents
	assert.NoError(t, err)
}

func TestValidate_InvalidStatus(t *testing.T) {
	// Create a contents object with an invalid status (not defined in the enum)
	contents := &Contents{
		Status: Status(100), // Invalid status
	}

	err := contents.Validate()

	// Assert that the error is returned for invalid status
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid status in contents")
}

func TestValidate_InvalidExperimentID(t *testing.T) {
	// Create ExperimentStartContentData with nil UUID (invalid ID)
	contents := &Contents{
		Status: Success,
		Data:   &ExperimentStartContentData{ExperimentID: uuid.Nil},
	}

	err := contents.Validate()

	// Assert that the error is returned due to nil UUID
	assert.Error(t, err)

	var contentErr *ContentNotValidError
	assert.True(t, errors.As(err, &contentErr))
	assert.Equal(t, contentErr.ContentType, "ExperimentStartContentData")
	assert.Contains(t, contentErr.Reasons, "Experiment Id is nil")
}

func TestValidate_ExperimentStopContentData_InvalidID(t *testing.T) {
	// Create ExperimentStopContentData with nil UUID
	contents := &Contents{
		Status: Success,
		Data:   &ExperimentStopContentData{ExperimentID: uuid.Nil},
	}

	err := contents.Validate()

	// Assert error is returned due to nil UUID
	assert.Error(t, err)

	var contentErr *ContentNotValidError
	assert.True(t, errors.As(err, &contentErr))
	assert.Equal(t, contentErr.ContentType, "ExperimentStopContentData")
	assert.Contains(t, contentErr.Reasons, "Experiment Id is nil")
}

func TestValidate_UnknownContentType(t *testing.T) {
	// Create contents with a data type that doesn't match any known types
	contents := &Contents{
		Status: Success,
		Data:   "string", // Invalid data type
	}

	err := contents.Validate()

	// Assert error is returned with Unknown content type
	assert.Error(t, err)

	var contentErr *ContentNotValidError
	assert.True(t, errors.As(err, &contentErr))
	assert.Equal(t, contentErr.ContentType, "Unknown")
	assert.Contains(t, contentErr.Reasons, "Unknown Contents Data Type")
}

func TestValidate_NilData(t *testing.T) {
	// Create contents with nil data and a valid status
	contents := &Contents{
		Status: Success,
		Data:   nil, // Nil data
	}

	err := contents.Validate()

	// Assert that no error is returned when data is nil
	assert.NoError(t, err)
}

func TestValidate_InvalidExperimentStopRequestContentData_ID(t *testing.T) {
	// Create ExperimentStopRequestContentData with nil UUID
	contents := &Contents{
		Status: Success,
		Data:   &ExperimentStopRequestContentData{ExperimentID: uuid.Nil},
	}

	err := contents.Validate()

	// Assert error is returned due to nil UUID
	assert.Error(t, err)

	var contentErr *ContentNotValidError
	assert.True(t, errors.As(err, &contentErr))
	assert.Equal(t, contentErr.ContentType, "ExperimentStopRequestContentData")
	assert.Contains(t, contentErr.Reasons, "Experiment Id is nil")
}

func TestStatusString(t *testing.T) {
	// Test that the Status enum correctly converts to a string
	assert.Equal(t, Success.String(), "Success")
	assert.Equal(t, Warn.String(), "Warning")
	assert.Equal(t, Fail.String(), "Fail")
	assert.Equal(t, Cancel.String(), "Cancel")
	assert.Equal(t, Status(100).String(), "Unknown")
}
