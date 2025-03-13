package parser

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/tree"
)

func ParseBlocks(block models.Block) (models.Command, error) {
	if len(block.Arguments) == 0 {
		return nil, fmt.Errorf("block %s is missing arguments", block.Id)
	}

	arg := block.Arguments[0]
	args, err := ParseArguments(arg)
	if err != nil {
		return nil, err
	}

	var parsedCommands []models.Command
	for _, subBlock := range block.Body {
		cmd, err := ParseBlocks(subBlock)
		if err != nil {
			return nil, err
		}
		parsedCommands = append(parsedCommands, cmd)
	}
	return models.CommandFactory(block.Id, parsedCommands, tree.InitTree(arg.Type, args))
}
