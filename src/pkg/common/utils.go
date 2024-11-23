package common

import (
	"encoding/json"
	"fmt"
)

// Generic unmarshall to attempt to unmarshall any interface into type
func GenericUnmarshal[T any](source interface{}) (T, error) {
	var target T

	// Marshal the source interface into JSON
	data, err := json.Marshal(source)
	if err != nil {
		return target, fmt.Errorf("failed to marshal source: %w", err)
	}

	// Unmarshal the JSON data into the specified type T
	err = json.Unmarshal(data, &target)
	if err != nil {
		return target, fmt.Errorf("failed to unmarshal into target type: %w", err)
	}

	return target, nil
}
