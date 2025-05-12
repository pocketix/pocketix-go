package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/statements"
	"github.com/stretchr/testify/assert"
)

// MockWhileExecute is a mock implementation of the While.Execute method.
// This funcion uses the same logic as the original implementation,
// but it has been modified with iterations to make it testable.
func MockWhileExecute(variableStore *models.VariableStore, w statements.While) (bool, int, error) {
	iterations := 0
	result := true

	for {
		// result, err := w.Arguments.Evaluate(variableStore)
		// if err != nil {
		// 	services.Logger.Println("Error executing while arguments", err)
		// 	return false, err
		// }
		if result {
			if _, success, err := statements.ExecuteStatements(w.Block, variableStore, nil, nil); err != nil {
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

	referencedValueStore := models.NewReferencedValueStore()
	whileStatement := statements.While{
		Id: "while",
		Arguments: &models.TreeNode{
			Value: "boolean_expression",
			Children: []*models.TreeNode{
				{Value: false, ResultValue: false, Type: "boolean"},
			},
		},
		Block: []statements.Statement{},
	}

	_, result, err := whileStatement.Execute(nil, referencedValueStore, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.Equal(true, result, "Result should be true")
}

func TestExecuteWhileWithFalseCondition(t *testing.T) {
	assert := assert.New(t)

	whileStatement := statements.While{
		Id: "while",
		Arguments: &models.TreeNode{
			Value: "boolean_expression",
			Children: []*models.TreeNode{
				{Value: "===", Children: []*models.TreeNode{
					{Value: "a", ResultValue: "a"},
					{Value: "ab", ResultValue: "ab"},
				}},
			},
		},
		Block: []statements.Statement{},
	}

	variableStore := models.VariableStore{}
	referencedValueStore := models.NewReferencedValueStore()
	_, result, err := whileStatement.Execute(&variableStore, referencedValueStore, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.Equal(true, result, "Result should be true")
}

func TestExecuteWhileWithIterations(t *testing.T) {
	assert := assert.New(t)

	// Arguments are not used in this test case, mock execute function will be used instead.
	whileStatement := statements.While{
		Id:    "while",
		Block: []statements.Statement{},
	}

	variableStore := models.VariableStore{}
	result, iterations, err := MockWhileExecute(&variableStore, whileStatement)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.Equal(true, result, "Result should be true")
	assert.Equal(10, iterations, "Iterations should be 10")
}

func TestExecuteWhileWithVariable(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()
	variable := models.Variable{
		Name:  "var",
		Type:  "boolean",
		Value: &models.TreeNode{Value: false, Type: "boolean", ResultValue: false},
	}
	variableStore.AddVariable(variable)

	whileStatement := statements.While{
		Id: "while",
		Arguments: &models.TreeNode{
			Value: "boolean_expression", Children: []*models.TreeNode{
				{Value: "var", ResultValue: "var", Type: "variable"},
			},
		},
		Block: []statements.Statement{},
	}

	referencedValueStore := models.NewReferencedValueStore()
	_, result, err := whileStatement.Execute(variableStore, referencedValueStore, nil)

	assert.Nil(err, "Error should be nil")
	assert.NotNil(result, "Result should not be nil")
	assert.Equal(true, result, "Result should be true")
}
