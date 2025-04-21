package parser

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/commands"
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/types"
)

func ParseBlockWithoutExecuting(block types.Block, variableStore *models.VariableStore, procedureStore *models.ProcedureStore, commandHandlingStore *models.CommandsHandlingStore) ([]commands.Command, error) {
	argumentTree := make([]*models.TreeNode, len(block.Arguments))

	if len(block.Arguments) == 0 {
		services.Logger.Println("Block has no arguments")
	} else {
		err := ParseArguments(block.Arguments, argumentTree, variableStore, commandHandlingStore)
		if err != nil {
			return nil, err
		}
	}

	var parsedCommands []commands.Command
	var previousSubCommand commands.Command
	for _, subBlock := range block.Body {
		cmd, err := ParseBlockWithoutExecuting(subBlock, variableStore, procedureStore, commandHandlingStore)
		if err != nil {
			return nil, err
		}

		if cmd[0].GetId() == "if" {
			previousSubCommand = cmd[0]
		} else if cmd[0].GetId() == "else" {
			if previousSubCommand != nil {
				previousSubCommand.(*commands.If).AddElseBlock(cmd[0])
				parsedCommands = append(parsedCommands, previousSubCommand)
				previousSubCommand = nil
			} else {
				services.Logger.Println("Error: Else without if")
				return nil, fmt.Errorf("else without if")
			}
		} else if cmd[0].GetId() == "elseif" {
			if previousSubCommand != nil {
				previousSubCommand.(*commands.If).AddElseIfBlock(cmd[0])
			} else {
				services.Logger.Println("Error: Elseif without if")
				return nil, fmt.Errorf("elseif without if")
			}
		} else {
			if previousSubCommand != nil {
				previousSubCommand = nil
			}

			parsedCommands = append(parsedCommands, cmd[0])
		}
	}

	if previousSubCommand != nil {
		parsedCommands = append(parsedCommands, previousSubCommand)
	}

	if procedureStore != nil && procedureStore.Has(block.Id) {
		procedure := procedureStore.Get(block.Id)
		commandList, err := ParseProcedureBody(procedure, variableStore, procedureStore, commandHandlingStore)
		if err != nil {
			return nil, err
		}
		return commandList, nil
	}
	cmd, err := commands.CommandFactory(block.Id, parsedCommands, argumentTree, procedureStore, commandHandlingStore.CommandInvocationStore)
	if err != nil {
		return nil, err
	}
	if cmd == nil {
		services.Logger.Println("Command is nil, therefore it is device command")
		return nil, nil
	}
	err = cmd.Validate(variableStore, commandHandlingStore.ReferencedValueStore)
	return []commands.Command{cmd}, err
}

func ParseBlocks(block types.Block, variableStore *models.VariableStore, procedureStore *models.ProcedureStore, commandHandlingStore *models.CommandsHandlingStore) ([]commands.Command, error) {
	argumentTree := make([]*models.TreeNode, len(block.Arguments))

	if len(block.Arguments) == 0 {
		services.Logger.Println("Block has no arguments")
	} else {
		err := ParseArguments(block.Arguments, argumentTree, variableStore, commandHandlingStore)
		if err != nil {
			return nil, err
		}
		// argumentTree = tree.InitTree(arg.Type, args, variableStore)
	}

	var parsedCommands []commands.Command
	var previousSubCommand commands.Command

	for _, subBlock := range block.Body {
		commandList, err := ParseBlocks(subBlock, variableStore, procedureStore, commandHandlingStore)
		if err != nil {
			return nil, err
		}
		if len(commandList) != 1 {
			continue
		}
		cmd := commandList[0]

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
				_, err := previousSubCommand.Execute(variableStore, commandHandlingStore)
				if err != nil {
					return nil, err
				}
				previousSubCommand = nil
			}

			parsedCommands = append(parsedCommands, cmd)
		}
	}

	if previousSubCommand != nil {
		parsedCommands = append(parsedCommands, previousSubCommand)
	}

	if procedureStore != nil && procedureStore.Has(block.Id) {
		procedure := procedureStore.Get(block.Id)
		commandList, err := ParseProcedureBody(procedure, variableStore, procedureStore, commandHandlingStore)
		if err != nil {
			return nil, err
		}
		return commandList, nil
	}
	cmd, err := commands.CommandFactory(block.Id, parsedCommands, argumentTree, procedureStore, commandHandlingStore.CommandInvocationStore)
	if err != nil {
		services.Logger.Println("Error creating command", err)
		return nil, err
	}
	if cmd == nil {
		services.Logger.Println("Command is nil, therefore it is device command")
		return nil, nil
	}
	err = cmd.Validate(variableStore, commandHandlingStore.ReferencedValueStore)
	return []commands.Command{cmd}, err
}
