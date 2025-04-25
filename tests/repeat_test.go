package tests

import (
	"fmt"
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/statements"
	"github.com/stretchr/testify/assert"
)

// MockRepeatExecute is a mock implementation of the Repeat.Execute method.
// This funcion uses the same logic as the original implementation,
// but it has been modified with iterations to make it testable.
func MockRepeatExecute(r statements.Repeat, variableStore *models.VariableStore) (bool, int, error) {
	iterations := 0

	var count int
	switch r.Count.(type) {
	case float64:
		count = int(r.Count.(float64))
	case int:
		count = r.Count.(int)
	case string:
		variable, err := variableStore.GetVariable(r.Count.(string))
		if err != nil {
			return false, -1, err
		}
		count = int(variable.Value.Value.(float64))
	}

	if count < 0 {
		return false, -1, fmt.Errorf("count cannot be negative")
	}

	for range count {
		iterations++
		result, err := statements.ExecuteStatements(r.Block, nil, nil)
		if err != nil {
			return result, -1, err
		}
	}

	return true, iterations, nil
}

func TestRepeatZeroTimes(t *testing.T) {
	assert := assert.New(t)

	repeatStatement := statements.Repeat{
		Id:    "repeat",
		Count: float64(0),
		Block: []statements.Statement{},
	}

	result, err := repeatStatement.Execute(nil, nil)
	assert.True(result)
	assert.Nil(err)
}

func TestRepeatTenTimes(t *testing.T) {
	assert := assert.New(t)

	repeatStatement := statements.Repeat{
		Id:    "repeat",
		Count: float64(10),
		Block: []statements.Statement{},
	}

	result, iterations, err := MockRepeatExecute(repeatStatement, nil)
	assert.True(result)
	assert.Nil(err)
	assert.Equal(10, iterations)
}

func TestRepeatNegativeCount(t *testing.T) {
	assert := assert.New(t)

	repeatStatement := statements.Repeat{
		Id:    "repeat",
		Count: float64(-1),
		Block: []statements.Statement{},
	}

	result, err := repeatStatement.Execute(nil, nil)
	assert.False(result)
	assert.NotNil(err)
	assert.Equal("count cannot be negative", err.Error())
}

func TestRepeatWithVariable(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()
	variable := models.Variable{
		Name:  "count",
		Type:  "number",
		Value: &models.TreeNode{Type: "number", Value: float64(5), ResultValue: float64(5)},
	}
	variableStore.AddVariable(variable)

	repeatStatement := statements.Repeat{
		Id:        "repeat",
		Count:     "count",
		CountType: "variable",
		Block:     []statements.Statement{},
	}

	result, iterations, err := MockRepeatExecute(repeatStatement, variableStore)

	assert.True(result, "Expected true, got false")
	assert.Nil(err, "Expected nil, got %v", err)
	assert.Equal(5, iterations, "Expected 5, got %v", iterations)
}
