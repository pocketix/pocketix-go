package parser

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/commands"
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/tree"
	"github.com/pocketix/pocketix-go/src/types"
)

func ParseBlocks(block types.Block, variableStore *models.VariableStore) (commands.Command, error) {
	var argumentTree *tree.TreeNode = nil

	if len(block.Arguments) == 0 {
		services.Logger.Println("Block has no arguments")
	} else {
		arg := block.Arguments[0]
		args, err := ParseArguments(arg)
		if err != nil {
			return nil, err
		}
		argumentTree = tree.InitTree(arg.Type, args, variableStore)
	}

	var parsedCommands []commands.Command
	var previousSubCommand commands.Command

	for _, subBlock := range block.Body {
		cmd, err := ParseBlocks(subBlock, variableStore)

		if cmd.GetId() == "if" {
			previousSubCommand = cmd
		} else if cmd.GetId() == "else" {
			if previousSubCommand != nil {
				previousSubCommand.(*commands.If).AddElseBlock(cmd)
				parsedCommands = append(parsedCommands, previousSubCommand)
				previousSubCommand = nil
			} else {
				services.Logger.Println("Error: Else without if")
				return nil, fmt.Errorf("else without if")
			}
		} else if cmd.GetId() == "elseif" {
			if previousSubCommand != nil {
				previousSubCommand.(*commands.If).AddElseIfBlock(cmd)
			} else {
				services.Logger.Println("Error: Elseif without if")
				return nil, fmt.Errorf("elseif without if")
			}
		} else {
			if previousSubCommand != nil {
				previousSubCommand.Execute(variableStore)
				previousSubCommand = nil
			}

			parsedCommands = append(parsedCommands, cmd)
		}

		if err != nil {
			return nil, err
		}
		// parsedCommands = append(parsedCommands, cmd)
	}

	if previousSubCommand != nil {
		parsedCommands = append(parsedCommands, previousSubCommand)
	}

	cmd, err := commands.CommandFactory(block.Id, parsedCommands, argumentTree)
	return cmd, err
}
