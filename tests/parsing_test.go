package tests

import (
	"encoding/json"
	"testing"

	"github.com/pocketix/pocketix-go/src/commands"
	"github.com/pocketix/pocketix-go/src/parser"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/types"
	"github.com/stretchr/testify/assert"
)

func TestParseEmptyBlock(t *testing.T) {
	assert := assert.New(t)

	program := services.OpenFile("../programs/basic/empty_block.json")
	assert.NotNil(program, "Program should not be nil")

	programResult, err := parser.Parse(program, nil)

	assert.NotNil(programResult, "Command should not be nil")
	assert.Nil(err, "Error should be nil")

	assert.Equal(0, len(programResult.Blocks), "Expected 0 block, got %d", len(programResult.Blocks))
}

func TestParseIfWithoutArguments(t *testing.T) {
	assert := assert.New(t)

	block := types.Block{
		Id:        "if",
		Body:      []types.Block{},
		Arguments: []types.Argument{},
	}

	cmd, err := parser.ParseBlocks(block, nil)

	assert.NotNil(cmd, "Command should be nil")
	assert.Nil(err, "Error should not be nil")
}

func TestParseSimpleIf(t *testing.T) {
	assert := assert.New(t)

	block := types.Block{
		Id: "if",
		Arguments: []types.Argument{
			{
				Type: "boolean_expression",
				Value: json.RawMessage(`[
                        {
                            "type": "boolean",
                            "value": true
                        }
                    ]`),
			},
		},
		Body: []types.Block{},
	}

	cmd, err := parser.ParseBlocks(block, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(cmd, "Command should not be nil")

	ifStatement := cmd.(*commands.If)
	assert.Equal(0, len(ifStatement.Block), "Expected 0 block, got %d", len(ifStatement.Block))

	arguments := ifStatement.GetArguments()
	assert.NotNil(arguments, "Arguments should not be nil")
	assert.Equal(arguments.Type, "boolean_expression", "Expected boolean_expression, got %v", arguments.Value)

	child := arguments.Children[0]
	assert.NotNil(child, "Child should not be nil")
	assert.Equal(child.Value, true, "Expected true, got %v", child.Value)
}

func TestParseIfWithCondition(t *testing.T) {
	assert := assert.New(t)

	block := types.Block{
		Id: "if",
		Arguments: []types.Argument{
			{
				Type: "boolean_expression",
				Value: json.RawMessage(`[
                        {
                            "value": [
                                {
                                    "type": "string",
                                    "value": "a"
                                },
                                {
                                    "type": "string",
                                    "value": "abc"
                                }
                            ], 
							"type": "==="
                        }
                    ]`),
			},
		},
		Body: []types.Block{},
	}

	cmd, err := parser.ParseBlocks(block, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(cmd, "Command should not be nil")

	ifStatement := cmd.(*commands.If)
	arguments := ifStatement.GetArguments()

	assert.NotNil(arguments, "Arguments should not be nil")
	assert.Equal(arguments.Type, "boolean_expression", "Expected boolean_expression, got %v", arguments.Value)

	child := arguments.Children[0]
	assert.NotNil(child, "Child should not be nil")
	assert.Equal(child.Value, "===", "Expected operator ===, got %v", child.Value)

	operand1, operand2 := child.Children[0], child.Children[1]
	assert.NotNil(operand1, "Operand1 should not be nil")
	assert.NotNil(operand2, "Operand2 should not be nil")

	assert.Equal(operand1.Value, "a", "Expected a, got %v", operand1.Value)
	assert.Equal(operand2.Value, "abc", "Expected abc, got %v", operand2.Value)
}

func TestParseIfWithComplexCondition(t *testing.T) {
	assert := assert.New(t)

	block := types.Block{
		Id: "if",
		Arguments: []types.Argument{
			{
				Type: "boolean_expression",
				Value: json.RawMessage(`[
                        {
                            "value": [
                                {
                                    "value": [
                                        {
                                            "type": "number",
                                            "value": 1
                                        },
                                        {
                                            "type": "boolean",
                                            "value": false
                                        }
                                    ],
                                    "type": "==="
                                },
                                {
                                    "type": "boolean",
                                    "value": true
                                }
                            ],
                            "type": "==="
                        }
                    ]`),
			},
		},
		Body: []types.Block{},
	}

	cmd, err := parser.ParseBlocks(block, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(cmd, "Command should not be nil")

	ifStatement := cmd.(*commands.If)
	arguments := ifStatement.GetArguments()

	assert.NotNil(arguments, "Arguments should not be nil")
	assert.Equal(arguments.Type, "boolean_expression", "Expected boolean_expression, got %v", arguments.Value)

	child := arguments.Children[0]
	assert.NotNil(child, "Child should not be nil")
	assert.Equal(child.Value, "===", "Expected operator ===, got %v", child.Value)

	operand1, operand2 := child.Children[0], child.Children[1]
	assert.NotNil(operand1, "Operand1 should not be nil")
	assert.NotNil(operand2, "Operand2 should not be nil")

	assert.Equal(operand1.Value, "===", "Expected operator ===, got %v", operand1.Value)
	assert.Equal(operand2.Value, true, "Expected true, got %v", operand2.Value)

	operand1_1, operand1_2 := operand1.Children[0], operand1.Children[1]

	assert.NotNil(operand1_1, "Operand1_1 should not be nil")
	assert.NotNil(operand1_2, "Operand1_2 should not be nil")

	assert.Equal(operand1_1.Value, float64(1), "Expected 1, got %v", operand1_1.Value)
	assert.Equal(operand1_2.Value, false, "Expected false, got %v", operand1_2.Value)
}

func TestParseElse(t *testing.T) {
	assert := assert.New(t)

	block := types.Block{
		Id:        "else",
		Arguments: []types.Argument{},
		Body:      []types.Block{},
	}

	cmd, err := parser.ParseBlocks(block, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(cmd, "Command should not be nil")

	elseStatement := cmd.(*commands.Else)
	assert.Equal(0, len(elseStatement.Block), "Expected 0 block, got %d", len(elseStatement.Block))

	arguments := elseStatement.GetArguments()
	assert.Nil(arguments, "Arguments should be nil")
}

func TestParseElseif(t *testing.T) {
	assert := assert.New(t)

	block := types.Block{
		Id: "elseif",
		Arguments: []types.Argument{
			{
				Type: "boolean_expression",
				Value: json.RawMessage(`[
                        {
                            "type": "boolean",
                            "value": true
                        }
                    ]`),
			},
		},
		Body: []types.Block{},
	}

	cmd, err := parser.ParseBlocks(block, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(cmd, "Command should not be nil")

	elseifStatement := cmd.(*commands.ElseIf)
	assert.Equal(0, len(elseifStatement.Block), "Expected 0 block, got %d", len(elseifStatement.Block))

	arguments := elseifStatement.GetArguments()
	assert.NotNil(arguments, "Arguments should not be nil")

	child := arguments.Children[0]
	assert.NotNil(child, "Child should not be nil")
	assert.Equal(child.Value, true, "Expected true, got %v", child.Value)
}

func TestParseElseifWithCondition(t *testing.T) {
	assert := assert.New(t)

	block := types.Block{
		Id: "elseif",
		Arguments: []types.Argument{
			{
				Type: "boolean_expression",
				Value: json.RawMessage(`[
						{
							"value": [
								{
									"type": "string",
									"value": "a"
								},
								{
									"type": "string",
									"value": "abc"
								}
							], 
							"type": "==="
						}
					]`),
			},
		},
		Body: []types.Block{},
	}

	cmd, err := parser.ParseBlocks(block, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(cmd, "Command should not be nil")

	elseifStatement := cmd.(*commands.ElseIf)
	arguments := elseifStatement.GetArguments()

	assert.NotNil(arguments, "Arguments should not be nil")
	assert.Equal(arguments.Type, "boolean_expression", "Expected boolean_expression, got %v", arguments.Value)

	child := arguments.Children[0]
	assert.NotNil(child, "Child should not be nil")
	assert.Equal(child.Value, "===", "Expected operator ===, got %v", child.Value)

	operand1, operand2 := child.Children[0], child.Children[1]
	assert.NotNil(operand1, "Operand1 should not be nil")
	assert.NotNil(operand2, "Operand2 should not be nil")

	assert.Equal(operand1.Value, "a", "Expected a, got %v", operand1.Value)
	assert.Equal(operand2.Value, "abc", "Expected abc, got %v", operand2.Value)
}

func TestParseWhile(t *testing.T) {
	assert := assert.New(t)

	block := types.Block{
		Id: "while",
	}

	cmd, err := parser.ParseBlocks(block, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(cmd, "Command should not be nil")

	whileStatement := cmd.(*commands.While)
	assert.Equal(0, len(whileStatement.Block), "Expected 0 block, got %d", len(whileStatement.Block))

	arguments := whileStatement.GetArguments()
	assert.Nil(arguments, "Arguments should be nil")
}

func TestParseWhileWithCondition(t *testing.T) {
	assert := assert.New(t)

	block := types.Block{
		Id: "while",
		Arguments: []types.Argument{
			{
				Type: "boolean_expression",
				Value: json.RawMessage(`[
						{
							"value": [
								{
									"type": "string",
									"value": "a"
								},
								{
									"type": "string",
									"value": "abc"
								}
							], 
							"type": "==="
						}
					]`),
			},
		},
		Body: []types.Block{},
	}

	cmd, err := parser.ParseBlocks(block, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(cmd, "Command should not be nil")

	whileStatement := cmd.(*commands.While)
	arguments := whileStatement.GetArguments()

	assert.NotNil(arguments, "Arguments should not be nil")
	assert.Equal(arguments.Type, "boolean_expression", "Expected boolean_expression, got %v", arguments.Value)

	child := arguments.Children[0]
	assert.NotNil(child, "Child should not be nil")
	assert.Equal(child.Value, "===", "Expected operator ===, got %v", child.Value)

	operand1, operand2 := child.Children[0], child.Children[1]
	assert.NotNil(operand1, "Operand1 should not be nil")
	assert.NotNil(operand2, "Operand2 should not be nil")

	assert.Equal(operand1.Value, "a", "Expected a, got %v", operand1.Value)
	assert.Equal(operand2.Value, "abc", "Expected abc, got %v", operand2.Value)
}

func TestParseRepeat(t *testing.T) {
	assert := assert.New(t)

	block := types.Block{
		Id: "repeat",
		Arguments: []types.Argument{
			{
				Type:  "number",
				Value: json.RawMessage(`10`),
			},
		},
		Body: []types.Block{},
	}

	cmd, err := parser.ParseBlocks(block, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(cmd, "Command should not be nil")

	repeatStatement := cmd.(*commands.Repeat)
	assert.Equal(0, len(repeatStatement.Block), "Expected 0 block, got %d", len(repeatStatement.Block))

	arguments := repeatStatement.GetArguments()
	assert.Nil(arguments, "Arguments should be nil")

	count := repeatStatement.GetCount()
	assert.Equal(count, 10, "Expected 10, got %d", count)
}

func TestParseSetVariable(t *testing.T) {
	assert := assert.New(t)

	block := types.Block{
		Id: "setvar",
		Arguments: []types.Argument{
			{
				Type:  "variable",
				Value: json.RawMessage(`"foo"`),
			},
			{
				Type:  "number",
				Value: json.RawMessage(`10`),
			},
		},
		Body: []types.Block{},
	}

	cmd, err := parser.ParseBlocks(block, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(cmd, "Command should not be nil")

	setVariableStatement := cmd.(*commands.SetVariable)
	assert.Nil(setVariableStatement.GetArguments(), "Arguments should be nil")
	assert.Nil(setVariableStatement.GetBody(), "Body should be nil")

	LVal, LValType := setVariableStatement.GetLVal(), setVariableStatement.GetLValType()
	assert.Equal(LVal, "foo", "Expected foo, got %v", LVal)
	assert.Equal(LValType, "variable", "Expected variable, got %v", LValType)

	RVal, RValType := setVariableStatement.GetRVal(), setVariableStatement.GetRValType()
	assert.Equal(RVal, float64(10), "Expected 10, got %v", RVal)
	assert.Equal(RValType, "number", "Expected number, got %v", RValType)
}

func TestParseSwitch(t *testing.T) {
	assert := assert.New(t)

	block := types.Block{
		Id: "switch",
		Arguments: []types.Argument{
			{
				Type:  "string",
				Value: json.RawMessage(`"foo"`),
			},
		},
		Body: []types.Block{
			{
				Id: "case",
				Arguments: []types.Argument{
					{
						Type:  "string",
						Value: json.RawMessage(`"foo"`),
					},
				},
			},
			{
				Id: "case",
				Arguments: []types.Argument{
					{
						Type:  "string",
						Value: json.RawMessage(`"bar"`),
					},
				},
			},
		},
	}

	cmd, err := parser.ParseBlocks(block, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(cmd, "Command should not be nil")

	switchStatement := cmd.(*commands.Switch)

	assert.Equal(2, len(switchStatement.Block), "Expected 2 block, got %d", len(switchStatement.Block))

	selectorValue, selectorType := switchStatement.GetSelector()

	assert.Equal(selectorValue, "foo", "Selector should be foo, got %v", selectorValue)
	assert.Equal(selectorType, "string", "Selector type should be string, got %v", selectorType)

	case1 := switchStatement.Block[0].(*commands.Case)
	case1Value := case1.GetValue()

	assert.Equal(0, len(case1.Block), "Expected 0 block, got %d", len(case1.Block))
	assert.Equal(case1Value, "foo", "Expected foo, got %v", case1Value)

	case2 := switchStatement.Block[1].(*commands.Case)
	case2Value := case2.GetValue()

	assert.Equal(0, len(case2.Block), "Expected 0 block, got %d", len(case2.Block))
	assert.Equal(case2Value, "bar", "Expected bar, got %v", case2Value)
}

func TestParseSwitchWithVariableSelector(t *testing.T) {
	assert := assert.New(t)

	block := types.Block{
		Id: "switch",
		Arguments: []types.Argument{
			{
				Type:  "variable",
				Value: json.RawMessage(`"foo"`),
			},
		},
		Body: []types.Block{
			{
				Id: "case",
				Arguments: []types.Argument{
					{
						Type:  "string",
						Value: json.RawMessage(`"foo"`),
					},
				},
			},
			{
				Id: "case",
				Arguments: []types.Argument{
					{
						Type:  "string",
						Value: json.RawMessage(`"bar"`),
					},
				},
			},
		},
	}

	cmd, err := parser.ParseBlocks(block, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(cmd, "Command should not be nil")

	switchStatement := cmd.(*commands.Switch)

	assert.Equal(2, len(switchStatement.Block), "Expected 2 block, got %d", len(switchStatement.Block))

	selectorValue, selectorType := switchStatement.GetSelector()

	assert.Equal(selectorValue, "foo", "Selector should be foo, got %v", selectorValue)
	assert.Equal(selectorType, "variable", "Selector type should be variable, got %v", selectorType)

	case1 := switchStatement.Block[0].(*commands.Case)
	case1Value := case1.GetValue()

	assert.Equal(0, len(case1.Block), "Expected 0 block, got %d", len(case1.Block))
	assert.Equal(case1Value, "foo", "Expected foo, got %v", case1Value)

	case2 := switchStatement.Block[1].(*commands.Case)
	case2Value := case2.GetValue()

	assert.Equal(0, len(case2.Block), "Expected 0 block, got %d", len(case2.Block))
	assert.Equal(case2Value, "bar", "Expected bar, got %v", case2Value)
}

func TestParseAlert(t *testing.T) {
	assert := assert.New(t)

	block := types.Block{
		Id: "alert",
		Arguments: []types.Argument{
			{
				Type:  "str_opt",
				Value: json.RawMessage(`"phone_number"`),
			},
			{
				Type:  "string",
				Value: json.RawMessage(`"1234567890"`),
			},
			{
				Type:  "string",
				Value: json.RawMessage(`"Hello, World!"`),
			},
		},
	}

	cmd, err := parser.ParseBlocks(block, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(cmd, "Command should not be nil")

	alertStatement := cmd.(*commands.Alert)

	method := alertStatement.GetMethod()
	assert.Equal(method, "phone_number", "Expected phone_number, got %v", method)

	receiver, receiverType := alertStatement.GetReceiver()
	assert.Equal(receiver, "1234567890", "Expected 1234567890, got %v", receiver)
	assert.Equal(receiverType, "string", "Expected string, got %v", receiverType)

	message, messageType := alertStatement.GetMessage()
	assert.Equal(message, "Hello, World!", "Expected Hello, World!, got %v", message)
	assert.Equal(messageType, "string", "Expected string, got %v", messageType)

}
