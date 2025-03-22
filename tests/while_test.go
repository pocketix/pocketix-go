package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/commands"
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/tree"
	"github.com/stretchr/testify/assert"
)

// MockWhileExecute is a mock implementation of the While.Execute method.
// This funcion uses the same logic as the original implementation,
// but it has been modified with iterations to make it testable.
func MockWhileExecute(variableStore *models.VariableStore, w commands.While) (bool, int, error) {
	iterations := 0
	result := true

	for {
		// result, err := w.Arguments.Evaluate(variableStore)
		// if err != nil {
		// 	services.Logger.Println("Error executing while arguments", err)
		// 	return false, err
		// }
		if result {
			if success, err := commands.ExecuteCommands(w.Block, variableStore); err != nil {
				return success, -1, err
			} else if !success {
				return success, -1, nil
			}
			iterations++
			if iterations == 10 {
				break
			}
			// variableStore.SetVariable("foo", "a") // Test setting of a variable
		} else {
			services.Logger.Println("While is false, breaking")
			break
		}
	}

	return true, iterations, nil
}

func TestExecuteWhileFalse(t *testing.T) {
	assert := assert.New(t)

	whileStatement := commands.While{
		Id: "while",
		Arguments: &tree.TreeNode{
			Value: "boolean_expression",
			Children: []*tree.TreeNode{
				{Value: false, ResultValue: false},
			},
		},
		Block: []commands.Command{},
	}

	result, err := whileStatement.Execute(nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.Equal(true, result, "Result should be true")
}

func TestExecuteWhileWithFalseCondition(t *testing.T) {
	assert := assert.New(t)

	whileStatement := commands.While{
		Id: "while",
		Arguments: &tree.TreeNode{
			Value: "boolean_expression",
			Children: []*tree.TreeNode{
				{Value: "===", Children: []*tree.TreeNode{
					{Value: "a", ResultValue: "a"},
					{Value: "b", ResultValue: "b"},
				}},
			},
		},
		Block: []commands.Command{},
	}

	variableStore := models.VariableStore{}
	result, err := whileStatement.Execute(&variableStore)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.Equal(true, result, "Result should be true")
}

func TestExecuteWhileWithIterations(t *testing.T) {
	assert := assert.New(t)

	// Arguments are not used in this test case, mock execute function will be used instead.
	whileStatement := commands.While{
		Id:    "while",
		Block: []commands.Command{},
	}

	variableStore := models.VariableStore{}
	result, iterations, err := MockWhileExecute(&variableStore, whileStatement)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.Equal(true, result, "Result should be true")
	assert.Equal(10, iterations, "Iterations should be 10")
}
