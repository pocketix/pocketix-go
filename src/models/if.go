package models

import (
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/tree"
)

type If struct {
	Id        string
	Block     []Command
	Arguments *tree.TreeNode
}

func (i *If) Execute() error {
	services.Logger.Println("Executing if")
	if result, err := i.Arguments.Evaluate(); err != nil {
		services.Logger.Println("Error executing if arguments", err)
	} else {
		if result {
			services.Logger.Println("If is true, can execute body")
		} else {
			services.Logger.Println("If is false, skip body")
		}
	}

	return nil
}

func (i *If) GetId() string {
	return i.Id
}

func (i *If) GetBody() []Command {
	return i.Block
}

func (i *If) GetArguments() *tree.TreeNode {
	return i.Arguments
}
