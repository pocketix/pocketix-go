package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/commands"
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/stretchr/testify/assert"
)

func MockSwitchExecute(variableStore *models.VariableStore, s commands.Switch) (bool, any, error) {
	for _, c := range s.Block {
		caseCommand := c.(*commands.Case)
		selectorValue := s.Selector
		caseValue := caseCommand.Value

		if s.SelectorType == "variable" {
			variable, err := variableStore.GetVariable(s.Selector.(string))
			if err != nil {
				return false, nil, err
			}
			selectorValue = variable.Value
		}

		if caseCommand.Type == "variable" {
			variable, err := variableStore.GetVariable(caseValue.(string))
			if err != nil {
				return false, nil, err
			}
			caseValue = variable.Value
		}
		if caseValue == selectorValue {
			caseCommand.Execute(variableStore)
			return true, caseValue, nil
		}
	}
	return true, nil, nil
}

func TestExecuteEmptySwitch(t *testing.T) {
	assert := assert.New(t)

	switchStatement := commands.Switch{
		Id:           "switch",
		Block:        []commands.Command{},
		SelectorType: "string",
		Selector:     "value",
	}

	result, caseValue, err := MockSwitchExecute(nil, switchStatement)

	assert.True(result, "Expected true, got false")
	assert.Nil(err, "Expected nil, got %v", err)
	assert.Nil(caseValue, "Expected nil, got %v", caseValue)
}

func TestExecuteSwitchWithOneCase(t *testing.T) {
	assert := assert.New(t)

	switchStatement := commands.Switch{
		Id: "switch",
		Block: []commands.Command{
			&commands.Case{
				Id:    "case",
				Block: []commands.Command{},
				Value: "value",
				Type:  "string",
			},
		},
		SelectorType: "string",
		Selector:     "value",
	}

	result, caseValue, err := MockSwitchExecute(nil, switchStatement)

	assert.True(result, "Expected true, got false")
	assert.Nil(err, "Expected nil, got %v", err)
	assert.Equal("value", caseValue, "Expected value, got %v", caseValue)
}
