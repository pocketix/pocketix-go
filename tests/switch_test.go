package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/statements"
	"github.com/stretchr/testify/assert"
)

func MockSwitchExecute(variableStore *models.VariableStore, s statements.Switch) (bool, any, error) {
	for _, c := range s.Block {
		caseCommand := c.(*statements.Case)
		selectorValue := s.Selector
		caseValue := caseCommand.Value

		if s.SelectorType == "variable" {
			variable, err := variableStore.GetVariable(s.Selector.(string))
			if err != nil {
				return false, nil, err
			}
			selectorValue = variable.Value.Value
		}

		if caseCommand.Type == "variable" {
			variable, err := variableStore.GetVariable(caseValue.(string))
			if err != nil {
				return false, nil, err
			}
			caseValue = variable.Value
		}
		if caseValue == selectorValue {
			caseCommand.Execute(variableStore, nil)
			return true, caseValue, nil
		}
	}
	return true, nil, nil
}

func TestExecuteEmptySwitch(t *testing.T) {
	assert := assert.New(t)

	switchStatement := statements.Switch{
		Id:           "switch",
		Block:        []statements.Statement{},
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

	switchStatement := statements.Switch{
		Id: "switch",
		Block: []statements.Statement{
			&statements.Case{
				Id:    "case",
				Block: []statements.Statement{},
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

func TestExecuteSwitchWithVariableSelector(t *testing.T) {
	assert := assert.New(t)

	variableStore := models.NewVariableStore()
	variable := models.Variable{
		Name:  "selector",
		Value: &models.TreeNode{Type: "string", Value: "value", ResultValue: "value"},
		Type:  "string",
	}
	variableStore.AddVariable(variable)

	switchStatement := statements.Switch{
		Id: "switch",
		Block: []statements.Statement{
			&statements.Case{
				Id:    "case",
				Block: []statements.Statement{},
				Value: "test",
				Type:  "string",
			},
			&statements.Case{
				Id:    "case",
				Block: []statements.Statement{},
				Value: "value",
				Type:  "string",
			},
		},
		SelectorType: "variable",
		Selector:     "selector",
	}

	result, caseValue, err := MockSwitchExecute(variableStore, switchStatement)

	assert.True(result, "Expected true, got false")
	assert.Nil(err, "Expected nil, got %v", err)
	assert.Equal("value", caseValue, "Expected value, got %v", caseValue)
}
