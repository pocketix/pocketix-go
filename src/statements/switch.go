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

func (s *Switch) Execute(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore) (bool, error) {
	services.Logger.Println("Executing switch")
	for _, c := range s.Block {
		caseStatement := c.(*Case)
		selectorValue := s.Selector
		caseValue := caseStatement.Value

		if s.SelectorType == "variable" {
			variable, err := variableStore.GetVariable(s.Selector.(string))
			if err != nil {
				return false, err
			}
			selectorValue = variable.Value
		}

		if caseStatement.Type == "variable" {
			variable, err := variableStore.GetVariable(caseValue.(string))
			if err != nil {
				return false, err
			}
			caseValue = variable.Value
		}
		if caseValue == selectorValue {
			caseStatement.Execute(variableStore, referencedValueStore)
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

func (s *Switch) Validate(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, args ...any) error {
	if s.SelectorType != "variable" && s.SelectorType != "boolean_expression" {
		return fmt.Errorf("invalid selector type: %s", s.SelectorType)
	}

	for _, statement := range s.Block {
		caseStatement, ok := statement.(*Case)
		if !ok {
			return fmt.Errorf("invalid statement in switch block: %T", statement)
		}

		// Can only validate variable type, boolean_expression could be validated only at runtime
		if s.SelectorType == "variable" {
			variable, err := variableStore.GetVariable(s.Selector.(string))
			if err != nil {
				return err
			}
			err = caseStatement.Validate(variableStore, referencedValueStore, variable.Type)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
