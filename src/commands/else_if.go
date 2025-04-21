package commands

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/utils"
)

type ElseIf struct {
	Id        string
	Block     []Command
	Arguments *models.TreeNode
}

func (e *ElseIf) Execute(variableStore *models.VariableStore, commandHandlingStore *models.CommandsHandlingStore) (bool, error) {
	services.Logger.Println("Executing else if")
	if result, err := e.Arguments.Evaluate(variableStore, commandHandlingStore.ReferencedValueStore); err != nil {
		services.Logger.Println("Error executing else if arguments", err)
	} else {
		if boolResult, boolErr := utils.ToBool(result); boolErr != nil {
			services.Logger.Println("Error converting else if result to bool", boolErr)
			return false, boolErr
		} else if boolResult {
			services.Logger.Println("Else if is true, can execute body")
			for _, cmd := range e.Block {
				if success, err := cmd.Execute(variableStore, commandHandlingStore); err != nil {
					return false, err
				} else if success {
					return true, nil
				}
			}
			return true, nil
		}
	}
	return false, nil
}

func (e *ElseIf) GetId() string {
	return e.Id
}

func (e *ElseIf) GetBody() []Command {
	return e.Block
}

func (e *ElseIf) GetArguments() *models.TreeNode {
	return e.Arguments
}

func (e *ElseIf) Validate(variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore, args ...any) error {
	return nil
}
