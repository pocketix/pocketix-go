package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/commands"
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/tree"
	"github.com/stretchr/testify/assert"
)

func TestEvaluateIf_SingleValue(t *testing.T) {
	assert := assert.New(t)

	data := []any{true, false, 1, 0, "true", "false", "", nil, 0.0, 1.0}
	types := []string{"boolean", "boolean", "number", "number", "string", "string", "string", "", "number", "number"}
	expected := []bool{true, false, true, false, true, true, false, false, false, true}

	for i, value := range data {
		ifStatement := commands.If{
			Id: "if",
			Arguments: &tree.TreeNode{
				Value: "boolean_expression", Children: []*tree.TreeNode{
					{Value: value, ResultValue: value, Type: types[i]},
				},
			},
			Block: []commands.Command{},
		}

		result, err := ifStatement.Execute(nil)

		assert.Nil(err, "Error should be nil")
		assert.NotNil(result, "Result should not be nil")
		assert.Equal(expected[i], result, "Result should be %v", expected[i])
	}
}

func TestEvaluateIf_SimpleCondition(t *testing.T) {
	assert := assert.New(t)

	type Pair struct {
		a, b any
	}

	operators := []string{"===", "!=="} // TODO add more operators

	data := []Pair{
		{true, true},
		{true, false},
		{false, false},
		{1, 1},
		{1, 0},
		{0, 0},
		{"true", "true"},
		{"true", "false"},
		{"", ""},
		{nil, nil},
		{0.0, 0.0},
		{1.0, 1.0},
		{"", nil},
		{"", "true"},
	}

	types := []Pair{
		{"boolean", "boolean"},
		{"boolean", "boolean"},
		{"boolean", "boolean"},
		{"number", "number"},
		{"number", "number"},
		{"number", "number"},
		{"string", "string"},
		{"string", "string"},
		{"string", "string"},
		{"", ""},
		{"number", "number"},
		{"number", "number"},
		{"string", ""},
		{"string", "string"},
	}

	expected := [][]bool{
		{true, false, true, true, false, true, true, false, true, true, true, true, false, false},
		{false, true, false, false, true, false, false, true, false, false, false, false, true, true},
	}

	for i, operator := range operators {
		for j, pair := range data {
			ifStatement := commands.If{
				Id: "if",
				Arguments: &tree.TreeNode{
					Value: "boolean_expression", Children: []*tree.TreeNode{
						{Value: operator, Children: []*tree.TreeNode{
							{Value: pair.a, ResultValue: pair.a, Type: types[i].a.(string)},
							{Value: pair.b, ResultValue: pair.b, Type: types[i].b.(string)},
						}},
					},
				},
				Block: []commands.Command{},
			}

			result, err := ifStatement.Execute(nil)

			assert.Nil(err, "Error should be nil")
			assert.NotNil(result, "Result should not be nil")
			assert.Equal(expected[i][j], result, "Result of %v %v %v should be %v", pair.a, operator, pair.b, expected[i][j])
		}
	}
}

func TestEvaluateIfWithVariable(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()
	variable := models.Variable{
		Name:  "foo",
		Type:  "string",
		Value: "abc",
	}
	variableStore.AddVariable(variable)

	ifStatement := commands.If{
		Id: "if",
		Arguments: &tree.TreeNode{
			Value: "boolean_expression", Children: []*tree.TreeNode{
				{Value: "===", Children: []*tree.TreeNode{
					{Value: "foo", ResultValue: "foo", Type: "variable"},
					{Value: "abc", ResultValue: "abc", Type: "string"},
				}},
			},
		},
		Block: []commands.Command{},
	}

	result, err := ifStatement.Execute(variableStore)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.True(result, "Result should be true")
}
