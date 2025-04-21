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

	err := variableStore.SetVariable("variable", "new value", "string", nil)

	setVariable, _ := variableStore.GetVariable("variable")

	assert.Nil(err, "Expected nil, got %v", err)
	assert.Equal(setVariable.Value.Value, "new value", "Expected new value, got %v", setVariable.Value.Value)
	assert.Equal(setVariable.Value.Type, "string", "Expected string, got %v", setVariable.Value.Type)
}

func TestSetVariableNotFound(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()

	err := variableStore.SetVariable("variable", "value", "string", nil)

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

	err := variableStore.SetVariable("variable", 1, "number", nil)

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

	expressionResult, err := result.Value.Evaluate(variableStore, nil)
	assert.Nil(err, "Expected nil, got %v", err)

	boolResult, boolErr := utils.ToBool(expressionResult)
	assert.Nil(boolErr, "Expected nil, got %v", boolErr)
	assert.True(boolResult, "Expected true, got false")
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

	_, err = result.Value.Evaluate(variableStore, nil)

	assert.NotNil(err, "Expected not nil, got nil")
}

func TestExpressionVariableNotFound(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()
	referencedValueStore := models.NewReferencedValueStore()

	variable := models.Variable{
		Name: "variable",
		Type: "string",
		Value: &models.TreeNode{Type: "boolean_expression", Children: []*models.TreeNode{
			{Value: "===", Children: []*models.TreeNode{
				{Type: "string", Value: "value", ResultValue: "value"},
				{Type: "variable", Value: "not_found", ResultValue: "not_found"},
			}},
		}},
	}

	variableStore.AddVariable(variable)
	result, err := variableStore.GetVariable("variable")

	assert.Nil(err, "Expected nil, got %v", err)

	_, err = result.Value.Evaluate(variableStore, referencedValueStore)

	assert.NotNil(err, "Expected not nil, got nil")
}

func TestExpressionVariableWithAnotherVariable(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()
	referencedValueStore := models.NewReferencedValueStore()

	fooVar := models.Variable{
		Name:  "foo",
		Type:  "string",
		Value: &models.TreeNode{Type: "string", Value: "value", ResultValue: "value"},
	}

	variable := models.Variable{
		Name: "bar",
		Type: "string",
		Value: &models.TreeNode{Type: "boolean_expression", Children: []*models.TreeNode{
			{Value: "===", Children: []*models.TreeNode{
				{Type: "string", Value: "value", ResultValue: "value"},
				{Type: "variable", Value: "foo", ResultValue: "foo"},
			}},
		}},
	}

	variableStore.AddVariable(variable)
	variableStore.AddVariable(fooVar)

	result, err := variableStore.GetVariable("bar")
	assert.Nil(err, "Expected nil, got %v", err)

	expressionResult, err := result.Value.Evaluate(variableStore, referencedValueStore)
	assert.Nil(err, "Expected nil, got %v", err)

	boolResult, boolErr := utils.ToBool(expressionResult)
	assert.Nil(boolErr, "Expected nil, got %v", boolErr)
	assert.True(boolResult, "Expected true, got false")
}

func TestExpressionVariableWithAnotherVariableNested(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()
	referencedValueStore := models.NewReferencedValueStore()

	fooVar := models.Variable{
		Name:  "foo",
		Type:  "boolean",
		Value: &models.TreeNode{Type: "boolean", Value: true, ResultValue: true},
	}

	barVar := models.Variable{
		Name:  "bar",
		Type:  "number",
		Value: &models.TreeNode{Type: "number", Value: 10, ResultValue: 10},
	}

	variable := models.Variable{
		Name: "var",
		Type: "boolean_expression",
		Value: &models.TreeNode{Type: "boolean_expression", Children: []*models.TreeNode{
			{Value: "===", Children: []*models.TreeNode{
				{Type: "variable", Value: "foo", ResultValue: "foo"},
				{Value: "===", Children: []*models.TreeNode{
					{Type: "number", Value: 10, ResultValue: 10},
					{Type: "variable", Value: "bar", ResultValue: "bar"},
				},
				}},
			}},
		},
	}

	variableStore.AddVariable(variable)
	variableStore.AddVariable(fooVar)
	variableStore.AddVariable(barVar)

	result, err := variableStore.GetVariable("var")
	assert.Nil(err, "Expected nil, got %v", err)

	expressionResult, err := result.Value.Evaluate(variableStore, referencedValueStore)
	assert.Nil(err, "Expected nil, got %v", err)

	boolResult, boolErr := utils.ToBool(expressionResult)
	assert.Nil(boolErr, "Expected nil, got %v", boolErr)
	assert.True(boolResult, "Expected true, got false")
}
