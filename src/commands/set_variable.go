package commands

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/tree"
)

type SetVariable struct {
	Id            string
	VariableToSet models.Variable
}

func (s *SetVariable) Execute(variableStore *models.VariableStore) (bool, error) {
	services.Logger.Println("Setting variable", s.VariableToSet.Name)
	variableStore.SetVariable(s.VariableToSet.Name, s.VariableToSet.Value)
	return true, nil
}

func (s *SetVariable) GetId() string {
	return s.Id
}

func (s *SetVariable) GetBody() []Command {
	return nil
}

func (s *SetVariable) GetArguments() *tree.TreeNode {
	return nil
}

func (s *SetVariable) GetVariableToSet() models.Variable {
	return s.VariableToSet
}
