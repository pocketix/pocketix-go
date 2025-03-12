package models

import (
	"github.com/pocketix/pocketix-go/src/interfaces"
	"github.com/pocketix/pocketix-go/src/tree"
)

func CommandFactory(id string, blocks []interfaces.Command, tree *tree.TreeNode) (interfaces.Command, error) {
	switch id {
	case "if":
		return &If{Id: id, Block: blocks, Arguments: tree}, nil
	}
	return nil, nil
}
