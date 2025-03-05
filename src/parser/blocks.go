package parser

import (
	"github.com/pocketix/pocketix-go/src/factories"
	"github.com/pocketix/pocketix-go/src/interfaces"
	"github.com/pocketix/pocketix-go/src/models"
)

func ParseBlocks(block models.Block) (interfaces.Command, error) {
	var parsedCommands []interfaces.Command
	for _, block := range block.Body {
		cmd, err := ParseBlocks(block)
		if err != nil {
			return nil, err
		}
		parsedCommands = append(parsedCommands, cmd)
	}
	return factories.CommandFactory(block.Id, parsedCommands, block.Arguments)
}
