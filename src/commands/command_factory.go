package commands

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/tree"
)

func CommandFactory(id string, blocks []Command, tree []*tree.TreeNode) (Command, error) {
	switch id {
	case "if":
		if len(tree) == 0 {
			return &If{Id: id, Block: blocks}, nil
		}
		return &If{Id: id, Block: blocks, Arguments: tree[0]}, nil
	case "else":
		return &Else{Id: id, Block: blocks}, nil
	case "elseif":
		return &ElseIf{Id: id, Block: blocks, Arguments: tree[0]}, nil
	case "while":
		if len(tree) == 0 {
			return &While{Id: id, Block: blocks}, nil
		}
		return &While{Id: id, Block: blocks, Arguments: tree[0]}, nil
	case "setvar":
		argument := models.Variable{
			Name:  tree[0].Value.(string),
			Type:  tree[1].Type,
			Value: tree[1].Value,
		}
		return &SetVariable{Id: id, VariableToSet: argument}, nil
	case "repeat":
		return &Repeat{Id: id, Count: int(tree[0].Value.(float64)), CountType: tree[0].Type, Block: blocks}, nil
	case "switch":
		return &Switch{Id: id, Block: blocks, Selector: tree[0].Value, SelectorType: tree[0].Type}, nil
	case "case":
		return &Case{Id: id, Block: blocks, Value: tree[0].Value, Type: tree[0].Type}, nil
	}
	return nil, nil
}
