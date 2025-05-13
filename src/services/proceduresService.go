package services

import (
	"encoding/json"
	"fmt"
)

// AddProceduresToProgram adds procedures to the userProcedures section of a program
// Parameters:
//   - programData: the raw JSON program data
//   - procedures: the procedures to add to the program
//
// Returns:
//   - []byte: the modified program data
//   - error: nil if successful, or an error if there was a problem
func AddProceduresToProgram(programData []byte, procedures json.RawMessage) ([]byte, error) {
	// Parse the program
	var program map[string]interface{}
	if err := json.Unmarshal(programData, &program); err != nil {
		return nil, fmt.Errorf("failed to unmarshal program: %w", err)
	}

	// Get the header
	header, ok := program["header"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("header not found or invalid format")
	}

	// Parse the procedures to add
	var proceduresToAdd map[string]interface{}
	if err := json.Unmarshal(procedures, &proceduresToAdd); err != nil {
		return nil, fmt.Errorf("failed to unmarshal procedures to add: %w", err)
	}

	// Update header.userProcedures
	// Get the existing header.userProcedures or create an empty map
	var headerProcedures map[string]interface{}
	if userProcsRaw, exists := header["userProcedures"]; exists && userProcsRaw != nil {
		// Convert to JSON first
		userProcsJSON, err := json.Marshal(userProcsRaw)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal existing header procedures: %w", err)
		}

		// Then unmarshal to map
		if err := json.Unmarshal(userProcsJSON, &headerProcedures); err != nil {
			return nil, fmt.Errorf("failed to unmarshal existing header procedures: %w", err)
		}
	} else {
		headerProcedures = make(map[string]interface{})
	}

	// Add the procedures to header.userProcedures
	for name, proc := range proceduresToAdd {
		headerProcedures[name] = proc
	}

	// Update the header with the modified userProcedures
	header["userProcedures"] = headerProcedures
	program["header"] = header

	// Update the top-level userProcedures field (which is always present in the program structure)
	// Get the existing userProcedures or create an empty map if it doesn't exist
	var topLevelProcedures map[string]interface{}
	if topLevelProcs, exists := program["userProcedures"]; exists && topLevelProcs != nil {
		// Convert to JSON first
		topLevelProcsJSON, err := json.Marshal(topLevelProcs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal existing top-level procedures: %w", err)
		}

		// Then unmarshal to map
		if err := json.Unmarshal(topLevelProcsJSON, &topLevelProcedures); err != nil {
			return nil, fmt.Errorf("failed to unmarshal existing top-level procedures: %w", err)
		}
	} else {
		// Create an empty map if userProcedures doesn't exist
		topLevelProcedures = make(map[string]interface{})
	}

	// Add the procedures to top-level userProcedures
	for name, proc := range proceduresToAdd {
		topLevelProcedures[name] = proc
	}

	// Update the top-level userProcedures
	program["userProcedures"] = topLevelProcedures

	// Marshal the program back to JSON
	modifiedProgram, err := json.Marshal(program)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal modified program: %w", err)
	}

	return modifiedProgram, nil
}
