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

	// Get the userProcedures
	var userProcedures map[string]interface{}
	
	// If userProcedures exists, unmarshal it
	if userProcsRaw, exists := header["userProcedures"]; exists && userProcsRaw != nil {
		// Convert to JSON first
		userProcsJSON, err := json.Marshal(userProcsRaw)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal existing userProcedures: %w", err)
		}
		
		// Then unmarshal to map
		if err := json.Unmarshal(userProcsJSON, &userProcedures); err != nil {
			return nil, fmt.Errorf("failed to unmarshal existing userProcedures: %w", err)
		}
	} else {
		// If userProcedures doesn't exist, create an empty map
		userProcedures = make(map[string]interface{})
	}

	// Parse the procedures to add
	var proceduresToAdd map[string]interface{}
	if err := json.Unmarshal(procedures, &proceduresToAdd); err != nil {
		return nil, fmt.Errorf("failed to unmarshal procedures to add: %w", err)
	}

	// Add the procedures to userProcedures
	for name, proc := range proceduresToAdd {
		userProcedures[name] = proc
	}

	// Update the header with the modified userProcedures
	header["userProcedures"] = userProcedures
	program["header"] = header

	// Marshal the program back to JSON
	modifiedProgram, err := json.Marshal(program)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal modified program: %w", err)
	}

	return modifiedProgram, nil
}
