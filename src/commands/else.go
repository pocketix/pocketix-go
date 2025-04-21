package commands

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

type Else struct {
	Id    string
	Block []Command
}

func (e *Else) Execute(variableStore *models.VariableStore, commandHandlingStore *models.CommandsHandlingStore) (bool, error) {
	services.Logger.Println("Executing else")
	for _, cmd := range e.Block {
		if success, err := cmd.Execute(variableStore, commandHandlingStore); err != nil {
			return false, err
		} else if success {
			return true, nil
		}
	}
	return true, nil
}

func (e *Else) GetId() string {
	return e.Id
}

func (e *Else) GetBody() []Command {
	return e.Block
}

func (e *Else) GetArguments() *models.TreeNode {
	return nil
}

func (e *Else) Validate(variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore, args ...any) error {
	return nil
}
