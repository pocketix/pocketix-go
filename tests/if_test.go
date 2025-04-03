package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/commands"
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/stretchr/testify/assert"
)

func TestEvaluateIf_SingleValue(t *testing.T) {
	assert := assert.New(t)

	data := []any{true, false, 1, 0, "true", "false", 0.0, 1.0}
	types := []string{"boolean", "boolean", "number", "number", "string", "string", "number", "number"}
	expected := []bool{true, false, true, false, true, false, false, true}

	for i, value := range data {
		ifStatement := commands.If{
			Id: "if",
			Arguments: &models.TreeNode{
				Value: "boolean_expression", Children: []*models.TreeNode{
					{Value: value, ResultValue: value, Type: types[i]},
				},
			},
			Block: []commands.Command{},
		}

		result, err := ifStatement.Execute(nil, nil)

		assert.Nil(err, "Error should be nil")
		assert.NotNil(result, "Result should not be nil")
		assert.Equal(expected[i], result, "Result should be %v", expected[i])
	}
}

func TestEvaluateIf_SimpleCondition(t *testing.T) {
	assert := assert.New(t)

	ifStatement := commands.If{
		Id: "if",
		Arguments: &models.TreeNode{
			Value: "boolean_expression", Children: []*models.TreeNode{
				{Value: "===", Children: []*models.TreeNode{
					{Value: "a", ResultValue: "a", Type: "string"},
					{Value: "b", ResultValue: "b", Type: "string"},
				}},
			},
		},
		Block: []commands.Command{},
	}

	result, err := ifStatement.Execute(nil, nil)
	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.False(result, "Result should be false")
}

func TestEvaluateIfWithVariable(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()
	variable := models.Variable{
		Name:  "foo",
		Type:  "string",
		Value: &models.TreeNode{Type: "string", Value: "abc", ResultValue: "abc"},
	}
	variableStore.AddVariable(variable)

	ifStatement := commands.If{
		Id: "if",
		Arguments: &models.TreeNode{
			Value: "boolean_expression", Children: []*models.TreeNode{
				{Value: "===", Children: []*models.TreeNode{
					{Value: "foo", ResultValue: "foo", Type: "variable"},
					{Value: "abc", ResultValue: "abc", Type: "string"},
				}},
			},
		},
		Block: []commands.Command{},
	}

	result, err := ifStatement.Execute(variableStore, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.True(result, "Result should be true")
}
