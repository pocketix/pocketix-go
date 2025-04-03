package commands

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

type SetVariable struct {
	Id       string
	LVal     string
	LValType string
	RVal     any
	RValType string
}

func (s *SetVariable) Execute(variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore) (bool, error) {
	services.Logger.Println("Setting variable", s.LVal)
	err := variableStore.SetVariable(s.LVal, s.RVal, s.RValType, referenceValueStore)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *SetVariable) GetId() string {
	return s.Id
}

func (s *SetVariable) GetBody() []Command {
	return nil
}

func (s *SetVariable) GetArguments() *models.TreeNode {
	return nil
}

func (s *SetVariable) GetLVal() string {
	return s.LVal
}

func (s *SetVariable) GetRVal() any {
	return s.RVal
}

func (s *SetVariable) GetLValType() string {
	return s.LValType
}

func (s *SetVariable) GetRValType() string {
	return s.RValType
}

func (s *SetVariable) Validate(variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore, args ...any) error {
	return nil
}
