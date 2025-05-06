package tests

import (
	"encoding/json"
	"testing"

	"github.com/pocketix/pocketix-go/src/services"
	"github.com/stretchr/testify/assert"
)

func TestAddProceduresToProgram(t *testing.T) {
	assert := assert.New(t)

	// Test case 1: Add procedures to a program with existing userProcedures
	programWithProcedures := []byte(`{
		"header": {
			"userVariables": {},
			"userProcedures": {
				"existingProc": [
					{
						"id": "alert",
						"arguments": [
							{
								"type": "str_opt",
								"value": "phone_number"
							},
							{
								"type": "string",
								"value": "123"
							},
							{
								"type": "string",
								"value": "message"
							}
						]
					}
				]
			}
		},
		"block": []
	}`)

	proceduresToAdd := []byte(`{
		"newProc": [
			{
				"id": "alert",
				"arguments": [
					{
						"type": "str_opt",
						"value": "email"
					},
					{
						"type": "string",
						"value": "test@example.com"
					},
					{
						"type": "string",
						"value": "new message"
					}
				]
			}
		]
	}`)

	modifiedProgram, err := services.AddProceduresToProgram(programWithProcedures, proceduresToAdd)
	assert.Nil(err, "Error should be nil, but got: %v", err)
	assert.NotNil(modifiedProgram, "Modified program should not be nil")

	// Verify the modified program contains both procedures
	var program map[string]interface{}
	err = json.Unmarshal(modifiedProgram, &program)
	assert.Nil(err, "Error should be nil when unmarshaling modified program")

	header := program["header"].(map[string]interface{})
	userProcedures := header["userProcedures"].(map[string]interface{})
	
	assert.Contains(userProcedures, "existingProc", "Modified program should contain existing procedure")
	assert.Contains(userProcedures, "newProc", "Modified program should contain new procedure")

	// Test case 2: Add procedures to a program without userProcedures
	programWithoutProcedures := []byte(`{
		"header": {
			"userVariables": {}
		},
		"block": []
	}`)

	modifiedProgram, err = services.AddProceduresToProgram(programWithoutProcedures, proceduresToAdd)
	assert.Nil(err, "Error should be nil, but got: %v", err)
	assert.NotNil(modifiedProgram, "Modified program should not be nil")

	// Verify the modified program contains the new procedure
	err = json.Unmarshal(modifiedProgram, &program)
	assert.Nil(err, "Error should be nil when unmarshaling modified program")

	header = program["header"].(map[string]interface{})
	userProcedures = header["userProcedures"].(map[string]interface{})
	
	assert.Contains(userProcedures, "newProc", "Modified program should contain new procedure")
}
