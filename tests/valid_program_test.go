package tests

import (
	"encoding/json"
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/parser"
	"github.com/stretchr/testify/assert"
)

func TestEmptyProgram(t *testing.T) {
	assert := assert.New(t)

	program := json.RawMessage(`
	{
		"header": {
        	"userVariables": {},
        	"userProcedures": {}
    	},
		"block": []
	}
	`)

	err := parser.ParseWithoutExecuting(program, nil, nil, nil)

	assert.Nil(err, "Error should be nil, but got: %v", err)
}

func TestWithoutBlock(t *testing.T) {
	assert := assert.New(t)

	program := json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		}
	}
	`)

	err := parser.ParseWithoutExecuting(program, nil, nil, nil)

	assert.NotNil(err, "Error should not be nil, but got: %v", err)
}

func TestWithoutHeader(t *testing.T) {
	assert := assert.New(t)

	program := json.RawMessage(`
	{
		"block": []
	}
	`)

	err := parser.ParseWithoutExecuting(program, nil, nil, nil)

	assert.NotNil(err, "Error should not be nil, but got: %v", err)
}

func TestValidVariables(t *testing.T) {
	assert := assert.New(t)

	// Test valid variable types
	variableStore := models.NewVariableStore()
	program := json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "string",
					"value": "abc"
				}
			},
			"userProcedures": {}
		},
		"block": []
	}
	`)

	err := parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "number",
					"value": 123
				}
			},
			"userProcedures": {}
		},
		"block": []
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "boolean",
					"value": true
				}
			},
			"userProcedures": {}
		},
		"block": []
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "boolean_expression",
					"value": [
						{
							"type": "===",
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
			},
			"userProcedures": {}
		},
		"block": []
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	// Test unknown variable type
	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "unknown",
					"value": 123
				}
			},
			"userProcedures": {}
		},
		"block": []
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	// Test wrong value for string type
	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "string",
					"value": 123
				}
			},
			"userProcedures": {}
		},
		"block": []
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "string",
					"value": false
				}
			},
			"userProcedures": {}
		},
		"block": []
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	// Test wrong value for number type
	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "number",
					"value": "0"
				}
			},
			"userProcedures": {}
		},
		"block": []
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "number",
					"value": false
				}
			},
			"userProcedures": {}
		},
		"block": []
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	// Test wrong value for boolean type
	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "boolean",
					"value": "0"
				}
			},
			"userProcedures": {}
		},
		"block": []
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "boolean",
					"value": 1
				}
			},
			"userProcedures": {}
		},
		"block": []
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	// Test expression variable with nonexistent variable
	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"bar": {
					"type": "boolean_expression",
					"value": [
						{
							"type": "===",
							"value": [
								{
									"type": "variable",
									"value": "nonexistent"
								},
								{
									"type": "boolean",
									"value": true
								}
							]
						}
					]
				}
			},
			"userProcedures": {}
		},
		"block": []
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	// Test expression variable with wrong type
	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"bar": {
					"type": "boolean_expression",
					"value": [
						{
							"type": "===",
							"value": [
								{
									"type": "boolean",
									"value": true
								},
								{
									"type": "string",
									"value": "test"
								}
							]
						}
					]
				}
			},
			"userProcedures": {}
		},
		"block": []
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)
}

func TestValidIfStatement(t *testing.T) {
	assert := assert.New(t)

	program := json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "if",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "string",
										"value": "a"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
        	}
		]
	}
	`)

	err := parser.ParseWithoutExecuting(program, nil, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	variableStore := models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "string",
					"value": "a"
				}
			},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "if",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "variable",
										"value": "foo"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
        	}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "if",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "variable",
										"value": "foo"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
        	}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	// Test if statement with invalid operand type
	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "if",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "unknown",
										"value": "foo"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
        	}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "number",
					"value": 1
				}
			},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "if",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "variable",
										"value": "foo"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
        	}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "if",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "string",
										"value": "a"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
        	},
			{
				"id": "else",
				"block": []
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "else",
				"block": []
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "elseif",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "string",
										"value": "a"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "if",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "string",
										"value": "a"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
        	},
			{
				"id": "elseif",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "string",
										"value": "a"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should not be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "if",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "string",
										"value": "a"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
        	},
			{
				"id": "elseif",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "string",
										"value": "a"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
			},
			{
				"id": "else",
				"block": []
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "if",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "string",
										"value": "a"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
        	},
			{
				"id": "elseif",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "string",
										"value": "a"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
			},
			{
				"id": "elseif",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "string",
										"value": "a"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
			},
			{
				"id": "else",
				"block": []
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "if",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "string",
										"value": "a"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
        	},
			{
				"id": "elseif",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "number",
										"value": 10
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
			},
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "if",
				"block": [
					{
						"id": "if",
						"block": [],
						"arguments": [
							{
								"type": "boolean_expression",
								"value": [
									{
										"value": [
											{
												"type": "string",
												"value": "b"
											},
											{
												"type": "string",
												"value": "b"
											}
										],
										"type": "==="
									}
								]
							}
						]
					}
				],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "string",
										"value": "a"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
        	}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "if",
				"block": [
					{
						"id": "if",
						"block": [],
						"arguments": [
							{
								"type": "boolean_expression",
								"value": [
									{
										"value": [
											{
												"type": "number",
												"value": 10
											},
											{
												"type": "string",
												"value": "b"
											}
										],
										"type": "==="
									}
								]
							}
						]
					}
				],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "string",
										"value": "a"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
        	}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)
}

func TestValidWhileStatement(t *testing.T) {
	assert := assert.New(t)

	program := json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "while",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "string",
										"value": "a"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
			}
		]
	}
	`)

	err := parser.ParseWithoutExecuting(program, nil, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "while",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "number",
										"value": 2
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, nil, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	variableStore := models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "string",
					"value": "a"
				}
			},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "while",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "variable",
										"value": "foo"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)
}

func TestValidRepeatStatement(t *testing.T) {
	assert := assert.New(t)

	program := json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "repeat",
				"block": [],
				"arguments": [
					{
						"type": "number",
						"value": 10
					}
				]
			}
		]
	}
	`)

	err := parser.ParseWithoutExecuting(program, nil, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "repeat",
				"block": [],
				"arguments": [
					{
						"type": "string",
						"value": "abc"
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, nil, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "repeat",
				"block": [],
				"arguments": [
					{
						"type": "boolean",
						"value": true
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, nil, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	variableStore := models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "number",
					"value": 10
				}
			},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "repeat",
				"block": [],
				"arguments": [
					{
						"type": "variable",
						"value": "foo"
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "string",
					"value": "abc"
				}
			},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "repeat",
				"block": [],
				"arguments": [
					{
						"type": "variable",
						"value": "foo"
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "boolean",
					"value": true
				}
			},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "repeat",
				"block": [],
				"arguments": [
					{
						"type": "variable",
						"value": "foo"
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)
}

func TestValidSwitch(t *testing.T) {
	assert := assert.New(t)

	program := json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "switch",
				"block": [],
				"arguments": [
					{
						"type": "string",
						"value": "abc"
					}
				]
			}
		]
	}
	`)

	err := parser.ParseWithoutExecuting(program, nil, nil, nil)
	assert.NotNil(err, "Error should be nil, but got: %v", err)

	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "switch",
				"block": [],
				"arguments": [
					{
						"type": "boolean",
						"value": true
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, nil, nil, nil)
	assert.NotNil(err, "Error should be nil, but got: %v", err)

	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "switch",
				"block": [],
				"arguments": [
					{
						"type": "number",
						"value": 10
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, nil, nil, nil)
	assert.NotNil(err, "Error should be nil, but got: %v", err)

	variableStore := models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "string",
					"value": "abc"
				}
			},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "switch",
				"block": [],
				"arguments": [
					{
						"type": "variable",
						"value": "foo"
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "string",
					"value": "abc"
				}
			},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "switch",
				"block": [
					{
						"id": "case",
						"block": [],
						"arguments": [
							{
								"type": "string",
								"value": "a"
							}
						]
					}
				],
				"arguments": [
					{
						"type": "variable",
						"value": "foo"
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "string",
					"value": "abc"
				}
			},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "switch",
				"block": [
					{
						"id": "case",
						"block": [],
						"arguments": [
							{
								"type": "string",
								"value": "a"
							}
						]
					},
					{
						"id": "case",
						"block": [],
						"arguments": [
							{
								"type": "string",
								"value": "b"
							}
						]
					}
				],
				"arguments": [
					{
						"type": "variable",
						"value": "foo"
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "number",
					"value": 10
				}
			},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "switch",
				"block": [
					{
						"id": "case",
						"block": [],
						"arguments": [
							{
								"type": "number",
								"value": 10
							}
						]
					}
				],
				"arguments": [
					{
						"type": "variable",
						"value": "foo"
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "string",
					"value": "abc"
				}
			},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "switch",
				"block": [
					{
						"id": "case",
						"block": [],
						"arguments": [
							{
								"type": "boolean",
								"value": true
							}
						]
					}
				],
				"arguments": [
					{
						"type": "variable",
						"value": "foo"
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "string",
					"value": "abc"
				}
			},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "switch",
				"block": [
					{
						"id": "case",
						"block": [],
						"arguments": [
							{
								"type": "number",
								"value": 10
							}
						]
					}
				],
				"arguments": [
					{
						"type": "variable",
						"value": "foo"
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)

	variableStore = models.NewVariableStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {
				"foo": {
					"type": "string",
					"value": "abc"
				}
			},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "switch",
				"block": [
					{
						"id": "case",
						"block": [],
						"arguments": [
							{
								"type": "string",
								"value": "a"
							}
						]
					},
					{
						"id": "case",
						"block": [],
						"arguments": [
							{
								"type": "number",
								"value": 10
							}
						]
					}
				],
				"arguments": [
					{
						"type": "variable",
						"value": "foo"
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, nil)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)
}

func TestValidProgramWithReferencedValue(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()
	referencedValueStore := models.NewReferencedValueStore()
	program := json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "if",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "device_variable",
										"value": "DistanceSensor-1.waterLevel"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
			}
		]
	}
	`)

	err := parser.ParseWithoutExecuting(program, variableStore, nil, referencedValueStore)
	assert.Nil(err, "Error should not be nil, but got: %v", err)

	referencedValues := referencedValueStore.GetReferencedValues()
	assert.Equal(1, len(referencedValues), "Expected 1 referenced value, but got: %d", len(referencedValues))

	referencedValue := referencedValues["DistanceSensor-1"]
	assert.NotNil(referencedValue, "Expected referenced value to be not nil, but got: %v", referencedValue)
	assert.Equal("waterLevel", referencedValue, "Expected referenced value name to be 'waterLevel', but got: %s", referencedValue)

	variableStore = models.NewVariableStore()
	referencedValueStore = models.NewReferencedValueStore()
	program = json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "if",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "device_variable",
										"value": "DistanceSensor-1.waterLevel"
									},
									{
										"type": "string",
										"value": "a"
									}
								],
								"type": "==="
							}
						]
					}
				]
			},
			{
				"id": "if",
				"block": [],
				"arguments": [
					{
						"type": "boolean_expression",
						"value": [
							{
								"value": [
									{
										"type": "device_variable",
										"value": "DistanceSensor-1.waterLevel"
									},
									{
										"type": "number",
										"value": 3
									}
								],
								"type": "==="
							}
						]
					}
				]
			}
		]
	}
	`)

	err = parser.ParseWithoutExecuting(program, variableStore, nil, referencedValueStore)
	assert.Nil(err, "Error should not be nil, but got: %v", err)

	referencedValues = referencedValueStore.GetReferencedValues()
	assert.Equal(1, len(referencedValues), "Expected 1 referenced value, but got: %d", len(referencedValues))
	referencedValue = referencedValues["DistanceSensor-1"]
	assert.NotNil(referencedValue, "Expected referenced value to be not nil, but got: %v", referencedValue)
	assert.Equal("waterLevel", referencedValue, "Expected referenced value name to be 'waterLevel', but got: %s", referencedValue)
}

func TestValidDeviceCommand(t *testing.T) {
	assert := assert.New(t)

	program := json.RawMessage(`
	{
		"header": {
			"userVariables": {},
			"userProcedures": {}
		},
		"block": [
			{
				"id": "OutDoorLight.state",
				"arguments": [
					{
						"type": "str_opt",
						"value": "off"
					}
				]
			}
		]
	}
	`)

	err := parser.ParseWithoutExecuting(program, nil, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)
}
