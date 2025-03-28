package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestMakeVariableStore(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()

	assert.NotNil(variableStore, "Expected not nil, got nil")
	assert.NotNil(variableStore.Variables, "Expected not nil, got nil")
	assert.Zero(len(variableStore.Variables), "Expected 0, got %d", len(variableStore.Variables))
}

func TestAddVariable(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()

	variable := models.Variable{
		Name:  "variable",
		Type:  "string",
		Value: &models.TreeNode{Type: "string", Value: "value", ResultValue: "value"},
	}

	variableStore.AddVariable(variable)

	assert.Equal(len(variableStore.Variables), 1, "Expected 1, got %d", len(variableStore.Variables))
	assert.Equal(variableStore.Variables["variable"], variable, "Expected %v, got %v", variable, variableStore.Variables["variable"])
}

func TestGetVariable(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()

	variable := models.Variable{
		Name:  "variable",
		Type:  "string",
		Value: &models.TreeNode{Type: "string", Value: "value", ResultValue: "value"},
	}

	variableStore.AddVariable(variable)

	newVariable, err := variableStore.GetVariable("variable")

	assert.Nil(err, "Expected nil, got %v", err)
	assert.Equal(newVariable.Value.Value, "value", "Expected value, got %v", newVariable.Value.Value)
	assert.Equal(newVariable.Value.Type, "string", "Expected string, got %v", newVariable.Value.Type)
}

func TestGetVariableNotFound(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()

	value, err := variableStore.GetVariable("variable")

	assert.NotNil(err, "Expected not nil, got nil")
	assert.Nil(value, "Expected nil, got %v", value)
}

func TestSetVariable(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()

	variable := models.Variable{
		Name:  "variable",
		Type:  "string",
		Value: &models.TreeNode{Type: "string", Value: "value", ResultValue: "value"},
	}

	variableStore.AddVariable(variable)

	err := variableStore.SetVariable("variable", "new value", "string")

	setVariable, _ := variableStore.GetVariable("variable")

	assert.Nil(err, "Expected nil, got %v", err)
	assert.Equal(setVariable.Value.Value, "new value", "Expected new value, got %v", setVariable.Value.Value)
	assert.Equal(setVariable.Value.Type, "string", "Expected string, got %v", setVariable.Value.Type)
}

func TestSetVariableNotFound(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()

	err := variableStore.SetVariable("variable", "value", "string")

	assert.NotNil(err, "Expected not nil, got nil")
}

func TestSetVariableTypeMismatch(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()

	variable := models.Variable{
		Name:  "variable",
		Type:  "string",
		Value: &models.TreeNode{Type: "string", Value: "value", ResultValue: "value"},
	}

	variableStore.AddVariable(variable)

	err := variableStore.SetVariable("variable", 1, "number")

	assert.NotNil(err, "Expected not nil, got nil")
}

func TestExpressionVariable(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()

	variable := models.Variable{
		Name: "variable",
		Type: "string",
		Value: &models.TreeNode{Type: "boolean_expression", Children: []*models.TreeNode{
			{Value: "===", Children: []*models.TreeNode{
				{Type: "string", Value: "value", ResultValue: "value"},
				{Type: "string", Value: "value", ResultValue: "value"},
			}},
		}},
	}

	variableStore.AddVariable(variable)
	result, err := variableStore.GetVariable("variable")

	assert.Nil(err, "Expected nil, got %v", err)

	expressionResult, err := result.Value.Evaluate(variableStore)

	assert.Nil(err, "Expected nil, got %v", err)
	assert.True(utils.ToBool(expressionResult), "Expected true, got false")
}

func TestExpressionVariableTypeMismatch(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()

	variable := models.Variable{
		Name: "variable",
		Type: "string",
		Value: &models.TreeNode{Type: "boolean_expression", Children: []*models.TreeNode{
			{Value: "===", Children: []*models.TreeNode{
				{Type: "string", Value: "value", ResultValue: "value"},
				{Type: "number", Value: 1, ResultValue: 1},
			}},
		}},
	}

	variableStore.AddVariable(variable)
	result, err := variableStore.GetVariable("variable")

	assert.Nil(err, "Expected nil, got %v", err)

	_, err = result.Value.Evaluate(variableStore)

	assert.NotNil(err, "Expected not nil, got nil")
}
