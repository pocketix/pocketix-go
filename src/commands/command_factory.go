package commands

import (
	"github.com/pocketix/pocketix-go/src/models"
)

func CommandFactory(id string, blocks []Command, tree []*models.TreeNode) (Command, error) {
	switch id {
	case "if":
		if len(tree) == 0 {
			return &If{
				Id:    id,
				Block: blocks,
			}, nil
		}
		return &If{
			Id:        id,
			Block:     blocks,
			Arguments: tree[0],
		}, nil
	case "else":
		return &Else{
			Id:    id,
			Block: blocks,
		}, nil
	case "elseif":
		return &ElseIf{
			Id:        id,
			Block:     blocks,
			Arguments: tree[0],
		}, nil
	case "while":
		if len(tree) == 0 {
			return &While{
				Id:    id,
				Block: blocks,
			}, nil
		}
		return &While{
			Id:        id,
			Block:     blocks,
			Arguments: tree[0],
		}, nil
	case "setvar":
		return &SetVariable{
			Id:       id,
			LVal:     tree[0].Value.(string),
			LValType: tree[0].Type,
			RVal:     tree[1].Value,
			RValType: tree[1].Type,
		}, nil
	case "repeat":
		return &Repeat{
			Id:        id,
			Count:     tree[0].Value,
			CountType: tree[0].Type,
			Block:     blocks,
		}, nil
	case "switch":
		return &Switch{
			Id:           id,
			Block:        blocks,
			Selector:     tree[0].Value,
			SelectorType: tree[0].Type,
		}, nil
	case "case":
		return &Case{
			Id:    id,
			Block: blocks,
			Value: tree[0].Value,
			Type:  tree[0].Type,
		}, nil
	case "alert":
		return &Alert{
			Id:           id,
			Method:       tree[0].Value.(string),
			Receiver:     tree[1].Value.(string),
			ReceiverType: tree[1].Type,
			Message:      tree[2].Value.(string),
			MessageType:  tree[2].Type,
		}, nil
	}

	return nil, nil
}
