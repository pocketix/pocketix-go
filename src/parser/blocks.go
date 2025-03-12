package parser

import (
	"github.com/pocketix/pocketix-go/src/interfaces"
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/tree"
)

func ParseBlocks(block models.Block) (interfaces.Command, error) {
	args, err := ParseArguments(block.Arguments[0])
	if err != nil {
		return nil, err
	}

	var parsedCommands []interfaces.Command
	for _, block := range block.Body {
		cmd, err := ParseBlocks(block)
		if err != nil {
			return nil, err
		}
		parsedCommands = append(parsedCommands, cmd)
	}
	return models.CommandFactory(block.Id, parsedCommands, tree.InitTree(block.Arguments[0].Type, args))
}
