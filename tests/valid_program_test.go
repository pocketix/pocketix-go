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

	_, err := parser.ParseWithoutExecuting(program, nil)

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

	_, err := parser.ParseWithoutExecuting(program, nil)

	assert.NotNil(err, "Error should not be nil, but got: %v", err)
}

func TestWithoutHeader(t *testing.T) {
	assert := assert.New(t)

	program := json.RawMessage(`
	{
		"block": []
	}
	`)

	_, err := parser.ParseWithoutExecuting(program, nil)

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

	_, err := parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err := parser.ParseWithoutExecuting(program, nil)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err := parser.ParseWithoutExecuting(program, nil)
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

	_, err = parser.ParseWithoutExecuting(program, nil)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err := parser.ParseWithoutExecuting(program, nil)
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

	_, err = parser.ParseWithoutExecuting(program, nil)
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

	_, err = parser.ParseWithoutExecuting(program, nil)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err := parser.ParseWithoutExecuting(program, nil)
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

	_, err = parser.ParseWithoutExecuting(program, nil)
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

	_, err = parser.ParseWithoutExecuting(program, nil)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
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

	_, err = parser.ParseWithoutExecuting(program, variableStore)
	assert.NotNil(err, "Error should not be nil, but got: %v", err)
}
