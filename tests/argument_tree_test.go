package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/tree"
	"github.com/stretchr/testify/assert"
)

func TestInitTree(t *testing.T) {

}

func TestParseChildren_EmptyArgs(t *testing.T) {
	assert := assert.New(t)
	root := &tree.TreeNode{}
	children := root.ParseChildren([]any{})

	assert.Nil(children, "Children should be nil", children)
	assert.Equal(0, len(children), "Children should be empty")
}

func TestParseChildren_SingleArg(t *testing.T) {
	assert := assert.New(t)

	root := &tree.TreeNode{}
	children := root.ParseChildren([]any{
		map[string]any{"type": "string", "value": "test"},
	})

	assert.NotNil(children, "Children should not be nil")
	assert.Equal(1, len(children), "Children should have one element")

	child := children[0]
	assert.NotNil(child, "Child should not be nil")
	assert.Equal("test", child.Value, "Child value should be 'test'")
}

func TestParseChildren_SimpleCondition(t *testing.T) {
	assert := assert.New(t)

	root := &tree.TreeNode{}
	children := root.ParseChildren([]any{
		map[string]any{"type": "===", "value": []any{
			map[string]any{"type": "string", "value": "test"},
			map[string]any{"type": "string", "value": "test2"},
		}},
	})

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

	root := &tree.TreeNode{}
	factory := tree.NewOperatorFactory()
	result, err, ok := root.EvaluateWithFactory(factory)

	assert.False(ok, "Evaluation should not be possible")
	assert.False(result.(bool), "Result should be false")
	assert.NotNil(err, "Error should not be nil")
}

func TestEvaluate_SingleValue(t *testing.T) {
	assert := assert.New(t)

	root := tree.InitTree("boolean_expresstion", []any{map[string]any{"type": "string", "value": "test"}})

	result, err := root.Evaluate()

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")
}

func TestEvaluate_SimpleCondition(t *testing.T) {
	assert := assert.New(t)

	root := tree.InitTree("boolean_expresstion", []any{
		map[string]any{"type": "===", "value": []any{
			map[string]any{"type": "string", "value": "test"},
			map[string]any{"type": "string", "value": "test"},
		}},
	})

	result, err := root.Evaluate()

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")

	root = tree.InitTree("boolean_expresstion", []any{
		map[string]any{"type": "===", "value": []any{
			map[string]any{"type": "string", "value": "test"},
			map[string]any{"type": "string", "value": "test2"},
		}},
	})

	result, err = root.Evaluate()

	assert.False(result, "Result should be false")
	assert.Nil(err, "Error should be nil")
}

func TestEvaluate_NestedCondition(t *testing.T) {
	assert := assert.New(t)

	root := tree.InitTree("boolean_expresstion", []any{
		map[string]any{"type": "===", "value": []any{
			map[string]any{"type": "boolean", "value": false},
			map[string]any{"type": "===", "value": []any{
				map[string]any{"type": "number", "value": 1},
				map[string]any{"type": "string", "value": "test"},
			}},
		}},
	})

	result, err := root.Evaluate()

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")

	root = tree.InitTree("boolean_expresstion", []any{
		map[string]any{"type": "===", "value": []any{
			map[string]any{"type": "boolean", "value": true},
			map[string]any{"type": "string", "value": []any{
				map[string]any{"type": "number", "value": 1},
				map[string]any{"type": "string", "value": "test"},
			}},
		}},
	})

	result, err = root.Evaluate()

	assert.False(result, "Result should be false")
	assert.Nil(err, "Error should be nil")
}
