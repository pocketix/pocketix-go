package parser

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/commands"
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/types"
)

// ParseBlocks parses a block of code.
//
// Parameters:
//   - block: the block to parse.
//   - variableStore: store for variables to use for parsing.
//   - procedureStore: store for procedure definitions.
//   - commandHandlingStore: store for command-related services.
//
// Returns:
//   - an AST of statements.
//   - an error if any occurred during parsing.
func ParseBlocks(
	block types.Block,
	variableStore *models.VariableStore,
	procedureStore *models.ProcedureStore,
	commandHandlingStore *models.CommandsHandlingStore,
	collector commands.Collector,
) ([]commands.Command, error) {
	argumentTree := make([]*models.TreeNode, len(block.Arguments))

	if len(block.Arguments) != 0 {
		err := ParseArguments(block.Arguments, argumentTree, variableStore, commandHandlingStore)
		if err != nil {
			return nil, err
		}
	}

	// var parsedCommands []commands.Command
	var previousSubCommand commands.Command

	for _, subBlock := range block.Body {
		// Parse nested blocks
		subAst := make([]commands.Command, 0)
		blockCollector := collector.NewCollectorBasedOnType(collector.Type(), &subAst)

		statementList, err := ParseBlocks(subBlock, variableStore, procedureStore, commandHandlingStore, blockCollector)
		if err != nil {
			return nil, err
		}
		if len(statementList) != 1 {
			for _, statement := range statementList {
				collector.Collect(statement)
			}
			continue
		}
		cmd := statementList[0]

		if cmd.GetId() == "if" {
			previousSubCommand = cmd
		} else if cmd.GetId() == "else" {
			if previousSubCommand != nil {
				previousSubCommand.(*commands.If).AddElseBlock(cmd)
				collector.Collect(previousSubCommand)
				// parsedCommands = append(parsedCommands, previousSubCommand)
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
				collector.Collect(previousSubCommand)
				// parsedCommands = append(parsedCommands, previousSubCommand)
				previousSubCommand = nil
			}

			collector.Collect(cmd)
			// parsedCommands = append(parsedCommands, cmd)
		}
	}

	if previousSubCommand != nil {
		collector.Collect(previousSubCommand)
		// parsedCommands = append(parsedCommands, previousSubCommand)
	}

	if procedureStore != nil && procedureStore.Has(block.Id) {
		procedure := procedureStore.Get(block.Id)
		statementList, err := ParseProcedureBody(procedure, variableStore, procedureStore, commandHandlingStore, collector)
		if err != nil {
			return nil, err
		}
		return statementList, nil
	}
	cmd, err := commands.CommandFactory(block.Id, *collector.GetTarget(), argumentTree, procedureStore, commandHandlingStore.CommandInvocationStore)
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
