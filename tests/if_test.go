package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/tree"
	"github.com/stretchr/testify/assert"
)

func TestEvaluateIf_SingleValue(t *testing.T) {
	assert := assert.New(t)

	data := []any{true, false, 1, 0, "true", "false", "", nil, 0.0, 1.0}
	expected := []bool{true, false, true, false, true, true, false, false, false, true}

	for i, value := range data {
		ifStatement := models.If{
			Id: "if",
			Arguments: &tree.TreeNode{
				Value: "boolean_expression", Children: []*tree.TreeNode{
					{Value: value},
				},
			},
			Block: []models.Command{},
		}

		result, err := ifStatement.Execute()

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

	expected := [][]bool{
		{true, false, true, true, false, true, true, false, true, true, true, true, false, false},
		{false, true, false, false, true, false, false, true, false, false, false, false, true, true},
	}

	for i, operator := range operators {
		for j, pair := range data {
			ifStatement := models.If{
				Id: "if",
				Arguments: &tree.TreeNode{
					Value: "boolean_expression", Children: []*tree.TreeNode{
						{Value: operator, Children: []*tree.TreeNode{
							{Value: pair.a},
							{Value: pair.b},
						}},
					},
				},
				Block: []models.Command{},
			}

			result, err := ifStatement.Execute()

			assert.Nil(err, "Error should be nil")
			assert.NotNil(result, "Result should not be nil")
			assert.Equal(expected[i][j], result, "Result of %v %v %v should be %v", pair.a, operator, pair.b, expected[i][j])
		}
	}
}
