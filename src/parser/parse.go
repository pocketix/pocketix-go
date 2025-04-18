package parser

import (
	"encoding/json"
	"fmt"

	"github.com/pocketix/pocketix-go/src/commands"
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/types"
)

func CheckMissingBlock(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if _, ok := raw["header"]; !ok {
		return fmt.Errorf("header is required")
	}

	if _, ok := raw["block"]; !ok {
		return fmt.Errorf("at least one block is required")
	}
	return nil
}

func ParseHeader(data []byte, variableStore *models.VariableStore, procedureStore *models.ProcedureStore, commandHandlingStore *models.CommandsHandlingStore) (*types.Program, error) {
	var program types.Program

	if err := CheckMissingBlock(data); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &program); err != nil {
		return nil, err
	}

	if err := ParseVariables(program.Header.Variables, variableStore, commandHandlingStore); err != nil {
		return nil, err
	}

	if err := ParseProcedures(program.Header.Procedures, procedureStore, commandHandlingStore); err != nil {
		return nil, err
	}

	return &program, nil
}

func ParseProcedureBody(procedure models.Procedure, variableStore *models.VariableStore, procedureStore *models.ProcedureStore, commandHandlingStore *models.CommandsHandlingStore) ([]commands.Command, error) {
	var blocks []types.Block
	if err := json.Unmarshal(procedure.Program, &blocks); err != nil {
		return nil, err
	}

	var commandList []commands.Command
	for _, block := range blocks {
		cmd, err := ParseBlockWithoutExecuting(block, variableStore, procedureStore, commandHandlingStore)
		if err != nil {
			return nil, err
		}
		commandList = append(commandList, cmd...)
	}
	return commandList, nil
}

func ParseWithoutExecuting(data []byte, variableStore *models.VariableStore, procedureStore *models.ProcedureStore, commandHandlingStore *models.CommandsHandlingStore) error {
	program, err := ParseHeader(data, variableStore, procedureStore, commandHandlingStore)
	if err != nil {
		return err
	}

	var previousCommand commands.Command
	for _, block := range program.Blocks {
		commandList, err := ParseBlockWithoutExecuting(block, variableStore, procedureStore, commandHandlingStore)
		if err != nil {
			return err
		}
		if len(commandList) != 1 {
			continue
		}
		cmd := commandList[0]
		if cmd.GetId() == "if" {
			previousCommand = cmd
		} else if cmd.GetId() == "else" {
			if previousCommand != nil {
				previousCommand.(*commands.If).AddElseBlock(cmd)
				previousCommand = nil
			} else {
				services.Logger.Println("Error: Else without if")
				return fmt.Errorf("else without if")
			}
		} else if cmd.GetId() == "elseif" {
			if previousCommand != nil {
				previousCommand.(*commands.If).AddElseIfBlock(cmd)
			} else {
				services.Logger.Println("Error: Elseif without if")
				return fmt.Errorf("elseif without if")
			}
		} else {
			if previousCommand != nil {
				previousCommand = nil
			}
		}
	}

	return nil
}

func Parse(data []byte, variableStore *models.VariableStore, procedureStore *models.ProcedureStore, commandHandlingStore *models.CommandsHandlingStore) ([]commands.Command, error) {
	program, err := ParseHeader(data, variableStore, procedureStore, commandHandlingStore)
	if err != nil {
		return nil, err
	}

	var commandList []commands.Command
	var previousCommand commands.Command
	for _, block := range program.Blocks {
		blockList, err := ParseBlocks(block, variableStore, procedureStore, commandHandlingStore)
		if err != nil {
			return nil, err
		}
		if len(blockList) != 1 {
			commandList = append(commandList, blockList...)
			continue
		}
		if blockList == nil {
			continue
		}
		cmd := blockList[0]

		if cmd.GetId() == "if" {
			previousCommand = cmd
		} else if cmd.GetId() == "else" {
			if previousCommand != nil {
				previousCommand.(*commands.If).AddElseBlock(cmd)
				commandList = append(commandList, previousCommand)
				previousCommand = nil
			} else {
				services.Logger.Println("Error: Else without if")
				return nil, fmt.Errorf("else without if")
			}
		} else if cmd.GetId() == "elseif" {
			if previousCommand != nil {
				previousCommand.(*commands.If).AddElseIfBlock(cmd)
			} else {
				services.Logger.Println("Error: Elseif without if")
				return nil, fmt.Errorf("elseif without if")
			}
		} else {
			if previousCommand != nil {
				commandList = append(commandList, previousCommand)
				previousCommand = nil
			}

			commandList = append(commandList, cmd)
		}
	}

	if previousCommand != nil {
		commandList = append(commandList, previousCommand)
	}

	return commandList, nil
}
