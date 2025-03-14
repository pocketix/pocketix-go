package models

import (
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/tree"
)

type ElseIf struct {
	Id        string
	Block     []Command
	Arguments *tree.TreeNode
}

func (e *ElseIf) Execute() (bool, error) {
	services.Logger.Println("Executing else if")
	if result, err := e.Arguments.Evaluate(); err != nil {
		services.Logger.Println("Error executing else if arguments", err)
	} else {
		if result {
			services.Logger.Println("Else if is true, can execute body")
			for _, cmd := range e.Block {
				if success, err := cmd.Execute(); err != nil {
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

func (e *ElseIf) GetArguments() *tree.TreeNode {
	return e.Arguments
}
