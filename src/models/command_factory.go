package models

import (
	"github.com/pocketix/pocketix-go/src/tree"
)

func CommandFactory(id string, blocks []Command, tree *tree.TreeNode) (Command, error) {
	switch id {
	case "if":
		return &If{Id: id, Block: blocks, Arguments: tree}, nil
	}
	return nil, nil
}
