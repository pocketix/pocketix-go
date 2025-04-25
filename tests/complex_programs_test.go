package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/parser"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/statements"
	"github.com/stretchr/testify/assert"
)

func TestExecuteWhileSetVar(t *testing.T) {
	assert := assert.New(t)

	data := services.OpenFile("../programs/complex/prog2.json")
	variableStore := models.NewVariableStore()
	procedureStore := models.NewProcedureStore()
	commandHandlingStore := models.NewCommandsHandlingStore()

	statementAST := make([]statements.Statement, 0)
	err := parser.Parse(data, variableStore, procedureStore, commandHandlingStore, &statements.ASTCollector{Target: &statementAST})
	assert.Nil(err, "Error should be nil, but got: %v", err)
	assert.NotNil(statementAST, "Commands list should not be nil")

	for _, statement := range statementAST {
		_, err := statement.Execute(variableStore, commandHandlingStore)
		assert.Nil(err, "Error should be nil, but got: %v", err)
	}

	variable, err := variableStore.GetVariable("count")
	assert.Nil(err, "Error should be nil, but got: %v", err)
	assert.Equal(float64(5), variable.Value.Value, "Variable value should be 5, but got: %v", variable.Value.Value)
}
