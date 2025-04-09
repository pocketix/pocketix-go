package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/parser"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/stretchr/testify/assert"
)

func TestExecuteWhileSetVar(t *testing.T) {
	assert := assert.New(t)

	data := services.OpenFile("../programs/complex/prog2.json")
	variableStore := models.NewVariableStore()

	commandsList, err := parser.Parse(data, variableStore, nil, nil)
	assert.Nil(err, "Error should be nil, but got: %v", err)
	assert.NotNil(commandsList, "Commands list should not be nil")

	for _, command := range commandsList {
		_, err := command.Execute(variableStore, nil)
		assert.Nil(err, "Error should be nil, but got: %v", err)
	}

	variable, err := variableStore.GetVariable("count")
	assert.Nil(err, "Error should be nil, but got: %v", err)
	assert.Equal(float64(5), variable.Value.Value, "Variable value should be 5, but got: %v", variable.Value.Value)
}
