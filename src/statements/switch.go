package statements

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

type Switch struct {
	Id           string
	Block        []Statement
	SelectorType string
	Selector     any
}

func (s *Switch) Execute(variableStore *models.VariableStore, commandHandlingStore *models.CommandsHandlingStore) (bool, error) {
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
			caseCommand.Execute(variableStore, commandHandlingStore)
			return true, nil
		}
	}
	return true, nil
}

func (s *Switch) GetId() string {
	return s.Id
}

func (s *Switch) GetBody() []Statement {
	return []Statement{}
}

func (s *Switch) GetSelector() (any, string) {
	return s.Selector, s.SelectorType
}

func (s *Switch) Validate(variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore, args ...any) error {
	if s.SelectorType != "variable" && s.SelectorType != "boolean_expression" {
		return fmt.Errorf("invalid selector type: %s", s.SelectorType)
	}

	for _, command := range s.Block {
		caseCommand, ok := command.(*Case)
		if !ok {
			return fmt.Errorf("invalid command in switch block: %T", command)
		}

		// Can only validate variable type, boolean_expression could be validated only at runtime
		if s.SelectorType == "variable" {
			variable, err := variableStore.GetVariable(s.Selector.(string))
			if err != nil {
				return err
			}
			err = caseCommand.Validate(variableStore, referenceValueStore, variable.Type)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
