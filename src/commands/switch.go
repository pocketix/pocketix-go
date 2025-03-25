package commands

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/tree"
)

type Switch struct {
	Id           string
	Block        []Command
	SelectorType string
	Selector     any
}

func (s *Switch) Execute(variableStore *models.VariableStore) (bool, error) {
	services.Logger.Println("Executing switch")
	for _, c := range s.Block {
		caseCommand := c.(*Case)
		selectorValue := s.Selector
		caseValue := caseCommand.Value

		if s.SelectorType == "variable" {
			variable, err := variableStore.GetVariable(s.Selector.(string))
			if err != nil {
				return false, err
			}
			selectorValue = variable.Value
		}

		if caseCommand.Type == "variable" {
			variable, err := variableStore.GetVariable(caseValue.(string))
			if err != nil {
				return false, err
			}
			caseValue = variable.Value
		}
		if caseValue == selectorValue {
			caseCommand.Execute(variableStore)
			return true, nil
		}
	}
	return true, nil
}

func (s *Switch) GetId() string {
	return s.Id
}

func (s *Switch) GetBody() []Command {
	return []Command{}
}

func (s *Switch) GetArguments() *tree.TreeNode {
	return nil
}

func (s *Switch) GetSelector() (any, string) {
	return s.Selector, s.SelectorType
}
