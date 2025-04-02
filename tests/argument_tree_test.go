package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseChildren_EmptyArgs(t *testing.T) {
	assert := assert.New(t)
	root := &models.TreeNode{}
	children, err := root.ParseChildren([]any{}, nil)
	assert.Nil(err, "Error should be nil")

	assert.NotNil(children, "Children should not be nil", children)
	assert.Equal(0, len(children), "Children should be empty")
}

func TestParseChildren_SingleArg(t *testing.T) {
	assert := assert.New(t)

	root := &models.TreeNode{}
	children, err := root.ParseChildren([]any{map[string]any{"type": "string", "value": "test"}}, nil)
	assert.Nil(err, "Error should be nil")

	assert.NotNil(children, "Children should not be nil")
	assert.Equal(1, len(children), "Children should have one element")

	child := children[0]
	assert.NotNil(child, "Child should not be nil")
	assert.Equal("test", child.Value, "Child value should be 'test'")
}

func TestParseChildren_SimpleCondition(t *testing.T) {
	assert := assert.New(t)

	root := &models.TreeNode{}
	children, err := root.ParseChildren([]any{
		map[string]any{"type": "===", "value": []any{
			map[string]any{"type": "string", "value": "test"},
			map[string]any{"type": "string", "value": "test2"},
		}},
	}, nil)
	assert.Nil(err, "Error should be nil")

	assert.NotNil(children, "Children should not be nil")
	assert.Equal(1, len(children), "Children should have one element")

	child := children[0]
	assert.NotNil(child, "Child should not be nil")
	assert.Equal("===", child.Value, "Child value should be '==='")

	assert.NotNil(child.Children, "Child children should not be nil")
	assert.Equal(2, len(child.Children), "Child children should have two elements")

	assert.Equal("test", child.Children[0].Value, "First child value should be 'test'")
	assert.Equal("test2", child.Children[1].Value, "Second child value should be 'test2'")
}

func TestEvaluate_NoChildren(t *testing.T) {
	assert := assert.New(t)

	root := &models.TreeNode{}
	result, err := root.Evaluate(nil)
	_, boolErr := utils.ToBool(result)

	assert.Nil(result, "Result should be nil")
	assert.Nil(err, "Error should be nil")
	assert.NotNil(boolErr, "Bool error should be nil")
}

func TestEvaluate_SingleValue(t *testing.T) {
	assert := assert.New(t)

	root, err := models.InitTree("boolean_expresstion", "", []any{map[string]any{"type": "string", "value": "true"}}, nil)
	assert.Nil(err, "Error should be nil")

	result, err := root.Evaluate(nil)
	boolResult, boolErr := utils.ToBool(result)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.NotNil(boolResult, "Bool result should not be nil")
	assert.Nil(boolErr, "Bool error should be nil")
	assert.True(boolResult, "Bool result should be true")
}

func TestEvaluate_SimpleCondition(t *testing.T) {
	assert := assert.New(t)

	root, err := models.InitTree("boolean_expresstion", "", []any{
		map[string]any{"type": "===", "value": []any{
			map[string]any{"type": "string", "value": "test"},
			map[string]any{"type": "string", "value": "test"},
		}},
	}, nil)
	assert.Nil(err, "Error should be nil")

	result, err := root.Evaluate(nil)
	boolResult, boolErr := utils.ToBool(result)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.NotNil(boolResult, "Bool result should not be nil")
	assert.Nil(boolErr, "Bool error should be nil")
	assert.True(boolResult, "Bool result should be true")

	root, err = models.InitTree("boolean_expresstion", "", []any{
		map[string]any{"type": "===", "value": []any{
			map[string]any{"type": "string", "value": "test"},
			map[string]any{"type": "string", "value": "test2"},
		}},
	}, nil)
	assert.Nil(err, "Error should be nil")

	result, err = root.Evaluate(nil)
	boolResult, boolErr = utils.ToBool(result)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.NotNil(boolResult, "Bool result should not be nil")
	assert.Nil(boolErr, "Bool error should be nil")
	assert.False(boolResult, "Bool result should be false")
}

func TestEvaluate_NestedCondition(t *testing.T) {
	assert := assert.New(t)

	root, err := models.InitTree("boolean_expresstion", "", []any{
		map[string]any{"type": "===", "value": []any{
			map[string]any{"type": "boolean", "value": false},
			map[string]any{"type": "===", "value": []any{
				map[string]any{"type": "string", "value": "testing"},
				map[string]any{"type": "string", "value": "test"},
			}},
		}},
	}, nil)
	assert.Nil(err, "Error should be nil")

	result, err := root.Evaluate(nil)
	boolResult, boolErr := utils.ToBool(result)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.NotNil(boolResult, "Bool result should not be nil")
	assert.Nil(boolErr, "Bool error should be nil")
	assert.True(boolResult, "Bool result should be true")

	root, err = models.InitTree("boolean_expresstion", "", []any{
		map[string]any{"type": "===", "value": []any{
			map[string]any{"type": "boolean", "value": true},
			map[string]any{"type": "===", "value": []any{
				map[string]any{"type": "string", "value": "testing"},
				map[string]any{"type": "string", "value": "test"},
			}},
		}},
	}, nil)
	assert.Nil(err, "Error should be nil")

	result, err = root.Evaluate(nil)
	boolResult, boolErr = utils.ToBool(result)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.NotNil(boolResult, "Bool result should not be nil")
	assert.Nil(boolErr, "Bool error should be nil")
	assert.False(boolResult, "Bool result should be false")
}

func TestEvaluateWithVariable(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()
	variable := models.Variable{
		Name:  "foo",
		Type:  "number",
		Value: &models.TreeNode{Type: "number", Value: 1, ResultValue: 1},
	}
	variableStore.AddVariable(variable)

	root, err := models.InitTree("boolean_expresstion", "", []any{
		map[string]any{"type": "===", "value": []any{
			map[string]any{"type": "variable", "value": "foo"},
			map[string]any{"type": "number", "value": 1},
		}},
	}, variableStore)
	assert.Nil(err, "Error should be nil")

	result, err := root.Evaluate(variableStore)
	boolResult, boolErr := utils.ToBool(result)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.NotNil(boolResult, "Bool result should not be nil")
	assert.Nil(boolErr, "Bool error should be nil")
	assert.True(boolResult, "Bool result should be true")
}
