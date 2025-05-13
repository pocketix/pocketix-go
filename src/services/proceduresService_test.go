package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddProceduresToProgram(t *testing.T) {
	// Test case 1: Adding procedures to a program with existing header.userProcedures
	t.Run("Add procedures to program with existing header.userProcedures", func(t *testing.T) {
		// Create a test program with existing header.userProcedures
		programData := []byte(`{
			"header": {
				"userProcedures": {
					"existingProc1": [
						{
							"id": "alert",
							"arguments": [
								{
									"type": "str_opt",
									"value": "phone_number"
								}
							]
						}
					]
				},
				"userVariables": {}
			},
			"block": []
		}`)

		// Create procedures to add
		proceduresToAdd := []byte(`{
			"newProc": [
				{
					"id": "if",
					"block": [],
					"arguments": [
						{
							"type": "boolean_expression",
							"value": [
								{
									"type": "||",
									"value": [
										{
											"type": "boolean",
											"value": true
										},
										{
											"type": "boolean",
											"value": false
										}
									]
								}
							]
						}
					]
				}
			]
		}`)

		// Add procedures to the program
		modifiedProgram, err := AddProceduresToProgram(programData, proceduresToAdd)
		assert.NoError(t, err)

		// Parse the modified program
		var program map[string]interface{}
		err = json.Unmarshal(modifiedProgram, &program)
		assert.NoError(t, err)

		// Check that the header.userProcedures contains both the existing and new procedures
		header, ok := program["header"].(map[string]interface{})
		assert.True(t, ok)
		headerUserProcs, ok := header["userProcedures"].(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, headerUserProcs, "existingProc1")
		assert.Contains(t, headerUserProcs, "newProc")
	})

	// Test case 2: Adding procedures to a program with empty header.userProcedures
	t.Run("Add procedures to program with empty header.userProcedures", func(t *testing.T) {
		// Create a test program with empty header.userProcedures
		programData := []byte(`{
			"header": {
				"userProcedures": {},
				"userVariables": {}
			},
			"block": []
		}`)

		// Create procedures to add
		proceduresToAdd := []byte(`{
			"newProc": [
				{
					"id": "if",
					"block": [],
					"arguments": [
						{
							"type": "boolean_expression",
							"value": [
								{
									"type": "||",
									"value": [
										{
											"type": "boolean",
											"value": true
										},
										{
											"type": "boolean",
											"value": false
										}
									]
								}
							]
						}
					]
				}
			]
		}`)

		// Add procedures to the program
		modifiedProgram, err := AddProceduresToProgram(programData, proceduresToAdd)
		assert.NoError(t, err)

		// Parse the modified program
		var program map[string]interface{}
		err = json.Unmarshal(modifiedProgram, &program)
		assert.NoError(t, err)

		// Check that the header.userProcedures contains the new procedure
		header, ok := program["header"].(map[string]interface{})
		assert.True(t, ok)
		headerUserProcs, ok := header["userProcedures"].(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, headerUserProcs, "newProc")
	})

	// Test case 3: Adding multiple procedures to a program
	t.Run("Add multiple procedures to program", func(t *testing.T) {
		// Create a test program
		programData := []byte(`{
			"header": {
				"userProcedures": {
					"existingProc": [
						{
							"id": "alert",
							"arguments": [
								{
									"type": "str_opt",
									"value": "phone_number"
								}
							]
						}
					]
				},
				"userVariables": {}
			},
			"block": []
		}`)

		// Create multiple procedures to add
		proceduresToAdd := []byte(`{
			"newProc1": [
				{
					"id": "if",
					"block": [],
					"arguments": [
						{
							"type": "boolean_expression",
							"value": [
								{
									"type": "||",
									"value": [
										{
											"type": "boolean",
											"value": true
										},
										{
											"type": "boolean",
											"value": false
										}
									]
								}
							]
						}
					]
				}
			],
			"newProc2": [
				{
					"id": "alert",
					"arguments": [
						{
							"type": "str_opt",
							"value": "phone_number"
						}
					]
				}
			]
		}`)

		// Add procedures to the program
		modifiedProgram, err := AddProceduresToProgram(programData, proceduresToAdd)
		assert.NoError(t, err)

		// Parse the modified program
		var program map[string]interface{}
		err = json.Unmarshal(modifiedProgram, &program)
		assert.NoError(t, err)

		// Check that the header.userProcedures contains all procedures
		header, ok := program["header"].(map[string]interface{})
		assert.True(t, ok)
		headerUserProcs, ok := header["userProcedures"].(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, headerUserProcs, "existingProc")
		assert.Contains(t, headerUserProcs, "newProc1")
		assert.Contains(t, headerUserProcs, "newProc2")
	})
}
