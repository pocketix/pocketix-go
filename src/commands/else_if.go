package commands

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

type ElseIf struct {
	Id        string
	Block     []Command
	Arguments *models.TreeNode
}

func (e *ElseIf) Execute(variableStore *models.VariableStore) (bool, error) {
	services.Logger.Println("Executing else if")
	if result, _, err := e.Arguments.Evaluate(variableStore); err != nil {
		services.Logger.Println("Error executing else if arguments", err)
	} else {
		if result {
			services.Logger.Println("Else if is true, can execute body")
			for _, cmd := range e.Block {
				if success, err := cmd.Execute(variableStore); err != nil {
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
