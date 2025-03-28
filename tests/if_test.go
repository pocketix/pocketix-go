package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/commands"
	"github.com/pocketix/pocketix-go/src/models"
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
			Arguments: &models.TreeNode{
				Value: "boolean_expression", Children: []*models.TreeNode{
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

	type TestCase struct {
		Operator      string
		Pair          Pair
		Types         Pair
		expected      bool
		expectedError bool
	}

	testCases := []TestCase{
		{"===", Pair{true, true}, Pair{"boolean", "boolean"}, true, false},
		{"===", Pair{true, false}, Pair{"boolean", "boolean"}, false, false},
		{"===", Pair{false, false}, Pair{"boolean", "boolean"}, true, false},
		{"===", Pair{1, 1}, Pair{"number", "number"}, true, false},
		{"===", Pair{1, 0}, Pair{"number", "number"}, false, false},
		{"===", Pair{0, 0}, Pair{"number", "number"}, true, false},
		{"===", Pair{"true", "true"}, Pair{"string", "string"}, true, false},
		{"===", Pair{"true", "false"}, Pair{"string", "string"}, false, false},
		{"===", Pair{"", ""}, Pair{"string", "string"}, true, false},
		{"===", Pair{nil, nil}, Pair{"", ""}, true, false},
		{"===", Pair{0.0, 0.0}, Pair{"number", "number"}, true, false},
		{"===", Pair{1.0, 1.0}, Pair{"number", "number"}, true, false},
		{"===", Pair{"", nil}, Pair{"string", ""}, false, true},
		{"===", Pair{"", "true"}, Pair{"string", "string"}, false, false},

		{"!==", Pair{true, true}, Pair{"boolean", "boolean"}, false, false},
		{"!==", Pair{true, false}, Pair{"boolean", "boolean"}, true, false},
		{"!==", Pair{false, false}, Pair{"boolean", "boolean"}, false, false},
		{"!==", Pair{1, 1}, Pair{"number", "number"}, false, false},
		{"!==", Pair{1, 0}, Pair{"number", "number"}, true, false},
		{"!==", Pair{0, 0}, Pair{"number", "number"}, false, false},
		{"!==", Pair{"true", "true"}, Pair{"string", "string"}, false, false},
		{"!==", Pair{"true", "false"}, Pair{"string", "string"}, true, false},
		{"!==", Pair{"", ""}, Pair{"string", "string"}, false, false},
		{"!==", Pair{nil, nil}, Pair{"", ""}, false, false},
		{"!==", Pair{0.0, 0.0}, Pair{"number", "number"}, false, false},
		{"!==", Pair{1.0, 1.0}, Pair{"number", "number"}, false, false},
		{"!==", Pair{"", nil}, Pair{"string", ""}, true, true},
		{"!==", Pair{"", "true"}, Pair{"string", "string"}, true, false},
	}

	for _, testCase := range testCases {
		ifStatement := commands.If{
			Id: "if",
			Arguments: &models.TreeNode{
				Value: "boolean_expression", Children: []*models.TreeNode{
					{Value: testCase.Operator, Children: []*models.TreeNode{
						{Value: testCase.Pair.a, ResultValue: testCase.Pair.a, Type: testCase.Types.a.(string)},
						{Value: testCase.Pair.b, ResultValue: testCase.Pair.b, Type: testCase.Types.b.(string)},
					}},
				},
			},
			Block: []commands.Command{},
		}
		result, err := ifStatement.Execute(nil)
		if testCase.expectedError {
			assert.NotNil(err, "Expected error, got nil")
		} else {
			assert.Nil(err, "Error should be nil")
			assert.NotNil(result, "Result should not be nil")
			assert.Equal(testCase.expected, result, "Result should be %v", testCase.expected)
		}
	}
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

	result, err := ifStatement.Execute(variableStore)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.True(result, "Result should be true")
}
