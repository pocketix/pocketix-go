package models

import (
	"github.com/pocketix/pocketix-go/src/tree"
)

func CommandFactory(id string, blocks []Command, tree *tree.TreeNode) (Command, error) {
	switch id {
	case "if":
		return &If{Id: id, Block: blocks, Arguments: tree}, nil
	case "else":
		return &Else{Id: id, Block: blocks}, nil
	case "elseif":
		return &ElseIf{Id: id, Block: blocks, Arguments: tree}, nil
	}
	return nil, nil
}
