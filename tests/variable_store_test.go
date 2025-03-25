package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
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
		Value: "value",
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
		Value: "value",
	}

	variableStore.AddVariable(variable)

	newVariable, err := variableStore.GetVariable("variable")

	assert.Nil(err, "Expected nil, got %v", err)
	assert.Equal(newVariable.Value, "value", "Expected value, got %v", newVariable.Value)
	assert.Equal(newVariable.Type, "string", "Expected string, got %v", newVariable.Type)
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
		Value: "value",
	}

	variableStore.AddVariable(variable)

	err := variableStore.SetVariable("variable", "new value")

	setVariable, _ := variableStore.GetVariable("variable")

	assert.Nil(err, "Expected nil, got %v", err)
	assert.Equal(setVariable.Value, "new value", "Expected new value, got %v", setVariable.Value)
	assert.Equal(setVariable.Type, "string", "Expected string, got %v", setVariable.Type)
}

func TestSetVariableNotFound(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()

	err := variableStore.SetVariable("variable", "value")

	assert.NotNil(err, "Expected not nil, got nil")
}

func TestSetVariableTypeMismatch(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()

	variable := models.Variable{
		Name:  "variable",
		Type:  "string",
		Value: "value",
	}

	variableStore.AddVariable(variable)

	err := variableStore.SetVariable("variable", 1)

	assert.NotNil(err, "Expected not nil, got nil")
}
